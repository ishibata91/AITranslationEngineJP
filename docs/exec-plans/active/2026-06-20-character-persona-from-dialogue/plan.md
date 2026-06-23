# Task Plan: 2026-06-20-character-persona-from-dialogue

- `workflow`: work
- `status`: スライス 1・2・4 完了。既知の問題 1〜3（性別・年齢の役割語、一人称・語尾、保留較正）を実装・検証・実画面目視まで完了。問題 4 は許容誤差。スライス 3 未着手（2026-06-23）
- `task_id`: 2026-06-20-character-persona-from-dialogue
- `request_summary`: 話者の口調の素（性質文・ペルソナ）の持ち方を仕様から設計し直す。種族・汎用声型には固定のペルソナを割り当て、キャラ専用声型はそのキャラの台詞群から個別ペルソナを生成する。
- `goal`: 翻訳の口調指示に差し込む性質文・ペルソナを、抽出データ由来の属性へ割り当て・生成でき、再抽出後も保持し、翻訳実行へ反映できる。
- `source_branch`: `master`
- `source_commit`: `6721b7a7`
- `work_branch`: `claude/2026-06-20-character-persona-from-dialogue`
- `target_branch`: `master`

## 完了定義（preparation-module で固定、2026-06-20）

本 task は、次の 4 つの振る舞いが観測できる状態を「完了」とする。差込点や空テーブルを置くだけで振る舞いが観測できない状態は「動く」と書かない。仕様の細部（生成方式、永続化テーブル構造、プリセット粒度など）は design-module で確定するが、完了定義の振る舞い自体は狭めない。

動かす範囲（task 後に検証者が観測できる振る舞い）:

1. 種族・汎用声型に対し、UI 上で性質文・ペルソナを割り当て・編集でき、保存される。
2. キャラ専用声型に対し、そのキャラの台詞群からペルソナを生成でき、生成結果が UI に表示され、保存される。
3. 割り当て・生成したペルソナが、抽出データ本体とは別のテーブルへ属性 EditorID 紐付けで永続化され、再抽出（属性一覧の入れ替え）後も保持される。
4. 割り当て・生成したペルソナが、翻訳実行時の口調指示へ実際に差し込まれ、翻訳結果の口調へ反映される。

観測点:

- 1, 2 → 実画面。UI を操作してペルソナの割り当て・生成・表示を確認する。
- 3 → 実データ。再抽出を実行し、永続化テーブルにペルソナが残ることを確認する。
- 4 → 実画面と実データ。翻訳を実行し、実プロンプトへ口調指示が差し込まれ、口調が割り当て・生成内容どおりに変わることを確認する。口調指示の合成・差し込みに不変ルールがある場合は、純粋 IO クラスへ分離し単体テストで確認する。

## close_conditions（観測点で検証する）

- `close-1`: 実画面で、種族・汎用声型に性質文・ペルソナを割り当て・編集・保存できる。
- `close-2`: 実画面で、キャラ専用声型のペルソナをそのキャラの台詞群から生成し、表示・保存できる。
- `close-3`: 実データで、再抽出後も割り当て・生成したペルソナが属性 EditorID 紐付けで残る。
- `close-4`: 実画面と実データで、翻訳実行時の口調指示へ割り当て・生成したペルソナが差し込まれ、口調へ反映される。口調合成の不変ルールは純粋 IO クラスの単体テスト（カバレッジ 100%）で確認する。

## 軽 / 重判定（preparation-module で固定、2026-06-20）

- 画面が動くか: Y。種族・声型へのペルソナ割り当て編集 UI と、キャラ専用声型の生成・表示 UI を新設する。layout・文言・表示構造・svelte 表示コンポーネントを変える。
- `docs/architecture.md` 反映が要るか: Y（design-module で確定）。性質文・ペルソナの専用永続化テーブル新設、台詞群からのペルソナ生成経路（LLM 呼び出し）、`internal/engine/persona_rule.go` のデータ駆動化が、永続化層・依存方向・Wails 境界（新 binding）を変える可能性が高い。最終要否は design-module の入口で構造変化の有無を見て確定する。
- 判定結果: 重 task（画面 Y だけで重 task 確定）。経路は `preparation-module` → `design-module` → `storybook-module`（画面が動くため）→ `implementation-module` → `finalization-module`。

