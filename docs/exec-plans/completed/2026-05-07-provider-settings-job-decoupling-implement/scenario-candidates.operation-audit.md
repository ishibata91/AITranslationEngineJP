# Scenario Candidates: 2026-05-07-provider-settings-job-decoupling-implement / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./task-frame.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `PSJD-OA`

## Generator Scope

- `viewpoint`: `operation-audit`
- `included_sources`:
  - `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/task-frame.md`
  - `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling/plan.md`
  - `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling/light-change-planning.md`
  - `docs/detail-specs/ai-provider-settings-management.md`
  - `docs/detail-specs/translation-job-setup.md`
  - `docs/er.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本本文の変更。
- `generation_notes`: 採否、統合、最終シナリオ設計、実装範囲は扱わない。secret 本体、raw payload、復号可能値を記録する候補は主案にしない。

## Candidate Scenarios

### CAND-PSJD-OA-001 実行開始時の provider settings 再解決要約を確認する

- `source requirement`: `task-frame.md` の「実行開始時に provider settings を再解決する契約」、`ai-provider-settings-management.md` の「Ready job は実行開始前に最新 provider settings を再解決する」、`translation-job-setup.md` の「Ready job を作成できる」。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-PSJD-OA-001`
- `actor`: 運用者、利用者。
- `trigger`: Ready job の翻訳実行を開始する。
- `operation need`: 実行開始時に Job 所有の credential、endpoint ではなく、共通 provider settings を参照したことを後から確認できる。
- `audit record`: job id、phase type、provider、model、execution mode、batch mode、provider settings 再解決時刻、credential 参照状態、接続確認状態、再解決結果分類。
- `forbidden output`: APIキー本体、secret store の実参照値、復号可能値、raw request、raw response、raw prompt、Job 側に永続所有される endpoint 値。
- `expected outcome`: 実行開始要約から、各 phase が最新 provider settings を再解決した事実と結果分類を確認できる。
- `observable point`: job detail、phase run summary、保存要約、structured log のいずれかで、再解決結果の分類だけを確認できる。
- `related detail requirement type`: provider settings 参照、Ready job 再解決、保存禁止、実行要約。
- `adoption hint`: designer は実行開始時の再解決を最終シナリオへ入れるか判断する。provider settings revision を持つかどうかは人間判断候補として残す。
- `conflict hint`: provider settings の更新履歴は保存しない仕様と、再解決の再現材料保存が衝突する可能性がある。

### CAND-PSJD-OA-002 Job 作成後の保存要約から Job 側所有値の混入を確認する

- `source requirement`: `task-frame.md` の「secret store 情報と endpoint は Job 側に永続所有させない」、`translation-job-setup.md` の「作成後の設定内容には、翻訳段階ごとの AIサービス、model、APIキー状態、一括処理の有無だけを表示する」。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-PSJD-OA-002`
- `actor`: 利用者、運用者。
- `trigger`: Job Setup で Ready job を作成する。
- `operation need`: Job 作成時の保存要約が、Job 側に credential、secret store 参照、endpoint を所有させていないことを確認できる。
- `audit record`: job id、input id、phase type、provider、model、execution mode、batch mode、APIキー状態分類、作成時刻、作成結果分類。
- `forbidden output`: APIキー文字列、secret、外部サービス raw data、内部ログ用識別子、credential_ref の実値、endpoint の Job 側 snapshot。
- `expected outcome`: 作成後の設定内容と保存要約は、phase ごとの AI 設定状態だけを表示する。
- `observable point`: Job Setup の作成後表示、job summary、保存結果要約で、endpoint と secret store 参照値が出ないことを確認できる。
- `related detail requirement type`: Job Setup 保存契約、保存要約、表示禁止、DB 所有境界。
- `adoption hint`: designer は Job 作成シナリオへ、保存要約の許可項目と禁止項目を明記するか判断する。
- `conflict hint`: 既存仕様では各翻訳段階が credential 参照を持つため、表示用の APIキー状態と永続所有値の境界を統合時に分ける必要がある。

### CAND-PSJD-OA-003 Running phase の開始時要約を再現材料として残す

- `source requirement`: `task-frame.md` の「Running phase が共通設定更新へ追従するか、開始時 snapshot を継続するか」、`ai-provider-settings-management.md` の「Running phase は開始時 snapshot の endpoint と credential 参照状態を使う」、`light-change-planning.md` の不足情報。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-PSJD-OA-003`
- `actor`: 運用者、障害調査者。
- `trigger`: 各翻訳 phase が Running へ入る。
- `operation need`: Running phase が開始時の設定を継続したか、共通設定更新へ追従したかを後から説明できる。
- `audit record`: phase run id、phase type、開始時刻、provider、model、execution mode、batch mode、credential 状態分類、設定参照方針、開始時 snapshot 有無、provider settings revision または revision 未採用の分類。
- `forbidden output`: APIキー本体、復号可能値、raw payload、secret store 参照の実値、Job 側に所有される endpoint 原文。
- `expected outcome`: Running phase の開始時要約から、phase が使った設定の由来と追従方針を確認できる。
- `observable point`: phase runtime snapshot、phase run detail、実行ログ要約で、保存可能な分類値だけを確認できる。
- `related detail requirement type`: phase runtime snapshot、再現材料、設定追従方針、保存禁止。
- `adoption hint`: designer は Running phase の追従方針を確定しないまま、候補として比較できる形で扱う。
- `conflict hint`: 開始時 snapshot 継続と Job 側 endpoint 非所有の両立には、snapshot に保存する値の範囲を別途決める必要がある。

