# Scenario Candidates: 2026-05-07-model-settings-card-controller / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./task-frame.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `MSC-OA`

## Generator Scope

- `viewpoint`: operation-audit
- `included_sources`: `task-frame.md`, `light-change-planning.md`, `docs/detail-specs/translation-job-setup.md`, `docs/detail-specs/ai-provider-settings-management.md`, `docs/architecture.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本本文の変更、他観点候補の生成、採否、統合、最終シナリオ表
- `generation_notes`: モデル設定カード controller 集約に対して、後追い確認、保存要約、伏せ字規則、観測点だけを候補化する。永続監査ログを追加するかは未確定のため、人間判断候補へ残す。

## Candidate Scenarios

### CAND-MSC-OA-001 モデル一覧更新結果を要約で確認できる

- `source requirement`: `task-frame.md:12`, `task-frame.md:20`, `light-change-planning.md:10`, `translation-job-setup.md:39-40`, `translation-job-setup.md:62-65`, `ai-provider-settings-management.md:29`, `ai-provider-settings-management.md:37`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-MSC-OA-001`
- `actor`: 利用者
- `trigger`: 利用者が共有モデル設定カードで provider を選び、model list 更新を実行する。
- `expected outcome`: 画面は更新中、取得済み、取得失敗、APIキー未設定で更新不可を区別して表示する。保存または観測される要約は、provider、対象画面、結果分類、短いエラー要約、取得できた model 数に限定する。
- `audit event`: model list 更新開始、更新成功、更新失敗、更新不可を後から区別できる事象として扱う。
- `saved summary`: provider、対象画面、結果分類、model 数、短いエラー要約、再更新が必要な状態だけを候補にする。
- `redaction rule`: APIキー、raw request、raw response、raw prompt、外部 provider 応答原文、内部ログ用識別子は保存要約、UI、DTO、ログに出さない。
- `observable point`: model list 更新後の表示状態、エラー文、保存または一時状態の要約、fake transport での結果分類を確認する。
- `related detail requirement type`: `observability_requirement`, `security_requirement`, `testability_requirement`, `data_requirement`
- `adoption hint`: designer が model list 更新の受け入れ観点へ統合するか判断するための候補であり、採否は確定しない。
- `conflict hint`: raw payload 非保存と、障害調査用の再現材料をどこまで持つかが競合しうる。

### CAND-MSC-OA-002 provider 変更後の遅延応答破棄を後追い確認できる

- `source requirement`: `task-frame.md:20`, `light-change-planning.md:46-47`, `architecture.md:101-117`, `architecture.md:119-130`, `translation-job-setup.md:64-65`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-MSC-OA-002`
- `actor`: 利用者
- `trigger`: 利用者が provider を変更した直後に、変更前 provider の model list 応答が遅れて返る。
- `expected outcome`: 共有 controller / usecase / store は、現在選択中の provider と一致しない応答を画面状態へ反映しない。利用者から見る状態は、現在 provider の更新状態、model 選択状態、失敗状態だけで説明できる。
- `audit event`: 遅延応答受信、破棄判定、現在 provider の状態維持を後追いできる事象として扱う。
- `saved summary`: 破棄された provider、現在 provider、破棄理由、画面状態を更新しなかった結果だけを候補にする。
- `redaction rule`: request token、内部ログ用識別子、外部応答原文、secret は保存要約へ出さない。必要な場合も短い分類名に丸める。
- `observable point`: provider 変更後に古い応答が返っても、model 候補、model 選択、エラー表示が現在 provider へ戻らないことを確認する。
- `related detail requirement type`: `observability_requirement`, `consistency_requirement`, `concurrency_requirement`, `security_requirement`
- `adoption hint`: designer が状態遷移候補または失敗候補と統合するか判断するための候補であり、統合は確定しない。
- `conflict hint`: 再現材料として内部 request 識別子を残したくなるが、内部ログ用識別子の表示禁止と衝突しうる。

### CAND-MSC-OA-003 model 選択、保存、取得の結果を再現材料として確認できる

- `source requirement`: `task-frame.md:5-6`, `task-frame.md:12`, `task-frame.md:18-20`, `light-change-planning.md:10-12`, `translation-job-setup.md:38-45`, `ai-provider-settings-management.md:13-15`, `ai-provider-settings-management.md:34`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-MSC-OA-003`
- `actor`: 利用者
- `trigger`: 利用者がマスターペルソナまたは翻訳ジョブ設定で model を選択し、保存後に再取得する。
- `expected outcome`: 保存結果と再取得結果は、対象画面、provider、model、credential 参照状態、保存成否、取得成否の要約で確認できる。Job Setup は master-persona の AI 設定を既定値または保存元として扱わない。
- `audit event`: model 選択、保存開始、保存成功、保存失敗、再取得成功、再取得失敗を後追いできる事象として扱う。
- `saved summary`: 対象画面、provider、model、credential 参照状態、保存結果、取得結果だけを候補にする。
- `redaction rule`: APIキー文字列、secret、endpoint 以外の接続詳細、raw payload、内部ログ用識別子は保存しない。Job Setup の作成後表示も APIキー状態だけにする。
- `observable point`: 保存前後のカード表示、再読込後の選択状態、Job Setup と master-persona の相互混入がないことを確認する。
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `consistency_requirement`, `security_requirement`, `compatibility_requirement`
- `adoption hint`: designer が保存取得シナリオへ採用するか、画面別候補へ分けるか判断するための候補であり、採否は確定しない。
- `conflict hint`: 共有 controller 化により保存要約の形を揃える必要がある一方、Job Setup は master-persona を保存元にしない既存仕様を維持する必要がある。