## 背景（T4 storybook レビューで判明、2026-06-20、人間指示）

T4（`2026-06-14-prompt-persona-customization`）の storybook レビュー中に、ルールの持ち方を見直す指示が出た。要点は次の 4 点とする。

- 属性キーは抽出データ由来とする。ユーザーは属性キーを追加・削除せず、抽出された一覧へ性質文（口調の素）を割り当てる。
- 性質文は属性 EditorID に紐付けた別テーブルへ永続化する。再抽出で属性一覧が入れ替わっても性質文を残す。
- 口調の素は種族（抽出 RACE）と声型に絞る。勢力 FACT は数百で個別管理が非現実的のため、話者のメタデータ（参照表示）に回す。
- 声型は 2 層で持つ。汎用声型（FemaleNord など多数 NPC 共有）は声型名のパターングループへ口調ペルソナを割り当てる。キャラ専用声型（DLC1SeranaVoice など単一キャラに 1:1）は個別ペルソナを割り当てる。根拠: `docs/skyrim-structure-model.md` L237 が VoiceType を汎用人間・生物・キャラ専用の 3 種に分ける。

さらに人間が次を指示した。キャラ専用声型の個別ペルソナは、ユーザーが手で書くのではなく、そのキャラの台詞群から生成したくなる。生成は仕組みが大きいため、T4 から外して本 plan で仕様から設計し直す。

## T4 からの切り出し範囲

T4 から本 plan へ移す対象は次とする。

- 属性（種族・声型）→ 性質文・ペルソナの割り当て編集 UI（旧 T4 Slice 3「ルール編集・永続化・反映」）。
- 性質文・ペルソナの専用テーブルへの永続化と、`internal/engine/persona_rule.go` のハードコードのデータ駆動化。
- キャラ専用声型ペルソナの台詞群からの生成（新規）。

T4 に残す対象は次とする。本 plan の対象外とする。

- プロンプトテンプレート編集 UI（base 翻訳指示文と口調指示テンプレートの編集）。
- 実プロンプト参照（翻訳実行後の結果行で実プロンプトを目視確認）。
- 機械置換内訳（原語 → 確定訳語）の結果行表示。

## Scope 素案（仕様で確定する）

含む（候補）:
- 種族・汎用声型への固定ペルソナ割り当て（数種類のプリセットから選ぶか、性質文を書くか、仕様で決める）。
- キャラ専用声型ペルソナの台詞群からの生成。
- 性質文・ペルソナの専用テーブルへの永続化（再抽出耐性）。
- 割り当て・生成結果の翻訳実行への反映。

含まない（候補）:
- T4 の対象（テンプレート編集・実プロンプト参照・機械置換内訳）。
- 勢力 FACT への個別ペルソナ割り当て（メタデータ参照表示にとどめる）。

## 確定した口調固定ルール（仕様検討、2026-06-21）

人間との仕様検討で合意した、口調付けの固定方針を記録する。口調種別（基底口調の軸）自体は未確定で、`tone-concept-model.md` の概念分析から起こす。

