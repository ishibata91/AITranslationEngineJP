# AGENTS.md

## 目的

AITranslationEngineJp は、Skyrim Mod 向け翻訳エンジンを構築する `agent-first` リポジトリです。
このファイルは最初に参照する地図として扱います。
作業方法と役割契約の正本は `.codex/` に置きます。

## 参照マップ

このファイルは入口だけを示す。必要がある場合に、以下を読む。
まず最初に、これから使う skill の `references/permissions.json` を読み、その後に lane の `SKILL.md` と関連文書へ進む。

- 作業方法と役割契約の正本: `.codex/README.md`
- 実装の進め方: `.codex/skills/directing-implementation/SKILL.md`
- 修正の進め方: `.codex/skills/directing-fixes/SKILL.md`
- テスト設計の進め方: `.codex/skills/architecting-tests/SKILL.md`
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
- 実作業に入る前に、選択した skill の `references/permissions.json` を最優先で読む
- 各 skill は、自身の `references/permissions.json` に書かれた権限の範囲だけで動く
- skill の権限にないことはしない。権限解釈が曖昧な場合は停止し、適切な lane、skill、または human へ handoff する
- 用語は `docs/spec.md` の用語集に合わせる
- 非自明な変更は、実装前に `docs/exec-plans/active/` へ計画を置く
- タスク完了後は計画を `docs/exec-plans/completed/` へ移し、結果を記録する
- 振る舞いが変わる変更では、関連する仕様文書や設計文書も同じ変更内で更新する
- 暗黙運用より、機械的に検証できる規約を優先する
- できるだけ指示代名詞は使わない｡

## 実装前に確認すること

1. これから使う skill の `references/permissions.json` を最初に確認する
2. `docs/index.md` から、対象タスクに関係する文書だけを確認する
3. 既存の active / completed plan に同種タスクがないか確認する
4. `powershell -File scripts/harness/run.ps1 -Suite structure` を実行する
5. 文書契約や役割契約に触れるなら `powershell -File scripts/harness/run.ps1 -Suite design` も実行する

## 実装後にやること

1. `powershell -File scripts/harness/run.ps1 -Suite all` を実行する
2. 必要な文書、負債項目、品質スコアを更新する
3. タスク完了時は計画を `docs/exec-plans/completed/` へ移す

## 検証入口

- Structure harness: `powershell -File scripts/harness/run.ps1 -Suite structure`
- Design harness: `powershell -File scripts/harness/run.ps1 -Suite design`
- Execution harness: `powershell -File scripts/harness/run.ps1 -Suite execution`
