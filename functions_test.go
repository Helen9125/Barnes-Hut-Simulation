// Author: Yu-Lun Chen
// Date: 2025-10-24
// Description: Testing functions for eight subroutines in function.go.
// There are at least four testing cases for each test functions (directory: Tests/[function_name].txt)
// Each txt file contains input testing cases and the expected output for each cases.

package main

import (
	"bufio"
	"os"
	"math"
	"strconv"
	"strings"
	"testing"
)




//// Difinition for some struct used in testing ////

type SubdivideTestCase struct {
    node *Node
    expected   [4]Quadrant
}

type IsInsideTestCases struct {
	star Star
	width float64
	expected bool
}

type ComputeCenterAndMassTestCase struct {
	node          *Node
	expectedX     float64
	expectedY     float64
	expectedMass  float64
}

type IsLeafTestCases struct {
	id string
	children []*Node
	expected bool
}

type DistanceTestCases struct {
	id string
	x1, y1, x2, y2 float64
	expectedDeltaX, expectedDeltaY, expectedDistance float64
}

type VelocityTestCases struct {
	id string
	star Star
	oldAcceleration OrderedPair
	time float64
	expected OrderedPair
}

type PositionTestCases struct {
	id string
	star Star
	oldAcceleration OrderedPair
	oldVelocity OrderedPair
	time float64
	expected OrderedPair
}




//// Functions for reading testing data from txt files ////

// ReadFindQuadrantData reads test data for the FindQuadrant function from a file.
// Input: filename (string) - path to the test data file.
// Output: slice of pointers to Star, width (float64), and slice of expected quadrant indices ([]int).
func ReadFindQuadrantData(fileName string) ([]*Star, float64, []int) {
	file, err := os.Open(fileName)
	Check(err)
	defer file.Close()

    scanner := bufio.NewScanner(file)
	var stars []*Star
    var width float64
    var expected []int
    var lineCount int
    readingExpected := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		lineCount ++

		// read width
		if width == 0.0 {
			width, _ = strconv.ParseFloat(line, 64)
			continue
		}

		// we are in to expected result when reading_expected is True
		if readingExpected {
			val, _ := strconv.Atoi(line)
			expected = append(expected, val)
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 1 && (parts[0] == "0" || parts[0] == "1" || parts[0] == "2" || parts[0] == "3") {
			// we go in to expected result region
			readingExpected = true
            val, _ := strconv.Atoi(parts[0])
            expected = append(expected, val)
            continue

		}

		// reading star information
		if len(parts) == 9 {
			x, _ := strconv.ParseFloat(parts[0], 64)
        	y, _ := strconv.ParseFloat(parts[1], 64)
        	vx, _ := strconv.ParseFloat(parts[2], 64)
        	vy, _ := strconv.ParseFloat(parts[3], 64)
        	m, _ := strconv.ParseFloat(parts[4], 64)
        	r, _ := strconv.ParseFloat(parts[5], 64)
        	red, _ := strconv.Atoi(parts[6])
        	green, _ := strconv.Atoi(parts[7])
        	blue, _ := strconv.Atoi(parts[8])

        	s := &Star{
            	position: OrderedPair{x, y},
            	velocity: OrderedPair{vx, vy},
            	mass:     m,
            	radius:   r,
            	red:      uint8(red),
            	green:    uint8(green),
            	blue:     uint8(blue),
        	}
        	stars = append(stars, s)
		}
	}	

	return stars, width, expected
}


