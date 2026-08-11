<script lang="ts">
  /* eslint-disable import-x/no-duplicates -- 値と型を別のimportとして明示する。 */
  import { onMount } from "svelte"
import {
    createTermDictionary,
    deleteTermDictionary,
    listTermDictionary,
  patchTermDictionary
} from "../../../gateway/term-dictionary-gateway"
  import type { TermDictionaryEntry } from "../../../gateway/term-dictionary-gateway"
  import PrebuiltDictionaryEditorScreen from "./PrebuiltDictionaryEditorScreen.svelte"
  import { PREBUILT_DICTIONARY_PAGE_SIZE } from "./prebuilt-dictionary-editor.fixtures"
  import type {
    PrebuiltDictionaryFilters,
    PrebuiltDictionaryRow
  } from "./prebuilt-dictionary-editor-view"

  const emptyFilters: PrebuiltDictionaryFilters = {
    source: "",
    destination: "",
    partOfSpeech: "",
    category: ""
  }

  type EditableField = "source" | "destination" | "partOfSpeech"
  type RowChanges = Partial<Pick<PrebuiltDictionaryRow, EditableField>>

  let rows = $state<PrebuiltDictionaryRow[]>([])
  let saved = $state<PrebuiltDictionaryRow[]>([])
  let revisions = $state<Record<string, number>>({})
  let filters = $state<PrebuiltDictionaryFilters>({ ...emptyFilters })
  let expandedRowIds = $state<string[]>([])
  let editingRowId = $state("")
  let pageNumber = $state(1)
  let totalCount = $state(0)
  let newRowSequence = 0

  const hasPendingChanges = $derived(rows.some((row) => row.pending != null || row.deletePending))

  function clone(source: PrebuiltDictionaryRow[]): PrebuiltDictionaryRow[] {
    return source.map((row) => ({ ...row, categories: [...row.categories] }))
  }

  function toViewRow(entry: TermDictionaryEntry): PrebuiltDictionaryRow {
    return {
      id: entry.id,
      source: entry.source,
      destination: entry.destination,
      partOfSpeech: entry.partOfSpeech,
      categories: entry.categories
    }
  }

  function findSavedRow(id: string): PrebuiltDictionaryRow | undefined {
    return saved.find((row) => row.id === id)
  }

  function changed(row: PrebuiltDictionaryRow, baseline: PrebuiltDictionaryRow): boolean {
    return row.source !== baseline.source ||
      row.destination !== baseline.destination ||
      row.partOfSpeech !== baseline.partOfSpeech
  }

  function rowChanges(row: PrebuiltDictionaryRow, baseline: PrebuiltDictionaryRow): RowChanges {
    const changes: RowChanges = {}
    if (row.source !== baseline.source) changes.source = row.source
    if (row.destination !== baseline.destination) changes.destination = row.destination
    if (row.partOfSpeech !== baseline.partOfSpeech) changes.partOfSpeech = row.partOfSpeech
    return changes
  }

  async function load(page = pageNumber): Promise<void> {
    const result = await listTermDictionary(filters, page)
    rows = result.entries.map(toViewRow)
    saved = clone(rows)
    revisions = Object.fromEntries(result.entries.map((entry) => [entry.id, entry.revision]))
    editingRowId = ""
    pageNumber = result.pageNumber
    totalCount = result.totalCount
  }

  function onFilterInput(field: keyof PrebuiltDictionaryFilters, value: string): void {
    filters = { ...filters, [field]: value }
  }

  function onRowInput(id: string, field: EditableField, value: string): void {
    rows = rows.map((row) => {
      if (row.id !== id) return row
      const next = { ...row, [field]: value }
      const baseline = findSavedRow(id)
      return {
        ...next,
        pending: baseline == null ? "created" : changed(next, baseline) ? "edited" : undefined
      }
    })
  }

  function onDelete(id: string): void {
    rows = rows.map((row) => row.id === id ? { ...row, deletePending: true } : row)
    if (editingRowId === id) editingRowId = ""
  }

  function onCancelRow(id: string): void {
    const baseline = findSavedRow(id)
    if (baseline == null) {
      rows = rows.filter((row) => row.id !== id)
      if (editingRowId === id) editingRowId = ""
      return
    }
    rows = rows.map((row) => {
      if (row.id !== id) return row
      if (row.deletePending) return { ...row, deletePending: undefined }
      return clone([baseline])[0]
    })
    if (editingRowId === id) editingRowId = ""
  }

  function onCreate(): void {
    newRowSequence += 1
    const id = `new-${newRowSequence}`
    rows = [{
      id,
      source: "",
      destination: "",
      partOfSpeech: "noun",
      categories: [],
      pending: "created"
    }, ...rows]
    editingRowId = id
  }

  function onStartEdit(id: string): void {
    if (rows.some((row) => row.id === id && !row.deletePending)) editingRowId = id
  }

  function onCancelChanges(): void {
    rows = clone(saved)
    editingRowId = ""
  }

  async function confirm(): Promise<void> {
    try {
      for (const row of rows) {
        if (row.deletePending) {
          if (row.pending !== "created") await deleteTermDictionary(row.id, revisions[row.id])
          continue
        }
        if (row.pending === "created") {
          await createTermDictionary(row)
          continue
        }
        if (row.pending === "edited") {
          const baseline = findSavedRow(row.id)
          if (baseline != null) await patchTermDictionary(row.id, revisions[row.id], rowChanges(row, baseline))
        }
      }
      await load()
    } catch (error: unknown) {
      window.alert(error instanceof Error ? error.message : "変更を保存できませんでした。")
    }
  }

  function toggleCategories(id: string): void {
    expandedRowIds = expandedRowIds.includes(id)
      ? expandedRowIds.filter((value) => value !== id)
      : [...expandedRowIds, id]
  }

  onMount(() => {
    void load().catch(() => window.alert("用語辞書を読み込めませんでした。"))
  })
</script>

<PrebuiltDictionaryEditorScreen
  {rows}
  {filters}
  {expandedRowIds}
  {hasPendingChanges}
  {editingRowId}
  {pageNumber}
  {totalCount}
  canPrev={pageNumber > 1}
  canNext={pageNumber * PREBUILT_DICTIONARY_PAGE_SIZE < totalCount}
  {onFilterInput}
  onSearch={() => { pageNumber = 1; void load(1) }}
  {onCreate}
  {onDelete}
  {onStartEdit}
  {onCancelRow}
  {onRowInput}
  onToggleCategories={toggleCategories}
  onConfirmChanges={() => { void confirm() }}
  {onCancelChanges}
  onPrev={() => { void load(pageNumber - 1) }}
  onNext={() => { void load(pageNumber + 1) }}
/>