- R1 機械分類を第一層にする。台詞の言語特徴から口調を決定的に分類する。LLM 生成は補助に降ろし、機械分類が拮抗するキャラと人手指定キャラだけに使う。
- R2 口調を 2 レイヤーに分ける。基底口調（排他・単一選択、一人称と文末を決める）と、態度修飾（合成・上位 N、語尾を変えない形容）に分離する。
- R3 基底口調の一人称・文末は自前の few-shot 例（英語原文 → その口調の訳文）で固定する。性質文の説明でなく例示で語尾を決定的に固定し、テスト可能にする。
- R4 態度修飾の入力辞書は Empath（MIT ライセンスの語彙カテゴリ辞書）と NRC（感情語彙辞書）を取り込む。Skyrim 固有語・古英語・register（フォーマリティ段階）は自作の小辞書で補う。
- R5 出力類型の体系は役割語（金水敏の言語ステレオタイプ理論）を土台にする。戯画化を避け、語尾固定は穏当にする。
- R6 職業は口調そのものでなく、口調軸の値を後天的に動かす決定要素の一つとして扱う。種族・性別・年齢・階層・信仰・性格も同じ立場の決定要素にする。
- R7 種族訛り（Khajiit の三人称自称など）は種族 EditorID から固定付与する別レイヤーにし、基底口調軸と直交させる。
- R8 分類の入力は英語原文の台詞群とし、翻訳前に行う。訳文からの口調抽出は鶏卵関係になるため却下する。
- R9 ペルソナ・ロールは数百でなく、ゲーム的に自然な十数種程度に厳選する。
- R10 口調分類の不変ルールは単一の純粋 IO クラスへ分離し、ユニットテストカバレッジ 100% を基準にする。
- R11 口調軸は対人態度軸（謙る・丁寧 ⇔ 粗野・尊大、一人称と敬語を決める）と感情表出軸（抑制・冷静 ⇔ 激情・熱、文末の強さを決める）の 2 軸に確定する。裏付けは Halliday の tenor と Appraisal の affect。古風さと種族訛りは軸でなく語彙マーカー層として基底口調へ重ねる。詳細は `tone-concept-model.md`。

## 仕様で決める未決事項（2026-06-21 更新）

- 各口調軸の段階数。対人態度軸と感情表出軸をそれぞれ何段階に切るか（3×3=9 か 3×2=6 か）。詳細は `tone-concept-model.md`。
- 語彙マーカー（古風・種族訛り）の初期範囲。Skyrim 世界観でどこまでを v1 に入れるか。
- 基底口調の類型数と線引き（十数種の具体）。口調軸の確定後に決める。
- 各基底口調の few-shot 例文の中身。
- 役割語類型のうち Skyrim 世界観で採用する範囲。
- 永続化テーブルの構造と、抽出データ本体との分離方法（属性 EditorID 紐付け、再抽出耐性）。
- キャラ専用声型の台詞が `line_speaker`（名指し話者の連関）で十分に引けるか。実データ検証が要る。引けない場合の extractor 改修の可否。
- LLM 補助層を使う場合の、手修正と再生成の優先順位。

## 関連 task

- T2（ペルソナ機構）。口調指示の合成経路。
- T3（マスター辞書）。固有名の機械置換。
- T4（`2026-06-14-prompt-persona-customization`）。テンプレート編集・実プロンプト参照・機械置換内訳。本 plan は T4 のルール編集を引き継ぐ。

## 進捗（仕様検討・PoC、2026-06-21）

### 到達点

