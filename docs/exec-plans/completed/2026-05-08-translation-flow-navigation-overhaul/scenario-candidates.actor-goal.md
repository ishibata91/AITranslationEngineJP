# Scenario Candidates: 2026-05-08-translation-flow-navigation-overhaul / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TFN`

## Generator Scope

- `viewpoint`: 利用者の目的、開始経路、成功結果、観測点から候補を作る。
- `included_sources`: `./plan.md`, `./navigation-state-machine.puml`, `../../../spec.md`, `../../../architecture.md`, `../../../detail-specs/translation-job-management.md`, `../../../detail-specs/translation-job-setup.md`, `../../../detail-specs/term-translation-phase.md`, `../../../detail-specs/persona-generation-phase.md`, `../../../detail-specs/body-translation-phase.md`, `../../../detail-specs/translation-output-artifact.md`
- `excluded_sources`: 最終シナリオ表、候補採否、候補統合、競合解消、実装指示、ツール権限、プロダクトコード、プロダクトテスト、docs 正本変更。
- `generation_notes`: 旧 `Job Run` は大箱として残さない前提で読む。ただし、既存詳細仕様に残る `Job Run` 表記は競合候補として残す。

## Candidate Scenarios

### CAND-TFN-001 新規 job 作成直後に単語翻訳ページへ進む

- `source requirement`: [plan.md](./plan.md) line 14-19, 23-33, 67-76, 121-130 / [navigation-state-machine.puml](./navigation-state-machine.puml) line 27-31, 70-72 / [translation-job-setup.md](../../../detail-specs/translation-job-setup.md) line 19-22, 37-47 / [term-translation-phase.md](../../../detail-specs/term-translation-phase.md) line 21-24
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TFN-001`
- `actor`: 翻訳入力データから新しい翻訳 job を作りたい利用者。
- `trigger`: 翻訳セクション入口から新規翻訳を開始し、入力データを選び、`Job Setup` で job 作成を完了する。
- `expected outcome`: 作成された job と初期フェーズ状態が単語翻訳ページへ引き継がれ、利用者は単語翻訳フェーズの開始可否と summary を確認できる。
- `observable point`: 単語翻訳ページに対象 job、単語翻訳 summary、開始操作、次へ進めない理由が表示される。旧 `Job Run` のセッション取得操作は表示されない。
- `related detail requirement type`: `translation-job-setup` の対象、仕様、UI 契約。`term-translation-phase` の対象、UI 契約。
- `adoption hint`: 新規開始経路の主要正常系候補として扱える。
- `conflict hint`: 既存詳細仕様の `Job Run` 表記を、分解後の単語翻訳ページへ読み替える必要がある。

### CAND-TFN-002 未完了 job 一覧から対象 job を選んで再開する

- `source requirement`: [plan.md](./plan.md) line 34-44, 67-76, 132-138 / [navigation-state-machine.puml](./navigation-state-machine.puml) line 32-35, 66-76 / [translation-job-management.md](../../../detail-specs/translation-job-management.md) line 13-22, 26-45, 70-78 / [spec.md](../../../spec.md) line 53-59, 132-133
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TFN-002`
- `actor`: 作成済みの未完了翻訳 job を再開したい利用者。
- `trigger`: 翻訳セクション入口から途中再開を選び、未完了 job 一覧で対象 job を選ぶ。
- `expected outcome`: 選択した jobId と現在フェーズが固定され、利用者は対応するフェーズページから作業を再開できる。
- `observable point`: 未完了 job 一覧に `Completed` 以外の job、状態、現在フェーズ、進捗、操作可否、無効理由が表示される。選択後のフェーズページは選択済み job を表示対象にする。
- `related detail requirement type`: `translation-job-management` の対象、仕様、UI 契約。`spec.md` の AI 実行基盤、業務フロー。
- `adoption hint`: 途中再開経路の主要正常系候補として扱える。
- `conflict hint`: `translation-job-management` は選択先を `Job Run` と表現しているため、分解後のフェーズページへ置き換える必要がある。

