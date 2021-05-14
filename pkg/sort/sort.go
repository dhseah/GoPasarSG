package sort

import "ProjectGoLive/pkg/models"

// sortByPriceA ranks two Product according to Price
// from lowest to highest. sortByPriceA compares
// Price > UnitSold > Rating > RatingNum > Inventory
func sortByPriceA(p1, p2 *models.Product) bool {
	if p1.Price == p2.Price {
		if p1.UnitSold == p2.UnitSold {
			if p1.Rating == p2.Rating {
				if p1.RatingNum == p2.RatingNum {
					return p1.Inventory > p2.Inventory
				}
				return p1.RatingNum > p2.RatingNum
			}
			return p1.Rating > p2.Rating
		}
		return p1.UnitSold > p2.UnitSold
	}
	return p1.Price < p2.Price
}

// sortByPriceD ranks two Product according to Price
// from highest to lowest. sortByPriceA compares
// Price > UnitSold > Rating > RatingNum > Inventory
func sortByPriceD(p1, p2 *models.Product) bool {
	if p1.Price == p2.Price {
		if p1.UnitSold == p2.UnitSold {
			if p1.Rating == p2.Rating {
				if p1.RatingNum == p2.RatingNum {
					return p1.Inventory > p2.Inventory
				}
				return p1.RatingNum > p2.RatingNum
			}
			return p1.Rating > p2.Rating
		}
		return p1.UnitSold > p2.UnitSold
	}
	return p1.Price > p2.Price
}

// sortByPriceV is an experiment to reduce the number
// of comparisons that are performed. sortByPriceV
// introduces inaccuracies to the sorting results
// compared to the previous two sorting logics.
func sortByPriceV(p1, p2 *models.Product) bool {
	p1Vector := (float64(p1.UnitSold) + p1.Rating + float64(p1.RatingNum) + float64(p1.Inventory)) / p1.Price
	p2Vector := (float64(p2.UnitSold) + p2.Rating + float64(p2.RatingNum) + float64(p2.Inventory)) / p2.Price
	return p1Vector > p2Vector
}

// sortByPop ranks two Product according to UnitSold
// from highest to lowest. sortByPop compares
// UnitSold > Price > Rating > RatingNum > Inventory
func sortByPop(p1, p2 *models.Product) bool {
	if p1.UnitSold == p2.UnitSold {
		if p1.Price == p2.Price {
			if p1.Rating == p2.Rating {
				if p1.RatingNum == p2.RatingNum {
					return p1.Inventory > p2.Inventory
				}
				return p1.RatingNum > p2.RatingNum
			}
			return p1.Rating > p2.Rating
		}
		return p1.Price < p2.Price
	}
	return p1.UnitSold > p2.UnitSold
}

// sortByRating ranks two Product according to Rating
// from highest to lowest. sortByRating compares
// Rating > RatingNum > Price > UnitSold > Inventory
func sortByRatings(p1, p2 *models.Product) bool {
	if p1.Rating == p2.Rating {
		if p1.RatingNum == p2.RatingNum {
			if p1.Price == p2.Price {
				if p1.UnitSold == p2.UnitSold {
					return p1.Inventory > p2.Inventory
				}
				return p1.UnitSold > p2.UnitSold
			}
			return p1.Price < p2.Price
		}
		return p1.RatingNum > p2.RatingNum
	}
	return p1.Rating > p2.Rating
}
