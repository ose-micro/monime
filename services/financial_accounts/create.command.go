package financial_accounts

import (
	"fmt"
	"strings"

	"github.com/ose-micro/cqrs"
)

// CreateCommand represents the command to create a financial account.
type CreateCommand struct {
	Name           string            `json:"name"`
	Currency       string            `json:"currency"`
	Reference      string            `json:"reference"`
	Metadata       map[string]string `json:"metadata"`
}

// CommandName implements cqrs.Command.
func (c CreateCommand) CommandName() string {
	return CREATED_COMMAND
}

// Validate implements cqrs.Command.
func (c CreateCommand) Validate() error {
	fields := make([]string, 0)

	if c.Currency == "" {
		fields = append(fields, "currency is required")
	}

	if c.Name == "" {
		fields = append(fields, "name is required")
	}

	if c.Reference == "" {
		fields = append(fields, "reference is required")
	}

	if len(fields) > 0 {
		return fmt.Errorf("%s", strings.Join(fields, ", "))
	}

	return nil
}

var _ cqrs.Command = CreateCommand{}
