package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

type Client struct {
	baseURL    string
	access     string
	version    string
	space      string
	timeoutSec int
	client     *http.Client
	log        logger.Logger
	tracer     tracing.Tracer
}

func New(baseURL, access, space, version string, timeout int, log logger.Logger,
	tracer tracing.Tracer) *Client {
	return &Client{
		baseURL:    baseURL,
		access:     access,
		space:      space,
		version:    version,
		timeoutSec: timeout,
		log:        log,
		tracer:     tracer,
		client:     &http.Client{Timeout: time.Duration(timeout) * time.Second},
	}
}

func (c *Client) Get(ctx context.Context, path string, body any, unmarshal func([]byte) (any, error)) (any, error) {
	var BUF io.Reader
	METHOD := "GET"
	TOKEN := fmt.Sprintf("Bearer %s", c.access)

	ctx, span := c.tracer.Start(ctx, "HttpClient.GET", trace.WithAttributes(
		attribute.String("method", METHOD),
		attribute.String("path", path),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// Marshal body if present
	if body != nil {
		v := reflect.ValueOf(body)
		if v.Kind() == reflect.Ptr && v.IsNil() {
			c.log.Error("request body is nil",
				zap.String("trace_id", traceId),
				zap.String("method", METHOD),
				zap.String("path", path))
		}

		b, err := json.Marshal(body)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.log.Error("failed to marshal request body",
				zap.String("trace_id", traceId),
				zap.Error(err),
			)
			return nil, err
		}
		BUF = bytes.NewBuffer(b)
	}

	URL := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequestWithContext(ctx, METHOD, URL, BUF)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create HTTP request",
			zap.String("trace_id", traceId),
			zap.Error(err),
		)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", TOKEN)
	req.Header.Set("Monime-Version", c.version)
	req.Header.Set("Monime-Space-Id", c.space)

	c.log.Info("starting HTTP request",
		zap.String("trace_id", traceId),
		zap.String("method", METHOD),
		zap.String("path", path),
	)

	res, err := c.client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("HTTP request failed",
			zap.String("trace_id", traceId),
			zap.Error(err),
		)
		return nil, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to read response body",
			zap.String("trace_id", traceId),
			zap.Error(err),
		)
		return nil, err
	}

	if res.StatusCode >= 400 {
		err := fmt.Errorf("request to %s failed with status=%d, body=%s", path, res.StatusCode, string(bodyBytes))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("HTTP response error",
			zap.String("trace_id", traceId),
			zap.Int("status_code", res.StatusCode),
			zap.String("body", string(bodyBytes)),
		)
		return nil, err
	}

	// Use the provided unmarshal function to decode the response
	out, err := unmarshal(bodyBytes)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to decode response",
			zap.String("trace_id", traceId),
			zap.Error(err),
		)
		return nil, err
	}

	c.log.Info("completed HTTP request and decoded response",
		zap.String("trace_id", traceId),
		zap.String("method", METHOD),
		zap.String("path", path),
	)

	return &out, nil
}

func (c *Client) POST(ctx context.Context, path string, body any, headers map[string]string, unmarshal func([]byte) (any, error)) (any, error) {
	var buf io.Reader
	method := "POST"
	token := fmt.Sprintf("Bearer %s", c.access)

	ctx, span := c.tracer.Start(ctx, "HttpClient.POST", trace.WithAttributes(
		attribute.String("method", method),
		attribute.String("path", path),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	// Marshal body if present
	if body != nil {
		v := reflect.ValueOf(body)
		if v.Kind() == reflect.Ptr && v.IsNil() {
			c.log.Error("request body is nil",
				zap.String("trace_id", traceId),
				zap.String("method", method),
				zap.String("path", path))
		}

		b, err := json.Marshal(body)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.log.Error("failed to marshal request body",
				zap.String("trace_id", traceId),
				zap.String("method", method),
				zap.Error(err),
			)
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to create HTTP request",
			zap.String("trace_id", traceId),
			zap.String("method", method),
			zap.Error(err),
		)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	req.Header.Set("Monime-Version", c.version)
	req.Header.Set("Monime-Space-Id", c.space)

	for k, v := range headers {
		if v != "" {
			req.Header.Set(k, v)
		}
	}

	c.log.Info("starting HTTP request",
		zap.String("trace_id", traceId),
		zap.String("method", method),
		zap.String("path", path),
	)

	res, err := c.client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("HTTP request failed",
			zap.String("trace_id", traceId),
			zap.String("method", method),
			zap.Error(err),
		)
		return nil, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to read response body",
			zap.String("trace_id", traceId),
			zap.String("method", method),
			zap.Error(err),
		)
		return nil, err
	}

	if res.StatusCode >= 400 {
		err := fmt.Errorf("request to %s failed with status=%d, body=%s", path, res.StatusCode, string(bodyBytes))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("HTTP response error",
			zap.String("trace_id", traceId),
			zap.Int("status_code", res.StatusCode),
			zap.String("body", string(bodyBytes)),
		)
		return nil, err
	}

	// Use the provided unmarshal function to decode the response
	out, err := unmarshal(bodyBytes)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.log.Error("failed to decode response",
			zap.String("trace_id", traceId),
			zap.Error(err),
		)
		return nil, err
	}

	c.log.Info("completed HTTP request and decoded response",
		zap.String("trace_id", traceId),
		zap.String("method", method),
		zap.String("path", path),
	)

	return &out, nil
}

