# 観察記録: master-dictionary 境界結合テスト pilot

- `skill`: tests-scenario（暫定。docs draft 追加のため最終的に finalization-module で取り扱う）
- `source_handoff`: handoff-observation-record
- `source_impl_scope`: implementation-scope.md の `### handoff_id: handoff-observation-record` 節
- `source_detail_spec`: detail-spec-diff.md の `boundary-test-REQ-006`
- `status`: draft（正本反映は finalization-module で扱う）
- `wave`: wave-3

本書は、master-dictionary 境界結合テスト pilot の実装中に得られた事実を記録する観察記録 draft である。
別 task `workflow-contract-maintenance` 起動時に、既存 skill 改訂（案 A）か新規 skill 追加（案 B）か保留（案 C）の判断材料として読む形式で書く。
各観察項目は「観察事実 → 判断材料 → 暫定推奨」の 3 段構造とする。

---

## 1. test-design の入力フォーマットに不足項目が出たか

### 観察事実

`test-design` skill が想定する通常の観点表は CSV 形式（`ID, 関連UC, 対象画面, 前提条件, 手順, 期待値, 備考`）を基本とする。
境界結合テストの観点表を作成する過程で、次の 2 列が境界結合テストに適合しないことが明らかになった。

- `対象画面` 列: 境界結合テストには UI 画面という概念がない。Wails controller method が入口であり、「対象画面」に当たる概念は「対象 method」に置き換えが必要になる。
- `手順` 列: 通常の E2E 観点表における「手順」はユーザー操作手順を指すが、境界結合テストでは「method 呼び出しと入力相当の説明」に置き換えが必要になる。

代替として、本 task では通常の CSV 列を使わず、独自の観点表形式（列: `ID, 関連仕様, 対象 method, case 分類, 入力条件, 観測 field, backend 期待値, frontend 期待値, 備考`）を `test-design.md` 内に採用した。
この代替形式は test-design skill に定義された形式ではなく、本 pilot 独自の対応である。

### 判断材料

- 不足した項目: `対象画面`（→ `対象 method` に置き換えが必要）、`手順`（→ method 呼び出しの説明に置き換えが必要）の 2 列。
- 情報欠落の深刻度: `対象画面` 列の欠如は境界結合テストの本質的な特性（UI を持たない）に起因するため、既存 CSV 形式の拡張では吸収できない。
- 代替で成立したか: 独自の観点表形式により観点の記述は成立した。ただし、`test-design` skill の形式規約との乖離が発生している。

### 暫定推奨

`workflow-contract-maintenance` 起動時に「test-design skill の入力フォーマットを、UI を持たない test 種別（境界結合テストなど）に拡張する」か「境界結合テスト専用の観点表形式を別途定義する」かを判断することを推奨する。
不足項目が 2 件であり、内容が構造的（UI を持たない test 種別）であるため、既存 skill の拡張で吸収しきれない可能性がある。

---

## 2. implementation-module の decision table に追加工程が必要だったか

### 観察事実

`implementation-module` の decision table には「シナリオテスト / 単体テスト / 観測ログ / 最終検証」の工程が並ぶ。
境界結合テスト pilot を実装する過程で、次の工程が既存の decision table に存在しないことが明らかになった。

- **境界結合テスト（backend 側）**: `internal/apitest/master_dictionary_boundary_test.go` を追加する工程。既存の「シナリオテスト」工程に最も近い区分として暫定的に寄せて運用した。
- **境界結合テスト（frontend 側）**: `frontend/src/controller/wails/master-dictionary.boundary.test.ts` を追加する工程。同様に「シナリオテスト」寄せで暫定運用した。
- **golden 管理工程**: 共有 golden の作成（`internal/apitest/testdata/boundary/master_dictionary/*.golden.json`）と、両側 helper（`boundary_golden_loader.go`、`boundary-golden-loader.ts`）の追加が「統合境界実装」として扱われたが、この区分は既存 decision table にあるものの golden 管理という独立工程の性格を持つ。

工程の入出力として明確になった点:

| 工程 | 入力 | 出力 |
| --- | --- | --- |
| golden 作成 | 境界 API 仕様 draft、pilot 対象 method と代表 case | `*.golden.json` ファイル群、両側 helper |
| 境界結合テスト backend 側 | golden ファイル群、backend 側 helper | `*_boundary_test.go`、`go test` 通過確認 |
| 境界結合テスト frontend 側 | golden ファイル群、frontend 側 helper | `*.boundary.test.ts`、`vitest` 通過確認 |
| golden 更新（仕様変更時） | 改訂済み境界 API 仕様 draft、旧 golden | 新 golden、両側 test の更新 |

