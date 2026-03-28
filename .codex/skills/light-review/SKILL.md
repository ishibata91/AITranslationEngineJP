---
name: light-review
description: AITranslationEngineJp 専用。light flow の差分を read-only で確認し、plan 適合性、退行、検証不足だけを短く返したいときに使う。
---

# Light Review

この skill は軽量フローで Architect が使うレビュー用 checklist です。
差分、short plan、検証結果を照らし合わせて、最小限の required delta を返します。

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
