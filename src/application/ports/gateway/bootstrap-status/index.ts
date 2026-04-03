import type { FeatureScreenGateway } from "@application/ports/gateway/feature-screen";
import type { BootstrapStatus } from "@shared/contracts/bootstrap-status";

export type BootstrapStatusGateway = FeatureScreenGateway<
  undefined,
  BootstrapStatus
>;
