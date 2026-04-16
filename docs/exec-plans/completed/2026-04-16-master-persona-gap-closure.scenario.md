# Scenario テスト一覧: master-persona-gap-closure

- `task_id`: `persona-management-gap-closure`
- `usecase`: `persona-management`
- `final_artifact_path`: `docs/scenario-tests/master-persona-management.md`
- `対象`: master persona の production persistence、keyring secret access、provider interface / transport seam、real provider set、JSON auto preview、restart persistence
- `topic_abbrev`: `MPG`

## ケース一覧

### SCN-MPG-001 AI 設定と keyring secret access の再起動復元

- `分類`: 正常系
- `観点`: 保存済み provider、model、API key secret が再起動後に復元される
- `事前条件`: provider と model を保存済みで、API key が `github.com/99designs/keyring` backed `SecretStore` concrete に格納されている
- `手順`:
  1. マスターペルソナ画面で API key を入力し、保存ボタンを押す。
  2. OS キーチェーンの認可ダイアログを確認して保存する。
  3. アプリを終了して再起動する。
  4. 画面再表示後の AI 設定欄を確認する。
- `期待結果`:
  1. macOS では Keychain、Windows では Windows Credential Manager backed の secret store が使われる。
  2. provider と model が保存済み値で復元される。
  3. API key は keyring backed secret seam から読み戻され、通常利用では再入力や再認証を要求しない。
  4. API key は error message、一覧、detail、run status に露出しない。
  5. アプリ内の追加確認 modal は表示されず、明示認証は OS キーチェーンの認可ダイアログで行われる。
- `再実行条件`: 同じ設定を保存し直しても状態が破綻しない。

### SCN-MPG-002 JSON 選択で auto preview が起動し AI 設定完了時だけ生成可能へ遷移する

- `分類`: 正常系
- `観点`: JSON 選択だけで preview が開始され、AI 設定完了時は成功後に生成ボタンが有効になる
- `事前条件`: AI 設定が完了しており、対象ありの `extractData.pas` JSON を選択できる
- `手順`:
  1. `JSON を選ぶ` で有効な JSON を選択する。
  2. 選択直後の run 表示を確認する。
  3. preview 完了後の stats と生成ボタン状態を確認する。
- `期待結果`:
  1. 選択直後に `入力検証中` が表示される。
  2. preview card に file 名、対象 plugin、総 NPC 数、生成対象数、skip 内訳が反映される。
  3. AI 設定完了かつ preview 成功時だけ `preview.status = 生成可能` となり、追加クリックなしで生成ボタンが有効になる。
- `再実行条件`: `preview を更新` を押すと同じ JSON で再計算できる。

### SCN-MPG-003 AI 設定未完了でも auto preview の集計だけは表示する

- `分類`: 主要例外系
- `観点`: AI 設定未完了でも JSON 集計は観測でき、生成許可は出ない
- `事前条件`: AI 設定が未完了で、対象ありの `extractData.pas` JSON を選択できる
- `手順`:
  1. 設定未完了のまま `JSON を選ぶ` で有効な JSON を選択する。
  2. preview 完了後の message、stats、status を確認する。
  3. AI 設定を保存して `preview を更新` できるか確認する。
- `期待結果`:
  1. preview card に file 名、対象 plugin、総 NPC 数、生成対象数、skip 内訳が反映される。
  2. status は `設定未完了` のまま維持される。
  3. 生成ボタンは無効のままである。
  4. AI 設定を補うと、同じ JSON の集計を保ったまま生成可能判定を再評価できる。
- `再実行条件`: JSON を差し替えるか設定を補って再試行できる。

### SCN-MPG-004 fake は provider list に出さず provider interface / transport seam で差し替える

- `分類`: 正常系
- `観点`: fake path が provider 選択肢ではなく DI された provider transport seam として動作する
- `事前条件`: test mode が有効で、real provider list が `gemini` / `lm_studio` / `xai` で構成されている
- `手順`:
  1. test mode で master persona 生成を開始する。
  2. provider list と provider validation の表示を確認する。
  3. provider interface と request / SDK transport seam に fake が注入される条件を確認する。
- `期待結果`:
  1. provider list に fake provider は表示されず、Gemini、LM Studio、xAI だけが表示される。
  2. prompt 組み立て、provider validation、run orchestration は real provider と共通の処理を通る。
  3. real provider concrete は provider interface 経由で共通 response を返す。
  4. 実際の外部 request / SDK transport だけが DI fake に差し替わり、paid な real AI API は呼ばれない。
  5. service は provider 文字列 switch と本文固定生成に依存しない。
- `再実行条件`: `AITRANSLATIONENGINEJP_MASTER_PERSONA_AI_MODE=fake` で再実行しても決定論的な fake response が返る。

### SCN-MPG-005 生成結果と直近 run status が再起動後も残る

- `分類`: 状態遷移
- `観点`: 生成完了後の entry と run status が再起動で失われない
- `事前条件`: 有効 JSON で少なくとも 1 件生成できる
- `手順`:
  1. 生成を完了させる。
  2. 一覧に新規 entry が追加されたことを確認する。
  3. アプリを再起動し、同じ plugin を表示する。
- `期待結果`:
  1. 新規 entry が一覧と detail に残る。
  2. run panel に直近の完了 status と件数が残る。
  3. 再起動後も update / delete lock は残らない。
- `再実行条件`: 同じ JSON を再度 preview すると既存分が skip として観測される。

### SCN-MPG-006 再起動時に stale running status を補正する

- `分類`: 再実行系
- `観点`: `生成中` のままアプリを閉じても永続 lock を残さない
- `事前条件`: 生成中にアプリを終了できる
- `手順`:
  1. 生成を開始し、run status が `生成中` になったことを確認する。
  2. 途中でアプリを終了する。
  3. 再起動後の run status と detail action を確認する。
- `期待結果`:
  1. 保存済み `生成中` status は起動時に `中断済み` または同等の非 active 状態へ補正される。
  2. run panel には再起動起因の中断が分かる message が残る。
  3. detail の update / delete は再び操作可能になる。
- `再実行条件`: 同じ JSON で preview と generate を再開できる。

## 受け入れ観点への対応

- `production persistence replacement`: `SCN-MPG-001`, `SCN-MPG-005`, `SCN-MPG-006`
- `keyring secret access`: `SCN-MPG-001`
- `provider interface common response`: `SCN-MPG-004`
- `fake request / SDK transport seam via DI`: `SCN-MPG-004`
- `JSON selection -> auto preview -> generation eligibility`: `SCN-MPG-002`, `SCN-MPG-003`
- `restart persistence verification`: `SCN-MPG-001`, `SCN-MPG-005`, `SCN-MPG-006`
- `no in-memory production wiring accepted`: `SCN-MPG-001` の前提と acceptance で production concrete を要求する

## 未確定事項

- なし。implementation-scope は human review 後に確定済み。
