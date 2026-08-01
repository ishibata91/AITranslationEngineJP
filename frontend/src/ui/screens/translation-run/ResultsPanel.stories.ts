import type { Meta, StoryObj } from "@storybook/svelte-vite"
import ResultsPanel from "./ResultsPanel.svelte"
import type { NarrationResultRow } from "./translation-run-view"

const ROWS: NarrationResultRow[] = [
  {
    edid: "DLC1BookSerana",
    source:
      "I have walked these halls for centuries, and still the cold of Castle Volkihar finds me.",
    dest: "私は何世紀もこの広間を歩いてきたが、それでもヴォルキハル城の冷気は私を捉えて離さない。",
    statusLabel: "仮訳"
  },
  {
    edid: "DLC1BookAetherium",
    source:
      "The Aetherium Forge lies beneath Blackreach, its fires long since dimmed by the dwarves who fled.",
    dest: "",
    statusLabel: "未訳"
  }
]

// ページャ確認用の 1 ページ分（話者ありの口調差を含む）。総件数 121 を別に渡してページングを表す。
const PAGE_ROWS: NarrationResultRow[] = [
  {
    edid: "AventusAretinoInnocenceLost",
    source: "Mother? Are you there? It's so cold without you.",
    dest: "母さん？ そこにいるの？ あなたがいないと、こんなに寒いんだ。",
    statusLabel: "仮訳",
    personaLabel: "声質: 幼い少年の声",
    directive:
      "この台詞の話者の人物像:\n- 声質: 幼い少年の声\n- 種族の気質: ノルド（実直で粘り強い北方の気質）\nこの人物像に合う口調と人称で訳すこと。"
  },
  {
    edid: "GrelodTheKindScold",
    source:
      "You are all here to serve and honor your matron. Now off to bed, the lot of you!",
    dest: "お前たちは院母に仕え、敬うためにここにいるんだよ。さあ、とっとと寝な、お前たち全員！",
    statusLabel: "仮訳",
    personaLabel: "声質: しわがれた老女の声",
    directive:
      "この台詞の話者の人物像:\n- 声質: しわがれた老女の声\n- 種族の気質: ノルド（実直で粘り強い北方の気質）\nこの人物像に合う口調と人称で訳すこと。"
  },
  {
    edid: "HonorhallDoorActivate",
    source: "The door is locked.",
    dest: "扉には鍵がかかっている。",
    statusLabel: "仮訳"
  }
]

const meta = {
  title: "UI Components/ResultsPanel",
  component: ResultsPanel,
  parameters: { layout: "padded" },
  args: {
    onPrev: () => {},
    onNext: () => {},
    onUntranslatedOnlyChange: () => {},
    onExport: () => {}
  }
} satisfies Meta<typeof ResultsPanel>

export default meta
type Story = StoryObj<typeof meta>

export const EmptyIdle: Story = {
  name: "空（未実行）",
  args: { phase: "idle", results: [] }
}

export const EmptyRunning: Story = {
  name: "空（実行中）",
  args: { phase: "running", results: [] }
}

export const WithResults: Story = {
  name: "結果あり（書き出し可能）",
  args: { phase: "done", results: ROWS }
}

// 対象 plugin 全体から未訳だけを取得した状態。一覧、件数、チェック状態が同じ条件になる。
export const UntranslatedOnly: Story = {
  name: "未訳のみ",
  args: {
    phase: "done",
    results: ROWS.filter((row) => row.statusLabel === "未訳"),
    total: 1,
    untranslatedOnly: true,
    hasUnfilteredResults: true
  }
}

// 絞り込み前には訳のある結果が存在するが、未訳は 0 件。書き出し操作は残す。
export const NoUntranslatedResults: Story = {
  name: "未訳なし（書き出し可能）",
  args: {
    phase: "done",
    results: [],
    total: 0,
    untranslatedOnly: true,
    hasUnfilteredResults: true
  }
}

// 書き出し実行中。ボタンは無効化され、スピナーと「書き出し中…」を出す。
export const Exporting: Story = {
  name: "書き出し中",
  args: { phase: "done", results: ROWS, exporting: true }
}

// 複数ページの先頭。総件数 121、前へ無効・次へ有効。
export const PagingFirst: Story = {
  name: "複数ページ・先頭（前へ無効）",
  args: {
    phase: "done",
    results: PAGE_ROWS,
    total: 121,
    pageNumber: 1,
    canPrev: false,
    canNext: true
  }
}

// 複数ページの中間。前へ・次へとも有効。
export const PagingMiddle: Story = {
  name: "複数ページ・中間（両方有効）",
  args: {
    phase: "done",
    results: PAGE_ROWS,
    total: 121,
    pageNumber: 2,
    canPrev: true,
    canNext: true
  }
}

// 複数ページの末尾。次へ無効。
export const PagingLast: Story = {
  name: "複数ページ・末尾（次へ無効）",
  args: {
    phase: "done",
    results: PAGE_ROWS,
    total: 121,
    pageNumber: 3,
    canPrev: true,
    canNext: false
  }
}
