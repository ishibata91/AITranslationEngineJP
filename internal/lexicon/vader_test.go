package lexicon

import (
	"os"
	"path/filepath"
	"testing"
)

// LoadVADER は valence 絶対値が閾値以上の語だけを集合に入れ、弱い語や大文字は正しく扱うこと。
func TestLoadVADERSelectsStrongValence(t *testing.T) {
	// タブ区切り（語・valence 平均・標準偏差・生評定）の最小辞書。
	// love は正の強語、hate は負の強語、meh は弱語、GREAT は大文字の強語（小文字照合を確認）。
	content := "love\t3.2\t0.4\t[3, 3, 3]\n" +
		"hate\t-2.7\t1.0\t[-4, -3]\n" +
		"meh\t0.5\t0.9\t[0, 1, 1]\n" +
		"GREAT\t3.1\t0.7\t[3, 4]\n"
	path := filepath.Join(t.TempDir(), "vader.txt")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("辞書ファイルの作成: %v", err)
	}

	v, err := LoadVADER(path)
	if err != nil {
		t.Fatalf("LoadVADER: %v", err)
	}
	if !v.IsStrongEmotion("love") {
		t.Error("正の強語 love が集合に無い")
	}
	if !v.IsStrongEmotion("hate") {
		t.Error("負の強語 hate が集合に無い")
	}
	if v.IsStrongEmotion("meh") {
		t.Error("弱語 meh が集合に入った")
	}
	if !v.IsStrongEmotion("great") {
		t.Error("大文字の強語 GREAT が小文字 great で引けない")
	}
}

// LoadVADER は閾値 strongThreshold の境界を含む（>=）で判定すること。
func TestLoadVADERThresholdBoundary(t *testing.T) {
	// edge は閾値ちょうど（採用）、below は閾値未満（除外）。
	content := "edge\t1.5\t0.3\t[1, 2]\n" +
		"below\t1.4\t0.3\t[1, 2]\n"
	path := filepath.Join(t.TempDir(), "vader.txt")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("辞書ファイルの作成: %v", err)
	}

	v, err := LoadVADER(path)
	if err != nil {
		t.Fatalf("LoadVADER: %v", err)
	}
	if !v.IsStrongEmotion("edge") {
		t.Error("閾値ちょうどの edge が集合に無い")
	}
	if v.IsStrongEmotion("below") {
		t.Error("閾値未満の below が集合に入った")
	}
}

// LoadVADER は辞書ファイルが無ければエラーを返すこと。
func TestLoadVADERMissingFile(t *testing.T) {
	if _, err := LoadVADER(filepath.Join(t.TempDir(), "absent.txt")); err == nil {
		t.Fatal("辞書ファイルが無いのにエラーを返さない")
	}
}
