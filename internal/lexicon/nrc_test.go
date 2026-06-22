package lexicon

import (
	"os"
	"path/filepath"
	"testing"
)

// LoadNRC は強い感情が 1 の語だけを集合に入れ、弱い感情や 0 の行は無視すること。
func TestLoadNRCSelectsStrongEmotion(t *testing.T) {
	// タブ区切り（語・感情・0|1）の最小辞書。angry は怒り 1、calm は怒り 0、trust は対象外の感情。
	content := "angry\tanger\t1\n" +
		"calm\tanger\t0\n" +
		"hope\ttrust\t1\n" +
		"joyful\tjoy\t1\n"
	path := filepath.Join(t.TempDir(), "nrc.txt")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("辞書ファイルの作成: %v", err)
	}

	nrc, err := LoadNRC(path)
	if err != nil {
		t.Fatalf("LoadNRC: %v", err)
	}
	if !nrc.IsStrongEmotion("angry") {
		t.Error("強い感情(怒り 1)の angry が集合に無い")
	}
	if !nrc.IsStrongEmotion("joyful") {
		t.Error("強い感情(喜び 1)の joyful が集合に無い")
	}
	if nrc.IsStrongEmotion("calm") {
		t.Error("怒り 0 の calm が集合に入った")
	}
	if nrc.IsStrongEmotion("hope") {
		t.Error("対象外の感情(trust)の hope が集合に入った")
	}
}

// LoadNRC は辞書ファイルが無ければエラーを返すこと。
func TestLoadNRCMissingFile(t *testing.T) {
	if _, err := LoadNRC(filepath.Join(t.TempDir(), "absent.txt")); err == nil {
		t.Fatal("辞書ファイルが無いのにエラーを返さない")
	}
}
