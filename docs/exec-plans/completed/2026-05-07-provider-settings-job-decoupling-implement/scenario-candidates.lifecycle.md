# Scenario Candidates: 2026-05-07-provider-settings-job-decoupling-implement / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./task-frame.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `PSJD`

## Generator Scope

- `viewpoint`: lifecycle
- `included_sources`: `task-frame.md`, `../2026-05-07-provider-settings-job-decoupling/plan.md`, `../2026-05-07-provider-settings-job-decoupling/light-change-planning.md`, `docs/detail-specs/ai-provider-settings-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/er.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本本文の変更、採否判断、統合シナリオ設計、implementation-scope
- `generation_notes`: Job 作成、Ready、実行開始、Running、完了、再実行、provider settings 更新の時間順だけを候補化する。開始時 snapshot を継続するかどうかは未決論点として相反候補に分離する。

## Candidate Scenarios

### CAND-PSJD-001 Job 作成時に provider settings 所有を Job へ移さない

- `source requirement`: `task-frame.md:10`, `task-frame.md:15-18`, `translation-job-setup.md:20-22`, `translation-job-setup.md:37-45`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-PSJD-001`
- `actor`: 翻訳 job を作成したい利用者
- `trigger`: 利用者が登録済み入力データを選び、3 つの翻訳段階の provider、model、execution mode、batch mode を確認して job 作成を実行する。
- `expected outcome`: Ready job が作成される。Job 側は provider、model、execution mode、batch mode を扱い、secret store 情報と endpoint を永続所有しない。
- `observable point`: 作成後の設定内容に、翻訳段階ごとの AIサービス、model、APIキー状態、一括処理の有無だけが表示される。APIキー文字列、secret、外部サービスの raw data、内部ログ用識別子は表示されない。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `security_requirement`, `compatibility_requirement`
- `adoption hint`: Job Setup の作成完了シナリオとして採用候補にする。Job 作成時点の provider settings revision を保持するかは `CAND-PSJD-002` と競合しうる。
- `conflict hint`: Job 側に provider settings revision を保持する場合、Job が所有しない情報と監査用参照状態の境界を designer が分ける必要がある。

### CAND-PSJD-002 Ready job は実行開始前に最新 provider settings を再解決する

- `source requirement`: `task-frame.md:18`, `plan.md:28-31`, `light-change-planning.md:28-29`, `ai-provider-settings-management.md:34-36`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-PSJD-002`
- `actor`: 翻訳 job を実行開始する利用者または実行処理
- `trigger`: Ready job の実行開始操作または実行開始処理が走る。
- `expected outcome`: 実行開始前に最新 provider settings が再解決される。Job Setup 時点の endpoint または secret store 参照へ fallback しない。
- `observable point`: 実行開始前検証または実行開始ログで、Ready job が provider settings 参照状態を再取得したことを確認できる。raw secret は観測対象に含めない。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: Ready から実行開始へ進む lifecycle の基準候補にする。Job 作成時 snapshot を持たない設計と相性がよい。
- `conflict hint`: provider settings 更新履歴を保存しない仕様があるため、再解決した設定のどの要約を実行証跡に残すかは別途判断が必要である。

### CAND-PSJD-003 実行開始時に phase runtime snapshot の保存境界を固定する

- `source requirement`: `task-frame.md:25`, `plan.md:28-31`, `light-change-planning.md:31`, `ai-provider-settings-management.md:35-38`, `er.md:25`, `er.md:63-67`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-PSJD-003`
- `actor`: 翻訳フェーズ実行処理
- `trigger`: 再解決済み provider settings を使って、各翻訳フェーズの実行が開始される。
- `expected outcome`: phase runtime snapshot は、開始時の実行に必要な provider settings 由来の許可済み値だけを保持する。保存してはいけない secret、APIキー平文、復号可能値、raw request、raw response、raw prompt は保持しない。
- `observable point`: phase runtime snapshot または実行要約で、endpoint の扱い、credential 参照状態の扱い、保存禁止値の不在を確認できる。
- `related detail requirement type`: `data_requirement`, `security_requirement`, `observability_requirement`, `consistency_requirement`
- `adoption hint`: 実行開始時 snapshot の内容を切る候補にする。endpoint を snapshot に含めるかどうかは `CAND-PSJD-004` と `CAND-PSJD-005` の未決論点に接続する。
- `conflict hint`: 既存仕様は Running phase が開始時 snapshot の endpoint と credential 参照状態を使うと定義している。一方、今回の依頼は endpoint を Job 側に永続所有させない。

