package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const numItems = 500

const minValue = 1
const maxValue = 10
const minWeight = 4
const maxWeight = 10

var allowedWeight int

type Item struct {
	id, blockedBy int
	blockList     []int // Other items that this one blocks.
	value, weight int
	isSelected    bool
}

// Test results:
// *** Parameters ***
// # items: 200
// Total value: 1177
// Total weight: 1393
// Allowed weight: 696
//
// *** Rod's technique Sorted ***
// Elapsed: 1.059725
// Value: 847, Weight: 696, Calls: 40878809
//
// *** Dynamic programming ***
// Elapsed: 0.001068
// Value: 847, Weight: 696, Calls: 1
//
//
// *** Parameters ***
// # items: 500
// Total value: 2733
// Total weight: 3552
// Allowed weight: 1776
//
// Too many items for Rod's technique
//
// *** Dynamic programming ***
// Elapsed: 0.007502
// Value: 2059, Weight: 1776, Calls: 1

func main() {
	items := makeItems(numItems, minValue, maxValue, minWeight, maxWeight)
	allowedWeight = sumWeights(items, true) / 2

	// Display basic parameters.
	fmt.Println("*** Parameters ***")
	fmt.Printf("# items: %d\n", numItems)
	fmt.Printf("Total value: %d\n", sumValues(items, true))
	fmt.Printf("Total weight: %d\n", sumWeights(items, true))
	fmt.Printf("Allowed weight: %d\n", allowedWeight)
	fmt.Println()

	// Rod's technique sorted
	if numItems > 350 { // Only use Rod's technique if numItems <= 350.
		fmt.Println("Too many items for Rod's technique\n")
	} else {
		fmt.Println("*** Rod's technique Sorted ***")
		runAlgorithm(rodsTechniqueSorted, items, allowedWeight)
	}

	// Dynamic programming
	fmt.Println("*** Dynamic programming ***")
	runAlgorithm(dynamicProgramming, items, allowedWeight)
}

// Make some random items.
func makeItems(numItems, minValue, maxValue, minWeight, maxWeight int) []Item {
	// Initialize a pseudorandom number generator.
	random := rand.New(rand.NewSource(time.Now().UnixNano())) // Initialize with a changing seed
	//random := rand.New(rand.NewSource(1337)) // Initialize with a fixed seed

	items := make([]Item, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = Item{
			i, -1, nil,
			random.Intn(maxValue-minValue+1) + minValue,
			random.Intn(maxWeight-minWeight+1) + minWeight,
			false}
	}
	return items
}

// Build the items' block lists.
func makeBlockLists(items []Item) {
	for i := range items {
		items[i].blockList = []int{}
		for j := range items {
			if i != j {
				if items[i].value >= items[j].value && items[i].weight <= items[j].weight {
					items[i].blockList = append(items[i].blockList, items[j].id)
				}
			}
		}
	}
}

// Block items on this item's blocks list.
func blockItems(source Item, items []Item) {
	for _, otherId := range source.blockList {
		if items[otherId].blockedBy < 0 {
			items[otherId].blockedBy = source.id
		}
	}
}

// Unblock items on this item's blocks list.
func unblockItems(source Item, items []Item) {
	for _, otherId := range source.blockList {
		if items[otherId].blockedBy == source.id {
			items[otherId].blockedBy = -1
		}
	}
}

// Return a copy of the items slice.
func copyItems(items []Item) []Item {
	newItems := make([]Item, len(items))
	copy(newItems, items)
	return newItems
}

// Return the total value of the items.
// If addAll is false, only add up the selected items.
func sumValues(items []Item, addAll bool) int {
	total := 0
	for i := 0; i < len(items); i++ {
		if addAll || items[i].isSelected {
			total += items[i].value
		}
	}
	return total
}

// Run the algorithm. Display the elapsed time and solution.
func runAlgorithm(alg func([]Item, int) ([]Item, int, int), items []Item, allowedWeight int) {
	// Copy the items so the run isn't influenced by a previous run.
	testItems := copyItems(items)

	start := time.Now()

	// Run the algorithm.
	solution, totalValue, functionCalls := alg(testItems, allowedWeight)

	elapsed := time.Since(start)

	fmt.Printf("Elapsed: %f\n", elapsed.Seconds())
	printSelected(solution)
	fmt.Printf("Value: %d, Weight: %d, Calls: %d\n",
		totalValue, sumWeights(solution, false), functionCalls)
	fmt.Println()
}

// Print the selected items.
func printSelected(items []Item) {
	numPrinted := 0
	for i, item := range items {
		if item.isSelected {
			fmt.Printf("%d(%d, %d) ", i, item.value, item.weight)
		}
		numPrinted += 1
		if numPrinted > 100 {
			fmt.Println("...")
			return
		}
	}
	fmt.Println()
}

// Return the value of this solution.
// If the solution is too heavy, return -1 so we prefer an empty solution.
func solutionValue(items []Item, allowedWeight int) int {
	// If the solution's total weight > allowedWeight,
	// return 0 so we won't use this solution.
	if sumWeights(items, false) > allowedWeight {
		return -1
	}

	// Return the sum of the selected values.
	return sumValues(items, false)
}

