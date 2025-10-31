// Author: Yu-Lun Chen
// Date: 2025-10-24
// Description: Functions using in the BarnesHut simulation.

package main

import (
	"math"
	"os"
	"bufio"
	"strconv"
	"strings"
)

//BarnesHut is our highest level function.
//Input: initial Universe object, a number of generations, and a time interval.
//Output: collection of Universe objects corresponding to updating the system
//over indicated number of generations every given time interval.
func BarnesHut(initialUniverse *Universe, numGens int, time float64, theta float64) []*Universe {
	timePoints := make([]*Universe, numGens + 1)
	timePoints[0] = CopyUniverse(initialUniverse)

	for i := 1; i < (numGens + 1); i++ {
		currentUniverse := timePoints[i-1]
		// for each universe
		// first, build a QuadTree
		tree := GenerateQuadTree(currentUniverse)

		// then we can update the universe
		newUniverse := UpdateUniverse(currentUniverse, time, tree, theta)
		timePoints[i] = newUniverse
	}

    return timePoints
}




//// Functions for Preprocessing the universe: GeneraQuadTree and its subroutines ////

// GenerateQuadTree constructs a QuadTree representation of the given universe.
// It initializes the root node covering the entire universe, inserts all stars
// that are within the universe bounds, and computes the mass and center of mass for each internal node recursively.
// Input: current_universe is a pointer to a Universe struct containing the width and stars.
// Output: a pointer to the constructed QuadTree with the root node.
func GenerateQuadTree(currentUniverse *Universe) *QuadTree {
	// Create root (type: pointer)
	root := &Node{sector: Quadrant{x: 0, y: 0, width: currentUniverse.width}}

	// Insert stars to root (recursively)
	for _, s := range currentUniverse.stars {
		// check if the star s is in the universe
		// Only insert the star if it is in the universe
		if IsInsideUniverse(s, currentUniverse.width) {
			InsertStar(root, s)
		}	
	}

	// After completing building the quadtree, calculate the mass and center position for each internal node
	// This is a recursive function
	ComputeCenterAndMass(root)

    // Create a QuadTree and return the address (type: pointer)
	return &QuadTree{root: root}
}


// InsertStar inserts a star into the given node of the QuadTree, subdividing the node if necessary.
// Input:
//   - node: pointer to the Node in the QuadTree where the star should be inserted.
//   - s: pointer to the Star to be inserted.
// Output:
//   - None (the function modifies the QuadTree in place).
func InsertStar(node *Node, s *Star) {
	// Case 1: no star in this node
	if node.star == nil && len(node.children) == 0 {
		node.star = s

		return
	}

	// Case 2: The node contains a star, need to subdivide
	if len(node.children) == 0 {
		Subdivide(node)
		
		// Copy the old star and insert both old star and new star
		old_star := node.star
		node.star = nil

		InsertStar(node.children[FindQuadrant(node.sector, old_star)], old_star)
		InsertStar(node.children[FindQuadrant(node.sector, s)], s)

		return
	}

	// Case 3: The node has children
	// Directly find the quadrant for the new star and insert it
	idx := FindQuadrant(node.sector, s)
	InsertStar(node.children[idx], s)
}


// Subdivide divide the square into four quadrant(NW, NE, SW, SE) and creates child nodes for each sub-quadrant.
// Input:
//   - node: pointer to the Node to be subdivided. The node's sector is split into four quadrants,
//           and its children field is populated with four new Nodes representing these quadrants.
// Output:
//   - None (modifies the node in place by adding its children).
func Subdivide(node *Node) {
	half := node.sector.width / 2.0
	x := node.sector.x
	y := node.sector.y

	node.children = []*Node{
		&Node{sector: Quadrant{x: x, y: y + half, width: half}},
		&Node{sector: Quadrant{x: x + half, y: y + half, width: half}},
		&Node{sector: Quadrant{x: x, y: y, width: half}},
		&Node{sector: Quadrant{x: x + half, y: y, width: half}},
	}
}


