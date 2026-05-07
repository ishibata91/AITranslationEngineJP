# Scenario Design: 2026-05-07-model-settings-card-controller

- `skill`: scenario-design
- `status`: approved
- `source_plan`: `./plan.md`
- `task_frame`: `./task-frame.md`
- `ui_source`: `./ui-design.md`
- `final_artifact_path`: `./scenario-design.md`
- `topic_abbrev`: `MSCC`
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## Fixed Requirements

- `must_pass_requirements`:
  - モデル設定カードの表示部品は、provider、model、model list、保存、取得、選択状態の制御を直接持たない。
  - マスターペルソナと翻訳ジョブ設定は、同じモデル設定カード制御を使う。
  - provider settings は endpoint と credential 参照状態の正本であり、model と処理方法の保存元にはしない。
  - 利用者向け provider list は `gemini`、`lm_studio`、`xai` だけを表示し、`fake` provider ID を表示しない。
  - fake mode では、選択中の通常 provider ID の model list 結果として `fake-model` を扱う。
  - frontend に fake mode 判定や `fake-model` 固有分岐を置かない。
  - APIキー本体、secret、raw request、raw response、raw prompt、内部ログ用識別子は UI、DTO、要約、log に出さない。
  - 遅延した model list 応答は、現在 provider と現在要求へ反映しない。
- `non_goals`:
  - AIサービス設定画面へ model、処理方法、Batch API 切り替えを追加しない。
  - `fake` provider ID を利用者向け選択肢や保存要約へ追加しない。
  - 有料の実 AI API 呼び出しを受け入れテストの前提にしない。
  - 永続監査ログまたは provider settings 更新履歴を新規要件にしない。
  - プロダクトコード、プロダクトテスト、docs 正本本文は変更しない。

## Scenario Candidate Coverage

正本: `./scenario-design.candidate-coverage.json`

6 種の候補成果物は作業計画フォルダに揃っている。
候補生成 agent は追加起動していない。

`needs_human_decision` は 0 件である。
未解決 conflict は 0 件である。
人間回答は `./scenario-design.questions.md` に反映済みである。

## Detail Requirement Coverage

正本: `./scenario-design.requirement-coverage.json`

抽象要件は詳細要求タイプへ展開した。
`needs_human_decision` は 0 件である。
人間回答により、保存単位、失敗回復、空一覧、APIキー未設定導線を固定した。

## Scenario Matrix

この表は人間回答を反映済みである。
`Q-MSCC-001` から `Q-MSCC-004` の回答により、参照側ごとの保存、保存失敗後の未保存維持、取得済み 0 件表示、共有カード内導線なしを固定した。

