package domain

type IProduct interface {
	SaveProduct(name string, description string, price float64, category string, imageURL string) error
	GetAll() ([]Product, error)
	Delete(id string) error
	Update(id string, name string, description string, price float64, category string, imageURL string) error
}

type Product struct {
	ID          int32   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
}

func NewProduct(name string, description string, price float64, category string, imageURL string) *Product {
	return &Product{
		Name:        name,
		Description: description,
		Price:       price,
		Category:    category,
		ImageURL:    imageURL,
	}
}

func (t *Product) SetPrice(price float64) {
	t.Price = price
}
