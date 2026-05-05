# Scenario Design: frontend-fake-api-review-foundation

- `skill`: scenario-design
- `status`: approved
- `human_review`: approved
- `source_plan`: `./plan.md`
- `ui_source`: `N/A`
- `final_artifact_path`: `docs/scenario-tests/frontend-fake-api-review-foundation.md`
- `topic_abbrev`: `FFARF`
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## Fixed Requirements

- `must_pass_requirements`:
  - fakeAPI 起動モードは、フロントエンドレビュー用の DI 差し替えである。
  - fakeAPI 起動モードを、利用者向けのプロバイダー選択肢にしない。
  - fakeAPI 起動では、Wails バインディング、バックエンド、永続化に依存せずに実画面を開ける。
  - 本番起動では、fakeAPI とレビューモックデータを選ばない。
  - 空状態、読み込み中、成功状態、進行中状態、失敗状態、設定不足状態をレビュー状態パターンとして再現できる。
  - 状態パターンの既定値はレビュー起動条件で指定し、実画面確認では URL パラメータで上書きできる。
  - URL パラメータによる状態パターン指定は fakeAPI 起動中だけ有効である。
  - 画面固有のモックデータは、後続のユースケース task 側で追加できる。
  - `agent-browser` で、状態パターンごとの実画面を確認できる。
  - coverage harness の例外は、局所テスト結果と例外理由の記録がある場合だけ扱える。
- `non_goals`:
  - バックエンドの本番挙動変更は含めない。
  - 生成済み `wailsjs` の手編集は含めない。
  - fakeAPI、レビューモックデータ、状態パターンを本番初期状態へ入れない。
  - 画面固有の詳細モックデータ本体は後続のユースケース task 側で扱う。
  - レビュー専用 UI、状態パターン選択 UI、表示文言設計は含めない。
  - プロダクトコード、プロダクトテスト、docs 正本、`.codex`、implementation-scope は扱わない。

## Candidate Coverage

正本: `./scenario-design.candidate-coverage.json`

6 件の候補成果物は揃っている。
候補網羅 JSON では、生成 agent 名を付けた `generator:CAND-...` を一意な識別子として扱う。

`needs_human_decision` は 0 件である。
未解決の競合は 0 件である。
全候補は、採用、統合、不採用のいずれかへ分類済みである。

## Requirement Coverage

正本: `./scenario-design.requirement-coverage.json`

各抽象要件の詳細要求タイプは別 JSON に分離する。
この成果物は人間レビュー前の draft である。
人間レビュー承認後に implementation-scope を作る。

### `REQ-FFARF-001` レビュー起動で fakeAPI へ差し替える

- `source_requirement`: 起動モードで、フロントエンドの API 接続先を fakeAPI に切り替える。
- `requirement_kind`: external_integration
- `needs_human_decision`: なし
- `fixed_decisions`: fakeAPI は利用者向けのプロバイダー選択肢ではない。フロントエンドの composition root で、レビュー起動時だけゲートウェイを差し替える。Wails バインディングとバックエンドがなくても、実画面確認へ進める。

### `REQ-FFARF-002` レビュー状態パターンを実画面で確認する

- `source_requirement`: 空状態、読み込み中、成功状態、進行中状態、失敗状態、設定不足状態を fakeAPI で再現する。
- `requirement_kind`: display
- `needs_human_decision`: なし
- `fixed_decisions`: 状態パターンはレビュー起動時だけ有効である。既定値はレビュー起動条件で指定する。実画面確認では URL パラメータで状態パターンを上書きできる。画面内の選択 UI は作らない。

### `REQ-FFARF-003` 画面固有モックデータを後続 task で追加する

- `source_requirement`: 画面固有のモックデータを、ユースケース task 側で追加できる。
- `requirement_kind`: operation
- `needs_human_decision`: なし
- `fixed_decisions`: 共通基盤は、状態パターンとゲートウェイ契約の差し込み境界だけを提供する。個別画面のモックデータ本体と表示期待値は、後続のユースケース task 側で固定する。

### `REQ-FFARF-004` 本番起動と永続化へ混入させない

- `source_requirement`: 本番起動では fakeAPI が選ばれず、モックデータが本番 API、永続化、本番初期状態に混入しない。
- `requirement_kind`: security
- `needs_human_decision`: なし
- `fixed_decisions`: 本番起動では本番ゲートウェイだけを使う。URL パラメータなどのレビュー状態パターン指定が残っていても、fakeAPI へ遷移しない。モックデータは永続化へ書き込まない。

### `REQ-FFARF-005` レビュー入力不備を安全に拒否する

