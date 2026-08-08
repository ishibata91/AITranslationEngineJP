# R-4 引き継ぎ

## 次のセッションへ渡す指示

次の文を新しいCodexセッションへ渡す。

```text
docs/exec-plans/active/prebuilt-translation-dictionary/r4-handoff.md を読み、スキルを使わずにR-4の残作業を続けてください。辞書MCPを呼ぶ前に「最初に行うこと」を実行してください。
```

## 目的

`plan.md`のR-4を完了する。

R-4の要求は次の4点である。

| 要求 | 現在の状態 |
|---|---|
| `master_term`由来の17,346件をR-3の規則で分類する。 | 完了 |
| 一般辞書に同じ意味と訳がある項目を収録対象外にする。 | 完了 |
| 自動判定だけで削除せず、判断と理由を保存する。 | 完了 |
| 分類、一般辞書の判断、件数確認に必要な操作をMCPへ追加する。 | 実装と自動テストは完了。Codexへ接続するMCPの切り替えだけが未完了。 |

## 作業場所

| 項目 | 値 |
|---|---|
| repository | `/Users/iorishibata/Repositories/AITranslationEngineJP` |
| branch | `codex/prebuilt-translation-dictionary` |
| 2026-08-03時点のHEAD | `e28b5e22` |
| 正規DB path | `dictionary/dictionary.sqlite3` |
| 検分済みDB | `dictionary/dictionary.r4.sqlite3` |
| 変更前DBの退避先 | `dictionary/dictionary.pre-r4.sqlite3` |

作業treeにはR-1からR-4の未commit差分がある。既存差分を戻さない。辞書DBと日本語WordNetのDBは`.gitignore`の対象であり、Gitへ追加しない。

## 最初に行うこと

辞書MCPを最初に呼んではならない。MCPを先に呼ぶと、`dictionary/run-mcp.sh`が正規DBを開いて切り替えを妨げる可能性がある。

最初にshellから次を確認する。

```sh
cd /Users/iorishibata/Repositories/AITranslationEngineJP
git branch --show-current
git status --short
lsof dictionary/dictionary.sqlite3
```

期待するbranchは`codex/prebuilt-translation-dictionary`である。

`lsof`に処理が表示された場合はDBを移動しない。処理を勝手に終了しない。Codexを終了して、辞書MCPの処理が終了した後に再度確認する。

`lsof`に処理が表示されない場合は、検分済みDBを正規DBへ切り替える。

```sh
mv dictionary/dictionary.sqlite3 dictionary/dictionary.pre-r4.sqlite3
mv dictionary/dictionary.r4.sqlite3 dictionary/dictionary.sqlite3
sqlite3 dictionary/dictionary.sqlite3 'PRAGMA integrity_check;'
```

`PRAGMA integrity_check`の期待値は`ok`である。`dictionary/dictionary.pre-r4.sqlite3`が既に存在する場合は上書きせず、作業を止めて実体を確認する。

## 分類結果

検分済みDBには次の結果が保存されている。

| 対象 | 件数 |
|---|---:|
| 英語表記 | 15,917 |
| 意味候補 | 17,346 |
| 使用箇所 | 17,346 |
| 意味へ割り当て済みの使用箇所 | 17,346 |
| 一般辞書との照合結果 | 19,492 |
| 同じ意味の候補として検分した意味候補 | 400 |
| 同じ意味と訳のため収録対象外 | 303 |
| 人間判断が必要 | 97 |
| 未検分の同じ意味の候補 | 0 |

人間判断が必要な97件は収録対象外にしていない。自動分類だけで収録対象外にした項目はない。

## 新しいMCP

新しいMCPは次の11操作を公開する。

| 操作 | 用途 |
|---|---|
| `dictionary_search` | 意味候補を条件付きで検索する。 |
| `dictionary_get` | 意味、使用箇所、一般辞書照合、レビューを取得する。 |
| `dictionary_sense_add` | 意味候補を追加する。 |
| `dictionary_sense_update` | 意味候補を理由付きで更新する。 |
| `dictionary_occurrence_assign` | 使用箇所を意味へ割り当てる。 |
| `dictionary_general_match_update` | 一般辞書との一致状態を理由付きで更新する。 |
| `dictionary_general_match_queue` | 一般辞書の候補を意味単位で取得する。 |
| `dictionary_review_add` | 収録判断と理由を保存する。 |
| `dictionary_classify` | 日本語WordNetとの照合を実行する。 |
| `dictionary_history` | 変更履歴を取得する。 |
| `dictionary_status` | 各状態の件数を取得する。 |

