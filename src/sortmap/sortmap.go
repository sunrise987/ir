// XXX: Remove unused commented code.
//package main
package sortmap

import (
	"fmt"
	"sort"
)

// A data structure to hold a Key/Value Pair.
type Pair struct {
	Key   string
	Value float64
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value }

func (list PairList) Print() {
	for _, pair := range list {
		fmt.Printf("%s\t%.2f\n", pair.Key, pair.Value)
	}
}

// A function to turn a map into a PairList, then sort and return it. 
// any Value below cutValue will be deleted.
func sortMapByValue(m map[string]float64, cutValue int) PairList {
	p := make(PairList, 0, len(m))
	for k, v := range m {
		if v > float64(cutValue) {
			p = append(p, Pair{k, v})
		}
	}
	sort.Sort(p)
	return p
}

func SortMapByValue(m map[string]float64, cutValue int) PairList {
	list := sortMapByValue(m, cutValue)
	return list
}

func getMapFromList(list PairList) map[string]float64 {
	returnMap := make(map[string]float64, len(list))
	for _, pair := range list {
		returnMap[pair.Key] = pair.Value
	}
	return returnMap
}

// XXX: Don't hard code topTen into the code. Pass a parameter top_n instead.
// Also rename the function to SortMapByValueTopN. Don't put underscores into
// function names.

// A function to turn a map into a PairList, then sort and return it.       
// any Value below cutValue will be deleted.                   
func SortMapByValue_topTen(m map[string]float64, cutValue int) PairList {
	p := sortMapByValue(m, cutValue)
	topTen := make(PairList, 10)
	topTen = p[0:10]
	return topTen
}

// XXX: Put this into a sortmap_test.go.
func main() {
	mymap := make(map[string]float64)
	mymap["bar1"] = 1.001
	mymap["bar2"] = 3.701
	mymap["bar3"] = 0.001
	mymap["bar4"] = 0.901
	mymap["bar5"] = 6.001
	mymap["bar6"] = 6.301
	mymap["bar7"] = 6.031
	mymap["foo1"] = 1.053
	mymap["foo2"] = 1.083
	mymap["foo3"] = 1.093
	mymap["foo4"] = 1.23
	mymap["foo5"] = 1.23
	mymap["foo6"] = 1.028
	mymap["foo7"] = 1.123
	fmt.Println(len(mymap))
	fmt.Println()
	cutValue := 1
	myPairs := SortMapByValue(mymap, cutValue)
	fmt.Println(myPairs)
	fmt.Println(len(myPairs))
	fmt.Println()
	myPairs = SortMapByValue_topTen(mymap, cutValue)
	fmt.Println(myPairs)
	fmt.Println(len(myPairs))
}