### CAND-PSJD-004 Running phase は開始時 snapshot を継続する

- `source requirement`: `task-frame.md:23`, `plan.md:30`, `light-change-planning.md:29`, `ai-provider-settings-management.md:36`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-PSJD-004`
- `actor`: 実行中の翻訳フェーズ処理
- `trigger`: Running phase の途中で、provider settings の endpoint または credential 参照状態が更新される。
- `expected outcome`: Running phase は開始時 snapshot を継続して使う。更新後の provider settings は、すでに Running である phase の処理結果へ途中反映されない。
- `observable point`: Running phase 中の provider 呼び出し要約または fake transport 観測で、実行開始時に確定した endpoint と credential 参照状態が継続して使われたことを確認できる。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `concurrency_requirement`, `testability_requirement`
- `adoption hint`: 既存詳細仕様を維持する候補として扱う。長時間実行中の設定揺れを避けたい場合に採用しやすい。
- `conflict hint`: endpoint と credential 参照状態を Job 側に永続所有しない方針と衝突しないように、snapshot の保存場所、保存期間、保存可能値を `CAND-PSJD-003` で切る必要がある。

### CAND-PSJD-005 Running phase は provider settings 更新へ追従する

- `source requirement`: `task-frame.md:18`, `task-frame.md:23`, `light-change-planning.md:51`, `ai-provider-settings-management.md:33-36`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-PSJD-005`
- `actor`: 実行中の翻訳フェーズ処理
- `trigger`: Running phase の途中で、provider settings の endpoint または credential 参照状態が更新される。
- `expected outcome`: Running phase は provider 呼び出し前に最新 provider settings を参照し直し、更新後の共通設定へ追従する。
- `observable point`: Running phase 中の provider 呼び出し単位で、参照した provider settings の状態要約を確認できる。更新前後の呼び出しが混在する場合は、混在した理由を要約で追える。
- `related detail requirement type`: `state_requirement`, `concurrency_requirement`, `observability_requirement`, `security_requirement`
- `adoption hint`: 開始時 snapshot 継続をやめる代替候補として扱う。provider settings を全 job provider 共通設定の正本に寄せる設計では検討対象になる。
- `conflict hint`: `CAND-PSJD-004` と相反する。既存詳細仕様の Running phase 定義を変更するため、人間判断が必要である。

### CAND-PSJD-006 provider settings 更新後は Ready job と Job Setup が最新状態を参照する

