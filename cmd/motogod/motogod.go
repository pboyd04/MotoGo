package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/pboyd04/MotoGo/internal/moto"
	"github.com/pboyd04/MotoGo/internal/moto/mototrbo"
)

func main() {
	//Default config values
	viper.SetDefault("master", "192.168.0.100:50000")
	viper.SetDefault("influx", "http://localhost:8086")
	viper.SetDefault("id", 1)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	//Parse the command line flags
	flag.String("master", "192.168.0.100:50000", "The motrobo master with IP and port")
	flag.String("influx", "http://localhost:8086", "The influx DB instance address")
	flag.Int("id", 1, "The radio ID for this node to use")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.ReadInConfig()
	viper.BindPFlags(pflag.CommandLine)

	//Connect to influx db
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: viper.GetString("influx"),
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	sys, err := moto.NewRadioSystem(mototrbo.RadioID(viper.GetInt("id")), mototrbo.CapacityPlus)
	if err != nil {
		panic(err)
	}
	err = sys.ConnectToMaster(viper.GetString("master"))
	if err != nil {
		panic(err)
	}
	defer sys.Close()
	fmt.Printf("Master ID = %d\n", sys.GetMasterID())
	master := sys.GetMaster()
	master.InitXNL()
	aliases := viper.GetStringMapString("aliases.radios")
	_, ok := aliases[strconv.Itoa(int(master.ID))]
	if !ok {
		//Repeater isn't in config file. Write it...
		aliases[strconv.Itoa(int(master.ID))] = master.GetRadioAlias()
		viper.Set("aliases.radios", aliases)
		viper.WriteConfig()
	}
	fmt.Printf("    XNL ID = %d\n", sys.GetMasterXNLID())
	fmt.Printf("    XCMP Version = %s\n", master.GetXCMPVersion())
	fmt.Printf("    Serial Number = %s\n", master.GetSerialNumber())
	fmt.Printf("    Firmware Version = %s\n", master.GetFirmwareVersion())
	fmt.Printf("    Model Number = %s\n", master.GetModelNumber())
	fmt.Printf("    Alarms\n")
	alarms := master.GetAlarmStatus()
	for name, state := range alarms {
		fmt.Printf("        %s: %t\n", name, state)
	}
	calls := make(chan *moto.RadioCall, 10)
	masterCallCount := make(chan int)
	go master.ListenForCalls(calls, masterCallCount)
	go logCountChanges(c, master.ID, masterCallCount)
	go logRSSI(c, master)
	peers := sys.PeerList()
	for index, peer := range peers {
		peerCallCount := make(chan int)
		fmt.Printf("Peer %d ID: %d\n", index, peer.ID)
		peer.InitXNL()
		aliases := viper.GetStringMapString("aliases.radios")
		_, ok := aliases[strconv.Itoa(int(peer.ID))]
		if !ok {
			//Repeater isn't in config file. Write it...
			aliases[strconv.Itoa(int(master.ID))] = peer.GetRadioAlias()
			viper.Set("aliases.radios", aliases)
			viper.WriteConfig()
		}
		fmt.Printf("    XNL ID = %d\n", peer.GetXNLID())
		fmt.Printf("    XCMP Version = %s\n", peer.GetXCMPVersion())
		fmt.Printf("    Serial Number = %s\n", peer.GetSerialNumber())
		fmt.Printf("    Firmware Version = %s\n", peer.GetFirmwareVersion())
		fmt.Printf("    Model Number = %s\n", peer.GetModelNumber())
		fmt.Printf("    Radio Alias = %s\n", peer.GetRadioAlias())
		fmt.Printf("    Alarms\n")
		//alarms := peer.GetAlarmStatus()
		//for name, state := range alarms {
		//	fmt.Printf("        %s: %t\n", name, state)
		//}
		go peer.ListenForCalls(calls, peerCallCount)
		go logCountChanges(c, peer.ID, peerCallCount)
		go logRSSI(c, peer)
	}
	for {
		call := <-calls
		writeCallToDB(c, call)
	}
}

func logCountChanges(c client.Client, id mototrbo.RadioID, countChan chan int) {
	tags := map[string]string{
		"Radio": radioIDToString(id, false),
	}
	for {
		count := <-countChan
		point, err := client.NewPoint("count", tags, map[string]interface{}{"value": count}, time.Now())
		if err != nil {
			fmt.Printf("Error creating point %v\n", err)
		}
		batch, err := client.NewBatchPoints(client.BatchPointsConfig{Precision: "s", Database: "radios"})
		if err != nil {
			fmt.Printf("Error creating batch %v\n", err)
		}
		batch.AddPoint(point)
		err = c.Write(batch)
		if err != nil {
			fmt.Printf("Error writing batch %v\n", err)
		}
	}
}

func logRSSI(c client.Client, radio *moto.RemoteRadio) {
	tags := map[string]string{
		"Radio": radioIDToString(radio.ID, false),
	}
	for {
		if radio.GetActiveCallCount() == 0 {
			rssi1, rssi2 := radio.GetRSSI()
			point, err := client.NewPoint("rssi", tags, map[string]interface{}{"rssi1": rssi1, "rssi2": rssi2}, time.Now())
			if err != nil {
				fmt.Printf("Error creating point %v\n", err)
			}
			batch, err := client.NewBatchPoints(client.BatchPointsConfig{Precision: "s", Database: "radios"})
			if err != nil {
				fmt.Printf("Error creating batch %v\n", err)
			}
			batch.AddPoint(point)
			err = c.Write(batch)
			if err != nil {
				fmt.Printf("Error writing batch %v\n", err)
			}
		}
		<-time.After(5 * time.Minute)
	}
}

func writeCallToDB(c client.Client, call *moto.RadioCall) {
	tags := map[string]string{
		"To":     radioIDToString(call.To, call.Group),
		"From":   radioIDToString(call.From, false),
		"Length": fmt.Sprintf("%f", call.EndTime.Sub(call.StartTime).Seconds()),
	}
	fmt.Printf("%s: Got call from %s to %s (%f seconds)\n", call.StartTime, tags["From"], tags["To"], call.EndTime.Sub(call.StartTime).Seconds())
	point, err := client.NewPoint("calls", tags, map[string]interface{}{"value": 1, "length": call.EndTime.Sub(call.StartTime).Seconds()}, call.StartTime)
	if err != nil {
		fmt.Printf("Error creating point %v\n", err)
	}
	batch, err := client.NewBatchPoints(client.BatchPointsConfig{Precision: "s", Database: "radios"})
	if err != nil {
		fmt.Printf("Error creating batch %v\n", err)
	}
	batch.AddPoint(point)
	err = c.Write(batch)
	if err != nil {
		fmt.Printf("Error writing batch %v\n", err)
	}
}

func radioIDToString(id mototrbo.RadioID, group bool) string {
	if group {
		aliases := viper.GetStringMapString("aliases.groups")
		alias, ok := aliases[strconv.Itoa(int(id))]
		if ok {
			return alias
		}
		return fmt.Sprintf("Group %d", id)
	}
	aliases := viper.GetStringMapString("aliases.radios")
	alias, ok := aliases[strconv.Itoa(int(id))]
	if ok {
		return alias
	}
	return fmt.Sprintf("Radio %d", id)
}
