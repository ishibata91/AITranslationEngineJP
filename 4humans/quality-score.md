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
| Codex workflow source of truth | Green | `.codex/README.md` と workflow skills で `impl-direction` / `fix-direction`、embedded `UI` / `Scenario` / `Logic`、single-pass review を固定した |
| Agent role contracts | Green | `.codex/agents/*.toml` で distill、work planning、implementation、trace、logging、review の責務境界を固定した |
| System of record | Green | 恒久契約は `docs/` と `.codex/` に固定し、task-local な詳細設計は active plan、永続化すべき詳細は tests / acceptance checks / validation commands に分離した |
| Structure harness | Green | 必須ファイル、必須ディレクトリ、Markdown リンクの検査入口を追加した |
| Design harness | Green | lane workflow、embedded `UI` / `Scenario` / `Logic`、single-pass review 契約まで検査し、設計契約の欠落を検出できる |
| Execution harness | Yellow | Rust / frontend の実行対象を追加し、skip 状態は解消したが、Tauri 含む本格的な acceptance checks はまだ未整備 |
| Test-level constraints and checks | Yellow | bootstrap 用の lint / test / build / cargo command は接続したが、翻訳ドメイン固有の tests / fixtures / acceptance checks はこれから追加する |

## Reserved Future Verification Tracks

- fixture-based execution checks
- scenario regression set
- contract-level tests

## Exit Criteria For Next Upgrade

- 詳細な振る舞いと制約が実際の tests / acceptance checks / validation commands に対応づいている
- `.codex/` の workflow、single-pass review、role 契約が実運用で手戻りなく使える
- 実装コードが追加され、実行ハーネスが test / lint / build を実行できる
- 設計ハーネスが用語不整合、主要な境界逸脱、lane 契約漏れを検出できる