// ReadSubdivideData reads test data for the Subdivide function from a file.
// Input: filename (string) - path to the test data file.
// Output: slice of SubdivideTestCase structs containing nodes and expected quadrants.
func ReadSubdivideData(fileName string) []SubdivideTestCase {
	file, err := os.Open(fileName)
	Check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var tests []SubdivideTestCase
	var n *Node
	var expected [4]Quadrant
	childIndex := 0
	readingExpected := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if ! readingExpected {
			// read node sector for current node
			if len(parts) == 3 {
				x, _ := strconv.ParseFloat(parts[0], 64)
                y, _ := strconv.ParseFloat(parts[1], 64)
                width, _ := strconv.ParseFloat(parts[2], 64)

				n = &Node{
					sector: Quadrant{x, y, width},
					children: nil,
				}
                readingExpected = true
				childIndex = 0
			}
		} else {
			// read expected result
			if len(parts) == 3 && childIndex < 4 {
				x, _ := strconv.ParseFloat(parts[0], 64)
                y, _ := strconv.ParseFloat(parts[1], 64)
                width, _ := strconv.ParseFloat(parts[2], 64)
                expected[childIndex] = Quadrant{x, y, width}
                childIndex ++
			}
		}

		// finish reading all expected results
		if childIndex == 4 {
			tests = append(tests, SubdivideTestCase{
				node: n,
				expected: expected,
			})
			// set reading_expected back to false to read next test data
			readingExpected = false
		}
	}

	return tests
}


// ReadIsInsideUniverse reads test data for the IsInsideUniverse function from a file.
// Input: file_name (string) - path to the test data file.
// Output: slice of IsInsideTestCases structs containing star, width, and expected result.
func ReadIsInsideUniverse(fileName string) []IsInsideTestCases {
	file, err := os.Open(fileName)
	Check(err)
	defer file.Close()

	var tests []IsInsideTestCases
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 4 {
			continue
		}

		x, _ := strconv.ParseFloat(parts[0], 64)
		y, _ := strconv.ParseFloat(parts[1], 64)
		width, _ := strconv.ParseFloat(parts[2], 64)
		expected, _ := strconv.ParseBool(parts[3])

		tests = append(tests, IsInsideTestCases{
			star: Star{
				position: OrderedPair{x, y},
			},
			width: width,
			expected: expected,
		})
	}

	return tests
}


// ReadComputeCenterAndMass reads test data for the ComputeCenterAndMass function from a file.
// Input: file_name (string) - path to the test data file.
// Output: slice of ComputeCenterAndMassTestCase structs containing node and expected center/mass values.
func ReadComputeCenterAndMass(fileName string) []ComputeCenterAndMassTestCase {
	file, err := os.Open(fileName)
	Check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var tests []ComputeCenterAndMassTestCase

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// to check if it is leaf or internal node
		if strings.Contains(line, "|") {
			// is internal node
			parts := strings.Split(line, "|")

			// first, get the expected values
			expectedParts := strings.Fields(parts[len(parts) - 1])
			expectedX, _ := strconv.ParseFloat(expectedParts[0], 64)
			expectedY, _ := strconv.ParseFloat(expectedParts[1], 64)
			expectedMass, _ := strconv.ParseFloat(expectedParts[2], 64)

			var children []*Node

			// extract information for children nodes
			for _, childPart := range parts[: len(parts) - 1] {
				fields := strings.Fields(childPart)

				if len(fields) != 3 {
					continue
				}

				x, _ := strconv.ParseFloat(fields[0], 64)
				y, _ := strconv.ParseFloat(fields[1], 64)
				mass, _ := strconv.ParseFloat(fields[2], 64)

				child := &Node{
					star: &Star{
						position: OrderedPair{x, y},
						mass: mass,
					},
				}
				children = append(children, child)
			}

			root := &Node{children: children}

			tests = append(tests, ComputeCenterAndMassTestCase{
				node: root,
				expectedX: expectedX,
				expectedY: expectedY,
				expectedMass: expectedMass,
			})
		} else {
			// is leaf
			parts := strings.Fields(line)
			x, _ := strconv.ParseFloat(parts[0], 64)
			y, _ := strconv.ParseFloat(parts[1], 64)
			mass, _ := strconv.ParseFloat(parts[2], 64)

			// first, get the expected value
			expectedX, _ := strconv.ParseFloat(parts[3], 64)
			expectedY, _ := strconv.ParseFloat(parts[4], 64)
			expectedMass, _ := strconv.ParseFloat(parts[5], 64)

			// extract value for node itself
			tests = append(tests, ComputeCenterAndMassTestCase{
				node: &Node{
					star: &Star{
							position: OrderedPair{x, y},
							mass: mass,}},
				expectedX: expectedX,
				expectedY: expectedY,
				expectedMass: expectedMass,
			})
		}
	}

	return tests
}