// FindQuadrant determines which quadrant of a sector a given star belongs to.
// Input:
//   - sector: Quadrant representing the current node's region.
//   - s: pointer to the Star to be located.
// Output:
//   - Integer index (0: NW, 1: NE, 2: SW, 3: SE) indicating the quadrant.
func FindQuadrant(sector Quadrant, s *Star) int {
	midX := sector.x + sector.width / 2.0
	midY := sector.y + sector.width / 2.0
	sX := s.position.x 
	sY := s.position.y 

	// NW
	if sX < midX && sY >= midY {
		return 0
	}
	// NE
	if sX >= midX && sY >= midY {
		return 1
	}
	// SW
	if sX < midX && sY < midY {
		return 2
	}
	// SE
	return 3
}


// ComputeCenterAndMass recursively computes the total mass and center of mass for each internal node in the QuadTree.
// Input:
//   - node: pointer to the Node for which to compute mass and center of mass.
// Output:
//   - None (modifies the node in place).
func ComputeCenterAndMass(node *Node) {
	totalMass := 0.0
	xCm, yCm := 0.0, 0.0

	if node == nil {
		return
	}

	if len(node.children) == 0 {
		return
	}

	for _, child := range node.children {
		// Calculate for all children node before calculate for parent nodes
		ComputeCenterAndMass(child)

		// Calculate for parent node (current node) with results from children nodes
		if child.star != nil {
			m := child.star.mass
			totalMass += m 
			xCm += m * child.star.position.x 
			yCm += m * child.star.position.y
		}
	}


	if totalMass > 0 {
		node.star = &Star{
			position: OrderedPair{x: xCm / totalMass, y: yCm / totalMass},
			mass: totalMass,
		}
	}
}


// IsInsideUniverse checks if a star is within the bounds of the universe.
// Input:
//   - s: pointer to the Star to check.
//   - width: width of the universe.
// Output:
//   - Boolean indicating whether the star is inside the universe.
func IsInsideUniverse(s *Star, width float64) bool {
	return s.position.x >= 0 && s.position.x <= width && s.position.y >= 0 && s.position.y <= width
}


// CalculateNetForce computes the net force on a star using the Barnes-Hut approximation.
// Input:
//   - node: pointer to the current Node in the QuadTree.
//   - curr_star: pointer to the Star for which to calculate the force.
//   - theta: threshold parameter for Barnes-Hut approximation.
// Output:
//   - OrderedPair representing the net force vector.
func CalculateNetForce(node *Node, currStar *Star,theta float64) OrderedPair {
    var force OrderedPair

	// no force cases
	if node == nil || node.star == nil || node.star.mass == 0 {
		return force
	}

	// if it is a leaf and contains a real star: calculate the force
	if IsLeaf(node) && node.star != nil && node.star != currStar {
		dX, dY, d := Distance(node.star.position, currStar.position)
		if d != 0 {
			f := G  * currStar.mass * node.star.mass / (d * d)
			fX := f * (dX / d)
			fY := f * (dY / d)

			force.x += fX
			force.y += fY	
		}
		return force
	}

	
	if node.star != currStar && node.star != nil {
		_, _, d := Distance(node.star.position, currStar.position)

		if d != 0 {
			s := node.sector.width
			if (s/d) < theta {
				// far enough to be a dummy body
				// we do not consider the force given by dummy star
				force.x += 0.0
				force.y += 0.0
			}
		}		
	}

	// if d is too small, indicating the node should be expanded
	// expand the node and run recursively on their children
	if node.children != nil {
		for _, child := range node.children {
			if child != nil {
				f := CalculateNetForce(child, currStar, theta)
				force.x += f.x
				force.y += f.y 				
			}
		}
	}

    return force
}


// ComputeForce calculates the gravitational force between two stars.
// Input:
//   - b: pointer to the first Star.
//   - b2: pointer to the second Star.
// Output:
//   - OrderedPair representing the force vector.
func ComputeForce(b, b2 *Star) OrderedPair{
	var force OrderedPair

	dX, dY, d := Distance(b.position, b2.position)
	
	// check if denominator == 0
	if d == 0.0 {
		return force
	}
	F := (G * b.mass * b2.mass) / (d * d)

	force.x = F * dX/d 
	force.y = F * dY/d

	return force
}


// Distance computes the difference in x, y, and Euclidean distance between two points.
// Input:
//   - p1: first OrderedPair.
//   - p2: second OrderedPair.
// Output:
//   - delta_x, delta_y, and Euclidean distance between p1 and p2.
func Distance(p1, p2 OrderedPair) (float64, float64, float64) {
	// this is the distance formula from days of precalculus long ago ...
	deltaX := p1.x - p2.x
	deltaY := p1.y - p2.y
	return deltaX, deltaY, math.Sqrt(deltaX * deltaX + deltaY * deltaY)
}


