# Core Beliefs

関連文書: [`index.md`](./index.md), [`spec.md`](./spec.md), [`architecture.md`](./architecture.md), [`tech-selection.md`](./tech-selection.md)

本プロジェクトは `agent-first` を強く採用するが、人間の責務は消えない。
人間は主に方針、受け入れ条件、境界、記録、ハーネス改善を担う。
作業方法と役割契約の正本は `.codex/` に置く。

## 1. 基本原則

- repo はエージェントが読める構造を優先する
- `AGENTS.md` は短い地図として保ち、作業方法は `.codex/`、プロダクト判断は `docs/` に置く
- `docs/` は説明資料ではなく、判断と制約の正本として扱う
- 仕様変更を伴う実装は、コード変更と同時に文書も更新する
- 非自明な変更は、実装前に短い plan を残す
- 暗黙知より機械検証を優先する
- 品質は review の往復回数ではなく、plan / checks / evidence / harness の強さで担保する
- review は single-pass で `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足` を見る
- desktop 基盤を作り直す時でも、プロダクト要件と用語は安易に揺らさない
- 過去実装の directory やファイルは、現行文書より優先されない

## 2. 記録システムの扱い

- 永続要件は `spec.md` に記録する
- 内部境界と依存方向は `architecture.md` に記録する
- 実装技術の選択は `tech-selection.md` に記録する
- 実装規約は `coding-guidelines.md` に記録する
- lint と static checks の責務分担は `lint-policy.md` に記録する
- データモデルと ER は `er.md` に記録する
- 詳細な振る舞いと制約は対応する tests / acceptance checks / validation commands に記録する
- 外部仕様と参照方針は `references/` に記録する
- 作業フローは `.codex/skills/` に記録する
- 役割契約は `.codex/agents/` に記録する
- 一時的な作業単位は `exec-plans/` に記録する

## 3. ルール化するべき失敗

次のものは、見つけた時点でユーザーに報告する。

- 同じ前提説明を毎回要求する曖昧な文書構成
- 参照先が存在しない文書リンク
- 用語集と異なる名称の使用
- 同じ責務を複数文書で別定義している状態
- エージェントが検証入口を発見できない状態
- role handoff が曖昧で、lane と補助 skill の責務が崩れる状態
- 生成物と hand-written code の境界が曖昧で、どこを編集すべきか読めない状態
- frontend から desktop transport を直接呼び、gateway を迂回する状態
- review で繰り返し検出される契約違反が harness や executable specs へ昇格されない状態
