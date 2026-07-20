package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// DefaultXAIEndpoint は xAI API の既定 base。Connection.Endpoint が空のとき使う。
// batch クライアントと api のモデル一覧取得が共有する xAI 既定 base の正本（重複定義を避ける）。
const DefaultXAIEndpoint = "https://api.x.ai"

// XAIBatch は xAI の batch API（非同期の大量翻訳）のクライアント。非同期 port BatchTranslator を実装する。
// 同期 OpenAICompatible とは別実装。openai_compatible と同じく httpDoer を狭い interface で受け、
// fake HTTP で単体テストできる。実 xAI API へは自動テストで触れない。
type XAIBatch struct {
	client httpDoer
}

// NewXAIBatch は xAI batch クライアントを生成する。bootstrap が唯一の生成点。
func NewXAIBatch(client httpDoer) *XAIBatch {
	return &XAIBatch{client: client}
}

// xaiBase は Connection.Endpoint を xAI の /v1 配下 base へ正規化する。空なら既定の xAI endpoint を使う。
func xaiBase(conn Connection) string {
	ep := strings.TrimSpace(conn.Endpoint)
	if ep == "" {
		ep = DefaultXAIEndpoint
	}
	return normalizeBase(ep)
}

// SubmitBatch は要求群を JSONL 化して Files API へ upload し、その file を参照する batch を作って外部 batch ID を返す。
// JSONL file 方式を採る理由: 各行 body に同期経路と同一の chat completions payload（response_format の
// json_schema strict を含む）をそのまま載せられ、訳文の取り出しも同期と同じ形で扱えるため。
func (c *XAIBatch) SubmitBatch(ctx context.Context, conn Connection, model string, requests []BatchRequest) (string, error) {
	jsonl, err := buildBatchJSONL(model, requests)
	if err != nil {
		return "", err
	}
	fileID, err := c.uploadFile(ctx, conn, jsonl)
	if err != nil {
		return "", err
	}
	return c.createBatch(ctx, conn, fileID)
}

// buildBatchJSONL は要求群を xAI batch の JSONL（1 行 1 リクエスト）へ組む。
// 各行は custom_id / method / url(/v1/chat/completions) / body。body は同期と同一の chat completions payload。
func buildBatchJSONL(model string, requests []BatchRequest) ([]byte, error) {
	var buf bytes.Buffer
	for _, r := range requests {
		line := xaiJSONLLine{
			CustomID: r.CustomID,
			Method:   http.MethodPost,
			URL:      "/v1/chat/completions",
			Body: chatRequest{
				Model: model,
				Messages: []chatMessage{
					{Role: "system", Content: r.Prompt.System},
					{Role: "user", Content: r.Prompt.User},
				},
				ResponseFormat: translationResponseFormat(),
			},
		}
		encoded, err := json.Marshal(line)
		if err != nil {
			return nil, fmt.Errorf("batch JSONL 行の生成: %w", err)
		}
		buf.Write(encoded)
		buf.WriteByte('\n')
	}
	return buf.Bytes(), nil
}

// xaiJSONLLine は JSONL の 1 行。body は同期と同一の chat completions payload（chatRequest）。
type xaiJSONLLine struct {
	CustomID string      `json:"custom_id"`
	Method   string      `json:"method"`
	URL      string      `json:"url"`
	Body     chatRequest `json:"body"`
}

