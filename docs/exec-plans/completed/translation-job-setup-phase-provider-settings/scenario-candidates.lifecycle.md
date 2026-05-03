# Scenario Candidates: translation-job-setup-phase-provider-settings / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJSPPS`

## Generator Scope

- `viewpoint`: lifecycle
- `included_sources`: `./plan.md`, `docs/exec-plans/completed/translation-job-setup/scenario-design.md`, `docs/exec-plans/completed/translation-job-setup/ui-design.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- `excluded_sources`: product code、product test、docs 正本変更、最終 scenario matrix、採否判断、統合判断
- `generation_notes`: Job Setup の phase 別 provider 設定が作成され、変更され、validation を通り、job create 後に read-only 再表示と phase 実行時参照へ進む流れだけを候補化する。

## Candidate Scenarios

### CAND-TJSPPS-LC-001 phase 別 provider 設定の Draft を作成する

- `source requirement`: `plan.md` の Design Requirements。Job Setup は phase ごとに provider、model、credential 参照、execution mode を持つ。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSPPS-LC-001`
- `lifecycle phase`: 設定作成
- `start condition`: `translation-input-intake` 完了後の入力があり、同一入力に既存 job がない状態で Job Setup を開く。
- `actor`: 翻訳ジョブを作成するユーザー
- `trigger`: Job Setup の新規 Draft を開始する。
- `expected outcome`: 単語翻訳、NPC ペルソナ生成、本文翻訳の 3 phase に独立した provider 設定欄が作られる。master-persona の provider 設定は既定値または保存元として使われない。
- `observable point`: Job Setup UI の phase 別 runtime 設定、validation 実行前状態、create job 無効状態。
- `related detail requirement type`: display、workflow、external_integration
- `adoption hint`: phase 別設定の最初の生成状態を扱う候補である。
- `conflict hint`: 初期値を完全空にするか、前回 Job Setup の値から復元するかは designer 統合時の判断候補である。

### CAND-TJSPPS-LC-002 phase 別 provider 設定を変更して validation を失効させる

- `source requirement`: `ui-design.md` の button_enablement。設定変更後は create job を無効にし、再 validation が必要であることを表示する。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSPPS-LC-002`
- `lifecycle phase`: 設定変更
- `start condition`: phase 別 provider 設定が入力済み、または validation pass 済みである。
- `actor`: 翻訳ジョブを作成するユーザー
- `trigger`: いずれかの phase で provider、model、credential 参照、execution mode、batch mode を変更する。
- `expected outcome`: 変更された phase の設定断面が dirty になり、直近 validation は失効する。create job は再 validation が通るまで実行できない。
- `observable point`: validation summary の dirty-validation 表示、create job 無効状態、変更された phase 名。
- `related detail requirement type`: workflow、display
- `adoption hint`: validation stale 条件を lifecycle 上の変更イベントとして扱う候補である。
- `conflict hint`: phase 1 件だけの変更で全体 validation を失効させるか、phase 単位 validation を持つかは designer 統合時の判断候補である。

### CAND-TJSPPS-LC-003 model list API 取得成功で model 候補を更新する

- `source requirement`: `plan.md` の Design Requirements。model 候補は provider ごとの model list API から取得する。API key が設定済みの場合だけ外部取得を試みる。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSPPS-LC-003`
- `lifecycle phase`: provider 選択後の model 候補取得
- `start condition`: API key が必要な provider を phase 設定で選び、その provider の credential 参照が設定済みである。
- `actor`: 翻訳ジョブを作成するユーザー
- `trigger`: provider 選択後に model list 取得を実行する。
- `expected outcome`: 選択 phase の model 候補が provider の model list API 結果で更新される。他 phase の model 選択と validation 断面は混線しない。
- `observable point`: model list の loading、success、候補件数、対象 phase 名、fake transport の request 証跡。
- `related detail requirement type`: external_integration、display
- `adoption hint`: model list API 成功時の設定 lifecycle 候補である。
- `conflict hint`: 取得成功後に既存 model 選択を維持するか初期化するかは designer 統合時の判断候補である。

### CAND-TJSPPS-LC-004 model list API 取得失敗から再取得待ちへ戻る

