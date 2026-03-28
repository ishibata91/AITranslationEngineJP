---
name: light-review
description: AITranslationEngineJp 専用。workflow-gate だけでは判定しきれない設計論点が残るときに、light flow の差分を補助 review したいときに使う。
---

# Light Review

この skill は軽量フローで Architect が使う補助 review checklist です。
It is a supplemental checklist, not the standard gate.
標準 gate は `workflow-gate` とし、この skill は gate では判定できない設計論点が残る時だけ使います。

## レビュー観点

- short plan 適合性
- 退行リスク
- 検証不足
- 記録更新漏れ

## 入力契約

- short plan
- 実装差分
- 検証結果

## 出力

- score
- severity
- required_delta
- recheck
- docs_sync_needed

## 禁止

- 実装しない
- 好みの改善を主目的にしない
- plan 外の論点を膨らませない
- standard gate の代わりにしない
- `workflow-gate` の代わりとして常用しない
