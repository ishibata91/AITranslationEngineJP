import type { ComponentProps } from "svelte"
import ConfirmDangerModal from "../ConfirmDangerModal.svelte"
import EmptyStatePanel from "../EmptyStatePanel.svelte"
import FileSelectionDisplay from "../FileSelectionDisplay.svelte"
import PaginationControls from "../PaginationControls.svelte"
import ProgressBar from "../ProgressBar.svelte"
import SearchFilterBar from "../SearchFilterBar.svelte"
import StatusPill from "../StatusPill.svelte"

type StatusPillProps = ComponentProps<typeof StatusPill>
type SearchFilterBarProps = ComponentProps<typeof SearchFilterBar>
type FileSelectionDisplayProps = ComponentProps<typeof FileSelectionDisplay>
type EmptyStatePanelProps = ComponentProps<typeof EmptyStatePanel>
type ProgressBarProps = ComponentProps<typeof ProgressBar>
type PaginationControlsProps = ComponentProps<typeof PaginationControls>
type ConfirmDangerModalProps = ComponentProps<typeof ConfirmDangerModal>

const noop = (): void => {}

export const statusPillDefaultFixture: StatusPillProps = {
  label: "準備完了",
  tone: "success"
}

export const statusPillBusyFixture: StatusPillProps = {
  label: "処理中",
  tone: "info"
}

export const statusPillFailureFixture: StatusPillProps = {
  label: "確認が必要",
  tone: "danger"
}

export const statusPillLongFixture: StatusPillProps = {
  label: "状態説明が長い場合でも折り返して一覧幅からはみ出さない",
  tone: "warning"
}

export const searchFilterDefaultFixture: SearchFilterBarProps = {
  searchId: "shared-search",
  searchLabel: "検索",
  searchValue: "dragon",
  filterId: "shared-filter",
  filterLabel: "状態",
  filterValue: "active",
  filterOptions: [
    { value: "all", label: "すべて" },
    { value: "active", label: "有効" },
    { value: "failed", label: "失敗" }
  ],
  placeholder: "名前で検索",
  onSearchInput: noop,
  onFilterChange: noop
}

export const searchFilterEmptyFixture: SearchFilterBarProps = {
  ...searchFilterDefaultFixture,
  searchValue: "",
  filterValue: "all"
}

export const searchFilterLongFixture: SearchFilterBarProps = {
  ...searchFilterDefaultFixture,
  searchValue: "長い検索語句が入力されても入力欄と条件欄の配置が崩れない状態",
  searchHelp: "一覧に表示される synthetic label だけを対象にします。"
}

export const fileSelectionDefaultFixture: FileSelectionDisplayProps = {
  title: "入力ファイル",
  fileName: "sample-plugin-translation.json",
  pathLabel: "project://workspace/input/sample-plugin-translation.json",
  hashLabel: "sha256:0000-synthetic-hash",
  statusLabel: "読み込み済み",
  statusTone: "success",
  actionLabel: "差し替える",
  onAction: noop
}

export const fileSelectionEmptyFixture: FileSelectionDisplayProps = {
  title: "入力ファイル",
  emptyMessage: "入力ファイルは未選択です",
  actionLabel: "選択する",
  onAction: noop
}

export const fileSelectionFailureFixture: FileSelectionDisplayProps = {
  ...fileSelectionDefaultFixture,
  statusLabel: "失敗",
  statusTone: "danger",
  error: "ファイル情報を表示できませんでした。"
}

export const fileSelectionBusyFixture: FileSelectionDisplayProps = {
  ...fileSelectionDefaultFixture,
  statusLabel: "確認中",
  statusTone: "warning",
  busy: true
}

export const fileSelectionLongFixture: FileSelectionDisplayProps = {
  ...fileSelectionDefaultFixture,
  fileName: "very-long-synthetic-file-name-for-layout-review-and-wrapping-check.json",
  pathLabel: "project://workspace/input/nested/folder/with/long/path/very-long-synthetic-file-name-for-layout-review-and-wrapping-check.json"
}

