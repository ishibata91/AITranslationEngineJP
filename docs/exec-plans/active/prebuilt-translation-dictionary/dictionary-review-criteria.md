# 翻訳辞書レビュー基準

## 判定すること

英語の原語と日本語の訳語を一組ずつ確認し、その対応を翻訳辞書として固定すると今後の翻訳に役立つかを判定する。

判定は `include` または `exclude` のどちらかを必ず返す。

## 入力

```json
{
  "id": 11593,
  "source": "Accidentally dropped this",
  "dest": "うっかり落としてしまいました"
}
```

判定には `id`、`source`、`dest` だけを使う。DBの属性、レコード種別、抽出方法、出現回数は判定条件にしない。

## 判定の中心

次の質問へ答える。

> `source` が今後の翻訳文に現れたとき、`dest` へ固定することで、翻訳の正確性または表記の一貫性が上がるか。

明確に上がる場合だけ `include` とする。それ以外は `exclude` とする。

## `include` にするもの

| 観点 | 条件 | 例 |
| --- | --- | --- |
| 固有名詞 | 人物、地名、組織、種族、神格などの名前である | `Whiterun` → `ホワイトラン` |
| 固有の対象名 | 武器、防具、アイテム、呪文、能力などの名前である | `Aetherium Forge` → `エセリウムの鋳造器具` |
| ゲーム固有語 | 一般語とは異なるゲーム内の意味や定訳がある | `Dragonborn` → `ドラゴンボーン` |
| 音写 | カタカナ表記を固定する価値がある | `Abacean Longfin` → `アバキアン・ロングフィン` |
| 非直訳 | 単語ごとの通常翻訳では既存訳を再現しにくい | 固有の称号、技能名、名称 |
| 表記統一 | 複数の自然な訳が考えられ、採用訳を固定する価値がある | 固有の役職名やゲーム用語 |

## `exclude` にするもの

| 観点 | 条件 | 例 |
| --- | --- | --- |
| 文 | 原語全体が発話文、叙述文、通知文、命令文である | `Accidentally dropped this` → `うっかり落としてしまいました` |
| 文脈依存 | 指示語、代名詞、話者、時制、前後関係がないと訳を固定できない | `Take this` → `これを取れ` |
| 内部説明 | 条件、テスト、トリガー、スクリプトなどを説明する文字列である | `CONDITIONED NOT TO START` → `開始しない条件` |
| 一般語 | 通常の翻訳で十分であり、固有の訳を固定する必要がない | `King` → `王` |
| 一般的な句 | 語の組み合わせどおりに訳せて、辞書で固定する効果がない | `Wooden Door` → `木の扉` |
| 誤った対応 | 原語と訳語が同じ意味または対象を表していない | 原語と無関係な訳語 |
| 危険な固定 | 文法や文脈によって訳が変わるため、一つの訳へ固定すると誤訳を生む | 一般動詞、代名詞を含む句 |

## 判定時の注意

- 既存訳があるという理由だけで `include` にしない。
- 大文字、語数、引用符、句読点だけで判定しない。
- 固有の題名や名称に動詞が含まれる場合は `include` にできる。
- 一度しか使われない固有名詞でも、表記を固定する価値があれば `include` にできる。
- Skyrimに存在する文字列という理由だけで `include` にしない。
- 判断材料は原語と訳語だけなので、固有名詞または固有の名称だと読み取れない場合は `exclude` にする。

## reason code

| code | 意味 |
| --- | --- |
| `INCLUDE_PROPER_NAME` | 固有名詞である。 |
| `INCLUDE_NAMED_OBJECT` | 固有の対象名である。 |
| `INCLUDE_GAME_TERM` | ゲーム固有語または定訳が必要な語である。 |
| `INCLUDE_TRANSLITERATION` | 音写表記を固定する価値がある。 |
| `INCLUDE_NON_COMPOSITIONAL` | 通常の直訳では再現しにくい対応である。 |
| `INCLUDE_NORMALIZATION` | 複数の自然な訳から採用訳を固定する価値がある。 |
| `EXCLUDE_SENTENCE` | 文全体である。 |
| `EXCLUDE_CONTEXT_DEPENDENT` | 文脈によって訳が変わる。 |
| `EXCLUDE_INTERNAL_DESCRIPTION` | 内部処理の説明である。 |
| `EXCLUDE_GENERIC_WORD` | 固定不要の一般語である。 |
| `EXCLUDE_GENERIC_PHRASE` | 固定不要の一般的な句である。 |
| `EXCLUDE_MISALIGNED` | 原語と訳語の対応が誤っている。 |
| `EXCLUDE_UNSAFE_FIXED_TRANSLATION` | 一つの訳へ固定すると誤訳を生む。 |
| `EXCLUDE_NO_DICTIONARY_VALUE` | ほかの採用条件に該当せず、辞書として固定する効果がない。 |

## 一段目AIの出力

```json
{
  "id": 11593,
  "decision": "exclude",
  "reason_code": "EXCLUDE_SENTENCE",
  "reason": "原語全体がthisの指示対象に依存する発話文であり、一つの訳へ固定しても今後の翻訳に役立たない。"
}
```

出力はJSONだけとする。`decision`、`reason_code`、`reason` を必須にする。

## 検査AIの確認項目

検査AIには、同じ `source`、`dest`、一段目AIの出力、この基準を渡す。

