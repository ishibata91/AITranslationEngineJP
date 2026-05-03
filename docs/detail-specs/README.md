# Detail Specs

この directory は、上位シナリオごとの詳細仕様正本を置く。
active exec-plan で承認された恒久仕様を close 前に製本し、後続 task の参照起点にする。

## Naming

- 詳細仕様は `docs/detail-specs/<upper-scenario-id>.md` を正本とする
- `<upper-scenario-id>` は、複数の scenario を束ねる利用者またはシステムの大きな目的を表す
- 画面名だけを理由にファイルを作らない
- 個別ユースケースごとにファイルを分けない
- スキーマ移行、DB 移行、基盤移行、cutover は詳細仕様にしない

## Specs

- [`body-translation-phase.md`](./body-translation-phase.md)
- [`master-dictionary.md`](./master-dictionary.md)
- [`persona-generation-phase.md`](./persona-generation-phase.md)
- [`template.md`](./template.md)
- [`term-translation-phase.md`](./term-translation-phase.md)
- [`translation-input-intake.md`](./translation-input-intake.md)
- [`translation-output-artifact.md`](./translation-output-artifact.md)

## Notes

- plan close 前に、承認済みの恒久仕様だけを移す
- 承認済み `scenario-design` の上位シナリオを起点にし、`ui-design`、実装結果、レビュー結果は補強根拠として使う
- UI 契約は恒久仕様だけを移し、実装前の見た目 artifact を正本化しない
- 実装手順や一時判断は active / completed exec-plan に残す
- この詳細仕様は実装と乖離している場合がある。