// IsLeaf checks if a node is a leaf node (has no children).
// Input:
//   - node: pointer to the Node to check.
// Output:
//   - Boolean indicating if the node is a leaf.
func IsLeaf(node *Node) bool {
	for _, child := range node.children {
		if child != nil {
			return false
		}
	}
	return true
}




//// subroutines for the higest function BarnesHut ////

// UpdateUniverse updates the positions, velocities, and accelerations of all stars in the universe for one timestep.
// Input:
//   - current_universe: pointer to the current Universe.
//   - time: time interval for the update.
//   - tree: pointer to the QuadTree representing the current universe.
//   - theta: threshold parameter for Barnes-Hut approximation.
// Output:
//   - Pointer to the updated Universe.
func UpdateUniverse(currentUniverse *Universe, time float64, tree *QuadTree, theta float64) *Universe{
	newUniverse := CopyUniverse(currentUniverse)

	for i, b := range newUniverse.stars {
		oldAcceleration, oldVelocity := b.acceleration, b.velocity

		newUniverse.stars[i].acceleration = UpdateAcceleration(b, tree, theta)
		newUniverse.stars[i].velocity = UpdateVelocity(newUniverse.stars[i], oldAcceleration, time)
		newUniverse.stars[i].position = UpdatePosition(newUniverse.stars[i], oldAcceleration, oldVelocity, time)
	}

	return newUniverse
}


// UpdateAcceleration computes the new acceleration for a star based on the net force from the QuadTree.
// Input:
//   - s: pointer to the Star.
//   - tree: pointer to the QuadTree.
//   - theta: threshold parameter for Barnes-Hut approximation.
// Output:
//   - OrderedPair representing the new acceleration.
func UpdateAcceleration(s *Star, tree *QuadTree, theta float64) OrderedPair {
	var accel OrderedPair

	// calculate the net force with QuadTree and the given theta
	force := CalculateNetForce(tree.root, s, theta)
	accel.x = force.x / s.mass
	accel.y = force.y / s.mass

	return accel
}


// UpdateVelocity updates the velocity of a star using the previous and current acceleration.
// Input:
//   - s: pointer to the Star.
//   - old_acceleration: OrderedPair of the previous acceleration.
//   - time: time interval for the update.
// Output:
//   - OrderedPair representing the new velocity.
func UpdateVelocity(s *Star, oldAcceleration OrderedPair, time float64) OrderedPair {
	var velo OrderedPair

	velo.x = s.velocity.x + 0.5 * (s.acceleration.x + oldAcceleration.x) * time
	velo.y = s.velocity.y + 0.5 * (s.acceleration.y + oldAcceleration.y) * time

	return velo
}


// UpdatePosition updates the position of a star using its previous acceleration and velocity.
// Input:
//   - s: pointer to the Star.
//   - old_acceleration: OrderedPair of the previous acceleration.
//   - old_velocity: OrderedPair of the previous velocity.
//   - time: time interval for the update.
// Output:
//   - OrderedPair representing the new position.
func UpdatePosition(s *Star, oldAcceleration, oldVelocity OrderedPair, time float64) OrderedPair {
	var pos OrderedPair

	pos.x = s.position.x + oldVelocity.x * time + 0.5 * oldAcceleration.x * time * time
	pos.y = s.position.y + oldVelocity.y * time + 0.5 * oldAcceleration.y * time * time

	return pos
}


// CopyUniverse creates a deep copy of the given Universe.
// Input:
//   - u: pointer to the Universe to copy.
// Output:
//   - Pointer to the new, copied Universe.
func CopyUniverse(u *Universe) *Universe {
	newUniverse := &Universe{width: u.width}

	for _, s := range u.stars {
		copy_s := &Star{
			position: OrderedPair{x: s.position.x, y: s.position.y},
			velocity: OrderedPair{x: s.velocity.x, y: s.velocity.y},
			acceleration: OrderedPair{x: s.acceleration.x, y: s.acceleration.y},
			mass: s.mass,
			radius: s.radius,
			red: s.red,
			blue: s.blue,
			green: s.green,
		}
		
		newUniverse.stars = append(newUniverse.stars, copy_s)
	}

	return newUniverse
}




