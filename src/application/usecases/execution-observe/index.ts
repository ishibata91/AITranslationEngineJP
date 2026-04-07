type ExecutionObserveControlStateValue =
  | "Running"
  | "Paused"
  | "Retrying"
  | "RecoverableFailed"
  | "Failed"
  | "Canceled"
  | "Completed";

type ExecutionObserveFailureCategory =
  | "RecoverableProviderFailure"
  | "UnrecoverableProviderFailure"
  | "ValidationFailure"
  | "UserCanceled";

type ExecutionObserveFailure = {
  category: ExecutionObserveFailureCategory;
  message: string;
};

export type ExecutionObserveSnapshot = {
  controlState: ExecutionObserveControlStateValue;
  failure: ExecutionObserveFailure | null;
  footerMetadata: {
    lastEventAt: string;
    manualRecoveryGuidance: string;
    providerRunId: string;
    runHash: string;
  };
  phaseRuns: Array<{
    endedAt: string | null;
    phaseKey: string;
    startedAt: string;
    statusLabel: string;
  }>;
  phaseTimeline: Array<{
    isCurrent: boolean;
    label: string;
    statusLabel: string;
  }>;
  selectedUnit: {
    destText: string;
    formId: string;
    sourceText: string;
    statusLabel: string;
  } | null;
  summary: {
    currentPhase: string;
    jobName: string;
    providerLabel: string;
    startedAt: string;
    statusLabel: string;
  };
  translationProgress: {
    completedUnits: number;
    queuedUnits: number;
    runningUnits: number;
    totalUnits: number;
  };
};

export type ExecutionObserveScreenState = {
  error: string | null;
  loading: boolean;
  snapshot: ExecutionObserveSnapshot | null;
};

type ExecutionObserveSubscriber = (state: ExecutionObserveScreenState) => void;

export interface ExecutionObserveScreenStore {
  getState(): ExecutionObserveScreenState;
  setState(state: ExecutionObserveScreenState): void;
  subscribe(run: ExecutionObserveSubscriber): () => void;
}

export interface ExecutionObserveScreenInput {
  initialize(): Promise<void>;
  refresh(): Promise<void>;
}

type CreateExecutionObserveScreenUsecaseOptions = {
  loadSnapshot: () => Promise<ExecutionObserveSnapshot>;
  store: ExecutionObserveScreenStore;
  toErrorMessage?: (error: unknown) => string;
};

function defaultToErrorMessage(): string {
  return "Execution observation failed. Try again.";
}

function createInitialState(): ExecutionObserveScreenState {
  return {
    error: null,
    loading: false,
    snapshot: null,
  };
}

export function createExecutionObserveScreenStore(): ExecutionObserveScreenStore {
  let state = createInitialState();
  const subscribers = new Set<ExecutionObserveSubscriber>();

  function notify(): void {
    subscribers.forEach((subscriber) => subscriber(state));
  }

  return {
    getState() {
      return state;
    },
    setState(nextState) {
      state = nextState;
      notify();
    },
    subscribe(run) {
      subscribers.add(run);
      run(state);

      return () => {
        subscribers.delete(run);
      };
    },
  };
}

export function createExecutionObserveScreenUsecase({
  loadSnapshot,
  store,
  toErrorMessage = defaultToErrorMessage,
}: CreateExecutionObserveScreenUsecaseOptions): ExecutionObserveScreenInput {
  let confirmedSnapshot = null as ExecutionObserveSnapshot | null;

  async function initialize(): Promise<void> {
    store.setState({
      error: null,
      loading: true,
      snapshot: confirmedSnapshot,
    });

    try {
      const snapshot = await loadSnapshot();
      confirmedSnapshot = snapshot;
      store.setState({
        error: null,
        loading: false,
        snapshot,
      });
    } catch (error) {
      store.setState({
        error: toErrorMessage(error),
        loading: false,
        snapshot: confirmedSnapshot,
      });
    }
  }

  async function refresh(): Promise<void> {
    const currentSnapshot = confirmedSnapshot ?? store.getState().snapshot;

    store.setState({
      error: null,
      loading: true,
      snapshot: currentSnapshot,
    });

    try {
      const snapshot = await loadSnapshot();
      confirmedSnapshot = snapshot;
      store.setState({
        error: null,
        loading: false,
        snapshot,
      });
    } catch (error) {
      store.setState({
        error: toErrorMessage(error),
        loading: false,
        snapshot: currentSnapshot,
      });
    }
  }

  return {
    initialize,
    refresh,
  };
}
