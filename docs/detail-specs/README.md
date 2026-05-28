# 詳細仕様

この directory は、目的ごとの詳細仕様正本を置く。
詳細仕様は、親要件と、その親要件を満たすための仕様を正本として残す。
active exec-plan で承認された恒久仕様は、close 前にここへ製本し、後続 task の参照起点にする。

## 命名

- 詳細仕様は `docs/detail-specs/<detail-spec-id>.md` を正本とする
- `<detail-spec-id>` は、利用者またはシステムの大きな目的を表す
- 画面名だけを理由にファイルを作らない
- 個別ユースケースごとにファイルを分けない
- スキーマ移行、データベース移行、基盤移行、切り替え手順は詳細仕様にしない

## 境界

詳細仕様は、親要件を満たすために必要な仕様境界を固定する。
後続の設計、実装、テストが再解釈すると仕様差分になる判断を、詳細仕様の正本にする。

次の内容は詳細仕様にしない。

- 実装方式: データベース、公開接口、転送形式、保存担当、移行手順、層配置、通信方式、状態管理方式
- 画面表現: 配置、文言、ボタン、補足説明、非活性表示、画面遷移、表示幅ごとの表現
- 検証方式: テストケース、検証コマンド、検証データ、代替実装、網羅率条件
- 作業運用: agent 引き継ぎ、作業順序、branch、commit、review 手順、正本化手順
- 一時判断: task 内試作、探索ログ、暫定回避、未承認の仮説

本文は日本語を基本にする。
英語の固定名は、状態値、AIサービス名、ファイル形式、既存成果物 key、外部仕様の列名だけに限定する。
一般概念としての英語は本文に使わない。
例として `list`、`summary`、`row`、`field`、`status`、`phase run`、`action`、`screen` を説明語にしない。
英語の固定名を残す場合は、日本語で意味または扱いを補う。

詳細仕様は、仕様として成立する条件、利用者またはシステムが判断できる状態、処理結果を固定する。
詳細仕様は、情報が画面、データベース列、公開応答、転送形式、ログのどこに置かれるかを固定しない。
表示項目、一覧項目、要約項目、保存項目の列挙は、画面設計、ER、公開契約、implementation-scope に委ねる。
外部形式の互換条件は詳細仕様に残せる。

詳細仕様は否定形を標準文体にしない。
禁止、除外、拒否が仕様になる場合は、許可範囲、対象範囲、成立条件、拒否結果として書く。

## 詳細仕様一覧

- [`ai-provider-settings-management.md`](./ai-provider-settings-management.md)
- [`body-translation-phase.md`](./body-translation-phase.md)
- [`master-dictionary.md`](./master-dictionary.md)
- [`persona-generation-phase.md`](./persona-generation-phase.md)
- [`template.md`](./template.md)
- [`term-translation-phase.md`](./term-translation-phase.md)
- [`translation-job-management.md`](./translation-job-management.md)
- [`translation-input-intake.md`](./translation-input-intake.md)
- [`translation-output-artifact.md`](./translation-output-artifact.md)

## 注意

- plan close 前に、承認済みの恒久仕様だけを移す
- 承認済み `detail-spec-diff.md` を起点にし、画面設計差分、実装結果、レビュー結果は補強根拠として使う
- 画面設計差分は恒久仕様だけを移し、実装前の見た目 artifact を正本化しない
- 実装手順や一時判断は active / completed exec-plan に残す
- 未決事項と回答欄は active plan の `detail-spec-diff.md` に置く
