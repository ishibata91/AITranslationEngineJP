# Tasks Catalog

この directory は、粗い実装スライスを順番付きで管理する task catalog を置く。

## 方針

- task は画面単位ではなく、実装向きの大きめなスライスで定義する
- `index.yaml` は既定順を持つ catalog として使う
- 各 task card は `depends_on` だけを持ち、並列可否や owned scope は持たない
- 画面は `related_screens` で紐付ける
- 実装順、担当分割、validation commands は `docs/exec-plans/` の `Implementation Plan` で管理する
- `tasks/` は「何を作るか」の参照用 catalog に留める

## 構成

- `index.yaml`: task catalog の索引と既定順
- `usecases/*.yaml`: 各実装スライスの task card
