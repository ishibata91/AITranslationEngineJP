# Scenario Candidates: term-translation-phase / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TTP`

## Generator Scope

- `viewpoint`: actor-goal
- `included_sources`: `./plan.md`, `tasks/usecases/term-translation-phase.yaml`, `tasks/index.yaml`, `docs/spec.md`, `docs/er.md`, `docs/architecture.md`, `docs/exec-plans/completed/translation-job-setup/plan.md`, `docs/exec-plans/completed/translation-job-setup/scenario-design.md`, `docs/exec-plans/completed/translation-job-setup/implementation-scope.md`
- `excluded_sources`: product code, product test, docs 正本更新, final scenario matrix, candidate adoption decision, other generator outputs
- `generation_notes`: 本文翻訳フェーズの前に、用語や固有名詞の訳語を確定し、ジョブ内辞書へ反映する成功体験だけを候補化する。状態遷移、失敗回復、外部 provider 失敗、監査保存の最終判断は designer と他観点に残す。

## Candidate Scenarios

### CAND-TTP-001 Ready job から単語翻訳フェーズを開始する

- `source requirement`: `term-translation-phase` の precondition は翻訳ジョブ作成完了であり、manual check は Job Run から単語翻訳フェーズを開始する。`translation-job-setup` は Ready 未成立では phase run を開始しない境界を固定している。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TTP-001`
- `actor`: 翻訳ユーザー
- `goal`: Ready job から単語翻訳フェーズを開始し、実行中であることを確認する。
- `trigger`: Job Run で Ready の翻訳ジョブを開き、単語翻訳フェーズ開始を実行する。
- `expected outcome`: 単語翻訳フェーズが現在フェーズとして開始され、progress を確認できる。
- `observable point`: Job Run の current phase、progress、phase run 開始結果。
- `related detail requirement type`: workflow, display
- `adoption hint`: phase 開始の代表正常系として採用候補にできる。
- `conflict hint`: Ready 判定、Running 遷移、phase run 作成条件は state-transition 観点と統合する。

### CAND-TTP-002 共通辞書の完全一致語を翻訳対象から除外する

- `source requirement`: `docs/spec.md` は共通辞書として登録済みの対象を翻訳対象とせず置き換えること、完全一致のみを構築済みとして扱うこと、内部出力ステータスとして `cached` を保持できることを求める。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TTP-002`
- `actor`: 翻訳ユーザー
- `goal`: 共通辞書に既にある語を再翻訳せず、一貫した訳語で置き換える。
- `trigger`: 共通辞書と翻訳対象語を持つ Ready job で単語翻訳フェーズを実行する。
- `expected outcome`: 共通辞書に完全一致する語は AI 翻訳対象から外れ、置換対象として結果に反映される。
- `observable point`: phase result の除外件数、置換対象の判定結果、`cached` 相当の内部観測情報。
- `related detail requirement type`: workflow, persistence
- `adoption hint`: 共通辞書再利用の主要成功系として採用候補にできる。
- `conflict hint`: 完全一致の正規化範囲、大小文字、空白、記号差の扱いは failure または state-transition 観点と競合しうる。

### CAND-TTP-003 用語や固有名詞の訳語を確定する

