# translation-job-list-fix

## 依頼要約

翻訳管理のジョブ一覧の修正と UX 改善を行う。論点は次の三点である。

1. 削除不可バグ
   - 症状: ジョブを作成した直後から削除できない。
   - 表示メッセージ: 「削除: 状態不整合があるため、一覧を再読込してから削除してください。」
   - 期待: 作成直後でも削除できる状態を返す。

2. 再開ボタンと翻訳段階へ進むボタンの役割重複
   - 観察: 「再開」ボタンと「翻訳段階へ進む」ボタンが並ぶ。
   - 評価: 「翻訳段階へ進む」だけで再開動線は成り立つ。
   - 方針: どちらか一方だけを残す。残す側の文言を「再開」とする案を後続モジュールで判断する。

3. 選択不可理由の表示重複
   - 観察: ホバー時の tooltip に選択不可理由が出るのに、ボタン直下にも同じ理由が羅列される。
   - 方針: ボタン直下の選択不可理由の表示を削除する。

## 分岐元

- 分岐元 branch: master
- 分岐元 commit: c2643f876e466123894287dd635b85915a03b214
- 作業 branch: claude/translation-job-list-fix

## 人間観測記録（論点1）

- 観測者: 人間（task 入力提出者）。
- 観測対象画面: 翻訳管理ジョブ一覧（`TranslationJobManagementScreen` 系）。
- 観測手順:
  - 翻訳ジョブを新規作成する。
  - 一覧に表示された直後の対象ジョブで削除ボタンを操作しようとする。
- 観測結果:
  - 削除ボタンが押せない状態で表示される。
  - 選択不可理由として「削除: 状態不整合があるため、一覧を再読込してから削除してください。」が表示される。
- 期待との差分:
  - 仕様 REQ-002「`Ready` の翻訳ジョブは、フェーズ開始前には実行中の翻訳段階を持たない。」と REQ-004「非実行中の翻訳ジョブを削除しても、入力データと抽出 JSON 正本は残る。」に基づき、作成直後の Ready ジョブは削除可能である必要がある。
- 観測仮説の取扱: 本セクションは観測事実のみで、原因仮説は `修正方針判断` で固定する。

## 想定 Y/N 評価（論点1 削除不可バグ対象、investigation-module 入口）

- 仕様変更または仕様追加がある: N。`docs/detail-specs/translation-job-management.md` REQ-002 仕様「`Ready` の翻訳ジョブは、フェーズ開始前には実行中の翻訳段階を持たない。」と REQ-004 仕様「非実行中の翻訳ジョブを削除しても、入力データと抽出 JSON 正本は残る。」で「作成直後の Ready job は削除可」が既に固定済み。
- 画面変更がある: N。仕様で固定済みの削除可否を表示するだけで、画面構造は変えない。
- 内部構造変更がある: Y。`internal/service/translation_job_management_service.go` の warnings 生成または `buildTranslationJobManagementDeleteAvailability` の判断ロジックを修正する想定。
- 画面の表示変更がある: N。`storybook-module` の TRIGGER 軸に該当しない。
- frontend ロジック変更がある: N。frontend は backend から返る削除可否表示を読むだけ。仮説検証次第で N から Y へ変わる可能性がある。
- backend 変更がある: Y。`internal/service/translation_job_management_service.go` の判断ロジック修正想定。
- frontend と backend を接続する: N。既存契約のまま。
- 実装済み責務を独立に証明したい: Y。Ready 直後の削除可否、warning 分類の境界を単体テストで証明する。
- 実行時にしか確定しない値または原因分離が要る分岐がある: Y。観測ログで「freshly created Ready job の warnings 内容」を確定させる。

判定: 仕様変更 N、画面変更 N、内部構造変更 Y、backend 変更 Y。investigation-module を継続する。

## 修正方針判断（論点1 削除不可バグ）

### 判断結果

- 判定: 完了

### 観測済み問題

- 問題: Ready 状態の翻訳ジョブを作成した直後に削除ボタンが無効になり、「削除: 状態不整合があるため、一覧を再読込してから削除してください。」が表示される。
- 期待との差分: 仕様 REQ-002 と REQ-004 に基づき、Ready ジョブはフェーズ開始前であり削除可能でなければならない。

### 画面再現確認

- Wails 接続対象: http://localhost:34115
- 再現手順: ダッシュボードから翻訳管理画面を開き、一覧に表示されたジョブ #1（状態: 実行前）の削除ボタンを確認する。
- 操作結果: 削除ボタンが `disabled` 状態で表示される。
- 画面状態: ジョブ #1 および #2 の両方で「削除: 状態不整合があるため、一覧を再読込してから削除してください。」が表示される。「再開: runtime snapshot が存在しないため、再開可否を安全側で評価します。」も同時に表示される。
- 証跡: a11y snapshot（uid=3_36 button "削除" disableable disabled、uid=3_39 StaticText 削除不可理由）

### 原因仮説と観測ログ検証