### CAND-MSC-OA-004 secret 非露出のまま更新不可と失敗理由を確認できる

- `source requirement`: `translation-job-setup.md:40-45`, `translation-job-setup.md:61-66`, `ai-provider-settings-management.md:28-32`, `ai-provider-settings-management.md:60-62`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-MSC-OA-004`
- `actor`: 利用者
- `trigger`: 利用者が APIキー未設定の provider で model list 更新または保存を試みる。
- `expected outcome`: APIキーが必要な provider では、APIキー未設定のため更新できない状態を確認できる。LM Studio では APIキー未設定 warning と credential select を出さない。
- `audit event`: 更新不可判定、保存不可判定、失敗表示、再更新要求を後追いできる事象として扱う。
- `saved summary`: provider、credential 参照状態、結果分類、短い日本語エラー要約、次に必要な操作だけを候補にする。
- `redaction rule`: APIキー本体、復号可能値、raw request、raw response、raw prompt、fake transport log の raw data は保存しない。
- `observable point`: APIキー未設定時の disabled 状態、LM Studio の例外表示、保存失敗文言、接続確認または保存結果の短い要約を確認する。
- `related detail requirement type`: `security_requirement`, `observability_requirement`, `failure_handling_requirement`, `testability_requirement`
- `adoption hint`: designer が失敗候補と統合するか判断するための候補であり、採否は確定しない。
- `conflict hint`: 失敗調査に詳しい外部応答を残す要望が出る場合、raw payload 非露出の既存仕様と衝突する。

### CAND-MSC-OA-005 fake mode の provider 表示と保存要約が漏れない

- `source requirement`: `task-frame.md:10`, `task-frame.md:13-14`, `task-frame.md:19`, `light-change-planning.md:11-12`, `light-change-planning.md:25`, `ai-provider-settings-management.md:26`, `ai-provider-settings-management.md:38`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-MSC-OA-005`
- `actor`: 利用者
- `trigger`: fake mode で、利用者が通常 provider ID のまま model list を更新し、`fake-model` を選択する。
- `expected outcome`: 利用者向け provider list に `fake` provider ID は出ない。frontend は fake mode 判定や `fake-model` 固有分岐を持たず、通常 provider ID と model list 結果だけを表示する。
- `audit event`: fake mode での model list 取得、`fake-model` 表示、model 選択、保存または取得を後追いできる事象として扱う。
- `saved summary`: 通常 provider ID、model 名、結果分類、対象画面だけを候補にする。fake provider ID を user-facing な保存要約へ出さない。
- `redaction rule`: fake mode 判定値、fake provider ID、fake transport の raw log、secret は user-facing provider list、保存要約、UI 表示へ出さない。
- `observable point`: provider list に `fake` がないこと、通常 provider ID のまま `fake-model` を選べること、保存取得の要約が fake 固有分岐を示さないことを確認する。
- `related detail requirement type`: `observability_requirement`, `security_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: designer が external-integration 候補または actor-goal 候補と統合するか判断するための候補であり、統合は確定しない。
- `conflict hint`: テスト再現性のため fake mode を観測したい要求と、frontend や user-facing provider list に fake 固有分岐を出さない制約が競合しうる。

## Open Notes

- `human decision candidate`: モデル設定カードの保存結果、取得結果、model list 更新結果、遅延応答破棄を永続履歴として残すか、画面状態とテスト観測だけにするかは未確定である。
- `human decision candidate`: 障害調査用の再現材料として、破棄された応答の識別情報をどの粒度で残せるかは未確定である。secret と内部ログ用識別子を出さない既存仕様との境界判断が必要である。
- `conflict candidate`: provider settings の更新履歴は保存しない既存仕様がある。共有モデル設定カードで model 選択履歴を保存する案は競合しうる。
- `conflict candidate`: fake mode の再現性を高める観測情報は、`fake` provider ID を user-facing provider list に追加しない制約と競合しうる。
- `merge candidate`: designer が扱う。operation-audit generator は統合判断をしない。
- `rejection candidate`: designer が扱う。operation-audit generator は採否判断をしない。
