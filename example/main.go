package main

import (
	"context"
	"fmt"

	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/monime"
)

func main() {
	config := monime.Config{
		BaseURL: "https://api.monime.io/v2",
		Access:  "mon_ojgk7W5vWiroQxJdnLuJJjos3yMAmXHXp0GNr1tsRlPNcK5Vlo1exAv5jRBLVQ2q",
		Space:   "spc-k6CqF5G1HDnyHKtquDxqPKRt4p7",
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

	res, err := mme.Services().FinancialAccount.List(context.Background())
	if err != nil {
		log.Error(err.Error())
	}

	log.Info(fmt.Sprintf("%+v", res))
}
