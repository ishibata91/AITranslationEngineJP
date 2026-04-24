---
name: codex-review-behavior
description: Codex 実装後 review の挙動正しさグループ知識 package。
---

# Codex Review Behavior

## 目的

変更後のコードが PR の目的どおりに振る舞うかを見る。
diff から取得した実コードを、正解の挙動ベクトルにどの程度近いかで score 化する。

## 見るもの

- 正常系の挙動
- 条件分岐と境界値
- 例外系
- 既存挙動との差分
- bug 修正の場合の原因対応

## 見ないもの

- 命名
- 関数分割
- 読みやすさ
- テスト網羅性
- コードスタイル

## 判定

`score > 0.85` を pass とする。
仕様にない入力や不明な期待値は confidence を下げ、score と混同しない。

## Checklist

- [codex-review-behavior-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-behavior/references/checklists/codex-review-behavior-checklist.md) を参照する。
