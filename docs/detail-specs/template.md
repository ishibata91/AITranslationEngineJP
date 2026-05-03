# 詳細仕様: <上位シナリオ名>

- `upper_scenario_id`: `<upper-scenario-id>`
- `status`: `draft` / `approved`
- `source_plan`: `<docs/exec-plans/completed/<task-id>/plan.md>`
- `scenario_source`: `<scenario-design.md>`
- `ui_source`: `<ui-design.md>` または `N/A`
- `implementation_source`: `<implementation-result または work report>`
- `review_source`: `<reviewback または review summary>`

## 要約

- 上位シナリオの目的を書く。
- 利用者またはシステムが開始する大きな作業単位を書く。
- 下位 scenario を個別ユースケースとして増やさず、仕様本文へ統合する。

## 対象

- 対象利用者または対象システム。
- 開始条件。
- 完了状態。
- 関係する主要データ。

## 仕様

- 主要操作または主要処理。
- 表示または返却する情報。
- 状態と遷移。
- 拒否条件、無効条件、再試行条件。
- 保存条件、非保存条件、露出禁止条件。

## 受け入れ根拠

- 対応する `scenario-design` の scenario ID。
- 対応する human decision ID。
- 対応する検証結果。

## UI 契約由来の恒久仕様

- 表示項目。
- 操作可否。
- 状態差分。
- overflow、アクセシビリティ、実装後確認で恒久仕様に残す制約。

## 対象外

- 別上位シナリオへ渡す範囲。
- 未承認の候補。
- 実装手順。
- task-local HTML モック。
