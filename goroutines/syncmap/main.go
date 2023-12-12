package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/Saddham/learning-go-concurrency/models"
)

func main() {
	pdb := models.NewProductDB()

	var wg sync.WaitGroup
	wg.Add(10)

	var lock sync.Mutex

	productProducer(pdb)

	// Create 10 consumers who concurrently uses and modifies same product
	// As each consumer decrements it only once, product quantity should
	// never go below zero
	for i := 0; i < 10; i++ {
		go productConsumer(&wg, &lock, pdb, 1)
	}

	wg.Wait()

	product, err := pdb.Find(strconv.Itoa(1))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Product quantity: %s\n", strconv.Itoa(product.Quantity))
}

func productProducer(odb *models.ProductDB) {
	product := &models.Product{
		ID:       strconv.Itoa(1),
		Quantity: 10,
	}

	odb.Upsert(*product)

}

func productConsumer(wg *sync.WaitGroup, lock *sync.Mutex, pdb *models.ProductDB, productId int) {
	lock.Lock()
	defer lock.Unlock()

	product, err := pdb.Find(strconv.Itoa(productId))
	if err != nil {
		panic(err)
	}

	product.Quantity--
	pdb.Upsert(product)

	wg.Done()
}
