package main

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/raster"
	"code.google.com/p/freetype-go/freetype/truetype"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const WIDTH int = 450
const HEIGHT int = 200
const FS int = 30

var font *truetype.Font

func drawHandler(rw http.ResponseWriter, req *http.Request) {
	qry, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing query: {}", err.Error())
		return
	}
	var str string
	width, height, fs := WIDTH, HEIGHT, FS
	if len(qry["txt"]) > 0 {
		str = qry["txt"][0]
	} else {
		str = "Plz insert text via 'txt' GET param"
	}
	if len(qry["w"]) > 0 {
		width, _ = strconv.Atoi(qry["w"][0])
	}
	if len(qry["h"]) > 0 {
		height, _ = strconv.Atoi(qry["h"][0])
	}
	if len(qry["fs"]) > 0 {
		fs, _ = strconv.Atoi(qry["fs"][0])
	}
	drawString(rw, str, width, height, fs)
}

func drawString(w io.Writer, str string, width, height, fs int) {
	bg, fg := image.Black, image.White
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetFont(font)
	c.SetFontSize(float64(fs))
	curs := freetype.Pt(fs, fs)
	scale := int(curs.X) / fs
	for i, _ := range str {
		if int(curs.X)/scale > width-fs {
			curs.Y += raster.Fix32(fs * scale)
			curs.X = raster.Fix32(fs * scale)
			if str[i] == ' ' {
				continue
			}
		}
		off, _ := c.DrawString(str[i:i+1], curs)
		curs.X = off.X
		i++
	}
	png.Encode(w, rgba)
}

func main() {
	http.HandleFunc("/", drawHandler)
	var fontFile string
	if len(os.Args) > 1 {
		fontFile = os.Args[1]
	} else {
		fontFile = "/usr/share/cups/fonts/FreeMono.ttf"
	}
	f, err := ioutil.ReadFile(fontFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading file: {}", err.Error())
		return
	}
	font, err = freetype.ParseFont(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing font: {}", err.Error())
		return
	}
	http.ListenAndServe(":9090", nil)
}
