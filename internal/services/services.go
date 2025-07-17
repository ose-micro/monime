package services

import (
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/monime/internal"
	"github.com/ose-micro/monime/internal/services/checkout"
	"github.com/ose-micro/monime/internal/services/financial_accounts"
)

type Service struct {
	FinancialAccount financial_accounts.Service
	Checkout checkout.Service
}

func NewService(client *internal.HttpClient, log logger.Logger, tracer tracing.Tracer) *Service {
	return &Service{
		FinancialAccount: financial_accounts.NewService(client, log, tracer),
		Checkout: checkout.NewService(client, log, tracer),
	}
}