# Scenario Candidates: observability-log-addition / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `OBSLOG-OA`

## Generator Scope

- `viewpoint`: operation-audit
- `included_sources`: `./plan.md`, `docs/observability-logging.md`, `docs/architecture.md`, `docs/spec.md`, `docs/er.md`, `docs/detail-specs/translation-input-intake.md`, `docs/detail-specs/master-dictionary.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`, `docs/detail-specs/translation-output-artifact.md`, `docs/detail-specs/ai-provider-settings-management.md`, `docs/scenario-tests/README.md`, `docs/screen-design/README.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本変更、`.codex/` 変更、他 agent の候補成果物
- `generation_notes`: 採否、統合、競合解消は `designer` に残す。候補は、運用確認、監査、再現、禁止ログ確認に必要な観測ログだけを扱う。

## Candidate Scenarios

### CAND-OBSLOG-OA-001 状態変更と拒否理由を後から確認する

- `source requirement`: `./plan.md` の原因分離ログ追加方針、`docs/observability-logging.md` の状態遷移ログ対象、`docs/architecture.md` の `StateMachine` / `JobIOService` 境界
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-OBSLOG-OA-001`
- `actor`: 運用確認者、障害調査者
- `trigger`: job または phase の状態変更、状態変更拒否、危険操作の無効化が起きる。
- `expected outcome`: 変更前状態、変更後状態、結果分類、拒否理由、代表 ID を後から確認できる。成功、拒否、未変更を区別できる。
- `observable point`: backend の `Backend UseCase`、`StateMachine`、`JobIOService` のうち、変更前状態、変更後状態、拒否理由が同じ場所で取れる境界
- `saved summary`: `event`、`where`、`result`、必要最小の `id`、`reason`
- `redaction rule`: job ID や phase run ID は代表 ID として残してよい。入力本文、provider 応答原文、credential 参照実値は残さない。
- `disappearing information`: 変更前状態、拒否時の判定理由、状態を変えなかった事実、集約不能を安全側に倒した理由
- `forbidden log`: 全 command の start / finish log、DTO 全体、translation field 全文、secret、API key、endpoint
- `related detail requirement type`: `observability_requirement`, `state_requirement`, `consistency_requirement`, `security_requirement`
- `adoption hint`: 状態遷移の不整合調査を受け入れ条件に含める場合に候補にする。
- `conflict hint`: 全 command start / finish log を求める案と衝突する。状態変更ログが履歴保存や監査テーブル作成に拡張される場合は、今回のログ仕様を超える。
- `open issue`: どの状態変更境界を最初の実装対象にするかは未決である。

### CAND-OBSLOG-OA-002 provider、file、DB、Wails 境界の失敗分類を再現材料にする

