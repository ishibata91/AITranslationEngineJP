# AGENTS.md

## 目的

AITranslationEngineJp は、Skyrim Mod 向け翻訳エンジンを構築する `agent-first` リポジトリです。
最初のハーネス対象は開発基盤であり、目的は Codex などのエージェントがこのリポジトリで迷わず作業し、
同じ判断を再現し、変更後に自分で検証できる状態を作ることです。
このファイルは百科事典ではなく、最初に参照する地図として扱います。
作業方法と役割契約の正本は `.codex/` に置きます。

## 参照マップ

このファイルは入口だけを示す。必要がある場合に、以下を読む。

- 作業方法と役割契約の正本: `.codex/README.md`
- 実装の進め方: `.codex/skills/impl-direction/SKILL.md`
- 修正の進め方: `.codex/skills/fix-direction/SKILL.md`
- テスト設計の進め方: `.codex/skills/test-architect/SKILL.md`
- エージェントの役割契約: `.codex/agents/`
- エージェントの作業フロー: `.codex/skills/`
- 仕様の入口: `docs/index.md`
- 長期原則: `docs/core-beliefs.md`
- 恒久要件: `docs/spec.md`
- 内部境界と依存方向: `docs/architecture.md`
- 実装技術の選定: `docs/tech-selection.md`
- データ構造と ER: `docs/er-draft.md`
- 詳細な振る舞いと制約: 対応する tests / acceptance checks / validation commands
- 作業計画: `docs/exec-plans/` と `docs/exec-plans/templates/`

## 強い制約

- `.codex/` は、エージェントの作業方法と役割契約の正本とする
- `docs/` は、スコープ、アーキテクチャ、技術選定を記録する正本とする
- `4humans/` は、人間向けの品質状態と負債整理を記録する正本とする
- 用語は `docs/spec.md` の用語集に合わせる
- live workflow は `.codex/README.md` と `impl-direction` / `fix-direction` に従う
- 非自明な変更は、実装前に `docs/exec-plans/active/` へ計画を置く
- タスク完了後は計画を `docs/exec-plans/completed/` へ移し、結果を記録する
- 振る舞いが変わる変更では、関連する仕様文書や設計文書も同じ変更内で更新する
- 細かな仕様や制約は対応する test / acceptance checks / validation commands に寄せる
- 暗黙運用より、機械的に検証できる規約を優先する

## 実装前に確認すること

1. `.codex/README.md` と、対象作業に関係する direction skill を確認する
2. `docs/index.md` から、対象タスクに関係する文書だけを確認する
3. 既存の active / completed plan に同種タスクがないか確認する
4. 実装や設計内包タスクなら `docs/exec-plans/templates/impl-plan.md`、修正タスクなら `docs/exec-plans/templates/fix-plan.md` を使って計画を追加または更新する
5. `powershell -File scripts/harness/run.ps1 -Suite structure` を実行する
6. 文書契約や役割契約に触れるなら `powershell -File scripts/harness/run.ps1 -Suite design` も実行する

## 実装後にやること

1. `powershell -File scripts/harness/run.ps1 -Suite all` を実行する
2. 必要な文書、負債項目、品質スコアを更新する
3. タスク完了時は計画を `docs/exec-plans/completed/` へ移す

## 検証入口

- Structure harness: `powershell -File scripts/harness/run.ps1 -Suite structure`
- Design harness: `powershell -File scripts/harness/run.ps1 -Suite design`
- Execution harness: `powershell -File scripts/harness/run.ps1 -Suite execution`

## 作業スタイル

- 隠れた前提を増やさず、短く明示的な文書更新を優先する
- タスクが仕様変更を求めていない限り、既存仕様は不用意に書き換えない
- 新しいルールは短く、見つけやすく保つ
- 実装系の標準は `impl-direction -> impl-distill -> impl-workplan -> test-architect -> impl-work -> impl-review -> impl-direction close` とする
- 修正系の標準は `fix-direction -> fix-distill -> fix-trace -> (必要時 fix-logging / fix-analysis) -> test-architect -> fix-work -> fix-review -> fix-direction close` とする
- 過去 repo 由来で今の repo に合わない skill / agent / artifact 前提は、互換維持より削除を優先する
- 実装コードがまだ存在しない段階では、推測で public API を増やすより、ハーネスと文書を改善する