仮説の出発点として 4 経路を立てた。

1. `buildTranslationJobManagementRuntimeSummary`（snapshot 空）が `state_projection_inconsistent` を返す経路
2. `buildInputSourceSummary`（入力参照 NotFound）が `state_projection_inconsistent` を返す経路
3. `buildCacheState`（cacheReader 未設定）が `state_projection_inconsistent` を返す経路
4. `buildTranslationJobManagementProgress`（Ready 以外で phase run 無し）が `state_projection_inconsistent` を返す経路

観測ログを 2 箇所に追加し、一覧表示時に backend が返す warnings の内容を観測した。

観測結果（job #1、job #2 の両方で同一）：

```json
{"event":"runtime_summary_no_snapshots","where":"backend.service.translation_job_management","result":"warning_added","warning_category":"state_projection_inconsistent"}
{"event":"delete_availability_pre_check","job_state":"ready","snapshot_count":0,"warning_count":1,"warning_categories":["state_projection_inconsistent"]}
```

- 仮説1（`buildTranslationJobManagementRuntimeSummary` 経路）: 支持。snapshot_count が 0 のため空 snapshots 分岐が実行され、`state_projection_inconsistent` warning が 1 件追加された。
- 仮説2（`buildInputSourceSummary` 経路）: 否定。warning_count が 1 のため、入力参照 NotFound 経路の追加はなかった。
- 仮説3（`buildCacheState` 経路）: 否定。warning_count が 1 のため、cacheReader 未設定経路の追加はなかった。
- 仮説4（`buildTranslationJobManagementProgress` 経路）: 否定。job.State が ready のため `else if` 条件を満たさず、warning は追加されなかった。

追加した一時観測ログは検証完了後に削除し、ビルドが通ることを確認済み。

### 確定原因

- 原因: `buildTranslationJobManagementRuntimeSummary` は snapshots が空の場合に `state_projection_inconsistent` category の warning を返す。この warning は本来「snapshot が無い時は再開不可（外部 API 設定を確認できないため再開を許可しない）」を表すための再開判断専用の信号である。しかし `buildTranslationJobManagementDeleteAvailability` は warnings スライス全体をループして `state_projection_inconsistent` が 1 件でも存在すれば削除をブロックする。Ready job は一度もフェーズを開始していないため runtime snapshot が存在しないことが正常であり、削除は仕様 REQ-004 により可能である必要がある。しかし削除判断が再開専用の warning を区別せずブロック条件として使うため、Ready 直後の job が削除不可になる。
- 観測根拠: 観測ログで `warning_count: 1`、`warning_categories: ["state_projection_inconsistent"]`、`snapshot_count: 0`、`job_state: "ready"` を確認。唯一の warning は `runtime_summary_no_snapshots` イベントから発生。

### 採用する修正方針

- 方針A: `buildTranslationJobManagementDeleteAvailability` 内の warnings ループで、`state_projection_inconsistent` warning が `buildTranslationJobManagementRuntimeSummary` 由来（すなわち snapshot 欠落）であり、かつ job state が Ready の場合は削除をブロックしない。
  - 実現方法: warnings のループ判断に「job.State が ready かつ warning が snapshot 欠落起因の場合はスキップ」の条件を追加する。ただしこの方法は warnings に発生元情報がなく、発生元を判別するために snapshot_count や job state の追加引数が必要になる。
- 方針B（採用）: `buildTranslationJobManagementRuntimeSummary` が返す warning の category を `state_projection_inconsistent` から別の category（例: `runtime_snapshot_missing`）へ変更し、`buildTranslationJobManagementDeleteAvailability` のブロック判断は既存の `state_projection_inconsistent` のみを対象とする。`buildTranslationJobManagementResumeBlockedReasons` は `runtime_snapshot_missing` を再開ブロック理由として扱う。
  - 理由: warning の category が表す意味を発生元の文脈に合わせて分離する。削除判断と再開判断それぞれが必要な category だけを判断対象にする。既存の warning 構造を大きく変えず、category 定数の追加と参照箇所の修正で完結する。
- 方針C: `buildTranslationJobManagementRuntimeSummary` から warnings 返却を廃止し、空 snapshots を正常状態として扱い、表示上の「未設定」だけを返す。再開判断側で snapshots 空を独立に判断する。
  - 理由: `buildTranslationJobManagementRuntimeSummary` の warnings が再開文脈専用であることを明示できる。ただし再開判断ロジックの変更範囲が広くなる。

採用: 方針B。`state_projection_inconsistent` と `runtime_snapshot_missing` を分離することで、削除ブロック判断と再開ブロック判断が参照する category を明確に分ける。変更範囲は `buildTranslationJobManagementRuntimeSummary` の返す category 定数の変更、`buildTranslationJobManagementDeleteAvailability` のループ条件への影響確認（`state_projection_inconsistent` のみを対象とすることを維持）、`buildTranslationJobManagementResumeBlockedReasons` に `runtime_snapshot_missing` を追加する修正に限定できる。

