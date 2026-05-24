# frontend-backend-connection-refactor リファクタ範囲確認

- 記録日: 2026-05-25
- 人間入力: 全部承認
- 判断対象: `structure-quality-investigation.md` と `test-quality-investigation.md` の候補

## 承認済み構造品質候補

| ID | 承認状態 | 実装範囲候補 | 検証要件 |
| --- | --- | --- | --- |
| `SQ-FBC-001` | `承認` | frontend gateway の Wails binding 解決経路を `generated wailsjs` と正規 binding 面に寄せる。 | frontend gateway test と接続境界の検証で回帰を確認する。 |
| `SQ-FBC-002` | `承認` | Wails bridge 戻り値を gateway 内で runtime shape 検証し、無検証の DTO 型変換を縮小する。 | frontend gateway test で成功系、失敗系、検証失敗時の error を確認する。 |
| `SQ-FBC-003` | `承認` | screen controller factory から gateway DTO 型依存を外し、依存方向を application contract 側へ寄せる。 | frontend unit test と typecheck で controller factory の依存境界を確認する。 |

## 承認済みテスト品質候補

| ID | 承認状態 | 実装範囲候補 | 検証要件 |
| --- | --- | --- | --- |
| `TQI-FBC-001` | `承認` | frontend gateway test の観測点を `globalThis.go` 探索順から public seam へ寄せる。 | frontend gateway test で request、response、未接続、検証失敗を確認する。 |
| `TQI-FBC-002` | `承認` | backend controller test を public method 単位へ拡張し、DTO 写像と error wrap の未観測面を埋める。 | backend controller unit test で method ごとの request、response、error wrap を確認する。 |
| `TQI-FBC-003` | `承認` | fake API ではない接続境界専用の scenario test または integration test を追加する。 | frontend と backend の接続境界の最短経路を検証する。 |

## 除外項目

なし。

## 判断保留項目

なし。

## 後続扱い

`designer` は `implementation-scope` で backend、frontend、統合境界、シナリオテスト、単体テストの引き継ぎへ分割する。
`refactor_lane` は承認済み ID だけを実装引き継ぎ入力へ渡す。
