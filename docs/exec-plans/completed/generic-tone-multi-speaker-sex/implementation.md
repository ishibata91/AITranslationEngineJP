# Implementation: generic-tone-multi-speaker-sex

## 変更したfile

- `internal/model/persona.go`: 台詞ごとの話者性別集合を表す値を追加した。
- `internal/store/persona_character.go`: 台詞ごとの全話者の性別集合を一括で読む処理を追加した。
- `internal/engine/engine.go`: 複数話者の INFO:NAM1 を汎用口調へ振り分け，条件性別が無い場合に同性集合を使うよう変更した。
- `internal/engine/engine_test.go`: 男性だけ，女性だけ，混在，条件優先，PC 発話のテストを追加した。
- `internal/harness/oracle_test.go`: 男女混在の複数話者が汎用口調へ入る統合オラクルへ更新した。
- `test-oracle/specs.json`: 統合オラクルの期待値を更新した。
- `spec.md`: 仕様と実テストを対応付けた。

## 仕様との対応

- R-1-1: 全男性の複数話者を男性口調へ渡すことを確認した。
- R-1-2: 複数話者の INFO:NAM1 を汎用口調へ渡すことを確認した。
- R-1-3: 全女性の複数話者を女性口調へ渡すことを確認した。
- R-1-4: 男女混在の複数話者へ性別を渡さないことを確認した。
- R-1-5: INFO:RNAM の PC 発話を PC 経路へ残すことを確認した。
- R-1-6: 条件由来の性別を話者性別集合より優先することを確認した。

## 検証結果

- `npm run verify:backend`: 成功した。
- `git diff --check`: 成功した。

## 未確認事項と停止理由

- なし

## 人間の指摘

