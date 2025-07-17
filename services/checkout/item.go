package checkout

import (
	"fmt"
	"strings"

	"github.com/ose-micro/cqrs"
)

// Item represents a line item for checkout.
type ItemPrice struct {
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}

type Item struct {
	Type        string   `json:"type"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Quantity    int      `json:"quantity"`
	Reference   string   `json:"reference"`
	Price ItemPrice `json:"price"`
}

// CommandName implements cqrs.Command.
func (i Item) CommandName() string {
	return CREATED_COMMAND
}

// Validate implements cqrs.Command.
func (i Item) Validate() error {
	var fields []string

	if i.Type == "" {
		fields = append(fields, "type is required")
	}

	if i.ID == "" {
		fields = append(fields, "id is required")
	}

	if i.Name == "" {
		fields = append(fields, "name is required")
	}

	if i.Quantity <= 0 {
		fields = append(fields, "quantity must be greater than 0")
	}

	if i.Reference == "" {
		fields = append(fields, "reference is required")
	}

	if i.Price.Currency == "" {
		fields = append(fields, "price.currency is required")
	}

	if i.Price.Value <= 0 {
		fields = append(fields, "price.value must be greater than 0")
	}

	if len(fields) > 0 {
		return fmt.Errorf("%s", strings.Join(fields, ", "))
	}

	return nil
}

var _ cqrs.Command = Item{}
