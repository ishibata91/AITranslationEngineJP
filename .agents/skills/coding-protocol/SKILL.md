---
name: coding-protocol
description: メインエージェント使用禁止。実装時に使うスキル。
---
# Coding Protocol

# 実行に必要な情報

- 実装計画
	- ない場合は差し戻す

### 作業前に読み込む資料

変更対象の言語に対応する資料だけを読む。

- Go: `.agents/skills/references/go.md`
- TypeScript / Svelte: `.agents/skills/references/typescript-svelte.md`
- C#: `.agents/skills/references/csharp.md`

フォルダ固有の責務と実装規約は、変更先に対応する `protocols/` から注入される内容を正とする。

---

## 実装

読み込んだ言語規約と、変更先に対応するプロトコルに従って実装を進める。

## テスト

変更する振る舞いを、変更先のプロトコルが定める観測点から検証する。

## 最終検証

- backend 実装後 `npm run verify:backend` を実行する。
- frontend 実装後 `npm run test:frontend` と `npm run lint:frontend` を実行する。
- `tools/` の C# 実装後 `dotnet test tools/extractor.Tests` を実行する。
- テスト失敗，lint失敗を残さない。別タスクであろうと，特別な判断なき見逃しは許さない。

---

## 作業を完了できる条件

- テストを実行した
- lintを実行した
- 赤が残っていない

## エスカレーション

- 実装計画と実コードが食い違っている
- ３回以上の試行で解決できない問題が発生した
