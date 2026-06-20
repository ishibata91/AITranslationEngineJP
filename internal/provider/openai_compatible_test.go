package provider

import (
	"context"
	"encoding/json"
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
		_, _ = io.WriteString(w, `{"choices":[{"message":{"content":"古代ノルドの文章"}}]}`)
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
}
