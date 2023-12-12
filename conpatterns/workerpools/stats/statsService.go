package stats

import (
	"fmt"
	"time"

	"github.com/Saddham/learning-go-concurrency/models"
)

const WorkerCount = 3

type statsService struct {
	result    models.Result
	processed <-chan models.Order
	done      <-chan struct{}
	pStats    chan models.Statistics
}

type StatsService interface {
	GetStats() models.Statistics
}

func NewStatsService(processed <-chan models.Order, done <-chan struct{}) StatsService {
	s := statsService{
		result:    models.NewResult(),
		processed: processed,
		done:      done,
		pStats:    make(chan models.Statistics, WorkerCount),
	}

	for i := 0; i < WorkerCount; i++ {
		go s.processStats()
	}

	go s.reconcile()

	return &s
}

// processStats is the overall processing method that listens to incoming orders
func (s *statsService) processStats() {
	fmt.Println("Stats processing started")

	for {
		select {
		case order := <-s.processed:
			stats := s.processOrder(order)
			s.pStats <- stats
		case <-s.done:
			fmt.Println("Stats processing stopped")
			return
		}
	}
}

// reconcile is a helper method which saves stats object back into the statisticsService
func (s *statsService) reconcile() {
	fmt.Println("Reconcile started")

	for {
		select {
		case p := <-s.pStats:
			s.result.Combine(p)
		case <-s.done:
			fmt.Println("Reconcile stopped")
			return
		}
	}
}

func (s *statsService) processOrder(order models.Order) models.Statistics {
	// Simulate processing as a costly operation
	time.Sleep(1 * time.Second)

	if order.Status == models.OrderStatus_Completed {
		return models.Statistics{
			CompletedOrders: 1,
			Revenue:         *order.Total,
		}
	}

	// Otherwise the order is rejected
	return models.Statistics{
		RejectedOrders: 1,
	}
}

func (s *statsService) GetStats() models.Statistics {
	return s.result.Get()
}
