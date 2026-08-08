# UI store

- 複数画面で共有する画面状態だけを置く。
- 一つの screen で完結する状態は対象の container に置く。
- gateway 呼出と backend DTO の変換を置かない。
- state の更新関数は新しい値を返し、呼出順に依存させない。
