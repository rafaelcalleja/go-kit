package domain

type ProductRepository interface {
	Save(*Product) error
	Of(id *ProductId) (*Product, error)
}
