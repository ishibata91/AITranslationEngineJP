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
- 実行キャッシュ構築
- ジョブ生成
- ステージ実行
- 中断
- 再開
- リトライ
- ジョブ完了時キャッシュ削除
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
- SQLite 実行キャッシュ
- 永続マスターデータ保管
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

- プラグイン入力キャッシュ
- 翻訳単位
- 会話グループ
- NPC ペルソナ
- 用語辞書
- クエスト文脈
- 翻訳ジョブ
- マスターペルソナ
- マスター辞書

`dialogue_groups` は文脈付き翻訳の中心単位として独立したドメインモデルで扱う。

`PLUGIN_EXPORT` は xEdit JSON から構築する入力キャッシュの親モデルとして扱う。
入力キャッシュは再構築可能な一時データとし、JSON 原本はファイルシステムに保持する。
`MASTER_PERSONA` と `MASTER_DICTIONARY` は入力キャッシュとは別の永続基盤データとして扱う。

## 5. DTO 境界

フロントエンドとバックエンド間のデータ受け渡しは DTO を明示的に定義して行う。

DTO の対象は少なくとも以下を含む。

- プラグイン取込結果
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
- ジョブフェーズ種別は DB テーブルではなくアプリケーション定数として定義する
- DB の内部主キーはシーケンシャル整数を採用し、外部 FormID は別列で保持する

## 7. 永続化方針

- `PLUGIN_EXPORT` 配下の入力データは SQLite 上の実行キャッシュとして保持する
- 実行キャッシュは `TRANSLATION_JOB` が参照する
- `TRANSLATION_JOB` が `Completed`, `Canceled`, `Failed` のいずれかになり、同一 `PLUGIN_EXPORT` に未完了ジョブが残っていない場合は入力キャッシュを削除する
- JSON 原本は削除せず、必要時に再取り込み可能とする
- `MASTER_PERSONA`, `MASTER_PERSONA_ENTRY`, `MASTER_DICTIONARY`, `MASTER_DICTIONARY_ENTRY` はジョブ完了後も保持する
