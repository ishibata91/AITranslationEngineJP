# Scenario Candidates: translation-job-management / external-integration

- `generator`: `external-integration`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJM`

## Generator Scope

- `viewpoint`: `external-integration`
- `included_sources`: `tasks/usecases/translation-job-management.yaml`, `docs/spec.md`, `docs/exec-plans/completed/translation-job-setup/scenario-design.md`
- `excluded_sources`: `引き継いでいない会話文脈`, `final scenario matrix`, `product code`, `product test`, `docs 正本変更`, `他 viewpoint の採否判断`
- `generation_notes`: `入力キャッシュ、保存済み AI 設定、secret 参照、provider / network 状態、fake transport に差し替えられる境界だけを external-integration 候補として残す。`

## Candidate Scenarios

### CAND-TJM-001 入力キャッシュ欠落時に再開不可理由を表示する

- `source requirement`: `tasks/usecases/translation-job-management.yaml` の「入力キャッシュ状態」「再開可否」「再開不可の job では入力キャッシュ欠落や terminal state などの理由を確認できる」、`docs/spec.md` の「抽出 JSON を正本として保持しつつ実行キャッシュへ取り込めること」「未完了ジョブが参照していない入力キャッシュを削除して再構築可能な状態を維持できること」、`translation-job-setup` の `REQ-TJS-004`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJM-001`
- `actor`: `ユーザー`
- `trigger`: `入力キャッシュが欠落した未完了 job を未完了ジョブ一覧または Job Run で開く。`
- `expected outcome`: `対象 job は表示できるが、再開可否は不可になり、入力キャッシュ欠落が再開不可理由として観測できる。`
- `observable point`: `未完了ジョブ一覧、Job Run の再開可否、再開不可理由、入力出自表示、cache rebuild 導線の有無`
- `related detail requirement type`: `failure_handling_requirement`, `file_boundary_requirement`, `display_requirement`
- `fake_or_stub`: `cache missing fixture`, `fixed input provenance fixture`, `temp DB`
- `adoption hint`: `入力キャッシュを外部ファイル境界として扱い、再開不可理由を user-facing scenario に残すなら採用候補。`
- `conflict hint`: `cache rebuild の実行導線は input-intake 側や failure viewpoint と統合される可能性がある。`

### CAND-TJM-002 入力出自を失わず既存 job を開く

- `source requirement`: `tasks/usecases/translation-job-management.yaml` の「job ごとの入力出自、状態、現在フェーズ、進捗を確認できる」「既存 job がある入力では新規作成ではなく既存 job を開ける」、`docs/spec.md` の「入力ファイルの出自を失わずに保持できること」、`translation-job-setup` の `REQ-TJS-001` と `REQ-TJS-004`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJM-002`
- `actor`: `ユーザー`
- `trigger`: `既存 job がある入力データを選ぶ、または未完了ジョブ一覧から job を選択する。`
- `expected outcome`: `新規 job 作成ではなく既存 job が開き、xEdit 抽出データの入力出自と job の対応を確認できる。`
- `observable point`: `未完了ジョブ一覧、選択中 job、入力データ名、入力出自、job detail の input 参照`
- `related detail requirement type`: `compatibility_requirement`, `file_boundary_requirement`, `consistency_requirement`
- `fake_or_stub`: `existing job fixture`, `xEdit extracted data fixture`, `temp DB`
- `adoption hint`: `入力ファイル由来の外部境界と job 管理 UI の対応を final scenario へ残すなら採用候補。`
- `conflict hint`: `同一入力の既存 job を開く導線は actor-goal viewpoint と merge 候補になる。`

### CAND-TJM-003 保存済み AI 設定と secret 参照状態を平文なしで表示する

- `source requirement`: `tasks/usecases/translation-job-management.yaml` の「現在フェーズ」「再開可否」、`docs/spec.md` の「翻訳ジョブ、API の実行進捗を確認できること」「各フェーズの API 選択、APIKey は再入力不要で保存ができること」「APIKey は暗号化して保存すること」、`translation-job-setup` の `REQ-TJS-003`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJM-003`
- `actor`: `ユーザー`
- `trigger`: `保存済み provider / model / credential 参照を持つ未完了 job を開く。`
- `expected outcome`: `現在フェーズの provider、model、execution mode、credential 参照状態を確認できるが、API key 平文または復号可能な値は UI、要約、エラー表示に出ない。`
- `observable point`: `Job Run の現在フェーズ表示、provider / model 要約、credential 参照状態、secret redaction、error summary`
- `related detail requirement type`: `security_requirement`, `external_integration_requirement`, `display_requirement`
- `fake_or_stub`: `fake secret store`, `saved AI setting fixture`, `redaction assertion`
- `adoption hint`: `job management から実行中または再開待ちの AI 設定を観測させるなら採用候補。`
- `conflict hint`: `secret 参照不能を再開不可理由に含めるかは CAND-TJM-005 や failure viewpoint と重なる。`

### CAND-TJM-004 実行中 provider request がある job は削除ではなく停止操作に寄せる

