package reverseindex_test

import (
	"log"
	"path"
	"path/filepath"
	. "reverseindex"
	"testing"
	//"fmt"
)

func loadIndex() *Index {
	index := NewIndex()
	//count2, data2 := ReadFile("../../data/stopwords.txt")
	//index.ListStopWords(count2, data2)

	fileNames, err := filepath.Glob("../../data/IR_Project1_Documents/*.txt")
	if err != nil {
		log.Fatal(err)
	}

	for _, fileName := range fileNames {
		count, data := ReadFile(fileName)

		fileName = path.Base(fileName)
		fileName = fileName[0 : len(fileName)-4]
		//log.Println(fileName)		
		index.MakeReverseIndex(count, data, fileName)
	}
	index.PrintStatistics()
	return index
}

func testQuery(t *testing.T, fileName string, query []string, index *Index) {
	count2, data2 := ReadFile(fileName)
	correctList := GetWordsList(count2, data2)
	docs := index.Query(query)

	for doc, _ := range docs {
		//log.Println(doc)
		if _, exists := correctList[doc]; !exists {
			t.Error(doc, " does not exists in the solution!")
		}
	}

	if len(correctList) != len(docs) {
		t.Errorf("Size of solution is wrong. Correct size: %d, we have: %d", len(correctList), len(docs))
	}
}

func TestSomeQueries(t *testing.T) {
	var query []string
	index := loadIndex()

	query = []string{"viet", "and", "nam", "and", "coup"}
	testQuery(t, "../../data/IR_Project1_Phase1/VietNamCoup.txt", query, index)

	//fmt.Println("\n\n+++++++++++++++++++++++\n")
	query = []string{"premier", "not", "khrushchev"}
	testQuery(t, "../../data/IR_Project1_Phase1/premierNotKhrushchev.txt", query, index)

	//log.Println("\n\n+++++++++++++++++++++++\n")
	query = []string{"MOSCOW", "or", "SUPPORT", "or", "AUTONOMY"}
	query = []string{"moscow", "or", "support", "or", "autonomy"}
	testQuery(t, "../../data/IR_Project1_Phase1/MoscowSupportAutonomy.txt", query, index)

}

/*
func TestPositionIndex(t *testing.T) {
	data := []byte("the one who is tired is hassan. he is sick. poor hassan. muah my love")
	MakeReverseIndex(len(data), data, "111")
	t.FailNow()
}
*/
