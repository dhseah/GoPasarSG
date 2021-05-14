package sort

import (
	"ProjectGoLive/pkg/models"
	"math/rand"
	"testing"
	"time"
)

func Test_InsertionSort(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(list2), func(i, j int) { list2[i], list2[j] = list2[j], list2[i] })

	is := NewInsertionSort(list2, 2)
	is.InsertionSort()

	for i := 0; i < len(list2); i++ {
		if list1[i].ProductID != is.arr[i].ProductID {
			t.Errorf("sort didn't sort: got %v, want %v at index %v", list2[i].ProductID, list1[i].ProductID, i)
		}
	}
}

func Benchmark_InsertionSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// create a long list of products
		list := []*models.Product{}
		for j := 0; j < 500; j++ {
			list = append(list, list2...)
		}

		// shuffle the list
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })

		// initialize an InsertionSort
		is := NewInsertionSort(list, 2)

		// perform and time the sorting
		b.StartTimer()
		is.InsertionSort()
		b.StopTimer()
	}
}
