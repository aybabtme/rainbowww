package rainbowww

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"sync/atomic"
)

var (
	reset      = []byte("\033[0;00m")
	returnline = []byte("\r")

	colors = func() [][]byte {
		var c [][]byte
		for i := 16; i <= 231; i++ {
			str := fmt.Sprintf("\033[38;5;%dm", i)
			c = append(c, []byte(str))
		}
		return c
	}()

	nextColor int64
)

func rgb(r, g, b uint8) []byte {
	r6 := ((uint16(r) * 5) / 255)
	g6 := ((uint16(g) * 5) / 255)
	b6 := ((uint16(b) * 5) / 255)
	i := 36*r6 + 6*g6 + b6
	if int(i) >= len(colors) {
		log.Printf("rgb(%d,%d,%d) -> rgb6(%d,%d,%d) -> i=%d -> len(colors)=%d", r, g, b, r6, g6, b6, i, len(colors))
	}
	return colors[i]
}

func NextColor(in string) string {
	idx := int(atomic.AddInt64(&nextColor, 1))
	raw := color([]byte(in), idx%len(colors))
	return string(raw)
}

func colorRGB(in []byte, r, g, b uint8) []byte {
	return append(append(rgb(r, g, b), in...), reset...)
}

func colorbyteRGB(in byte, r, g, b uint8) []byte {
	return append(append(rgb(r, g, b), in), reset...)
}

func color(in []byte, cidx int) []byte {
	return append(append(colors[cidx], in...), reset...)
}

func colorbyte(in byte, cidx int) []byte {
	return append(append(colors[cidx], in), reset...)
}

type Rainbow struct {
	wrap io.Writer
	cidx int
}

func New(w io.Writer) *Rainbow {
	return &Rainbow{wrap: w}
}

func (r *Rainbow) Write(p []byte) (int, error) {

	buf := bytes.NewBuffer(nil)
	for i := range p {
		r.cidx = (r.cidx + 1) % len(colors)
		_, _ = buf.Write(colorbyteRGB(p[i], uint8(r.cidx%255), 0, 0))
		// r.cidx = (r.cidx + 1) % len(colors)
		// _, _ = buf.Write(colorbyte(p[i], r.cidx))
	}

	_, err := buf.WriteTo(r.wrap)
	return len(p), err
}

// stolen from gorilla colors

func RGBToHSL(r, g, b uint8) (h, s, l float64) {
	fR := float64(r) / 255
	fG := float64(g) / 255
	fB := float64(b) / 255
	max := math.Max(math.Max(fR, fG), fB)
	min := math.Min(math.Min(fR, fG), fB)
	l = (max + min) / 2
	if max == min {
		// Achromatic.
		h, s = 0, 0
	} else {
		// Chromatic.
		d := max - min
		if l > 0.5 {
			s = d / (2.0 - max - min)
		} else {
			s = d / (max + min)
		}
		switch max {
		case fR:
			h = (fG - fB) / d
			if fG < fB {
				h += 6
			}
		case fG:
			h = (fB-fR)/d + 2
		case fB:
			h = (fR-fG)/d + 4
		}
		h /= 6
	}
	return
}

func HSLToRGB(h, s, l float64) (r, g, b uint8) {
	var fR, fG, fB float64
	if s == 0 {
		fR, fG, fB = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - s*l
		}
		p := 2*l - q
		fR = hueToRGB(p, q, h+1.0/3)
		fG = hueToRGB(p, q, h)
		fB = hueToRGB(p, q, h-1.0/3)
	}
	r = uint8((fR * 255) + 0.5)
	g = uint8((fG * 255) + 0.5)
	b = uint8((fB * 255) + 0.5)
	return
}

// hueToRGB is a helper function for HSLToRGB.
func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	if t < 1.0/6 {
		return p + (q-p)*6*t
	}
	if t < 0.5 {
		return q
	}
	if t < 2.0/3 {
		return p + (q-p)*(2.0/3-t)*6
	}
	return p
}
