# vader-lexicon-swap

## 依頼要約

感情辞書の concrete 実装を、再配布不可の NRC EmoLex から MIT ライセンスの VADER lexicon へ差し替え、git 同梱可能にする。VADER の valence 絶対値を「感情表出の強さ」の二値ゲートに使う。差し替え後に単体テストで判定を確かめる。

- 分岐元 branch: master
- 分岐元 commit: c89afa68

## 完了定義

- 動かす範囲: `linefeatures.EmotionLexicon` 境界の concrete 実装が VADER lexicon 由来になる。強感情語判定 `IsStrongEmotion(word)` が VADER の valence 絶対値の閾値で行われる。VADER の lexicon データが `assets/` に MIT 表記付きで同梱され、実行時ダウンロードなしで読める。
- 観測点: 単体テスト（lexicon package）。既知の強語（valence 絶対値が高い語）が true、中立語が false になることを確かめる。bootstrap 経由の engine 起動が壊れないことは既存 backend テストで確かめる。
- goal 整合: goal は「NRC を VADER へ差し替え、git に載せられる状態にする」。完了定義は VADER 実装が実際に判定を返し、データが同梱されることを要求する。空実装・仮データで満たしたことにしない。
- close_conditions: 下記。

### close_conditions

- VADER lexicon が `assets/` に同梱され、ライセンス表記がある。
- `EmotionLexicon` 境界の VADER concrete 実装が単体テストで強語 true・中立語 false を返す。
- bootstrap が VADER 実装を注入し、`go build` と既存 backend テストが通る。

## 軽 / 重判定

- 画面が動くか: N。backend の感情辞書 concrete 実装の差し替えのみ。layout・文言・style・表示構造・story・fixture を変えない。
- docs/architecture.md 反映が要るか: N。`EmotionLexicon` 境界と層構成・依存方向・Wails 境界は不変。concrete 実装の中身だけ変える。
- 判定: 両方 N → 軽 task。design-module と storybook-module を bypass。preparation → implementation → finalization。

## 実装・最終検証結果

### 変更ファイル

- internal/lexicon/vader.go: VADER concrete 実装（LoadVADER・IsStrongEmotion・閾値 1.5）を追加。
- internal/lexicon/vader_test.go: 強/弱語・閾値境界・大文字小文字・欠損ファイルの単体テストを追加。
- internal/bootstrap/bootstrap.go: 感情辞書の concrete 注入を LoadNRC から LoadVADER へ、参照パスを assets/vader_lexicon.txt へ切り替え。
- assets/vader_lexicon.txt: VADER lexicon データ（MIT、7520 行）を同梱。
- assets/vader_lexicon.LICENSE: 出典・取得日・sha256・形式・MIT 全文を記録。
- assets/CLAUDE.md: assets マップに vader_lexicon.txt を追記。

### 残置

- internal/lexicon/nrc.go と cmd/poc-tone は NRC 依存のまま残す（プロダクト経路は bootstrap のみ。撤去は本 task スコープ外）。

### 検証結果

- npm run test:backend: 全 package ok。lexicon の新規単体テスト通過。
- 実データ観測（一時テスト、確認後削除）: angry/love/hate/fear/kill/joy = true、the/table/of/and/chair = false。期待どおり。
- lint: 今回変更ファイルは format/vet/static/arch/boundary/module すべて指摘なし。arch・boundary・module は OK。
- go build ./...: OK。

### lint 負債の解消（触った task で緑にする方針）

分岐点 master から存在した static lint の赤 6 件を、赤を残さず全て解消した。

- internal/engine/export.go: interface method の error を wrap（wrapcheck 4 件）。
- internal/core/termxml/export.go: error 文字列の大文字始まりを語順変更で回避（ST1005 1 件）。
- internal/api/export.go: 公開メソッド名 ExportXTranslatorXml を ExportXTranslatorXML へ改名（revive var-naming 1 件）。
- 波及追随: frontend/src/gateway/translation-gateway.ts の import・呼び出し名、Wails 生成 binding（frontend/wailsjs/go/api/App.js・App.d.ts、dev watch が再生成）。

解消後の検証:

- npm run lint:backend: OK - No warnings found（arch・boundary・module も OK）。
- npm run test:backend: 全 package ok。
- npm run frontendlint: passed。

### 結合テスト（統合オラクル 2 段）

当初は Go 側のみ走らせていた。指摘を受け C# 側も走らせ、両段の緑を確認した。

- Go 側（stage=integration、internal/harness/oracle_test.go、specs.json 照合）: npm run test:backend に含まれ ok。
- C# 側（stage=extraction、tools/extractor.Tests、dotnet test）: 34 合格・0 失敗。
- frontend（npm run test:frontend）: 2 passed。gateway・binding の改名追随後も緑。

## 正本化判断（finalization）

- docs/architecture.md 反映: 不要。EmotionLexicon 境界の concrete 実装差し替え（NRC→VADER）と Wails 公開メソッド名の cosmetic 変更（Xml→XML）で、層構成・依存方向・Wails 境界の構造は不変。
- 人間承認済みの恒久仕様: なし。正本反映は skip。
- frontend/wailsjs（生成 binding）は gitignore のため commit 対象外（実行時に再生成）。

## finalization 結果

- 作業 commit: f4c9a8ff（branch claude/vader-lexicon-swap）。
- local merge: git merge --no-ff で master へ取り込み。merge commit 6749bad8。conflict なし。
- merge 後検証（master 上、全緑）: npm run test:backend ok、npm run lint:backend 0 issues・boundary OK、npm run test:frontend 2 passed、dotnet test tools/extractor.Tests 34 合格・0 失敗。
- remote 変更なし（push・tag push は行わない）。

### 実 app 確認（npm run dev:wails:run、http://localhost:34115）

- 起動: VADER 辞書読み込みエラーで停止せず起動。bootstrap の LoadVADER が実 app で成功。
- 再翻訳: 既存の完了済みプラグイン（Innocence Lost - Quest Expansion.esp）を削除し、LM Studio（http://192.168.0.226:1234、モデル hy-mt2-7b）で再翻訳。197 件完了。
- 口調生成: 攻撃的な台詞（Grelod の "Riff-raff! ..."、"Do your worst!"）が「口調: ぞんざい」、素直な台詞（AventusAretino の "Finally! My prayers have been answered!"）が「口調: 平明」に分類。VADER の valence 絶対値ゲートが tone 分類経路で妥当に機能。
