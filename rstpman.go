package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
	"runtime"
	"errors"

	snmp "github.com/soniah/gosnmp"
)

func main() {
	var ip string = "127.0.0.1"
	var community string = "public"
	var interval int = 3
	var err error

	if len(os.Args) == 3 {
		ip = os.Args[1]
		community = os.Args[2]
	} else if len(os.Args) == 4 {
		ip = os.Args[1]
		community = os.Args[2]
		interval,err = strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
	} else {

		fmt.Println("Usage: rstpman <ip> <community> <interval>")
		return
	}

	target := &snmp.GoSNMP{
		Target:    ip,
		Community: community,
		Port:      161,
		Version:   snmp.Version2c,
		Timeout:   time.Duration(1) * time.Second,
	}

	err = target.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer target.Conn.Close()

	ifIndexMap, err := getInterfaces(target)
	if err != nil {
		log.Fatal(err)
	}


	for {
		err = clearScreen()
		if err != nil {
			log.Fatal(err)
		}

		portStatus := ""
		// dot1dStpPortState ... 1.3.6.1.2.1.17.2.15.1.3
		err = target.BulkWalk("1.3.6.1.2.1.17.2.15.1.3", func(pdu snmp.SnmpPDU) error {
			// export ifIndex from pdu
			portStatus += fmt.Sprintf("%-10v = ", ifIndexMap[pdu.Name[25:]])

			switch pdu.Value {
			case 1:
				portStatus += "disabled\n"
			case 2:
				portStatus += "blocking\n"
			case 3:
				portStatus += "listening\n"
			case 4:
				portStatus += "learning\n"
			case 5:
				portStatus += "forwarding\n"
			case 6:
				portStatus += "broken^n"
			}

			return nil
		})

		now := time.Now()
		fmt.Printf("Fetched from %s at %d-%02d-%02d %02d:%02d:%02d\n\n",
			target.Target,
			now.Year(), now.Month(), now.Day(),
			now.Hour(), now.Minute(), now.Second(),
		)

		fmt.Println(portStatus)

		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func getInterfaces(target *snmp.GoSNMP) (map[string]string, error) {
	ifIndexMap := map[string]string{}

	// Get mapping ifIndex from BRIDGE-MIB to IF-MIB
	// dot1dBasePortIfIndex ... 1.3.6.1.2.1.17.1.4.1.2
	err := target.BulkWalk("1.3.6.1.2.1.17.1.4.1.2", func(pdu snmp.SnmpPDU) error {
		ifIndexMap[pdu.Name[24:]] = strconv.Itoa(pdu.Value.(int))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create map from dot1dStpPortState etnry to ifName
	// ifName ... 1.3.6.1.2.1.31.1.1.1.1
	err = target.BulkWalk("1.3.6.1.2.1.31.1.1.1.1", func(pdu snmp.SnmpPDU) error {
		for k, v := range ifIndexMap {
			if v == pdu.Name[24:] {
				ifIndexMap[k] = string(pdu.Value.([]uint8))
			}
		}
		return nil
	})

	return ifIndexMap, err
}

func clearScreen() error {
	var clearCmd *exec.Cmd

	if runtime.GOOS == "windows" {
		clearCmd = exec.Command("cls")
	} else if runtime.GOOS == "linux" {
		clearCmd = exec.Command("clear")
	} else if runtime.GOOS == "darwin" {
		clearCmd = exec.Command("clear")
	} else {
		return cantDetectOSTypeError()
	}
	clearCmd.Stdout = os.Stdout
	clearCmd.Run()
	return nil
}

func cantDetectOSTypeError() error {
	return errors.New("can't detect OS type")
}