func (c *Client) PUT(ctx context.Context, path string, body any, headers map[string]string, unmarshal func([]byte) (any, error)) (any, error) {
	method := "PUT"
	token := fmt.Sprintf("Bearer %s", c.access)
	var buf io.Reader

	ctx, span := c.tracer.Start(ctx, "HttpClient.PUT", trace.WithAttributes(
		attribute.String("method", method),
		attribute.String("path", path),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	if body != nil {
		v := reflect.ValueOf(body)
		if v.Kind() == reflect.Ptr && v.IsNil() {
			c.log.Error("request body is nil",
				zap.String("trace_id", traceId),
				zap.String("method", method),
				zap.String("path", path))
		}

		b, err := json.Marshal(body)
		if err != nil {
			span.RecordError(err)
			c.log.Error("failed to marshal PUT body", zap.String("trace_id", traceId), zap.Error(err))
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		span.RecordError(err)
		c.log.Error("failed to create PUT request", zap.String("trace_id", traceId), zap.Error(err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	req.Header.Set("Monime-Version", c.version)
	req.Header.Set("Monime-Space-Id", c.space)
	for k, v := range headers {
		if v != "" {
			req.Header.Set(k, v)
		}
	}

	res, err := c.client.Do(req)
	if err != nil {
		span.RecordError(err)
		c.log.Error("PUT request failed", zap.String("trace_id", traceId), zap.Error(err))
		return nil, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		span.RecordError(err)
		c.log.Error("failed to read PUT response", zap.String("trace_id", traceId), zap.Error(err))
		return nil, err
	}

	if res.StatusCode >= 400 {
		err := fmt.Errorf("PUT request to %s failed: %d - %s", path, res.StatusCode, string(bodyBytes))
		span.RecordError(err)
		c.log.Error("PUT response error", zap.String("trace_id", traceId), zap.String("body", string(bodyBytes)))
		return nil, err
	}

	out, err := unmarshal(bodyBytes)
	if err != nil {
		span.RecordError(err)
		c.log.Error("PUT unmarshal failed", zap.String("trace_id", traceId), zap.Error(err))
		return nil, err
	}

	return &out, nil
}

func (c *Client) DELETE(ctx context.Context, path string, body any, headers map[string]string, unmarshal func([]byte) (any, error)) (any, error) {
	method := "DELETE"
	token := fmt.Sprintf("Bearer %s", c.access)
	var buf io.Reader

	ctx, span := c.tracer.Start(ctx, "HttpClient.DELETE", trace.WithAttributes(
		attribute.String("method", method),
		attribute.String("path", path),
	))
	defer span.End()

	traceId := trace.SpanContextFromContext(ctx).TraceID().String()

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			span.RecordError(err)
			c.log.Error("failed to marshal DELETE body", zap.String("trace_id", traceId), zap.Error(err))
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		span.RecordError(err)
		c.log.Error("failed to create DELETE request", zap.String("trace_id", traceId), zap.Error(err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	req.Header.Set("Monime-Version", c.version)
	req.Header.Set("Monime-Space-Id", c.space)
	for k, v := range headers {
		if v != "" {
			req.Header.Set(k, v)
		}
	}

	res, err := c.client.Do(req)
	if err != nil {
		span.RecordError(err)
		c.log.Error("DELETE request failed", zap.String("trace_id", traceId), zap.Error(err))
		return nil, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		span.RecordError(err)
		c.log.Error("failed to read DELETE response", zap.String("trace_id", traceId), zap.Error(err))
		return nil, err
	}

	if res.StatusCode >= 400 {
		err := fmt.Errorf("DELETE request to %s failed: %d - %s", path, res.StatusCode, string(bodyBytes))
		span.RecordError(err)
		c.log.Error("DELETE response error", zap.String("trace_id", traceId), zap.String("body", string(bodyBytes)))
		return nil, err
	}

	out, err := unmarshal(bodyBytes)
	if err != nil {
		span.RecordError(err)
		c.log.Error("DELETE unmarshal failed", zap.String("trace_id", traceId), zap.Error(err))
		return nil, err
	}

	return &out, nil
}
