// Package harness は api 公開面の非劣化テスト基盤を持つ。
// 実 SQLite（temp DB）・決定的 fake provider・抽出子 seam を組み、RunExtractAndTranslate を端から端まで通し、
// 送信プロンプト列と DB 最終状態を golden として捕獲する。合成入力（CI 恒久）と実 .esm 入力（local 限定）の両方で使う。
// 本 package は composition root 的に store・engine・provider・api を束ねるため、特定の層に属さない（arch-lint の component 外）。
package harness

import (
	"context"
	"fmt"
	"sync"

	"aitranslationenginejp/internal/core/runtimetag"
	"aitranslationenginejp/internal/provider"
)

// RecordingProvider は決定的な fake の翻訳 provider。受け取った Prompt を送信順に記録し、規則訳を返す。
// AI 出力の非決定性を排し、出力を入力とコードだけに依存させる（非劣化判定の前提）。
// engine と api の両方へ同一インスタンスを渡し、Run 中の全 Translate を 1 列に記録する。
type RecordingProvider struct {
	mu      sync.Mutex
	prompts []RecordedPrompt
	models  []string
	// Fixed は原文（user メッセージ）ごとの固定訳。連番訳では通せない訳形（中黒区切りのカタカナ人名など）を
	// 与えて、その訳形に依る経路を決定的に通すために使う。ここに無い原文は連番訳へ落ちる。
	Fixed map[string]string
}

// RecordedPrompt は 1 回の Translate 呼び出しの記録（送信モデルと完成プロンプト）。
type RecordedPrompt struct {
	Model  string
	System string
	User   string
}

// Translate は規則訳を返し、プロンプトを送信順に記録する。訳文は呼び出し順の連番で決定的にする
// （DB 最終状態の dest が安定し、プロンプト列とは独立に書き戻し順序を観測できる）。
// Fixed に原文がある場合はその固定訳を返す（訳形に依る経路を通すため）。
func (r *RecordingProvider) Translate(_ context.Context, _ provider.Connection, model string, prompt provider.Prompt) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prompts = append(r.prompts, RecordedPrompt{Model: model, System: prompt.System, User: prompt.User})
	if fixed, ok := r.Fixed[prompt.User]; ok {
		return fixed, nil
	}
	dest := fmt.Sprintf("訳%03d", len(r.prompts))
	// 送信プロンプトに含まれる生タグ（<...>）を保持するモデルを模す。engine の CountMissing が欠落 0 と判定でき、
	// タグ保護（機械置換から守る退避 → 生タグで送信 → モデルが保持）を結合オラクルで観測できるようにする。
	for _, tag := range runtimetag.Tags(prompt.User) {
		dest += tag
	}
	return dest, nil
}

// ListModels は固定の 1 モデルを返す（本基盤は GetModels を通さないため呼ばれない想定の最小実装）。
func (r *RecordingProvider) ListModels(_ context.Context, _ provider.Connection) ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.models = append(r.models, "list")
	return []string{"fake-model"}, nil
}

// Prompts は記録した送信プロンプト列を送信順で返す。
func (r *RecordingProvider) Prompts() []RecordedPrompt {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]RecordedPrompt, len(r.prompts))
	copy(out, r.prompts)
	return out
}
