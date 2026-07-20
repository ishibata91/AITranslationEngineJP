package provider

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// LM Studio など base URL が /v1 を含まない指定でも、/v1 配下のエンドポイントへ届くこと。
// API キーが空のときは Authorization ヘッダを付けないこと。
func TestListModelsNormalizesV1AndOmitsAuthWhenKeyEmpty(t *testing.T) {
	var gotPath string
	var hadAuth bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, hadAuth = r.Header["Authorization"]
		_, _ = io.WriteString(w, `{"data":[{"id":"qwen2.5-7b"},{"id":"llama-3.1-8b"}]}`)
	}))
	defer srv.Close()

	client := NewOpenAICompatible(http.DefaultClient)
	models, err := client.ListModels(context.Background(), Connection{Endpoint: srv.URL, APIKey: ""})
	if err != nil {
		t.Fatalf("ListModels error: %v", err)
	}

	if gotPath != "/v1/models" {
		t.Errorf("path = %q, want /v1/models", gotPath)
	}
	if hadAuth {
		t.Errorf("Authorization header sent though API key is empty")
	}
	if len(models) != 2 || models[0] != "qwen2.5-7b" || models[1] != "llama-3.1-8b" {
		t.Errorf("models = %v, want [qwen2.5-7b llama-3.1-8b]", models)
	}
}

// Translate は engine が組んだ完成 Prompt を /v1/chat/completions へ素通しで送ること。
// System を system メッセージ、User を user メッセージへ写し、内容を加工しないこと。
// API キーがあるときは Bearer で送り、応答本文を訳文として返すこと。
func TestTranslateSendsPromptAndReturnsContent(t *testing.T) {
	var gotPath, gotAuth, gotModel, gotSystem, gotUser string
	var gotRespFmtType, gotRespFmtName string
	var gotRespFmtStrict bool
	var gotSchemaLen int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		var req struct {
			Model    string `json:"model"`
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
			ResponseFormat *struct {
				Type       string `json:"type"`
				JSONSchema struct {
					Name   string          `json:"name"`
					Strict bool            `json:"strict"`
					Schema json.RawMessage `json:"schema"`
				} `json:"json_schema"`
			} `json:"response_format"`
		}
		_ = json.Unmarshal(body, &req)
		gotModel = req.Model
		for _, m := range req.Messages {
			switch m.Role {
			case "system":
				gotSystem = m.Content
			case "user":
				gotUser = m.Content
			}
		}
		if req.ResponseFormat != nil {
			gotRespFmtType = req.ResponseFormat.Type
			gotRespFmtName = req.ResponseFormat.JSONSchema.Name
			gotRespFmtStrict = req.ResponseFormat.JSONSchema.Strict
			gotSchemaLen = len(req.ResponseFormat.JSONSchema.Schema)
		}
		_, _ = io.WriteString(w, `{"choices":[{"message":{"content":"{\"translation\":\"古代ノルドの文章\"}"}}]}`)
	}))
	defer srv.Close()

	prompt := Prompt{
		System: "base 指示\n\nこの台詞の話者の人物像:\n- 声質: 幼い少年の声",
		User:   "Ancient Nord text",
	}
	client := NewOpenAICompatible(http.DefaultClient)
	dest, err := client.Translate(context.Background(),
		Connection{Endpoint: srv.URL, APIKey: "sk-test"}, "qwen2.5-7b", prompt)
	if err != nil {
		t.Fatalf("Translate error: %v", err)
	}

	if gotPath != "/v1/chat/completions" {
		t.Errorf("path = %q, want /v1/chat/completions", gotPath)
	}
	if gotAuth != "Bearer sk-test" {
		t.Errorf("Authorization = %q, want Bearer sk-test", gotAuth)
	}
	if gotModel != "qwen2.5-7b" {
		t.Errorf("model = %q, want qwen2.5-7b", gotModel)
	}
	// 完成 Prompt の System / User を加工せず素通しで送ること（合成は engine の責務）。
	if gotSystem != prompt.System {
		t.Errorf("system メッセージ = %q, want %q", gotSystem, prompt.System)
	}
	if gotUser != prompt.User {
		t.Errorf("user メッセージ = %q, want %q", gotUser, prompt.User)
	}
	if dest != "古代ノルドの文章" {
		t.Errorf("dest = %q, want 古代ノルドの文章", dest)
	}
	// 構造化出力指定（response_format の json_schema）がリクエストへ載ること。
	if gotRespFmtType != "json_schema" {
		t.Errorf("response_format.type = %q, want json_schema", gotRespFmtType)
	}
	if gotRespFmtName != "translation" {
		t.Errorf("json_schema.name = %q, want translation", gotRespFmtName)
	}
	if !gotRespFmtStrict {
		t.Errorf("json_schema.strict = false, want true")
	}
	if gotSchemaLen == 0 {
		t.Errorf("json_schema.schema が空")
	}
}

