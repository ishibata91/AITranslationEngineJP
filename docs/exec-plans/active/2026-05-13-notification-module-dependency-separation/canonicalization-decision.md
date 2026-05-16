# 正本化判断

- `task_id`: `2026-05-13-notification-module-dependency-separation`
- `status`: `completed`
- `decision`: `no_additional_docs_update_required`

## 判断

追加の `docs/detail-specs/` 正本反映は不要である。
理由は、通知 module の構造責務は `docs/architecture.md` と `docs/diagrams/backend/backend-architecture.puml` に初期反映済みであり、今回の実装後差分はその範囲内で閉じているためである。

## 根拠

- `docs/architecture.md`: `NotificationSinkPort`、`NotificationDispatcher`、`NotificationPort`、`Runtime adapter` の責務と依存方向を記録済みである。
- `docs/diagrams/backend/backend-architecture.puml`: UseCase / Service から `NotificationSinkPort` へ通知事実を渡し、`NotificationPort` から `Runtime adapter` へ送る構造を記録済みである。
- `rg -n "RuntimeEventPublisher|runtime event|master-dictionary:import|operation summary|NotificationSink|NotificationPort|NotificationDispatcher|RuntimeAdapter" docs/detail-specs docs/architecture.md docs/diagrams/backend/backend-architecture.puml`: detail spec には通知 module の追加反映対象は見つからなかった。
- `rg -n "RuntimeEventPublisher|RuntimeContextProvider|NewWailsMasterDictionaryRuntimeEventPublisher|NewImportProgressEmitter|fakeRuntimeEventPublisher|publishedCompleted" internal .go-arch-lint.yml`: 旧 runtime publisher 名の残存は解消済みである。

## 影響

- `詳細仕様正本反映` は条件不成立として扱う。
- docs 正本本文の追加変更は行わない。
- 今後 `Runner / Worker` が通知入口へ接続する時は、同じ `NotificationSinkPort` 経由を正本構造として扱う。
