package main

import (
	"context"
	"fmt"

	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/core/utils"
	"github.com/ose-micro/monime"
	"github.com/ose-micro/monime/internal/services/checkout"
)

func main() {
	config := monime.Config{
		BaseURL:    "https://api2.monime.io/v1",
		Access:     "mon_ojgk7W5vWiroQxJdnLuJJjos3yMAmXHXp0GNr1tsRlPNcK5Vlo1exAv5jRBLVQ2q",
		Space:      "spc-k6CqF5G1HDnyHKtquDxqPKRt4p7",
		Version:    "caph.2025-06-20",
		TimeoutSec: 30,
	}

	log, _ := logger.NewZap(logger.Config{
		Level:       "info",
		Environment: "development",
	})

	tracer, _ := tracing.NewOtel(tracing.Config{
		Endpoint:    "localhost:4317",
		ServiceName: "monime",
		SampleRatio: 1.0,
	}, log)

	mme := monime.New(config, log, tracer)

	ac, err := mme.Services().Checkout.Create(context.Background(), &checkout.CreateCommand{
		Name:        "Help Ishmael",
		Description: "On the night of July 10th, around 2am, fire tore through a compound inside Wellington. The whole area was dark, and a candle that was left unattended started the fire. Aunty Ramatu, a well-known akara seller, lost everything.",
		CancelURL:   "https://example.com/cancel",
		SuccessURL: "https://example.com/success",
		Reference: utils.GenerateUUID(),
		FinancialAccountID: "fac-k6CqF5HqTWmgr6DgfnMQphu818F",
		LineItems: []checkout.Item{
			{
				Type: "custom",
				Name: "Help Ishmael",
				Quantity: 1,
				Price: checkout.ItemPrice{
					Currency: "SLE",
					Value: 2 * 100,
				},
			},
		},
	})
	if err != nil {
		log.Error(err.Error())
	}

	log.Info(fmt.Sprintf("%+v", ac.Result.RedirectURL))
}
