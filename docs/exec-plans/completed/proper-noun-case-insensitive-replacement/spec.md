# Spec: proper-noun-case-insensitive-replacement

`spec.md` はこの task の確定仕様として、要求ごとの仕様を持つ。
要求は `plan.md`、設計理由と変更手順は `design.md` が持つ。
`design.md` と `spec.md` が食い違う場合は `spec.md` を優先する。

## R-1 固有名置換の安定化

- R-1-1（正常系）: 本文から固有名を見つける時は大文字小文字を区別せず、機械置換辞書の同じ訳語へ置き換えること
    - 前提条件: 機械置換辞書に固有名と訳語があり、本文では固有名の先頭と末尾に語境界があり、大文字小文字だけが一部異なる
    - 確かめ方: 結果画面で訳す前の本文と置換内訳を開き、本文の固有名が訳語になり、置換内訳に機械置換辞書の固有名と訳語が表示されることを確認する
    - 対応する実テスト: `internal/core/dictionary.TestDictionaryApplyCaseInsensitiveSpecifications/本文から固有名を見つける時は大文字小文字を区別せず機械置換辞書の同じ訳語へ置き換えること`、`internal/core/mention.TestDetectorCaseInsensitiveSpecifications/本文から固有名を見つける時は大文字小文字を区別しないこと`
- R-1-2（対象に入る側の境界）: 本文の固有名がすべて小文字の場合とすべて大文字の場合も、機械置換辞書の同じ訳語へ置き換えること
    - 前提条件: 機械置換辞書に固有名と訳語があり、本文では固有名の先頭と末尾に語境界があり、固有名がすべて小文字またはすべて大文字である
    - 確かめ方: 結果画面で小文字と大文字の本文をそれぞれ開き、両方の固有名が同じ訳語になっていることを確認する
    - 対応する実テスト: `internal/core/dictionary.TestDictionaryApplyCaseInsensitiveSpecifications/本文の固有名がすべて小文字の場合とすべて大文字の場合も機械置換辞書の同じ訳語へ置き換えること`、`internal/core/mention.TestDetectorCaseInsensitiveSpecifications/本文の固有名がすべて小文字の場合も見つけること`、`internal/core/mention.TestDetectorCaseInsensitiveSpecifications/本文の固有名がすべて大文字の場合も見つけること`
- R-1-3（対象に入らない側の境界）: 機械置換辞書にない固有名の内側で、機械置換辞書の固有名と一致した部分の先頭または末尾に語境界がない場合は置き換えないこと
    - 前提条件: 本文に機械置換辞書にない固有名があり、その内側で機械置換辞書の固有名と大文字小文字だけを除いて一致した部分の先頭または末尾に語境界がない
    - 確かめ方: 結果画面で訳す前の本文を開き、機械置換辞書にない固有名が本文のまま表示されることを確認する
    - 対応する実テスト: `internal/core/dictionary.TestDictionaryApplyCaseInsensitiveSpecifications/機械置換辞書にない固有名の内側で一致した部分の先頭または末尾に語境界がない場合は置き換えないこと`、`internal/core/mention.TestDetectRespectsWordBoundary`
- R-1-4（同じ固有名が複数ある場合）: 本文に大文字小文字が異なる同じ固有名が複数ある場合、置換内訳には機械置換辞書の固有名と訳語を1件表示すること
    - 前提条件: 機械置換辞書に固有名と訳語があり、本文には先頭と末尾に語境界がある同じ固有名が異なる大文字小文字で複数ある
    - 確かめ方: 結果画面で訳す前の本文と置換内訳を開き、本文では複数の固有名が同じ訳語になり、置換内訳には機械置換辞書の固有名と訳語が1件表示されることを確認する
    - 対応する実テスト: `internal/core/dictionary.TestDictionaryApplyCaseInsensitiveSpecifications/本文に大文字小文字が異なる同じ固有名が複数ある場合置換内訳には機械置換辞書の固有名と訳語を1件表示すること`、`internal/core/mention.TestDetectCollapsesCaseVariants`
