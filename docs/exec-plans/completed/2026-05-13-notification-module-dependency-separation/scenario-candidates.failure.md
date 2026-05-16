# Scenario Candidates: 2026-05-13-notification-module-dependency-separation / failure

- `generator`: `failure`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `NOTIFY-FAIL`

## Generator Scope

- `viewpoint`: 失敗観点。通知送信失敗、送信不可、redaction 漏れ、provider raw payload 混入、通知失敗と application result の分離を扱う。
- `included_sources`: `plan.md`, `docs/architecture.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本本文、他 agent の候補成果物
- `generation_notes`: 最終シナリオ表、採否、統合、競合解消は `designer` に残す。

## Candidate Scenarios

### CAND-NOTIFY-FAIL-001 通知送信失敗を application result の失敗へ混ぜない

- `source requirement`: `plan.md:58-62`, `plan.md:104-105`, `plan.md:143-144`, `docs/architecture.md:171-181`, `docs/architecture.md:209-217`
- `viewpoint`: 失敗観点
- `candidate scenario id`: `CAND-NOTIFY-FAIL-001`
- `actor`: Backend UseCase を起動する呼び出し元
- `trigger`: job / phase run の状態保存が成功した後、`NotificationPort` の送信だけが失敗する。
- `failure start condition`: 通知 module より下流の transport 送信が失敗する。
- `rejected operation`: 通知送信失敗を application result の失敗として返すことを拒否する。
- `expected error`: application result は保存済み操作結果を返し、通知失敗は通知境界の失敗として扱う。
- `expected outcome`: 保存済み job / phase run 状態を成功から失敗へ巻き戻さない。
- `observable point`: UseCase の戻り値、job / phase run の保存状態、NotificationPort fake の送信失敗記録。
- `acceptance condition`: 通知送信失敗後も UseCase の synchronous response は操作結果を表し、通知失敗 reason は application result の失敗理由へ混入しない。
- `exclusion condition`: 状態保存そのものが失敗した場合は、この候補の対象外とする。
- `related detail requirement type`: `failure_handling_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: 通知 module 分離の中核失敗条件として採用候補にする。
- `conflict hint`: 通知失敗を UI 表示上の失敗として扱う候補がある場合、application result の意味と衝突する。

### CAND-NOTIFY-FAIL-002 送信不可時は通知を抑止し、状態判断へ混ぜない

- `source requirement`: `plan.md:79-83`, `plan.md:97-105`, `docs/architecture.md:203-217`, `docs/architecture.md:281-288`
- `viewpoint`: 失敗観点
- `candidate scenario id`: `CAND-NOTIFY-FAIL-002`
- `actor`: NotificationDispatcher
- `trigger`: runtime handle が未接続、通知先が無効、または transport が送信可能状態ではない。
- `failure start condition`: 通知送信可否が false になる。
- `rejected operation`: 送信不可を理由に状態遷移可否、terminal guard、provider response validation を判断することを拒否する。
- `expected error`: 通知は送信されず、送信不可は通知境界の扱いとして閉じる。
- `expected outcome`: UseCase と Service は通知先の可否を知らず、確定済み状態事実だけを `NotificationSinkPort` へ渡す。
- `observable point`: NotificationPort fake の未呼び出し、UseCase / Service の戻り値、状態保存結果。
- `acceptance condition`: 送信不可時に NotificationPort は呼ばれず、application result と job / phase run 状態は通知送信可否に影響されない。
- `exclusion condition`: 起動時 DI の不備で UseCase 自体が構築できない場合は、この候補の対象外とする。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `testability_requirement`
- `adoption hint`: 送信不可を送信失敗と分けて検証する候補にする。
- `conflict hint`: 通知先未接続を起動失敗にする設計候補がある場合、実行時送信不可の扱いと分ける必要がある。

### CAND-NOTIFY-FAIL-003 redaction 漏れを送信前に拒否する

