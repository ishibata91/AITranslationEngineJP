# Frontend test harness

- test helper は DOM と公開 callback から観測できる結果を支援する。
- production の状態遷移、gateway、表示変換を test helper 内へ複製しない。
- clock と非同期処理は fake timer と明示的な完了待ちで固定する。
- test ごとの DOM と mock を後片付けし、実行順へ依存させない。
