# UI Design: 2026-05-07-provider-settings-job-decoupling-implement

- `skill`: `ui-design`
- `status`: `pending-human-review`
- `source_plan`: `./plan.md`
- `scenario_source`: `./scenario-design.md`
- `ux_standard_source`: `docs/UX-standard.md`
- `target_screen`: Job Setup
- `related_screen`: AIサービス設定
- `implementation_scope`: 未作成

## 判断結果

UI 要件契約は作成済みである。
状態は人間レビュー待ちである。

固定点は、Job Setup が provider、model、execution mode、batch mode の選択値だけを扱うことである。
Job Setup は endpoint、secret store 参照実値、credential 参照実値を表示しない。

AIサービス設定画面は全面再設計しない。
AIサービス設定画面は provider ごとの endpoint、APIキー状態、接続確認状態を扱う既存画面として維持する。

## 根拠参照

| 根拠 | UI 判断 |
| --- | --- |
| `./scenario-design.md:25` | Job Setup の固定要件を UI 契約の固定点にする。 |
| `./scenario-design.md:36` | Job から外す値を Job Setup の禁止表示にする。 |
| `./scenario-design.md:118` | APIキー未設定、model list 取得失敗、model 未選択を別状態にする。 |
| `./scenario-design.md:148` | Running phase の保存要約では非 secret 分類だけを許可する。 |
| `docs/detail-specs/translation-job-setup.md:37` | Job Setup は 3 つの翻訳段階を扱う。 |
| `docs/detail-specs/translation-job-setup.md:56` | 既存 UI 契約由来の表示項目と操作を差分更新する。 |
| `docs/detail-specs/ai-provider-settings-management.md:26` | AIサービス設定は provider ごとの共通設定を扱う。 |
| `docs/detail-specs/ai-provider-settings-management.md:54` | AIサービス設定画面の UI 契約は維持する。 |
| `docs/architecture.md:75` | UI Component は画面専用部品と共有部品の二層で扱う。 |
| `docs/screen-design/README.md:17` | 新規 task の UI 判断は task-local `ui-design.md` に置く。 |

## UI Contract

対象画面:
Job Setup を対象にする。
作成後の設定内容も対象に含める。

関連画面:
AIサービス設定は参照元の設定画面として確認対象にする。
endpoint 編集、APIキー保存、接続確認、未設定化の再設計は対象外にする。

表示項目:

- 入力カード一覧、共通辞書、共通ペルソナを表示する。
- 単語翻訳、NPC ペルソナ生成、本文翻訳の各翻訳段階を表示する。
- 各翻訳段階は AIサービス、モデル、実行方法、一括処理、APIキー状態分類を表示する。
- 作成後の設定内容は AIサービス、モデル、実行方法、一括処理、APIキー状態分類だけを表示する。
- 作成前確認は不足理由を分類して表示する。

操作:

- input を選択する。
- job 未作成 input を削除する。
- 各翻訳段階で AIサービスを選ぶ。
- モデル一覧を更新し、モデルを選ぶ。
- 対象 provider だけ一括処理を切り替える。
- 作成前確認を満たした後に job を作成する。

状態:

- 読み込み中: Job Setup option を取得中である。
- 空: job 未作成 input がない。
- 未設定: APIキーが必要な AIサービスで APIキー状態分類が未設定である。
- 更新不可: APIキー未設定のため、モデル一覧を更新できない。
- 更新中: モデル一覧を取得している。
- 取得失敗: モデル一覧を取得できない。
- 作成不可: input、AIサービス、model、APIキー状態分類のいずれかが不足している。
- 作成済み: Ready job の要約を表示する。

エラー分類:

- APIキー未設定: APIキーが必要な provider のみ表示する。
- モデル未選択: モデル一覧取得後にモデルが未選択である。
- モデル一覧未更新: モデル一覧をまだ取得していない。
- モデル一覧取得失敗: 外部取得または fake transport 取得に失敗した。
- provider 参照不能: provider settings の再解決が失敗した。
- 作成失敗: job 作成 API が失敗した。

禁止表示:

- endpoint 原文を Job Setup に表示しない。
- `credential_ref`、secret store 参照実値、secret store の key 名を表示しない。
- APIキー本体、伏せ字以外の復号可能値、raw request、raw response、raw prompt を表示しない。
- structured log 用 ID、内部 token、`modelListSourceToken` を利用者向け表示に出さない。
- provider settings の更新履歴、revision、Job 側 fallback の説明を出さない。

## 操作契約

