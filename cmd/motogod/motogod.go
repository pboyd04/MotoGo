package main

import (
	"flag"
	"fmt"
	"time"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/pboyd04/MotoGo/internal/moto"
	"github.com/pboyd04/MotoGo/internal/moto/mototrbo"
)

func main() {
	//Parse the command line flags
	masterAddrPtr := flag.String("master", "192.168.0.100:50000", "The motrobo master with IP and port")
	influxAddr := flag.String("influx", "http://localhost:8086", "The influx DB instance address")
	myIDPtr := flag.Int("id", 1, "The radio ID for this node to use")

	flag.Parse()

	//Connect to influx db
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: *influxAddr,
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	sys, err := moto.NewRadioSystem(mototrbo.RadioID(*myIDPtr), mototrbo.CapacityPlus)
	if err != nil {
		panic(err)
	}
	err = sys.ConnectToMaster(*masterAddrPtr)
	if err != nil {
		panic(err)
	}
	defer sys.Close()
	fmt.Printf("Master ID = %d\n", sys.GetMasterID())
	master := sys.GetMaster()
	master.InitXNL()
	fmt.Printf("    XNL ID = %d\n", sys.GetMasterXNLID())
	fmt.Printf("    XCMP Version = %s\n", master.GetXCMPVersion())
	fmt.Printf("    Serial Number = %s\n", master.GetSerialNumber())
	fmt.Printf("    Firmware Version = %s\n", master.GetFirmwareVersion())
	fmt.Printf("    Model Number = %s\n", master.GetModelNumber())
	fmt.Printf("    Radio Alias = %s\n", master.GetRadioAlias())
	rssi1, rssi2 := master.GetRSSI()
	fmt.Printf("    RSSI = %f %f\n", rssi1, rssi2)
	fmt.Printf("    Alarms\n")
	alarms := master.GetAlarmStatus()
	for name, state := range alarms {
		fmt.Printf("        %s: %t\n", name, state)
	}
	calls := make(chan *moto.RadioCall, 10)
	masterCallCount := make(chan int)
	go master.ListenForCalls(calls, masterCallCount)
	go logCountChanges(c, master.ID, masterCallCount)
	peers := sys.PeerList()
	for index, peer := range peers {
		peerCallCount := make(chan int)
		fmt.Printf("Peer %d ID: %d\n", index, peer.ID)
		peer.InitXNL()
		fmt.Printf("    XNL ID = %d\n", peer.GetXNLID())
		fmt.Printf("    XCMP Version = %s\n", peer.GetXCMPVersion())
		fmt.Printf("    Serial Number = %s\n", peer.GetSerialNumber())
		fmt.Printf("    Firmware Version = %s\n", peer.GetFirmwareVersion())
		fmt.Printf("    Model Number = %s\n", peer.GetModelNumber())
		fmt.Printf("    Radio Alias = %s\n", peer.GetRadioAlias())
		rssi1, rssi2 := peer.GetRSSI()
		fmt.Printf("    RSSI = %f %f\n", rssi1, rssi2)
		fmt.Printf("    Alarms\n")
		//alarms := peer.GetAlarmStatus()
		//for name, state := range alarms {
		//	fmt.Printf("        %s: %t\n", name, state)
		//}
		go peer.ListenForCalls(calls, peerCallCount)
		go logCountChanges(c, peer.ID, peerCallCount)
	}
	for {
		call := <-calls
		fmt.Printf("%s: Got call from %d to %d (%f seconds)\n", call.StartTime, call.From, call.To, call.EndTime.Sub(call.StartTime).Seconds())
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

func writeCallToDB(c client.Client, call *moto.RadioCall) {
	tags := map[string]string{
		"To":     radioIDToString(call.To, call.Group),
		"From":   radioIDToString(call.From, false),
		"Length": fmt.Sprintf("%f", call.EndTime.Sub(call.StartTime).Seconds()),
	}
	point, err := client.NewPoint("calls", tags, map[string]interface{}{"value": 1}, call.StartTime)
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
	//TODO Look up the id
	if group {
		return fmt.Sprintf("Group %d", id)
	}
	return fmt.Sprintf("Radio %d", id)
}
