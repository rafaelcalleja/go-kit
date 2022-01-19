package domain

type StockCreateService struct{}

func NewStockCreateService() StockCreateService {
	return StockCreateService{}
}

func (s StockCreateService) Create(id string, quantity int) error {
	return nil
}
