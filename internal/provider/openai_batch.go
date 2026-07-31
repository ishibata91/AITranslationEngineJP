package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// DefaultOpenAIEndpoint は OpenAI 公式 API の既定 base。
const DefaultOpenAIEndpoint = "https://api.openai.com/v1"

// OpenAIBatch は OpenAI Batch API のクライアント。完成プロンプトは共通 JSONL へ載せる。
type OpenAIBatch struct {
	client httpDoer
}

// NewOpenAIBatch は OpenAI batch クライアントを生成する。bootstrap が唯一の生成点。
func NewOpenAIBatch(client httpDoer) *OpenAIBatch {
	return &OpenAIBatch{client: client}
}

func openAIBase(conn Connection) string {
	ep := strings.TrimSpace(conn.Endpoint)
	if ep == "" {
		ep = DefaultOpenAIEndpoint
	}
	return normalizeBase(ep)
}

// SubmitBatch は共通 JSONL を Files API へ送り、Chat Completions の 24 時間 batch を作る。
func (c *OpenAIBatch) SubmitBatch(ctx context.Context, conn Connection, model string, requests []BatchRequest) (string, error) {
	jsonl, err := buildBatchJSONL(model, requests)
	if err != nil {
		return "", err
	}
	fileID, err := c.uploadFile(ctx, conn, jsonl)
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(map[string]string{
		"input_file_id":     fileID,
		"endpoint":          "/v1/chat/completions",
		"completion_window": "24h",
	})
	if err != nil {
		return "", fmt.Errorf("OpenAI batch 作成リクエスト生成: %w", err)
	}
	var decoded openAIBatchObject
	if err := c.doJSON(ctx, http.MethodPost, openAIBase(conn)+"/batches", conn, payload, "OpenAI batch 作成", &decoded); err != nil {
		return "", err
	}
	if decoded.ID == "" {
		return "", fmt.Errorf("OpenAI batch 作成: 外部 batch ID が空")
	}
	return decoded.ID, nil
}

