# Quality Score

関連文書: [`../docs/index.md`](../docs/index.md), [`../docs/core-beliefs.md`](../docs/core-beliefs.md), [`tech-debt-tracker.md`](./tech-debt-tracker.md)

このファイルは、現在の品質状態を見える化し、不足している検証面を明示する。

## Scoring Guide

- `Green`: 最低限の運用契約と検証が成立している
- `Yellow`: 骨格はあるが、内容や自動化が不足している
- `Red`: 未整備、または検証入口が存在しない

## Current Areas

| Area | Score | Reason |
|---|---|---|
| Repository entrypoint | Green | `AGENTS.md` と `docs/index.md` で読む順序を固定した |
| Codex workflow source of truth | Green | `.codex/README.md` と workflow skills で plan stabilization loop、workflow gate、heavy/light reroute を固定した |
| Agent role contracts | Green | `.codex/agents/*.toml` で Architect / Research / Coder の責務境界を固定した |
| System of record | Green | 恒久契約は `docs/` と `.codex/` に固定し、task-local な詳細設計は plan、永続化すべき詳細は executable specs と tests / acceptance checks に分離した |
| Structure harness | Green | 必須ファイル、必須ディレクトリ、Markdown リンクの検査入口を追加した |
| Design harness | Green | workflow gate、unknown 分類、plan template 必須項目まで検査し、設計契約の欠落を検出できる |
| Execution harness | Yellow | 標準入口を追加したが、実装対象が未存在のため多くは skip になる |
| Executable specs and constraints | Yellow | evidence と reroute 契約は追加したが、対応する tests / fixtures / validation commands はまだ未整備 |

## Reserved Future Verification Tracks

- fixture-based execution checks
- scenario regression set
- contract-level tests

## Exit Criteria For Next Upgrade

- `docs/executable-specs.md` の主要項目が実際の tests / acceptance checks に対応づいている
- `.codex/` の workflow、workflow gate、role 契約が実運用で手戻りなく使える
- 実装コードが追加され、実行ハーネスが test / lint / build を実行できる
- 設計ハーネスが用語不整合、主要な境界逸脱、workflow gate 契約漏れを検出できる
