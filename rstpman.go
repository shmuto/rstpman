package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	snmp "github.com/soniah/gosnmp"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage: rstpman <ip> <community>")
		return 
	}

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
	oids := []string{"1.3.6.1.2.1.17.2.15.1.3", "1.3.6.1.2.1.2.2.1.2"}

	for {
		clearScreen()

		now := time.Now()
		fmt.Printf("Fetched from %s at %d-%02d-%02d %02d:%02d:%02d\n\n",
			target.Target,
			now.Year(), now.Month(), now.Day(),
			now.Hour(), now.Minute(), now.Second(),
		)

		ifList, err := getInterfaces(target)
		if err != nil {
			log.Fatal(err)
		}



		err = target.BulkWalk(oids[0], printValue)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)
	}
}

func getInterfaces(target *snmp.GoSNMP) (map[string]int, error) {
	ifIndexMap := map[string]int{}

	// dot1dBasePortIfIndex ... 1.3.6.1.2.1.17.1.4.1.2
	err := target.BulkWalk("1.3.6.1.2.1.17.1.4.1.2", func (pdu snmp.SnmpPDU) error {
		ifIndexMap[pdu.Name[24:]] = pdu.Value.(int)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	//ifName ... 1.3.6.1.2.1.31.1.1.1.1 
	err := target.BulkWalk("1.3.6.1.2.1.31.1.1.1.1", func (pdu snmp.SnmpPDU) error {

	})
	return ifIndex, nil 
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

func clearScreen() {
	clearCmd := exec.Command("clear")
	clearCmd.Stdout = os.Stdout
	clearCmd.Run()
}