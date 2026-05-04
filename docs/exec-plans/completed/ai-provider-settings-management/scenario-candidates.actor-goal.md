# Scenario Candidates: ai-provider-settings-management / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `AIPSM`
- `candidate_count`: 8

## Generator Scope

- `viewpoint`: 利用者の目的、開始操作、成功体験を起点にする。
- `included_sources`:
  - `./plan.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/translation-job-setup-phase-provider-settings/scenario-design.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.implementation-scope.md`
- `excluded_sources`: 引き継ぎ入力外の会話文脈、product code、product test、docs 正本変更、implementation-scope 作成、他 generator 起動。
- `generation_notes`: 候補の採否、統合、最終シナリオ表の確定は designer に残す。

## Candidate Scenarios

### CAND-AIPSM-001 app-shell からプロバイダ設定へ移動する

- `source requirement`: `plan.md:8`, `plan.md:36`, `plan.md:77`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-AIPSM-001`
- `actor`: 利用者
- `trigger`: app-shell の設定導線から、プロバイダ設定画面を開く。
- `expected outcome`: Gemini、xAI、LM Studio のプロバイダ設定へ到達できる。fake provider は利用者向け候補に出ない。
- `observable point`: app-shell navigation、provider settings route、provider list、fake provider 非表示状態。
- `related detail requirement type`: `success_requirement`, `compatibility_requirement`, `testability_requirement`
- `acceptance viewpoint`: UI人間操作E2E。設定導線、画面到達、実 provider 表示、fake provider 非表示を確認する。
- `adoption hint`: provider settings の独立画面と app-shell route の主要正常系候補として扱える。
- `conflict hint`: provider id は `gemini`、`lm_studio`、`xai` のみとする既存判断と合わせる必要がある。

### CAND-AIPSM-002 プロバイダ別に API キーを保存する

- `source requirement`: `plan.md:8`, `plan.md:9`, `plan.md:37`, `plan.md:38`, `spec.md:57`, `spec.md:58`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-AIPSM-002`
- `actor`: 利用者
- `trigger`: Gemini または xAI のプロバイダ設定画面で API キーを入力し、保存する。
- `expected outcome`: API キーは再入力不要な状態で保存される。UI、DTO、log、エラー要約には API キー平文が出ない。
- `observable point`: 保存完了表示、API キー存在状態、再表示時のマスク表示、redacted DTO、redacted log。
- `related detail requirement type`: `success_requirement`, `security_requirement`, `data_requirement`, `observability_requirement`
- `acceptance viewpoint`: UI人間操作E2E と lower-level only。利用者の保存体験と、平文非露出の検査を分けて確認する。
- `adoption hint`: secret persistence の主要正常系候補として扱える。
- `conflict hint`: API キー保存先を DB 暗号化にするか keyring-backed secret store にするかは、DB 変更候補と secret store 境界の判断に影響する。

### CAND-AIPSM-003 プロバイダ別にエンドポイントを保存する

- `source requirement`: `plan.md:8`, `plan.md:9`, `plan.md:37`, `plan.md:40`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-AIPSM-003`
- `actor`: 利用者
- `trigger`: 各プロバイダ設定画面で endpoint を入力または変更し、保存する。
- `expected outcome`: provider ごとの endpoint が保存され、再起動後も同じ provider 設定から参照できる。
- `observable point`: endpoint 入力欄、保存完了表示、再起動後の復元表示、provider settings repository の保存値。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `compatibility_requirement`
- `acceptance viewpoint`: UI人間操作E2E と APIテスト。画面復元と repository 永続化を確認する。
- `adoption hint`: DB migration candidate を含む provider settings persistence の主要正常系候補として扱える。
- `conflict hint`: 旧 Job Setup 設計では LM Studio base URL 設定画面が非対象だったため、今回の独立 provider settings 側へ責務を移す前提を designer が確認する必要がある。

### CAND-AIPSM-004 利用モデルをプロバイダ設定として保存する

- `source requirement`: `plan.md:8`, `plan.md:9`, `plan.md:39`, `spec.md:55`, `spec.md:56`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-AIPSM-004`
- `actor`: 利用者
- `trigger`: プロバイダ設定画面で利用モデルを選択または入力し、保存する。
- `expected outcome`: provider ごとの利用モデルが保存され、翻訳フェーズやマスターペルソナ生成から参照できる。
- `observable point`: model control、保存完了表示、provider settings summary、参照側の model source。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `consistency_requirement`
- `acceptance viewpoint`: UI人間操作E2E と APIテスト。保存操作と参照元の分離を確認する。
- `adoption hint`: model setting の主要正常系候補として扱える。
- `conflict hint`: model 候補を外部 list から選ぶか、手動入力を許可するかは、既存 Job Setup の「手動 model 入力なし」判断と衝突しうる。model の保存単位を provider 共通にするか phase 個別に残すかも統合時に確認が必要である。

### CAND-AIPSM-005 Batch API 利用可否だけを切り替える

