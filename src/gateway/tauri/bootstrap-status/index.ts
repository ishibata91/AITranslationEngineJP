import type { BootstrapStatusGateway } from "@application/ports/gateway/bootstrap-status";
import { createTauriFeatureScreenGateway } from "@gateway/tauri/feature-screen";
import type { BootstrapStatus } from "@shared/contracts/bootstrap-status";

export function createTauriBootstrapStatusGateway(): BootstrapStatusGateway {
  return createTauriFeatureScreenGateway<undefined, BootstrapStatus>("get_bootstrap_status");
}
