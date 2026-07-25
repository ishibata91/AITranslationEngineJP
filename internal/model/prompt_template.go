package model

// PromptTemplate は翻訳 AI へ送る指示文の雛形。base 翻訳指示文と話者なし台詞の口調設定を持つ。
// 中心 DB の prompt_template テーブル（単一行 id=1）へ永続し、抽出データと別に保つ。
// 起動ごとの中心 DB 消去（抽出・翻訳を残さない方針）の対象から外し、編集したテンプレートを残す。
// 話者のいる台詞の口調雛形は directive テーブルの「口調」が唯一の正本のため、本型は持たない。
// db タグは db/migrations の prompt_template 列に対応する。
type PromptTemplate struct {
	// BaseDirective は叙述文・台詞のどちらにも付く冒頭の system 指示。
	BaseDirective string `db:"base_directive"`
	// GenericToneText は汎用台詞（話者なし）の口調指示文（自由記述）。tone_default テーブル由来。
	// 注入時、感情段階（本文 1 行）と条件由来の性別の一人称・語尾を自動で重ねる。
	GenericToneText string `db:"generic_tone_text"`
	// PcToneText は PC 発話（選択肢）の口調指示文（自由記述）。tone_default テーブル由来。
	// 注入時、感情段階と PcSex の一人称・語尾を自動で重ねる。
	PcToneText string `db:"pc_tone_text"`
	// PcSex は PC の性別（Female/Male/空）。PC 発話の一人称・語尾の根拠。tone_default テーブル由来。
	PcSex string `db:"pc_sex"`
}
