import { invoke } from "@tauri-apps/api/core";
import type { BootstrapStatusPort } from "@application/bootstrap/bootstrap-status.port";
import type { BootstrapStatus } from "@shared/contracts/bootstrap-status";

export const tauriBootstrapStatusGateway: BootstrapStatusPort = {
  async getBootstrapStatus(): Promise<BootstrapStatus> {
    return invoke<BootstrapStatus>("get_bootstrap_status");
  }
};
