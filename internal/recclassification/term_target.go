// Package recclassification は単語翻訳対象 REC の共通 config を公開する。
// 翻訳対象 REC 集合はこの package が単一情報源であり、
// service 層や repository 層の SQL IN 句はここから生成した値を使う。
package recclassification

// TermTargets は単語翻訳対象 13 種別の REC 集合を定義順で保持する。
// 形式は "RECORD:FIELD"（例: "NPC_:FULL"）。
// DOOR:FULL、FLOR:FULL、FURN:FULL は対象外であり、この集合に含まれない。
var TermTargets = []string{
	"BOOK:FULL",
	"NPC_:FULL",
	"NPC_:SHRT",
	"ARMO:FULL",
	"WEAP:FULL",
	"LCTN:FULL",
	"CELL:FULL",
	"CONT:FULL",
	"MISC:FULL",
	"ALCH:FULL",
	"RACE:FULL",
	"INGR:FULL",
	"SHOU:FULL",
}

// termTargetSet は IsTermTarget の O(1) 判定用の内部集合。
// TermTargets と同期して初期化する。
var termTargetSet map[string]struct{}

func init() {
	termTargetSet = make(map[string]struct{}, len(TermTargets))
	for _, rec := range TermTargets {
		termTargetSet[rec] = struct{}{}
	}
}

// IsTermTarget は rec が単語翻訳対象 13 種別のいずれかであれば true を返す。
// rec が空文字または 13 種別以外（例: "DOOR:FULL"、"FLOR:FULL"、"FURN:FULL"、未知）の場合は false を返す。
func IsTermTarget(rec string) bool {
	_, ok := termTargetSet[rec]
	return ok
}

// TermTargetRECList は SQL IN 句への埋め込みを目的として、
// 単語翻訳対象 13 種別を定義順で返す。
// SQL 側に REC リテラルを直書きする二重管理を防ぐため、
// repository 層はこの関数から取得した値を使って IN 句を組み立てる。
func TermTargetRECList() []string {
	result := make([]string, len(TermTargets))
	copy(result, TermTargets)
	return result
}
