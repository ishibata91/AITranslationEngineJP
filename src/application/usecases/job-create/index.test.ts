import { describe, expect, it, vi } from "vitest";
import {
  createDefaultJobCreateRequest,
  createJobCreateScreenStore,
  createJobCreateScreenUsecase,
} from "./index";

function createDeferred<T>() {
  let resolve = undefined as unknown as (value: T) => void;

  const promise = new Promise<T>((nextResolve) => {
    resolve = nextResolve;
  });

  return {
    promise,
    resolve,
  };
}

describe("createJobCreateScreenUsecase", () => {
  it("Given a valid request When submit succeeds Then result state is stored", async () => {
    const request = createDefaultJobCreateRequest();
    const store = createJobCreateScreenStore(request);
    const deferred = createDeferred<{
      jobId: string;
      state: "Ready";
    }>();
    const executor = vi.fn(() => deferred.promise);
    const usecase = createJobCreateScreenUsecase({
      executor,
      store,
    });

    const submitPromise = usecase.submit();

    expect(store.getState()).toEqual({
      error: null,
      isSubmitting: true,
      request,
      result: null,
    });

    deferred.resolve({
      jobId: "job-77",
      state: "Ready",
    });

    await submitPromise;

    expect(executor).toHaveBeenCalledWith(request);
    expect(store.getState()).toEqual({
      error: null,
      isSubmitting: false,
      request,
      result: {
        jobId: "job-77",
        state: "Ready",
      },
    });
  });

  it("Given a blank required field When submit runs Then validation prevents executor call", async () => {
    const store = createJobCreateScreenStore(createDefaultJobCreateRequest());
    const executor = vi.fn();
    const usecase = createJobCreateScreenUsecase({
      executor,
      store,
    });

    usecase.updateSourceGroupField(0, "sourceJsonPath", "");
    await usecase.submit();

    expect(executor).not.toHaveBeenCalled();
    expect(store.getState().error).toBe(
      "Fill in all required fields before creating a job.",
    );
    expect(store.getState().result).toBeNull();
  });

  it("Given executor failure When submit runs Then error is exposed and request is preserved", async () => {
    const request = createDefaultJobCreateRequest();
    const store = createJobCreateScreenStore(request);
    const usecase = createJobCreateScreenUsecase({
      executor: async () => {
        throw new Error("backend create failed");
      },
      store,
    });

    await usecase.submit();

    expect(store.getState()).toEqual({
      error: "Job creation failed. Try again.",
      isSubmitting: false,
      request,
      result: null,
    });
  });

  it("Given a pending create request When submit runs again Then the executor is not called twice", async () => {
    const deferred = createDeferred<{
      jobId: string;
      state: "Ready";
    }>();
    const executor = vi.fn(() => deferred.promise);
    const store = createJobCreateScreenStore(createDefaultJobCreateRequest());
    const usecase = createJobCreateScreenUsecase({
      executor,
      store,
    });

    const firstSubmit = usecase.submit();
    const secondSubmit = usecase.submit();

    expect(executor).toHaveBeenCalledTimes(1);

    deferred.resolve({
      jobId: "job-88",
      state: "Ready",
    });

    await Promise.all([firstSubmit, secondSubmit]);
  });
});
