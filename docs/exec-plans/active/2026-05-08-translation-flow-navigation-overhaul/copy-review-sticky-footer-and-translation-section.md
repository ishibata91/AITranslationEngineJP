# sticky footer / 翻訳セクション文言レビュー

- `status`: human-copy-review
- `purpose`: sticky footer と翻訳セクションの文言を 1 件ずつ見直す。
- `source`: 現行 frontend 実装から抽出。
- `scope`: 文言整理だけ。挙動変更、導線変更、仕様正本化はこの一覧では扱わない。

## 文言方針

- 画面に出す文言は、利用者の作業語彙を優先する。
- 実装都合の画面名、状態名、変数名は、画面表示の修正案に入れない。
- `job` は画面表示では `ジョブ` と書く。
- `phase` は原則として `翻訳段階` と書く。
- footer は、移動先と移動できない理由だけを伝える。

## sticky footer 共通

- 追加理由あり表示
  - 現行: `ほか {remainingCount} 件`
  - 修正案: `ほか {remainingCount} 件の確認が必要です`
  - 出典: `frontend/src/ui/components/StickyActionFooter.svelte`

- 追加理由 tooltip 見出し
  - 現行: `残りの不足`
  - 修正案: `確認が必要な項目`
  - 出典: `frontend/src/ui/components/StickyActionFooter.svelte`

- 理由なし default
  - 現行: `確認する項目はありません。`
  - 修正案: （何も表示しない）
  - 出典: `frontend/src/ui/components/StickyActionFooter.svelte`

## データロード footer

- title
  - 現行: `次の移動`
  - 修正案: `次の作業`
  - 出典: `frontend/src/ui/screens/translation-input/InputReviewPage.svelte`

- description
  - 現行: `登録済み入力データから Job Setup へ進む`
  - 修正案: `選択した入力データで、ジョブの作成確認へ進みます。`
  - 出典: `frontend/src/ui/screens/translation-input/InputReviewPage.svelte`

- emptyText
  - 現行: `入力データを選択済みです。Job Setup で job 作成条件を確認してください。`
  - 修正案: `入力データを選択済みです。次に翻訳設定を確認します。`
  - 出典: `frontend/src/ui/screens/translation-input/InputReviewPage.svelte`

- primaryLabel
  - 現行: `Job Setup へ進む`
  - 修正案: `翻訳設定へ進む`
  - 出典: `frontend/src/ui/screens/translation-input/InputReviewPage.svelte`

- 表示条件
  - 現行: 登録済み入力データを選択済み、かつ Job Setup へ進める状態
  - 修正案: 入力データを選択済みで、翻訳設定へ進める状態
  - 出典: `frontend/src/ui/screens/translation-input/InputReviewPage.svelte`

## Job Setup footer

- title
  - 現行: `作成前確認`
  - 修正案: `ジョブの作成確認`
  - 出典: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`

- description
  - 現行: `viewModel.createSectionText`
  - 修正案: `入力データと翻訳設定を確認し、最初の翻訳段階へ進む準備をします。`
  - 注意: presenter 側の生成文言も別途確認対象。
  - 出典: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`

- emptyText: 作成可能
  - 現行: `不足はありません。`
  - 修正案: `作成に必要な確認は完了しています。`
  - 出典: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`

- emptyText: 作成不可
  - 現行: `作成前確認はまだ未完了です。`
  - 修正案: `作成前に確認が必要な項目があります。`
  - 出典: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`

- primaryLabel
  - 現行: `次へ`
  - 修正案: `単語翻訳へ進む`
  - 出典: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`

- 補助文
  - 現行: `cache missing は Job Setup で再構築しません。Input Review の再構築導線へ戻ってください。`
  - 修正案: `入力データの再構築が必要です。入力データの確認画面に戻ってください。`
  - 出典: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`

