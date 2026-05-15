# Human Decision Questionnaire

## Q-TJSM-001 job state の正本

決める仕様:
画面責務ごとに、操作可否と表示状態の正本を決める。

決定済み:
`docs/spec.md` は job state を列挙している。
`docs/er.md` はジョブ状態を `JOB_PHASE_RUN` 群から集約すると定義している。
各フェーズ画面の操作可否は、現在フェーズの `JOB_PHASE_RUN.state` を正本にする。
大枠の一覧、導線、ジョブ全体の表示は、`TRANSLATION_JOB.state` を正本にする。

未確定:
保存時の整合性検査で、`TRANSLATION_JOB.state` と現在フェーズの `JOB_PHASE_RUN.state` が食い違った場合の拒否理由と復旧入口は決まっていない。

選択肢:
1. `JOB_PHASE_RUN` 群からの集約値を操作可否と表示状態の正本にする。
2. `TRANSLATION_JOB.state` を保存正本にし、フェーズ実行集約値は検証補助にする。
3. `TRANSLATION_JOB.state` はジョブ全体のライフサイクルだけに使い、操作可否は `JOB_PHASE_RUN` 群の集約値で決める。
4. 画面責務で正本を分ける。各フェーズ画面の操作可否は `JOB_PHASE_RUN.state` を使い、大枠の一覧、導線、ジョブ全体の表示は `TRANSLATION_JOB.state` を使う。

人間回答:
4 を採用する。理由は、各フェーズ画面ではフェーズ実行の状態が操作可否を決め、大枠ではジョブ全体の状態が導線と表示を決めるためである。

設計への反映:
ページ要求は状態値を直接変更しない。
大枠画面は `TRANSLATION_JOB.state` に要求イベントを送る。
各フェーズ画面は、対象フェーズの `JOB_PHASE_RUN.state` に要求イベントを送る。
ジョブ全体の状態更新が必要な時だけ、`translationjobpolicy` が `TRANSLATION_JOB.state` への作用を返す。

## Q-TJSM-002 Ready job と pending phase run

決める仕様:
`Ready` job が active でない `JOB_PHASE_RUN` を持ってよいかを決める。

決定済み:
`Ready` job は read-only の実行入口である。
Job Run 表示だけでは `Running` へ暗黙遷移しない。

人間回答:
1 を採用する。
`Ready` job には `JOB_PHASE_RUN` を事前作成しない。
`pending` は公開状態として未完了一覧やフェーズ画面へ出さない。

選択肢:
1. `Ready` job には `JOB_PHASE_RUN` を事前作成しない。
2. `Ready` job は `pending` の `JOB_PHASE_RUN` を持ってよいが、操作可否では active と別扱いにする。
3. `pending` を公開状態として扱い、未完了一覧と Job Run に表示する。
4. その他

設計への反映:
作成直後の job は `TRANSLATION_JOB.state=ready` だけを持つ。
開始要求が許可された時だけ、対象フェーズの `JOB_PHASE_RUN` を作る。

## Q-TJSM-003 RecoverableFailed から Ready への再実行準備

決める仕様:
`RecoverableFailed -> Ready` の再実行準備を残すか、同じ `JOB_PHASE_RUN` の retry / resume に統合するかを決める。

決定済み:
phase 詳細仕様は、再送、再開、リトライでは同じ `JOB_PHASE_RUN` を継続すると定義している。
`docs/spec.md` には `RecoverableFailed -> Ready` がある。

人間回答:
`RecoverableFailed -> Ready` は廃止する。
失敗から実行へ戻る遷移は自然だが、失敗から実行前へ戻る遷移は不自然である。
retry と resume は、同じ `JOB_PHASE_RUN` を継続する。

選択肢:
1. `RecoverableFailed -> Ready` を廃止し、retry / resume は同じ `JOB_PHASE_RUN` を継続する。
2. `RecoverableFailed -> Ready` を残し、利用者が再実行準備を選べるようにする。
3. phase 別に、再実行準備を許可する phase と禁止する phase を分ける。
4. その他

設計への反映:
`RecoverableFailed` の再試行要求は、同じ `JOB_PHASE_RUN` を `running` へ戻す。
`Ready` へ戻して新規開始に見せる経路は作らない。


## Q-TJSM-004 cancel の境界

決める仕様:
job-level cancel と phase-level cancel の境界を決める。

決定済み:
`docs/spec.md` は `Ready` と `Paused` から `Canceled` へ進む経路を持つ。
`body-translation-phase.md` は Running から直接 cancel しないと定義している。

人間回答:
1 を採用する。
`Ready` cancel は job-level 操作として残す。
phase 開始後の cancel は、`Paused` の対象フェーズからだけ許可する。

選択肢:
1. `Ready` cancel は job-level 操作として残し、phase 開始後は `Paused` phase からだけ cancel できる。
2. cancel は phase 開始後の `Paused` からだけ許可し、`Ready` job では削除だけを使う。
3. job state が `Paused` なら current phase に関係なく cancel できる。
4. その他

設計への反映:
`Running` から直接 cancel しない。
実行中に取り消す場合は、先に対象フェーズを `Paused` にする。

## Q-TJSM-005 resume と retry の共通規則

決める仕様:
`Paused`、`RecoverableFailed`、retryable failure に対して、resume と retry のどちらを出すかを決める。

