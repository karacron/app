package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/authuser-dev/karacron/apps/core/src/module/inference/domain"
	"go.uber.org/zap"
)

type LocalExecutor struct {
	defaultModel string
	defaultN     int
	defaultTemp  float64
	timeout      time.Duration
	serverHost   string
	preferredPort int
	portRangeMin  int
	portRangeMax  int
	serverURL    string
	httpClient   *http.Client
	log          *zap.Logger
	mu           sync.Mutex
	cmd          *exec.Cmd
	modelPath    string
	executable   string
}

func NewLocalExecutor(defaultModel string, defaultN int, defaultTemp float64, timeout time.Duration, serverHost string, preferredPort, portRangeMin, portRangeMax int, logger *zap.Logger) *LocalExecutor {
	host := strings.TrimSpace(serverHost)
	if host == "" {
		host = "127.0.0.1"
	}
	if preferredPort <= 0 {
		preferredPort = 32111
	}
	if portRangeMin <= 0 {
		portRangeMin = 32111
	}
	if portRangeMax <= 0 {
		portRangeMax = 32130
	}
	if portRangeMax < portRangeMin {
		portRangeMin, portRangeMax = portRangeMax, portRangeMin
	}
	serverURL := buildServerURL(host, preferredPort)

	return &LocalExecutor{
		defaultModel: defaultModel,
		defaultN:     defaultN,
		defaultTemp:  defaultTemp,
		timeout:      timeout,
		serverHost:   host,
		preferredPort: preferredPort,
		portRangeMin: portRangeMin,
		portRangeMax: portRangeMax,
		serverURL:    serverURL,
		httpClient:   &http.Client{Timeout: timeout},
		log:          logger,
	}
}

func (e *LocalExecutor) Run(ctx context.Context, executablePath string, req domain.InferenceRequest, modelPath string) (domain.InferenceResponse, error) {
	prompt := strings.TrimSpace(req.Prompt)
	if prompt == "" {
		return domain.InferenceResponse{}, domain.ErrEmptyPrompt
	}

	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = e.defaultN
	}
	if maxTokens <= 0 {
		maxTokens = 256
	}

	temp := req.Temperature
	if temp <= 0 {
		temp = e.defaultTemp
	}
	if temp <= 0 {
		temp = 0.7
	}

	if err := e.ensureServer(runCtxOrBackground(ctx), executablePath, modelPath); err != nil {
		return domain.InferenceResponse{}, err
	}

	runCtx := ctx
	if e.timeout > 0 {
		var cancel context.CancelFunc
		runCtx, cancel = context.WithTimeout(ctx, e.timeout)
		defer cancel()
	}

	start := time.Now()
	first, err := e.requestCompletion(runCtx, req.ModelName, prompt, maxTokens, temp)
	if err != nil {
		return domain.InferenceResponse{}, err
	}

	content := first.Content
	if content == "" && first.FinishReason == "length" && first.HasReasoning {
		retryPrompt := "Responde solo con la respuesta final, sin razonamiento interno ni explicacion del proceso.\n\n" + prompt
		retryMaxTokens := maxTokens * 2
		if retryMaxTokens < 512 {
			retryMaxTokens = 512
		}
		if retryMaxTokens > 1536 {
			retryMaxTokens = 1536
		}

		retryTemp := temp
		if retryTemp > 0.3 {
			retryTemp = 0.3
		}

		second, retryErr := e.requestCompletion(runCtx, req.ModelName, retryPrompt, retryMaxTokens, retryTemp)
		if retryErr == nil && second.Content != "" {
			content = second.Content
		} else if retryErr != nil {
			return domain.InferenceResponse{}, retryErr
		} else {
			return domain.InferenceResponse{}, fmt.Errorf("llama-server respondio sin contenido (body=%s)", second.RawPreview)
		}
	}

	if content == "" {
		return domain.InferenceResponse{}, fmt.Errorf("llama-server respondio sin contenido (body=%s)", first.RawPreview)
	}

	return domain.InferenceResponse{
		ModelName: req.ModelName,
		Content:   content,
		Latency:   time.Since(start),
	}, nil
}

