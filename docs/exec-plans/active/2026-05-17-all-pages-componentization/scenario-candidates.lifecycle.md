# Scenario Candidates: 2026-05-17-all-pages-componentization / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `APC-LC`
- `candidate_count`: `9`

## Generator Scope

- `viewpoint`: lifecycle
- `included_sources`: `./plan.md`, `docs/spec.md`, `docs/architecture.md`, `docs/er.md`, `docs/coding-guidelines-frontend.md`, `docs/UX-standard.md`, `docs/screen-design/README.md`, `docs/detail-specs/README.md`, `docs/scenario-tests/README.md`, `docs/exec-plans/completed/2026-05-18-storybook-foundation/scenario-design.md`, `docs/exec-plans/completed/2026-05-17-master-persona-componentization/scenario-design.md`, `docs/exec-plans/completed/2026-05-17-master-persona-componentization/ui-design.md`, `frontend/src/ui/screens/README.md`, `frontend/src/ui/screens/`, `frontend/src/ui/components/`
- `excluded_sources`: プロダクトコード変更、プロダクトテスト変更、docs 正本化、最終シナリオ採否、他観点の候補統合、他 agent の成果物
- `generation_notes`: 全ページの部品化と Storybook story 追加について、作成、更新、レビュー、完了、再実行の流れだけを候補化する。Storybook は固定 props または view model fixture に閉じ、backend、Wails runtime、Gateway、generated binding、AI provider、secret store、DB、実 filesystem flow を要求しない。

## Candidate Scenarios

### CAND-APC-LC-001 全ページの部品化候補を作成する

- `source requirement`: `plan.md` は、全ページの表示領域を事前に調べ、パネル、カード、モーダル単位の部品化候補を列挙するとしている。
- `viewpoint`: lifecycle / 作成
- `candidate scenario id`: `CAND-APC-LC-001`
- `lifecycle stage`: 作成
- `start condition`: Storybook 基盤と Master Persona POC の成果物を読める。`frontend/src/ui/screens/` と `frontend/src/ui/components/` の現行構成を確認できる。
- `actor`: design bundle を作る担当者
- `trigger`: 全ページ部品化 task の設計成果物作成を開始する。
- `expected outcome`: 画面ごとに、共有部品候補、画面専用部品候補、分けない候補が分類される。
- `observable point`: `ui-design.md` または `component-candidates.md` に、対象画面、候補名、候補分類、分ける理由、分けない理由が残る。
- `related detail requirement type`: `testability_requirement`, `compatibility_requirement`, `data_requirement`
- `adoption hint`: 実装前の棚卸し確認として採用候補になる。
- `conflict hint`: 部品化候補の画面順序は `plan.md` で未決であるため、最終シナリオでは順序を designer が固定する必要がある。

### CAND-APC-LC-002 ページ component を薄い合成役へ更新する

- `source requirement`: `plan.md` の `goal` は、全ページのページ component を薄くし、主要表示領域に props 境界を持たせるとしている。
- `viewpoint`: lifecycle / 更新
- `candidate scenario id`: `CAND-APC-LC-002`
- `lifecycle stage`: 更新
- `start condition`: 対象ページの部品化候補が承認され、画面専用部品の切り出し先が決まっている。
- `actor`: Codex implementation lane
- `trigger`: 対象ページの component 分割を実装する。
- `expected outcome`: ページ component は controller 接続、購読、dispose、画面専用部品の合成を主に持つ。切り出した部品は props と callback で表示できる。
- `observable point`: 対象 page の import、props 型、callback、controller 参照の有無、Storybook で描画できる props 境界。
- `related detail requirement type`: `compatibility_requirement`, `state_requirement`, `testability_requirement`
- `adoption hint`: 全ページ共通の更新 lifecycle として採用候補になる。
- `conflict hint`: 親画面状態を大量に読む部品、業務フロー全体を進める部品、props が条件分岐の塊になる部品は、分けない候補と競合する。

### CAND-APC-LC-003 共有部品候補を更新または維持する

