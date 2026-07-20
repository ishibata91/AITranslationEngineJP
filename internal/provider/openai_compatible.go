// Package provider は AI 翻訳クライアントの port と実装を持つ。本構成で唯一の interface 境界。
package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Connection は AI サービスへの接続情報。画面から都度渡し、永続化しない。
type Connection struct {
	Endpoint string
	APIKey   string
}

// Prompt は engine が組み立てた完成プロンプト。provider は送るだけで、base 指示や口調指示の合成はしない。
// System は system メッセージ（base 指示＋口調指示を合成済み）、User は user メッセージ（機械置換済み原文）。
type Prompt struct {
	System string
	User   string
}

// Translator は AI 翻訳クライアントの port。engine と api がこの interface に依存する。
// プロンプトの組み立て（base 指示・口調指示の合成）は engine が行い、provider は完成 Prompt を送るだけにする。
type Translator interface {
	Translate(ctx context.Context, conn Connection, model string, prompt Prompt) (string, error)
	// TranslateBatch は複数行を 1 リクエストにまとめ、行キーごとの訳文を返す（キー配列・部分成功）。
	// system は共有 system メッセージ、items は行キー付きの機械置換済み原文。
	// 返り値はモデルが正しく返したキーだけを持つ map。欠けたキー・壊れたキーは map に含めない（再送しない）。
	TranslateBatch(ctx context.Context, conn Connection, model string, system string, items []BatchItem) (map[string]string, error)
	ListModels(ctx context.Context, conn Connection) ([]string, error)
}

// BatchItem はキー配列翻訳の 1 件。Key は行キー（line.id の文字列）、User は機械置換済み原文。
type BatchItem struct {
	Key  string
	User string
}

// httpDoer は net/http の Client を抽象化する狭い interface（テストで差し替えるため）。
type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// OpenAICompatible は OpenAI 互換 API（OpenAI 本体・LM Studio など）のクライアント。
type OpenAICompatible struct {
	client httpDoer
}

// NewOpenAICompatible は OpenAI 互換クライアントを生成する。
func NewOpenAICompatible(client httpDoer) *OpenAICompatible {
	return &OpenAICompatible{client: client}
}

// normalizeBase は base URL を /v1 配下へ正規化する。
// LM Studio の "http://127.0.0.1:1234" でも、末尾 /v1 を補って /v1 配下へ届かせる。
func normalizeBase(endpoint string) string {
	base := strings.TrimRight(strings.TrimSpace(endpoint), "/")
	if !strings.HasSuffix(base, "/v1") {
		base += "/v1"
	}
	return base
}

func (c *OpenAICompatible) newRequest(ctx context.Context, method, url string, conn Connection, body []byte) (*http.Request, error) {
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
	if strings.TrimSpace(conn.APIKey) != "" {
		req.Header.Set("Authorization", "Bearer "+conn.APIKey)
	}
	return req, nil
}

// ListModels は /v1/models を引き、利用可能なモデル名の一覧を返す。
func (c *OpenAICompatible) ListModels(ctx context.Context, conn Connection) ([]string, error) {
	req, err := c.newRequest(ctx, http.MethodGet, normalizeBase(conn.Endpoint)+"/models", conn, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("モデル一覧取得: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("モデル一覧取得: status %d", resp.StatusCode)
	}
	var decoded struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("モデル一覧解析: %w", err)
	}
	models := make([]string, 0, len(decoded.Data))
	for _, m := range decoded.Data {
		models = append(models, m.ID)
	}
	return models, nil
}

// Translate は engine が組んだ完成 Prompt を /v1/chat/completions へ送り、応答本文を訳文として返す。
// System を system メッセージ、User を user メッセージとして送る。プロンプトの文面合成は engine の責務。
func (c *OpenAICompatible) Translate(ctx context.Context, conn Connection, model string, prompt Prompt) (string, error) {
	payload := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{Role: "system", Content: prompt.System},
			{Role: "user", Content: prompt.User},
		},
		ResponseFormat: translationResponseFormat(),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("翻訳リクエスト生成: %w", err)
	}
	req, err := c.newRequest(ctx, http.MethodPost, normalizeBase(conn.Endpoint)+"/chat/completions", conn, body)
	if err != nil {
		return "", err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("翻訳要求: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("翻訳要求: status %d", resp.StatusCode)
	}
	var decoded struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return "", fmt.Errorf("翻訳応答のデコード: %w", err)
	}
	if len(decoded.Choices) == 0 {
		return "", fmt.Errorf("翻訳応答に choices が無い")
	}
	translation, err := extractTranslation(decoded.Choices[0].Message.Content)
	if err != nil {
		return "", fmt.Errorf("翻訳応答解析: %w", err)
	}
	return translation, nil
}