// uploadFile は JSONL を Files API（POST /v1/files）へ multipart upload し、file ID を返す。
func (c *XAIBatch) uploadFile(ctx context.Context, conn Connection, jsonl []byte) (string, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	if err := w.WriteField("purpose", "batch"); err != nil {
		return "", fmt.Errorf("file upload の purpose 書き込み: %w", err)
	}
	part, err := w.CreateFormFile("file", "batch.jsonl")
	if err != nil {
		return "", fmt.Errorf("file upload の part 生成: %w", err)
	}
	if _, err = part.Write(jsonl); err != nil {
		return "", fmt.Errorf("file upload の本体書き込み: %w", err)
	}
	if err = w.Close(); err != nil {
		return "", fmt.Errorf("file upload の multipart close: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, xaiBase(conn)+"/files", &body)
	if err != nil {
		return "", fmt.Errorf("file upload リクエスト生成: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	setBearer(req, conn)
	var decoded struct {
		ID string `json:"id"`
	}
	if err = c.doDecode(req, "file upload", &decoded); err != nil {
		return "", err
	}
	if decoded.ID == "" {
		return "", fmt.Errorf("file upload: file ID が空")
	}
	return decoded.ID, nil
}

// createBatch は upload 済み file を参照する batch を作り、外部 batch ID を返す。
// 応答の外部 ID は doc の記述揺れ（batch_id / id）に備え両方を受ける。
func (c *XAIBatch) createBatch(ctx context.Context, conn Connection, fileID string) (string, error) {
	payload, err := json.Marshal(map[string]string{
		"name":          "aitranslation",
		"input_file_id": fileID,
	})
	if err != nil {
		return "", fmt.Errorf("batch 作成リクエスト生成: %w", err)
	}
	req, err := c.newJSONRequest(ctx, http.MethodPost, xaiBase(conn)+"/batches", conn, payload)
	if err != nil {
		return "", err
	}
	var decoded struct {
		BatchID string `json:"batch_id"`
		ID      string `json:"id"`
	}
	if err = c.doDecode(req, "batch 作成", &decoded); err != nil {
		return "", err
	}
	id := decoded.BatchID
	if id == "" {
		id = decoded.ID
	}
	if id == "" {
		return "", fmt.Errorf("batch 作成: 外部 batch ID が空")
	}
	return id, nil
}

// PollBatch は GET /v1/batches/{id} の state 件数を読み、終端判定（未処理が 0）を付けて返す。
// vendor 固有の全体 status 文字列は doc に明記が無いため、件数（num_pending）で終端を判定する。
func (c *XAIBatch) PollBatch(ctx context.Context, conn Connection, externalBatchID string) (BatchStatus, error) {
	req, err := c.newJSONRequest(ctx, http.MethodGet, xaiBase(conn)+"/batches/"+url.PathEscape(externalBatchID), conn, nil)
	if err != nil {
		return BatchStatus{}, err
	}
	var decoded struct {
		State struct {
			NumRequests  int `json:"num_requests"`
			NumPending   int `json:"num_pending"`
			NumSuccess   int `json:"num_success"`
			NumError     int `json:"num_error"`
			NumCancelled int `json:"num_cancelled"`
		} `json:"state"`
	}
	if err = c.doDecode(req, "batch 状態確認", &decoded); err != nil {
		return BatchStatus{}, err
	}
	s := decoded.State
	return BatchStatus{
		Total:     s.NumRequests,
		Pending:   s.NumPending,
		Succeeded: s.NumSuccess,
		Failed:    s.NumError,
		Cancelled: s.NumCancelled,
		Done:      s.NumPending == 0,
	}, nil
}

// FetchResults は GET /v1/batches/{id}/results をページングで最後までたどり、custom_id ごとの結果へ写す。
// 個別リクエストの失敗（error_message あり・choices 空・content が空や非 JSON）は失敗種別へ写す。
func (c *XAIBatch) FetchResults(ctx context.Context, conn Connection, externalBatchID string) ([]BatchResult, error) {
	var out []BatchResult
	token := ""
	for {
		page, err := c.fetchResultsPage(ctx, conn, externalBatchID, token)
		if err != nil {
			return nil, err
		}
		for _, r := range page.Results {
			out = append(out, toBatchResult(r))
		}
		if page.PaginationToken == nil || *page.PaginationToken == "" {
			return out, nil
		}
		token = *page.PaginationToken
	}
}

// fetchResultsPage は結果 1 ページ分を取得する。token が空なら先頭ページ。
func (c *XAIBatch) fetchResultsPage(ctx context.Context, conn Connection, externalBatchID, token string) (xaiResultsPage, error) {
	u := xaiBase(conn) + "/batches/" + url.PathEscape(externalBatchID) + "/results?limit=100"
	if token != "" {
		u += "&pagination_token=" + url.QueryEscape(token)
	}
	req, err := c.newJSONRequest(ctx, http.MethodGet, u, conn, nil)
	if err != nil {
		return xaiResultsPage{}, err
	}
	var page xaiResultsPage
	if err = c.doDecode(req, "batch 結果取得", &page); err != nil {
		return xaiResultsPage{}, err
	}
	return page, nil
}

// xaiResultsPage は結果ページの応答形。訳文は batch_result.response.chat_get_completion.choices[].message.content。
type xaiResultsPage struct {
	Results []xaiResultRow `json:"results"`
	// pagination_token は次ページ用。null（nil）または空文字で終端。
	PaginationToken *string `json:"pagination_token"`
}

type xaiResultRow struct {
	BatchRequestID string `json:"batch_request_id"`
	BatchResult    struct {
		Response struct {
			ChatGetCompletion struct {
				Choices []struct {
					Message struct {
						Content string `json:"content"`
					} `json:"message"`
				} `json:"choices"`
			} `json:"chat_get_completion"`
		} `json:"response"`
	} `json:"batch_result"`
	// error_message は個別リクエストの失敗理由。成功時は null。
	ErrorMessage *string `json:"error_message"`
}

// toBatchResult は xAI の 1 結果行を port の BatchResult へ写す。失敗は同期と同じ番兵種別へ寄せる。
// error_message あり → ErrServerTransient（一時失敗として再送信で回収）。
// choices 空 → ErrResponseUnreadable。content が空・非 JSON → extractTranslation の ErrStructuredParse。
func toBatchResult(r xaiResultRow) BatchResult {
	res := BatchResult{CustomID: r.BatchRequestID}
	if r.ErrorMessage != nil && strings.TrimSpace(*r.ErrorMessage) != "" {
		res.Err = fmt.Errorf("%w: %s", ErrServerTransient, *r.ErrorMessage)
		return res
	}
	choices := r.BatchResult.Response.ChatGetCompletion.Choices
	if len(choices) == 0 {
		res.Err = fmt.Errorf("%w: choices が無い", ErrResponseUnreadable)
		return res
	}
	translation, err := extractTranslation(choices[0].Message.Content)
	if err != nil {
		res.Err = err
		return res
	}
	res.Translation = translation
	return res
}

// newJSONRequest は JSON 本体（GET は nil）のリクエストを組み、Bearer 認証を付ける。
func (c *XAIBatch) newJSONRequest(ctx context.Context, method, url string, conn Connection, body []byte) (*http.Request, error) {
	var reader *bytes.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	} else {
		reader = bytes.NewReader(nil)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return nil, fmt.Errorf("リクエスト生成: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	setBearer(req, conn)
	return req, nil
}

// setBearer は API キーがあるとき Authorization: Bearer を付ける。
func setBearer(req *http.Request, conn Connection) {
	if strings.TrimSpace(conn.APIKey) != "" {
		req.Header.Set("Authorization", "Bearer "+conn.APIKey)
	}
}

// doDecode はリクエストを送り、200 応答の JSON を out へ decode する。非 200 は失敗にする。
// where は失敗メッセージの接頭（"batch 作成" など）。
func (c *XAIBatch) doDecode(req *http.Request, where string, out any) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %w", where, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: status %d", where, resp.StatusCode)
	}
	if err = json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("%s の応答解析: %w", where, err)
	}
	return nil
}
