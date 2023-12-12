package models

import "math"

type Statistics struct {
	CompletedOrders int
	RejectedOrders  int
	Revenue         float64
}

func Combine(this, that Statistics) Statistics {
	return Statistics{
		CompletedOrders: this.CompletedOrders + that.CompletedOrders,
		RejectedOrders:  this.RejectedOrders + that.RejectedOrders,
		Revenue:         math.Round((this.Revenue+that.Revenue)*100) / 100,
	}
}
