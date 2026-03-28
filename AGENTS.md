# AGENTS.md

## 目的

AITranslationEngineJp は、Skyrim Mod 向け翻訳エンジンを構築する `agent-first` リポジトリです。
最初のハーネス対象は開発基盤であり、目的は Codex などのエージェントがこのリポジトリで迷わず作業し、
同じ判断を再現し、変更後に自分で検証できる状態を作ることです。
このファイルは百科事典ではなく、最初に参照する地図として扱います。
作業方法と役割契約の正本は `.codex/` に置きます。

## 最初に読む順序

1. `.codex/README.md`
2. `.codex/skills/architect-direction/SKILL.md`
3. `.codex/agents/architect.toml`
4. 必要なら `.codex/agents/research.toml` または `.codex/agents/coder.toml`
5. `docs/index.md`
6. `docs/core-beliefs.md`
7. `docs/spec.md`
8. `docs/architecture.md`
9. `docs/tech-selection.md`
10. `docs/er-draft.md`
11. `docs/executable-specs.md`
12. 必要なら `docs/exec-plans/` 配下の計画

## 強い制約

- `.codex/` は、エージェントの作業方法と役割契約の正本とする
- `docs/` は、スコープ、アーキテクチャ、技術選定、実行可能仕様を記録する正本とする
- `4humans/` は、人間向けの品質状態と負債整理を記録する正本とする
- 用語は `docs/spec.md` の用語集に合わせる
- heavy / light の判定は `.codex/README.md` と `architect-direction` / `light-direction` に従う
- 非自明な変更は、実装前に `docs/exec-plans/active/` へ計画を置く
- タスク完了後は計画を `docs/exec-plans/completed/` へ移し、結果を記録する
- 振る舞いが変わる変更では、関連する仕様文書や設計文書も同じ変更内で更新する
- 細かな仕様や制約は `docs/executable-specs.md` と対応する test / acceptance checks に寄せる
- Architect が最終レビュー責任を持つ
- エージェントが繰り返し迷うなら、個別修正で終わらせず `.codex/` か `docs/` にルールを昇格させる
- 暗黙運用より、機械的に検証できる規約を優先する

## 実装前にやること

1. `.codex/README.md` と relevant agent / direction skill を読む
2. `docs/index.md` と対象タスクに関係する設計文書を読む
3. 既存の active / completed plan に同種タスクがないか確認する
4. heavy なら `docs/exec-plans/templates/heavy-plan.md`、light なら `docs/exec-plans/templates/light-plan.md` を使って計画を追加または更新する
5. `powershell -File scripts/harness/run.ps1 -Suite structure` を実行する
6. 文書契約や役割契約に触れるなら `powershell -File scripts/harness/run.ps1 -Suite design` も実行する

## 実装後にやること

1. `powershell -File scripts/harness/run.ps1 -Suite all` を実行する
2. 必要な文書、負債項目、品質スコアを更新する
3. タスク完了時は計画を `docs/exec-plans/completed/` へ移す

## 何をどこに記録するか

- 恒久要件: `docs/spec.md`
- 内部境界と依存方向: `docs/architecture.md`
- 実装技術の選定: `docs/tech-selection.md`
- データ構造と ER: `docs/er-draft.md`
- 実行可能仕様と制約: `docs/executable-specs.md`
- エージェントの役割契約: `.codex/agents/`
- エージェントの作業フロー: `.codex/skills/`
- 作業計画: `docs/exec-plans/`
- 長期的な原則とガードレール: `docs/core-beliefs.md`
- 未解消の課題や曖昧さ: `4humans/tech-debt-tracker.md`
- 現在の品質状態と不足: `4humans/quality-score.md`
- 外部仕様やベンダー資料: `docs/references/`

## 検証入口

- Structure harness: `powershell -File scripts/harness/run.ps1 -Suite structure`
- Design harness: `powershell -File scripts/harness/run.ps1 -Suite design`
- Execution harness: `powershell -File scripts/harness/run.ps1 -Suite execution`

## 作業スタイル

- 隠れた前提を増やさず、短く明示的な文書更新を優先する
- タスクが仕様変更を求めていない限り、既存仕様は不用意に書き換えない
- 新しいルールは短く、見つけやすく保つ
- heavy では `Architect -> Research -> Plan -> Coder -> Architect review` を標準とする
- light では `Architect -> Short plan -> Coder -> Architect review` を標準とする
- 実装コードがまだ存在しない段階では、推測で public API を増やすより、ハーネスと文書を改善する
