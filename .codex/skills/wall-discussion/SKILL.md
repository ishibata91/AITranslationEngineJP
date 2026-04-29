---
name: wall-discussion
description: Codex 側の設計壁打ち作業プロトコル。read-only で資料を読み、論点、質問、短い整理を返す基準を提供する。
---
# Wall Discussion

## 目的

`wall-discussion` は作業プロトコルである。
`designer` agent が人間と設計壁打ちをする時に、質問、深読み、事実と推測の分離、短いまとめの作り方を提供する。

実行境界、source of truth、handoff、stop / reroute は [design-bundle](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md) を参照する。

## 対応ロール

- `designer` が使う。
- 呼び出し元は `implement_lane` または人間とする。
- 返却先は human review または `implement_lane` とする。
- owner artifact は `wall-discussion` の出力規約で固定する。

## 入力規約

- 入力は caller から渡された task-local artifact、source_ref、必要な承認状態を含む。
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- エージェント実行定義とツール権限は [designer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.toml) の `allowed_write_paths` / `allowed_commands` とする。
- runtime skill: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md)
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- human の案をそのまま固定せず根拠と反例を確認する
- 資料から読める事実と AI の推測を分ける
- 同じ質問を繰り返さない
- 3〜6 往復ごとに短くまとめる

- 事実と推測を分ける
- human の明示判断だけを confirmed にする
- 次に聞くべき論点を絞る

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

## 完了規約

- 出力規約を満たし、次の actor が再解釈なしで判断できる。
- 不足情報または停止理由がある場合は明示されている。
- 指定資料を読み、事実と推測を分けた。
- 質問を 1〜3 個に絞った。
- confirmed decisions と open questions を分けた。

## 停止規約

- read-only 範囲で成果物を作らない
- 未確認の論点を設計として固定しない
- 実装や docs 正本更新へ進まない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- read-only 範囲で成果物を作らなかった場合は停止する。
- 未確認論点を設計として固定しなかった場合は停止する。
- 同じ質問を繰り返さなかった場合は停止する。
