package reverseindex

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"stemmer"
	"strings"
)

var doStemmingFlag *bool = flag.Bool("stem", false, "Do stemming.")
var doStopwordsFlag *bool = flag.Bool("stopw", false, "Remove stop words.")

/* A map from a document id to an array of index positions in this document.
 * These positions keep track of where a given term is found in this document.
 * The length of this array is the frequency count for the given term in this document.
 * Each term in our index will have a DocList. */
type DocList map[string][]int

type Statistics struct {
	numDocs  int
	numWords int
	numOnes  int
	/* Longest Posting List */
	longestPLword string
	longestPLsize int
	/* Shortest Posting List */
	shortestPLword string
	shortestPLsize int
}

type Index struct {
	reverseIndex map[string]DocList
	stopwords    map[string]bool
	stat         Statistics
}

func NewIndex() *Index {
	index := &Index{}
	index.init()
	return index
}

func (index *Index) init() {
	index.reverseIndex = make(map[string]DocList)
	index.stopwords = make(map[string]bool)
	index.stat.numDocs = 0
}

func (index *Index) getQueryWeights(query []string) map[string]float64 {
	//termWeight = log (N/df) = totalNumDocuments/ numDocumentsThatContainThisTerm.
	totalNumDocs := float64(index.stat.numDocs)
	termWeights := make(map[string]float64)
	for _, word := range query {
		docFreq_t := float64(len(index.reverseIndex[word]))
		termWeights[word] = totalNumDocs / docFreq_t
	}
	return termWeights
}

/*
 * Takes a query and returns the list of documents containing this qurey
 * For phrase queries, query format must be "[word1] cand [word2] cand ..." 
 */
func (index *Index) SearchQuery(query []string) map[string]float64 {
	// Find weights of query terms. Type map[string]float64
	queryTermWeight := index.getQueryWeights(query)

	// Stores the cosine normalization factor for each document.
	norm_d := make(map[string]float64)

	// Stores the cosine similarity measure between the query and document. Initialized to prevent aliacing.
	rankingList := make(map[string]float64)
	for d := range index.reverseIndex[query[0]] {
		rankingList[d] = 0
		norm_d[d] = 0
	}

	i := 2
	lastWord := ""
	for _, word := range query {
		word = strings.ToLower(string(word))

		// If the word is a connector
		if word == "not" {
			i = 0
			continue
		} else if word == "or" {
			i = 1
			continue
		} else if word == "and" {
			i = 2
			continue
		} else if word == "cand" {
			i = 3
			continue
		}
		// If the word is not a connector
		if *doStemmingFlag {
			word = string(stemmer.Stem([]byte(word)))
		}
		// WARNING: templist is an aliace variable pointing to the index.reverseIndex. 
		// Do NOT modify this variable.
		templist := index.reverseIndex[word]
		docFreq := float64(len(templist))
		invDocFreq := math.Log10(float64(index.stat.numDocs) / docFreq)
		switch i {
		case 3:
			// Perform 'cand' operation.
			for doc := range rankingList {
				if _, exists := templist[doc]; !exists {
					delete(rankingList, doc)
					delete(norm_d, doc)
				} else if lastWord != "" {
					// Check if it is consecutive with each other in the each document.
					if !consecutive(index.reverseIndex[lastWord][doc], templist[doc]) {
						delete(rankingList, doc)
						delete(norm_d, doc)
					} else {
						// Calculate score for this document.
						termFreq := float64(len(templist[doc]))
						termFreqLog := 1 + math.Log10(termFreq)
						docTermWeight := termFreqLog * invDocFreq
						rankingList[doc] += (docTermWeight * queryTermWeight[word])
						norm_d[doc] += (docTermWeight * docTermWeight)
					}
				}
			}
		case 2:
			// Perform 'and' operation.
			for doc := range rankingList {
				if _, exists := templist[doc]; !exists {
					delete(rankingList, doc)
					delete(norm_d, doc)
				} else {
					// Calculate score for this document.
					termFreq := float64(len(templist[doc]))
					termFreqLog := 1 + math.Log10(termFreq)
					docTermWeight := termFreqLog * invDocFreq
					rankingList[doc] += (docTermWeight * queryTermWeight[word])
					norm_d[doc] += (docTermWeight * docTermWeight)
				}
			}
		case 1:
			// Perform 'or' operation.
			for doc := range templist {
				if _, exists := rankingList[doc]; !exists {
					rankingList[doc] = 0
					norm_d[doc] = 0
				}
				// Calculate score for this document.
				termFreq := float64(len(templist[doc]))
				termFreqLog := 1 + math.Log10(termFreq)
				docTermWeight := termFreqLog * invDocFreq
				rankingList[doc] += (docTermWeight * queryTermWeight[word])
				norm_d[doc] += (docTermWeight * docTermWeight)
			}

		case 0:
			// Perform 'not' operation.
			for doc := range rankingList {
				if _, exists := templist[doc]; exists {
					delete(rankingList, doc)
					delete(norm_d, doc)
				}
			}
		}
		lastWord = word
		//fmt.Printf("last word : %s\n", lastWord)
	}

	// Calculate the cosine normalization factor for the query.
	norm_q := 0.0
	for _, word := range query {
		norm_q += (queryTermWeight[word] * queryTermWeight[word])
	}
	norm_q = math.Sqrt(norm_q)

	// Divide rank by normalization factors.
	for doc := range rankingList {
		norm_d[doc] = math.Sqrt(norm_d[doc])
		rankingList[doc] = (rankingList[doc] / (norm_d[doc] * norm_q))
	}

	return rankingList
}

