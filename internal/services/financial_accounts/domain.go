package financial_accounts

import "context"

type Balance struct {
	Available float64 `json:"available"`
}

type Domain struct {
	Id        string  `json:"id"`
	Name      float64 `json:"name"`
	Currency  string  `json:"currency"`
	Reference string  `json:"reference"`
	Balance   *Balance `json:"balance"`
	CreatedAt string  `json:"createdTime"`
	UpdatedAt string  `json:"updatedTime"`
}

const (
	CREATED_COMMAND string = "financial_accounts.create.command"
	UPDATED_COMMAND string = "financial_accounts.update.command"
)

type Service interface {
	Create(ctx context.Context, account CreateCommand) (*Domain, error)
	Get(reference string) (*Domain, error)
	Update(reference string, account Domain) (*Domain, error)
	Delete(reference string) error
	List(ctx context.Context) ([]Domain, error)
}