- `source requirement`: `ui-design.md` の state_variants。loading、error、retry を持ち、credential 解決、provider capability、ネットワーク到達性の失敗は blocking validation failure にする。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSPPS-LC-004`
- `lifecycle phase`: model 候補取得失敗
- `start condition`: API key が必要な provider の credential 参照が設定済みだが、model list API が失敗する。
- `actor`: 翻訳ジョブを作成するユーザー
- `trigger`: model list 取得が provider failure、network failure、invalid response のいずれかで失敗する。
- `expected outcome`: 失敗理由と retry 可能状態が表示される。model 候補が確定しない間は validation pass と create job へ進めない。
- `observable point`: model list error、retry action、validation failure reason、外部 request の redacted log。
- `related detail requirement type`: external_integration、workflow、display
- `adoption hint`: model list API 取得失敗後の lifecycle 停止点を扱う候補である。
- `conflict hint`: 失敗時に手入力 model を許可するか、API 取得成功だけを許可するかは designer 統合時の判断候補である。

### CAND-TJSPPS-LC-005 API key 未設定 provider では model list API を呼ばず credential 待ちへ遷移する

- `source requirement`: `plan.md` の Design Requirements。API key が設定済みの場合だけ外部取得を試みる。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSPPS-LC-005`
- `lifecycle phase`: API key 未設定時の provider 選択
- `start condition`: API key が必要な provider を phase 設定で選び、その provider の credential 参照がない。
- `actor`: 翻訳ジョブを作成するユーザー
- `trigger`: API key 未設定 provider を選択する。
- `expected outcome`: model list API は呼ばれない。credential missing として表示され、validation は blocking failure になる。credential 設定または provider 変更後に再取得へ進める。
- `observable point`: API key 未設定表示、model list 未実行証跡、validation summary、create job 無効状態。
- `related detail requirement type`: external_integration、workflow、display
- `adoption hint`: API key 未設定時の lifecycle 分岐を扱う候補である。
- `conflict hint`: API key 設定画面への導線を Job Setup 内に置くか、別設定画面へ送るかは designer 統合時の判断候補である。

### CAND-TJSPPS-LC-006 LM Studio は API key なしで provider 設定 lifecycle を進める

- `source requirement`: `plan.md` の Design Requirements。LM Studio は API key を要求しないため、API key 入力、API key 未設定 warning、credential select に出さない。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSPPS-LC-006`
- `lifecycle phase`: API key 不要 provider の設定
- `start condition`: いずれかの phase で LM Studio を provider として選ぶ。
- `actor`: 翻訳ジョブを作成するユーザー
- `trigger`: LM Studio の model 候補取得または validation を実行する。
- `expected outcome`: API key 入力、credential select、API key 未設定 warning は表示されない。credential missing は blocking failure にならず、LM Studio 固有の到達性または model 取得結果だけで状態が進む。
- `observable point`: LM Studio provider 表示、credential UI 非表示、API key warning 非表示、model list 状態、validation summary。
- `related detail requirement type`: external_integration、display、workflow
- `adoption hint`: LM Studio の API key なし状態を通常 provider と分けて扱う候補である。
- `conflict hint`: LM Studio の local endpoint 未設定を model list failure と同じ扱いにするか、別 failure category にするかは designer 統合時の判断候補である。

### CAND-TJSPPS-LC-007 Gemini の batch mode を明示的に切り替える

- `source requirement`: `plan.md` の Design Requirements。batch mode は暗黙推定にしない。対象 provider は Gemini と xAI だけに限定し、checkbox または select で明示する。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSPPS-LC-007`
- `lifecycle phase`: execution mode 変更
- `start condition`: Gemini を phase 設定の provider として選び、model 候補が選択済みである。
- `actor`: 翻訳ジョブを作成するユーザー
- `trigger`: Gemini の batch mode を on または off に切り替える。
- `expected outcome`: Gemini の batch mode 設定は明示値として phase 設定へ反映される。切り替え後は validation が失効し、再 validation が必要になる。
- `observable point`: Gemini の batch mode control、明示値、dirty-validation、create job 無効状態。
- `related detail requirement type`: workflow、external_integration、display
- `adoption hint`: Gemini の batch mode 明示切替を lifecycle 上の設定変更として扱う候補である。
- `conflict hint`: checkbox と select のどちらを使うかは UI design 側の判断候補である。