AIサービス選択:
利用者は各翻訳段階で provider を選ぶ。
provider 切り替え後は、モデル一覧の状態を未更新に戻す。

モデル一覧更新:
APIキーが必要な provider で APIキーが未設定なら、更新ボタンを押せない。
LM Studio は APIキー不要として扱い、APIキー未設定 warning を出さない。

モデル選択:
モデル一覧取得に成功し、候補が 1 件以上ある場合だけ model select を有効にする。
候補 0 件は取得成功と不足理由を分けて表示する。

一括処理:
一括処理は Gemini と xAI だけに表示する。
未対応 provider では一括処理の切り替えを表示しない。

job 作成:
3 つの翻訳段階で APIキー状態分類と model 選択が満たされた時だけ、job 作成ボタンを有効にする。
作成後の要約は選択値と APIキー状態分類だけを再掲する。

## 表示文言契約

固定名として残せる表示:

- `Gemini`
- `LM Studio`
- `xAI`
- `Batch API`
- provider の model 名

日本語にする表示:

- `provider` は `AIサービス` と表示する。
- `model` は `モデル` と表示する。
- `execution mode` は `実行方法` と表示する。
- `credential reference` は表示しない。
- `Validation status` は `作成前確認` と表示する。
- `blocking failure` は `作成できない理由` と表示する。
- `dirty state` は `再確認が必要` または `確認済み` と表示する。

状態文言:

- `configured` は `設定済み` と表示する。
- `missing` は `APIキー未設定` と表示する。
- `not_required` は `APIキー不要` と表示する。
- `not_updated` は `モデル一覧未更新` と表示する。
- `loading` は `更新中` と表示する。
- `failed` は `取得失敗` と表示する。
- `success` は `取得済み` と表示する。

## 構造

画面区画:

- 入力確認を最初に置く。
- 共通基盤を入力確認の後に置く。
- 翻訳段階別設定を 3 枚の同種カードで並べる。
- 作成前確認と job 作成操作を末尾に置く。
- 作成済み設定は job 作成後だけ表示する。

配置制約:

- 既存の page shell、card、grid、余白体系を維持する。
- AIサービス、モデル、状態、更新操作を同じ翻訳段階カード内に置く。
- 不足理由は該当カードの近くと作成前確認の両方で確認できるようにする。
- endpoint や credential 参照値を入れる行を Job Setup に作らない。

アクセシビリティ:

- エラーと警告は色だけで伝えず、日本語の状態文言を併記する。
- モデル一覧更新ボタンは翻訳段階名を含む `aria-label` を持つ。
- disabled の理由は作成前確認または該当カードの補足文で読めるようにする。
- キーボード順は入力カード、翻訳段階カード、作成前確認の順にする。

## UI 部品化判断

再利用する部品:

| 部品または契約 | 根拠 | 判断 |
| --- | --- | --- |
| `AIModelSelectionCard` | `frontend/src/ui/components/AIModelSelectionCard.svelte:21` | APIキー状態分類、モデル一覧更新、モデル選択、一括処理の既存部品として使う。 |
| `StickyActionFooter` | `frontend/src/ui/components/StickyActionFooter.svelte:34` | 作成前確認と job 作成操作の既存部品として使う。 |
| `ModelSettingsCardViewModel` | `frontend/src/application/gateway-contract/model-settings-card/model-settings-card-contract.ts:50` | 状態文言と有効条件の既存判断を確認対象にする。 |
| `buildModelSettingsCardViewModel` | `frontend/src/application/gateway-contract/model-settings-card/model-settings-card-policy.ts:365` | APIキー未設定時の更新不可と文言変換の確認対象にする。 |

画面専用で維持する部品:

| 対象 | 根拠 | 判断 |
| --- | --- | --- |
| `JobSetupPage` の入力カード一覧 | `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte:304` | input 選択と削除は Job Setup 専用の業務流れなので分けない。 |
| `JobSetupPage` の作成済み要約 | `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte:220` | Ready job 作成後の確認区画なので Job Setup 内に置く。 |
| `ProviderSettingsPage` | `frontend/src/ui/screens/provider-settings/ProviderSettingsPage.svelte:63` | 共通設定画面として維持し、Job Setup の UI 部品へ混ぜない。 |

分けない対象:

- Job Setup 全体の page shell は分けない。
- Provider settings の endpoint 入力部品は Job Setup へ流用しない。
- credential 参照選択 UI は Job Setup から外す表示契約にする。
- task-local prototype と mock は作らない。

確認対象:

- 現行 `JobSetupPage` には legacy fallback として `credential reference` select がある。
  対象行は `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte:466` である。
