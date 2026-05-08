# Scenario Candidates: 2026-05-08-translation-flow-navigation-overhaul / failure

- `generator`: `failure`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TFN`

## Generator Scope

- `viewpoint`: 失敗観点。対象 job 未確定、参照不能、禁止移動、次工程へ進めない条件、成果物出力との混線を候補化する。
- `included_sources`: `plan.md`, `navigation-state-machine.puml`, `docs/spec.md`, `docs/architecture.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`, `docs/detail-specs/translation-output-artifact.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本変更、`.codex` 変更、最終シナリオ表、候補の採否、候補の統合判断
- `generation_notes`: 旧 `Job Run` のセッション取得を回復経路として使わない前提で、表示拒否、操作拒否、復帰先、観測点を候補化する。

## Candidate Scenarios

### CAND-TFN-001 対象 job 未確定のフェーズ表示を拒否する

- `source requirement`: `plan.md:56-65`, `plan.md:67-76`, `navigation-state-machine.puml:130-132`, `navigation-state-machine.puml:145-148`, `docs/detail-specs/translation-job-management.md:43-45`
- `viewpoint`: 参照不能、設定不整合、回復動作
- `candidate scenario id`: `CAND-TFN-001`
- `actor`: 翻訳管理画面を開く利用者
- `trigger`: route state または復元状態に対象 job がなく、単語翻訳、NPC ペルソナ生成、本文翻訳のいずれかのページを表示しようとする。
- `rejected operation`: 対象 job が未確定のままフェーズ summary を取得し、フェーズ操作を表示する。
- `expected error`: 対象 job が未確定であるため、フェーズページを表示できない理由を出し、未完了 job 一覧へ戻す。
- `expected outcome`: フェーズ操作、summary、旧 `Job Run` のセッション取得導線を出さない。未完了 job 一覧で再開対象を選ばせる。
- `observable point`: 未完了 job 一覧が表示されること。フェーズ固有の start、pause、resume、retry、cancel が表示されないこと。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `consistency_requirement`, `recovery_requirement`
- `adoption hint`: 直リンク防止と旧セッション取得廃止を同時に確認する異常系候補として扱える。
- `conflict hint`: 旧詳細仕様の「選択した job を Job Run の表示対象にする」という表現は、新しいフェーズページ表示対象へ読み替える必要がある。

### CAND-TFN-002 参照不能 job をフェーズページの対象にしない

- `source requirement`: `plan.md:56-65`, `plan.md:67-76`, `navigation-state-machine.puml:32-35`, `docs/detail-specs/translation-job-management.md:40-45`, `docs/detail-specs/translation-job-management.md:77-78`
- `viewpoint`: 参照不能、設定不整合、回復動作
- `candidate scenario id`: `CAND-TFN-002`
- `actor`: 未完了 job 一覧から job を再開しようとする利用者
- `trigger`: 未完了 job 一覧の選択結果が、削除済み、集約不能、入力キャッシュ欠落、または phase progress 集約不能になっている。
- `rejected operation`: 参照不能または集約不能の job を単語翻訳、NPC ペルソナ生成、本文翻訳ページの表示対象にする。
- `expected error`: 参照不能または集約不能の理由カテゴリを表示し、危険操作を無効にする。
- `expected outcome`: 対象 job はフェーズページへ進まない。状態表示だけでは job 状態を変更しない。
- `observable point`: 一覧上で再開不可理由を確認できること。フェーズページへ遷移しないこと。空一覧または成功状態として扱わないこと。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: 未完了 job 一覧を再開入口に固定する時の参照不能候補として扱える。
- `conflict hint`: 集約不能時に一覧に残して理由を見せるか、一覧から除外するかは、最終統合時に UI 表示粒度が競合しうる。

### CAND-TFN-003 グローバルナビからフェーズページへ直接移動しない

