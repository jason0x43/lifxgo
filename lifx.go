package main

import "github.com/jason0x43/lifxgo/lifx"
import "log"
import "fmt"
import "flag"
import "io/ioutil"

func main() {
	var verbose = flag.Bool("verbose", false, "Be verbose")
	flag.Parse()

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	gateways, err := lifx.DiscoverGateways(false)
	if err != nil {
		fmt.Print("ERROR:", err)
		return
	}

	gateway := gateways[0]
	fmt.Println("Found gateway at", gateway.Address)

	err = gateway.Dial()
	if err == nil {
		fmt.Println("Connected!")
	}

	lights, _ := gateway.RefreshLights()
	if err == nil {
		fmt.Println("Found lights:", lights)
	}

	var cmd string

	for cmd != "q" {
		fmt.Print("> ")
		fmt.Scan(&cmd)

		switch cmd {
		case "on":
			fmt.Println("Turning", lights[0], "on...")
			err = lights[0].SetPower(50)
		case "off":
			fmt.Println("Turning", lights[0], "off...")
			err = lights[0].SetPower(0)
		}

		if err != nil {
			fmt.Println("ERROR:", err)
		}
	}
}
