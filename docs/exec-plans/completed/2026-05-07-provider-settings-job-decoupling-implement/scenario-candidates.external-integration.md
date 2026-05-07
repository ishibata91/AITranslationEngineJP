# Scenario Candidates: 2026-05-07-provider-settings-job-decoupling-implement / external-integration

- `generator`: `external-integration`
- `source_plan`: `../2026-05-07-provider-settings-job-decoupling/plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `PSJD-EI`

## Generator Scope

- `viewpoint`: `external-integration`
- `included_sources`: `./task-frame.md`, `../2026-05-07-provider-settings-job-decoupling/plan.md`, `../2026-05-07-provider-settings-job-decoupling/light-change-planning.md`, `docs/detail-specs/ai-provider-settings-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/er.md`
- `excluded_sources`: 引き継ぎ入力に含まれない会話文脈、プロダクトコード、プロダクトテスト、docs 正本本文の変更判断
- `generation_notes`: 外部 provider、secret store、adapter、fake、network 境界だけを候補化する。採否、統合、競合解消は行わない。有料の実 AI API 呼び出しは必須検証にしない。

## Candidate Scenarios

### CAND-PSJD-EI-001 provider settings 保存を Job 所有情報から分離する

- `source requirement`: `task-frame.md:10`, `task-frame.md:16-18`, `light-change-planning.md:10-11`, `ai-provider-settings-management.md:27-32`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-PSJD-EI-001`
- `actor`: AIサービス設定を保存してから翻訳 job を作成する利用者。
- `external boundary`: provider settings 保存境界、secret store 保存境界、Job 作成時の provider settings 参照境界。
- `trigger`: 利用者が provider ごとの endpoint と APIキーを保存し、その provider を Job Setup の翻訳段階で選んで Ready job を作成する。
- `expected outcome`: provider settings row は endpoint と credential 参照状態を保持する。Job 側 DB は secret store 情報と endpoint を永続所有しない。Job 作成後の表示は AIサービス、model、APIキー状態、一括処理の有無だけを表示し、APIキー文字列と secret を表示しない。
- `fake_or_stub`: fake secret store で APIキー本体の保存有無だけを返す。provider settings 保存結果は fake transport に渡す前の要約で観測する。
- `observable point`: provider settings 保存要約、Job 作成応答、Job 側永続情報に endpoint と secret 本体が含まれないこと、UI/DTO/log に APIキー文字列が出ないこと。
- `related detail requirement type`: provider settings 保存、secret 非露出、Job Setup 参照契約、DB 永続境界。
- `adoption hint`: Job 側の endpoint / secret store 所有を外す主経路候補として扱える。
- `conflict hint`: `credential_ref` を完全削除するか監査用状態だけ残すかは未確定である。

### CAND-PSJD-EI-002 provider settings 未設定化が secret 本体を削除し参照側の fallback を禁止する

