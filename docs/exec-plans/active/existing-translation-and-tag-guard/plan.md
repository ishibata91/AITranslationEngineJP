# existing-translation-and-tag-guard

## 依頼要約

`docs/known-issues.md` の項目7・項目8 を実装する。あわせて実装済みの項目4（xTranslator 書き出し）を `known-issues.md`・`roadmap.md` から除く。

- 項目7 既存訳との完全一致置換: 原文が既存訳（公式日本語版など、xTranslator 英日 XML が持つレコード単位の訳）と完全一致するレコードを、AI 翻訳へ回さず既訳へ機械的に置き換える。
- 項目8 実行時タグ保護: 叙述文・台詞の原文に含まれる実行時タグ（`<Alias=...>` 等）を、AI 翻訳の前後で退避・復元し、欠落・改変を検出する。

## 分岐

- 分岐元 branch: `master`
- 分岐元 commit: `8e1b274f`
- 作業 branch: `claude/existing-translation-and-tag-guard`

## 供給源の確定（項目7）

参照訳の供給源は、既に取り込んでいる xTranslator 英日 XML（`xTranslatorXMLs`）の拡張とする。現状は固有名（FULL 行）だけを master 辞書用に拾うが、これをレコード単位の既訳まで拾い、抽出レコードと突き合わせる。新規の参照訳 import 画面は作らない。

## match key の判断（実装で確定）

known-issues 項目7 の文面は照合キーを「plugin・form_id・原文」とするが、実装では **(rec, field, source)** で照合する。理由は 2 点。

- **form_id を使わない**: 供給源の xTranslator XML の `<String>` は `EDID`・`REC`・`Source`・`Dest` だけを持ち、FormID を持たない（実データ `dictionaries/xTranslatorXMLs/Skyrim_english_japanese.xml` で確認）。form_id では照合できない。同一 `rec:field` かつ原文が完全一致するレコードへ既訳を流用する。
- **plugin で絞らない**: 公式日本語版などの既訳は base ゲーム由来で、翻訳対象 plugin（Mod）とは plugin が一致しない。plugin で絞ると Mod 側で 1 件も一致しない。同一原文の既訳を対象横断で再利用するため、plugin 非依存にする。

## 完了定義

2 つの独立した縦切りを持つ。それぞれ観測できる振る舞いと観測点を固定する。

### スライスA（項目8 実行時タグ保護）

- 動かす範囲: 原文中の実行時タグを AI 翻訳へ渡す前に退避し、訳文受領後に原形へ復元する。タグを持つ本文には保護指示をプロンプトへ注入して欠落自体を減らす。訳文でタグが欠落したら、その行は壊れた訳を確定させず未訳（status=0）のまま残し、再実行で再翻訳させる（人間決定 2026-07-19）。
- 観測点: タグの退避・復元・欠落検出・保護指示注入を純粋ルール／プロンプトの単体テストで確かめる。engine 経由の結線と欠落時の未訳保留を harness テスト・engine 単体で確かめる。タグを含む実データを実画面で翻訳し、保持ケースは原形保持、欠落ケースは未訳保留を目視する。

補足（タグ欠落時の対応の判断）: 落ちたタグの自動差し戻しは差し込み位置が不明で不可のため採らない。検出後は未訳のまま残し（安全網）、モデルにはタグ保護指示を与えて欠落を減らす。per-record の画面上の拾い上げ・修正は結果画面編集（項目5＝別課題）に委ねる。

### スライスB（項目7 既存訳との完全一致置換）

- 動かす範囲: 参照訳と `plugin`・`form_id`・`rec`・`field`・原文で完全一致する叙述文・台詞は、AI 翻訳を呼ばず既訳を訳文へ書き、status を確定訳相当にする。一致しないレコードは従来どおり AI 翻訳へ回る。
- 観測点: 一致判定と置換を単体テストで確かめる。既訳一致行が provider を呼ばないことを harness テストで確かめる。既訳のあるレコードが AI を経ず既訳で表示されることを実画面で目視する。

## close_conditions

- スライスA: 純粋ルールの単体テストが退避・復元・欠落検出を網羅する。harness テストで engine 経由のタグ保持を確認する。実画面でタグ原形保持を目視する。
- スライスB: 単体テストで完全一致判定と置換を確認する。harness テストで既訳一致行の provider 未呼び出しを確認する。実画面で既訳表示を目視する。
- 項目4 の doc 反映: `known-issues.md` から項目4 を除き、番号ずれ（roadmap の参照・項目3 の相互参照）を整合させ、`changelog.md` に経緯を残す。

## 軽 / 重判定

- 画面が動くか: N。新規 UI を作らず、既存の翻訳実行・結果表示に乗せる。
- `docs/architecture.md` 反映が要るか: N。層構成・依存方向・Wails 境界を変えない。変更は `internal/core`・`internal/engine`・`internal/store` の既存層内に収まる。参照訳のスキーマ影響は `er.md` 同期の対象で、`architecture.md` の対象ではない。

判定: 両方 N のため軽 task。`design-module`・`storybook-module` を bypass し、`implementation-module` へ進む。file 単位の触り方は実装時に Claude 本体が決める。

## 実装成果（implementation-module）

スライスA（実行時タグ保護）:

