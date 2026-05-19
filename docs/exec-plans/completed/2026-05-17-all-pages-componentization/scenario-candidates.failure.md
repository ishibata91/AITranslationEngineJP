# Scenario Candidates: 2026-05-17-all-pages-componentization / failure

- `generator`: `failure`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `APF`

## Generator Scope

- `viewpoint`: 失敗観点
- `included_sources`: `plan.md`, `docs/spec.md`, `docs/architecture.md`, `docs/coding-guidelines-frontend.md`, `docs/UX-standard.md`, `docs/screen-design/README.md`, `docs/screen-design/screens/README.md`, `docs/screen-design/screens/*.md`, `docs/detail-specs/README.md`, `docs/scenario-tests/README.md`, `docs/exec-plans/completed/2026-05-18-storybook-foundation/scenario-design.md`, `docs/exec-plans/completed/2026-05-17-master-persona-componentization/scenario-design.md`, `docs/exec-plans/completed/2026-05-17-master-persona-componentization/ui-design.md`, `frontend/src/ui/screens/`, `frontend/src/ui/components/`
- `excluded_sources`: プロダクトコード変更、プロダクトテスト変更、docs 正本変更、他 agent の候補成果物、採否判断、統合判断、競合解消、最終シナリオ表
- `generation_notes`: 全ページの部品化と Storybook story 追加について、表示欠落、props 分割失敗、story fixture drift、既存表示項目の削除、`frontend-local` 失敗を候補化する。Storybook は fixed props または view model fixture に閉じ、Gateway mock、backend DTO mock、実行フロー再現は候補に含めない。

## Candidate Scenarios

### CAND-APF-001 主要表示領域の story が欠落する

- `source requirement`: `plan.md` の goal と close_conditions。各ページの主要表示領域は props 境界と Storybook story を持つ。
- `viewpoint`: 参照不能、保存失敗
- `candidate scenario id`: `CAND-APF-001`
- `actor`: 実装者
- `failure start condition`: 画面ごとのパネル、カード、モーダル候補のうち、主要表示領域に対応する story がない。
- `rejected operation`: 全ページの主要表示領域を Storybook review 済みとして扱う。
- `expected error`: story 欠落として記録し、対象ページの見た目レビューを完了扱いにしない。
- `expected outcome`: 欠落した表示領域、対象画面、未確認理由、再確認入口が task-local 成果物に残る。
- `observable point`: story 一覧、`storybook-review.md` または `ui-design.md#storybook-review`、対象画面の screen local component、Storybook review URL。
- `related detail requirement type`: `failure_handling_requirement`, `testability_requirement`, `observability_requirement`
- `adoption hint`: 採否は `designer` が判断する。全ページ一括 task の story 網羅不足を検出する候補である。
- `conflict hint`: page 合成 story を必須にする候補と、主要表示領域 story だけを必須にする候補で検証単位が競合する可能性がある。

### CAND-APF-002 props 分割後の component が controller または Store を要求する