### CAND-PSJD-OA-004 Ready 作成後の provider settings 変更を実行時要約で検出する

- `source requirement`: `task-frame.md` の「Job 側に provider settings revision を保持するか」、`ai-provider-settings-management.md` の「provider settings の更新履歴は保存しない」と「Ready job は実行開始前に最新 provider settings を再解決する」。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-PSJD-OA-004`
- `actor`: 運用者、利用者。
- `trigger`: Ready job 作成後、実行開始前に provider settings の endpoint または APIキー状態が変わる。
- `operation need`: Ready 作成時の表示内容と実行開始時の再解決結果に差があり得ることを、secret を残さず説明できる。
- `audit record`: job id、phase type、Ready 作成時の provider と model、実行開始時の provider settings 解決結果分類、credential 状態分類、接続確認状態、差分あり分類。
- `forbidden output`: provider settings 更新履歴そのもの、APIキー本体、credential_ref 実値、raw request、raw response、raw prompt。
- `expected outcome`: 実行開始時の要約で、Ready 作成後に最新 provider settings を見直したことを確認できる。
- `observable point`: 実行開始要約または job detail に、差分あり分類または再解決済み分類が表示される。
- `related detail requirement type`: Ready job 再解決、更新履歴非保存、監査表示、再現材料。
- `adoption hint`: designer は revision 保存を採用するか、差分分類だけにするかを人間判断候補として残す。
- `conflict hint`: provider settings 更新履歴を保存しない仕様のため、どの差分まで監査表示するかが競合候補になる。

### CAND-PSJD-OA-005 実行開始時の provider settings 不足を失敗要約で確認する

- `source requirement`: `translation-job-setup.md` の「APIキー不足と model 未選択がない時だけ、翻訳 job を作成できる」、`ai-provider-settings-management.md` の「失敗時も secret と raw payload を露出しない」、`task-frame.md` の「実行開始時に provider settings を再解決する契約」。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-PSJD-OA-005`
- `actor`: 利用者、運用者。
- `trigger`: Ready job の実行開始時に、provider settings が未設定、endpoint 参照不能、APIキー状態不足、接続確認未確定のいずれかになる。
- `operation need`: 実行開始に失敗した理由を、分類と短い要約で後から確認できる。
- `audit record`: job id、phase type、provider、model、失敗分類、credential 状態分類、接続確認状態、失敗発生時刻、利用者向け短文要約。
- `forbidden output`: APIキー本体、復号可能値、endpoint の詳細原文、raw request、raw response、raw prompt、外部 provider 応答原文。
- `expected outcome`: 失敗要約は、設定不足または参照不能の分類だけを示し、secret や raw payload を含まない。
- `observable point`: job detail、phase error summary、structured log、画面エラー表示で、同じ失敗分類を確認できる。
- `related detail requirement type`: 失敗要約、保存禁止、実行開始再解決、provider settings 参照。
- `adoption hint`: designer は失敗シナリオと統合できるが、operation-audit 側では後追い確認の保存項目だけを候補化する。
- `conflict hint`: endpoint は AIサービス設定画面では表示可能だが、Job 側の失敗要約へどこまで出すかは別判断になる。

