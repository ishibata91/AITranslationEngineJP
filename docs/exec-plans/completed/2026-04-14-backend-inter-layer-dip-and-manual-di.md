# 実装計画

- workflow: orchestrate
- task_mode: refactor
- status: completed
- lane_owner: orchestrate
- scope: backend-inter-layer-dip-and-manual-di
- task_id: backend-inter-layer-dip-and-manual-di
- task_catalog_ref: N/A
- parent_phase: implementation-lane

## 要求要約

- 現在実装済みの master-dictionary backend flow に限定して層間 DIP を導入する。
- DI コンテナは使わず、composition root で手動 DI を行う。
- coverage 目標のため、上位層 unit test が結合テスト化しない構造へ改める。

## 判断根拠

- 2026-04-14 継続判断として `task_mode: refactor` を固定した。主目的は既存 master-dictionary backend flow の依存方向整理、manual DI 化、unit test と integration test の責務分離であり、新規機能追加ではない。
- orchestrate の routing は `distill -> design -> implement/tests -> review` とし、existing active plan を継続利用する。


<!-- Decision Basis -->

- 既存 active plan は master dictionary と SQLite 永続化を中心とした局所 task だった。
- review diff は対象を implemented master-dictionary flow only に限定している。
- phase-2.5 reroute の前提として、scope は現在実装済みの master-dictionary backend flow のみとし、placeholder package、docs-only package、未実装 layer のための port は先回り定義しない。
- `docs/tech-selection.md` は `手動 DI` を採用済みであり、DI コンテナ導入は不要である。
- 高 coverage を unit test 主体で維持するには、master-dictionary の controller / usecase / service / repository 境界と Wails/XML runtime concern で concrete 実装への直結を避ける必要がある。
- `docs/architecture.md` の現行記述には repository concrete 許容が残るため、恒久仕様更新は行わず、この plan と実装で先行整合する。

## 対象範囲

- `internal/controller/wails/master_dictionary_controller.go`
- `internal/usecase/master_dictionary_usecase.go`
- `internal/usecase/master_dictionary_runtime_event_publisher.go`
- `internal/service/master_dictionary_query_service.go`
- `internal/service/master_dictionary_command_service.go`
- `internal/service/master_dictionary_import_service.go`
- `internal/service/master_dictionary_runtime_event_publisher.go`
- `internal/service/master_dictionary_runtime_event_publisher_wails.go`
- `internal/repository/master_dictionary_repository.go`
- `internal/controller/wails/app_controller.go` と `main.go` の default wiring / bootstrap
- master-dictionary backend flow の unit / integration test 責務再配置

## 対象外

- frontend の画面仕様変更
- docs 正本の恒久仕様更新
- DI コンテナ導入
- backend 以外の層構成変更

## 依存関係・ブロッカー

- `docs/architecture.md` の依存方向を壊さないこと。
- `docs/tech-selection.md` の `手動 DI` 方針を守ること。
- どの境界で interface を定義し、どこで default wiring するかを detail design で固定すること。
- Wails runtime、filesystem、`encoding/xml`、repository concrete の依存を upper layer test へ漏らさないこと。

## 並行安全メモ

- master-dictionary flow の controller / usecase / service / repository / runtime adapter へ波及するため、先に interface 境界と composition root を固定しないと実装 handoff が競合する。
- unit test と integration test の責務分離を先に決める。
- repository concrete と Wails/XML runtime adapter の実行証明は integration 側へ寄せる。

## 機能要件

- `summary`:
  - 現在実装済みの master-dictionary backend flow の層間呼び出しを、現行の層順序を維持したまま manual DI 前提の差し替え可能な依存へ揃える。
  - production の依存グラフは `main.go` と `AppController` 側 bootstrap の手動配線で組み立て、DI コンテナ、service locator、反射ベース自動配線は持ち込まない。
  - `controller` / `usecase` / `service` の unit test は test double だけで成立させ、上位層 test が実質 integration test 化しない境界を固定する。
  - `repository` concrete と Wails/XML runtime adapter の実行証明は lower-layer integration test または wiring test へ隔離する。
