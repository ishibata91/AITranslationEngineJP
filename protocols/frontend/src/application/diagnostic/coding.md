# Frontend diagnostic

- frontend の観測は browser console へ出し、backend へ転送しない。
- `event`、`where`、`result` を基本属性とし、必要な場合だけ `id` と `reason` を加える。
- runtime event の破棄、画面から消える非同期状態、外部境界の失敗だけを記録する。
- secret、API key、prompt 全文、翻訳本文全文、backend payload 全体を記録しない。
- `console.log` を production code に残さない。
