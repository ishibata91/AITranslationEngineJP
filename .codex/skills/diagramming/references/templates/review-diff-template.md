# <task-name> review diff

## 概要

- 図化目的: <人間レビューまたは実装着手判断で確認する内容>
- 根拠参照: <task 内成果物または docs 正本>
- 範囲: <予定変更箇所と接続先の最小範囲>

## コンポーネント図

```mermaid
flowchart TB
    Current["既存境界"]
    Added["追加予定"]
    Removed["削除予定"]
    Unchanged["変更なし"]
    Note["確認観点"]

    Current --> Unchanged
    Added --> Unchanged
    Removed -.削除.-> Unchanged
    Note -.注意.-> Unchanged

    class Added added
    class Removed removed
    class Current,Unchanged unchanged
    class Note note

    classDef added fill:#e8f5e9,stroke:#2e7d32,color:#1b5e20
    classDef removed fill:#ffebee,stroke:#c62828,color:#7f1d1d
    classDef unchanged fill:#fff8e1,stroke:#f9a825,color:#4e342e
    classDef note fill:#f5f5f5,stroke:#757575,color:#212121
```

## 差分凡例

- 赤: 削除する要素または経路を示す。
- 緑: 追加する要素または経路を示す。
- 黄色: 変更しない要素または経路を示す。

## 各箱の説明

- 既存境界: 変更しない接続先を示す。
- 追加予定: 新しく追加する要素または経路を示す。
- 削除予定: 削除または廃止する要素または経路を示す。
- 変更なし: 接続または責務を維持する要素を示す。
- 確認観点: 人間レビューで確認する制約または注意を示す。

## シーケンス図

```mermaid
sequenceDiagram
    participant Caller as 呼び出し元
    participant Changed as 変更予定
    participant Target as 接続先

    Caller->>Changed: 変更後の呼び出し
    Changed->>Target: 維持する接続
    Target-->>Changed: 応答
    Changed-->>Caller: 返却
```

## 検証

- Mermaid 記述確認: Mermaid コードブロック、図種別、箱または参加者、接続を確認する。
