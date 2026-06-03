# refactor-action-enablement-derive-on-frontend

## 依頼要約

backend が返す UX 遷移可否 flag（`canStart` / `canPause` / `canResume` / `canRetry` / `canStartNextPhase` など）と `BlockedReason` 群を、ドメイン状態・種別・条件（phase 状態 enum、ジョブ lifecycle、件数、設定構成有無 など）に置き換え、UX 遷移判断を frontend 側（presenter / state 選択子）に集約する横断リファクタ。

対象は term / persona / body の 3 フェーズ controller、対応する usecase の result 型、frontend gateway contract、runtime shape validator、presenter、関連テスト。

motivating bug: 単語翻訳フェーズで「開始」押下時に「単語翻訳段階の操作に失敗しました。」が出る不具合（frontend validator が `canStartNextPhase` を必須検証する一方、Go 側 DTO に存在しない）。本リファクタの結果として自然解消する見込み。

## branch 情報

- 作業 branch: `claude/refactor-action-enablement-derive-on-frontend`
- 分岐元 branch: `master`
- 分岐元 commit: `325f054a698083d68c3ceaa3a96d9e1a587eb71f`

## 参照

- 原則: `docs/coding-guidelines-backend.md`（`canProceed` / `nextEnabled` / `isReady` / `isStartable` 風 flag を禁止する記述）

## 想定 Y/N 評価（design-module 入口）

- 仕様変更または仕様追加がある: N。外部に見えるユーザー操作の結果と状態文言は変えず、UX 遷移可否の判定主体だけ frontend に移す。参照: `docs/coding-guidelines-backend.md` 第 7 節 禁止事項。
- 画面変更がある: N。画面構造、文言、layout、表示順は変えない。参照: 対象 presenter は内部の状態導出だけを差し替える。
- 内部構造変更がある: Y。backend controller DTO（`internal/controller/wails/{term,persona,body}_translation_phase_controller.go`）と usecase contract（`internal/usecase/{term,persona,body}_translation_phase_contract.go` / `persona_generation_phase_contract.go`）、frontend gateway contract（`frontend/src/application/gateway-contract/*/`）、runtime shape validator、presenter の責務再配置がある。
- 画面の表示変更がある: N。layout、文言、style、表示構造は変えない。
- frontend ロジック変更がある: Y。gateway contract 型、runtime shape validator、presenter（state 選択子）でアクション遷移可否の導出を新設する。
- backend 変更がある: Y。3 phase の controller DTO と usecase result 型から `canStart` / `canPause` / `canResume` / `canRetry` / `canStartNextPhase` / `BlockedReason` を除去し、判断根拠となる状態・種別・条件だけ残す。
- frontend と backend を接続する: Y。contract 差分が gateway を通って frontend に届くため、Wails bridge の型と shape validator の同時改修が要る。
- 実装済み責務を独立に証明したい: Y。frontend 側で新設する遷移可否導出関数（presenter / 選択子）を単体テストで独立に証明する。
- 実行時にしか確定しない値または原因分離が要る分岐がある: N。導出は phase 状態 enum、ジョブ lifecycle、件数、設定構成有無から決定論的に求まる。

### decision table 結果

- 詳細仕様差分: 不要（仕様変更 N）。省略理由: 外部仕様は変えない、内部契約変更だけ。
- 画面設計差分: 不要（画面変更 N）。省略理由: 画面の表示構造、文言、layout を変えない。
- 設計差分図: 要（内部構造変更 Y）。`design-diff.md` に固定済み。
- 人間設計レビュー: 承認済み（第 3 版 + H 節 + 下記論点確定をもって承認）。
- 実装範囲: 要。
- テスト設計: 要。

### finalization-module 記録

- 正本化判断: 不要（仕様変更 N、外部仕様を変えない内部契約リファクタ）。
- 詳細仕様正本反映: 不要（人間承認済み恒久仕様なし）。
- 作業 commit: `ca5f8480`（77 files changed, +4975 / -1197）。

### 最終検証結果

- `python3 scripts/harness/run.py --suite backend-local`: 通過（apitest / bootstrap / controller wails / service / usecase 等 13 パッケージ）。
- `python3 scripts/harness/run.py --suite frontend-local`: 通過（54 Test Files / 613 Tests）。
- 追加差し戻し:
  - BE-body `previousPhaseLifecycle` を空文字で返していた問題を解消（`bodyTranslationPersonaReadinessPort.ReadBodyReadiness` 経路で persona phase の lifecycle を取得）。
  - `gofmt` 整形漏れ修正。
  - apitest の `TestSCN_TJSR_002` を本 task の責務 4 分割方針に合わせて改題・projection 検証へ書き換え（旧名: `BodyPhaseActionEnablementUsesCommonStateRules` → 新名: `BodyPhaseProjectionExposesLifecycleForFrontendDerivation`）。
  - revive lint 4 件解消（receiver 名統一、unused parameter）。
- 単体テスト 29 件追加（RAEF-UNIT-001〜029）、シナリオテスト 16 件追加（RAEF-E2E-001〜016）。
- SonarCloud coverage 66.9% は本 task 範囲外の既存課題として残置（プロダクトコード未カバーは別 task）。

### 人間設計レビュー確定事項

- 責務 4 分割: ドメイン情報集合 / ドメイン状態射影 / summary（表示用集約） / UX 遷移可否。frontend presenter の判定入力は「ドメイン状態射影」のみ、summary を判定入力に使わない。
- backend response: 1 method の中で **projection field group** と **summary field group** を並列に置く（別 method にしない。往復削減）。
- persona → body 移行の readiness 判定:
  - 「ペルソナ生成 - 次段階へ」ボタンの有効化 ⇔ ¬terminal ∧ persona phaseLifecycle ∈ COMPLETED_PHASE。
  - ペルソナ成果物が本文翻訳から参照可能か（readiness P2）は body 側「本文翻訳 - 開始」ボタンだけで評価する。persona 側 projection に readiness を載せない。
  - design-diff.md H-11 と G-2-b、body projection の `personaBodyReadiness` がこの方針で書かれている。
