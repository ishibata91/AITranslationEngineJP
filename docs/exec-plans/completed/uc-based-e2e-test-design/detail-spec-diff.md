# 詳細仕様差分: UC ベース E2E テスト設計

## 親要件

- 利用者は、画面 UC 正本で定義された主要操作を UI 人間操作 E2E として確認できる。
- テスト設計は、既存の詳細仕様正本を変更せず、E2E 観点の根拠として参照する。

## 仕様

- `docs/usecases/uc-*.md` を関連 UC の根拠にする。
- `docs/screen-design/screens/*.md` を対象画面と selector の根拠にする。
- `docs/e2e-test-guidelines.md` の CSV header と selector 方針に従う。
- 前提条件は、各 E2E テストが単独実行に必要な状態として書く。
- 未確定の子要素 selector は、備考に未決として残す。

## 参照する詳細仕様正本

- `docs/detail-specs/ai-provider-settings-management.md`
- `docs/detail-specs/master-dictionary.md`
- `docs/detail-specs/translation-input-intake.md`
- `docs/detail-specs/translation-job-management.md`
- `docs/detail-specs/term-translation-phase.md`
- `docs/detail-specs/persona-generation-phase.md`
- `docs/detail-specs/body-translation-phase.md`
- `docs/detail-specs/translation-output-artifact.md`

## 未決

- 画面設計書で親領域の `data-testid` だけが固定されている操作は、子ボタン、入力欄、行の selector を E2E 実装前に固定する必要がある。

## 回答

- `2026-05-28 user request` を、UC 正本を根拠にした E2E テスト設計の作成承認として扱う。