- 理由: Ready job の runtime snapshot 欠落は正常状態であり `state_projection_inconsistent`（状態不整合）という category は誤った分類である。「snapshot 欠落は再開不可、削除は可能」という現象として正しい振る舞いを、category 分離で表現する。

補足: 既存メッセージ「runtime snapshot が存在しないため、再開可否を安全側で評価します。」「保存済み AI 設定要約が不足しています」も曖昧で利用者に意味が伝わらないため、implementation-module では「snapshot が無いため再開できません」「外部 API 設定の確認に必要な snapshot が無いため再開できません」など、結果として何が起きるかを直接示す表現へ書き換える。

### 禁止する修正

- 禁止1: `buildTranslationJobManagementRuntimeSummary` からの warnings 返却を削除し、削除判断の warning ループを全廃する。理由: 入力参照 NotFound や cacheReader 未設定などの他の `state_projection_inconsistent` 経路は、削除ブロックとして正当な理由となりうる。warnings ループを全廃すると、それらの正当なブロックも解除される。
- 禁止2: `buildTranslationJobManagementDeleteAvailability` に `job.State == ready` の特別分岐を追加して削除を強制許可する。理由: 表面症状だけを隠す対症療法であり、warnings の category が誤分類されている根本問題を残す。将来別の経路で同様の問題が再発しうる。
- 禁止3: 新しい job 状態値を追加して Ready 直後を別状態として扱う。理由: 既存の状態モデルで説明できる問題であり、状態値の追加は over-engineering となる。

### 影響ファイル候補

- `internal/service/translation_job_management_service.go`: `buildTranslationJobManagementRuntimeSummary`（category 変更）、`buildTranslationJobManagementResumeBlockedReasons`（`runtime_snapshot_missing` 追加）、定数定義（`translationJobManagementReasonRuntimeSnapshotMissing` 追加）
- 上記ファイルの単体テスト（存在する場合）: Ready job の削除可否と再開ブロック理由の境界を証明するテストを追加する。

## UC 差分候補（論点1）

### 判定

- 分類: 差分なし
- 結論: ユースケース正本 `docs/usecases/uc-translation-management.md` と詳細仕様正本 `docs/detail-specs/translation-job-management.md` REQ-002, REQ-004 で、論点1 の利用者観察可能な振る舞い（Ready 直後の削除可、Ready 直後の再開は snapshot 欠落を安全側で評価する）はすべて説明できる。
- design-module 迂回判断: 不要。本モジュール（investigation-module）の中で続行する。

### 判定根拠

| 観察可能な振る舞い | 既存正本の根拠 |
| --- | --- |
| Ready 直後の翻訳ジョブは削除可能である | REQ-002「`Ready` の翻訳ジョブは、フェーズ開始前には実行中の翻訳段階を持たない。」、REQ-004「非実行中の翻訳ジョブを削除しても、入力データと抽出 JSON 正本は残る。」 |
| 削除不可状態では無効理由を表示し削除を拒否する | UC「ジョブ情報を削除する」E1、REQ-004「実行不可理由を理由分類として判断できる」 |
| 翻訳ジョブ状態と現在フェーズ状態の不整合では危険操作を拒否する | REQ-004「翻訳ジョブ状態と現在フェーズ状態が食い違う場合、状態不整合として扱い、危険操作は拒否結果にする。」 |
| Ready 直後の再開は snapshot 欠落を安全側で評価する余地がある | UC「ジョブを再開する」E1「再開可能ではない、AI 設定が不足している、または送信中である」と REQ-004「実行不可理由を理由分類として判断できる」 |

### 判定の補足

採用方針B（`state_projection_inconsistent` と `runtime_snapshot_missing` の category 分離）は、削除と再開の拒否理由分類を内部で分離する実装変更である。category 名は内部識別子で、利用者向けには既存の「無効理由」表示として観察される。REQ-004 が「実行不可理由を理由分類として判断できる」と複数分類を許容しているため、UC 正本と詳細仕様正本の本文追加は不要と判断する。

## E2E テスト観点差分（論点1）

### 判定

- 分類: 追加候補あり
- 追記先想定: `docs/e2e-test-design/test-design.csv`（観点 ID は採番未確定）

### 既存観点の関連行（参考）

| 既存 ID | 関連UC | 既存観点が扱う範囲 |
| --- | --- | --- |
| E2E-UC-020 | ジョブ情報を削除する | Canceled ジョブの削除有効状態と削除確定の正常系 |
| E2E-UC-037 | ジョブを停止する | Ready 状態で停止無効である境界 |
| E2E-UC-038 | ジョブを再開する | 再開不可の失敗ジョブで無効理由を維持する例外 |
| E2E-UC-039 | ジョブ情報を削除する | 削除確認モーダルでの削除取り消し代替 |

### 追加候補テスト観点表

