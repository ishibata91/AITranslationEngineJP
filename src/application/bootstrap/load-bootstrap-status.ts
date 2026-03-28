import type { BootstrapStatus } from "@shared/contracts/bootstrap-status";
import type { BootstrapStatusPort } from "./bootstrap-status.port";
import { tauriBootstrapStatusGateway } from "@gateway/tauri/bootstrap-status.gateway";

export async function loadBootstrapStatus(
  port: BootstrapStatusPort = tauriBootstrapStatusGateway
): Promise<BootstrapStatus> {
  return port.getBootstrapStatus();
}
