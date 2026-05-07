# 実装引き継ぎ入力: integration-model-settings-wails-gateway

## 状態

- `handoff_id`: `integration-model-settings-wails-gateway`
- `implementation_artifact`: `統合境界実装`
- `implementation_skill`: `implement-integration`
- `ready_wave`: `wave-3`
- `depends_on`: `frontend-shared-model-card-controller`, `backend-reference-model-settings-core`
- `source_scope`: `./implementation-scope.md`
- `source_scenario`: `./scenario-design.md`
- `source_ui_design`: `./ui-design.md`
- `source_frontend_result`: `./frontend-implementation-result.md`
- `source_backend_result`: `./backend-implementation-result.md`
- `frontend_human_review`: 承認済み

## 目的

backend core と frontend shared controller を Wails、DTO、gateway、bootstrap 境界で接続する。
接続範囲は provider、model、credential 状態、model list 状態の意味合わせに限定する。

## 承認済み実装範囲

- Wails bind、backend controller DTO、frontend gateway DTO、frontend gateway contract は provider / model / credential state / model list status を同じ意味で扱う。
- generated binding は必要な生成結果だけを扱い、手編集しない。
- APIキー本体と raw payload は Wails DTO、frontend gateway DTO、UI、console へ出さない。
- fake mode で通常 provider ID のまま `fake-model` が取得結果として伝わる。
- 遅延した model list 応答は現在 provider と現在要求へ反映しない。
- 実装後に `agent-browser` で Job Setup とマスターペルソナの共有カード表示を確認する材料を揃える。

## 対象範囲

- `internal/controller/wails/*` のモデル設定カード用 bind または既存 controller の接続。
- `internal/bootstrap/*` の手動 DI 接続。
- `frontend/src/controller/wails/*` と `frontend/src/controller/wails/gateway-dto/*` の gateway 接続。
- `frontend/wailsjs/` の生成物が必要な場合は生成結果だけを扱う。
- frontend screen controller factory の production gateway wiring。

## 依存完了情報

- `frontend-implementation-result.md`: frontend shared model card controller は完了した。
- `frontend-human-review-input.md`: 人間が「フロントレビュー終わり」と回答したため承認済みである。
- `backend-implementation-result.md`: backend core は既存実装で承認済み範囲を満たすため、backend プロダクトコード変更は不要と判定済みである。

## 範囲外

- frontend 表示だけの修正は扱わない。
- backend core だけの修正は扱わない。
- docs 正本本文、`.codex`、プロダクト外の workflow は変更しない。
- プロダクトテスト、検証データ、snapshot、test helper は変更しない。
- APIキー本文、復号可能値、provider authorization、raw request、raw response、raw prompt、内部 request 識別子を Wails DTO、frontend gateway DTO、UI、console、structured log、error summary、fake transport log へ出さない。

## 初手

`internal/controller/wails/` の公開 bind と `frontend/src/controller/wails/` の gateway DTO の field 対応を 1 つ固定する。
対応する完了条件は、Wails DTO と frontend gateway DTO が同じ redaction 境界を持つことである。

理由: backend と frontend の接続不一致を早く検出するためである。

## 完了条件

- Wails bind、backend controller DTO、frontend gateway DTO、frontend gateway contract が provider / model / credential state / model list status を同じ意味で扱う。
- generated binding を手編集していない。
- APIキー本体と raw payload を Wails DTO、frontend gateway DTO、UI、console へ出していない。
- fake mode で通常 provider ID のまま `fake-model` が取得結果として伝わる。
- 遅延した model list 応答は現在 provider と現在要求へ反映されない。
- agent-browser で Job Setup とマスターペルソナの共有カード表示を確認した結果を残す。

## 検証コマンド

- `go test ./internal/controller/wails ./internal/bootstrap -run 'ProviderSettings|Model|MasterPersona|TranslationJobSetup'`
- `npm --prefix frontend run check`
- `npm --prefix frontend run test -- --run 'gateway|master-persona|translation-job-setup|model'`
- 変更が frontend と backend の両方へ及ぶ場合は `python3 scripts/harness/run.py --suite backend-local` と `python3 scripts/harness/run.py --suite frontend-local`

## UI 確認

- 起動済みでなければ `npm run dev:wails:agent-browser` を使う。
- `agent-browser open http://localhost:34115/#master-persona` でマスターペルソナの共有カードを確認する。
- `agent-browser open http://localhost:34115/#translation-management` から `セットアップ` tab で Job Setup の共有カードを確認する。
- fake mode の `fake-model` は backend または adapter 境界から伝わる結果として確認する。
- console error と Wails 呼び出し失敗が残っていないことを確認する。

## 期待する返却

- 実装結果を `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-settings-card-controller/integration-implementation-result.md` に作成する。
- 変更ファイル、接続内容、secret 非露出確認、検証結果、実画面確認結果、残留リスクを分けて記録する。
- 承認済み統合境界範囲外の失敗は停止理由として記録する。
