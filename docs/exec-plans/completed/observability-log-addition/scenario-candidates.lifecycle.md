# Scenario Candidates: observability-log-addition / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `OLA-LC`
- `candidate_count`: `8`

## Generator Scope

- `viewpoint`: lifecycle。作成、更新、実行、完了、再開、終了の流れで消える状態と分岐理由を候補化する。
- `included_sources`: `plan.md`, `docs/observability-logging.md`, `docs/architecture.md`, `docs/spec.md`, `docs/detail-specs/ai-provider-settings-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/translation-output-artifact.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本、`.codex/` の変更。
- `generation_notes`: 最終シナリオ表、採否、統合、競合解消は `designer` に残す。候補は観測ログの導入対象を確定しない。

## Candidate Scenarios

### CAND-OLA-LC-001 AIサービス設定の更新後に参照状態を再評価する

- `source requirement`: `docs/detail-specs/ai-provider-settings-management.md` の provider settings 更新、未設定化、接続確認状態の再評価。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-OLA-LC-001`
- `lifecycle stage`: 更新
- `actor`: AIサービス設定を更新する利用者、または後続調査をする開発者。
- `trigger`: endpoint、APIキー状態、未設定化のいずれかを保存する。
- `start condition`: provider settings row が存在し、更新前の接続確認状態または credential 状態分類がある。
- `expected outcome`: 更新後の参照状態、接続確認状態の無効化、保存結果分類を原因分離できる。
- `observable point`: `event`, `where`, `result` に加え、必要な場合だけ provider 名、credential 状態分類、接続確認状態、拒否または保存失敗の `reason` を出す。
- `reason`: endpoint 変更後の古い接続確認結果は画面上で消えるため、後続調査では更新前後の分岐理由を再構成しにくい。
- `disappearing information`: 更新前の接続確認状態、更新後に未確定へ戻した理由、未設定化で secret 本体を削除した分類、保存失敗の短い理由。
- `forbidden log`: APIキー本文、secret store key、復号可能値、raw request、raw response、raw prompt、endpoint を secret と同等に扱う箇所の詳細値。
- `related detail requirement type`: `observability_requirement`, `security_requirement`, `state_requirement`
- `adoption hint`: provider settings 更新の lifecycle を観測したい場合に採用候補にする。
- `conflict hint`: operation-audit 観点が履歴保存を求める場合、provider settings の更新履歴は保存しない仕様と衝突する可能性がある。

### CAND-OLA-LC-002 翻訳 job 作成前検証で stale 選択を拒否する

- `source requirement`: `docs/detail-specs/translation-job-setup.md` の Job Setup 作成前検証、モデル一覧鮮度 token、Ready job 作成。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-OLA-LC-002`
- `lifecycle stage`: 作成
- `actor`: 翻訳 job を作成する利用者、または後続調査をする開発者。
- `trigger`: Job Setup で翻訳 job 作成を実行する。
- `start condition`: 入力データ、共通基盤、3 つの翻訳段階の AI 設定が選択済みである。
- `expected outcome`: Ready job 作成成功、または作成拒否の分類を原因分離できる。
- `observable point`: 作成前検証の結果、拒否 `reason`、必要最小の input ID または job ID、phase 設定の count を出す。
- `reason`: stale なモデル一覧、APIキー不足、model 未選択は作成操作後に UI 状態として上書きされやすい。
- `disappearing information`: stale 判定、APIキー不足、model 未選択、既存 job summary を block 理由にしなかった判断、作成成功時の Ready job ID。
- `forbidden log`: モデル一覧鮮度 token、APIキー本文、credential 参照実値、secret store key、endpoint、provider raw data、内部ログ用識別子。
- `related detail requirement type`: `observability_requirement`, `state_requirement`, `data_requirement`, `security_requirement`
- `adoption hint`: job 作成 lifecycle の観測開始点を固定したい場合に採用候補にする。
- `conflict hint`: actor-goal 観点が UI 表示結果だけを観測点にする場合、backend 作成拒否理由との検証段階がずれる可能性がある。

### CAND-OLA-LC-003 Ready job から phase 実行を開始する

- `source requirement`: `docs/spec.md` の `Ready` から `Running`、`docs/detail-specs/term-translation-phase.md` の phase 開始条件、`docs/observability-logging.md` の状態遷移観測。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-OLA-LC-003`
- `lifecycle stage`: 実行
- `actor`: Job Run で phase を開始する利用者、または後続調査をする開発者。
- `trigger`: Ready job で単語翻訳フェーズ開始を実行する。
- `start condition`: job が Ready であり、active phase run が存在しない。
- `expected outcome`: job は Running へ進み、phase run の開始、拒否、または空完了を原因分離できる。
- `observable point`: 変更前 job state、変更後 job state、phase state、active phase run 有無、credential 状態分類、対象件数、拒否 `reason` を出す。
- `reason`: 実行開始では job state、phase state、credential 再解決、対象件数の判断が短時間で上書きされる。
- `disappearing information`: Ready 判定、terminal 判定、active phase run 既存判定、provider 未実行の空完了判定、credential 再解決の分類。
- `forbidden log`: 全 command の start / finish log、trace ID、endpoint、secret、API key、provider raw request、prompt 全文、翻訳本文全文。
- `related detail requirement type`: `observability_requirement`, `state_requirement`, `consistency_requirement`, `security_requirement`
- `adoption hint`: phase 開始時の状態遷移を最初の backend 対象にする場合に採用候補にする。
- `conflict hint`: failure 観点が provider 失敗を主対象にする場合、開始拒否と実行中失敗の境界を分ける必要がある。

