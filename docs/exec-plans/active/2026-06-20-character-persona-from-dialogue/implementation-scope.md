# 実装範囲とテスト設計（2026-06-22）

本書は design-module の `実装範囲` と `テスト設計` である。承認済み設計（`persona-design.md`）を、既存スキーマの規約（`id PRIMARY KEY`、自然キー `UNIQUE(plugin, form_id)`）に乗る DB の形で実現できる計画へ落とす。実装は implementation-module で Claude 本体が 1 文脈で縦通しに書く。

## DB の形（追加スキーマ）

抽出データ本体（line・speaker・race・voice_type・line_speaker）はキャッシュで、再抽出で再構築する。ペルソナと、最も重い行解析は、extractor が触らない別テーブルに安定キーで持つ。

追加テーブル 1: line_analysis（行解析のキャッシュ。最重要。prose の結果を本文ごとに 1 度だけ持つ）

```sql
CREATE TABLE IF NOT EXISTS line_analysis (
    id INTEGER PRIMARY KEY,
    source_hash TEXT NOT NULL,      -- 本文テキストのハッシュ（同一英文は 1 行）
    sentence_count INTEGER NOT NULL,
    polite_count INTEGER NOT NULL,
    insult_count INTEGER NOT NULL,
    is_imperative INTEGER NOT NULL, -- prose 由来（重い）。誘導命令は除外済み
    exclaim_count INTEGER NOT NULL,
    elong_count INTEGER NOT NULL,
    emotion_count INTEGER NOT NULL, -- 強感情語の数。感情辞書は差し替え可能境界で数えるため辞書名を列名に含めない
    UNIQUE (source_hash)
);
```

追加テーブル 2: persona_character（生成結果＋手修正）

```sql
CREATE TABLE IF NOT EXISTS persona_character (
    id INTEGER PRIMARY KEY,
    speaker_plugin TEXT NOT NULL,   -- 話者の安定識別（解決した base NPC の plugin）
    speaker_form_id TEXT NOT NULL,  -- 話者の安定識別（base NPC の form_id）
    attitude_band INTEGER NOT NULL, -- 対人段階 0尊大/1中立/2丁寧
    emotion_band INTEGER NOT NULL,  -- 感情段階 0抑制/1中/2激情
    marked INTEGER NOT NULL,        -- 印（信頼度）。UI 表示用
    decision_path TEXT NOT NULL,    -- 本文/voice/保留。UI 表示用
    hand_edited INTEGER NOT NULL DEFAULT 0,
    UNIQUE (speaker_plugin, speaker_form_id)
);
```

追加テーブル 3: persona_assignment（属性割当、手動）

```sql
CREATE TABLE IF NOT EXISTS persona_assignment (
    id INTEGER PRIMARY KEY,
    attr_kind TEXT NOT NULL,        -- 'race' | 'voice_group'
    attr_key TEXT NOT NULL,         -- race EditorID または voice グループ名
    attitude_band INTEGER NOT NULL,
    emotion_band INTEGER NOT NULL,
    UNIQUE (attr_kind, attr_key)
);
```

- 重い処理のキャッシュ: prose の品詞解析は line_analysis に本文ハッシュで持ち、ユニークな本文ごとに 1 度だけ実行する。生成は line_analysis を引いて集計するだけで安価。共有・プール台詞も 1 回で済み、再抽出でも本文が同じなら再利用する。
- 再抽出耐性: 3 テーブルとも extractor が書かない。persona_character は base NPC の `(plugin, form_id)`（同一 NPC で安定）をキーにし、speaker 行が再構築されても残る。
- 全プラグイン横断の集約: 同一 base NPC の台詞は、extractor が話者を master 連鎖で base へ解決するため、全プラグインを 1 DB に抽出すれば speaker 行 1 つに束なる。集約は line_speaker を speaker_id で束ねるだけで足りる。
- 非正規化の正当化（2 軸を列にした理由）: 口調軸の基数は固定で常に 2、かつ非対称（対人は印・decision_path を持つが感情は持たない）。子テーブルにすると感情行で印・path が NULL になり、join も増える。基底口調の位置は (対人段階, 感情段階) の 2 次元座標で、1 エンティティの 2 列が自然。
- 性質文カタログ: 基底口調 → 性質文・few-shot の固定参照。完了定義はカタログ編集を要求しないため v1 はコード定数で持つ（編集 UI は将来 task）。
- 語彙マーカー: 種族 EditorID から純粋関数で導出するためテーブルを持たない。
- 既存の speaker.tone・nature 列は使わない（旧機構の残骸。本 task で旧機構を作り直す）。

## 純粋 IO クラス（ToneClassifier）

