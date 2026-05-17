# Scenario Candidates: 2026-05-18-storybook-foundation / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `STORYBOOK-AUDIT`

## Generator Scope

- `viewpoint`: `operation-audit`
- `included_sources`: `plan.md`, `docs/tech-selection.md`, `docs/coding-guidelines-frontend.md`, `docs/coding-guidelines-tests.md`, `docs/lint-policy.md`, `docs/frontend-fake-api.md`, `frontend/package.json`, `frontend/vite.config.ts`, `docs/spec.md`, `docs/architecture.md`, `docs/er.md`, `docs/screen-design/README.md`, `docs/detail-specs/README.md`, `docs/scenario-tests/README.md`
- `excluded_sources`: プロダクトコード変更、プロダクトテスト変更、docs 正本化、`.codex/` 変更、Master Persona 部品化、画面再設計、gateway mock、backend DTO mock、Storybook での実行フロー再現
- `generation_notes`: Storybook 最小基盤で後から確認できるべき検証、review URL、fixture 方針、後続 task への証跡だけを候補化する。採否、統合、最終シナリオ ID は designer に残す。

## Candidate Scenarios

### CAND-STORYBOOK-AUDIT-001 Storybook 検証コマンド結果を後から確認できる

- `source requirement`: `plan.md` は Storybook dev / build の入口、`python3 scripts/harness/run.py --suite frontend-local`、`npm --prefix frontend run build-storybook` を検証対象にしている。`lint-policy.md` は検証 command と失敗条件を acceptance checks 側へ残す責務にしている。
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-STORYBOOK-AUDIT-001`
- `actor`: 実装担当者、レビュー担当者、merge 担当者
- `trigger`: Storybook 最小基盤の実装後に検証を実行する。
- `expected outcome`: 実行した command、通過または失敗、未実行理由、失敗時の要約が task 成果物に残る。
- `audit event`: Storybook dev 起動確認、Storybook build、frontend-local の検証結果記録。
- `stored summary`: command 名、実行場所、結果、失敗時の短い原因、再実行 command。
- `redaction rule`: command 出力全文、依存 cache の絶対 path、secret、token、ローカル環境固有の長い path は保存しない。
- `observable point`: `plan.md` または実装結果 artifact に検証結果があり、再実行 command が `frontend/package.json` の script と矛盾しない。
- `related detail requirement type`: `observability_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: Storybook build を `frontend-local` に含めるか別 gate にするかが未決でも、実行結果の記録シナリオとして採用しやすい。
- `conflict hint`: `frontend-local` へ Storybook build を含める判断と、別 gate に分ける判断が競合しうる。designer が検証段階を確定する。

### CAND-STORYBOOK-AUDIT-002 Storybook review URL と確認状態を後から追跡できる

