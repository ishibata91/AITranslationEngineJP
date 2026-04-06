use async_trait::async_trait;
use serde::{Deserialize, Serialize};

use crate::application::dto::{
    ExecutionControlFailureCategoryDto, ExecutionControlFailureDto, ProviderExecutionModeDto,
    ProviderSelectionDto,
};
use crate::application::ports::provider_runtime::{ProviderRuntimeFailure, ProviderRuntimePort};

const GEMINI_PROVIDER_ID: &str = "gemini";
const GEMINI_API_KEY_ENV: &str = "GEMINI_API_KEY";
const GEMINI_BASE_URL: &str = "https://generativelanguage.googleapis.com";
const DEFAULT_GEMINI_MODEL_PATH: &str = "models/gemini-1.5-flash";

#[derive(Debug, Clone)]
struct GeminiConfig {
    endpoint: String,
    api_key: Option<String>,
}

impl Default for GeminiConfig {
    fn default() -> Self {
        let endpoint = format!(
            "{}/v1beta/{}:generateContent",
            GEMINI_BASE_URL, DEFAULT_GEMINI_MODEL_PATH
        );

        let api_key = std::env::var(GEMINI_API_KEY_ENV)
            .ok()
            .map(|value| value.trim().to_string())
            .filter(|value| !value.is_empty());

        Self { endpoint, api_key }
    }
}

#[derive(Debug, Clone, Serialize)]
struct GeminiRequest {
    contents: Vec<GeminiContent>,
}

#[derive(Debug, Clone, Serialize)]
struct GeminiContent {
    parts: Vec<GeminiPart>,
}

