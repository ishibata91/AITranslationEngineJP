# 翻訳フロー移動オーバーホール

## 状態

- `task_id`: `2026-05-08-translation-flow-navigation-overhaul`
- `lane`: `design-bundle`
- `status`: `planned`
- `current_artifact`: `plan`
- `human_review_status`: `pending-after-plan`
- `source`: 人間との壁打ち

## 依頼要約

- 翻訳セクション間の移動を整合性重視で見直す。
- フェーズページへのグローバル直移動を禁止する。
- 新規 job 作成直後は `Job Setup` から単語翻訳ページへ進める。
- 途中再開は未完了 job 一覧から選択して開始する。
- 成果物出力は翻訳管理とは別セクションとして扱う。
- 旧 `Job Run` のセッション取得は廃止する。

## 決定済み

### D-01 旧 `Job Run` の分解

決定: 旧 `Job Run` は大箱として残さず、単語翻訳、NPC ペルソナ生成、本文翻訳の各ページへ分解する。
各フェーズページは、翻訳セクション直下で、入力データページ、`Job Setup` ページ、未完了 job 一覧ページと同列に扱う。

理由: 旧 `Job Run` にフェーズ表示、実行状態、操作、再開導線が集まり、画面の意味が読みにくくなっている。
フェーズをページとして分けると、利用者は現在の翻訳段階と次の移動先を直接理解できる。

影響: `Job Run` 固有の表示単位は廃止する。
実装では、フェーズごとの route または screen state を翻訳セクション配下に持つ。

### D-02 翻訳セクション内の移動制限

決定: 翻訳セクションの初期ページは未完了 job 一覧ページにする。
翻訳セクション内の基本移動は `新規翻訳を開始`、`次へ進む`、`一覧へ戻る`、`出力管理へ移動` に絞る。
フェーズページから `Job Setup` や入力データページへ戻る導線は出さない。

理由: job 作成後の入力と設定は、作成済み job の前提として固定される。
前工程へ戻れる UI にすると、既存 job を編集できるように見え、状態整合性を崩す。

影響: 入力や設定を変えたい場合は、未完了 job 一覧ページから新規翻訳を開始し、既存 job を巻き戻さず新しい job を作る。
途中再開は、未完了 job 一覧から対象 job を選ぶ形に固定する。

### D-03 `sticky footer` の責務

決定: `sticky footer` は、フェーズ変更、次へ進む、一覧へ戻るための移動導線だけを持つ。
実行、一時停止、再開、再試行、取消は、各フェーズページ本文の操作として扱う。

理由: `sticky footer` にフェーズ実行操作まで集めると、移動と実行の責務が混ざる。
footer は画面遷移の判断に集中させ、フェーズの実行操作は対象フェーズの文脈内に置く方が誤操作を避けやすい。

影響: `sticky footer` の移動語彙は全翻訳管理画面で統一する。
`sticky footer` は、既存 sticky footer コンポーネントのエラー表示方式に従い、次へ進めない理由をリアクティブに表示する。

### D-04 直リンク防止

決定: フェーズページへの直リンクは防止する。
直接フェーズページへ入った場合は、未完了 job 一覧へ戻す。

理由: フェーズページは対象 job が確定していることを前提にする。
対象 job が未確定のままフェーズページを表示すると、summary 取得や操作可否の判断が曖昧になる。

影響: Wails 前提では通常 URL 直リンクは起きにくい。
それでも route state や復元状態が不整合な場合は、未完了 job 一覧を復帰先にする。

### D-05 セッション取得の廃止

決定: 旧 `Job Run` のセッション取得ボタンまたはセッション取得処理は廃止する。
フェーズページは、`Job Setup` 完了直後または未完了 job 一覧の選択結果から job を受け取る。

理由: 未完了 job 一覧があるため、利用者が扱う job は一覧選択で固定できる。
セッション取得を別操作として残すと、一覧選択と別の対象取得経路が併存し、責務が重複する。

影響: フェーズページは「job を探す画面」ではなく「選ばれた job を進める画面」になる。
実装では、セッション取得に依存した表示更新を、一覧選択または作成結果からの状態引き継ぎへ置き換える。

### D-06 翻訳完了ページ

決定: 翻訳完了ページは、原文と訳文をページング表示し、出力管理への移動ボタンを持つ。
翻訳完了ページは、出力処理そのものを扱わない。

