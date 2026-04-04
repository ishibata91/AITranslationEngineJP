import { invoke } from "@tauri-apps/api/core";
import type { DictionaryObserveResult } from "@application/usecases/dictionary-observe";

type DictionaryObserveRequest = {
  sourceTexts: string[];
};

export function createTauriDictionaryObserveExecutor(): (
  request: DictionaryObserveRequest,
) => Promise<DictionaryObserveResult> {
  return (request) => {
    return invoke<DictionaryObserveResult>("lookup_dictionary", {
      request,
    });
  };
}
