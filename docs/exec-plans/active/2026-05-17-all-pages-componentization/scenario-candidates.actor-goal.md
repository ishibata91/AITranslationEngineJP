# Scenario Candidates: 2026-05-17-all-pages-componentization / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `APC`
- `candidate_count`: 10

## Generator Scope

- `viewpoint`: 利用者目的、レビュー目的、画面確認目的から、全ページ部品化と Storybook story 追加の候補を作る。
- `running_task_artifact_location`: `docs/exec-plans/active/2026-05-17-all-pages-componentization/`
- `target_diff`: Master Persona POC 後に、全ページの page component を薄くし、主要なパネル、カード、モーダルへ props 境界と Storybook story を付ける。
- `candidate_artifact_path`: `docs/exec-plans/active/2026-05-17-all-pages-componentization/scenario-candidates.actor-goal.md`
- `included_sources`: `./plan.md`, `docs/architecture.md`, `docs/coding-guidelines-frontend.md`, `docs/UX-standard.md`, 完了済み Storybook 基盤 scenario, 完了済み Master Persona 部品化 scenario と UI design, `frontend/src/ui/screens/`, `frontend/src/ui/components/`。
- `excluded_sources`: プロダクトコード変更、テスト変更、docs 正本変更、最終シナリオ表、採否、統合判断、他観点の候補生成。
- `generation_notes`: Storybook は固定 props または view model fixture による表示確認に限定する。backend、Wails runtime、Gateway、generated binding、AI provider、secret store、DB、実 filesystem flow の成立証明は扱わない。

## Candidate Scenarios

### CAND-APC-001 翻訳入力を確認する利用者がロード結果を判断できる

- `source requirement`: `plan.md` は全ページの主要表示領域へ props 境界と Storybook story を付ける。現行 `translation-input` には `DataLoadHero`, `DataLoadImportPanel`, `LoadedInputList`, `LoadedInputDetail` が存在する。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-APC-001`
- `actor`: Skyrim Mod 翻訳作業の利用者。
- `purpose`: 利用者は、XML または読み込み済み入力の状態、カテゴリ件数、サンプル項目、問題と次の手順を確認できる。
- `trigger`: 利用者が翻訳入力確認画面を開く。Storybook 確認者は翻訳入力の主要パネル story を開く。
- `expected outcome`: ロード準備、読み込み済み一覧、選択データ詳細、問題と再構築の表示が、画面と Storybook story の両方で同じ判断単位として見える。
- `observable point`: `translation-input-review-screen-status-header`, `translation-input-review-load-preparation-region`, `translation-input-review-loaded-input-list`, `translation-input-review-selected-input-region`, story ID, fixed fixture。
- `related detail requirement type`: `success_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 入力確認画面の既存部品構成を基準に、他画面の分割粒度を決める候補にできる。
- `conflict hint`: lifecycle 観点が実ファイル選択や実 XML 読み込みを Storybook の受け入れ条件に含める場合、props fixture 限定の前提と競合する。

### CAND-APC-002 ジョブ作成者が入力、共通資産、AI 設定を作成前に確認できる

