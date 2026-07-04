package langfuse

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/namnv2496/post-service/configs"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// ensureBucketFromEnv will try to create the S3 bucket specified by environment
// variables `LANGFUSE_S3_EVENT_UPLOAD_ENDPOINT`, `LANGFUSE_S3_EVENT_UPLOAD_ACCESS_KEY_ID`,
// `LANGFUSE_S3_EVENT_UPLOAD_SECRET_ACCESS_KEY`, and `LANGFUSE_S3_EVENT_UPLOAD_BUCKET`.
// This is a best-effort helper for local setups (MinIO).
func ensureBucketFromEnv(ctx context.Context) error {
	endpoint := os.Getenv("LANGFUSE_S3_EVENT_UPLOAD_ENDPOINT")
	accessKey := os.Getenv("LANGFUSE_S3_EVENT_UPLOAD_ACCESS_KEY_ID")
	secretKey := os.Getenv("LANGFUSE_S3_EVENT_UPLOAD_SECRET_ACCESS_KEY")
	bucket := os.Getenv("LANGFUSE_S3_EVENT_UPLOAD_BUCKET")
	forcePath := os.Getenv("LANGFUSE_S3_EVENT_UPLOAD_FORCE_PATH_STYLE")

	if endpoint == "" || accessKey == "" || secretKey == "" || bucket == "" {
		return nil // nothing to do
	}

	// lazy import of minio to avoid hard failure when not used in other envs
	// create client
	secure := strings.HasPrefix(endpoint, "https")

	// Use reflection-like local implementation with real minio library
	// to avoid compile errors; actually import minio packages below.
	_ = secure

	return tryEnsureBucket(ctx, endpoint, accessKey, secretKey, bucket, forcePath == "true")
}

// tryEnsureBucket uses the MinIO client to ensure the bucket exists.
func tryEnsureBucket(ctx context.Context, endpoint, accessKey, secretKey, bucket string, forcePath bool) error {
	parsed, err := url.Parse(endpoint)
	host := endpoint
	secure := false
	if err == nil && parsed.Host != "" {
		host = parsed.Host
		secure = parsed.Scheme == "https"
	}

	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	}
	if forcePath {
		opts.BucketLookup = minio.BucketLookupPath
	}

	mc, err := minio.New(host, opts)
	if err != nil {
		return fmt.Errorf("minio client init: %w", err)
	}

	exists, err := mc.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("bucket exists check: %w", err)
	}
	if exists {
		fmt.Printf("langfuse: S3 bucket %s already exists\n", bucket)
		return nil
	}

	// Try to create the bucket
	if err := mc.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: "us-east-1"}); err != nil {
		return fmt.Errorf("make bucket: %w", err)
	}
	fmt.Printf("langfuse: created S3 bucket %s\n", bucket)
	return nil
}

type Client struct {
	cfg         configs.CommentRule
	traceClient *sdktrace.TracerProvider
	httpClient  *http.Client
}

func buildEndpointURL(raw string) string {
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
			raw = "http://" + raw
		}
		parsed, err = url.Parse(raw)
		if err != nil {
			return raw
		}
	}
	path := strings.TrimSuffix(parsed.Path, "/")
	if path == "" || path == "/" {
		parsed.Path = "/api/public/otel"
	} else if strings.HasSuffix(path, "/api/public/otel") {
		parsed.Path = path
	} else {
		parsed.Path = path + "/api/public/otel"
	}
	return parsed.String()
}

