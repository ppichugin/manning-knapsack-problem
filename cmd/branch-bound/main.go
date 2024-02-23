package main

import (
	"fmt"
	"math/rand"
	"time"
)

// const numItems = 20    // A reasonable value for exhaustive search.
// const numItems = 40    // A reasonable value for branch and bound.
const numItems = 25

const minValue = 1
const maxValue = 10
const minWeight = 4
const maxWeight = 10

var allowedWeight int

type Item struct {
	value, weight int
	isSelected    bool
}

// TEST RESULTs:
// *** Parameters ***
// # items: 25
// Total value: 142
// Total weight: 158
// Allowed weight: 79
//
// *** Exhaustive Search ***
// Elapsed: 4.897523
// 0(9, 5) 1(10, 4) 2(7, 5) 3(5, 6) 6(4, 6) 7(9, 6) 9(10, 7) 15(6, 5) 17(8, 8) 18(9, 7) 19(7, 5) 20(9, 4) 21(4, 5) 22(6, 6)
// Value: 103, Weight: 79, Calls: 67108863
//
// *** Branch and Bound ***
// Elapsed: 0.001524
// 0(9, 5) 1(10, 4) 2(7, 5) 3(5, 6) 4(4, 5) 6(4, 6) 7(9, 6) 9(10, 7) 15(6, 5) 17(8, 8) 18(9, 7) 19(7, 5) 20(9, 4) 22(6, 6)
// Value: 103, Weight: 79, Calls: 589017

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

	// Exhaustive search
	if numItems > 25 { // Only run exhaustive search if numItems <= 25.
		fmt.Println("Too many items for exhaustive search\n")
	} else {
		fmt.Println("*** Exhaustive Search ***")
		runAlgorithm(exhaustiveSearch, items, allowedWeight)
	}

	// Branch and bound
	if numItems > 45 { // Only run branch and bound if numItems <= 45.
		fmt.Println("Too many items for branch and bound\n")
	} else {
		fmt.Println("*** Branch and Bound ***")
		runAlgorithm(branchAndBound, items, allowedWeight)
	}
}

// Make some random items.
func makeItems(numItems, minValue, maxValue, minWeight, maxWeight int) []Item {
	// Initialize a pseudorandom number generator.
	//random := rand.New(rand.NewSource(time.Now().UnixNano())) // Initialize with a changing seed
	random := rand.New(rand.NewSource(1337)) // Initialize with a fixed seed

	items := make([]Item, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = Item{
			random.Intn(maxValue-minValue+1) + minValue,
			random.Intn(maxWeight-minWeight+1) + minWeight,
			false}
	}
	return items
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

// Recursively assign values in or out of the solution.
// Return the best assignment, value of that assignment,
// and the number of function calls we made.
func exhaustiveSearch(items []Item, allowedWeight int) ([]Item, int, int) {
	return doExhaustiveSearch(items, allowedWeight, 0)
}

func doExhaustiveSearch(items []Item, allowedWeight, nextIndex int) ([]Item, int, int) {
	if nextIndex >= len(items) {
		copiedItems := copyItems(items)
		solutionVal := solutionValue(copiedItems, allowedWeight)
		return copiedItems, solutionVal, 1
	}

	items[nextIndex].isSelected = true
	withItem, withValue, withCalls := doExhaustiveSearch(items, allowedWeight, nextIndex+1)

	items[nextIndex].isSelected = false
	withoutItem, withoutValue, withoutCalls := doExhaustiveSearch(items, allowedWeight, nextIndex+1)

	if withValue > withoutValue {
		return withItem, withValue, withCalls + withoutCalls + 1
	}
	return withoutItem, withoutValue, withCalls + withoutCalls + 1
}

// Use branch and bound to find a solution.
// Return the best assignment, value of that assignment,
// and the number of function calls we made.
func branchAndBound(items []Item, allowedWeight int) ([]Item, int, int) {
	bestValue := 0
	currentValue := 0
	currentWeight := 0
	remainingValue := sumValues(items, true)

	return doBranchAndBound(items, allowedWeight, 0,
		bestValue, currentValue, currentWeight, remainingValue)
}

func doBranchAndBound(items []Item, allowedWeight, nextIndex,
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
	var test1Value int
	var test1Calls int
	if currentWeight+items[nextIndex].weight <= allowedWeight {
		items[nextIndex].isSelected = true
		test1Solution, test1Value, test1Calls = doBranchAndBound(items, allowedWeight, nextIndex+1,
			bestValue, currentValue+items[nextIndex].value, currentWeight+items[nextIndex].weight, remainingValue-items[nextIndex].value)
		if test1Value > bestValue {
			bestValue = test1Value
		}
	} else {
		test1Solution = nil
		test1Value = 0
		test1Calls = 1
	}

	// Try not adding the next item.
	var test2Solution []Item
	var test2Value int
	var test2Calls int
	// See if there is a chance of improvement without this item's value.
	if currentValue+remainingValue-items[nextIndex].value > bestValue {
		items[nextIndex].isSelected = false
		test2Solution, test2Value, test2Calls = doBranchAndBound(items, allowedWeight, nextIndex+1,
			bestValue, currentValue, currentWeight, remainingValue-items[nextIndex].value)
		if test2Value > bestValue {
			bestValue = test2Value
		}
	} else {
		test2Solution = nil
		test2Value = 0
		test2Calls = 1
	}

	// Return the solution that is better.
	if test1Value >= test2Value {
		return test1Solution, test1Value, test1Calls + test2Calls + 1
	} else {
		return test2Solution, test2Value, test1Calls + test2Calls + 1
	}
}
