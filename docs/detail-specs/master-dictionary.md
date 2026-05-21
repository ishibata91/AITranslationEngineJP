# 詳細仕様: マスター辞書

- `detail_spec_id`: `master-dictionary`
- `status`: `approved`
- `source_artifacts`: `docs/exec-plans/active/2026-04-11-master-dictionary-management.md`
- `implementation_artifacts`: `N/A`
- `review_artifacts`: `docs/scenario-tests/master-dictionary-management.md`

## 親要件と仕様

### `master-dictionary-REQ-001` マスター辞書管理へ到達できる

親要件:
利用者はマスター辞書を独立した管理対象として開き、辞書データを管理できる。

仕様:
- 利用者はマスター辞書を他のマスター情報と独立して管理できる。
- マスター辞書管理は一覧参照、検索、詳細確認、新規作成、編集、削除、XML 取り込みを対象にする。
- マスターペルソナ管理はマスター辞書管理の対象外にする。
- 利用者向け説明は、辞書管理の目的と結果を示す情報として扱う。

### `master-dictionary-REQ-002` 辞書エントリを検索して詳細を判断できる

親要件:
利用者はマスター辞書から辞書エントリを探し、詳細を判断できる。

仕様:
- 利用者はマスター辞書の辞書エントリ集合を参照できる。
- 利用者は辞書エントリを検索できる。
- 利用者は選択した辞書エントリの詳細情報を判断できる。
- マスター辞書は数万件規模の辞書レコードでも、一覧参照、検索、選択、詳細確認、編集を継続できる。

### `master-dictionary-REQ-003` 辞書エントリを作成、編集、削除できる

親要件:
利用者はマスター辞書の辞書データを新規作成、編集、削除できる。

仕様:
- 利用者は辞書データを新規作成できる。
- 利用者は選択した辞書データを編集できる。
- 利用者は選択した辞書データを削除できる。
- XML 取り込み、新規作成、編集、削除は、現在状態を判断しながら実行できる。

### `master-dictionary-REQ-004` XML から辞書データを取り込める

親要件:
利用者は XML ファイルを選択し、許可された REC だけから辞書データを取り込める。

仕様:
- 利用者は XML ファイルを選択して辞書データを取り込める。
- XML 取り込みは選択中ファイルを利用者が識別できる状態で開始する。
- XML 取り込みは xTranslator 形式の辞書 XML から単語を抽出できる。
- XML 取り込み時は `BOOK:FULL`, `NPC_:FULL`, `NPC_:SHRT`, `ARMO:FULL`, `WEAP:FULL`, `LCTN:FULL`, `CELL:FULL`, `CONT:FULL`, `MISC:FULL`, `ALCH:FULL`, `FURN:FULL`, `DOOR:FULL`, `RACE:FULL`, `INGR:FULL`, `FLOR:FULL`, `SHOU:FULL` のみを単語抽出対象とする。
- 許可リスト外の REC は抽出対象外として扱う。

### `master-dictionary-REQ-005` 利用者向け情報を仕様用語へ合わせる

親要件:
利用者は日本語で一貫した用語を通じて辞書管理を操作できる。

仕様:
- 利用者向け用語は `docs/spec.md` の用語に合わせ、日本語で一貫させる。
- xTranslator 写像、利用状況、要件説明は、辞書エントリ管理の標準情報とは分離する。

## 根拠

- `docs/exec-plans/active/2026-04-11-master-dictionary-management.md` を元に製本している。
- `docs/scenario-tests/master-dictionary-management.md` は検証観点の参照として残っている。
