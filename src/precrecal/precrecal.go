//package main
package precrecal

import (
	"fmt"
	"sortmap"
	"strconv"
	"strings"
	"io/ioutil"
)

type PRGraph struct {
	testSample map[string]map[string]bool
}

func NewPRGraph() *PRGraph {
	prGraph := &PRGraph{}
	prGraph.init("data/RelevancyLists.txt")
	return prGraph
}

func NewPRGraph___forTesting() *PRGraph {
	prGraph := &PRGraph{}
	prGraph.testSample = readTestSampleToMap___forTesting()
	return prGraph
}

func (prGraph *PRGraph) init(fileName string) {
	prGraph.testSample = readTestSampleToMap(fileName)
	//fmt.Println(len(prGraph.testSample))
}

func (prGraph *PRGraph) MakeAvgInterpolatedPRTable(retrievedLists map[string]sortmap.PairList) {
	numBins := 11
	avgIntrplPRTable := make([]float64, numBins)
	sampleSize := len(prGraph.testSample)
	
	for sample := range prGraph.testSample {
		fmt.Printf("\nQueryNumber: %s\n", sample)
		precision, recall, size :=makePrecRecallTable(retrievedLists[sample], prGraph.testSample[sample])
		intrplPRTable := makeInterpolatedPRTable(precision, recall, size, numBins)
		for i := 0 ; i < numBins ; i++ {
			avgIntrplPRTable[i] += intrplPRTable[i]/float64(sampleSize)
		}
		fmt.Println("---------------------------------")
	}
	fmt.Println("\n\n\nFinal Avg. Interpolated PR Table:\n")
	for j:=0 ; j < numBins ; j++ {
		fmt.Printf("%.2f  %.2f\n", float64(j)/10, avgIntrplPRTable[j])
	}
}

func makePrecRecallTable(retrievedList sortmap.PairList, testData map[string]bool) ([]float64, []float64, int) {
	totalExpected := float64(len(testData))
	precision := make([]float64, len(retrievedList))
	recall := make([]float64, len(retrievedList))
	correctCount := 0.0
	/*
	fmt.Printf("retrieved: %v, expected: %v\n", len(retrievedList), totalExpected)
	fmt.Println("Test Data:")
	for n := range testData {
		fmt.Printf("%v\n", n)
	}
	fmt.Printf("----\n")
	 */
	var num int
	var pair sortmap.Pair
	fmt.Printf("DocNum\tRecall\tPrecision\n")
	for num, pair = range retrievedList {

		// No need to continue because all correct results were already	displayed.
		if len(testData) == 0 {	break }

		// Calculate Precision-Recall Table:
		if testData[pair.Key] {
			correctCount++
			// Delete to keep track of when no more correct resluts are left.
			delete(testData, pair.Key)
		}
		recall[num] = correctCount/totalExpected
		precision[num] = correctCount/float64(num+1) // num retrieved sofar.
		fmt.Printf("%s\t%.2f\t%.2f\n", pair.Key, recall[num], precision[num])
	}
	return precision, recall, num
}

func makeInterpolatedPRTable(precision, recall []float64, numRetrievals int, numBins int) []float64 {

	interpolatedPR := make([]float64, numBins)
	periods := make([]float64, numBins)
	currentRecall := 0.0
	maxPrecision := 0.0

	fmt.Printf("num retrievals: %d\n", numRetrievals)
	for num := 0 ; num < numRetrievals ; num++ { 
		// Find max precision for this period.
		if maxPrecision < precision[num] {
			maxPrecision = precision[num]
		}
		// Do not split on zero.
		if currentRecall  == 0.0 {
			currentRecall =recall[num]
		}
		// Split a new period.
		if currentRecall != recall[num] {
			// Save this period.
			periods[int(currentRecall*10)] = maxPrecision
			// Restart counting for the next period.
			maxPrecision = precision[num]
			currentRecall = recall[num]

			// Don't forget the last period.
				if num == (numRetrievals-1) {
				periods[int(currentRecall*10)] = maxPrecision
			}
		}
	}
	max := 0.0
	for i := (numBins - 1) ; i >= 0 ; i-- {
		if periods[i] != 0.0 {
			max = periods[i]
		}
		interpolatedPR[i] = max
	}
	fmt.Println("\nInterpolated Table:")
	for j:=0; j < numBins ; j++ {
		fmt.Printf("%.1f  %.2f\n", float64(j)/10, interpolatedPR[j])
	}
	return interpolatedPR
}


func readTestSampleToMap(file string) map[string]map[string]bool {
	testSample := make(map[string]map[string]bool)
	data, _ := ioutil.ReadFile(file)
	lines := strings.Split(string(data), "\r\n")
	
	for _, line := range lines {
		tokens := strings.Split(line, " ")

		if tokens != nil && tokens[0] != "" {
			//To handle the null-charachter case
			v, _, _, _ := strconv.UnquoteChar(tokens[0], 0)
			if v != 0 {
				maplist := make(map[string]bool, len(tokens) -1)

				for i := 1 ; i < len(tokens) ; i++ {
					if tokens[i] != "" { 
						maplist[tokens[i]] = true
					}
				}
				testSample[tokens[0]] = maplist 
			}
		}
	}
/*
	for key := range testSample {
		fmt.Printf("%v: ", key)
		fmt.Println(testSample[key])
	}
	fmt.Println(len(testSample))
*/
	return testSample
}

func readTestSampleToMap___forTesting() map[string]map[string]bool {
	testData := make(map[string]bool, 2)
	testData["1"] = true
	testData["3"] = true
	testData2 := make(map[string]bool, 3)
	testData2["2"] = true
	testData2["3"] = true
	testData2["4"] = true
	
	testDataSet := map[string]map[string]bool {"100":testData, "200":testData2}
	return testDataSet
}


func main() {
	/* Note: for Testing purposes: Change package name to main.*/

	// Test 1
	list := make(sortmap.PairList, 10)
	list[0] = sortmap.Pair{"2", 11.53}  
	list[1] = sortmap.Pair{"3", 9.30}  
	list[2] = sortmap.Pair{"8", 9.26}  
	list[3] = sortmap.Pair{"7", 5.26}  
	list[4] = sortmap.Pair{"9", 2.26}  
	list[5] = sortmap.Pair{"1", 1.26}  
	list[6] = sortmap.Pair{"20", 1.26}  
	list[7] = sortmap.Pair{"19", 1.26}  

	testData := make(map[string]bool, 2)
	testData["1"] = true
	testData["3"] = true


	// Test 2
	list2 := make(sortmap.PairList, 10)
	list2[0] = sortmap.Pair{"3", 11.53}  
	list2[1] = sortmap.Pair{"2", 9.30}  
	list2[2] = sortmap.Pair{"5", 9.26}  
	list2[3] = sortmap.Pair{"7", 5.26}  
	list2[4] = sortmap.Pair{"4", 2.26}  
	list2[5] = sortmap.Pair{"90", 1.26}  
	list2[6] = sortmap.Pair{"20", 1.26}  
	list2[7] = sortmap.Pair{"19", 1.26}  

	testData2 := make(map[string]bool, 3)
	testData2["2"] = true
	testData2["3"] = true
	testData2["4"] = true

	//numBins := 11
	//precision, recall, size := makePrecRecallTable(list, testData)
	//makeInterpolatedPRTable(precision, recall, size, numBins)

	prg := NewPRGraph___forTesting()
	retrievedLists := map[string]sortmap.PairList {"100":list, "200":list2}
	prg.MakeAvgInterpolatedPRTable(retrievedLists)
}

