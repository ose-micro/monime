package monime

import (
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/monime/rest"
	"github.com/ose-micro/monime/services"
)

type Monime struct {
	httpClient *rest.Client
	log        logger.Logger
	tracer     tracing.Tracer
	services   *services.Service
}

func (m Monime) Services() *services.Service {
	return m.services
}

func New(conf Config, log logger.Logger, tracer tracing.Tracer) *Monime {
	client := rest.New(conf.BaseURL, conf.Access, conf.Space, conf.Version, conf.TimeoutSec, log, tracer)
	svc := services.NewService(client, log, tracer)

	return &Monime{
		httpClient: client,
		log:        log,
		tracer:     tracer,
		services:   svc,
	}
}
