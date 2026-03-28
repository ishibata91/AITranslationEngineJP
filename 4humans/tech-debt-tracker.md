# Tech Debt Tracker

関連文書: [`../docs/index.md`](../docs/index.md), [`quality-score.md`](./quality-score.md), [`../docs/references/index.md`](../docs/references/index.md)

このファイルは、既知だが未解消の構造負債、設計負債、運用負債を集約する。

## Open Items

### 1. Legacy reference directory name

- Status: Open
- Area: Documentation structure
- Detail: 既存の API 参照ファイルは `docs/api-refrences/` にあり、ディレクトリ名に typo がある
- Impact: 新規参加者とエージェントが `references` と `api-refrences` の差を誤解しやすい
- Next step: `docs/references/` へ段階移行し、互換リンク方針を決めた上でリネームする

### 2. Test-level constraints are not established

- Status: Open
- Area: Design detail
- Detail: 細かな仕様と制約をテストや acceptance checks と結びつける方針は立てたが、実体はまだ少ない
- Impact: 実装判断の細部が文書だけでは固定しきれず、実行可能仕様が不足している
- Next step: plan の Acceptance Checks と将来の tests を対応づけ、詳細な制約を実行可能な形へ落とす

### 3. Design harness is rule-based, not semantic

- Status: Open
- Area: Harness
- Detail: 現在の設計ハーネスは主要キーワードと基本契約の確認に留まる
- Impact: 文書間の高度な矛盾や境界逸脱を自動検知できない
- Next step: 実装が進んだ段階で、用語辞書チェックや依存規約チェックを追加する

### 4. Execution harness does not cover translation-specific acceptance checks yet

- Status: Open
- Area: Verification
- Detail: Rust / Svelte の bootstrap target は追加されたが、翻訳ジョブ、importer、writer、persona 系の acceptance checks はまだ乗っていない
- Impact: execution harness は基盤 lint / test / build を回せるが、業務フロー成立の検証までは到達していない
- Next step: 実装進行に応じて fixture と acceptance checks を追加し、execution harness から段階的に呼ぶ

### 5. Test-level constraints coverage is not started

- Status: Open
- Area: Quality
- Detail: fixture-based checks、scenario regression、contract-level tests は方針化したが未実装
- Impact: 細かな仕様と制約を実行で確認する面がまだ弱い
- Next step: 実装追加時に最初の acceptance checks と test fixture を置く

### 6. Workflow skills are intentionally minimal

- Status: Open
- Area: Multi-agent workflow
- Detail: live workflow は `directing-implementation` / `directing-fixes` に寄せ直したが、旧 `.codex/.codex` 由来の不要 artifact や文言が作業ツリーに残る可能性がある
- Impact: live 正本と legacy 断片を取り違えると、古い packet 前提や review loop を再導入する危険がある
- Next step: live `.codex` に参照されない legacy artifact を段階的に整理し、lane 契約に関係ない古い references を残さない

