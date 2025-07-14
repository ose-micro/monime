package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type HttpClient struct {
	baseURL    string
	access     string
	space      string
	timeoutSec int
	client     *http.Client
	log        logger.Logger
	tracer     tracing.Tracer
}

func NewHttpClient(baseURL, access, space string, timeout int, log logger.Logger,
	tracer tracing.Tracer) *HttpClient {
	return &HttpClient{
		baseURL:    baseURL,
		access:     access,
		space:      space,
		timeoutSec: timeout,
		log:        log,
		tracer:     tracer,
		client:     &http.Client{Timeout: time.Duration(timeout) * time.Second},
	}
}

func (c *HttpClient) Do(ctx context.Context, method, path string, body any, out any) (*http.Response, error) {
	var buf io.Reader
	var idempotencyKey string

	ctx, span := c.tracer.Start(ctx, "HttpClient.Do", trace.WithAttributes(
		attribute.String("operation", "REQUEST"),
		attribute.String("method", method),
		attribute.String("path", path),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	if body != nil {
		// Try to extract idempotency key via reflection
		if v := reflect.ValueOf(body); v.Kind() == reflect.Ptr {
			if field := v.Elem().FieldByName("IdempotencyKey"); field.IsValid() && field.Kind() == reflect.String {
				idempotencyKey = field.String()
			}
		}

		// Marshal the body to JSON
		b, err := json.Marshal(body)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.log.Error("failed to marshal request body",
				zap.String("trace_id", traceId),
				zap.String("operation", "REQUEST"),
				zap.Error(err),
			)

			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}

	c.log.Info("starting HTTP request",
		zap.String("trace_id", traceId),
		zap.String("operation", "REQUEST"),
		zap.String("method", method),
		zap.String("path", path),
	)

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, buf)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create HTTP request",
			zap.String("trace_id", traceId),
			zap.String("operation", "REQUEST"),
			zap.Error(err),
		)

		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.access)
	req.Header.Set("Monime-Space-Id", c.space)

	if idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", idempotencyKey)
	}

	res, err := c.client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to execute HTTP request",
			zap.String("trace_id", traceId),
			zap.String("operation", "REQUEST"),
			zap.Error(err),
		)

		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, parseError(res)
	}

	if out != nil {
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(out); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.log.Error("failed to execute HTTP request",
				zap.String("trace_id", traceId),
				zap.String("operation", "REQUEST"),
				zap.Error(err),
			)

			return res, err
		}
	}

	c.log.Info("completed HTTP request",
		zap.String("trace_id", traceId),
		zap.String("operation", "REQUEST"),
		zap.String("method", method),
		zap.String("path", path),
	)

	return res, nil
}

func parseError(res *http.Response) error {
	switch res.StatusCode {
	case 400:
		return errors.New("bad request")
	case 401:
		return errors.New("unauthorized")
	case 403:
		return errors.New("forbidden")
	case 404:
		return errors.New("not found")
	case 409:
		return errors.New("conflict")
	case 500:
		return errors.New("internal server error")
	default:
		return errors.New("unexpected error")
	}
}