- `source requirement`: `docs/architecture.md` の UI Component 責務。`docs/coding-guidelines-frontend.md` は UI Component が `Store`、Gateway、generated binding を直接扱わないと定めている。
- `viewpoint`: 参照不能、設定不整合
- `candidate scenario id`: `CAND-APF-002`
- `actor`: 実装者
- `failure start condition`: 部品化した panel、card、modal が controller factory、`Store`、Gateway、generated `wailsjs`、RuntimeEventAdapter のいずれかを直接要求する。
- `rejected operation`: fixed props または view model fixture だけで story を render する。
- `expected error`: Storybook render または build が失敗する。失敗しない場合でも、props 境界違反として扱う。
- `expected outcome`: component は props と callback だけで表示できる境界へ戻す。
- `observable point`: component import、story import、fixture import、Storybook iframe render 結果、`npm --prefix frontend run build-storybook` の結果。
- `related detail requirement type`: `failure_handling_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 採否は `designer` が判断する。ページ component を薄くする目的に対する基本失敗候補である。
- `conflict hint`: page component 全体を story 化して runtime 依存を許す候補がある場合、この候補の props-only 境界と競合する。

### CAND-APF-003 props 分割で操作可否と disabled 理由がずれる

- `source requirement`: `docs/UX-standard.md` は状態表示、状態別操作、表示条件を high priority にしている。画面設計書は操作の有効条件、無効条件、失敗時表示を画面ごとに定義している。
- `viewpoint`: 設定不整合、失敗入力
- `candidate scenario id`: `CAND-APF-003`
- `actor`: 利用者
- `failure start condition`: 親画面の view model から小さい props へ分ける時に、操作可否、無効理由、処理中状態のいずれかが欠落または別値になる。
- `rejected operation`: 禁止状態の操作を有効表示する、または有効状態の操作を理由なしで無効表示する。
- `expected error`: 操作可否の不整合として story または `frontend-local` の確認を失敗扱いにする。
- `expected outcome`: disabled 状態、理由文、処理中表示、callback 発火可否が同じ fixture 内で整合する。
- `observable point`: 各画面の操作 button、disabled 属性、無効理由表示、callback stub の呼び出し有無、対象 story の args。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `failure_handling_requirement`, `testability_requirement`
- `adoption hint`: 採否は `designer` が判断する。props 分割で見た目だけが合い、操作条件が壊れる失敗を拾う候補である。
- `conflict hint`: 操作条件を component 内で再計算する候補と、親画面 view model で計算済みにする候補が責務境界で競合する。

### CAND-APF-004 story fixture が画面設計書の状態別表示から drift する

- `source requirement`: `docs/screen-design/screens/README.md` は表示項目、状態別表示、操作、依存情報を画面設計書に書くと定めている。完了済み Storybook 基盤は fixture を component 横の `__fixtures__` に置く方針を固定している。
- `viewpoint`: 参照不能、設定不整合
- `candidate scenario id`: `CAND-APF-004`
- `actor`: 実装者
- `failure start condition`: story fixture が画面設計書の loading、empty、error、disabled、selected、progress の状態を表現できない、または古い表示文言や古い項目構成を持つ。
- `rejected operation`: fixture で表せない状態を Storybook 確認済みとして扱う。
- `expected error`: fixture drift として未確認理由を残し、該当状態の story を確認済みにしない。
- `expected outcome`: fixture は画面設計書の状態別表示と対応し、古い fixture は更新または除外候補として分けられる。
- `observable point`: `frontend/src/ui/**/__fixtures__`、story args、画面設計書の状態別表示、Storybook canvas、review 記録。
- `related detail requirement type`: `state_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 採否は `designer` が判断する。全ページを一括で部品化する時の fixture drift 検出候補である。
- `conflict hint`: 画面設計書に未定義の見た目状態を fixture で補う候補は、docs 正本化未承認の状態追加と競合する可能性がある。

### CAND-APF-005 既存表示項目が人間承認なしに削除される

- `source requirement`: `plan.md` の未決事項。既存表示項目を削る場合の人間承認粒度は未決である。`docs/screen-design/screens/README.md` は表示項目を画面設計書へ固定する。
- `viewpoint`: 保存失敗、競合候補、人間判断候補
- `candidate scenario id`: `CAND-APF-005`
- `actor`: 人間レビュー担当者
- `failure start condition`: 部品化後の実画面または story から、画面設計書にある表示項目、状態文、操作、`aria-label` が消える。
- `rejected operation`: 表示項目の削除を refactor の副作用として完了扱いにする。
- `expected error`: 既存表示項目の削除候補として記録し、人間承認または screen design diff がない限り完了扱いにしない。
- `expected outcome`: 削除対象、根拠画面設計書、影響する story、削除承認の有無が task-local 成果物に残る。
- `observable point`: screen design の表示項目、component markup、Storybook screenshot、`aria-label`、review 記録。
- `related detail requirement type`: `compatibility_requirement`, `data_requirement`, `observability_requirement`
- `adoption hint`: 採否は `designer` が判断する。既存表示を維持する refactor の回帰候補である。
- `conflict hint`: UI 簡素化として表示項目を削る候補がある場合、人間承認粒度と docs 正本化の扱いが競合する。

