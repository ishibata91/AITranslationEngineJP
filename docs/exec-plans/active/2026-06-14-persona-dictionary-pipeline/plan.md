# Task Plan: 2026-06-14-persona-dictionary-pipeline

- `workflow`: work
- `status`: planned（未着手。T1 完了後に着手）
- `task_id`: 2026-06-14-persona-dictionary-pipeline
- `task_mode`: 未判定（着手時に preparation-module で確定）
- `request_summary`: engine の翻訳手続きに「ペルソナ指示の注入点」と「辞書解決の差込点」を作り、最小のルール 1 件・辞書 1 件で口調と用語が翻訳に効くことを通す。
- `goal`: 翻訳プロンプトへ、ルールベースのペルソナ指示と辞書解決を差し込む機構を 1 本通す。中身は最小。
- `constraints`: マスター辞書の本格実体（T3）、ルール編集 UI とプロンプト編集（T4）は対象外。
- `close_conditions`: ルール 1 件・辞書 1 件が翻訳プロンプトへ反映される経路が通る。
- `source_branch`: `master`
- `target_branch`: `master`

## Scope（含む / 含まない）

含む:
- engine の翻訳手続きに、ペルソナ指示（口調指示文）をプロンプトへ注入する差込点。
- 辞書解決（原語→訳語）をプロンプトまたは後処理に差し込む差込点。
- 属性 → 翻訳指示のテンプレート変換（最小、ルール 1 件）。`system_requirements.md` §3 のとおり機械的・AI 不使用。
- `er.md` の `speaker` / `race` / `faction` / `voice_type` を読み、話者属性からペルソナ指示を組む最小経路。

含まない:
- マスター辞書の登録・管理・適用の実体（T3）。
- ルール・プロンプトの編集 UI、実プロンプト参照（T4）。
- 会話履歴解析・AI ペルソナ生成（`system_requirements.md` §3 で将来検討）。

## 依存

- T1（`2026-06-14-extract-translate`）の翻訳手続きが存在すること。

## Routing Notes

- `required_reading`:
  - `docs/system_requirements.md`（§3 ペルソナ＝ルールベース、構造化属性＋翻訳ディレクティブ、永続境界）
  - `docs/concept-model.md`（話者・素材・性質の合成で口調が決まる関係）
  - `docs/er.md`（`speaker`/`race`/`faction`/`voice_type`）
  - `docs/architecture.md`（engine の責務）

## Outcome

- 未着手。
