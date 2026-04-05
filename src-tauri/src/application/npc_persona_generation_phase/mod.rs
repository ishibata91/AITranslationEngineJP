use async_trait::async_trait;

use crate::application::dto::{
    JobPersonaEntryDto, JobPersonaReadRequestDto, JobPersonaSaveRequestDto,
};
use crate::application::ports::persona_storage::JobPersonaStoragePort;

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct NpcPersonaGenerationPhaseRequestDto {
    pub job_id: String,
    pub source_type: String,
    pub npc_form_id: String,
    pub race: String,
    pub sex: String,
    pub voice: String,
    pub source_text: String,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct NpcPersonaGenerationRequestDto {
    pub npc_form_id: String,
    pub race: String,
    pub sex: String,
    pub voice: String,
    pub source_text: String,
}

#[async_trait]
pub trait NpcPersonaGenerationPort: Send + Sync {
    async fn generate_job_persona(
        &self,
        request: NpcPersonaGenerationRequestDto,
    ) -> Result<Option<JobPersonaEntryDto>, String>;
}

pub struct RunNpcPersonaGenerationPhaseUseCase<S, G>
where
    S: JobPersonaStoragePort,
    G: NpcPersonaGenerationPort,
{
    storage: S,
    generator: G,
}

impl<S, G> RunNpcPersonaGenerationPhaseUseCase<S, G>
where
    S: JobPersonaStoragePort,
    G: NpcPersonaGenerationPort,
{
    pub fn new(storage: S, generator: G) -> Self {
        Self { storage, generator }
    }

    pub async fn execute(
        &self,
        request: NpcPersonaGenerationPhaseRequestDto,
    ) -> Result<Option<JobPersonaEntryDto>, String> {
        if request.job_id.trim().is_empty() {
            return Err("job_id must not be empty".to_string());
        }
        if request.source_type.trim().is_empty() {
            return Err("source_type must not be empty".to_string());
        }
        if request.npc_form_id.trim().is_empty() {
            return Err("npc_form_id must not be empty".to_string());
        }

        let read_result = self
            .storage
            .read_job_persona(JobPersonaReadRequestDto {
                job_id: request.job_id.clone(),
            })
            .await?;
        if let Some(existing) = read_result
            .entries
            .into_iter()
            .find(|entry| entry.npc_form_id == request.npc_form_id)
        {
            return Ok(Some(existing));
        }

        let generated = self
            .generator
            .generate_job_persona(NpcPersonaGenerationRequestDto {
                npc_form_id: request.npc_form_id.clone(),
                race: request.race,
                sex: request.sex,
                voice: request.voice,
                source_text: request.source_text,
            })
            .await?;

        let Some(generated_persona) = generated else {
            return Ok(None);
        };

        self.storage
            .save_job_persona(JobPersonaSaveRequestDto {
                job_id: request.job_id.clone(),
                source_type: request.source_type,
                entries: vec![generated_persona.clone()],
            })
            .await?;

        let saved_result = self
            .storage
            .read_job_persona(JobPersonaReadRequestDto {
                job_id: request.job_id,
            })
            .await?;

        Ok(saved_result
            .entries
            .into_iter()
            .find(|entry| entry.npc_form_id == generated_persona.npc_form_id)
            .or(Some(generated_persona)))
    }
}

#[cfg(test)]
mod tests {
    use std::sync::{Arc, Mutex};

    use async_trait::async_trait;

    use super::*;
    use crate::application::dto::JobPersonaReadResultDto;

    #[derive(Default)]
    struct StubJobPersonaStorageState {
        read_requests: Mutex<Vec<JobPersonaReadRequestDto>>,
        save_requests: Mutex<Vec<JobPersonaSaveRequestDto>>,
        read_results: Mutex<Vec<JobPersonaReadResultDto>>,
    }

    #[derive(Clone)]
    struct StubJobPersonaStorage {
        state: Arc<StubJobPersonaStorageState>,
    }

    impl StubJobPersonaStorage {
        fn new(
            read_results: Vec<JobPersonaReadResultDto>,
        ) -> (Self, Arc<StubJobPersonaStorageState>) {
            let state = Arc::new(StubJobPersonaStorageState {
                read_requests: Mutex::new(vec![]),
                save_requests: Mutex::new(vec![]),
                read_results: Mutex::new(read_results),
            });

            (
                Self {
                    state: Arc::clone(&state),
                },
                state,
            )
        }
    }

    #[async_trait]
    impl JobPersonaStoragePort for StubJobPersonaStorage {
        async fn save_job_persona(&self, request: JobPersonaSaveRequestDto) -> Result<(), String> {
            self.state
                .save_requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .push(request);
            Ok(())
        }

        async fn read_job_persona(
            &self,
            request: JobPersonaReadRequestDto,
        ) -> Result<JobPersonaReadResultDto, String> {
            self.state
                .read_requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .push(request);

            let mut queue = self
                .state
                .read_results
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner());
            if queue.is_empty() {
                return Err("no configured read result".to_string());
            }

            Ok(queue.remove(0))
        }
    }

    #[derive(Clone)]
    struct StubNpcPersonaGenerator {
        requests: Arc<Mutex<Vec<NpcPersonaGenerationRequestDto>>>,
        generated: Option<JobPersonaEntryDto>,
    }

    impl StubNpcPersonaGenerator {
        fn new(
            generated: Option<JobPersonaEntryDto>,
        ) -> (Self, Arc<Mutex<Vec<NpcPersonaGenerationRequestDto>>>) {
            let requests = Arc::new(Mutex::new(vec![]));
            (
                Self {
                    requests: Arc::clone(&requests),
                    generated,
                },
                requests,
            )
        }
    }

    #[async_trait]
    impl NpcPersonaGenerationPort for StubNpcPersonaGenerator {
        async fn generate_job_persona(
            &self,
            request: NpcPersonaGenerationRequestDto,
        ) -> Result<Option<JobPersonaEntryDto>, String> {
            self.requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .push(request);
            Ok(self.generated.clone())
        }
    }

    #[test]
    fn given_cached_job_persona_when_running_npc_persona_generation_phase_then_cache_hit_skips_generation_and_save(
    ) {
        tauri::async_runtime::block_on(async {
            let cached_persona = build_job_persona("00013BA1", "cached persona");
            let (storage, storage_state) =
                StubJobPersonaStorage::new(vec![JobPersonaReadResultDto {
                    job_id: "job-00042".to_string(),
                    entries: vec![cached_persona.clone()],
                }]);
            let (generator, generator_requests) =
                StubNpcPersonaGenerator::new(Some(build_job_persona("00013BA1", "generated")));
            let usecase = RunNpcPersonaGenerationPhaseUseCase::new(storage, generator);

            let result = usecase
                .execute(build_phase_request())
                .await
                .expect("cache-hit path should succeed");

            assert_eq!(result, Some(cached_persona));
            assert!(generator_requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .is_empty());
            assert!(storage_state
                .save_requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .is_empty());
            assert_eq!(
                storage_state
                    .read_requests
                    .lock()
                    .unwrap_or_else(|poisoned| poisoned.into_inner())
                    .len(),
                1
            );
        });
    }

    #[test]
    fn given_missing_cache_and_generated_persona_when_running_npc_persona_generation_phase_then_generated_entry_is_saved(
    ) {
        tauri::async_runtime::block_on(async {
            let generated_persona = build_job_persona("00013BA1", "generated persona");
            let persisted_persona = build_job_persona("00013BA1", "persisted persona");
            let (storage, storage_state) = StubJobPersonaStorage::new(vec![
                JobPersonaReadResultDto {
                    job_id: "job-00042".to_string(),
                    entries: vec![],
                },
                JobPersonaReadResultDto {
                    job_id: "job-00042".to_string(),
                    entries: vec![persisted_persona.clone()],
                },
            ]);
            let (generator, generator_requests) =
                StubNpcPersonaGenerator::new(Some(generated_persona.clone()));
            let usecase = RunNpcPersonaGenerationPhaseUseCase::new(storage, generator);

            let result = usecase
                .execute(build_phase_request())
                .await
                .expect("generate-save path should succeed");

            assert_eq!(result, Some(persisted_persona));

            let generator_requests = generator_requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .clone();
            assert_eq!(
                generator_requests,
                vec![NpcPersonaGenerationRequestDto {
                    npc_form_id: "00013BA1".to_string(),
                    race: "NordRace".to_string(),
                    sex: "Male".to_string(),
                    voice: "MaleNord".to_string(),
                    source_text: "Welcome, <Alias=Player>.".to_string(),
                }]
            );

            let save_requests = storage_state
                .save_requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .clone();
            assert_eq!(
                save_requests,
                vec![JobPersonaSaveRequestDto {
                    job_id: "job-00042".to_string(),
                    source_type: "xedit_export".to_string(),
                    entries: vec![generated_persona],
                }]
            );
            assert_eq!(
                storage_state
                    .read_requests
                    .lock()
                    .unwrap_or_else(|poisoned| poisoned.into_inner())
                    .len(),
                2
            );
        });
    }

    #[test]
    fn given_generator_returns_none_when_running_npc_persona_generation_phase_then_none_is_returned_without_save(
    ) {
        tauri::async_runtime::block_on(async {
            let (storage, storage_state) =
                StubJobPersonaStorage::new(vec![JobPersonaReadResultDto {
                    job_id: "job-00042".to_string(),
                    entries: vec![],
                }]);
            let (generator, generator_requests) = StubNpcPersonaGenerator::new(None);
            let usecase = RunNpcPersonaGenerationPhaseUseCase::new(storage, generator);

            let result = usecase
                .execute(build_phase_request())
                .await
                .expect("none persona path should succeed");

            assert_eq!(result, None);
            assert_eq!(
                generator_requests
                    .lock()
                    .unwrap_or_else(|poisoned| poisoned.into_inner())
                    .len(),
                1
            );
            assert!(storage_state
                .save_requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .is_empty());
            assert_eq!(
                storage_state
                    .read_requests
                    .lock()
                    .unwrap_or_else(|poisoned| poisoned.into_inner())
                    .len(),
                1
            );
        });
    }

    fn build_phase_request() -> NpcPersonaGenerationPhaseRequestDto {
        NpcPersonaGenerationPhaseRequestDto {
            job_id: "job-00042".to_string(),
            source_type: "xedit_export".to_string(),
            npc_form_id: "00013BA1".to_string(),
            race: "NordRace".to_string(),
            sex: "Male".to_string(),
            voice: "MaleNord".to_string(),
            source_text: "Welcome, <Alias=Player>.".to_string(),
        }
    }

    fn build_job_persona(npc_form_id: &str, persona_text: &str) -> JobPersonaEntryDto {
        JobPersonaEntryDto {
            npc_form_id: npc_form_id.to_string(),
            race: "NordRace".to_string(),
            sex: "Male".to_string(),
            voice: "MaleNord".to_string(),
            persona_text: persona_text.to_string(),
        }
    }
}