| scenario id | 受け入れユースケース | 実行者 | 開始条件 | 主要操作 | 期待結果 | 主要観測点 | 実行テスト種別 | 実行段階 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `SCN-MSCC-001` | 共有モデル設定カードの初期取得と再取得 | 利用者 | マスターペルソナ画面または翻訳ジョブ設定画面を開く | 画面を開く、または設定を再読込する | 参照側ごとの保存済み provider、model、credential 参照状態がカードへ反映され、相互に混入しない | provider、model、APIキー状態、保存済みまたは未保存表示 | UI人間操作E2E | 実装後 |
| `SCN-MSCC-002` | provider 変更後の model list 更新と model 選択 | 利用者 | カードに現在 provider と model 状態がある | provider を変更し、model list を更新し、model を選ぶ | 変更前 provider の model list と model は現在 provider へ混入しない | 更新中、取得済み、model 未選択、選択済み表示 | UI人間操作E2E | 実装後 |
| `SCN-MSCC-003` | fake mode で通常 provider ID のまま `fake-model` を選ぶ | 利用者または人間レビュー担当者 | fake mode の環境でカードを開く | 通常 provider を選び、model list を更新する | `fake` provider は表示されず、`fake-model` を model list 結果として選べる | provider list、model list、保存要約に fake 固有分岐が出ないこと | UI人間操作E2E | 実装後 |
| `SCN-MSCC-004` | provider ごとの credential 条件で model list 外部取得を制御する | 利用者 | APIキーが必要な provider または LM Studio を選ぶ | model list 更新可否を確認する | APIキー未設定では外部取得を開始せず、LM Studio は credential なしで更新できる | 更新ボタン可否、APIキー状態、adapter 呼び出し有無、secret 非露出 | APIテスト | 実装後 |
| `SCN-MSCC-005` | 遅延した model list 応答を現在状態へ反映しない | 利用者 | model list 更新中に provider を変更する | provider A 更新中に provider B へ切り替える | provider A の遅延応答は破棄され、provider B の状態だけが残る | 現在 provider、model list、model 選択、古い候補が表示されないこと | UI人間操作E2E | 実装後 |
| `SCN-MSCC-006` | provider と model の選択状態を保存し再取得する | 利用者 | provider と model を選んでいる | 保存し、画面を再表示する | 参照側の provider と model 選択状態を保存し、再取得で復元する | 保存中、保存済み、未保存、再取得結果、secret 非露出 | UI人間操作E2E | 実装後 |
| `SCN-MSCC-007` | 翻訳ジョブ設定で 3 翻訳段階の不足がない時だけ job 作成へ進む | 利用者 | Job Setup で input と 3 翻訳段階が表示されている | 各段階の provider と model を選ぶ | APIキー不足と model 未選択がない時だけ Ready job を作成できる | job 作成可否、作成前確認、作成後設定内容 | UI人間操作E2E | 最終検証 |
| `SCN-MSCC-008` | model list 失敗、空一覧、更新中を job 作成拒否へ接続する | 利用者 | model list が未更新、更新中、取得失敗、または空である | model 選択または job 作成を試みる | 空の model list 成功は取得済み 0 件として表示し、model 未選択または更新不能状態では保存または job 作成を拒否する | 状態表示、拒否理由、保存要求または作成要求の未送信 | UI人間操作E2E | 実装後 |
| `SCN-MSCC-009` | model list 更新、保存、取得の結果を短い要約で確認する | 利用者 | model list 更新、保存、再取得が行われる | 結果表示を確認する | 結果分類と短い要約だけを表示し、raw payload と secret は出さない | provider、対象画面、結果分類、model 数、短い日本語要約 | APIテスト | 実装後 |
| `SCN-MSCC-010` | APIキー未設定時に状態表示だけで更新不可を示す | 利用者 | APIキーが必要な provider で APIキーが未設定である | カード上の状態と操作を確認する | 共有カード内には AIサービス設定を開く導線を出さず、APIキー未設定で更新不可であることだけを表示する | 更新不可表示、共有カード内導線なし、表示領域の折り返し | UI人間操作E2E | 実装後 |

## Acceptance Checks

- `SCN-MSCC-001`: 再読込後も参照側の provider、model、credential 参照状態が表示される。
- `SCN-MSCC-002`: provider 変更後に旧 provider の model を保存できない。
- `SCN-MSCC-003`: fake mode でも利用者向け provider list に `fake` が出ない。
- `SCN-MSCC-004`: APIキー未設定の Gemini または xAI は外部 model list 取得を開始しない。
- `SCN-MSCC-005`: 遅延応答後も現在 provider の model list と model 選択が維持される。
- `SCN-MSCC-006`: 保存失敗時は未保存変更として残し、保存済み扱いにしない。
- `SCN-MSCC-007`: 3 翻訳段階の不足が 1 件でもある場合、Ready job は作成できない。
- `SCN-MSCC-008`: 空の model list 成功は取得済み 0 件として表示し、raw payload と secret は表示されない。
- `SCN-MSCC-009`: 保存要約と取得要約は分類、対象、短い説明だけを含む。
- `SCN-MSCC-010`: APIキー未設定時の共有カード内表示は、内部状態名や AIサービス設定導線ではなく、更新不可の状態を示す。

## Fake / Stub Policy

- fake transport DI、fake secret store、Wails gateway stub、provider adapter stub を使う。
- 有料の実 AI API は呼ばない。
- fake mode は backend または adapter 境界で扱い、frontend の条件分岐にしない。
- 遅延応答は応答順序を制御できる stub で検証する。

## Risks

- 参照側ごとの保存 namespace を誤ると、マスターペルソナと Job Setup の状態が混入する。
- 空の model list 成功を取得失敗と混ぜると、利用者に通信失敗と候補 0 件の違いが伝わらない。
- 保存失敗後の未保存変更表示が弱いと、利用者が保存済みか未保存か判断できない。
- APIキー未設定時の既存バナー導線と共有カードの状態表示が離れすぎると、利用者が設定場所を探す可能性がある。

## Validation Commands

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/2026-05-07-model-settings-card-controller/scenario-design.md --json`
- `python3 scripts/harness/run.py --suite scenario-gate`

## Human Review State

- `design_review`: 人間回答反映済み
- `human_decision_required`: no
- `questionnaire`: `./scenario-design.questions.md`
- `implementation_scope`: `./implementation-scope.md`
