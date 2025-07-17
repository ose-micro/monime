package financial_accounts

import (
	"context"

	"github.com/ose-micro/monime/internal/common"
)

type Available struct {
	Currency string  `json:"currency"`
	Value    float32 `json:"value"`
}

type Balance struct {
	Available Available `json:"available"`
}

type Domain struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Currency  string   `json:"currency"`
	Reference string   `json:"reference"`
	Balance   *Balance `json:"balance"`
	CreatedAt string   `json:"createTime"`
	UpdatedAt string   `json:"updateTime"`
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
