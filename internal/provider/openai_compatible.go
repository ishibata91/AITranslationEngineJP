// Package provider は AI 翻訳クライアントの port と実装を持つ。本構成で唯一の interface 境界。
package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Connection は AI サービスへの接続情報。画面から都度渡し、永続化しない。
type Connection struct {
	Endpoint string
	APIKey   string
}

// Translator は AI 翻訳クライアントの port。engine と api がこの interface に依存する。
type Translator interface {
	Translate(ctx context.Context, conn Connection, model, source string) (string, error)
	ListModels(ctx context.Context, conn Connection) ([]string, error)
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

// Translate は原文を Japanese へ翻訳して返す。
func (c *OpenAICompatible) Translate(ctx context.Context, conn Connection, model, source string) (string, error) {
	payload := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{Role: "system", Content: translationDirective},
			{Role: "user", Content: source},
		},
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
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return "", fmt.Errorf("翻訳応答解析: %w", err)
	}
	if len(decoded.Choices) == 0 {
		return "", fmt.Errorf("翻訳応答に choices が無い")
	}
	return decoded.Choices[0].Message.Content, nil
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// translationDirective は本文翻訳の指示。Skyrim の英文を自然な日本語へ訳す。
const translationDirective = "あなたは Skyrim Mod の翻訳者です。" +
	"与えられた英語の本文を、原文の意味と語調を保った自然な日本語へ翻訳してください。" +
	"訳文だけを出力し、説明や注釈は加えないでください。"
