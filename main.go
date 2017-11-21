package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

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

var p *Player
var wg sync.WaitGroup

func main() {
	// Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	viper.SetDefault("Autostart", false)

	viper.SetDefault("BasePath", "/home/pi/Videos/simpsons/Simpsons*")
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
	if viper.GetBool("Shuffle") {
		ShuffleList(episodes)
	}

	go playEpisodes(episodes)

	// Robot
	r := raspi.NewAdaptor()
	button := gpio.NewButtonDriver(r, fmt.Sprintf("%X", PIN))

	work := func() {
		button.On(gpio.ButtonRelease, func(data interface{}) {
			if p == nil {
				if !viper.GetBool("Autostart") {
					wg.Done()
				}
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

func playEpisodes(episodes []string) {
	// wait until button is pressed to start playing.
	if !viper.GetBool("Autostart") {
		wg.Add(1)
		wg.Wait()
	}
	p = NewPlayer()
	// Play everything
	for _, episode := range episodes {
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
}
