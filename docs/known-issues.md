# 既知の課題

本書は docs 全体の「未解決の課題」「未確定の設計判断」「未実装の後続 task」を 1 か所へ集約した正本である。
各正本文書（`system_requirements.md`、`concept-model.md`、`skyrim-structure-model.md`、`er.md`）は現在状態だけを書き、開いている課題は本書へ集約する。
課題が解決したら本書から除き、解決の経緯は [`changelog.md`](./changelog.md) に残す。本書は「現在開いている課題」だけを保つ。

## 1. 固有名一貫性の後続 task（残りの言及関連）

本文中の言及（`narration_mention` e4・`line_mention` e5）と叙述文の説明対象（`narration_described` e3）は実装済みである（migration 0008。経緯は `changelog.md`）。固有名一貫性の関連で残る未実装は以下。

| 概念要素 | 概念モデル | 状態 |
|---|---|---|
| `line_sequence`（e7） | 台詞 → 台詞（会話の流れ） | 未実装 |
| `speaker_name`（e8）, `faction_name`（e14） | 話者・勢力 → 固有名（名乗る名） | 未実装。話者名は `master_term`・人名派生で代替 |
| `race.name_proper_noun_id`（e13） | 種族 → 固有名（名称） | 未実装の FK。現テーブルに列を持たない |

- 辞書に無い漏れ語（本文・会話文中にだけ現れ、名前付きレコードに出ない語）の拾い上げ方式は未確定である。LLM 不使用の決定的な候補検出は研究したが、精度（真の精度 54〜63%）と候補量の点で不採用にした（実装と実測は git 履歴の branch claude/dictionary-missing-term-detection、判断は `changelog.md` 2026-07-05）。代替候補は翻訳結果側の事後の訳揺れ検出である。言及テーブルは辞書に載る語の言及だけを持つ。
- 各テーブルの設計は [`er.md`](./er.md)、概念上の位置づけは [`concept-model.md`](./concept-model.md) にある。

## 2. 翻訳 runtime の未整備

- **Dialogue tree の context 長さ**: 現状は台詞を 1 件ずつ翻訳し（`engine` の台詞翻訳）、会話 tree の context を与えない。長い Dialogue の文脈利用（tree 分割・root path 抜粋・chunk 化）は未着手である。
- **固有名一貫性の事後検証**: 注入した確定訳語がモデル出力で保持されたかを照合する仕組みは無い。弱い小型モデルは注入したカタカナを書き換える場合があるため、検証方式は未確定である（`project-injected-token-fidelity` の観測を参照）。照合対象（どの行がどの語を含むか）は言及テーブル（`narration_mention`・`line_mention`）が持つため、残るのは照合方式の設計である。

## 3. クラウド AI プロバイダの未実装（Gemini・xAI・Claude）

`system_requirements.md` 業務要件1、`architecture.md` §3 は AI provider を Gemini・xAI・OpenAI 互換・Claude の 4 系統と定義するが、`internal/provider/` の実装は OpenAI 互換（`openai_compatible.go`）のみである。

- OpenAI 互換 API は OpenAI 本体とローカル実行（LM Studio・llama）を含むため、現状はこの経路だけで翻訳が動く。
- Gemini・xAI・Claude の各クライアント実装、および `provider.Translator` interface への適合が未着手である。
- Gemini・xAI の batch mode（`gemini-xai-batch-translation` task で設計中）は、個別行の翻訳がコンテンツフィルタ（暴力・差別表現等の判定）や不正なリクエストで、再送しても直らない形で恒久的に失敗する場合がある。失敗理由を保持する設計は決めたが、失敗した行を人間が見て手で直す編集手段は課題4（結果画面の編集機能未実装）に依存し、現状は存在しない。

## 4. 翻訳結果表示画面の編集・絞り込み機能未実装

`TranslationRunScreen` 配下の結果表示（`ResultsPanel`・`TranslationResultRow`）は一覧表示専用である。

- status・話者・種別によるフィルタ、原文・訳文の検索が無い。
- 訳文・status の編集、承認/却下操作が無い。
- 複数行の一括操作が無い。
- ページングは前へ／次への順送りのみで、ページ番号ジャンプが無い（`ResultsPager.svelte`）。
- `2026-06-14-results-paging-bulk-persona` plan は N+1 解消と keyset ページング追加が目的で、上記機能は scope 外として明示的に除外されている。

## 5. 機械置換辞書の誤爆対策の残り

一般語 stoplist（stopwords-en による機械置換辞書・言及語彙の供給源選別）は実装済みである（経緯は `changelog.md` 2026-07-04）。残る誤爆リスクと未確定の判断は以下。

- **stoplist 外の一般語 1 語の固有名**: stopwords-en に無い「名前全体が一般語 1 語」の固有名（例: Chest・Door・Summonable・Close・Health）は従来どおり辞書へ載る。本文の文頭・文中の大文字出現に当たり得る。stopword リストの独自拡張はしない方針（外部配布リストで賄う）のため、対策するなら注入方式の変更（本文置換をやめ、用語対訳ヒントをプロンプトへ添付する方式）の検討になる。プロンプト設計の再検討を伴うため未着手である。
- **管理用勢力・階級称号の機械判定基準**: 画面に出ない管理用の勢力名・階級称号（対話状態の内部フラグ用。実例 inigo.esp の FACT:MNAM "Yes"・"No"）を翻訳対象から機械的に除く基準は未確定である。候補だった FACT の Hidden from PC flag は実データ観測で不成立と確定した（Yes/No の供給源勢力に flag が無く、Skyrim.esm では既訳ありの実在勢力名 154 件が誤って落ちる。観測は `docs/exec-plans/completed/dictionary-false-positive-guard/plan.md`）。現状は stoplist が同綴りの一般語だけを止めている。
