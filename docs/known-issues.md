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

- 辞書に無い漏れ語（本文・会話文中にだけ現れ、名前付きレコードに出ない語）の拾い上げ方式は未確定である。AI 抽出・頻度抽出による第2層は保留する（`changelog.md` 参照）。言及テーブルは辞書に載る語の言及だけを持つ。
- 各テーブルの設計は [`er.md`](./er.md)、概念上の位置づけは [`concept-model.md`](./concept-model.md) にある。

## 2. 翻訳 runtime の未整備

- **Dialogue tree の context 長さ**: 現状は台詞を 1 件ずつ翻訳し（`engine` の台詞翻訳）、会話 tree の context を与えない。長い Dialogue の文脈利用（tree 分割・root path 抜粋・chunk 化）は未着手である。
- **固有名一貫性の事後検証**: 注入した確定訳語がモデル出力で保持されたかを照合する仕組みは無い。弱い小型モデルは注入したカタカナを書き換える場合があるため、検証方式は未確定である（`project-injected-token-fidelity` の観測を参照）。照合対象（どの行がどの語を含むか）は言及テーブル（`narration_mention`・`line_mention`）が持つため、残るのは照合方式の設計である。

## 3. クラウド AI プロバイダの未実装（Gemini・xAI・Claude）

`system_requirements.md` 業務要件1、`architecture.md` §3 は AI provider を Gemini・xAI・OpenAI 互換・Claude の 4 系統と定義するが、`internal/provider/` の実装は OpenAI 互換（`openai_compatible.go`）のみである。

- OpenAI 互換 API は OpenAI 本体とローカル実行（LM Studio・llama）を含むため、現状はこの経路だけで翻訳が動く。
- Gemini・xAI・Claude の各クライアント実装、および `provider.Translator` interface への適合が未着手である。

## 4. xTranslator 形式への書き出し未実装

業務要件4（xTranslator 形式で出力したい）に対応する書き出し処理が現状無い。

- `internal/core/termxml` は xTranslator XML の読み込み専用であり、人名部分形の派生（`DeriveMasterTerms`）の入力に使うだけである。
- 翻訳結果は `narration`・`line`・`proper_noun` テーブルへ dest・status 列として永続化済みだが、これを xTranslator 互換 XML へ変換して書き出す処理は未実装である。

## 5. 翻訳結果表示画面の編集・絞り込み機能未実装

`TranslationRunScreen` 配下の結果表示（`ResultsPanel`・`TranslationResultRow`）は一覧表示専用である。

- status・話者・種別によるフィルタ、原文・訳文の検索が無い。
- 訳文・status の編集、承認/却下操作が無い。
- 複数行の一括操作が無い。
- ページングは前へ／次への順送りのみで、ページ番号ジャンプが無い（`ResultsPager.svelte`）。
- `2026-06-14-results-paging-bulk-persona` plan は N+1 解消と keyset ページング追加が目的で、上記機能は scope 外として明示的に除外されている。
