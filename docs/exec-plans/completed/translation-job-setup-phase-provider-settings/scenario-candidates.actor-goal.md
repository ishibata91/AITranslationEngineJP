# Scenario Candidates: translation-job-setup-phase-provider-settings / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJSPPS`

## Generator Scope

- `viewpoint`: アクターの目的、開始操作、成功体験から候補を作る。
- `included_sources`: `./plan.md`, `../../completed/translation-job-setup/scenario-design.md`, `../../completed/translation-job-setup/ui-design.md`, `../../../detail-specs/term-translation-phase.md`, `../../../detail-specs/persona-generation-phase.md`, `../../../detail-specs/body-translation-phase.md`
- `excluded_sources`: product code、product test、docs 正本変更、final scenario matrix、採否判断、統合判断
- `generation_notes`: Job Setup 画面の利用者が、phase ごとの AI 実行設定を選び、validation と create の前に成功状態を確認できる候補だけを扱う。

## Candidate Scenarios

### CAND-TJSPPS-001 master-persona provider 設定に依存せず Job Setup を開始する

- `source requirement`: Job Setup は master-persona の provider 設定を既定値または保存元として扱わない。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TJSPPS-001`
- `actor`: Skyrim Mod 翻訳作業者
- `trigger`: Job Setup を開く。
- `expected outcome`: master-persona の provider / model / credential / execution mode が未設定または別値でも、Job Setup は phase 別の実行設定を選べる状態で表示される。
- `observable point`: Job Setup UI の phase 別 runtime selector。validation summary。create 前の設定要約。
- `related detail requirement type`: workflow、display、external_integration
- `adoption hint`: 既存 Job Setup の AI runtime 選択から、master-persona 依存を切り離す正常系として採用候補にする。
- `conflict hint`: 初期値を空にするか、前回 Job Setup Draft 相当から復元するかは designer 側で扱う。

### CAND-TJSPPS-002 単語翻訳フェーズ用の provider 設定を選ぶ

- `source requirement`: Job Setup は単語翻訳フェーズごとに provider、model、credential 参照、execution mode を持つ。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TJSPPS-002`
- `actor`: Skyrim Mod 翻訳作業者
- `trigger`: Job Setup で単語翻訳フェーズの runtime selector を操作する。
- `expected outcome`: 作業者は、単語翻訳フェーズで使う provider、model、credential 参照、execution mode を選び、共通辞書 hit 後に provider request へ送る対象語の実行設定を確認できる。
- `observable point`: 単語翻訳フェーズの設定欄。validation summary。create 後の phase 実行設定要約。
- `related detail requirement type`: external_integration、display、phase_input
- `adoption hint`: 単語翻訳フェーズは provider / model / execution mode の要約を phase result に表示するため、Job Setup 側の選択候補として独立させる。
- `conflict hint`: 単語翻訳フェーズの 1 対象語 1 request unit と Batch API 設定の扱いは、batch mode 候補と統合時に整理する。

### CAND-TJSPPS-003 NPC ペルソナ生成フェーズ用の provider 設定を選ぶ

- `source requirement`: Job Setup は NPC ペルソナ生成フェーズごとに provider、model、credential 参照、execution mode を持つ。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TJSPPS-003`
- `actor`: Skyrim Mod 翻訳作業者
- `trigger`: Job Setup で NPC ペルソナ生成フェーズの runtime selector を操作する。
- `expected outcome`: 作業者は、NPC ペルソナ生成フェーズで使う provider、model、credential 参照、execution mode を選び、persona snapshot 作成用の実行設定として確認できる。
- `observable point`: NPC ペルソナ生成フェーズの設定欄。validation summary。create 後の phase 実行設定要約。
- `related detail requirement type`: external_integration、display、phase_input
- `adoption hint`: NPC ペルソナ生成フェーズ仕様は Job Setup の persona 専用設定を継承するため、phase 別設定の必須候補にする。
- `conflict hint`: 共通ペルソナ hit で provider 未実行になる場合でも、設定を必須にするかは designer 側で整理する。

### CAND-TJSPPS-004 本文翻訳フェーズ用の provider 設定を選ぶ

- `source requirement`: Job Setup は本文翻訳フェーズごとに provider、model、credential 参照、execution mode を持つ。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TJSPPS-004`
- `actor`: Skyrim Mod 翻訳作業者
- `trigger`: Job Setup で本文翻訳フェーズの runtime selector を操作する。
- `expected outcome`: 作業者は、本文翻訳フェーズで使う provider、model、credential 参照、execution mode を選び、phase 開始時に再選択せず使われる設定として確認できる。
- `observable point`: 本文翻訳フェーズの設定欄。validation summary。create 後の phase 実行設定要約。
- `related detail requirement type`: external_integration、display、phase_input
- `adoption hint`: 本文翻訳フェーズ仕様は Job Setup で設定した本文翻訳用 provider / model / execution mode を使うため、独立候補にする。
- `conflict hint`: 本文翻訳フェーズの開始時再選択 UI を作らない前提と、Job Setup での変更可能期間を designer 側で一致させる。

### CAND-TJSPPS-005 API key 設定済み provider の model 候補を取得して選ぶ

