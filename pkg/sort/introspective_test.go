package sort

import (
	"ProjectGoLive/pkg/models"
	"math/rand"
	"testing"
	"time"
)

func Test_IntroSort(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(list2), func(i, j int) { list2[i], list2[j] = list2[j], list2[i] })

	is := NewIntroSort(list2, 2)
	is.IntroSort()

	for i := 0; i < len(list2); i++ {
		if list1[i].ProductID != is.arr[i].ProductID {
			t.Errorf("sort didn't sort: got %v, want %v at index %v", list2[i].ProductID, list1[i].ProductID, i)
		}
	}
}

func Benchmark_IntroSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// create a long list of products
		list := []*models.Product{}
		for j := 0; j < 500; j++ {
			list = append(list, list2...)
		}

		// shuffle the list
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })

		// initialize an IntroSort
		is := NewIntroSort(list, 2)

		// perform and time the sorting
		b.StartTimer()
		is.IntroSort()
		b.StopTimer()
	}
}

func Benchmark_IntroSort32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// shuffle the list
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(list2), func(i, j int) { list2[i], list2[j] = list2[j], list2[i] })

		// initialize an IntroSort
		is := NewIntroSort(list2, 2)

		// perform the sorting
		is.IntroSort()
	}
}