- 口調の概念モデルを確定した（`tone-concept-model.md`）。口調 = 基底口調（排他・キャラ固定）＋ 態度修飾（合成・台詞ごと）＋ 語彙マーカー（軸の外の固定型）。口調軸 = 対人態度 × 感情表出の 2 軸。裏付けは Halliday の tenor と Appraisal の affect。基底口調は 2 軸 3 段の 3×3 グリッド。
- 確定した口調固定ルール R1〜R11 を本書上部に記録した。
- 口調抽出 PoC を実装・実行した（`cmd/poc-tone`、結果は `poc-tone-report.md` v2）。skyrim.esm の実台詞（`db/poc-skyrim.sqlite3`、34427 台詞）で 2 軸抽出を検証した。
- PoC v2 で外部辞書 NRC・命令文検出（prose 品詞解析）・マーカーあり台詞比率を導入した。対人軸は「印」（マーカーを含む台詞数）が十分なキャラで機能した（Ulfric・Tullius・Astrid が尊大、Farengar・Paarthurnax が丁寧）。感情軸・分散・メタ語彙マーカーは機能した。
- 無作為検証を実施した（`go run ./cmd/poc-tone verify`、結果は `poc-tone-report.md` の「無作為検証」節）。手選びを外し、台詞 100 以上の信頼層 31 人を悉皆で対象にした。印 10 以上の 27 人で中立に潰れたのは 22% にとどまり、対人軸の判別が恣意選択なしでも一般へ及んだ。学者・神官は丁寧、冷血幹部（Delphine・Astrid・Mercer）は冷然へ妥当に分かれた。
- 無作為検証で新しい失敗を 2 つ見つけ、両方を直して再検証した。失敗 1（誘導命令の尊大誤判定）は包括命令 `let's` と道案内の定型（`come on`・`this way` 等）を威圧から外して解消した（Hadvar −0.50→+0.20、Ralof −0.64→−0.14。Ulfric・Tullius・Astrid の本物の威圧は保持）。失敗 2（Khajiit 三人称の誤爆）は種族レコードで gate して解消した（非 Khajiit 10 人から消え、実 Khajiit の Jzargo らは保持）。詳細は `poc-tone-report.md` の「修正と再検証」節。
- メタデータ信号を棚卸しした（`poc-tone-report.md` の「メタデータ信号の棚卸し」節）。voice_type の EditorID が全 912 話者の 786（86%）に temperament を符号化している（Brute・Coward・SlyCynical・OldKindly/OldGrumpy・Commander 等）。現状コードは Old/Child/Young/Condescending しか拾わず、温厚老人と気難老人を潰している。faction 所属（843 話者）も役割信号として未使用。occupation（class）と AI データ（aggression）は extractor が未抽出。空の `nature` 列は含意注釈の curation slot。
- voice_type の役割を訂正した。「汎用ボイス 86%」は汎用 NPC の割合ではない。台詞の多い汎用ボイス話者は上位が全員固有名（Aela・Tolfdir・Farkas）で、固有キャラの汎用ボイス再利用が大半。voice_type は汎用 NPC の穴埋めでなく、固有キャラへの第 2 の独立信号（相互検証・精緻化と、台詞少の固有キャラの分類）と位置付けた。真の汎用 NPC は話者リンクの無いプール台詞（line の 38%）に散り個別 speaker に現れにくい。ユニーク/テンプレの正確な分離は template 継承が要るが `template_speaker_id` は空で、extractor 拡張の対象。
- 含意の尊大（Maven 型）の機械検出を検証し、表層特徴では取れないと結論した（`poc-tone-report.md` の「含意の尊大の機械検出を試した結論」節）。候補 3 種を採用条件（過学習しない・他 NPC を尊大にしない）で測り、再現性のある信号（含意脅し）は温厚な Tolfdir を誤爆し、安全な信号（支配言明）は全 912 話者で 8 人しか出ず Maven 専用＝過学習だった。3 条件を同時に満たす表層手法は無い。LLM での特別対応も採らない（効果が小さく、LLM へ回す条件「含意か」の判定自体が含意の検出を要して循環する）。Maven 型は中立へ分類される既知の許容誤差として受け入れる、と決めた。
- 人物像の信号マップを図に起こした（`persona-signal-map.md`）。対人マーカー（印）の有無で本文 2 軸と voice 気質 prior を切り替える縮退合成。感情軸はマーカー無しでも常時測れる。マーカーの少ない固有キャラ（Nazeem 横柄・Grelod 気難＋激情・Edorfin 高慢）で voice＋感情が人物像になることを実証した。全経路が機械処理で LLM は使わない。
- extractor 拡張（template・occupation・AI データ）は本 task の対象外と決めた。穴は固有ボイス×台詞少の端役 106 人（12%）に限られ翻訳量も小さい。個別 NPC と真の汎用 NPC の分離は、台詞が line_speaker で本人に紐づくか（台詞所有）で達せられ、template 継承の抽出は不要。完全網羅（未聴取 NPC もメタだけで分類）が後で要るなら別 task に切り出す。
- メタと本文の融合を PoC へ実装・実証した（`cmd/poc-tone` に voice 気質辞書 18 種と `fuseAttitude`、`sparse`・`verify` モード。結果は `poc-tone-report.md` の「メタと本文の融合」節）。印 10 以上は本文、未満は voice 気質 prior、固有 voice で prior 無しは保留。感情軸は全経路で本文から測る。voice prior 経路は台詞少の固有キャラを妥当に分類した（Nazeem 横柄・Edorfin 高慢・Grelod 気難＋激情）。本文経路は無回帰で、voice との相互検証は多弁キャラで食い違い、本文が精緻と確認した（voice は粗い prior）。残課題は印 7〜9 の閾値の崖と prior 値の較正。
- 実 mod `inigo.esp` で動作確認した（`cmd/poc-tone` に POC_DB 環境変数と `all` モードを追加。`db/poc-inigo.sqlite3` を抽出、結果は `poc-tone-report.md` の「実 mod での動作確認」節）。追加 NPC の Inigo は 3980 台詞・印 452 で物腰やわと分類でき、新規 NPC の判定は対話で成立した。既存 NPC は追加台詞だけだとセリフ不足で、voice が気質を持つ Nazeem は skyrim.esm 版と一致したが、voice 中立の Mjoll は中立へ落ちた。設計含意は、既存 NPC の persona を解決した話者の同一性で全プラグインの台詞を束ねて出すこと（既存は base game の算出済みを引き、新規だけ算出）。creature（Cr 声型）は保留で雑音を出し、対象外にすべきと判明。

