# 作業計画

- workflow: work
- status: completed
- lane_owner: orchestrate
- scope: master-persona-gap-closure
- task_id: persona-management-gap-closure
- task_catalog_ref: /Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/persona-management.yaml
- parent_phase: implementation-lane

## 依頼要約

- 人間が指摘した未解決事項を解消するため、マスターペルソナ実装のギャップ解消計画を新規作成する。
- 対象は、インメモリ永続化、service 内の fake/mock 風実装、JSON 選択後に生成へ進めない導線の 3 点とする。
- 追加判断として、secret access は `github.com/99designs/keyring`、fake は request / SDK transport seam の DI 差し替え、JSON preview は生成許可と分離して扱う。
- 追加判断として、real provider は `gemini` / `lm_studio` / `xai` とし、provider 実装は interface で抽象化して共通 response を返す。
- docs 正本は更新せず、既存 active plan と scenario artifact の design bundle を修正する。

## 判断根拠

- 直前の実装では `internal/bootstrap/app_controller.go` が `NewInMemoryMasterPersonaRepository` と `NewInMemorySecretStore` を直接配線しており、再起動後に状態が残らない。
- `internal/repository/master_persona_repository.go` は `InMemoryMasterPersonaRepository` に list / detail / settings / run status / mutation を集約しており、SQLite 永続化 concrete が未実装である。
- `internal/service/master_persona_service.go` は provider 名の判定で fake path を許可しているが、生成本文は service 内の固定ロジックで組み立てており、request / SDK transport seam の DI fake になっていない。
- `frontend/src/application/presenter/master-persona/master-persona.presenter.ts` は `preview.status === "生成可能"` を満たすまで生成ボタンを有効化しないため、JSON 選択だけでは生成へ進めない。
- 人間から、`インメモリやめろ`、`fake は provider list に出さず request / SDK transport seam だけを DI で差し替える`、`JSON 選択後は AI 設定未完了でも集計表示だけは出すが生成ボタンは無効にする` という判断基準が追加で示された。
- 2026-04-16 に `python3 scripts/harness/run.py --suite structure` を再実行し、pass を確認した。

## タスクモード

- task_mode: `fix`
- goal: マスターペルソナ実装の暫定 concrete と UX 上の詰まりを除去し、人間が明示した判断基準に一致させる。
- constraints: orchestrate 自身では product 実装をしない。docs 正本は更新しない。human review 前に implementation-scope を確定しない。
- close_conditions: implementation-review と ui-check を通し、master persona でインメモリ concrete が除去され、request / SDK transport seam の DI fake と JSON 選択導線のギャップが解消されることを確認する。

## 事実

- 現在のマスターペルソナ画面は Wails 実機で動作するが、永続化はインメモリ concrete に依存している。
- 現在の fake path は provider interface 実装ではなく、service 内の条件分岐と固定文面生成で成立している。
- 現在の UI は preview 成功まで生成ボタンを押せず、`JSON を選んでもボタン押せない` という違和感が発生する。
- 人間から、インメモリ禁止は residual risk ではなく stop condition として扱うべき、という恒久ルールが追加された。
- 人間から、API key 保存時の明示認証はアプリ内 modal ではなく OS キーチェーンの認可ダイアログで出す、という判断が示された。

## 機能要件

- summary:
  - マスターペルソナの repository、settings persistence、run status persistence はインメモリ concrete をやめ、再起動後も残る concrete へ置き換える。
  - secret access は `github.com/99designs/keyring` backed concrete とし、macOS Keychain と Windows Credential Manager を対象にする。
  - fake path は provider list に表示せず、HTTP request / SDK transport seam だけを DI で差し替える。
  - real provider は `gemini` / `lm_studio` / `xai` の 3 種とし、各 concrete は provider interface 経由で共通 response を返す。
  - JSON 選択後は preview を自動実行し、AI 設定未完了でも集計表示だけは出す。生成ボタンは AI 設定完了と preview 成功の両方を満たす時だけ有効にする。