- `in_scope`:
  - 現在実装済みの master-dictionary backend flow に対して、`controller -> usecase -> service -> repository` と runtime event / XML import の各呼び出し境界で、上流が下流 concrete 実装へ直結しない抽象 seam を定義する。
  - production constructor と default wiring を manual DI 前提へ揃え、master-dictionary backend の既定起動経路を 1 つの composition root から再現できるようにする。
  - master-dictionary flow へ caller-owned seam ルールを適用する。
  - upper-layer orchestration / branching test を unit test へ残し、concrete resource を使う証明を integration test または wiring test へ再配置する。
  - frontend 契約を変えずに master-dictionary backend の default graph を起動可能なまま保つ。
- `non_functional_requirements`:
  - DI コンテナ、service locator、反射ベース自動配線、生成コード前提の配線機構を導入しない。
  - `controller` / `usecase` / `service` の unit test は実 filesystem、`encoding/xml` decoder、Wails runtime 初期化、repository concrete に依存しない。
  - 依存方向は `Controller -> UseCase -> Service -> Repository` と runtime adapter 補助境界を維持し、上位層から lower concrete への漏れを増やさない。
  - coverage 向上のための seam 導入で、production と test の業務ロジック重複や分岐乖離を作らない。
  - 後続実装は backend lint / execution / coverage の validation command を通過できる前提で設計する。
- `out_of_scope`:
  - frontend 画面仕様、frontend architecture、Wails bind 契約の変更。
  - `docs/` 正本更新による恒久仕様確定。
  - DB 種別や外部 provider 技術選定そのものの見直し。
  - DIP 目的を超える domain behavior の変更。
  - DI コンテナや plugin 型拡張機構の導入。
  - 未実装 package、placeholder package、docs-only package のための先回り port 定義。
- `design_decisions`:
  - interface ownership は review diff に合わせて caller-owned とする。`controller/wails` は `UsecasePort` と `RuntimeEmitterSource`、`usecase` は `QueryServicePort` / `CommandServicePort` / `ImportServicePort` / `RuntimeEventPublisherPort`、`service` は `RepositoryPort` / `XMLFilePort` / `XMLRecordReaderPort` / `RuntimeContextPort` を自 package 側に置く。
  - composition root は `main.go` と `internal/controller/wails/app_controller.go` の bootstrap に置く。default wiring、concrete constructor 呼び出し、Wails bind 組み立て、wiring test 用 entrypoint はこの bootstrap 側へ集約する。
  - runtime concern は review diff にある現在実装済みのものだけを扱う。Wails event emitter handle は `RuntimeEmitterSource` / Wails event adapter、filesystem access と XML reader は filesystem/XML adapter、repository concrete は repository adapter に閉じ込め、未実装の DB / HTTP / AI provider handle は本 plan の設計対象へ含めない。