| 仮ID | 関連UC | 対象画面 | 前提条件 | 手順 | 期待値 | 備考 | 分類 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| E2E-UC-TJL-FIX-1 | ジョブ情報を削除する | 未完了ジョブ一覧 | 画面表示: `[data-testid=translation-job-management-job-card]` に作成直後の Ready 状態ジョブが表示されている。画面表示: 当該ジョブの `[data-testid=translation-job-management-state-label]` は 実行前 を表示している。runtime snapshot は未生成である。入力参照と cacheReader の参照は成功する。 | `[data-testid=translation-job-management-job-card]` の text=Ready を含む card 内の `[data-testid=translation-job-management-job-actions]` 内の `button[name=削除]` の enabled を確認する; click; `[data-testid=translation-job-management-delete-confirmation-modal]` 内の `button[name=ジョブ情報だけを削除する]` を click | `button[name=削除]` が enabled である; `[data-testid=translation-job-management-feedback-notification]` が削除結果を表示する; `[data-testid=translation-job-management-job-card]` が削除対象を表示しない | 境界: Ready 直後（snapshot 欠落のみ）で削除可。修正の主検証点。削除ボタン単体 selector は未決 | 境界 |
| E2E-UC-TJL-FIX-2 | ジョブ情報を削除する | 未完了ジョブ一覧 | 画面表示: `[data-testid=translation-job-management-job-card]` に Ready 状態のジョブが表示されている。当該ジョブの入力参照は NotFound である（真の状態不整合）。 | `[data-testid=translation-job-management-job-actions]` 内の `button[name=削除]` の disabled を確認する; click は実行しない | `button[name=削除]` が disabled である; `[data-testid=translation-job-management-disabled-reason]` が状態不整合由来の削除不可理由を表示する; 削除確認モーダルは表示されない | 例外: 真の `state_projection_inconsistent` warning では Ready でも削除拒否を維持する。disabled-reason 子 selector は未決 | 例外 |
| E2E-UC-TJL-FIX-3 | ジョブを再開する | 未完了ジョブ一覧 | 画面表示: `[data-testid=translation-job-management-job-card]` に作成直後の Ready 状態ジョブが表示されている。runtime snapshot は未生成である。AI 設定は満たしている。 | `[data-testid=translation-job-management-job-actions]` 内の `button[name=再開]` の disabled を確認する; click は実行しない | `button[name=再開]` が disabled である; `[data-testid=translation-job-management-disabled-reason]` が runtime snapshot 欠落由来の再開不可理由（安全側評価）を表示する; `[data-testid=translation-job-management-state-label]` が 実行前 を維持する | 境界: `runtime_snapshot_missing` を再開ブロック理由として引き続き考慮することを証明する。disabled-reason 子 selector は未決 | 境界 |

### data-testid 必要 selector

- `[data-testid=translation-job-management-job-card]`（既存）
- `[data-testid=translation-job-management-job-actions]`（既存）
- `[data-testid=translation-job-management-state-label]`（既存）
- `[data-testid=translation-job-management-disabled-reason]`（既存。子要素として「削除不可理由」「再開不可理由」を区別する selector が必要。区別 selector は未決のため `data-testid-gaps.md` に記録する。）
- `[data-testid=translation-job-management-delete-confirmation-modal]`（既存）
- `[data-testid=translation-job-management-feedback-notification]`（既存）
- 削除ボタン単体 selector は未決（既存観点と同じ未決事項）

### 観点の網羅

- 正常: 既存 E2E-UC-020 で網羅済み（Canceled 削除）。
- 代替: 既存 E2E-UC-039 で網羅済み（削除取り消し）。
- 例外: 追加候補 E2E-UC-TJL-FIX-2 が「真の `state_projection_inconsistent` warning では削除拒否を維持」を証明する。
- 境界: 追加候補 E2E-UC-TJL-FIX-1 が「Ready 直後の `runtime_snapshot_missing` のみで削除可」、E2E-UC-TJL-FIX-3 が「Ready 直後の `runtime_snapshot_missing` で再開拒否を維持」を証明する。

## 人間修正レビュー（論点1）

- 承認日: 2026-05-31
- 承認者: 人間（task 入力提出者）
- 採用方針: 方針B（`state_projection_inconsistent` と `runtime_snapshot_missing` の category 分離）
- 表現指示: 既存メッセージの「安全側で評価」「保存済み AI 設定要約が不足」など曖昧表現を「snapshot 欠落のため再開できません」「削除は可能」など、結果として何が起きるかを直接示す表現へ書き換える。
- 承認済み UC 差分: 差分なし（UC 正本本文の追加不要）。
- 承認済み E2E テスト観点差分: 追加候補 E2E-UC-TJL-FIX-1, E2E-UC-TJL-FIX-2, E2E-UC-TJL-FIX-3 を採用する。
- 補助記録: `data-testid-gaps.md`（削除/再開で異なる無効理由を区別する子 selector 未決 3 件）を implementation-module へ申し送る。

## 修正実行入力（論点1）

