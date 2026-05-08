# fake mode のモデル選択 gate 修正

## 状態

- `task_id`: `2026-05-07-model-selection-fake-gate`
- `lane`: `fix-lane`
- `target`: モデル一覧取得とモデル設定カードの接続条件
- `current_artifact`: `xAI 実画面確認済み`
- `source_task`: [2026-05-07-fake-fixed-model-closed-path](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/plan.md)

## 人間観測

- fake mode なら、credential と endpoint の取得可否や内容に関係なく AI model を使えるべきである。
- fake 判定は frontend へ波及させない。
- frontend は `fake-model` 文字列を特別扱いしない。
- frontend は backend が返した通常の model list 結果だけでモデル選択可否を決める。
- fake は実 provider と同じ provider interface を実装する代替 provider として扱う。
- fake mode 判定は DI 構成でだけ使う。service 層へ `test-safe` や `fake` 判定を漏らさない。

## 成果物DAG

- `人間観測記録`: 完了
- `修正前調査`: 完了
- `修正実行入力`: 完了
- `実装証跡`: 完了
- `回帰テスト証跡`: 完了
- `レビュー通過根拠`: 5 観点通過
- `作業レポート入力`: 完了
- `作業計画完了移動`: 未着手

## 境界

- `fix_lane` はプロダクトコードとプロダクトテストを直接変更しない。
- `fake` provider ID を UI、provider catalog、Job Setup provider list に追加しない。
- frontend に fake mode 判定や `fake-model` 固有分岐を追加しない。
- backend の service 層は fake / test-safe を知らない。
- fake mode の正本は infra の provider interface 実装差し替えである。
- 実装が必要な場合は `implementation_implementer` に `implement-integration` 固定で渡す。