- `required_reading`:
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-14-backend-inter-layer-dip-and-manual-di.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/tech-selection.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/master-dictionary.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/master-dictionary-management.md`

## UI モック

- `artifact_path`: `N/A`
- `final_artifact_path`: `N/A`
- `summary`:
  - UI 追加要件はない。

## Scenario テスト一覧

- `artifact_path`: `docs/scenario-tests/master-dictionary-management.md`
- `final_artifact_path`: `docs/scenario-tests/master-dictionary-management.md`
- `template_path`: `docs/exec-plans/templates/scenario-tests.md`
- `summary`:
  - 既存 scenario のうち backend 依存注入変更で壊れ得る導線を回帰対象とする。

## 実装計画

<!-- Implementation Plan -->

- `ordered_scope`:
  1. review diff に合わせて caller-owned interface を固定する。`controller/wails` は `UsecasePort` と `RuntimeEmitterSource`、`usecase` は `QueryServicePort` / `CommandServicePort` / `ImportServicePort` / `RuntimeEventPublisherPort`、`service` は `RepositoryPort` / `XMLFilePort` / `XMLRecordReaderPort` / `RuntimeContextPort` を所有し、cross-layer 共通 `contracts` package は作らない。
  2. composition root は `main.go` と `AppController` bootstrap に寄せる。`MasterDictionaryController` と `MasterDictionaryUsecase` の self-wire を外し、default wiring、concrete 実装の new、Wails bind 組み立て、wiring test 用 entrypoint は bootstrap 側へ集約する。
  3. service 層は query / command / import の 3 seam と runtime/event 補助境界へ分割し、repository concrete、os、`encoding/xml`、Wails emitter handle を service core から直接引かない。
  4. runtime adapter は review diff の現在実装済み範囲に限定する。Wails event adapter と filesystem/XML adapter を bootstrap で生成し、上位層 unit test は fake port のみで成立させ、Wails と filesystem/XML と repository adapter の実証は wiring / integration へ退避する。
  5. validation は `structure` で計画前提を守り、`backend-lint` と `execution` で refactor の成立を確認し、`coverage` と最終 `all` で unit / integration 責務分離後の回帰を確認する。
- `parallel_task_groups`:
  - `group_id`: `backend-dip-contract-and-bootstrap`
    - `can_run_in_parallel_with`: なし
    - `blocked_by`: なし
    - `completion_signal`: interface ownership、constructor 注入方針、composition root 配置、default wiring 境界、unit と integration の責務分離が固定されている。
  - `group_id`: `backend-dip-upper-layer-refactor`
    - `can_run_in_parallel_with`: `backend-dip-lower-adapter-refactor`
    - `blocked_by`: `backend-dip-contract-and-bootstrap`
    - `completion_signal`: controller / usecase / service が caller-owned port 前提で constructor 注入へ移行し、upper-layer unit test が concrete 依存なしで通る。
  - `group_id`: `backend-dip-lower-adapter-refactor`
    - `can_run_in_parallel_with`: `backend-dip-upper-layer-refactor`
    - `blocked_by`: `backend-dip-contract-and-bootstrap`
    - `completion_signal`: repository concrete、Wails event adapter、filesystem/XML adapter が lower adapter として整理され、bootstrap からのみ default wiring される。
  - `group_id`: `backend-dip-validation`
    - `can_run_in_parallel_with`: なし
    - `blocked_by`: `backend-dip-upper-layer-refactor`, `backend-dip-lower-adapter-refactor`
    - `completion_signal`: wiring test、integration test、coverage 証跡が揃い、manual DI 化後も backend 回帰がない。
- `task_dependencies`:
  - `task_id`: `backend-dip-contract-and-bootstrap`
    - `depends_on`: なし
    - `enables`: `backend-dip-upper-layer-refactor`, `backend-dip-lower-adapter-refactor`
    - `reason`: interface ownership と composition root policy が未確定だと、各 layer が別々の abstraction を増やして競合し、manual DI の配線点も分裂する。
  - `task_id`: `backend-dip-upper-layer-refactor`
    - `depends_on`: `backend-dip-contract-and-bootstrap`
    - `enables`: `backend-dip-validation`
    - `reason`: controller / usecase / service の constructor 注入と fake-friendly port が先に揃わないと、upper-layer unit test を結合テストから切り離せない。
  - `task_id`: `backend-dip-lower-adapter-refactor`
    - `depends_on`: `backend-dip-contract-and-bootstrap`
    - `enables`: `backend-dip-validation`
    - `reason`: repository concrete と Wails/XML runtime adapter を bootstrap 配線へ寄せないと、default wiring と integration 証明の責務を upper layer から外せない。
  - `task_id`: `backend-dip-validation`
    - `depends_on`: `backend-dip-upper-layer-refactor`, `backend-dip-lower-adapter-refactor`
    - `enables`: `tests`, `implement`, `review`
    - `reason`: unit と integration の責務分離、manual DI 起動、既存 backend 導線の回帰なしは、upper layer と lower adapter の両方が揃ってからでないと証明できない。
- `tasks`:
  - `task_id`: `backend-dip-contract-and-bootstrap`
    - `owned_scope`: caller-owned interface rules for the master-dictionary flow, constructor signature alignment, bootstrap placement in main/AppController, default wiring entrypoint, and wiring test boundary
    - `depends_on`: なし
    - `parallel_group`: `backend-dip-contract-and-bootstrap`
    - `required_reading`:
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-14-backend-inter-layer-dip-and-manual-di.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/tech-selection.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/master-dictionary.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/master-dictionary-management.md`
    - `validation_commands`:
      - `python3 scripts/harness/run.py --suite structure`
  - `task_id`: `backend-dip-upper-layer-refactor`
    - `owned_scope`: master-dictionary controller / usecase / service constructor injection, caller-owned port placement, fake-friendly unit seam creation, and replacement of direct concrete imports across upper layers
    - `depends_on`: `backend-dip-contract-and-bootstrap`
    - `parallel_group`: `backend-dip-upper-layer-refactor`
    - `required_reading`:
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-14-backend-inter-layer-dip-and-manual-di.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/master-dictionary-management.md`
    - `validation_commands`:
      - `python3 scripts/harness/run.py --suite backend-lint`
      - `python3 scripts/harness/run.py --suite execution`
  - `task_id`: `backend-dip-lower-adapter-refactor`
    - `owned_scope`: repository adapter, Wails event adapter, and filesystem/XML adapter split, lower-layer concrete constructor cleanup, default wiring assembly in bootstrap, and isolation of current runtime concretes behind ports
    - `depends_on`: `backend-dip-contract-and-bootstrap`
    - `parallel_group`: `backend-dip-lower-adapter-refactor`
    - `required_reading`:
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/tech-selection.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-14-backend-inter-layer-dip-and-manual-di.md`
    - `validation_commands`:
      - `python3 scripts/harness/run.py --suite backend-lint`
      - `python3 scripts/harness/run.py --suite execution`
  - `task_id`: `backend-dip-validation`
    - `owned_scope`: upper-layer unit test relocation, wiring tests for default manual DI, concrete integration coverage for repository and Wails/XML runtime boundaries, and final evidence packaging for review
    - `depends_on`: `backend-dip-upper-layer-refactor`, `backend-dip-lower-adapter-refactor`
    - `parallel_group`: `backend-dip-validation`
    - `required_reading`:
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-14-backend-inter-layer-dip-and-manual-di.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/master-dictionary-management.md`
      - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/master-dictionary.md`
    - `validation_commands`:
      - `python3 scripts/harness/run.py --suite backend-lint`
      - `python3 scripts/harness/run.py --suite execution`
      - `python3 scripts/harness/run.py --suite coverage`
      - `python3 scripts/harness/run.py --suite all`
