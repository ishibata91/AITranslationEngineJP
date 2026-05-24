# 仕様乖離整理 起動入力: frontend-overall-refactor

## 対象成果物

`仕様乖離整理`

## 満たされた依存対象

- `task 枠`: `docs/exec-plans/active/frontend-overall-refactor/plan.md#task-frame`

## 読むファイル

- `.codex/skills/investigate/SKILL.md`
- `docs/index.md`
- `docs/spec.md`
- `docs/architecture.md`
- `docs/coding-guidelines-frontend.md`
- `docs/coding-guidelines-tests.md`
- `docs/screen-design/README.md`
- `docs/diagrams/frontend/`
- `docs/diagrams/components/frontend/`
- `frontend/src/main.ts`
- `frontend/src/bootstrap/`
- `frontend/src/application/`
- `frontend/src/controller/`
- `frontend/src/ui/`

## 禁止事項

- 仕様と実装のどちらが正しいかを AI だけで決めない。
- プロダクトコードを変更しない。
- プロダクトテストを変更しない。
- docs 正本本文を変更しない。
- 新しい親要件、受け入れ条件、公開契約、永続仕様を作らない。
- `frontend/wailsjs/` を hand-edit 前提にしない。

## 期待する成果物

`docs/exec-plans/active/frontend-overall-refactor/refactor-classification.md` の `仕様乖離整理` を埋める。
各行は仕様参照、実装参照、差分内容、影響範囲、人間判断待ちを持つ。

## 調査粒度

- `View` と `UI Component` の責務差分。
- `ScreenController`、`Frontend UseCase`、`Store`、`Gateway` の依存方向差分。
- `frontend/src/ui/screens/` の画面設計正本との差分。
- `frontend/src/application/` と `frontend/src/controller/` の公開契約差分。
- テスト品質調査へ回すべき疑いの一覧。

## 停止条件

- 仕様参照が見つからない場合は停止する。
- 実装参照が広すぎて 1 回の調査で根拠を保てない場合は、調査単位の分割案を返す。
- 仕様差分ではなく UX 改善や新機能要求になる場合は、`refactor-lane` へ戻す。