- `source requirement`: `plan.md:14-19`, `plan.md:56-65`, `plan.md:121-130`, `navigation-state-machine.puml:13-20`, `navigation-state-machine.puml:66-68`, `navigation-state-machine.puml:130-132`
- `viewpoint`: 失敗入力、設定不整合、回復動作
- `candidate scenario id`: `CAND-TFN-003`
- `actor`: グローバルナビまたは dashboard から翻訳セクションへ入る利用者
- `trigger`: グローバルナビ、保存済み UI 状態、または dashboard 導線がフェーズページを直接開こうとする。
- `rejected operation`: 翻訳セクション入口を経由せず、対象 job 未確定のままフェーズページへ入る。
- `expected error`: フェーズページへの直接移動は禁止され、翻訳セクション入口または未完了 job 一覧へ戻る。
- `expected outcome`: 新規開始または途中再開の入口を先に選ばせる。フェーズページは job 選択済みの場合だけ表示する。
- `observable point`: グローバルナビ上にフェーズ直リンクが存在しないこと。復元状態が不整合な場合は未完了 job 一覧へ戻ること。
- `related detail requirement type`: `authorization_requirement`, `state_requirement`, `consistency_requirement`, `recovery_requirement`
- `adoption hint`: route state と UI ナビゲーションの両方に効く異常系候補として扱える。
- `conflict hint`: Wails 前提では URL 直リンクが起きにくいという plan の記述と、UI 状態復元の異常系をどこまで扱うかが競合しうる。

### CAND-TFN-004 Job Setup の job 作成失敗時に単語翻訳ページへ進まない