- `source requirement`: `ai-provider-settings-management.md:28-34`, `ai-provider-settings-management.md:49-50`, `translation-job-setup.md:61-65`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-PSJD-EI-002`
- `actor`: 保存済み provider settings を未設定へ戻す利用者。
- `external boundary`: secret store 削除境界、provider settings row 維持境界、Job Setup の参照境界。
- `trigger`: 利用者が APIキー保存済み provider を未設定へ戻し、その provider を Job Setup で選択する。
- `expected outcome`: provider settings row は残る。endpoint と APIキー状態は未設定に戻る。secret 本体は削除される。Job Setup は個別 secret や endpoint へ fallback せず、APIキー未設定と model list 更新不可を表示する。
- `fake_or_stub`: fake secret store で削除済み状態を返す。model list API は呼ばれない stub で観測する。
- `observable point`: provider settings の未設定状態、secret store の削除済み状態、model list 外部呼び出し回数 0、APIキー未設定表示。
- `related detail requirement type`: secret lifecycle、provider settings 未設定化、Job Setup 外部呼び出し抑止。
- `adoption hint`: 未設定化と参照側 fallback 禁止を 1 つの候補として統合できる。
- `conflict hint`: secret 削除後に既存 Running phase が開始時 snapshot を使うかどうかは別候補と競合しうる。

### CAND-PSJD-EI-003 model list API は provider settings と secret store を再解決して呼び出す

- `source requirement`: `translation-job-setup.md:38-40`, `translation-job-setup.md:60-66`, `ai-provider-settings-management.md:34`, `ai-provider-settings-management.md:38`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-PSJD-EI-003`
- `actor`: Job Setup で model list を更新する利用者。
- `external boundary`: model list API 境界、provider settings 参照境界、secret store 参照境界、fake transport DI 境界。
- `trigger`: Gemini または xAI の APIキー状態が設定済みで、利用者が Job Setup の model list 更新を実行する。
- `expected outcome`: model list 外部取得は provider settings の endpoint と credential 参照状態から解決される。取得成功時だけ model 選択を表示する。APIキー文字列、raw request、raw response、内部ログ用識別子は表示しない。
- `fake_or_stub`: fake secret store が APIキー存在状態を返す。fake transport DI が model list 成功応答を返す。有料の実 AI API は使わない。
- `observable point`: model list 更新中、取得済み、model 選択表示、fake transport 入力要約、raw payload 非露出。
- `related detail requirement type`: model list API、adapter 変換、secret 非露出、UI 状態表示。
- `adoption hint`: model list 外部取得の正規成功経路として扱える。
- `conflict hint`: provider settings 更新履歴を持たないため、取得時の revision 記録要否は designer 判断が必要である。

### CAND-PSJD-EI-004 APIキー未設定の provider は model list API を呼ばない

- `source requirement`: `translation-job-setup.md:39-43`, `translation-job-setup.md:61-65`, `ai-provider-settings-management.md:29`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-PSJD-EI-004`
- `actor`: APIキーが未設定の provider を Job Setup で選ぶ利用者。
- `external boundary`: secret store 参照境界、model list API 呼び出し可否境界。
- `trigger`: Gemini または xAI の APIキー状態が未設定で、利用者が該当 provider を翻訳段階に選ぶ。
- `expected outcome`: model list 更新は押せない。model list 外部取得は実行されない。翻訳 job 作成は APIキー不足と model 未選択が解消するまでできない。
- `fake_or_stub`: fake secret store が未設定状態を返す。model list API stub は未呼び出しを記録する。
- `observable point`: APIキー未設定表示、model list 更新不可表示、model list 外部呼び出し回数 0、job 作成不可状態。
- `related detail requirement type`: secret 前提条件、network 呼び出し抑止、Job Setup 作成可否。
- `adoption hint`: 失敗観点候補と統合可能だが、外部境界では「呼び出さないこと」を主観測点にする。
- `conflict hint`: APIキー不足の UI 表示文言は UI 設計側の最終表と競合しうる。

### CAND-PSJD-EI-005 LM Studio は secret store なしで model list API を扱う

- `source requirement`: `ai-provider-settings-management.md:26-32`, `translation-job-setup.md:39-41`, `translation-job-setup.md:61-65`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-PSJD-EI-005`
- `actor`: LM Studio を Job Setup の翻訳段階で使う利用者。
- `external boundary`: ローカル endpoint 境界、model list API 境界、secret store 非使用境界。
- `trigger`: 利用者が LM Studio の endpoint を provider settings に保存し、Job Setup で LM Studio の model list 更新を実行する。
- `expected outcome`: LM Studio では API key 入力、API key 未設定 warning、credential select を出さない。model list 外部取得は endpoint だけで実行される。endpoint は表示できるが secret は存在しない。
- `fake_or_stub`: ローカル endpoint 用 fake transport が model list 成功または接続失敗を返す。fake secret store は未使用であることを記録する。
- `observable point`: APIキー関連 UI の非表示、model list 取得状態、fake secret store 呼び出し回数 0、endpoint 参照要約。
- `related detail requirement type`: provider 別認証差分、adapter 分岐、ローカル endpoint 境界。
- `adoption hint`: Gemini / xAI と LM Studio の外部境界差分を明確にする候補として扱える。
- `conflict hint`: LM Studio の endpoint 参照不能時の表示は failure 観点候補と統合しうる。

