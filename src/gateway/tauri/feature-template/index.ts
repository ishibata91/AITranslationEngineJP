import type { FeatureTemplateGateway } from "@application/ports/gateway/feature-template";
import { createTauriFeatureScreenGateway } from "@gateway/tauri/feature-screen";
import type {
  FeatureTemplateData,
  FeatureTemplateQuery,
} from "@shared/contracts/feature-template";

export function createTauriFeatureTemplateGateway(
  commandName = "replace_with_backend_command",
): FeatureTemplateGateway {
  return createTauriFeatureScreenGateway<
    FeatureTemplateQuery,
    FeatureTemplateData
  >(commandName);
}
