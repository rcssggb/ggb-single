package main

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/Arafatk/glot"
)

func main() {
	var err error

	filename := "returns.rln"
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	var returns []float32
	err = dec.Decode(&returns)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Plot self position
	plotReturns, _ := glot.NewPlot(1, true, false)

	plotReturns.SetTitle("Return per Episode")

	plotReturns.AddPointGroup("returns", "lines", returns)

	plotReturns.SavePlot("plotReturns.png")
}