- secondaryLabel
  - 現行: `Input Review へ戻る`
  - 修正案: `入力データの確認へ戻る`
  - 出典: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`

## 翻訳段階 footer 共通

- primaryLabel default
  - 現行: `次へ進む`
  - 修正案: `次へ進む`
  - 出典: `frontend/src/ui/screens/job-run/PhaseNavigationFooter.svelte`

- emptyText
  - 現行: `移動条件を満たしています。`
  - 修正案: `次の作業へ進めます。`
  - 出典: `frontend/src/ui/screens/job-run/PhaseNavigationFooter.svelte`

- secondaryLabel
  - 現行: `一覧へ戻る`
  - 修正案: `未完了一覧へ戻る`
  - 出典: `frontend/src/ui/screens/job-run/PhaseNavigationFooter.svelte`

- outputLabel
  - 現行: `出力管理へ移動`
  - 修正案: `出力管理へ進む`
  - 出典: `frontend/src/ui/screens/job-run/PhaseNavigationFooter.svelte`

- showPrimary false note
  - 現行: `次へ進む操作はありません。`
  - 修正案: `この画面で進む先はありません。`
  - 出典: `frontend/src/ui/screens/job-run/PhaseNavigationFooter.svelte`

## 単語翻訳 footer

- title
  - 現行: `単語翻訳の移動`
  - 修正案: `単語翻訳の次の作業`
  - 出典: `frontend/src/ui/screens/job-run/JobRunPage.svelte`

- description
  - 現行: `単語翻訳の完了と辞書参照が成立した時だけ次へ進みます。`
  - 修正案: `単語翻訳が完了し、辞書を参照できる場合だけ次へ進めます。`
  - 出典: `frontend/src/ui/screens/job-run/JobRunPage.svelte`

- blocked fallback
  - 現行: `次へ進めません。単語翻訳の完了と辞書参照を確認してください。`
  - 修正案: `次へ進めません。単語翻訳の完了状況と辞書の参照状態を確認してください。`
  - 出典: `frontend/src/ui/screens/job-run/JobRunPage.svelte`

## NPC ペルソナ生成 footer

- title
  - 現行: `NPC ペルソナ生成の移動`
  - 修正案: `NPC ペルソナ生成の次の作業`
  - 出典: `frontend/src/ui/screens/job-run/JobRunPage.svelte`

- description
  - 現行: `ペルソナ生成の完了と snapshot 参照が成立した時だけ次へ進みます。`
  - 修正案: `ペルソナ生成が完了し、生成結果を参照できる場合だけ次へ進めます。`
  - 出典: `frontend/src/ui/screens/job-run/JobRunPage.svelte`

- blocked fallback
  - 現行: `次へ進めません。ペルソナ参照を確認してください。`
  - 修正案: `次へ進めません。ペルソナ生成の完了状況と参照状態を確認してください。`
  - 出典: `frontend/src/ui/screens/job-run/JobRunPage.svelte`

## 本文翻訳 footer

- title
  - 現行: `本文翻訳の移動`
  - 修正案: `本文翻訳の次の作業`
  - 出典: `frontend/src/ui/screens/job-run/JobRunPage.svelte`

- description
  - 現行: `本文翻訳が Completed で、翻訳結果と出力状態が整合する時だけ完了確認へ進みます。`
  - 修正案: `本文翻訳が完了し、翻訳結果を確認できる場合だけ完了確認へ進めます。`
  - 出典: `frontend/src/ui/screens/job-run/JobRunPage.svelte`

- blocked fallback
  - 現行: `出力できません。本文翻訳の完了と翻訳結果の整合を確認してください。`
  - 修正案: `完了確認へ進めません。本文翻訳の完了状況と翻訳結果を確認してください。`
  - 出典: `frontend/src/ui/screens/job-run/JobRunPage.svelte`

## 翻訳完了 footer

- title
  - 現行: `翻訳完了後の移動`
  - 修正案: `翻訳完了後の次の作業`
  - 出典: `frontend/src/ui/screens/job-run/JobRunPage.svelte`

- description
  - 現行: `翻訳結果の確認後は出力管理へ移動できます。出力対象 job は出力管理で選びます。`
  - 修正案: `翻訳結果を確認した後は、出力管理で出力対象を選びます。`
  - 出典: `frontend/src/ui/screens/job-run/JobRunPage.svelte`

- primaryLabel
  - 現行: `出力管理へ移動`
  - 修正案: `出力管理へ進む`
  - 出典: `frontend/src/ui/screens/job-run/JobRunPage.svelte`

## 翻訳セクション共通

- セクション label
  - 現行: `翻訳セクション`
  - 修正案: `翻訳管理`
  - 出典: `frontend/src/ui/screens/translation-job-management/TranslationManagementStepper.svelte`

- セクション title
  - 現行: `進行ステップ`
  - 修正案: `ジョブの進み方`
  - 出典: `frontend/src/ui/screens/translation-job-management/TranslationManagementStepper.svelte`

- 現在位置 label
  - 現行: `現在位置`
  - 修正案: `現在の作業`
  - 出典: `frontend/src/ui/screens/translation-job-management/TranslationManagementStepper.svelte`

- 直接移動可 label
  - 現行: `直接移動可`
  - 修正案: `ここから開けます`
  - 出典: `frontend/src/ui/screens/translation-job-management/TranslationManagementStepper.svelte`

- 参照のみ label
  - 現行: `参照のみ`
  - 修正案: `順番に進む作業`
  - 出典: `frontend/src/ui/screens/translation-job-management/TranslationManagementStepper.svelte`

- aria label
  - 現行: `翻訳管理セクション`
  - 修正案: `翻訳管理の進行状況`
  - 出典: `frontend/src/ui/screens/translation-job-management/TranslationManagementStepper.svelte`

## 翻訳セクション項目

- step 1: `job-management`
  - label 現行: `未完了 job 一覧`
  - label 修正案: `未完了のジョブ`
  - description 現行: `新規翻訳の開始と未完了 job の途中再開を選びます。`
  - description 修正案: `新しい翻訳を始めるか、途中のジョブを再開します。`
  - 直接移動: はい
  - 出典: `frontend/src/ui/stores/shell-state.ts`

- step 2: `input-review`
  - label 現行: `データロード`
  - label 修正案: `入力データの確認`
  - description 現行: `新規翻訳開始後に入力ファイルの登録結果を確認します。`
  - description 修正案: `翻訳に使う入力データを選び、登録結果を確認します。`
  - 直接移動: いいえ
  - 出典: `frontend/src/ui/stores/shell-state.ts`

- step 3: `job-setup`
  - label 現行: `Job Setup`
  - label 修正案: `翻訳設定`
  - description 現行: `データロード後に AI 設定と ready job 作成を確認します。`
  - description 修正案: `入力データと AI 設定を確認し、ジョブを作成します。`
  - 直接移動: いいえ
  - 出典: `frontend/src/ui/stores/shell-state.ts`

- step 4: `term-translation`
  - label 現行: `単語翻訳`
  - label 修正案: `単語翻訳`
  - description 現行: `選択済み job で単語翻訳フェーズを実行します。`
  - description 修正案: `選択したジョブで、単語翻訳を実行します。`
  - 直接移動: いいえ
  - 出典: `frontend/src/ui/stores/shell-state.ts`

- step 5: `persona-generation`
  - label 現行: `NPC ペルソナ生成`
  - label 修正案: `NPC ペルソナ生成`
  - description 現行: `単語翻訳完了後にペルソナ生成フェーズを実行します。`
  - description 修正案: `単語翻訳の完了後に、NPC の話し方や役割を整理します。`
  - 直接移動: いいえ
  - 出典: `frontend/src/ui/stores/shell-state.ts`

- step 6: `body-translation`
  - label 現行: `本文翻訳`
  - label 修正案: `本文翻訳`
  - description 現行: `ペルソナ参照成立後に本文翻訳フェーズを実行します。`
  - description 修正案: `NPC ペルソナを参照できる状態で、本文の翻訳を実行します。`
  - 直接移動: いいえ
  - 出典: `frontend/src/ui/stores/shell-state.ts`

- step 7: `translation-complete`
  - label 現行: `翻訳完了`
  - label 修正案: `翻訳結果の確認`
  - description 現行: `本文翻訳 Completed 後に原文と訳文を確認します。`
  - description 修正案: `本文翻訳が完了した後に、原文と訳文を確認します。`
  - 直接移動: いいえ
  - 出典: `frontend/src/ui/stores/shell-state.ts`

- step 8: `output-management`
  - label 現行: `出力管理`
  - label 修正案: `出力管理`
  - description 現行: `翻訳完了後に出力管理へ移動し、Completed job を選びます。`
  - description 修正案: `翻訳結果を確認した後に、出力するジョブを選びます。`
  - 直接移動: いいえ
  - 出典: `frontend/src/ui/stores/shell-state.ts`
