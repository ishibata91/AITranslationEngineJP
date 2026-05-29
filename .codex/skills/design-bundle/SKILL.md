---
name: design-bundle
description: Codex 側の設計成果物進行 skill。詳細仕様差分、画面設計差分、implementation-scope を task 内成果物として固定するための正本、判断、引き継ぎを提供する。
---
# Design Bundle

## 目的

`design-bundle` は作業プロトコルである。
`designer` agent と Codex 本体が、詳細仕様差分、画面設計差分、implementation-scope を task 内成果物として固定する時の、人間可読な実行説明の正本として使う。

作業流れの次実行判断、作業計画フォルダ進行管理、人間向け Codex 実装系レーン引き継ぎの返却は呼び出し元レーンが担当する。
プロダクトコードとプロダクトテストは変更しない。

## 対応ロール

- `designer` が使う。
- 呼び出し元は `implement_lane`、`refactor_lane`、人間のいずれかとする。
- 返却先は人間レビューまたは呼び出し元レーンとする。
- 担当成果物は `design-bundle` の出力規約で固定する。

## 呼び出し元から渡される情報

- 呼び出し元: `implement_lane`、`refactor_lane`、人間のいずれかを受け取る。
- 引き継ぎ入力: 呼び出し元が渡す設計対象と根拠参照。
- 依頼要約: 新規実装または機能拡張として扱う依頼内容。
- 作業計画フォルダ: task 内成果物を置く `docs/exec-plans/active/<task-id>/`。
- 設計範囲: 設計成果物として固定する対象範囲。
- 非必須入力: 人間レビュー記録、対象 skill、既存設計成果物、既知不足を受け取る。
- 必須成果物: `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md` と作業計画フォルダを受け取る。
- 文脈独立条件: 引き継ぎ入力だけで作業でき、引き継いでいない会話文脈に依存しない。

## 作業前に読む正本

- エージェント実行定義と実行境界は [designer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.toml) に従う。
- 要件正本: [spec.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md) とする。
- architecture 正本: [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) とする。
- ER 正本: [er.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/er.md) と [diagrams/er](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/er/) とする。
- 画面設計書正本: [screen-design](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/README.md) とする。
- 詳細仕様正本: [detail-specs](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/README.md) とする。
- シナリオテスト正本: [scenario-tests](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/README.md) とする。
- 詳細仕様差分: [detail-spec-design](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/detail-spec-design/SKILL.md) に従う。
- implementation 対象範囲: [implementation-scope](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-scope/SKILL.md) に従う。
- 補助参照: 入力に明示された関連 docs、関連 skill、人間の現在指示。
- 外部成果物が不足または衝突する場合は停止し、衝突箇所を返す。

## skill 内の拘束条件

### 詳細仕様差分完備判定条件

`designer` は呼び出し元レーンから渡された設計対象、根拠参照、docs 正本、既存設計成果物から `detail-spec-diff.md` を作る。
`detail-spec-diff.md` は親要件、仕様、未決、回答を持つ。
仕様は親要件を満たすために必要な内容だけにし、責務配置、実装手順、テスト手順を混ぜない。

未決は同じ `detail-spec-diff.md` に置く。
未決を別ファイルの質問票へ分けない。
未回答の未決がある場合は、設計完了ではなく人間回答待ちで停止する。

### 画面設計差分条件

画面変更が関係する task では、active plan 内に `screen-design-diff.<screen-id>.md` を人間レビュー前に揃える。
画面設計差分は、`docs/screen-design/screens/` へ適用できる恒久的な画面内容だけを書く。
画面設計差分には、実装指示、テスト手順、agent handoff を書かない。
frontend 実装がある task では、承認済み画面設計差分を implementation-scope と frontend 実装の根拠にする。

### 図成果物条件

設計差分図の標準形は、予定変更箇所だけを示すコンポーネント図とする。
シーケンス図、状態遷移図、その他の図は、ユーザー要求または複雑性がある場合だけ作る。
図は正本ではなく、人間レビューと実装引き継ぎの補助成果物として扱う。

