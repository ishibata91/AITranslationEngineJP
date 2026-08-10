# References: prebuilt-dictionary-editor

| ID | 種類 | 所在 |
| --- | --- | --- |
| REF-1 | source | `dictionary/schema.sql`: `dictionary_term`、`dictionary_sense`、外部キー、FTS trigger |
| REF-2 | source | `dictionary/store.go`: `Store.Get`、`Store.Add`、`Store.Update` |
| REF-3 | source | `internal/store/prebuilt_dictionary.go`: `PrebuiltDictionary`、`OpenPrebuiltDictionary` |
| REF-4 | source | `dictionary/mcp.go`: `newMCPServer` |
| REF-5 | source | `frontend/src/App.svelte`: `NAV_ITEMS`、画面分岐 |
| REF-6 | source | `dictionary/search.go`: `searchFilter`、`Search` |
| REF-7 | source | `frontend/src/ui/components/translation-run/ResultsPager.svelte`: `ResultsPager` |
| REF-8 | source | `frontend/src/ui/components/Field.svelte`: `Field` |
| REF-9 | source | `internal/bootstrap/bootstrap.go`: `NewApp`、`AppCloser` |
| REF-10 | source | `internal/api/app.go`: `App`、`New` |
| REF-11 | source | `main.go`: Wails `Bind` |
| REF-12 | source | `frontend/src/ui/screens/prebuilt-dictionary-editor/PrebuiltDictionaryEditorScreen.svelte`: `PrebuiltDictionaryEditorScreen` |
