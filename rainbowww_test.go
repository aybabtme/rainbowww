package rainbowww

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func ExampleNew() {

	timeout := time.NewTimer(time.Second * 1)
	ticker := time.NewTicker(time.Millisecond * 10)

	log.SetOutput(New(os.Stderr, 252, 255, 43))

	for t := range ticker.C {

		log.Printf("it's %v", t)
		select {
		case <-timeout.C:
			log.Print("stopping!")
			return
		default:
		}
	}

	// Output:
	//
}

func TestRGB(t *testing.T) {
	data := []byte("█")

	h, s, l := rgbToHSL(252, 255, 43)

	timeout := time.NewTimer(time.Second * 1)
	ticker := time.NewTicker(time.Millisecond * 10)

	for _ = range ticker.C {

		select {
		case <-timeout.C:
			log.Print("stopping!")
			return
		default:
		}

		h += (6.0 / 360.0)
		if h > 1.0 {
			h = 0.0
		}
		r, g, b := hslToRGB(h, s, l)
		fmt.Printf("%s", colorRGB(data, r, g, b))

	}
}
