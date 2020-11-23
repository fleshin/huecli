package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/amimof/huego"
)

const (
	USER string = "TmyIhrQ4A3TmcduWcp6RGJliXkxEUpdXa97444HE"
)

func main() {
	var CLI struct {
		Register struct {
			Username string `arg:"" name:"username" help:"User to register on the bridge." type:"string"`
		} `cmd:"" help:"Detect bridge and authenticate a user. You should press the bridge button before."`

		Turn struct {
			On struct {
				ID int `arg:"" name:"id" help:"ID of the light to turn on." type:"int"`
			} `cmd:"" help:"Turn on light."`
			Off struct {
				ID int `arg:"" name:"id" help:"ID of the light to turn off." type:"int"`
			} `cmd:"" help:"Turn off light."`
		} `cmd:"" help:"Switch a light."`

		List struct {
		} `cmd:"" help:"List lights."`

		Dim struct {
			Intensity uint8 `arg:"" name:"intensity" help:"Intensity of the light from 0 to 255." type:"uint8"`
			ID        int   `arg:"" name:"id" help:"ID of the light to dim." type:"int"`
		} `cmd:"" help:"Dim lights."`

		Temp struct {
			Deg uint16 `arg:"" name:"deg" help:"Color of the light 153 to 500 mireks." type:"uint16"`
			ID  int    `arg:"" name:"id" help:"ID of the light to dim." type:"int"`
		} `cmd:"" help:"Set lights color temperature."`
	}

	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "list":
		ses := getSession()
		list(ses)
	case "turn on <id>":
		ses := getSession()
		turnlight(ses, true, CLI.Turn.On.ID)
	case "turn off <id>":
		ses := getSession()
		turnlight(ses, false, CLI.Turn.Off.ID)
	case "dim <intensity> <id>":
		ses := getSession()
		dimlight(ses, CLI.Dim.Intensity, CLI.Dim.ID)
	case "temp <deg> <id>":
		ses := getSession()
		templight(ses, CLI.Temp.Deg, CLI.Temp.ID)
	default:
		//flag.PrintDefaults()
		os.Exit(1)
	}

}

func list(bridge *huego.Bridge) {
	l, err := bridge.GetLights()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d lights\n", len(l))
	for _, light := range l {
		fmt.Println(light.ID, light.Name, light.Type, light.ManufacturerName)
	}
}

func turnlight(bridge *huego.Bridge, st bool, id int) {
	l, err := bridge.GetLight(id)
	if err != nil {
		panic(err)
	}
	if st {
		l.On()
		return
	}
	l.Off()
	return
}

func dimlight(bridge *huego.Bridge, intensity uint8, id int) {
	l, err := bridge.GetLight(id)
	if err != nil {
		panic(err)
	}
	l.Bri(intensity)
	return
}

func templight(bridge *huego.Bridge, deg uint16, id int) {
	l, err := bridge.GetLight(id)
	if err != nil {
		panic(err)
	}
	l.Ct(deg)
	return
}

func count() {
}

func register() {
	bridge, _ := huego.Discover()
	user, _ := bridge.CreateUser("huego") // Link button needs to be pressed
	fmt.Println("User: " + user)
	bridge = bridge.Login(user)
	light, _ := bridge.GetLight(3)
	light.Off()
}

func getSession() *huego.Bridge {
	bridge, err := huego.Discover()
	if err != nil {
		panic(err)
	}
	bridge = bridge.Login(USER)
	return bridge
}
