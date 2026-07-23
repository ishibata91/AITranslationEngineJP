package api

import "testing"

// classifyStringsPresence（純粋関数）の判定を境界条件込みで固定する。
// 対象拡張子 3 種・言語接尾・大文字小文字の揺れ・非対象ファイルの無視を確かめる。
func TestClassifyStringsPresence(t *testing.T) {
	cases := []struct {
		name  string
		files []string
		want  StringsPresenceView
	}{
		{
			name:  "英日そろい",
			files: []string{"skyrim_english.strings", "skyrim_japanese.strings"},
			want:  StringsPresenceView{English: true, Japanese: true},
		},
		{
			name:  "日本語欠け",
			files: []string{"skyrim_english.strings", "skyrim_english.dlstrings", "skyrim_english.ilstrings"},
			want:  StringsPresenceView{English: true},
		},
		{
			name:  "英語欠け（dlstrings だけでも日本語ありと判定する）",
			files: []string{"skyrim_japanese.dlstrings"},
			want:  StringsPresenceView{Japanese: true},
		},
		{
			name:  "大文字小文字の揺れを区別しない",
			files: []string{"Skyrim_English.STRINGS", "Skyrim_Japanese.ILStrings"},
			want:  StringsPresenceView{English: true, Japanese: true},
		},
		{
			name:  "他言語と非対象拡張子は数えない",
			files: []string{"skyrim_french.strings", "skyrim_japanese.txt", "readme.md"},
			want:  StringsPresenceView{},
		},
		{
			name:  "空",
			files: nil,
			want:  StringsPresenceView{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := classifyStringsPresence(tc.files); got != tc.want {
				t.Errorf("classifyStringsPresence(%v) = %+v, want %+v", tc.files, got, tc.want)
			}
		})
	}
}