- `source requirement`: `tasks/usecases/translation-job-management.yaml` の「実行中 job は削除ではなく停止操作を実行できる」「実行中 job は削除できず、停止または中断後に削除可否を再判定する」、`docs/spec.md` の「翻訳ジョブの中断、再開、失敗回復が継続的に行えること」「翻訳ジョブ、API の実行進捗を確認できること」
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJM-004`
- `actor`: `ユーザー`
- `trigger`: `provider request または Batch API 実行を含む Running job で削除または停止可否を確認する。`
- `expected outcome`: `削除は不可として表示され、停止操作だけが選べる。停止後は外部 request の停止中または中断結果を観測でき、削除可否は状態更新後に再判定される。`
- `observable point`: `停止可否、削除可否、runtime adapter stop result、phase state、late response の扱い`
- `related detail requirement type`: `network_boundary_requirement`, `state_requirement`, `failure_handling_requirement`
- `fake_or_stub`: `fake transport with cancellable request`, `running phase fixture`, `late response fixture`
- `adoption hint`: `実行中の外部通信を job 管理の操作制約へ反映するなら採用候補。`
- `conflict hint`: `停止後の状態遷移と late response rejection は state-transition / failure viewpoint と競合しうる。`

### CAND-TJM-005 再開前に provider capability と credential 参照を再確認する

- `source requirement`: `tasks/usecases/translation-job-management.yaml` の「再開入口」「再開可否」「再開不可理由」、`docs/spec.md` の「LMStudio、Gemini、xAI を利用できること」「Gemini, xAI は BatchAPI が利用できること」「各フェーズではいずれのモデルでも選択できる」、`translation-job-setup` の `REQ-TJS-003`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJM-005`
- `actor`: `ユーザー`
- `trigger`: `Paused または RecoverableFailed の job で再開入口を開く。`
- `expected outcome`: `保存済み provider、model、execution mode、credential 参照が再開可能条件を満たす場合だけ再開入口が有効になる。不整合または参照不能の場合は再開不可理由が表示される。`
- `observable point`: `再開可否、再開不可理由、provider capability check、credential resolution result、phase resume precheck`
- `related detail requirement type`: `external_integration_requirement`, `compatibility_requirement`, `failure_handling_requirement`
- `fake_or_stub`: `provider capability fixture`, `fake secret store`, `paused job fixture`, `recoverable failed job fixture`
- `adoption hint`: `再開入口を外部 provider 設定の現在状態と結び付けるなら採用候補。`
- `conflict hint`: `再開 precheck を job management で行うか phase 実行開始時に遅延するかは human decision 候補。`

### CAND-TJM-006 paid API なしで job management の provider 境界を検証する

- `source requirement`: `docs/spec.md` の「Gemini, xAI は BatchAPI が利用できること」「失敗しても特に別プロバイダフォールバックは必要ない」、`translation-job-setup` の「paid な real AI API を scenario validation の前提にしない」「test では外部 request / SDK transport だけを fake に差し替える」
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJM-006`
- `actor`: `テスト実行者`
- `trigger`: `job management の再開可否、停止可否、現在フェーズ表示を scenario validation で検証する。`
- `expected outcome`: `user-facing provider list は real provider のまま保ち、外部 request または SDK transport だけを fake に差し替えて、paid API を呼ばずに provider 境界の結果を観測できる。`
- `observable point`: `provider list、fake transport wiring、request log、resume precheck result、stop result、external request 未実行証跡`
- `related detail requirement type`: `testability_requirement`, `external_integration_requirement`, `security_requirement`
- `fake_or_stub`: `fake transport`, `fixed provider response`, `request log assertion`
- `adoption hint`: `job management の受け入れテストを外部 API 費用に依存させない条件として残すなら採用候補。`
- `conflict hint`: `user-facing scenario ではなく lower-level acceptance に落とす判断がありうる。`

### CAND-TJM-007 非実行中 job 削除で入力データと再構築可能性を保持する

- `source requirement`: `tasks/usecases/translation-job-management.yaml` の「実行中ではない job は入力データを残したまま削除できる」、`docs/spec.md` の「入力ファイルの出自を失わずに保持できること」「未完了ジョブが参照していない入力キャッシュを削除して再構築可能な状態を維持できること」、`translation-job-setup` の `REQ-TJS-004`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJM-007`
- `actor`: `ユーザー`
- `trigger`: `Paused、RecoverableFailed、Failed、Canceled など実行中ではない job の削除を実行する。`
- `expected outcome`: `job は一覧から外れるが、入力データ、入力出自、抽出 JSON 正本、再構築可能性は失われない。入力キャッシュ削除の扱いは参照状況に従って観測できる。`
- `observable point`: `job list、input data list、input provenance、cache reference state、delete result`
- `related detail requirement type`: `file_boundary_requirement`, `persistence_requirement`, `consistency_requirement`
- `fake_or_stub`: `non-running job fixture`, `input cache reference fixture`, `temp DB`
- `adoption hint`: `削除操作を外部入力ファイル境界と分離して固定するなら採用候補。`
- `conflict hint`: `削除可否そのものは lifecycle / state-transition viewpoint と merge 候補になる。`

## Open Notes

- `human decision candidate`: `CAND-TJM-005` の provider / credential 再確認を job management entry で行うか、phase resume boundary まで遅延するかは未確定である。
- `human decision candidate`: `CAND-TJM-004` の停止操作で外部 request を即時 cancel するか、停止要求を記録して late response を破棄するだけにするかは未確定である。
- `merge candidate`: `CAND-TJM-001` と `CAND-TJM-007` は入力キャッシュと入力出自の file boundary として統合される可能性がある。
- `merge candidate`: `CAND-TJM-003` と `CAND-TJM-005` は AI 設定復元と再開 precheck の 1 連シナリオへ統合される可能性がある。
- `rejection candidate`: `CAND-TJM-006` は user-facing scenario ではなく lower-level acceptance または testability 条件へ落とす判断がありうる。