func (c *OpenAIBatch) uploadFile(ctx context.Context, conn Connection, jsonl []byte) (string, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	if err := w.WriteField("purpose", "batch"); err != nil {
		return "", fmt.Errorf("OpenAI file upload の purpose 書き込み: %w", err)
	}
	part, err := w.CreateFormFile("file", "batch.jsonl")
	if err != nil {
		return "", fmt.Errorf("OpenAI file upload の part 生成: %w", err)
	}
	if _, err = part.Write(jsonl); err != nil {
		return "", fmt.Errorf("OpenAI file upload の本体書き込み: %w", err)
	}
	if err = w.Close(); err != nil {
		return "", fmt.Errorf("OpenAI file upload の multipart close: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIBase(conn)+"/files", &body)
	if err != nil {
		return "", fmt.Errorf("OpenAI file upload リクエスト生成: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	setBearer(req, conn)
	var decoded struct {
		ID string `json:"id"`
	}
	if err := c.doDecode(req, "OpenAI file upload", &decoded); err != nil {
		return "", err
	}
	if decoded.ID == "" {
		return "", fmt.Errorf("OpenAI file upload: file ID が空")
	}
	return decoded.ID, nil
}

// PollBatch は OpenAI の status と request_counts を共通状態へ写す。
func (c *OpenAIBatch) PollBatch(ctx context.Context, conn Connection, externalBatchID string) (BatchStatus, error) {
	obj, err := c.getBatch(ctx, conn, externalBatchID)
	if err != nil {
		return BatchStatus{}, err
	}
	pending := obj.RequestCounts.Total - obj.RequestCounts.Completed - obj.RequestCounts.Failed
	if pending < 0 {
		pending = 0
	}
	return BatchStatus{
		Total:     obj.RequestCounts.Total,
		Pending:   pending,
		Succeeded: obj.RequestCounts.Completed,
		Failed:    obj.RequestCounts.Failed,
		Done:      openAIBatchTerminal(obj.Status),
	}, nil
}

func openAIBatchTerminal(status string) bool {
	switch status {
	case "completed", "failed", "expired", "cancelled":
		return true
	default:
		return false
	}
}

// FetchResults は成功と失敗の JSONL を読み、成功行だけ訳文を返す。結果 file が無い終端は空結果とする。
func (c *OpenAIBatch) FetchResults(ctx context.Context, conn Connection, externalBatchID string) ([]BatchResult, error) {
	obj, err := c.getBatch(ctx, conn, externalBatchID)
	if err != nil {
		return nil, err
	}
	var out []BatchResult
	for _, fileID := range []string{obj.OutputFileID, obj.ErrorFileID} {
		if fileID == "" {
			continue
		}
		rows, err := c.fetchFileResults(ctx, conn, fileID)
		if err != nil {
			return nil, err
		}
		out = append(out, rows...)
	}
	return out, nil
}

type openAIBatchObject struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	OutputFileID  string `json:"output_file_id"`
	ErrorFileID   string `json:"error_file_id"`
	RequestCounts struct {
		Total     int `json:"total"`
		Completed int `json:"completed"`
		Failed    int `json:"failed"`
	} `json:"request_counts"`
}

func (c *OpenAIBatch) getBatch(ctx context.Context, conn Connection, externalBatchID string) (openAIBatchObject, error) {
	var decoded openAIBatchObject
	err := c.doJSON(ctx, http.MethodGet, openAIBase(conn)+"/batches/"+url.PathEscape(externalBatchID), conn, nil, "OpenAI batch 状態取得", &decoded)
	return decoded, err
}

type openAIResultRow struct {
	CustomID string `json:"custom_id"`
	Response *struct {
		StatusCode int `json:"status_code"`
		Body       struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		} `json:"body"`
	} `json:"response"`
	Error *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (c *OpenAIBatch) fetchFileResults(ctx context.Context, conn Connection, fileID string) ([]BatchResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, openAIBase(conn)+"/files/"+url.PathEscape(fileID)+"/content", nil)
	if err != nil {
		return nil, fmt.Errorf("OpenAI result file リクエスト生成: %w", err)
	}
	setBearer(req, conn)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OpenAI result file 取得: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI result file 取得: status %d", resp.StatusCode)
	}
	return decodeOpenAIResults(resp.Body)
}

func decodeOpenAIResults(r io.Reader) ([]BatchResult, error) {
	var out []BatchResult
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == "" {
			continue
		}
		var row openAIResultRow
		if err := json.Unmarshal(scanner.Bytes(), &row); err != nil {
			return nil, fmt.Errorf("OpenAI result JSONL の解析: %w", err)
		}
		out = append(out, toOpenAIBatchResult(row))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("OpenAI result JSONL の読み込み: %w", err)
	}
	return out, nil
}

func toOpenAIBatchResult(row openAIResultRow) BatchResult {
	res := BatchResult{CustomID: row.CustomID}
	if row.Error != nil {
		res.Err = fmt.Errorf("%w: %s: %s", ErrServerTransient, row.Error.Code, row.Error.Message)
		return res
	}
	if row.Response == nil {
		res.Err = fmt.Errorf("%w: response が無い", ErrResponseUnreadable)
		return res
	}
	if row.Response.StatusCode < 200 || row.Response.StatusCode >= 300 {
		res.Err = fmt.Errorf("%w: status %d", ErrServerTransient, row.Response.StatusCode)
		return res
	}
	if len(row.Response.Body.Choices) == 0 {
		res.Err = fmt.Errorf("%w: choices が無い", ErrResponseUnreadable)
		return res
	}
	translation, err := extractTranslation(row.Response.Body.Choices[0].Message.Content)
	if err != nil {
		res.Err = err
		return res
	}
	res.Translation = translation
	return res
}

func (c *OpenAIBatch) doJSON(ctx context.Context, method, target string, conn Connection, body []byte, where string, out any) error {
	req, err := http.NewRequestWithContext(ctx, method, target, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("%s リクエスト生成: %w", where, err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	setBearer(req, conn)
	return c.doDecode(req, where, out)
}

func (c *OpenAIBatch) doDecode(req *http.Request, where string, out any) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %w", where, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("%s: status %d", where, resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("%s の応答解析: %w", where, err)
	}
	return nil
}
