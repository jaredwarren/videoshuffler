package main

// GOARM=6 GOARCH=arm GOOS=linux go build
import (
	"fmt"
	"math/rand"
	"path/filepath"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

var episodes []string

func main() {
	r := raspi.NewAdaptor()
	button := gpio.NewButtonDriver(r, "17")

	work := func() {
		button.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("PUSH")
		})
		button.On(gpio.ButtonRelease, func(data interface{}) {
			episode := getRandomEpisode()
			playEpisode(episode)
		})
	}

	robot := gobot.NewRobot("simpsons",
		[]gobot.Connection{r},
		[]gobot.Device{button},
		work,
	)

	robot.Start()
}

func getRandomEpisode() string {
	matchLen := len(episodes)
	if matchLen == 0 {
		episodes = getEpisodes("videos/Simpsons*")
		matchLen = len(episodes)
	}

	i := rand.Intn(matchLen)
	randomEpisode := episodes[i]
	episodes = remove(episodes, i)
	return randomEpisode
}

func getEpisodes(path string) []string {
	matches, err := filepath.Glob(path)
	if err != nil {
		panic(err)
	}
	return matches
}

// remove this doesn't care about ordering
func remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func playEpisode(episode string) {
	if episode == "" {
		return
	}

	fmt.Println("Play:", episode)
	// kill process
	// play video
}
