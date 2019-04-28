package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/koszonetdoktor/goADS"
)

func main() {
	fmt.Println("hello main Torkoly szerver!")
	// Flags/*{{{*/
	debug := flag.Bool("debug", false, "print debugging messages.")
	ip := flag.String("ip", "", "the address to the AMS router")
	netid := flag.String("netid", "", "AMS NetID of the target")
	port := flag.Int("port", 801, "AMS Port of the target")

	flag.Parse()
	fmt.Println(*debug, *ip, *netid, *port)
	// Startup the connection/*{{{*/
	connection, e := goADS.NewConnection(*ip, *netid, *port)
	defer connection.Close() // Close the connection when we are done
	if e != nil {
		fmt.Println("ERROR: ", e)
		os.Exit(1)
	} /*}}}*/
	// Add a handler for Ctrl^C,  soft shutdown/*{{{*/
	go shutdownRoutine(connection) /*}}}*/
	// Check what device are we connected to/*{{{*/
	data, e := connection.ReadDeviceInfo()
	if e != nil {
		fmt.Println("ERROR in read device: ", e)
		os.Exit(1)
	}
	fmt.Printf("Successfully conncected to \"%s\" version %d.%d (build %d)", data.DeviceName, data.MajorVersion, data.MinorVersion, data.BuildVersion)
}

func shutdownRoutine(conn *goADS.Connection) {
	sigchan := make(chan os.Signal, 2)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, syscall.SIGTERM)
	<-sigchan

	conn.Close()
}
