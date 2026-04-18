---
name: distill-investigate
description: Codex 側の調査用文脈圧縮 skill。Codex investigate へ渡す観測対象と不足情報を整理するための知識を提供する。
---

# Distill Investigate

## 目的

`distill-investigate` は、Codex `investigate` へ渡す前提を整理するための知識である。
観測対象、既知事実、未観測情報、仮説候補の見方を提供する。

共通の圧縮粒度、重複除去、facts / inferred / gap の分離は `distill` を参照する。
この skill は調査向けの観点だけを持つ。

## いつ参照するか

- 観測計画を立てるための facts を集める時
- 何をまだ知らないかを明示する時
- Codex `investigate` へ渡す画面、API、service、log path を整理する時

## 参照しない場合

- requirements、UI、scenario、diagram の入口を整理する時
- human review 済み `implementation-scope` から実装前 context を作る時
- 深い trace、再現、temporary logging を実施する時

## 知識範囲

- 既知の再現手順
- 観測対象の画面、API、service、log path
- 仮説を立てるための最小 code pointer
- 未観測情報と observation target の分離

## 原則

- 観測済み事実、未観測対象、仮説候補を分ける
- trace に不要な背景説明は落とす
- log や UI 状態は、証跡 path と再現条件を優先して圧縮する

## 標準パターン

1. user request と active work plan から観測目的を catalog 化する。
2. 既知の再現手順、観測対象、未観測情報を分ける。
3. 根拠 path がある fact と仮説候補を分離する。
4. 再観測に必要な入口だけを `summary` または `full` に上げる。
5. `observation_targets` に入れるべき対象と理由を整理する。

## DO / DON'T

DO:
- 観測済み事実と未観測情報を分ける
- 証跡 path と再現条件を優先する
- observation target を次に調べる入口として返す

DON'T:
- 設計向けの入口整理に使わない
- 深い trace をこの skill の範囲で実施しない
- 憶測ベースの結論を facts に混ぜない

## Checklist

- [distill-investigate-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill-investigate/references/checklists/distill-investigate-checklist.md) を参照する。

## References

- 共通圧縮: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill/SKILL.md)

## Maintenance

- 設計向けの観点は `distill-design` に置く。
- 長い例や判断表は [references](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill-investigate/references/) に分離する。