正規DBへの切り替え後に、Codexへ接続された辞書MCPが11操作を公開していることを確認する。旧MCPの5操作だけが表示される場合は再接続できていない。

## MCPで行う最終確認

最初に`dictionary_status`を呼ぶ。期待する主要な値は次のとおりである。

```text
terms: 15917
senses: 17346
occurrences: 17346
assigned_occurrences: 17346
general_matches: 19492
classification_statuses.classified: 400
classification_statuses.general_dictionary_checked: 16946
general_match_statuses.same_mean_and_translation: 303
general_match_statuses.same_mean_candidate: 97
inclusion_decisions.exclude: 303
review_decisions.exclude: 303
review_decisions.needs_human: 97
review_stages.ai_reviewed: 400
```

次に`dictionary_general_match_queue`を次の条件で呼ぶ。

```json
{
  "status": "same_mean_candidate",
  "inclusion_decision": "undecided",
  "review_stage": "unreviewed",
  "after_sense_id": 0,
  "limit": 30
}
```

期待する結果は空の`entries`である。

## DBへ直接行う最終確認

次の照会を実行する。

```sh
sqlite3 -header -column dictionary/dictionary.sqlite3 "SELECT 'unreviewed_same_mean_candidates' AS check_name, COUNT(DISTINCT s.id) AS value FROM dictionary_sense s JOIN general_dictionary_match gm ON gm.sense_id=s.id WHERE gm.match_status='same_mean_candidate' AND s.review_stage='unreviewed' AND s.inclusion_decision='undecided' UNION ALL SELECT 'excluded_without_confirmed_match', COUNT(*) FROM dictionary_sense s WHERE s.inclusion_decision='exclude' AND NOT EXISTS (SELECT 1 FROM general_dictionary_match gm WHERE gm.sense_id=s.id AND gm.match_status='same_mean_and_translation') UNION ALL SELECT 'excluded_without_exclude_review', COUNT(*) FROM dictionary_sense s WHERE s.inclusion_decision='exclude' AND NOT EXISTS (SELECT 1 FROM dictionary_review r WHERE r.sense_id=s.id AND r.decision='exclude') UNION ALL SELECT 'excluded_unreviewed', COUNT(*) FROM dictionary_sense WHERE inclusion_decision='exclude' AND review_stage='unreviewed' UNION ALL SELECT 'empty_match_reasons', COUNT(*) FROM general_dictionary_match WHERE TRIM(reason)='' UNION ALL SELECT 'unclassified', COUNT(*) FROM dictionary_sense WHERE classification_status NOT IN ('classified','general_dictionary_checked');"
```

全行の期待値は0である。

## 検証command

正規DBへの切り替えとMCP確認後に次を再実行する。

```sh
./scripts/go/run.sh test ./dictionary
./scripts/go/run.sh run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2 run ./dictionary/...
npm run verify:backend
git diff --check
```

2026-08-03時点では、辞書単体テスト、辞書単体の静的検査、`npm run verify:backend`、`dictionary/`のLSP診断が通過している。

## 最後に更新する文書

最終確認後に次の2ファイルを更新する。

| file | 更新内容 |
|---|---|
| `r4-classification-results.md` | 旧MCPへ接続中という末尾の記述を、新しいMCPと正規DBで最終確認した記録へ置き換える。 |
| `plan.md` | R-4のMCP再接続を含めて完了したことを確認する。件数は変更しない。 |

R-5以降は本作業の対象外である。R-4の完了確認だけを行う。

## 禁止事項

- スキルを使わない。
- 接続中の辞書MCPを勝手に終了しない。
- DBが開かれている状態でDBを移動または上書きしない。
- `dictionary/dictionary.r4.sqlite3`を検分前のDBで上書きしない。
- `dictionary/*.sqlite3`と`dictionary/reference/wnjpn.db`をGitへ追加しない。
- 303件以外を自動で収録対象外にしない。
- 人間判断が必要な97件を収録対象外へ変更しない。
- R-5以降へ進まない。
- 既存の未commit差分を戻さない。
