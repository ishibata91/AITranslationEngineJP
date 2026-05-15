# 翻訳ジョブ状態 stale 廃止 plan

- `task_id`: `2026-05-16-translation-job-state-stale-retirement`
- `lane`: `light-change-lane`
- `status`: `active`
- `created_at`: `2026-05-16`
- `human_request`: 翻訳ジョブ状態関連の追加差分を、stale 廃止を主目的に整理する。

## 目的

翻訳ジョブ状態関連の実装から、旧設計名、空 package、重複した phase 別 wrapper、古い task-local 参照を廃止する。
新しい状態、永続値、公開 DTO、画面仕様は追加しない。
コード削減を主目的にし、状態可否の意味は現在の `TranslationJobPolicy` に合わせる。

この plan での `stale` は、古い設計名、使われていない package、同じ規則を繰り返す wrapper、完了済み判断とずれた active task-local 参照を意味する。
`stale_selection`、`validation_stale`、runtime event の stale 判定、translation output artifact の stale 表示はドメイン仕様なので廃止対象にしない。

## 判断

- 判定: `範囲内修正`
- 戻し先: `light_change_lane`
- 実装種別: `implement-backend`
- docs 正本化: 実装結果が architecture 正本の構造主語を変える場合だけ `docs_updater` に渡す

理由:

- 廃止対象は既存仕様の意味を広げない。
- 追加する状態遷移、DB カラム、Wails 公開契約はない。
- 目的は、既存の `TranslationJobPolicy` と現在の実装の不一致を減らすことである。
- `refactor_lane` は未定義なので、軽量変更として対象を限定する。

## 対象

廃止候補:

- `internal/statemachine/`: `doc.go` だけの旧設計 package。現在の正本名は `TranslationJobPolicy` である。
- `internal/jobio/`: `doc.go` だけの旧設計 package。実装がないまま architecture 正本だけに残っているため、廃止または正本再承認が必要である。
- `.go-arch-lint.yml`: `statemachine` と `jobio` component、許可依存。空 package 廃止後は残す理由がない。
- `internal/usecase/*_phase_usecase.go`: `evaluateTermPolicy`、`evaluatePersonaPolicy`、`evaluateBodyPolicy` と、同じ policy input 生成を繰り返す helper。
- `internal/service/*_phase_service.go`: `CanPause`、`CanResume`、`CanRetry`、`CanCancel` を phase 別に手書きする action enablement 分岐。
- `docs/exec-plans/active/observability-log-addition/`: `StateMachine` / `JobIOService` を前提にした active task-local 参照。

禁止対象:

- `docs/exec-plans/completed/**`: 完了済み履歴なので書き換えない。
- `stale_selection`、`validation_stale`、`model_selection_stale`: 利用者向けまたは API 向けの理由分類なので消さない。
- `docs/detail-specs/*` の状態意味: 実装削減だけを理由に変えない。
- provider 応答、credential、prompt、翻訳本文をログへ増やさない。

## 進め方

1. stale 候補を棚卸しする。
   - `rg` で旧設計名、空 package、phase 別 wrapper を列挙する。
   - ドメイン仕様の `stale_*` は削除候補から除外する。

2. 空 package と lint component を廃止する。
   - `internal/statemachine/` を削除する。
   - `internal/jobio/` は architecture 正本を同時に直す判断が通った場合だけ削除する。
   - `.go-arch-lint.yml` から廃止 package の component と許可依存を消す。

3. UseCase の phase 別 policy wrapper を畳む。
   - policy input 生成を共通 helper に寄せる。
   - phase ごとの差分は、phase 名、error kind 対応、service method 呼び出しだけにする。
   - `termPolicyPhaseRunMatches` のような phase 固有重複 helper を消す。

4. Service の操作可否分岐を policy 由来へ寄せる。
   - `pause`、`resume`、`retry`、`cancel` の可否は `TranslationJobPolicy` の共通操作規則から導出する。
   - phase service は開始前提、read model 集約、provider 実行だけを残す。
   - terminal job の操作不可理由は policy と read model で同じ分類にする。

5. active task-local の旧名を直す。
   - `observability-log-addition` の active 成果物に残る `StateMachine` は `TranslationJobPolicy` へ置き換える。
   - `JobIOService` を廃止する判断なら、状態事実の取得と保存の観測点を usecase/service/repository 境界へ書き換える。
   - completed archive は履歴として残す。

## 成果物依存表

| 成果物ID | 状態 | 担当 | 依存対象 | 出力 |
| --- | --- | --- | --- | --- |
| `task 枠` | 完了 | `light_change_lane` | なし | `plan.md` |
| `軽量変更計画` | 完了 | `light_change_planner` | `task 枠` | `light-change-planning.md` |
| `設計差分図` | 未着手 | `diagrammer` | `軽量変更計画` | component / sequence の差分図 |
| `実装証跡` | 未着手 | `backend_implementer` | `軽量変更計画`, `設計差分図` | stale 廃止差分 |
| `人間確認` | 未着手 | 人間 | `実装証跡` | 削減範囲の確認 |
| `テスト修正証跡` | 未着手 | `implementation_unit_tester` | `実装証跡`, `人間確認?` | policy / read model の単体テスト整理 |
| `実装後ブラウザ確認` | 該当なし | `light_change_lane` | `実装証跡`, `テスト修正証跡?` | UI 変更なし |
| `レビュー通過根拠` | 未着手 | `light_change_lane` | `実装証跡`, `テスト修正証跡?` | 5 観点 reviewback 集約 |
| `正本化判断` | 未着手 | `light_change_lane` | `レビュー通過根拠` | docs 正本反映要否 |
| `詳細仕様正本反映` | 条件付き未着手 | `docs_updater` | `正本化判断` | architecture / active docs の同期 |
| `作業レポート入力` | 未着手 | `light_change_lane` | 全完了または停止済み成果物 | work reporter 向け入力 |
| `作業計画完了移動` | 未着手 | `light_change_lane` | `作業レポート入力` | completed への移動 |

## 検証

実装 agent は次を通す。

- `gofmt -l internal/usecase internal/service`
- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite backend-lint`
- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite coverage`

追加の確認:

- `rg -n "internal/statemachine|StateMachine" internal docs .go-arch-lint.yml --glob '!docs/exec-plans/completed/**' --glob '!work_history/**'`
- `rg -n "internal/jobio|JobIOService" internal docs .go-arch-lint.yml --glob '!docs/exec-plans/completed/**' --glob '!work_history/**'`
- `rg -n "translationjobpolicy" internal/service internal/repository internal/controller internal/infra`

## 未決事項

- `JobIOService` を廃止して architecture 正本から外すか、別 task で実装するかを人間が確認する必要がある。
- `observability-log-addition` の active 成果物をこの task で更新するか、observability task 再開時に更新するかを決める必要がある。
- policy 由来の action enablement へ寄せた後、phase service 側に残す拒否理由文字列の粒度を決める必要がある。

## 着手可能成果物

`設計差分図` が着手可能である。
図では、削除する package、残す policy、UseCase から service への呼び出し範囲だけを示す。