- `source requirement`: `docs/observability-logging.md` の外部境界失敗分類、`docs/architecture.md` の adapter 境界、各 phase 詳細仕様の provider 失敗と保存失敗の扱い
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-OBSLOG-OA-002`
- `actor`: 障害調査者、運用確認者
- `trigger`: provider、file、DB、Wails 境界で失敗、拒否、応答不正、保存失敗、参照不能が起きる。
- `expected outcome`: 失敗した境界、結果分類、短い理由、再試行判断に必要な分類を後から確認できる。raw payload を見なくても原因候補を分離できる。
- `observable point`: backend adapter、backend controller、frontend gateway、Wails runtime event adapter のうち、失敗分類が確定する境界
- `saved summary`: `event`、`where`、`result`、必要最小の `id`、`reason`
- `redaction rule`: provider 名、model 名、credential 状態分類は要約として残してよい。endpoint、API key、secret store key、provider raw request / response は残さない。
- `disappearing information`: adapter で正規化された error kind、retryable 判定、参照不能理由、Wails 境界の失敗箇所
- `forbidden log`: raw request、raw response、raw prompt、fake transport log への secret 出力、frontend log の backend 転送
- `related detail requirement type`: `observability_requirement`, `failure_handling_requirement`, `recovery_requirement`, `security_requirement`
- `adoption hint`: provider 失敗と保存失敗を成功扱いにしない検証と結びつける。
- `conflict hint`: 実 API の詳細 payload を証跡に残す案は、保護仕様と衝突する。
- `open issue`: provider、file、DB、Wails のうち、最初に対象にする境界は未決である。

### CAND-OBSLOG-OA-003 翻訳 phase の監査要約を redaction 付きで確認する

- `source requirement`: `docs/detail-specs/term-translation-phase.md`、`docs/detail-specs/persona-generation-phase.md`、`docs/detail-specs/body-translation-phase.md` の監査要約と保護仕様
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-OBSLOG-OA-003`
- `actor`: 運用確認者、障害調査者
- `trigger`: 単語翻訳、NPC ペルソナ生成、本文翻訳の phase 開始、完了、失敗、retry、resume が起きる。
- `expected outcome`: phase ごとに provider、model、execution mode、batch mode、credential 状態分類、input count、output count、snapshot digest、prompt digest、error kind を確認できる。
- `observable point`: backend usecase と service のうち、phase run、入力 summary、provider 実行 summary、保存結果が集約済みになる境界
- `saved summary`: `event`、`where`、`result`、`id`、`count`、`reason`、phase 固有の digest または version
- `redaction rule`: prompt は全文ではなく digest または version だけにする。原文発話全文、会話文脈全文、翻訳本文全文、raw prompt は残さない。
- `disappearing information`: phase 開始時 snapshot、provider 未実行理由、partial failure の件数、retry 対象の分類、後続 phase 不可理由
- `forbidden log`: prompt 全文、翻訳本文全文、原文発話全文、会話文脈全文、provider raw payload、credential 参照実値
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `recovery_requirement`, `security_requirement`
- `adoption hint`: phase の監査要約と secret 非露出を同時に検証する候補にする。
- `conflict hint`: phase ごとの詳細ログを永続履歴として保存する案は、ER の attempt 履歴なし方針と衝突する可能性がある。
- `open issue`: digest の対象範囲と保持粒度は、人間判断が必要である可能性がある。

### CAND-OBSLOG-OA-004 frontend runtime event の破棄理由を画面操作後に確認する

- `source requirement`: `docs/observability-logging.md` の frontend runtime event 破棄理由、`docs/architecture.md` の `RuntimeEventAdapter` 境界
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-OBSLOG-OA-004`
- `actor`: frontend 調査者、運用確認者
- `trigger`: Wails runtime event が stale event、対象画面不一致、購読解除後到着、再読込不要などの理由で破棄される。
- `expected outcome`: browser console 上で破棄理由、対象境界、代表 ID、結果分類を確認できる。画面状態が変わらなかった理由を後から説明できる。
- `observable point`: frontend の `RuntimeEventAdapter`、`ScreenController`、`Frontend UseCase` のうち、event 採用可否を判定する境界
- `saved summary`: `event`、`where`、`result`、必要最小の `id`、`reason`
- `redaction rule`: frontend log は browser console にだけ出す。backend へ転送しない。runtime event payload 全体を残さない。
- `disappearing information`: stale 判定、画面 unmount 後到着、再読込不要判定、store 更新をスキップした理由
- `forbidden log`: frontend から backend へのログ送信、runtime event payload 全体、DTO 全体、翻訳本文全文
- `related detail requirement type`: `observability_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: UI 操作後に状態更新が起きない理由を確認する受け入れ条件に使う。
- `conflict hint`: backend log と frontend log を同一 file に集約する案と衝突する。
- `open issue`: UI を伴う観測確認が必要かは `./plan.md` で未決である。

### CAND-OBSLOG-OA-005 大量処理の集約結果だけを監査材料にする

