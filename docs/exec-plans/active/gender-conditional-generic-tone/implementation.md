# Implementation: gender-conditional-generic-tone

## 変更したfile

- `tools/extractor/InfoConditionSqliteWriter.cs`: 同一 INFO の性別条件を集合として判定するために変更した。
- `tools/synthetic-fixture/SyntheticEsmBuilder.cs`: 男性だけ，男女両方，性別条件なしの入力を追加した。
- `tools/extractor.Tests/OracleExtractionTests.cs`: SQLite に保存される性別を検証するオラクルテストを追加した。
- `test-oracle/specs.json`: 追加したオラクルテストに対応する仕様を追加した。
- `test-oracle/fixture/Synthetic.esm`: 追加した fixture から再生成した。
- `spec.md`: 各仕様に対応する実テストを記録した。

## 仕様との対応

- R-1-1: 男女の性別条件がある場合に性別を保存しないことを `condition-sex-both-generic` で確認した。
- R-1-2: 男性だけの性別条件から男性を保存することを `condition-sex-male-only` で確認した。
- R-1-3: 女性だけの性別条件から女性を保存することを `condition-sex-from-getissex` で確認した。
- R-1-4: 性別条件が無い場合に性別を保存しないことを `condition-sex-absent` で確認した。

## 検証結果

- `dotnet run --project tools/synthetic-fixture -- --out test-oracle/fixture`: 成功した。
- `dotnet test tools/extractor.Tests`: 44 件が成功した。
- `git diff --check`: 成功した。

## 未確認事項と停止理由

- なし

## 人間の指摘