type chatRequest struct {
	Model          string          `json:"model"`
	Messages       []chatMessage   `json:"messages"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// responseFormat は OpenAI 互換 API の構造化出力指定（response_format の json_schema）。
// LM Studio・xAI も同一形式で、対応モデルは content を schema 準拠の JSON 文字列で返す。
type responseFormat struct {
	Type       string     `json:"type"`
	JSONSchema jsonSchema `json:"json_schema"`
}

type jsonSchema struct {
	Name   string          `json:"name"`
	Strict bool            `json:"strict"`
	Schema json.RawMessage `json:"schema"`
}

// translationSchema は単一行訳文の出力スキーマ。translation 1 フィールドだけを許し、追加フィールドを拒む。
// 静的定義のため独立モジュールへ切り出さず、provider が固定で持つ。将来の複数対象一括は別 plan で配列版へ育てる。
var translationSchema = json.RawMessage(`{"type":"object","properties":{"translation":{"type":"string"}},"required":["translation"],"additionalProperties":false}`)

// translationResponseFormat は単一行訳文の構造化出力指定を組む。スキーマが固定のため常に同一値を返す。
func translationResponseFormat() *responseFormat {
	return &responseFormat{
		Type: "json_schema",
		JSONSchema: jsonSchema{
			Name:   "translation",
			Strict: true,
			Schema: translationSchema,
		},
	}
}

// ErrStructuredParse は構造化出力から訳文を取り出せなかったことを表す番兵エラー。
// engine は errors.Is でこのエラーを識別し、壊れた訳を確定させず該当行を未訳のまま残す（再実行で再翻訳）。
// モデルが response_format を守らず空・不完全な応答を返した場合に生じる。通信・HTTP 障害とは区別する。
var ErrStructuredParse = errors.New("構造化出力の解析失敗")

// extractTranslation は構造化出力の content 文字列から訳文を取り出す純粋関数。
// content は {"translation": "..."} 形式の JSON 文字列を期待する。訳文が取れない場合を型で分ける。
//   - 構文不正: content が JSON として解釈できない（空応答・途中切れを含む）。
//   - 必須欠落: translation フィールドが無い。
//   - 空値: translation が空文字。
//
// いずれも ErrStructuredParse でラップして返し、engine は該当行を未訳のまま skip する（再実行で再翻訳）。
func extractTranslation(content string) (string, error) {
	var parsed struct {
		Translation *string `json:"translation"`
	}
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return "", fmt.Errorf("%w: 構文不正: %w", ErrStructuredParse, err)
	}
	if parsed.Translation == nil {
		return "", fmt.Errorf("%w: translation フィールドが無い", ErrStructuredParse)
	}
	if *parsed.Translation == "" {
		return "", fmt.Errorf("%w: translation が空", ErrStructuredParse)
	}
	return *parsed.Translation, nil
}

// TranslateBatch は複数行を 1 リクエストにまとめて /v1/chat/completions へ送り、行キーごとの訳文を返す。
// system を system メッセージ、行キー付きの原文一覧を user メッセージとし、出力は行キーごとの構造化 JSON を求める。
// 返り値はモデルが正しく返したキーだけを持つ map（部分成功）。欠けたキー・空応答・壊れたキーは map に含めず、
// 呼び出し側（engine）が該当行を未訳のまま残して再実行で拾う。バッチ再送はしない。通信・HTTP 障害は error で返す。
func (c *OpenAICompatible) TranslateBatch(ctx context.Context, conn Connection, model string, system string, items []BatchItem) (map[string]string, error) {
	keys := make([]string, len(items))
	var user strings.Builder
	user.WriteString("次の各行を日本語へ訳し、行頭のキーごとに訳文を返す。キーは変えない。\n")
	for i, it := range items {
		keys[i] = it.Key
		user.WriteString(it.Key)
		user.WriteString(": ")
		user.WriteString(it.User)
		user.WriteString("\n")
	}
	payload := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: user.String()},
		},
		ResponseFormat: batchResponseFormat(keys),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("一括翻訳リクエスト生成: %w", err)
	}
	req, err := c.newRequest(ctx, http.MethodPost, normalizeBase(conn.Endpoint)+"/chat/completions", conn, body)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("一括翻訳要求: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("一括翻訳要求: status %d", resp.StatusCode)
	}
	var decoded struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("一括翻訳応答のデコード: %w", err)
	}
	if len(decoded.Choices) == 0 {
		return nil, fmt.Errorf("一括翻訳応答に choices が無い")
	}
	return extractBatchTranslations(decoded.Choices[0].Message.Content, keys), nil
}

// batchResponseFormat は行キーごとの訳文を求めるキー配列スキーマの構造化出力指定を組む。
// 各キーの値は {translation:string} オブジェクト。additionalProperties:false で余分なキーを禁じるが、
// 部分成功を許すため strict は false にし required も置かない（モデルが一部キーを落とすのを許容する）。
func batchResponseFormat(keys []string) *responseFormat {
	props := make(map[string]any, len(keys))
	for _, k := range keys {
		props[k] = map[string]any{
			"type":                 "object",
			"properties":           map[string]any{"translation": map[string]any{"type": "string"}},
			"required":             []string{"translation"},
			"additionalProperties": false,
		}
	}
	schema, _ := json.Marshal(map[string]any{
		"type":                 "object",
		"properties":           props,
		"additionalProperties": false,
	})
	return &responseFormat{
		Type: "json_schema",
		JSONSchema: jsonSchema{
			Name:   "translations",
			Strict: false,
			Schema: schema,
		},
	}
}

// extractBatchTranslations は一括応答の content から、要求キーごとの訳文を取り出す純粋関数。
// content は {"<key>": {"translation": "..."}, ...} 形式の JSON を期待する。
// 返り値には、キーが存在し translation が非空文字のものだけを載せる（部分成功）。
// 構文不正（空応答・途中切れ）は空 map を返し、全キーを欠落扱いにする。壊れたキー（欠落・空・型不正）は載せない。
func extractBatchTranslations(content string, keys []string) map[string]string {
	out := make(map[string]string, len(keys))
	var parsed map[string]struct {
		Translation *string `json:"translation"`
	}
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return out
	}
	for _, k := range keys {
		if v, ok := parsed[k]; ok && v.Translation != nil && *v.Translation != "" {
			out[k] = *v.Translation
		}
	}
	return out
}
