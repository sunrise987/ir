/*
 * To run, type: go run mymain.go [flags]
 * Flags: 
 * -stem       : Stems all words before storing them in reverse index. (ex. 'going' becomes 'go'.
 * -stopw      : Removes all stop words from reverse index. Stopwords are loaded from data file.
 * -numResults : Specifies how many top search results to print.
 * -dataDirectoryName : Can change to root path from 'data/'.
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
	"strings"
	"strconv"

	. "precrecal"
	. "reverseindex"
	. "sortmap"
)

// default is "or"
var andFlag *bool = flag.Bool("and", false, "Embed 'and' within each query.")
var orFlag *bool = flag.Bool("or", false, "Embed 'or' within each query.")

var numResults *int = flag.Int("numResults", 10, "Print top ranked numResults.")
var dataDirectoryName *string = flag.String("dataDirectoryName", "data/", "Data directory Name")

// Search results with scores below this value will be cut (deleted)
var cutValue int = 0

/**/
func main() {
	flag.Parse() // Scan the arguments list

	// Make reversed index:
	corpus := *dataDirectoryName + "IR_Project1_Documents/*.txt"
	stopwords := *dataDirectoryName + "stopwords.txt"
	index := RunMakeReverseIndex(corpus, stopwords)

	//runOnline_receivesQueryText(index)
	//runOffline_receivesFile(index, *dataDirectoryName + "/IR_Project2_Queries/Q1")
	runOffline_receivesFolder(index, *dataDirectoryName + "IR_Project2_Queries/*")
}

/**/
func RunMakeReverseIndex(corpus_pattern, stopwords_filename string) *Index {
	index := NewIndex()
	stopwords_count, stopwords := ReadFile(stopwords_filename)
	index.ListStopWords(stopwords_count, stopwords)

	fileNames, err := filepath.Glob(corpus_pattern)
	if err != nil {
		log.Fatal(err)
	}

	for _, fileName := range fileNames {
		count, data := ReadFile(fileName)
		docName := getDocNameFromFileName(fileName)
		index.MakeReverseIndex(count, data, docName)
	}
	index.PrintStatistics()
	return index
}

func getDocNameFromFileName(fileName string) string {
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
		docList := index.SearchQuery(tokens)
		fmt.Printf("Number of matches: %d\n", len(docList))
		for doc := range docList {
			if docList[doc] >= 1 {
				fmt.Printf("%s\t%.2f\n", doc, docList[doc])
			}
		}
	}
}

func runOffline_receivesFolder(index *Index, queryfiles_pattern string) {
	query_files, err := filepath.Glob(queryfiles_pattern)
	if err != nil { log.Fatal(err) }
	fmt.Sprintf("%d queries to process.\n", len(query_files))
	var queryName string
	tokens := make([]string, 50) // TODO: fix this: assuming max query size 50 
	retrievedLists := make(map[string]PairList, len(query_files))
	
	for _, queryFileName := range query_files {
		tokens = getWordsFromQueryFile(queryFileName)
		tokens = joinTokensWithOp(tokens)
		docList := index.SearchQuery(tokens)
		pairlist := SortMapByValue(docList, cutValue, *numResults)
		
		// Add resluts to map. This map will be sent to precrecal.go 
		// to create a Precision-Recall Graph.
		queryName = getQueryNumberFromQueryFileName(queryFileName)
		retrievedLists[queryName] = pairlist

		// Print top search results:
		fmt.Println()
		fmt.Println(queryFileName)
		fmt.Printf("Number of matches: %d\n", len(docList))
		fmt.Printf("DocName\tScore\n")
		pairlist.Print()
	}
	fmt.Printf("\n\n\nMaking Precision-Recall Graph...\n\n")
	// Make Interpolated Precision-Recall Graph
	prgraph := NewPRGraph()
	prgraph.MakeAvgInterpolatedPRTable(retrievedLists)
}

func runOffline_receivesFile(index *Index, queryFileName string) {
	// Get query.
	tokens := getWordsFromQueryFile(queryFileName)
	tokens = joinTokensWithOp(tokens)

	// Get query results and rank them.
	docList := index.SearchQuery(tokens)
	pairlist := SortMapByValue(docList, cutValue, *numResults)
	
	// Print top search results:
	fmt.Println(queryFileName)
	fmt.Printf("Number of matches: %d\n", len(docList))
	pairlist.Print()
	fmt.Println()

	// Make Interpolated Precision-Recall Graph
	queryName := getQueryNumberFromQueryFileName(queryFileName)
	prgraph := NewPRGraph()
	prgraph.MakeOneInterpolatedPRTable(pairlist, queryName)
}

/* Returns Query Words, out of order and without repetitions. */
func getWordsFromQueryFile(file string) []string {
	tokenMap := make(map[string]bool)
	_, data := ReadFile(file)
	lowcaseData := strings.ToLower(string(data))
	lines := strings.Split(string(lowcaseData), "\n")
	for _, line := range lines {
		tokens := strings.Split(line, " ")
		if tokens != nil {
			for i := 0; i < len(tokens); i++ {
				if tokens[i] != "" && tokens[i] != "." {
					//To handle the null-charachter case
					v, _, _, _ := strconv.UnquoteChar(tokens[i], 0)
					if v != 0 {
						tokenMap[tokens[i]] = true
					}
				}
			}
		}
	}
	returnArray := make([]string, len(tokenMap))
	index := 0
	for word := range tokenMap {
		returnArray[index] = word
		index++
	}
	return returnArray
}

// XXX: A simpler way of doing the same thing is to use strings.Join
// and then strings.Split again.
func joinTokensWithOp(tokens []string) []string {
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
