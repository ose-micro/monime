package financial_accounts

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/monime/internal"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type financialAccountService struct {
	client *internal.HttpClient
	log    logger.Logger
	tracer tracing.Tracer
}

// Create implements Service.
func (f *financialAccountService) Create(ctx context.Context, command CreateCommand) (*Domain, error) {
	ctx, span := f.tracer.Start(ctx, "app.financial_account.create.command.handler", trace.WithAttributes(
		attribute.String("operation", "CREATE"),
		attribute.String("payload", fmt.Sprintf("%v", command)),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	// validate command payload
	if err := command.Validate(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		f.log.Error("validation process failed",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return nil, err
	}

	res, err := f.client.Do(ctx, "POST", "/v1/financial-accounts", command, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		f.log.Error("failed to create financial account",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return nil, err
	}
	if res.StatusCode >= 400 {
		err := fmt.Errorf("failed to create financial account, status code: %d", res.StatusCode)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		f.log.Error("failed to create financial account",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)
	}

	var data Domain
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		f.log.Error("failed to decode response body",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return nil, err
	}

	f.log.Info("financial account created successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "CREATE"),
		zap.String("payload", fmt.Sprintf("%v", data)))

	return &data, nil
}

// Delete implements Service.
func (f *financialAccountService) Delete(reference string) error {
	panic("unimplemented")
}

// Get implements Service.
func (f *financialAccountService) Get(reference string) (*Domain, error) {
	panic("unimplemented")
}

// List implements Service.
func (f *financialAccountService) List(ctx context.Context) ([]Domain, error) {
	ctx, span := f.tracer.Start(ctx, "app.financial_account.create.command.handler", trace.WithAttributes(
		attribute.String("operation", "CREATE"),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	res, err := f.client.Do(ctx, "GET", "/financial-accounts", nil, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		f.log.Error("failed to get financial accounts",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return nil, err
	}

	var data []Domain
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		f.log.Error("failed to decode response body",
			zap.String("trace_id", traceId),
			zap.String("operation", "CREATE"),
			zap.Error(err),
		)

		return nil, err
	}

	f.log.Info("financial account created successfully",
		zap.String("trace_id", traceId),
		zap.String("operation", "CREATE"),
		zap.String("payload", fmt.Sprintf("%v", data)))

	return data, nil
}

// Update implements Service.
func (f *financialAccountService) Update(reference string, account Domain) (*Domain, error) {
	panic("unimplemented")
}

func NewService(client *internal.HttpClient, log logger.Logger, tracer tracing.Tracer) Service {
	return &financialAccountService{
		client: client,
		log: log,
		tracer: tracer,
	}

}