- `source requirement`: model 候補は provider ごとの model list API から取得する。API key が設定済みの場合だけ外部取得を試みる。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TJSPPS-005`
- `actor`: Skyrim Mod 翻訳作業者
- `trigger`: API key 設定済み provider を phase の provider として選ぶ。
- `expected outcome`: 作業者は、provider の model list API から取得された model 候補を確認し、対象 phase の model として選べる。API key 平文は表示されない。
- `observable point`: model selector の候補一覧。credential 参照状態。model list 更新状態。validation summary。
- `related detail requirement type`: external_integration、display、secret_redaction
- `adoption hint`: provider / model 選択を手入力ではなく候補取得へ寄せる正常系として採用候補にする。
- `conflict hint`: model list API の取得失敗、stale 候補、再取得操作は failure または external-integration 観点と統合する。

### CAND-TJSPPS-006 LM Studio を credential なしの provider として選ぶ

- `source requirement`: LM Studio は API key 入力、API key 未設定 warning、credential select を出さない。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TJSPPS-006`
- `actor`: ローカル LLM を使う Skyrim Mod 翻訳作業者
- `trigger`: phase の provider として LM Studio を選ぶ。
- `expected outcome`: 作業者は、LM Studio の model または endpoint 由来の候補を選べる。画面には API key 入力、API key 未設定 warning、credential select が表示されない。
- `observable point`: LM Studio 選択時の provider 設定欄。credential 欄の非表示状態。warning 表示領域。validation summary。
- `related detail requirement type`: display、external_integration、secret_redaction
- `adoption hint`: API key provider と credential 不要 provider の UI 差分を利用者目的として切り出す。
- `conflict hint`: LM Studio の model 候補取得元、接続先設定、到達性 validation は external-integration 観点と統合する。

### CAND-TJSPPS-007 Gemini または xAI で batch mode を明示切替する

- `source requirement`: batch mode は暗黙推定にしない。対象 provider は Gemini と xAI だけに限定し、checkbox または select で明示する。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TJSPPS-007`
- `actor`: 大量レコードを効率よく処理したい Skyrim Mod 翻訳作業者
- `trigger`: phase の provider として Gemini または xAI を選ぶ。
- `expected outcome`: 作業者は、batch mode を使うかどうかを checkbox または select で明示的に切り替えられる。切替結果は対象 phase の execution mode として表示される。
- `observable point`: Gemini / xAI 選択時の batch mode 操作。execution mode 表示。validation summary。create 後の設定要約。
- `related detail requirement type`: external_integration、display、execution_mode
- `adoption hint`: Gemini / xAI の batch mode を利用者が意図して選ぶ正常系として採用候補にする。
- `conflict hint`: checkbox と select のどちらにするか、phase ごとに別設定にするか共通操作にするかは designer 側で扱う。

### CAND-TJSPPS-008 Gemini / xAI 以外では batch mode を選ばない

- `source requirement`: batch mode の対象 provider は Gemini と xAI だけに限定する。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TJSPPS-008`
- `actor`: Job Setup を確認する Skyrim Mod 翻訳作業者
- `trigger`: Gemini / xAI 以外の provider を phase の provider として選ぶ。
- `expected outcome`: 作業者は、対象外 provider で batch mode 操作を見ない。execution mode は対象 provider が対応する通常実行として確認できる。
- `observable point`: provider 切替後の execution mode 欄。batch mode 操作の非表示状態。validation summary。
- `related detail requirement type`: display、execution_mode、external_integration
- `adoption hint`: batch mode を暗黙推定せず、対象 provider だけに明示操作を出す境界候補として残す。
- `conflict hint`: 非対応 provider で既存 batch 設定が残っていた場合の stale 表示は state-transition または failure 観点と統合する。

### CAND-TJSPPS-009 phase 別設定を validation と create 前に確認する

- `source requirement`: create job 前に、AI runtime と実行方式の validation 結果を確認できる。AI 基盤設定は provider、model、credential 参照、実行方式を区別する。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TJSPPS-009`
- `actor`: Skyrim Mod 翻訳作業者
- `trigger`: 3 phase の provider / model / credential / execution mode を選んだ後に validation を実行する。
- `expected outcome`: 作業者は、単語翻訳、NPC ペルソナ生成、本文翻訳の各設定が validation 対象になったことを確認し、pass かつ未失効の場合だけ create job へ進める。API key 平文は表示されない。
- `observable point`: phase 別 validation summary。create button state。create 後の read-only 設定要約。secret redaction。
- `related detail requirement type`: workflow、display、secret_redaction、external_integration
- `adoption hint`: Job Setup の最終確認体験として、phase 別設定と既存 validation / create 条件をつなぐ候補にする。
- `conflict hint`: phase 別設定の一部だけ validation failure の時に create 全体を止めるかは designer 側で固定する。

## Open Notes

- `human decision candidate`: batch mode 操作を checkbox にするか select にするか。
- `human decision candidate`: phase 別設定の初期値を空にするか、前回 Job Setup の作業途中状態から復元するか。
- `human decision candidate`: 共通ペルソナ hit や辞書 hit で provider 未実行になる phase でも、provider 設定を必須にするか。
- `merge candidate`: model list API の取得失敗、credential 不足、stale 候補は failure / external-integration 観点と統合する。
- `merge candidate`: provider 切替後の validation 失効と batch mode stale 条件は state-transition 観点と統合する。
- `rejection candidate`: product code、product test、docs 正本変更、final scenario matrix の確定は本候補成果物では扱わない。
