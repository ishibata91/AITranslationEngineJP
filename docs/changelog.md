# 変更・判断履歴

正本（`requirements.md`、`system_requirements.md`、`architecture.md` など）には現在の状態だけを書く。
「なぜ変えたか」「何を落としたか」などの判断履歴は本ファイルに残し、正本へ混ぜない。
新しい entry を上に追加する。1 entry は date 見出しで区切る。

## 2026-06-13 業務要件2・3 のシステム要件を一部確定（対応策の発散と絞り込み）

### 変更

- `docs/system_requirements.md` 業務要件2「単語の一貫性」: `TBD` から、確定部分（Mod 横断マスター辞書、用語特定はレコード固有名詞の機械抽出）と未確定部分（訳語供給、漏れ語対応、適用方式、検証）を分けて記述。
- `docs/system_requirements.md` 業務要件3「NPC の口調」: 「属性と会話履歴から AI でペルソナ生成」から「属性からルールベースで生成」へ方針変更。表現（構造化属性＋翻訳ディレクティブ）、変換（テンプレート・機械的）、永続境界（ルール集を永続・翻訳前設定／個々ペルソナ非永続）、適用（プロンプト注入）を記述。

### 判断

- 一貫性のスコープは Mod 横断（永続マスター辞書）を選択。ジョブ内・Mod 内は不採用。
- 用語特定は、誤検出が無く既訳ヒット率が高いレコード固有名詞の機械抽出のみを採用。AI 抽出・頻度抽出・辞書マッチによる漏れ語対応（第2層）は保留。
- 頻度抽出は単言語では対訳を出せず、対訳コーパスがある場合の統計アラインメントでのみ対訳を出せると整理。用語特定（どの語を揃えるか）と訳語供給（訳語をどう得るか）は別軸と確認。
- 業務要件3 は機械的抽出（属性 → ルール）をまず採用し、AI 生成と会話履歴解析は保留。
- ペルソナ表現は構造化属性（内部表現）＋翻訳ディレクティブ（適用形）の 2 段とし、変換はテンプレートで機械的に行う。
- 機械的抽出ではペルソナを NPC レコードから都度導出できるため、個々ペルソナは永続化しない。永続資産は「属性 → 翻訳指示のルール集」とし、翻訳前にユーザーが設定可能とする。過去構想の `master-persona`（個々ペルソナを永続・編集）とは性質が異なる。
- VoiceType（声タイプ。Skyrim が音声収録のため NPC を声・性格でグループ化した分類）は属性の中で口調と相関が高く、ルールの主軸候補とする。

### 残課題

- 業務要件2: 訳語供給方式（既訳流用のみか AI 併用か）、本文・会話中の漏れ語対応、辞書の適用方式、一貫性検証が未確定。
- 業務要件3: 使用属性の選定、属性の分類、ルール合成の衝突優先順位を概念モデルで整理する。
- ペルソナ属性・Skyrim 構造の概念モデルの置き場が未定（`index.md` Read Order の `skyrim-structure-model.md` は実体なし）。

## 2026-06-13 screen-design 廃止、画面の正本を Storybook へ。tech-selection に Storybook / Tailwind / daisyUI

### 変更

- `docs/screen-design/` を削除（README.md、design-system-ethereal-archive.md、code.html、screens/ 配下すべて）。
- `docs/tech-selection.md`: フロントエンドに Tailwind CSS、daisyUI、Storybook を追加。CSS framework 不採用の行を削除。公式参照 3 件を追加。
- `docs/index.md`: Read Order・Directory Contract・Choose The Right Record から screen-design を削除。画面・表示の設計は Storybook（`frontend/`）と明記。
- `.claude/skills/`: design-module、storybook-module、finalization-module、coding-protocol、fix-decision、investigation-module、implementation-module、diagramming の 8 件を「画面の正本 = Storybook」へ作り替え。
- `docs/exec-plans/active/README.md`、`templates/work-plan.md`、`templates/task-folder/{README.md, plan.md, detail-spec-diff.md}`: `screen-design-diff` と「Storybook 後画面設計差分整合」前提を Storybook 正本へ付け替え。
- memory（`feedback_boundary_responsibility_separation`、`feedback_storybook_module_trigger`、`feedback_implementer_no_agent_split`、`MEMORY.md`）: 画面設計参照を Storybook へ更新。

### 判断

