<script lang="ts">
  import PrebuiltDictionaryEditorScreen from "./PrebuiltDictionaryEditorScreen.svelte"
  import { emptyPrebuiltDictionaryFilters, prebuiltDictionaryRows } from "./prebuilt-dictionary-editor.fixtures"
  import type { PrebuiltDictionaryFilters, PrebuiltDictionaryRow } from "./prebuilt-dictionary-editor-view"

  function copyRows(source: PrebuiltDictionaryRow[]): PrebuiltDictionaryRow[] {
    return source.map((row) => ({ ...row, categories: [...row.categories] }))
  }

  let originalRows = $state<PrebuiltDictionaryRow[]>(copyRows(prebuiltDictionaryRows))
  let rows = $state<PrebuiltDictionaryRow[]>(copyRows(prebuiltDictionaryRows))
  let filters = $state<PrebuiltDictionaryFilters>({ ...emptyPrebuiltDictionaryFilters })
  let appliedFilters = $state<PrebuiltDictionaryFilters>({ ...emptyPrebuiltDictionaryFilters })
  let expandedRowIds = $state<string[]>([])
  let editingRowId = $state("")
  let pageNumber = $state(1)
  let hasPendingChanges = $derived(rows.some((row) => row.pending != null || row.deletePending))
  let filteredRows = $derived(rows.filter((row) =>
    row.source.toLocaleLowerCase().includes(appliedFilters.source.toLocaleLowerCase()) &&
    row.destination.includes(appliedFilters.destination) &&
    row.partOfSpeech.includes(appliedFilters.partOfSpeech) &&
    row.categories.some((category) => category.includes(appliedFilters.category))
  ))

  function updateRow(id: string, field: "source" | "destination" | "partOfSpeech", value: string): void {
    rows = rows.map((row) => {
      if (row.id !== id) return row
      const original = originalRows.find((candidate) => candidate.id === id)
      const next = { ...row, [field]: value }
      if (original == null) return { ...next, pending: "created" }
      const changed = next.source !== original.source || next.destination !== original.destination || next.partOfSpeech !== original.partOfSpeech
      return { ...next, pending: changed ? "edited" : undefined }
    })
  }
  function toggleDelete(id: string): void {
    rows = rows.map((row) => row.id === id ? { ...row, deletePending: true } : row)
    if (editingRowId === id) editingRowId = ""
  }
  function cancelRow(id: string): void {
    const original = originalRows.find((row) => row.id === id)
    if (original == null) {
      rows = rows.filter((row) => row.id !== id)
      if (editingRowId === id) editingRowId = ""
      return
    }
    rows = rows.map((row) => {
      if (row.id !== id) return row
      if (row.deletePending) return { ...row, deletePending: undefined }
      return { ...original }
    })
    if (editingRowId === id) editingRowId = ""
  }
  function toggleCategories(id: string): void { expandedRowIds = expandedRowIds.includes(id) ? expandedRowIds.filter((rowId) => rowId !== id) : [...expandedRowIds, id] }
  function cancelChanges(): void {
    rows = copyRows(originalRows)
    editingRowId = ""
  }
  function confirmChanges(): void {
    rows = rows.filter((row) => !row.deletePending).map((row) => ({ ...row, pending: undefined, deletePending: undefined }))
    originalRows = copyRows(rows)
    editingRowId = ""
  }
  function createRow(): void {
    const id = `new-${rows.length + 1}`
    rows = [{ id, source: "", destination: "", partOfSpeech: "noun", categories: [], pending: "created" }, ...rows]
    editingRowId = id
  }
  function updateFilter(field: keyof PrebuiltDictionaryFilters, value: string): void {
    if (field === "source") filters = { ...filters, source: value }
    if (field === "destination") filters = { ...filters, destination: value }
    if (field === "partOfSpeech") filters = { ...filters, partOfSpeech: value }
    if (field === "category") filters = { ...filters, category: value }
  }
  function search(): void {
    appliedFilters = { ...filters }
    editingRowId = ""
  }
  function startEdit(id: string): void { editingRowId = id }
</script>

<PrebuiltDictionaryEditorScreen
  rows={filteredRows} {filters} {expandedRowIds} {hasPendingChanges} {editingRowId}
  {pageNumber} totalCount={15917} canPrev={pageNumber > 1} canNext={true}
  onFilterInput={updateFilter}
  onSearch={search}
  onCreate={createRow}
  onDelete={toggleDelete}
  onStartEdit={startEdit}
  onCancelRow={cancelRow}
  onRowInput={updateRow}
  onToggleCategories={toggleCategories}
  onConfirmChanges={confirmChanges}
  onCancelChanges={cancelChanges}
  onPrev={() => { pageNumber -= 1 }}
  onNext={() => { pageNumber += 1 }}
/>
