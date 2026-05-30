---
name: detail-spec-design
description: Codex 側の詳細仕様差分作業プロトコル。親要件、仕様、未決、回答を task 内成果物として固定する。
---
# Detail Spec Design

## 目的

`detail-spec-design` は作業プロトコルである。
`designer` agent が、`docs/detail-specs/` へ反映する前の詳細仕様差分を、親要件、仕様、未決、回答として固定する時に使う。

実行境界、正本、引き継ぎ、停止 / 戻し は [design-bundle](/Users/iorishibata/Repositories/AITranslationEngineJP/.claude/skills/design-bundle/SKILL.md) を参照する。

## 対応ロール

- `designer` が使う。
- 呼び出し元は `implement_lane`、`refactor_lane`、人間のいずれかとする。
- 返却先は 人間レビュー または呼び出し元レーンとする。
- 担当成果物は `detail-spec-diff.md` とする。

## 呼び出し元から渡される情報

- task 内成果物: 呼び出し元から渡された設計対象と既存成果物。
- 根拠参照: 詳細仕様差分の根拠にする要件、既存詳細仕様、画面設計差分、観測事実。
- 承認状態: 呼び出し元が渡す承認済みまたは未承認の状態。

## 作業前に読む正本

- エージェント実行定義と実行境界は [designer.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.claude/agents/designer.md) に従う。
- 要件正本: [spec.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md) とする。
- architecture 正本: [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) とする。
- ER 正本: [er.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/er.md) と [diagrams/er](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/er/) とする。
- 画面正本: [screen-design](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/README.md) とする。
- 詳細仕様正本: [detail-specs](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/README.md) とする。
- 詳細仕様差分雛形: [detail-spec-diff.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.claude/skills/detail-spec-design/assets/detail-spec-diff.md)
- 実行定義 skill: [design-bundle](/Users/iorishibata/Repositories/AITranslationEngineJP/.claude/skills/design-bundle/SKILL.md)
- 外部成果物が不足または衝突する場合は停止し、衝突箇所を返す。

## skill 内の拘束条件

### 詳細仕様の境界

詳細仕様は、親要件を満たすために必要な仕様境界を固定する。
詳細仕様は、後続の設計、実装、テストが再解釈すると仕様差分になる判断を固定する。

### 詳細仕様差分の構成

| 項目 | 拘束する内容 |
| --- | --- |
| 親要件 | `spec.md` に出てくる要件、またはタスクで追加されたと見なせる業務要件だけを書く。 |
| 要件扱い | 親要件が既存要件か、タスクで追加された業務要件かを書く。 |
| 仕様 | 親要件を満たすために必要な仕様だけを書く。 |
| 未決 | 人間判断が必要な仕様境界だけを書く。 |
| 回答 | 未決ごとの人間回答、回答者、回答日、反映先を書く。 |

### 未決と回答の扱い

未決は親要件の近くに置く。
未決を別ファイルの質問票へ分けない。
回答は未決 1 件ごとに書く。
未回答の未決がある場合は、`implementation-scope` を作らない。

### 詳細仕様にしない内容

| 分類 | 内容 |
| --- | --- |
| 実装方式 | データベース、公開接口、転送形式、保存担当、移行手順、層配置、通信方式、状態管理方式。 |
| 画面表現 | 配置、文言、ボタン、補足説明、非活性表示、画面遷移、表示幅ごとの表現。 |
| 検証方式 | テストケース、検証コマンド、検証データ、代替実装、網羅率条件。 |
| 作業運用 | agent 引き継ぎ、作業順序、branch、commit、review 手順、正本化手順。 |
| 一時判断 | task 内試作、探索ログ、暫定回避、未承認の仮説。 |

### 文体

詳細仕様の本文は日本語を基本にする。
英語の固定名は、状態値、AIサービス名、ファイル形式、既存成果物 key、外部仕様の列名だけに限定する。
一般概念としての英語は本文に使わない。
例として `list`、`summary`、`row`、`field`、`status`、`phase run`、`action`、`screen` を説明語にしない。
英語の固定名を残す場合は、日本語で意味または扱いを補う。

