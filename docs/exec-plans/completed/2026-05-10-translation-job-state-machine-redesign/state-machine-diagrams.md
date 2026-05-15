# 翻訳ジョブ状態機械図

## 目的

既存状態と理想状態を比較する。
この図は task 内の判断補助であり、docs 正本ではない。

## 記法方針

状態遷移図は PlantUML の state diagram 記法に合わせる。
開始と終了は `[*]` で示す。

実装予定図は PlantUML の component diagram 記法に合わせる。
追加予定、変更予定、既存接続先、削除予定を色で分ける。

状態は丸角の箱で示す。
分岐条件は `<<choice>>` の疑似状態で示す。

ジョブ状態とフェーズ実行状態は別の状態機械として扱う。
`TRANSLATION_JOB.state` はジョブ全体のライフサイクルを表す。
`JOB_PHASE_RUN.state` は 1 フェーズ実行のライフサイクルを表す。
大枠の一覧、導線、ジョブ全体の表示は `TRANSLATION_JOB.state` を読む。
各フェーズ画面の操作可否は、対象フェーズの `JOB_PHASE_RUN.state` を読む。

状態の中には、その状態の振る舞いを書く。
`entry / ...` は状態に入った時の作用である。
`do / ...` はその状態にいる間の継続処理である。
`表示要求 / ...` のような内部反応は、状態を変えない処理である。

遷移ラベルは、`要求イベント [判定条件] / 作用` の形で書く。
ページは状態変更を要求しない。
ページは要求イベントだけを送る。

図の注記は、状態そのものではなく判断補助として扱う。

## 図一覧

- `state-machine-existing.puml`: `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` が混ざる問題を示す。
- `state-machine-ideal-draft.puml`: `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` を分け、操作可否を共通規則へ寄せる理想案を示す。
- `state-machine-implementation-components.puml`: `translationjobpolicy` の共通操作規則と phase 別開始前提で実装する予定コンポーネントを示す。
- `state-machine-implementation-sequence.puml`: UseCase だけが `translationjobpolicy` を呼ぶ実行順序を示す。

## 既存状態の見方

既存状態図は、左を `TRANSLATION_JOB.state` として読む。
右を `JOB_PHASE_RUN.state` として読む。

問題は、`ready` ジョブが `pending` フェーズ実行を同時に持てる点である。
この同居により、実行前ジョブなのか開始待ちフェーズなのかが曖昧になる。

ページ要求が両方の状態を直接読ませる点も問題である。
そのため、どちらが操作可否の正本かが曖昧になる。

## 理想状態の見方

理想状態は人間回答を反映した案である。
task 内の人間設計レビュー用に扱う。

ページは状態変更を直接要求しない。
ページは `表示要求`、`開始要求`、`中断要求`、`再開要求`、`再試行要求`、`取り消し要求` を送る。

`translationjobpolicy` は、現在状態、要求イベント、判定条件から、許可、拒否、作用を返す。
`translationjobpolicy` は、共通操作規則を先に評価する。
`start` の時だけ、対象 phase の開始前提を評価する。
JobIOService は、UseCase が確定した状態事実の取得と保存を行う。

大枠画面の要求は `TRANSLATION_JOB.state` に置く。
各フェーズ画面の要求は、対象フェーズの `JOB_PHASE_RUN.state` に置く。
ジョブ全体の状態更新が必要な時だけ、`translationjobpolicy` が `TRANSLATION_JOB.state` への作用を返す。

`ready` ジョブでは、フェーズ実行を事前作成しない。
開始要求が許可された時だけ `JOB_PHASE_RUN` を作る。

retry、resume、pause、cancel の可否は phase type で分けない。
phase type で分ける対象は、開始前提データ、完了判定、呼び出す service method だけである。

## 実装予定コンポーネント図の見方

実装予定コンポーネント図は、`internal/usecase/translationjobpolicy` と `internal/jobio` の責務分離を示す。
`translationjobpolicy` は UseCase 専用のポリシーとして扱う。
`internal/jobio` は状態判定用 snapshot の取得と、確定済み状態事実の保存だけを扱う。

`OperationRuleCatalog` は、retry、resume、pause、cancel、terminal guard などの共通操作規則を持つ。
`PhasePrerequisiteCatalog` は、`start` の時だけ対象 phase の開始前提を評価する。
`ServiceCommandMap` は、許可済み操作に対応する service method を選ぶ。
`PolicyResult` は UseCase 内の一時値であり、DB には保存しない。
UseCase は `PolicyResult` に従って Service を呼び、保存直前に確定済み状態事実へ変換する。
JobIOService は policy の rule 名、判定履歴、`PolicyResult` を保存しない。

廃止対象は、`internal/statemachine`、`StateMachineFacade`、状態 class 群、`pending` phase run の事前作成、operation summary の DB 永続保存、policy 判断結果の DB 永続保存、各 Service 内の状態遷移可否判断である。

## 根拠参照

- PlantUML state diagram 公式記法: `[*]`、複合状態、`<<choice>>`、note、遷移ラベル。
- OMG PSSM: UML State Machine の精密意味論仕様。
- `docs/spec.md`: 翻訳ジョブ状態遷移。
- `docs/er.md`: ジョブ状態を `JOB_PHASE_RUN` 群から集約する方針。
- `docs/architecture.md`: `StateMachine` と `JobIOService` の既存責務。今後は `StateMachine` を `TranslationJobPolicy` に置き換える。
- `scenario-design.md`: 人間回答済みの状態語彙、入出力境界、受け入れ観点。
- `scenario-design.questions.md`: 人間回答。
- `implementation-targets.md`: docs 更新対象、ハーネス / lint 更新対象、廃止対象。

## 回答済み事項

- `Ready` ジョブには `pending` フェーズ実行を事前作成しない。
- `RecoverableFailed -> Ready` は廃止する。
- `Ready` cancel は job-level 操作として残す。
- phase 開始後の cancel は `Paused` の対象フェーズからだけ許可する。
- 保存済み `TRANSLATION_JOB.state` と対象フェーズの `JOB_PHASE_RUN.state` が食い違う時は、表示だけで状態を書き換えず、危険操作を無効化する。
- retry、resume、pause、cancel の可否は phase type で分けない。
- policy の判断結果は永続化しない。
- JobIOService は確定済み状態事実だけを保存する。
