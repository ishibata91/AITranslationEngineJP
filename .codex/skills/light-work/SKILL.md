---
name: light-work
description: AITranslationEngineJp 専用。light flow の short plan に従い、最小範囲のコードと文書更新だけを実装したいときに使う。
---

# Light Work

この skill は、軽量フローの実装担当です。
short plan で指定された範囲だけを変更し、必要最小限の検証を返します。

## 入力契約

- short plan
- 変更対象
- 非対象
- 最小検証方法
- required evidence
- docs sync

## Required Reading

- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- 必要なら `docs/executable-specs.md` と該当 plan の `Acceptance Checks`

## 手順

1. short plan を読み、変更対象と非対象を固定する。
2. 必要なファイルだけを読む。
3. 最小差分で実装する。
4. 指定された checks を実行する。
5. 変更点、検証結果、docs 更新有無、残る non-blocking unknown、未解消点を返す。

## 禁止

- plan 外の仕様判断を増やさない
- 無関係な整理や広範囲 refactor を混ぜない
- 他 agent の変更を巻き戻さない
- blocking unknown が出たのにそのまま進めない
