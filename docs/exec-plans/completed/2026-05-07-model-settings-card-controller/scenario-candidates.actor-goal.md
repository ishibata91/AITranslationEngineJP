# Scenario Candidates: 2026-05-07-model-settings-card-controller / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./task-frame.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `MSCC`

## Generator Scope

- `viewpoint`: actor-goal
- `included_sources`: `task-frame.md`, `light-change-planning.md`, `docs/spec.md`, `docs/architecture.md`, `docs/detail-specs/translation-job-setup.md`, `docs/detail-specs/ai-provider-settings-management.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本本文の変更、最終シナリオ表、候補の採否、候補の統合判断
- `generation_notes`: 利用者の目的、開始操作、成功判定から候補を分ける。状態遷移、失敗回復、外部連携の網羅は他観点へ残す。

## Candidate Scenarios

### CAND-MSCC-001 共通ペルソナ構築で保存済みモデル設定を再利用する

- `source requirement`: `task-frame.md` の「マスターペルソナと翻訳ジョブ設定が、同じモデル設定カード制御を使う。」、`ai-provider-settings-management.md` の「Job Setup と master-persona は provider settings を参照し、個別の secret や endpoint を fallback にしない。」、`docs/spec.md` の「目的に沿ったAIを選択可能であること。」
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-MSCC-001`
- `actor`: 共通ペルソナを構築したい利用者
- `trigger`: 利用者が master-persona 側のモデル設定カードを開く。
- `expected outcome`: 利用者は provider、model、保存済み選択状態を同じカード制御で確認できる。利用者は保存済み provider settings を参照し、secret や endpoint を個別 fallback として扱わずにモデルを選べる。
- `observable point`: カード表示に provider、model、モデル一覧状態、保存状態が出る。再読込後も保存済み provider と model が復元される。
- `related detail requirement type`: master-persona 参照側モデル選択、provider settings 参照、保存取得
- `adoption hint`: 採否判断なし。master-persona 側の共有カード利用は、両参照先を扱う候補として `designer` へ渡す。
- `conflict hint`: master-persona 固有の保存単位と、AIサービス設定の provider settings 保存単位を混同しない。

### CAND-MSCC-002 翻訳ジョブ設定で各翻訳段階のモデルを選び job 作成へ進む

- `source requirement`: `translation-job-setup.md` の「Job Setup は単語翻訳、NPC ペルソナ生成、本文翻訳の 3 つの翻訳段階を扱う。」、「各翻訳段階は provider、model、credential 参照、execution mode を持つ。」、「3 つの翻訳段階で APIキー不足と model 未選択がない時だけ、翻訳 job を作成できる。」
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-MSCC-002`
- `actor`: 翻訳 job を作成したい利用者
- `trigger`: 利用者が Job Setup で入力データを選び、翻訳段階ごとのモデル設定カードを操作する。
- `expected outcome`: 利用者は単語翻訳、NPC ペルソナ生成、本文翻訳の provider と model を選択できる。全段階で APIキー不足と model 未選択がなくなると、利用者は Ready job を作成できる。
- `observable point`: 各翻訳段階カードに AIサービス、model、APIキー状態、設定済みまたは未設定の状態が出る。作成前確認と作成後の設定内容に、段階ごとの AIサービスと model が出る。
- `related detail requirement type`: Job Setup phase runtime settings、model list 更新、model 選択、job 作成成立条件
- `adoption hint`: 採否判断なし。3 段階すべてを扱う候補として `designer` へ渡す。
- `conflict hint`: execution mode の選択目的は actor-goal の補助情報に留める。Batch API 条件の詳細は別観点で扱う。

### CAND-MSCC-003 provider 変更後に選択可能な model だけを保存する

- `source requirement`: `task-frame.md` の「provider 変更、model list 更新、model 選択、保存、取得、遅延応答破棄の責務境界が設計成果物で固定されている。」、`translation-job-setup.md` の「model 候補は provider ごとの model list API から取得する。」
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-MSCC-003`
- `actor`: 利用する AIサービスを切り替えたい利用者
- `trigger`: 利用者が共有モデル設定カードで provider を変更し、モデル一覧を更新する。
- `expected outcome`: 利用者は変更後 provider に対応する model list から model を選べる。利用者は変更前 provider の model を誤って保存しない。
- `observable point`: provider 変更後に model list 状態と model 選択状態が変更後 provider 基準で表示される。保存結果は変更後 provider と選択 model の組として観測できる。
- `related detail requirement type`: provider 変更、model list 更新、model 選択、保存
- `adoption hint`: 採否判断なし。provider 変更時に利用者が達成したい結果だけを `designer` へ渡す。
- `conflict hint`: actor-goal では利用者の成功体験だけを扱う。禁止遷移や stale response の詳細は state-transition 観点と競合しうる。

