---
name: wall-discussion
description: Codex 側の設計壁打ち知識 package。read-only で資料を読み、論点、質問、短い整理を返す基準を提供する。
---

# Wall Discussion

## 目的

`wall-discussion` は知識 package である。
`designer` agent が人間と設計壁打ちをする時に、質問、深読み、事実と推測の分離、短いまとめの作り方を提供する。

実行権限、agent contract、handoff、stop / reroute は [designer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.agent.md) が持つ。

## いつ参照するか

- human の設計判断を深掘りする時
- 資料を読み、前提、制約、矛盾、未決事項を見つける時
- 早すぎる設計固定を避け、質問を主行動にする時

## 参照しない場合

- 成果物作成、実装、docs 更新が明示された時
- 設計判断を artifact に固定する段階に入った時
- 必要資料が読めず、質問だけでは前進しない時

## 知識範囲

- 1〜3 個の深掘り質問
- current understanding の短い整理
- confirmed decisions と open questions の分離
- risks / tensions の見つけ方

## 原則

- human の案をそのまま固定せず根拠と反例を確認する
- 資料から読める事実と AI の推測を分ける
- 同じ質問を繰り返さない
- 3〜6 往復ごとに短くまとめる

## 標準パターン

1. user が指定した資料を読む。
2. 前提、制約、矛盾、未決事項を分ける。
3. 目的、制約、利用者、失敗条件、検証方法から質問を選ぶ。
4. 質問は 1〜3 個に絞る。
5. 固定が必要になったら `propose-plans` へ進む候補を示す。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は `designer` agent contract に従う。

## DO / DON'T

DO:
- 事実と推測を分ける
- human の明示判断だけを confirmed にする
- 次に聞くべき論点を絞る

DON'T:
- read-only 範囲で成果物を作らない
- 未確認の論点を設計として固定しない
- 実装や docs 正本更新へ進まない

## Checklist

- [wall-discussion-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/wall-discussion/references/checklists/wall-discussion-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は `designer` agent contract が決める。

## References

- agent spec: [designer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.agent.md)
- agent contract: [designer.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/designer/contracts/designer.contract.json)

## Maintenance

- read-only skill として保つ。
- artifact 作成が必要なら `propose-plans` の設計 flow へ分ける。
- long examples は references に分離する。
