# 詳細仕様: マスター辞書

- `page_name`: `master-dictionary`
- `source_plan`: `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-11-master-dictionary-management.md`
- `related_mock`: `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/mocks/master-dictionary/index.html`
- `related_scenario`: `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/master-dictionary-management.md`

## 要約

- 今回の task はマスター辞書だけを対象にし、マスターペルソナは別 task へ切り分ける。
- マスター辞書ページは独立ページとして扱い、一覧参照、検索、詳細確認、新規作成、編集、削除、XML 取り込みを同一 task の機能要件に含める。
- XML 取り込みはファイル選択 UI から開始し、選択後は同一画面内に選択中ファイル名と取込開始操作を持つ取込バーを表示する。
- XML 取り込みでは `/Users/iorishibata/Repositories/AITranslationEngineJP/dictionaries/Dawnguard_english_japanese.xml` を読み込んで単語を抽出できることと、抽出対象 REC を許可リストに限定することを固定する。
- 詳細上部の `更新` は編集モーダル、`削除` は確認モーダルで扱う。

## 機能要件

- マスター辞書ページへ独立して到達できること。
- マスター辞書一覧を参照できること。
- マスター辞書一覧で辞書エントリを検索できること。
- 一覧で選択した辞書エントリの詳細情報を参照できること。
- 一覧上から辞書データを新規作成できること。
- 詳細上の `更新` から編集モーダルを開き、辞書データを編集できること。
- 詳細上の `削除` から確認モーダルを開き、辞書データを削除できること。
- `XMLから取り込み` からファイル選択 UI を開き、選択後の取込バー経由で辞書データを取り込めること。
- XML 取り込み時は `BOOK:FULL`, `NPC_:FULL`, `NPC_:SHRT`, `ARMO:FULL`, `WEAP:FULL`, `LCTN:FULL`, `CELL:FULL`, `CONT:FULL`, `MISC:FULL`, `ALCH:FULL`, `FURN:FULL`, `DOOR:FULL`, `RACE:FULL`, `INGR:FULL`, `FLOR:FULL`, `SHOU:FULL` のみを単語抽出対象とし、それ以外の REC は抽出しないこと。
- UI 上の文言、ラベル、説明には実装方針が見える表現を持ち込まないこと。

## 非機能要件

- マスター辞書ページは数万件規模の辞書レコードを保持しても、一覧参照、検索、選択、詳細確認、編集導線が破綻せず継続して操作できること。
- XML 取り込み、新規作成、編集、削除の各導線は、同一ページ内で現在状態が把握できること。
- UI 文言は `docs/spec.md` の用語に合わせ、日本語で一貫していること。

## 対象外

- マスターペルソナ画面および `extractData.pas` 由来 JSON を扱う導線。
- 基盤データ以外の翻訳ジョブ画面、設定画面、翻訳成果物画面の詳細導線追加。
- `docs/` 正本の恒久仕様変更、および usecase 完了条件を超える新規業務要件の追加。
- 一覧サマリ表示、xTranslator 写像表示、利用状況表示、要件説明の露出。

## 未確定事項

- 詳細表示の項目粒度は、原文、訳語、由来、最終更新を基準とし、追加項目要否は後続 task で再判定する。
