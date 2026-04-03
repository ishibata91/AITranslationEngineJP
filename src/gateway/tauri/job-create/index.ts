import { invoke } from "@tauri-apps/api/core";
import type {
  JobCreateRequest,
  JobCreateResult,
} from "@application/usecases/job-create";

export function createTauriJobCreateExecutor(): (
  request: JobCreateRequest,
) => Promise<JobCreateResult> {
  return (request) => {
    return invoke<JobCreateResult>("create_job", {
      request,
    });
  };
}
