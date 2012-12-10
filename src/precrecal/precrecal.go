package precrecal

import (
	"fmt"
	"io/ioutil"
	"sortmap"
	"strconv"
	"strings"
)

var numBins int = 11

type SearchResults map[string]bool
type PRGraph struct {
	TestQueries map[string]SearchResults
}

func NewPRGraph() *PRGraph {
	prGraph := &PRGraph{}
	prGraph.init("data/RelevancyLists.txt")
	return prGraph
}


func (prGraph *PRGraph) init(fileName string) {
	prGraph.TestQueries = readTestSampleToMap(fileName)
}

func (prGraph *PRGraph) MakeAvgInterpolatedPRTable(retrievedLists map[string]sortmap.PairList) []float64 {
	avgIntrplPRTable := make([]float64, numBins)
	sampleSize := len(prGraph.TestQueries)
	
	for sample := range prGraph.TestQueries {
		fmt.Printf("\nQueryNumber: %s\n", sample)
		precision, recall, size := makePrecRecallTable(retrievedLists[sample], prGraph.TestQueries[sample])
		if size == 0 { continue }
		intrplPRTable := makeInterpolatedPRTable(precision, recall, size)
		for i := 0; i < numBins; i++ {
			avgIntrplPRTable[i] += intrplPRTable[i] / float64(sampleSize)
		}
	}
	fmt.Println("\n\nAverage Interpolated Precision Recall Table.")
	fmt.Println("Recall\tPrecision")
	for j := 0; j < numBins; j++ {
		fmt.Printf("%.2f  %.2f\n", float64(j)/10, avgIntrplPRTable[j])
	}
	return avgIntrplPRTable
}

func (prGraph *PRGraph) MakeOneInterpolatedPRTable(retrievedList sortmap.PairList, queryNum string) {
	fmt.Printf("\nQueryNumber: %s\n", queryNum)
	precision, recall, size := makePrecRecallTable(retrievedList, prGraph.TestQueries[queryNum])
	makeInterpolatedPRTable(precision, recall, size)
}

func makePrecRecallTable(retrievedList sortmap.PairList, testData map[string]bool) ([]float64, []float64, int) {
	if len(retrievedList) == 0 {
		return nil, nil, 0
	}
	totalExpected := float64(len(testData))
	precision := make([]float64, len(retrievedList))
	recall := make([]float64, len(retrievedList))
	correctCount := 0.0
	var num int
	var pair sortmap.Pair
	fmt.Printf("DocNum\tRecall\tPrecision\n")
	for num, pair = range retrievedList {

		// No need to continue because all correct results were already	displayed.
		if len(testData) == 0 {
			num--
			break
		}

		// Calculate Precision-Recall Table:
		if testData[pair.Key] {
			correctCount++
			// Delete to keep track of when no more correct resluts are left.
			delete(testData, pair.Key)
		}
		recall[num] = correctCount / totalExpected
		precision[num] = correctCount / float64(num+1) // num retrieved sofar.
		fmt.Printf("%s\t%.2f\t%.2f\n", pair.Key, recall[num], precision[num])
	}
	size := num + 1
	return precision, recall, size
}

func makeInterpolatedPRTable(precision, recall []float64, numRetrievals int) []float64 {
	interpolatedPR := make([]float64, numBins)
	if numRetrievals == 0 { return nil }
	if numRetrievals != len(precision) && numRetrievals != len(recall) {
		fmt.Printf("Error: numRetrievals: %v, precisionlen: %v, recallen: %v\n", 
			numRetrievals, len(precision), len(recall))
	}
	periods := make([]float64, numBins)
	currentRecall := 0.0
	maxPrecision := 0.0
	fmt.Printf("num retrievals: %d\n", numRetrievals)
	for num := 0; num < numRetrievals; num++ {
		// Find max precision for this period.
		if maxPrecision < precision[num] {
			maxPrecision = precision[num]
		}
		// Do not split on zero.
		if currentRecall == 0.0 {
			currentRecall = recall[num]
		}
		// Split a new period.
		if currentRecall != recall[num] {
			// Save this period.
			periods[int(currentRecall*10)] = maxPrecision
			// Restart counting for the next period.
			maxPrecision = precision[num]
			currentRecall = recall[num]

			// Don't forget the last period.
			if num == (numRetrievals - 1) {
				periods[int(currentRecall*10)] = maxPrecision
			}
		}
	}
	max := 0.0
	for i := (numBins - 1); i >= 0; i-- {
		if periods[i] != 0.0 {
			max = periods[i]
		}
		interpolatedPR[i] = max
	}
	fmt.Println("Interpolated Table:")
	for j := 0; j < numBins; j++ {
		fmt.Printf("%.1f  %.2f\n", float64(j)/10, interpolatedPR[j])
	}
	return interpolatedPR
}

func readTestSampleToMap(file string) map[string]SearchResults {
	testQueries := make(map[string]SearchResults)
	data, _ := ioutil.ReadFile(file)
	lines := strings.Split(string(data), "\r\n")

	for _, line := range lines {
		tokens := strings.Split(line, " ")

		if tokens != nil && tokens[0] != "" {
			//To handle the null-charachter case
			v, _, _, _ := strconv.UnquoteChar(tokens[0], 0)
			if v != 0 {
				maplist := make(map[string]bool, len(tokens)-1)

				for i := 1; i < len(tokens); i++ {
					if tokens[i] != "" {
						maplist[tokens[i]] = true
					}
				}
				testQueries[tokens[0]] = maplist
			}
		}
	}
	return testQueries
}

func (prgraph *PRGraph) Print() {
	for key := range prgraph.TestQueries {
		fmt.Printf("%v: ", key)
		fmt.Println(prgraph.TestQueries[key])
	}
	fmt.Println(len(prgraph.TestQueries))
}