### CAND-TFN-003 Ready job を一覧から実行入口として開く

- `source requirement`: [plan.md](./plan.md) line 67-76 / [navigation-state-machine.puml](./navigation-state-machine.puml) line 32-42, 74 / [translation-job-management.md](../../../detail-specs/translation-job-management.md) line 33-35, 40-47 / [term-translation-phase.md](../../../detail-specs/term-translation-phase.md) line 21-31, 71-78
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TFN-003`
- `actor`: 作成済みでまだ単語翻訳を開始していない job を進めたい利用者。
- `trigger`: 未完了 job 一覧から Ready job を選ぶ。
- `expected outcome`: Ready job は再編集画面ではなく、単語翻訳フェーズの read-only な実行入口として開く。
- `observable point`: 単語翻訳ページで対象 job が Ready であること、active phase run がないこと、開始操作が有効または無効理由付きで表示されることを確認できる。
- `related detail requirement type`: `translation-job-management` の仕様。`term-translation-phase` の対象、仕様、UI 契約。
- `adoption hint`: 新規作成直後以外の Ready job 開始候補として扱える。
- `conflict hint`: フェーズページから `Job Setup` へ戻って設定を再編集できるように見せると、対象差分の移動制限と競合する。

### CAND-TFN-004 単語翻訳完了後だけ NPC ペルソナ生成ページへ進む

- `source requirement`: [plan.md](./plan.md) line 34-54, 121-130 / [navigation-state-machine.puml](./navigation-state-machine.puml) line 37-43, 78-80 / [term-translation-phase.md](../../../detail-specs/term-translation-phase.md) line 55-57, 71-78 / [spec.md](../../../spec.md) line 129-130, 246-248
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TFN-004`
- `actor`: 単語翻訳フェーズを終えて次の翻訳段階へ進みたい利用者。
- `trigger`: 単語翻訳ページの `次へ進む` を使う。
- `expected outcome`: 単語翻訳フェーズが Completed で、ジョブ内辞書参照が成立している場合だけ、NPC ペルソナ生成ページへ移動できる。
- `observable point`: `sticky footer` に `次へ進む` と `一覧へ戻る` が表示される。未完了または辞書参照不能の時は、次へ進めない理由が近接表示される。
- `related detail requirement type`: `term-translation-phase` の仕様、UI 契約、受け入れ根拠。`spec.md` の業務フロー、用語。
- `adoption hint`: フェーズ間移動制約の主要正常系候補として扱える。
- `conflict hint`: 見出しクリック移動やグローバル直移動を許すと、開始経路が複数化して候補の前提と競合する。

### CAND-TFN-005 NPC ペルソナ生成完了後だけ本文翻訳ページへ進む

