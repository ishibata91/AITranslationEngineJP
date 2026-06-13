# Task Plan: 2026-06-14-master-dictionary

- `workflow`: work
- `status`: planned（未着手。T2 完了後に着手）
- `task_id`: 2026-06-14-master-dictionary
- `task_mode`: 未判定（着手時に preparation-module で確定）
- `request_summary`: Mod 横断マスター辞書（原語 → 確定訳語）のテーブル・登録・適用を実装する。T2 で空けた辞書解決の差込点に実体を入れる。
- `goal`: 同一原語へ常に同一訳語を当て、Mod 横断で一貫性を保つ機構を実装する。
- `constraints`: ルール・プロンプト編集 UI（T4）は対象外。
- `close_conditions`: 抽出した固有名詞が辞書に登録され、同一原語に同一訳語が当たる。
- `source_branch`: `master`
- `target_branch`: `master`

## Scope（含む / 含まない）

含む:
- マスター辞書テーブルの設計と実装。`er.md` のスコープ外だったため、ここで ER に追加する（`er.md` 更新または辞書専用設計）。
- 揃える対象用語の特定（名前付きレコードの固有名詞の機械抽出。`system_requirements.md` §2）。`er.md` の `proper_noun`（訳の単位・重複排除 e1）と接続する。
- 辞書の登録・適用（プロンプト注入または機械置換。方式は着手時に確定）。

含まない:
- ルール・プロンプトの編集 UI（T4）。
- 本文・会話文中にだけ現れる語の拾い上げ（`system_requirements.md` §2 で未確定）。

## 依存

- T2（`2026-06-14-persona-dictionary-pipeline`）の辞書解決の差込点。
- `er.md` の `proper_noun`（固有名の訳の単位）。

## 未確定（着手時に詰める）

- 訳語の供給方式（既訳流用のみか AI 併用か。`system_requirements.md` §2）。
- 辞書の適用方式（プロンプト注入 / 機械置換）と一貫性の検証方式。

## Routing Notes

- `required_reading`:
  - `docs/system_requirements.md`（§2 一貫性＝Mod 横断マスター辞書）
  - `docs/concept-model.md`（固有名・重複排除 e1）
  - `docs/er.md`（`proper_noun`/`set_phrase`/`placement`）
  - `docs/architecture.md`（store / engine の責務）

## Outcome

- 未着手。
