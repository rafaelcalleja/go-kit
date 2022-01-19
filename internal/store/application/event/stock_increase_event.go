package event

import (
	"context"
	"errors"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/events"
	"github.com/rafaelcalleja/go-kit/internal/store/domain"
)

var (
	ErrStockEventUnhandled = errors.New("unexpected event")
)

type IncreaseStockOnProductCreated struct {
	stockCreateService domain.StockCreateService
}

func NewIncreaseStockOnProductCreated(stockCreateService domain.StockCreateService) IncreaseStockOnProductCreated {
	return IncreaseStockOnProductCreated{
		stockCreateService: stockCreateService,
	}
}

func (e IncreaseStockOnProductCreated) Handle(_ context.Context, evt events.Event) error {
	productCreatedEvent, ok := evt.(domain.ProductCreatedEvent)
	if !ok {
		return ErrStockEventUnhandled
	}

	return e.stockCreateService.Create(productCreatedEvent.ID(), 1)
}
