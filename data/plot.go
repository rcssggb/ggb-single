package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Arafatk/glot"
)

func main() {

	estJSON, _ := ioutil.ReadFile("estJSON.json")
	absJSON, _ := ioutil.ReadFile("absJSON.json")
	estTJSON, _ := ioutil.ReadFile("estTJSON.json")
	absTJSON, _ := ioutil.ReadFile("absTJSON.json")

	var estPoints [][]float64
	var absPoints [][]float64
	var estTPoints []float64
	var absTPoints []float64

	err := json.Unmarshal(estJSON, &estPoints)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(absJSON, &absPoints)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(estTJSON, &estTPoints)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(absTJSON, &absTPoints)
	if err != nil {
		fmt.Println(err)
		return
	}

	plotPos, _ := glot.NewPlot(2, false, false)

	plotPos.AddPointGroup("estimated", "lines", estPoints)
	plotPos.AddPointGroup("absolute", "lines", absPoints)

	plotPos.SavePlot("plotPos.png")

	plotAng, _ := glot.NewPlot(1, false, false)

	plotAng.AddPointGroup("estimated", "lines", estTPoints)
	plotAng.AddPointGroup("absolute", "lines", absTPoints)

	plotAng.SavePlot("plotAng.png")
}
