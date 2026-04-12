# 実装計画テンプレート

- workflow: impl
- status: planned
- lane_owner: orchestrating-implementation
- scope: master-dictionary-management
- task_id: master-dictionary-management
- task_catalog_ref: /Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/master-dictionary-management.yaml
- parent_phase: implementation-lane

## 要求要約

- マスター辞書ページで、辞書の取り込み、参照、作成、更新、削除導線を成立させる。

## 判断根拠

<!-- Decision Basis -->

- `/Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/master-dictionary-management.yaml` に completion criteria と manual check steps が定義されている。
- `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md` では、マスター辞書が基盤データとして UI から観測可能であることを要求している。
- user review-back により、マスターペルソナは `extractData.pas` 由来 JSON 前提で別 task へ切り分け、今回はマスター辞書に集中する。
- user review-back により、マスター辞書はダッシュボード配下ではなく独立ページとして扱う。

## 対象範囲

- `tasks/usecases/master-dictionary-management.yaml`
- `docs/exec-plans/active/2026-04-11-master-dictionary-management.md`
- `docs/exec-plans/active/master-dictionary-management.ui.html`
- `docs/exec-plans/active/master-dictionary-management.scenario.md`
- master dictionary 関連の frontend / backend 実装一式

## 対象外

- マスター辞書以外の画面詳細実装
- `docs/` 正本の恒久仕様変更

## 依存関係・ブロッカー

- 前段 HITL と後段 HITL の承認前は実装へ進めない。
- 個別 screen-design 正本が未整備のため、task-local UI mock で不足分を補う必要がある。

## 並行安全メモ

- 詳細設計確定前は plan 本文と task-local artifact 以外へ変更を広げない。
- frontend / backend 実装の並列化は `実装計画` section で task group 固定後に判断する。

## 機能要件

- `summary`:
  - 今回の task はマスター辞書だけを対象にし、マスターペルソナは別 task へ切り分ける。
  - マスター辞書ページは独立ページとして扱い、一覧参照、検索、詳細確認、新規作成、編集、削除、XML 取り込みを同一 task の機能要件に含める。
  - XML 取り込みはファイル選択 UI から開始し、選択後は同一画面内に選択中ファイル名と取込開始操作を持つ取込バーを表示する前提で固定する。
  - XML 取り込みでは `/Users/iorishibata/Repositories/AITranslationEngineJP/dictionaries/Dawnguard_english_japanese.xml` を読み込んで単語を抽出できること、および抽出対象 REC を許可リストに限定することを固定する。
  - 詳細上部の `更新` は編集モーダル、`削除` は確認モーダルで扱う前提を固定する。
  - 一覧サマリ表示、xTranslator 写像表示、実装方針が見える UI 表現は要求しない。
- `in_scope`:
  - マスター辞書ページへ独立して到達できること。
  - マスター辞書一覧を参照できること。
  - マスター辞書一覧で辞書エントリを検索できること。
  - 一覧で選択した辞書エントリの詳細情報を参照できること。
  - 一覧上から辞書データを新規作成できること。
  - 詳細上の `更新` から編集モーダルを開き、辞書データを編集できること。
  - 詳細上の `削除` から確認モーダルを開き、辞書データを削除できること。
  - `XMLから取り込み` からファイル選択 UI を開き、選択後の取込バー経由で辞書データを取り込めること。
  - XML 取り込み時は `BOOK:FULL`, `NPC_:FULL`, `NPC_:SHRT`, `ARMO:FULL`, `WEAP:FULL`, `LCTN:FULL`, `CELL:FULL`, `CONT:FULL`, `MISC:FULL`, `ALCH:FULL`, `FURN:FULL`, `DOOR:FULL`, `RACE:FULL`, `INGR:FULL`, `FLOR:FULL`, `SHOU:FULL` のみを単語抽出対象とし、それ以外の REC は抽出しないこと。
  - UI 上の文言、ラベル、説明には実装方針が見える表現を持ち込まないこと。