- 承認済み修正方針:
  - `buildTranslationJobManagementRuntimeSummary` が snapshot 空時に返す warning の category を `state_projection_inconsistent` から新規定数 `runtime_snapshot_missing` へ変更する。
  - 新規定数 `translationJobManagementReasonRuntimeSnapshotMissing = "runtime_snapshot_missing"` を追加する。
  - `buildTranslationJobManagementResumeBlockedReasons` で `runtime_snapshot_missing` warning を再開ブロック理由として扱う追加分岐を加える（再開不可表示を維持する）。
  - `buildTranslationJobManagementDeleteAvailability` の warnings ループは `state_projection_inconsistent` のみを対象とし、`runtime_snapshot_missing` は削除をブロックしないことを保証する（実装上は既存の category 比較条件を変えない形でよい）。
  - 表現修正: snapshot 欠落由来の warning の `Title` と `Detail` を「snapshot が無いため再開できません」「削除は可能」と直接示す文言へ書き換える。曖昧表現「安全側で評価」は採用しない。
- 禁止する修正:
  - warnings ループの全廃。
  - `job.State == ready` の特別分岐で削除を強制許可する対症療法。
  - 新規 job 状態値の追加。
- 影響ファイル候補:
  - `internal/service/translation_job_management_service.go`（`buildTranslationJobManagementRuntimeSummary`、`buildTranslationJobManagementResumeBlockedReasons`、定数定義、メッセージ表現）。
  - 上記の単体テスト（境界: Ready 直後の snapshot 欠落のみで削除可、再開不可。例外: 真の `state_projection_inconsistent` で削除拒否を維持）。
- 承認済み UC 差分: 差分なし。
- 承認済み E2E テスト観点差分: E2E-UC-TJL-FIX-1（境界: Ready 直後の削除可）、E2E-UC-TJL-FIX-2（例外: 真の状態不整合で削除拒否）、E2E-UC-TJL-FIX-3（境界: Ready 直後の再開不可は維持）。
- 画面再現確認の再現手順と修正後に満たすべき期待状態:
  - 手順: ダッシュボードから翻訳管理画面を開き、新規作成直後の Ready ジョブカードで削除ボタンと再開ボタンを確認する。
  - 期待状態（修正後）:
    - 削除ボタンが enabled になり、押下で削除確認モーダルが表示される。
    - 再開ボタンは disabled のまま維持される。
    - 再開不可理由は「snapshot が無いため再開できません」など直接的表現で表示される。
    - 削除不可理由は表示されない（warning が削除判断に作用しない）。

## 後続モジュール引き継ぎ

- task-id: translation-job-list-fix
- 後続: 論点1 は人間修正レビュー承認済み。implementation-module で `internal/service/translation_job_management_service.go` の修正実行と単体テスト追加に進む。E2E 観点の正本反映は finalization-module へ引き継ぐ。論点2, 3（UX 改善）は別途 design-module または storybook-module で扱う。論点の振り分けは後続入口で判断する。

## backend 実装結果（論点1）

- 判定: 完了
- ハーネス結果: `python3 scripts/harness/run.py --suite backend-local` 全スイート通過（lint:backend 通過、test:backend 通過）
- `go build ./...`: 通過

### 変更ファイル

- `internal/service/translation_job_management_service.go`: 定数 `translationJobManagementReasonRuntimeSnapshotMissing = "runtime_snapshot_missing"` を `translationJobManagementReason*` 定数ブロックへ追加した。`buildTranslationJobManagementRuntimeSummary` の snapshot 空時 warning の `Category` を `state_projection_inconsistent` から `runtime_snapshot_missing` へ、`Title` を「snapshot が無いため再開できません」へ、`Detail` を「保存済み AI 設定要約 (runtime snapshot) が無いため、Paused/RecoverableFailed からの再開時に外部 API 設定を確認できません。削除は可能です。」へ変更した。`buildTranslationJobManagementResumeBlockedReasons` の warnings ループを `switch` 文へ書き換え、`runtime_snapshot_missing` を再開ブロック理由として追加した（既存の `state_projection_inconsistent`・`phase_progress_aggregation_failed` 分岐は維持した）。

### 修正意図の確認

- `buildTranslationJobManagementDeleteAvailability` の warnings ループ（line 877〜883）は変更していない。対象 category は `state_projection_inconsistent` のみのため、`runtime_snapshot_missing` は削除をブロックしない。
- `buildTranslationJobManagementResumeBlockedReasons` は `runtime_snapshot_missing` を再開ブロック理由として返す。Ready 状態でも `runtime_snapshot_missing` がブロック理由リストに積まれるが、`buildTranslationJobManagementResumeAvailability` は Paused/RecoverableFailed 以外では `Enabled = true` にしないため、UI 振る舞いに矛盾しない。

### 後続への申し送り

- 単体テスト追加（境界: Ready 直後の snapshot 欠落のみで削除可・再開不可。例外: 真の `state_projection_inconsistent` で削除拒否を維持）は `implementation_unit_tester` 担当。

