# ジョブセットアップ未開始 phase run 修正 人間観測記録

## 対象

- 作業計画フォルダ: `docs/exec-plans/completed/2026-05-10-job-setup-unstarted-phase-run-fix/`
- 修正対象: ジョブセットアップ完了時点の未開始 phase 情報の扱い
- 呼び出し元: 人間

## 人間観測

- ジョブセットアップ時点で `TRANSLATION_JOB` は作成されている。
- ジョブセットアップ時点で未開始 phase の情報が `JOB_PHASE_RUN` に無い。
- 単語翻訳 summary と next phase readiness は phase 情報を必要とする。
- 現在の修正は、`JOB_PHASE_RUN` が無い場合を読み取り側で補う形に見える。
- 根本的には、ジョブセットアップ完了時点で未開始 phase を `JOB_PHASE_RUN` に作る方が自然である。

## 期待との差分

- 期待: ジョブセットアップ完了時点で、job に紐づく未開始 phase が `JOB_PHASE_RUN` として存在する。
- 期待: 未開始 phase は削除、開始、summary、readiness の guard で実行中扱いされない。
- 期待: phase 開始時は、既存の未開始 `JOB_PHASE_RUN` を実行中状態へ遷移させる。
- 実際: 直近の修正では、ready job かつ `JOB_PHASE_RUN` 0 件を読み取り側で許容している。

## 修正対象の仮置き

- `JOB_PHASE_RUN` を「実行開始後の記録」だけではなく、「未開始を含む phase 状態」として扱う。
- 未開始 phase の state 名は、既存状態との衝突を避けて調査で確認する。
- `pending` を未開始として使う場合は、delete guard と start guard が危険扱いしないことを確認する。

## 既存根拠

- 過去の不整合では、`pending` phase run が実行中相当として扱われ、ready job の削除や開始を妨げた。
- 直近の修正では、`JOB_PHASE_RUN` が 0 件でも summary と readiness が読めるようにした。
- 今回は、読み取り側の fallback ではなく、未開始 phase を `JOB_PHASE_RUN` に持つ方向へ戻す。

## 禁止事項

- 新しい中間 table を追加する前提にしない。
- `JOB_PHASE_RUN` と job の間に別の phase plan table を作る前提にしない。
- 未開始 phase を実行中、停止中、再開待ちとして扱わない。