### CAND-TJSPPS-LC-008 xAI の batch mode を明示的に切り替える

- `source requirement`: `plan.md` の Design Requirements。batch mode は暗黙推定にしない。対象 provider は Gemini と xAI だけに限定し、checkbox または select で明示する。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSPPS-LC-008`
- `lifecycle phase`: execution mode 変更
- `start condition`: xAI を phase 設定の provider として選び、model 候補が選択済みである。
- `actor`: 翻訳ジョブを作成するユーザー
- `trigger`: xAI の batch mode を on または off に切り替える。
- `expected outcome`: xAI の batch mode 設定は明示値として phase 設定へ反映される。切り替え後は validation が失効し、再 validation が必要になる。
- `observable point`: xAI の batch mode control、明示値、dirty-validation、create job 無効状態。
- `related detail requirement type`: workflow、external_integration、display
- `adoption hint`: xAI の batch mode 明示切替を lifecycle 上の設定変更として扱う候補である。
- `conflict hint`: provider を Gemini または xAI から非 batch 対象へ変えた時に既存 batch 明示値を破棄するか保持するかは designer 統合時の判断候補である。

### CAND-TJSPPS-LC-009 validation pass 後に phase 別 provider 設定を含む job を作成する

- `source requirement`: `scenario-design.md` の REQ-TJS-001、REQ-TJS-002、REQ-TJS-003。validation pass なしで Ready job を作成しない。AI 基盤設定は provider、model、credential 参照、実行方式を区別する。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSPPS-LC-009`
- `lifecycle phase`: create
- `start condition`: 3 phase の provider、model、credential 参照、execution mode が valid であり、validation pass が未失効である。
- `actor`: 翻訳ジョブを作成するユーザー
- `trigger`: create job を実行する。
- `expected outcome`: `Ready` job が 1 件作成される。job は入力出自と validation pass 断面に加え、単語翻訳、NPC ペルソナ生成、本文翻訳の provider 設定要約を保持する。
- `observable point`: create result、Ready 状態、phase 別 provider/model/mode 要約、API key 平文非表示、job 未重複。
- `related detail requirement type`: workflow、persistence、external_integration
- `adoption hint`: create 時に phase 別 provider 設定を確定断面として保存する候補である。
- `conflict hint`: phase 別設定の保存単位と DTO 名は implementation-scope 側で固定する必要がある。

### CAND-TJSPPS-LC-010 作成済み Ready job を read-only で再表示する

- `source requirement`: `scenario-design.md` の REQ-TJS-004、REQ-TJS-005。Ready job は再表示だけ許可し、入力、基盤参照、AI runtime、実行方式の再編集は許可しない。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSPPS-LC-010`
- `lifecycle phase`: 終了後利用
- `start condition`: phase 別 provider 設定を含む Ready job が作成済みである。
- `actor`: 翻訳ジョブを確認するユーザー
- `trigger`: 作成済み job の Job Setup 表示を開き直す。
- `expected outcome`: 入力出自、validation 結果、単語翻訳、NPC ペルソナ生成、本文翻訳の provider 設定要約が read-only で表示される。再編集 action と create action は表示または実行されない。
- `observable point`: read-only summary、phase 別 provider/model/mode、credential 参照状態、API key 平文非表示、編集 action 非表示。
- `related detail requirement type`: display、persistence
- `adoption hint`: create 後の read-only 再表示を扱う候補である。
- `conflict hint`: read-only 再表示を Job Setup 画面に残すか、job detail 側へ寄せるかは designer 統合時の判断候補である。

### CAND-TJSPPS-LC-011 単語翻訳フェーズ開始時に単語翻訳用 provider 設定を参照する

- `source requirement`: `term-translation-phase.md`。主要データは provider / model / execution mode の要約を含む。provider 実行は phase run 内で扱う。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSPPS-LC-011`
- `lifecycle phase`: phase 実行開始
- `start condition`: Ready job があり、active な単語翻訳 phase run が存在しない。
- `actor`: 翻訳ジョブを実行するユーザー
- `trigger`: 単語翻訳フェーズを開始する。
- `expected outcome`: 単語翻訳 phase run は Job Setup で確定した単語翻訳用 provider、model、credential 参照、execution mode を読む。master-persona provider や他 phase の設定は参照しない。
- `observable point`: 単語翻訳 phase run summary、provider/model/mode 要約、credential 参照状態、provider request unit。
- `related detail requirement type`: workflow、external_integration、persistence
- `adoption hint`: job 作成後に単語翻訳 phase が phase 専用設定を使う流れの候補である。
- `conflict hint`: phase 開始時に provider 設定参照が失効していた場合の blocking 表示は failure 観点との統合候補である。

