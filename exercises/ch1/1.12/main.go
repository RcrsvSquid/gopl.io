/*
Modify the Lissajous server to read parameter values from the URL. For example,
you might arrange it so that a URL like http://localhost:8000/?cycles=20 sets
the number of cycles to 20 instead of the default 5. Use the strconv.Atoi
function to convert the string parameter into an integer. You can see its
documentation with go doc strconv.Atoi. */

// Package main provides a web server which sends a Lissajous gif
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

func main() {
	http.HandleFunc("/", lissajousHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func logReq(r *http.Request) {
	log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		log.Printf("Header[%q] = %q\n", k, v)
	}
	log.Printf("Host = %q\n", r.Host)
	log.Printf("RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		log.Printf("Form[%q] = %q\n", k, v)
	}

	fmt.Println("")
}

func lissajousHandler(w http.ResponseWriter, r *http.Request) {
	logReq(r)
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	lissajous(createLissajousConf(r.Form), w)
}

func lissajous(conf LissaJousConf, out io.Writer) {
	log.Printf("config = %+v\n", conf)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: conf.nframes}
	phase := 0.0 // phase difference
	for i := 0; i < conf.nframes; i++ {
		rect := image.Rect(0, 0, int(2*conf.size+1), int(2*conf.size+1))
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < conf.cycles*2.0*math.Pi; t += conf.res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(int(conf.size+(x*conf.size+0.5)), int(conf.size+(y*conf.size+0.5)),
				blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, conf.delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

// LissaJousConf is input for lissajous
type LissaJousConf struct {
	cycles  float64 // number of complete x oscillator revolutions
	res     float64 // angular resolution
	size    float64 // image canvas covers [-size..+size]
	nframes int     // number of animation frames
	delay   int     // delay between frames in 10ms units
}

func createLissajousConf(form url.Values) LissaJousConf {
	cycles := parseOrDefaultFloat(form, "cycles", 5.0)
	res := parseOrDefaultFloat(form, "res", 0.001)
	size := parseOrDefaultFloat(form, "size", 100)

	nframes := parseOrDefaultInt(form, "nframes", 64)
	delay := parseOrDefaultInt(form, "delay", 8)

	return LissaJousConf{cycles, res, size, nframes, delay}
}

func parseOrDefaultFloat(form url.Values, key string, defVal float64) float64 {
	val, err := strconv.ParseFloat(form.Get(key), 64)
	if err != nil || val == 0.0 {
		val = defVal
	}

	return val
}

func parseOrDefaultInt(form url.Values, key string, defVal int) int {
	val, err := strconv.Atoi(form.Get(key))
	if err != nil || val == 0 {
		val = defVal
	}

	return val
}
