package services

import (
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/monime/rest"
	"github.com/ose-micro/monime/services/checkout"
	"github.com/ose-micro/monime/services/financial_accounts"
)

type Service struct {
	FinancialAccount financial_accounts.Service
	Checkout         checkout.Service
}

func NewService(client *rest.Client, log logger.Logger, tracer tracing.Tracer) *Service {
	return &Service{
		FinancialAccount: financial_accounts.NewService(client, log, tracer),
		Checkout:         checkout.NewService(client, log, tracer),
	}
}
