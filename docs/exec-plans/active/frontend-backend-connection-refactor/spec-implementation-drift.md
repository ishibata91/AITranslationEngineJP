# frontend-backend-connection-refactor 仕様乖離整理

- 調査 mode: `仕様乖離整理`
- 判断結果: 完了
- 根拠参照:
  - `docs/detail-specs/ai-provider-settings-management.md:47-66`
  - `docs/detail-specs/term-translation-phase.md:25-32`
  - `docs/screen-design/screens/term-translation-phase.md:10-24`
  - `docs/screen-design/screens/provider-settings.md:54-58`
  - `frontend/src/controller/wails/term-translation-phase.gateway.ts:46-97`
  - `frontend/src/controller/wails/translation-job-management.gateway.ts:33-89`
  - `frontend/src/controller/term-translation-phase/term-translation-phase-screen-controller-factory.ts:1-36`
  - `frontend/src/controller/wails/translation-job-management.gateway.test.ts:5-103`
  - `frontend/src/controller/wails/term-translation-phase.gateway.test.ts:5-170`
- 不足情報:
  - なし。task 指定の参照範囲だけで、既存 `DRIFT-FBC-001` から `DRIFT-FBC-004` が要件、詳細仕様、画面仕様に基づく判断対象かどうかを整理できた。
- 次判断材料:
  - 要件、詳細仕様、画面仕様に基づく `人間判断待ち` は残っていない。
  - `DRIFT-FBC-001` から `DRIFT-FBC-003` は後続の `構造品質調査` 候補である。
  - `DRIFT-FBC-004` は後続の `テスト品質調査` 候補である。
- 引き継ぎ先: `designer`

## 人間判断待ち

- なし

## DRIFT-FBC-001

- 判定: `判断対象外`
- 仕様参照:
  - `docs/detail-specs/ai-provider-settings-management.md:47-56`
  - `docs/detail-specs/term-translation-phase.md:25-31`
- 実装参照:
  - `frontend/src/controller/wails/term-translation-phase.gateway.ts:46-97`
  - `frontend/src/controller/wails/translation-job-management.gateway.ts:33-89`
- 差分内容:
  - 実装は `globalThis.go.wails.*` を探索し、個別 controller と `AppController` を候補にして binding を解決している。
  - 既存 `DRIFT-FBC-001` は `generated wailsjs` を正規入口とみなすかどうかの判断を求めていた。
- 判断対象外とした理由:
  - 今回参照を許可された要件、詳細仕様、画面仕様は、利用者が判断できる接続情報、接続状態、秘密値の非表示は定めている。
  - 今回参照を許可された要件、詳細仕様、画面仕様は、`generated wailsjs` を唯一の呼び出し経路にすることや、`globalThis.go.wails.*` 探索を禁止することまでは定めていない。
  - したがって、この論点は仕様実装差分ではなく、接続境界の構造方針として後続で扱う。
- 影響範囲:
  - `frontend/src/controller/wails/*.gateway.ts`
  - `frontend/src/main.ts`
  - `frontend/wailsjs/`
- 後続候補:
  - `構造品質調査`

## DRIFT-FBC-002

- 判定: `判断対象外`
- 仕様参照:
  - `docs/detail-specs/ai-provider-settings-management.md:52-56`
  - `docs/detail-specs/term-translation-phase.md:28-32`
- 実装参照:
  - `frontend/src/controller/wails/term-translation-phase.gateway.ts:83-97`
  - `frontend/src/controller/wails/translation-job-management.gateway.ts:71-89`
- 差分内容:
  - 実装は binding 戻り値を `response as ResponseDto` で型変換して返している。
  - 既存 `DRIFT-FBC-002` は、Wails bridge 戻り値に runtime shape 検証を要求するかどうかの判断を求めていた。
- 判断対象外とした理由:
  - 今回参照を許可された要件、詳細仕様、画面仕様は、利用者へ見せる接続情報の種類と秘密値の扱いを定めている。
  - 今回参照を許可された要件、詳細仕様、画面仕様は、frontend gateway が runtime validation を必須にすることや、型 assertion を禁止することまでは定めていない。
  - したがって、この論点は仕様実装差分ではなく、DTO 境界の実装方針として後続で扱う。
