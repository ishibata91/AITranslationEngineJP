# Scenario Candidates: frontend-fake-api-review-foundation / external-integration

- `generator`: `external-integration`
- `source_plan`: `./task-frame.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `FFAR-EXT`

## Generator Scope

- `viewpoint`: `external-integration`
- `included_sources`: `task-frame.md`, `tasks/usecases/frontend-fake-api-review-foundation.yaml`, `docs/architecture.md`, `docs/coding-guidelines-frontend.md`, `tmp/code-map/index.json`
- `excluded_sources`: プロダクトコード, プロダクトテスト, docs 正本更新, `.codex` 変更
- `generation_notes`: fakeAPI は プロバイダー選択肢 ではなく、フロントエンド composition root の DI による差し替え候補として扱う。採否、統合、最終シナリオ表は `designer` に残す。

## Candidate Scenarios

### CAND-FFAR-EXT-001 fakeAPI 起動では Wails バインディングを呼ばない

- `根拠要件`: task-frame 完了条件「起動モードで フロントエンドの API 接続先を fakeAPI に切り替えられる」、設計前提「本番ゲートウェイは Wails バインディング adapter として `フロントエンド/src/controller/wails/` に閉じ込める」
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-FFAR-EXT-001`
- `external boundary`: Wails バインディング adapter
- `actor`: フロントエンドレビューer
- `trigger`: fakeAPI 起動モードで フロントエンド を起動する
- `期待結果`: 画面は `GatewayContract` に接続された fake ゲートウェイ から状態を取得し、生成済み `wailsjs` と バックエンド Controller を呼ばない
- `fake_or_stub`: フロントエンド DI で差し替える fake ゲートウェイ
- `観測点`: agent-browser で画面が開き、Wails バインディング 未接続でも対象状態が表示される
- `関連詳細要求タイプ`: Wails バインディング 境界、フロントエンド DI 境界
- `採用判断材料`: Wails バインディング と バックエンド 非依存のレビュー起動を証明する候補として採用候補
- `競合注意`: lifecycle 観点では起動成功として扱われる可能性があるが、external-integration 観点では Wails 呼び出し遮断を主目的にする

### CAND-FFAR-EXT-002 本番起動では fakeAPI が選ばれない

- `根拠要件`: task-frame 完了条件「本番起動では fakeAPI が選ばれない」、usecase 完了条件「fakeAPI と モックデータが本番 API、永続化、本番初期状態に混入しない」
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-FFAR-EXT-002`
- `external boundary`: dev 起動モードと 本番 wiring の境界
- `actor`: フロントエンドレビューer
- `trigger`: 本番相当の起動条件で フロントエンド composition root を評価する
- `期待結果`: 本番ゲートウェイ だけが生成され、fake ゲートウェイ と モックデータ は 本番 graph に入らない
- `fake_or_stub`: なし。fakeAPI 選択が無効であることを観測する
- `観測点`: 局所テストまたは import 境界確認で 本番 wiring が fakeAPI を参照しない
- `関連詳細要求タイプ`: dev 起動モード境界、本番 混入防止
- `採用判断材料`: 本番混入防止の最低限候補として採用候補
- `競合注意`: state-transition 観点では起動モード遷移として扱われる可能性がある

### CAND-FFAR-EXT-003 バックエンド 未接続でも レビュー状態を再現できる

- `根拠要件`: task-frame 目的「Wails バインディング と バックエンド に依存せずに実画面を確認できる状態にする」、usecase goal「実フロントのままレビュー用状態を再現できるようにする」
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-FFAR-EXT-003`
- `external boundary`: バックエンド Controller と Wails transport 境界
- `actor`: フロントエンドレビューer
- `trigger`: バックエンド が起動していない、または Wails バインディング が利用できない状態で fakeAPI 起動モードを開く
- `期待結果`: 空状態、成功状態、失敗状態、設定不足状態などの レビュー状態が バックエンド 応答なしで表示される
- `fake_or_stub`: バックエンド 応答相当の DTO を返す フロントエンド fakeAPI
- `観測点`: agent-browser で 状態パターン 表示を確認し、バックエンド 起動有無に依存しないことを観測する
- `関連詳細要求タイプ`: バックエンド 非依存、Wails transport 境界
- `採用判断材料`: 実画面レビュー基盤の中核候補として採用候補
- `競合注意`: failure 観点の バックエンド 接続失敗と競合しうる。external-integration 観点では接続失敗をユーザー向けエラーにせず、レビュー起動成立として扱う