### 判断材料

- 追加が必要だった独立工程の数: 「境界結合テスト（backend 側）」「境界結合テスト（frontend 側）」「golden 管理（作成・検証・更新）」の 3 工程。
- 暫定運用のギャップ: 「シナリオテスト」区分では「業務シナリオの受け入れ条件証明」が責務であり、「golden 一致 assert による DTO field 値 semantic の固定」とは責務が異なる。decision table を「シナリオテスト」で読み替えても、golden 管理工程は表現しきれない。

### 暫定推奨

`workflow-contract-maintenance` 起動時に「境界結合テスト（backend 側）」「境界結合テスト（frontend 側）」「golden 管理」の 3 工程を decision table に追加するかどうかを判断することを推奨する。
暫定運用として「シナリオテスト」寄せで動いたが、golden 管理工程の独立性が高く、既存工程区分で表現しきれない。

---

## 3. 既存 skill 本文の文言が境界結合テストの責務と衝突したか

### 観察事実

`tests-scenario` skill の次の文言と、境界結合テストの実作業が衝突した。

**衝突箇所 1: `tests-scenario` SKILL.md — 「担当ロールが判断してよい範囲」節**

> `UIが入口のシナリオでは、画面操作、ファイル選択、フォーム入力などのユーザー入力を開始点にする`

境界結合テストの入口は「Wails controller method の直接呼び出し」または「wails runtime mock への golden 供給」であり、UI 操作を開始点にする責務とは異なる。本 pilot では「入力開始点 = wails mock 応答、主要観測点 = gateway が返す field 値」と読み替えて対応したが、skill 文言どおりに適用すると停止条件に抵触する形になる。

**衝突箇所 2: `tests-scenario` SKILL.md — 「作業を止める条件」節**

> `シナリオ結果、公開接点、入力開始点、主要観測点のいずれかが不足している時`

「入力開始点 = wails mock 応答（UI なし）」「主要観測点 = DTO field 値の golden 一致」という境界結合テスト固有の形を skill の停止条件と突き合わせると、形式的には「UI が入口の場合の UI 人間操作 E2E」の定義に合致しない。本 pilot では読み替え対応を行い停止しなかったが、strict に解釈すると停止対象になりうる。

`tests-unit` skill との衝突は発生しなかった。境界結合テストは「公開振る舞い・分岐・エラー経路」の証明ではなく「DTO 値 semantic の両側整合」であり、`tests-unit` の責務範囲とは明確に異なるため、`tests-unit` を適用する判断には至らなかった。

### 判断材料

- 衝突した skill: `tests-scenario`（SKILL.md）
- 衝突箇所: 「UI が入口の場合の開始点定義」と「作業を止める条件（入力開始点不足）」の 2 箇所。
- 対応方法: 「wails mock を入口」と「DTO field 値 = 主要観測点」と読み替えて暫定運用した。読み替えの根拠は implementation-scope.md の `handoff-boundary-test-frontend` の `notes` に記録されている。

### 暫定推奨

`workflow-contract-maintenance` 起動時に `tests-scenario` skill の「UI が入口の場合の開始点定義」と「停止条件の入力開始点不足」が境界結合テストに適合するかを確認し、必要に応じて改訂または新規 skill でカバーする判断を推奨する。

---

## 4. golden 更新手順が既存 skill 群と独立工程として成立したか

### 観察事実

既存 `tests-unit` および `tests-scenario` は期待値の更新を「仕様変更に伴うテスト修正」として一元的に扱う。
境界結合テストでは、golden という中間成果物が存在し、更新手順が次の独立した 2 段構造になることが確認された。

1. **golden 更新**: 仕様変更の内容を `*.golden.json` ファイルに反映する。この時点では両側 test（backend / frontend）は変更していない。
2. **両側 test の更新**: golden が変わった場合、backend test（`go test` で失敗）と frontend test（`vitest` で失敗）の両方を更新する必要がある。

この 2 段更新は「仕様変更 → test 修正」という既存の一段手順とは異なり、「仕様変更 → golden 更新 → backend test 更新 → frontend test 更新」という独立工程として成立した。

golden を read-only として扱う運用は本 pilot で採用した設計であり、golden は backend test も frontend test も「読むだけ」とした。更新権限は golden 管理工程が持つ形になる。

**片側書き換え検出の成立状況**: wave-2 で手動確認を行い、`list_normal.golden.json` の field 値を 1 件書き換えた後に backend `go test` および frontend `vitest` の両方が失敗し、revert で再 pass することを確認した。この成立条件は次のとおり。

