package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
	"strconv"

	snmp "github.com/soniah/gosnmp"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage: rstpman <ip> <community>")
		return 
	}

	target := &snmp.GoSNMP{
		Target:    os.Args[1],
		Community: os.Args[2],
		Port:      161,
		Version:   snmp.Version2c,
		Timeout:   time.Duration(1) * time.Second,
	}

	err := target.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer target.Conn.Close()

	for {
		clearScreen()

		now := time.Now()
		fmt.Printf("Fetched from %s at %d-%02d-%02d %02d:%02d:%02d\n\n",
			target.Target,
			now.Year(), now.Month(), now.Day(),
			now.Hour(), now.Minute(), now.Second(),
		)

		ifIndexMap, err := getInterfaces(target)
		if err != nil {
			log.Fatal(err)
		}

		// dot1dStpPortState ... 1.3.6.1.2.1.17.2.15.1.3
		err = target.BulkWalk("1.3.6.1.2.1.17.2.15.1.3", func(pdu snmp.SnmpPDU) error {
			fmt.Printf("%s = ", ifIndexMap[pdu.Name[25:]])

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
		})
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(2 * time.Second)
	}
}

func getInterfaces(target *snmp.GoSNMP) (map[string]string, error) {
	ifIndexMap := map[string]string{}

	// Get mapping ifIndex from BRIDGE-MIB to IF-MIB
	// dot1dBasePortIfIndex ... 1.3.6.1.2.1.17.1.4.1.2
	err := target.BulkWalk("1.3.6.1.2.1.17.1.4.1.2", func (pdu snmp.SnmpPDU) error {
		ifIndexMap[pdu.Name[24:]] = strconv.Itoa(pdu.Value.(int))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create map from dot1dStpPortState etnry to ifName
	//ifName ... 1.3.6.1.2.1.31.1.1.1.1 
	err = target.BulkWalk("1.3.6.1.2.1.31.1.1.1.1", func (pdu snmp.SnmpPDU) error {
		for k, v := range ifIndexMap {
			if v == pdu.Name[24:] {
				ifIndexMap[k] = string(pdu.Value.([]uint8))
			}
		}
		return nil
	})

	return ifIndexMap, err
}

func clearScreen() {
	clearCmd := exec.Command("clear")
	clearCmd.Stdout = os.Stdout
	clearCmd.Run()
}