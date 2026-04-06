use async_trait::async_trait;
use serde::{Deserialize, Serialize};

use crate::application::dto::{
    ExecutionControlFailureCategoryDto, ExecutionControlFailureDto, ProviderExecutionModeDto,
    ProviderSelectionDto,
};
use crate::application::ports::provider_runtime::{ProviderRuntimeFailure, ProviderRuntimePort};

const LMSTUDIO_PROVIDER_ID: &str = "lmstudio";

#[derive(Debug, Clone)]
struct LmstudioConfig {
    endpoint: String,
    model: String,
}

impl Default for LmstudioConfig {
    fn default() -> Self {
        Self {
            endpoint: "http://127.0.0.1:1234/v1/chat/completions".to_string(),
            model: "local-model".to_string(),
        }
    }
}

#[derive(Debug, Clone, Serialize)]
struct LmstudioRequest {
    model: String,
    messages: Vec<LmstudioMessage>,
}

#[derive(Debug, Clone, Serialize)]
struct LmstudioMessage {
    role: String,
    content: String,
}

#[derive(Debug, Deserialize)]
struct LmstudioResponse {
    choices: Vec<LmstudioChoice>,
}

#[derive(Debug, Deserialize)]
struct LmstudioChoice {
    message: LmstudioResponseMessage,
}

#[derive(Debug, Deserialize)]
struct LmstudioResponseMessage {
    content: String,
}

#[derive(Debug)]
enum LmstudioTransportError {
    Connection,
    Response,
}

#[async_trait]
trait LmstudioTransport: Send + Sync {
    async fn execute(
        &self,
        endpoint: &str,
        request: &LmstudioRequest,
    ) -> Result<String, LmstudioTransportError>;
}

struct ReqwestLmstudioTransport {
    client: reqwest::Client,
}

impl ReqwestLmstudioTransport {
    fn new() -> Self {
        Self {
            client: reqwest::Client::new(),
        }
    }
}

#[async_trait]
impl LmstudioTransport for ReqwestLmstudioTransport {
    async fn execute(
        &self,
        endpoint: &str,
        request: &LmstudioRequest,
    ) -> Result<String, LmstudioTransportError> {
        let response = self
            .client
            .post(endpoint)
            .json(request)
            .send()
            .await
            .map_err(|_| LmstudioTransportError::Connection)?;

        if !response.status().is_success() {
            if response.status().is_client_error() {
                return Err(LmstudioTransportError::Response);
            }
            return Err(LmstudioTransportError::Connection);
        }

        response
            .text()
            .await
            .map_err(|_| LmstudioTransportError::Response)
    }
}

pub struct LmstudioProviderRuntimeAdapter {
    config: LmstudioConfig,
    transport: Box<dyn LmstudioTransport>,
}

impl LmstudioProviderRuntimeAdapter {
    pub fn new() -> Self {
        Self {
            config: LmstudioConfig::default(),
            transport: Box::new(ReqwestLmstudioTransport::new()),
        }
    }

    fn map_request(&self) -> LmstudioRequest {
        LmstudioRequest {
            model: self.config.model.clone(),
            messages: vec![LmstudioMessage {
                role: "system".to_string(),
                content: "Run one provider step.".to_string(),
            }],
        }
    }

    fn validate_selection(selection: &ProviderSelectionDto) -> Result<(), ProviderRuntimeFailure> {
        if selection.provider_id != LMSTUDIO_PROVIDER_ID {
            return Err(ExecutionControlFailureDto {
                category: ExecutionControlFailureCategoryDto::ValidationFailure,
                message: "LMStudio adapter cannot run the selected provider.".to_string(),
            });
        }

        if selection.execution_mode != ProviderExecutionModeDto::Batch {
            return Err(ExecutionControlFailureDto {
                category: ExecutionControlFailureCategoryDto::ValidationFailure,
                message: "LMStudio adapter supports only batch execution mode.".to_string(),
            });
        }

        Ok(())
    }

    fn normalize_transport_failure(error: LmstudioTransportError) -> ProviderRuntimeFailure {
        match error {
            LmstudioTransportError::Connection => ExecutionControlFailureDto {
                category: ExecutionControlFailureCategoryDto::RecoverableProviderFailure,
                message: "LMStudio connection failed. Retry may recover.".to_string(),
            },
            LmstudioTransportError::Response => ExecutionControlFailureDto {
                category: ExecutionControlFailureCategoryDto::UnrecoverableProviderFailure,
                message: "LMStudio response was malformed.".to_string(),
            },
        }
    }

    fn parse_response(body: &str) -> Result<(), ProviderRuntimeFailure> {
        let response = serde_json::from_str::<LmstudioResponse>(body)
            .map_err(|_| Self::normalize_transport_failure(LmstudioTransportError::Response))?;

        let has_content = response
            .choices
            .first()
            .map(|choice| !choice.message.content.trim().is_empty())
            .unwrap_or(false);

        if !has_content {
            return Err(Self::normalize_transport_failure(
                LmstudioTransportError::Response,
            ));
        }

        Ok(())
    }

    #[cfg(test)]
    fn with_transport(transport: Box<dyn LmstudioTransport>) -> Self {
        Self {
            config: LmstudioConfig::default(),
            transport,
        }
    }
}

impl Default for LmstudioProviderRuntimeAdapter {
    fn default() -> Self {
        Self::new()
    }
}