- `source requirement`: `plan.md` は全ページ一括で主要表示領域を棚卸しする。現行 `JobSetupPage` は入力データ、共通辞書と共通ペルソナ、AI サービスとモデル、作成前確認を同一 page に持つ。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-APC-002`
- `actor`: 翻訳ジョブを作成する利用者。
- `purpose`: 利用者は、入力データ、共通辞書、共通ペルソナ、翻訳段階別 AI 設定、作成前確認を順に確認し、ジョブを作成できるか判断できる。
- `trigger`: 利用者がジョブ作成画面を開く。Storybook 確認者は入力データ、共通資産、段階別設定、作成前確認の story を開く。
- `expected outcome`: 作成前に必要な判断単位が分かれ、作成可能状態と作成不可理由が story fixture で確認できる。
- `observable point`: `translation-job-setup-input-data-region`, `translation-job-setup-shared-dictionary-persona-region`, `translation-job-setup-ai-service-model-region`, `translation-job-setup-compatibility-precheck-region`, Storybook canvas。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `testability_requirement`
- `adoption hint`: 複数の業務判断を 1 page に持つ画面の代表候補として使える。
- `conflict hint`: state-transition 観点が作成可能条件の詳細を固定する場合、この候補の利用者目的と結合される可能性がある。

### CAND-APC-003 未完了ジョブを探す利用者が状態と操作可否を判断できる

- `source requirement`: `plan.md` は既存表示項目の維持確認を必要とする。現行 `TranslationJobManagementPage` は検索、状態 filter、job card、操作、無効理由、削除確認 modal を持つ。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-APC-003`
- `actor`: 翻訳作業を再開または整理する利用者。
- `purpose`: 利用者は、未完了ジョブを検索し、現在状態、再開可否、削除可否、無効理由を見て次操作を選べる。
- `trigger`: 利用者が翻訳管理画面を開く。Storybook 確認者は一覧、job card、action button、削除 modal の story を開く。
- `expected outcome`: 一覧の空状態、読み込み失敗、一覧あり、操作無効理由、削除確認が個別 story で確認できる。
- `observable point`: `translation-job-management-job-list-region`, `translation-job-management-job-card`, `translation-job-management-job-actions`, `translation-job-management-disabled-reason`, `translation-job-management-delete-confirmation-modal`。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: job card と削除 modal は、画面専用 component として story 化する候補にできる。
- `conflict hint`: failure 観点が削除失敗や再開不能理由を詳細化する場合、この候補の story 状態が増える可能性がある。

### CAND-APC-004 翻訳段階の作業者が進行状況と次段階 readiness を確認できる

- `source requirement`: `plan.md` は全ページをパネル、カード、モーダル単位へ部品化する。現行 `term-translation-phase`, `persona-generation-phase`, `body-translation-phase` は状態 header、操作、進行状況、実行設定、結果、失敗情報を持つ。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-APC-004`
- `actor`: 翻訳段階を実行または監視する利用者。
- `purpose`: 利用者は、単語翻訳、NPC ペルソナ生成、本文翻訳で、進行状況、処理件数、実行設定、結果、失敗情報、次段階 readiness を確認できる。
- `trigger`: 利用者が対象 phase 画面を開く。Storybook 確認者は phase ごとの状態 header、操作 card、進行 card、結果 card、失敗情報 card の story を開く。
- `expected outcome`: phase 画面の共通判断単位が story fixture で比較でき、実行中、完了、失敗、次段階不可が画面上で区別できる。
- `observable point`: `term-translation-phase-progress-region`, `persona-generation-phase-progress-card`, `body-translation-phase-progress`, `body-translation-phase-output-readiness`, phase story ID。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `alternative_success_requirement`, `testability_requirement`
- `adoption hint`: 3 つの phase 画面で似た表示単位を共有化候補として調べる入口にできる。
- `conflict hint`: 共通化を優先しすぎると、phase 固有の表示規則が props 分岐の塊になる可能性がある。

### CAND-APC-005 ジョブ実行画面の利用者が選択ジョブと現在 phase を見失わない

- `source requirement`: `docs/architecture.md` は page view を `UI Component` の合成へ寄せる。現行 `JobRunPage` は選択ジョブ summary、phase screen region、未選択 guidance、phase navigation footer を持つ。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-APC-005`
- `actor`: 翻訳ジョブを段階的に進める利用者。
- `purpose`: 利用者は、選択中ジョブ、現在 phase、未選択時の案内、phase 間 navigation を混同せずに作業できる。
- `trigger`: 利用者が job run 画面を開く。Storybook 確認者は選択ジョブ summary、未選択 guidance、phase navigation footer の story を開く。
- `expected outcome`: `JobRunPage` は controller 接続と phase 合成に寄り、利用者に見える job summary と navigation は props で確認できる。
- `observable point`: `job-run-selected-job-summary`, `job-run-phase-screen-region`, `job-run-job-unselected-guidance`, `PhaseNavigationFooter` の story。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `compatibility_requirement`
- `adoption hint`: page 合成 story を密度確認と配置確認に限定する候補として使える。
- `conflict hint`: lifecycle 観点が phase 実行そのものを `JobRunPage` story に含める場合、page story の限定前提と競合する。

### CAND-APC-006 出力確認者が候補、差分、出力操作を判断できる

