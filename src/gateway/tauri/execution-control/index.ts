import { invoke } from "@tauri-apps/api/core";
import type { ExecutionControlScreenInput } from "@application/usecases/execution-control";

type ExecutionControlSnapshot =
  Awaited<ReturnType<ExecutionControlScreenInput["initialize"]>> extends never
    ? never
    : {
        failure: {
          category:
            | "RecoverableProviderFailure"
            | "UnrecoverableProviderFailure"
            | "ValidationFailure"
            | "UserCanceled";
          message: string;
        } | null;
        state:
          | "Running"
          | "Paused"
          | "Retrying"
          | "RecoverableFailed"
          | "Failed"
          | "Canceled"
          | "Completed";
      };

export function createTauriExecutionControlGateway() {
  return {
    cancelCommand: () => invoke<ExecutionControlSnapshot>("cancel_execution"),
    loadSnapshot: () =>
      invoke<ExecutionControlSnapshot>("get_execution_control_snapshot"),
    pauseCommand: () => invoke<ExecutionControlSnapshot>("pause_execution"),
    resumeCommand: () => invoke<ExecutionControlSnapshot>("resume_execution"),
    retryCommand: () => invoke<ExecutionControlSnapshot>("retry_execution"),
  };
}
