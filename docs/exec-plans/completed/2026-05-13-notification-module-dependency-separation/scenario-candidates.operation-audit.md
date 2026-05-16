# Scenario Candidates: 2026-05-13-notification-module-dependency-separation / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `NOTIF-AUDIT`

## Generator Scope

- `viewpoint`: 通知 module の運用確認、監査、観測ログ、履歴、再現性。
- `included_sources`: `./plan.md`, `docs/architecture.md`, `docs/observability-logging.md`, `docs/spec.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本本文、他 agent の候補成果物。
- `generation_notes`: 候補は採否未確定である。最終シナリオ表、統合、競合解消は `designer` に残す。

## Candidate Scenarios

### CAND-NOTIF-AUDIT-001 通知 dispatch 結果を後から確認できる

- `source requirement`: `./plan.md` の責務境界、`docs/architecture.md` の `4.6 Notification`、`docs/observability-logging.md` の共通 payload。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-NOTIF-AUDIT-001`
- `actor`: 運用者または実装後調査者。
- `trigger`: UseCase、Service、将来の Runner / Worker が `NotificationSinkPort` へ進捗事実、完了事実、破棄事実を渡す。
- `expected outcome`: 通知処理の結果は、`event`、`where`、`result`、必要最小の `id` または `reason` で後追い確認できる。
- `acceptance condition`: backend 観測ログから、通知事実の受領、送信、破棄、送信失敗の結果分類を区別できる。
- `acceptance condition`: 観測ログは全 command の start / finish log や loop 内の 1 件ごとの log を増やさない。
- `exclusion condition`: 通知の詳細 payload 全体、翻訳本文全文、provider raw payload を監査材料として保存する要件は含めない。
- `saved summary`: 通知種別、代表 ID、件数、結果分類、拒否または失敗理由。
- `redaction rule`: payload 全体ではなく、原因分離に必要な最小要約だけを残す。
- `observable point`: `NotificationDispatcher` 境界の backend JSON log。
- `related detail requirement type`: `observability_requirement`, `audit_requirement`
- `adoption hint`: 通知 module の運用確認を固定する候補として採用しやすい。
- `conflict hint`: 観測ログを増やしすぎる案は、`docs/observability-logging.md` の禁止事項と競合する。

### CAND-NOTIF-AUDIT-002 通知 payload に秘密値と翻訳本文全文を含めない

- `source requirement`: `./plan.md` の設計上の注意、`docs/architecture.md` の `4.5 JobIOService` と `4.6 Notification`、`docs/observability-logging.md` の禁止事項、`docs/spec.md` の APIKey 暗号化要件。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-NOTIF-AUDIT-002`
- `actor`: 運用者、実装後調査者、設計 reviewer。
- `trigger`: 通知 module が通知 payload を整形し、Wails runtime event へ送信する。
- `expected outcome`: 通知 payload、通知結果 log、再現材料には secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文、XML 全文を含めない。
- `acceptance condition`: 通知 payload と観測ログは、job id、phase run id、通知種別、件数、失敗分類などの要約だけを含む。
- `acceptance condition`: 翻訳本文は全文ではなく、必要な場合でも件数、対象種別、状態分類などの要約に留める。
- `exclusion condition`: prompt、provider 応答原文、翻訳本文全文を redaction 済みではない形で保持する確認は含めない。
- `saved summary`: 実行対象の代表 ID、通知種別、件数、結果分類。
- `redaction rule`: 秘密値、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文、XML 全文は保存禁止とする。
- `observable point`: `NotificationDispatcher` の redaction 後 payload と backend JSON log。
- `related detail requirement type`: `security_requirement`, `data_requirement`, `observability_requirement`
- `adoption hint`: 通知分離で最も重要な監査条件として採用候補にする。
- `conflict hint`: 障害再現のために本文全文を保存する案は、保存禁止条件と競合する。

### CAND-NOTIF-AUDIT-003 通知送信失敗を状態巻き戻しと分離して確認できる

- `source requirement`: `./plan.md` の設計上の注意、`docs/architecture.md` の `4.6 Notification`、`docs/spec.md` の翻訳ジョブ状態。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-NOTIF-AUDIT-003`
- `actor`: 運用者または障害調査者。
- `trigger`: `NotificationPort` または Runtime adapter で Wails runtime event の送信に失敗する。
- `expected outcome`: 通知送信失敗は観測ログで確認できるが、保存済み job / phase run 状態を成功から失敗へ巻き戻さない。
- `acceptance condition`: 通知送信失敗の log は `result` と `reason` で失敗分類を示す。
- `acceptance condition`: 同じ操作の `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` は、UseCase が確定済みの状態事実を維持する。
- `exclusion condition`: 通知送信失敗を翻訳フェーズの失敗として扱う要件は含めない。
- `saved summary`: 送信失敗した通知種別、代表 ID、失敗分類。
- `redaction rule`: adapter error は分類だけを残し、transport payload 全体と外部応答原文は残さない。
- `observable point`: `NotificationDispatcher` または `RuntimeAdapter` 境界の backend JSON log、job / phase run 状態の読み取り結果。
- `related detail requirement type`: `resilience_requirement`, `observability_requirement`, `state_requirement`
- `adoption hint`: 通知 module が状態判断の正本にならないことを検証する候補になる。
- `conflict hint`: state-transition 観点の候補と重なる可能性があるため、最終統合時は責務境界で整理する。

