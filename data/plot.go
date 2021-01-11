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

	estXVelJSON, _ := ioutil.ReadFile("estXVelJSON.json")
	estYVelJSON, _ := ioutil.ReadFile("estYVelJSON.json")
	absXVelJSON, _ := ioutil.ReadFile("absXVelJSON.json")
	absYVelJSON, _ := ioutil.ReadFile("absYVelJSON.json")

	ballEstJSON, _ := ioutil.ReadFile("ballEstJSON.json")
	ballAbsJSON, _ := ioutil.ReadFile("ballAbsJSON.json")

	var estPoints, absPoints [][]float64
	var estTPoints, absTPoints []float64

	var estXVelPoints, absXVelPoints []float64
	var estYVelPoints, absYVelPoints []float64

	var ballEstPoints, ballAbsPoints [][]float64

	var err error
	// Unmarshal self pos
	err = json.Unmarshal(estJSON, &estPoints)
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

	// Unmarshal self vel
	err = json.Unmarshal(estXVelJSON, &estXVelPoints)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(estYVelJSON, &estYVelPoints)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(absXVelJSON, &absXVelPoints)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(absYVelJSON, &absYVelPoints)
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

	// Plot self position
	plotPos, _ := glot.NewPlot(2, false, false)

	plotPos.SetTitle("Self Pos")

	plotPos.AddPointGroup("estimated", "lines", estPoints)
	plotPos.AddPointGroup("absolute", "lines", absPoints)

	plotPos.SavePlot("plotPos.png")

	// Plot Self Angle
	plotAng, _ := glot.NewPlot(1, false, false)
	plotAng.SetTitle("Self Angle")

	plotAng.AddPointGroup("estimated", "lines", estTPoints)
	plotAng.AddPointGroup("absolute", "lines", absTPoints)

	plotAng.SavePlot("plotAng.png")

	// Plot self X Vel
	plotXVel, _ := glot.NewPlot(1, false, false)
	plotXVel.SetTitle("Self X Velocity")

	plotXVel.AddPointGroup("estimated", "lines", estXVelPoints)
	plotXVel.AddPointGroup("absolute", "lines", absXVelPoints)

	plotXVel.SavePlot("plotXVel.png")

	// Plot self X Vel
	plotYVel, _ := glot.NewPlot(1, false, false)
	plotYVel.SetTitle("Self Y Velocity")

	plotYVel.AddPointGroup("estimated", "lines", estYVelPoints)
	plotYVel.AddPointGroup("absolute", "lines", absYVelPoints)

	plotYVel.SavePlot("plotYVel.png")

	// Plot Ball X
	plotBallX, _ := glot.NewPlot(1, false, false)
	plotBallX.SetTitle("Ball Pos X")

	plotBallX.AddPointGroup("estimated", "lines", ballEstPoints[0])
	plotBallX.AddPointGroup("absolute", "lines", ballAbsPoints[0])

	plotBallX.SavePlot("plotBallX.png")

	// Plot Ball Y
	plotBallY, _ := glot.NewPlot(1, false, false)
	plotBallY.SetTitle("Ball Pos Y")

	plotBallY.AddPointGroup("estimated", "lines", ballEstPoints[1])
	plotBallY.AddPointGroup("absolute", "lines", ballAbsPoints[1])

	plotBallY.SavePlot("plotBallY.png")
}