理由: 翻訳完了直後には、利用者が翻訳結果を確認できる必要がある。
一方で XML 出力、preview、再出力、互換性確認は既存の出力管理の責務である。

影響: 翻訳完了ページは、確認と出力管理への案内に限定する。
出力対象 job を自動選択するかどうかは、出力管理側の仕様として未決に残す。

### D-07 成果物出力の分離

決定: 成果物出力は、翻訳セクションの続きではなく、既存の出力管理として扱う。
完了済み job は翻訳管理一覧には出さず、出力管理側の completed job 一覧で扱う。

理由: 翻訳セクションは job 作成、未完了一覧、フェーズ実行、再開までを扱う。
成果物出力は Completed job を素材に XML などを生成する別作業であり、翻訳実行とは責務が異なる。

影響: 翻訳セクションから成果物出力処理を直接開始しない。
翻訳完了ページのボタンは、既存の出力管理へ利用者を送るための導線に限定する。

## 成果物

- `navigation-state-machine.puml`: 翻訳管理と成果物出力の移動状態を示すステートマシン図。
- `navigation-state-machine.svg`: ステートマシン図の描画確認用 SVG。
- `ui-design.md`: 未作成。人間レビュー後に UI 要件契約として作る。
- `scenario-design.md`: 未作成。移動制約と受け入れ条件を固定する時に作る。
- `implementation-scope.md`: 未作成。人間レビュー後だけ作る。

## 状態遷移図

正本 source は [navigation-state-machine.puml](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-08-translation-flow-navigation-overhaul/navigation-state-machine.puml) とする。
描画結果は [navigation-state-machine.svg](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-08-translation-flow-navigation-overhaul/navigation-state-machine.svg) とする。

図は次の状態を持つ。

- 翻訳セクションの初期ページである未完了 job 一覧ページ、入力データページ、`Job Setup` ページ、出口。
- 翻訳セクション直下の単語翻訳ページ、NPC ペルソナ生成ページ、本文翻訳ページ、翻訳完了ページ。
- 成果物出力の入口、Completed job 一覧ページ、`Output Review` ページ、出口。
- 赤破線で禁止移動を示す。
- 直リンク防止として、直接フェーズページへ入った場合は未完了 job 一覧へ戻す。

## UI 方針

- 翻訳管理の `sticky footer` は、現在地点からの移動操作だけを出す。
- 単語翻訳、NPC ペルソナ生成、本文翻訳は、翻訳セクション直下の同列ページとして表示する。
- 完了済みフェーズは要約だけ表示し、未来フェーズはブロック理由だけ表示する。
- フェーズ実行操作は、各フェーズページ本文に置く。
- フェーズ変更と次へ進む操作は、`sticky footer` に集約する。
- フェーズページ間は `次へ進む` で移動し、見出しクリック移動は作らない。
- `sticky footer` は、次へ進めない理由を既存コンポーネントのエラー表示方式でリアクティブに反映する。
- 旧 `Job Run` のセッション取得ボタンまたはセッション取得処理は廃止する。一覧があり、どの job を選んだかが常に固定になるため。

## 成果物出力の分離

- 成果物出力は翻訳管理の続きではなく、Completed job を素材にする別作業である。
- 成果物出力側は completed job 一覧を持つ。
- 成果物出力側は XML 出力、preview、再出力、互換性確認を扱う。
- 翻訳セクション側は job 作成、未完了一覧、フェーズ実行、再開までを扱う。
- 翻訳完了ページのボタンは、既存の出力管理へ利用者を送るための導線である。

## 根拠参照

- [translation-job-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-management.md): 未完了 job 一覧と `Job Run` 表示対象を扱う。
- [translation-output-artifact.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-output-artifact.md): Completed job から成果物を出力する別仕様を扱う。
- [body-translation-phase.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/body-translation-phase.md): 本文翻訳完了時に job 全体が `Completed` になる。

## 未決事項

- 出力管理へ移動した後、出力対象 job を自動選択するか、出力管理側で選ばせるか。

## 次の作業

- `ui-design.md` で UI 要件契約を作る。
- `scenario-design.md` で移動制約の受け入れ条件を作る。
- 人間レビュー後に `implementation-scope.md` を作る。
- 実装対象は frontend route、screen state、フェーズページ分解、navigation guard に分ける。

## 検証

- `plantuml -tsvg docs/exec-plans/active/2026-05-08-translation-flow-navigation-overhaul/navigation-state-machine.puml`: pass
