/*
 * To run, type: go run mymain.go -stem -stopw -ten
 */
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	. "precrecal"
	. "reverseindex"
	. "sortmap"
	"strings"
)

const APP_VERSION = "0.1"

// default is "or"
var andFlag *bool = flag.Bool("and", false, "Embed 'and' within each query.")
var orFlag *bool = flag.Bool("or", false, "Embed 'or' within each query.")

var tenFlag *bool = flag.Bool("ten", false, "Print only 10 results.")

/**/
func main() {

	// Make reversed index:
	//root := "../"
	root := ""
	corpus := root + "data/IR_Project1_Documents/*.txt"
	stopwords := root + "data/stopwords.txt"
	index := RunMakeReverseIndex(corpus, stopwords)

	//runOnline_receivesQueryText(index)
	runOffline_receivesFile(index, root+"data/IR_Project2_Queries/Q1")
	//runOffline_receivesFolder(index, root + "data/IR_Project2_Queries/*")

}

/**/
func RunMakeReverseIndex(corpus, stopwords string) *Index {
	flag.Parse() // Scan the arguments list

	index := NewIndex()

	count2, data2 := ReadFile(stopwords)
	index.ListStopWords(count2, data2)

	fileNames, err := filepath.Glob(corpus)
	if err != nil {
		log.Fatal(err)
	}

	for _, fileName := range fileNames {
		count, data := ReadFile(fileName)

		queryNum := getQueryNumberFromFileName(fileName)
		//fmt.Println(fileName)		
		index.MakeReverseIndex(count, data, queryNum)
	}
	index.PrintStatistics()
	return index
}

func getQueryNumberFromFileName(fileName string) string {
	fileName = path.Base(fileName)
	queryNum := fileName[3 : len(fileName)-4]
	return queryNum
}

func getQueryNumberFromQueryFileName(fileName string) string {
	fileName = path.Base(fileName)
	queryNum := fileName[1:len(fileName)]
	return queryNum
}

/**
 * on the colsole: type a query after you see "Search:"
 * type "quit" to exit.
 */
func runOnline_receivesQueryText(index *Index) {
	for {
		// Get input form the Console.
		fmt.Println("\nSearch:")
		stdin := bufio.NewReaderSize(os.Stdin, 100)
		line, _, _ := stdin.ReadLine()

		if string(line) == "quit" {
			break
		}

		lowcaseline := strings.ToLower(string(line))
		tokens := strings.Split(lowcaseline, " ")
		docList := index.Query(tokens)
		fmt.Printf("Number of matches: %d\n", len(docList))
		for doc := range docList {
			if docList[doc] >= 1 {
				fmt.Printf("%s\t%.2f\n", doc, docList[doc])
			}
		}
	}
}

func runOffline_receivesFolder(index *Index, address string) {
	QueryFiles, err := filepath.Glob(address)
	fmt.Printf("%d queries to process.\n", len(QueryFiles))
	if err != nil {
		log.Fatal(err)
	}
	//os.Remove("data/BasicResults.txt")
	//resultFile, e := os.Create("data/BasicResults.txt")
	//if e != nil {log.Fatal(e)

	var queryNum string
	cutValue := 1                // To delete scores below 1
	tokens := make([]string, 50) // TODO: fix this: assuming max query size 50 
	var pairlist PairList
	retrievedLists := make(map[string]PairList, len(QueryFiles))
	fmt.Printf("Double Check: there are %v queries proccessed.", len(QueryFiles))

	for _, queryFileName := range QueryFiles {
		tokens = getTokensFromFile(queryFileName)
		tokens = addANDS(tokens)
		docList := index.Query(tokens)

		if *tenFlag {
			pairlist = SortMapByValue_topTen(docList, cutValue)
		} else {
			pairlist = SortMapByValue(docList, cutValue)
		}

		// Add resluts to map. This map will be sent to precrecal.go 
		// to create a Precision-Recall Graph.
		queryNum = getQueryNumberFromQueryFileName(queryFileName)
		retrievedLists[queryNum] = pairlist

		// Print top search results:
		fmt.Println()
		fmt.Println(queryFileName)
		fmt.Printf("Number of matches: %d\n\n", len(docList))
		fmt.Printf("DocName\tScore\n")
		pairlist.Print()
		fmt.Println()
	}

	// Make Interpolated Precision-Recall Graph
	prgraph := NewPRGraph()
	prgraph.MakeAvgInterpolatedPRTable(retrievedLists)
}

func runOffline_receivesFile(index *Index, queryFileName string) {
	// Get Query.
	tokens := getTokensFromFile(queryFileName)
	tokens = addANDS(tokens)

	// Get Query Results.
	docList := index.Query(tokens)

	// Rank Results.
	cutValue := 1 // To delete scores below 1.
	var pairlist PairList
	if *tenFlag {
		pairlist = SortMapByValue_topTen(docList, cutValue)
	} else {
		pairlist = SortMapByValue(docList, cutValue)
	}

	// Print top search results:
	fmt.Println(queryFileName)
	fmt.Printf("Number of matches: %d\n\n", len(docList))
	fmt.Printf("DocName\tScore\n")
	pairlist.Print()
	fmt.Println()

	// Make Interpolated Precision-Recall Graph
	queryNum := getQueryNumberFromQueryFileName(queryFileName)
	prgraph := NewPRGraph()
	prgraph.MakeOneInterpolatedPRTable(pairlist, queryNum)
}

func getTokensFromFile(fileName string) []string {
	_, data := ReadFile(fileName)
	lowcaseline := strings.ToLower(string(data))
	tokens := strings.Split(lowcaseline, " ")

	return tokens
}

func addANDS(tokens []string) []string {
	op := "or"
	if *andFlag {
		op = "and"
	} else if *orFlag {
		op = "or"
	}
	newTokens := make([]string, (len(tokens) * 2))
	i := 0
	for _, token := range tokens {
		newTokens[i] = token
		newTokens[i+1] = op
		i += 2
	}
	return newTokens
}