- `source requirement`: `task-frame.md:15-18`, `plan.md:12-16`, `light-change-planning.md:10-12`, `ai-provider-settings-management.md:27-35`, `translation-job-setup.md:20-22`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-PSJD-006`
- `actor`: AIサービス設定を更新する利用者、Job Setup を開く利用者
- `trigger`: 利用者が provider settings の endpoint または APIキー状態を更新する。続けて Job Setup または Ready job 実行開始が provider settings を参照する。
- `expected outcome`: provider settings row が provider ごとの共通設定正本として更新される。Job Setup と Ready job は最新の provider settings 参照状態を使い、個別の secret や endpoint へ fallback しない。
- `observable point`: endpoint 変更後に接続確認状態が未確定へ戻る。Job Setup の APIキー状態、model list 更新可否、Ready 実行開始前の再解決結果が最新状態と一致する。
- `related detail requirement type`: `data_requirement`, `state_requirement`, `consistency_requirement`, `security_requirement`, `compatibility_requirement`
- `adoption hint`: provider settings 更新から Job lifecycle へ波及する候補として採用候補にする。Job 作成済みかつ未実行の Ready job を重点観測点にする。
- `conflict hint`: provider settings 更新履歴は保存しないため、過去設定の復元や差分監査をこの候補へ混ぜない。

### CAND-PSJD-007 phase 完了後は設定由来の secret と raw payload を残さず完了状態を観測する

- `source requirement`: `ai-provider-settings-management.md:28-38`, `translation-job-setup.md:43-45`, `er.md:22-25`, `er.md:71-77`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-PSJD-007`
- `actor`: 翻訳フェーズ実行処理、完了結果を確認する利用者
- `trigger`: 各翻訳フェーズまたは翻訳 job が完了する。
- `expected outcome`: job 状態は `JOB_PHASE_RUN` 群から集約される。翻訳結果と出力ステータスは `JOB_TRANSLATION_FIELD` と出力成果物で確認できる。secret、APIキー平文、raw payload は完了要約へ残らない。
- `observable point`: 完了画面、保存要約、structured log、fake transport log で、分類と要約だけを確認できる。出力成果物は xTranslator 互換 XML として確認できる。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `observability_requirement`, `security_requirement`
- `adoption hint`: 実行完了後の lifecycle 候補にする。provider settings の保存履歴や raw payload 監査は対象外として分ける。
- `conflict hint`: 完了時に endpoint 要約を残す場合は、ローカル運用で表示可能な endpoint と、Job 側に永続所有させない情報の境界を designer が決める必要がある。

### CAND-PSJD-008 フェーズ再実行は同じ phase run を戻し、開始時に provider settings を再解決する

- `source requirement`: `task-frame.md:18`, `light-change-planning.md:28-31`, `er.md:63-67`, `ai-provider-settings-management.md:35-36`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-PSJD-008`
- `actor`: フェーズ再実行を開始する利用者または実行処理
- `trigger`: 完了済みまたは失敗済みのフェーズを再実行する。
- `expected outcome`: フェーズ再実行は同じ `JOB_PHASE_RUN` の状態を戻す扱いにする。Attempt 履歴テーブルは持たない。再実行開始時は provider settings を再解決し、再解決結果から新しい実行開始条件を満たす。
- `observable point`: 再実行後も phase run は同一単位として観測できる。再実行開始時の provider settings 参照状態と、再実行後の状態遷移を確認できる。
- `related detail requirement type`: `state_requirement`, `recovery_requirement`, `冪等性_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: 終了後利用と再開可否の候補にする。再実行で前回 snapshot を使うか、再実行開始時に最新 provider settings を使うかは designer の統合対象にする。
- `conflict hint`: `CAND-PSJD-004` を採用する場合でも、再実行開始は新しい開始点とみなすかどうかを分けて判断する必要がある。

## Open Notes

- `human decision candidate`: Job 作成時に provider settings revision を保持するか。根拠は `task-frame.md:22` と `light-change-planning.md:50` である。
- `human decision candidate`: Running phase が共通設定更新へ追従するか、開始時 snapshot を継続するか。`CAND-PSJD-004` と `CAND-PSJD-005` は相反候補である。
- `human decision candidate`: `credential_ref` を完全削除するか、監査用の参照状態だけ残すか。根拠は `task-frame.md:24` と `light-change-planning.md:52` である。
- `merge candidate`: `CAND-PSJD-002` と `CAND-PSJD-003` は、実行開始の最終シナリオで統合される可能性がある。
- `rejection candidate`: `CAND-PSJD-005` は既存詳細仕様を変更する候補であるため、人間が開始時 snapshot 継続を選ぶ場合は不採用候補になる。

## Completion Material

- `viewpoint`: lifecycle
- `candidate_count`: 8
- `artifact_path`: `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/scenario-candidates.lifecycle.md`
- `handoff_target`: designer
- `remaining_risk`: 開始時 snapshot 継続、provider settings revision 保持、credential 参照状態の保存範囲は AI だけで確定しない。
