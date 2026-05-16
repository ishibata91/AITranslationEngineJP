# Scenario Candidates: 2026-05-16-dev-fake-secret-store / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `DFSS-OA`

## Generator Scope

- `viewpoint`: operation-audit
- `included_sources`: `./plan.md`, `internal/bootstrap/app_controller.go`, `internal/repository/provider_settings_keyring_secret_store.go`, `internal/repository/provider_settings_cached_secret_store.go`, `internal/repository/provider_settings_repository.go`, `internal/repository/master_persona_repository.go`, `scripts/dev/run-wails-agent-browser.sh`, `docs/detail-specs/ai-provider-settings-management.md`, `docs/architecture.md`, `docs/er.md`
- `excluded_sources`: プロダクトコード変更、プロダクトテスト変更、docs 正本変更、最終シナリオ表、候補の採否、候補の統合判断
- `generation_notes`: 候補は開発実行、agent-browser 証跡、log、review evidence、再現手順の後追い確認に限定する。secret 平文、復号可能値、credential 参照実値、secret store key は記録対象にしない。fake secret store は user-facing 設定として扱わない。

## Candidate Scenarios

### CAND-DFSS-OA-001 開発起動時の secret store 選択を後追い確認できる

- `source requirement`: `plan.md` の D-01、D-02、D-03、D-06。`app_controller.go` は Wails 起動時に provider settings secret store を構築する。`provider_settings_keyring_secret_store.go` は既定で OS keyring backend を選ぶ。
- `viewpoint`: operation-audit / 後追い確認 / 監査ログ / 保存禁止
- `candidate scenario id`: `CAND-DFSS-OA-001`
- `actor`: 開発者、または agent-browser 確認を実行する AI
- `trigger`: 開発用 Wails 起動を行う。
- `expected outcome`: 開発起動が OS keyring prompt を待たずに進んだ事実を、後から確認できる。production 既定の secret store が OS keyring のままである事実も、別の確認結果として区別できる。
- `observable point`: 開発起動結果、起動時刻、起動経路、secret store の分類、OS keyring prompt が発生しなかった結果、起動失敗時の短い error kind。
- `record`: 起動 command の分類、開発用 secret store が有効だったかの分類、process-local かどうかの分類、起動成功または失敗の結果、失敗時の error kind。
- `do not record`: APIキー平文、復号可能値、credential 参照実値、secret store key、secret store 内の key 一覧、開発用環境変数の値、OS keyring item 名、外部 provider raw payload。
- `related detail requirement type`: `observability_requirement`, `security_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: designer は、開発起動の証跡を scenario に採用する場合、secret store の分類だけを観測点にする。具体的な secret 識別子は受け入れ条件に入れない。
- `conflict hint`: secret store の分類を log へ残す粒度は、保存禁止情報と衝突する可能性がある。分類名が user-facing 設定へ漏れる場合も D-01 と衝突する。

### CAND-DFSS-OA-002 agent-browser 証跡で fake secret store が画面に露出しないことを確認できる

- `source requirement`: `plan.md` の D-01、D-03、D-06。`docs/detail-specs/ai-provider-settings-management.md` は利用者向け provider list に fake provider を表示しないこと、secret を UI、DTO、log、保存要約へ出さないことを固定している。
- `viewpoint`: operation-audit / 後追い確認 / 再現材料 / 保存禁止
- `candidate scenario id`: `CAND-DFSS-OA-002`
- `actor`: agent-browser 確認を実行する AI
- `trigger`: `agent-browser` で開発中アプリを開き、AIサービス設定または関連画面を確認する。
- `expected outcome`: agent-browser 証跡から、画面が起動できたこと、fake secret store が user-facing 設定として表示されないこと、APIキー本文が表示されないことを確認できる。
- `observable point`: 画面到達結果、provider list の表示分類、APIキー状態の表示分類、password prompt による停止がない結果。
- `record`: 確認 URL、到達した画面名、provider list に利用者向け provider だけが表示された結果、APIキー状態が存在状態または伏せ字で表示された結果、確認時刻、スクリーンショットまたは操作ログの保存先。
- `do not record`: APIキー入力値、復号可能値、credential 参照実値、secret store key、OS keyring item 名、fake secret store を設定項目として見せる文言、外部 provider raw request、外部 provider raw response。
- `related detail requirement type`: `security_requirement`, `observability_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: UI 人間操作 E2E 候補に統合する場合、観測点は「表示されないこと」と「prompt で止まらないこと」に寄せる。
- `conflict hint`: fake secret store の有効状態を画面へ表示する受け入れ条件にすると、D-01 と user-facing 非表示方針に衝突する。