### CAND-MSCC-004 fake mode で通常 provider ID のまま fake-model を選ぶ

- `source requirement`: `task-frame.md` の「fake mode では、通常 provider ID のまま model list から `fake-model` を選べる。」、「`fake` provider ID を user-facing provider list へ追加しない。」、`light-change-planning.md` の「frontend に fake mode 判定や `fake-model` 固有分岐を追加しない。」
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-MSCC-004`
- `actor`: fake mode で UI を確認する利用者または人間レビュー担当者
- `trigger`: `AITRANSLATIONENGINEJP_MASTER_PERSONA_AI_MODE=fake` が設定された状態で、利用者が共有モデル設定カードの model list を更新する。
- `expected outcome`: 利用者は user-facing provider list に `fake` provider を見ない。利用者は通常 provider ID のまま返された `fake-model` を選び、保存または選択状態として扱える。
- `observable point`: provider list は `gemini`、`lm_studio`、`xai` の利用者向け provider だけを表示する。model list には `fake-model` が出る。frontend に fake 固有の分岐を示す UI は出ない。
- `related detail requirement type`: fake mode model list、user-facing provider list、frontend 責務境界
- `adoption hint`: 採否判断なし。fake mode で利用者に見える provider と model の候補として `designer` へ渡す。
- `conflict hint`: fake mode の切替位置と通常 provider ID の選定は、外部連携観点または実装範囲で追加判断が必要になる可能性がある。

### CAND-MSCC-005 遅い model list 応答があっても最後の操作結果を見る

- `source requirement`: `task-frame.md` の「provider 変更、model list 更新、model 選択、保存、取得、遅延応答破棄の責務境界が設計成果物で固定されている。」、`translation-job-setup.md` の「モデル一覧未更新、更新中、取得済み、取得失敗、APIキー未設定で更新不可を分けて表示する。」
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-MSCC-005`
- `actor`: 連続して provider 変更または model list 更新を行う利用者
- `trigger`: 利用者が model list 更新中に provider を変更する、または別の model list 更新を開始する。
- `expected outcome`: 利用者は最後に行った操作に対応する provider と model list 状態を見る。遅れて返った古い応答は、利用者の現在の選択状態を上書きしない。
- `observable point`: カードは更新中、取得済み、取得失敗の状態を最後の操作基準で表示する。古い応答の到着後も現在 provider と現在 model 選択が維持される。
- `related detail requirement type`: model list 更新、遅延応答破棄、選択状態表示
- `adoption hint`: 採否判断なし。actor-goal では利用者に見える成功判定だけを `designer` へ渡す。
- `conflict hint`: 状態遷移候補では、更新中に許可する操作と破棄条件をより厳密に扱う必要がある。

### CAND-MSCC-006 APIキー状態を再入力せずにモデル選択の可否を判断する

- `source requirement`: `translation-job-setup.md` の「API key が設定済みの場合だけ、API key が必要な provider の model list 外部取得を試みる。」、「model list 更新は、APIキーが必要な AIサービスで APIキー未設定の場合は押せない。」、`ai-provider-settings-management.md` の「APIキー本体と raw payload は、UI、DTO、要約、log、debug 出力に出さない。」
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-MSCC-006`
- `actor`: APIキーを再入力せずに provider と model を選びたい利用者
- `trigger`: 利用者が共有モデル設定カードで APIキーを必要とする provider を選ぶ。
- `expected outcome`: 利用者は provider settings の credential 参照状態から model list 更新可否を判断できる。利用者は APIキー本体をカード上で再入力または閲覧せずに、必要な場合だけ AIサービス設定側で設定が必要だと分かる。
- `observable point`: カードに APIキー状態、model list 更新可否、APIキー未設定で更新不可の表示が出る。APIキー文字列、secret、raw request、raw response は出ない。
- `related detail requirement type`: credential 参照、model list 更新可否、secret 非表示
- `adoption hint`: 採否判断なし。APIキー状態に基づく利用者の判断点だけを `designer` へ渡す。
- `conflict hint`: secret 非表示と raw payload 非露出は trust-boundary 観点と競合しうるため、最終統合時に重複整理が必要である。

## Open Notes

- `human decision candidate`: 通常 provider ID のまま `fake-model` を返す時、どの provider ID を代表入力にするか。
- `human decision candidate`: master-persona 側の model 保存単位を、共通ペルソナ構築全体の設定として扱うか、画面内のカード単位として扱うか。
- `human decision candidate`: Job Setup の 3 翻訳段階を、共有 controller の 3 インスタンスとして扱うか、phase settings collection として扱うか。
- `human decision candidate`: APIキー未設定時に共有カードから AIサービス設定画面へ遷移させるか、状態表示だけに留めるか。
