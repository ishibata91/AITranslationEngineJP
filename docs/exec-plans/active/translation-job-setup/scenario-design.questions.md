# Human Decision Questionnaire

- `status`: answered
- `reflected_in`: `scenario-design.md`, `scenario-design.requirement-coverage.json`, `scenario-design.candidate-coverage.json`, `ui-design.md`
- `note`: Q-001 から Q-008 までの回答履歴として保持する。

## [Q-001] validation blocking 条件

質問:
Job Setup の validation failure を、どこまで create job 禁止の blocking として扱うかを決めてください。

やりたいこと:
翻訳担当者が create job 前に修正すべき問題と、warning のまま進めてよい問題を区別する。

背景:
source は validation failure の理由表示を要求しているが、blocking / warning の分類は固定していない。ここを決めないと create job button の有効条件と scenario acceptance が確定しない。

選択肢:
1. 必須設定不足、参照不能、provider / mode 不整合、credential 参照不能をすべて blocking にする
2. 必須設定不足、参照不能、provider / mode 不整合を blocking にし、credential reachability は warning にする
3. 必須設定不足だけを blocking にし、他は phase 実行時に失敗させる
4. その他

AI推奨:
1

推奨理由:
Job Setup の goal は実行前 validation 完了であるため、create 前に検出できる不整合は fail closed にした方が後続 phase の失敗が減る。

不確実性:
provider reachability を毎回検査すると offline 作業や provider 側一時障害で setup が止まりすぎる可能性がある。

回答形式:
選択肢番号を選んでください。
4 の場合は、採用したい業務ルールを1〜3文で記入してください。
1

## [Q-002] provider 検証の範囲

質問:
provider reachability、credential 参照不能、model / execution mode mismatch を Job Setup validation の必須失敗条件にするかを決めてください。

やりたいこと:
AI runtime と実行方式を選んだ時点で、create job してよい状態かを判断する。

背景:
source は AI 基盤選択と API key 保存を要求しているが、外部 provider へ接続できるかを setup で検査するかは固定していない。paid API 非依存検証も必要である。

選択肢:
1. credential 解決と provider capability だけを blocking にし、実ネットワーク到達性は phase 実行時に扱う
2. credential 解決、provider capability、ネットワーク到達性をすべて blocking にする
3. provider / model 選択だけを保存し、credential と到達性は phase 実行時に扱う
4. その他

AI推奨:
1

推奨理由:
paid API や外部 network に依存せず、create 前に機械的に判定できる範囲を blocking にできる。

不確実性:
運用上、create job 時点で実接続確認まで済ませたい可能性がある。

回答形式:
選択肢番号を選んでください。
4 の場合は、採用したい業務ルールを1〜3文で記入してください。
2

## [Q-003] 共通基盤参照の固定方式

質問:
Job Setup で選んだ共通辞書と共通ペルソナを、job 作成時に snapshot、参照 bind、lock のどれで扱うかを決めてください。

やりたいこと:
validation pass 後から create / 実行までの間に、共通基盤が変わった時の一貫性を保つ。

背景:
共通辞書と共通ペルソナは job setup では管理しない。参照だけにするか、作成時点の断面を保存するかで、再 validation 条件と監査表示が変わる。

選択肢:
1. job 作成時に参照 ID と検証断面を bind し、参照不能や更新検知で再 validation 必須にする
2. job 作成時に必要な共通基盤を snapshot し、以後の共通側更新から切り離す
3. job 作成から実行完了まで共通基盤を lock し、削除や更新を禁止する
4. その他

AI推奨:
1

推奨理由:
共通管理 UI と job setup の責務を分けたまま、古い validation pass を誤用しない条件を作れる。

不確実性:
翻訳品質の再現性を最優先するなら snapshot の方が明確になる可能性がある。

回答形式:
選択肢番号を選んでください。
4 の場合は、採用したい業務ルールを1〜3文で記入してください。
4 phase実行時は使ってる最中なのでlockされるべきだが，今回のプランであるジョブ作成時はこの考慮自体必要ないと思う

## [Q-004] 入力 cache 欠落時の扱い

質問:
Job Setup 対象の入力 cache が欠落している場合、自動再構築するか、Job Setup をブロックするかを決めてください。

やりたいこと:
抽出 JSON 正本から再構築可能な入力データを、job setup で安全に扱う。

背景:
恒久仕様は抽出 JSON 正本からの再構築を要求している。一方で、Job Setup が自動再構築まで担うかは固定されていない。

選択肢:
1. 抽出 JSON 正本があれば Job Setup 内で自動再構築してから validation する
2. cache 欠落時は Job Setup をブロックし、Input Review の再構築導線へ戻す
3. cache 欠落時は warning にし、create 時に再構築を試みる
4. その他

AI推奨:
2

推奨理由:
入力取り込みと再構築の責務を Input Review 側に残せるため、Job Setup の責務が広がりにくい。

