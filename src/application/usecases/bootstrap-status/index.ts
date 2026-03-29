import type { BootstrapStatusGateway } from "@application/ports/gateway/bootstrap-status";
import type {
  BootstrapStatusField,
  BootstrapStatusScreenInput
} from "@application/ports/input/bootstrap-status";
import type { FeatureScreenStorePort } from "@application/ports/input/feature-screen";
import { createFeatureScreenUsecase } from "@application/usecases/feature-screen";
import type { BootstrapStatus } from "@shared/contracts/bootstrap-status";

type BootstrapStatusStorePort = FeatureScreenStorePort<
  BootstrapStatus,
  BootstrapStatusField,
  undefined
>;

type CreateBootstrapStatusScreenUsecaseOptions = {
  gateway: BootstrapStatusGateway;
  store: BootstrapStatusStorePort;
};

export function createBootstrapStatusScreenUsecase({
  gateway,
  store
}: CreateBootstrapStatusScreenUsecaseOptions): BootstrapStatusScreenInput {
  return createFeatureScreenUsecase({
    createRequest: () => undefined,
    gateway,
    reconcileSelection: ({ currentSelection }) => currentSelection,
    store
  });
}