## 観測ログ追加結果（論点1）

- 判定: 完了
- ハーネス結果: `python3 scripts/harness/run.py --suite backend-local` 全スイート通過（lint:backend 通過、test:backend 通過）

### 追加した恒久ログ

- `logTranslationJobDeleteAvailabilityEvaluated`（`internal/service/translation_job_management_service.go`）: `buildJobDetail` 内で `buildTranslationJobManagementDeleteAvailability` の結果を受け取った直後に呼ぶ。`result: "allowed"` または `result: "rejected"` と `reason` を出力する。目的: `runtime_snapshot_missing` が削除をブロックしないこと（`result: "allowed"` が出る）と `state_projection_inconsistent` がブロックした事実（`reason: "state_projection_inconsistent"` が出る）を backend ログで一次切り分けできる。
- 再開 blocked ログ（`ResumeJob` 内、同ファイル）: `detail.ResumeAvailability.Enabled` が false の分岐で `event: "translation_job_resume_blocked"`、`result: "rejected"`、`reason` を出力する。目的: `runtime_snapshot_missing` が再開をブロックした事実を backend ログで確認できる。別経路の `state_projection_inconsistent` が再混入して再開ブロックになった場合も category で区別できる。

### 削除した一時ログ

- なし（一時ログは investigation-module で削除済み）

### 残留観測点

- なし

## 単体テスト結果（論点1）

- 判定: 完了
- ハーネス結果: `python3 scripts/harness/run.py --suite backend-local` 全スイート通過（lint:backend 通過、test:backend 通過）
- 網羅率: `python3 scripts/harness/run.py --suite coverage` は frontend 61.07%（既存水準）で 70.0% 未達。今回の backend テスト追加が直接壊した箇所ではなく、承認済み実装範囲内では改善不能なため未実行扱いとして記録する。

### 変更ファイル

- `internal/service/translation_job_management_service_test.go`: 以下 5 テストを追加した。

### 追加テスト一覧

| テスト関数名 | 証明対象 | 分類 |
| --- | --- | --- |
| `TestBuildTranslationJobManagementRuntimeSummaryReturnsRuntimeSnapshotMissingWhenSnapshotsEmpty` | snapshot 空時の warning Category が `runtime_snapshot_missing`、Title が「snapshot が無いため再開できません」、Detail に「削除は可能」を含む | 分類 |
| `TestBuildTranslationJobManagementDeleteAvailabilityEnabledWhenOnlyRuntimeSnapshotMissingWarning` | `runtime_snapshot_missing` warning のみの Ready ジョブで削除可（Enabled=true） | 境界1 |
| `TestBuildTranslationJobManagementDeleteAvailabilityBlockedWhenStateProjectionInconsistentWarning` | `state_projection_inconsistent` warning が存在する Ready ジョブで削除拒否（Enabled=false）を維持 | 例外 |
| `TestBuildTranslationJobManagementResumeBlockedReasonsContainsRuntimeSnapshotMissingForPausedJob` | Paused ジョブ・snapshot 欠落で再開ブロック理由に `runtime_snapshot_missing` を含む | 境界2 |
| `TestBuildTranslationJobManagementResumeBlockedReasonsContainsRuntimeSnapshotMissingForRecoverableFailedJob` | RecoverableFailed ジョブ・snapshot 欠落で再開ブロック理由に `runtime_snapshot_missing` を含む | 境界2 |

## 最終検証（論点1）

- 判定: 通過
- 実行 suite: `python3 scripts/harness/run.py --suite backend-local`
- 結果: Backend lint harness passed、Backend test harness passed、All requested harness suites passed
- 差し戻し履歴: なし

## Y/N 評価訂正（論点1）

- 訂正項目: 「frontend と backend を接続する」を N → Y。
- 訂正理由: 新規 backend category `runtime_snapshot_missing` を frontend gateway DTO validator が未知として弾き、`BlockedReason` 検証失敗が一覧 load 失敗を引き起こした。category 定数追加は frontend/backend の契約境界に変更が及ぶ。
- 追加修正:
  - `frontend/src/controller/wails/translation-job-management.gateway.ts`: `isReasonCategory` validator に `runtime_snapshot_missing` を追加。
  - `frontend/src/application/gateway-contract/translation-job-management/translation-job-management-gateway-contract.ts`: `TranslationJobManagementReasonCategory` union に追加。
  - `frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts`: `toReasonCategoryLabel` に「AI 設定要約欠落」case を追加。
- 検証: `python3 scripts/harness/run.py --suite frontend-local` 全 557 テスト通過。

## 実画面確認（論点1）

- Wails 接続対象: http://localhost:34115/#translation-management
- 確認結果:
  - ジョブ #1, #2（Ready 状態）で「削除」ボタンが enabled。修正の主検証点通過。
  - 「再開」ボタンは disabled のまま維持。
  - 再開不可理由が「保存済み AI 設定要約 (runtime snapshot) が無いため、Paused/RecoverableFailed からの再開時に外部 API 設定を確認できません。削除は可能です。」と直接表現で表示されている。
  - 削除不可理由は表示されない（warning が削除判断に作用しない）。