### 実装対象範囲判定条件

implementation-scope を扱う時は、Codex 実装レーンの agent 起動のトークン量を事前計算しない。
代わりに [implementation-scope](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-scope/SKILL.md) の引き継ぎ分割規約と規模判定条件に従い、論理境界と規模の目安で分割する。

各引き継ぎは原則として `1 受け入れユースケース × 1 検証意図` に収める。
Codex 実装系レーンから対象範囲過大で戻された場合は、既存承認を維持せず `pending-human-review` に戻す。

## 担当ロールが判断してよい範囲

- 判断は入力成果物、外部参照規約、対象 agent の責務境界に従う。
- 詳細仕様差分は、親要件を先に固定してから仕様を列挙する。
- 対象外と成立条件は標準項目にしない。
- 人間レビューが必要な判断を AI だけで確定しない。

## skill が扱わない対象

- 作業流れ順序決定、作業計画フォルダ進行管理、作業前確認は扱わない。
- 画面設計根拠の取得を除く実画面 observation、docs 正本化、プロダクト実装は扱わない。
- ツール権限、agent 実行定義、プロダクト仕様正本は変更しない。

## 返す成果物

- 判断結果: 設計成果物の完了、未完了、停止の判定を返す。
- 引き継ぎ先: 呼び出し元レーンを返す。
- 渡す対象範囲: 設計成果物、人間レビュー状態、未回答の未決を返す。
- 対象成果物: 扱った詳細仕様差分、画面設計差分、implementation-scope の状態を返す。
- 変更成果物: 作成または更新した task 内成果物パスを返す。
- 人間レビュー状態: 人間レビューが必要な判断、承認待ち、承認済みの状態を返す。
- 確認結果: 実行した確認と未実行理由を返す。
- 引き継ぎまたは停止理由: 呼び出し元レーンへ戻す理由または停止理由を返す。
- 未決事項: 設計継続に必要な未決事項を返す。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコード変更の指示を含めない。

## 作業を完了できる条件

- task 内成果物が承認状態、根拠参照、未決事項を含んでいる。
- `detail-spec-diff.md` が親要件、仕様、未決、回答を持つ。
- 未回答の未決が 0 件である。
- 人間レビューが必要な判断を AI だけで完了扱いにしていない。
- 必須根拠として、根拠成果物パス、必要な人間承認記録、実行した検証結果がある。
- 完了判断材料として、呼び出し元レーンが次の作業流れ実行判断、人間レビュー、人間向け Codex 実装系レーン引き継ぎを判断できる情報が返っている。
- 残留リスクとして、設計継続に必要な未決事項が返っている。

## 作業を止める条件

- `detail-spec-diff.md` に未回答の未決または未解決競合が残る場合は、人間回答待ちにする。
- 作業流れ順序決定や作業計画フォルダ進行管理が主目的なら呼び出し元レーンへ戻す。
- 作業前の影響範囲、実行計画、検証方法の確認が不足する場合は呼び出し元レーンへ戻す。
- 画面設計根拠の範囲外で実画面 observation が必要なら `investigator` を使う前提で呼び出し元レーンへ戻す。
- docs 正本化が必要なら人間承認後に `docs_updater` を使う前提で呼び出し元レーンへ戻す。
- プロダクト実装が必要なら呼び出し元レーンへ戻し、人間向け Codex 実装系レーン引き継ぎの扱いを判断させる。
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 作業流れの進行管理要求は停止する。
- 不足引き継ぎ入力では停止する。
- プロダクト実装依頼では停止する。
- 未承認 docs 正本化では停止する。
- 実装レーン責務の実装時作業では停止する。
- 人間レビューが必要な判断を AI だけで確定しそうな場合は停止する。
- 作業計画フォルダが不足する場合は停止する。
- 設計範囲が不明な場合は停止する。
- 引き継ぎ入力だけでは作業できない場合は停止する。
