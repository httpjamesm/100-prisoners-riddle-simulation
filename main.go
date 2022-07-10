package main

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"time"

	"math/rand"

	"github.com/gookit/color"
	"github.com/schollz/progressbar/v3"
)

func main() {

	runs := flag.Int("runs", 1000, "number of simulation runs to perform")

	flag.Parse()

	if *runs < 1 {
		fmt.Println("runs must be greater than 0")
		return
	}

	if *runs < 100 {
		color.Yellow.Println("[WARN] Less than 100 runs may return inaccurate results.\n")
	}

	run(*runs)

}

func run(runs int) {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])

	var newSeed int64

	if err != nil {
		color.Yellow.Println("[WARN] Unable to generate random seed with cryptographically secure source. Falling back to time seed.\n")
		newSeed = time.Now().UnixNano()
	} else {
		newSeed = int64(binary.LittleEndian.Uint64(b[:]))
	}

	rand.Seed(newSeed)

	max := runs

	attempts := 0
	successes := 0
	fails := 0

	bar := progressbar.Default(int64(max))

	startTime := time.Now().UnixNano()

	for i := 0; i < max; i++ {
		attempts++
		if runSimulation() {
			successes++
		} else {
			fails++
		}

		bar.Add(1)
	}

	endTime := time.Now().UnixNano()

	fmt.Println("Attempts:", attempts)
	fmt.Println("Successes:", successes)
	fmt.Println("Fails:", fails)
	fmt.Println("Success Rate:", float64(successes)/float64(attempts)*100, "%")
	fmt.Println("Time:", (endTime-startTime)/1000000, "ms")
}

// SIMULATION INSTRUCTIONS:
// there are 100 prisoners and 100 boxes
// each box has a random prisoner number
// each prisoner goes into the room and looks into the box labeled with their number
// if the prisoner's number is not within that box, they go to the box with the number within that box
// they continue until they find a box with their number, or until they exceed 50 attempts
// if they exceed 50 attempts, the entire prisoner group is unsuccessful
// if all prisoners succeed in finding a box with their number, the simulation is successful

func runSimulation() (success bool) {
	room := map[int]int{}

	boxes := []int{}

	for i := 1; i < 101; i++ {
		boxes = append(boxes, i)
	}

	// shuffle boxes
	for i := len(boxes) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		boxes[i], boxes[j] = boxes[j], boxes[i]
	}

	// assign prisoner numbers to boxes
	for i := 1; i < 101; i++ {
		room[i] = boxes[i-1]
	}

	successfulPrisoners := []int{}

	for prisonerNumber := 1; prisonerNumber < 101; prisonerNumber++ {
		// first prisoner goes to the box with their number
		initialBoxContents := room[prisonerNumber]
		if initialBoxContents == prisonerNumber {
			// if that box has the prisoner's number, they're successful
			successfulPrisoners = append(successfulPrisoners, prisonerNumber)
			continue
		}

		// otherwise follow the numbers
		foundBox := false

		attempts := 0

		currentBoxContents := initialBoxContents

		trail := []int{
			prisonerNumber,
		}

		for !foundBox {
			if attempts == 49 {
				// the first box counted as 1 attempt, and 50 is the max, so stop after 49 (49 + 1 = 50)
				return false
			}

			attempts++
			trail = append(trail, currentBoxContents)

			nextBoxContents := room[currentBoxContents]
			if nextBoxContents == prisonerNumber {
				// if that box has the prisoner's number, they're successful
				successfulPrisoners = append(successfulPrisoners, prisonerNumber)
				foundBox = true
				break
			}

			currentBoxContents = nextBoxContents
		}
	}

	return len(successfulPrisoners) == 100
}
