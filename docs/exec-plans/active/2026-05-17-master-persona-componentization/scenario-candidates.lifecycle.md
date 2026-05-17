# Scenario Candidates: 2026-05-17-master-persona-componentization / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `MPC-LC`
- `candidate_count`: `8`

## Generator Scope

- `viewpoint`: lifecycle
- `included_sources`: `./plan.md`, `docs/index.md`, `docs/spec.md`, `docs/architecture.md`, `docs/er.md`, `docs/coding-guidelines-frontend.md`, `docs/coding-guidelines-tests.md`, `docs/UX-standard.md`, `docs/screen-design/README.md`, `docs/detail-specs/README.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/scenario-tests/README.md`, `docs/exec-plans/completed/2026-05-18-storybook-foundation/implementation-scope.md`, `docs/exec-plans/completed/2026-05-18-storybook-foundation/storybook-review.md`, `frontend/src/ui/screens/master-persona/`
- `excluded_sources`: プロダクトコード変更、プロダクトテスト変更、docs 正本化、他観点の候補成果物、最終シナリオ採否
- `generation_notes`: Storybook 基盤完了後に、マスターペルソナ生成画面の部品化と Storybook POC が同じ props 境界を共有できるかを、作成前、preview、生成中、完了後確認、編集、削除、Storybook 確認の時間順で候補化した。

## Candidate Scenarios

### CAND-MPC-LC-001 初期表示を部品 props へ渡せる

- `source requirement`: `MasterPersonaPage` は controller 接続と部品合成へ寄せる。`UI Component` は backend DTO、generated binding、`Store`、`Gateway` を直接扱わない。
- `viewpoint`: lifecycle / 作成前
- `candidate scenario id`: `CAND-MPC-LC-001`
- `actor`: マスターペルソナ生成画面を開くユーザー
- `trigger`: 画面を初回表示し、controller が `入力待ち` の view model を返す。
- `expected outcome`: 生成準備、進行状況、生成結果、操作モーダルの各部品は、小さい props と callback だけで初期状態を表示する。
- `observable point`: Storybook fixture または画面表示で、`入力待ち`、未選択ファイル、空の preview、空一覧、閉じた modal が同時に確認できる。
- `related detail requirement type`: `compatibility_requirement`, `testability_requirement`, `state_requirement`
- `adoption hint`: 初期 story と画面回帰の共通 fixture にしやすい。
- `conflict hint`: page story を作るか、review-only wrapper を作るかの未決と競合する可能性がある。

### CAND-MPC-LC-002 JSON 選択後の preview が生成準備部品へ反映される

- `source requirement`: 既存画面は JSON 選択後に preview を取得し、`candidateCount`、`newlyAddableCount`、`existingCount` を表示する。
- `viewpoint`: lifecycle / 入力選択と preview
- `candidate scenario id`: `CAND-MPC-LC-002`
- `actor`: 生成対象 JSON を選ぶユーザー
- `trigger`: ユーザーが JSON を選択し、preview 状態が `生成可能` または `設定未完了` として返る。
- `expected outcome`: 生成準備部品は候補数、新規作成数、既存スキップ数を表示し、生成開始可否を AI 設定と preview 状態から決める。
- `observable point`: `GenerationSetupPanel` の story または画面で、preview 集計、ファイル名、生成ボタンの enabled / disabled が確認できる。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: preview 成功、設定未完了、ファイル未選択の 3 fixture に分けると採用しやすい。
- `conflict hint`: 現行 view model だけで未設定と生成可能を十分に表現できるかという未決と競合する。

### CAND-MPC-LC-003 生成開始後の進行状態が実行中 props で固定される

- `source requirement`: 画面は `runStatus` と `isRunActive` から進行状況、進捗、現在対象、一時停止、中止を表示する。
- `viewpoint`: lifecycle / 実行中
- `candidate scenario id`: `CAND-MPC-LC-003`
- `actor`: ペルソナ生成を開始したユーザー
- `trigger`: 生成開始後に `runState` が `生成中` になり、処理済み件数、作成済み件数、既存スキップ数、現在対象が更新される。
- `expected outcome`: 進行状況部品は生成中の状態を表示し、編集と削除は不可になり、一時停止と中止は有効になる。
- `observable point`: `RunStatusPanel` と結果確認部品の story または画面で、進捗 bar、`生成中` 状態、編集削除 disabled、一時停止中止 enabled が確認できる。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: 実行中 fixture は画面全体 story と個別部品 story の両方に使える。
- `conflict hint`: 生成中の操作可否を画面全体 fixture で見るか、各部品 story で見るかの統合判断が必要である。

### CAND-MPC-LC-004 生成完了後に一覧と詳細を確認できる

- `source requirement`: 生成完了または実行状態が active から inactive へ遷移した時、ページを再取得し、生成結果一覧と詳細を確認できる。
- `viewpoint`: lifecycle / 完了後確認
- `candidate scenario id`: `CAND-MPC-LC-004`
- `actor`: 生成完了後に結果を確認するユーザー
- `trigger`: `runState` が `完了` になり、一覧と選択中 detail が view model に反映される。
- `expected outcome`: 結果確認部品は件数範囲、検索、プラグイン filter、一覧、詳細、編集削除操作を表示する。
- `observable point`: `PersonaReviewPanel` の story または画面で、一覧 row、選択中 detail、ページャ、`canMutate=true` の編集削除ボタンが確認できる。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: 生成完了後の review fixture は人間見た目レビューの主要確認対象にできる。
- `conflict hint`: 生成完了の story が product gateway の再取得を証明するものではない点を designer が分ける必要がある。

