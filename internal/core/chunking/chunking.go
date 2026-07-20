// Package chunking は台詞のバルク翻訳で、行をトークン予算でチャンクへ分ける純粋ルール。
// 正確な tokenizer を持てない OpenAI 互換の任意モデル向けに、トークン数は文字数ベースで概算する。
// 予算は厳密上限でなく目安に使う。副作用を持たず、入出力で完結する。
package chunking

// EstimateTokens は文字列のトークン数を安価に概算する純粋関数。
// 英語主体の原文で 1 トークン ≈ 4 文字の経験則を使い、最低 1 を返す。
// 正確さより速さを取る。予算判定の目安に使う値で、厳密なトークン数ではない。
func EstimateTokens(s string) int {
	n := len([]rune(s))
	est := (n + 3) / 4
	if est < 1 {
		return 1
	}
	return est
}

// Item はチャンク分割の入力 1 件。
// Group は同一まとまり（同一 INFO・同一口調指示）の識別子で、Group が変わる境界ではチャンクを切る。
// Tokens は EstimateTokens で求めた概算トークン数。
type Item struct {
	Group  string
	Tokens int
}

// Plan は Item 列を、同一 Group 連続かつトークン予算内のチャンク（元 index の並び）へ分ける純粋関数。
// 返り値の各要素は 1 チャンク分の index 列で、入力順を保つ。全 index がちょうど 1 回ずつ現れる。
//   - tokenBudget <= 0: バルクしない。各 Item を単独チャンクにする。
//   - Group が変わる、または Item を足すと予算を超える場合にチャンクを切る。
//   - 単独で予算を超える Item は、その 1 件だけのチャンクになる。
func Plan(items []Item, tokenBudget int) [][]int {
	chunks := make([][]int, 0, len(items))
	if len(items) == 0 {
		return chunks
	}
	if tokenBudget <= 0 {
		for i := range items {
			chunks = append(chunks, []int{i})
		}
		return chunks
	}

	var cur []int
	curGroup := ""
	curTokens := 0
	flush := func() {
		if len(cur) > 0 {
			chunks = append(chunks, cur)
			cur = nil
			curTokens = 0
		}
	}
	for i, it := range items {
		// 既存チャンクがあり、Group が変わる、または追加で予算超過なら、先に切る。
		if len(cur) > 0 && (it.Group != curGroup || curTokens+it.Tokens > tokenBudget) {
			flush()
		}
		if len(cur) == 0 {
			curGroup = it.Group
		}
		cur = append(cur, i)
		curTokens += it.Tokens
	}
	flush()
	return chunks
}