#[derive(Debug, Clone, Serialize)]
struct GeminiPart {
    text: String,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct GeminiResponse {
    #[serde(default)]
    candidates: Vec<GeminiCandidate>,
    prompt_feedback: Option<GeminiPromptFeedback>,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct GeminiPromptFeedback {
    block_reason: Option<String>,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct GeminiCandidate {
    content: Option<GeminiCandidateContent>,
}

#[derive(Debug, Deserialize)]
struct GeminiCandidateContent {
    #[serde(default)]
    parts: Vec<GeminiCandidatePart>,
}

#[derive(Debug, Deserialize)]
struct GeminiCandidatePart {
    text: Option<String>,
}

#[derive(Debug, Clone, Copy)]
enum GeminiTransportError {
    Connection,
    Api { status: u16 },
    Response,
}

#[async_trait]
trait GeminiTransport: Send + Sync {
    async fn execute(
        &self,
        endpoint: &str,
        request: &GeminiRequest,
    ) -> Result<String, GeminiTransportError>;
}

struct ReqwestGeminiTransport {
    client: reqwest::Client,
    api_key: String,
}

impl ReqwestGeminiTransport {
    fn new(api_key: String) -> Self {
        Self {
            client: reqwest::Client::new(),
            api_key,
        }
    }
}

#[async_trait]
impl GeminiTransport for ReqwestGeminiTransport {
    async fn execute(
        &self,
        endpoint: &str,
        request: &GeminiRequest,
    ) -> Result<String, GeminiTransportError> {
        let response = self
            .client
            .post(endpoint)
            .header("x-goog-api-key", &self.api_key)
            .json(request)
            .send()
            .await
            .map_err(|_| GeminiTransportError::Connection)?;

        let status = response.status();
        if !status.is_success() {
            return Err(GeminiTransportError::Api {
                status: status.as_u16(),
            });
        }

        response
            .text()
            .await
            .map_err(|_| GeminiTransportError::Response)
    }
}

pub struct GeminiProviderRuntimeAdapter {
    config: GeminiConfig,
    transport: Box<dyn GeminiTransport>,
}

impl GeminiProviderRuntimeAdapter {
    pub fn new() -> Self {
        let config = GeminiConfig::default();
        let transport: Box<dyn GeminiTransport> = match config.api_key.clone() {
            Some(api_key) => Box::new(ReqwestGeminiTransport::new(api_key)),
            None => Box::new(NoopGeminiTransport),
        };

        Self { config, transport }
    }

    fn map_request(&self) -> GeminiRequest {
        GeminiRequest {
            contents: vec![GeminiContent {
                parts: vec![GeminiPart {
                    text: "Run one provider step.".to_string(),
                }],
            }],
        }
    }

    fn validate_selection_and_config(
        &self,
        selection: &ProviderSelectionDto,
    ) -> Result<(), ProviderRuntimeFailure> {
        if selection.provider_id != GEMINI_PROVIDER_ID {
            return Err(ExecutionControlFailureDto {
                category: ExecutionControlFailureCategoryDto::ValidationFailure,
                message: "Gemini adapter cannot run the selected provider.".to_string(),
            });
        }

        if selection.execution_mode != ProviderExecutionModeDto::Batch {
            return Err(ExecutionControlFailureDto {
                category: ExecutionControlFailureCategoryDto::ValidationFailure,
                message: "Gemini adapter supports only batch execution mode.".to_string(),
            });
        }

        if self.config.api_key.is_none() {
            return Err(ExecutionControlFailureDto {
                category: ExecutionControlFailureCategoryDto::ValidationFailure,
                message: "Gemini API credentials are not configured.".to_string(),
            });
        }

        Ok(())
    }

    fn normalize_transport_failure(error: GeminiTransportError) -> ProviderRuntimeFailure {
        match error {
            GeminiTransportError::Connection => ExecutionControlFailureDto {
                category: ExecutionControlFailureCategoryDto::RecoverableProviderFailure,
                message: "Gemini provider connection failed. Retry may recover.".to_string(),
            },
            GeminiTransportError::Api {
                status: 429 | 500 | 503 | 504,
            } => ExecutionControlFailureDto {
                category: ExecutionControlFailureCategoryDto::RecoverableProviderFailure,
                message: "Gemini provider is temporarily unavailable. Retry may recover."
                    .to_string(),
            },
            GeminiTransportError::Api { .. } | GeminiTransportError::Response => {
                ExecutionControlFailureDto {
                    category: ExecutionControlFailureCategoryDto::UnrecoverableProviderFailure,
                    message: "Gemini provider returned an unusable response.".to_string(),
                }
            }
        }
    }

    fn parse_response(body: &str) -> Result<(), ProviderRuntimeFailure> {
        let response = serde_json::from_str::<GeminiResponse>(body)
            .map_err(|_| Self::normalize_transport_failure(GeminiTransportError::Response))?;

        if response
            .prompt_feedback
            .as_ref()
            .and_then(|feedback| feedback.block_reason.as_deref())
            .is_some()
        {
            return Err(Self::normalize_transport_failure(
                GeminiTransportError::Response,
            ));
        }

        let has_usable_candidate = response.candidates.iter().any(|candidate| {
            candidate
                .content
                .as_ref()
                .map(|content| {
                    content.parts.iter().any(|part| {
                        part.text
                            .as_deref()
                            .map(|text| !text.trim().is_empty())
                            .unwrap_or(false)
                    })
                })
                .unwrap_or(false)
        });

        if !has_usable_candidate {
            return Err(Self::normalize_transport_failure(
                GeminiTransportError::Response,
            ));
        }

        Ok(())
    }

    #[cfg(test)]
    fn with_transport(transport: Box<dyn GeminiTransport>) -> Self {
        Self {
            config: GeminiConfig::default(),
            transport,
        }
    }
}

impl Default for GeminiProviderRuntimeAdapter {
    fn default() -> Self {
        Self::new()
    }
}

#[async_trait]
impl ProviderRuntimePort for GeminiProviderRuntimeAdapter {
    async fn run_provider_step(
        &self,
        selection: ProviderSelectionDto,
    ) -> Result<(), ProviderRuntimeFailure> {
        self.validate_selection_and_config(&selection)?;

        let request = self.map_request();
        let body = self
            .transport
            .execute(&self.config.endpoint, &request)
            .await
            .map_err(Self::normalize_transport_failure)?;

        Self::parse_response(&body)
    }
}

struct NoopGeminiTransport;

#[async_trait]
impl GeminiTransport for NoopGeminiTransport {
    async fn execute(
        &self,
        _endpoint: &str,
        _request: &GeminiRequest,
    ) -> Result<String, GeminiTransportError> {
        Err(GeminiTransportError::Connection)
    }
}

#[cfg(test)]
mod tests;