- `source requirement`: `docs/architecture.md` と `docs/coding-guidelines-frontend.md` は、複数画面で使う部品だけを `frontend/src/ui/components/` に置くとしている。
- `viewpoint`: lifecycle / 更新
- `candidate scenario id`: `CAND-APC-LC-003`
- `lifecycle stage`: 更新
- `start condition`: 複数画面で同じ表示規則、操作、状態表示を持つ候補が見つかっている。
- `actor`: Codex implementation lane
- `trigger`: 画面専用部品を共有部品へ上げるか判断する。
- `expected outcome`: 共有部品に上げる候補は、画面固有条件を props 分岐として増やし続けない範囲に限られる。共有にしない候補は screen local に残る。
- `observable point`: `frontend/src/ui/components/` の追加または既存部品利用、screen local に残す理由、`AIModelSelectionCard` と `StickyActionFooter` の再利用可否。
- `related detail requirement type`: `compatibility_requirement`, `testability_requirement`, `consistency_requirement`
- `adoption hint`: shared component への昇格判断を検査する候補として採用できる。
- `conflict hint`: 共通化候補を shared component へ上げる具体基準は `plan.md` で未決である。

### CAND-APC-LC-004 主要部品の Storybook story と fixture を作成する

- `source requirement`: `plan.md` の `close_conditions` は、各ページの主要表示領域が props 境界と Storybook story を持つことを求めている。Storybook 基盤の完了成果物は、fixture を固定 props または view model fixture に限定している。
- `viewpoint`: lifecycle / 作成
- `candidate scenario id`: `CAND-APC-LC-004`
- `lifecycle stage`: 作成
- `start condition`: 対象部品が props と callback で描画できる状態になっている。
- `actor`: Codex implementation lane
- `trigger`: パネル、カード、モーダル単位の story と fixture を追加する。
- `expected outcome`: 主要部品の代表状態が Storybook で表示できる。fixture は外部接続や実データを持たない。
- `observable point`: story source、fixture source、story ID、props 入力、callback stub、禁止 import の不在。
- `related detail requirement type`: `testability_requirement`, `data_requirement`, `compatibility_requirement`
- `adoption hint`: 部品化と Storybook story 追加の中心シナリオとして採用候補になる。
- `conflict hint`: page story は密度確認と配置確認に限定されるため、主要表示領域 story の不足を page story で代替できない。

### CAND-APC-LC-005 既存表示項目と操作状態を更新後に維持する

- `source requirement`: `docs/UX-standard.md` は画面目的、状態、操作、制約、結果を判断できる構造を求めている。`plan.md` は既存表示項目の維持確認を必要としている。
- `viewpoint`: lifecycle / 更新後確認
- `candidate scenario id`: `CAND-APC-LC-005`
- `lifecycle stage`: 更新後確認
- `start condition`: 対象ページの部品化と story 作成が完了している。
- `actor`: 実装後確認者
- `trigger`: Storybook または画面確認で、更新後の表示項目と操作状態を確認する。
- `expected outcome`: 部品化前に画面で確認できた主要表示項目、主要 CTA、補助 CTA、disabled 状態、empty、loading、error、success の表現が失われない。
- `observable point`: screen design、対象 component、Storybook story、画面表示の表示項目対応。
- `related detail requirement type`: `compatibility_requirement`, `state_requirement`, `consistency_requirement`
- `adoption hint`: 見た目リファクタによる表示欠落を防ぐ候補として採用できる。
- `conflict hint`: 既存表示項目を削る場合の人間承認粒度は `plan.md` で未決である。

### CAND-APC-LC-006 Storybook で人間見た目レビューを行う

- `source requirement`: `plan.md` は、実装後の見た目レビューを Storybook で行うとしている。Master Persona POC は、review URL、story ID、確認状態、未確認理由を task-local に残す方針を持つ。
- `viewpoint`: lifecycle / レビュー
- `candidate scenario id`: `CAND-APC-LC-006`
- `lifecycle stage`: レビュー
- `start condition`: 主要部品 story と fixture が作成され、Storybook dev server または static preview を開ける。
- `actor`: 人間レビュアー
- `trigger`: Storybook review URL から対象 story を確認する。
- `expected outcome`: 人間レビュアーは、主要部品の見た目、表示密度、長文、狭い幅、状態差分を確認できる。
- `observable point`: Storybook URL、story ID、確認状態、未確認理由、必要な再確認対象。
- `related detail requirement type`: `testability_requirement`, `observability_requirement`, `compatibility_requirement`
- `adoption hint`: frontend human review 前の Storybook 確認として採用候補になる。
- `conflict hint`: review 記録先は `ui-design.md#storybook-review` または `storybook-review.md` と plan に揺れがある。Storybook 基盤と Master Persona POC は `storybook-review.md` を使っている。

