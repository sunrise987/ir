package sortmap_test

import (
	"testing"
	"log"

	. "sortmap"
)

func test1(t *testing.T) {
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


	//sortedPairs := make([]Pairs, 14)
//	sortedPairs[0] = {"
	t.Error("")

	cutValue := 1
	topN := 10

	myPairs := SortMapByValue(mymap, cutValue, -1)
	log.Println(myPairs)

	if len(mymap) != len(myPairs) {
		t.Error("Error. Origional list size is %v. Sorted list size is %v.", len(mymap), len(myPairs))
	}
	myPairs = SortMapByValue(mymap, cutValue, topN)
	log.Println(myPairs)

	if topN != len(myPairs) {
		t.Error("Error. Requested top %v elements. Sorted list size is %v.", len(mymap), len(myPairs))
	}
}
