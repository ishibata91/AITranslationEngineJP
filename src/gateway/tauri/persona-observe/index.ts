import { invoke } from "@tauri-apps/api/core";
import type { PersonaObserveResult } from "@application/usecases/persona-observe";

type PersonaObserveRequest = {
  personaName: string;
};

export function createTauriPersonaObserveExecutor(): (
  request: PersonaObserveRequest,
) => Promise<PersonaObserveResult> {
  return (request) => {
    return invoke<PersonaObserveResult>("read_master_persona", {
      request,
    });
  };
}
