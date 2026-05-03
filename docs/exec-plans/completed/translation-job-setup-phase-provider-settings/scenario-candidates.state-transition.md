# Scenario Candidates: translation-job-setup-phase-provider-settings / state-transition

- `generator`: `state-transition`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJSPPS-ST`
- `candidate_artifact_path`: `docs/exec-plans/active/translation-job-setup-phase-provider-settings/scenario-candidates.state-transition.md`
- `candidate_count`: 7

## Generator Scope

- `viewpoint`: state-transition
- `task_artifact_location`: `docs/exec-plans/active/translation-job-setup-phase-provider-settings/`
- `target_diff`: Job Setup が phase ごとの provider、model、credential 参照、execution mode、batch mode、model 候補取得状態を持つ。
- `included_sources`:
  - `./plan.md`
  - `docs/exec-plans/completed/translation-job-setup/scenario-design.md`
  - `docs/exec-plans/completed/translation-job-setup/ui-design.md`
  - `docs/detail-specs/term-translation-phase.md`
  - `docs/detail-specs/persona-generation-phase.md`
  - `docs/detail-specs/body-translation-phase.md`
- `excluded_sources`: product code、product test、docs 正本、他 generator 成果物。
- `generation_notes`: 最終 scenario matrix、採否、統合、競合解消は designer に残す。

## Candidate Scenarios

### CAND-TJSPPS-ST-001 設定変更で validation を失効させる

- `source requirement`:
  - `plan.md:57-64`: Job Setup は phase ごとの provider、model、credential 参照、execution mode を持つ。
  - `docs/exec-plans/completed/translation-job-setup/scenario-design.md:57-63`: create 前 validation は必須設定不足、参照不能、provider / mode 不整合、credential 参照不能を blocking failure にする。
  - `docs/exec-plans/completed/translation-job-setup/ui-design.md:22-27`: 設定変更後は create job を無効にし、再 validation が必要であることを表示する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSPPS-ST-001`
- `actor`: Job Setup 利用者
- `trigger`: 任意の phase で provider、model、credential 参照、execution mode のいずれかを変更する。
- `pre-transition state`: Draft があり、直近 validation は pass かつ未失効である。
- `start condition`: create 前で、入力データと phase 別実行設定が選択済みである。
- `post-transition state`: validation は `dirty-validation` または失効状態になり、create job は無効になる。
- `expected outcome`: 古い validation 断面では `Ready` job を作成できない。再 validation 後だけ create 可否が再判定される。
- `observable point`:
  - validation summary に失効状態と再 validation 必要理由が表示される。
  - create job action が disabled になり、disabled 理由が表示される。
  - 最新設定断面で validation を再実行できる。
- `related detail requirement type`: `workflow`、`display`、`external_integration`
- `adoption hint`: validation stale 条件の主要候補として扱う。
- `conflict hint`: validation 表示を job 全体 1 件にするか、phase 別に表示するかは designer 統合時に整理する。

### CAND-TJSPPS-ST-002 phase 別設定の選択状態を混線させない

- `source requirement`:
  - `plan.md:8-10`: 単語翻訳、NPC ペルソナ生成、本文翻訳の各 phase が独立した実行設定を持つ。
  - `plan.md:59-60`: Job Setup は master-persona の provider 設定を既定値または保存元として扱わず、phase ごとの設定を持つ。
  - `docs/detail-specs/persona-generation-phase.md:31`: persona 生成は Job Setup の persona 専用設定を継承する。
  - `docs/detail-specs/body-translation-phase.md:27`: 本文翻訳は Job Setup の本文翻訳用 provider、model、execution mode を使う。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSPPS-ST-002`
- `actor`: Job Setup 利用者
- `trigger`: 単語翻訳 phase の provider、model、credential 参照、execution mode を変更する。
- `pre-transition state`: 3 phase の設定 Draft があり、それぞれ別の選択状態を持つ。
- `start condition`: 少なくとも 2 phase に異なる provider または model が設定されている。
- `post-transition state`: 変更した phase だけが新しい選択状態になる。他 phase の provider、model、credential 参照、execution mode は変わらない。
- `expected outcome`: phase 実行時に、対象 phase は対象 phase 専用の設定だけを参照する。
- `observable point`:
  - UI の phase 別設定欄で、変更対象外の phase の値が保持される。
  - validation 断面に 3 phase の設定が分離して記録される。
  - phase detail の provider / model / execution mode 要約が phase ごとに一致する。
- `related detail requirement type`: `persistence`、`external_integration`、フェーズ実行設定参照
- `adoption hint`: phase 間の状態分離候補として扱う。
- `conflict hint`: 保存構造や DTO 名は implementation-scope に残す。

### CAND-TJSPPS-ST-003 create 後の Ready 要約を read-only に固定する

- `source requirement`:
  - `docs/exec-plans/completed/translation-job-setup/scenario-design.md:71-77`: Ready job は再表示だけ許可し、入力、基盤参照、AI runtime、実行方式の再編集を許可しない。
  - `docs/exec-plans/completed/translation-job-setup/scenario-design.md:78-83`: UI は作成後の read-only 要約を表示する。
  - `docs/exec-plans/completed/translation-job-setup/ui-design.md:22-27`: Ready job は再表示だけ許可し、再編集導線は表示しない。
  - `docs/exec-plans/completed/translation-job-setup/ui-design.md:94-97`: success state では Ready job を read-only 要約として表示する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSPPS-ST-003`
