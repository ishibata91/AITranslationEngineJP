# Tech Debt Tracker

関連文書: [`../docs/index.md`](../docs/index.md), [`quality-score.md`](./quality-score.md), [`../docs/references/index.md`](../docs/references/index.md)

このファイルは、既知だが未解消の構造負債、設計負債、運用負債を集約する。

## Open Items

### 1. Legacy reference directory name

- Status: Closed (2026-03-29)
- Area: Documentation structure
- Detail: ベンダー API 参照ファイルを `docs/references/vendor-api/` へ移し、typo のある `docs/api-refrences/` を廃止した
- Impact: 参照資料の導線が `docs/references/` に一本化され、`references` と `api-refrences` の混同が解消された
- Resolution: `docs/references/index.md` と `docs/index.md` を新配置へ同期した

### 2. Test-level constraints are not established

- Status: Open
- Area: Design detail
- Detail: 細かな仕様と制約をテストや acceptance checks と結びつける方針は立てたが、実体はまだ少ない
- Impact: 実装判断の細部が文書だけでは固定しきれず、実行可能仕様が不足している
- Next step: plan の Acceptance Checks と将来の tests を対応づけ、詳細な制約を実行可能な形へ落とす

### 3. Design harness is rule-based, not semantic

- Status: Closed (2026-03-29)
- Area: Harness
- Detail: 設計ハーネスに canonical phrase 整合チェックと architecture bootstrap boundary チェックを追加した
- Impact: 主要文書間で executable-spec の正本表現が崩れた場合や、frontend / backend の初期境界が設計契約から逸脱した場合を自動検知できる
- Resolution: `scripts/harness/check-design.ps1` で semantic checks を実行し、`docs/core-beliefs.md` と `4humans/quality-score.md` を同期した

### 4. Execution harness does not cover translation-specific acceptance checks yet

- Status: Open
- Area: Verification
- Detail: Rust / Svelte の bootstrap target に加え xEdit importer validation の Rust tests は入ったが、翻訳ジョブ、writer、persona 系の acceptance checks はまだ乗っていない
- Impact: execution harness は基盤 lint / test / build と importer failure path の一部を回せるが、業務フロー成立の検証までは到達していない
- Next step: 実装進行に応じて importer 以外の fixture と acceptance checks を追加し、execution harness から段階的に呼ぶ

### 5. Test-level constraints coverage is not started

- Status: Open
- Area: Quality
- Detail: fixture-based checks、scenario regression、contract-level tests は方針化済みで、xEdit importer validation と canonical `TRANSLATION_UNIT` contract から着手したが coverage はまだ限定的
- Impact: importer と translation-unit canonical boundary の一部制約は実行で確認できるようになった一方、細かな仕様と制約を広く確認する面はまだ弱い
- Next step: 実装追加時に translation flow 全体へ acceptance checks と test fixture を広げる

### 6. Workflow skills are intentionally minimal

- Status: Open
- Area: Multi-agent workflow
- Detail: live workflow は `directing-implementation` / `directing-fixes` に寄せ直したが、旧 `.codex/.codex` 由来の不要 artifact や文言が作業ツリーに残る可能性がある
- Impact: live 正本と legacy 断片を取り違えると、古い packet 前提や review loop を再導入する危険がある
- Next step: live `.codex` に参照されない legacy artifact を段階的に整理し、lane 契約に関係ない古い references を残さない

### 7. Execution-cache SQLite path is still temporary

- Status: Open
- Area: Persistence
- Detail: `PER-41` で xEdit import の execution cache を SQLite へ保存できるようにしたが、command 境界の既定 path は `%TEMP%/ai-translation-engine-jp-execution-cache.sqlite` のままで、正式な app-owned cache location と retention policy は未固定である
- Impact: 実行キャッシュの寿命と観測先が OS の temp 運用に依存し、後続の job-linking や cleanup 方針を固める前提がまだ弱い
- Next step: app-owned execution-cache location、retention policy、job cleanup と `PLUGIN_EXPORT` 削除条件を後続 persistence task で固定する

### 8. xTranslator SST compatibility is still fixture-scoped

- Status: Open
- Area: Dictionary foundation
- Detail: `P2-I01` で xTranslator dictionary importer を追加し、shared reusable-entry snapshot と whitespace-sensitive fixture を通す最小 SST layout は検証したが、metadata 優先規則、旧 version 互換、追加 field semantics はまだ固定していない
- Impact: 現状の importer は Phase 2 foundation-data の最小 ingest path には使える一方、広い xTranslator SST 実データ互換を前提にすると解釈差分が残る
- Next step: foundation-data 後続 task で実データ fixture を追加し、`dictionaryName` metadata precedence と version 互換の acceptance checks を固定する