- `source requirement`: `plan.md` は実装後の Storybook review を close condition に含める。現行 `TranslationOutputArtifactPage` は出力候補、選択 job summary、出力操作、差分 preview を持つ。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-APC-006`
- `actor`: 翻訳結果の出力を確認する利用者。
- `purpose`: 利用者は、出力候補を選び、選択 job の状態、拒否理由、最新結果、差分 preview、出力操作を確認できる。
- `trigger`: 利用者が出力成果物画面を開く。Storybook 確認者は候補一覧、選択 summary、操作 card、差分 preview の story を開く。
- `expected outcome`: 出力候補なし、候補選択済み、出力操作可、拒否理由あり、差分ありの表示状態を story fixture で確認できる。
- `observable point`: `output-management-output-candidate-list`, `output-management-selected-job`, `output-management-output-actions`, `output-management-diff-preview`, Storybook canvas。
- `related detail requirement type`: `success_requirement`, `alternative_success_requirement`, `testability_requirement`
- `adoption hint`: 出力レビューの UI 表示確認を backend の実出力生成と分離できる。
- `conflict hint`: external-integration 観点が xTranslator 出力生成の成立を Storybook に含める場合、この候補の表示確認境界と競合する。

### CAND-APC-007 マスター辞書管理者が取り込み、一覧、詳細、編集、削除を確認できる

- `source requirement`: `docs/UX-standard.md` は一覧、詳細、フォーム、危険操作、空状態、エラー状態の検査可能性を重視する。現行 `MasterDictionaryPage` は取り込み導線、辞書一覧、詳細、作成編集 modal、削除確認 modal を持つ。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-APC-007`
- `actor`: マスター辞書を管理する利用者。
- `purpose`: 利用者は、XML 取り込み結果、辞書一覧、検索またはページング、選択詳細、作成編集、削除確認を同じ画面で判断できる。
- `trigger`: 利用者がマスター辞書画面を開く。Storybook 確認者は取り込み、一覧、詳細、作成編集 modal、削除 modal の story を開く。
- `expected outcome`: 長い辞書項目、空一覧、取り込み結果、編集失敗、削除確認が story fixture で確認できる。
- `observable point`: `master-dictionary-xml-import-region`, `master-dictionary-dictionary-list-panel`, `master-dictionary-detail-panel`, `master-dictionary-create-edit-modal`, `master-dictionary-delete-confirmation-modal`。
- `related detail requirement type`: `success_requirement`, `failure_handling_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 一覧、詳細、modal が厚い page に混在する画面の分割候補にできる。
- `conflict hint`: failure 観点が import 失敗や保存失敗の扱いを詳細化する場合、story 状態の粒度が増える可能性がある。

### CAND-APC-008 AI provider 設定者が資格情報と接続状態を確認できる

- `source requirement`: 完了済み Storybook 基盤 scenario は fixture に secret、API key、token、実 endpoint を含めない。現行 `ProviderSettingsPage` は provider 一覧、設定 detail、API key 状態、API key 入力、接続確認、保存操作を持つ。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-APC-008`
- `actor`: AI provider 設定を行う利用者。
- `purpose`: 利用者は、provider ごとの設定状態を比較し、選択 provider の endpoint、API key 状態、接続確認結果、保存可否を判断できる。
- `trigger`: 利用者が provider 設定画面を開く。Storybook 確認者は provider 一覧、設定 detail、API key panel、接続確認 panel、保存 action の story を開く。
- `expected outcome`: 未設定、設定済み、入力中、接続確認中、接続失敗、保存不可の表示状態が story fixture で確認できる。fixture は実 secret と実 endpoint を持たない。
- `observable point`: `provider-settings-ai-service-list-region`, `provider-settings-settings-detail-region`, `provider-settings-api-key-status-region`, `provider-settings-api-key-input-region`, `provider-settings-connection-check-region`。
- `related detail requirement type`: `success_requirement`, `security_requirement`, `state_requirement`, `testability_requirement`
- `adoption hint`: secret 表示禁止と UI 表示確認を分ける actor-goal 候補として使える。
- `conflict hint`: external-integration 観点が接続確認の実 network 成否を Storybook に含める場合、固定 fixture の前提と競合する。

### CAND-APC-009 人間レビュアーが全ページの主要部品を Storybook で順に確認できる