- `source requirement`: `plan.md` は review URL 記録方針を Storybook 基盤の目的に含めている。`docs/frontend-fake-api.md` は旧 review URL 運用で実行 URL、確認状態、未確認理由を task 成果物へ残すとしている。
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-STORYBOOK-AUDIT-002`
- `actor`: 実装担当者、人間レビュー担当者
- `trigger`: Storybook dev server で最小 story または空確認用 story を開く。
- `expected outcome`: review URL、story ID、確認した状態、未確認状態、未確認理由が task 成果物に残る。
- `audit event`: Storybook review URL の提示と人間確認対象の固定。
- `stored summary`: localhost URL、story ID、確認状態、未確認状態、確認者へ渡す対象成果物。
- `redaction rule`: URL query に secret、API key、token、ローカル絶対 path、実ユーザーデータを含めない。
- `observable point`: task 成果物に Storybook review URL があり、後続 Master Persona task が同じ記録形式を参照できる。
- `related detail requirement type`: `observability_requirement`, `testability_requirement`
- `adoption hint`: fakeAPI の review URL 運用を Storybook へ置き換える監査候補として使える。
- `conflict hint`: review URL を `ui-design.md` に置くか、専用 `storybook-review.md` に置くかは未決である。designer が成果物置き場を確定する。

### CAND-STORYBOOK-AUDIT-003 fixture 方針が後続 task で監査できる

- `source requirement`: `plan.md` は fixture 配置方針を最小単位で固定するとしている。`coding-guidelines-frontend.md` と `architecture.md` は UI Component が backend DTO、generated binding、Gateway、Store を直接扱わないと定めている。`docs/frontend-fake-api.md` は API key、token、secret を mock data や画面表示へ入れないとしている。
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-STORYBOOK-AUDIT-003`
- `actor`: 後続 Storybook story 追加担当者、レビュー担当者
- `trigger`: 最小 story と fixture 配置方針を追加する。
- `expected outcome`: fixture の置き場所、fixture が表す状態、Storybook に入れてよいデータ境界が task 成果物に残る。
- `audit event`: Storybook fixture 方針の固定と禁止保存対象の確認。
- `stored summary`: fixture root、story が使う props の概要、サンプルデータの由来、禁止した mock 種別。
- `redaction rule`: secret、API key、token、外部 provider 応答原文、実翻訳本文の過剰な原文、ローカル絶対 path を fixture に保存しない。
- `observable point`: story または fixture が generated `wailsjs`、backend DTO、gateway mock、実行フロー再現を前提にしていないことを task 成果物から確認できる。
- `related detail requirement type`: `data_requirement`, `security_requirement`, `observability_requirement`, `testability_requirement`
- `adoption hint`: fixture 置き場が未決でも、保存してよい要約と保存禁止対象の候補として使える。
- `conflict hint`: fixture を `frontend/src/ui/**/__fixtures__` に置くか、Storybook 専用 directory に置くかが競合しうる。designer が採用先を決める。

### CAND-STORYBOOK-AUDIT-004 後続 Master Persona task が基盤完了証跡を再利用できる

- `source requirement`: `plan.md` は Storybook 基盤完了後に Master Persona task が story を追加できることを close condition にしている。`2026-05-17-master-persona-componentization/plan.md` は Storybook review URL と build-storybook の記録を後続 task の前提にしている。
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-STORYBOOK-AUDIT-004`
- `actor`: 後続 Master Persona task の実装担当者、レビュー担当者
- `trigger`: Storybook foundation task を close し、後続 task へ渡す。
- `expected outcome`: 後続 task が参照する Storybook scripts、config、最小 story、build 検証結果、review URL 記録形式が task 成果物に残る。
- `audit event`: Storybook 基盤の closeout 証跡と後続 task への引き継ぎ。
- `stored summary`: script 名、config path、最小 story path、build 検証結果、review URL 記録形式、残留リスク。
- `redaction rule`: build artifact の全一覧、依存 package の長いログ、ローカル絶対 path、secret を closeout 証跡に保存しない。
- `observable point`: 後続 Master Persona task が `2026-05-18-storybook-foundation` の完了成果物を参照して、追加 story の検証入口を特定できる。
- `related detail requirement type`: `observability_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 後続 task への監査可能性を保証する候補として使える。
- `conflict hint`: Storybook 運用をどの docs 正本へ反映するかは未決である。候補段階では task-local closeout 証跡だけを扱う。

## Open Notes

- `human decision candidate`: Storybook build を `frontend-local` に含めるか、別 gate にするか。
- `human decision candidate`: fixture を `frontend/src/ui/**/__fixtures__` に置くか、Storybook 専用 directory に置くか。
- `human decision candidate`: review URL 記録を `ui-design.md` に置くか、専用 `storybook-review.md` に置くか。
- `human decision candidate`: Storybook 運用をどの docs 正本へ反映するか。
- `merge candidate`: CAND-STORYBOOK-AUDIT-001 と CAND-STORYBOOK-AUDIT-004 は、closeout 証跡シナリオとして統合できる可能性がある。
- `merge candidate`: CAND-STORYBOOK-AUDIT-002 と後続 Master Persona の Storybook review 候補は、review URL 記録形式として統合できる可能性がある。
- `rejection candidate`: 実行フロー再現、gateway mock、backend DTO mock を Storybook に入れる候補は、この task の禁止事項により不採用候補である。
