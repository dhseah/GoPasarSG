package sort

import (
	"ProjectGoLive/pkg/models"
	"math/rand"
	"testing"
	"time"
)

func TestTimSort(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(list2), func(i, j int) { list2[i], list2[j] = list2[j], list2[i] })

	ts := NewTimSort(list2, 2)
	ts.TimSort()

	for i := 0; i < len(list2); i++ {
		if list1[i].ProductID != ts.arr[i].ProductID {
			t.Errorf("sort didn't sort: got %v, want %v at index %v", list2[i].ProductID, list1[i].ProductID, i)
		}
	}
}

func BenchmarkTimSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// create a long list of products
		list := []*models.Product{}
		for j := 0; j < 500; j++ {
			list = append(list, list2...)
		}

		// shuffle the list
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })

		// initialize a TimSort
		ts := NewTimSort(list, 2)

		// perform and time the sorting
		b.StartTimer()
		ts.TimSort()
		b.StopTimer()
	}
}

func Benchmark_TimSort32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// shuffle the list
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(list2), func(i, j int) { list2[i], list2[j] = list2[j], list2[i] })

		// initialize an IntroSort
		ts := NewTimSort(list2, 2)

		// perform the sorting
		ts.TimSort()

	}
}