### CAND-APF-006 共有 component 化で画面固有条件が props 分岐の塊になる

- `source requirement`: `docs/architecture.md` の UI Component 判断表。共有 component は複数画面で使う部品だけにする。`docs/coding-guidelines-frontend.md` は画面固有条件を shared component の props 分岐として増やし続けないと定めている。
- `viewpoint`: 設定不整合、競合候補
- `candidate scenario id`: `CAND-APF-006`
- `actor`: 実装者
- `failure start condition`: 画面専用の大きなレイアウト、行、カード、modal を `frontend/src/ui/components/` へ上げ、画面固有条件を variant や props 分岐で吸収する。
- `rejected operation`: 共有 component として全ページの story 確認対象へ含める。
- `expected error`: 共有化判断の失敗として扱い、screen local component へ戻す候補にする。
- `expected outcome`: 共有 component は UI 規則の集約に限り、画面固有の業務条件は screen local component に残る。
- `observable point`: `frontend/src/ui/components/` の props 型、variant 数、画面固有文言、story args、利用箇所。
- `related detail requirement type`: `consistency_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 採否は `designer` が判断する。全ページ一括部品化で過剰共通化を検出する候補である。
- `conflict hint`: 共通化候補を広く取る候補と、画面専用部品へ残す候補が責務境界で競合する。

### CAND-APF-007 story または fixture に保存禁止情報が混入する

- `source requirement`: 完了済み Storybook 基盤と Master Persona 成果物は、story と fixture に secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータ、raw request、raw response、raw prompt を含めない方針を固定している。
- `viewpoint`: 保存失敗、設定不整合
- `candidate scenario id`: `CAND-APF-007`
- `actor`: 実装者
- `failure start condition`: story、fixture、review URL、screenshot 名、確認記録に保存禁止情報が入る。
- `rejected operation`: Storybook review 証跡として保存する。
- `expected error`: 安全条件違反として扱い、該当 story または review 記録を確認済みにしない。
- `expected outcome`: story と fixture は synthetic data、状態値、callback stub だけを持つ。
- `observable point`: story source、fixture source、Storybook URL、review 記録、screenshot path、command 記録。
- `related detail requirement type`: `security_requirement`, `observability_requirement`, `testability_requirement`
- `adoption hint`: 採否は `designer` が判断する。全ページ fixture の安全境界を固定する候補である。
- `conflict hint`: 実在入力サンプルを fixture に使う候補と、synthetic fixture だけを使う候補が安全条件で競合する。

### CAND-APF-008 `build-storybook` が story 追加後に失敗する

- `source requirement`: `plan.md` の validation_commands。`npm --prefix frontend run build-storybook` が検証入口に含まれる。
- `viewpoint`: 設定不整合、参照不能
- `candidate scenario id`: `CAND-APF-008`
- `actor`: 実装者
- `failure start condition`: 新規 story、fixture、component import、Vite alias、Svelte compile のいずれかが不整合になる。
- `rejected operation`: Storybook static build を通過扱いにして人間見た目レビューへ進める。
- `expected error`: build command が失敗し、対象 story の静的出力を作れない。
- `expected outcome`: 失敗 command、短い原因、対象 story、再実行 command、未確認理由が task-local 成果物に残る。
- `observable point`: `npm --prefix frontend run build-storybook` の exit code、error log、対象 story path、Storybook static output。
- `related detail requirement type`: `failure_handling_requirement`, `testability_requirement`, `observability_requirement`
- `adoption hint`: 採否は `designer` が判断する。Storybook review を close condition に含める場合の検証失敗候補である。
- `conflict hint`: `build-storybook` を `frontend-local` に含めるか別 gate にするかで検証段階が競合する可能性がある。

### CAND-APF-009 `frontend-local` が props 型、lint、既存画面テストで失敗する

- `source requirement`: `plan.md` の close_conditions と validation_commands。実装後の Storybook review と `frontend-local` は通過または未通過理由を持つ。
- `viewpoint`: 設定不整合、保存失敗
- `candidate scenario id`: `CAND-APF-009`
- `actor`: 実装者
- `failure start condition`: 部品化と story 追加後に、TypeScript 型、Svelte check、lint、既存 frontend test のいずれかが失敗する。
- `rejected operation`: `frontend-local` 未通過のまま実装完了または人間レビュー準備完了として扱う。
- `expected error`: `frontend-local` 失敗として、通過条件を満たさないことを記録する。
- `expected outcome`: 失敗 suite、短い原因、影響画面、再実行 command、未通過理由が task-local 成果物に残る。
- `observable point`: `python3 scripts/harness/run.py --suite frontend-local` の exit code、失敗 log、対象 file、既存 screen test の結果。
- `related detail requirement type`: `failure_handling_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 採否は `designer` が判断する。表示が Storybook で見えても既存 frontend gate を壊す失敗候補である。
- `conflict hint`: `frontend-local` と `build-storybook` のどちらを必須 gate にするかで完了条件の表現が競合する可能性がある。

