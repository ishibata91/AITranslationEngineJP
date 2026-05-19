# Storybook Review: 2026-05-17-all-pages-componentization

- `status`: browser-confirmed
- `review_target`: 全ページ部品化後の主要 panel、card、modal story
- `review_url`: `http://127.0.0.1:6006`
- `iframe_url_pattern`: `http://127.0.0.1:6006/iframe.html?id=<story-id>&viewMode=story`
- `storybook_build`: `pass`
- `storybook_build_command`: `npm --prefix frontend run build-storybook`
- `static_preview_command`: `python3 -m http.server 6006 --bind 127.0.0.1`
- `browser_confirmation`: `pass`
- `browser_confirmation_command`: `node -e "<playwright story smoke>"`
- `browser_confirmation_note`: 初回の Node REPL 経由 Playwright 起動は sandbox 権限で失敗した。権限付きの同等確認では対象 story がすべて描画された。
- `frontend_human_review`: `pending`
- `ux_review`: `pending`

## Review Boundary

- Storybook は props、callback stub、view model fixture だけで表示する。
- Storybook review URL は localhost の Storybook static preview だけを指す。
- Storybook は backend、Wails runtime、Gateway、RuntimeEventAdapter、AI provider、secret store、DB、実 filesystem flow を要求しない。
- Page 合成 story は密度確認と配置確認に限定する。主要部品 story の不足は Page 合成 story で代替しない。
- review 証跡には secret、API key、token、実 endpoint、実ユーザーデータ、raw request、raw response、raw prompt、provider 応答原文を含めない。

## Page Review Order

| page | primary story ID | review URL | status | unconfirmed reason |
| --- | --- | --- | --- | --- |
| 翻訳入力確認 | `screens-translation-input-dataloadhero--ready` | `http://127.0.0.1:6006/iframe.html?id=screens-translation-input-dataloadhero--ready&viewMode=story` | `confirmed` | `N/A` |
| ジョブ作成 | `screens-translation-job-setup-inputsourcepanel--selected` | `http://127.0.0.1:6006/iframe.html?id=screens-translation-job-setup-inputsourcepanel--selected&viewMode=story` | `confirmed` | `N/A` |
| 翻訳ジョブ管理 | `screens-translation-job-management-jobcard--selected-running` | `http://127.0.0.1:6006/iframe.html?id=screens-translation-job-management-jobcard--selected-running&viewMode=story` | `confirmed` | `N/A` |
| ジョブ実行 | `screens-job-run-jobruntargetsummary--term-phase` | `http://127.0.0.1:6006/iframe.html?id=screens-job-run-jobruntargetsummary--term-phase&viewMode=story` | `confirmed` | `N/A` |
| 単語翻訳段階 | `screens-term-translation-phase-termresultsummarycard--default` | `http://127.0.0.1:6006/iframe.html?id=screens-term-translation-phase-termresultsummarycard--default&viewMode=story` | `confirmed` | `N/A` |
| NPC ペルソナ生成段階 | `screens-persona-generation-phase-personaresultsummarycard--default` | `http://127.0.0.1:6006/iframe.html?id=screens-persona-generation-phase-personaresultsummarycard--default&viewMode=story` | `confirmed` | `N/A` |
| 本文翻訳段階 | `screens-body-translation-phase-bodyinputsummarycard--default` | `http://127.0.0.1:6006/iframe.html?id=screens-body-translation-phase-bodyinputsummarycard--default&viewMode=story` | `confirmed` | `N/A` |
| 翻訳完了 | `screens-job-run-translationcompletesummarypanel--default` | `http://127.0.0.1:6006/iframe.html?id=screens-job-run-translationcompletesummarypanel--default&viewMode=story` | `confirmed` | `N/A` |
| 出力成果物 | `screens-translation-output-artifact-outputactionpanel--ready` | `http://127.0.0.1:6006/iframe.html?id=screens-translation-output-artifact-outputactionpanel--ready&viewMode=story` | `confirmed` | `N/A` |
| マスター辞書 | `screens-master-dictionary-dictionarylistpanel--normal` | `http://127.0.0.1:6006/iframe.html?id=screens-master-dictionary-dictionarylistpanel--normal&viewMode=story` | `confirmed` | `N/A` |
| Provider 設定 | `screens-provider-settings-providerlistpanel--mixed` | `http://127.0.0.1:6006/iframe.html?id=screens-provider-settings-providerlistpanel--mixed&viewMode=story` | `confirmed` | `N/A` |
| Dashboard | `screens-dashboard-dashboardentrygrid--standard` | `http://127.0.0.1:6006/iframe.html?id=screens-dashboard-dashboardentrygrid--standard&viewMode=story` | `confirmed` | `N/A` |
| Master Persona | `screens-master-persona-personareviewpanel--with-list` | `http://127.0.0.1:6006/iframe.html?id=screens-master-persona-personareviewpanel--with-list&viewMode=story` | `confirmed` | `N/A` |

## Additional Story Coverage

- shared primitive story は `ui-components-*` として `ActionButton`、form、feedback、status、file、progress、pagination、confirm modal、phase primitive を持つ。
- 翻訳入力確認は file 未選択、選択済み、読み込み中、失敗、登録後導線を持つ。
- ジョブ作成は入力 source、基盤データ、phase 設定、互換性確認、作成済み summary、設定 summary を持つ。
- 翻訳ジョブ管理は job card、job list、stepper、delete modal の通常、空、失敗、処理中を持つ。
- マスター辞書と Provider 設定は list、detail、modal、action panel の通常、空、失敗、長文を持つ。

## Validation Evidence

| command | result | note |
| --- | --- | --- |
| `npm --prefix frontend run build-storybook` | `pass` | Storybook static build が成功した。 |
| `python3 scripts/harness/run.py --suite frontend-local` | `pass` | lint、type check、export check、boundary test、frontend test が通過した。 |
| `npm --prefix frontend run lint:boundaries` | `pass` | 30 tests が通過した。 |
| `npm --prefix frontend run lint:types` | `pass` | TypeScript type check が通過した。 |
| `node -e "<playwright story smoke>"` | `pass` | 主要 13 story は本文あり、page error なしで描画された。 |

## Human Review State

- UX review は未実施である。
- frontend human review は未承認である。
- 人間レビューでは、この `storybook-review.md` の review URL と story ID を入口にする。