- 現行 gateway contract は `credentialRef` を request、selection、summary に含む。
  対象行は `frontend/src/application/gateway-contract/translation-job-setup/translation-job-setup-gateway-contract.ts:104`、`frontend/src/application/gateway-contract/translation-job-setup/translation-job-setup-gateway-contract.ts:167`、`frontend/src/application/gateway-contract/translation-job-setup/translation-job-setup-gateway-contract.ts:184` である。
- AIサービス設定 contract は endpoint と credentialReferenceId を持つ。
  対象行は `frontend/src/application/gateway-contract/provider-settings/provider-settings-gateway-contract.ts:35` である。

## UX Standard Review

画面構造の高優先度:

- 画面目的: Job Setup は Ready job 作成に必要な選択値を固定する画面である。
- 主要CTA: 主要操作は job 作成に絞る。
- 情報階層: input、共通基盤、翻訳段階、作成前確認の順を維持する。
- 状態表示: APIキー状態分類、モデル一覧状態、作成可否を分ける。
- 禁止条件: endpoint、secret store 参照実値、APIキー本体を出さない。

配置とレスポンシブの高優先度:

- 関連情報の近接: AIサービス、モデル、APIキー状態分類、モデル一覧更新を同じカード内に置く。
- セクション構造: 3 つの翻訳段階を同じ構造で比較できるようにする。
- 折り返し順: mobile では翻訳段階カードを縦積みにし、単語翻訳、NPC ペルソナ生成、本文翻訳の順を維持する。
- 長文耐性: 長い model 名、長い input 名、長いエラー文は折り返す。
- 色以外の表現: warning と success は状態文言も併記する。

延期:

- 幅別スクリーンショットは実装後確認で取得する。
- 状態別スクリーンショットは実装後確認で取得する。
- visual polish は人間レビュー後の実装確認で扱う。

## Desktop Review Points

- 1280px 以上で、入力カード一覧、翻訳段階カード、作成前確認が読める。
- 3 つの翻訳段階カードでラベル位置と操作位置が揃う。
- 作成前確認の sticky footer が翻訳段階カードの補足文を隠さない。
- endpoint と credential 参照値が Job Setup に表示されない。
- 長い model 名がカード幅からはみ出さない。

## Mobile Review Points

- 390px 幅で入力確認、共通基盤、翻訳段階別設定、作成前確認の順に読める。
- 翻訳段階カードは 1 カラムで縦に並ぶ。
- モデル一覧更新ボタン、select、checkbox のタップ領域が狭すぎない。
- sticky footer がキーボード表示時に入力欄を隠さない。
- エラー文と不足理由が横スクロールなしで読める。

## Agent Browser Review

- `command_source`: `agent-browser`
- `app_start_command`: 未実行
- `checked_url`: 未確認
- `checked_viewports`: 未確認
- `screenshot_or_snapshot_refs`: なし
- `console_errors`: 未確認
- `layout_breaks`: 未確認
- `ambiguous_interactions`: 未確認

未実施理由:
この作業は UI 設計成果物の作成であり、実装後画面はまだ存在しない。
入力根拠と既存 UI ファイルから表示契約を固定できるため、アプリ起動と実画面確認は実装後確認へ回す。

表示文言レビュー:

- `review_timing`: 実画面確認未実施のため、設計時レビューとして実施した。
- `fixed_names_preserved`: provider 名、model 名、`Batch API` だけ固定名として残す。
- `business_japanese_terms`: `AIサービス`、`モデル`、`実行方法`、`作成前確認` を使う。
- `internal_state_names_hidden`: `credential_ref`、`modelListSourceToken`、`dirty`、`blocking failure` は表示しない。
- `next_action_wording`: 不足理由は `APIキーを設定してください`、`モデル一覧を更新してください`、`モデルを選んでください` のように次操作で書く。
- `allowed_english_labels`: 利用者が設定画面で見る provider 名と model 名だけを許可する。
- `plain_language_next_action_judgement`: 設計上は通過である。実画面表示は実装後に確認する。

## 完了条件

- Job Setup の表示項目、操作、状態、エラー分類、禁止表示を固定した。
- AIサービス設定画面の全面再設計を対象外にした。
- desktop と mobile の確認観点を固定した。
- 既存 component 再利用または確認対象を実在ファイル根拠付きで固定した。
- 実画面確認の未実施理由を明示した。

## 未決事項

人間判断が必要な未決事項はない。
人間レビュー後、implementation-scope で backend、frontend、integration の実装範囲を分ける必要がある。

