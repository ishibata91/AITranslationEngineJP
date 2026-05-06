# 正本化判断

## 判断

- 仕様変更または仕様追加: あり。
- human 承認済み恒久仕様: 未確認。
- docs 正本本文反映: 未実施。

## 理由

fake provider を user-facing provider ID ではなく、通常 provider interface を DI で差し替える偽物として整理した。
fake mode の model list は通常の `getModels` 契約で `fake-model` 1 件を返す。
これは provider 境界と開発時の課金防止仕様に関わるため、恒久仕様にするか人間判断が必要である。

## 正本化候補

- `docs/detail-specs/ai-provider-settings-management.md`
- `docs/detail-specs/translation-job-setup.md`
- `docs/detail-specs/body-translation-phase.md`

## 停止理由

human 承認済みの恒久仕様であることが未確認である。
そのため `docs_updater` は起動しない。
