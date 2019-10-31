package main

import (
        "fmt"
        "log"
        "os"
        "strings"
        snmp "github.com/soniah/gosnmp"
)

func main() {
        target, err := os.Open(os.Args[1])
        if err != nil {
                log.Fatal(err)
        }

        snmp.Default.Target = target
        err := snmp.Defaul.Connect()
        if err != nil {
                log.Fatal(err)
        }
        defer snmp.Default.Conn.Close()

		// *dot1dStpPortState	... 1.3.6.1.2.1.17.2.15.1.3
		// *ifDescr 			... 1.3.6.1.2.1.2.2.1.2
		oids = []string{"1.3.6.1.2.1.17.2.15.1.3", "1.3.6.1.2.1.2.2.1.2"}
		
}