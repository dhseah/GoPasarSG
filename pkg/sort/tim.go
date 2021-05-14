package sort

import "ProjectGoLive/pkg/models"

const run = 16

type TimSort struct {
	arr []*models.Product
	lt  func(*models.Product, *models.Product) bool
}

func NewTimSort(data []*models.Product, sortBy int) *TimSort {
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

	return &TimSort{
		arr: data,
		lt:  sortLogic,
	}
}

// TimSort is a naive implementation of the TimSort
// which combines insertion sort and merge sort but
// does not implement further optimizations during
// the merge phase.
func (s *TimSort) TimSort() {
	n := len(s.arr)
	for i := 0; i < n; i += run {
		if (i + run - 1) < (n - 1) {
			s.insertionSort(i, i+run-1)
		} else {
			s.insertionSort(i, n-1)
		}
	}

	for size := run; size < n; size *= 2 {
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

func (s *TimSort) insertionSort(l, r int) {
	for i := l + 1; i <= r; i++ {
		temp := s.arr[i]
		j := i - 1
		for j >= l && s.lt(temp, s.arr[j]) {
			s.arr[j+1] = s.arr[j]
			j--
		}
		s.arr[j+1] = temp
	}
}

func (s *TimSort) merge(l, m, r int) {
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
