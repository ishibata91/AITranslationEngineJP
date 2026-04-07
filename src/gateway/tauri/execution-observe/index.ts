import { invoke } from "@tauri-apps/api/core";
import type { ExecutionObserveSnapshot } from "@application/usecases/execution-observe";

export function createTauriExecutionObserveLoader(): () => Promise<ExecutionObserveSnapshot> {
  return () => {
    return invoke<ExecutionObserveSnapshot>("get_execution_observe_snapshot");
  };
}