### CAND-DFSS-OA-003 開発 log で secret 非露出と起動失敗理由を確認できる

- `source requirement`: `plan.md` の D-03、D-06。`run-wails-agent-browser.sh` は開発起動 log を `tmp/logs/wails-dev.log` に出す。`docs/architecture.md` と `docs/detail-specs/ai-provider-settings-management.md` は secret、APIキー、credential 参照実値、raw payload を log や保存要約へ出さないことを固定している。
- `viewpoint`: operation-audit / 監査ログ / 保存禁止 / 再現材料
- `candidate scenario id`: `CAND-DFSS-OA-003`
- `actor`: 開発者、review 担当者
- `trigger`: 開発起動後に log を確認する。
- `expected outcome`: log から、開発起動が成功したか、失敗したか、失敗した場合は分類だけを確認できる。log に secret 平文や復号可能値が含まれない。
- `observable point`: log file の存在、起動結果、secret store 分類、password prompt で停止しなかった結果、失敗時の error kind、redaction 確認結果。
- `record`: log file path、起動成功または失敗、短い error kind、redaction 確認の pass / fail、確認者、確認時刻。
- `do not record`: APIキー平文、復号可能値、credential 参照実値、secret store key、開発用環境変数の値、credential を含む endpoint 値、provider raw request、provider raw response、raw prompt。
- `related detail requirement type`: `observability_requirement`, `security_requirement`, `failure_handling_requirement`, `testability_requirement`
- `adoption hint`: designer は、log 内容の固定ではなく、log に残す分類と log に残さない情報を受け入れ条件へ分ける。
- `conflict hint`: provider settings の更新履歴は保存しない正本がある。log 候補が更新履歴の永続保存へ広がる場合は data requirement と衝突する。

### CAND-DFSS-OA-004 review evidence で secret 境界の維持を確認できる

- `source requirement`: `plan.md` の D-02、D-06。`docs/detail-specs/ai-provider-settings-management.md` は実装後検証で fake transport DI と fake secret store を使い、有料の実 AI API を呼ばないことを固定している。
- `viewpoint`: operation-audit / 後追い確認 / 履歴 / 保存禁止
- `candidate scenario id`: `CAND-DFSS-OA-004`
- `actor`: contract review 担当者、trust-boundary review 担当者
- `trigger`: 実装後 review evidence を確認する。
- `expected outcome`: review evidence から、production 既定が OS keyring のままであること、開発起動だけが fake secret store を使うこと、UI、DTO、log、browser evidence に secret が出ないことを確認できる。
- `observable point`: review 対象、review 結果、secret 境界確認項目、production 既定確認、user-facing 非表示確認。
- `record`: review 種別、review 結果、確認した境界の分類、production 既定の維持結果、fake secret store が user-facing 設定に出ていない結果、未解決リスク。
- `do not record`: APIキー平文、復号可能値、credential 参照実値、secret store key、具体 credential、外部 provider raw payload、review 用に採取した secret 値。
- `related detail requirement type`: `security_requirement`, `compatibility_requirement`, `observability_requirement`, `testability_requirement`
- `adoption hint`: designer は、review evidence を最終検証の補助観測点として扱う。review evidence 自体を product behavior の代替にしない。
- `conflict hint`: review evidence が secret 値の提示を要求する場合、D-06 と保存禁止観点に衝突する。

### CAND-DFSS-OA-005 再現手順で password prompt 回避を再確認できる