- 画面・表示の設計判断の置き場を Storybook の story と svelte コンポーネントへ移す（ユーザー選択「Storybook を正本にする」）。
- `画面設計差分`（`screen-design-diff.<screen-id>.md`）doc を廃止。design-module は画面表示 doc を作らず、画面表示の設計は storybook-module が story とコンポーネントで直接行う。
- storybook-module の「Storybook 後画面設計差分整合」成果物を廃止。実装範囲を越える画面変更が要る場合は design-module へ差し戻して `実装範囲` を見直す経路に統一。
- finalization-module の docs 正本反映対象を `docs/architecture.md` のみに限定。画面の正本（Storybook）は frontend source として作業 commit に含める。
- fix-decision / investigation-module の画面再現確認の selector 正本を「実装済み画面の `data-testid` またはセレクタ」へ変更。
- AI プロバイダ指定「OpenAI API(llama, lmstudio)」は OpenAI 互換 API でローカル実行（llama、LM Studio）と OpenAI 本体を含む意味と解釈。

### 残課題

- exec-plans templates には screen-design と無関係の旧表記（Codex 名称、`docs/detail-specs/`、`agent-browser`）が残る。本変更の対象外。整理は別途判断する。
- `index.md` Read Order の `skyrim-structure-model.md`、`core-beliefs.md` の `er.md` は実体が無い既存リンク（本変更前から）。

## 2026-06-13 要件文書を業務要件とシステム要件へ分割

### 変更

- `docs/spec.md` を `docs/requirements.md` へ rename（git mv、業務要件の内容は不変）。
- `docs/system_requirements.md`: 新規作成。業務要件 1〜4 に対応するシステム要件を記述。1 = AI 利用（Gemini / xAI / OpenAI 互換 API（OpenAI、llama、LM Studio）/ Claude）、2 = TBD、3 = NPC 属性と会話履歴からペルソナ生成、4 = 機能要件＝業務要件。
- `docs/index.md`: Read Order、Directory Contract、Choose The Right Record を `requirements.md` / `system_requirements.md` へ更新。
- `docs/core-beliefs.md`: 関連文書リンクと記録方針を「業務要件 = `requirements.md`、システム要件 = `system_requirements.md`」へ更新。
- `docs/architecture.md` / `docs/tech-selection.md`: 関連文書リンクを `requirements.md` / `system_requirements.md` へ更新。
- `docs/screen-design/README.md`: `spec.md` 参照を `requirements.md` へ更新。

### 判断

- 業務要件（何をしたいか）とシステム要件（どう達成するか）を別文書に分ける。
- 単語の一貫性のシステム要件は TBD として明示的に保留（ユーザー判断）。
- AI プロバイダ指定「OpenAI API(llama, lmstudio)」は、OpenAI 互換 API でローカル実行（llama、LM Studio）と OpenAI 本体を含む意味と解釈した。

### 残課題

- `.claude/skills/coding-protocol/SKILL.md` の 2 箇所が `docs/spec.md` を参照（system 要件の参照行、docs 正本一覧）。auto mode で skill 編集が拒否されたため未修正。ユーザー承認後に `requirements.md` / `system_requirements.md` へ張り替える。
- `requirements.md` は用語集を廃止済みのため、`screen-design/README.md` の「用語」参照は形式的に古い。画面設計の用語運用を決める時に整理する。

## 2026-06-13 spec.md を業務要件専用へ書き換え

### 変更

- `docs/spec.md`: 恒久要件・用語集・状態機械を全削除し、4 つの業務要件（翻訳したい / 単語の一貫性 / NPC の口調 / xTranslator 形式出力）へ全面書き換え。各要件に目的を併記し、成功条件は不記載。

### 判断

- `spec.md` は業務要件（何をしたいか）だけにする。システム要件（どう実現するか）は別文書で扱う。
- 入力取得手段（xEdit 抽出など）はシステム要件側へ回す。業務要件側は対象を「Skyrim Mod のテキスト」とだけ書く。
- xTranslator 形式出力は、ツール固有だがユーザーの明示要望のため業務要件として残す。
- 成功条件は記載しない（ユーザー判断）。目的は記載する。

### 残課題

- システム要件の置き場が未定。入力取得手段、AI 基盤、ジョブ運用は置き場を決めてから書き起こす。
- `index.md`（`spec.md` を「恒久要件と用語集」と記述）と `core-beliefs.md`（「永続要件は `spec.md` に記録する」と記述）の文言が、業務要件専用化により古くなった。システム要件の文書を決める時に併せて直す。
