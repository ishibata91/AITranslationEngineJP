# line-emotion-trdt

## 依頼要約

INFO 応答が持つ感情メタ（TRDT の感情型。台詞を喋る時の顔の表情演出として制作者が付ける）を、「その台詞を喋る時の感情」として翻訳プロンプトへ加算する。話者ペルソナ（`persona_character`、`tone.Classifier` の9セル、話者集計）は据え置き、変更は加算のみとする。

- 分岐元 branch: `master`
- 分岐元 commit: `68b20e21`

## 背景（この方針に至った決定の要点）

現状の感情表現が2つの意味で狭い、という観測から出発する。

- **種別を持たない**: 感情は強度3段（抑制/中/激情）だけで、怒り・悲しみ等の種別は `NRC` 辞書の段階で真偽値へ潰して捨てている。
- **話者単位でしか持たない**: `persona_generate.Generate` が話者ごとに `Classify` を1回呼び、話者に1段階を保存する。名指し話者は全台詞が同じ感情段階を共有し、台詞ごとの差が消える。

TRDT はこの2つを同時に埋める。応答（台詞1単位）ごとに感情型8種を持ち、制作者が明示した強い信号のため話者集計での補強を要しない。実データ計測では型はよく埋まる（バニラ非 Neutral 約55%、follower mod は90%超）。値（0-100）は既定50が大半で強度として弱いため、型のみを使う。

感情は「その台詞の状態」であり、ペルソナ（話者の恒常的特性）と直交する。よってペルソナへ畳まず、台詞ごとにプロンプトへ加算する。

## goal

実 plugin の INFO 応答の TRDT 感情型が抽出・取込され、翻訳実行時に台詞ごとの感情1行が翻訳プロンプトへ載り、非 Neutral の台詞でプロンプト内容が台詞ごとに変わる。

## 含まない（除外範囲）

- ペルソナ（`persona_character`、`tone.Classifier`、9セル、話者集計）の改変。据え置く。
- 感情の UI 表示（翻訳結果一覧への感情チップ等）。プロンプトへ載せるまでを対象にする。
- TRDT 値（強度 0-100）の利用。型のみを使う。
- Neutral 既定台詞の推測補完。上乗せなしで基調（ペルソナ）へ委ねる。
- 音声からの感情分類（SER）。将来の任意 prior として別 task に置く。

## 完了定義

### 動かす範囲

実 plugin の INFO 応答の TRDT 感情型が SQLite 契約経由で取り込まれ、翻訳実行時に「この台詞の感情」1行が翻訳プロンプトへ載る。非 Neutral の台詞でプロンプト内容が台詞ごとに変わり、Neutral 既定の台詞では加算しない。ペルソナ側（`persona_character` 由来の口調指示）の出力は本 task 前後で不変とする。

### 観測点

- 単体テスト（Go core）: `personatone` が非 Neutral で感情行を加算し Neutral で加算しないこと、感情型→日本語感情語の写像。
- C# テスト: 抽出器が TRDT 感情型を SQLite へ書くこと（抽出の正しさは C# テストで担保する。§6）。
- 実データ: 実 plugin（`Skyrim.esm` 等）の INFO 応答から感情型が中心 DB へ取り込まれること。
- 実画面: 翻訳実行で、同一話者の異なる台詞に別々の感情がプロンプトへ載ること（実 app、`http://localhost:34115`）。

## close_conditions

- [x] 実 plugin 抽出で INFO 応答の TRDT 感情型が中心 DB に入る（実データ観測）。
- [x] 非 Neutral 台詞でプロンプトに感情1行が載り、Neutral 台詞では載らない（単体テスト＋実画面）。
- [x] ペルソナ（`persona_character`）由来の口調指示が本 task 前後で不変（単体テスト）。

## 軽 / 重判定

- 画面が動くか: **N**。表示構造・文言・style・svelte 表示コンポーネント・story を変えない。翻訳結果の中身は変わるが、表示コンポーネントは変えない。
- `docs/architecture.md` 反映が要るか: **Y**。C#↔Go の SQLite 契約（§6）に TRDT 感情の抽出経路が加わり、抽出器の書き込み責務（§8）が広がる。
- 判定: **重 task**。`preparation-module` → `design-module` → `implementation-module` → `finalization-module`。画面が動かないため `storybook-module` は bypass する。

## 後続への引き継ぎ

`design-module` の入口へ完了定義と除外範囲を渡す。`design-module` で次を詰める。感情の保存単位（`line` 隣接）、C#↔Go 抽出契約の拡張形、感情型8種→日本語感情語の対応、`personatone` への加算行の形、テスト設計。

## 実装・検証結果（implementation-module）

- **Go test**: 全通過（`harness` の実データ合成 E2E 含む）。`personatone` に TRDT 写像・加算・Neutral 無加算・二重回避のテストを追加。
- **Go lint**: format / vet / arch / boundary / module 通過。static の6件は変更外の既存負債（`api`・`termxml`・`engine` の export.go）で本 task の対象外。
- **C# build 成功、C# test 20 passed**。
- **実データ観測**: `Innocence Lost - Quest Expansion.esp` 抽出で `extracted_info_emotion`=110・`line.emotion_type` 非空=110（INFO:NAM1 121 のうち Neutral 11 を除外、ordinal 全一致）。感情分布は独立パース（TRDT 直読）と一致。
- **実画面観測**: 翻訳プロンプトに感情行「- 感情: この台詞は〈X〉を込めた口調で話す。」が助言調で載る。同一話者 AventusAretino で台詞ごとに喜び／恐れが切り替わり、ペルソナ（口調・人称・対人/感情段階）は不変。
- **既知の環境問題**: `/tmp` GOMODCACHE 破損で run.sh 検証が失敗。`GOMODCACHE=~/go/pkg/mod` で回避（実装と無関係）。

## 保存形（Q3 承認通り line 列）

Q3 の承認通り `line.emotion_type` 列で持つ。当初 migration 慣習（0007・0008 が ALTER を避け専用テーブル化）を尊重して `line_emotion` 専用テーブルにしたが、dev は DB をリセットして作り直す運用のため ALTER の再実行問題は起きない（SchemaMigrator も db.Apply も user_version で 0011 を 1 度だけ適用する）。承認通り line 列へ戻し、dev DB をリセットして ALTER 適用・`UPDATE line` 取込・プロンプト感情行を実画面で再検証した（感情 110 件・分布一致・プロンプトに感情行、ペルソナ不変）。`extracted_info_emotion`（抽出 staging）は line 作成前に C# が書くため専用テーブルのまま残す。

## 正本化判断（finalization-module）

- **docs/architecture.md 反映**: 不要。層構成・依存方向・Wails 境界は不変（既存の抽出→取込→core→プロンプト経路にデータが1つ増えるのみ）。§6 は契約方針（SQL schema 1 本）で具体テーブルを列挙せず、新テーブル追加で方針は変わらない。§8 は状態スナップショットで、構造不変のため churn しない（feedback-architecture-reflection-structural-only）。人間承認は不要。
- **論理 ER（docs/er.md）**: 同期する（feature commit に含める）。`extracted_info_emotion` テーブルと `line.emotion_type` 列を追記した。