不確実性:
利用者体験を優先するなら、正本がある場合の自動再構築が望まれる可能性がある。

回答形式:
選択肢番号を選んでください。
4 の場合は、採用したい業務ルールを1〜3文で記入してください。
2

## [Q-005] terminal job 後の再作成

質問:
Completed、Canceled、Failed 済み job がある入力データで、再度 create job を許可するかを決めてください。

やりたいこと:
`1 input = 1 job` の不変条件と、やり直し作業の両方を破綻させない。

背景:
source は 1 入力データに対して 1 翻訳ジョブを要求しているが、終了済み job を履歴として残したまま新規 job を作るかは未固定である。

選択肢:
1. 状態に関係なく同一入力への 2 件目 job 作成を禁止する
2. Completed / Canceled / Failed の terminal job がある場合だけ、新しい job 作成を許可する
3. 同一入力の再作成は旧 job を上書きまたは再利用し、新規 job ID は作らない
4. その他

AI推奨:
1

推奨理由:
現在の恒久仕様と ER の `TRANSLATION_JOB -> X_EDIT_EXTRACTED_DATA` 1:1 を最も保ちやすい。

不確実性:
実運用では同じ入力を別設定で翻訳し直したい需要があり、禁止が強すぎる可能性がある。

回答形式:
選択肢番号を選んでください。
4 の場合は、採用したい業務ルールを1〜3文で記入してください。
1でいいが，過去のジョブを廃棄できる手段が必要だ。

## [Q-006] Draft の保存単位

質問:
Job Setup の `Draft` を永続化済み job 状態として扱うか、UI 上の未保存 setup 状態として扱うかを決めてください。

やりたいこと:
validation failure 中や作成前の設定を、状態遷移と UI 保存のどちらで表すかを固定する。

背景:
spec は `Draft -> Ready` を状態遷移として持つが、Job Setup の draft が DB に保存されるかは未固定である。ここを決めないと `Draft` の観測点がぶれる。

選択肢:
1. Draft は UI 未保存状態とし、DB job は validation pass 後の create で初めて作る
2. Draft を永続化状態として保存し、validation や設定変更を Draft job に更新する
3. Draft は session 内だけ保存し、アプリ再起動では破棄する
4. その他

AI推奨:
1

推奨理由:
partial job や未確定設定の永続化を避けられ、create 成功時だけ job を作る境界が明確になる。

不確実性:
長い setup 作業を中断・再開したい場合は、Draft 永続化が必要になる可能性がある。

回答形式:
選択肢番号を選んでください。
4 の場合は、採用したい業務ルールを1〜3文で記入してください。
1

## [Q-007] Ready job の再編集可否

質問:
作成済み `Ready` job の入力、基盤参照、AI runtime、実行方式を、実行前に再編集できるかを決めてください。

やりたいこと:
create 後から phase 実行前までの設定見直しと、監査可能な設定断面を両立する。

背景:
lifecycle candidate は Ready job の再表示を要求しているが、再編集を許可するかは未固定である。許可する場合、validation 結果の失効と更新履歴も必要になる。

選択肢:
1. Ready job は再表示だけ許可し、設定変更は新規作成または cancel 後の別 workflow にする
2. Ready job の実行前再編集を許可し、変更後は validation を失効させる
3. input は固定し、基盤参照と runtime だけ再編集を許可する
4. その他

AI推奨:
1

推奨理由:
作成時の validation 断面と Ready job の監査内容がずれにくい。

不確実性:
利用者は create 後に provider や辞書を見直したい可能性がある。

回答形式:
選択肢番号を選んでください。
4 の場合は、採用したい業務ルールを1〜3文で記入してください。
1

## [Q-008] validation 履歴の保持粒度

質問:
validation の履歴や失敗理由を、どの粒度で保持するかを決めてください。

やりたいこと:
作成前後に、なぜ create 可能または不可だったかを再確認できるようにする。

背景:
operation-audit candidate は validation failure 理由と設定断面の再確認を要求している。一方で attempt table は持たない判断があり、保存しすぎると ER と衝突する。

選択肢:
1. 直近 validation の結果、対象設定断面、失敗カテゴリだけを保持する
2. validation 実行ごとの履歴を複数件保持する
3. validation 結果は UI 表示だけにし、job 作成時の pass 断面だけを保存する
4. その他

AI推奨:
1

推奨理由:
attempt table なしでも、create 可否と監査に必要な最小断面を残せる。

不確実性:
運用監査や不具合再現を重視するなら、複数履歴が必要になる可能性がある。

回答形式:
選択肢番号を選んでください。
4 の場合は、採用したい業務ルールを1〜3文で記入してください。
1。
ただし validation 実行ごとの履歴は business table ではなく structured log に残す。
アプリ状態として保持するのは、直近 validation 結果、対象設定断面、失敗カテゴリ、job 作成時の pass 断面だけにする。
