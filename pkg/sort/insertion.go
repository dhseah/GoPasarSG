package sort

import "ProjectGoLive/pkg/models"

type InsertionSort struct {
	arr []*models.Product
	lt  func(*models.Product, *models.Product) bool
}

func NewInsertionSort(data []*models.Product, sortBy int) *InsertionSort {
	var sortLogic func(*models.Product, *models.Product) bool
	switch sortBy {
	case 0:
		sortLogic = sortByPop
	case 1:
		sortLogic = sortByRatings
	case 2:
		sortLogic = sortByPriceA
	case 3:
		sortLogic = sortByPriceD
	}

	return &InsertionSort{
		arr: data,
		lt:  sortLogic,
	}
}

func (s InsertionSort) InsertionSort() {
	for i := 1; i < len(s.arr); i++ {
		key := s.arr[i]
		j := i - 1
		for j >= 0 && s.lt(key, s.arr[j]) {
			s.arr[j+1] = s.arr[j]
			j--
		}
		s.arr[j+1] = key
	}
}