//// Load data from jupiterMoons.txt ////

// LoadJupiterMoons loads star data from a file and constructs a Universe.
// Input:
//   - file_name: string path to the data file.
// Output:
//   - Pointer to the constructed Universe.
func LoadJupiterMoons(file_name string) *Universe {
	file, err := os.Open(file_name)
	Check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	
	var lines []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		if line != "" {
			lines = append(lines, line)
		}
	}

	width, err := strconv.ParseFloat(lines[0], 64)
	Check(err)

	u := &Universe {
		width: width,
		stars: make([]*Star, 0),
	}

	var currStar *Star

	for i := 2; i < len(lines); i++ {
		line := lines[i]

		if strings.HasPrefix(line, ">") {
			
			// add the previous moon to universe
			if currStar != nil {
				u.stars = append(u.stars, currStar)
			}
			// start manage the current moon
			currStar = &Star{}
			continue
		}

		// if it is the first moon
		if currStar == nil {
			continue
		}

		// manage color information
		if strings.Count(line, ",") == 2 {
			fields := strings.Split(line, ",")
			r, _ := strconv.Atoi(strings.TrimSpace(fields[0]))
			g, _ := strconv.Atoi(strings.TrimSpace(fields[1]))
			b, _ := strconv.Atoi(strings.TrimSpace(fields[2]))
			currStar.red = uint8(r)
			currStar.green = uint8(g)
			currStar.blue = uint8(b)
			continue
		}

		// mamage position, velocity
		if strings.Contains(line, ",") && strings.Count(line, ",") == 1 {
			
			fields := strings.Split(line, ",")
			x, _ := strconv.ParseFloat(strings.TrimSpace(fields[0]), 64)
			y, _ := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)

			if currStar.position == (OrderedPair{}) {
				currStar.position = OrderedPair{x, y}
			} else {
				currStar.velocity = OrderedPair{x, y}
			}
			continue
		}

		// manage mass, radius
		val, _ := strconv.ParseFloat(line, 64)
		if currStar.mass == 0.0 {
			currStar.mass = val
		} else {
			currStar.radius = val
		}
	}

	// add the last moon to the universe
	if currStar != nil {
		u.stars = append(u.stars, currStar)
	}

	return u
}




//// Push functions for pushing galaxies in collision command ////

// GalaxyPush applies a velocity "push" to two galaxies in opposite directions along the line connecting their centers.
// Input:
//   - g0: first Galaxy (slice of *Star).
//   - g1: second Galaxy (slice of *Star).
//   - v: magnitude of the velocity to apply.
// Output:
//   - None (modifies the velocities of the stars in place).
func GalaxyPush(g0, g1 Galaxy, v float64) {
	// center of the galaxies is needed for computing the distance
	center_0 := GalaxyCenter(g0)
	center_1 := GalaxyCenter(g1)

	d_x := center_1.x - center_0.x
	d_y := center_1.y - center_0.y 
	distance := math.Sqrt(d_x * d_x + d_y * d_y)

	// if two galaxies are at same position
	if distance == 0 {
		// slightly change the position
		d_x, d_y = 1e-3, 0
		distance = 1e-3
	}

	// else, simply calculate the pushing direction and velocity
	// the pushing directions for two galaxies are opposite.
	dir_0 := OrderedPair{d_x / distance, d_y / distance}
	dir_1 := OrderedPair{-d_x / distance, -d_y / distance}

	// update the velocities
	for _, s := range g0 {
		s.velocity.x += v * dir_0.x
		s.velocity.y += v * dir_0.y
	}

	for _, s := range g1 {
		s.velocity.x += v * dir_1.x
		s.velocity.y += v * dir_1.y
	}

}


// GalaxyCenter computes the center (average position) of a galaxy.
// Input:
//   - g: Galaxy (slice of *Star).
// Output:
//   - OrderedPair representing the center position.
func GalaxyCenter(g Galaxy) OrderedPair {
	var c_x, c_y float64

	for _, s := range g {
		c_x += s.position.x 
		c_y += s.position.y 
	}
	n := float64(len(g))

	return OrderedPair{x: c_x / n, y: c_y / n}
}