### CAND-PSJD-EI-006 Ready job 実行開始時に最新 provider settings を再解決する

- `source requirement`: `task-frame.md:18`, `light-change-planning.md:28`, `ai-provider-settings-management.md:34-36`, `plan.md:29-30`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-PSJD-EI-006`
- `actor`: Ready job 作成後に AIサービス設定を更新してから翻訳実行を開始する利用者。
- `external boundary`: Ready job 実行開始時の provider settings 再解決境界、secret store 参照境界、provider execution adapter 境界。
- `trigger`: Ready job 作成後、利用者が provider settings の endpoint または APIキー状態を更新し、その後に翻訳実行を開始する。
- `expected outcome`: 実行開始前に最新 provider settings を再解決する。provider execution は再解決結果を使う。Job 側は個別 secret や endpoint を fallback にしない。
- `fake_or_stub`: fake secret store と fake provider execution transport で、開始時に解決された endpoint 要約と credential 存在状態を観測する。
- `observable point`: 実行開始時の provider settings 解決要約、provider execution fake transport の受信設定要約、Job 側 fallback 不使用、secret 非露出。
- `related detail requirement type`: Ready job 実行開始、provider execution、secret 参照、adapter 入力。
- `adoption hint`: Job 作成時保存ではなく実行開始時再解決を検証する中心候補として扱える。
- `conflict hint`: Job 作成時に provider settings revision を保持するかどうかは未確定である。

### CAND-PSJD-EI-007 Running phase は開始時 snapshot と provider settings 更新の扱いを外部境界として分ける

- `source requirement`: `task-frame.md:22-25`, `light-change-planning.md:29`, `light-change-planning.md:50-52`, `ai-provider-settings-management.md:33-36`, `er.md:22-25`, `er.md:63-66`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-PSJD-EI-007`
- `actor`: 翻訳実行中に AIサービス設定を変更する利用者。
- `external boundary`: Running phase の provider execution 境界、phase runtime snapshot 境界、provider settings 更新境界。
- `trigger`: 翻訳 phase が Running の間に、同じ provider の endpoint または APIキー状態が変更される。
- `expected outcome`: Running phase が共通設定更新へ追従するか、開始時 snapshot を継続するかを観測できる。phase runtime snapshot は保存してよい値と保存してはいけない値を分ける。APIキー本体と raw payload は保存しない。
- `fake_or_stub`: 長時間実行の fake provider execution transport を使う。fake secret store は更新前後の存在状態だけを返す。
- `observable point`: Running phase の provider execution 入力要約、phase runtime snapshot の保存項目、secret / raw payload 非保存、更新前後の設定利用差分。
- `related detail requirement type`: Running phase、runtime snapshot、provider settings 更新、secret 非露出。
- `adoption hint`: designer が snapshot 継続または追従の最終方針を選ぶための競合候補として残す。
- `conflict hint`: 既存仕様は Running phase が開始時 snapshot を使うと定義する一方、対象 task は共通設定更新への追従可否を重点論点にしている。

## Open Notes

- `human decision candidate`: Job 作成時に provider settings revision を保持するかは未確定である。
- `human decision candidate`: Running phase が provider settings 更新へ追従するか、開始時 snapshot を継続するかは未確定である。
- `human decision candidate`: `credential_ref` を完全削除するか、監査用状態だけ残すかは未確定である。
- `merge candidate`: `CAND-PSJD-EI-003` と `CAND-PSJD-EI-004` は model list API の成功経路と呼び出し抑止経路として統合候補である。
- `merge candidate`: `CAND-PSJD-EI-006` と `CAND-PSJD-EI-007` は実行開始時再解決と Running phase snapshot の統合候補である。
- `rejection candidate`: 有料の実 AI API 呼び出しを必須にする候補は除外する。
