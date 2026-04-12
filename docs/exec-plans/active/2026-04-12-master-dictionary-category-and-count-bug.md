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
- Playwright MCP で `http://host.docker.internal:34115` に到達し、Wails frontend の backend 接続済み表示を確認した。
- `Dawnguard_english_japanese.xml` 取り込み後、UI は `完了` と `XML取り込みを一覧と詳細へ反映しました。` を表示した。
- 取り込み結果カードには `700 件`、`一覧件数 740`、`選択状態 Ancient Vampire`、`詳細表示 太古の吸血鬼` が表示された。
- 取り込み前の辞書一覧は `40 件から絞り込みます。` で、取り込み後は `740 件から絞り込みます。` に変化した。
- 取り込み後のカテゴリ候補として観測できたのは `すべて`、`NPC`、`装備`、`地名` だった。
- 取り込み後のカテゴリ候補に `アイテム` は観測できなかった。

## Trace Plan

- Playwright MCP で XML 取込後のカテゴリ候補と件数表示を再現し、console と Wails ログを採取した。
- `distilling-fixes` で import、カテゴリ集計、件数集計の関連コードを絞り、カテゴリ候補が `state.entries` 依存で生成される事実を確認した。
- `tracing-fixes` と追加観測で、取込結果件数表示が `importedCount` と `totalCount` の別指標を混在表示していることを確認した。
- XML 実体には `CONT:FULL` と `BOOK:FULL` が含まれており、少なくとも `アイテム` と `書籍` が候補から欠ける UI は不自然であると判断した。

## Fix Plan

- `reproduce-issues` で画面再現と証跡取得を行った。
- `distilling-fixes` で既知事実、関連コード、関連仕様、open gap を整理した。
- `tracing-fixes` で原因仮説と観測点を決め、追加 logging 不要と判断した。
- frontend の accepted scope として、カテゴリ候補生成を現在ページ依存から外し、取込結果件数の表示ラベルを誤読しにくい形へ修正する。
- frontend phase-6 実装後に UI check、回帰 test、review を順に実施する。

## Acceptance Checks

- `dictionaries/dawnguar_english_japanese.xml` 取り込み後、カテゴリ検索候補が実データに含まれるカテゴリ集合と整合する。
- 取り込み後の件数表示が実際に保存されたデータ件数と整合する。
- マスター辞書画面の主要導線でカテゴリ検索と件数表示の回帰がない。

## Required Evidence

- Playwright MCP による再現結果。
- カテゴリ候補 UI と件数表示 UI の screen capture または同等証跡。
- console と既定では `tmp/logs/wails-dev.log` を使う Wails ログ。
- browser console では `favicon.ico` の 404 以外に import 失敗や runtime exception を観測していない。
- Wails ログでは起動時の binding 生成関連の `Not found: struct ...` が複数回出ているが、今回の import 成功・失敗に直接結びつく専用ログは観測していない。
- Wails ログには `runtime:ready -> Unknown message from front end: runtime:ready` があるが、UI は操作継続可能だった。
- 再取込時の direct bridge response は `importedCount: 0`、`updatedCount: 947`、`skippedCount: 7235`、`lastEntryId: 740`、`page.totalCount: 740` を返した。
- XML 実体に `CONT:FULL` と `BOOK:FULL` が存在することを確認した。

## Closeout Notes

- 必要な残留リスクや follow-up はここに記録する。

## Outcome

- planned