export const emptyStateDefaultFixture: EmptyStatePanelProps = {
  title: "表示する項目がありません",
  message: "条件を変更すると一覧に表示される可能性があります。",
  actionLabel: "条件を戻す",
  onAction: noop
}

export const emptyStateWarningFixture: EmptyStatePanelProps = {
  title: "入力が未選択です",
  message: "先に対象を選択すると次の操作へ進めます。",
  tone: "warning"
}

export const emptyStateFailureFixture: EmptyStatePanelProps = {
  title: "一覧を表示できません",
  message: "読み込み失敗の内容を確認してください。",
  tone: "danger",
  actionLabel: "再試行",
  onAction: noop
}

export const emptyStateBusyFixture: EmptyStatePanelProps = {
  ...emptyStateDefaultFixture,
  actionLabel: "確認中",
  busy: true
}

export const emptyStateLongFixture: EmptyStatePanelProps = {
  title: "長い説明を持つ空状態",
  message: "空状態の説明文が長くなっても、文章は枠内で折り返され、操作ボタンと重ならない必要があります。",
  actionLabel: "長い文言の操作を実行する",
  onAction: noop
}

export const progressZeroFixture: ProgressBarProps = {
  label: "処理準備",
  value: 0,
  helperText: "まだ開始していません。"
}

export const progressHalfFixture: ProgressBarProps = {
  label: "翻訳処理",
  value: 50,
  helperText: "半分まで完了しました。"
}

export const progressDoneFixture: ProgressBarProps = {
  label: "出力生成",
  value: 100,
  tone: "success",
  helperText: "完了しました。"
}

export const progressFailureFixture: ProgressBarProps = {
  label: "読み込み処理",
  value: 35,
  tone: "danger",
  helperText: "途中で失敗しました。"
}

export const progressLongFixture: ProgressBarProps = {
  label: "長い処理名が付いた進捗表示でも横幅の中で読める状態",
  value: 72,
  tone: "warning",
  helperText: "説明文も長い場合に、進捗バーの下で折り返して表示されます。"
}

export const paginationDefaultFixture: PaginationControlsProps = {
  page: 2,
  pageCount: 5,
  totalLabel: "42 件",
  onPrevious: noop,
  onNext: noop
}

export const paginationFirstFixture: PaginationControlsProps = {
  ...paginationDefaultFixture,
  page: 1
}

export const paginationLastFixture: PaginationControlsProps = {
  ...paginationDefaultFixture,
  page: 5
}

export const paginationBusyFixture: PaginationControlsProps = {
  ...paginationDefaultFixture,
  busy: true
}

export const paginationLongFixture: PaginationControlsProps = {
  ...paginationDefaultFixture,
  totalLabel: "検索結果 123456 件中 synthetic page を表示中"
}

export const confirmDangerDefaultFixture: ConfirmDangerModalProps = {
  open: true,
  title: "対象を削除しますか",
  targetLabel: "Synthetic Job 001",
  message: "削除した対象は一覧に戻りません。",
  onConfirm: noop,
  onCancel: noop
}

export const confirmDangerBusyFixture: ConfirmDangerModalProps = {
  ...confirmDangerDefaultFixture,
  busy: true
}

export const confirmDangerFailureFixture: ConfirmDangerModalProps = {
  ...confirmDangerDefaultFixture,
  error: "削除処理に失敗しました。状態を確認してから再試行してください。"
}

export const confirmDangerLongFixture: ConfirmDangerModalProps = {
  ...confirmDangerDefaultFixture,
  targetLabel: "Synthetic Job With Very Long Name For Modal Layout Review And Wrapping",
  message: "対象識別情報と説明文が長い場合でも、確認 modal は横幅を超えず、操作ボタンを押せる状態を維持します。"
}