### CAND-APC-LC-007 Storybook build と frontend-local で完了判定する

- `source requirement`: `plan.md` の `validation_commands` は `python3 scripts/harness/run.py --suite frontend-local` と `npm --prefix frontend run build-storybook` を求めている。
- `viewpoint`: lifecycle / 完了
- `candidate scenario id`: `CAND-APC-LC-007`
- `lifecycle stage`: 完了
- `start condition`: 全対象ページの部品化、主要 story、review 記録の作成が完了している。
- `actor`: Codex implementation lane
- `trigger`: 実装後検証を実行する。
- `expected outcome`: Storybook static build と frontend-local が通過する。または未通過理由が task-local に残る。
- `observable point`: command、exit code、未通過理由、Storybook build output、frontend-local 結果。
- `related detail requirement type`: `testability_requirement`, `compatibility_requirement`
- `adoption hint`: task close condition の検証候補として採用できる。
- `conflict hint`: Storybook build は通常 frontend-local と別 gate で扱われるため、検証順序と失敗時の戻し先を designer が整理する必要がある。

### CAND-APC-LC-008 完了後に実装結果を design bundle へ戻す

- `source requirement`: `plan.md` は、Storybook review、frontend human review、approved frontend protection、implementation result を closeout に残すとしている。
- `viewpoint`: lifecycle / 完了後引き継ぎ
- `candidate scenario id`: `CAND-APC-LC-008`
- `lifecycle stage`: 完了後引き継ぎ
- `start condition`: 実装、Storybook review、frontend-local、build-storybook の結果が揃っている。
- `actor`: implement_lane
- `trigger`: 全ページ部品化 task の完了判断材料を designer または人間へ戻す。
- `expected outcome`: 変更された部品、story ID、fixture 方針、review 結果、未確認理由、残留リスクが task-local 成果物から追跡できる。
- `observable point`: `storybook-review.md` または `ui-design.md#storybook-review`、実装結果欄、review 証跡、検証結果。
- `related detail requirement type`: `observability_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: ライフサイクルの終点と後続判断材料として採用候補になる。
- `conflict hint`: docs 正本化が必要な差分は、この候補では確定しない。docs 正本化は別レーンへ渡す。

### CAND-APC-LC-009 レビュー指摘後に story と検証を再実行する

- `source requirement`: `plan.md` は、frontend human review と approved frontend protection を実装後に要求している。
- `viewpoint`: lifecycle / 再実行
- `candidate scenario id`: `CAND-APC-LC-009`
- `lifecycle stage`: 再実行
- `start condition`: Storybook review、UX review、frontend human review、frontend-local、build-storybook のいずれかで指摘または未通過理由が残っている。
- `actor`: Codex implementation lane
- `trigger`: 指摘に対応して部品、fixture、story、review 記録を更新する。
- `expected outcome`: 指摘対象の story と関連する検証を再実行し、更新後の確認状態と残留リスクが task-local に残る。
- `observable point`: 修正対象 story ID、再実行 command、再実行結果、未解決指摘、更新後 review 記録。
- `related detail requirement type`: `recovery_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: レビュー戻りの lifecycle 候補として採用できる。
- `conflict hint`: 再実行で表示項目削除や画面設計差分が発生する場合は、人間承認または docs updater の判断へ戻す必要がある。

## Open Notes

- `human decision candidate`: 部品化候補を確認する画面順序。
- `human decision candidate`: 共通化候補を `frontend/src/ui/components/` へ上げる具体基準。
- `human decision candidate`: 既存表示項目を削る場合の人間承認粒度。
- `human decision candidate`: Storybook review 記録先を `ui-design.md#storybook-review` と `storybook-review.md` のどちらへ固定するか。
- `merge candidate`: `CAND-APC-LC-002` と `CAND-APC-LC-004` は、部品更新と story 作成を同じ実装サイクルとして統合できる可能性がある。
- `merge candidate`: `CAND-APC-LC-006`、`CAND-APC-LC-007`、`CAND-APC-LC-008` は、実装後 review から完了判定までの closeout シナリオとして統合できる可能性がある。
- `rejection candidate`: backend DTO mock、Gateway mock、generated binding、Wails runtime、AI provider、secret store、DB、実 filesystem flow を Storybook story の前提にする候補は除外する。
- `rejection candidate`: docs 正本化またはプロダクトコード変更の指示を含む候補は、この lifecycle 候補成果物から除外する。
