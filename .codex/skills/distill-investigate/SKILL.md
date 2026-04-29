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

## 対応ロール

- `distiller` が使う。
- 返却先は caller または次 agent とする。
- owner artifact は `distill-investigate` の出力規約で固定する。

## 入力規約

- 観測計画を立てるための facts を集める時
- 何をまだ知らないかを明示する時
- Codex `investigate` へ渡す画面、API、service、log path を整理する時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- agent runtime と tool policy は [distiller.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/distiller.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 共通圧縮: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill/SKILL.md)
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

### 拘束観点

- 既知の再現手順
- 観測対象の画面、API、service、log path
- 仮説を立てるための最小 code pointer
- 未観測情報と observation target の分離

## 判断規約

- 観測済み事実、未観測対象、仮説候補を分ける
- trace に不要な背景説明は落とす
- log や UI 状態は、証跡 path と再現条件を優先して圧縮する

- 観測済み事実と未観測情報を分ける
- 証跡 path と再現条件を優先する
- observation target を次に調べる入口として返す

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力に tool policy、agent runtime、product code の変更義務を含めない。

## 完了規約

- 出力規約を満たし、次の actor が再解釈なしで判断できる。
- 不足情報または停止理由がある場合は明示されている。
- 観測済み事実、未観測対象、仮説候補が分かれている。
- 証跡 path と再現条件が残っている。
- Codex `investigate` が次に見る observation target が明確である。

## 停止規約

- requirements、UI、scenario の入口を整理する時
- human review 済み `implementation-scope` から実装前 context を作る時
- 深い trace、再現、temporary logging を実施する時
- requirements、UI、scenario そのものの入口整理に使わない
- 深い trace をこの skill の範囲で実施しない
- 憶測ベースの結論を facts に混ぜない
- 停止時は不足項目、衝突箇所、reroute 先を返す。
- 深い trace や再現実施へ進んでいない場合は停止する。
- 憶測ベースの結論を facts に混ぜていない場合は停止する。
- implementation-scope 承認後の実装前整理を扱っていない場合は停止する。