- `source requirement`: `term-translation-phase` の goal は本文翻訳フェーズの前に、用語や固有名詞の訳語を確定することである。`docs/spec.md` は単語翻訳フェーズを、単語や固有名詞を個別に翻訳し訳語を辞書化するフェーズとして定義している。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TTP-003`
- `actor`: 翻訳ユーザー
- `goal`: 共通辞書にない用語や固有名詞の訳語を本文翻訳前に確定する。
- `trigger`: 共通辞書に未登録の翻訳対象語を含む Ready job で単語翻訳フェーズを実行する。
- `expected outcome`: 未登録の用語や固有名詞について確定訳語が生成される。
- `observable point`: phase result の確定訳語一覧、原語と訳語の対応、対象語の処理件数。
- `related detail requirement type`: operation, display
- `adoption hint`: 単語翻訳フェーズの中心目的として採用候補にできる。
- `conflict hint`: AI 自動確定か人間確認を挟むかは対象差分だけでは確定できない。

### CAND-TTP-004 確定訳語をジョブ内辞書へ反映する

- `source requirement`: `term-translation-phase` の output はジョブ内辞書と確定訳語である。`docs/er.md` は `DICTIONARY_ENTRY` が共通辞書とジョブ内辞書を同じテーブルで扱い、ジョブ内データでは `translation_job_id` を持つと定義している。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TTP-004`
- `actor`: 単語翻訳フェーズ処理
- `goal`: 確定訳語を対象ジョブのジョブ内辞書として後続フェーズへ残す。
- `trigger`: 単語翻訳フェーズで訳語が確定する。
- `expected outcome`: 確定訳語が対象ジョブのジョブ内辞書として保存される。
- `observable point`: ジョブ内辞書の辞書項目、対象 job ID、作成経路、phase result。
- `related detail requirement type`: persistence
- `adoption hint`: 本文翻訳フェーズへの再利用を成立させる永続化成功系として採用候補にできる。
- `conflict hint`: 辞書項目の一意性、既存ジョブ内辞書との上書き可否、flush 時期は state-transition または lifecycle 観点と統合する。

### CAND-TTP-005 確定訳語を本文翻訳フェーズの入力として参照する

- `source requirement`: `term-translation-phase` の completion criteria は確定訳語を本文翻訳フェーズの入力として参照できることである。`docs/spec.md` は本文翻訳フェーズが単語翻訳フェーズで確定した訳語を再利用すると定義している。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TTP-005`
- `actor`: 本文翻訳フェーズ処理
- `goal`: 本文翻訳フェーズで単語翻訳フェーズの確定訳語を入力として使う。
- `trigger`: 単語翻訳フェーズが完了し、本文翻訳フェーズの入力準備が始まる。
- `expected outcome`: 本文翻訳フェーズは対象ジョブのジョブ内辞書を入力として参照できる。
- `observable point`: 本文翻訳フェーズ入力 summary、参照辞書項目数、対象 job ID の一致。
- `related detail requirement type`: workflow, persistence
- `adoption hint`: phase 間 handoff の正常系として採用候補にできる。
- `conflict hint`: 本文翻訳フェーズ側の final scenario と重複しうるため、designer が境界を調整する。

### CAND-TTP-006 単語翻訳フェーズ結果を Job Run で確認する

- `source requirement`: `term-translation-phase` の completion criteria と manual check は progress、phase result、確定訳語、ジョブ内辞書を確認することである。`docs/spec.md` は翻訳に利用する辞書と共通基盤データを UI から観測可能にすることを求める。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TTP-006`
- `actor`: 翻訳ユーザー
- `goal`: 単語翻訳フェーズの結果を Job Run で確認し、後続作業へ進める状態を判断する。
- `trigger`: 単語翻訳フェーズの実行結果が返る。
- `expected outcome`: ユーザーは current phase、progress、phase result、確定訳語、ジョブ内辞書反映を確認できる。
- `observable point`: Job Run の phase result、progress 表示、確定訳語一覧、ジョブ内辞書 summary。
- `related detail requirement type`: display
- `adoption hint`: UI 観測の代表成功系として採用候補にできる。
- `conflict hint`: Job Run 画面の詳細 UI 契約は ui-design で固定する。

### CAND-TTP-007 再利用語を本文翻訳で置換対象として扱う

