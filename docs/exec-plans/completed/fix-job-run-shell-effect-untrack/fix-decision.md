# Fix Decision

## 判断結果

- 判定: 完了
- 停止理由: なし

## 観測済み問題

- 問題: `JobRunPage.svelte` の `$effect` が `selectedJobTarget` の変化だけでなく `currentPhasePage`（`$state`）の変化にも反応して再実行される。ユーザーが「次へ進む」操作で `setCurrentPhasePage("persona")` を呼ぶと、`currentPhasePage` の変化が `$effect` を再実行し、`$effect` が `resolveInitialPhasePage` で "term" を返して `setCurrentPhasePage("term")` を呼ぶ。これにより persona 画面が即座に term 画面に戻り、遷移が成立しない。
- 期待との差分: `clickNext()` 後に `persona-generation-phase-screen` が visible になるべきだが、`$effect` の再実行で term 画面に戻るため visible にならず timeout する（E2E-UC-048）。`clickBodyCompleteNext()` 後に `translation-complete-translation-complete-screen` が visible になるべきだが、同じ理由で body 画面に戻るため visible にならず timeout する（E2E-UC-049）。

## 画面再現確認

- Wails 接続対象: `http://localhost:34115`（`npm run dev:wails:run`）
- 再現手順:
  1. 翻訳管理画面（`#translation-management`）でジョブ #3（`term_translation` Running 状態）のカードの選択領域（`translation-job-management-job-selection-region`）をクリックし、`#translation-management/job-run` へ遷移
  2. `job-run-next-action-footer` 内の「次へ進む」ボタンを JavaScript で `disabled` 解除してクリック（実アプリでは Running ジョブのため非活性だが、操作自体は `onPrimary` → `setCurrentPhasePage("persona")` を呼ぶ経路と同一）
- 操作結果:
  - `setCurrentPhasePage("persona")` 呼び出し（msgid=89）: `previous:"term"` → `result:"persona"` へ変化
  - `$effect` 再実行（msgid=91）: `currentPhasePageBeforeSet:"persona"` で起動（`currentPhasePage` 変化が依存として追跡されていたことが確定）
  - `resolveInitialPhasePage` が "term" を返す（msgid=92）: `currentPhase:"term_translation"` のため
  - `setCurrentPhasePage("term")` 呼び出し（msgid=93）: `previous:"persona"` → `result:"term"` に戻る
  - `$effect` が再び実行（msgid=95）: `currentPhasePageBeforeSet:"term"` で起動
  - `setCurrentPhasePage("term")` 呼び出し（msgid=97）: `previous:"term"` のため変化なし、安定
- 画面状態: term 画面が維持される（persona 画面への遷移が成立しない）
- 証跡: コンソールログ msgid=82〜103（`[TMP-OBS]` プレフィックスのログで観測、観測ログは削除済み）

## 確定原因

- 原因: `setCurrentPhasePage` 関数内の `const previous = currentPhasePage`（JobRunPage.svelte:113）が `currentPhasePage`（`$state`）を読み取っている。この関数が `$effect` 内（line 168）から呼ばれるため、Svelte 5 の reactivity がその読み取りを `$effect` の依存として登録する。その結果、ユーザー操作で `currentPhasePage` が変化すると `$effect` が再実行され、`resolveInitialPhasePage` で初期 phase に戻す処理が走る。
- 観測根拠:
  - `$effect` 実行ログ（msgid=91）で `currentPhasePageBeforeSet:"persona"` を観測。これは `selectedJobTarget` に変化がないにもかかわらず `$effect` が再実行されたことを示す。`selectedJobTarget` の変化だけが依存なら `$effect` は再実行されない。
  - `setCurrentPhasePage("term")` 呼び出しログ（msgid=93）で `stacktrace` が `JobRunPage.svelte:190`（`$effect` 内の `setCurrentPhasePage(initialPage)` 呼び出し行）を指している。
  - 循環が 1 回で止まる（msgid=95→97 で `previous:"term"` のため変化なし）ことも観測で確認した。

## 採用する修正方針

- 方針: `setCurrentPhasePage` 内の `const previous = currentPhasePage` の読み取りを Svelte 5 の `untrack()` でラップし、`$effect` の依存追跡から除外する。
  具体的には `import { untrack } from "svelte"` を追加し、`setCurrentPhasePage` の先頭を次のように変更する。
  ```ts
  function setCurrentPhasePage(phasePage: PhasePageId): void {
    const previous = untrack(() => currentPhasePage)
    currentPhasePage = phasePage
    // ...
  }
  ```
- 理由: `setCurrentPhasePage` の `previous` 比較は「同じ phase への重複呼び出しをスキップする」ための guard であり、この読み取りは `$effect` の依存に含める必要がない。`untrack()` を使うと `$effect` の実行コンテキストであっても依存追跡を抑制できるため、`currentPhasePage` の変化が `$effect` を再実行しなくなる。`selectedJobTarget`（prop）の変化だけが `$effect` を起動する、意図した依存関係に戻る。

## 禁止する修正

- 禁止修正 1: `$effect` 内で `setCurrentPhasePage` の呼び出しを削除し、`currentPhasePage` を直接代入する変更。
  - 理由: `setCurrentPhasePage` は page 変化に応じた controller の `setJobId` 呼び出しも担っており、`$effect` 内でも初回ジョブ選択時の初期 fetch をトリガーする責務がある。削除すると初期 fetch が走らなくなる。
- 禁止修正 2: `isUserNavigated`、`isEffectSyncing` のような新しい boolean `$state` を追加して `$effect` の再実行を条件分岐でスキップする対症療法。
  - 理由: 既存の `$state` モデルで説明できる原因（`$effect` の依存追跡への不正な登録）を、新しい状態値で隠す対症療法である。状態数が増えて条件の追跡が困難になる。
- 禁止修正 3: `$effect` を `$effect.pre` または `afterUpdate` 相当に変更して実行タイミングをずらす変更。
  - 理由: 依存追跡の問題を根本的に解決しない。タイミングをずらしても `currentPhasePage` が依存に残る限り循環は発生しうる。
- 禁止修正 4: `resolveInitialPhasePage` の戻り値をキャッシュして `$effect` の再実行時に比較する変更。
  - 理由: `currentPhasePage` の依存追跡が残ったまま `$effect` が再実行される問題を解決しない。

## 影響ファイル候補

- 候補: `frontend/src/ui/screens/job-run/JobRunPage.svelte`
  - 理由: `setCurrentPhasePage`（line 112-139）と `$effect`（line 156-177）が修正対象箇所。`untrack` の import 追加（line 1 付近の `import { onMount } from "svelte"` と同行または別行）と `setCurrentPhasePage` 内の 1 行変更のみで修正できる。
- 候補（単体テスト追加）: `frontend/src/ui/screens/job-run/` 配下の新規テストファイル
  - 理由: `setCurrentPhasePage` 呼び出し後に `$effect` が再実行されないことを Svelte testing で証明する単体テストを追加することが plan.md の影響範囲節に記載されている。既存テストファイルの修正ではなく新規追加なので実装 agent の判断範囲。