- `source requirement`: `plan.md:8`, `plan.md:39`, `spec.md:50`, `spec.md:51`, `translation-job-setup-phase-provider-settings/scenario-design.md:27`, `translation-job-setup-phase-provider-settings/scenario-design.md:28`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-AIPSM-005`
- `actor`: 利用者
- `trigger`: Gemini または xAI のプロバイダ設定画面で Batch API 利用可否を切り替え、保存する。
- `expected outcome`: Gemini と xAI では Batch API を明示的に on / off できる。LM Studio など対象外 provider では Batch API を設定対象にしない。
- `observable point`: Batch API toggle、provider capability 表示、保存完了表示、provider settings summary。
- `related detail requirement type`: `success_requirement`, `boundary_requirement`, `consistency_requirement`
- `acceptance viewpoint`: UI人間操作E2E。対象 provider と対象外 provider の表示差を確認する。
- `adoption hint`: provider capability と利用者の明示設定をつなぐ候補として扱える。
- `conflict hint`: 旧 Job Setup 設計の phase 別 batch mode と、今回の provider 単位 batch toggle の優先順位は統合時に確認が必要である。

### CAND-AIPSM-006 保存済み provider settings を翻訳フェーズから参照する

- `source requirement`: `plan.md:9`, `plan.md:37`, `spec.md:55`, `spec.md:57`, `translation-job-setup-phase-provider-settings/scenario-design.md:20`, `translation-job-setup-phase-provider-settings/scenario-design.md:22`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-AIPSM-006`
- `actor`: 利用者
- `trigger`: provider settings 保存後に、翻訳ジョブ設定または翻訳フェーズ設定を開く。
- `expected outcome`: 翻訳フェーズは endpoint と API キーを独自に保存せず、provider settings の保存状態を参照する。利用者は API キーを再入力しない。
- `observable point`: translation phase settings、credential 参照状態、endpoint source、API キー再入力不要状態。
- `related detail requirement type`: `success_requirement`, `compatibility_requirement`, `consistency_requirement`
- `acceptance viewpoint`: APIテスト と UI人間操作E2E。参照元が provider settings であることと、再入力不要の体験を確認する。
- `adoption hint`: 翻訳フェーズとの永続仕様分離を確認する候補として扱える。
- `conflict hint`: 既存 Job Setup は phase 別 provider / model / credential 参照を持つため、何を provider settings に寄せ、何を phase 設定に残すかの整理が必要である。

### CAND-AIPSM-007 保存済み provider settings をマスターペルソナ生成から参照する

- `source requirement`: `plan.md:9`, `plan.md:37`, `plan.md:86`, `2026-04-16-master-persona-gap-closure.implementation-scope.md:8`, `2026-04-16-master-persona-gap-closure.implementation-scope.md:69`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-AIPSM-007`
- `actor`: 利用者
- `trigger`: provider settings 保存後に、マスターペルソナ生成の AI 設定または生成開始導線を開く。
- `expected outcome`: マスターペルソナ生成は endpoint と API キーを独自に保存せず、provider settings の保存状態を参照する。AI 設定完了後だけ生成開始できる。
- `observable point`: master persona AI settings、provider settings 参照状態、生成ボタン状態、secret 存在状態。
- `related detail requirement type`: `success_requirement`, `compatibility_requirement`, `consistency_requirement`
- `acceptance viewpoint`: UI人間操作E2E と APIテスト。参照元分離と生成可否を確認する。
- `adoption hint`: マスターペルソナ生成との永続仕様分離を確認する候補として扱える。
- `conflict hint`: 既存 master persona の keyring secret store 判断を、全 provider settings へ拡張するか別 namespace にするかは未決になりうる。

### CAND-AIPSM-008 アプリ再起動後も provider settings を利用できる

- `source requirement`: `plan.md:9`, `plan.md:37`, `plan.md:40`, `spec.md:57`, `spec.md:58`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-AIPSM-008`
- `actor`: 利用者
- `trigger`: provider settings を保存し、アプリを再起動して provider settings 画面を再度開く。
- `expected outcome`: endpoint、model、Batch API 利用可否、API キー存在状態が復元される。API キー平文は復元表示されない。
- `observable point`: 再起動後の provider settings UI、repository 保存値、secret 存在状態、redacted summary。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `security_requirement`, `recovery_requirement`
- `acceptance viewpoint`: UI人間操作E2E と APIテスト。永続化復元と API キー非露出を確認する。
- `adoption hint`: DB migration candidate と secret persistence を利用者体験から束ねる候補として扱える。
- `conflict hint`: API キーの保存先と endpoint/model/batch の保存先が分離される場合、復元時の不整合表示が別観点候補と統合対象になる。

## Open Notes

- `human decision candidate`:
  - API キー保存先を DB 暗号化にするか、keyring-backed secret store にするか、既存 master persona secret namespace から分離するかは designer の質問候補である。
  - provider settings の model control を外部 model list 選択だけにするか、手動入力を許可するかは designer の質問候補である。
  - model の保存単位を provider 共通にするか、既存 Job Setup と同じ phase 個別に残すかは designer の質問候補である。
  - provider 単位 Batch API toggle と phase 単位 execution mode の優先順位は designer の質問候補である。
- `merge candidate`:
  - `CAND-AIPSM-002` と `CAND-AIPSM-008` は secret persistence と redaction の候補として統合されうる。
  - `CAND-AIPSM-003`、`CAND-AIPSM-004`、`CAND-AIPSM-005` は provider settings repository / DB migration の候補として統合されうる。
  - `CAND-AIPSM-006` と `CAND-AIPSM-007` は参照側永続仕様分離の候補として統合されうる。
- `rejection candidate`:
  - actor-goal 観点だけでは候補却下を判断しない。
- `conflict candidate`:
  - 旧 Job Setup は phase 別 provider / model / credential / execution mode / batch mode を保持する。一方、今回の plan は endpoint と API キーを provider settings 側へ寄せるため、保存単位の切り分けが競合候補になる。
  - 旧 Job Setup は LM Studio base URL 設定画面を非対象にしていた。一方、今回の plan は各 provider endpoint 保存を対象にするため、LM Studio endpoint の扱いが競合候補になる。