- `internal/core/runtimetag/runtimetag.go`: 新規 pure package。`<...>` タグを ⟦連番⟧ へ退避 `Mask`／復元 `Unmask`（欠落数を返す）。保護指示 `GuardInstruction`・退避有無 `HasPlaceholder`。
- `internal/core/runtimetag/runtimetag_test.go`: 単体テスト（カバレッジ 100%）。
- `internal/core/prompt/prompt.go`: `ComposePrompt` が退避トークンを持つ本文に system へタグ保護指示を足す（送信・再構成の両方で自動一致）。
- `internal/engine/engine.go`: `translateNarrations`・`translateLines` で `dict.Apply` の前に Mask、翻訳後に Unmask。欠落した行は未訳（status=0）のまま残す（保存せず再実行で再翻訳）。欠落は `slog`（result=skipped）で phase 集約 1 件。
- `internal/api/app.go`: 結果画面の実プロンプト再構成も送信と同順（Mask→`dict.Apply`）にして一致させる。
- `main.go`: `slog` の JSON handler を composition root で設定（observability-logging.md）。

スライスB（既存訳との完全一致置換）:

- `internal/core/termxml/reference.go`: 新規 pure。record 単位の既訳を XML から抽出（`ParseReferences`・`ReferencesFromFiles`）。
- `internal/core/termxml/reference_test.go`: 単体テスト。
- `db/migrations/0012_reference_translation.sql`: 新テーブル `reference_translation`（rec, field, source, dest, UNIQUE(rec,field,source)）。
- `internal/model/reference_translation.go`: 対応 model。
- `internal/store/reference_translation.go`: `InsertReferenceTranslations`・`ListReferenceTranslations`。
- `internal/engine/reference.go`: `ReferenceStore` interface、`LoadReferenceTranslations`（XML→表）、`referenceIndex`（表→照合 map）。
- `internal/engine/engine.go`: `Run` で `referenceIndex` を組み、翻訳ループの先頭で完全一致なら `statusTranslated`(1) で既訳を書き provider を呼ばずスキップ。
- `internal/api/app.go`: `RunExtractAndTranslate` で `DeriveMasterTerms` の後に `LoadReferenceTranslations` を呼ぶ。

テスト・検証基盤:

- `internal/engine/engine_test.go`: fakeStore に参照訳メソッド、既訳流用の Run テスト追加。
- `internal/harness/provider.go`: `RecordingProvider` が退避プレースホルダを保持（タグ往復の観測）。
- `internal/harness/fixture.go`: タグ入り叙述文（0x900）と既訳一致台詞（0x910）＋ TermsXML の参照訳を追加。
- `test-oracle/specs.json`・`internal/harness/oracle_test.go`: integration spec 2 件（`runtime-tag-preserved`・`existing-translation-reused`）と対応オラクル。
- `.go-arch-lint.yml`: `runtimetag` component 登録、engine・harness の依存追加。

## 最終検証

- `npm run verify:backend`（go test＋arch-lint＋境界走査）: 通過。
- `gofmt -l` 空、`go vet` 無警告。
- 純粋ルール単体（runtimetag 100%・termxml reference）、engine 単体（既訳流用スキップ）、harness 結合オラクル 8 件（新規 2 件含む）通過。
- C# 側: `SchemaMigrator` が `*.sql` を動的列挙し user_version で適用するため 0012 は自動適用（追加のみ・C# コード変更なし）。
- 実データ検証（スライスB）: 本物の base ゲーム XML（`dictionaries/xTranslatorXMLs`）を `LoadReferenceTranslations` で取込み、81,358 件の参照訳を materialize。既知 record（CELL:FULL "The Ratway Vaults"→"ラットウェイ・ウォーレンズ"）が照合表に載ることを確認（一時テストで検証後、削除）。
- 実画面 e2e（実 LLM）: `npm run dev:wails:run` で実 app を起動し、`Innocence Lost - Quest Expansion.esp` を LM Studio（`hy-mt2-7b`、`http://192.168.0.226:1234`）で削除→再翻訳し、結果画面で目視した。dev DB は version 11→12 へ migration。
    - スライスB: 参照 corpus 81,358 件を取込み、67 件が既訳を確定訳（status=1・画面「訳済」）で流用（AI 未呼び出し）。例「Kill Grelod the Kind → 親切者のグレロッドを殺す」。残 113 件は AI 経路（status=3・「仮訳」）。
    - スライスA: `<dur>`・`<BribeCost>` の実行時タグを退避→翻訳→復元。保持ケースはタグ原形保持（`<dur>` narration）。弱いモデル `hy-mt2-7b` が保護指示ありでも落とした `<BribeCost>` の 2 行は、壊れた訳を保存せず未訳（status=0・dest 空）のまま残し、集約ログ `{"event":"runtime_tag_lost","result":"skipped","count":2}` が発火。保護・注入・欠落時の未訳保留を確認。

判定: 完了定義の全観測点（単体・結合・実データ・実 app・実 LLM e2e）を満たした。タグ欠落時の対応（未訳保留＋保護指示注入）は人間決定（2026-07-19）を反映済み。

## 正本化判断（finalization）

- architecture.md 反映: **不要**。`runtimetag` は既存 core 構造（1 ルール 1 package・一方向 leaf）に収まる新 leaf で、層・依存方向・Bootstrap・Wails 境界・強い制約は不変。enforced な依存契約 `.go-arch-lint.yml` のみ更新済み（`feedback-architecture-reflection-structural-only`）。
- feature doc 同期（作業 commit に含める、正本反映 ceremony 対象外）: `known-issues.md`（項目4/7/8 除去＋番号整合）、`roadmap.md`（完了項目除去＋参照整合）、`er.md`（`reference_translation` 追加）、`changelog.md`（entry 追加）。
- 人間承認: architecture.md 反映が無いため承認 ceremony なし。項目7 の照合キー・項目8 の欠落時対応は本作業中に人間へ確認済み。
