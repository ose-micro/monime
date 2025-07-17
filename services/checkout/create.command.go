package checkout

import (
	"fmt"
	"strings"

	"github.com/ose-micro/cqrs"
)

// CreateCommand represents the command to create a financial account.
type CreateCommand struct {
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	CancelURL          string     `json:"cancelUrl"`
	SuccessURL         string     `json:"successUrl"`
	CallbackState      string     `json:"callbackState"`
	Reference          string     `json:"reference"`
	FinancialAccountID string     `json:"financialAccountId"`
	LineItems          []Item `json:"lineItems"`
}

// CommandName implements cqrs.Command.
func (c CreateCommand) CommandName() string {
	return CREATED_COMMAND
}

// Validate implements cqrs.Command.
func (c CreateCommand) Validate() error {
	var fields []string

	if c.Name == "" {
		fields = append(fields, "name is required")
	}

	if c.Reference == "" {
		fields = append(fields, "reference is required")
	}

	if c.FinancialAccountID == "" {
		fields = append(fields, "financialAccountId is required")
	}

	if c.SuccessURL == "" {
		fields = append(fields, "successUrl is required")
	}

	if c.CancelURL == "" {
		fields = append(fields, "cancelUrl is required")
	}

	if len(c.LineItems) == 0 {
		fields = append(fields, "at least one line item is required")
	} else {
		for i, item := range c.LineItems {
			if err := item.Validate(); err != nil {
				fields = append(fields, fmt.Sprintf("lineItems[%d]: %v", i, err))
			}
		}
	}

	if len(fields) > 0 {
		return fmt.Errorf("%s", strings.Join(fields, ", "))
	}

	return nil
}

var _ cqrs.Command = CreateCommand{}