- `source requirement`: `plan.md:67-76`, `navigation-state-machine.puml:27-30`, `navigation-state-machine.puml:70-72`, `docs/detail-specs/translation-job-setup.md:19-22`, `docs/detail-specs/translation-job-setup.md:41-49`, `docs/detail-specs/translation-job-setup.md:67-70`
- `viewpoint`: 失敗入力、設定不整合、保存失敗
- `candidate scenario id`: `CAND-TFN-004`
- `actor`: Job Setup で翻訳 job を作成する利用者
- `trigger`: API key 不足、model 未選択、stale なモデル一覧、保存失敗、または作成結果の jobId 欠落が発生する。
- `rejected operation`: Ready job が作成されていない状態で、単語翻訳ページへ移動する。
- `expected error`: job 作成に失敗した理由を Job Setup 上で表示し、単語翻訳ページへの遷移を止める。
- `expected outcome`: Ready job、初期 phase 状態、対象 input が固定された時だけ単語翻訳ページへ進む。
- `observable point`: 作成失敗時に Job Setup に留まること。対象 job 未確定のフェーズ summary を読まないこと。
- `related detail requirement type`: `failure_handling_requirement`, `boundary_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: 新規 job 作成直後に単語翻訳ページへ進む正常系の失敗境界として扱える。
- `conflict hint`: 作成済み job を確認する導線と作成失敗時の留まり方は、UI 上の近接表示位置が競合しうる。

### CAND-TFN-005 単語翻訳未完了または辞書参照不能なら次へ進まない

- `source requirement`: `plan.md:34-44`, `plan.md:45-55`, `navigation-state-machine.puml:37-43`, `navigation-state-machine.puml:78-80`, `docs/detail-specs/term-translation-phase.md:50-57`, `docs/detail-specs/term-translation-phase.md:71-79`
- `viewpoint`: 参照不能、設定不整合、回復動作
- `candidate scenario id`: `CAND-TFN-005`
- `actor`: 単語翻訳ページで次へ進もうとする利用者
- `trigger`: 単語翻訳フェーズが未完了、失敗中、辞書参照不能、または active phase run が存在する。
- `rejected operation`: NPC ペルソナ生成ページへ進む。後続 phase run を作成する。
- `expected error`: 次へ進めない理由を `sticky footer` にリアクティブに表示する。
- `expected outcome`: 単語翻訳ページに留まり、実行、再開、再試行などの回復操作はページ本文で扱う。
- `observable point`: `次へ進む` が無効で、理由が近接表示されること。NPC ペルソナ生成の summary や run が作成されないこと。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `consistency_requirement`, `recovery_requirement`
- `adoption hint`: phase 間移動を `sticky footer` に集約した時の代表的なブロック候補として扱える。
- `conflict hint`: 旧 `Job Run` の状態表示語彙を各フェーズページ本文と footer のどちらへ出すかは、UI 統合時に競合しうる。

### CAND-TFN-006 ペルソナ snapshot 参照不能なら本文翻訳へ進まない

- `source requirement`: `plan.md:34-44`, `plan.md:45-55`, `navigation-state-machine.puml:44-50`, `navigation-state-machine.puml:82-84`, `docs/detail-specs/persona-generation-phase.md:37-45`, `docs/detail-specs/persona-generation-phase.md:56-64`
- `viewpoint`: 参照不能、設定不整合、回復動作
- `candidate scenario id`: `CAND-TFN-006`
- `actor`: NPC ペルソナ生成ページで本文翻訳へ進もうとする利用者
- `trigger`: NPC ペルソナ生成フェーズが未完了、失敗中、snapshot 参照不能、または本文翻訳 readiness が成立していない。
- `rejected operation`: 本文翻訳ページへ進む。本文翻訳 phase run を作成する。
- `expected error`: 本文翻訳へ進めない理由を `sticky footer` に表示し、ページ本文の回復操作へ誘導する。
- `expected outcome`: NPC ペルソナ生成ページに留まる。成功済み persona を重複作成しない。
- `observable point`: `次へ進む` が無効で、snapshot 参照状態と body readiness の不足理由を確認できること。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `data_requirement`, `recovery_requirement`
- `adoption hint`: 後続本文翻訳の入力前提を壊さない異常系候補として扱える。
- `conflict hint`: snapshot ID または digest の表示粒度は、失敗理由表示と redaction の境界で競合しうる。

### CAND-TFN-007 本文翻訳が出力可能状態でなければ出力管理へ進めない

- `source requirement`: `plan.md:78-88`, `plan.md:89-99`, `navigation-state-machine.puml:51-61`, `navigation-state-machine.puml:86-92`, `docs/detail-specs/body-translation-phase.md:34-45`, `docs/detail-specs/body-translation-phase.md:70-75`, `docs/detail-specs/translation-output-artifact.md:26-29`
- `viewpoint`: 設定不整合、参照不能、回復動作、競合候補
- `candidate scenario id`: `CAND-TFN-007`
- `actor`: 本文翻訳ページまたは翻訳完了ページから出力管理へ移動しようとする利用者
- `trigger`: 本文翻訳フェーズが Failed、Canceled、RecoverableFailed、field result 不整合、status 不整合、または output readiness false である。
- `rejected operation`: 成功済み Completed job として成果物出力を開始する。
- `expected error`: 出力管理へ渡せない理由を表示し、成果物出力処理を開始しない。
- `expected outcome`: 出力管理は Completed job だけを出力候補にする。翻訳管理側は出力処理を直接開始しない。
- `observable point`: 出力管理ボタンまたは遷移先で拒否理由を確認できること。成功状態の artifact や output row が作られないこと。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `consistency_requirement`, `data_requirement`
- `adoption hint`: 翻訳完了ページと出力管理の責務分離を確認する異常系候補として扱える。
- `conflict hint`: 状態図は Body から翻訳完了ページへ `job Completed / Canceled / Failed` で遷移すると書く。一方で出力仕様は Completed job だけを候補にするため、Failed / Canceled を完了ページでどう見せるかは競合しうる。

### CAND-TFN-008 Completed job を翻訳管理の未完了一覧に混ぜない

- `source requirement`: `plan.md:89-99`, `plan.md:132-138`, `navigation-state-machine.puml:95-118`, `docs/detail-specs/translation-job-management.md:26-28`, `docs/detail-specs/translation-job-management.md:68-82`, `docs/detail-specs/translation-output-artifact.md:19-28`
- `viewpoint`: 設定不整合、データ整合性、回復動作
- `candidate scenario id`: `CAND-TFN-008`
- `actor`: 翻訳管理または成果物出力を開く利用者
- `trigger`: Completed job が翻訳管理の未完了一覧に混入する、または未完了 job が成果物出力の completed job 一覧に混入する。
- `rejected operation`: 翻訳管理から Completed job を再開対象として扱う。成果物出力から未完了 job の artifact 生成を開始する。
- `expected error`: 一覧の責務と job 状態が一致しないため、該当操作を出さない、または無効理由を表示する。
- `expected outcome`: 翻訳管理は未完了 job だけを扱い、成果物出力は completed job だけを扱う。
- `observable point`: Completed job が未完了一覧に表示されないこと。未完了、失敗中、Canceled、不整合 job では出力 action が無効になること。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`, `failure_handling_requirement`
- `adoption hint`: 翻訳管理と成果物出力のセクション分離を確認する異常系候補として扱える。
- `conflict hint`: 翻訳完了ページから一覧へ戻る時に、Completed job をどちらの一覧で見せるかが統合時に競合しうる。

### CAND-TFN-009 旧セッション取得に依存した再開を拒否する