- `actor`: Job Setup 利用者
- `trigger`: validation pass 後に create job を実行し、作成済み job を再表示する。
- `pre-transition state`: Draft は validation pass かつ未失効である。
- `start condition`: 同一入力に既存 job がなく、create job が許可されている。
- `post-transition state`: job は `Ready` になり、Job Setup は read-only 要約状態になる。
- `expected outcome`: 作成後の input、foundation、phase 別 provider、model、credential 参照、execution mode は再編集できない。
- `observable point`:
  - 入力、基盤参照、AI runtime、実行方式の編集 action が表示されない。
  - create result に job ID、`Ready` 状態、設定要約、validation 結果が表示される。
  - 再表示しても保存済み設定要約が変化しない。
- `related detail requirement type`: `persistence`、`display`
- `adoption hint`: create 後の禁止遷移候補として扱う。
- `conflict hint`: 作成済み job の再設定機能が将来追加される場合、別 task の状態遷移として扱う。

### CAND-TJSPPS-ST-004 model 候補取得状態を loading / success / failure / skipped で分ける

- `source requirement`:
  - `plan.md:61`: model 候補は provider ごとの model list API から取得し、API key が設定済みの場合だけ外部取得を試みる。
  - `plan.md:100`: designer 入力は model list 取得失敗時の UI 状態を含める。
  - `docs/exec-plans/completed/translation-job-setup/ui-design.md:28-34`: loading、success、error、disabled、retry、dirty-validation と blocking validation failure を区別する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSPPS-ST-004`
- `actor`: Job Setup 利用者
- `trigger`: 任意の phase で provider を選ぶ、credential 参照を変更する、または model 候補取得を再試行する。
- `pre-transition state`: 対象 phase の provider は未選択、または前回の model 候補取得状態を持つ。
- `start condition`: provider 選択後に model 候補の再評価が必要である。
- `post-transition state`:
  - `loading`: model list API の取得中である。
  - `success`: 選択中 provider に一致する model 候補を表示できる。
  - `failure`: 取得失敗理由を表示し、retry 可能性を示す。
  - `skipped`: 外部取得の前提がなく、model 候補取得を実行しない。
- `expected outcome`: model 候補取得状態が phase 別に観測でき、古い provider の model 候補を現在選択として扱わない。
- `observable point`:
  - 対象 phase の model 欄に loading、success、failure、skipped が表示される。
  - failure では secret や API key 平文を表示しない。
  - skipped では外部取得未実行の理由が表示される。
- `related detail requirement type`: `display`、`external_integration`
- `adoption hint`: model list API 状態遷移候補として扱う。
- `conflict hint`: `skipped` の理由を credential 未設定だけに限定するか、provider 側の model list 非対応にも使うかは designer 統合時に整理する。

### CAND-TJSPPS-ST-005 provider 変更時に batch mode 選択可否を遷移させる

- `source requirement`:
  - `plan.md:63-64`: batch mode は暗黙推定にせず、対象 provider は Gemini と xAI だけに限定する。
  - `plan.md:91-92`: Gemini / xAI batch mode 明示切替を must_include とし、他 generator 起動や product code 変更を禁止する。
  - `plan.md:100`: designer 入力は batch mode の対象 provider を含める。
  - `docs/exec-plans/completed/translation-job-setup/scenario-design.md:57-63`: provider / mode 不整合は blocking validation failure にする。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSPPS-ST-005`
