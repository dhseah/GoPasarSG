package sort

import (
	"ProjectGoLive/pkg/models"
	"testing"
)

func Test_sortByPriceA(t *testing.T) {
	type args struct {
		p1 *models.Product
		p2 *models.Product
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Price of p1 < p2",
			args: args{
				&models.Product{ProductID: 1, Name: "Product 1", Desc: "This is Product 1", Price: 1.00, Inventory: 100, Rating: 3.0, RatingNum: 80, UnitSold: 100},
				&models.Product{ProductID: 2, Name: "Product 2", Desc: "This is Product 2", Price: 2.00, Inventory: 100, Rating: 3.0, RatingNum: 80, UnitSold: 100},
			},
			want: true,
		},
		{
			name: "UnitSold of p1 > p2",
			args: args{
				&models.Product{ProductID: 1, Name: "Product 1", Desc: "This is Product 1", Price: 1.00, Inventory: 100, Rating: 3.0, RatingNum: 80, UnitSold: 100},
				&models.Product{ProductID: 2, Name: "Product 2", Desc: "This is Product 2", Price: 1.00, Inventory: 100, Rating: 3.0, RatingNum: 80, UnitSold: 90},
			},
			want: true,
		},
		{
			name: "Rating of p1 > p2",
			args: args{
				&models.Product{ProductID: 1, Name: "Product 1", Desc: "This is Product 1", Price: 1.00, Inventory: 100, Rating: 3.5, RatingNum: 80, UnitSold: 100},
				&models.Product{ProductID: 2, Name: "Product 2", Desc: "This is Product 2", Price: 1.00, Inventory: 100, Rating: 3.0, RatingNum: 80, UnitSold: 100},
			},
			want: true,
		},
		{
			name: "RatingNum of p1 > p2",
			args: args{
				&models.Product{ProductID: 1, Name: "Product 1", Desc: "This is Product 1", Price: 1.00, Inventory: 100, Rating: 3.5, RatingNum: 80, UnitSold: 100},
				&models.Product{ProductID: 2, Name: "Product 2", Desc: "This is Product 2", Price: 1.00, Inventory: 100, Rating: 3.5, RatingNum: 60, UnitSold: 100},
			},
			want: true,
		},
		{
			name: "Inventory of p1 > p2",
			args: args{
				&models.Product{ProductID: 1, Name: "Product 1", Desc: "This is Product 1", Price: 1.00, Inventory: 200, Rating: 3.5, RatingNum: 80, UnitSold: 100},
				&models.Product{ProductID: 2, Name: "Product 2", Desc: "This is Product 2", Price: 1.00, Inventory: 100, Rating: 3.5, RatingNum: 80, UnitSold: 100},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortByPriceA(tt.args.p1, tt.args.p2); got != tt.want {
				t.Errorf("sortByPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortByPriceV(t *testing.T) {
	type args struct {
		p1 *models.Product
		p2 *models.Product
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Price of p1 < p2",
			args: args{
				&models.Product{ProductID: 1, Name: "Product 1", Desc: "This is Product 1", Price: 1.00, Inventory: 100, Rating: 3.0, RatingNum: 80, UnitSold: 100},
				&models.Product{ProductID: 2, Name: "Product 2", Desc: "This is Product 2", Price: 2.00, Inventory: 100, Rating: 3.0, RatingNum: 80, UnitSold: 100},
			},
			want: true,
		},
		{
			name: "UnitSold of p1 > p2",
			args: args{
				&models.Product{ProductID: 1, Name: "Product 1", Desc: "This is Product 1", Price: 1.00, Inventory: 100, Rating: 3.0, RatingNum: 80, UnitSold: 100},
				&models.Product{ProductID: 2, Name: "Product 2", Desc: "This is Product 2", Price: 1.00, Inventory: 100, Rating: 3.0, RatingNum: 80, UnitSold: 90},
			},
			want: true,
		},
		{
			name: "Rating of p1 > p2",
			args: args{
				&models.Product{ProductID: 1, Name: "Product 1", Desc: "This is Product 1", Price: 1.00, Inventory: 100, Rating: 3.5, RatingNum: 80, UnitSold: 100},
				&models.Product{ProductID: 2, Name: "Product 2", Desc: "This is Product 2", Price: 1.00, Inventory: 100, Rating: 3.0, RatingNum: 80, UnitSold: 100},
			},
			want: true,
		},
		{
			name: "RatingNum of p1 > p2",
			args: args{
				&models.Product{ProductID: 1, Name: "Product 1", Desc: "This is Product 1", Price: 1.00, Inventory: 100, Rating: 3.5, RatingNum: 80, UnitSold: 100},
				&models.Product{ProductID: 2, Name: "Product 2", Desc: "This is Product 2", Price: 1.00, Inventory: 100, Rating: 3.5, RatingNum: 60, UnitSold: 100},
			},
			want: true,
		},
		{
			name: "Inventory of p1 > p2",
			args: args{
				&models.Product{ProductID: 1, Name: "Product 1", Desc: "This is Product 1", Price: 1.00, Inventory: 200, Rating: 3.5, RatingNum: 80, UnitSold: 100},
				&models.Product{ProductID: 2, Name: "Product 2", Desc: "This is Product 2", Price: 1.00, Inventory: 100, Rating: 3.5, RatingNum: 80, UnitSold: 100},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortByPriceV(tt.args.p1, tt.args.p2); got != tt.want {
				t.Errorf("sortByPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_sortByPriceA(b *testing.B) {
	p1 := &models.Product{ProductID: 1, Name: "Product 1", Desc: "This is Product 1", Price: 1.00, Inventory: 100, Rating: 3.0, RatingNum: 80, UnitSold: 100}
	p2 := &models.Product{ProductID: 2, Name: "Product 2", Desc: "This is Product 2", Price: 1.00, Inventory: 80, Rating: 3.0, RatingNum: 80, UnitSold: 100}

	for i := 0; i < b.N; i++ {
		sortByPriceA(p1, p2)
	}
}

func Benchmark_sortByPriceV(b *testing.B) {
	p1 := &models.Product{ProductID: 1, Name: "Product 1", Desc: "This is Product 1", Price: 1.00, Inventory: 100, Rating: 3.0, RatingNum: 80, UnitSold: 100}
	p2 := &models.Product{ProductID: 2, Name: "Product 2", Desc: "This is Product 2", Price: 1.00, Inventory: 80, Rating: 3.0, RatingNum: 80, UnitSold: 100}

	for i := 0; i < b.N; i++ {
		sortByPriceV(p1, p2)
	}
}
