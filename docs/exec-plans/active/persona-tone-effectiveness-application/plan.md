# Task Plan: persona-tone-effectiveness-application

## branch 情報

- 作業 branch: `codex/persona-tone-effectiveness-application`
- 統合先 branch: `master`
- 分岐元 commit: `93a4e3eebf7000aed23c1c19fd54f800d673a61e`
- worktree: `/Users/iorishibata/.codex/worktrees/persona-tone-effectiveness`

## やることの要点

### R-1 ペルソナ・性別・年齢・種族の組み合わせに応じた採用済みfew-shotを口調指示へ適用する

実験 `persona-tone-effectiveness` で人間が採用した、平明・ぞんざい・物腰やわの基底口調セル×性別×年齢×Khajiitまたは特別扱いなしのF1・F2・F3の英日対を、話者を解決できた台詞の口調指示へ入力順で適用する。3例の直後には、例を語句の写し替えに使わず、同じ自称・終助詞・命令形を反復せず、性別や年齢を示すためだけに「来い」「ぞ」「おらん」「おくれ」を選ばない例の使い方を適用する。

### R-2 汎用台詞では性別に応じてfew-shotを変え、PC発話にはfew-shotを追加しない

性別を取得でき、自由記述の口調がある汎用台詞では、成人・特別扱いなし・平明のF1・F2・F3の日本語訳文を男性と女性で変え、入力順で口調指示へ適用する。PC発話は性別・感情・言い回しを維持するが、F1・F2・F3を口調指示へ適用しない。

### R-3 汎用台詞の既定指示から衛兵の前提を外す

新しいDBの `tone_default` に入る汎用台詞の既定指示を「話者を特定できない汎用的な台詞。特定の職業や立場を仮定せず、原文に合う自然な口調で訳す。」へ変更する。既存DBでは、旧既定指示と完全一致する `tone_default` の保存値だけを新しい既定指示へ更新し、利用者が編集した保存値は保持する。保存値は既存のテンプレート編集画面へ表示される値で確認できる状態を維持する。

## やらないことの要点

- 平明・ぞんざい・物腰やわ以外の6種類の既存例文は変更しない。
- ペルソナ、性別、年齢、種族の判定方法は変更しない。
- Khajiit以外の種族へ専用few-shotを追加しない。
- 画面の文言・構造・style、表示用fixture、API契約、DB schema、翻訳モデルの選択は変更しない。