- `source requirement`: [plan.md](./plan.md) line 34-54, 121-130 / [navigation-state-machine.puml](./navigation-state-machine.puml) line 44-50, 82-84 / [persona-generation-phase.md](../../../detail-specs/persona-generation-phase.md) line 21-24, 44-64 / [body-translation-phase.md](../../../detail-specs/body-translation-phase.md) line 20-23 / [spec.md](../../../spec.md) line 129-131, 246-248
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TFN-005`
- `actor`: NPC ペルソナ生成フェーズを終えて本文翻訳へ進みたい利用者。
- `trigger`: NPC ペルソナ生成ページの `次へ進む` を使う。
- `expected outcome`: persona phase が Completed で、snapshot 参照が成立している場合だけ、本文翻訳ページへ移動できる。
- `observable point`: NPC ペルソナ生成ページに phase state、progress、snapshot 参照状態、body phase readiness、無効理由が表示される。
- `related detail requirement type`: `persona-generation-phase` の対象、仕様、UI 契約。`body-translation-phase` の対象。
- `adoption hint`: フェーズ間移動制約の主要正常系候補として扱える。
- `conflict hint`: 実行操作を `sticky footer` に置くと、移動導線とフェーズ操作の責務分離と競合する。

### CAND-TFN-006 本文翻訳完了後に翻訳完了ページで結果を確認する

- `source requirement`: [plan.md](./plan.md) line 78-88, 132-138 / [navigation-state-machine.puml](./navigation-state-machine.puml) line 51-61, 86-92 / [body-translation-phase.md](../../../detail-specs/body-translation-phase.md) line 13-23, 42-43, 64-80 / [spec.md](../../../spec.md) line 43, 65-67
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TFN-006`
- `actor`: 本文翻訳を終えた結果を確認したい利用者。
- `trigger`: 本文翻訳ページで job が Completed になる。
- `expected outcome`: 翻訳完了ページへ移動し、原文と訳文をページング表示で確認できる。
- `observable point`: 翻訳完了ページに原文、訳文、ページング、一覧へ戻る導線、出力管理への移動ボタンが表示される。XML 出力や再出力の操作は表示されない。
- `related detail requirement type`: `body-translation-phase` の対象、仕様、UI 契約。`translation-output-artifact` の対象外境界。
- `adoption hint`: 翻訳管理内の結果確認候補として扱える。
- `conflict hint`: 状態遷移図では `Canceled` と `Failed` も翻訳完了ページへ向かうため、完了ページの対象状態を designer が分離判断する必要がある。

### CAND-TFN-007 翻訳完了ページから出力管理へ移動する

- `source requirement`: [plan.md](./plan.md) line 78-98, 132-148 / [navigation-state-machine.puml](./navigation-state-machine.puml) line 58-61, 90-99, 102-118, 121-123 / [translation-output-artifact.md](../../../detail-specs/translation-output-artifact.md) line 13-22, 26-28, 61-75 / [body-translation-phase.md](../../../detail-specs/body-translation-phase.md) line 15-16, 42-43
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TFN-007`
- `actor`: 完了済み翻訳 job から xTranslator 互換 XML を出力したい利用者。
- `trigger`: 翻訳完了ページの出力管理への移動ボタンを使う。
- `expected outcome`: 翻訳管理では出力処理を開始せず、成果物出力セクションで Completed job 一覧または Output Review を確認できる。
- `observable point`: 成果物出力セクションに completed job list、selected job summary、output readiness、拒否理由、preview、出力 action が表示される。
- `related detail requirement type`: `translation-output-artifact` の対象、仕様、UI 契約。`body-translation-phase` の出力 readiness。
- `adoption hint`: 翻訳管理と成果物出力の境界候補として扱える。
- `conflict hint`: 出力対象 job を自動選択するか、成果物出力側で選ばせるかは plan の未決事項である。

### CAND-TFN-008 グローバル直移動でフェーズページへ入れない

- `source requirement`: [plan.md](./plan.md) line 56-66, 113-119, 121-130 / [navigation-state-machine.puml](./navigation-state-machine.puml) line 66-68, 124-132, 145-148 / [translation-job-management.md](../../../detail-specs/translation-job-management.md) line 43-45, 70-78
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TFN-008`
- `actor`: 対象 job が未確定のままフェーズページへ入ろうとした利用者。
- `trigger`: グローバルナビ、復元状態、または不整合な route state から単語翻訳、NPC ペルソナ生成、本文翻訳のいずれかへ入ろうとする。
- `expected outcome`: 対象 job が曖昧なフェーズページは表示されず、未完了 job 一覧へ戻される。
- `observable point`: 未完了 job 一覧が表示され、利用者は job の状態、現在フェーズ、操作可否、再開不可理由を確認して選び直せる。
- `related detail requirement type`: `translation-job-management` の一覧表示、安全側表示、UI 契約。
- `adoption hint`: 直リンク防止の actor-goal 候補として扱える。
- `conflict hint`: route state の不整合と本来のグローバルナビ導線の扱いは、state-transition 観点と重なる可能性がある。