### CAND-FFAR-EXT-004 モックデータ は永続化へ書き込まれない

- `根拠要件`: task-frame 完了条件「fakeAPI と モックデータが本番 API、永続化、本番初期状態に混入しない」、architecture「Repository は SQLite などの永続化実装を持つ」
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-FFAR-EXT-004`
- `external boundary`: バックエンド repository と 本番初期状態 の境界
- `actor`: フロントエンドレビューer
- `trigger`: fakeAPI 起動モードで状態パターンを切り替え、画面固有モックデータを表示する
- `期待結果`: モックデータ は フロントエンド fakeAPI 内に閉じ、SQLite、keyring、本番初期状態 へ保存されない
- `fake_or_stub`: 永続化を持たない in-memory モックデータ 根拠
- `観測点`: 再読み込み後の状態、局所テスト、または依存境界確認で repository 参照が発生しない
- `関連詳細要求タイプ`: 永続化境界、モックデータ 混入防止
- `採用判断材料`: fakeAPI が本番データを汚さないことを証明する候補として採用候補
- `競合注意`: operation-audit 観点では保存証跡を期待する可能性があるが、fakeAPI の モックデータ は監査保存対象外にする前提が必要

### CAND-FFAR-EXT-005 runtime event なしで進行中状態を再現できる

- `根拠要件`: task-frame 完了条件「進行中状態を fakeAPI で再現できる」、architecture「Wails event は push 通知専用に限定し、通常の query / command を置き換えない」
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-FFAR-EXT-005`
- `external boundary`: Wails runtime event adapter 境界
- `actor`: フロントエンドレビューer
- `trigger`: fakeAPI 起動モードで進行中 状態パターンを選ぶ
- `期待結果`: Wails runtime event の購読なしに進行中表示が再現され、query / command の主経路は fake ゲートウェイ の 契約応答として扱われる
- `fake_or_stub`: fake runtime event adapter または progress 状態を返す fake ゲートウェイ
- `観測点`: agent-browser で進行中表示を確認し、runtime event adapter への本番依存が発生しないことを局所確認する
- `関連詳細要求タイプ`: runtime event 境界、進行中 状態パターン
- `採用判断材料`: runtime event に依存する画面の レビュー状態を固定する候補として採用候補
- `競合注意`: lifecycle 観点では進行中状態の時間経過を扱う可能性がある。external-integration 観点では runtime event 境界の置換可否だけを扱う

### CAND-FFAR-EXT-006 agent-browser で 状態パターンごとの外部境界を確認できる