### 媒体非依存

詳細仕様は、仕様として成立する条件、利用者またはシステムが判断できる状態、処理結果を固定する。
詳細仕様は、情報が画面、データベース列、公開応答、転送形式、ログのどこに置かれるかを固定しない。
表示項目、一覧項目、要約項目、保存項目の列挙は、画面設計、ER、公開契約、implementation-scope に委ねる。
ただし、xTranslator 互換 XML の列名など外部形式の互換条件は詳細仕様に残せる。

### 否定形制約

詳細仕様は否定形を標準文体にしない。
禁止、除外、拒否が仕様になる場合は、許可範囲、対象範囲、成立条件、拒否結果として書く。

## 担当ロールが判断してよい範囲

- 親要件を先に固定する。
- 親要件は、`spec.md` に出てくる要件、またはタスクで追加されたと見なせる業務要件だけにする。
- 根拠参照に存在しない親要件を作らない。
- タスクで追加されたと判断できる業務要件は、要件扱いを `追加要件` と明示する。
- `追加要件` は、タスク内成果物または人間依頼に根拠がある場合だけ書く。
- 後続の設計、実装、テストが再解釈すると仕様差分になる判断を固定する。
- 詳細仕様にしない内容を仕様へ混ぜない。
- 詳細仕様の本文は日本語を基本にし、英語の一般語を説明用語にしない。
- 情報の置き場所を画面、データベース列、公開応答、転送形式、ログとして暗黙固定しない。
- 対象外と成立条件は標準項目にしない。
- 禁止、除外、拒否は許可範囲、対象範囲、成立条件、拒否結果として書く。
- 誤読されやすい境界だけ、未決として人間へ返す。
- 未決を AI 判断で埋めない。
- 詳細仕様正本へ反映する内容だけを差分に書く。

## skill が扱わない対象

- 実装方針、implementation-scope の承認済み実装範囲、プロダクトテスト実装詳細は扱わない。
- docs 正本本文の直接更新は扱わない。
- 画面設計書正本へ反映する画面別差分は扱わない。
- 未回答の未決を AI 判断で確定しない。

## 返す成果物

- 判断結果: 詳細仕様差分の完了、未完了、停止の判定を返す。
- 根拠参照: 詳細仕様差分の根拠にした要件、既存詳細仕様、画面設計差分、観測事実を返す。
- 不足情報: 詳細仕様差分を固定できない不足項目を返す。
- 変更成果物: 作成または更新した `detail-spec-diff.md` を返す。
- 追加要件: タスクで追加された業務要件として扱った親要件と根拠参照を返す。
- 未決事項: 人間回答が必要な未決を返す。
- 次判断材料: `designer` または呼び出し元レーンが次を判断できる材料を返す。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコード変更の指示を含めない。

## 作業を完了できる条件

- `detail-spec-diff.md` が親要件、仕様、未決、回答を持つ。
- 各親要件に仕様が 1 件以上ある。
- 未決がある場合は、未決ごとに回答欄がある。
- 未回答の未決がない場合だけ、人間レビュー後に `implementation-scope` へ進める。
- `追加要件` がある場合は、`updating-docs` が要件追加を判断できる根拠参照がある。
- 根拠参照として、根拠成果物パス、必要な人間承認記録、実行した検証結果がある。
- 完了判断材料として、呼び出し元レーンが人間レビュー、実装範囲作成、正本化判断を判断できる情報が返っている。

## 作業を止める条件

- 未回答の未決がある状態で `implementation-scope` を作る必要がある場合は停止する。
- 詳細仕様差分を作る根拠参照が不足する場合は停止する。
- 親要件が既存要件か追加要件か判断できない場合は停止する。
- 外部正本同士が衝突し、人間判断なしで解消できない場合は停止する。
- docs 正本本文を未承認で更新する必要がある場合は停止する。
- プロダクト実装が必要な場合は停止する。
- 停止時は不足項目、衝突箇所、戻し先を返す。