// ReadIsLeaf reads test data for the IsLeaf function from a file.
// Input: file_name (string) - path to the test data file.
// Output: slice of IsLeafTestCases structs containing node children and expected boolean result.
func ReadIsLeaf(fileName string) []IsLeafTestCases {
	file, err := os.Open(fileName)
	Check(err)
	defer file.Close()

	var tests []IsLeafTestCases
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "|")

		if len(parts) != 3 {
			continue
		}

		id := strings.TrimSpace(parts[0])
		childrenStr := strings.TrimSpace(parts[1])
		expectedStr := strings.TrimSpace(parts[2])

		// make children slice []*Node
		children := make([]*Node, 4)
		if strings.Contains(childrenStr, "Node") {
			items := strings.Split(strings.Trim(childrenStr, "[]"), ",")
			for i, item := range items {
				item = strings.TrimSpace(item)
				if item == "Node" {
					 // give a non-nil Node if string is Node
					children[i] = &Node{}
				} else {
					children[i] = nil
				}
			}
		}

		expected := false
		if strings.Contains(expectedStr, "true") {
			expected = true
		}

		tests = append(tests, IsLeafTestCases{
			id:       id,
			children: children,
			expected: expected,
		})

	}

	return tests
}


// ReadDistance reads test data for the Distance function from a file.
// Input: file_name (string) - path to the test data file.
// Output: slice of DistanceTestCases structs containing points and expected deltas/distances.
func ReadDistance(fileName string) []DistanceTestCases {
	file, err := os.Open(fileName)
	Check(err)
	defer file.Close()

	var tests []DistanceTestCases
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 3 {
			continue
		}

		id := strings.TrimSpace(parts[0])
		points := strings.Fields(strings.TrimSpace(parts[1]))
		expectedParts := strings.Fields(parts[2])

		if len(points) != 4 {
			continue
		}

		x1, err := strconv.ParseFloat(points[0], 64)
		Check(err)
		y1, err := strconv.ParseFloat(points[1], 64)
		Check(err)
		x2, err := strconv.ParseFloat(points[2], 64)
		Check(err)
		y2, err := strconv.ParseFloat(points[3], 64)
		Check(err)
		expectedDeltaX, err := strconv.ParseFloat(expectedParts[0], 64)
		Check(err)
		expectedDeltaY, err := strconv.ParseFloat(expectedParts[1], 64)
		Check(err)
		expectedDistance, err := strconv.ParseFloat(expectedParts[2], 64)
		Check(err)

		tests = append(tests, DistanceTestCases{
			id:       id,
			x1:       x1,
			y1:       y1,
			x2:       x2,
			y2:       y2,
			expectedDeltaX: expectedDeltaX,
			expectedDeltaY: expectedDeltaY,
			expectedDistance: expectedDistance,
		})
	}
	return tests
}


// ReadVelocity reads test data for the UpdateVelocity function from a file.
// Input: file_name (string) - path to the test data file.
// Output: slice of VelocityTestCases structs containing star, old acceleration, time, and expected velocity.
func ReadVelocity(fileName string) []VelocityTestCases {
	file, err := os.Open(fileName)
	Check(err)
	defer file.Close()

	var tests []VelocityTestCases

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

        if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
            continue
        }
		
		parts := strings.Fields(line)

		id := strings.TrimSpace(parts[0])
		vx, err := strconv.ParseFloat(parts[1], 64)
		Check(err)
		vy, err := strconv.ParseFloat(parts[2], 64)
		Check(err)
		ax, err := strconv.ParseFloat(parts[3], 64)
		Check(err)
		ay, err := strconv.ParseFloat(parts[4], 64)
		Check(err)
		oldAx, err := strconv.ParseFloat(parts[5], 64)
		Check(err)
		oldAy, err := strconv.ParseFloat(parts[6], 64)
		Check(err)
		t, err := strconv.ParseFloat(parts[7], 64)
		Check(err)
		expVx, err := strconv.ParseFloat(parts[8], 64)
		Check(err)
		expVy, err := strconv.ParseFloat(parts[9], 64)
		Check(err)

		test := VelocityTestCases{
			id: id,
			star: Star{
				velocity: OrderedPair{vx, vy},
				acceleration: OrderedPair{ax, ay},
			},
			oldAcceleration: OrderedPair{oldAx, oldAy},
			time: t,
			expected: OrderedPair{expVx, expVy},
		}
		tests = append(tests, test)
	}
	return tests
}


