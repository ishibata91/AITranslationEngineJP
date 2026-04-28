# Scenario Candidates: translation-job-setup / external-integration

- `generator`: `external-integration`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJS`

## Generator Scope

- `viewpoint`: `external-integration`
- `included_sources`: `tasks/usecases/translation-job-setup.yaml`, `docs/spec.md`, `docs/er.md`, `docs/architecture.md`, `docs/exec-plans/completed/translation-input-intake/scenario-design.md`, `docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.scenario.md`
- `excluded_sources`: `引き継いでいない会話文脈`, `final scenario matrix`, `product code`, `product test`, `他 viewpoint の採否判断`
- `generation_notes`: `AI 基盤設定の保存と参照、provider capability、secret store、transport seam、paid API 非依存検証だけを external-integration 候補として残す。`

## Candidate Scenarios

### CAND-TJS-001 保存済み AI 設定と secret store 参照を Job Setup へ復元する

- `source requirement`: `docs/spec.md` の「各フェーズの API 選択、APIKey は再入力不要で保存できること」「基盤参照、AI 基盤、実行方式を選択できること」、`docs/er.md` の `credential_ref` は secret store 参照だけを保持すること、`docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.scenario.md` の `SCN-MPG-001`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJS-001`
- `actor`: `ユーザー`
- `trigger`: `保存済み provider / model / API key secret がある状態で Job Setup を開く。`
- `expected outcome`: `Job Setup は保存済み provider と model を復元し、API key 平文を再表示せずに secret store 参照で validation を開始できる。`
- `observable point`: `Job Setup の AI 設定表示、backend settings query、secret store read 成否、validation 開始可否、UI と error surface に API key 平文が出ないこと`
- `related detail requirement type`: `success_requirement`, `security_requirement`, `compatibility_requirement`
- `adoption hint`: `Job Setup の再入力不要要件と secret 非露出要件を final scenario へ残すなら採用候補。`
- `conflict hint`: `保存済み設定の復元を draft 復元や state-transition と一体化するなら lifecycle / state-transition 側候補と merge 候補。secret 参照不能時の扱いは failure 候補へ分離される可能性がある。`

### CAND-TJS-002 provider capability に応じて execution mode と validation 結果を切り替える

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の「基盤参照、AI 基盤、実行方式を選択できる」「create job の前に validation を通せる」、`docs/spec.md` の `LMStudio` / `Gemini` / `xAI` が利用可能で `Gemini, xAI` は `BatchAPI` を利用できること
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJS-002`
- `actor`: `ユーザー`
- `trigger`: `Job Setup で provider、model、execution mode を変更して validation を実行する。`
- `expected outcome`: `provider list は real provider だけを表示し、Batch API 非対応 provider では batch 実行方式を通さず、対応 provider では capability に沿った validation 結果を返す。`
- `observable point`: `provider dropdown、execution mode 選択肢、validation pass / fail 理由、create job 可否`
- `related detail requirement type`: `success_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: `runtime 選択と execution mode 選択を Job Setup の受け入れ条件へ含めるなら採用候補。`
- `conflict hint`: `unsupported provider / execution mode の失敗理由は failure viewpoint と重なる。provider 選択 UI の詳細は actor-goal / ui-design 側候補と merge される可能性がある。`

### CAND-TJS-003 fake transport seam で Job Setup validation を paid API なしに検証できる

- `source requirement`: `docs/architecture.md` の `Service` は AI 実行が必要な機能で `AIProvider` を使うこと、`AIProvider` が provider 差異を吸収すること、`docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.scenario.md` の `SCN-MPG-004`、`.codex/skills/scenario-design/SKILL.md` の paid な real AI API を system test 前提にしないこと
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJS-003`
- `actor`: `テスト実行者`
- `trigger`: `test mode または scenario validation で Job Setup の provider validation を実行する。`
- `expected outcome`: `Job Setup は real provider list を保ったまま AIProvider 経由の共通 validation 経路を通し、外部 request / SDK transport だけを fake に差し替えて paid API を呼ばない。`
- `observable point`: `provider list、test mode の adapter wiring、validation response、外部 request 未実行の証跡`
- `related detail requirement type`: `testability_requirement`, `compatibility_requirement`, `security_requirement`
- `adoption hint`: `Job Setup の scenario を implementation 後に system-level で固定したいなら採用候補。`
- `conflict hint`: `fake mode の可視化方法は operation-audit と競合する可能性がある。user-visible scenario にするか lower-level 寄り受け入れ条件に留めるかは designer 判断になる。`

### CAND-TJS-004 secret 参照不能または provider mismatch では validation failure を返し create job を止める

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の「create job の前に validation を通せる」「validation failure の理由を確認できる」、`docs/spec.md` の API key 保存と provider 選択要件、`docs/er.md` の `credential_ref` は secret store 参照だけを保持すること
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJS-004`
- `actor`: `ユーザー`
- `trigger`: `credential_ref が解決できない、provider に対して model / execution mode が不整合、または外部接続前提が不足した状態で validation を実行する。`
- `expected outcome`: `validation は失敗理由を UI で確認可能に返し、create job は許可されず、API key 平文や secret 本体を露出しない。`
- `observable point`: `validation error kind、UI の failure reason、create job button state、job 未作成の永続化結果`
- `related detail requirement type`: `failure_handling_requirement`, `security_requirement`, `state_requirement`
- `adoption hint`: `create job 前 validation の hard gate を scenario へ固定するなら採用候補。`
- `conflict hint`: `失敗分類と retry 導線は failure / lifecycle viewpoint と衝突しやすい。どこまで Job Setup で検知し、どこから phase 実行時 failure に回すかは human decision 候補。`

## Open Notes

- `human decision candidate`: `provider reachability や secret store 参照不能を Job Setup validation の必須失敗条件にするか、初回 phase 実行まで遅延させるかは未確定。`
- `merge candidate`: `CAND-TJS-001` と `CAND-TJS-004` は AI 設定復元と validation gate の 1 連シナリオへ統合される可能性がある。`
- `rejection candidate`: `CAND-TJS-003` は user-facing scenario ではなく lower-level acceptance に落とす判断がありうる。`