func NewClient(cfg configs.CommentRule) *Client {
	if !cfg.Langfuse.Enabled || cfg.Langfuse.Endpoint == "" || cfg.Langfuse.PublicKey == "" || cfg.Langfuse.SecretKey == "" || cfg.LangfuseTrace.Endpoint == "" || cfg.LangfuseTrace.PublicKey == "" || cfg.LangfuseTrace.SecretKey == "" {
		fmt.Printf("langfuse: tracing disabled or misconfigured\n")
		return nil
	}

	auth := base64.StdEncoding.EncodeToString(
		[]byte(cfg.LangfuseTrace.PublicKey + ":" + cfg.LangfuseTrace.SecretKey),
	)

	// Normalize base host (scheme + host)
	baseURL := cfg.LangfuseTrace.Endpoint
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "http://" + baseURL
	}

	// Probe common OTLP ingestion paths and pick the first non-404 endpoint.
	candidatePaths := []string{"/api/public/otel", "/api/public/otel/v1/traces", "/api/public/otel/v1", "/v1/traces", "/api/otel/v1/traces", "/api/otel"}
	endpointURL := probeOTLPEndpoint(baseURL, candidatePaths)
	// logging transport used for exporter and HTTP clients (redacts Authorization header)
	baseTransport := http.DefaultTransport
	loggingTransport := &loggingRoundTripper{base: baseTransport, verbose: cfg.LangfuseTrace.VerboseDump}

	// HTTP client used by the OTLP exporter (do NOT wrap with otelhttp transport to avoid exporter recursion)
	exporterHTTPClient := &http.Client{Transport: loggingTransport, Timeout: 10 * time.Second}

	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpointURL(endpointURL),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization":                "Basic " + auth,
			"x-langfuse-ingestion-version": "4",
		}),
		otlptracehttp.WithTimeout(10 * time.Second),
		otlptracehttp.WithHTTPClient(exporterHTTPClient),
	}

	exporter, err := otlptracehttp.New(context.Background(), opts...)
	if err != nil {
		fmt.Printf("langfuse: failed to create OTLP exporter: %v\n", err)
		return nil
	}
	res := resource.NewWithAttributes(
		"",
		attribute.String("service.name", "post-service"),
		attribute.String("langfuse.trace_name", cfg.LangfuseTrace.TraceName),
	)

	traceClient := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
	)

	// instrumented HTTP client for Langfuse public API calls (prompts)
	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(loggingTransport),
		Timeout:   10 * time.Second,
	}

	otel.SetTracerProvider(traceClient)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	fmt.Printf("langfuse: OTLP exporter initialized for endpoint=%s\n", endpointURL)

	// Best-effort: if S3 upload env vars are present (local MinIO), ensure the bucket exists.
	if err := ensureBucketFromEnv(context.Background()); err != nil {
		fmt.Printf("langfuse: ensure bucket error: %v\n", err)
	}

	return &Client{
		cfg:         cfg,
		traceClient: traceClient,
		httpClient:  httpClient,
	}
}

// probeOTLPEndpoint tries a list of candidate paths on the baseURL and returns
// the first URL that does not return 404. If all candidates return 404, it
// returns baseURL + "/api/public/otel" as a fallback.
func probeOTLPEndpoint(baseURL string, paths []string) string {
	client := &http.Client{Timeout: 5 * time.Second}
	for _, p := range paths {
		u := strings.TrimSuffix(baseURL, "/") + p
		req, err := http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			// network error — skip
			continue
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			fmt.Printf("langfuse: probed OTLP endpoint %s -> status=%d\n", u, resp.StatusCode)
			return u
		}
	}
	fallback := strings.TrimSuffix(baseURL, "/") + "/api/public/otel"
	fmt.Printf("langfuse: no OTLP probe succeeded, falling back to %s\n", fallback)
	return fallback
}

// loggingRoundTripper logs request and response details for debugging.
// It redacts Authorization headers to avoid printing secrets.
type loggingRoundTripper struct {
	base    http.RoundTripper
	verbose bool
}

func (l *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// copy request body for logging
	var reqBody []byte
	if req.Body != nil {
		b, _ := io.ReadAll(req.Body)
		reqBody = b
		req.Body = io.NopCloser(bytes.NewReader(b))
	}

	// redact headers copy
	hdrs := map[string][]string{}
	for k, v := range req.Header {
		if strings.ToLower(k) == "authorization" {
			hdrs[k] = []string{"REDACTED"}
		} else {
			hdrs[k] = v
		}
	}

	// short human-readable preview
	fmt.Printf("[langfuse-debug] HTTP Request: %s %s headers=%v body=%s\n", req.Method, req.URL.String(), hdrs, truncate(reqBody, 1024))
	if l.verbose && len(reqBody) > 0 {
		fmt.Printf("[langfuse-debug] HTTP Request BODY (base64, len=%d): %s\n", len(reqBody), base64.StdEncoding.EncodeToString(reqBody))
	}

	resp, err := l.base.RoundTrip(req)
	if err != nil {
		fmt.Printf("[langfuse-debug] HTTP Request error: %v\n", err)
		return resp, err
	}

	// copy response body for logging
	var respBody []byte
	if resp.Body != nil {
		b, _ := io.ReadAll(resp.Body)
		respBody = b
		resp.Body = io.NopCloser(bytes.NewReader(b))
	}
	fmt.Printf("[langfuse-debug] HTTP Response: %s %s status=%d body=%s\n", req.Method, req.URL.String(), resp.StatusCode, truncate(respBody, 4096))
	if l.verbose && len(respBody) > 0 {
		fmt.Printf("[langfuse-debug] HTTP Response BODY (base64, len=%d): %s\n", len(respBody), base64.StdEncoding.EncodeToString(respBody))
	}
	return resp, nil
}

