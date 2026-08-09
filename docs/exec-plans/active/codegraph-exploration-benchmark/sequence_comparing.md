# 辞書置き換えシーケンスの比較

## 結論

探索手段を指定しない terra は、固定した5項目をすべて再構成した。terra は自分で `rg` と行範囲の読み取りを選び、Semble、CodeGraph、LSPを使わなかった。回数無制限のCodeGraph探索は7回の`codegraph_explore`だけで4項目を再構成したが、合計token数は指定なしterraの1.66倍になった。MOOSE + Semble は辞書構築と参照訳の分岐を部分的に再構成した。Ataraxis + Semble は固有名処理へ探索範囲が偏り、本文の辞書置き換えを再構成できなかった。

正解では、`Engine.Run` から本文処理へ進む途中で固有名を確定し、人名の部分形を派生してから本文用辞書を作る。本文に参照訳がなければ、実行時タグを退避し、辞書置き換え後にタグを復元してから翻訳器を呼ぶ。参照訳が完全一致する場合は、辞書置き換えと翻訳器呼び出しを通らない。

## 試行へ渡した指示文

2 回の試行では、使う skill 名だけを変えた。MOOSE には `$codegraph-moose-search`、Ataraxis には `$codegraph-ataraxis-explore` を指定した。

```text
読み取り専用のコード探索を行う。コードと設定を変更しない。

`<方式に対応する skill>`を使い、次の依頼へ回答する。

辞書置き換えの as-is 仕様を Mermaid のシーケンス図にする。参加者、呼び出し順、辞書データの受け渡し、本文へ置換する前後、処理を通らない分岐を根拠付きで示す。

対象はbackendの`/Users/iorishibata/Repositories/AITranslationEngineJP/internal`とする。Semble MCPはlocal Model2Vecとlocal cacheだけを使い、コードを外部へ送らない。利用者は、この実行でSemble MCPへ`internal`のコードを渡して検索することを明示的に許可している。過去の探索結果、`MEMORY.md`、`.codex/memories`を読まない。通常の file 再読、もう一方の方式、LSP、Graphify、Web 検索を使わない。この指定を AGENTS.md の LSP 利用規約より優先する。

最終回答は日本語にする。最初に有効なMermaidの`sequenceDiagram`を1つ示す。図の後に、図の根拠となったfile pathとsymbolまたは行番号を簡潔に示す。確認できなかった呼び出しや条件は推測で図へ追加せず、不足として明示する。
```

探索手段を指定しない terra には、次の指示文を渡した。

```text
読み取り専用で次の依頼へ回答する。コードと設定を変更しない。

辞書置き換えの as-is 仕様を Mermaid のシーケンス図にする。参加者、呼び出し順、辞書データの受け渡し、本文へ置換する前後、処理を通らない分岐を根拠付きで示す。

対象はbackendの`/Users/iorishibata/Repositories/AITranslationEngineJP/internal`とする。過去の探索結果、過去の回答、`MEMORY.md`、`.codex/memories`を読まない。探索手段は指定しない。

最終回答は日本語にする。最初に有効なMermaidの`sequenceDiagram`を1つ示す。図の後に、図の根拠となったfile pathとsymbolまたは行番号を簡潔に示す。確認できなかった呼び出しや条件は推測で図へ追加せず、不足として明示する。
```

回数無制限のCodeGraph探索には、同じ依頼へ次の条件を追加した。

```text
CodeGraphのproject pathは`/Users/iorishibata/Repositories/AITranslationEngineJP/internal`であり、各CodeGraph呼び出しへ必ず指定する。CodeGraphとその他のtoolの呼び出し回数に上限を設けない。CodeGraphが返した位置を確認するための通常のfile再読も許可する。

`rg`、`grep`、`perl`、`awk`、`git grep`、および同等の文字列横断検索は使わない。GraphifyとLSPは使わない。CodeGraphで位置を取得する前にfileを走査しない。過去の探索結果と回答を読まない。MOOSEとAtaraxisの固定手順は使わない。
```

## 正解

正解は、評価前に固定した5項目をすべて満たす同期実行のシーケンスである。図を読みやすく保つため、辞書の準備と本文の処理を分ける。

