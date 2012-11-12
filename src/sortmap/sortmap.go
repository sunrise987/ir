package sortmap

import (
	"sort"
	"fmt"
)


// A data structure to hold a key/value pair.
type Pair struct {
  Key string
  Value float64
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair
func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value }

// A function to turn a map into a PairList, then sort and return it. 
func SortMapByValue(m map[string]float64) PairList {
   p := make(PairList, len(m))
   i := 0
   for k, v := range m {
      p[i] = Pair{k, v}
      i++     
   }
   sort.Sort(p)
   return p
}


func main() {
	mymap := make(map[string]float64)
	mymap["bar"] = 0.001
	mymap["foo"] = 1.023
	fmt.Println(mymap)
	mypairs := SortMapByValue(mymap)
	fmt.Println(mypairs)
}
