# Frontend application

- `application/` には画面から独立した技術機能だけを置く。
- usecase、domain、presenter、gateway contract の層を作らない。
- 画面状態は対象 screen の container か、共有が必要な場合だけ `ui/stores/` に置く。
