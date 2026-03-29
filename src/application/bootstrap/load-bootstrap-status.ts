import type { BootstrapStatus } from "@shared/contracts/bootstrap-status";
import type { BootstrapStatusPort } from "./bootstrap-status.port";

let configuredBootstrapStatusPort: BootstrapStatusPort | null = null;

export function configureBootstrapStatusPort(port: BootstrapStatusPort | null): void {
  configuredBootstrapStatusPort = port;
}

export async function loadBootstrapStatus(
  port: BootstrapStatusPort | null = configuredBootstrapStatusPort
): Promise<BootstrapStatus> {
  if (port === null) {
    throw new Error("BootstrapStatusPort is not configured.");
  }

  return port.getBootstrapStatus();
}