### CAND-OLA-LC-004 実行中 phase の進捗と部分失敗を集約する

- `source requirement`: `docs/detail-specs/term-translation-phase.md` と `docs/detail-specs/persona-generation-phase.md` の progress、失敗分類、partial state、retryable failure。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-OLA-LC-004`
- `lifecycle stage`: 実行
- `actor`: phase 実行を監視する利用者、または後続調査をする開発者。
- `trigger`: provider 実行、辞書反映、persona 保存、progress 更新のいずれかが進む。
- `start condition`: job が Running であり、phase run が進行中である。
- `expected outcome`: 大量処理の進捗、分類、最初の失敗、最後の失敗を集約して原因分離できる。
- `observable point`: phase run ID、input count、output count、skipped count、failed count、最初と最後の error kind、retryable 分類を集約して出す。
- `reason`: loop 内の個別状態は成功分の保存や retryable failure への遷移で消え、後続調査では失敗分布を再構成しにくい。
- `disappearing information`: 共通辞書 hit、provider 実行対象件数、成功済み件数、未処理件数、一部失敗の error kind、保存途中失敗の stage。
- `forbidden log`: loop 内 1 件ごとのログ、provider raw payload、prompt 全文、原文発話全文、会話文脈全文、翻訳フィールド本文全文、XML 全文。
- `related detail requirement type`: `observability_requirement`, `recovery_requirement`, `performance_requirement`, `security_requirement`
- `adoption hint`: 大量処理の原因分離価値が高い boundary を選ぶ場合に採用候補にする。
- `conflict hint`: operation-audit 観点が監査要約を求める場合、観測ログと保存要約の責務境界を分ける必要がある。

### CAND-OLA-LC-005 Running job を中断またはキャンセルする

- `source requirement`: `docs/spec.md` の操作系状態遷移、`docs/detail-specs/translation-job-management.md` の停止、削除拒否、削除再判定。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-OLA-LC-005`
- `lifecycle stage`: 終了
- `actor`: job を停止またはキャンセルする利用者、または後続調査をする開発者。
- `trigger`: Running job の停止、Ready または Paused job のキャンセル、非実行中 job の削除を実行する。
- `start condition`: job が Ready、Running、Paused のいずれかであり、操作対象として選択されている。
- `expected outcome`: Paused、Canceled、削除成功、削除拒否のいずれになったかを原因分離できる。
- `observable point`: 操作前 job state、操作後 job state、削除可否、停止要求中かどうか、削除拒否 `reason`、入力データ保持の分類を出す。
- `reason`: 停止要求中は削除可否が変化し、停止完了後に再判定されるため、拒否理由が画面操作後に消えやすい。
- `disappearing information`: Running job の削除拒否理由、停止要求中の一時状態、停止完了後の再判定結果、非実行中 job 削除後に input data を残した分類。
- `forbidden log`: 入力 XML 全文、抽出 JSON 全文、API key、secret、provider raw response、全 command の start / finish log。
- `related detail requirement type`: `observability_requirement`, `state_requirement`, `data_requirement`, `compatibility_requirement`
- `adoption hint`: 操作系 lifecycle の安全確認を観測したい場合に採用候補にする。
- `conflict hint`: state-transition 観点が削除とキャンセルを同一遷移として扱う場合、入力保持と job 終了の意味が混ざる可能性がある。

### CAND-OLA-LC-006 Paused または RecoverableFailed job を再開する

