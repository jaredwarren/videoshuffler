package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/rapidloop/skv"
	"github.com/spf13/viper"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// PIN decimal pin number
const PIN = 17

var (
	p  *Player
	db *skv.KVStore

	// Shuffle ...
	Shuffle = true

	// AutoStart ...
	AutoStart = false
)

func main() {
	var startFile = ""
	if db, err := skv.Open("/home/pi/player.db"); err == nil {
		db.Get("StartFile", &startFile)
		defer db.Close()
	} else {
		// If we can't open the db, just use default
		fmt.Fprintf(os.Stderr, "error opening db: %v\n", err)
	}

	// Load Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/home/pi/")
	if err := viper.ReadInConfig(); err != nil {
		// ignore errors and use default
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	viper.SetDefault("Autostart", false)
	AutoStart = viper.GetBool("AutoStart")

	viper.SetDefault("BasePath", "/home/pi/Videos/*")
	BasePath := viper.GetString("BasePath")
	episodes, err := filepath.Glob(BasePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if len(episodes) == 0 {
		fmt.Fprintf(os.Stderr, "error: No video files found at %s\n", BasePath)
		os.Exit(1)
	}

	viper.SetDefault("Shuffle", true)
	Shuffle = viper.GetBool("Shuffle")
	if Shuffle {
		ShuffleList(episodes)
	}

	if AutoStart {
		go playEpisodes(episodes, startFile)
	}

	// Robot
	r := raspi.NewAdaptor()
	button := gpio.NewButtonDriver(r, fmt.Sprintf("%X", PIN))

	work := func() {
		button.On(gpio.ButtonRelease, func(data interface{}) {
			if p == nil {
				go playEpisodes(episodes, startFile)
			} else {
				p.End()
			}
		})
		button.On(gpio.Error, func(data interface{}) {
			fmt.Fprintf(os.Stderr, "error: %v\n", data)
		})
	}

	robot := gobot.NewRobot("simpsons",
		[]gobot.Connection{r},
		[]gobot.Device{button},
		work,
	)

	robot.Start()
}

// ShuffleList randomly reorder slice of strings
func ShuffleList(slc []string) {
	N := len(slc)
	for i := 0; i < N; i++ {
		// choose index uniformly in [i, N-1]
		r := i + rand.Intn(N-i)
		slc[r], slc[i] = slc[i], slc[r]
	}
}

func playEpisodes(episodes []string, startFile string) {
	p = NewPlayer()
	// Play everything
	for _, episode := range episodes {
		// Skip until startFile
		if startFile != "" && startFile != episode {
			continue
		}
		startFile = ""
		if db != nil {
			db.Put("StartFile", episode)
		}

		err := p.Start(episode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "start error: %v\n", err)
		}
		err = p.Handler.Wait()
		if err != nil {
			fmt.Fprintf(os.Stderr, "wait error: %v\n", err)
		}
		err = p.End()
		if err != nil {
			fmt.Fprintf(os.Stderr, "end error: %v\n", err)
		}
	}

	// cleanup player
	p = nil
	if db != nil {
		db.Delete("StartFile")
	}
}