#[async_trait]
impl ProviderRuntimePort for LmstudioProviderRuntimeAdapter {
    async fn run_provider_step(
        &self,
        selection: ProviderSelectionDto,
    ) -> Result<(), ProviderRuntimeFailure> {
        Self::validate_selection(&selection)?;

        let request = self.map_request();
        let body = self
            .transport
            .execute(&self.config.endpoint, &request)
            .await
            .map_err(Self::normalize_transport_failure)?;

        Self::parse_response(&body)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::sync::{
        atomic::{AtomicUsize, Ordering},
        Arc,
    };

    struct SpyTransport {
        calls: Arc<AtomicUsize>,
        response: Result<String, LmstudioTransportError>,
    }

    #[async_trait]
    impl LmstudioTransport for SpyTransport {
        async fn execute(
            &self,
            _endpoint: &str,
            _request: &LmstudioRequest,
        ) -> Result<String, LmstudioTransportError> {
            self.calls.fetch_add(1, Ordering::SeqCst);
            self.response
                .as_ref()
                .map(Clone::clone)
                .map_err(|error| match error {
                    LmstudioTransportError::Connection => LmstudioTransportError::Connection,
                    LmstudioTransportError::Response => LmstudioTransportError::Response,
                })
        }
    }

    fn selection(provider_id: &str, mode: ProviderExecutionModeDto) -> ProviderSelectionDto {
        ProviderSelectionDto {
            provider_id: provider_id.to_string(),
            execution_mode: mode,
            runtime_settings: crate::application::dto::ProviderRuntimeSettingsDto {
                retry_limit: 1,
                max_concurrency: 1,
                pause_supported: false,
            },
        }
    }

    #[tokio::test]
    async fn lmstudio_rejects_provider_mismatch_without_transport() {
        let calls = Arc::new(AtomicUsize::new(0));
        let adapter = LmstudioProviderRuntimeAdapter::with_transport(Box::new(SpyTransport {
            calls: Arc::clone(&calls),
            response: Ok("{\"choices\":[{\"message\":{\"content\":\"ok\"}}]}".to_string()),
        }));

        let error = adapter
            .run_provider_step(selection("gemini", ProviderExecutionModeDto::Batch))
            .await
            .expect_err("provider mismatch must fail validation");

        assert_eq!(
            error.category,
            ExecutionControlFailureCategoryDto::ValidationFailure
        );
        assert_eq!(calls.load(Ordering::SeqCst), 0);
    }

    #[tokio::test]
    async fn lmstudio_rejects_unsupported_mode_without_transport() {
        let calls = Arc::new(AtomicUsize::new(0));
        let adapter = LmstudioProviderRuntimeAdapter::with_transport(Box::new(SpyTransport {
            calls: Arc::clone(&calls),
            response: Ok("{\"choices\":[{\"message\":{\"content\":\"ok\"}}]}".to_string()),
        }));

        let error = adapter
            .run_provider_step(selection("lmstudio", ProviderExecutionModeDto::Streaming))
            .await
            .expect_err("unsupported mode must fail validation");

        assert_eq!(
            error.category,
            ExecutionControlFailureCategoryDto::ValidationFailure
        );
        assert_eq!(calls.load(Ordering::SeqCst), 0);
    }

    #[tokio::test]
    async fn lmstudio_normalizes_connection_failure_without_endpoint_leak() {
        let adapter = LmstudioProviderRuntimeAdapter::with_transport(Box::new(SpyTransport {
            calls: Arc::new(AtomicUsize::new(0)),
            response: Err(LmstudioTransportError::Connection),
        }));

        let error = adapter
            .run_provider_step(selection("lmstudio", ProviderExecutionModeDto::Batch))
            .await
            .expect_err("connection failure must return recoverable provider failure");

        assert_eq!(
            error.category,
            ExecutionControlFailureCategoryDto::RecoverableProviderFailure
        );
        assert!(!error.message.contains("127.0.0.1"));
        assert!(!error.message.contains("chat/completions"));
    }

    #[tokio::test]
    async fn lmstudio_normalizes_malformed_response_without_raw_body_leak() {
        let malformed = "{\"oops\":true}";
        let adapter = LmstudioProviderRuntimeAdapter::with_transport(Box::new(SpyTransport {
            calls: Arc::new(AtomicUsize::new(0)),
            response: Ok(malformed.to_string()),
        }));

        let error = adapter
            .run_provider_step(selection("lmstudio", ProviderExecutionModeDto::Batch))
            .await
            .expect_err("malformed response must return unrecoverable provider failure");

        assert_eq!(
            error.category,
            ExecutionControlFailureCategoryDto::UnrecoverableProviderFailure
        );
        assert!(!error.message.contains("oops"));
        assert!(!error.message.contains("true"));
    }

    #[tokio::test]
    async fn lmstudio_completes_when_transport_returns_valid_response() {
        let adapter = LmstudioProviderRuntimeAdapter::with_transport(Box::new(SpyTransport {
            calls: Arc::new(AtomicUsize::new(0)),
            response: Ok("{\"choices\":[{\"message\":{\"content\":\"ok\"}}]}".to_string()),
        }));

        let result = adapter
            .run_provider_step(selection("lmstudio", ProviderExecutionModeDto::Batch))
            .await;

        assert!(result.is_ok(), "valid lmstudio response should complete");
    }

    #[tokio::test]
    async fn lmstudio_normalizes_permanent_http_failure_to_unrecoverable() {
        let adapter = LmstudioProviderRuntimeAdapter::with_transport(Box::new(SpyTransport {
            calls: Arc::new(AtomicUsize::new(0)),
            response: Err(LmstudioTransportError::Response),
        }));

        let error = adapter
            .run_provider_step(selection("lmstudio", ProviderExecutionModeDto::Batch))
            .await
            .expect_err("permanent http failure must return unrecoverable provider failure");

        assert_eq!(
            error.category,
            ExecutionControlFailureCategoryDto::UnrecoverableProviderFailure
        );
    }
}