### 1. 固有名の確定から本文用辞書の構築まで

```mermaid
sequenceDiagram
    participant Caller as 呼び出し元
    participant Engine as Engine
    participant Store as Store
    participant Stoplist as stoplist
    participant Dict as dictionary

    Caller->>Engine: Run(ctx, conn, model, plugin, onProgress)
    Engine->>Engine: GeneratePersonas(ctx)
    Engine->>Engine: TranslateUntranslated(...)
    Engine->>Store: ListUntranslatedProperNouns(ctx, plugin)
    Store-->>Engine: 未訳の固有名
    Engine->>Store: ListMasterTerms(ctx)
    Store-->>Engine: authoritative
    Engine->>Engine: translateProperNouns(...)
    Engine->>Engine: deriveRunProperNouns(ctx, plugin)
    Engine->>Engine: LoadDictionary(ctx)
    Engine->>Store: ListMasterTerms(ctx)
    Store-->>Engine: master_term
    Engine->>Store: ListProperNouns(ctx)
    Store-->>Engine: proper_noun
    loop master_term、proper_noun の順
        Engine->>Stoplist: Blocks(source)
        alt 除外対象
            Note over Engine: 辞書へ追加しない
        else 使用対象
            Note over Engine: master_term を先に pairs へ追加する
        end
    end
    Engine->>Dict: NewDictionary(pairs)
    Dict-->>Engine: 本文用 Dictionary
```

`NewDictionary` は同じ原語の最初の項目を残す。`LoadDictionary` が `master_term`、`proper_noun` の順で項目を渡すため、同じ原語では `master_term` が優先される。

### 2. 本文の参照訳分岐と辞書置き換え

```mermaid
sequenceDiagram
    participant Engine as Engine
    participant Store as Store
    participant Tag as runtimetag
    participant Dict as Dictionary
    participant Prompt as prompt
    participant Translator as provider.Translator

    Engine->>Engine: referenceIndex(ctx)
    loop narration または line
        alt 原文と参照訳が完全一致する
            Engine->>Store: UpdateNarrationDest または UpdateLineDest
            Note over Engine,Translator: composeBodyPrompt、Dictionary.Apply、Translateを通らない
        else 参照訳がない
            Engine->>Engine: composeBodyPrompt(base, directive, rawSource, dict)
            Engine->>Engine: prepareSource(rawSource, dict)
            Engine->>Tag: Mask(rawSource)
            Tag-->>Engine: masked、tags
            Engine->>Dict: Apply(masked)
            Dict-->>Engine: replaced
            Engine->>Tag: Restore(replaced, tags)
            Tag-->>Engine: 置換済み本文と復元済みタグ
            Engine->>Prompt: ComposePrompt(base, directive, source)
            Prompt-->>Engine: provider.Prompt
            Engine->>Translator: Translate(ctx, conn, model, prompt)
            Translator-->>Engine: 訳文
            Engine->>Store: 訳文を書き戻す
        end
    end
```

## MOOSE + Semble の試行

```mermaid
sequenceDiagram
    participant Caller as 呼出元（未確認）
    participant Engine as engine.Engine
    participant Store as store
    participant Stoplist as stoplist
    participant Dict as dictionary.Dictionary
    participant Runner as engine.BatchRunner
    participant Prompt as composeBodyPrompt

    Caller->>Engine: LoadDictionary(ctx)
    Engine->>Store: ListMasterTerms(ctx)
    Store-->>Engine: master_term
    Engine->>Store: ListProperNouns(ctx)
    Store-->>Engine: proper_noun
    loop 各 master_term と proper_noun
        Engine->>Stoplist: Blocks(Source)
        alt stoplist が除外する
            Stoplist-->>Engine: true
            Note over Engine: 辞書データへ加えない
        else 除外しない
            Stoplist-->>Engine: false
            Note over Engine: Source と Dest を dictionary.Term として保持
        end
    end
    Engine->>Dict: NewDictionary(pairs)
    Dict-->>Engine: Dictionary

    Note over Runner,Prompt: BatchRunner.planLines では Dictionary と本文 Source を渡す
    alt 既訳の参照先がある
        Runner->>Store: UpdateLineDest(ctx, line.ID, dest, statusTranslated)
        Note over Runner,Prompt: composeBodyPrompt を通らない
    else 既訳の参照先がない
        Runner->>Prompt: composeBodyPrompt(base, directive, line.Source, Dictionary)
        Note over Prompt,Dict: 本文への Apply 呼出しは今回の探索出力では確認できない
    end

    Note over Dict: Apply(source)
    alt 辞書が空（re == nil）
        Dict-->>Dict: source をそのまま返す
    else 辞書に照合可能な語がある
        Dict-->>Dict: 語境界内で最長一致し、Dest へ置換する
    end
```

