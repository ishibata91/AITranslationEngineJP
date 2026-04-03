type JobState = "Draft" | "Ready" | "Running" | "Completed";

type JobCreateTranslationUnitRequest = {
  sourceEntityType: string;
  formId: string;
  editorId: string;
  recordSignature: string;
  fieldName: string;
  extractionKey: string;
  sourceText: string;
  sortKey: string;
};

type JobCreateSourceGroupRequest = {
  sourceJsonPath: string;
  targetPlugin: string;
  translationUnits: JobCreateTranslationUnitRequest[];
};

export type JobCreateRequest = {
  sourceGroups: JobCreateSourceGroupRequest[];
};

export type JobCreateResult = {
  jobId: string;
  state: JobState;
};

export type JobCreateScreenState = {
  error: string | null;
  isSubmitting: boolean;
  request: JobCreateRequest;
  result: JobCreateResult | null;
};

type JobCreateSubscriber = (state: JobCreateScreenState) => void;
type JobCreateSourceGroupField = "sourceJsonPath" | "targetPlugin";
type JobCreateTranslationUnitField = keyof JobCreateTranslationUnitRequest;

export interface JobCreateScreenStore {
  getState(): JobCreateScreenState;
  replaceState(nextState: JobCreateScreenState): void;
  subscribe(run: JobCreateSubscriber): () => void;
}

export interface JobCreateScreenInput {
  initialize(): Promise<void>;
  resetResult(): void;
  submit(): Promise<void>;
  updateSourceGroupField(
    groupIndex: number,
    field: JobCreateSourceGroupField,
    value: string,
  ): void;
  updateTranslationUnitField(
    groupIndex: number,
    unitIndex: number,
    field: JobCreateTranslationUnitField,
    value: string,
  ): void;
}

type CreateJobCreateScreenUsecaseOptions = {
  executor: (request: JobCreateRequest) => Promise<JobCreateResult>;
  store: JobCreateScreenStore;
  toErrorMessage?: (error: unknown) => string;
};

function defaultToErrorMessage(): string {
  return "Job creation failed. Try again.";
}

function createInitialState(request: JobCreateRequest): JobCreateScreenState {
  return {
    error: null,
    isSubmitting: false,
    request: cloneRequest(request),
    result: null,
  };
}

export function createDefaultJobCreateRequest(): JobCreateRequest {
  return {
    sourceGroups: [
      {
        sourceJsonPath: "F:/imports/xedit-export-minimal.json",
        targetPlugin: "ExampleMod.esp",
        translationUnits: [
          {
            sourceEntityType: "item",
            formId: "00012345",
            editorId: "ExampleSword",
            recordSignature: "WEAP",
            fieldName: "name",
            extractionKey: "item:00012345:name",
            sourceText: "Iron Sword",
            sortKey: "item:00012345:name",
          },
        ],
      },
    ],
  };
}

export function createJobCreateScreenStore(
  request = createDefaultJobCreateRequest(),
): JobCreateScreenStore {
  let state = createInitialState(request);
  const subscribers = new Set<JobCreateSubscriber>();

  function notify(): void {
    subscribers.forEach((subscriber) => subscriber(state));
  }

  return {
    getState() {
      return state;
    },
    replaceState(nextState) {
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

export function createJobCreateScreenUsecase({
  executor,
  store,
  toErrorMessage = defaultToErrorMessage,
}: CreateJobCreateScreenUsecaseOptions): JobCreateScreenInput {
  function updateRequest(mutator: (request: JobCreateRequest) => void): void {
    const currentState = store.getState();
    const nextRequest = cloneRequest(currentState.request);

    mutator(nextRequest);

    store.replaceState({
      ...currentState,
      error: null,
      request: nextRequest,
      result: null,
    });
  }

  return {
    async initialize() {},
    resetResult() {
      const currentState = store.getState();

      store.replaceState({
        ...currentState,
        result: null,
      });
    },
    async submit() {
      const currentState = store.getState();

      if (currentState.isSubmitting) {
        return;
      }

      const request = cloneRequest(currentState.request);
      const validationError = validateRequest(request);

      if (validationError !== null) {
        store.replaceState({
          ...currentState,
          error: validationError,
          result: null,
        });

        return;
      }

      store.replaceState({
        ...currentState,
        error: null,
        isSubmitting: true,
        result: null,
      });

      try {
        const result = await executor(request);
        const nextState = store.getState();

        store.replaceState({
          ...nextState,
          error: null,
          isSubmitting: false,
          result,
        });
      } catch (error) {
        const nextState = store.getState();

        store.replaceState({
          ...nextState,
          error: toErrorMessage(error),
          isSubmitting: false,
          result: null,
        });
      }
    },
    updateSourceGroupField(groupIndex, field, value) {
      updateRequest((request) => {
        request.sourceGroups[groupIndex][field] = value;
      });
    },
    updateTranslationUnitField(groupIndex, unitIndex, field, value) {
      updateRequest((request) => {
        request.sourceGroups[groupIndex].translationUnits[unitIndex][field] =
          value;
      });
    },
  };
}

function cloneRequest(request: JobCreateRequest): JobCreateRequest {
  return {
    sourceGroups: request.sourceGroups.map((sourceGroup) => ({
      sourceJsonPath: sourceGroup.sourceJsonPath,
      targetPlugin: sourceGroup.targetPlugin,
      translationUnits: sourceGroup.translationUnits.map((translationUnit) => ({
        editorId: translationUnit.editorId,
        extractionKey: translationUnit.extractionKey,
        fieldName: translationUnit.fieldName,
        formId: translationUnit.formId,
        recordSignature: translationUnit.recordSignature,
        sortKey: translationUnit.sortKey,
        sourceEntityType: translationUnit.sourceEntityType,
        sourceText: translationUnit.sourceText,
      })),
    })),
  };
}

function validateRequest(request: JobCreateRequest): string | null {
  const blankFieldMessage =
    "Fill in all required fields before creating a job.";

  if (request.sourceGroups.length === 0) {
    return blankFieldMessage;
  }

  for (const sourceGroup of request.sourceGroups) {
    if (isSourceGroupInvalid(sourceGroup)) {
      return blankFieldMessage;
    }
  }

  return null;
}

function isSourceGroupInvalid(
  sourceGroup: JobCreateSourceGroupRequest,
): boolean {
  if (sourceGroup.sourceJsonPath.trim().length === 0) {
    return true;
  }

  if (sourceGroup.targetPlugin.trim().length === 0) {
    return true;
  }

  if (sourceGroup.translationUnits.length === 0) {
    return true;
  }

  return sourceGroup.translationUnits.some((translationUnit) =>
    hasBlankRequiredTranslationField(translationUnit),
  );
}

function hasBlankRequiredTranslationField(
  translationUnit: JobCreateTranslationUnitRequest,
): boolean {
  const requiredFields: Record<JobCreateTranslationUnitField, string> = {
    editorId: translationUnit.editorId,
    extractionKey: translationUnit.extractionKey,
    fieldName: translationUnit.fieldName,
    formId: translationUnit.formId,
    recordSignature: translationUnit.recordSignature,
    sortKey: translationUnit.sortKey,
    sourceEntityType: translationUnit.sourceEntityType,
    sourceText: translationUnit.sourceText,
  };

  return Object.values(requiredFields).some(
    (fieldValue) => fieldValue.trim().length === 0,
  );
}