- `actor`: Job Setup 利用者
- `trigger`: 任意の phase で provider を Gemini、xAI、または対象外 provider へ変更する。
- `pre-transition state`: 対象 phase に provider と execution mode の選択状態がある。
- `start condition`: batch mode の選択可否を再判定できる provider が選ばれている。
- `post-transition state`: Gemini または xAI では batch mode が選択可能になる。対象外 provider では batch mode が選択不能になり、既存の batch mode 選択は解除または invalid として扱われる。
- `expected outcome`: 対象外 provider では batch mode を選べず、古い batch mode 選択で validation pass にならない。
- `observable point`:
  - execution mode 欄で batch mode の enable / disabled が provider に連動する。
  - disabled 理由に対象 provider ではないことが表示される。
  - provider / mode 不整合は validation summary で blocking として表示される。
- `related detail requirement type`: `workflow`、`display`、`external_integration`
- `adoption hint`: batch mode の禁止遷移候補として扱う。
- `conflict hint`: batch mode 解除を自動解除にするか invalid 表示にするかは UI 統合時に整理する。

### CAND-TJSPPS-ST-006 LM Studio 選択時は credential 状態を skipped にする

- `source requirement`:
  - `plan.md:61-64`: model list 取得、LM Studio の API key 非要求、batch mode 対象 provider を分ける。
  - `plan.md:91`: must_include は LM Studio の API key 非表示を含める。
  - `docs/exec-plans/completed/translation-job-setup/ui-design.md:82-84`: credential 不備と secret 値の非表示を区別する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSPPS-ST-006`
- `actor`: Job Setup 利用者
- `trigger`: 任意の phase で provider として LM Studio を選択する。
- `pre-transition state`: 対象 phase の provider は未選択、または API key を必要とする provider である。
- `start condition`: provider 変更により credential 必要性を再評価できる。
- `post-transition state`: 対象 phase の credential 状態は `skipped` になり、API key 入力、API key 未設定 warning、credential select は表示されない。
- `expected outcome`: LM Studio の credential 未設定は blocking validation failure にならない。provider、model、execution mode の不整合は引き続き validation 対象になる。
- `observable point`:
  - LM Studio 選択時に API key 入力と credential select が表示されない。
  - validation summary に API key 未設定 warning が出ない。
  - 他 phase の credential 参照状態は変わらない。
- `related detail requirement type`: `display`、`external_integration`
- `adoption hint`: credential 不要 provider の状態遷移候補として扱う。
- `conflict hint`: LM Studio の model 候補取得を local endpoint へ試行するか、手入力または skipped とするかは designer 統合時に整理する。

### CAND-TJSPPS-ST-007 model 候補取得の遅延結果を現在状態へ混入させない

- `source requirement`:
  - `plan.md:59-61`: phase ごとの provider、model、credential 参照、execution mode と model list API 取得を扱う。
  - `docs/exec-plans/completed/translation-job-setup/ui-design.md:28-34`: loading、success、error、retry、dirty-validation と provider / mode 不整合を区別する。
  - `docs/detail-specs/body-translation-phase.md:72-73`: provider skipped、provider running、provider partial failure、late response rejected を区別して表示する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSPPS-ST-007`
- `actor`: Job Setup 利用者
- `trigger`: model 候補取得が `loading` の間に、同じ phase の provider または credential 参照を変更する。
- `pre-transition state`: 対象 phase の model 候補取得が `loading` で、取得要求の provider 断面が残っている。
- `start condition`: 利用者が取得完了前に provider または credential 参照を変更できる。
- `post-transition state`: 変更前 provider の遅延 success / failure は現在 phase の model 候補状態へ反映されない。現在 provider の model 候補状態だけが有効になる。
- `expected outcome`: 遅れて返った model list によって model 選択、validation 状態、batch mode 可否が巻き戻らない。
- `observable point`:
  - provider 変更後の model 欄は現在 provider の loading、success、failure、skipped だけを表示する。
  - 遅延結果は current selection を更新しない。
  - validation は失効状態のままで、古い取得結果によって pass に戻らない。
- `related detail requirement type`: `display`、`external_integration`、状態不変条件
- `adoption hint`: async 状態混入の禁止遷移候補として扱う。
- `conflict hint`: 遅延結果を audit / debug に残すかどうかは operation-audit 観点へ残す。

## Open Notes

- `human decision candidate`: なし。候補内の conflict hint は designer の統合判断候補であり、最終仕様として確定しない。
- `merge candidate`: `CAND-TJSPPS-ST-001` と `CAND-TJSPPS-ST-005` は validation stale と provider / mode 不整合で統合可能である。
- `merge candidate`: `CAND-TJSPPS-ST-002` と `CAND-TJSPPS-ST-007` は phase 間混線防止と async 状態混入防止で統合可能である。
- `rejection candidate`: なし。