- `owned_scope`:
  - backend のうち現在実装済みの master-dictionary flow のみ。upper layer は caller-owned port と constructor injection への移行を担当し、lower layer は repository concrete、Wails event adapter、filesystem/XML adapter を閉じ込める。composition root は `main.go` と `AppController` bootstrap が唯一の default wiring 所有者となる。
- `required_reading`:
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-14-backend-inter-layer-dip-and-manual-di.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/tech-selection.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/master-dictionary.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/master-dictionary-management.md`
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite structure`
  - `python3 scripts/harness/run.py --suite backend-lint`
  - `python3 scripts/harness/run.py --suite execution`
  - `python3 scripts/harness/run.py --suite coverage`
  - `python3 scripts/harness/run.py --suite all`
- `implementation_plan_updates`:
  - exact scope を currently implemented master-dictionary backend flow only に縮小し、未実装 package や future layer の port は plan から外した。
  - caller-owned interface を review diff と一致する port 名と責務へ固定し、cross-layer 共通 contracts package を禁止する方針を fixed した。
  - composition root は `main.go` と `AppController` bootstrap に固定し、controller/usecase の self-wire を外す方針を fixed した。
  - runtime concern は現在実装済みの Wails event adapter、filesystem/XML adapter、repository adapter だけへ絞り、DB / HTTP / AI provider handle は plan 対象外へ外した。
  - upper-layer refactor と lower-adapter refactor を contract 固定後に並列化し、最後に wiring / integration / coverage を集約する dependency order を fixed した。

