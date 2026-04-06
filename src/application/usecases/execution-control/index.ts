export type ExecutionControlAction = "pause" | "resume" | "retry" | "cancel";

export type ExecutionControlStateValue =
  | "Running"
  | "Paused"
  | "Retrying"
  | "RecoverableFailed"
  | "Failed"
  | "Canceled"
  | "Completed";

type ExecutionControlFailureCategory =
  | "RecoverableProviderFailure"
  | "UnrecoverableProviderFailure"
  | "ValidationFailure"
  | "UserCanceled";

export type ExecutionControlFailure = {
  category: ExecutionControlFailureCategory;
  message: string;
};

type ExecutionControlSnapshot = {
  failure: ExecutionControlFailure | null;
  state: ExecutionControlStateValue;
};

export type ExecutionControlScreenState = {
  canCancel: boolean;
  canPause: boolean;
  canResume: boolean;
  canRetry: boolean;
  controlState: ExecutionControlStateValue;
  error: string | null;
  failure: ExecutionControlFailure | null;
  pendingAction: ExecutionControlAction | null;
};

type ExecutionControlSubscriber = (state: ExecutionControlScreenState) => void;

export interface ExecutionControlScreenStore {
  getState(): ExecutionControlScreenState;
  setState(state: ExecutionControlScreenState): void;
  subscribe(run: ExecutionControlSubscriber): () => void;
}

export interface ExecutionControlScreenInput {
  cancel(): Promise<void>;
  initialize(): Promise<void>;
  pause(): Promise<void>;
  resume(): Promise<void>;
  retry(): Promise<void>;
}

type CreateExecutionControlScreenUsecaseOptions = {
  cancelCommand: () => Promise<ExecutionControlSnapshot>;
  loadSnapshot: () => Promise<ExecutionControlSnapshot>;
  pauseCommand: () => Promise<ExecutionControlSnapshot>;
  resumeCommand: () => Promise<ExecutionControlSnapshot>;
  retryCommand: () => Promise<ExecutionControlSnapshot>;
  store: ExecutionControlScreenStore;
  toErrorMessage?: (error: unknown) => string;
};

function defaultToErrorMessage(): string {
  return "Execution control failed. Try again.";
}

function createInitialState(): ExecutionControlScreenState {
  return {
    canCancel: false,
    canPause: false,
    canResume: false,
    canRetry: false,
    controlState: "Running",
    error: null,
    failure: null,
    pendingAction: null,
  };
}

function deriveAvailability(snapshot: ExecutionControlSnapshot) {
  switch (snapshot.state) {
    case "Running":
      return {
        canCancel: true,
        canPause: true,
        canResume: false,
        canRetry: false,
      };
    case "Paused":
      return {
        canCancel: true,
        canPause: false,
        canResume: true,
        canRetry: false,
      };
    case "RecoverableFailed":
      return {
        canCancel: true,
        canPause: false,
        canResume: false,
        canRetry: true,
      };
    case "Retrying":
      return {
        canCancel: true,
        canPause: false,
        canResume: false,
        canRetry: false,
      };
    case "Failed":
    case "Canceled":
    case "Completed":
      return {
        canCancel: false,
        canPause: false,
        canResume: false,
        canRetry: false,
      };
  }
}

function toViewState(args: {
  error: string | null;
  fallbackFailure: ExecutionControlFailure | null;
  pendingAction: ExecutionControlAction | null;
  snapshot: ExecutionControlSnapshot;
}): ExecutionControlScreenState {
  const availability =
    args.pendingAction === null
      ? deriveAvailability(args.snapshot)
      : {
          canCancel: false,
          canPause: false,
          canResume: false,
          canRetry: false,
        };

  return {
    ...availability,
    controlState: args.snapshot.state,
    error: args.error,
    failure: args.snapshot.failure ?? args.fallbackFailure,
    pendingAction: args.pendingAction,
  };
}

export function createExecutionControlScreenStore(): ExecutionControlScreenStore {
  let state = createInitialState();
  const subscribers = new Set<ExecutionControlSubscriber>();

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

export function createExecutionControlScreenUsecase({
  cancelCommand,
  loadSnapshot,
  pauseCommand,
  resumeCommand,
  retryCommand,
  store,
  toErrorMessage = defaultToErrorMessage,
}: CreateExecutionControlScreenUsecaseOptions): ExecutionControlScreenInput {
  let confirmedSnapshot: ExecutionControlSnapshot | null = null;

  async function initialize(): Promise<void> {
    try {
      const snapshot = await loadSnapshot();
      confirmedSnapshot = snapshot;
      store.setState(
        toViewState({
          error: null,
          fallbackFailure: null,
          pendingAction: null,
          snapshot,
        }),
      );
    } catch (error) {
      store.setState({
        ...store.getState(),
        error: toErrorMessage(error),
        pendingAction: null,
      });
    }
  }

  async function runAction(
    action: ExecutionControlAction,
    command: () => Promise<ExecutionControlSnapshot>,
  ): Promise<void> {
    const currentSnapshot =
      confirmedSnapshot ??
      ({
        failure: null,
        state: store.getState().controlState,
      } satisfies ExecutionControlSnapshot);
    const previousFailure = store.getState().failure;

    store.setState(
      toViewState({
        error: null,
        fallbackFailure: previousFailure,
        pendingAction: action,
        snapshot: currentSnapshot,
      }),
    );

    try {
      const snapshot = await command();
      confirmedSnapshot = snapshot;
      store.setState(
        toViewState({
          error: null,
          fallbackFailure: previousFailure,
          pendingAction: null,
          snapshot,
        }),
      );
    } catch (error) {
      store.setState(
        toViewState({
          error: toErrorMessage(error),
          fallbackFailure: previousFailure,
          pendingAction: null,
          snapshot: currentSnapshot,
        }),
      );
    }
  }

  return {
    cancel() {
      return runAction("cancel", cancelCommand);
    },
    initialize,
    pause() {
      return runAction("pause", pauseCommand);
    },
    resume() {
      return runAction("resume", resumeCommand);
    },
    retry() {
      return runAction("retry", retryCommand);
    },
  };
}
