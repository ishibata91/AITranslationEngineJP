# Scenario Candidates: 2026-05-17-all-pages-componentization / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `APCOA`

## Generator Scope

- `viewpoint`: operation-audit
- `included_sources`: `./plan.md`, `docs/architecture.md`, `docs/coding-guidelines-frontend.md`, `docs/UX-standard.md`, completed Storybook foundation scenario, completed Master Persona componentization scenario and UI design, `frontend/src/ui/screens/`, `frontend/src/ui/components/`
- `excluded_sources`: product code changes, product tests, docs canonicalization, final scenario adoption, candidate integration
- `generation_notes`: 全ページ部品化と Storybook story 追加について、後追い確認、検証履歴、review 証跡、保存禁止情報を候補化する。採否と統合は `designer` に残す。

## Candidate Scenarios

### CAND-APCOA-001 Storybook review 対象と確認状態を追跡する

- `source requirement`: plan の close conditions。各ページの主要表示領域が props 境界と Storybook story を持ち、実装後の Storybook review が通過または未通過理由を持つ。
- `viewpoint`: 後追い確認 / 履歴
- `candidate scenario id`: `CAND-APCOA-001`
- `actor`: 人間レビュアーまたは実装後確認者
- `trigger`: 全ページ部品化後に Storybook review を開始する。
- `expected outcome`: story ID、review URL、確認状態、未確認理由、再実行 command が task-local 成果物に残る。
- `audit event`: Storybook 上の対象 story を開き、表示確認を完了または未確認にする事象。
- `saved summary`: 対象画面、component 名、story ID、Storybook URL、確認状態、未確認理由、確認日時相当の記録。
- `redaction rule`: URL query、確認メモ、fixture 名に secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータ、raw request、raw response、raw prompt を残さない。
- `observable point`: `storybook-review.md` または `ui-design.md#storybook-review` に残る review URL、story ID、確認状態、未確認理由。
- `related detail requirement type`: `observability_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: 最終シナリオでは Storybook review 証跡の必須項目として統合できる。
- `conflict hint`: review 証跡の保存先が `storybook-review.md` と `ui-design.md#storybook-review` で揺れる可能性がある。

### CAND-APCOA-002 `frontend-local` の検証結果を後から確認できる