口調分類の不変ルールを単一の純粋 IO クラスへ閉じる（R10）。

- 置き場所: `internal/engine` 配下の新パッケージ（例 `internal/engine/tone`）。
- 入力: 話者の台詞群の言語特徴（line_analysis から引いたキャッシュ済み特徴の列）、voice 気質ラベル。
- 出力: attitude_band、emotion_band、attitude_score、emotion_score、marked、decision_path。
- 内側の不変ルール: 特徴量採点、印集計、帯分けしきい値、voice prior 融合。PoC `cmd/poc-tone` のロジックを移植する。
- クラス外: prose の品詞解析（命令文判定）を含む行ごとの特徴抽出は line_analysis のキャッシュ段に置き、クラスへはキャッシュ済み特徴を渡す。これでクラスは prose に依存せず 100% テストできる。DB 読み書き・ハッシュ計算もクラス外。
- 感情辞書の境界: 強感情語の判定は差し替え可能な境界（`engine.EmotionLexicon`）にする。dev は NRC 実装（研究用ライセンス）、製品化時に MIT ライセンスの実装へ差し替える。line_analysis の `emotion_count` はこの境界で数えるため、列名に辞書名を含めない。

## 実装スライス（縦切り・順序・各々観測可能）

- スライス 1: ToneClassifier（純粋 IO）と単体テスト 100%。観測点は単体テスト通過（close-4 の不変ルール検証）。
- スライス 2: line_analysis ＋ persona_character ＋ バッチ生成 ＋ 翻訳注入。migration 0005（line_analysis・persona_character）、行解析のキャッシュ（本文ハッシュで prose を 1 度だけ）、store の集約読みと upsert、engine のバッチ生成（キャッシュ済み特徴を集計 → ToneClassifier、hand_edited は保護）、基底口調 → 性質文カタログ → `{traits}` 注入（`ComposePrompt` 再利用）。観測点は実翻訳で多弁キャラ（Inigo 等）の基底口調がプロンプトへ入り、再抽出後も本文ハッシュ一致で prose を回さず保持されること。
- スライス 3: persona_assignment ＋ 属性割当 UI ＋ fallback 解決。migration 0006、store の CRUD、engine の解決順（persona_character があれば使い、無ければ persona_assignment の voice グループ・種族へ畳む）、api binding（属性一覧と割当の取得・保存）、画面（storybook-module）。観測点は UI で種族・声型グループへ基底口調を割り当て保存し、印不足の話者の翻訳へ反映されること。
- スライス 4: キャラペルソナの表示・編集 UI ＋ 手修正保護。api binding（キャラペルソナ取得・編集）、画面（storybook-module）、編集で hand_edited を立てる。観測点は生成ペルソナを UI で見て編集でき、次のバッチ生成で上書きされないこと。

## テスト設計

- 単体テスト（純粋・100% 基準）: ToneClassifier（採点・融合・帯分け）、基底口調 → 性質文カタログの写像、source_hash 計算（行の本文ハッシュ）、fallback 解決ロジック（store を mock で切る）、voice 気質 prior 辞書、few-shot とマーカーの合成順（マーカー優先で一人称・語尾粒子を上書き）。
- 網羅回帰（ESP 抽出 fixture）: ToneClassifier・融合・マーカーが PoC の結果を実抽出データで再現することを保証する backend テスト。詳細は `test-plan.md`。
- E2E（UI 起点）: 属性割当の編集・保存、キャラペルソナの表示・編集、翻訳実行での口調反映。
- 単体で書かない: store の SQL、Wails bridge、画面ロジック、バッチ生成の DB 込み経路（E2E に任せる）。

## 既存コードの置き換え

- `internal/engine/persona_rule.go`（`voiceNature`・`raceNatureByEDID`・`factionNatureByEDID`）と `persona.go` の性質文組み立てを、新機構（ToneClassifier ＋ persona_character ＋ persona_assignment ＋ 性質文カタログ）へ作り直す。
- voice 気質の対応（横柄・臆病・温厚老 等）は ToneClassifier の voice prior 辞書へ移す。
- `buildPersonaDirective` の `{traits}` 差し込みと `ComposePrompt`・`RenderPrompt` は再利用する。
- `engine.LinePersonas` の解決を、ハードコード表から persona_character・persona_assignment 参照へ差し替える。

## アーキテクチャ反映（finalization で §8 へ）

- 追加: ToneClassifier（engine 内の純粋 IO）、line_analysis（行解析キャッシュ）・persona_character・persona_assignment（store＋migration 0005・0006）、バッチ生成（engine）、ペルソナ取得・割当・編集の api binding、画面（storybook-module）。
- 不変: 層・port（provider が唯一）・Wails 境界の方向・Bootstrap の DI。
