# Scenario Candidates: 2026-05-17-master-persona-componentization / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `MPCOA`
- `task_artifact_location`: `docs/exec-plans/active/2026-05-17-master-persona-componentization/`
- `target_diff`: マスターペルソナ生成画面の部品化と、同じ props 境界を使う Storybook POC。

## Generator Scope

- `viewpoint`: `operation-audit`
- `included_sources`: `plan.md`, `docs/index.md`, `docs/architecture.md`, `docs/coding-guidelines-frontend.md`, `docs/coding-guidelines-tests.md`, `docs/UX-standard.md`, `docs/exec-plans/completed/2026-05-18-storybook-foundation/implementation-scope.md`, `docs/exec-plans/completed/2026-05-18-storybook-foundation/storybook-review.md`, `docs/spec.md`, `docs/er.md`, `docs/screen-design/screens/master-persona.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/scenario-tests/README.md`, `frontend/src/ui/screens/master-persona/`
- `excluded_sources`: プロダクトコード変更、プロダクトテスト変更、docs 正本化、他観点の候補成果物。
- `generation_notes`: 候補は Storybook で後追い確認できる状態、review 証跡、props fixture の再現性、保存禁止情報の境界に限定する。最終採否と統合は `designer` が行う。

## Candidate Scenarios

### CAND-MPCOA-001 Storybook 確認記録で対象 story を後追い確認できる

- `source requirement`: `plan.md` の close conditions は、主要状態を Storybook 上で確認できることを求める。Storybook foundation の `storybook-review.md` は review URL、story ID、確認状態、build 結果を記録する前例を持つ。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-MPCOA-001`
- `actor`: 実装後確認者、見た目レビュー担当者。
- `trigger`: マスターペルソナ部品の Storybook story を追加し、Storybook dev 表示または static build を確認する。
- `expected outcome`: 確認記録から、対象 component、story ID、review URL、確認状態、build 結果、未確認理由を後から追える。
- `audit event`: Storybook review 証跡の作成または更新。
- `saved summary`: story ID、story title、story name、component path、確認済み状態、未確認理由、再実行 command、build 結果。
- `redaction rule`: review URL、story ID、確認記録に secret、API key、token、ローカル絶対 path、実ユーザーデータ、入力 JSON 本文を含めない。
- `observable point`: `storybook-review.md` または同等の task-local 証跡で、対象 story と確認結果を読める。
- `related detail requirement type`: `observability_requirement`, `testability_requirement`, `security_requirement`
- `adoption hint`: `SBF-ST-03-storybook-review-evidence` の記録形を踏襲すると、Storybook foundation との整合を保てる。
- `conflict hint`: review 証跡の置き場所を `ui-design.md` に統合するか、別 artifact にするかは `plan.md` の未決事項である。

### CAND-MPCOA-002 主要状態の story fixture から UI 状態を再現できる

- `source requirement`: `plan.md` は未設定、生成中、生成成功、生成失敗、編集中を現行 view model で表現できるかを未決事項にしている。`docs/screen-design/screens/master-persona.md` は生成準備、進行状況、一覧、詳細、編集モーダル、削除モーダルの状態表示を定義している。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-MPCOA-002`
- `actor`: 実装後確認者、回帰確認者。
- `trigger`: `GenerationSetupPanel`, `RunStatusPanel`, `PersonaReviewPanel`, `PersonaActionModal` の story fixture を作る。
- `expected outcome`: Storybook 上で、入力待ち、JSON 選択済み、生成中、生成失敗表示、一覧空状態、選択済み詳細、編集モーダル、削除モーダルを再現できる。
- `audit event`: component props fixture による UI 状態の固定。
- `saved summary`: component 名、fixture 名、対象状態、表示される状態値、操作可否の要約。
- `redaction rule`: fixture には実ユーザーの NPC データ、実 JSON 内容、長い原文発話、会話文脈全文、provider raw request / response、raw prompt を入れない。
- `observable point`: Storybook の story registry と各 story の表示内容で、状態別 fixture が存在することを確認できる。
- `related detail requirement type`: `observability_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: 状態ごとの story は、部品化後の props 境界が過剰に親画面へ依存していないことの確認材料になる。
- `conflict hint`: 主要状態の最小セットは他観点の正常系、失敗系、状態遷移候補と重なる可能性がある。

### CAND-MPCOA-003 Storybook fixture が secret と実データを保存しない

- `source requirement`: Storybook foundation の `implementation-scope.md` は、fixture、story、review URL、command 記録に secret、API key、token、実ユーザーデータを出さないことを固定している。`docs/detail-specs/persona-generation-phase.md` は secret、API key、credential 参照実値、secret store key、endpoint、provider raw request / response、raw prompt、原文発話全文、会話文脈全文の露出禁止を定義している。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-MPCOA-003`
- `actor`: 実装後確認者、セキュリティ確認者。
- `trigger`: マスターペルソナ story と fixture を追加し、Storybook build または review 証跡を残す。
- `expected outcome`: fixture と確認記録は、表示用の短いサンプル値、件数、状態値だけを持ち、secret や実データを保存しない。
- `audit event`: fixture 内容と review 証跡の redaction 確認。
- `saved summary`: fixture path、story path、使用するサンプル値の種類、保存禁止情報が含まれない確認結果。
- `redaction rule`: API key 本文、credential 参照実値、secret store key、endpoint、token、provider raw payload、raw prompt、実ファイル path、実 NPC 本文、入力 JSON 本文を保存しない。
- `observable point`: fixture file、story file、review 証跡、build log 要約に保存禁止情報がないことを確認できる。
- `related detail requirement type`: `security_requirement`, `data_requirement`, `observability_requirement`
- `adoption hint`: AI 設定カードとペルソナ本文の fixture は、表示値を固定サンプルに限定すると監査しやすい。
- `conflict hint`: ペルソナ本文の表示確認には文章サンプルが必要になるため、どこまでを「過剰な本文」とみなすかは人間判断候補になる。

