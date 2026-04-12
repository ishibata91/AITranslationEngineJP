# Fix Plan

- workflow: fix
- status: planned
- lane_owner: orchestrating-fixes
- scope: master-dictionary-category-and-count-bug

## Request Summary

- マスター辞書画面で `dictionaries/dawnguar_english_japanese.xml` を取り込むと、インポート自体は成功するがカテゴリ検索に `装備` と `NPC` しか現れない。
- 実際の XML には `アイテム` 相当を含む他カテゴリも存在している想定で、カテゴリ集計またはカテゴリ抽出に異常が疑われる。
- 画面上のデータ量集計も実データ量と一致していない疑いがある。

## Decision Basis

- 画面起点で再現確認できる不具合であり、`orchestrating-fixes` の required workflow に従って `reproduce-issues` を先行させる。
- カテゴリ検索と件数集計の不整合は frontend 表示だけでなく import 後の backend 集計値やカテゴリ正規化の異常でも起きうるため、再現証跡と関連コードの切り分けが必要である。
- 対象画面は既存の master dictionary 実装範囲に含まれるため、先行 plan と実装成果物を seed area に含める。

## Known Facts

- user 報告では `dictionaries/dawnguar_english_japanese.xml` の取り込み自体は成功している。
- user 報告ではカテゴリ検索に `装備` と `NPC` しか表示されない。
- user 報告では本来は `アイテム` など他カテゴリも存在する。
- user 報告では画面上のデータ量集計も実データ量と食い違っている。
- `playwright MCP` で確認した再現結果があればここに記録する。

## Trace Plan

- まず Playwright MCP で XML 取込後のカテゴリ候補と件数表示を再現し、console と必要なら Wails ログを採取する。
- 次に `distilling-fixes` で import、カテゴリ集計、件数集計の関連コードと仕様の入口を絞る。
- その後 `tracing-fixes` で frontend 側集計異常か backend 側集計異常かを分離する最小観測点を決める。

## Fix Plan

- `reproduce-issues` で画面再現と証跡取得を行う。
- `distilling-fixes` で既知事実、関連コード、関連仕様、open gap を整理する。
- `tracing-fixes` で原因仮説と観測点を決める。
- 必要時のみ logging を挟み、その後に再度 `reproduce-issues` を行う。
- ownership に応じて backend または frontend の phase-6 実装へ handoff する。

## Acceptance Checks

- `dictionaries/dawnguar_english_japanese.xml` 取り込み後、カテゴリ検索候補が実データに含まれるカテゴリ集合と整合する。
- 取り込み後の件数表示が実際に保存されたデータ件数と整合する。
- マスター辞書画面の主要導線でカテゴリ検索と件数表示の回帰がない。

## Required Evidence

- Playwright MCP による再現結果。
- カテゴリ候補 UI と件数表示 UI の screen capture または同等証跡。
- console と既定では `tmp/logs/wails-dev.log` を使う Wails ログ。
- `playwright MCP` の確認結果、console、network、screen capture のうち必要な証跡をここに記録する。

## Closeout Notes

- 必要な残留リスクや follow-up はここに記録する。

## Outcome

- planned
