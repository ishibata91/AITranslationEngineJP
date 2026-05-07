# Scenario Candidates: 2026-05-07-provider-settings-job-decoupling-implement / state-transition

- `generator`: `state-transition`
- `source_plan`: `../2026-05-07-provider-settings-job-decoupling/plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `PSJD-ST`
- `candidate_count`: `7`

## Generator Scope

- `viewpoint`: state-transition
- `included_sources`: `./task-frame.md`, `../2026-05-07-provider-settings-job-decoupling/plan.md`, `../2026-05-07-provider-settings-job-decoupling/light-change-planning.md`, `docs/detail-specs/ai-provider-settings-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/er.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本本文の変更、採否判断、統合シナリオ設計、implementation-scope
- `generation_notes`: Job から credential、secret store 参照、endpoint を外し、実行開始時に provider settings を再解決する状態遷移だけを候補化する。

## Shared State Invariants

- Job 側 DB は secret store 情報と endpoint を永続所有しない。
- Job Setup は provider、model、execution mode、batch mode の選択を保持し、credential / endpoint の値は保持しない。
- Ready job の実行開始は、provider settings の最新状態を再解決する。
- Running phase の開始後に使う provider settings は、開始時 snapshot 継続か更新追従かが未確定である。
- `credential_ref` は、完全削除するか監査用参照状態だけ残すかが未確定である。

## Candidate Scenarios

### CAND-PSJD-ST-001 Ready job 開始時に provider settings を再解決して Running へ遷移する

- `source requirement`: `task-frame.md:15-18`, `light-change-planning.md:10-12`, `ai-provider-settings-management.md:34-36`, `translation-job-setup.md:37-45`, `er.md:22-25`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-PSJD-ST-001`
- `actor`: 利用者
- `start state`: Ready job が存在し、Job には provider、model、execution mode、batch mode だけが保存されている。
- `event`: 利用者が翻訳 job の実行開始を要求する。
- `end state`: 実行対象 phase は Running になり、開始時に解決した provider settings の参照状態だけを runtime 側で観測できる。
- `forbidden transition`: Job に保存された credential、secret store 参照、endpoint を使って Ready から Running へ遷移してはならない。
- `expected outcome`: provider settings が最新状態で再解決され、Job 側 DB は credential / endpoint の所有者にならない。
- `observable point`: phase 開始結果、runtime snapshot または実行要約、Job 系 DB の credential / endpoint 非保持。
- `related detail requirement type`: provider settings 参照、Ready job 再解決、Job 永続境界
- `adoption hint`: Ready から Running への正常系候補として扱える。
- `conflict hint`: Running phase が開始時 snapshot を継続するか、共通設定更新へ追従するかは統合時に判断が必要である。

### CAND-PSJD-ST-002 provider settings 未設定なら Ready から Running へ遷移しない

- `source requirement`: `task-frame.md:10-18`, `light-change-planning.md:28-31`, `ai-provider-settings-management.md:27-32`, `translation-job-setup.md:40-45`, `er.md:84`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-PSJD-ST-002`
- `actor`: 利用者
- `start state`: Ready job が存在し、選択 provider の endpoint または credential 参照状態が未設定である。
- `event`: 利用者が翻訳 job の実行開始を要求する。
- `end state`: Job は Ready のまま残る、または対象 phase は Failed になる。
- `forbidden transition`: provider settings 未設定のまま Running を作成してはならない。
- `expected outcome`: 未設定理由が分類と要約で観測でき、secret 本体や raw payload は出ない。
- `observable point`: 実行開始エラー、phase 状態、provider settings の未設定状態、Job 系 DB の credential / endpoint 非保持。
- `related detail requirement type`: 未設定防止、secret 非露出、開始前再解決
- `adoption hint`: 未設定化と Ready 実行開始の境界候補として扱える。
- `conflict hint`: Ready 維持と Failed 遷移のどちらを正にするかは designer 側で人間判断候補にする。

### CAND-PSJD-ST-003 Running phase は開始時に解決した provider settings 境界を維持する

- `source requirement`: `task-frame.md:22-25`, `plan.md:28-31`, `light-change-planning.md:29-31`, `ai-provider-settings-management.md:33-36`, `er.md:63-67`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-PSJD-ST-003`
- `actor`: システム
- `start state`: phase が Running であり、開始時に provider settings を解決済みである。
- `event`: 利用者が AIサービス設定で同じ provider の endpoint または APIキー状態を更新する。
- `end state`: Running phase は開始時に解決した provider settings 境界を維持して処理を続ける。
- `forbidden transition`: Running phase の途中で更新後 provider settings へ暗黙に切り替わってはならない。
- `expected outcome`: Running phase の外部接続条件は途中で変わらず、完了または失敗の判定対象が開始時条件に固定される。
- `observable point`: Running phase の runtime snapshot、provider settings 更新後の未確定表示、phase 完了結果。
- `related detail requirement type`: Running snapshot、更新競合、状態不変条件
- `adoption hint`: 開始時 snapshot 継続を採る場合の候補である。
- `conflict hint`: 共通設定更新へ追従する案と競合する。

### CAND-PSJD-ST-004 Running phase 中の provider settings 未設定化は進行中 phase を破壊しない

- `source requirement`: `task-frame.md:22-25`, `light-change-planning.md:50-52`, `ai-provider-settings-management.md:30-36`, `er.md:63-67`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-PSJD-ST-004`
- `actor`: 利用者
- `start state`: phase が Running であり、同じ provider の provider settings は設定済みである。
- `event`: 利用者が AIサービス設定で同じ provider を未設定へ戻す。
- `end state`: provider settings row は残り、endpoint と APIキー状態は未設定になり、Running phase は定義済み方針に従って継続または Failed へ遷移する。
- `forbidden transition`: 未設定化によって Job 側へ credential / endpoint を再保存して補完してはならない。
- `expected outcome`: 未設定化は provider settings 側の状態遷移として完結し、Job 側 DB は secret store 情報を所有しない。
- `observable point`: provider settings row の残存、endpoint 未設定、APIキー状態未設定、Running phase の状態変化。
- `related detail requirement type`: 未設定化、Running 競合、Job 分離
- `adoption hint`: provider settings 未設定化と実行中 phase の競合候補として扱える。
- `conflict hint`: Running phase 継続と Failed 遷移のどちらを正にするかは未確定である。

