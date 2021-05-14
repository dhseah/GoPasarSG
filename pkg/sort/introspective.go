package sort

import (
	"ProjectGoLive/pkg/models"
	"math"
)

type IntroSort struct {
	arr        []*models.Product
	lt         func(*models.Product, *models.Product) bool
	depthLimit int
}

func NewIntroSort(data []*models.Product, sortBy int) *IntroSort {
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

	return &IntroSort{
		arr: data,
		lt:  sortLogic,
	}
}

func (s IntroSort) IntroSort() {
	begin := 0
	end := len(s.arr) - 1

	s.depthLimit = 2 * int(math.Round(math.Log2(float64(end))))
	s.introSortUtil(begin, end)
}

func (s IntroSort) introSortUtil(begin, end int) {
	size := end - begin
	if size <= 16 {
		s.insertionSort(begin, end)
		return
	}

	if s.depthLimit == 0 {
		s.heapSort(begin, end)
		return
	}
	s.depthLimit--

	// perform quick sort
	pivot := s.medianOfThree(begin, begin+size/2, end)
	s.arr[pivot], s.arr[end] = s.arr[end], s.arr[pivot]

	pivotIdx := s.partition(begin, end)

	s.introSortUtil(begin, pivotIdx-1)
	s.introSortUtil(pivotIdx+1, end)
}

func (s IntroSort) partition(low, high int) int {
	pivot := s.arr[high]
	i := low - 1
	for j := low; j < high; j++ {
		if s.lt(s.arr[j], pivot) {
			i++
			s.arr[i], s.arr[j] = s.arr[j], s.arr[i]
		}
	}
	s.arr[i+1], s.arr[high] = s.arr[high], s.arr[i+1]
	return i + 1
}

func (s IntroSort) medianOfThree(a, b, c int) int {
	da := s.arr[a]
	db := s.arr[b]
	dc := s.arr[c]

	if s.lt(da, db) && s.lt(db, dc) {
		return b
	}

	if s.lt(dc, db) && s.lt(db, da) {
		return b
	}

	if s.lt(db, da) && s.lt(da, dc) {
		return a
	}
	if s.lt(dc, da) && s.lt(da, db) {
		return a
	}

	if s.lt(da, dc) && s.lt(dc, db) {
		return c
	}

	if s.lt(db, dc) && s.lt(dc, da) {
		return c
	}

	return a
}

func (s IntroSort) insertionSort(begin, end int) {
	left := begin

	for i := left + 1; i <= end; i++ {
		key := s.arr[i]
		j := i - 1
		for j >= left && s.lt(key, s.arr[j]) {
			s.arr[j+1] = s.arr[j]
			j--
		}
		s.arr[j+1] = key
	}
}

func (s IntroSort) maxHeap(i, n, begin int) {
	temp := s.arr[begin+i-1]
	child := 0

	for i <= n/2 {
		child = 2 * i
		if child < n && s.lt(s.arr[begin+child-1], s.arr[begin+child]) {
			child++
		}

		if s.lt(s.arr[begin+child-1], temp) {
			break
		}

		s.arr[begin+i-1] = s.arr[begin+child-1]
		i = child
	}
	s.arr[begin+i-1] = temp
}

func (s IntroSort) heapify(begin, end, n int) {
	for i := n / 2; i >= 1; i-- {
		s.maxHeap(i, n, begin)
	}
}

func (s IntroSort) heapSort(begin, end int) {
	n := end - begin
	s.heapify(begin, end, n)
	for i := n; i >= 1; i-- {
		s.arr[begin], s.arr[begin+1] = s.arr[begin+1], s.arr[begin]
		s.maxHeap(1, i, begin)
	}
}