// Return the total weight of the items.
// If addAll is false, only add up the selected items.
func sumWeights(items []Item, addAll bool) int {
	total := 0
	for i := 0; i < len(items); i++ {
		if addAll || items[i].isSelected {
			total += items[i].weight
		}
	}
	return total
}

func rodsTechniqueSorted(items []Item, allowedWeight int) ([]Item, int, int) {
	makeBlockLists(items)

	// Sort so items with longer blocked lists come first.
	sort.Slice(items, func(i, j int) bool {
		return len(items[i].blockList) > len(items[j].blockList)
	})

	// Reset the items' IDs.
	for i := range items {
		items[i].id = i
	}

	// Rebuild the blocked lists with the new indices.
	makeBlockLists(items)

	bestValue := 0
	currentValue := 0
	currentWeight := 0
	remainingValue := sumValues(items, true)

	return doRodsTechnique(items, allowedWeight, 0,
		bestValue, currentValue, currentWeight, remainingValue)
}

func doRodsTechnique(items []Item, allowedWeight, nextIndex,
	bestValue, currentValue, currentWeight, remainingValue int,
) ([]Item, int, int) {
	// See if we have a full assignment.
	if nextIndex >= len(items) {
		copiedItems := copyItems(items)
		solutionVal := solutionValue(copiedItems, allowedWeight)
		if solutionVal > bestValue {
			bestValue = solutionVal
		}
		return copiedItems, solutionVal, 1
	}

	// We do not have a full assignment.
	// See if we can improve this solution enough to be worth pursuing.
	if currentValue+remainingValue < bestValue {
		// We cannot improve on the best solution found so far.
		return nil, 0, 1
	}

	// Try adding the next item.
	var test1Solution []Item
	test1Solution = nil
	test1Value := 0
	test1Calls := 1
	if currentWeight+items[nextIndex].weight <= allowedWeight && items[nextIndex].blockedBy < 0 {
		items[nextIndex].isSelected = true
		test1Solution, test1Value, test1Calls = doRodsTechnique(items, allowedWeight, nextIndex+1,
			bestValue, currentValue+items[nextIndex].value, currentWeight+items[nextIndex].weight, remainingValue-items[nextIndex].value)
		if test1Value > bestValue {
			bestValue = test1Value
		}
	}

	// Try not adding the next item.
	blockItems(items[nextIndex], items)
	items[nextIndex].isSelected = false
	test2Solution, test2Value, test2Calls := doRodsTechnique(items, allowedWeight, nextIndex+1,
		bestValue, currentValue, currentWeight, remainingValue-items[nextIndex].value)
	unblockItems(items[nextIndex], items)
	if test2Value > bestValue {
		bestValue = test2Value
	}

	// Return the solution that is better.
	if test1Value >= test2Value {
		return test1Solution, test1Value, test1Calls + test2Calls + 1
	} else {
		return test2Solution, test2Value, test1Calls + test2Calls + 1
	}
}

// Use dynamic programming to find a solution.
// Return the best assignment, value of that assignment,
// and the number of function calls we made.
func dynamicProgramming(items []Item, allowedWeight int) ([]Item, int, int) {
	numItems := len(items)

	// Allocate the arrays.
	solutionValue := make([][]int, numItems)
	prevWeight := make([][]int, numItems)
	for i := 0; i < numItems; i++ {
		solutionValue[i] = make([]int, allowedWeight+1)
		prevWeight[i] = make([]int, allowedWeight+1)
	}

	// Initialize the row item 0.
	for w := 0; w <= allowedWeight; w++ {
		if items[0].weight <= w {
			// items[0] fits.
			solutionValue[0][w] = items[0].value
			prevWeight[0][w] = -1
		} else {
			// items[0] does not fit.
			solutionValue[0][w] = 0
			prevWeight[0][w] = w
		}
	}

	// Fill in the remaining table rows.
	for i := 1; i < numItems; i++ {
		for w := 0; w <= allowedWeight; w++ {
			// Calculate the value if we do not use the new item i.
			valueWithoutI := solutionValue[i-1][w]

			// Calculate the value if we do use the new item i.
			valueWithI := 0
			if items[i].weight <= w { // Make sure it fits.
				valueWithI = solutionValue[i-1][w-items[i].weight] + items[i].value
			}

			// See which is better.
			if valueWithoutI >= valueWithI {
				// We're better off omitting item i.
				solutionValue[i][w] = valueWithoutI
				prevWeight[i][w] = w
			} else {
				// We're better off including item i.
				solutionValue[i][w] = valueWithI
				prevWeight[i][w] = w - items[i].weight
			}
		}
	}

	// Reconstruct the solution.
	// Get the row and column for the final solution.
	i := numItems - 1
	w := allowedWeight

	// Work backwards until we reach an initial solution.
	for i >= 0 {
		// Check prevWeight for the current solution.
		prevW := prevWeight[i][w]
		if w == prevW {
			// We skipped item i.
			// Leave w unchanged.
		} else {
			// We added item i.
			items[i].isSelected = true // Select this item in the solution.
			w = prevW                  // Move to the previous solution's weight.
		}
		i -= 1 // Move to the previous row.
	}

	return items, solutionValue[numItems-1][allowedWeight], 1
}