- frontend 側 test では、golden の `expected` を wails mock 応答として登録し、assert 期待値はテストソースコードにリテラルでハードコードする設計を採用している。mock 応答のみ golden から供給するため、golden を書き換えると mock 応答が変わり、ハードコードされたリテラル期待値との assert で test が落ちる。
- backend 側 test では、golden の `expected` を JSON decode して controller 応答と突き合わせるため、golden を書き換えると期待値が変わり test が落ちる。

### 判断材料

- 独立工程として成立したか: 成立した。「golden 更新 → 両側 test 更新」の 2 段が独立工程として機能する。
- 既存 skill との差異: 既存 `tests-unit` / `tests-scenario` の「仕様変更に伴うテスト修正」は 1 段だが、境界結合テストは golden という中間成果物を介して 2 段になる。この差異は既存 skill の文言では表現されていない。
- 片側書き換え検出の成立条件: frontend 側 test は「golden を mock 応答として使い、assert 期待値はリテラルでハードコード」する設計が必要。この設計を採用していない場合（golden を mock 応答と assert 期待値の両方に使う設計）では、golden を書き換えても両方が同時に変わるため検出経路が成立しない。

### 暫定推奨

`workflow-contract-maintenance` 起動時に「golden 更新手順（2 段更新）」を skill の規約として明文化するかどうかを判断することを推奨する。
独立工程として成立しているため、golden 管理の規約（意図的変更時の手順、壊れた時の対応 flow）を skill または正本 docs に含める候補になる。

---

## 5. 既存 skill 改訂で吸収しきれない数が閾値を超えるか

### 観察事実

`test-design.md` の観察ポイント（5.5 節）で定めた閾値判定基準は次のとおり。

> 5.1〜5.4 のうち 2 件以下の衝突・不足: 既存 skill 改訂（案 A）で対応可能と判定する根拠になる
> 5.1〜5.4 のうち 3 件以上の衝突・不足: 新規 skill `tests-boundary` の追加（案 B）を検討する根拠になる

本 pilot における衝突・不足の集計（観察項目 1〜4 から）:

| 観察項目 | 衝突・不足の有無 | 内容 |
| --- | --- | --- |
| 1. test-design 入力フォーマット | あり | `対象画面` / `手順` 列が境界結合テストに不適合（2 列） |
| 2. implementation-module decision table | あり | 3 工程（境界結合テスト backend 側 / frontend 側 / golden 管理）が未定義 |
| 3. 既存 skill 本文との衝突 | あり | `tests-scenario` の開始点定義と停止条件の 2 箇所 |
| 4. golden 更新手順の独立性 | あり（独立工程として成立） | 既存 skill の「仕様変更に伴うテスト修正」と異なる 2 段更新が必要 |

衝突・不足の件数: 4 件（観察項目 1〜4 全て）。閾値（3 件以上）を超えている。

### 判断材料

- 閾値: `test-design.md` 5.5 節の判定基準「3 件以上で新規 skill `tests-boundary` の追加を検討する根拠になる」。
- 現在の件数: 4 件。閾値超過。
- 案 A（既存 skill 改訂）: `tests-scenario`、`test-design`、`implementation-module` の 3 skill を改訂する必要がある。改訂箇所が複数 skill にまたがり、境界結合テストの概念が既存 skill の責務構造に収まらない部分がある。
- 案 B（新規 skill `tests-boundary` 追加）: golden 管理、両側 test の独立工程、観点表の形式を新規 skill で定義することで、既存 skill を改訂せずに境界結合テストの責務を独立した形で扱える。
- 案 C（保留）: 現状の暫定運用（`tests-scenario` 寄せ）を継続し、次の feature で境界結合テストを適用する際に再判断する。pilot 1 件で判断するリスクを回避できるが、skill との乖離が積み上がる。

### 暫定推奨

閾値（3 件以上）を超えたため、案 B（新規 skill `tests-boundary` 追加）を `workflow-contract-maintenance` 起動時の起点推奨とする。

ただし案 A（既存 skill 改訂）も完全に排除しない。理由は次のとおり。
- `tests-scenario` の衝突箇所は「UI が入口の場合の開始点定義」であり、条件文として「UI を持たない境界結合テストの場合は〜」を追記する改訂で吸収できる可能性がある。
- `implementation-module` の decision table については、「境界結合テスト」を独立行として追加する改訂で対応できる可能性がある。

最終的な案 A / 案 B の選択は `workflow-contract-maintenance` 起動時に、pilot の次 feature（AI service 境界または SQL repository 境界への適用）の見通しも踏まえて判断することを推奨する。

---