- `non_functional_requirements`:
  - マスター辞書ページは数万件規模の辞書レコードを保持しても、一覧参照、検索、選択、詳細確認、編集導線が破綻せず継続して操作できること。
  - XML 取り込み、新規作成、編集、削除の各導線は、同一ページ内で現在状態が把握できること。
  - UI 文言は `docs/spec.md` の用語に合わせ、日本語で一貫していること。
- `out_of_scope`:
  - マスターペルソナ画面および `extractData.pas` 由来 JSON を扱う導線。
  - 基盤データ以外の翻訳ジョブ画面、設定画面、翻訳成果物画面の詳細導線追加。
  - `docs/` 正本の恒久仕様変更、および usecase 完了条件を超える新規業務要件の追加。
  - 一覧サマリ表示、xTranslator 写像表示、利用状況表示、要件説明の露出。
  - UI モック、Scenario テスト一覧、実装計画、review 用差分図の確定。
- `open_questions`:
  - 辞書 detail に表示する項目を、原文、訳語、由来、最終更新のどこまでに絞るか。
- `required_reading`:
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-11-master-dictionary-management.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/master-dictionary-management.yaml`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/code.html`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/dictionaries/Dawnguard_english_japanese.xml`

## UI モック

- `artifact_path`: `docs/exec-plans/active/master-dictionary-management.ui.html`
- `final_artifact_path`: `docs/mocks/master-dictionary/index.html`
- `summary`:
  - マスター辞書は独立ページとして扱い、見出し、XML 取り込みカード、一覧と詳細の 2 カラムで同一画面内の現在状態を把握できる構成にする。
  - UI 全体は本文と見出しを明朝系へそろえ、ラベルと本文の見え方を統一する。
  - XML 取り込みカードの右上に `ファイルを選択` を置き、ファイル選択後は同一カード内の取込バーで選択ファイル名、状態表示、進捗、取込開始操作を確認できる。
  - XML 取り込み完了後は、同一カード内で追加件数と反映結果を示し、一覧件数更新、先頭 30 件への再着座、追加エントリの選択状態、詳細表示の切替を同じ画面で確認できる。
  - `新規登録` は一覧ヘッダーへ移し、検索とカテゴリ絞り込みの直上で一覧操作と並べて扱う。
  - 一覧行は訳語、原文、区分、ID の横並び高密度リストとして細く保ち、1ページ 30 件で選択中行を維持したまま詳細確認へ繋げる。
  - `新規登録` と `更新` は同系統の編集モーダルで扱い、保存後は同じページ内で一覧と詳細の両方に反映された状態を見せる。
  - `削除` は確認モーダルで扱い、完了後は一覧から対象を外し、詳細は同一画面内で次の表示対象または空状態へ切り替える。

## Scenario テスト一覧

- `artifact_path`: `docs/exec-plans/active/master-dictionary-management.scenario.md`
- `final_artifact_path`: `docs/scenario-tests/master-dictionary-management.md`
- `template_path`: `docs/exec-plans/templates/scenario-tests.md`
- `summary`:
  - Scenario artifact を `docs/exec-plans/active/master-dictionary-management.scenario.md` に固定し、最終正本を `docs/scenario-tests/master-dictionary-management.md` とした。
  - 一覧と検索は 1 ページ 30 件表示、同一画面での一覧選択と詳細同期を証明対象として固定した。
  - `新規登録` / `更新` は編集モーダル、`削除` は確認モーダルとして導線ごとの正常系を固定し、成功後は同一ページで一覧と詳細が同期更新されることを証明対象にした。
  - `XML取り込み` はファイル選択起点とし、ファイル未選択時は取込バーと取込開始操作を表示しないゲートを証明対象に固定した。
  - 取込バー状態は `待機中` `取込待ち` `取込中` `完了` の順序遷移を証明対象に固定した。
  - XML 取り込み完了後の同一ページ再同期として、検索/絞り込み解除、先頭 30 件への復帰、取込対象エントリの再選択と詳細表示を期待結果へ明記した。
  - 画面責務境界として、独立ページかつ同一画面内で状態可視性を維持する観点を追加した。

## 実装計画

<!-- Implementation Plan -->

- `ordered_scope`:
  1. 独立ページ route、screen controller、frontend gateway 契約を先に固定し、一覧取得、詳細取得、CRUD 応答、XML のファイル選択後取込開始、import 完了時 same-page refresh の DTO 境界を UI mock と scenario に合わせて揃える。
  2. backend でマスター辞書 CRUD と一覧検索を成立させ、30 件ページング、検索条件、詳細取得、作成、更新、削除後に再表示へ必要な識別子と最新値を返せる状態まで永続化を通す。
  3. frontend で一覧、検索、詳細、編集モーダル、削除確認モーダルを接続し、CRUD 成功後に同一ページ内で一覧件数、選択中行、詳細表示を更新する state 遷移を成立させる。
  4. frontend のファイル選択導線と import bar 表示条件を先に固定した上で、backend XML import pipeline と接続し、`ファイル選択` -> `取込バー表示` -> `取込開始` -> `進捗表示` -> `完了後の一覧/詳細/選択再同期` を同一ページ内で完結させる。
  5. integration と validation を揃え、approved flow と scenario artifact の最終内容に対して CRUD 成功後反映、XML import の file-selection-first 順序制約、import 完了後の same-page refresh 成立を証明する。
- `parallel_task_groups`:
  - `group_id`: `mdm-contract-shell`
    - `can_run_in_parallel_with`: なし
    - `blocked_by`: dashboard-and-app-shell の独立導線前提以外なし
    - `completion_signal`: route、screen controller、gateway contract、backend controller/usecase の entrypoint 名称、DTO 境界、mutation 成功時と import 完了時の same-page refresh 契約が UI mock と usecase に沿って固定されている。
  - `group_id`: `mdm-backend-crud`
    - `can_run_in_parallel_with`: `mdm-frontend-crud`
    - `blocked_by`: `mdm-contract-shell`
    - `completion_signal`: 一覧検索、詳細取得、作成、更新、削除 API/service/repository が同一 DTO 契約で通り、30 件ページングと検索条件を保持したまま、作成/更新/削除後に一覧再取得と詳細再同期へ必要な結果を返せる。
  - `group_id`: `mdm-frontend-crud`
    - `can_run_in_parallel_with`: `mdm-backend-crud`
    - `blocked_by`: `mdm-contract-shell`
    - `completion_signal`: 独立ページ、一覧、検索、詳細、編集モーダル、削除確認モーダルが approved UI mock に沿って接続され、CRUD 成功後に一覧件数、選択中行、詳細表示が同一ページで即時に整合する。
  - `group_id`: `mdm-import-flow`
    - `can_run_in_parallel_with`: なし
    - `blocked_by`: `mdm-contract-shell`, `mdm-backend-crud`, `mdm-frontend-crud`
    - `completion_signal`: XML は `ファイルを選択` でのみ対象を確定し、選択後にだけ取込バーが表示され、そこから取込開始でき、進捗と完了結果が同一ページで表示され、完了時に imported entry を含む一覧、詳細、選択状態が same-page refresh される。
  - `group_id`: `mdm-validation`
    - `can_run_in_parallel_with`: なし
    - `blocked_by`: `mdm-backend-crud`, `mdm-frontend-crud`, `mdm-import-flow`
    - `completion_signal`: approved flow を covering する integration / harness / scenario 対応 validation が実行可能で、CRUD 成功後反映、XML import の file-selection-first 順序制約、import 完了後の一覧/詳細/選択再同期の証跡が後段 review に渡せる。
- `task_dependencies`:
  - `task_id`: `mdm-contract-shell`
    - `depends_on`: `dashboard-and-app-shell`
    - `enables`: `mdm-backend-crud`, `mdm-frontend-crud`
    - `reason`: 独立ページ導線と DTO 契約が先に固定されないと frontend / backend を安全に並列化できず、import 完了後の same-page refresh 責務も曖昧になる。
  - `task_id`: `mdm-backend-crud`
    - `depends_on`: `mdm-contract-shell`
    - `enables`: `mdm-frontend-crud`, `mdm-import-flow`, `mdm-validation`
    - `reason`: 一覧検索、詳細、CRUD の永続化境界がないと frontend の再同期設計、import 完了時に返す affected entries、選択再評価ルールを確定できない。
  - `task_id`: `mdm-frontend-crud`
    - `depends_on`: `mdm-contract-shell`, `mdm-backend-crud`
    - `enables`: `mdm-import-flow`, `mdm-validation`
    - `reason`: CRUD 成功後の一覧件数更新、選択中行維持、詳細再同期、および XML 選択後にだけ取込バーを出す画面状態は backend 契約確定後にしか固定できない。
  - `task_id`: `mdm-import-flow`
    - `depends_on`: `mdm-contract-shell`, `mdm-backend-crud`, `mdm-frontend-crud`
    - `enables`: `mdm-validation`
    - `reason`: XML 取り込みは file-selection-first の画面制御を前提に CRUD と同じ辞書保存先へ反映し、完了直後に imported entries を一覧、詳細、選択へ反映する必要があるため、repository と frontend state shell の両方の確定が前提になる。
  - `task_id`: `mdm-validation`
    - `depends_on`: `mdm-backend-crud`, `mdm-frontend-crud`, `mdm-import-flow`
    - `enables`: `phase-8-review`
    - `reason`: approved flow 全体の成立証明は CRUD と import の両系統、および import 完了後を含む同一ページ反映の証跡が揃ってから実施する。
- `tasks`:
  - `task_id`: `mdm-contract-shell`
    - `owned_scope`: frontend route and page shell, frontend gateway contract, backend controller / usecase entrypoint, shared DTO naming, and response contract for CRUD/import same-page refresh
    - `depends_on`: `dashboard-and-app-shell`
    - `parallel_group`: `mdm-contract-shell`
    - `required_reading`:
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-11-master-dictionary-management.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/master-dictionary-management.yaml`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/master-dictionary-management.ui.html`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/master-dictionary-management.scenario.md`
    - `validation_commands`:
      - `python3 scripts/harness/run.py --suite structure`
  - `task_id`: `mdm-backend-crud`
    - `owned_scope`: backend controller/usecase/service/repository for list, search, detail, create, update, delete, and refresh payloads for same-page list/detail/selection updates
    - `depends_on`: `mdm-contract-shell`
    - `parallel_group`: `mdm-backend-crud`
    - `required_reading`:
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/master-dictionary-management.yaml`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/master-dictionary-management.scenario.md`
    - `validation_commands`:
      - `python3 scripts/harness/run.py --suite structure`
      - `python3 scripts/harness/run.py --suite execution`
  - `task_id`: `mdm-frontend-crud`
    - `owned_scope`: frontend state handling for list, search, detail, create modal, edit modal, delete confirmation modal, same-page refresh after create/update/delete, and file-selection-driven import bar visibility
    - `depends_on`: `mdm-contract-shell`, `mdm-backend-crud`
    - `parallel_group`: `mdm-frontend-crud`
    - `required_reading`:
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/master-dictionary-management.ui.html`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/master-dictionary-management.scenario.md`
    - `validation_commands`:
      - `python3 scripts/harness/run.py --suite structure`
      - `python3 scripts/harness/run.py --suite execution`
  - `task_id`: `mdm-import-flow`
    - `owned_scope`: backend XML import pipeline and persistence, frontend file-selection-first import bar state/progress/completion reflection, and imported-entry same-page refresh into list/detail/selection
    - `depends_on`: `mdm-contract-shell`, `mdm-backend-crud`, `mdm-frontend-crud`
    - `parallel_group`: `mdm-import-flow`
    - `required_reading`:
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/master-dictionary-management.yaml`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/master-dictionary-management.ui.html`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/master-dictionary-management.scenario.md`
    - `validation_commands`:
      - `python3 scripts/harness/run.py --suite structure`
      - `python3 scripts/harness/run.py --suite execution`
  - `task_id`: `mdm-validation`
    - `owned_scope`: integration wiring, scenario coverage alignment, final evidence for review, and proof of CRUD/import same-page state transitions including list/detail/selection refresh
    - `depends_on`: `mdm-backend-crud`, `mdm-frontend-crud`, `mdm-import-flow`
    - `parallel_group`: `mdm-validation`
    - `required_reading`:
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-11-master-dictionary-management.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/master-dictionary-management.scenario.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/master-dictionary-management.yaml`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/master-dictionary-management.ui.html`
    - `validation_commands`:
      - `python3 scripts/harness/run.py --suite structure`
      - `python3 scripts/harness/run.py --suite execution`
      - `python3 scripts/harness/run.py --suite all`
- `owned_scope`:
  - split。frontend は独立ページ shell、一覧検索状態、詳細表示、編集 / 削除 modal、CRUD 成功後の同一ページ再同期、ファイル選択による import bar 表示制御、import 完了後の一覧 / 詳細 / 選択再同期、進捗 / 完了表示を担当する。backend は controller / usecase / service / repository の CRUD、mutation 後と import 完了後の再表示に必要な応答、XML import pipeline を担当し、両者は Wails DTO 境界で接続する。

## 受け入れ確認

- マスター辞書一覧を参照できる。
- マスター辞書を作成できる。
- 選択中の基盤エントリの詳細を確認できる。
- 選択中の基盤エントリを更新できる。
- 選択中の基盤エントリを削除できる。
- `XMLから取り込み` 導線を確認できる。

## 必要な証跡

<!-- Required Evidence -->

- phase-1 以降の artifact path と要約
- 前段 HITL と後段 HITL の承認記録
- 実装レビュー結果
- `python3 scripts/harness/run.py --suite all` の最終結果

## 機能要件 HITL 状態

- approved

## 機能要件 承認記録

- 2026-04-11 orchestrating-implementation: 機能要件と UI モックを前段 HITL へ回付。承認待ち。
- 2026-04-12 human: 前段 HITL を承認。phase-2-scenario、phase-2-logic、review 用構造差分図作成へ進める。

## 詳細設計 HITL 状態

- approved

## 詳細設計 承認記録

- 2026-04-12 orchestrating-implementation: 詳細設計 AI review は pass。後段 HITL へ回付。承認待ち。
- 2026-04-12 human review-back: `/Users/iorishibata/Repositories/AITranslationEngineJP/dictionaries/Dawnguard_english_japanese.xml` を読み込んで単語を取れることと、単語抽出対象 REC を許可リストに限定することを明記するよう指示。
- 2026-04-12 human: 後段 HITL を承認。phase-5-test-implementation へ進める。

## review 用差分図

- `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/master-dictionary-management-structure-diff.d2`
- `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/master-dictionary-management-structure-diff.svg`
- frontend shell 更新と frontend / backend detail 図新規作成対象を 1 枚の review diff に集約した。

## 差分正本適用先

- `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/components/frontend/dashboard-and-app-shell.d2`
- `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/components/frontend/master-dictionary-management.d2`
- `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/components/backend/master-dictionary-management.d2`

## Closeout Notes

- review 用に active exec-plan 配下へ置いた差分 D2 / SVG copy は、`diagrams/backend/` または `diagrams/frontend/` 正本適用後に削除し、completed plan へ持ち越さない。
- 第1.6段階で作った UI モック working copy は、完了前に `docs/mocks/master-dictionary/index.html` へ移す。
- 第2段階で作った Scenario artifact working copy は、完了前に `docs/scenario-tests/master-dictionary-management.md` へ移す。

## 結果

<!-- Outcome -->

- in_progress
