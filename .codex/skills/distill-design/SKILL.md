---
name: distill-design
description: Codex 側の設計用文脈圧縮 skill。requirements、UI、scenario、diagram の入口を整理するための知識を提供する。
---

# Distill Design

## 目的

`distill-design` は、設計 bundle の前提を整理するための知識である。
requirements、UI、scenario、diagram へ渡す情報の見方を提供する。

共通の圧縮粒度、重複除去、facts / inferred / gap の分離は `distill` を参照する。
この skill は設計向けの観点だけを持つ。

## いつ参照するか

- request を設計可能な facts と constraints に落とす時
- requirements、UI、scenario、diagram の入口を明示する時
- design bundle 作成前に読むべき正本と不足情報を整理する時

## 参照しない場合

- Codex `investigate` へ渡す観測対象を整理する時
- human review 済み `implementation-scope` から実装前 context を作る時
- 実装案、owned_scope、対象ファイルを確定する時

## 知識範囲

- 既存仕様の正本
- 影響を受ける画面、usecase、service、scenario
- 変更してはいけない既存境界
- downstream が最初に読むべき file / doc

## 原則

- request を設計可能な facts と constraints に落とす
- 実装案ではなく、design bundle 作成に必要な事実だけを残す
- [AGENTS.md](/Users/iorishibata/Repositories/AITranslationEngineJP/AGENTS.md)、[README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/README.md)、関連 skill の重複制約は canonical source へ寄せる

## 標準パターン

1. request と active work plan から設計対象を catalog 化する。
2. requirements、UI、scenario、diagram に関係する入口を分ける。
3. 正本、制約、未確認事項を分離する。
4. 必要な正本だけ `summary` または `full` に上げる。
5. `related_design_pointers` に入れるべき path と理由を整理する。

## DO / DON'T

DO:
- requirements、UI、scenario、diagram の入口を分ける
- 変更禁止の境界を constraints として残す
- downstream が読む順番を明示する

DON'T:
- 実装案を確定しない
- owned_scope や対象ファイルを確定しない
- UI モックや scenario 本文を作成しない

## Checklist

- [distill-design-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill-design/references/checklists/distill-design-checklist.md) を参照する。

## References

- 共通圧縮: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill/SKILL.md)

## Maintenance

- 調査向けの観点は `distill-investigate` に置く。
- 長い例や判断表は [references](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill-design/references/) に分離する。
