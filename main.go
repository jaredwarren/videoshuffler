package main

// GOARM=6 GOARCH=arm GOOS=linux go build
import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	//led := gpio.NewLedDriver(r, "7")
	button := gpio.NewButtonDriver(r, "17")

	work := func() {
		// gobot.Every(1*time.Second, func() {
		// 	led.Toggle()
		// })

		button.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("PUSH")
		})
		button.On(gpio.ButtonRelease, func(data interface{}) {
			fmt.Println("Release")
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{r},
		[]gobot.Device{button},
		work,
	)

	robot.Start()
}