### 成果物

設計・PoC（確定記録）:
- `tone-concept-model.md`: 口調概念モデル。
- `poc-tone-report.md`: PoC の確定記録（結果・確定判断・再現手順）。
- `persona-signal-map.md`: 人物像の信号マップ（縮退判定フロー）。対人マーカーの有無で本文 2 軸と voice 気質 prior を切り替える戦略。

設計成果物（design-module、人間レビュー待ち）:
- `persona-design.md`: 本実装の設計（Speaker に生える口調の箱、純粋 IO クラスの境界、行ごとの prose をキャッシュする永続化、few-shot とマーカーの合成順）。
- `implementation-scope.md`: 実装範囲とテスト設計（追加 3 テーブル、ToneClassifier、縦切り 4 スライス）。
- `test-plan.md`: ESP 抽出 fixture での網羅回帰（PoC の golden を実抽出データで再現することを保証）。

実装・検証用（正本でない）:
- `cmd/poc-tone/main.go`: PoC 実装（prose 依存を go.mod に追加）。引数 `verify`・`sparse`・`all` で再現する。
- `db/poc-skyrim.sqlite3`・`db/poc-inigo.sqlite3`: 検証用抽出 DB。
- `dictionaries/nrc-emolex.txt`: NRC EmoLex 辞書（研究用ライセンス）。

### 次段

- 設計成果物の人間レビューを受ける（design-module の出口）。承認後に実装へ入る。
- 縦切り 4 スライスを実装する（`implementation-scope.md`）。スライス 1（ToneClassifier・純粋 IO・テスト 100%）から始め、画面のあるスライス（属性割当 UI・キャラ表示編集 UI）は storybook-module、実装本体は implementation-module で縦通しに書く。
- 実装で詰める較正は `poc-tone-report.md` の「本実装で残る較正」にまとめた（閾値の崖・prior 値・頑健統計・語彙マーカー拡張）。creature（Cr 声型）の対象外化、話者同一性での台詞集約も実装で入れる。

### 注意

- NRC EmoLex は研究用途のみのライセンス。本番組み込み時は商用確認か代替（Empath は MIT）へ差し替える。
- PoC コード・抽出 DB・辞書は検証用。本実装の書き換えは implementation-module で行う。
- 対人スコアは「印」が十分（目安 10 以上）なキャラだけで評価する。サンプル少のキャラ（Grelod 印 1、Nazeem 印 2）は ±1.00 の極端値で参考外。

## 実装進捗（implementation-module、2026-06-23）

### スライス 1: ToneClassifier（純粋 IO・テスト 100%）完了

`implementation-scope.md` の縦切りスライス 1 を実装した。口調分類の不変ルールを単一の純粋 IO クラスへ閉じ、prose・DB・ハッシュへ依存させない（R10）。

変更ファイル:
- `internal/engine/tone/classifier.go`: Classifier（採点・印集計・帯分け・融合）と Features・Persona・段階定数。
- `internal/engine/tone/voice_traits.go`: voice 気質 prior 辞書（18 種）と voicePrior・voiceLabel・isUniqueVoice。
- `internal/engine/tone/classifier_test.go`: 単体テスト（各規則と Classify 統合）。

入力はキャッシュ済みの行特徴（line_analysis 由来の Features 列）と voice 気質ラベル。出力は基底口調（対人段階・感情段階・セル名）と品質指標（印・決定経路・voice 気質名）。融合は印 10 以上で本文、未満で voice 気質 prior、固有 voice で prior 無しは保留（PoC のロジックを移植）。

