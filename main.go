package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/alecthomas/kong"
	"github.com/amimof/huego"
)

//Config stores user preferences
type Config struct {
	User string
}

func main() {
	var CLI struct {
		Register struct{} `cmd:"" help:"Detect bridge and authenticate a user. You should press the bridge button before."`

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
	case "register":
		register()
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
	var conf Config
	bridge, _ := huego.Discover()
	usr, err := bridge.CreateUser("huecli") // Link button needs to be pressed
	if err != nil {
		panic(err)
	}
	conf.User = usr
	err = writeConfig(&conf)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully registered huecli.")
	}
}

func getSession() *huego.Bridge {
	conf, err := readConfig()
	if err != nil {
		log.Fatal("Cannot open configuration. Please register with the hub. Error:" + err.Error())
	}
	bridge, err := huego.Discover()
	if err != nil {
		panic(err)
	}
	bridge = bridge.Login(conf.User)
	return bridge
}

func readConfig() (*Config, error) {
	// initialize conf with default values.
	conf := &Config{User: ""}
	osusr, err := user.Current()
	b, err := ioutil.ReadFile(osusr.HomeDir + "/.huecli.conf")
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func writeConfig(conf *Config) error {
	var jsonData []byte
	jsonData, err := json.Marshal(conf)
	if err != nil {
		log.Println(err)
	}
	osusr, err := user.Current()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(osusr.HomeDir+"/.huecli.conf", jsonData, 0600)
	if err != nil {
		return err
	}

	return nil
}
