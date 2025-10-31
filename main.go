// Author: Yu-Lun Chen
// Date: 2025-10-24
// Description: Main code for running and visualizing the universe simulation.

package main

import (
	"fmt"
	"gifhelper"
	"os"
)

// main is the entry point of the Barnes-Hut simulation program
func main() {
	// read parameters from command line
	// the command should be: ./BarnesHut "jupiter/galaxy/collision"
	// as mention on cogniterra
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./BarnesHut [jupiter|galaxy|collision]")
		os.Exit(1)
	}

	command := os.Args[1]

	// initialize parameters, will be customerized for each command
	width := 0.0
	numGens := 0
	time := 0.0
	theta := 0.0

	canvasWidth := 0
	frequency := 0
	scalingFactor := 0.0

	var initialUniverse *Universe

	// set different parameters for different command
	switch command {

	// set parameters for argument "jupiter"
	case "jupiter":
		// The "jupiter" scenario uses much smaller parameters (such as width, time, and scaling factors) 
		// because Jupiter's moons occur on a much smaller spatial and temporal scale than galactic interactions.
		width = 1.0e23
		numGens = 100000
		time = 1e1
		theta = 0.5

		canvasWidth = 1000
		frequency = 1000
		scalingFactor = 5.0

		// "Data/jupiterMoons.txt" is copy from "ProgrammingforScientists2025Grad/Starter_Code/gravity/data"
		initialUniverse = LoadJupiterMoons("Data/jupiterMoons.txt")
		fmt.Println("Loaded", len(initialUniverse.stars), "bodies from file.")
		for _, s := range initialUniverse.stars {
    		fmt.Printf("star at (%.2f, %.2f)\n", s.position.x, s.position.y)
			fmt.Printf("star velocity (%.2f, %.2f)\n", s.velocity.x, s.velocity.y)
			fmt.Printf("star mass (%.2f)\n", s.mass)
			fmt.Printf("star radius (%.2f)\n", s.radius)
		}
		

	// set parameters for argument "galaxy"
	case "galaxy":
		width = 1.0e23
		numGens = 100000
		time = 2e15
		theta = 0.5

		canvasWidth = 1000
		frequency = 1000
		scalingFactor = 5e11

		g := InitializeGalaxy(500, 1e22, 5e22, 5e22)
		initialUniverse = InitializeUniverse([]Galaxy{g}, width)

	// set parameters for argument "collision"
	case "collision":
		width = 1.0e23
		numGens = 100000
		time = 2e14
		theta = 0.5

		canvasWidth = 1000
		frequency = 1000
		scalingFactor = 1e11
		// the following sample parameters may be helpful for the "collide" command
		// all units are in SI (meters, kg, etc.)
		// but feel free to change the positions of the galaxies.

		g0 := InitializeGalaxy(500, 4e21, 7e22, 2e22)
		g1 := InitializeGalaxy(500, 4e21, 3e22, 7e22)

		// you probably want to apply a "push" function at this point to these galaxies to move
		// them toward each other to collide.
		// be careful: if you push them too fast, they'll just fly through each other.
		// too slow and the black holes at the center collide and hilarity ensues.

		// Push galaxy by simple push function
		v := 5e3      // 5e3 found to be a proper speed value after multiple tests
		GalaxyPush(g0, g1, v)

		galaxies := []Galaxy{g0, g1}
		initialUniverse = InitializeUniverse(galaxies, width)

	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)

	}

	// === Run Simulation ===
	timePoints := BarnesHut(initialUniverse, numGens, time, theta)

	fmt.Println("Simulation run. Now drawing images.")

	imageList := AnimateSystem(timePoints, canvasWidth, frequency, scalingFactor)

	fmt.Println("Images drawn. Now generating GIF.")
	gifhelper.ImagesToGIF(imageList, "galaxy")
	fmt.Println("GIF drawn.")
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

