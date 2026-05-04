# Scenario Design Questions: ai-provider-settings-management

- `skill`: scenario-design
- `status`: answered
- `source`: `./scenario-design.md`

## [Q-AIPSM-001] endpoint 既定値と検証条件

状態:
回答済み。

質問:
Gemini と xAI の endpoint は必須入力にするか、既定 endpoint を保存値なしで使うか。

やりたいこと:
利用者が provider settings で接続先を管理し、endpoint 更新後の検証状態を迷わず判断できるようにする。

背景:
LM Studio は local endpoint が実質必須である。一方で Gemini と xAI は公式 endpoint を既定値として持てる可能性がある。保存値なしの既定 endpoint を許すかどうかで、未設定状態、DB 保存値、接続確認条件が変わる。

選択肢:
1. LM Studio だけ endpoint 必須にし、Gemini と xAI は既定 endpoint を使う。
2. 全 provider で endpoint 入力または確認済み endpoint を必須にする。
3. Gemini と xAI は既定 endpoint を表示して保存対象にし、利用者が変更できるようにする。
4. その他

AI推奨:
3

推奨理由:
provider settings の保存画面として endpoint を明示できる。既定値を見せるため、利用者は何が使われるか確認できる。空 endpoint と既定 endpoint の意味も分けやすい。

不確実性:
公式 endpoint を変更させたくない運用なら、選択肢 1 の方が保守しやすい可能性がある。

回答形式:
選択肢番号を選んでください。
4 の場合は、採用したい業務ルールを1〜3文で記入してください。
3

反映結果:
Gemini と xAI は既定 endpoint を表示して保存対象にする。
利用者は既定 endpoint を確認し、必要な場合は変更できる。

## [Q-AIPSM-002] 保存失敗時の補償方針

状態:
回答済み。

質問:
secret store と DB 設定値の片方だけが保存成功した場合、どの状態へ戻すか。

やりたいこと:
APIキー保存と endpoint 保存の不整合を、利用者と後続参照へ漏らさないようにする。

背景:
APIキーは secret store、endpoint などは DB に分かれる。片方だけ成功した場合に完了表示すると、参照側が存在しない secret や古い endpoint を使う可能性がある。

選択肢:
1. 保存単位を transaction 相当に扱い、片方が失敗したら provider settings 全体を失敗にする。
2. secret 保存成功後に DB 保存が失敗した場合、secret を補償削除する。
3. secret 保存成功後に DB 保存が失敗した場合、未参照 secret として残して再保存だけ許可する。
4. その他

AI推奨:
1

推奨理由:
利用者に見える状態と後続参照の状態を一致させやすい。補償削除の失敗や未参照 secret の掃除を、初期実装の必須複雑度にしなくて済む。

不確実性:
OS keychain の削除失敗を確実に扱えるなら、選択肢 2 の方が秘密値の残留を減らせる。

回答形式:
選択肢番号を選んでください。
4 の場合は、採用したい業務ルールを1〜3文で記入してください。
1

反映結果:
secret store と DB 設定値は transaction 相当の保存単位として扱う。
片方だけが成功した場合、provider settings 全体を保存失敗にする。

## [Q-AIPSM-003] model 保存単位と参照側の優先順位

状態:
回答済み。

質問:
provider settings の model と、Job Setup / master-persona 側の model 選択はどちらを優先するか。

やりたいこと:
共通の接続設定と、用途ごとの model 選択を混同しないようにする。

背景:
依頼では provider settings で model を設定できる。一方で既存 Job Setup は phase 別 model を持つ。provider settings の model を全用途の既定値にするか、参照側で用途別に上書きするかを決める必要がある。

選択肢:
1. provider settings の model を既定値にし、Job Setup と master-persona は用途別に上書きできる。
2. provider settings の model を唯一の model source にし、参照側の model 選択は削除する。
3. provider settings は接続確認用 model だけを持ち、実行時 model は参照側で必ず選ぶ。
4. その他

AI推奨:
1

推奨理由:
provider settings で model を保存する要件を満たしつつ、翻訳 phase ごとに異なる model を選びたい既存要件も残せる。

不確実性:
設定画面を単純化したい場合は、選択肢 2 の方が利用者負荷を下げる可能性がある。

回答形式:
選択肢番号を選んでください。
4 の場合は、採用したい業務ルールを1〜3文で記入してください。
4 AIプロバイダ設定画面はAPIキーとエンドポイントだけ設定できる
モデルは設定できない。
各フェーズ・マスターペルソナではモデルと処理方法・どのプロバイダ使うか，しか設定できない。

