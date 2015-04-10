package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"os"

	"github.com/gamedevja/backend/git"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/testpush", testpush)
	fmt.Println("listening...")
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		panic(err)
	}
}

func index(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "index")
}

func testpush(w http.ResponseWriter, r *http.Request) {
	m := image.NewRGBA(image.Rect(0, 0, 40, 40))
	c := color.RGBA{0, 0, 0, 255}
	draw.Draw(m, image.Rect(m.Rect.Min.X, m.Rect.Min.Y, m.Rect.Max.X, m.Rect.Max.Y), &image.Uniform{c}, image.ZP, draw.Src)

	var b bytes.Buffer
	if err := png.Encode(&b, m); err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	bin := b.String()
	if err := git.Push(&bin, "assets/testimage/40x40.png", "remote"); err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	text := "# md\nmd sample alt\n"
	if err := git.Push(&text, "assets/testimage/sample.md", "remote"); err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	fmt.Fprintln(w, "done")
}
