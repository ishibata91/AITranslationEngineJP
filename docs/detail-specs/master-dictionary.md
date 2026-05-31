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
- REC は、レコード種別とフィールド名の組として `RECORD:FIELD` 形式で識別する。
- XML 取り込み対象 REC は、`BOOK:FULL`、`NPC_:FULL`、`NPC_:SHRT`、`ARMO:FULL`、`WEAP:FULL`、`LCTN:FULL`、`CELL:FULL`、`CONT:FULL`、`MISC:FULL`、`ALCH:FULL`、`RACE:FULL`、`INGR:FULL`、`SHOU:FULL` の 13 種別とする。
  - 書籍: `BOOK:FULL`
  - NPC: `NPC_:FULL`、`NPC_:SHRT`、`RACE:FULL`
  - 装備: `ARMO:FULL`、`WEAP:FULL`
  - 地名: `LCTN:FULL`、`CELL:FULL`
  - アイテム: `CONT:FULL`、`MISC:FULL`、`INGR:FULL`、`ALCH:FULL`
  - シャウト: `SHOU:FULL`
- 13 種別の外の REC は、XML 取り込みの対象外として扱う。
- XML 辞書取り込み対象 REC 集合は、単語翻訳フェーズの対象 REC 集合と同一の 13 種別とする。
- XML 辞書取り込みと単語翻訳フェーズは、同一の単語翻訳対象 REC 判定（`IsTermTarget`）を共有して対象判定を行う。
- 既存に保存されている `DOOR:FULL`、`FLOR:FULL`、`FURN:FULL` を含む 13 種別外の master_dictionary レコードは、既存翻訳データ reset の対象として扱い、互換 migration を作らない。

### `master-dictionary-REQ-005` 利用者向け情報を仕様用語へ合わせる

親要件:
利用者は日本語で一貫した用語を通じて辞書管理を操作できる。

仕様:
- 利用者向け用語は `docs/spec.md` の用語に合わせ、日本語で一貫させる。
- xTranslator 写像、利用状況、要件説明は、辞書エントリ管理の標準情報とは分離する。

## 根拠

- `docs/exec-plans/active/2026-04-11-master-dictionary-management.md` を元に製本している。
- `docs/scenario-tests/master-dictionary-management.md` は検証観点の参照として残っている。
- 許可 REC を 16 種別から 13 種別へ縮小し、両集合同一性と `IsTermTarget` 共有を明示する詳細仕様差分は、2026-05-31 の人間設計レビュー承認に基づいて反映済みである（`docs/exec-plans/active/term-target-rec-config/detail-spec-diff.md`）。