反映結果:
AIサービス設定画面は APIキーと endpoint だけを設定対象にする。
model、処理方法、利用 provider、Batch API 切り替えは、各翻訳フェーズと master-persona 側の設定として扱う。
この回答により、model 保存単位と provider settings 側の Batch API 保存に関する質問は不要になった。

## [Q-AIPSM-004] 既存 job と実行中処理への反映時点

状態:
回答済み。

質問:
provider settings 更新後、作成済み job や実行中 phase は開始時 snapshot と最新設定のどちらを使うか。

やりたいこと:
再開、失敗回復、障害調査で、どの接続設定を使ったかを一貫して説明できるようにする。

背景:
常に最新設定を読むと、更新後の再開で挙動が変わる。開始時 snapshot を残すと、古い endpoint 要約や credential 参照状態を保持する必要がある。

選択肢:
1. 作成済み job は作成時 snapshot を使い、実行中 phase は開始時 snapshot を使う。
2. Ready job は最新 provider settings を再解決し、Running phase は開始時 snapshot を使う。
3. すべて最新 provider settings を読む。
4. その他

AI推奨:
2

推奨理由:
実行中処理の安定性を守りつつ、まだ開始していない job は設定修正を反映できる。利用者が設定ミスを直した後に再作成だけを強いられにくい。

不確実性:
監査の単純さを最優先するなら、選択肢 1 の方が説明しやすい。

回答形式:
選択肢番号を選んでください。
4 の場合は、採用したい業務ルールを1〜3文で記入してください。
2

反映結果:
Ready job は最新 provider settings を再解決する。
Running phase は開始時 snapshot を使う。

## [Q-AIPSM-005] 削除と secret metadata の扱い

状態:
回答済み。

質問:
provider settings の削除またはリセットでは、DB 設定値、secret 参照、secret 本体をどう扱うか。

やりたいこと:
設定を消した後に外部 request が送られず、APIキー平文も再表示されない状態に戻す。

背景:
削除を hard delete にするか、未設定状態に戻すか、secret だけ削除するかで、復元、監査、再保存、DB migration の扱いが変わる。
Q-AIPSM-003 の回答により、model と Batch API は provider settings の削除対象ではない。

選択肢:
1. provider settings row は残し、endpoint と APIキー状態を未設定へ戻し、secret 本体は削除する。
2. provider settings row と secret 本体を両方 hard delete する。
3. secret 本体だけ削除し、endpoint は残す。
4. その他

AI推奨:
1

推奨理由:
provider list の表示設定状態を安定させやすい。削除後も「この provider は未設定」と説明でき、secret は残さない。

不確実性:
endpoint の再入力負荷を下げるなら、選択肢 3 が合う可能性がある。

回答形式:
選択肢番号を選んでください。
4 の場合は、採用したい業務ルールを1〜3文で記入してください。
1

反映結果:
provider settings row は残す。
endpoint と APIキー状態は未設定へ戻す。
secret 本体は削除する。

## [Q-AIPSM-006] 監査履歴と endpoint 要約の粒度

状態:
回答済み。

質問:
provider settings の保存、検証、実行参照について、どの監査要約を残すか。

やりたいこと:
障害調査に必要な再現材料を残しつつ、APIキー、raw payload、endpoint の過剰露出を避ける。

背景:
operation-audit 候補は、更新履歴、接続確認結果、実行時設定断面、endpoint 要約を挙げている。保持期間と伏せ字粒度を決めないと、observability と data minimization が衝突する。

選択肢:
1. 直近更新要約だけを残し、endpoint は host のみ表示する。
2. 更新履歴を一定件数残し、endpoint は hash または fingerprint のみ表示する。
3. 実行時 snapshot / version も残し、endpoint は host と path まで表示する。
4. その他

AI推奨:
2

推奨理由:
履歴の追跡性と秘密値の最小露出のバランスがよい。endpoint 全文を出さず、変更有無を比較できる。

不確実性:
初期実装を小さくしたい場合は、選択肢 1 が適している可能性がある。

回答形式:
選択肢番号を選んでください。
4 の場合は、採用したい業務ルールを1〜3文で記入してください。
4　ローカル運用なので，endpointは見れるけど，シークレットは伏せ字，履歴保存は不要

反映結果:
ローカル運用のため、endpoint は画面と保存要約で表示できる。
secret は伏せ字または存在状態だけを表示する。
provider settings の更新履歴は保存しない。
