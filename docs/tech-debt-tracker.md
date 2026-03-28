# Tech Debt Tracker

関連文書: [`index.md`](./index.md), [`quality-score.md`](./quality-score.md), [`references/index.md`](./references/index.md)

このファイルは、既知だが未解消の構造負債、設計負債、運用負債を集約する。

## Open Items

### 1. Legacy reference directory name

- Status: Open
- Area: Documentation structure
- Detail: 既存の API 参照ファイルは `docs/api-refrences/` にあり、ディレクトリ名に typo がある
- Impact: 新規参加者とエージェントが `references` と `api-refrences` の差を誤解しやすい
- Next step: `docs/references/` へ段階移行し、互換リンク方針を決めた上でリネームする

### 2. Executable specs are not established

- Status: Open
- Area: Design detail
- Detail: 細かな仕様と制約をテストや acceptance checks と結びつける方針は立てたが、実体はまだ少ない
- Impact: 実装判断の細部が文書だけでは固定しきれず、実行可能仕様が不足している
- Next step: `docs/executable-specs.md` を起点に、plan の Acceptance Checks と将来の tests を対応づける

### 3. Design harness is rule-based, not semantic

- Status: Open
- Area: Harness
- Detail: 現在の設計ハーネスは主要キーワードと基本契約の確認に留まる
- Impact: 文書間の高度な矛盾や境界逸脱を自動検知できない
- Next step: 実装が進んだ段階で、用語辞書チェックや依存規約チェックを追加する

### 4. Execution harness has no code targets yet

- Status: Open
- Area: Verification
- Detail: 現時点では Rust / Svelte の実装ルートや test / lint / build 対象がまだ存在しない
- Impact: 実行ハーネスは標準入口のみ先行整備され、現在は主に skip を返す
- Next step: 実装追加時に package manager と cargo ルールを確定する

### 5. Test-level constraints coverage is not started

- Status: Open
- Area: Quality
- Detail: fixture-based checks、scenario regression、contract-level tests は方針化したが未実装
- Impact: 細かな仕様と制約を実行で確認する面がまだ弱い
- Next step: 実装追加時に最初の acceptance checks と test fixture を置く

### 6. Workflow skills are intentionally minimal

- Status: Open
- Area: Multi-agent workflow
- Detail: `.codex/skills/` は `architect-direction` と light flow を中心に最小構成へ寄せ、agent も Architect / Research / Coder の 3 役へ絞っている
- Impact: 複雑な lane 分割は減るが、必要な helper を再導入する際は current workflow に合わせて書き直す必要がある
- Next step: 実運用で本当に必要な helper だけを、Architect 起点の flow に沿って追加し直す
