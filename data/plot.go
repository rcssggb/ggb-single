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

	ballEstJSON, _ := ioutil.ReadFile("ballEstJSON.json")
	ballAbsJSON, _ := ioutil.ReadFile("ballAbsJSON.json")

	var estPoints, absPoints [][]float64
	var estTPoints, absTPoints []float64

	var ballEstPoints, ballAbsPoints [][]float64

	// Unmarshal self pos
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

	// Unmarshal ball pos
	err = json.Unmarshal(ballEstJSON, &ballEstPoints)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(ballAbsJSON, &ballAbsPoints)
	if err != nil {
		fmt.Println(err)
		return
	}

	plotPos, _ := glot.NewPlot(2, false, false)

	plotPos.SetTitle("Self Pos")

	plotPos.AddPointGroup("estimated", "lines", estPoints)
	plotPos.AddPointGroup("absolute", "lines", absPoints)

	plotPos.SavePlot("plotPos.png")

	plotAng, _ := glot.NewPlot(1, false, false)
	plotAng.SetTitle("Self Angle")

	plotAng.AddPointGroup("estimated", "lines", estTPoints)
	plotAng.AddPointGroup("absolute", "lines", absTPoints)

	plotAng.SavePlot("plotAng.png")

	plotBallPos, _ := glot.NewPlot(2, false, false)
	plotBallPos.SetTitle("Ball Pos")

	plotBallPos.AddPointGroup("estimated", "lines", ballEstPoints)
	plotBallPos.AddPointGroup("absolute", "lines", ballAbsPoints)

	plotBallPos.SavePlot("plotBallPos.png")
}