/* Both arrays 'last' and 'current' store positional values (int) for two different 
 * words in the same file. This method checks if the 'last' word is located at the
 * immediate previous position then the 'current' word. So, this function searches 
 * to find two consecutive indexes. If found, it returns true. Otherwise, false. */
func consecutive(last []int, current []int) bool {
	lastI := 0
	curI := 0
	for curI < len(current) && lastI < len(last) {
		li := int(math.Min(float64(lastI), float64(len(last)-1)))
		ci := int(math.Min(float64(curI), float64(len(current)-1)))

		if last[li]+1 == current[ci] {
			return true
		} else if last[li] > current[ci] {
			curI++
		} else if last[li] < current[ci] {
			lastI++
		} else if last[li] == current[ci] {
			log.Fatal("two words exist in the same location at the same document!")
		}

	}
	return false
}

func GetWordsList(count int, data []byte) map[string]bool {
	i := 0
	word := ""
	mymap := make(map[string]bool)

	for i < count {
		word, i = GetNextWord(data, count, i)
		word = strings.ToLower(word)
		if word != "" {
			mymap[word] = true
		}
	}
	return mymap
}

func (index *Index) ListStopWords(count int, data []byte) {
	i := 0
	word := ""
	for i < count {
		word, i = GetNextWord(data, count, i)
		word = strings.ToLower(word)
		if word != "" {
			index.stopwords[word] = true
		}
	}
	fmt.Print()
	//fmt.Println(stopwords)
}

func ReadFile(fileName string) (int, []byte) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 500000) // TODO: Fix the size of the data array
	count, err := file.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("- read %d bytes: %q\n\n", count, data[:count])	
	return count, data
}