- `source_requirement`: 未登録状態パターン、モックデータ欠落、fakeAPI 未選択、Wails バインディング混入を検証で検出できる。
- `requirement_kind`: workflow
- `needs_human_decision`: なし
- `fixed_decisions`: 未登録状態パターンは本番初期状態へフォールバックしない。モックデータ欠落は成功状態に見せない。バインディング混入は、利用者向けメッセージではなく検証結果または内部診断として扱う。

### `REQ-FFARF-006` 局所テストと coverage 例外理由を記録する

- `source_requirement`: fakeAPI 起動モードが壊れていないことを局所テストで確認し、coverage harness 例外理由を記録する。
- `requirement_kind`: non_functional
- `needs_human_decision`: なし
- `fixed_decisions`: coverage 数値判定の代替にできるのは、fakeAPI 起動、DI 差し替え、本番非選択、状態パターン供給を確認する局所テスト結果と例外理由が対応づく場合だけである。

## Human Decision Questionnaire

正本: `N/A`

未回答質問はない。
実装手段の細部は、人間レビュー後の implementation-scope で分割して固定する。

## Risks

- レビュー起動モードの指定方法をプロバイダー設定に混ぜると、fakeAPI が利用者向け選択肢として見える。
- URL パラメータとレビュー起動条件の優先順が曖昧な場合、実画面レビューの再現性が落ちる。
- モックデータを共有初期状態に置くと、本番起動と永続化へ混入する。
- 進行中画面が runtime event に強く依存している場合、fake ゲートウェイ応答だけでは再現できない可能性がある。
- 画面固有モックデータは本 task では揃わないため、後続のユースケース task がテストデータを追加する必要がある。
- coverage 例外を広く扱うと、本番経路の検証不足を隠す可能性がある。

## Rules

- ケース ID は `SCN-FFARF-NNN` 形式にする。
- Markdown table は使わず、1 ケースごとの縦型ブロックで書く。
- 受け入れテストは全ケースで先に固定する。
- `実行テスト種別` は `APIテスト | UI人間操作E2E | lower-level only` に固定する。
- `実行段階` は `実装後 | 最終検証` に固定する。
- `期待結果` は観測可能な結果にする。
- `needs_human_decision` が残る場合はシナリオ完了にしない。
- 未解決競合が残る場合はシナリオ完了にしない。
- 有料の実 AI API を前提にしない。

## Scenario Matrix

このシナリオ表は人間レビュー前の draft である。
人間レビュー承認後に implementation-scope を作る。

### SCN-FFARF-001 レビュー起動で fakeAPI へ差し替える

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: fakeAPI レビュー起動が、Wails バインディングとバックエンドなしで成立する。
- `受け入れ条件`: レビュー起動モードでフロントエンドを起動すると、GatewayContract に fake ゲートウェイが注入される。
- `事前条件`: 基盤データ管理が成立し、レビュー起動モードが指定されている。
- `public_seam_or_api_boundary`: フロントエンドの composition root とゲートウェイ契約。詳細 API 名は implementation-scope で固定する。
- `contract_freeze`: あり。fakeAPI はプロバイダー選択肢ではなく、レビュー起動用の DI 差し替えである。
- `入力開始点`: fakeAPI レビュー起動 command。
- `主要 outcome`: 実画面がバックエンド未起動でも開く。
- `開始操作`: レビュー起動モードでアプリを起動する。
- `入力方法`: レビュー起動条件で fakeAPI を指定する。
- `主要操作列`: 起動条件で fakeAPI を指定し、`agent-browser` でレビューURLを開き、fakeAPI 由来の状態を確認する。
- `期待結果`:
  1. fake ゲートウェイが選ばれる。
  2. 生成済み `wailsjs` とバックエンド Controller を呼ばずに画面確認へ進める。
  3. View、ScreenController、Frontend UseCase は生成済み `wailsjs` を直接参照しない。
- `観測点`: `agent-browser` 表示、ゲートウェイ選択の局所テスト、Wails adapter 境界の検査結果。
- `UI-visible outcome`: fakeAPI 由来の対象画面状態が表示される。
- `fake_or_stub`: フロントエンド DI で差し替える fake ゲートウェイ。
- `責務境界メモ`: 本番ゲートウェイは `frontend/src/controller/wails/` に閉じる。

