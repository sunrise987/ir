package sortmap

import (
	"fmt"
	"sort"
	"bytes"
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

func (list PairList) PrintToString() string{
	var buffer bytes.Buffer
	for _, pair := range list {
		buffer.WriteString(fmt.Sprintf("%s\t%.2f\n", pair.Key, pair.Value))
	}
	return buffer.String()
}

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
// A function to turn a map into a PairList, then sort and return it.       
// any Value below cutValue will be deleted.                   
// This function only return topN values ranked.
func SortMapByValue(m map[string]float64, cutValue int, topN int) PairList {
	list := sortMapByValue(m, cutValue)
	if topN != -1 {
		topNlist := make(PairList, topN)
		topNlist = list[0:topN]
		return topNlist

	}
	return list
}

func getMapFromList(list PairList) map[string]float64 {
	returnMap := make(map[string]float64, len(list))
	for _, pair := range list {
		returnMap[pair.Key] = pair.Value
	}
	return returnMap
}

