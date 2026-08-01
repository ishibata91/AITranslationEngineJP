# Spec: persona-tone-effectiveness-application

---

## R-1 ペルソナ・性別・年齢・種族の組み合わせに応じた採用済みfew-shotを口調指示へ適用する

- R-1-1（正常系）: 平明・ぞんざい・物腰やわの話者を解決できた台詞では、ペルソナ・性別・年齢・種族の組み合わせに対応する採用済みfew-shotが3例とも口調指示へ入ること
    - 前提条件: 話者から性別と種族を取得でき、基底口調セルが平明・ぞんざい・物腰やわのいずれかである
    - 確かめ方: 送信する口調指示に、該当するF1・F2・F3の英日対が入力順で3件含まれることを見る
    - 対応する実テスト: `internal/core/personatone.TestBuildToneTraitsIncludesKhajiitExamplesAndUsageInstruction`、`internal/core/rolespeech.TestRealTableCoversAllNamedSpeakerKeys`
- R-1-2（対象に入る側の境界）: Khajiitの話者を解決できた台詞では、F1とF3が「この者」を使い、F2が「この者」を足さない採用済みfew-shotになること
    - 前提条件: 種族がKhajiitであり、基底口調セルが平明・ぞんざい・物腰やわのいずれかである
    - 確かめ方: 送信する口調指示のF1とF3に「この者」があり、F2に「この者」がないことを見る
    - 対応する実テスト: `internal/core/personatone.TestBuildToneTraitsIncludesKhajiitExamplesAndUsageInstruction`、`internal/core/rolespeech.TestRealTableCoversAllNamedSpeakerKeys`
- R-1-3（対象に入らない側の境界）: 平明・ぞんざい・物腰やわ以外の基底口調セルでは、既存のfew-shotが1例だけ口調指示へ入ること
    - 前提条件: 基底口調セルが冷然・見下し、淡々・実務、慇懃・端正、居丈高・罵倒、率直・興奮、情に厚い懇願のいずれかである
    - 確かめ方: 送信する口調指示に現行と同じ英日対が1件だけ含まれることを見る
    - 対応する実テスト: `internal/core/rolespeech.TestRealTableCoversAllNamedSpeakerKeys`
- R-1-4（例の使い方）: 採用済みfew-shotが3例入る口調指示では、例を語句の写し替えに使わず、同じ自称・終助詞・命令形を1台詞で反復せず、性別や年齢を示すためだけに「来い」「ぞ」「おらん」「おくれ」を選ばない指示が3例の直後へ入ること
    - 前提条件: 口調指示へ採用済みfew-shotが3例入る
    - 確かめ方: 送信する口調指示で、3件目の英日対の直後に例の使い方を示す文があることを見る
    - 対応する実テスト: `internal/core/personatone.TestBuildToneTraitsIncludesKhajiitExamplesAndUsageInstruction`

---

## R-2 汎用台詞では性別に応じてfew-shotを変え、PC発話にはfew-shotを追加しない

- R-2-1（正常系）: 性別を取得できる汎用台詞では、成人・特別扱いなし・平明の採用済みfew-shotが性別に応じて3例とも口調指示へ入ること
    - 前提条件: 汎用台詞の性別が男性または女性であり、自由記述の口調が空ではない
    - 確かめ方: 男性と女性へ送信する口調指示を比べ、F1・F2・F3の日本語訳文が性別に応じて異なり、入力順で並ぶことを見る
    - 対応する実テスト: `internal/core/personatone.TestGenericAndPCToneTraitsSeparateExamples`、`internal/engine.TestLinePersonasGenericAndPC`
- R-2-2（対象に入る側の境界）: PC発話では性別・感情・言い回しを維持し、F1・F2・F3が口調指示へ入らないこと
    - 前提条件: PC発話の性別が男性または女性であり、自由記述の口調が空ではない
    - 確かめ方: PC発話へ送信する口調指示に性別・感情・言い回しがあり、F1・F2・F3が含まれないことを見る
    - 対応する実テスト: `internal/core/personatone.TestGenericAndPCToneTraitsSeparateExamples`、`internal/engine.TestLinePersonasGenericAndPC`
- R-2-3（対象に入らない側の境界）: 性別を取得できない汎用台詞では、男性または女性に固定した採用済みfew-shotが口調指示へ入らないこと
    - 前提条件: 汎用台詞の性別が空であり、自由記述の口調が空ではない
    - 確かめ方: 送信する口調指示に男性または女性のF1・F2・F3が含まれないことを見る
    - 対応する実テスト: `internal/core/personatone.TestGenericAndPCToneTraitsSeparateExamples`、`internal/core/rolespeech.TestRealTableCoversFreeLinePath`

---

## R-3 汎用台詞の既定指示から衛兵の前提を外す

- R-3-1（正常系）: 新しいDBでは汎用台詞の既定指示が「話者を特定できない汎用的な台詞。特定の職業や立場を仮定せず、原文に合う自然な口調で訳す。」になること
    - 前提条件: `tone_default`に保存値がない
    - 確かめ方: テンプレート編集画面へ表示される汎用台詞の保存値を見る
    - 対応する実テスト: `internal/store.TestPromptTemplateToneDefaultsRoundTrip`
- R-3-2（対象に入る側の境界）: 旧既定指示と完全一致する保存値は新しい汎用台詞の既定指示へ更新されること
    - 前提条件: `tone_default`の保存値が旧既定指示と完全一致する
    - 確かめ方: 更新後にテンプレート編集画面へ表示される汎用台詞の保存値を見る
    - 対応する実テスト: `internal/store.TestMigration21UpdatesUneditedGenericToneDefault`
- R-3-3（対象に入らない側の境界）: 利用者が編集した汎用台詞の保存値は更新されないこと
    - 前提条件: `tone_default`の保存値が旧既定指示と一致しない
    - 確かめ方: 更新前後でテンプレート編集画面へ表示される汎用台詞の保存値が同じであることを見る
    - 対応する実テスト: `internal/store.TestMigration21PreservesEditedGenericToneDefault`