### CAND-MPCOA-004 部品化後も import 境界の検査結果を確認できる

- `source requirement`: `docs/architecture.md` は UI Component が backend DTO、generated binding、Store、Gateway を直接扱わないことを定義している。Storybook foundation の `implementation-scope.md` は story と fixture が generated `wailsjs`、backend DTO、Gateway、RuntimeEventAdapter、AI provider、secret store、DB、filesystem business flow を import しないことを定義している。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-MPCOA-004`
- `actor`: 実装後確認者、責務境界レビュー担当者。
- `trigger`: マスターペルソナ画面を表示部品へ分割し、Storybook story と fixture を追加する。
- `expected outcome`: lint または boundary check の結果から、Storybook と UI Component が禁止 import を持たないことを後から確認できる。
- `audit event`: import 境界検査の実行と結果記録。
- `saved summary`: 実行 command、対象範囲、結果、失敗時の禁止 import 種別、未実行理由。
- `redaction rule`: lint 出力や記録にローカル絶対 path 全文、secret、API key、token、実ユーザーデータを保存しない。
- `observable point`: `npm --prefix frontend run lint`、`npm --prefix frontend run lint:boundaries`、または `frontend-local` の結果記録で境界検査の有無を追える。
- `related detail requirement type`: `observability_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: componentization の成否を見た目だけでなく責務境界の後追い確認へ接続できる。
- `conflict hint`: 実行 command の最終採用は `designer` と implementation-scope が決めるため、この候補では検査入口を固定しない。

### CAND-MPCOA-005 生成実行中の表示と操作不可理由を story で確認できる

- `source requirement`: `docs/screen-design/screens/master-persona.md` は生成中に編集と削除を操作不可にし、進行状況、処理件数、現在の対象、一時停止、中止を表示する。`docs/UX-standard.md` は状態表示、状態別操作、禁止条件、検査可能性を高優先度として扱う。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-MPCOA-005`
- `actor`: 実装後確認者、見た目レビュー担当者。
- `trigger`: 生成中状態の view model fixture を Storybook で表示する。
- `expected outcome`: 生成中状態の story で、進行状況、処理済み件数、作成済み件数、既存スキップ数、現在の対象、編集不可、削除不可、一時停止可、中止可を後から確認できる。
- `audit event`: 生成中 UI 状態の再現証跡。
- `saved summary`: run state、processed count、success count、existing skip count、current actor label、操作可否の要約。
- `redaction rule`: current actor label には実 NPC 名や実プラグイン名を使わず、固定サンプル名にする。入力 JSON の path と本文は保存しない。
- `observable point`: `RunStatusPanel` と `PersonaReviewPanel` の story で、生成中の操作可否と状態表示を確認できる。
- `related detail requirement type`: `observability_requirement`, `state_requirement`, `testability_requirement`
- `adoption hint`: 生成中の story は、部品化で親画面の状態更新手順を変えずに表示だけ切り出せたかを確認する材料になる。
- `conflict hint`: 一時停止と中止の実処理は state-transition 観点または failure 観点と重なる可能性がある。

### CAND-MPCOA-006 生成失敗や AI 設定不足の表示が保存禁止情報を含まない

- `source requirement`: `MasterPersonaPage.svelte` は `errorMessage` と `aiSettingsMessage` を通知として表示する。`docs/detail-specs/persona-generation-phase.md` は障害調査用の要約で provider、model、execution mode、batch mode、credential 状態分類、input count、output count、prompt digest、error kind を確認できる一方、secret と raw payload を出さないことを定義している。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-MPCOA-006`
- `actor`: 実装後確認者、セキュリティ確認者。
- `trigger`: AI 設定不足、モデル一覧未取得、生成失敗の表示状態を story fixture で再現する。
- `expected outcome`: UI と story fixture は利用者が判断できる短いエラー種別や対応文だけを表示し、provider raw request / response、raw prompt、API key、endpoint を出さない。
- `audit event`: 失敗表示 fixture の redaction 確認。
- `saved summary`: error kind、user-facing message、credential 状態分類、対象件数、確認した story ID。
- `redaction rule`: API key、credential 参照実値、secret store key、endpoint、provider raw payload、raw prompt、入力 JSON 本文、原文発話全文を保存しない。
- `observable point`: 失敗状態の story と review 証跡で、表示文言と保存禁止情報の不在を確認できる。
- `related detail requirement type`: `security_requirement`, `observability_requirement`, `failure_handling_requirement`
- `adoption hint`: 失敗状態の fixture は、見た目確認と redaction 確認を同じ review URL で扱える。
- `conflict hint`: 失敗種別の粒度は failure 観点の候補と重なるため、最終シナリオでは統合対象になる可能性がある。

## Open Notes

- `human decision candidate`: ペルソナ本文の story サンプルで許容する文字量と内容。短い架空サンプルを使う案が安全だが、見た目確認に必要な長文耐性との境界は人間判断が必要である。
- `human decision candidate`: Storybook review 証跡を `ui-design.md` 内へ残すか、`storybook-review.md` として分けるか。`plan.md` では未決事項である。
- `merge candidate`: `CAND-MPCOA-001` と `CAND-MPCOA-002` は Storybook 確認シナリオとして統合できる可能性がある。
- `merge candidate`: `CAND-MPCOA-003` と `CAND-MPCOA-006` は redaction 確認シナリオとして統合できる可能性がある。
- `rejection candidate`: プロダクト実行ログの追加、DB 監査テーブル、operation summary 永続化は今回の部品化と Storybook POC の対象外である。