// extractTranslation は構造化出力 content から訳文を取り出し、取れない場合を型で分ける。
// 正常・構文不正・空応答・必須欠落・空値の各分岐を網羅する（純粋関数の 100% カバレッジ）。
// 失敗はすべて ErrStructuredParse でラップされ、engine が errors.Is で識別して skip できること。
// 「空応答」は実 LLM（7B）が response_format 下で空文字を返した実ケースの再発防止（json.Unmarshal が unexpected end of JSON input）。
func TestExtractTranslation(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
		wantErr bool
	}{
		{name: "正常", content: `{"translation":"古代ノルドの文章"}`, want: "古代ノルドの文章", wantErr: false},
		{name: "構文不正", content: `not json`, want: "", wantErr: true},
		{name: "空応答", content: ``, want: "", wantErr: true},
		{name: "必須欠落", content: `{}`, want: "", wantErr: true},
		{name: "空値", content: `{"translation":""}`, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractTranslation(tt.content)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, ErrStructuredParse) {
				t.Errorf("err = %v, want ErrStructuredParse でラップ", err)
			}
			if got != tt.want {
				t.Errorf("got = %q, want %q", got, tt.want)
			}
		})
	}
}

// statusSkippable は非 200 応答の status を、その行だけ飛ばせる一時失敗(true)か run を止める失敗(false)かに分ける。
// 429（rate limit）と 5xx（サーバ一時）だけを飛ばし、その他の 4xx（認証・不正・不明）と異常な 3xx は止める。純粋関数の網羅。
func TestStatusSkippable(t *testing.T) {
	tests := []struct {
		status int
		want   bool
	}{
		{http.StatusTooManyRequests, true},     // 429
		{http.StatusInternalServerError, true}, // 500
		{http.StatusBadGateway, true},          // 502
		{http.StatusServiceUnavailable, true},  // 503
		{599, true},                            // 5xx 上端
		{http.StatusBadRequest, false},         // 400
		{http.StatusUnauthorized, false},       // 401
		{http.StatusForbidden, false},          // 403
		{http.StatusNotFound, false},           // 404
		{http.StatusTeapot, false},             // 418（429 以外の 4xx）
		{http.StatusMovedPermanently, false},   // 301（異常な 3xx は止める）
	}
	for _, tt := range tests {
		if got := statusSkippable(tt.status); got != tt.want {
			t.Errorf("statusSkippable(%d) = %v, want %v", tt.status, got, tt.want)
		}
	}
}

// 非 200 応答のうち 429・5xx はその行だけ飛ばせる ErrServerTransient でラップし、
// その他の 4xx（認証・不正など）は skippable 番兵を付けず run を止める失敗として返すこと。
func TestTranslateClassifiesNon200(t *testing.T) {
	tests := []struct {
		name          string
		status        int
		wantTransient bool // ErrServerTransient でラップされるべきか
	}{
		{"429 rate limit", http.StatusTooManyRequests, true},
		{"500 server error", http.StatusInternalServerError, true},
		{"503 unavailable", http.StatusServiceUnavailable, true},
		{"401 unauthorized", http.StatusUnauthorized, false},
		{"400 bad request", http.StatusBadRequest, false},
		{"404 not found", http.StatusNotFound, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
			}))
			defer srv.Close()

			client := NewOpenAICompatible(http.DefaultClient)
			_, err := client.Translate(context.Background(), Connection{Endpoint: srv.URL}, "m", Prompt{})
			if err == nil {
				t.Fatalf("非 200 応答でエラーを返していない")
			}
			if errors.Is(err, ErrServerTransient) != tt.wantTransient {
				t.Errorf("err = %v, ErrServerTransient ラップ = %v, want %v", err, errors.Is(err, ErrServerTransient), tt.wantTransient)
			}
			// fatal 側（4xx）は skippable のどの番兵にも該当せず、engine が run を止められること。
			if !tt.wantTransient {
				if errors.Is(err, ErrServerTransient) || errors.Is(err, ErrResponseUnreadable) || errors.Is(err, ErrStructuredParse) {
					t.Errorf("run を止めるべき 4xx が skippable 番兵でラップされた: %v", err)
				}
			}
		})
	}
}

// 応答エンベロープが JSON として読めない、または choices が空の応答は、
// その行だけ飛ばせる ErrResponseUnreadable でラップすること（engine が errors.Is で識別して skip できる）。
func TestTranslateClassifiesUnreadableResponse(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"エンベロープ decode 失敗", `not json`},
		{"choices 無し", `{"choices":[]}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = io.WriteString(w, tt.body)
			}))
			defer srv.Close()

			client := NewOpenAICompatible(http.DefaultClient)
			_, err := client.Translate(context.Background(), Connection{Endpoint: srv.URL}, "m", Prompt{})
			if !errors.Is(err, ErrResponseUnreadable) {
				t.Errorf("err = %v, want ErrResponseUnreadable でラップ", err)
			}
		})
	}
}

// content が schema 非準拠（地の文など）で返った場合、二段パースの段2 が失敗し、
// Translate は訳文を返さずエラーを返すこと（フォールバックしない）。
func TestTranslateReturnsErrorOnBrokenStructuredContent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"choices":[{"message":{"content":"Here is the translation."}}]}`)
	}))
	defer srv.Close()

	client := NewOpenAICompatible(http.DefaultClient)
	dest, err := client.Translate(context.Background(),
		Connection{Endpoint: srv.URL, APIKey: ""}, "qwen2.5-7b",
		Prompt{System: "base", User: "text"})
	if err == nil {
		t.Fatalf("段2 失敗時に error を返していない")
	}
	// engine が errors.Is で構造化パース失敗を識別し行を skip できるよう、Translate の err が ErrStructuredParse チェーンを保つこと。
	if !errors.Is(err, ErrStructuredParse) {
		t.Errorf("err = %v, want ErrStructuredParse でラップ", err)
	}
	if dest != "" {
		t.Errorf("段2 失敗時に dest = %q（空であるべき）", dest)
	}
}
