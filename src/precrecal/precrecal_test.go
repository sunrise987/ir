package precrecal_test

import (
	"math"
	"testing"
	//"log"
	//"fmt"
	. "precrecal"
	. "sortmap"
)

func NewPRGraph___forTesting() *PRGraph {
	prGraph := &PRGraph{}
	prGraph.TestQueries = readTestSampleToMap___forTesting()
	return prGraph
}

func readTestSampleToMap___forTesting() map[string]SearchResults {
	testData := make(SearchResults, 2)
	testData["1"] = true
	testData["3"] = true
	testData2 := make(SearchResults, 3)
	testData2["2"] = true
	testData2["3"] = true
	testData2["4"] = true
	testDataSet := map[string]SearchResults{"100": testData, "200": testData2}
	return testDataSet
}

func Test1(t *testing.T) {
	// Test 1
	list := make(PairList, 10)
	list[0] = Pair{"2", 11.53}
	list[1] = Pair{"3", 9.30}
	list[2] = Pair{"8", 9.26}
	list[3] = Pair{"7", 5.26}
	list[4] = Pair{"9", 2.26}
	list[5] = Pair{"1", 1.26}
	list[6] = Pair{"20", 1.26}
	list[7] = Pair{"19", 1.26}

	testData := make(SearchResults, 2)
	testData["1"] = true
	testData["3"] = true

	// Test 2
	list2 := make(PairList, 10)
	list2[0] = Pair{"3", 11.53}
	list2[1] = Pair{"2", 9.30}
	list2[2] = Pair{"5", 9.26}
	list2[3] = Pair{"7", 5.26}
	list2[4] = Pair{"4", 2.26}
	list2[5] = Pair{"90", 1.26}
	list2[6] = Pair{"20", 1.26}
	list2[7] = Pair{"19", 1.26}

	testData2 := make(SearchResults, 3)
	testData2["2"] = true
	testData2["3"] = true
	testData2["4"] = true

	prg := NewPRGraph___forTesting()
	retrievedLists := map[string]PairList{"100": list, "200": list2}
	avgInt := prg.MakeAvgInterpolatedPRTable(retrievedLists)


	// Check the exact values of the Average Interpolated Precision-Recall Table.

	for i := 0 ; i < 6 ; i++ {
		// Round avgInt table Values to two decimal values.
		val := int(math.Floor((avgInt[i] * 100) + 0.5))

		if val != 75 {
			t.Error("Wrong Calculation. Expected value 0.75. Calculated value %.2f", avgInt[i])
		}
	}
	// Round avgInt table Values to two decimal values.
	val := int(math.Floor((avgInt[6] * 100) + 0.5))
	if val != 67 {
		t.Error("Wrong Calculation. Expected value 0.67. Calculated value %.2f", avgInt[6])
	}

	for i := 7 ; i < 11 ; i++ {
		// Round avgInt table Values to two decimal values.
		val := int(math.Floor((avgInt[i] * 100) + 0.5))

		if val != 47 {
			t.Error("Wrong Calculation. Expected value 0.47. Calculated value %.2f", avgInt[i])
		}
	}
}
