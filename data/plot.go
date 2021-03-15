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
	var returns []float64
	err = dec.Decode(&returns)
	if err != nil {
		fmt.Println(err)
		return
	}

	const windowSize = 50
	var returnMeans []float64
	var acc float64
	k := 0
	for _, r := range returns {
		if k < windowSize {
			k++
			acc += r
		} else {
			returnMeans = append(returnMeans, acc/windowSize)
			k = 0
			acc = 0
		}
	}

	// Plot return per episode
	plotReturns, _ := glot.NewPlot(1, false, false)

	plotReturns.SetTitle("Return per Episode")

	plotReturns.AddPointGroup("returns", "lines", returns)

	plotReturns.SavePlot("plotReturns.png")

	// Plot avg return per windowSize (default 50) episodes
	plotMeanReturns, _ := glot.NewPlot(1, false, false)

	plotMeanReturns.SetTitle(fmt.Sprintf("Avg Return per %d Episodes", windowSize))

	plotMeanReturns.AddPointGroup("returns", "lines", returnMeans)

	plotMeanReturns.SavePlot("plotMeanReturns.png")

}
