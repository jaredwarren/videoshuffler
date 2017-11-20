package main

// GOARM=6 GOARCH=arm GOOS=linux go build
import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

var episodes []string
var episodeIndex int

func main() {
	r := raspi.NewAdaptor()
	button := gpio.NewButtonDriver(r, "11") // 17 in hex

	// TODO: add path config
	episodes = getEpisodes("/home/pi/Videos/simpsons/Simpsons*")

	// TODO: add shuffle config
	ShuffleStrings(episodes)

	work := func() {
		button.On(gpio.ButtonRelease, func(data interface{}) {
			// start over
			if episodeIndex+1 >= len(episodes) {
				episodeIndex = 0
			} else {
				episodeIndex++
			}

			playEpisode(episodes[episodeIndex])
		})
		button.On(gpio.Error, func(data interface{}) {
			fmt.Println("Error:", data)
		})
	}

	robot := gobot.NewRobot("simpsons",
		[]gobot.Connection{r},
		[]gobot.Device{button},
		work,
	)

	robot.Start()
}

func getEpisodes(path string) []string {
	matches, err := filepath.Glob(path)
	if err != nil {
		panic(err)
	}
	return matches
}

// ShuffleStrings ...
func ShuffleStrings(slc []string) {
	N := len(slc)
	for i := 0; i < N; i++ {
		// choose index uniformly in [i, N-1]
		r := i + rand.Intn(N-i)
		slc[r], slc[i] = slc[i], slc[r]
	}
}

var p *Player

func playEpisode(episode string) {
	if episode == "" {
		return
	}
	if p != nil {
		p.End()
	}
	fmt.Println("Playing:", episode)
	p = NewPlayer(episode)
	p.Start()
}
