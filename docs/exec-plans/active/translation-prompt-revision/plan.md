# Task Plan: translation-prompt-revision

`plan.md` は branch 情報と、この task でやること・やらないことの要点を持つ。
設計判断、判断履歴、検証結果、実装結果は持たない。設計は `design.md`、恒久的に残す判断は `docs/changelog.md` に書く。

## やること

- 翻訳 AI へ送る指示文の既定値を全体として見直す。対象は次の 4 つ。
    - base 翻訳指示文（`prompt_template.base_directive`）の文面。
    - 指示文（`directive`）7 種の粒度と文面。REC:FIELD への割り当て（`record_type_master`）の変更を含む。
    - 口調の与え方。役割語テンプレート（`assets/role-speech.tsv`）の欠落、口調を例文で固定する仕組みの有無、基底口調 9 セルの性質文、汎用台詞と PC 発話の自由記述口調を含む。
    - 口調の雛形が 2 箇所にある状態の解消。`prompt_template.persona_template` と `directive` の「口調」が同じ役割を持つ。

## branch 情報

- `execution_branch`: `claude/translation-prompt-revision`
- `target_branch`: `master`
- `source_commit`: `5eed5813714ff7d33850a1d13e1b06b21ead4263`

## やらないこと

- 送信単位の変更。台詞・叙述文を 1 件 1 リクエストで送る形は変えない。会話の前後や同一話者の複数台詞をまとめて送る案は扱わない。
- 翻訳 provider の実装と構造化出力スキーマ（`translationSchema`）の変更。
- 抽出（C# 抽出器）と取込段の変更。`extracted_field` から `narration`・`line`・`proper_noun` へ振り分ける規則は変えない。
- 機械置換辞書の供給規則と一般語の選別（stoplist）の変更。
- 対人段階・感情段階の判定しきい値（`internal/core/tone` の `Classifier`）と、段階を決める本文特徴量（`internal/core/linefeatures`）の変更。
- 公式日本語既訳の形態素解析による一人称・語尾の観測。`japanese-tone-persona` として設計したが保留し、`docs/exec-plans/rejected/` へ移した。
- 実行時タグの退避・復元の仕組み（`internal/core/runtimetag` の `Mask`・`Restore`・`CountMissing`）の変更。タグ保護指示の文面は見直しの対象に含める。