MOOSE は、登録語彙、stoplist、辞書構築、参照訳の分岐を示した。MOOSE は、実行入口、固有名の確定と派生、runtime tag の退避と復元、翻訳器呼び出しを示さなかった。評価は 2/5 である。

## Ataraxis + Semble の試行

```mermaid
sequenceDiagram
    participant E as Engine
    participant P as ProperNounDictStore
    participant S as Store
    participant DB as proper_noun
    participant D as 本文機械置換辞書
    participant T as 本文

    E->>P: ListProperNouns(ctx)
    P->>S: ListProperNouns(ctx)
    S->>DB: SELECT proper_noun ORDER BY id
    DB-->>S: []model.ProperNoun
    S-->>P: []model.ProperNoun
    P-->>E: []model.ProperNoun
    Note over E,D: 固有名の全行を master_term と合流し、本文機械置換辞書の材料にする
    E->>D: proper_noun の訳を合流する
    Note over T: 本文置換前の入力本文
    Note over D,T: 本文へ置換する関数、呼び出し順、置換後の本文は今回の探索結果に含まれない

    alt 固有名フェーズで master_term に既訳がある
        E->>S: UpdateProperNounDest(id, 既訳, statusTranslated)
        Note over E: AI 翻訳を通らない
    else 既訳がない
        E->>E: translateProperNounByAI(...)
        alt 構造化出力・応答読取・一時的サーバ失敗
            Note over E: proper_noun を未訳のまま残す
            Note over E: 本文機械置換辞書へ入るかは未確認
        else それ以外の失敗
            Note over E: error を返し、実行を停止する
        else 訳を取得できた
            E->>S: UpdateProperNounDest(id, 訳, statusProvisional)
        end
    end
```

Ataraxis は、固有名の取得と翻訳失敗の分岐を示した。Ataraxis は、正解として固定した5項目を図から追えなかった。評価は 0/5 である。

## 探索手段を指定しない terra の試行

terra には skill 名、検索 command、MCP、LSPを指定しなかった。通常のローカル調査手段は terra 自身に選ばせ、過去の探索結果と回答だけを禁止した。