観測点（最終検証、go ツール直実行）:
- `gofmt -l internal/engine/tone/`: 差分なし。
- `go vet ./internal/engine/tone/`: 警告なし。
- `go test -cover ./internal/engine/tone/`: 通過、coverage 100.0%（全 11 関数 100%）。
- `go build ./...`: 全パッケージ成功。
- harness の `backend-local` suite は `scripts/harness/run.py` に未定義（定義は all・frontend-lint・frontend-local・frontend-test・structure）のため、backend 検証は go ツール直実行で行った。

close-4 の不変ルール（口調合成の純粋 IO クラス・カバレッジ 100%）は本スライスで満たした。残る close 定義（実画面・実データの振る舞い）はスライス 2〜4 で観測する。

### スライス 2: line_analysis ＋ persona_character ＋ 生成 ＋ 翻訳注入 完了

対話由来の基底口調を生成し、翻訳プロンプトへ注入する backend を縦通しに入れた。旧メタデータ専用ペルソナ機構（声質・種族・所属の気質をハードコード表で引く形）は作り直した。感情辞書は差し替え可能な境界（`EmotionLexicon`）にし、dev は NRC、製品化時に MIT 実装へ差し替える。

変更ファイル:
- `db/migrations/0005_persona_character.sql`: line_analysis（行解析キャッシュ）・persona_character（生成結果＋手修正保護）を追加。
- `internal/model/persona.go`: line_analysis・persona_character・生成入力・注入入力のモデル。
- `internal/store/persona_character.go`: 行解析キャッシュの一括取得・upsert、生成入力クエリ、手修正保護 upsert、注入入力クエリ。
- `internal/engine/linefeatures.go`: 感情辞書境界、prose と curated 辞書での特徴抽出、本文ハッシュ。
- `internal/engine/tone_catalog.go`: 基底口調セル → 性質文カタログ、種族訛りマーカー、注入ビルダー。
- `internal/engine/persona_generate.go`: 一括生成（キャッシュ集計 → ToneClassifier、手修正保護）。
- `internal/engine/engine.go`: Store interface・生成呼び出し・注入を新機構へ差し替え。
- `internal/lexicon/nrc.go`: 感情辞書の NRC 実装（差し替え可能境界の concrete）。
- `internal/bootstrap/bootstrap.go`: 感情辞書を読み engine へ注入。
- 旧機構を削除: `internal/engine/persona.go`・`persona_rule.go`（と各テスト）、`model.SpeakerIdentity`・`SpeakerPersona`、`store.LoadLineSpeakers`。

観測点（最終検証）:
- 単体テスト: engine・tone・store・lexicon すべて通過。tone は 100%。生成（fake store）・特徴抽出（fake 辞書）・性質文カタログ・注入を単体で固定。
- 実データ観測（使い捨てツールで poc-inigo へ生成、検証後に削除）: 23 話者を生成し、PoC の inigo 結果を再現した（Inigo=丁寧・本文・印452、Nazeem=尊大・voice 横柄、Mjoll・Jzargo=中立・voice で印不足）。
- キャッシュ観測: 1 回目（prose）3m41s、2 回目（line_analysis 命中・prose なし）31ms、line_analysis 4595 件・分類は同一。重い prose を本文ごと 1 度に畳む設計が実データで効いた。
- `gofmt -l`・`go vet ./...`・`go build ./...`: 差分・警告・失敗なし。harness の backend suite は未定義のため go ツール直実行。

残る観測（close-4 の実翻訳）: 実 app（UI）＋ LLM での実翻訳でプロンプトへ口調が入ることの目視は、app 起動時の UI 検証（chrome-devtools）で行う。注入ロジックは単体で、生成は実データで確認済み。

### スライス 4（表示）: 結果行に口調メタデータ — storybook-module 承認済み

方針変更: 専用のキャラ管理画面・編集はやめ、既存の翻訳結果行を開いた詳細に、台詞の話者の生成済み基底口調を「口調」メタデータとして出す（表示のみ）。判定結果（基底口調セル＋性質文）を強調し、根拠（決定経路・対人段階・感情段階・印）は小さくする。

