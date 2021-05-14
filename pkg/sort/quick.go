package sort

import "ProjectGoLive/pkg/models"

type QuickSort struct {
	arr []*models.Product
	lt  func(*models.Product, *models.Product) bool
}

func NewQuickSort(data []*models.Product, sortBy int) *QuickSort {
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

	return &QuickSort{
		arr: data,
		lt:  sortLogic,
	}
}

func (s QuickSort) QuickSort(l, r int) {
	if l < r {
		pivotIdx := s.partition(l, r)
		s.QuickSort(l, pivotIdx-1)
		s.QuickSort(pivotIdx+1, r)
	}
}

func (s QuickSort) partition(l, r int) int {
	pivot := s.arr[r]
	i := l - 1
	for j := l; j <= r-1; j++ {
		if s.lt(s.arr[j], pivot) {
			i++
			s.arr[i], s.arr[j] = s.arr[j], s.arr[i]
		}
	}
	s.arr[i+1], s.arr[r] = s.arr[r], s.arr[i+1]
	return i + 1
}