// ReadPosition reads test data for the UpdatePosition function from a file.
// Input: file_name (string) - path to the test data file.
// Output: slice of PositionTestCases structs containing star, old acceleration, old velocity, time, and expected position.
func ReadPosition(fileName string) []PositionTestCases {
	file, err := os.Open(fileName)
	Check(err)
	defer file.Close()

	var tests []PositionTestCases

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

        if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
            continue
        }
		
		parts := strings.Fields(line)

		id := strings.TrimSpace(parts[0])
		px, err := strconv.ParseFloat(parts[1], 64)
		Check(err)
		py, err := strconv.ParseFloat(parts[2], 64)
		Check(err)
		oldVx, err := strconv.ParseFloat(parts[3], 64)
		Check(err)
		oldVy, err := strconv.ParseFloat(parts[4], 64)
		Check(err)
		oldAx, err := strconv.ParseFloat(parts[5], 64)
		Check(err)
		oldAy, err := strconv.ParseFloat(parts[6], 64)
		Check(err)
		t, err := strconv.ParseFloat(parts[7], 64)
		Check(err)
		expPx, err := strconv.ParseFloat(parts[8], 64)
		Check(err)
		expPy, err := strconv.ParseFloat(parts[9], 64)
		Check(err)

		test := PositionTestCases{
			id: id,
			star: Star{
				position: OrderedPair{px, py},
			},
			oldAcceleration: OrderedPair{oldAx, oldAy},
			oldVelocity: OrderedPair{oldVx, oldVy},
			time: t,
			expected: OrderedPair{expPx, expPy},
		}
		tests = append(tests, test)
	}
	return tests
}




//// Test functions for eight subroutines in functions.go ////

// TestFindQuadrant tests the FindQuadrant function using data from a file.
// Input: t (*testing.T) - testing context.
// Output: None. Reports errors via t.Errorf if results do not match expected.
func TestFindQuadrant(t *testing.T) {
	stars, width, expected := ReadFindQuadrantData("Tests/FindQuadrant.txt")

	q := Quadrant{x: 0.0, y:0.0, width: width}

	for i, s := range stars {
		result := FindQuadrant(q, s)
		expectedResult := expected[i]

		if result != expectedResult {
			t.Errorf("TestFindQuadrant(test %v) = %v, want %v",
        		i, result, expectedResult)
		}
	}
}


// TestSubdivide tests the Subdivide function using data from a file.
// Input: t (*testing.T) - testing context.
// Output: None. Reports errors via t.Errorf if results do not match expected.
func TestSubdivide(t *testing.T) {
	tests := ReadSubdivideData("Tests/Subdivide.txt")

	for i, test := range tests {
		Subdivide(test.node)

		for j, child := range test.node.children {
			result := child.sector
			expectedResult := test.expected[j]

			if result != expectedResult {
				t.Errorf("TestSubdivide(test %v, children %v) = %v, want %v",
        			i, j, result, expectedResult)	
			}

		}
	}
}


// TestIsInsideUniverse tests the IsInsideUniverse function using data from a file.
// Input: t (*testing.T) - testing context.
// Output: None. Reports errors via t.Errorf if results do not match expected.
func TestIsInsideUniverse(t *testing.T) {
	tests := ReadIsInsideUniverse("Tests/IsInsideUniverse.txt")

	for i, test := range tests {
		result := IsInsideUniverse(&test.star, test.width)
		expectedResult := test.expected

		if result != expectedResult {
			t.Errorf("TestIsInsideUniverse(test %v) = %v, want %v",
				i, result, expectedResult)
		}
	}
}


