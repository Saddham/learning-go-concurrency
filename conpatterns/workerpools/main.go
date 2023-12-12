package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Saddham/learning-go-concurrency/conpatterns/workerpools/stats"
	"github.com/Saddham/learning-go-concurrency/models"
)

const SIMULATION_CNT int = 10

func main() {
	processed := make(chan models.Order, stats.WorkerCount)
	done := make(chan struct{})

	orderDB := models.NewOrderDB()

	r := stats.NewRepo(orderDB, processed, done)
	defer r.Close()

	var wg sync.WaitGroup
	wg.Add(SIMULATION_CNT)

	for i := 1; i <= 10; i++ {
		createOrder(i, r, &wg)
	}

	wg.Wait()

	// Wait for stats service to finish collecting stats
	time.Sleep(10 * time.Second)

	orderStats := r.GetOrderStats()
	orderStatsStr, err := json.Marshal(orderStats)

	if err != nil {
		panic("Failed to marshall order stats")
	}

	fmt.Printf("Stats: %s\n", orderStatsStr)
}

func createOrder(i int, r stats.Repo, wg *sync.WaitGroup) {
	product := &models.Product{
		ID:       strconv.Itoa(i),
		Quantity: i * 10,
	}

	r.CreateOrder(product)

	wg.Done()
}
