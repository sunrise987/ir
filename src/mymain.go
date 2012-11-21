/*
 * To run, type: go run mymain.go -stem -stopw -ten
 */
package main 

import (
    "flag"
    . "reverseindex"
    "os"
    "fmt"
    "path/filepath"
    "path"
    "log"
    "bufio"
    "strings"
	. "sortmap"
)

const APP_VERSION = "0.1"
// default is "or"
var andFlag *bool = flag.Bool("and", false, "Embed 'and' within each query.")
var orFlag *bool = flag.Bool("or", false, "Embed 'or' within each query.")

var tenFlag *bool = flag.Bool("ten", false, "Print only 10 results.")

func main() {
/*
	//root := "../"
	root := ""
	corpus := root + "data/IR_Project1_Documents/*.txt"
	stopwords := root + "data/stopwords.txt"
    index := RunMakeReverseIndex(corpus, stopwords)

	
	
	//runOnline(index)	
	runOffline_receivesFolder(index, root + "data/IR_Project2_Queries/*")
	*/
	
	readCorrectAnswersToMap("data/RelevancyLists.txt")
}

func RunMakeReverseIndex(corpus, stopwords string) *Index {
	flag.Parse() // Scan the arguments list
   	
   	index := NewIndex()
    
	count2, data2 := ReadFile(stopwords)
	index.ListStopWords(count2, data2)
	
	
	fileNames, err := filepath.Glob(corpus)
	if err != nil {
		log.Fatal(err)
	}
	
	for _, fileName := range fileNames{	
		count, data := ReadFile(fileName)
		
		fileName = path.Base(fileName)
		fileName = fileName[3:len(fileName)-4]
		//fmt.Println(fileName)		
		index.MakeReverseIndex(count, data, fileName)
	}
	index.PrintStatistics()	
	return index
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
   		line, _, _ := stdin.ReadLine();
   		
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
	fileNames, err := filepath.Glob(address)
	if err != nil {log.Fatal(err)}
	
	//os.Remove("data/BasicResults.txt")
	//resultFile, e := os.Create("data/BasicResults.txt")
	//if e != nil {log.Fatal(e)
	
	fmt.Printf("%d queries processed.\n", len(fileNames))
	
	for _, fileName := range fileNames{
		fmt.Println()
		fmt.Println(fileName)		
		_, data := ReadFile(fileName)		
		lowcaseline := strings.ToLower(string(data))
    	tokens := strings.Split(lowcaseline, " ")
    	tokens = addANDS(tokens)

       	docList := index.Query(tokens)    	
		fmt.Printf("Number of matches: %d\n\n", len(docList))
		fmt.Printf("DocName\tScore\n")
		pairlist := SortMapByValue(docList)
		counter := 0
		for _, pair := range pairlist {
			if pair.Value >= 1 {
			 	if *tenFlag {
			 		if counter <= 10{
						fmt.Printf("%s\t%.2f\n", pair.Key, pair.Value)
					}
				} else {
					fmt.Printf("%s\t%.2f\n", pair.Key, pair.Value)
				}
			}
			counter++
		}
		fmt.Println()
										
	}
}

func addANDS(tokens []string) []string {
	op := "or"
	if *andFlag {
		op = "and"
	} else if *orFlag {
		op = "or"
	}
	newTokens := make([]string, (len(tokens)*2))
	i := 0
	for _,token := range tokens {
		newTokens[i] = token
		newTokens[i+1] = op
		i += 2
	}
	return newTokens
}

func readCorrectAnswersToMap(file string) map[string][]string {
	correctAnswers := make(map[string][]string)
	_, data := ReadFile(file)
	lines := strings.Split(string(data), "\r\n")
	
	
	for _, line := range lines {
		tokens := strings.Split(line, " ")
		if tokens != nil {
			fmt.Println(tokens)
		}
		correctAnswers[tokens[0]] = tokens[1:len(tokens)]		
	}
	//fmt.Println("\n\n\n -----------------")
	//fmt.Println(correctAnswers)
	//fmt.Println("\n\n\n -----------------")
	return correctAnswers
}


