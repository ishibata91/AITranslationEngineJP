---
name: distill
description: Codex 側の共通文脈圧縮 skill。入口情報を facts、constraints、gaps、required_reading へ整理するための共通知識と圧縮パターンを提供する。
---
# Distill

## 目的

`distill` は、入口情報を短く整理するための共通知識である。
圧縮粒度、重複除去、facts / inferred / gap の分離、downstream が読む順番の作り方を扱う。

設計向けの観点は `distill-design` が持つ。
調査向けの観点は `distill-investigate` が持つ。
`distill` 本体は、どちらにも共通する圧縮の見方だけを持つ。

## 対応ロール

- `distiller` が使う。
- 返却先は caller または次 agent とする。
- owner artifact は `distill` の出力規約で固定する。

## 入力規約

- `implement-lane` の次判断に必要な repo 文脈を圧縮する時
- 設計向けまたは調査向けの詳細観点を読む前に、共通の圧縮粒度をそろえる時
- user request、active plan、docs、関連 skill の重複を短く整理する時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。
- 必須入力: caller, user_instruction_or_task_summary, entrance_artifacts, lane_owner
- 任意入力: active_work_plan, candidate_source_paths, known_gaps, knowledge_focus
- input_notes: {"entrance_artifacts": "次の grouped input を想定する: task_frame、canonical_evidence、code_evidence、effective_prior_decisions、observation_evidence。caller は利用可能なものだけを渡し、欠ける場合は path catalog と短い note で代替してよい。effective_prior_decisions は active / completed plan の検索と object 抽出で作り、decision、rationale、source_path、source_lines を持つ。未確認または失効した判断は採用しない。", "knowledge_focus": "設計向けまたは調査向けの知識参照を選ぶための任意ヒント。未指定なら caller context から判断し、迷う場合は両方の focused skill を読む。"}
- 必須 artifact: user の現在指示または task summary, active work plan または caller-provided task context, /Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md

## 外部参照規約

- agent runtime と tool policy は [distiller.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/distiller.toml) の `allowed_write_paths` / `allowed_commands` とする。
- binding: [distiller.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/distiller.toml)
- agent runtime: [distiller.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/distiller.toml)
- tool policy: agent runtime の `allowed_write_paths` / `allowed_commands` に従う
- 圧縮判断: [compression-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill/references/compression-patterns.md)
- 設計向け観点: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill-design/SKILL.md)
- 調査向け観点: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill-investigate/SKILL.md)
- binding: [distiller.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/distiller.toml)
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。
- common: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill/SKILL.md
- focused_optional: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill-design/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill-investigate/SKILL.md

## 内部参照規約

### 拘束観点

- `catalog`、`summary`、`full` の圧縮粒度
- canonical source への重複寄せ
- `confirmed`、`inferred`、`gap` の分離
- active / completed plan から有効な過去判断だけを抽出すること
- downstream が読む順番の整理

### 圧縮方針

- 先に対象を機械的に棚卸しし、その後で判断する
- 読む粒度は `catalog`、`summary`、`full` の順で上げる
- 重複 instruction は正本だけを残し、重複元は path で退避する
- downstream の次判断に必要な情報だけを残す
- 出力ごとに `confirmed`、`inferred`、`gap` の状態を明示する

## 判断規約

- 先に対象を機械的に棚卸しし、その後で判断する
- 読む粒度は `catalog`、`summary`、`full` の順で上げる
- 重複 instruction は正本だけを残し、重複元は path で退避する
- downstream の次判断に必要な情報だけを残す
- 出力ごとに `confirmed`、`inferred`、`gap` の状態を明示する

- 確認済み事実と推測を分ける
- 重要な fact には根拠 path を付ける
- 過去判断は有効なものだけを残し、出典 line を付ける
- 必要な参照先を読む順番つきで返す

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力に tool policy、agent runtime、product code の変更義務を含めない。

### Handoff

- handoff 先: `implement_lane`
- 渡す scope: 次の設計または調査を判断するための圧縮済み facts と gaps
- 必須出力: facts, constraints, gaps, required_reading, related_design_pointers, observation_targets, recommended_next_skill
- 出力 field 要件: {"facts": "confirmed / inferred / gap の状態と根拠 path を含める。過去判断を含める時は、有効な判断だけを採用する", "constraints": "正本 path と、重複元があれば source note を含める", "gaps": "未確認事項を事実として混ぜずに列挙する。有効な過去判断が見つからないことだけでは gap にしない", "required_reading": "downstream が読む順番が分かる形で返す", "related_design_pointers": "requirements、UI、scenario に関係する path があれば返す。関係しない場合は空にする", "observation_targets": "観測対象、入口、未観測情報があれば返す。関係しない場合は空にする", "recommended_next_skill": "implement-lane が次に呼ぶ skill または停止理由を返す"}

## 完了規約

- 出力規約を満たし、次の actor が再解釈なしで判断できる。
- 不足情報または停止理由がある場合は明示されている。
- 圧縮粒度が `catalog`、`summary`、`full` のどれかとして説明できる。
- 重要な fact に根拠 path がある。
- `confirmed`、`inferred`、`gap` が混ざっていない。
- 必須 evidence: 重要な fact の根拠 path, 採用した有効な過去判断の source_path と source_lines, 参照した正本 path, 未確認事項の理由
- completion signal: implement-lane が次の設計または調査を判断できる
- residual risk key: gaps

## 停止規約

- requirements、UI、scenario、diagram の詳細観点だけが必要な時
- 観測対象、再現条件、未観測情報の詳細観点だけが必要な時
- human review 済み `implementation-scope` から実装前 context を作る時
- fix、refactor、product code 実装のために文脈を整理する時
- active work plan や関連 docs が不足している場合は停止する。
- 重要な fact の根拠 path を確認できない場合は停止する。
- 主要な設計判断が未確定で事実整理だけでは前進しない場合は `implement_lane` へ戻す。
- 実装前の文脈整理が目的なら、Codex implementation lane [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-distill/SKILL.md) を使う前提で `implement_lane` へ戻す。
- 設計向けまたは調査向けの詳細観点を共通 skill に戻さない
- broad な repo tour をしない
- product code、product test、docs 正本の変更に進まない
- 停止時は不足項目、衝突箇所、reroute 先を返す。
- focused skill が持つ詳細観点を共通 skill に戻していない場合は停止する。
- implementation-scope 承認後の実装前整理を扱っていない場合は停止する。
- product code、product test、docs 正本の変更に進んでいない場合は停止する。
- 拒否条件: implementation-scope 承認後の実装前整理である
- 拒否条件: caller が implement-lane ではなく実装 lane である
- 拒否条件: source artifact が不足し、根拠 path を確認できない
- 拒否条件: 圧縮ではなく設計判断、調査実行、実装が求められている
- 停止条件: active work plan や関連 docs が不足している
- 停止条件: 重要な fact の根拠 path を確認できない
- 停止条件: 主要な設計判断が未確定で事実整理だけでは前進しない
- 停止条件: 作業が Codex implementation lane implementation-distill の責務である
