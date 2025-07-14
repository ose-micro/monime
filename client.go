package monime

import (
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/monime/internal"
	"github.com/ose-micro/monime/internal/services"
)

type Monime struct {
	httpClient *internal.HttpClient
	log        logger.Logger
	tracer     tracing.Tracer
	services   *services.Service
}

func (m Monime) Services() *services.Service {
	return m.services
}

func New(conf Config, log logger.Logger, tracer tracing.Tracer) *Monime {
	client := internal.NewHttpClient(conf.BaseURL, conf.Access, conf.Space, conf.TimeoutSec, log, tracer)
	svc := services.NewService(client, log, tracer)

	return &Monime{
		httpClient: client,
		log:        log,
		tracer:     tracer,
		services:   svc,
	}
}