## 6. 片側書き換え検出を自動化する場合の必要工程・必要場所

### 観察事実

**手動確認した片側書き換え検出経路（wave-2 で確認済み）**:

1. `list_normal.golden.json` の `entries[0].source` 値を書き換える（例: `"Whiterun"` → `"TestCity"`）。
2. backend: `go test -run TestBoundary_MasterDictionary ./internal/apitest/...` が失敗する。`assertSummaryField` が `entries[0].source: want "Whiterun", got "TestCity"` を報告する。
3. frontend: `vitest run src/controller/wails/master-dictionary.boundary.test.ts` が失敗する。mock 応答の `entries[0].source` が `"TestCity"` になるが、ハードコードされたリテラル期待値 `"Whiterun"` との assert が失敗する。
4. revert 後、両 test が再 pass することを確認した。

この成立条件は観察項目 4 で記述したとおり、frontend 側 test が「golden を mock 応答として使い、assert 期待値はリテラルでハードコード」する設計を採用していることに依存する。

**自動化候補の観察事実（本 task では自動化判断を行わず、材料として列挙する）**:

- **CI hook（go test）**: backend `go test ./internal/apitest/...` は CI step として存在し、golden を `go:embed` で読み込む。CI が通る条件として、backend 側の golden 一致が自動的に検証される。CI step 上では現在の手動確認経路が自動化されている状態に近い。
- **CI hook（vitest）**: frontend `vitest run` は CI step として存在する。frontend 側の golden 一致が自動的に検証される。CI step 上では手動確認経路が自動化されている状態に近い。
- **pre-commit hook**: golden 更新時に両側 test を実行する pre-commit hook を設定することで、片側のみ更新して commit しようとした場合に警告できる。ただし pre-commit hook は回避可能であり、強制力は CI より低い。
- **golden generator**: 仕様変更時に golden を自動生成する仕組みを導入することで、golden 更新を手動操作ではなく generate step に変える。この場合、generator が両側の期待値構造と整合する golden を生成することが成立条件になる。

**各候補の実装場所候補**:

| 候補 | 実装場所 | 強制力 | 備考 |
| --- | --- | --- | --- |
| CI hook（go test） | `.github/workflows/` 等の CI 設定、または既存 CI step への追加 | 高（CI pass が merge 条件） | backend 側は現在の CI で既に網羅している可能性がある |
| CI hook（vitest） | 同上 | 高（CI pass が merge 条件） | frontend 側も同様 |
| pre-commit hook | `.git/hooks/pre-commit` または husky 等の設定 | 中（回避可能） | 開発者の手元での早期検出 |
| golden generator | `Makefile` または `scripts/` 配下のスクリプト、または `go generate` directive | 低（実行は任意） | generator の設計が複雑になる可能性がある |

### 判断材料

- 手動確認は wave-2 で成立済みであり、片側書き換え検出の経路は物理的に存在する。
- CI hook（go test / vitest）は既存の CI step と重複する可能性が高く、追加コストが低い候補である。
- pre-commit hook は開発者の手元での早期検出に効果があるが、CI ほどの強制力を持たない。
- golden generator は更新手順を自動化する最も強力な手段だが、設計コストが最も高い。

### 暫定推奨

自動化の判断は本 task では行わない。`workflow-contract-maintenance` 起動時に次の観点で判断することを推奨する。

1. CI 上の既存 `go test` / `vitest` step が境界結合テストを網羅しているかを確認し、網羅している場合は自動化は既に成立していると判定できる。
2. 手動確認の経路（片側書き換え → 両側失敗 → revert で再 pass）を CI step として組み込むかどうかは、境界結合テストの件数が増えた時点で判断する。
3. golden generator の導入は、境界結合テストを 3 feature 以上に適用した後の負荷を見てから判断することを推奨する。

---

## 7. 結合テストは API 仕様 + 画面設計から「導出される」のではないか（user 気付き 2026-06-05）

### 観察事実

本 task の implementation-module 完了直後（2026-06-05）、user から次の指摘があった。

> 「境界 API テスト設計じゃなくて、API 仕様書じゃないかな必要なのは」
> 「API 仕様と画面設計から、本来は導出されるのが結合テストなのではないかと思った」

本 pilot では、独立した `test-design.md`（BCT-MDC-001..010）として境界結合テストの観点表を起こした。
ただし、観点表の中身は実質的に「pilot 6 method × 代表 case」の組み合わせであり、上流が `boundary-api-draft.master-dictionary.md`（境界 API 仕様 draft）に集約されていた。
画面設計の差分は本 task では発生していない（plan 想定 2 N）が、画面変更が伴う task では画面設計差分も導出元になる可能性がある。

