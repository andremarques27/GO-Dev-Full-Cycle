package usercase

import "github.com/devfullcycle/go-intensivo-andre/internal/entity"

type OrderInput struct {
	ID    string  `json: "id"`
	Price float64 `json: "price"`
	Tax   float64 `json: "tax"`
}

type OrderOutput struct {
	ID         string
	Price      float64
	Tax        float64
	FinalPrice float64
}

type CalculateFinalPrice struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewCalculateFinalPrice(orderRepository entity.OrderRepositoryInterface) *CalculateFinalPrice {
	return &CalculateFinalPrice{
		OrderRepository: orderRepository,
	}
}

func (c *CalculateFinalPrice) Execute(imput OrderInput) (*OrderOutput, error) {
	order, err := entity.NewOrder(imput.ID, imput.Price, imput.Tax)
	if err != nil {
		return nil, err
	}
	err = order.CalculateFinalPrice()
	if err != nil {
		return nil, err
	}
	err = c.OrderRepository.Save(order)
	if err != nil {
		return nil, err
	}
	return &OrderOutput{
		ID:         order.ID,
		Price:		order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}, nil
}