### CAND-MPC-LC-005 完了後の編集 modal は保存後に閉じる

- `source requirement`: 選択済みペルソナは編集 modal で更新でき、保存成功時は `modalState` が `null` になる。
- `viewpoint`: lifecycle / 完了後編集
- `candidate scenario id`: `CAND-MPC-LC-005`
- `actor`: 生成済みペルソナを補正するユーザー
- `trigger`: ユーザーが選択済み detail から編集 modal を開き、要約、話し方、本文を保存する。
- `expected outcome`: 編集 modal は選択中 identity と編集フォームを表示し、保存成功後に閉じ、一覧または detail の更新状態を表示できる。
- `observable point`: `PersonaActionModal` の edit story または画面で、modal open、入力項目、保存 callback、保存後の modal closed 状態が確認できる。
- `related detail requirement type`: `alternative_success_requirement`, `data_requirement`, `state_requirement`
- `adoption hint`: 編集 modal は部品 story の対象として独立させやすい。
- `conflict hint`: modal 保存後の一覧再取得まで story で扱うか、usecase 側シナリオへ分けるかが競合する。

### CAND-MPC-LC-006 完了後の削除 modal は確認後に閉じる

- `source requirement`: 選択済みペルソナは削除 modal で確認して削除でき、削除成功時は `modalState` が `null` になる。
- `viewpoint`: lifecycle / 完了後削除
- `candidate scenario id`: `CAND-MPC-LC-006`
- `actor`: 不要な生成済みペルソナを削除するユーザー
- `trigger`: ユーザーが選択済み detail から削除 modal を開き、削除を確定する。
- `expected outcome`: 削除 modal は対象識別情報を表示し、削除成功後に閉じ、一覧から対象が外れた状態を表示できる。
- `observable point`: `PersonaActionModal` の delete story または画面で、modal open、対象識別情報、削除 callback、削除後の modal closed 状態が確認できる。
- `related detail requirement type`: `alternative_success_requirement`, `data_requirement`, `state_requirement`
- `adoption hint`: 削除 modal は危険操作として edit modal と別 story にする価値がある。
- `conflict hint`: 削除後にどの row を選択するかは designer の最終シナリオで固定する必要がある。

### CAND-MPC-LC-007 preview または生成開始の失敗後も同じ段階へ戻れる

- `source requirement`: preview 失敗時は preview を消して error message を保持する。生成開始失敗時は error message を保持し、ページ再取得を行わない。
- `viewpoint`: lifecycle / 段階内失敗から再試行
- `candidate scenario id`: `CAND-MPC-LC-007`
- `actor`: 失敗後に入力や AI 設定を見直すユーザー
- `trigger`: preview 取得または生成開始が失敗し、画面が error message を受け取る。
- `expected outcome`: notice banner は失敗理由を表示し、生成準備部品は同じ入力段階を維持し、ユーザーが JSON 選び直しまたは AI 設定更新へ戻れる。
- `observable point`: 画面 story または部品 story で、error banner、preview null、生成ボタン disabled、入力選び直し操作が確認できる。
- `related detail requirement type`: `failure_handling_requirement`, `recovery_requirement`, `state_requirement`
- `adoption hint`: 失敗系の採否は failure 観点と重複しやすいため、lifecycle では戻り先だけを扱うとよい。
- `conflict hint`: 詳細な provider 失敗や secret 失敗は external-integration 観点と競合する可能性がある。

### CAND-MPC-LC-008 Storybook 確認は部品 fixture の作成から review 記録まで完了する

- `source requirement`: Storybook 基盤では story と fixture を `frontend/src/ui/**/__fixtures__` に置き、Storybook review URL と確認状態を task-local `storybook-review.md` に記録する。
- `viewpoint`: lifecycle / Storybook POC の作成から確認
- `candidate scenario id`: `CAND-MPC-LC-008`
- `actor`: 実装後に Storybook で見た目を確認する人間レビュアー
- `trigger`: マスターペルソナ画面の主要部品 story と fixture が作成され、Storybook dev server または static build が実行される。
- `expected outcome`: 主要表示状態の story ID、review URL、確認状態、未確認理由、build-storybook 結果が task-local に残る。
- `observable point`: Storybook の対象 story、`npm --prefix frontend run build-storybook` 結果、`storybook-review.md` の review URL と確認状態が確認できる。
- `related detail requirement type`: `testability_requirement`, `observability_requirement`, `compatibility_requirement`
- `adoption hint`: Storybook 基盤 task の `AIModelSelectionCard` 確認経路を踏襲できる。
- `conflict hint`: Storybook review 記録を `ui-design.md` へ置くか別 artifact へ置くかの未決と競合する。

## Open Notes

- `human decision candidate`: 現行 view model だけで、未設定、生成中、生成成功、生成失敗、編集中をすべて表現できるか。
- `human decision candidate`: ページ全体 story の最小合成単位を page component にするか、review-only wrapper にするか。
- `human decision candidate`: 既存部品を残す範囲と、新規部品として切り出す範囲。
- `human decision candidate`: Storybook review URL、確認状態、未確認状態を `ui-design.md` と別 artifact のどちらに残すか。
- `merge candidate`: CAND-MPC-LC-001 と CAND-MPC-LC-008 は Storybook POC の初期確認として統合できる可能性がある。
- `merge candidate`: CAND-MPC-LC-005 と CAND-MPC-LC-006 は modal lifecycle として統合できる可能性がある。
- `rejection candidate`: CAND-MPC-LC-007 の失敗詳細は、failure 観点の候補へ寄せる場合は lifecycle 側から除外できる。
