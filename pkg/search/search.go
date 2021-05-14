package search

import (
	"ProjectGoLive/pkg/models"
	"sort"
	"sync"
)

// Search looks through the IndexSlice for matches for
// every search term provided by the user. If a match is
// found, the ProductID is inserted into the results map
// that maps ProductID to relevance score.
func (idx IndexSlice) Search(text string) ([]int, map[int]int) {
	var countMap = make(map[int]int)
	var wg sync.WaitGroup

	for _, token := range analyze(text) {
		wg.Add(3)
		go func() {
			defer wg.Done()
			if ids, ok := idx[0][token]; ok {
				for _, v := range ids {
					countMap[v] += 2
				}
			}
		}()
		go func() {
			defer wg.Done()
			if ids, ok := idx[1][token]; ok {
				for _, v := range ids {
					countMap[v]++
				}
			}
		}()
		go func() {
			defer wg.Done()
			if ids, ok := idx[2][token]; ok {
				for _, v := range ids {
					countMap[v] += 4
				}
			}
		}()
		wg.Wait()
	}
	unsorted := []int{}

	for k := range countMap {
		unsorted = append(unsorted, k)
	}
	return unsorted, countMap
}

// RankedProduct
// RankedProduct is called to prepare a list of products
// for writing to the http response in sorted order.
func RankedProducts(input []*models.Product, countMap map[int]int) []*models.Product {
	for i := 0; i < len(input); i++ {
		input[i].Score = countMap[input[i].ProductID]
	}
	sort.SliceStable(input, func(i, j int) bool { return input[i].Score > input[j].Score })
	return input
}