### CAND-APF-010 Storybook review 証跡が不足して人間レビュー対象を再現できない

- `source requirement`: `plan.md` の close_conditions。完了済み Storybook 基盤と Master Persona 成果物は、review URL、story ID、確認状態、未確認理由、再実行 command を task-local に残す方針を固定している。
- `viewpoint`: 保存失敗、回復動作
- `candidate scenario id`: `CAND-APF-010`
- `actor`: 人間レビュー担当者
- `failure start condition`: story は存在するが、review URL、story ID、確認状態、未確認理由、再実行 command が task-local 成果物に残らない。
- `rejected operation`: 人間レビュー担当者が同じ Storybook 表示へ到達する。
- `expected error`: 見た目レビュー未準備として扱い、確認済み状態へ進めない。
- `expected outcome`: 対象 story、URL、確認状態、未確認理由、再実行 command が後から追跡できる。
- `observable point`: `storybook-review.md` または `ui-design.md#storybook-review`、Storybook URL、iframe URL、story ID、確認状態。
- `related detail requirement type`: `observability_requirement`, `testability_requirement`, `recovery_requirement`
- `adoption hint`: 採否は `designer` が判断する。全ページ一括 task の人間見た目レビュー入口を守る候補である。
- `conflict hint`: review 記録を `ui-design.md` に置くか別 artifact に置くかで記録先が競合する可能性がある。

## Open Notes

- `human decision candidate`: 既存表示項目を削る場合の人間承認粒度は `plan.md` の未決事項である。
- `human decision candidate`: page component の story を必須にするか、panel、card、modal の story だけを必須にするかは最終シナリオ統合時に決める必要がある。
- `human decision candidate`: Storybook review 記録先を `ui-design.md#storybook-review` と別 artifact のどちらにするかは `plan.md` 上で両方許容されている。
- `merge candidate`: `CAND-APF-001` と `CAND-APF-010` は Storybook review 網羅と証跡欠落として統合される可能性がある。
- `merge candidate`: `CAND-APF-002` と `CAND-APF-003` は props 境界失敗として統合される可能性がある。
- `merge candidate`: `CAND-APF-004` と `CAND-APF-005` は既存 screen design との drift として統合される可能性がある。
- `rejection candidate`: 正常系の裏返しだけで、失敗開始条件、拒否される操作、期待エラー、観測点を持たない候補は採否判断前に除外対象になりうる。
- `rejection candidate`: Storybook に Gateway mock、backend DTO mock、実行フロー再現を入れる候補は、完了済み Storybook 基盤の境界と衝突するため除外候補である。

## Completion Summary

- `viewpoint`: 失敗観点
- `candidate_count`: 10
- `artifact_path`: `docs/exec-plans/active/2026-05-17-all-pages-componentization/scenario-candidates.failure.md`
- `task_artifact_root`: `docs/exec-plans/active/2026-05-17-all-pages-componentization/`
- `target_diff`: 全ページのページ component を薄くし、パネル、カード、モーダル単位へ部品化し、主要部品へ Storybook story を追加する。
- `remaining_risk`: 採否、統合、競合解消、質問票化は `designer` が行う。
