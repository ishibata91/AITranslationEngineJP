<script lang="ts">
  import type { CreateBodyTranslationPhaseScreenController } from "@application/contract/body-translation-phase"
  import type { CreateMasterDictionaryScreenController } from "@application/contract/master-dictionary"
  import type { CreateMasterPersonaScreenController } from "@application/contract/master-persona"
  import type { CreatePersonaGenerationPhaseScreenController } from "@application/contract/persona-generation-phase"
  import type { CreateProviderSettingsScreenController } from "@application/contract/provider-settings"
  import type { CreateTermTranslationPhaseScreenController } from "@application/contract/term-translation-phase"
  import type { CreateTranslationJobManagementScreenController } from "@application/contract/translation-job-management"
  import type { CreateTranslationOutputArtifactScreenController } from "@application/contract/translation-output-artifact"
  import type { CreateTranslationInputScreenController } from "@application/contract/translation-input"
  import { createShellState } from "@ui/stores/shell-state"
  import AppShell from "@ui/views/AppShell.svelte"

  interface Props {
    createBodyTranslationPhaseScreenController?: CreateBodyTranslationPhaseScreenController | null
    createMasterDictionaryScreenController?: CreateMasterDictionaryScreenController | null
    createMasterPersonaScreenController?: CreateMasterPersonaScreenController | null
    createPersonaGenerationPhaseScreenController?: CreatePersonaGenerationPhaseScreenController | null
    createProviderSettingsScreenController?: CreateProviderSettingsScreenController | null
    createTermTranslationPhaseScreenController?: CreateTermTranslationPhaseScreenController | null
    createTranslationJobManagementScreenController?: CreateTranslationJobManagementScreenController | null
    createTranslationOutputArtifactScreenController?: CreateTranslationOutputArtifactScreenController | null
    createTranslationInputScreenController?: CreateTranslationInputScreenController | null
  }

  let {
    createBodyTranslationPhaseScreenController = null,
    createMasterDictionaryScreenController = null,
    createMasterPersonaScreenController = null,
    createPersonaGenerationPhaseScreenController = null,
    createProviderSettingsScreenController = null,
    createTermTranslationPhaseScreenController = null,
    createTranslationJobManagementScreenController = null,
    createTranslationOutputArtifactScreenController = null,
    createTranslationInputScreenController = null
  }: Props = $props()

  const shellState = createShellState()
</script>

<AppShell
  defaultRouteId={shellState.defaultRouteId}
  defaultTranslationManagementViewId={shellState.defaultTranslationManagementViewId}
  {createBodyTranslationPhaseScreenController}
  {createMasterDictionaryScreenController}
  {createMasterPersonaScreenController}
  {createPersonaGenerationPhaseScreenController}
  {createProviderSettingsScreenController}
  {createTermTranslationPhaseScreenController}
  {createTranslationJobManagementScreenController}
  {createTranslationOutputArtifactScreenController}
  {createTranslationInputScreenController}
  routes={shellState.routes}
/>
