# generic-voice-tone-fallback（汎用ボイス NPC の話者未解決台詞への口調 fallback）

## 状態

backlog（未着手）。着手時に preparation-module で作業 branch と完了定義を固定する。本書は観測記録と暫定 goal だけを持つ。

## 依頼要約

話者連関（line_speaker）で特定 NPC に解決できない汎用ボイス NPC（衛兵など）の台詞へ、voice_type の気質 prior から口調を fallback で付ける。現状は口調生成が「話者解決済み」を前提とするため、汎用ボイスの台詞が「口調なし」になる。

## 背景（観測事実）

backend-violation-cleanup の実画面確認（`dictionaries/Data/Innocence Lost - Quest Expansion.esp`、実 LLM hy-mt2-7b）で観測した。

- 台詞 151 件のうち話者解決は 113 件で、これらには口調が付く。未解決は 38 件で口調なし。
  - DIAL（プレイヤー発話）26 件: 話者が NPC でないため口調なしが正しい。本 backlog の対象外。
  - INFO（NPC 応答）で未解決 12 件: すべて衛兵の汎用台詞（護送・逮捕・「他の衛兵の問題だ」など）。本 backlog の対象。
- voice fallback の規則は実在し、正常動作している。`internal/core/tone/classifier.go` の `fuseAttitude` が、本文の印が不足する話者へ voice 気質 prior（`internal/core/tone/voice_traits.go`）で対人軸を補う（DecisionPath="voice"）。
- ただし `classifier.Classify(feats, voice)` の voice は speaker に紐づく VoiceEDID から渡される。話者が解決できた台詞だけが対象で、話者未解決の台詞には届かない。
- `line` テーブルに voice カラムは無く、voice は speaker 経由でのみ得られる。衛兵台詞は line_speaker に話者が無いため voice も引けず、fallback の手前で口調なしになる。
- persona_character は 7 話者（Grelod・AventusAretino・ConstanceMichel・Samuel・Runa・Hroar・Francois）に限られ、衛兵は speaker 化されていない。

## goal（暫定）

話者未解決でも voice_type が判明する台詞へ、voice 気質 prior から口調を付ける。衛兵の汎用台詞に口調が付くことを実 app で確認する。プレイヤー発話（DIAL）は対象外に保つ。

## 未確定事項（着手時に確定する）

- 衛兵の voice_type が C# 抽出器で抽出され DB に入っているか。入っていなければ抽出器の改修要否を判断する。
- voice の渡し方: line に voice 列を持たせるか、汎用ボイスでも line_speaker を作るか、別経路で voice を注入するか。
- 公開境界（DTO、Store Port、Wails Bind）の変更要否。
- 対象範囲: 衛兵以外の汎用ボイス（種族ボイスなど）も含めるか。

## 起動条件

backend-violation-cleanup の merge 後に別途着手する。
