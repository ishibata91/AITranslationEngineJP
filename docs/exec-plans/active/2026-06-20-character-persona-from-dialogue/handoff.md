# 引き継ぎ（2026-06-23、compact 前）

post-compact の継続用。現在地・未 commit・実行中状態・既知の問題・次の手・検証・制約をまとめる。正本は `plan.md`。

## 現在地

- task: 2026-06-20-character-persona-from-dialogue（口調ペルソナを話者メタの表引きから、対話由来の機械分類へ作り直す）。
- 完了: スライス 1（ToneClassifier・純粋 IO・テスト 100%）、スライス 2（生成・キャッシュ・翻訳注入の backend）、スライス 4（結果行に口調メタ表示＋配線、実 app で close-4 目視済み）。
- 未着手: スライス 3（属性割当 UI＝印不足話者の口調を種族・声型グループから引く fallback）。
- 最優先: 既知の問題 1〜4（ペルソナの中身の不足）。詳細 `persona-known-issues.md`。

## commit 状態

- commit 済み `f131ffc1`: スライス 1・2 の backend ＋ 設計 docs ＋ PoC（cmd/poc-tone）。
- 未 commit（作業ツリーに残る）:
  - スライス 4 表示: `TranslationResultRow.svelte`、`translation-run-view.ts`、`translation-run-presentation.ts`、`TranslationResultRow.stories.ts`（旧レビュー story は削除済み）。
  - スライス 4 配線: `internal/engine/engine.go`・`tone_catalog.go`（`Persona` に口調メタ）、`internal/api/app.go`（`PersonaView` DTO ＋ `ResultView.Persona`）、`frontend/wailsjs/go/models.ts`（再生成）、`frontend/src/gateway/translation-gateway.ts`（`PersonaRow` ＋ 写像）。
  - docs: `persona-known-issues.md`（新規）、`plan.md`・`implementation-scope.md`・`storybook-review-loop.md` の更新、本 `handoff.md`。
  - `.DS_Store`・`.claude/skills/presentation/SKILL.md` は本 task と無関係。commit に含めない。

## 実行中の状態

- Wails app が `http://localhost:34115/` で起動中（`npm run dev:wails:run`、background task `btck5m0ol`）。chrome-devtools 接続済み。
- dev DB `db/aitranslation.dev.sqlite3` に persona_character を生成済み（使い捨てツールで生成、ツールは削除済み）。121 台詞・7 話者（Grelod・AventusAretino・ConstanceMichel・子供 4）。
- app 再起動: `npm run dev:wails:run`（run-wails.sh が既存リスナを pkill して起動し直す）。
- 生成のやり直し（dev DB）: `engine.GeneratePersonas` を呼ぶ使い捨てを再作成して `go run` するか、翻訳 Run（Run 冒頭で生成）。

## 既知の問題（最優先、persona-known-issues.md）

1. 性別・年齢・世代が口調指示から落ちる。Grelod（老婆）→「失せろ」「お前」で山賊風、Aventus（子供）→「俺」、子供→大人口調、Constance（女性）→中性。
2. 一人称・語尾が同一話者内で揺れる。Aventus（1 話者・物腰やわ）が「俺の祈り」「感謝してるよ」「聞いたことがねえ」と揺れる。few-shot（R3）未実装。
3. 保留経路の誤分類。Aventus（興奮した子供）が「物腰やわ」になる（印 6 の薄い本文値を保持）。
4. 基底口調の取り違え。Constance（優しい世話役）が「淡々・実務」になる。
- 共通根: 口調指示が口調を言葉で「説明」するだけで、一人称・語尾・属性（性別/年齢）を「固定」しない。語彙マーカー（性別/年齢/役割語）と few-shot が未実装。
- 効いている点: 対人態度軸は概ね正しい（HOW は取れる、WHO と語形固定が抜け）。

## 次の手（推奨順）