- `根拠要件`: task-frame 完了条件「実画面を `agent-browser` で開き、状態パターンごとの表示を確認できる」、usecase manual_check_steps「実画面で状態パターンを切り替える」
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-FFAR-EXT-006`
- `external boundary`: agent-browser 確認入口と dev 起動モードの境界
- `actor`: フロントエンドレビューer
- `trigger`: fakeAPI 起動モードを agent-browser で開き、レビュー状態パターンを切り替える
- `期待結果`: 空状態、読み込み中、成功状態、進行中状態、失敗状態、設定不足状態を実画面で確認できる
- `fake_or_stub`: agent-browser から確認可能な fakeAPI 起動時 状態パターン指定
- `観測点`: agent-browser の画面表示証跡に、各 状態パターンの表示結果が残る
- `関連詳細要求タイプ`: agent-browser 確認境界、レビューstate 状態パターン
- `採用判断材料`: 人間UIレビュー前の外部確認入口として採用候補
- `競合注意`: 最終シナリオ表の確定に踏み込まない。表示内容の採否は画面固有 task 側に残す

### CAND-FFAR-EXT-007 画面固有モックデータを ユースケース task 側で追加できる

- `根拠要件`: task-frame 完了条件「画面固有のモックデータを ユースケース task 側で追加できる」、usecase outputs「レビューstate モックデータ」
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-FFAR-EXT-007`
- `external boundary`: shared fakeAPI 基盤と画面固有モックデータの差し込み境界
- `actor`: ユースケース task implementer
- `trigger`: 画面単位の レビュー状態を追加する task が モックデータを登録する
- `期待結果`: 画面固有モックデータ は shared fakeAPI 基盤へ差し込めるが、本番ゲートウェイ、生成済み `wailsjs`、バックエンド DTO へ依存しない
- `fake_or_stub`: 画面別 fake テストデータ と GatewayContract 実装
- `観測点`: 局所テストで画面別 テストデータ が fakeAPI 起動モードだけに接続されることを確認する
- `関連詳細要求タイプ`: adapter 境界、画面別 テストデータ 境界
- `採用判断材料`: 後続 ユースケース task が レビュー状態を増やすための拡張候補として採用候補
- `競合注意`: actor-goal 観点では画面レビュー体験として扱われる可能性がある。external-integration 観点では テストデータ 差し込み境界だけを扱う

### CAND-FFAR-EXT-008 coverage harness 例外は局所テスト根拠つきで記録される

- `根拠要件`: usecase 完了条件「coverage harness では fakeAPI 基盤を数値判定の例外として扱い、例外理由と局所テスト結果を記録できる」、task-frame 完了条件「fakeAPI 起動モードが壊れていないことを局所テストで確認できる」
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-FFAR-EXT-008`
- `external boundary`: coverage harness と fakeAPI dev-only 境界
- `actor`: implement_lane
- `trigger`: fakeAPI 基盤の局所テスト結果を確認し、coverage harness の数値判定例外を扱う
- `期待結果`: fakeAPI 基盤は本番挙動の coverage 対象と混同されず、例外理由と局所テスト結果が記録される
- `fake_or_stub`: dev-only fakeAPI 基盤。coverage 数値判定の代替として局所テスト証跡を使う
- `観測点`: 局所テスト結果と coverage 例外理由が作業成果物上で対応づく
- `関連詳細要求タイプ`: dev-only 外部境界、検証境界
- `採用判断材料`: fakeAPI 基盤の検証責務を過剰な本番 coverage と混同しない候補として採用候補
- `競合注意`: operation-audit 観点では記録形式を扱う可能性がある。external-integration 観点では coverage harness との境界だけを扱う

## Open Notes

- `人間判断候補`: agent-browser で 状態パターンを選ぶ UI を常設するか、起動時パラメータだけにするかは、最終シナリオ統合時に判断が必要である。
- `人間判断候補`: runtime event の進行中状態を fake runtime event adapter で再現するか、fake ゲートウェイ の応答状態で再現するかは、対象画面の既存 controller 境界を見て判断が必要である。
- `統合候補`: `CAND-FFAR-EXT-001` と `CAND-FFAR-EXT-003` は Wails/バックエンド 非依存の同一シナリオへ統合できる可能性がある。
- `統合候補`: `CAND-FFAR-EXT-002` と `CAND-FFAR-EXT-004` は本番混入防止の同一シナリオへ統合できる可能性がある。
- `不採用候補`: 有料 real API provider の選択、secret 保存、バックエンド provider 実装方針は対象外である。

## Completion Evidence

- `candidate_count`: 8
- `artifact`: `docs/exec-plans/completed/frontend-fake-api-review-foundation/scenario-candidates.external-integration.md`
- `viewpoint`: `external-integration`
- `source_task_folder`: `docs/exec-plans/completed/frontend-fake-api-review-foundation/`
