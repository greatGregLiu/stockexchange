package stockexchange

import "fmt"

// Share that you own
type Share struct {
	// Symbol that identifies the stock
	Symbol string `json:"symbol"`
	// Name of the company
	Name string `json:"name"`
	// PaidPrice is the price that you bought
	PaidPrice float32 `json:"paidPrice"`
	// Quantity of shares
	Quantity int `json:"quantity"`
}

// Invoice to sell shares
type Invoice struct {
	// Symbol of the share
	Symbol string `json:"symbol"`
	// Price on the stock market
	Price float32 `json:"price"`
	// Quantity of the shares
	Quantity int `json:"quantity"`
}

// Portfolio is a stock trader that buys and sell shares
type Portfolio struct {
	// Balance that the trader has
	Balance float32 `json:"balance"`
	// Shares that the trader has
	Shares []*Share `json:"shares"`
}

// Buy performs a buy operation and adds the share to the portfolio
func (p *Portfolio) Buy(stock *Stock, quantity int) error {
	if quantity < 0 {
		return fmt.Errorf("The quantity cannot be negative number")
	}

	price := stock.AskPrice * float32(quantity)
	if price > p.Balance {
		return fmt.Errorf("Insufficient funds")
	}

	p.Balance -= price

	for _, share := range p.Shares {
		if share.Symbol == stock.Symbol {
			share.Quantity += quantity
			return nil
		}
	}

	share := &Share{
		Symbol:    stock.Symbol,
		Name:      stock.Name,
		PaidPrice: stock.AskPrice,
		Quantity:  quantity,
	}

	p.Shares = append(p.Shares, share)
	return nil
}

// Sell performs a sell operation and remove the share from the portfolio
func (p *Portfolio) Sell(invoice *Invoice) error {
	price := invoice.Price
	if price < 0 {
		return fmt.Errorf("The price cannot be negative number")
	}

	quantity := invoice.Quantity
	if quantity < 0 {
		return fmt.Errorf("The quantity cannot be negative number")
	}

	symbol := invoice.Symbol

	for index, share := range p.Shares {
		if share.Symbol == symbol {
			if quantity > share.Quantity {
				return fmt.Errorf("The desired quantity is greater than share quantity")
			}
			share.PaidPrice = price
			share.Quantity -= quantity
			if share.Quantity == 0 {
				p.Shares = deleteAt(p.Shares, index)
			}
			p.Balance += price * float32(quantity)
			return nil
		}
	}
	return fmt.Errorf("The desired share '%s' does not exist in this portfolio", symbol)
}

func deleteAt(shares []*Share, index int) []*Share {
	return append(shares[:index], shares[index+1:]...)
}
