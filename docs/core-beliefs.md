# Core Beliefs

関連文書: [`index.md`](./index.md), [`spec.md`](./spec.md), [`architecture.md`](./architecture.md), [`../4humans/quality-score.md`](../4humans/quality-score.md)

本プロジェクトは `agent-first` を強く採用するが、人間の責務は消えない。
人間は主に方針、受け入れ条件、境界、記録、ハーネス改善を担う。
作業方法と役割契約の正本は `.codex/` に置く。

## 1. 基本原則

- リポジトリは、エージェントが理解しやすい構造を優先する
- `AGENTS.md` は短い地図として保ち、作業方法は `.codex/`、プロダクト判断は `docs/` に置く
- `docs/` は説明資料ではなく、判断履歴と制約の正本として扱う
- `.codex/` は、`impl-direction` / `fix-direction` と補助 agent 契約の正本として扱う
- 非自明な変更は、実装前に短い計画を残す
- 仕様変更を伴う実装は、コード変更と同時に文書も更新する
- 繰り返し起こる失敗は、個別修正ではなくルールへ昇格させる
- 暗黙知より機械検証を優先する
- 品質は review の往復回数ではなく、plan / checks / evidence / harness の強さで担保する
- review は single-pass で `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足` を見る

## 2. 記録システムの扱い

- 永続要件は `spec.md` に記録する
- 内部境界と依存方向は `architecture.md` に記録する
- 実装技術の選択は `tech-selection.md` に記録する
- データモデルと ER は `er-draft.md` に記録する
- 実行可能仕様と制約は `executable-specs.md` と対応する test / acceptance checks に記録する
- 作業フローは `.codex/skills/` に記録する
- 役割契約は `.codex/agents/` に記録する
- 一時的な作業単位は `exec-plans/` に記録する
- 未整理の課題は `../4humans/tech-debt-tracker.md` に集約する
- 品質の現在地と不足は `../4humans/quality-score.md` に集約する

## 3. ルール化するべき失敗

次のものは、見つけた時点でリポジトリルールに引き上げる。

- 同じ前提説明を毎回要求する曖昧な文書構成
- 参照先が存在しない文書リンク
- 用語集と異なる名称の使用
- 同じ責務を複数文書で別定義している状態
- エージェントが検証入口を発見できない状態
- role handoff が曖昧で、lane と補助 skill の責務が崩れる状態
- review で繰り返し検出される契約違反が harness や executable specs へ昇格されない状態

## 4. 実装前後の標準動作

実装前:

- `AGENTS.md` から入り、`.codex/README.md` と relevant direction skill を読む
- 実装では `impl-direction` を使い、task-local design が要る時だけ active plan に `UI` / `Scenario` / `Logic` を埋める
- 修正では `fix-direction` を使い、事実不足なら trace と optional logging で scope を狭める
- 非自明な変更なら template を使って `docs/exec-plans/active/` に計画を置く
- 構造ハーネスを先に通す

実装後:

- `scripts/harness/run.ps1 -Suite all` を実行する
- 必要な記録を更新する
- 計画を `completed/` へ移す
- 最終 accept は Architect が持つ

## 5. 今はやらないこと

- まだ翻訳品質ゲートを repository gate にはしない
- まだ新しいコード上の public API を先回りして増やさない
- まだ巨大な単一マニュアルを作らない