- 判定: 通過。

## 次工程

- 論点1 は実装、検証、実画面確認すべて通過。
- 残り: 論点2, 3（UX 改善）の振り分け（design-module / storybook-module）と、E2E テスト観点正本反映（finalization-module）。

## 論点2, 3 storybook 表示実装結果

- 判定: 完了
- ハーネス結果: `python3 scripts/harness/run.py --suite frontend-local` 全 555 テスト通過（lint:frontend 通過、test:frontend 通過）
- Storybook ビルド結果: `npm --prefix frontend run build-storybook` 成功（build completed successfully）
- 実画面確認: a11y スナップショット取得済み（`link "ジョブ 1 を選択して再開する"`、`button "再開"` の存在を確認）
- UI 証跡: `/Users/iorishibata/Repositories/AITranslationEngineJP/test-results/translation-job-list-fix-screen.png`

### 変更ファイル

- `frontend/src/ui/screens/translation-job-management/JobOperationGroup.svelte`: Props から `resumeOperation` を削除し、each ループを廃止して `stopOperation` の ActionButton だけを直書きにした。continue-button のラベルを「現在の翻訳段階へ進む」から「再開」へ変更した。
- `frontend/src/ui/screens/translation-job-management/JobCard.svelte`: link の aria-label を「ジョブ N を選択して現在の翻訳段階へ進む」から「ジョブ N を選択して再開する」へ変更した。JobOperationGroup への `resumeOperation` prop 渡しを削除した。job-card-reasons div（disabled-reason 表示塊）を削除した。
- `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.test.ts`: ボタンセレクタ `getByRole("button", { name: "現在の翻訳段階へ進む" })` を「再開」へ変更した。`translation-job-management-resume-button` の testid を参照するテスト 2 件（resume success feedback、resume warning feedback）を削除した。未使用になった `waitFor` import を削除した。
- `frontend/src/ui/views/AppShell.test.ts`: ボタンセレクタ `getByRole("button", { name: "現在の翻訳段階へ進む" })` を「再開」へ変更した（2 箇所）。

### Storybook 確認資源

- 変更対象コンポーネント: `JobCard.svelte`、`JobOperationGroup.svelte`
- story ファイル: `frontend/src/ui/screens/translation-job-management/stories/JobCard.stories.ts`（既存 3 story: SelectedRunning、PausedWithDisabledReason、RecoverableFailed）
- fixture: `frontend/src/ui/screens/translation-job-management/__fixtures__/translation-job-management-fixtures.ts`（変更なし。resumeOperation の type 定義と presenter 出力は残存）
- Storybook 分類: 本変更は Storybook 人間レビュー前の実装であるが、入力指示で「Storybook 起動は本変更が単純なため省略可」と明示されているため、story の作業中分類への移動は省略する。

### 画面設計根拠確認結果

- 論点2 一致確認: a11y スナップショットで `button "再開"` が continue-button として確認できた。`link "ジョブ 1 を選択して再開する"` が aria-label 通りに存在する。「再開」ActionButton（旧 resumeOperation）は DOM から消えている。
- 論点3 一致確認: a11y スナップショットにボタン下の disabled-reason テキスト（「停止: ...」「再開: ...」「削除: ...」）が存在しない。tooltip 機構（data-tooltip）は引き続き保持されている。
- UX-standard.md 対応: 変更は承認済み画面設計根拠の範囲内であり、UX-standard.md で禁止された変更には該当しない。
- console エラー: 404 Not Found が 1 件（既存の wails dev 環境での静的リソース欠落。今回の変更とは無関係）。

## 合意済み frontend 保護（論点2, 3）

- 承認済み画面: 翻訳管理 → 未完了ジョブカード（JobCard.svelte + JobOperationGroup.svelte）。
- 表示規則:
  - 操作ボタンは「再開」（旧 continue-button、enabled 時 phase 画面へ遷移）、「停止」、「削除」の 3 種。
  - resumeOperation の ActionButton は表示しない（store/usecase/contract コードは保持）。
  - 選択不可理由は hover tooltip でのみ表示し、ボタン下に羅列しない。
- 確認済み Storybook 状態: build-storybook 成功。Storybook 人間レビューループは入力指示により省略。
- 変更禁止範囲: resumeOperation の type 定義、store、usecase、contract、backend は触らない（フィールドの完全削除は別 task 扱い）。tooltip 機構（FloatingTooltipTrigger、ActionButton の title）は変更しない。
- 反映先 frontend ファイル:
  - `frontend/src/ui/screens/translation-job-management/JobOperationGroup.svelte`
  - `frontend/src/ui/screens/translation-job-management/JobCard.svelte`
  - 影響テスト: `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.test.ts`、`frontend/src/ui/views/AppShell.test.ts`