- 影響範囲:
  - `frontend/src/controller/wails/*.gateway.ts`
  - gateway DTO 変換全般
- 後続候補:
  - `構造品質調査`

## DRIFT-FBC-003

- 判定: `判断対象外`
- 仕様参照:
  - `docs/screen-design/screens/term-translation-phase.md:10-24`
  - `docs/detail-specs/term-translation-phase.md:28-32`
- 実装参照:
  - `frontend/src/controller/term-translation-phase/term-translation-phase-screen-controller-factory.ts:1-36`
- 差分内容:
  - screen controller factory は `@controller/wails/gateway-dto/term-translation-phase` から DTO 型を import している。
  - 既存 `DRIFT-FBC-003` は、gateway 境界外での DTO 型参照を許容するかどうかの判断を求めていた。
- 判断対象外とした理由:
  - 今回参照を許可された画面仕様と詳細仕様は、画面が表示する `Gateway 状態`、段階状態、接続情報、秘匿情報の扱いを定めている。
  - 今回参照を許可された画面仕様と詳細仕様は、screen controller factory の import 境界や DTO 型参照の配置までは定めていない。
  - したがって、この論点は仕様実装差分ではなく、依存方向と責務配置の構造方針として後続で扱う。
- 影響範囲:
  - `frontend/src/controller/term-translation-phase/`
  - controller factory と gateway DTO の依存方向
- 後続候補:
  - `構造品質調査`

## DRIFT-FBC-004

- 判定: `判断対象外`
- 仕様参照:
  - `docs/detail-specs/ai-provider-settings-management.md:52-56`
  - `docs/detail-specs/term-translation-phase.md:28-31`
  - `docs/screen-design/screens/provider-settings.md:54-55`
  - `docs/screen-design/screens/term-translation-phase.md:16-24`
- 実装参照:
  - `frontend/src/controller/wails/translation-job-management.gateway.test.ts:5-103`
  - `frontend/src/controller/wails/term-translation-phase.gateway.test.ts:5-170`
- 差分内容:
  - frontend gateway テストは `globalThis.go` を差し替えて binding 呼び出しを観測している。
  - 既存 `DRIFT-FBC-004` は、public seam test をどこに置くか、`generated wailsjs` を観測点に含めるかどうかの判断を求めていた。
- 判断対象外とした理由:
  - 今回参照を許可された要件、詳細仕様、画面仕様は、利用者が画面で判断する接続状態、接続情報、秘匿情報の扱いを定めている。
  - 今回参照を許可された要件、詳細仕様、画面仕様は、接続境界テストの観測点、test helper の形、public seam test の粒度までは定めていない。
  - したがって、この論点は仕様実装差分ではなく、テスト観測点の品質方針として後続で扱う。
- 影響範囲:
  - `frontend/src/controller/wails/*.gateway.test.ts`
  - backend controller の接続境界テスト整理
- 後続候補:
  - `テスト品質調査`

## 残り不足

- 未確認事項:
  - 要件、詳細仕様、画面仕様の参照範囲内では、`DRIFT-FBC-001` から `DRIFT-FBC-004` 以外に新しい仕様実装差分は確認できていない。
- 理由:
  - 今回の成果物は既存 `仕様乖離整理` の修正に限定されており、構造品質調査とテスト品質調査の掘り下げは対象外である。

## 残留リスク

- 接続境界の内部方針に関する論点は残っているが、今回の成果物では人間へ `仕様が正` / `実装が正` の判断を求めない。
- 後続で `構造品質調査` または `テスト品質調査` を行う時は、今回除外した 4 件を再利用できる。

## 推奨 next step

- `refactor-lane` は、この成果物を `人間判断待ち: なし` として扱う。
- 必要なら `DRIFT-FBC-001` から `DRIFT-FBC-003` を `構造品質調査` へ渡す。
- 必要なら `DRIFT-FBC-004` を `テスト品質調査` へ渡す。
