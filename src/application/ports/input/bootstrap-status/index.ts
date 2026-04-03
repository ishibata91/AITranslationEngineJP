import type {
  FeatureScreenState,
  FeatureScreenUsecase,
} from "@application/ports/input/feature-screen";
import type { BootstrapStatus } from "@shared/contracts/bootstrap-status";

export type BootstrapStatusField = keyof BootstrapStatus;
export type BootstrapStatusScreenState = FeatureScreenState<
  BootstrapStatus,
  BootstrapStatusField,
  undefined
>;
export type BootstrapStatusScreenInput = FeatureScreenUsecase<
  BootstrapStatusField,
  undefined
>;
