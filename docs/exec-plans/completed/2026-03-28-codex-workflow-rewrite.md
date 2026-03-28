# Codex Workflow Rewrite

- Date: 2026-03-28
- Status: Completed

## Goal

持ち込まれた旧 `.codex/` を、このリポジトリで使うマルチエージェント作業フローに合わせて書き換える。

## Scope

- `.codex/README.md` の追加
- Architect / light flow の入口 skill 追加
- `architect.toml` の追加
- 既存 direction / work / review skill を現在の workflow と `docs/` に合わせて更新
- `AGENTS.md` と harness を現行 `.codex/` 構成へ追従

## Acceptance Criteria

- User 要求の入口が `architect-direction` として定義される
- 低リスク修正の軽量フローが `light-direction -> light-work -> light-review` で定義される
- heavy flow が `Investigation -> Plan -> Coder -> Review` で読める
- `.codex` の主要 skill が現行 `docs/` を参照する
- harness が現在の `.codex` 構成を検査できる

## Outcome

- `.codex/` を旧構成のまま使うのではなく、上位入口を今回の workflow に合わせて再設計した
- 別プロジェクト由来の主要 docs 参照を、この repo の `docs/` へ差し替えた
- `AGENTS.md` と harness を現行 `.codex` に追従させた

## Follow-up

- 残る legacy skill の description と references を段階的に current repo 向けへ揃える
- 実運用の handoff 失敗を見て architect / light フローを微調整する
