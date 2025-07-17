package financial_accounts

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/core/utils"
	"github.com/ose-micro/monime/common"
	"github.com/ose-micro/monime/rest"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type financialAccountService struct {
	client *rest.Client
	log    logger.Logger
	tracer tracing.Tracer
}

// Create implements Service.
func (f *financialAccountService) Create(ctx context.Context, command *CreateCommand) (*common.OneResponse[Domain], error) {
	var data common.OneResponse[Domain]

	ctx, span := f.tracer.Start(ctx, "app.financial_account.create.handler", trace.WithAttributes(
		attribute.String("operation", "CREATE"),
		attribute.String("payload", fmt.Sprintf("%+v", command))))

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	if _, err := f.client.POST(ctx, "/financial-accounts", command, map[string]string{
		"Idempotency-Key": utils.GenerateUUID(),
	}, func(b []byte) (any, error) {
		if err := json.Unmarshal(b, &data); err != nil {
			return nil, err
		}
		return data, nil
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		f.log.Error("failed to create financial account",
			zap.String("trace_id", traceId),
			zap.Error(err),
		)
		return nil, err
	}

	f.log.Info("financial account created",
		zap.String("trace_id", traceId),
		zap.String("payload", fmt.Sprintf("%+v", data)),
	)

	return &data, nil
}

// List implements Service.
func (f *financialAccountService) List(ctx context.Context) (*common.Response[Domain], error) {
	var data common.Response[Domain]

	ctx, span := f.tracer.Start(ctx, "app.financial_account.list.handler", trace.WithAttributes(
		attribute.String("operation", "LIST"),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	if _, err := f.client.Get(ctx, "/financial-accounts", nil, func(b []byte) (any, error) {
		if err := json.Unmarshal(b, &data); err != nil {
			return nil, err
		}
		return data, nil
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		f.log.Error("failed to list financial accounts",
			zap.String("trace_id", traceId),
			zap.Error(err),
		)
		return nil, err
	}

	f.log.Info("financial accounts fetched",
		zap.String("trace_id", traceId),
		zap.String("payload", fmt.Sprintf("%+v", data.Result)),
		zap.Int("count", data.Pagination.Count),
		zap.String("next", data.Pagination.Next),
	)

	return &data, nil
}

func (f *financialAccountService) Get(ctx context.Context, id string) (*common.OneResponse[Domain], error) {
	var data common.OneResponse[Domain]

	ctx, span := f.tracer.Start(context.Background(), "app.financial_account.get.handler", trace.WithAttributes(
		attribute.String("operation", "GET"),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	if _, err := f.client.Get(ctx, fmt.Sprintf("/financial-accounts/%s", id), nil, func(b []byte) (any, error) {
		if err := json.Unmarshal(b, &data); err != nil {
			return nil, err
		}
		return data, nil
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		f.log.Error("failed to get financial account",
			zap.String("trace_id", traceId),
			zap.String("payload", id),
			zap.Error(err),
		)
		return nil, err
	}

	f.log.Info("financial account fetched",
		zap.String("trace_id", traceId),
		zap.String("payload", fmt.Sprintf("%+v", data)),
	)

	return &data, nil
}

func (f *financialAccountService) Update(ctx context.Context, cmd *UpdateCommand) (*common.OneResponse[Domain], error) {
	var data common.OneResponse[Domain]

	ctx, span := f.tracer.Start(context.Background(), "app.financial_account.update.handler", trace.WithAttributes(
		attribute.String("operation", "UPDATE"),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	url := fmt.Sprintf("/financial-accounts/%s", cmd.Id)
	if _, err := f.client.PUT(ctx, url, cmd, map[string]string{
		"Idempotency-Key": utils.GenerateUUID(),
	}, func(b []byte) (any, error) {
		if err := json.Unmarshal(b, &data); err != nil {
			return nil, err
		}
		return data, nil
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		f.log.Error("failed to update financial account",
			zap.String("trace_id", traceId),
			zap.String("payload", fmt.Sprintf("%+v", cmd)),
			zap.Error(err),
		)
		return nil, err
	}

	f.log.Info("financial account updated",
		zap.String("trace_id", traceId),
		zap.String("payload", fmt.Sprintf("%+v", cmd)),
		zap.String("payload", fmt.Sprintf("%+v", data)),
	)

	return &data, nil
}

func NewService(client *rest.Client, log logger.Logger, tracer tracing.Tracer) Service {
	return &financialAccountService{
		client: client,
		log:    log,
		tracer: tracer,
	}
}
