package sort

import "ProjectGoLive/pkg/models"

type MergeSort struct {
	arr []*models.Product
	lt  func(*models.Product, *models.Product) bool
}

func NewMergeSort(data []*models.Product, sortBy int) *MergeSort {
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

	return &MergeSort{
		arr: data,
		lt:  sortLogic,
	}
}

func (s *MergeSort) MergeSort() {
	n := len(s.arr)
	for size := 1; size < n; size *= 2 {
		for left := 0; left < n; left += 2 * size {
			mid := left + size - 1
			right := left + 2*size - 1
			if right > n-1 {
				right = n - 1
			}

			if mid < right {
				s.merge(left, mid, right)
			}
		}
	}
}

func (s *MergeSort) merge(l, m, r int) {
	len1 := m - l + 1
	len2 := r - m
	left := make([]*models.Product, len1)
	right := make([]*models.Product, len2)

	for i := 0; i < len1; i++ {
		left[i] = s.arr[l+i]
	}
	for i := 0; i < len2; i++ {
		right[i] = s.arr[m+1+i]
	}

	i := 0
	j := 0
	k := l

	for i < len1 && j < len2 {
		if s.lt(left[i], right[j]) {
			s.arr[k] = left[i]
			i++
		} else {
			s.arr[k] = right[j]
			j++
		}
		k++
	}
	for i < len1 {
		s.arr[k] = left[i]
		k++
		i++
	}
	for j < len2 {
		s.arr[k] = right[j]
		k++
		j++
	}
}