- `source requirement`: `plan.md` の validation commands、D-03、D-04、D-05、Scenario Seeds。`run-wails-agent-browser.sh` は `.env` を読み込んだ後に Wails dev を起動する。
- `viewpoint`: operation-audit / 再現材料 / 後追い確認 / 人間判断候補
- `candidate scenario id`: `CAND-DFSS-OA-005`
- `actor`: 開発者、実装後確認担当者
- `trigger`: 別環境または別セッションで、開発起動と agent-browser 確認を再実行する。
- `expected outcome`: 再現手順から、OS keyring prompt が出ないこと、ブラウザ確認が進むこと、process-local secret store の再起動後消失を確認できる。
- `observable point`: 再現手順、事前条件、実行 command、起動 URL、再起動前後の APIキー状態分類、prompt 発生有無。
- `record`: 再現手順の command 名、事前条件の分類、期待する観測結果、実行結果、再起動前後の secret 状態分類、失敗時の error kind。
- `do not record`: `.env` の secret 値、開発用環境変数の値、APIキー平文、復号可能値、credential 参照実値、secret store key、OS keyring item 名。
- `related detail requirement type`: `testability_requirement`, `observability_requirement`, `security_requirement`, `recovery_requirement`
- `adoption hint`: designer は、再現手順を受け入れテストの入口にする場合、実 secret を使わない fake provider / fake secret store の前提を固定する。
- `conflict hint`: `.env` と起動 script の優先順位は plan 上で未確定である。再現手順にどちらを優先するか固定するには人間判断が必要になる可能性がある。

### CAND-DFSS-OA-006 production 既定が OS keyring から変わっていないことを確認できる

- `source requirement`: `plan.md` の D-02、D-06。`provider_settings_keyring_secret_store.go` は空または default の backend 指定で OS 別 keyring backend を選ぶ。`docs/er.md` は `credential_ref` が API key そのものではなく secret store への参照であることを固定している。
- `viewpoint`: operation-audit / 後追い確認 / 監査ログ / 競合候補
- `candidate scenario id`: `CAND-DFSS-OA-006`
- `actor`: 実装後確認担当者、review 担当者
- `trigger`: production 相当の起動条件、または開発用 override なしの構成確認を行う。
- `expected outcome`: fake secret store が production 既定へ混入していないことを確認できる。DB や UI が APIキー平文や復号可能値を持たないことも確認できる。
- `observable point`: secret store 既定分類、provider settings row の保存内容分類、credential 状態分類、APIキー本文が UI と DTO に出ない結果。
- `record`: production 相当の起動条件分類、secret store 既定分類、DB に保存される credential 参照状態の分類、UI / DTO の secret 非露出確認結果。
- `do not record`: APIキー平文、復号可能値、credential 参照実値、secret store key、secret store 内の key 一覧、OS keyring item 名、外部 provider raw payload。
- `related detail requirement type`: `compatibility_requirement`, `security_requirement`, `data_requirement`, `observability_requirement`
- `adoption hint`: designer は、production 既定維持を回帰シナリオへ統合できる。fake secret store の有効条件とは別シナリオにしてもよい。
- `conflict hint`: production 確認を実 OS keyring 操作必須にすると、agent-browser 安定化の目的と衝突する可能性がある。検証段階を APIテスト、lower-level only、最終検証のどこへ置くかは designer が決める。

## Open Notes

- `human decision candidate`: `.env` と起動 script のどちらが開発用 secret store 設定の優先権を持つかは、候補だけでは確定しない。
- `human decision candidate`: log に secret store の分類名を残すか、redaction pass / fail だけを残すかは、人間判断が必要になる可能性がある。
- `human decision candidate`: process-local secret store の再起動後消失を受け入れ条件に含めるか、将来の file backend fallback に残すかは、designer が判断する。
- `merge candidate`: `CAND-DFSS-OA-001` と `CAND-DFSS-OA-005` は、開発起動の再現手順シナリオへ統合できる可能性がある。
- `merge candidate`: `CAND-DFSS-OA-002` と user-facing 非表示の外部連携候補は、同じ UI 人間操作 E2E へ統合できる可能性がある。
- `merge candidate`: `CAND-DFSS-OA-003` と `CAND-DFSS-OA-004` は、secret redaction の review / log 証跡シナリオへ統合できる可能性がある。
- `rejection candidate`: provider settings の更新履歴を永続保存する候補は、既存詳細仕様の対象外と衝突するため不採用候補になる。
- `rejection candidate`: secret 値の提示で keyring 非使用を証明する候補は、D-06 と保存禁止観点に衝突するため不採用候補になる。
