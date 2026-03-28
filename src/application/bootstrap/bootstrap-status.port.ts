import type { BootstrapStatus } from "@shared/contracts/bootstrap-status";

export interface BootstrapStatusPort {
  getBootstrapStatus(): Promise<BootstrapStatus>;
}