- in_scope:
  - `internal/bootstrap/` の master persona 配線から `NewInMemoryMasterPersonaRepository` と `NewInMemorySecretStore` を除去する。
  - master persona 向けの persistence concrete を実装し、entry、AI settings、run status の責務を分離する。
  - secret access は `github.com/99designs/keyring` で macOS Keychain / Windows Credential Manager に保存し、一度保存した API key は通常利用で再認証不要にする。
  - API key 保存ボタンクリック時だけ OS キーチェーンの認可ダイアログを明示的に出し、アプリ内の追加確認 modal は置かない。
  - test mode では provider list に fake provider を表示せず、real provider と同じ prompt 組み立て、provider validation、run orchestration を通した後、request / SDK transport seam を DI fake に差し替える。
  - `JSON を選ぶ` 後に preview が自動で走り、AI 設定未完了でも file 名、対象 plugin、総 NPC 数、生成対象数、skip 内訳を表示する。
  - AI 設定未完了時は status を `設定未完了` のまま維持し、生成ボタンは無効にする。
  - preview 失敗時は error message を出し、生成ボタンは無効のまま維持する。
  - 再起動後に entry、AI settings、run status の保持と API key 再入力不要を確認できる。
- out_of_scope:
  - master persona 以外の画面や `JOB_PERSONA_ENTRY` 側の機能変更。
  - docs 正本の恒久仕様更新。
  - prompt template の human editable 化。
- open_questions:
  - なし
- required_reading:
  - /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-15-master-persona-management.md
  - /Users/iorishibata/Repositories/AITranslationEngineJP/internal/bootstrap/app_controller.go
  - /Users/iorishibata/Repositories/AITranslationEngineJP/internal/repository/master_persona_repository.go
  - /Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/master_persona_service.go
  - /Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/presenter/master-persona/master-persona.presenter.ts
  - /Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/controller/master-persona/master-persona-screen-controller.ts
  - /Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte

## 作業ブリーフ

- implementation_target: master persona gap closure
- accepted_scope: design bundle、human review、implementation-scope、実装、検証、review closeout まで。
- tasks:
  1. インメモリ concrete が抱える責務を distill で圧縮する。
  2. design で、keyring backed secret concrete、request / SDK transport seam の DI fake、自動 preview と生成許可の分離を requirements / scenario / implementation-brief に揃える。
  3. design-review を通し、`インメモリ禁止` が stop condition と close 不可条件に落ちていることを確認する。
  4. design bundle 完了後に human review を待つ。
  5. human review 後にのみ implementation-scope を確定する。
- implementation_brief_background:
  - 現状は UI と test が通る最小 concrete を優先した結果、master persona だけインメモリ repository とインメモリ secret store に依存している。
  - 現状の fake path は provider の独立 concrete ではなく、service 内の固定ロジックで成立している。
  - 現状の JSON 選択導線は preview 完了を待たないと生成へ進めず、人間の期待とずれている。
- implementation_brief_recommendation:
  - master persona の persistence は SQLite concrete と `github.com/99designs/keyring` backed secret concrete へ分離し、bootstrap で明示的に配線する。
  - fake は provider list の選択肢にせず、HTTP request / SDK transport seam を DI で差し替える。prompt 組み立て、provider validation、run orchestration は real provider と共通にする。
  - JSON 選択イベントで preview を自動起動し、AI 設定未完了でも集計表示だけは出す。生成ボタンは AI 設定完了と preview 成功の両方を満たす時だけ有効化する。
  - implementation-scope では backend persistence、backend secret/keyring、backend provider transport seam、frontend state/UX、tests/review を分離して handoff する。
- implementation_brief_unresolved_items: なし
- validation_commands: `python3 scripts/harness/run.py --suite all` / `python3 scripts/harness/run.py --suite coverage`

## 受け入れ確認