変更ファイル（表示・storybook-module）:
- `TranslationResultRow.svelte`: 展開詳細に口調メタ節を追加。
- `translation-run-view.ts`・`translation-run-presentation.ts`: `PersonaMeta`（cell・trait・段階・印・経路）・段階型・ラベル・補足を追加。
- `TranslationResultRow.stories.ts`: 新表示の story（本文・声質・保留・叙述文・口調なし）。

承認・検証: 人間レビュー承認済み。通常分類（UI Components）へ統合、レビュー story 削除。`build-storybook` 成功、`frontend-local` 通過。詳細は `storybook-review-loop.md`。宙に浮いていた独立コンポーネント案（CharacterPersonaRow/Screen・編集 UI）は破棄した。

### スライス 4（配線）: 口調メタを結果行へ流す — implementation-module 完了・close-4 目視済み

承認済みの表示を実 app で動かす配線を縦通しに入れた。

変更ファイル:
- `internal/engine/engine.go`・`tone_catalog.go`: `Persona` に口調メタ（cell・trait・段階・印・経路）を足し、`LinePersonas` で埋める。
- `internal/api/app.go`: `PersonaView` DTO と `ResultView.Persona` を足し、`buildResultsPage` で話者ありの台詞へ載せる。
- `frontend/wailsjs/go/models.ts`: `wails generate module` で再生成（`api.PersonaView`・`persona?`）。
- `frontend/src/gateway/translation-gateway.ts`: `PersonaRow` と `toPersonaRow` を足し、結果行へ写す（段階・経路は境界で union へ確定）。

観測（close-4、実 app）: `npm run dev:wails:run`（:34115）で結果一覧を開き、台詞行に口調チップ、展開で口調メタ節（判定結果＋性質文を強調、根拠を小さく）が出ることを目視した。Grelod=ぞんざい・声質・印6 等、生成結果が表示へ正しく流れた。dev DB は使い捨てツールで persona_character を生成して観測した（ツールは削除済み）。検証: `go test ./...`・`go vet`・`gofmt` 通過、frontend `check`（当方 0 エラー）・`frontend-local` 通過。

### 既知の問題 1〜3 の対応（implementation-module、2026-06-23）

`persona-known-issues.md` の問題 1〜3 を実装した。問題 4 は許容誤差として受け入れた。

設計の出どころ（人間修正）: 性別は声型でなく NPC の Female flag、年齢は race EditorID（`ElderRace`／`*Child`）、声型は対人 prior フォールバック専用。役割語は決定経路（本文／voice／保留）に依らず常に付く。一人称・語尾は戯画的にせず（フィクションの老人語 わし・〜じゃ を使わない）現実的な register にする。

変更ファイル:
- `tools/extractor/LineSpeakerSqliteWriter.cs`: NPC の Female flag を読み speaker.sex を書く（問題 1）。
- `assets/role-speech.tsv`（新規）: 一人称・語尾テンプレート。`race`×`sex`×`cell` キー、ワイルドカード優先。見直しはファイル編集＋再 Run。
- `internal/engine/role_speech.go`＋`role_speech_test.go`（新規）: `RoleSpeechTable`・loader・純粋照合。カバレッジ 100%。
- `internal/engine/tone_catalog.go`: `buildToneTraits` で 性質文→役割語→種族訛り を合成（問題 1・2）。
- `internal/engine/engine.go`・`bootstrap.go`: テンプレートを読み Engine へ DI（`EmotionLexicon` と同じ差し替え可能境界）。
- `internal/model/persona.go`・`internal/store/persona_character.go`: 注入入力へ `Sex` を通す。
- `internal/engine/tone/classifier.go`＋`classifier_test.go`: 保留経路で対人を中立へ寄せる（問題 3）。

検証: `gofmt -l internal/`（差分なし）・`go vet ./...`・`go test ./...` 通過。tone・role_speech のカバレッジ 100%。`dotnet build tools/extractor` 成功。

