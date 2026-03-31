import { invoke } from "@tauri-apps/api/core";
import type { JobListResult } from "@application/usecases/job-list";

export function createTauriJobListExecutor(): () => Promise<JobListResult> {
  return () => {
    return invoke<JobListResult>("list_jobs");
  };
}
