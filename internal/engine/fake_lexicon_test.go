package engine

// fakeLexicon は engine テスト用の最小感情辞書。常に非強感情を返し、core/linefeatures.EmotionLexicon を満たす。
// 口調生成の感情経路を決定的に通すための代替で、実 NRC 辞書から engine テストを切り離す。
type fakeLexicon struct{}

// IsStrongEmotion は常に false を返す（engine テストは感情強度に依存しない経路を確かめる）。
func (fakeLexicon) IsStrongEmotion(string) bool { return false }
