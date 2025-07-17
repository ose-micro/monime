package checkout

import (
	"context"

	"github.com/ose-micro/monime/common"
)

type LineItems struct {
	Data []Item `json:"data"`
}

type PaymentOptions struct {
	Card   map[string]interface{} `json:"card"`
	Bank   map[string]interface{} `json:"bank"`
	Momo   map[string]interface{} `json:"momo"`
	Wallet map[string]interface{} `json:"wallet"`
}

type BrandingOptions struct {
	PrimaryColor string `json:"primaryColor"`
}

type Domain struct {
	ID                 string                 `json:"id"`
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	CancelURL          string                 `json:"cancelUrl"`
	SuccessURL         string                 `json:"successUrl"`
	CallbackState      string                 `json:"callbackState"`
	Reference          string                 `json:"reference"`
	RedirectURL        string                 `json:"redirectUrl"`
	FinancialAccountID string                 `json:"financialAccountId"`
	LineItems          LineItems              `json:"lineItems"`
	PaymentOptions     PaymentOptions         `json:"paymentOptions"`
	BrandingOptions    BrandingOptions        `json:"brandingOptions"`
	Metadata           map[string]interface{} `json:"metadata"`
}

const (
	CREATED_COMMAND string = "financial_accounts.create.command"
	UPDATED_COMMAND string = "financial_accounts.update.command"
)

type Service interface {
	Create(ctx context.Context, cmd *CreateCommand) (*common.OneResponse[Domain], error)
	Get(ctx context.Context, id string) (*common.OneResponse[Domain], error)
	Update(ctx context.Context, cmd *UpdateCommand) (*common.OneResponse[Domain], error)
	List(ctx context.Context) (*common.Response[Domain], error)
}
