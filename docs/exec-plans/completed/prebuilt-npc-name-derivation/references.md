# References: prebuilt-npc-name-derivation

`references.md` は source の path・symbol と外部資料の所在だけを持つ。

| ID | 種類 | 所在 |
| --- | --- | --- |
| REF-1 | source | `internal/store/prebuilt_dictionary.go`: `PrebuiltDictionary.References` |
| REF-2 | source | `internal/model/translation_reference.go`: `PrebuiltDictionaryReference` |
| REF-3 | source | `internal/engine/engine.go`: `Engine.bodyReferences`、`toTranslationReferences`、`referencesForSource` |
| REF-4 | source | `internal/core/dictionary/dictionary.go`: `Dictionary.Extract` |
| REF-5 | source | `internal/core/termderive/termderive.go`: `DeriveTerms`、`deriveByname`、`deriveTwo`、`safePair` |
| REF-6 | source | `internal/engine/engine.go`: `Engine.dialogueUsage` |
| REF-7 | source | `internal/engine/batch.go`: `BatchRunner.planBodyRequests` |
| REF-8 | source | `internal/engine/engine_test.go`: `TestBodyReferencesUsesPrebuiltAndTargetPluginProperNouns` |
| REF-9 | source | `internal/store/prebuilt_dictionary_test.go`: `TestPrebuiltDictionaryReferencesAreReadOnlyAndIncludeAllSenses` |
| REF-10 | source | `internal/engine/engine_test.go`: `TestBodyReferencesDerivesPrebuiltNPCNameParts` |
| REF-11 | source | `internal/engine/engine_test.go`: `TestPlanBodyRequestsUsesPrebuiltNPCDerivedReferences` |
| REF-12 | source | `internal/core/termderive/termderive.go`: `NamePair`、`DerivedTerm`、`DeriveTerms` |