### CAND-PSJD-OA-006 監査表示で Job と provider settings の所有境界を確認する

- `source requirement`: `plan.md` の「Job 系 DB、Job Setup、翻訳フェーズ実行が provider settings を共通設定として参照する境界を設計し直せること」、`ai-provider-settings-management.md` の「Job Setup と master-persona は provider settings を参照し、個別の secret や endpoint を fallback にしない」。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-PSJD-OA-006`
- `actor`: 運用者、利用者。
- `trigger`: job detail または実行履歴で、作成済み job の AI 設定を確認する。
- `operation need`: Job が provider settings を参照していることと、個別の secret や endpoint へ fallback していないことを監査表示で確認できる。
- `audit record`: job id、phase type、provider、model、execution mode、batch mode、credential 状態分類、provider settings 参照済み分類、fallback 未使用分類。
- `forbidden output`: APIキー文字列、secret store 参照の実値、credential_ref 実値、raw payload、Job DB に所有された endpoint snapshot。
- `expected outcome`: 監査表示は、Job 側に保存された AI 設定要約と provider settings 参照結果を分けて表示する。
- `observable point`: job detail、phase run detail、監査用 summary で、Job 所有値と共通設定参照値の境界が読める。
- `related detail requirement type`: 監査表示、所有境界、provider settings 参照、保存禁止。
- `adoption hint`: designer は画面表示シナリオまたは運用確認シナリオに統合するか判断する。
- `conflict hint`: Job Setup の作成後表示と Running phase の実行時要約は、同じ値に見えても所有者が異なるため、表示名の混同が競合候補になる。

### CAND-PSJD-OA-007 フェーズ再実行時の履歴粒度を監査候補として残す

- `source requirement`: `er.md` の「フェーズ再実行は同じ JOB_PHASE_RUN の状態を戻す扱い」と「Attempt 履歴テーブルは持たない」、`task-frame.md` の「phase runtime snapshot が保存してよい値と保存してはいけない値を分ける」。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-PSJD-OA-007`
- `actor`: 運用者、障害調査者。
- `trigger`: 失敗または中断した phase を再実行する。
- `operation need`: Attempt 履歴テーブルを持たない前提でも、最新の再実行結果と設定要約を後から確認できる。
- `audit record`: phase run id、phase type、再実行時刻、前状態、後状態、provider、model、execution mode、batch mode、credential 状態分類、再解決結果分類、直近失敗分類。
- `forbidden output`: 過去 attempt の raw payload、APIキー本体、復号可能値、secret store 参照の実値、endpoint 原文の履歴化。
- `expected outcome`: 同じ phase run の最新状態として、再実行後の監査要約を確認できる。
- `observable point`: phase run detail、job status summary、直近エラー要約で、最新結果と直近失敗分類を確認できる。
- `related detail requirement type`: 履歴粒度、再実行、phase run 状態、再現材料。
- `adoption hint`: designer は attempt 履歴なしで足りるか、別の監査粒度が必要かを人間判断候補として残す。
- `conflict hint`: 履歴を増やす判断は ER の Attempt 履歴テーブルなし方針と衝突する可能性がある。

## Open Notes

- `human decision candidate`: provider settings revision を Job または phase に保持するか。
- `human decision candidate`: Running phase は共通設定更新へ追従するか、開始時要約を継続するか。
- `human decision candidate`: `credential_ref` を完全削除するか、監査用の状態分類だけ残すか。
- `human decision candidate`: endpoint を Job 側の要約、失敗要約、監査表示へ出す範囲をどこまで許可するか。
- `merge candidate`: `CAND-PSJD-OA-001` と `CAND-PSJD-OA-004` は、実行開始時再解決のシナリオへ統合できる可能性がある。
- `merge candidate`: `CAND-PSJD-OA-002` と `CAND-PSJD-OA-006` は、Job 作成後表示と監査表示の境界として統合できる可能性がある。
- `rejection candidate`: provider settings の更新履歴保存を主目的にする候補は、既存詳細仕様の対象外と衝突する。