- `NewInMemoryMasterPersonaRepository` と `NewInMemorySecretStore` が production wiring から除去される implementation-scope が作られる。
- インメモリ concrete 禁止が residual risk ではなく stop condition として明記される。
- secret access が `github.com/99designs/keyring` backed concrete として定義され、macOS Keychain / Windows Credential Manager、保存後の再認証不要、保存時の OS 認可ダイアログを acceptance に含める。
- fake path が provider list に出ず、HTTP request / SDK transport seam の DI fake として定義される design が示される。
- prompt 組み立て、provider validation、run orchestration は real provider と共通で、service 内固定生成に依存しないことを acceptance に含める。
- JSON 選択後に自動 preview を走らせ、AI 設定未完了でも集計表示だけは出し、生成ボタンは AI 設定完了と preview 成功の両方を満たす時だけ有効化する導線が scenario と implementation-brief に反映される。
- 再起動後も entry / AI settings / run status が残り、API key の再入力が不要である確認を validation と acceptance に含める。

## 必要エビデンス

- structure harness pass 記録
- active plan
- distill 結果
- requirements
- scenario
- implementation-brief
- keyring backed secret concrete の macOS / Windows 対応確認
- request / SDK transport seam fake の design-review pass 記録

## HITL 状態

- functional_or_design_hitl: `completed-after-design-bundle`
- approval_record: `approved-after-design-bundle`
- approved_at: `2026-04-16`
- approver: `human`
- design_review_result: `pass`
- approval_note: `設計レビュー OK として implementation-scope 以降へ進める`

## 実装スコープ

- artifact: `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.implementation-scope.md`
- split_order:
  1. `backend-master-persona-persistence-and-wiring`
  2. `backend-master-persona-keyring-secret-store`
  3. `backend-master-persona-provider-transport-seam`
  4. `frontend-master-persona-json-preview-gate`
  5. `tests-master-persona-gap-closure`
  6. `review-master-persona-gap-closure`
- split_connections:
  - `backend-master-persona-keyring-secret-store` は `backend-master-persona-persistence-and-wiring` の後に実行する。
  - `backend-master-persona-provider-transport-seam` は `backend-master-persona-keyring-secret-store` の後に実行する。
  - `frontend-master-persona-json-preview-gate` は `backend-master-persona-provider-transport-seam` の後に実行する。
  - `tests-master-persona-gap-closure` は backend 3 split と frontend split の完了後に実行する。
  - `review-master-persona-gap-closure` は tests split の完了後に実行する。
- open_questions: なし

## クローズアウトメモ

- canonicalized_artifacts:
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.implementation-scope.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.scenario.md`
- implementation_review_result: `pass`
- ui_check_result: `pass`
- provider_result: `gemini` / `lm_studio` / `xai` を real provider とし、`internal/infra/ai` の provider interface が共通 response を返す構造へ更新した。
- fake_result: `fake` は provider option / saved provider value に出さず、env `AITRANSLATIONENGINEJP_MASTER_PERSONA_AI_MODE=fake` と test-safe transport seam で扱う。
- validation_result:
  - `go test ./internal/infra/ai ./internal/service ./internal/bootstrap`: pass
  - `npm run lint:backend`: pass
  - `npm --prefix frontend run check`: pass
  - `npm --prefix frontend run test`: pass
  - `npm run lint:frontend`: pass
  - `python3 scripts/harness/run.py --suite coverage`: pass, coverage=71.2%, line=73.7%, branch=45.7%
  - Sonar MCP: HIGH/BLOCKER open 0, reliability open 0, security open 0
- validation_note: `python3 scripts/harness/run.py --suite all` は 1 回目の coverage harness で Sonar API fetch が SSL EOF になった。scanner までは pass し、直後の coverage suite 再実行は pass した。

## 成果

- master persona gap closure の active plan を新規作成した。
- インメモリ concrete、request / SDK transport seam の DI fake、自動 preview 導線を今回の修正対象として固定した。
- SQLite persistence、keyring secret store、AI provider transport seam、JSON auto preview gate を実装した。
- AI request 実装を `internal/infra/ai` に集約し、provider interface と共通 response へ整理した。
- provider option を `gemini` / `lm_studio` / `xai` の real provider だけにし、fake は env / DI seam に閉じた。