1. 問題 1・2 を直す（複合・最優先）。
   - 語彙マーカー層に性別・年齢（世代）の役割語を足す。voice_type の EditorID（`FemaleOldGrumpy` 等）から 老女/老人/男性/女性/子供 を引く。足場: `internal/engine/tone_catalog.go` の `raceMarkerTrait` と同じレイヤー、`buildToneTraits` で基底口調へ重ねる。
   - few-shot（R3）で一人称・語尾を固定。合成順（マーカー優先で一人称・語尾を上書き）は `persona-design.md`「few-shot とマーカーの矛盾の検討」で設計済み。例文の具体が未実装。
   - 役割語の範囲・few-shot 例文は design 寄りの判断を含む。人間と相談しつつ進める可能性。
2. 保留経路の較正（問題 3）。中立へ寄せる／感情軸を主に使う／UI で低信頼を明示、のいずれか。
3. スライス 4 を commit（表示＋配線）。
4. スライス 3（persona_assignment 属性割当 UI、storybook-module 経由）。
5. 残較正（`poc-tone-report.md`「本実装で残る較正」: 閾値の崖・prior 値・頑健統計・語彙マーカー拡張）。

## 主要ファイル

- backend: `internal/engine/tone/`（ToneClassifier 純粋 IO・classifier.go・voice_traits.go）、`internal/engine/{linefeatures,tone_catalog,persona_generate,engine}.go`、`internal/store/persona_character.go`、`internal/lexicon/nrc.go`、`internal/api/app.go`、`internal/model/persona.go`、`db/migrations/0005_persona_character.sql`、`internal/bootstrap/bootstrap.go`。
- frontend: `frontend/src/ui/screens/translation-run/{TranslationResultRow.svelte,translation-run-view.ts,translation-run-presentation.ts,TranslationResultRow.stories.ts}`、`frontend/src/gateway/translation-gateway.ts`、`frontend/wailsjs/go/models.ts`（生成）。
- docs（plan フォルダ）: `plan.md`（正本）、`persona-design.md`、`implementation-scope.md`、`test-plan.md`、`poc-tone-report.md`、`persona-signal-map.md`、`tone-concept-model.md`、`persona-known-issues.md`。

## 検証コマンド

- backend（harness に backend suite 無し。go 直実行）: `gofmt -l internal/` / `go vet ./...` / `go test ./...`。tone は 100% カバレッジ基準。
- frontend: `npm --prefix frontend run check`（`node_modules/@storybook/svelte` の型宣言 1 件は既存・無視）／`python3 scripts/harness/run.py --suite frontend-local`（lint・test）／`npm --prefix frontend run build-storybook`。
- Wails バインディング再生成（Go DTO 変更後）: `/Users/iorishibata/go/bin/wails generate module`。
- 実 app（実画面確認）: `npm run dev:wails:run`（:34115、chrome-devtools）。

## 設計の要点（思い出し用）

- 口調 = 基底口調（対人態度軸×感情表出軸の 3×3 セル）＋ 語彙マーカー（種族訛り＝実装済、性別/年齢/役割語＝未実装）。
- 融合: 印 ≥10 で本文 2 軸 / 印 <10 で voice 気質 prior / 固有 voice で prior 無しは保留。感情軸は常に本文。
- ToneClassifier は純粋 IO（prose・DB・hash の外）。特徴抽出（prose ＋ 感情辞書）は line_analysis に本文ハッシュでキャッシュ（最も重い処理を本文ごと 1 度）。感情辞書は差し替え可能境界 `EmotionLexicon`（dev は NRC、製品化で MIT）。
- persona_character は `(speaker_plugin, speaker_form_id)` キー、`hand_edited` 保護。
- 含意の尊大（Maven 型）は既知の許容誤差（対応せず、LLM 分岐も設けない）。

## プロセス制約（CLAUDE.md / memory）

- 日本語・常体。実装・テスト・設計は Claude 本体が 1 文脈で書く（agent 分割は広範な探索のみ。subagent は haiku/sonnet）。
- module skill 中は「人間レビュー/承認」以外は確認せず自動進行。storybook-module は表示のみ＋人間レビュー、implementation-module は backend・frontend ロジック・配線。
- 動作確認は実画面 UI（binding 直叩きで代替しない）。LLM 実機は別マシン `http://192.168.0.226:1234`（dev 画面の :1234 表示は既定値）。
- task 構造・task-id は人間に確認せず Claude 側で決める。後回し禁止（scope か起動条件を明示）。