- `source requirement`: `plan.md:67-76`, `plan.md:121-130`, `navigation-state-machine.puml:32-40`, `navigation-state-machine.puml:134-138`, `docs/detail-specs/translation-job-management.md:33-35`, `docs/architecture.md:125-130`, `docs/architecture.md:207-214`
- `viewpoint`: 設定不整合、参照不能、回復動作、競合候補
- `candidate scenario id`: `CAND-TFN-009`
- `actor`: 未完了 job を再開したい利用者
- `trigger`: 画面が旧 `Job Run` のセッション取得ボタンまたはセッション取得処理に頼って job を探そうとする。
- `rejected operation`: フェーズページ側から別操作で対象 job を取得し、一覧選択や Job Setup 作成結果とは別の対象取得経路を作る。
- `expected error`: セッション取得操作を出さず、未完了 job 一覧の選択または Job Setup 作成結果から job を受け取る。
- `expected outcome`: フェーズページは「選ばれた job を進める画面」として表示される。対象 job がない場合は未完了 job 一覧へ戻る。
- `observable point`: セッション取得ボタンが表示されないこと。対象 job なしでは summary 更新が始まらないこと。
- `related detail requirement type`: `compatibility_requirement`, `state_requirement`, `consistency_requirement`, `failure_handling_requirement`
- `adoption hint`: 旧 `Job Run` 分解による回帰防止候補として扱える。
- `conflict hint`: 既存詳細仕様の `Job Run` 表現と、新しいフェーズページ分解の語彙が競合する。最終シナリオでは呼称と表示対象を designer が整理する必要がある。

### CAND-TFN-010 出力管理へ移動後の出力対象 job 未選択を成功扱いしない

- `source requirement`: `plan.md:78-88`, `plan.md:146-149`, `navigation-state-machine.puml:58-61`, `navigation-state-machine.puml:102-116`, `docs/detail-specs/translation-output-artifact.md:26-29`, `docs/detail-specs/translation-output-artifact.md:61-69`
- `viewpoint`: 参照不能、設定不整合、回復動作、人間判断候補
- `candidate scenario id`: `CAND-TFN-010`
- `actor`: 翻訳完了ページから出力管理へ移動する利用者
- `trigger`: 出力管理へ移動したが、出力対象 job が自動選択されていない、または選択対象が Completed job として確認できない。
- `rejected operation`: 出力対象 job 未選択のまま preview、XML 出力、再出力を開始する。
- `expected error`: completed job 一覧で出力対象 job を選ぶ必要があること、または自動選択できない理由を表示する。
- `expected outcome`: 出力対象 job が固定され、output readiness を確認できるまで出力 action を無効にする。
- `observable point`: selected job summary が空または未確定の場合、出力 action が無効で理由を確認できること。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `data_requirement`, `recovery_requirement`
- `adoption hint`: plan の未決事項を人間判断候補として designer に渡すための異常系候補として扱える。
- `conflict hint`: 出力対象 job の自動選択を採るか、出力管理側で選ばせるかで、開始条件と観測点が変わる。

## Open Notes

- `human decision candidate`:
  - `CAND-TFN-007`: Body から `TranslationCompletePage` へ `Completed / Canceled / Failed` で遷移する状態図と、成果物出力が Completed job だけを対象にする仕様の見せ方。
  - `CAND-TFN-010`: 出力管理へ移動した後、出力対象 job を自動選択するか、出力管理側で選ばせるか。
  - `CAND-TFN-002`: 参照不能 job を一覧に残して理由表示するか、再開対象から除外して別表示にするか。
- `merge candidate`:
  - `CAND-TFN-001` と `CAND-TFN-003`: 直リンク防止と対象 job 未確定の復帰先。
  - `CAND-TFN-005` と `CAND-TFN-006`: 次工程へ進めない理由の `sticky footer` 表示。
  - `CAND-TFN-007`, `CAND-TFN-008`, `CAND-TFN-010`: 翻訳管理と成果物出力の責務分離。
- `rejection candidate`:
  - 正常系の「Job Setup 完了後に単語翻訳ページへ進む」だけを確認する候補は、failure 観点からは対象外にする。
  - 旧 `Job Run` のセッション取得を復活させる候補は、対象差分に反するため対象外にする。

## Evidence Summary

- `candidate_count`: 10
- `conflict_candidates`: `CAND-TFN-001`, `CAND-TFN-002`, `CAND-TFN-003`, `CAND-TFN-005`, `CAND-TFN-006`, `CAND-TFN-007`, `CAND-TFN-008`, `CAND-TFN-009`, `CAND-TFN-010`
- `human_decision_candidates`: `CAND-TFN-002`, `CAND-TFN-007`, `CAND-TFN-010`
- `required_basis`: 実行中タスク成果物場所は `docs/exec-plans/active/2026-05-08-translation-flow-navigation-overhaul/`、対象差分は翻訳フロー移動オーバーホール、候補成果物パスは `scenario-candidates.failure.md`、観点は failure。