func (index *Index) MakeReverseIndex(count int, data []byte, fileName string) {
	// This represent the count inside each document. 
	// Will be used to note the position of a term in a document.
	numWords := 0
	i := 0
	word := ""
	fileNameAdded := false

	for i < count {
		word, i = GetNextWord(data, count, i)

		if _, exists := index.stopwords[word]; (exists && *doStopwordsFlag) || word == "" {
			// Ignore the word and not put in the index
			continue
		}
		numWords++
		word = strings.ToLower(word)

		if *doStemmingFlag {
			word = string(stemmer.Stem([]byte(word)))
		}

		if _, exists := index.reverseIndex[word]; exists {
			if _, entryExists := index.reverseIndex[word][fileName]; entryExists {
				index.reverseIndex[word][fileName] = append(index.reverseIndex[word][fileName], numWords)
			} else {
				index.reverseIndex[word][fileName] = make([]int, 0, 10)
				index.reverseIndex[word][fileName] = append(index.reverseIndex[word][fileName], numWords)
				// Compute statistics
				fileNameAdded = true
			}
			//fmt.Println(index.reverseIndex[word][fileName])
		} else {
			index.reverseIndex[word] = make(map[string][]int)
			index.reverseIndex[word][fileName] = make([]int, 0, 10)
			index.reverseIndex[word][fileName] = append(index.reverseIndex[word][fileName], numWords)
			// Compute statistics
			fileNameAdded = true
		}
		//fmt.Println(index.reverseIndex)		
	}
	if fileNameAdded {
		index.stat.numDocs++
	}
	//fmt.Println(index.reverseIndex)	
}

func (index *Index) computeStats() {
	index.stat.shortestPLsize = 100000
	index.stat.longestPLsize = 0
	index.stat.numOnes = 0
	index.stat.numWords = 0

	for word := range index.reverseIndex {
		if len(index.reverseIndex[word]) < index.stat.shortestPLsize {
			index.stat.shortestPLsize = len(index.reverseIndex[word])
			index.stat.shortestPLword = word
		}
		if len(index.reverseIndex[word]) > index.stat.longestPLsize {
			index.stat.longestPLsize = len(index.reverseIndex[word])
			index.stat.longestPLword = word
		}
		index.stat.numWords++
		index.stat.numOnes += len(index.reverseIndex[word])
	}
}

func (index *Index) PrintStatistics() {
	index.computeStats()
	fmt.Printf("The number of ones in this matrix: \t%d.\n", index.stat.numOnes)
	fmt.Printf("The size of the term-document matrix: \t%d x %d.\n", index.stat.numWords, index.stat.numDocs)
	fmt.Printf("The longest posting list: \t%d\n", index.stat.longestPLsize)
	fmt.Printf("'%s'\n", index.stat.longestPLword)
	fmt.Printf("The shortest posting list: \t%d\n", index.stat.shortestPLsize)
	fmt.Printf("'%s'\n", index.stat.shortestPLword)
}

/*
 * Returns the next word from the []byte data. Splits according to spaces and Punktuation marks. 
 * It includes numbers in words. It returns "" if called at the end of the []byte array. 
 */
func GetNextWord(data []byte, count, index int) (string, int) {
	word := ""

	for _, c := range data[index:count] {
		index++
		//fmt.Println("word: ", word)
		//fmt.Println("index: ", index)
		//fmt.Println("c: ", c)

		// Any new line or space
		if (8 <= c && c <= 10) || (32 <= c && c <= 47) || (58 <= c && c <= 64) {
			//fmt.Println("Exit at IF statement")			
			return word, index

			// A word starts or continues
		} else {
			word = word + string(c)
			//fmt.Println("\nNo exit at ELSE statement")					
		}
	}
	//fmt.Println("Exit at FOR loop")	
	return word, index
}

/**
 * Okt 31, 2012 
 */

/* Takes a list of qury terms and scores the document for this query. */
func (index *Index) score(query []string, docs DocList) float64 {
	for doc, _ := range docs {
		score := 0.0
		for _, term := range query {
			fmt.Println(term)
			invDocFreq := float64(index.stat.numDocs) / float64(len(index.reverseIndex[term]))
			termFreq := float64(len(index.reverseIndex[term][doc]))
			score += (1 + math.Log10(termFreq)) * math.Log10(invDocFreq)
		}
	}

	return 0.0
}

/**
 * Nov 30, 2012
 */
func (index *Index) calculateWeightVector() {}//[]float64 {}

func (index *Index) createNewQuery(alpha, beta float64) {}
