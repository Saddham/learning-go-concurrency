package stats

import (
	"fmt"

	"github.com/Saddham/learning-go-concurrency/models"
)

type repo struct {
	orders    *models.OrderDB
	stats     StatsService
	incoming  chan models.Order
	processed chan models.Order
	done      chan struct{}
}

type Repo interface {
	CreateOrder(product *models.Product) (*models.Order, error)
	GetOrderStats() models.Statistics
	Close()
}

func NewRepo(orders *models.OrderDB, processed chan models.Order, done chan struct{}) Repo {
	stats := NewStatsService(processed, done)

	r := repo{
		orders:    orders,
		stats:     stats,
		incoming:  make(chan models.Order),
		processed: processed,
		done:      done,
	}

	go r.processOrders()

	return &r
}

func (r *repo) CreateOrder(product *models.Product) (*models.Order, error) {
	order := models.NewOrder(product)

	select {
	case r.incoming <- order:
		r.orders.Upsert(order)
		return &order, nil
	case <-r.done:
		return nil, fmt.Errorf("Orders app is closed, try again later")
	}
}

func (r *repo) processOrders() {
	fmt.Println("Order processing started!")

	for {
		select {
		case order := <-r.incoming:
			if *order.Total < 500 {
				order.Reject()
			} else {
				order.Complete()
			}

			r.orders.Upsert(order)
			fmt.Printf("Processing order %s completed\n", order.ID)
			r.processed <- order
		case <-r.done:
			fmt.Println("Order processing stopped!")
			return
		}
	}
}

func (r *repo) GetOrderStats() models.Statistics {
	return r.stats.GetStats()
}

func (r *repo) Close() {
	close(r.done)
}
