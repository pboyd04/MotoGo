# MotoGo
A go interface to talk to mototrbo repeaters

## Primary Use Case
I'm using this to monitor 2 repeaters used for an event keeping track of how many radio calls each talk group makes, how many talk
groups are in use at a time, any repeater faults, etc.

## XNL/XCMP
To do this I had to figure out XNL/XCMP. However, in an effort not to get a takedown notice I haven't included the constants needed by the
XNL encrypter. With the algorithm I have here, Wireshark, and a copy of RDAC or similar software you should be able to reverse engineer the
constants after a few itterations. 