決定済み:
body phase は resume を `Paused` の時だけ有効にし、retry を `RecoverableFailed` かつ retryable failure の時だけ有効にする。
persona phase は `Paused` または `RecoverableFailed` で resume を有効にすると定義している。

人間回答:
1 を採用する。
全 phase で resume は `Paused` だけに許可する。
全 phase で retry は `RecoverableFailed` だけに許可する。

選択肢:
1. 全 phase で resume は `Paused` だけ、retry は `RecoverableFailed` だけにそろえる。
2. persona phase だけ `RecoverableFailed` の resume を残す。
3. resume と retry を UI では同じ回復操作として表示し、内部 reason だけ分ける。
4. その他

設計への反映:
中断再開と失敗再試行を別イベントにする。
persona phase だけ `RecoverableFailed` の resume を許可する分岐は作らない。
retry、resume、pause、cancel の可否は phase type で分けない。

## Q-TJSM-006 provider 境界失敗の分類

決める仕様:
credential 未設定、provider timeout、invalid response、correlation error を `RecoverableFailed` と `Failed` のどちらへ分類するかを決める。

決定済み:
provider failure、invalid response、保存失敗、検証失敗は successful `Completed` として扱わない。
secret 実値と provider raw payload は表示、保存、ログに出さない。

人間回答:
4 を採用する。
credential 未設定は、そもそも開始できない状態として扱う。
provider timeout は、ジョブセットアップ時のヘルスチェックで扱うため、`translationjobpolicy` の主分岐には置かない。
invalid response は provider 応答不正として扱う。
correlation error は `RecoverableFailed` として扱う。

選択肢:
1. 設定修正や再送で回復できる失敗は `RecoverableFailed` にし、データ破損や仕様違反だけ `Failed` にする。
2. provider 境界失敗はすべて `RecoverableFailed` にする。
3. invalid response と correlation error は `Failed` にし、timeout と credential 未設定だけ `RecoverableFailed` にする。
4. その他

設計への反映:
開始前に credential が成立しない場合は、phase run を作らず開始拒否にする。
実行中の invalid response は成功として保存しない。
correlation error は再試行可能な失敗として、対象フェーズを `recoverable_failed` にする。


## Q-TJSM-007 operation summary の保存粒度

決める仕様:
状態変更、拒否、削除、再開不可、provider 境界、runtime event 破棄に残す operation summary の粒度を決める。

決定済み:
状態名、必要最小の ID、件数、理由カテゴリは残してよい。
secret、API key、credential 参照実値、endpoint、provider raw payload、prompt 全文、翻訳本文全文、DTO 全体、DB row 全体は残さない。

人間回答:
operation summary を DB に永続保存しない。
operation summary は、必要な時にロジックで導出する。
DB に operation summary が残る設計なら削除対象にする。

選択肢:
1. source file path は永続 summary に残さず、input data ID と短い出自分類だけ残す。
2. source file path の basename だけを残す。
3. source file path を残すが、ユーザー home や credential 相当の path 要素を伏せる。
4. その他

設計への反映:
状態変更の永続値は、状態値、必要な ID、reason category などの正規化された値に限定する。
source file path や prompt 全文のような表示用 summary を保存しない。

## Q-TJSM-008 phase 別 ruleset の要否

決める仕様:
retry、resume、pause、cancel、phase 開始可否を phase type ごとの ruleset に分けるかを決める。

決定済み:
各 phase 詳細仕様は、開始条件と呼び出す service method が異なる。
一方で、retry、resume、pause、cancel の可否は `JOB_PHASE_RUN.state` と要求イベントで同じ形にできる。

人間回答:
phase type ごとの ruleset は作らない。
`translationjobpolicy` は、共通操作規則と phase 別開始前提を分ける。
phase type で分ける対象は、開始前提データ、完了判定、呼び出す service method だけにする。

選択肢:
1. retry、resume、pause、cancel の可否は共通操作規則にし、start の前提だけ phase type ごとに分ける。
2. phase type ごとに `canRetry`、`canResume`、`canPause`、`canCancel` を作る。
3. 全操作を service 内判断へ戻す。
4. その他

設計への反映:
`translationjobpolicy` は、最初に共通操作規則を評価する。
`start` の時だけ、対象 phase の開始前提を評価する。
retry、resume、pause、cancel は phase type を主キーにしない。

## Q-TJSM-009 policy 判断結果の永続化

決める仕様:
`translationjobpolicy` の判断結果、rule 名、判定履歴、`PolicyResult` を DB に永続化するかを決める。

決定済み:
`translationjobpolicy` は DB を読まない pure rule として扱う。
`JobIOService` は状態の取得と保存だけを扱う。
operation summary は DB に永続保存しない。

人間回答:
policy の判断結果は永続化しない。
`PolicyResult` は UseCase 内でだけ消費する一時値にする。
`JobIOService` が保存する対象は、確定済みの job / phase run 状態と、仕様で保存対象にした安全な状態事実だけにする。

選択肢:
1. policy の判断結果を永続化しない。確定済み状態事実だけを保存する。
2. 適用 rule 名だけを監査用に保存する。
3. `PolicyResult` 全体を operation summary として保存する。
4. その他

設計への反映:
`translationjobpolicy` は保存対象を持たない。
`PolicyResult` は DB schema、repository、DTO の永続契約に出さない。
`JobIOService` は policy を永続化せず、UseCase が確定した状態変更だけを保存する。