### CAND-TFN-009 フェーズページから入力データや Job Setup へ戻らない

- `source requirement`: [plan.md](./plan.md) line 34-44, 121-130 / [navigation-state-machine.puml](./navigation-state-machine.puml) line 124-129 / [translation-job-management.md](../../../detail-specs/translation-job-management.md) line 29-35 / [translation-job-setup.md](../../../detail-specs/translation-job-setup.md) line 19-22, 26-36
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TFN-009`
- `actor`: 作成済み job の入力または設定を変えたい利用者。
- `trigger`: 単語翻訳、NPC ペルソナ生成、本文翻訳のいずれかのフェーズページを開いている。
- `expected outcome`: 作成済み job の前提を編集する導線は出ず、利用者は既存 job を巻き戻さずに一覧へ戻るか、新規 job 作成を別経路で始める。
- `observable point`: フェーズページの移動導線は `次へ進む` と `一覧へ戻る` に絞られ、入力データページや `Job Setup` へ戻る操作は表示されない。
- `related detail requirement type`: `translation-job-management` の job 識別と read-only 実行入口。`translation-job-setup` の job 作成条件。
- `adoption hint`: 作成済み job の整合性維持候補として扱える。
- `conflict hint`: 同一 input から複数 job を作れる仕様と併せて、新規作成導線の所在を designer が確認する必要がある。

### CAND-TFN-010 旧 Job Run のセッション取得なしで選択済み job を表示する

- `source requirement`: [plan.md](./plan.md) line 67-76, 121-130 / [navigation-state-machine.puml](./navigation-state-machine.puml) line 30-42, 70-76, 134-138 / [translation-job-management.md](../../../detail-specs/translation-job-management.md) line 33-35, 43-45 / [architecture.md](../../../architecture.md) line 119-130, 207-214
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TFN-010`
- `actor`: 選択した job のフェーズ状態をすぐ確認したい利用者。
- `trigger`: `Job Setup` 完了結果、または未完了 job 一覧の選択結果からフェーズページへ入る。
- `expected outcome`: フェーズページは jobId を受け取り、summary を表示する。利用者は別操作でセッション取得を実行しない。
- `observable point`: フェーズページの表示対象は作成結果または一覧選択結果と一致する。セッション取得ボタン、セッション取得待ち空状態、別 job を探す操作は表示されない。
- `related detail requirement type`: `translation-job-management` の job 選択、参照不能時の安全側表示。`architecture.md` の screen state と query / command 境界。
- `adoption hint`: 旧 `Job Run` 廃止差分の actor-goal 候補として扱える。
- `conflict hint`: 既存詳細仕様の `Job Run` 表記や stepper 連携と、セッション取得廃止の関係を designer が整理する必要がある。

## Open Notes

- `human decision candidate`: 翻訳完了ページで `Canceled` と `Failed` を扱うか、別の終了表示へ分けるか。
- `human decision candidate`: 出力管理へ移動した後、対象 job を自動選択するか、成果物出力側で選ばせるか。
- `human decision candidate`: 旧 `Job Run` 表記を全詳細仕様でフェーズページへ読み替えるか、互換名として一時的に残すか。
- `merge candidate`: `CAND-TFN-001` と `CAND-TFN-003` は、単語翻訳ページ開始の候補として統合対象になりうる。
- `merge candidate`: `CAND-TFN-004` と `CAND-TFN-005` は、フェーズ間 `次へ進む` の共通候補として統合対象になりうる。
- `merge candidate`: `CAND-TFN-007` と `CAND-TFN-010` は、翻訳管理と外部表示対象の受け渡し候補として接続されうる。
- `rejection candidate`: グローバルナビから各フェーズページへ直接移動できる候補は、対象差分と状態遷移図の禁止移動に反する。
- `rejection candidate`: フェーズページから入力データページや `Job Setup` へ戻って作成済み job を編集できる候補は、対象差分と状態遷移図の禁止移動に反する。
