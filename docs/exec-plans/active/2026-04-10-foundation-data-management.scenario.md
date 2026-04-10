# Scenario テスト一覧: foundation-data-management

- task_id: `2026-04-10-foundation-data-management`
- usecase: `tasks/usecases/foundation-data-management.yaml`
- 対象: Foundation Data の観測導線（マスターペルソナ / マスター辞書 / 詳細 / 編集導線 / Rebuild 導線）

## ケース一覧

| ID | 観点 | 事前条件 | 手順 | 期待結果 |
| --- | --- | --- | --- | --- |
| SCN-FDM-001 | app-shell から Foundation Data へ遷移できる | app-shell を表示済み | 1. app-shell の Foundation Data 導線を選択する | 1. Foundation Data 画面が表示される 2. 画面内にコレクション領域、詳細領域、操作領域が表示される |
| SCN-FDM-002 | Persona / Dictionary のコレクション切替ができる | Foundation Data 画面を表示済み | 1. Persona を選択する 2. Dictionary を選択する 3. Persona に戻す | 1. 選択中コレクション表示が切り替わる 2. 一覧領域が選択中コレクションの内容へ更新される |
| SCN-FDM-003 | 選択変更で詳細領域が追従する | 選択中コレクションに2件以上のエントリがある | 1. 一覧の先頭エントリを選択する 2. 別エントリを選択する | 1. 詳細領域が現在の選択エントリに一致して更新される 2. 直前エントリの詳細表示が残留しない |
| SCN-FDM-004 | エントリ未選択時の詳細表示 | 選択中コレクションの一覧を表示済み | 1. 画面初期表示または選択解除状態を確認する | 1. 詳細領域はプレースホルダ表示になる 2. 編集導線は非表示または無効状態である |
| SCN-FDM-005 | 編集導線の露出条件 | 選択中コレクションにエントリがある | 1. エントリを選択する 2. 編集導線を確認する | 1. 選択中エントリに対する編集導線が表示される 2. 編集導線は現在選択中のエントリ文脈に結びつく |
| SCN-FDM-006 | Rebuild 導線の露出条件 | Foundation Data 画面を表示済み | 1. Persona を選択する 2. Rebuild 導線を確認する 3. Dictionary を選択する 4. Rebuild 導線を確認する | 1. 各コレクションで Rebuild 導線が表示される 2. 導線ラベルまたは対象表示が選択中コレクションに一致する |
| SCN-FDM-007 | 基盤データ取得失敗時の観測可能性（主要例外） | Foundation Data 初期ロードで取得失敗を注入する | 1. Foundation Data を開く | 1. 観測不能状態が画面上で識別できる（エラー表示など） 2. 成功時表示（一覧/詳細）と誤認しない |
| SCN-FDM-008 | 責務境界: UI から backend 参照は Wails 境界経由 | Foundation Data 画面のデータロードをトレース可能 | 1. Foundation Data を開く 2. コレクションを切り替える | 1. UI は Wails の公開境界を経由して基盤データを取得する 2. UI が backend 内部構造へ直接依存しない |

## 受け入れ観点への対応

- `app-shell から Foundation Data を開く`: SCN-FDM-001
- `Persona と Dictionary を切り替えて観測する`: SCN-FDM-002
- `選択中の基盤エントリ詳細を確認できる`: SCN-FDM-003, SCN-FDM-004
- `編集導線と Rebuild 導線を確認する`: SCN-FDM-005, SCN-FDM-006
- 主要例外系と責務境界確認: SCN-FDM-007, SCN-FDM-008

## 未確定事項（次工程へ引き継ぎ）

- 詳細領域で表示する具体フィールド名は未確定のため、ケースは「詳細領域の追従」と「選択文脈一致」で固定する
- Foundation Data の取得 API 名称と DTO 形状は未確定のため、境界確認は「Wails 公開境界を経由すること」を観測対象にする