func truncate(b []byte, n int) string {
	if len(b) == 0 {
		return ""
	}
	s := string(b)
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}

func (c *Client) IsEnabled() bool {
	return c.cfg.Langfuse.Enabled && c.cfg.Langfuse.Endpoint != "" && c.cfg.Langfuse.PublicKey != "" && c.cfg.Langfuse.SecretKey != ""
}

func (c *Client) GetPrompt(ctx context.Context, promptName string) (string, error) {
	if !c.IsEnabled() {
		return "", fmt.Errorf("langfuse is disabled")
	}
	reqURL := fmt.Sprintf("%s/api/public/v2/prompts/%s", c.cfg.Langfuse.Endpoint, url.PathEscape(promptName))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set(
		"Authorization",
		"Basic "+base64.StdEncoding.EncodeToString(
			[]byte(c.cfg.Langfuse.PublicKey+":"+c.cfg.Langfuse.SecretKey),
		),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"get prompt failed: status=%d body=%s",
			resp.StatusCode,
			string(body),
		)
	}

	var result struct {
		Name    string `json:"name"`
		Version int    `json:"version"`
		Prompt  string `json:"prompt"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.Prompt, nil
}

// GetPromptByProjectAndName fetches a prompt by project and name from Langfuse public API.
func (c *Client) GetPromptByProjectAndName(ctx context.Context, project, name string) (string, error) {
	if !c.IsEnabled() {
		return "", fmt.Errorf("langfuse is disabled")
	}
	reqURL := fmt.Sprintf("%s/api/public/v2/prompts?project=%s&name=%s", c.cfg.Langfuse.Endpoint, url.QueryEscape(project), url.QueryEscape(name))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set(
		"Authorization",
		"Basic "+base64.StdEncoding.EncodeToString(
			[]byte(c.cfg.Langfuse.PublicKey+":"+c.cfg.Langfuse.SecretKey),
		),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get prompt failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	// Langfuse may return either a raw array or a wrapped object { data: [...] }
	var wrapped struct {
		Data []struct {
			Name    string `json:"name"`
			Version int    `json:"version"`
			Prompt  string `json:"prompt"`
			Label   string `json:"label"`
			Project string `json:"project"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &wrapped); err == nil && len(wrapped.Data) > 0 {
		if wrapped.Data[0].Prompt != "" {
			return wrapped.Data[0].Prompt, nil
		}
		if wrapped.Data[0].Name != "" {
			return c.GetPrompt(ctx, wrapped.Data[0].Name)
		}
	}

	// fallback to raw array
	var results []struct {
		Name    string `json:"name"`
		Version int    `json:"version"`
		Prompt  string `json:"prompt"`
		Label   string `json:"label"`
		Project string `json:"project"`
	}
	if err := json.Unmarshal(body, &results); err != nil {
		return "", err
	}
	if len(results) == 0 {
		// debug: print raw response body so callers can inspect why Langfuse returned no data
		fmt.Printf("langfuse: no prompt found for project=%s name=%s; raw response=%s\n", project, name, string(body))
		return "", fmt.Errorf("no prompt found for project=%s name=%s", project, name)
	}
	if results[0].Prompt != "" {
		return results[0].Prompt, nil
	}
	if results[0].Name != "" {
		return c.GetPrompt(ctx, results[0].Name)
	}

	fmt.Printf("langfuse: no prompt found for project=%s name=%s; raw response=%s\n", project, name, string(body))
	return "", fmt.Errorf("no prompt found for project=%s name=%s", project, name)
}