- `source requirement`: `plan.md:79-83`, `plan.md:104-105`, `plan.md:143-156`, `docs/architecture.md:209-214`, `docs/architecture.md:271-276`
- `viewpoint`: 失敗観点
- `candidate scenario id`: `CAND-NOTIFY-FAIL-003`
- `actor`: NotificationDispatcher
- `trigger`: 通知事実に secret、API key、credential 参照実値、prompt 全文、翻訳本文全文が含まれる。
- `failure start condition`: redaction 前の通知事実に送信禁止値が含まれる。
- `rejected operation`: 送信禁止値を含む通知 payload を NotificationPort へ渡すことを拒否する。
- `expected error`: 送信前に redaction される。redaction できない値は送信しない。
- `expected outcome`: transport payload には送信禁止値が含まれない。
- `observable point`: NotificationPort fake が受け取った payload、送信抑止結果、redaction 後 payload。
- `acceptance condition`: NotificationPort へ渡される payload は送信禁止値を含まず、redaction 不能時は NotificationPort が呼ばれない。
- `exclusion condition`: DB 永続化時の secret 保存禁止は、この候補の対象外とする。
- `related detail requirement type`: `security_requirement`, `failure_handling_requirement`, `testability_requirement`
- `adoption hint`: redaction の単体テスト候補として扱いやすい。
- `conflict hint`: prompt 全文や翻訳本文全文を UI に出す仕様候補がある場合、通知 payload 禁止条件と衝突する。

### CAND-NOTIFY-FAIL-004 provider raw payload の通知混入を拒否する