- `source requirement`: `docs/spec.md` の再開と失敗回復、`docs/detail-specs/translation-job-management.md` の再開不可理由、`docs/detail-specs/term-translation-phase.md` と `docs/detail-specs/persona-generation-phase.md` の同じ `JOB_PHASE_RUN` 継続。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-OLA-LC-006`
- `lifecycle stage`: 再開
- `actor`: job を再開または retry する利用者、または後続調査をする開発者。
- `trigger`: Paused または RecoverableFailed の job で再開または retry を実行する。
- `start condition`: 再開対象 job が Paused または RecoverableFailed であり、再開不可理由がない。
- `expected outcome`: 同じ `JOB_PHASE_RUN` を継続し、重複作成を避けたこと、または再開拒否理由を原因分離できる。
- `observable point`: job state、phase run ID、未処理 count、retryable 分類、credential 再解決結果、再開不可 `reason` を出す。
- `reason`: 再開では最新の provider settings を再解決し、成功済み成果物を維持したまま未処理分だけ進めるため、分岐理由が永続 summary だけでは不足しやすい。
- `disappearing information`: 入力キャッシュ欠落、terminal state、状態不整合、credential 再解決失敗、成功済み成果物を重複作成しなかった判断。
- `forbidden log`: API key、credential 参照実値、secret store key、endpoint、provider raw request / response、prompt 全文、翻訳本文全文。
- `related detail requirement type`: `observability_requirement`, `recovery_requirement`, `冪等性_requirement`, `security_requirement`
- `adoption hint`: 再開と retry の冪等性を観測したい場合に採用候補にする。
- `conflict hint`: failure 観点が retryable failure の詳細を扱う場合、lifecycle 側は再開入口と同一 phase run 継続に限定する必要がある。

### CAND-OLA-LC-007 phase 完了後に後続 phase readiness を確定する

- `source requirement`: `docs/detail-specs/term-translation-phase.md` の後続 phase 入力 summary、`docs/detail-specs/persona-generation-phase.md` の body phase readiness、`docs/spec.md` の Running から Completed。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-OLA-LC-007`
- `lifecycle stage`: 完了
- `actor`: Job Run で phase 結果を確認する利用者、または後続調査をする開発者。
- `trigger`: 単語翻訳フェーズまたは NPC ペルソナ生成フェーズが Completed になる。
- `start condition`: 実行中 phase の保存と参照整合が完了している。
- `expected outcome`: phase 完了、後続 phase run 作成可否、job-level Completed 到達可否を原因分離できる。
- `observable point`: phase state、result summary count、参照成立状態、後続 phase readiness、後続作成拒否 `reason` を出す。
- `reason`: 完了直後の readiness 判定は UI では次の操作可否へ変換され、辞書参照不能や snapshot 参照不能の理由が消えやすい。
- `disappearing information`: 辞書参照成立、snapshot 参照成立、provider 未実行の空完了、後続 phase を作らなかった理由、terminal job 後書き拒否。
- `forbidden log`: raw prompt、原文発話全文、会話文脈全文、翻訳本文全文、provider raw payload、secret、API key。
- `related detail requirement type`: `observability_requirement`, `consistency_requirement`, `state_requirement`, `security_requirement`
- `adoption hint`: phase 完了から後続開始までの境界を観測したい場合に採用候補にする。
- `conflict hint`: state-transition 観点が job-level Completed だけを見る場合、phase-level readiness の候補と統合が必要になる。

### CAND-OLA-LC-008 Completed job から成果物を生成または再出力する

- `source requirement`: `docs/detail-specs/translation-output-artifact.md` の output readiness、artifact 生成、再出力、失敗 stage。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-OLA-LC-008`
- `lifecycle stage`: 終了後利用
- `actor`: 完了済み job から xTranslator 互換 XML を出力する利用者、または後続調査をする開発者。
- `trigger`: completed job を選択し、artifact 生成または再出力を実行する。
- `start condition`: body phase と job-level 状態が Completed であり、出力準備状態を評価できる。
- `expected outcome`: 成果物生成成功、再出力成功、または生成拒否と失敗 stage を原因分離できる。
- `observable point`: output readiness、row count、artifact status、failed stage、retryable flag、stale reason、再出力分類を出す。
- `reason`: Completed job は未完了一覧から消えるため、終了後利用の拒否理由や再出力理由は Job Management の lifecycle から追いにくい。
- `disappearing information`: 未完了 job 除外理由、field result 不整合、status 不整合、row validation 失敗、XML serialization 失敗、file write 失敗、artifact 保存失敗。
- `forbidden log`: Source 全文、Dest 全文、XML 全文、provider raw payload、secret、API key、復号可能値。
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `recovery_requirement`, `security_requirement`
- `adoption hint`: job 完了後の再利用 lifecycle を観測したい場合に採用候補にする。
- `conflict hint`: operation-audit 観点が artifact digest や operation kind の保存を求める場合、structured log と監査要約の保持先を分ける必要がある。

## Open Notes

- `human decision candidate`: 最初の実装対象 boundary は未決である。plan では backend 境界、frontend runtime event、UI を伴う観測確認が未決である。
- `human decision candidate`: frontend runtime event の lifecycle 候補を独立シナリオにするか、各画面 lifecycle の観測点へ含めるかは未決である。
- `human decision candidate`: `observability-logger-lightweight` を完了扱いにするか、本 task へ統合するかは未決である。
- `merge candidate`: `CAND-OLA-LC-003` と `CAND-OLA-LC-006` は、phase 開始と retry が同じ credential 再解決を使うため統合候補になりうる。
- `merge candidate`: `CAND-OLA-LC-004` と operation-audit 観点の候補は、集約 count と監査要約の境界で統合確認が必要である。
- `rejection candidate`: 全 command の start / finish log、trace ID、loop 内 1 件ごとのログを前提にした候補は、観測ログ仕様に反するため不採用候補にする。
- `conflict candidate`: provider settings の更新履歴保存を前提にした候補は、AIサービス設定管理の「更新履歴は保存しない」仕様と衝突する。
- `handoff target`: `designer`