func (e *LocalExecutor) requestCompletion(ctx context.Context, modelName, prompt string, maxTokens int, temperature float64) (llamaResponseMeta, error) {
	body := map[string]interface{}{
		"model": modelName,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens":  maxTokens,
		"temperature": temperature,
		"stream":      false,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return llamaResponseMeta{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, e.serverURL+"/v1/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return llamaResponseMeta{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := e.httpClient.Do(httpReq)
	if err != nil {
		return llamaResponseMeta{}, fmt.Errorf("llama-server no disponible: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return llamaResponseMeta{}, fmt.Errorf("llama-server devolvio status %d", resp.StatusCode)
	}

	rawResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return llamaResponseMeta{}, err
	}

	meta, err := extractLlamaResponseMeta(rawResp)
	if err != nil {
		return llamaResponseMeta{}, err
	}

	return meta, nil
}

func (e *LocalExecutor) ensureServer(ctx context.Context, executablePath, modelPath string) error {
	absExecPath, err := filepath.Abs(executablePath)
	if err != nil {
		return err
	}
	absModelPath, err := filepath.Abs(modelPath)
	if err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.cmd != nil && e.cmd.Process != nil {
		if e.executable == absExecPath && e.modelPath == absModelPath {
			if err := pingServer(ctx, e.httpClient, e.serverURL); err == nil {
				return nil
			}
		}
		_ = stopProcess(e.cmd)
		e.cmd = nil
	}

	ports := e.candidatePorts()
	var lastErr error

	for _, port := range ports {
		if !isPortAvailable(e.serverHost, port) {
			if e.log != nil {
				e.log.Debug("puerto ocupado para llama-server",
					zap.String("host", e.serverHost),
					zap.Int("port", port),
				)
			}
			continue
		}

		serverURL := buildServerURL(e.serverHost, port)
		args := []string{
			"-m", absModelPath,
			"--port", strconv.Itoa(port),
			"--host", e.serverHost,
		}

		cmd := exec.Command(absExecPath, args...)
		cmd.Dir = filepath.Dir(absExecPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			lastErr = err
			continue
		}

		startupCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
		err := waitForServer(startupCtx, e.httpClient, serverURL)
		cancel()
		if err != nil {
			_ = stopProcess(cmd)
			lastErr = err
			continue
		}

		e.cmd = cmd
		e.executable = absExecPath
		e.modelPath = absModelPath
		e.serverURL = serverURL
		if e.log != nil {
			e.log.Info("llama-server iniciado",
				zap.String("host", e.serverHost),
				zap.Int("port", port),
				zap.Int("preferredPort", e.preferredPort),
				zap.Bool("usedFallback", port != e.preferredPort),
				zap.String("url", serverURL),
			)
		}
		return nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("no hay puertos disponibles en rango %d-%d", e.portRangeMin, e.portRangeMax)
	}

	return fmt.Errorf("no se pudo iniciar llama-server: %w", lastErr)
}

func waitForServer(ctx context.Context, client *http.Client, baseURL string) error {
	ticker := time.NewTicker(400 * time.Millisecond)
	defer ticker.Stop()

	for {
		if err := pingServer(ctx, client, baseURL); err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("llama-server no inicio a tiempo")
		case <-ticker.C:
		}
	}
}

func pingServer(ctx context.Context, client *http.Client, baseURL string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/health", nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("health status %d", resp.StatusCode)
	}

	return nil
}

func runCtxOrBackground(ctx context.Context) context.Context {
	if ctx != nil {
		return ctx
	}
	return context.Background()
}

func (e *LocalExecutor) candidatePorts() []int {
	rangeStart := e.portRangeMin
	rangeEnd := e.portRangeMax
	if rangeStart <= 0 {
		rangeStart = 32111
	}
	if rangeEnd <= 0 {
		rangeEnd = 32130
	}
	if rangeEnd < rangeStart {
		rangeStart, rangeEnd = rangeEnd, rangeStart
	}

	ports := make([]int, 0, (rangeEnd-rangeStart+2))
	seen := map[int]struct{}{}

	add := func(p int) {
		if p <= 0 {
			return
		}
		if _, ok := seen[p]; ok {
			return
		}
		seen[p] = struct{}{}
		ports = append(ports, p)
	}

	add(e.preferredPort)
	for p := rangeStart; p <= rangeEnd; p++ {
		add(p)
	}

	return ports
}

func buildServerURL(host string, port int) string {
	return fmt.Sprintf("http://%s:%d", host, port)
}

func isPortAvailable(host string, port int) bool {
	ln, err := net.Listen("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

type llamaResponseMeta struct {
	Content      string
	FinishReason string
	HasReasoning bool
	RawPreview   string
}

func extractLlamaResponseMeta(raw []byte) (llamaResponseMeta, error) {
	meta := llamaResponseMeta{}
	trimmed := strings.TrimSpace(string(raw))
	if len(trimmed) > 280 {
		trimmed = trimmed[:280] + "..."
	}
	meta.RawPreview = trimmed

	var payload map[string]interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return meta, err
	}

	// OpenAI chat format: choices[0].message.content
	if choices, ok := payload["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice0, ok := choices[0].(map[string]interface{}); ok {
			if finishReason, ok := choice0["finish_reason"].(string); ok {
				meta.FinishReason = strings.TrimSpace(finishReason)
			}

			if message, ok := choice0["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					meta.Content = strings.TrimSpace(content)
				}
				if reasoning, ok := message["reasoning_content"].(string); ok {
					meta.HasReasoning = strings.TrimSpace(reasoning) != ""
				}
			}

			// Completion-like fallback: choices[0].text
			if text, ok := choice0["text"].(string); ok {
				meta.Content = strings.TrimSpace(text)
			}
		}
	}

	// Generic fallback sometimes used by wrappers
	if meta.Content == "" {
		if content, ok := payload["content"].(string); ok {
			meta.Content = strings.TrimSpace(content)
		}
	}

	return meta, nil
}

func (e *LocalExecutor) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.cmd == nil {
		return nil
	}

	err := stopProcess(e.cmd)
	e.cmd = nil
	e.modelPath = ""
	e.executable = ""
	return err
}

func stopProcess(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}

	_ = cmd.Process.Signal(os.Interrupt)
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		if err == nil {
			return nil
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.ExitStatus() == 0 {
				return nil
			}
		}
		return nil
	case <-time.After(3 * time.Second):
		_ = cmd.Process.Kill()
		<-done
		return nil
	}
}