### CAND-PSJD-ST-005 Failed phase の再実行は最新 provider settings を再解決して Running へ戻す

- `source requirement`: `task-frame.md:18`, `light-change-planning.md:28-31`, `ai-provider-settings-management.md:34-36`, `er.md:66-67`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-PSJD-ST-005`
- `actor`: 利用者
- `start state`: phase が Failed であり、同じ `JOB_PHASE_RUN` を再実行できる状態である。
- `event`: 利用者が失敗 phase の再実行を要求する。
- `end state`: 同じ `JOB_PHASE_RUN` が Running へ戻り、provider settings は再実行開始時に再解決される。
- `forbidden transition`: Failed phase の古い credential / endpoint snapshot を Job 側から再利用してはならない。
- `expected outcome`: Attempt 履歴を増やさず、同じ phase run の状態を戻す。
- `observable point`: `JOB_PHASE_RUN` の状態、再実行開始時の provider settings 解決結果、Attempt 履歴テーブル非作成。
- `related detail requirement type`: 再実行、冪等性、provider settings 再解決
- `adoption hint`: Failed から Running への再実行候補として扱える。
- `conflict hint`: 再実行時に runtime snapshot を置き換えるか、監査用に古い参照状態を残すかは統合時に判断が必要である。

### CAND-PSJD-ST-006 Completed phase は provider settings 更新や未設定化で再評価されない

- `source requirement`: `task-frame.md:10-18`, `light-change-planning.md:10-12`, `ai-provider-settings-management.md:31-36`, `er.md:22-25`, `er.md:71-77`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-PSJD-ST-006`
- `actor`: 利用者
- `start state`: phase が Completed であり、翻訳結果または出力成果物が作成済みである。
- `event`: 利用者が同じ provider の endpoint 更新、APIキー更新、または未設定化を行う。
- `end state`: Completed phase と生成済み成果物は Completed のまま残る。
- `forbidden transition`: provider settings 変更だけを理由に Completed phase を Running、Failed、Ready へ戻してはならない。
- `expected outcome`: provider settings は将来の開始または再実行で参照され、完了済み結果は暗黙に再評価されない。
- `observable point`: Completed phase 状態、生成済み成果物、provider settings 更新後状態。
- `related detail requirement type`: Completed 不変、provider settings 更新、成果物保持
- `adoption hint`: 完了済み phase と共通設定更新の分離候補として扱える。
- `conflict hint`: 完了後の手動再実行が許可される場合は、再実行候補と統合が必要である。

### CAND-PSJD-ST-007 Job 作成後に provider settings が未設定化されても Job の保存済み設定は資格情報を保持しない

- `source requirement`: `task-frame.md:15-25`, `plan.md:28-37`, `light-change-planning.md:26-32`, `translation-job-setup.md:37-45`, `ai-provider-settings-management.md:31-35`, `er.md:20-25`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-PSJD-ST-007`
- `actor`: 利用者
- `start state`: Job Setup で Ready job が作成済みであり、provider settings は設定済みである。
- `event`: 利用者が AIサービス設定で選択 provider を未設定へ戻す。
- `end state`: Ready job の provider、model、execution mode、batch mode は残り、credential / endpoint は Job 側に生成されない。
- `forbidden transition`: 未設定化の前後で Job 側へ `credential_ref`、secret store 参照、endpoint summary を補助保存してはならない。
- `expected outcome`: Ready job は実行開始時まで provider settings の未設定を抱えず、開始要求時に未設定を検出する。
- `observable point`: Ready job の保存済み phase runtime settings、provider settings 未設定状態、Job 系 DB の credential / endpoint 非保持。
- `related detail requirement type`: Job 作成後設定変更、未設定化、保存契約
- `adoption hint`: Job 作成時 snapshot を残さない場合の状態不変条件候補である。
- `conflict hint`: Job 側に provider settings revision を保持する案と競合しうる。

## Open Notes

- `human decision candidate`: Ready job の実行開始失敗を Ready 維持にするか、phase Failed にするか。
- `human decision candidate`: Running phase が provider settings 更新へ追従するか、開始時 snapshot を継続するか。
- `human decision candidate`: `credential_ref` を完全削除するか、監査用の参照状態だけ残すか。
- `merge candidate`: Ready 開始時再解決、Failed 再実行時再解決、Job 作成後未設定化は同じ provider settings 再解決契約へ統合できる。
- `rejection candidate`: Job 側 credential / endpoint snapshot を前提にする候補は、今回の分離依頼と衝突する。
