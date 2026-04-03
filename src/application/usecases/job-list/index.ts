import type {
  FeatureScreenState,
  FeatureScreenStorePort,
  FeatureScreenUsecase,
} from "@application/ports/input/feature-screen";

export type ObservableJobState = "Ready" | "Running" | "Completed";

export type JobListItem = {
  jobId: string;
  state: ObservableJobState;
};

export type JobListResult = {
  jobs: JobListItem[];
};

export type JobListScreenState = FeatureScreenState<
  JobListResult,
  string,
  undefined
>;

type JobListSubscriber = (state: JobListScreenState) => void;

export interface JobListScreenStore extends FeatureScreenStorePort<
  JobListResult,
  string,
  undefined
> {
  subscribe(run: JobListSubscriber): () => void;
}

export interface JobListScreenInput extends FeatureScreenUsecase<
  string,
  undefined
> {
  initialize(): Promise<void>;
  refresh(): Promise<void>;
  retry(): Promise<void>;
  select(selection: string | null): void;
}

type CreateJobListScreenUsecaseOptions = {
  executor: () => Promise<JobListResult>;
  store: JobListScreenStore;
  toErrorMessage?: (error: unknown) => string;
};

function defaultToErrorMessage(): string {
  return "Job list failed to load. Try again.";
}

function reconcileSelection(args: {
  currentSelection: string | null;
  data: JobListResult;
}): string | null {
  if (args.data.jobs.length === 0) {
    return null;
  }

  if (
    args.currentSelection !== null &&
    args.data.jobs.some((job) => job.jobId === args.currentSelection)
  ) {
    return args.currentSelection;
  }

  return args.data.jobs[0].jobId;
}

export function createJobListScreenStore(): JobListScreenStore {
  let state: JobListScreenState = {
    data: null,
    error: null,
    filters: undefined,
    loading: false,
    selection: null,
  };
  const subscribers = new Set<JobListSubscriber>();

  function notify(): void {
    subscribers.forEach((subscriber) => subscriber(state));
  }

  return {
    subscribe(run) {
      subscribers.add(run);
      run(state);

      return () => {
        subscribers.delete(run);
      };
    },
    getState() {
      return state;
    },
    setError(message) {
      state = {
        ...state,
        error: message,
        loading: false,
      };
      notify();
    },
    setFilters(filters) {
      state = {
        ...state,
        filters,
      };
      notify();
    },
    setLoaded(payload) {
      state = {
        ...state,
        data: payload.data,
        error: null,
        loading: false,
        selection: payload.selection,
      };
      notify();
    },
    setLoading() {
      state = {
        ...state,
        error: null,
        loading: true,
      };
      notify();
    },
    setSelection(selection) {
      state = {
        ...state,
        selection,
      };
      notify();
    },
  };
}

export function createJobListScreenUsecase({
  executor,
  store,
  toErrorMessage = defaultToErrorMessage,
}: CreateJobListScreenUsecaseOptions): JobListScreenInput {
  async function loadCurrent(): Promise<void> {
    store.setLoading();

    try {
      const data = await executor();
      const currentState = store.getState();
      const selection = reconcileSelection({
        currentSelection: currentState.selection,
        data,
      });

      store.setLoaded({
        data,
        selection,
      });
    } catch (error) {
      store.setError(toErrorMessage(error));
    }
  }

  return {
    initialize() {
      return loadCurrent();
    },
    refresh() {
      return loadCurrent();
    },
    retry() {
      return loadCurrent();
    },
    select(selection) {
      store.setSelection(selection);
    },
    async updateFilters(filters, options) {
      store.setFilters(filters);

      if (options?.reload === true) {
        await loadCurrent();
      }
    },
  };
}
