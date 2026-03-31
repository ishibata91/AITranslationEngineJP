import { describe, expect, it, vi } from "vitest";
import { createJobListScreenStore, createJobListScreenUsecase } from "./index";

type ObservableJobState = "Ready" | "Running" | "Completed";

type JobListResult = {
  jobs: Array<{
    jobId: string;
    state: ObservableJobState;
  }>;
};

function createDeferred<T>() {
  let resolve = undefined as unknown as (value: T) => void;

  const promise = new Promise<T>((nextResolve) => {
    resolve = nextResolve;
  });

  return {
    promise,
    resolve
  };
}

describe("createJobListScreenUsecase", () => {
  it("Given the first observation is pending When initialize runs Then loading starts and the first returned job is selected", async () => {
    const deferred = createDeferred<JobListResult>();
    const executor = vi.fn(() => deferred.promise);
    const store = createJobListScreenStore();
    const usecase = createJobListScreenUsecase({
      executor,
      store
    });

    const initializePromise = usecase.initialize();

    expect(store.getState()).toEqual({
      data: null,
      error: null,
      filters: undefined,
      loading: true,
      selection: null
    });

    deferred.resolve({
      jobs: [
        {
          jobId: "job-101",
          state: "Ready"
        },
        {
          jobId: "job-202",
          state: "Running"
        }
      ]
    });

    await initializePromise;

    expect(executor).toHaveBeenCalledTimes(1);
    expect(store.getState()).toEqual({
      data: {
        jobs: [
          {
            jobId: "job-101",
            state: "Ready"
          },
          {
            jobId: "job-202",
            state: "Running"
          }
        ]
      },
      error: null,
      filters: undefined,
      loading: false,
      selection: "job-101"
    });
  });

  it("Given a selected observable job When refresh returns the same job id Then the current selection is preserved", async () => {
    const executor = vi
      .fn<() => Promise<JobListResult>>()
      .mockResolvedValueOnce({
        jobs: [
          {
            jobId: "job-101",
            state: "Ready"
          },
          {
            jobId: "job-202",
            state: "Running"
          }
        ]
      })
      .mockResolvedValueOnce({
        jobs: [
          {
            jobId: "job-202",
            state: "Completed"
          },
          {
            jobId: "job-303",
            state: "Ready"
          }
        ]
      });
    const store = createJobListScreenStore();
    const usecase = createJobListScreenUsecase({
      executor,
      store
    });

    await usecase.initialize();
    usecase.select("job-202");
    await usecase.refresh();

    expect(store.getState()).toEqual({
      data: {
        jobs: [
          {
            jobId: "job-202",
            state: "Completed"
          },
          {
            jobId: "job-303",
            state: "Ready"
          }
        ]
      },
      error: null,
      filters: undefined,
      loading: false,
      selection: "job-202"
    });
  });

  it("Given the selected job disappears When refresh returns a different list Then selection falls back to the first returned job", async () => {
    const executor = vi
      .fn<() => Promise<JobListResult>>()
      .mockResolvedValueOnce({
        jobs: [
          {
            jobId: "job-101",
            state: "Ready"
          },
          {
            jobId: "job-202",
            state: "Running"
          }
        ]
      })
      .mockResolvedValueOnce({
        jobs: [
          {
            jobId: "job-404",
            state: "Running"
          },
          {
            jobId: "job-505",
            state: "Completed"
          }
        ]
      });
    const store = createJobListScreenStore();
    const usecase = createJobListScreenUsecase({
      executor,
      store
    });

    await usecase.initialize();
    usecase.select("job-202");
    await usecase.refresh();

    expect(store.getState()).toEqual({
      data: {
        jobs: [
          {
            jobId: "job-404",
            state: "Running"
          },
          {
            jobId: "job-505",
            state: "Completed"
          }
        ]
      },
      error: null,
      filters: undefined,
      loading: false,
      selection: "job-404"
    });
  });

  it("Given the backend returns zero observable jobs When refresh succeeds Then the list stays successful and the selected job is cleared", async () => {
    const executor = vi
      .fn<() => Promise<JobListResult>>()
      .mockResolvedValueOnce({
        jobs: [
          {
            jobId: "job-101",
            state: "Ready"
          }
        ]
      })
      .mockResolvedValueOnce({
        jobs: []
      });
    const store = createJobListScreenStore();
    const usecase = createJobListScreenUsecase({
      executor,
      store
    });

    await usecase.initialize();
    await usecase.refresh();

    expect(store.getState()).toEqual({
      data: {
        jobs: []
      },
      error: null,
      filters: undefined,
      loading: false,
      selection: null
    });
  });

  it("Given a loaded list When select runs Then the selected job id changes without another query", async () => {
    const executor = vi.fn<() => Promise<JobListResult>>().mockResolvedValue({
      jobs: [
        {
          jobId: "job-101",
          state: "Ready"
        },
        {
          jobId: "job-202",
          state: "Running"
        }
      ]
    });
    const store = createJobListScreenStore();
    const usecase = createJobListScreenUsecase({
      executor,
      store
    });

    await usecase.initialize();
    usecase.select("job-202");

    expect(executor).toHaveBeenCalledTimes(1);
    expect(store.getState()).toEqual({
      data: {
        jobs: [
          {
            jobId: "job-101",
            state: "Ready"
          },
          {
            jobId: "job-202",
            state: "Running"
          }
        ]
      },
      error: null,
      filters: undefined,
      loading: false,
      selection: "job-202"
    });
  });

  it("Given a user-facing load failure When initialize and retry run Then the generic error is shown and retry can recover", async () => {
    const executor = vi
      .fn<() => Promise<JobListResult>>()
      .mockRejectedValueOnce(new Error("transport timeout at F:/imports/secret.json"))
      .mockResolvedValueOnce({
        jobs: [
          {
            jobId: "job-909",
            state: "Completed"
          }
        ]
      });
    const store = createJobListScreenStore();
    const usecase = createJobListScreenUsecase({
      executor,
      store
    });

    await usecase.initialize();

    expect(store.getState()).toEqual({
      data: null,
      error: "Job list failed to load. Try again.",
      filters: undefined,
      loading: false,
      selection: null
    });

    await usecase.retry();

    expect(executor).toHaveBeenCalledTimes(2);
    expect(store.getState()).toEqual({
      data: {
        jobs: [
          {
            jobId: "job-909",
            state: "Completed"
          }
        ]
      },
      error: null,
      filters: undefined,
      loading: false,
      selection: "job-909"
    });
  });

  it("Given a successful list load When refresh fails Then the prior list and selection stay visible behind the generic error", async () => {
    const executor = vi
      .fn<() => Promise<JobListResult>>()
      .mockResolvedValueOnce({
        jobs: [
          {
            jobId: "job-101",
            state: "Ready"
          },
          {
            jobId: "job-202",
            state: "Running"
          }
        ]
      })
      .mockRejectedValueOnce(new Error("transport timeout at F:/imports/secret.json"));
    const store = createJobListScreenStore();
    const usecase = createJobListScreenUsecase({
      executor,
      store
    });

    await usecase.initialize();
    usecase.select("job-202");
    await usecase.refresh();

    expect(store.getState()).toEqual({
      data: {
        jobs: [
          {
            jobId: "job-101",
            state: "Ready"
          },
          {
            jobId: "job-202",
            state: "Running"
          }
        ]
      },
      error: "Job list failed to load. Try again.",
      filters: undefined,
      loading: false,
      selection: "job-202"
    });
  });
});