- `source requirement`: `plan.md:66-83`, `plan.md:104-105`, `plan.md:153-156`, `docs/architecture.md:197-214`, `docs/architecture.md:222-239`
- `viewpoint`: 失敗観点
- `candidate scenario id`: `CAND-NOTIFY-FAIL-004`
- `actor`: Service と NotificationDispatcher
- `trigger`: provider 呼び出し後の通知事実に provider raw payload が渡される。
- `failure start condition`: provider raw payload が通知候補値へ混入する。
- `rejected operation`: provider raw payload を Wails event payload または operation summary として送ることを拒否する。
- `expected error`: provider raw payload は通知 payload から除外される。必要な通知は確定済み状態事実から組み立てる。
- `expected outcome`: NotificationPort が受け取る payload は provider raw payload を含まない。
- `observable point`: Service fake の通知入力、NotificationDispatcher の整形結果、NotificationPort fake の payload。
- `acceptance condition`: provider raw payload が通知入力へ来ても、NotificationPort へ渡る payload には混入しない。
- `exclusion condition`: provider response validation の成否判断は、この候補の対象外とする。
- `related detail requirement type`: `security_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: external-integration 観点の provider 失敗候補と統合される可能性がある。
- `conflict hint`: provider raw payload を監査保存または通知表示する候補がある場合、保存禁止と通知禁止の条件が衝突する。

### CAND-NOTIFY-FAIL-005 unsafe payload の送信拒否は DB 保存対象を増やさない

- `source requirement`: `plan.md:97-105`, `plan.md:155-156`, `docs/architecture.md:197-214`
- `viewpoint`: 失敗観点
- `candidate scenario id`: `CAND-NOTIFY-FAIL-005`
- `actor`: NotificationDispatcher
- `trigger`: redaction 不能、provider raw payload 混入、または送信禁止値混入により通知送信を拒否する。
- `failure start condition`: 通知 payload を安全に作れない。
- `rejected operation`: notification result、operation summary、Wails event payload を DB に永続化することを拒否する。
- `expected error`: 通知送信は行わず、永続化対象は UseCase が確定した状態事実に限定する。
- `expected outcome`: notification result や transport payload の永続化が発生しない。
- `observable point`: repository fake の保存呼び出し、NotificationPort fake の呼び出し有無、job / phase run の保存内容。
- `acceptance condition`: unsafe payload 拒否後も、DB には notification result、operation summary、Wails event payload が保存されない。
- `exclusion condition`: 状態事実として定義済みの失敗 reason category の保存は、この候補の対象外とする。
- `related detail requirement type`: `data_requirement`, `security_requirement`, `failure_handling_requirement`
- `adoption hint`: 通知失敗と保存境界を結びつける候補にする。
- `conflict hint`: 運用・監査観点が通知失敗詳細の永続保存を求める場合、保存禁止条件との競合候補にする。

### CAND-NOTIFY-FAIL-006 UseCase が failure payload を直接組み立てる経路を拒否する

- `source requirement`: `plan.md:17-21`, `plan.md:56-62`, `plan.md:146-156`, `docs/architecture.md:169-181`, `docs/architecture.md:243-258`
- `viewpoint`: 失敗観点
- `candidate scenario id`: `CAND-NOTIFY-FAIL-006`
- `actor`: Backend UseCase
- `trigger`: 失敗結果または完了結果を返す処理で、UseCase が Wails event payload を直接組み立てようとする。
- `failure start condition`: application result と transport payload が同じ構築処理に混ざる。
- `rejected operation`: UseCase が Runtime adapter、NotificationDispatcher、NotificationPort、Wails event payload 形式を直接扱うことを拒否する。
- `expected error`: UseCase は application result と通知事実だけを扱い、通知 payload 化は NotificationDispatcher へ残す。
- `expected outcome`: 失敗系でも application result と transport payload が分離される。
- `observable point`: usecase package の依存、UseCase の戻り値、NotificationSinkPort fake が受け取る通知事実。
- `acceptance condition`: UseCase は Wails event payload を生成せず、通知事実だけを NotificationSinkPort へ渡す。
- `exclusion condition`: Runtime adapter 内の event name と transport payload 形式の検証は、この候補の対象外とする。
- `related detail requirement type`: `compatibility_requirement`, `consistency_requirement`, `failure_handling_requirement`
- `adoption hint`: import 境界 lint と単体テストの両方で検証できる候補にする。
- `conflict hint`: UseCase が通知表示文言を確定する候補がある場合、通知 module の責務境界と衝突する。

### CAND-NOTIFY-FAIL-007 通知失敗時に Controller へ途中経過通知を戻さない

- `source requirement`: `plan.md:32-35`, `plan.md:127-134`, `plan.md:140-151`, `docs/architecture.md:159-167`, `docs/architecture.md:203-220`, `docs/architecture.md:250-258`
- `viewpoint`: 失敗観点
- `candidate scenario id`: `CAND-NOTIFY-FAIL-007`
- `actor`: Controller
- `trigger`: NotificationPort 送信失敗または送信不可が発生する。
- `failure start condition`: 通知失敗の扱いを Controller 経由の同期応答で補おうとする。
- `rejected operation`: Controller が途中経過通知の戻り先になり、NotificationDispatcher を直接呼ぶことを拒否する。
- `expected error`: Controller は request / response DTO の写像と synchronous response に限定される。
- `expected outcome`: 通知失敗時も途中経過通知の経路は NotificationSinkPort から NotificationDispatcher へ閉じる。
- `observable point`: controller package の依存、Controller の response、NotificationDispatcher fake の呼び出し元。
- `acceptance condition`: Controller は NotificationDispatcher を呼ばず、通知失敗を補うための途中経過通知を response 経路へ混ぜない。
- `exclusion condition`: Bind call 自体の request validation 失敗は、この候補の対象外とする。
- `related detail requirement type`: `compatibility_requirement`, `failure_handling_requirement`, `testability_requirement`
- `adoption hint`: architecture lint 候補と対応する失敗シナリオにする。
- `conflict hint`: UI 表示が通知失敗を同期応答で受け取りたい場合、Controller 責務との競合候補にする。

## Open Notes

- `human decision candidate`: 通知送信失敗、送信不可、unsafe payload 拒否を恒久ログへ残すかは、現時点の根拠だけでは確定しない。
- `merge candidate`: provider raw payload 混入防止は external-integration 観点の候補と統合される可能性がある。
- `rejection candidate`: 正常系の通知成功確認だけの候補は、failure 観点からは除外する。
- `conflict candidate`: 通知失敗を利用者へ同期応答で表示する候補は、application result 分離と Controller 非通知経路の条件に衝突する可能性がある。
- `candidate count`: 7