- `source requirement`: plan の validation commands。`python3 scripts/harness/run.py --suite frontend-local` を検証入口にする。
- `viewpoint`: 監査ログ / 再現材料
- `candidate scenario id`: `CAND-APCOA-002`
- `actor`: 実装者または merge 前確認者
- `trigger`: 全ページ部品化と story 追加後に frontend-local を実行する。
- `expected outcome`: command、結果、失敗時の未通過理由、再実行条件が task-local 成果物に残る。
- `audit event`: frontend-local を実行し、通過または未通過を確定する事象。
- `saved summary`: 実行 command、exit 状態、失敗時の要約、未実行理由、再実行に必要な前提。
- `redaction rule`: command output から secret、token、環境変数の実値、実ユーザーデータ、不要な長文ログを保存しない。
- `observable point`: task-local closeout 証跡または review 証跡の frontend-local 結果欄。
- `related detail requirement type`: `observability_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: 最終シナリオでは実装後の lower-level 検証証跡として扱える。
- `conflict hint`: `frontend-local` の失敗を UI review の未通過と同一視すると、原因調査の入口が曖昧になる。

### CAND-APCOA-003 `build-storybook` の静的 build 証跡を残す

- `source requirement`: plan の validation commands と Storybook foundation の固定要件。`npm --prefix frontend run build-storybook` は backend、Wails runtime、AI provider、secret store、DB に依存しない。
- `viewpoint`: 監査ログ / 再現材料
- `candidate scenario id`: `CAND-APCOA-003`
- `actor`: 実装者または merge 前確認者
- `trigger`: Storybook story 追加後に static build を実行する。
- `expected outcome`: build command、結果、生成物の扱い、失敗時の未通過理由が task-local 成果物に残る。
- `audit event`: Storybook static build を実行し、通過または未通過を確定する事象。
- `saved summary`: 実行 command、exit 状態、生成物の有無、失敗時の要約、backend 非依存であることの確認。
- `redaction rule`: build log へ secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータを残さない。
- `observable point`: task-local closeout 証跡または `storybook-review.md` の build-storybook 結果欄。
- `related detail requirement type`: `observability_requirement`, `testability_requirement`, `security_requirement`
- `adoption hint`: 最終シナリオでは Storybook 専用 gate の証跡として扱える。
- `conflict hint`: app build の生成物と Storybook build の生成物を同じ証跡として扱うと、確認対象がずれる。

### CAND-APCOA-004 部品化変更証跡で props 境界を追跡する

- `source requirement`: plan の goal と Pre-design Investigation。全ページの主要表示領域をパネル、カード、モーダル単位へ部品化し、共有部品と画面専用部品を分類する。
- `viewpoint`: 履歴 / 再現材料
- `candidate scenario id`: `CAND-APCOA-004`
- `actor`: 実装者、reviewer、designer
- `trigger`: 画面 page component を薄くし、screen local component または shared component へ分ける。
- `expected outcome`: 変更された画面、component、story、fixture、分けない候補の理由が task-local 成果物に残る。
- `audit event`: ページ表示領域を component 境界へ移動し、story 対象を確定する事象。
- `saved summary`: 画面名、component 名、配置先、story path、fixture path、共有化理由、画面専用に残す理由、分けない理由。
- `redaction rule`: fixture と変更メモに実ユーザーデータ、実ファイル内容、raw prompt、raw response、provider 応答原文を残さない。
- `observable point`: component candidate 証跡、implementation result、storybook-review の story 対象一覧。
- `related detail requirement type`: `observability_requirement`, `compatibility_requirement`, `data_requirement`
- `adoption hint`: 最終シナリオでは component 境界と Storybook story の対応確認として統合できる。
- `conflict hint`: 共有 component へ上げる基準が未確定のまま証跡化されると、後続 review で責務境界の採否判断が必要になる。

### CAND-APCOA-005 人間 review と UX review の証跡を残す

- `source requirement`: plan の HITL Status。ux review と frontend human review は実装後に必要である。
- `viewpoint`: 後追い確認 / 履歴
- `candidate scenario id`: `CAND-APCOA-005`
- `actor`: UX reviewer、人間レビュアー、implement_lane
- `trigger`: Storybook review、UX review、frontend human review を実装後に実施する。
- `expected outcome`: review 対象、review 結果、未確認理由、承認状態、残留リスクが task-local 成果物に残る。
- `audit event`: 人間または UX review が Storybook 表示と検証結果を確認し、承認または未承認を記録する事象。
- `saved summary`: review 種別、対象 story、確認結果、指摘、未確認理由、承認状態、再確認 command。
- `redaction rule`: review メモに secret、token、実ユーザーデータ、外部 provider 応答原文、過剰な本文引用を残さない。
- `observable point`: `ux-review` 成果物、frontend human review 証跡、`storybook-review.md` の確認結果。
- `related detail requirement type`: `observability_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: 最終シナリオでは merge readiness の review evidence 条件として統合できる。
- `conflict hint`: UX review と frontend human review の役割を混ぜると、承認対象と見た目確認対象が曖昧になる。

### CAND-APCOA-006 Storybook fixture と review 証跡の保存禁止情報を監査する

- `source requirement`: Storybook foundation と Master Persona componentization の固定要件。story と fixture は fixed props または view model fixture で表示し、secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータ、raw request、raw response、raw prompt を含めない。
- `viewpoint`: 保存禁止 / 競合候補
- `candidate scenario id`: `CAND-APCOA-006`
- `actor`: 実装者または reviewer
- `trigger`: story、fixture、review 証跡、検証ログを追加または更新する。
- `expected outcome`: 保存禁止情報が story、fixture、review 証跡、検証結果に含まれていないことを後から確認できる。
- `audit event`: Storybook 関連成果物と review 証跡の保存内容を確認する事象。
- `saved summary`: 確認対象 path、確認した禁止情報カテゴリ、検出結果、検出時の扱い。
- `redaction rule`: secret、API key、token、credential 実値、実 endpoint、ローカル絶対 path、実ユーザーデータ、raw request、raw response、raw prompt、provider 応答原文は保存しない。
- `observable point`: story source、fixture source、`storybook-review.md`、検証結果要約。
- `related detail requirement type`: `security_requirement`, `data_requirement`, `observability_requirement`
- `adoption hint`: 最終シナリオでは保存禁止情報の回帰確認として扱える。
- `conflict hint`: 保存対象を広げる候補は `security_requirement` または `data_requirement` と衝突する可能性がある。

## Open Notes

- `human decision candidate`: Storybook review 証跡の正式保存先が `storybook-review.md` か `ui-design.md#storybook-review` かは最終設計で固定する必要がある。
- `human decision candidate`: 全ページの review 対象 story をどこまで必須にするかは、component candidate と UI design の確定後に判断する必要がある。
- `merge candidate`: `CAND-APCOA-001` と `CAND-APCOA-005` は review 証跡の最終シナリオへ統合できる可能性がある。
- `merge candidate`: `CAND-APCOA-002` と `CAND-APCOA-003` は検証結果証跡の最終シナリオへ統合できる可能性がある。
- `rejection candidate`: プロダクト実装ログ形式、ツール権限、エージェント実行定義の変更は operation-audit 候補から除外する。