### SCN-FFARF-002 状態パターンを実画面で確認する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: レビュー実行者が 6 種の状態差分を同じ実画面上で確認できる。
- `受け入れ条件`: 空状態、読み込み中、成功状態、進行中状態、失敗状態、設定不足状態が状態パターンとして選べる。
- `事前条件`: fakeAPI レビュー起動が成立し、状態パターン一覧が登録済みである。
- `public_seam_or_api_boundary`: レビュー状態パターン指定境界。詳細指定方法は implementation-scope で固定する。
- `contract_freeze`: あり。状態パターンはレビュー起動時だけ有効である。
- `入力開始点`: `agent-browser` で開いた実画面。
- `主要 outcome`: 指定した状態パターンと実画面表示が一致する。
- `開始操作`: 状態パターン指定済みのレビューURLを開く。
- `入力方法`: レビュー起動条件で既定の状態パターンを指定し、必要に応じて URL パラメータで上書きして画面を開く。
- `主要操作列`: 状態パターンごとに URL パラメータを変えて再読み込みし、表示、主要操作可否、メッセージを確認する。
- `期待結果`:
  1. 空状態はデータなしと次操作を示す。
  2. 読み込み中と進行中は区別して表示される。
  3. 成功状態はモックデータに対応する表示モデルを表示する。
  4. 失敗状態は利用者向けメッセージだけを表示し、内部診断を出さない。
  5. 設定不足状態は不足項目と次操作を失敗状態と分けて表示する。
- `観測点`: 状態パターン一覧、Store の画面状態、Presenter の表示モデル、`agent-browser` の snapshot または screenshot。
- `UI-visible outcome`: 状態名、主要表示、操作可否、次操作が状態パターンごとに分かる。
- `fake_or_stub`: レビュー状態モックデータ、fake ゲートウェイ応答、必要時の fake runtime event adapter。
- `責務境界メモ`: 進行中状態は、Wails runtime event の本番購読なしで再現できることを優先する。

### SCN-FFARF-003 画面固有モックデータを後続 task で追加できる

- `分類`: 拡張系
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: fakeAPI 基盤が画面固有のレビューデータ追加を受け入れる。
- `受け入れ条件`: 後続のユースケース task が、対象画面のゲートウェイ契約に適合するモックデータと状態パターンを追加できる。
- `事前条件`: 共通 fakeAPI 基盤と状態パターン登録境界がある。
- `public_seam_or_api_boundary`: fakeAPI テストデータ登録境界。詳細名は implementation-scope で固定する。
- `contract_freeze`: あり。画面固有モックデータは、本番ゲートウェイ、生成済み `wailsjs`、バックエンド DTO に依存しない。
- `入力開始点`: 後続のユースケース task のテストデータ追加。
- `主要 outcome`: 追加した状態パターンは、レビュー起動時だけ選べる。
- `開始操作`: 画面固有モックデータを登録する。
- `入力方法`: 対象画面、状態パターン id、モックデータ要約を fakeAPI 基盤へ渡す。
- `主要操作列`: モックデータを登録し、局所テストで fakeAPI 起動時だけ解決されることを確認する。
- `期待結果`:
  1. 登録済み状態パターンが対象画面の成功または失敗状態へ反映される。
  2. 別画面と本番初期状態へモックデータが広がらない。
  3. モックデータの本文全量や secret に見える値を証跡へ保存しない。
- `観測点`: 状態パターン登録、テストデータ解決結果、本番構成からの非参照確認。
- `UI-visible outcome`: 後続 task の実画面で対象画面固有の状態パターンを確認できる。
- `fake_or_stub`: 画面別 fake テストデータ。
- `責務境界メモ`: モックデータの配置単位と命名規則は、人間レビュー後の implementation-scope で固定する。

### SCN-FFARF-004 本番起動で fakeAPI とモックデータを選ばない

- `分類`: 禁止遷移
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `最終検証`
- `観点`: レビュー用の fakeAPI とモックデータが、本番起動、永続化、本番初期状態へ混入しない。
- `受け入れ条件`: 本番起動相当では本番ゲートウェイが選ばれ、fake ゲートウェイとレビューモックデータは本番構成に入らない。
- `事前条件`: レビュー起動モードと本番起動相当を区別できる。
- `public_seam_or_api_boundary`: 本番構成境界。詳細確認方法は implementation-scope で固定する。
- `contract_freeze`: あり。本番起動から fakeAPI ゲートウェイ注入状態へ遷移しない。
- `入力開始点`: 本番起動相当の composition root 評価。
- `主要 outcome`: fakeAPI 非選択とモックデータ非混入を確認できる。
- `開始操作`: 本番起動相当でフロントエンドを起動、または構成を評価する。
- `入力方法`: レビュー起動モードを指定しない。
- `主要操作列`: 本番ゲートウェイ選択、fakeAPI 非参照、永続化非接続、初期状態非注入を確認する。
- `期待結果`:
  1. 本番ゲートウェイだけが選ばれる。
  2. fakeAPI 状態パターン指定の URL パラメータが残っても、本番表示へモックデータが混入しない。
  3. SQLite、keyring、本番初期状態へモックデータを書き込まない。
  4. 同一状態パターンの再実行でモックデータが重複挿入されない。
