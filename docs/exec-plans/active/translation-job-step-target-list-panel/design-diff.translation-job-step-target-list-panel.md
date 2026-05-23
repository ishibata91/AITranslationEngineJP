# 設計差分図: 翻訳ジョブステップ処理対象一覧表示パネル

## 概要

- 図化目的: `JobRunPage` の下に共通 `ProcessingTargetListPanel` を追加し、4 段階が同じ表示用情報を渡す予定変更箇所だけを固定する。
- 根拠参照: `detail-spec-diff.md`、`screen-design-diff.job-run.md`、`screen-design-diff.term-translation-phase.md`、`screen-design-diff.persona-generation-phase.md`、`screen-design-diff.body-translation-phase.md`、`screen-design-diff.translation-complete.md`
- 範囲: `JobRunPage`、4 段階の既存画面、追加予定の `ProcessingTargetListPanel`、段階ごとの表示用処理対象情報、ページング状態、現在ページ表示範囲に限定する。

## コンポーネント図

```mermaid
flowchart TB
    classDef added fill:#e8f5e9,stroke:#2e7d32,color:#1b5e20
    classDef removed fill:#ffebee,stroke:#c62828,color:#7f1d1d
    classDef unchanged fill:#fff8e1,stroke:#f9a825,color:#4e342e

    JobRunPage["JobRunPage<br/>選択ジョブ概要の下に共通パネルを追加する親画面"]:::unchanged
    ProcessingTargetListPanel["ProcessingTargetListPanel<br/>処理対象名と処理対象詳細を同じ構造で表示する共通パネル"]:::added
    PagingState["ページング状態<br/>既定ページサイズ: 約50件<br/>ページ切替で表示範囲を操作"]:::added
    CurrentPageRange["現在ページ表示範囲<br/>画面要素にする範囲"]:::added

    TermPhase["単語翻訳段階画面<br/>用語と固有名詞の情報を渡す"]:::unchanged
    PersonaPhase["NPC ペルソナ生成段階画面<br/>NPC の情報を渡す"]:::unchanged
    BodyPhase["本文翻訳段階画面<br/>翻訳対象項目の情報を渡す"]:::unchanged
    CompletePhase["翻訳完了確認画面<br/>翻訳結果の確認情報を渡す"]:::unchanged

    TermTarget["表示用処理対象情報<br/>段階名: 単語翻訳<br/>対象名: 共通辞書対象外の用語と固有名詞<br/>詳細: 確定訳語として保存するもの"]:::added
    PersonaTarget["表示用処理対象情報<br/>段階名: NPC ペルソナ生成<br/>対象名: NPC ごとのペルソナ生成入力<br/>詳細: ペルソナ参照情報を作るもの"]:::added
    BodyTarget["表示用処理対象情報<br/>段階名: 本文翻訳<br/>対象名: 辞書置換対象外の翻訳項目<br/>詳細: 辞書とペルソナを参照して訳文を作るもの"]:::added
    CompleteTarget["表示用処理対象情報<br/>段階名: 翻訳結果の確認<br/>対象名: 翻訳項目単位の訳文<br/>詳細: 出力前に確認するもの"]:::added

    JobRunPage --> ProcessingTargetListPanel
    JobRunPage --> TermPhase
    JobRunPage --> PersonaPhase
    JobRunPage --> BodyPhase
    JobRunPage --> CompletePhase

    TermPhase --> TermTarget
    PersonaPhase --> PersonaTarget
    BodyPhase --> BodyTarget
    CompletePhase --> CompleteTarget

    TermTarget --> ProcessingTargetListPanel
    PersonaTarget --> ProcessingTargetListPanel
    BodyTarget --> ProcessingTargetListPanel
    CompleteTarget --> ProcessingTargetListPanel
    PagingState --> CurrentPageRange
    CurrentPageRange --> ProcessingTargetListPanel
```

## 差分凡例

- 赤: 削除する要素または経路を示す。今回の設計差分では削除予定はない。
- 緑: 追加する要素または経路を示す。
- 黄色: 維持する要素または経路を示す。

## 各箱の説明

- `JobRunPage`: 共通パネルの配置先であり、4 段階画面の既存接続先を維持する。
- `ProcessingTargetListPanel`: 4 段階で共通の見た目を使い、処理対象名と処理対象詳細だけを表示する追加予定部品である。
- ページング状態: 約 50 件を既定ページサイズとして扱い、ページ切替で現在ページ表示範囲を操作する追加予定状態である。
- 現在ページ表示範囲: 数万件レベルの処理対象でも、画面要素にする範囲を現在ページへ限定する追加予定範囲である。
- 4 段階画面: 既存の段階画面として残り、現在段階に応じた表示用処理対象情報を共通パネルへ渡す。
- 表示用処理対象情報: 段階名、処理対象名、処理対象詳細を持つ表示データである。

## 追加予定

- `JobRunPage` の選択ジョブ概要の下に `ProcessingTargetListPanel` を追加する。
- 単語翻訳、NPC ペルソナ生成、本文翻訳、翻訳完了確認の 4 段階は、同じ共通パネルへ段階ごとの表示用処理対象情報を渡す。
- `ProcessingTargetListPanel` は、ページング状態から得た現在ページ表示範囲を表示する。
- 処理対象一覧にページング状態を追加する。
- 共通パネルは約 50 件を現在ページ表示範囲として扱う。
- 数万件レベルの処理対象では、現在ページ表示範囲を画面要素として扱う。

## 削除予定

- なし。

## 維持する接続先

- `JobRunPage` が 4 段階画面を切り替える既存構造は維持する。
- 各段階画面の本来の状態表示、操作、結果表示の責務は維持する。

## 検証

- Markdown 確認: 見出し、Mermaid コードブロック、差分凡例、各箱の説明、追加予定、削除予定、維持する接続先が揃っていることを確認した。
- Mermaid 記述確認: `flowchart TB`、箱、接続、`classDef` のみで構成し、差分凡例に必要な赤、緑、黄色のクラスを定義していることを確認した。

## 未決事項

- なし。