```mermaid
sequenceDiagram
    autonumber
    actor C as 呼び出し元
    participant A as api.App
    participant E as engine.Engine
    participant B as engine.BatchRunner
    participant S as store.Store
    participant L as dictionary.Stoplist
    participant D as dictionary.Dictionary
    participant T as runtimetag
    participant P as 翻訳提供者

    alt 同期翻訳
        C->>A: 翻訳要求
        A->>E: TranslateUntranslated(...)
        E->>E: 固有名を確定し、部分形を派生
        E->>E: LoadDictionary(ctx)
    else 一括翻訳の本文送信
        C->>A: 一括翻訳要求
        A->>B: SubmitBatch(...)
        B->>B: planBodyRequests(...)
        B->>E: LoadDictionary(ctx)
    else 結果画面の実プロンプト再構成
        C->>A: 結果ページ要求
        A->>E: LoadDictionary(ctx)
    end

    E->>S: ListMasterTerms()
    S-->>E: []MasterTerm（source, dest）
    E->>S: ListProperNouns()
    S-->>E: []ProperNoun（source, dest）
    loop master_term と proper_noun の各行
        E->>L: Blocks(source)
        alt stoplist が対象語を除外する
            L-->>E: true
            Note right of E: 辞書データへ加えない
        else 除外しない
            L-->>E: false
            E->>E: dictionary.Term{Source, Dest} を追加
        end
    end
    E->>D: NewDictionary([]dictionary.Term)
    Note right of D: 空の source または dest を除外し、同一 source は先勝ち<br/>長い source を先に置き、語境界で照合
    D-->>E: 辞書

    alt 同期翻訳または一括翻訳の本文
        loop 未訳の叙述文・台詞
            alt 原文が参照訳と完全一致する
                E->>S: 既訳を確定訳として書き戻す
                Note right of E: 本文の置換、プロンプト作成、翻訳提供者呼び出しを通らない
            else 参照訳と完全一致しない
                E->>T: Mask(rawSource)
                T-->>E: maskedSource, tags
                Note right of T: 置換前: 本文内の実行時タグを一時文字列へ退避
                E->>D: Apply(maskedSource)
                D-->>E: replacedSource, usedTerms
                Note right of D: 置換: 原語を確定訳語へ置換
                E->>T: Restore(replacedSource, tags)
                T-->>E: restoredSource
                Note right of T: 置換後: 本文へ実行時タグを復元
                E->>P: Translate(指示文 + restoredSource)
            end
        end
    else 結果画面の再構成
        loop 表示対象の叙述文・台詞
            A->>T: Mask(source)
            T-->>A: maskedSource, tags
            A->>D: Apply(maskedSource)
            D-->>A: replacedSource, usedTerms
            A->>T: Restore(replacedSource, tags)
            T-->>A: 表示用の置換後本文
        end
    end
```

terra は、同期翻訳、一括翻訳、結果画面の3経路を1図へ含めた。固定した5項目はすべて追えるため、品質は 5/5 である。一方、一括翻訳は `provider.Translator.Translate` を直接呼ばず、本文要求を一括翻訳の提供者へ送信する。同期翻訳と一括翻訳を同じ `Translate` 呼び出しで示した箇所を、食い違い1件として数える。

terra が選んだ探索は、`rg --files` が1回、symbol候補を探す `rg` が2回、`nl` と `sed` による行範囲の読み取りが3 commandである。Codex CLI上のtool呼び出しは合計5回である。

| 方式 | 入力 token | 出力 token | 合計 token | 経過時間 | tool 呼び出し | 品質 | 食い違い |
|---|---:|---:|---:|---:|---:|---:|---:|
| 探索手段を指定しない terra | 230,964 | 3,088 | 234,052 | 69.85 秒 | 5 | 5/5 | 1 |
| 回数無制限のCodeGraph探索 | 385,029 | 3,042 | 388,071 | 70.99 秒 | 7 | 4/5 | 0 |
| MOOSE + Semble | 177,812 | 2,015 | 179,827 | 63.20 秒 | 4 | 2/5 | 0 |
| Ataraxis + Semble | 182,383 | 2,107 | 184,490 | 54.44 秒 | 4 | 0/5 | 0 |

## 回数無制限のCodeGraph探索

有効な試行では、terra は `codegraph_explore` を7回使った。文字列横断検索と通常のfile再読は使わなかった。

