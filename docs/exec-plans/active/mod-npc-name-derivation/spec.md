# Spec: mod-npc-name-derivation

`spec.md` はこの task の確定仕様として、要求ごとの仕様を持つ。要求は `plan.md`、設計理由・変更手順・図は `design.md` が持つ。
`design.md` と `spec.md` が食い違う場合は `spec.md` を優先する。

仕様は合成 fixture（`internal/harness/fixture.go` の `SyntheticFixture`）の上で確かめる。実データ（`dictionaries/Data`）の中身に依存させない。R-1 と R-2 は、日本語訳を持たない plugin を fixture へ 1 本足して確かめる。

---

## R-1 Data フォルダ全体の英日対を既訳の辞書にする

- R-1-1（正常系）: 日本語訳を持たない plugin を翻訳対象に選んだ実行で、同じ Data フォルダにある日本語訳を持つ別 plugin の英日対から、横断辞書と既訳の照合表が組まれること。
    - 前提条件: Data フォルダに、日本語訳を持つ plugin と持たない plugin が 1 本ずつある。翻訳対象は日本語訳を持たない側。
    - 確かめ方: 本文へ送る指示の文面に、別 plugin 由来の訳語が機械置換辞書の行として出る。訳された本文で、その原文が日本語になっている。
    - 対応する実テスト: `internal/harness/oracle_test.go` の `TestGoOracles/reference-supply-cross-plugin`
- R-1-2（対象に入る側の境界）: 翻訳対象に選んだ plugin 自身が日本語訳を持つ時も、その plugin の英日対が既訳として集まること。
    - 前提条件: 翻訳対象に、日本語訳を持つ plugin を選ぶ。
    - 確かめ方: 供給の観測ログ `reference_supply_built` の英日対の件数が、その plugin 単体で取れる対の件数以上になる。訳された本文で、その plugin 自身が持つ既訳の語が日本語になっている。
    - 対応する実テスト: `internal/harness/harness_test.go` の `TestReferenceCollectionScansPluginDataFolder`（収集の走査対象が対象 plugin のフォルダであること）、`internal/api/app_test.go` の `TestBuildReferenceScanArgs`（走査を 1 本へ絞らないこと）
- R-1-3（対象に入らない側の境界）: 日本語側を解決できなかった record を、既訳として供給しないこと。
    - 前提条件: Data フォルダのどの plugin も日本語訳を持たない。
    - 確かめ方: 供給の観測ログ `reference_supply_built` の英日対が 0 件になる。本文へ送る指示の文面に機械置換辞書の行が出ない。
    - 対応する実テスト: `tools/extractor.Tests/ReferenceTranslationSqliteWriterTests.cs` の `日本語を解決できた_field_だけを既訳として書く`
- R-1-4（対象に入らない側の境界）: 英日対を集める段を 2 回続けて走らせても、英日対の件数が増えないこと。
    - 前提条件: 同じ Data フォルダで翻訳を 2 回実行する。
    - 確かめ方: 1 回目と 2 回目の観測ログ `reference_supply_built` の英日対の件数が等しい。
    - 対応する実テスト: `tools/extractor.Tests/ReferenceTranslationSqliteWriterTests.cs` の `同じ英日対を二度書いても件数が増えない`

---

## R-2 実行内で確定した NPC 名から部分形の訳を作る

- R-2-1（正常系）: 実行内で氏名の訳が確定した NPC について、本文に名だけで出た人名が訳語へ置き換わること。
    - 前提条件: 日本語訳を持たない plugin に、空白で区切られた 2 語の氏名を持つ NPC がいる。その氏名は固有名として実行内に訳が確定し、訳は中黒で区切られた 2 語のカタカナである。本文の台詞に、その NPC の名だけが出る。
    - 確かめ方: 訳された台詞で、名だけの人名が日本語になっている。
    - 対応する実テスト: `internal/harness/oracle_test.go` の `TestGoOracles/run-name-part-derived`
- R-2-2（対象に入る側の境界）: 二つ名を持つ NPC について、二つ名を除いた名だけが本文に出た時も訳語へ置き換わること。
    - 前提条件: `Grelod the Kind` のように ` the ` を含む氏名の NPC がいて、本文の台詞に `Grelod` だけが出る。
    - 確かめ方: 訳された台詞で、その名が日本語になっている。
    - 対応する実テスト: `internal/harness/oracle_test.go` の `TestGoOracles/run-byname-part-derived`