## 受け入れ確認

- master-dictionary の controller / usecase / service の unit test が concrete 実装を起動しない。
- composition root の手動 DI だけで master-dictionary backend flow が起動可能である。
- repository concrete と Wails/XML runtime adapter の実行証明が integration または wiring test へ分離される。
- backend lint / execution / coverage を満たす。

## 必要な証跡

<!-- Required Evidence -->

- master-dictionary inter-layer DIP 導入後の unit test 証跡
- manual DI wiring の起動証跡
- repository/Wails/XML adapter の concrete integration test 証跡
- backend lint / execution / coverage の結果
- implementation review の結果

## 機能要件 HITL 状態

- approved

## 機能要件 承認記録

- 2026-04-14 human: 全てのレイヤー間を DIP し、DI コンテナを使わず手動 DI とする。coverage 目標のため unit test を結合テスト化させない方針を要求。

## 詳細設計 HITL 状態

- approved

## 詳細設計 承認記録

- 2026-04-14 orchestrate: `legacy MasterDictionaryService` は削除する。`internal/controller/wails/app_controller.go` を唯一の backend composition root とし、`main.go` は Wails 起動だけを担う。caller-owned port は `controller/wails` に `MasterDictionaryUsecasePort`、`usecase` に `QueryServicePort` / `CommandServicePort` / `ImportServicePort` / `RuntimeEventPublisherPort`、`service` に `RepositoryPort` / `XMLFilePort` / `XMLRecordReaderPort` / `ImportProgressEmitterPort` を置く。bind 名、runtime event 名、payload 互換は維持する。

## review 用差分図

- `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-14-backend-inter-layer-dip-and-manual-di.review-structure-diff.d2`
- `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-14-backend-inter-layer-dip-and-manual-di.review-structure-diff.svg`
- caller-owned port、manual composition root、test boundary shift を review 用差分図へ集約した。

## 差分正本適用先

- `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/components/backend/master-dictionary-management.d2`

## Implementation Brief

- Slice A `contract-and-bootstrap`: constructor 署名、caller-owned port、default wiring 入口、wiring test 境界を固定する。
- Slice B `upper-layer-refactor`: `controller / usecase / service` を constructor injection 化し、upper-layer unit test を fake port 前提へ移す。
- Slice C `lower-adapter-refactor`: repository concrete、Wails event adapter、filesystem/XML adapter を lower adapter として閉じ込める。
- Slice D `validation`: legacy service test を分解し、unit / integration / wiring の証跡を集約する。
- 依存順序は `A -> (B || C) -> D` とする。

## Closeout Notes

- `docs/` 正本は human 先行でのみ更新する。
- 現行 architecture 記述との差分は review と residual risk に残す。

## 検証結果

- `python3 scripts/harness/run.py --suite structure`: pass
- `python3 scripts/harness/run.py --suite backend-lint`: pass
- `python3 scripts/harness/run.py --suite execution`: pass
- `python3 scripts/harness/run.py --suite coverage`: pass
- `python3 scripts/harness/run.py --suite all`: pass

## Sonar Gate 結果

- `HIGH` / `BLOCKER` open issue: 0 件
- open reliability issue: 0 件
- open security issue: 0 件
- project quality gate status: `NONE`
- close 判定は issue-count gate を採用し、quality gate status は補足扱いとした。

## 実装レビュー結果

- 初回 implementation-review は `reroute` だった。理由は fake-only service unit test 未達と bootstrap wiring test 不足だった。
- reroute 指摘に対して、service unit test の concrete 依存除去、bootstrap wiring test の追加、runtime emitter 契約修正、generated binding の `TS2305` 修正を入れた。
- 修正後の implementation-review 再実行で `pass / no findings` を確認した。
- caller-owned port、manual DI in bootstrap、controller/usecase self-wire 排除、fake-only service unit test、bootstrap wiring test を最終確認した。

## 結果

<!-- Outcome -->

- completed