- `source requirement`: `docs/spec.md` は用語翻訳で翻訳済みの対象を翻訳対象とせず置き換えること、再利用語を会話文やクエスト文へ流用する対象として定義している。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TTP-007`
- `actor`: 翻訳ユーザー
- `goal`: 確定済みの再利用語を本文翻訳で再翻訳せず、置換対象として使う。
- `trigger`: 単語翻訳フェーズで本文内に出現する用語の訳語が確定する。
- `expected outcome`: 確定済みの再利用語は本文翻訳で再翻訳されず、置換対象として扱われる。
- `observable point`: phase result の置換対象判定、再利用語フラグ、本文翻訳入力の辞書参照 summary。
- `related detail requirement type`: workflow, display
- `adoption hint`: 一貫した単語訳というユーザー価値を直接示す候補として採用候補にできる。
- `conflict hint`: 置換実行を単語翻訳フェーズで行うか本文翻訳フェーズで行うかは designer が統合時に決める。

### CAND-TTP-008 複数ジョブのジョブ内辞書を混線させない

- `source requirement`: `docs/spec.md` は複数の入力データを独立した翻訳ジョブとして管理することを求める。`docs/er.md` はジョブ内辞書のときだけ `translation_job_id` を設定すると定義している。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TTP-008`
- `actor`: 翻訳ユーザー
- `goal`: 複数ジョブで単語翻訳フェーズを実行しても、ジョブ内辞書を対象ジョブごとに分離する。
- `trigger`: 複数の Ready job で、それぞれ単語翻訳フェーズを実行する。
- `expected outcome`: 各ジョブのジョブ内辞書は対象ジョブだけに紐づき、別ジョブの本文翻訳入力へ混入しない。
- `observable point`: job A と job B のジョブ内辞書 summary、対象 job ID、本文翻訳入力の参照辞書。
- `related detail requirement type`: persistence
- `adoption hint`: 複数ジョブ管理時の actor 目的として採用候補にできる。
- `conflict hint`: DB 制約、参照整合性、並行実行中の分離は state-transition または operation-audit 観点と統合する。

### CAND-TTP-009 共通辞書とジョブ内訳語の重複時に結果を確認する

- `source requirement`: `docs/spec.md` は共通辞書として登録済み、または用語翻訳で翻訳済みの対象を翻訳対象とせず置き換えることを求める。共通辞書とジョブ内辞書の優先順位は指定資料だけでは確定しない。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TTP-009`
- `actor`: 翻訳ユーザー
- `goal`: 共通辞書とジョブ内辞書の候補が重複した時、本文翻訳へ渡る訳語を確認する。
- `trigger`: 同じ原語について共通辞書とジョブ内辞書の候補が存在する。
- `expected outcome`: ユーザーはどちらの訳語が本文翻訳に使われるかを確認できる。
- `observable point`: phase result の採用訳語、辞書 source、置換対象判定。
- `related detail requirement type`: workflow, display
- `adoption hint`: human decision 後に、辞書優先順位の確認シナリオとして採用候補にできる。
- `conflict hint`: 優先順位の業務判断が未確定であり、designer は人間判断候補として扱う必要がある。

### CAND-TTP-010 翻訳対象語が全て共通辞書で解決済みでも完了する

- `source requirement`: `docs/spec.md` は共通辞書に登録済みの対象を翻訳対象としないことを求める。`term-translation-phase` の completion criteria は progress と phase result、確定訳語とジョブ内辞書の確認である。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TTP-010`
- `actor`: 翻訳ユーザー
- `goal`: 新規訳語が不要なジョブでも単語翻訳フェーズを完了させ、本文翻訳へ進める。
- `trigger`: 翻訳対象語が全て共通辞書の完全一致で解決する Ready job で単語翻訳フェーズを実行する。
- `expected outcome`: AI に新規訳語を依頼せず、単語翻訳フェーズが完了し、本文翻訳フェーズへ進める。
- `observable point`: phase result の新規訳語 0 件、共通辞書置換件数、current phase / progress、次フェーズ開始可否。
- `related detail requirement type`: workflow, display
- `adoption hint`: 代替成功系として採用候補にできる。
- `conflict hint`: 0 件時にジョブ内辞書を作成するか、共通辞書参照だけにするかは lifecycle または persistence 観点と統合する。

## Open Notes

- `human decision candidate`: `CAND-TTP-003` の訳語確定に人間確認を挟むか、AI 結果を自動確定するか。
- `human decision candidate`: `CAND-TTP-009` の共通辞書とジョブ内辞書が同じ原語を持つ場合の優先順位。
- `human decision candidate`: `CAND-TTP-002` の完全一致判定で大小文字、空白、記号差、正規化をどこまで許可するか。
- `human decision candidate`: `CAND-TTP-010` の新規訳語 0 件時にジョブ内辞書 row を作るか、共通辞書参照だけで完了させるか。
- `merge candidate`: `CAND-TTP-005` と `CAND-TTP-007` は本文翻訳フェーズ側の入力参照シナリオと統合候補である。
- `rejection candidate`: `CAND-TTP-008` は state-transition / persistence 側で十分に覆われる場合、actor-goal 側では補助候補に下げられる。
