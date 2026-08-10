# References: hyphenated-npc-name-derivation

| ID | 種類 | 所在 |
| --- | --- | --- |
| REF-1 | source | `internal/core/termderive/termderive.go`: `DeriveTerms`、`deriveTwo`、`safePair` |
| REF-2 | source | `internal/core/termderive/termderive_test.go`: `TestDeriveTerms` |
| REF-3 | source | `internal/engine/engine.go`: `Engine.bodyReferences`、`appendPrebuiltNPCDerivedReferences` |
| REF-4 | source | `internal/engine/proper_noun.go`: `Engine.deriveRunProperNouns` |
| REF-5 | source | `internal/engine/batch.go`: `BatchRunner.planBodyRequests` |
| REF-6 | source | `internal/engine/engine.go`: `Engine.DeriveMasterTerms` |
| REF-7 | source | `internal/core/termderive/termderive_test.go`: `TestDeriveTermsHyphenatedNameParts` |
| REF-8 | source | `internal/engine/engine_test.go`: `TestBodyReferencesDerivesPrebuiltNPCNameParts`、`TestPlanBodyRequestsUsesPrebuiltNPCDerivedReferences`、`TestDeriveMasterTermsKeepsHyphenatedNamePartsDisabled` |
| REF-9 | source | `internal/engine/proper_noun_test.go`: `TestDeriveRunProperNounsDerivesHyphenatedNameParts` |
| REF-10 | source | `internal/api/process_other.go`、`internal/api/process_other_test.go`: `hideChildProcessWindow` |