- `source requirement`: `plan.md` は実装後の Storybook review と frontend-local を close condition にする。完了済み Master Persona UI design は review URL、story ID、確認状態、未確認理由、再実行 command を `storybook-review.md` に残す方針を持つ。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-APC-009`
- `actor`: 人間見た目レビュー担当者。
- `purpose`: レビュアーは、全ページの主要なパネル、カード、モーダルを Storybook 上で画面順に確認し、確認済み、未確認、未確認理由を追跡できる。
- `trigger`: 実装後にレビュアーが `storybook-review.md` の review URL と story ID を開く。
- `expected outcome`: 各 screen の主要 story が一覧化され、表示崩れ、長文耐性、空状態、エラー状態、モーダル状態を確認できる。
- `observable point`: task-local `storybook-review.md`, Storybook review URL, story ID, 確認状態, 未確認理由, `npm --prefix frontend run build-storybook` の結果。
- `related detail requirement type`: `success_requirement`, `observability_requirement`, `testability_requirement`
- `adoption hint`: 利用者向け正常系候補とは別に、人間レビューの受け入れ入口として使える。
- `conflict hint`: operation-audit 観点が review 証跡を docs 正本へ即時反映する前提を置く場合、この task-local 記録前提と競合する。

### CAND-APC-010 画面確認者が props 境界と既存表示維持を同時に確認できる

- `source requirement`: `docs/architecture.md` は `UI Component` が backend DTO、generated binding、`Store`、`Gateway` を直接扱わないと定義する。`docs/coding-guidelines-frontend.md` は意味単位、明確な props、明確な event、閉じた内部状態を component 分割基準にする。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-APC-010`
- `actor`: 実装後の画面確認者または frontend レビュー担当者。
- `purpose`: 確認者は、page component が controller 接続と部品合成へ薄くなり、主要表示部品が props と callback だけで表示でき、既存表示項目が失われていないことを確認できる。
- `trigger`: 確認者が実装後の source、Storybook story、fixture、画面表示証跡を見る。
- `expected outcome`: screen local component は `Store`、`Gateway`、generated binding を直接 import しない。共有化した component は複数画面で使う UI 規則だけを持つ。画面専用表示は screen local に残る。
- `observable point`: component import 境界、props 型、story args、fixture source、既存 screen の `data-testid` または `aria-label`、Storybook 表示。
- `related detail requirement type`: `compatibility_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: 全ページ横断の props 境界確認シナリオとして使える。
- `conflict hint`: 共有化候補を広く取りすぎると、画面固有条件を props 分岐として増やす実装と競合する可能性がある。

## Open Notes

- `human decision candidate`: 画面順序は未決である。`plan.md` は事前調査で列挙する部品化候補の画面順序を未決事項にしている。
- `human decision candidate`: 共通化候補を `frontend/src/ui/components/` へ上げる具体基準は未決である。`docs/architecture.md` の判断表は正本だが、今回の対象候補ごとの最終分類は designer が固定する必要がある。
- `human decision candidate`: 既存表示項目を削る場合の人間承認粒度は未決である。actor-goal 候補は既存表示維持を前提に置く。
- `human decision candidate`: Master Persona は POC 後の基準として扱うか、全ページ対象として再確認するかを designer が明示する必要がある。
- `human decision candidate`: page 合成 story をどの画面で任意にするか、どの画面で密度確認に必要とするかは最終シナリオ統合で固定する必要がある。
- `merge candidate`: `CAND-APC-004` と `CAND-APC-005` は、job run 系画面の story 確認として統合される可能性がある。
- `merge candidate`: `CAND-APC-009` と `CAND-APC-010` は、review 証跡と props 境界確認を 1 つの横断シナリオへ統合できる可能性がある。
- `rejection candidate`: Storybook story が Gateway、generated binding、backend DTO mock、実 AI provider、secret store、DB、実 filesystem flow を要求する候補は、完了済み Storybook 基盤の境界に反する。
- `rejection candidate`: 親画面の状態を大量に読む部品、業務フロー全体の進行状態を持つ部品、props が条件分岐の塊になる共有 component は、部品化候補から外す可能性がある。
- `rejection candidate`: docs 正本変更、`.codex` 変更、screen-design 正本反映は、この候補生成成果物では扱わない。
