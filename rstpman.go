package main

import (
	"fmt"
	"log"
	"os"
	"time"

	snmp "github.com/soniah/gosnmp"
)

func main() {

	target := &snmp.GoSNMP{
		Target:    os.Args[1],
		Community: "public",
		Port:      161,
		Version:   snmp.Version2c,
		Timeout:   time.Duration(1) * time.Second,
	}

	err := target.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer target.Conn.Close()

	// dot1dStpPortState ... 1.3.6.1.2.1.17.2.15.1.3
	// ifDescr           ... 1.3.6.1.2.1.2.2.1.2
	oids := []string{"1.3.6.1.2.1.17.2.15.1.3", "1.3.6.1.2.1.2.2.1.2"}

	for {
		err = target.BulkWalk(oids[0], printValue)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)		
	}
}

func printValue(pdu snmp.SnmpPDU) error {
	fmt.Printf("%s = ", pdu.Name)

	switch pdu.Value {
	case 1:
		fmt.Println("disabled")
	case 2:
		fmt.Println("blocking")
	case 3:
		fmt.Println("listening")
	case 4:
		fmt.Println("learning")
	case 5:
		fmt.Println("forwarding")
	case 6:
		fmt.Println("broken")
	}

	return nil
}