### CAND-NOTIF-AUDIT-004 通知結果と operation summary を DB に永続化しない

- `source requirement`: `./plan.md` の設計上の注意と禁止事項、`docs/architecture.md` の `4.5 JobIOService` と `4.6 Notification`。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-NOTIF-AUDIT-004`
- `actor`: 運用者、設計 reviewer、実装後調査者。
- `trigger`: 通知 module が通知結果、operation summary、Wails event payload を生成または送信する。
- `expected outcome`: 通知結果、operation summary、Wails event payload は DB に永続化されない。
- `acceptance condition`: 再起動後の DB には、通知結果、operation summary、Wails event payload を保持する新規永続項目が残らない。
- `acceptance condition`: 必要な再現材料は、保存済み job / phase run 状態事実と観測ログの要約から確認する。
- `exclusion condition`: 通知履歴一覧や payload 再送のための永続保存要件は含めない。
- `saved summary`: DB には保存しない。観測ログには結果分類と代表 ID だけを残す。
- `redaction rule`: 永続化しない対象に secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文を混ぜない。
- `observable point`: DB schema または repository 読み取り結果、backend JSON log。
- `related detail requirement type`: `data_requirement`, `audit_requirement`
- `adoption hint`: 通知分離が保存対象を増やさないことを固定する候補になる。
- `conflict hint`: 監査履歴を画面に表示する案が出た場合、保持期間と保存対象の人間判断が必要になる。

### CAND-NOTIF-AUDIT-005 通知送信可否の拒否理由を再現材料として残す

- `source requirement`: `./plan.md` の責務境界、`docs/architecture.md` の `4.6 Notification`、`docs/observability-logging.md` の任意 payload。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-NOTIF-AUDIT-005`
- `actor`: 運用者または障害調査者。
- `trigger`: `NotificationDispatcher` が送信可否を判定し、送信しない通知を破棄または skip する。
- `expected outcome`: 送信しない理由は、状態判断や provider response validation と混同せず、通知 module の送信可否理由として後追い確認できる。
- `acceptance condition`: 観測ログは `event`、`where`、`result`、`reason` で送信不可の分類を示す。
- `acceptance condition`: `reason` は runtime 未接続、通知種別未対応、redaction 後 payload 不足などの通知送信可否に限定する。
- `exclusion condition`: phase start 可否、terminal guard、provider response validation の拒否理由を通知 module の判断として記録する要件は含めない。
- `saved summary`: 送信しなかった通知種別、代表 ID、送信不可理由。
- `redaction rule`: 送信不可理由に payload 原文、provider raw payload、翻訳本文全文を含めない。
- `observable point`: `NotificationDispatcher` 境界の backend JSON log。
- `related detail requirement type`: `observability_requirement`, `responsibility_boundary_requirement`
- `adoption hint`: 通知 module の責務境界を運用ログで確認する候補になる。
- `conflict hint`: responsibility-boundary 観点の候補と重なる可能性がある。

### CAND-NOTIF-AUDIT-006 Wails event は push 通知専用として監査できる

- `source requirement`: `docs/architecture.md` の `3.7 Gateway と RuntimeEventAdapter`、`4.2 Controller`、`4.6 Notification`、`./plan.md` の採用方針。
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-NOTIF-AUDIT-006`
- `actor`: 運用者、frontend 調査者、設計 reviewer。
- `trigger`: Runtime adapter が `NotificationPort` から受け取った通知を Wails runtime event として送信する。
- `expected outcome`: Wails event は push 通知専用として観測でき、通常の query / command の結果正本として扱われない。
- `acceptance condition`: Controller は途中経過通知の戻り先にならず、synchronous response は Bind call の結果として残る。
- `acceptance condition`: frontend 側で runtime event を破棄する場合、画面操作後に消える破棄理由は frontend runtime log で確認できる。
- `exclusion condition`: Wails event だけで query / command の正本状態を復元する要件は含めない。
- `saved summary`: event 送信種別、代表 ID、送信結果、frontend 側の破棄理由。
- `redaction rule`: Wails event payload と frontend runtime log は、翻訳本文全文、prompt 全文、provider raw payload を含めない。
- `observable point`: Runtime adapter の backend JSON log、frontend runtime event adapter の browser console log。
- `related detail requirement type`: `integration_requirement`, `observability_requirement`
- `adoption hint`: runtime event 消費側の変更がある場合、UI 設計や browser confirmation の候補と接続できる。
- `conflict hint`: frontend 差分が無い実装範囲では、frontend runtime log の検証は採用しない可能性がある。

## Open Notes

- `human decision candidate`: 監査用の観測ログをどの粒度まで恒久化するかは、人間判断が必要である。
- `human decision candidate`: 通知履歴を UI に表示する場合、保持期間、保存対象、伏せ字範囲は別途判断が必要である。
- `merge candidate`: CAND-NOTIF-AUDIT-003 は state-transition 観点の候補と統合される可能性がある。
- `merge candidate`: CAND-NOTIF-AUDIT-005 は responsibility-boundary 観点の候補と統合される可能性がある。
- `rejection candidate`: payload 原文、prompt 全文、翻訳本文全文を再現材料として保存する候補は保存禁止条件により棄却候補である。