- R-2-3（対象に入らない側の境界）: 本文で一般語としても使われる語を、人名の部分形として置き換えないこと。
    - 前提条件: NPC の名が、本文の台詞に小文字始まりでも出る語である。
    - 確かめ方: 訳された台詞で、小文字始まりで出た語が原文のまま残る。文頭で大文字になった同じ語も原文のまま残る（破壊置換の形）。
    - 対応する実テスト: `internal/harness/oracle_test.go` の `TestGoOracles/run-common-word-part-skipped`
- R-2-4（対象に入らない側の境界）: 実行内で確定した氏名から作った部分形を、別の plugin を翻訳する実行へ持ち越さないこと。
    - 前提条件: 日本語訳を持たない plugin を翻訳して部分形を作った後、別の plugin を翻訳する。
    - 確かめ方: 2 本目の本文へ送る指示の文面に、1 本目で作った部分形が機械置換辞書の行として出ない。
    - 対応する実テスト: `internal/engine/proper_noun_test.go` の `TestDeriveRunProperNounsStaysInTargetPlugin`
- R-2-5（対象に入らない側の境界）: 横断辞書に既にある原語について、実行内で確定した氏名から別の訳を作らないこと。
    - 前提条件: 実行内で確定した NPC の氏名の一部が、横断辞書に原語として既にある。
    - 確かめ方: 本文へ送る指示の文面で、その原語に対応する機械置換辞書の行が 1 つだけ出て、訳語が横断辞書の側になっている。
    - 対応する実テスト: `internal/engine/proper_noun_test.go` の `TestDeriveRunProperNounsSkipsExistingSources`
- R-2-6（対象に入らない側の境界）: 派生した部分形を、翻訳対象の固有名として数えないこと。
    - 前提条件: 実行内で確定した氏名から部分形が作られている。
    - 確かめ方: 画面の固有名の一覧に部分形の行が出ない。対象 plugin の進捗件数が部分形のぶん増えない。
    - 対応する実テスト: `internal/store/proper_noun_test.go` の `TestDerivedProperNounsAreNotTranslationTargets`

---

## R-3 台詞の口調指示から句点の禁止を外す

仕様を立てず、実テストも置かない。人間の決定による。

変更は翻訳プロンプトへ載る既定値のテキスト 2 件（口調の指示文と役割語の例文）で、どちらも実行前に読める確定した文字列である。文字列そのものを照合するテストは書けるが、要求が求めるのは訳文に句読点が戻ることであり、それを決めるのは翻訳モデルなので、テストで固定できない。変更の中身は `design.md` の変更点が持ち、訳文に句点が戻るかは実 LLM で訳して人間が見て判定する。

**満たさない部分**: 読点が入らない件は原因が確定していないので扱わない。

---

## R-4 話者の性別を口調指示の行として出す

- R-4-1（正常系）: 性別が取れる名指し話者の台詞について、口調指示に性別を示す行が出ること。
    - 前提条件: 話者を解決でき、その話者の性別が取れる台詞。
    - 確かめ方: その台詞へ送る指示の文面に、男性か女性かを示す行が、一人称・言い回しの行とは別に出る。
    - 対応する実テスト: `internal/core/personatone/personatone_test.go` の `TestBuildToneTraitsHasSexLine`
- R-4-2（対象に入る側の境界）: 話者を解決できない汎用台詞とプレイヤーの選択肢についても、性別が取れる時は名指し話者と同じ形の行が出ること。
    - 前提条件: 話者を解決できない台詞、またはプレイヤーの選択肢で、性別が設定されている。
    - 確かめ方: その台詞へ送る指示の文面に出る性別の行が、名指し話者の場合と同じ形になっている。
    - 対応する実テスト: `internal/core/personatone/personatone_test.go` の `TestBuildFreeToneTraitsHasSameSexLine`
- R-4-3（対象に入らない側の境界）: 性別を取れない話者の台詞に、性別を示す行が出ないこと。
    - 前提条件: 話者の性別を取れない台詞、またはプレイヤーの性別が未設定の選択肢。
    - 確かめ方: その台詞へ送る指示の文面に、男性か女性かを示す行が無い。
    - 対応する実テスト: `internal/core/personatone/personatone_test.go` の `TestSexLineAbsentWhenSexUnknown`

---

前提条件に「なし」と書いた仕様が、状況によらず成立させる振る舞いになる。
「対応する実テスト」は設計段階では空にする。`implementation-module` が最終検証で埋め、対応する実テストを置けなかった仕様は停止理由または残課題として人間へ上げる。
