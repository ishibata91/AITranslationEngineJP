# 実装引き継ぎ入力: frontend-shared-model-card-controller

## 状態

- `handoff_id`: `frontend-shared-model-card-controller`
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `ready_wave`: `wave-1`
- `depends_on`: なし
- `source_scope`: `./implementation-scope.md`
- `source_scenario`: `./scenario-design.md`
- `source_ui`: `./ui-design.md`
- `human_review_status`: 承認済み

## 目的

マスターペルソナと翻訳ジョブ設定が使うモデル設定カード制御を、frontend 側の共有 controller / usecase / store / presenter へ集約する。
`AIModelSelectionCard.svelte` は表示部品のまま維持する。

## 承認済み実装範囲

- 参照側ごとに provider、model、model list、保存状態を保持する。
- マスターペルソナと翻訳ジョブ設定は同じ状態規則を使う。
- provider 変更後は旧 provider の model list と model を現在 provider の保存済み状態へ混入しない。
- 空の model list 成功は取得済み 0 件として表示する。
- 保存失敗後は未保存変更として残し、再試行できる。
- APIキー未設定時は共有カード内に AIサービス設定導線を出さず、更新不可状態だけを表示する。

## 対象範囲

- `frontend/src/application/gateway-contract/*` のモデル設定カード用 contract、または既存 provider settings / Job Setup contract の参照側 model list 拡張。
- `frontend/src/application/store/*` の共有モデル設定状態。
- `frontend/src/application/usecase/*` の provider 変更、model list 更新、model 選択、保存状態制御。
- `frontend/src/application/presenter/*` の共有カード view model。
- `frontend/src/controller/master-persona/*` と `frontend/src/controller/translation-job-setup/*` の screen controller 接続。
- `frontend/src/ui/components/AIModelSelectionCard.svelte` と両画面の利用箇所。

## 範囲外

- backend、Wails 公開 method、永続化、DB schema は変更しない。
- generated binding を変更しない。
- docs 正本本文、`.codex`、プロダクト外の workflow は変更しない。
- `fake` provider ID を user-facing provider として追加しない。
- frontend に fake mode 判定や `fake-model` 固有分岐を追加しない。
- APIキー本文、復号可能値、provider raw request / response、raw prompt、内部 request 識別子を UI、DTO、console、log へ出さない。

## 初手

`frontend/src/application/store/` 配下に、参照側ごとの provider / model / model list / 保存状態を保持する最小 state 型を追加する。
対応する完了条件は、参照側ごとの provider / model / model list / 保存状態を 1 つの状態規則で保持することである。

理由: 両画面の controller と presenter が同じ状態規則に依存するためである。

## 完了条件

- マスターペルソナと Job Setup は同じ frontend 状態規則で provider、model、model list、保存状態を扱う。
- provider 変更後は旧 provider の model list と model を現在 provider へ混入しない。
- 空の model list 成功は取得済み 0 件として表示される。
- 保存失敗後は未保存変更として残り、保存済み設定として表示または利用しない。
- APIキー未設定時は共有カード内に AIサービス設定導線を出さず、更新不可状態だけを表示する。
- fake mode 判定、`fake` provider ID、`fake-model` 固有分岐を frontend に追加していない。

## 検証コマンド

- `npm --prefix frontend run check`
- `npm --prefix frontend run test -- --run 'master-persona|translation-job-setup|AIModelSelectionCard|model'`
- `python3 scripts/harness/run.py --suite frontend-local`

## 期待する返却

- 実装結果を `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-settings-card-controller/frontend-implementation-result.md` に作成する。
- 変更ファイル、実装内容、検証結果、残留リスクを分けて記録する。
- 実画面確認ができない場合は、未確認理由と次に確認する URL / command を記録する。
