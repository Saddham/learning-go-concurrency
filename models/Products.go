package models

import (
	"fmt"
	"sync"
)

type Product struct {
	ID       string
	Quantity int
}

type ProductDB struct {
	products sync.Map
}

func NewProductDB() *ProductDB {
	return &ProductDB{}
}

func (odb *ProductDB) Find(id string) (Product, error) {
	product, ok := odb.products.Load(id)
	if !ok {
		return Product{}, fmt.Errorf("No product found for product id %s", id)
	}

	return product.(Product), nil
}

func (odb *ProductDB) Upsert(product Product) {
	odb.products.Store(product.ID, product)
}