- 表示範囲外の残課題: resumeOperation の contract / presenter / store / backend ResumeJob を完全に整理するかは別 task で扱う。本 task では「表示から外す」だけにとどめた。

## 実画面確認（論点2, 3）

- Wails 接続対象: http://localhost:34115/#translation-management
- 確認結果:
  - 操作ボタンが「再開」「停止」「削除」の 3 種に減った（旧「再開」ActionButton は消えた）。
  - 「再開」ボタンが enabled、aria-label「ジョブ 1 を選択して再開する」になっている。
  - ボタン下の選択不可理由テキストが表示されない（tooltip は維持）。
- 判定: 通過。

## 三論点完了サマリ

- 論点1（削除不可バグ）: backend category 分離（`runtime_snapshot_missing`）、frontend gateway DTO 拡張、単体テスト 5 件、観測ログ恒久化、実画面確認すべて通過。
- 論点2（再開ボタン削除）: 「再開」ActionButton 削除と continue-button のラベル「再開」化、実画面確認通過。
- 論点3（選択不可理由表示削除）: ボタン下の disabled-reason 表示削除、tooltip は維持、実画面確認通過。
- harness 結果: backend-local 通過、frontend-local 通過（555 テスト）。
- 次工程: 論点1 の E2E テスト観点正本反映（finalization-module）と、作業 commit、マージ準備入力（finalization-module）。

## 正本化判断

- 詳細仕様正本（`docs/detail-specs/translation-job-management.md`）への反映: 不要。論点1 の UC 差分は「差分なし」判定（REQ-002, REQ-004 で既存仕様が説明できる）。論点2, 3 は表示変更のみで詳細仕様に新規仕様文の追加は無い。
- E2E テスト観点正本（`docs/e2e-test-design/test-design.csv`）への反映: 後続課題に切り出す。
  - 理由: E2E-UC-TJL-FIX-1〜3 は `data-testid-gaps.md` に記録された削除/再開で異なる無効理由を区別する子 selector の未決 3 件に依存する。selector 未決のままで観点正本に投入すると「実行不能な観点」が混入するため、selector 確定（別 task）後にまとめて反映する。
  - 申し送り先: 次の selector 確定 task で `data-testid-gaps.md` の解消と E2E-UC-TJL-FIX-1〜3 の正本投入を同時に行う。
- 人間承認状態: 承認済み（task 入力提出者の継続承認）。
- 判定: 詳細仕様正本反映 なし、E2E 観点正本反映 後続課題に切り出す、`docs_updater` 起動 なし。


## 作業 commit

- commit hash: 0e16117d0a4c0972eb18534e2c5b8ec74c347683
- 作業 branch: claude/translation-job-list-fix
- 分岐元: master @ c2643f876e466123894287dd635b85915a03b214
- 変更ファイル: 11 files changed, 620 insertions(+), 125 deletions(-)
  - backend: `internal/service/translation_job_management_service.go`、`internal/service/translation_job_management_service_test.go`
  - frontend gateway/contract/presenter: `frontend/src/controller/wails/translation-job-management.gateway.ts`、`frontend/src/application/gateway-contract/translation-job-management/translation-job-management-gateway-contract.ts`、`frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts`
  - frontend 表示: `frontend/src/ui/screens/translation-job-management/JobCard.svelte`、`frontend/src/ui/screens/translation-job-management/JobOperationGroup.svelte`
  - 影響テスト: `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.test.ts`、`frontend/src/ui/views/AppShell.test.ts`
  - active plan: `docs/exec-plans/active/translation-job-list-fix/plan.md`、`docs/exec-plans/active/translation-job-list-fix/data-testid-gaps.md`（新規）
- 検証結果: backend-local 通過、frontend-local 通過（555 テスト）、実画面 GUI 確認通過。
- 残留リスク:
  - `resumeOperation` の store / usecase / contract / backend `ResumeJob` 経路は表示から外しただけで、フィールドと関連経路は残置している。完全削除は別 task で扱う。
  - `data-testid-gaps.md` の selector 未決 3 件と、E2E 観点 E2E-UC-TJL-FIX-1〜3 の正本投入は後続 task に切り出し。

## マージ準備入力

- active plan folder: `docs/exec-plans/active/translation-job-list-fix/`
- source branch: `claude/translation-job-list-fix`
- target branch: `master`
- 作業 commit hash: `0e16117d0a4c0972eb18534e2c5b8ec74c347683`
- 最終検証結果: backend-local 通過、frontend-local 通過（555 テスト）、実画面 GUI 確認通過。
- 残留リスク（merge_lane への申し送り）:
  - `resumeOperation` の経路は frontend 表示からのみ外している。usecase/store/contract/backend ResumeJob は残置のため、merge 後の動作に影響しないことを念のため確認する。
  - 後続 task: selector 未決 3 件の解消と E2E 観点 E2E-UC-TJL-FIX-1〜3 の `docs/e2e-test-design/test-design.csv` 反映。

