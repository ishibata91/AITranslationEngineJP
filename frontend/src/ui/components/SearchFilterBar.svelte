<script lang="ts">
  import SelectField from "./SelectField.svelte"
  import TextInputField from "./TextInputField.svelte"

  export interface SearchFilterOption {
    value: string
    label: string
  }

  interface Props {
    searchId: string
    searchLabel: string
    searchValue: string
    filterId?: string
    filterLabel?: string
    filterValue?: string
    filterOptions?: SearchFilterOption[]
    placeholder?: string
    disabled?: boolean
    searchHelp?: string
    filterHelp?: string
    onSearchInput: (value: string) => void
    onFilterChange?: ((value: string) => void) | null
  }

  let {
    searchId,
    searchLabel,
    searchValue,
    filterId = "shared-filter",
    filterLabel = "表示条件",
    filterValue = "",
    filterOptions = [],
    placeholder = "",
    disabled = false,
    searchHelp = "",
    filterHelp = "",
    onSearchInput,
    onFilterChange = null
  }: Props = $props()

  const hasFilter = $derived(
    filterOptions.length > 0 && onFilterChange !== null
  )

  const handleFilterChange = (value: string): void => {
    onFilterChange?.(value)
  }
</script>

<div class="search-filter-bar" role="search">
  <TextInputField
    id={searchId}
    label={searchLabel}
    value={searchValue}
    type="search"
    {placeholder}
    help={searchHelp}
    {disabled}
    autocomplete="off"
    onInput={onSearchInput}
  />
  {#if hasFilter}
    <SelectField
      id={filterId}
      label={filterLabel}
      value={filterValue}
      options={filterOptions}
      help={filterHelp}
      {disabled}
      onChange={handleFilterChange}
    />
  {/if}
</div>

<style>
  .search-filter-bar {
    align-items: end;
    display: grid;
    gap: 0.8rem;
    grid-template-columns: minmax(14rem, 1fr) minmax(12rem, 18rem);
    width: 100%;
  }

  @media (max-width: 720px) {
    .search-filter-bar {
      grid-template-columns: 1fr;
    }
  }
</style>