本 pilot の `test-design.md` は次の点で「上流からの導出」性を持っていた事後的に確認できる。

- `関連仕様` 列は `boundary-api-draft` および `detail-spec-diff` の REQ ID を参照していた。
- `対象 method` / `観測 field` / `backend 期待値` / `frontend 期待値` は境界 API 仕様 draft の DTO 定義から決定された。
- `case 分類`（正常 / 空集合 / 不在 / 状態遷移 3 連鎖 / REC 取捨）は detail-spec-diff の `master-dictionary-REQ-pilot-001` で固定済みの代表 case と一致していた。

つまり、`test-design.md` は独立した設計成果物として書かれているが、内容は「API 仕様書から機械的に展開可能な検証点」と「画面設計差分から導出される UI 起点シナリオ」の組み合わせとして再構成できる可能性がある。

### 判断材料

- 境界結合テストの「観点」は本来、上流（API 仕様書、画面設計差分）から導出される性格を持つ。pilot で独立 `test-design.md` を起こしたのは、API 仕様書が未確立だったための暫定対応であった可能性がある。
- 「導出」の枠組みを skill 上に固定するには、(a) API 仕様書の正本フォーマット、(b) 画面設計差分の正本フォーマット、(c) 両者からの境界結合テスト導出ルール、の 3 つを整える必要がある。
- 現状の本 task 出口では `boundary-api-draft.master-dictionary.md` が API 仕様書相当の draft として固定されている。finalization-module で正本化する候補先は `docs/detail-specs/master-dictionary.md` 内「境界 API」節（detail-spec-diff `boundary-test-REQ-005` の方針通り）。
- 画面設計差分は本 task では発生していないため、画面設計差分からの導出ルールは本 pilot では検証されていない。次 feature（UI 変更を伴う境界結合テスト適用 task）で確認する余地がある。
- 観察項目 1（`test-design` 入力フォーマット不足）と密接に関連する。観点 1 が「形式の不足」を扱うのに対し、観点 7 は「test-design という独立成果物の存在意義」を扱う上位の論点になる。

### 暫定推奨

`workflow-contract-maintenance` 起動時に次の判断を行うことを推奨する。

1. **API 仕様書の正本フォーマット固定**: `boundary-api-draft.master-dictionary.md` の構造を `docs/detail-specs/<feature>.md` 内「境界 API」節の標準節構成として固定し、新規 feature でも同じフォーマットで書ける形にする。
2. **境界結合テスト導出ルールの skill 上固定**: 「API 仕様書の DTO 定義 + 列挙値 + 状態遷移規約 + エラー応答契約」から、境界結合テストの観点（method × case × field 一致 assert）を機械的に展開するルールを skill として整備するか判断する。
3. **`test-design` 独立成果物の位置付け再判断**: 境界結合テストに限って、`test-design.md` を独立成果物として要求するのではなく、「API 仕様書 + 画面設計差分から導出される観点表として自動生成される」位置付けに変える判断を検討する。

観察項目 5 で固定した「案 B（新規 skill `tests-boundary` 追加）暫定推奨」と矛盾しない。新規 skill `tests-boundary` の責務に「API 仕様書 + 画面設計差分から境界結合テストを導出する手順」を含める設計が、本観察項目の暫定推奨と整合する。

---

## 根拠

- `source`:
  - `docs/exec-plans/active/master-dictionary-boundary-test/implementation-scope.md`（handoff-observation-record 節、人間レビュー観点 4）
  - `docs/exec-plans/active/master-dictionary-boundary-test/detail-spec-diff.md`（boundary-test-REQ-006）
  - `docs/exec-plans/active/master-dictionary-boundary-test/test-design.md`（5.1〜5.5 節）
  - `internal/apitest/master_dictionary_boundary_test.go`（wave-2 backend test 成果物）
  - `frontend/src/controller/wails/master-dictionary.boundary.test.ts`（wave-2 frontend test 成果物）
  - `internal/apitest/boundary_golden_loader.go`（wave-1 backend helper 成果物）
  - `frontend/src/controller/wails/boundary-golden-loader.ts`（wave-1 frontend helper 成果物）
  - `internal/apitest/testdata/boundary/master_dictionary/*.golden.json`（wave-1 共有 golden 11 ファイル）
  - `docs/exec-plans/active/master-dictionary-boundary-test/boundary-api-draft.master-dictionary.md`（wave-1 境界 API draft）
- `validation`: wave-2 手動確認（片側書き換え検出）を含む。
- `正本反映先`: 別 task `workflow-contract-maintenance` 起動後に finalization-module で判断する。
