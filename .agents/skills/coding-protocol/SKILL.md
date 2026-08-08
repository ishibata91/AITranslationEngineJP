---
name: coding-protocol
description: 実装時に使うスキル。TRIGGER WHEN 実装時　メインエージェント使用禁止。
---
# Coding Protocol

# 実行に必要な情報

- 実装計画
	- ない場合は差し戻す

### 作業前に読み込む資料

```
.agents/skills/references/coding-guidelines-backend.md
.agents/skills/references/coding-guidelines-frontend.md
.agents/skills/references/observability-logging.md
.agents/skills/references/coding-guidelines-tests.md
.agents/skills/references/coding-guidelines.md
```

---

## 実装

読み込んだ規約とフォルダのAGENTS.mdに従って実装を進める。

## テスト

AGENTS.mdがテストを要求する場合み実装する。

## 最終検証

- backend 実装後 `npm run test:backend` を実行する。
- frontend 実装後 `npm run test:frontend` を実行する。
- テスト失敗，lint失敗を残さない。別タスクであろうと，特別な判断なき見逃しは許さない。

---

## 作業を完了できる条件

- テストを実行した
- lintを実行した
- 赤が残っていない

## エスカレーション

- 実装計画と実コードが食い違っている
- ３回以上の試行で解決できない問題が発生した