- `source requirement`: `docs/observability-logging.md` の大量処理ログ対象、`docs/detail-specs/translation-input-intake.md`、`docs/detail-specs/master-dictionary.md`、`docs/detail-specs/body-translation-phase.md`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-OBSLOG-OA-005`
- `actor`: 障害調査者、運用確認者
- `trigger`: XML import、辞書抽出、本文 field 処理、provider request unit 処理など、複数件をまとめて処理する。
- `expected outcome`: 件数、分類、最初の失敗分類、最後の失敗分類、成功件数、失敗件数を集約済みログで確認できる。1 件ごとの loop log は出さない。
- `observable point`: service または adapter のうち、処理件数と分類が集約済みになる境界
- `saved summary`: `event`、`where`、`result`、`count`、必要なら `reason`
- `redaction rule`: record type、status 分布、件数は残してよい。XML 全文、翻訳本文全文、個別 field payload は残さない。
- `disappearing information`: 抽出対象外理由の分布、未定義組み合わせの件数、保存前 validation の失敗分布、最初と最後の失敗分類
- `forbidden log`: loop 内 1 件ごとのログ、XML 全文、Source / Dest 全文、provider raw payload
- `related detail requirement type`: `observability_requirement`, `boundary_requirement`, `data_requirement`, `performance_requirement`, `security_requirement`
- `adoption hint`: 大量処理の原因分離を、過剰ログなしで確認する候補にする。
- `conflict hint`: 個別 row 監査を求める案は、loop 内 1 件ごとのログ禁止と衝突する。
- `open issue`: 最初の失敗と最後の失敗に含める安全な要約項目は未決である。

### CAND-OBSLOG-OA-006 翻訳成果物出力と再出力の操作要約を監査する

- `source requirement`: `docs/detail-specs/translation-output-artifact.md` の監査要約、`docs/observability-logging.md` の file 境界失敗分類
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-OBSLOG-OA-006`
- `actor`: 利用者、運用確認者、障害調査者
- `trigger`: xTranslator 互換 XML の出力、再出力、row validation 失敗、XML serialization 失敗、file write 失敗、artifact 保存失敗が起きる。
- `expected outcome`: artifact id、operation kind、row count、digest、error kind、status を中心に、出力処理の成否と失敗 stage を確認できる。
- `observable point`: output artifact usecase、XML serialization adapter、file write adapter、repository 保存境界
- `saved summary`: `event`、`where`、`result`、`id`、`count`、`reason`、digest
- `redaction rule`: file path は UI 表示要件と衝突しない範囲で扱う。Source / Dest の全文、XML 全文、過剰な本文全文は残さない。
- `disappearing information`: readiness 拒否理由、row validation の失敗分類、serialization 失敗 stage、file write 失敗分類、再出力で更新または置換した事実
- `forbidden log`: XML 全文、Source / Dest 全文、provider raw payload、secret、API key
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `failure_handling_requirement`, `recovery_requirement`, `security_requirement`
- `adoption hint`: 成果物出力が provider、network、secret を使わないことを観測確認する候補にする。
- `conflict hint`: 複数世代の artifact revision 履歴を正本化する案は、対象外仕様と衝突する。
- `open issue`: file path をログへ残す範囲は、人間判断が必要である可能性がある。

## Open Notes

- `human decision candidate`: 最初に対象にする backend 境界、frontend runtime event、UI を伴う観測確認の有無は未決である。
- `human decision candidate`: digest の対象範囲、file path の扱い、最初の失敗と最後の失敗に含める要約項目は未決である。
- `merge candidate`: 状態変更ログと phase 監査要約は、同じ phase start / retry / resume 操作で統合候補になる。
- `merge candidate`: 外部境界失敗分類と大量処理集約は、file import と artifact 出力で統合候補になる。
- `rejection candidate`: trace ID、全 command start / finish log、loop 内 1 件ごとのログ、backend への frontend log 転送を前提にする候補は不採用候補である。
- `conflict candidate`: 保存対象が secret、API key、credential 参照実値、endpoint、provider raw payload、prompt 全文、翻訳本文全文、XML 全文を含む場合は `security_requirement` / `data_requirement` と衝突する。
- `candidate count`: 6