// TestComputeCenterAndMass tests the ComputeCenterAndMass function using data from a file.
// Input: t (*testing.T) - testing context.
// Output: None. Reports errors via t.Errorf if results do not match expected.
func TestComputeCenterAndMass(t *testing.T) {
	tests := ReadComputeCenterAndMass("Tests/ComputeCenterAndMass.txt")

	for i, test := range tests {

		ComputeCenterAndMass(test.node)
		result := test.node.star

		if math.Abs(result.position.x - test.expectedX) > 1e-3 ||
			math.Abs(result.position.y - test.expectedY) > 1e-3 ||
			math.Abs(result.mass - test.expectedMass) > 1e-3 {
				t.Errorf("TestComputeCenterAndMass (test %v) = (x: %v, y: %v, mass: %v), want (x: %v, y: %v, mass: %v)",
					i, result.position.x, result.position.y, result.mass, test.expectedX, test.expectedY, test.expectedMass)
			}
	}
}


// TestIsLeaf tests the IsLeaf function using data from a file.
// Input: t (*testing.T) - testing context.
// Output: None. Reports errors via t.Errorf if results do not match expected.
func TestIsLeaf(t *testing.T) {
	tests := ReadIsLeaf("Tests/IsLeaf.txt")

	for _, test := range tests {
		node := &Node{children: test.children}
		result := IsLeaf(node)

		if result != test.expected {
			t.Errorf("TestIsLeaf (test %v) = %v, want %v",
				test.id, result, test.expected)
		}
	}
}


// TestDistance tests the Distance function using data from a file.
// Input: t (*testing.T) - testing context.
// Output: None. Reports errors via t.Errorf if results do not match expected.
func TestDistance(t *testing.T) {
	tests := ReadDistance("Tests/Distance.txt")

	for _, test := range tests {
		p1 := OrderedPair{x:test.x1, y:test.y1}
		p2 := OrderedPair{x:test.x2, y:test.y2}

		deltaX, deltaY, distance := Distance(p1, p2)

		if math.Abs(deltaX - test.expectedDeltaX) > 1e-3 ||
			math.Abs(deltaY - test.expectedDeltaY) > 1e-3 ||
			math.Abs(distance - test.expectedDistance) > 1e-3 {
				t.Errorf("TestDistance(test %v) = (deltaX: %v, deltaY: %v, distance: %v), want (x: %v, y:%v, distance: %v)",
					test.id, deltaX, deltaY, distance, test.expectedDeltaX, test.expectedDeltaY, test.expectedDistance)
			}
	}
}


// TestVelocity tests the UpdateVelocity function using data from a file.
// Input: t (*testing.T) - testing context.
// Output: None. Reports errors via t.Errorf if results do not match expected.
func TestVelocity(t *testing.T) {
	tests := ReadVelocity("Tests/UpdateVelocity.txt")

	for _, test := range tests {
		// need an address for the star!!!
		result := UpdateVelocity(&test.star, test.oldAcceleration, test.time)

		if math.Abs(result.x - test.expected.x) > 1e-3 ||
			math.Abs(result.y - test.expected.y) > 1e-3 {
				t.Errorf("TestVelocity(test %v) = (x: %v, y: %v), want (x: %v, y: %v)",
					test.id, result.x, result.y, test.expected.x, test.expected.y)
			}
	}
}


// TestPosition tests the UpdatePosition function using data from a file.
// Input: t (*testing.T) - testing context.
// Output: None. Reports errors via t.Errorf if results do not match expected.
func TestPosition(t *testing.T) {
	tests := ReadPosition("Tests/UpdatePosition.txt")

	for _, test := range tests {
		result := UpdatePosition(&test.star, test.oldAcceleration, test.oldVelocity, test.time)

		if math.Abs(result.x - test.expected.x) > 1e-3 ||
			math.Abs(result.y - test.expected.y) > 1e-3 {
				t.Errorf("TestPosition(test %v) = (x: %v, y: %v), want (x: %v, y: %v)",
					test.id, result.x, result.y, test.expected.x, test.expected.y)
			}
	}
}
