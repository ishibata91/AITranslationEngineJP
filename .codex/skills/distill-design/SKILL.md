---
name: distill-design
description: Codex 側の設計用文脈圧縮 skill。必須要件、UI、scenario の入口を整理するための知識を提供する。
---
# Distill Design

## 目的

`distill-design` は、設計 bundle の前提を整理するための知識である。
必須要件、UI、scenario へ渡す情報の見方を提供する。

共通の圧縮粒度、重複除去、facts / inferred / gap の分離は `distill` を参照する。
この skill は設計向けの観点だけを持つ。

## 対応ロール

- `distiller` が使う。
- 返却先は caller または次 agent とする。
- owner artifact は `distill-design` の出力規約で固定する。

## 入力規約

- 入力は caller から渡された task-local artifact、source_ref、必要な承認状態を含む。
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- エージェント実行定義とツール権限は [distiller.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/distiller.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 共通圧縮: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill/SKILL.md)
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- request を設計可能な facts と constraints に落とす
- 実装案ではなく、design bundle 作成に必要な事実だけを残す

- 必須要件、UI、scenario の入口を分ける
- 変更禁止の境界を constraints として残す
- downstream が読む順番を明示する

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

## 完了規約

- 出力規約を満たし、次の actor が再解釈なしで判断できる。
- 不足情報または停止理由がある場合は明示されている。
- 必須要件、UI、scenario の入口が分かれている。
- design bundle 作成に必要な正本 path が残っている。
- constraints と gaps が facts に混ざっていない。

## 停止規約

- 実装案を確定しない
- owned_scope や対象ファイルを確定しない
- UI モックや scenario 本文を作成しない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 実装案を設計前の事実として固定していない場合は停止する。
- owned_scope や対象ファイルを確定していない場合は停止する。
- UI モックや scenario 本文の作成へ進んでいない場合は停止する。