実画面の最終目視まで完了した。dev DB を作り直し、`Innocence Lost - Quest Expansion.esp` を実 LLM（`hy-mt2-7b`）で再抽出→再生成→翻訳して目視した。性別は 7 話者すべてに充填され、Grelod（老女）は老年女性の口調になり山賊風が消えた。子供は一人称が性別どおりで安定し、Aventus（保留）は対人が中立へ寄り「平明」になった。Constance（成人女）は一人称「わたし」＋女性 register が付いた（基底口調は淡々のまま＝問題 4 許容）。詳細は `persona-known-issues.md`「対応状況」。

### 次

- スライス 3: persona_assignment（属性割当）＋ 属性割当 UI ＋ fallback 解決（印不足話者の口調を種族・声型グループから引く）。画面のため storybook-module 経由。
- 役割語の精緻化: 基底口調セル別の行を `assets/role-speech.tsv` へ後追いで足す（語尾の揺れが残るセルが見つかった場合）。few-shot 例文（R3）も同じ起動条件。
- 残較正（`poc-tone-report.md`「本実装で残る較正」）は別途。

## Outcome

口調の概念モデルと 2 軸を確定し、機械分類が LLM なしで成立することを実データで実証した（無作為検証・メタと本文の融合・実 mod inigo.esp 確認）。含意の尊大は既知の許容誤差と決め、extractor 拡張は対象外とした。設計成果物（`persona-design.md`・`implementation-scope.md`・`test-plan.md`）が揃い、人間レビュー後に storybook-module・implementation-module へ渡せる。計画資料は確定記録へ整理した（superseded な PoC 計画・v1 説明・第 2 段作業記録を削除し、再現手順を `poc-tone-report.md` へ一本化）。

## Finalization（2026-06-23）

### 正本化判断

- 反映対象: `docs/architecture.md`。
- 判断: 反映不要。
- 影響範囲: 本 session の変更は層・依存方向・Wails 境界の構造を変えない。`engine`・`store`・`api` 内へ機能追加した（`DeriveMasterTerms`、`RoleSpeechTable`、`LoadLineSpeakers`、`SpeakerView`）が、新層も依存方向の変更も無い。`SpeakerView` は既存 DTO へのフィールド追加であり、`assets/role-speech.tsv` と nrc 辞書は `bootstrap` 読み込みの既存パターンに乗る。
- 根拠: `feedback-architecture-reflection-structural-only`（層・依存・Wails 境界が不変なら §8 を churn せず承認も求めない）。
- 人間承認状態: 構造不変のため承認不要（恒久仕様の正本反映は行わない）。
- 対象 docs パス候補: なし。

### 最終検証（commit 前、role-speech 移動後の取り直し）

- `gofmt -l internal/`: 差分なし。
- `go build ./...`: 通過。
- `go vet ./...`: 通過（先行 session）。
- `go test ./...`: 全パッケージ ok。
- frontend `npm run check`（svelte-check）: 当方コード 0 error。既存の `node_modules/@storybook/svelte/dist/index.d.ts` 1 件のみ。
- 実 app（:34115）再起動: `assets/role-speech.tsv` から読めて clean 起動（`bootstrap.NewApp` の `LoadRoleSpeech` 失敗時はエラー伝播で起動失敗するため、起動成功が読み込み成功を裏づける）。

### 作業 commit

- branch: `claude/2026-06-20-character-persona-from-dialogue`。
- commit hash: `3bb42ea7`（32 files changed, 1323 insertions(+), 330 deletions(-)）。
- 内容: ペルソナ既知問題 1-3 の解消、固有名ファーストネーム派生の Run 組み込み、結果行の話者名＋属性表示、役割語テンプレートの `assets/role-speech.tsv` 外部化。
- commit 除外（本 task と無関係、または生成物）: `.DS_Store`、`.claude/skills/presentation/SKILL.md`、`cmd/CLAUDE.md`・`db/CLAUDE.md`・`docs/CLAUDE.md`・`scripts/CLAUDE.md`・`tools/CLAUDE.md`、`frontend/wailsjs/go/models.ts`（gitignore・再生成物）、`dictionaries/`（gitignore・nrc は研究ライセンス）。
- 残留リスク: スライス 3（属性割当 UI）は未着手。役割語の語尾揺れが残るセルは後追い拡張（`assets/role-speech.tsv` 編集＋再 Run）。