- `観測点`: ゲートウェイ選択結果、import 境界確認、永続化 write 有無、局所テスト結果。
- `UI-visible outcome`: 本番起動では、レビューモックデータ由来の状態を表示しない。
- `fake_or_stub`: なし。fakeAPI が選ばれないことを確認する。
- `責務境界メモ`: 本番安全性は通過条件であり、他の UI 成功で相殺しない。

### SCN-FFARF-005 レビュー入力不備を安全に拒否する

- `分類`: 主要失敗系
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: 状態パターンやモックデータの入力不備が、成功表示や本番フォールバックへ化けない。
- `受け入れ条件`: 未登録状態パターン、欠落モックデータ、fakeAPI 未選択、バインディング混入を成功状態として扱わない。
- `事前条件`: fakeAPI レビュー起動または局所テストで不備を再現できる。
- `public_seam_or_api_boundary`: レビュー状態パターン検証境界。詳細 API 名またはテストデータ名は implementation-scope で固定する。
- `contract_freeze`: あり。不備時は本番初期状態へフォールバックしない。
- `入力開始点`: レビュー状態パターン指定または fakeAPI 起動判定。
- `主要 outcome`: レビュー不備が明示的に拒否または検証失敗として観測できる。
- `開始操作`: 未登録状態パターンまたは欠落モックデータを指定する。
- `入力方法`: レビュー起動条件または局所テストのテストデータで、不正な状態パターンを指定する。
- `主要操作列`: 不備を指定し、局所テスト結果を確認する。
- `期待結果`:
  1. 未登録状態パターンは成功状態として表示されない。
  2. 欠落モックデータは設定不足状態そのものと区別して扱う。
  3. fakeAPI 未選択時は、レビュー開始不能または検証失敗として分かる。
  4. Wails バインディング混入は、利用者向けメッセージではなく検証結果または内部診断に分ける。
- `観測点`: 状態パターン検証、テストデータ解決結果、責務境界検査。
- `UI-visible outcome`: なし。入力不備の利用者向け表示は画面固有 task で扱う。
- `fake_or_stub`: 不正状態パターン用テストデータ、欠落モックデータ用テストデータ。
- `責務境界メモ`: 設定不足状態パターンは業務表示であり、モックデータ欠落によるレビュー不能とは分ける。

### SCN-FFARF-006 局所テストと coverage 例外理由を完了根拠にする

- `分類`: 検証根拠
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `最終検証`
- `観点`: fakeAPI 基盤の検証を広い coverage 数値だけに依存しない。
- `受け入れ条件`: coverage harness 例外を扱う場合、例外対象、理由、代替局所テスト結果、再確認条件が記録される。
- `事前条件`: fakeAPI 起動モードの局所テストが存在する。
- `public_seam_or_api_boundary`: 検証証跡境界。記録先は implementation-scope と run 証跡で固定する。
- `contract_freeze`: あり。coverage 例外は fakeAPI レビュー基盤に限定し、本番経路の検証不足を隠さない。
- `入力開始点`: fakeAPI 起動モードの局所テスト。
- `主要 outcome`: DI 差し替え、本番非選択、状態パターン供給、coverage 例外理由が対応づく。
- `開始操作`: 局所テストを実行する。
- `入力方法`: fakeAPI 起動、本番起動、状態パターン登録のテストデータを使う。
- `主要操作列`: 局所テストを実行し、coverage 例外理由と結果を作業成果物へ対応づける。
- `期待結果`:
  1. fakeAPI 起動モードが壊れていないことを局所テストで確認できる。
  2. 本番非選択の局所テスト結果がある。
  3. coverage 例外理由と代替テスト結果が同じ検証証跡から追える。
  4. `agent-browser` 証跡は状態パターン id と確認結果に対応づく。
- `観測点`: 局所テスト結果、coverage 例外理由、`agent-browser` 証跡参照。
- `UI-visible outcome`: なし。実画面表示そのものは SCN-FFARF-002 で扱う。
- `fake_or_stub`: fake ゲートウェイ、状態パターンテストデータ、本番構成テストデータ。
- `責務境界メモ`: operation-audit 候補の保存内容は product audit log ではなく、作業成果物と検証証跡に限定する。
