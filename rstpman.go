package main

import (
        "fmt"
        "log"
		"os"
		"time"
        snmp "github.com/soniah/gosnmp"
)

func main() {
		target := os.Args[1]

		params := &snmp.GoSNMP{
			Target: target,
			Community: "public",
			Port: 161,
			Version: snmp.Version2c,
			Timeout: time.Duration(1) * time.Second,
			Logger: log.New(os.Stdout, "", 0),
		}

        err := params.Connect()
        if err != nil {
                log.Fatal(err)
        }
        defer params.Conn.Close()

		// dot1dStpPortState ... 1.3.6.1.2.1.17.2.15.1.3
		// ifDescr           ... 1.3.6.1.2.1.2.2.1.2
		oids := []string{"1.3.6.1.2.1.17.2.15.1.3", "1.3.6.1.2.1.2.2.1.2"}
		result, err := params.Get(oids)
		if err != nil {
			log.Fatal(err)
		}

		for i, variable := range result.Variables {
			fmt.Printf("%d: oid: %s ", i, variable.Name)

			// the Value of each variable returned by Get() implements
			// interface{}. You could do a type switch...
			switch variable.Type {
			case snmp.OctetString:
				fmt.Printf("string: %s\n", string(variable.Value.([]byte)))
			default:
				// ... or often you're just interested in numeric values.
				// ToBigInt() will return the Value as a BigInt, for plugging
				// into your calculations.
				fmt.Printf("number: %d\n", snmp.ToBigInt(variable.Value))
		}
	}
		
}