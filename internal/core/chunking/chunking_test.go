package chunking

import (
	"reflect"
	"testing"
)

// EstimateTokens は文字数ベースでトークン数を概算し、最低 1 を返すこと。
func TestEstimateTokens(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"", 1},         // 空でも最低 1
		{"a", 1},        // 1 文字 → (1+3)/4 = 1
		{"abcd", 1},     // 4 文字 → 1
		{"abcde", 2},    // 5 文字 → (5+3)/4 = 2
		{"abcdefgh", 2}, // 8 文字 → 2
		{"あいう", 1},      // 3 runes → 1（rune 単位で数える）
	}
	for _, c := range cases {
		if got := EstimateTokens(c.in); got != c.want {
			t.Errorf("EstimateTokens(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

// Plan は空入力で空を返すこと。
func TestPlanEmpty(t *testing.T) {
	if got := Plan(nil, 100); len(got) != 0 {
		t.Errorf("Plan(nil) = %v, want empty", got)
	}
}

// tokenBudget <= 0 は各 Item を単独チャンクにすること（バルクしない）。
func TestPlanNoBudgetIsAllSingletons(t *testing.T) {
	items := []Item{{Group: "a", Tokens: 1}, {Group: "a", Tokens: 1}, {Group: "b", Tokens: 1}}
	for _, budget := range []int{0, -5} {
		got := Plan(items, budget)
		want := [][]int{{0}, {1}, {2}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Plan(budget=%d) = %v, want %v", budget, got, want)
		}
	}
}

// 同一 Group はトークン予算内でまとまり、予算を超えると切れること。
func TestPlanGroupsByBudget(t *testing.T) {
	items := []Item{
		{Group: "a", Tokens: 3},
		{Group: "a", Tokens: 3},
		{Group: "a", Tokens: 3},
	}
	// 予算 6: [3,3] で 1 チャンク、次の 3 で切れる。
	got := Plan(items, 6)
	want := [][]int{{0, 1}, {2}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Plan = %v, want %v", got, want)
	}
}

// Group が変わる境界でチャンクを切ること（予算に余裕があっても混ぜない）。
func TestPlanSplitsOnGroupBoundary(t *testing.T) {
	items := []Item{
		{Group: "a", Tokens: 1},
		{Group: "a", Tokens: 1},
		{Group: "b", Tokens: 1},
	}
	got := Plan(items, 100)
	want := [][]int{{0, 1}, {2}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Plan = %v, want %v", got, want)
	}
}

// 単独で予算を超える Item は、その 1 件だけのチャンクになること。
func TestPlanSingleOverBudgetIsAlone(t *testing.T) {
	items := []Item{
		{Group: "a", Tokens: 100},
		{Group: "a", Tokens: 1},
	}
	got := Plan(items, 10)
	want := [][]int{{0}, {1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Plan = %v, want %v", got, want)
	}
}

// 全 index がちょうど 1 回ずつ現れること（取りこぼし・重複が無い）。
func TestPlanCoversAllIndicesOnce(t *testing.T) {
	items := []Item{
		{Group: "a", Tokens: 2}, {Group: "a", Tokens: 2}, {Group: "a", Tokens: 2},
		{Group: "b", Tokens: 5}, {Group: "b", Tokens: 5},
	}
	got := Plan(items, 4)
	seen := map[int]int{}
	for _, chunk := range got {
		for _, idx := range chunk {
			seen[idx]++
		}
	}
	for i := range items {
		if seen[i] != 1 {
			t.Errorf("index %d appeared %d times, want 1", i, seen[i])
		}
	}
}