func (c *Client) Trace(ctx context.Context, traceName string, input, output map[string]any) error {
	if c == nil || c.traceClient == nil {
		fmt.Printf("langfuse: trace skipped because tracer is not initialized\n")
		return nil
	}

	tracer := c.traceClient.Tracer("langfuse-client")
	ctx, span := tracer.Start(ctx, traceName, oteltrace.WithSpanKind(oteltrace.SpanKindInternal))
	defer span.End()

	var attrs []attribute.KeyValue
	attrs = append(attrs, attribute.String("langfuse.trace.name", traceName))

	if len(input) > 0 {
		inputJSON, err := json.Marshal(input)
		if err == nil {
			attrs = append(attrs, attribute.String("langfuse.trace.input", string(inputJSON)))
		} else {
			for k, v := range input {
				attrs = append(attrs, attribute.String(fmt.Sprintf("langfuse.trace.metadata.%s", k), fmt.Sprint(v)))
			}
		}
	}

	if len(output) > 0 {
		outputJSON, err := json.Marshal(output)
		if err == nil {
			attrs = append(attrs, attribute.String("langfuse.trace.output", string(outputJSON)))
		} else {
			attrs = append(attrs, attribute.String("langfuse.trace.output", fmt.Sprint(output)))
		}
	}

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.SetStatus(codes.Ok, "ok")
	fmt.Printf("langfuse: trace emitted name=%s input=%d output=%d\n", traceName, len(input), len(output))

	return nil
}

// Shutdown flushes and stops the tracer provider.
func (c *Client) Shutdown(ctx context.Context) error {
	if c == nil || c.traceClient == nil {
		return nil
	}
	return c.traceClient.Shutdown(ctx)
}

// EmitTrace is a convenience wrapper to emit a trace using the Langfuse client.
// It is safe to call with a nil client.
func EmitTrace(client *Client, ctx context.Context, traceName string, input, output map[string]any) error {
	if client == nil || client.traceClient == nil {
		return nil
	}
	tracer := client.traceClient.Tracer("langfuse-client")
	ctx, span := tracer.Start(ctx, traceName, oteltrace.WithSpanKind(oteltrace.SpanKindInternal))
	defer span.End()

	var attrs []attribute.KeyValue
	attrs = append(attrs, attribute.String("langfuse.trace.name", traceName))

	if len(input) > 0 {
		inputJSON, err := json.Marshal(input)
		if err == nil {
			attrs = append(attrs, attribute.String("langfuse.trace.input", string(inputJSON)))
		} else {
			for k, v := range input {
				attrs = append(attrs, attribute.String(fmt.Sprintf("langfuse.trace.metadata.%s", k), fmt.Sprint(v)))
			}
		}
	}

	if len(output) > 0 {
		outputJSON, err := json.Marshal(output)
		if err == nil {
			attrs = append(attrs, attribute.String("langfuse.trace.output", string(outputJSON)))
		} else {
			attrs = append(attrs, attribute.String("langfuse.trace.output", fmt.Sprint(output)))
		}
	}

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.SetStatus(codes.Ok, "ok")
	fmt.Printf("langfuse: trace emitted name=%s input=%d output=%d\n", traceName, len(input), len(output))
	return nil
}

// DoAndTrace performs an outgoing HTTP request using the client's instrumented HTTP client,
// starts a client span, injects propagation headers, and emits a Langfuse trace with the
// provided traceName, input and the HTTP response as output (body and status).
func (c *Client) DoAndTrace(ctx context.Context, traceName string, req *http.Request, input map[string]any) (int, []byte, error) {
	if c == nil || c.traceClient == nil {
		// fallback: perform a normal request without tracing
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return 0, nil, err
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, b, nil
	}

	tracer := c.traceClient.Tracer("langfuse-client")
	ctx, span := tracer.Start(ctx, traceName, oteltrace.WithSpanKind(oteltrace.SpanKindClient))
	defer span.End()

	// inject propagation headers
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	// perform request using the client's httpClient (instrumented)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		// emit trace with error
		_ = EmitTrace(c, ctx, traceName, input, map[string]any{"error": err.Error()})
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// set span attributes
	span.SetAttributes(attribute.String("http.method", req.Method), attribute.String("http.url", req.URL.String()), attribute.Int("http.status_code", resp.StatusCode))

	// emit a Langfuse trace summarizing the request/response
	out := map[string]any{"http_status": resp.StatusCode, "http_body": truncate([]byte(body), 1024)}
	_ = EmitTrace(c, ctx, traceName, input, out)

	return resp.StatusCode, body, nil
}