| check code | 確認すること |
| --- | --- |
| `CHECK_DECISION` | `include` または `exclude` が基準に合っている。 |
| `CHECK_REASON_CODE` | reason code が原語と訳語から確認できる。 |
| `CHECK_REASON` | reason が原語と訳語の具体的な特徴を説明している。 |
| `CHECK_DICTIONARY_VALUE` | 「今後の翻訳に役立つか」を判断しており、既存訳の存在だけを根拠にしていない。 |
| `CHECK_NO_SURFACE_SHORTCUT` | 語数、記号、大文字などの表面だけで判断していない。 |

## 検査AIの出力

```json
{
  "id": 11593,
  "verdict": "pass",
  "failed_checks": [],
  "confirmed_decision": "exclude",
  "reason": "原語は指示語を含む発話文であり、EXCLUDE_SENTENCEの根拠を原語と訳語から確認できる。"
}
```

`verdict` は `pass` または `fail` とする。`pass` は `confirmed_decision` が一段目AIの `decision` と一致する場合だけ許可する。

## バッチ結果の通過条件

| 検査 | 通過条件 |
| --- | --- |
| 件数 | 入力件数、一段目の成功件数、検査AIの成功件数が一致する。 |
| ID | 入力、一段目、検査AIの `id` が一対一で一致する。 |
| JSON | 全件が定義したJSONとして読み取れる。 |
| 値域 | decision、reason code、verdict、check codeが許可値だけである。 |
| 一致 | `pass` のconfirmed decisionが一段目のdecisionと一致する。 |
| 反映 | 一段目と検査AIが一致した `pass` だけを辞書判定として採用する。 |

## 一段目AIへの指示

```text
あなたはSkyrim Mod向け翻訳辞書の候補をレビューする。

入力は次のJSONである。
{
  "id": 整数,
  "source": "英語の原語",
  "dest": "日本語の訳語"
}

確認することは一つだけである。
sourceが今後の翻訳文に現れたとき、destへ固定することで、翻訳の正確性または表記の一貫性が上がるかを判断する。

判定はincludeまたはexcludeのどちらかを必ず返す。
明確に役立つ場合だけincludeにする。
採用する根拠がsourceとdestから確認できない場合はexcludeにする。
入力にない使用箇所、カテゴリ、設定、意味を推測してはならない。

includeにする条件:

1. INCLUDE_PROPER_NAME
   人物、地名、組織、種族、神格などの固有名詞である。
   例: Whiterun → ホワイトラン

2. INCLUDE_NAMED_OBJECT
   武器、防具、アイテム、呪文、能力などの固有の対象名である。
   例: Aetherium Forge → エセリウムの鋳造器具

3. INCLUDE_GAME_TERM
   一般語とは異なるゲーム内の意味または定訳がある。
   例: Dragonborn → ドラゴンボーン

4. INCLUDE_TRANSLITERATION
   カタカナ音写の表記を固定する価値がある。
   例: Abacean Longfin → アバキアン・ロングフィン

5. INCLUDE_NON_COMPOSITIONAL
   単語ごとの通常翻訳ではdestを再現しにくい固有の称号、技能名、名称である。

6. INCLUDE_NORMALIZATION
   複数の自然な訳が考えられ、採用訳を固定することで表記を統一できるゲーム用語または固有の役職名である。

includeは、上の条件を一つ以上sourceとdestから確認できる場合だけ許可する。

excludeにする条件:

1. EXCLUDE_SENTENCE
   source全体が発話文、叙述文、通知文、命令文である。
   例: Accidentally dropped this → うっかり落としてしまいました

2. EXCLUDE_CONTEXT_DEPENDENT
   指示語、代名詞、話者、時制、前後関係によって訳が変わる。
   例: Take this → これを取れ

3. EXCLUDE_INTERNAL_DESCRIPTION
   条件、テスト、トリガー、スクリプトなどを説明する文字列である。
   例: CONDITIONED NOT TO START → 開始しない条件

4. EXCLUDE_GENERIC_WORD
   通常の翻訳で十分であり、固有の訳を固定する必要がない一般語である。
   例: King → 王

5. EXCLUDE_GENERIC_PHRASE
   語の組み合わせどおりに訳せる一般的な句であり、辞書で固定する効果がない。
   例: Wooden Door → 木の扉

6. EXCLUDE_MISALIGNED
   sourceとdestが同じ意味または対象を表していない。

7. EXCLUDE_UNSAFE_FIXED_TRANSLATION
   文法または文脈によって訳が変わるため、destへ固定すると誤訳を生む。

8. EXCLUDE_NO_DICTIONARY_VALUE
   ほかのinclude条件に該当せず、辞書として固定する効果を確認できない。

判定時の禁止事項:

- 既存訳があるという理由だけでincludeにしない。
- Skyrimに存在する文字列という理由だけでincludeにしない。
- 大文字、語数、引用符、句読点の一つだけで判定しない。
- 動詞があるという理由だけでexcludeにしない。固有の題名や名称の場合がある。
- 語が一度しか使われないと推測してexcludeにしない。
- sourceとdest以外の情報を補って判定しない。

出力は次のJSONだけとし、Markdownや説明文を追加しない。
{
  "id": 入力と同じ整数,
  "decision": "include" または "exclude",
  "reason_code": "上で定義したcodeを一つ",
  "reason": "sourceとdestの具体的な特徴に基づく日本語の理由"
}

decision=includeの場合はINCLUDE_で始まるreason_codeを使う。
decision=excludeの場合はEXCLUDE_で始まるreason_codeを使う。
idは入力から変更しない。
```
