package main
//package precrecal

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

func (prGraph *PRGraph) init(fileName string) {
	prGraph.testSample = readTestSampleToMap(fileName)
	//fmt.Println(len(prGraph.testSample))
}

func (prGraph *PRGraph) MakePrecRecalTable(list sortmap.PairList, queryNum string) {
	retrieved := float64(len(list))
	size := int(len(prGraph.testSample[queryNum]))
	expected := float64(size)
	fmt.Printf("retrieved: %v, expected: %v\n", retrieved, expected)

	for n := range prGraph.testSample[queryNum] {
		fmt.Printf("%v\n", n)
	}
	fmt.Printf("----\n")
	recall := make([]float64, len(list))
	precision := make([]float64, len(list))
	interpolatedPR := make([]float64, 11)
	currentRecall := 0.0
	maxPrecision := 0.0
	correctCount := 0.0
	j:=0
	for i, pair := range list {
		fmt.Println(pair.Key)
		if prGraph.testSample[queryNum][pair.Key] == true {
			correctCount++
		}
		recall[i] = correctCount/expected
		precision[i] = correctCount/retrieved
		if precision[i] > maxPrecision {
			maxPrecision = precision[i]
		}
		if recall[i] != currentRecall { // start a new period 
			// Save this Period:
			for ; j <= int(currentRecall*10) && j < 11 ; j++ {
				interpolatedPR[j] = maxPrecision
			}
			maxPrecision = 0.0
			currentRecall = recall[i]
		}
		fmt.Printf("recall: %v,\tPrecision: %v\n", recall[i], precision[i])
	}
	fmt.Printf("recall tablesize: %v\n", len(recall))
	fmt.Printf("precision tablesize: %v\n", len(precision))
	for j=0; j < 11 ; j++ {
		fmt.Println(interpolatedPR[j])
	}
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

/* For Testing purposes. Note: to run main func, 
 * change package name to main.*/
func main() {
	prg := NewPRGraph()
	list := make(sortmap.PairList, 3)
	list[0] = sortmap.Pair{"326", 11.53}  
	list[1] = sortmap.Pair{"304", 9.30}  
	list[2] = sortmap.Pair{"308", 9.26}  
	prg.MakePrecRecalTable(list, "1")

}