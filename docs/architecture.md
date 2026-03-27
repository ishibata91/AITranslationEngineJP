# アーキテクチャ仕様

関連文書: [`spec.md`](./spec.md), [`tech-selection.md`](./tech-selection.md), [`interface-spec.md`](./external-design/interface-spec.md), [`execution-spec.md`](./external-design/execution-spec.md), [`ui-spec.md`](./external-design/ui-spec.md)

本書は、システムの内部構成と責務分割を定義する。

## 1. 層構成

システムは以下の 4 層で構成する。

- `UI層`
- `アプリケーション層`
- `ドメイン層`
- `インフラ層`

## 2. 各層の責務

### 2.1 UI層

- ジョブ作成
- 辞書閲覧と編集
- ペルソナ閲覧
- 補助知識表示
- 進捗表示
- ログビュー表示
- 出力操作

### 2.2 アプリケーション層

- xEdit JSON 取込
- ジョブ生成
- ステージ実行
- 中断
- 再開
- リトライ
- xTranslator XML 出力

### 2.3 ドメイン層

- 翻訳単位
- 会話グループ
- ペルソナ
- 用語辞書
- クエスト文脈
- 保護トークン

### 2.4 インフラ層

- LMStudio アダプタ
- Gemini アダプタ
- xAI アダプタ
- SQLite 永続化
- XML writer
- ファイルシステム

## 3. AI プロバイダ抽象化

AI プロバイダ境界は Rust の trait で定義する。

```rust
trait TranslationProvider {
    async fn translate(&self, request: TranslateRequest) -> Result<TranslateResponse>;
    async fn translate_batch(&self, request: TranslateBatchRequest) -> Result<TranslateBatchResponse>;
    async fn build_persona(&self, request: PersonaRequest) -> Result<PersonaResponse>;
    async fn build_dictionary(&self, request: DictionaryRequest) -> Result<DictionaryResponse>;
}
```

各プロバイダ実装は以下の機能境界に従う。

- 単発翻訳
- Batch 翻訳
- ペルソナ生成
- 辞書構築

## 4. ドメインモデル方針

内部モデルは少なくとも以下を持つ。

- 翻訳単位
- 会話グループ
- NPC ペルソナ
- 用語辞書
- クエスト文脈
- 保護トークン

`dialogue_groups` は文脈付き翻訳の中心単位として独立したドメインモデルで扱う。

## 5. DTO 境界

フロントエンドとバックエンド間のデータ受け渡しは DTO を明示的に定義して行う。

DTO の対象は少なくとも以下を含む。

- ジョブ作成入力
- ジョブ状態
- 進捗情報
- ログイベント
- 辞書表示データ
- ペルソナ表示データ
- 翻訳対象一覧データ

## 6. 型安全方針

- バックエンドの中核ロジックは Rust の型で定義する
- UI は TypeScript の型で定義する
- xEdit JSON はロード時に型検証する
- xTranslator XML は内部ドメインモデルから生成する

`<10gold>` のような保護対象は独立したトークン型として扱う。