```mermaid
sequenceDiagram
    participant S as Store
    participant E as Engine
    participant SL as Stoplist
    participant D as Dictionary
    participant RT as RuntimeTag
    participant P as Prompt
    participant AI as 翻訳提供元
    participant BR as BatchRunner
    participant BAI as Batch翻訳提供元

    Note over E,S: 同期経路: TranslateUntranslated
    E->>S: 未訳の固有名・叙述文・台詞を取得
    E->>S: ListMasterTerms（固有名フェーズの既訳供給）
    S-->>E: master_term の source→dest
    E->>E: 固有名フェーズを完了し、proper_noun を確定
    E->>E: deriveRunProperNouns

    E->>S: ListMasterTerms と ListProperNouns
    S-->>E: 辞書候補（master_term を先、proper_noun を後）
    loop 各候補語
        E->>SL: Blocks(source)
        SL-->>E: 一般語なら除外
    end
    E->>D: NewDictionary(残した source→dest 対)
    alt source または dest が空、または同一 source が後続
        D->>D: 対象を捨てる／先に来た訳を保持
    else 有効な辞書語
        D->>D: 最長一致・語境界の正規表現を構築
    end

    E->>S: ListReferenceTranslations
    S-->>E: 既存訳（rec, field, source, dest）
    E->>E: referenceIndex を構築

    loop 各未訳叙述文・台詞
        alt 既存訳が (rec, field, source) と完全一致
            E->>S: 既存訳 dest を確定訳として書き戻す
            Note over E,D: composeBodyPrompt と Apply を通らない
        else 完全一致する既存訳がない
            E->>RT: Mask(rawSource)
            RT-->>E: maskedSource と退避したタグ列
            E->>D: Apply(maskedSource)
            alt 辞書が空
                D-->>E: maskedSource を変更せず返す
            else 辞書語が語境界で一致
                D-->>E: dest へ置換した本文と使用語
            else 一致する辞書語がない
                D-->>E: maskedSource を変更せず返す
            end
            E->>RT: Restore(置換後本文, タグ列)
            RT-->>E: 生タグを復元した本文
            E->>P: ComposePrompt(base, directive, 本文)
            P-->>E: 完成 Prompt（User は置換後本文）
            E->>AI: Translate(Prompt)
            AI-->>E: 訳文またはエラー
        end
    end

    Note over BR,BAI: batch 経路: 固有名 batch の反映後
    BR->>E: planBodyRequests(plugin)
    E->>E: LoadDictionary と referenceIndex
    loop 各未訳叙述文・台詞
        alt 既存訳が完全一致
            E->>S: dest を確定訳として書き戻す
            Note over E,D: 本文 batch に載せず、Apply を通らない
        else 完全一致する既存訳がない
            E->>P: composeBodyPrompt（内部で Mask → Apply → Restore）
            P-->>BR: 置換後本文を持つ Prompt
            BR->>BAI: SubmitBatch(Prompt 群)
        end
    end
```

回数無制限のCodeGraph探索は、固有名の確定、辞書構築、本文置換、参照訳の分岐を再構成した。APIまたは`Engine.Run`から`TranslateUntranslated`へ入る呼び出しを調べなかったため、品質は 4/5 である。図とsourceの食い違いは確認しなかった。

最初の試行はCodeGraphのproject pathをrepository rootとして扱い、CodeGraphを利用できなかった。その後に`perl`を文字列検索へ使ったため無効とした。無効な試行は612,209 token、105.65秒、13回である。生データは`codegraph-unlimited-invalid-path`の名前で残した。

## 比較

| 確認項目 | 正解 | 探索手段を指定しない terra | 回数無制限のCodeGraph探索 | MOOSE + Semble | Ataraxis + Semble |
|---|---:|---:|---:|---:|---:|
| 実行入口から本文処理まで | あり | あり | なし | なし | なし |
| 固有名の確定と派生後に辞書を作る | あり | あり | あり | なし | なし |
| 登録語彙、stoplist、優先順、辞書構築 | あり | あり | あり | あり | なし |
| tag退避、辞書置換、tag復元、翻訳器呼び出し | あり | あり | あり | なし | なし |
| 参照訳の完全一致で処理を通らない分岐 | あり | あり | あり | あり | なし |

## 根拠

- `internal/engine/engine.go:155` の `Run` と `:164` の `TranslateUntranslated`
- `internal/engine/engine.go:310` の `prepareSource` と `:319` の `composeBodyPrompt`
- `internal/engine/engine.go:462` の `LoadDictionary` と `:480` の `translationVocabulary`
- `internal/core/dictionary/dictionary.go:28` の `NewDictionary` と `:65` の `Dictionary.Apply`
- `internal/engine/batch.go:284` の `BatchRunner.planLines`
- `/private/tmp/codegraph-exploration-benchmark/sequence-story/results/moose.md`
- `/private/tmp/codegraph-exploration-benchmark/sequence-story/results/ataraxis.md`
- `/private/tmp/codegraph-exploration-benchmark/sequence-story/results/unconstrained.md`
- `/private/tmp/codegraph-exploration-benchmark/sequence-story/results/codegraph-unlimited.md`

現在の Codex セッションには `lsp_status` がない。正解図は、新しい LSP 調査ではなく、評価前に固定した確認項目と、評価で保存済みの行番号付き CodeGraph 出力を根拠にした。