### CAND-TJSPPS-LC-012 NPC ペルソナ生成フェーズ開始時に persona 用 provider 設定を参照する

- `source requirement`: `persona-generation-phase.md`。provider、model、execution mode は Job Setup の persona 専用設定を継承する。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSPPS-LC-012`
- `lifecycle phase`: phase 実行開始
- `start condition`: 単語翻訳フェーズが Completed であり、active phase run がない。
- `actor`: 翻訳ジョブを実行するユーザー
- `trigger`: NPC ペルソナ生成フェーズを開始する。
- `expected outcome`: NPC ペルソナ生成 phase run は Job Setup で確定した persona 用 provider、model、credential 参照、execution mode を読む。単語翻訳用 provider と本文翻訳用 provider は参照しない。
- `observable point`: NPC ペルソナ生成 phase run summary、provider/model/mode 要約、credential ref、persona snapshot readiness。
- `related detail requirement type`: workflow、external_integration、persistence
- `adoption hint`: job 作成後に persona phase が phase 専用設定を使う流れの候補である。
- `conflict hint`: persona phase 実行時に provider 設定を再 validation するか、create 時の pass 断面だけを信用するかは designer 統合時の判断候補である。

### CAND-TJSPPS-LC-013 本文翻訳フェーズ開始時に本文翻訳用 provider 設定を参照する

- `source requirement`: `body-translation-phase.md`。本文翻訳フェーズは Job Setup で設定した本文翻訳用 provider、model、execution mode を使う。開始時の再選択 UI は作らない。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSPPS-LC-013`
- `lifecycle phase`: phase 実行開始
- `start condition`: NPC ペルソナ生成フェーズが Completed であり、辞書と persona snapshot の参照が成立している。
- `actor`: 翻訳ジョブを実行するユーザー
- `trigger`: 本文翻訳フェーズを開始する。
- `expected outcome`: 本文翻訳 phase run は Job Setup で確定した本文翻訳用 provider、model、credential 参照、execution mode を読む。開始時に provider 再選択 UI は出さない。
- `observable point`: 本文翻訳 phase run summary、provider/model/mode 要約、credential 参照状態、request unit count、output readiness。
- `related detail requirement type`: workflow、external_integration、persistence
- `adoption hint`: job 作成後に本文翻訳 phase が phase 専用設定を使う流れの候補である。
- `conflict hint`: body phase の batch mode と field 単位 request の関係は external-integration 観点との統合候補である。

## Open Notes

- `human decision candidate`: 初期 Draft 値を完全空にするか、過去の Job Setup 値から復元するか。
- `human decision candidate`: model list API 失敗時に手入力 model を許可するか、取得成功済み model だけを許可するか。
- `human decision candidate`: phase 実行開始時に provider 設定を再 validation するか、create 時の pass 断面だけを使うか。
- `merge candidate`: `CAND-TJSPPS-LC-002` は validation stale 条件として state-transition 観点と統合される可能性がある。
- `merge candidate`: `CAND-TJSPPS-LC-003`、`CAND-TJSPPS-LC-004`、`CAND-TJSPPS-LC-005`、`CAND-TJSPPS-LC-006` は model list / credential の external-integration 観点と統合される可能性がある。
- `merge candidate`: `CAND-TJSPPS-LC-011`、`CAND-TJSPPS-LC-012`、`CAND-TJSPPS-LC-013` は phase 実行シナリオの既存 detail-spec と統合される可能性がある。
