package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// const numItems = 20    // A reasonable value for exhaustive search.
// const numItems = 40    // A reasonable value for branch and bound.
const numItems = 80

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
// # items: 80
// Total value: 470
// Total weight: 534
// Allowed weight: 267
//
// Too many items for exhaustive search
//
// Too many items for branch and bound
//
// *** Rod's technique ***
// Elapsed: 23.307932
// 0(7, 8) 1(7, 4) 3(4, 4) 4(9, 5) 7(4, 5) 9(8, 5) 10(8, 6) 12(10, 4) 15(10, 7) 17(8, 9) 19(10, 7) 20(6, 7) 21(9, 4) 23(6, 4) 24(6, 5) 27(6, 5) 28(10, 4) 30(7, 8) 31(7, 6) 32(8, 7) 33(7, 8) 36(6, 4) 37(6, 5) 38(4, 4) 39(8, 9) 40(10, 5) 42(8, 9) 43(6, 6) 44(8, 8) 46(8, 4) 48(10, 6) 49(8, 9) 50(8, 5) 51(9, 6) 57(10, 6) 60(4, 4) 62(6, 6) 65(10, 9) 66(10, 6) 70(8, 5) 75(6, 5) 76(8, 7) 77(5, 5) 78(9, 7) 79(8, 5)
// Value: 340, Weight: 267, Calls: 1138854623
//
// *** Rod's technique Sorted ***
// Elapsed: 0.002952
// 0(10, 4) 1(10, 4) 2(9, 4) 3(8, 4) 4(10, 5) 5(9, 5) 6(8, 5) 7(8, 5) 8(8, 5) 9(8, 5) 10(7, 4) 11(10, 6) 12(10, 6) 13(10, 6) 14(6, 4) 15(6, 4) 16(9, 6) 17(8, 6) 18(10, 7) 19(10, 7) 20(6, 5) 21(9, 7) 22(6, 5) 23(6, 5) 24(6, 5) 25(8, 7) 26(8, 7) 27(7, 6) 28(4, 4) 29(4, 4) 30(6, 6) 31(6, 6) 32(4, 4) 33(5, 5) 34(8, 8) 35(6, 7) 36(10, 9) 37(4, 5) 38(8, 9) 39(7, 8) 40(8, 9) 41(8, 9) 42(8, 9) 43(7, 8) 44(7, 8)
// Value: 340, Weight: 267, Calls: 209529

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

	// Rod's technique
	if numItems > 85 { // Only use Rod's technique if numItems <= 85.
		fmt.Println("Too many items for Rod's technique\n")
	} else {
		fmt.Println("*** Rod's technique ***")
		runAlgorithm(rodsTechnique, items, allowedWeight)
	}

	// Rod's technique sorted
	if numItems > 350 { // Only use Rod's technique if numItems <= 350.
		fmt.Println("Too many items for Rod's technique\n")
	} else {
		fmt.Println("*** Rod's technique Sorted ***")
		runAlgorithm(rodsTechniqueSorted, items, allowedWeight)
	}
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

func rodsTechnique(items []Item, allowedWeight int) ([]Item, int, int) {
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
