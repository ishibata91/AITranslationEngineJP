# アイデア評価: dictionary-article-exclusion-and-derived-extraction

## 固定条件

評価対象は `idea.md` の R-1、R-2、R-3、CI-1 だけとする。

R-4 と R-5、製品コード、辞書DB、prompt、snapshot、永続化、画面は変更しない。

砂場は `tmp/dictionary-article-exclusion-and-derived-extraction/` に置いた。

砂場パスは `.gitignore` により追跡対象外であることを確認した。

入力manifestは辞書DBと固定本文DBのSHA-256を記録した。

## 比較実行

固定本文は `tmp/prebuilt-translation-dictionary/r5/source.sqlite3` の87,005本文とした。

比較前照合器は製品の `dictionary.NewDictionary(...).Extract` と同じ語境界、最長一致、本文順を砂場側で実装した。

比較前の全件同値検証は実行時間制約により未完了である。

自動照合件数は次のとおりである。

| 辞書範囲 | 比較前 | CS | CS-CI適格 | CI-1 |
| --- | ---: | ---: | ---: | ---: |
| all-senses | 57,934 | 58,796 | 58,631 | 84,202 |
| included-senses | 47,767 | 48,736 | 48,463 | 54,151 |

## 品質標本

OpenAI Batch APIへ送る本文位置は，評価用本文から安定ハッシュで抽出した150件に限定した。

各層は最大15件とした。

対象層はR-1/R-2追加、CI-1追加、ケース衝突除外による消失である。

各要求は本文全文、照合所有者、候補由来、資格語義の訳・意味・品詞・採否を含む。

モデル名は `gpt-5.6-luna` とした。

Batch IDは `batch_6a786c2223e081909967b4b3a1a6dfa0` とした。

送信直後の状態は `validating` である。

## モデル品質標本の結果

再送信Batch `batch_6a786ef188ec81909c5f771db7ba9613` は150件すべて成功した。

| 辞書範囲 | 層 | 意図した一致 | 誤一致 | 文脈不足 |
| --- | --- | ---: | ---: | ---: |
| all-senses | R-1追加 | 14 | 1 | 0 |
| included-senses | R-1追加 | 15 | 0 | 0 |
| all-senses | CI-1原語追加 | 12 | 3 | 0 |
| included-senses | CI-1原語追加 | 11 | 4 | 0 |
| all-senses | CI-1 R-1別名追加 | 6 | 9 | 0 |
| included-senses | CI-1 R-1別名追加 | 6 | 9 | 0 |
| all-senses | ケース衝突除外による原語消失 | 15 | 0 | 0 |
| included-senses | ケース衝突除外による原語消失 | 15 | 0 | 0 |
| all-senses | ケース衝突除外によるR-1別名消失 | 14 | 0 | 1 |
| included-senses | ケース衝突除外によるR-1別名消失 | 15 | 0 | 0 |

R-2由来の追加一致は固定標本に現れなかった。

## 判定

CI-1は両辞書範囲と両候補由来で誤一致を確認したため、固定コーパス範囲では `棄却` とする。

ケース衝突除外は意図した既存一致を失うため、CI-1の副作用を上回る防御にならない。

R-1はall-sensesで誤一致を確認したため、all-senses条件では `棄却` とする。

R-1はincluded-sensesの固定15標本で誤一致を確認しなかったが、R-3の除外候補層と比較前の全件同値検証が未完了であるため、included-senses条件では `人間判断` とする。

R-2は追加一致の観測がないため、固定コーパス範囲では `棄却` とする。

R-3は除外候補層を評価していないため、`判断不能` とする。

各判定は固定87,005本文と固定150件標本だけに当てはめる。
