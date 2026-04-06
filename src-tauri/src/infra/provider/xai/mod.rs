use async_trait::async_trait;
use serde::{Deserialize, Serialize};

use crate::application::dto::{
    ExecutionControlFailureCategoryDto, ExecutionControlFailureDto, ProviderExecutionModeDto,
    ProviderSelectionDto,
};
use crate::application::ports::provider_runtime::{ProviderRuntimeFailure, ProviderRuntimePort};

const XAI_PROVIDER_ID: &str = "xai";
const XAI_API_KEY_ENV: &str = "XAI_API_KEY";
const XAI_BASE_URL_ENV: &str = "XAI_BASE_URL";
const XAI_BATCH_NAME_ENV: &str = "XAI_BATCH_NAME";
const XAI_MODEL_ENV: &str = "XAI_BATCH_MODEL";
const XAI_BASE_URL_DEFAULT: &str = "https://api.x.ai";
const XAI_BATCH_NAME_DEFAULT: &str = "ai_translation_engine_batch";
const XAI_MODEL_DEFAULT: &str = "grok-4.20-beta-latest-non-reasoning";
const XAI_POLL_LIMIT_DEFAULT: u32 = 20;
const XAI_POLL_INTERVAL_MS_DEFAULT: u64 = 1_000;
const XAI_RESULTS_PAGE_SIZE_DEFAULT: u32 = 100;

#[derive(Debug, Clone)]
struct XaiConfig {
    base_url: String,
    api_key: Option<String>,
    batch_name: String,
    model: String,
    poll_limit: u32,
    poll_interval_ms: u64,
    results_page_size: u32,
}

impl Default for XaiConfig {
    fn default() -> Self {
        let base_url = read_trimmed_env(XAI_BASE_URL_ENV)
            .unwrap_or_else(|| XAI_BASE_URL_DEFAULT.to_string())
            .trim_end_matches('/')
            .to_string();

        let api_key = read_trimmed_env(XAI_API_KEY_ENV);

        let batch_name = read_trimmed_env(XAI_BATCH_NAME_ENV)
            .unwrap_or_else(|| XAI_BATCH_NAME_DEFAULT.to_string());

        let model =
            read_trimmed_env(XAI_MODEL_ENV).unwrap_or_else(|| XAI_MODEL_DEFAULT.to_string());

        Self {
            base_url,
            api_key,
            batch_name,
            model,
            poll_limit: XAI_POLL_LIMIT_DEFAULT,
            poll_interval_ms: XAI_POLL_INTERVAL_MS_DEFAULT,
            results_page_size: XAI_RESULTS_PAGE_SIZE_DEFAULT,
        }
    }
}

fn read_trimmed_env(key: &str) -> Option<String> {
    std::env::var(key)
        .ok()
        .map(|value| value.trim().to_string())
        .filter(|value| !value.is_empty())
}

#[derive(Debug, Clone, Serialize)]
struct XaiCreateBatchRequest {
    name: String,
}

#[derive(Debug, Deserialize)]
struct XaiCreateBatchResponse {
    #[serde(default)]
    batch_id: Option<String>,
    #[serde(default)]
    id: Option<String>,
}

#[derive(Debug, Clone, Serialize)]
struct XaiAddRequestsRequest {
    batch_requests: Vec<XaiBatchRequestItem>,
}

#[derive(Debug, Clone, Serialize)]
struct XaiBatchRequestItem {
    batch_request_id: String,
    batch_request: XaiBatchRequestPayload,
}

#[derive(Debug, Clone, Serialize)]
struct XaiBatchRequestPayload {
    chat_get_completion: XaiChatGetCompletionRequest,
}

#[derive(Debug, Clone, Serialize)]
struct XaiChatGetCompletionRequest {
    model: String,
    messages: Vec<XaiMessage>,
}

#[derive(Debug, Clone, Serialize)]
struct XaiMessage {
    role: String,
    content: String,
}

#[derive(Debug, Deserialize)]
struct XaiBatchStatusResponse {
    state: XaiBatchState,
}

#[derive(Debug, Deserialize)]
struct XaiBatchState {
    num_requests: u32,
    num_pending: u32,
    num_success: u32,
    num_error: u32,
    num_cancelled: u32,
}

#[derive(Debug, Deserialize)]
struct XaiBatchResultsResponse {
    #[serde(default)]
    succeeded: Vec<XaiSucceededResultItem>,
    #[serde(default)]
    failed: Vec<XaiFailedResultItem>,
    #[serde(default)]
    pagination_token: Option<String>,
}

#[derive(Debug, Deserialize)]
struct XaiFailedResultItem {
    #[allow(dead_code)]
    batch_request_id: Option<String>,
    #[allow(dead_code)]
    error_message: Option<String>,
}

#[derive(Debug, Deserialize)]
struct XaiSucceededResultItem {
    response: XaiChatCompletion,
}

#[derive(Debug, Deserialize)]
struct XaiChatCompletion {
    #[serde(default)]
    choices: Vec<XaiChatChoice>,
}

#[derive(Debug, Deserialize)]
struct XaiChatChoice {
    message: XaiChatMessage,
}

#[derive(Debug, Deserialize)]
struct XaiChatMessage {
    content: String,
}

#[derive(Debug, Clone, Copy)]
enum XaiTransportError {
    Connection,
    Api { status: u16 },
    Response,
}

#[async_trait]
trait XaiBatchTransport: Send + Sync {
    async fn create_batch(
        &self,
        endpoint: &str,
        request: &XaiCreateBatchRequest,
    ) -> Result<String, XaiTransportError>;

    async fn add_requests(
        &self,
        endpoint: &str,
        request: &XaiAddRequestsRequest,
    ) -> Result<(), XaiTransportError>;

    async fn get_batch_status(&self, endpoint: &str) -> Result<String, XaiTransportError>;

    async fn get_batch_results(
        &self,
        endpoint: &str,
        page_size: u32,
        pagination_token: Option<&str>,
    ) -> Result<String, XaiTransportError>;
}

struct ReqwestXaiBatchTransport {
    client: reqwest::Client,
    api_key: String,
}

impl ReqwestXaiBatchTransport {
    fn new(api_key: String) -> Self {
        Self {
            client: reqwest::Client::new(),
            api_key,
        }
    }
}

#[async_trait]
impl XaiBatchTransport for ReqwestXaiBatchTransport {
    async fn create_batch(
        &self,
        endpoint: &str,
        request: &XaiCreateBatchRequest,
    ) -> Result<String, XaiTransportError> {
        let response = self
            .client
            .post(endpoint)
            .header("Authorization", format!("Bearer {}", self.api_key))
            .json(request)
            .send()
            .await
            .map_err(|_| XaiTransportError::Connection)?;

        let status = response.status();
        if !status.is_success() {
            return Err(XaiTransportError::Api {
                status: status.as_u16(),
            });
        }

        response
            .text()
            .await
            .map_err(|_| XaiTransportError::Response)
    }

    async fn add_requests(
        &self,
        endpoint: &str,
        request: &XaiAddRequestsRequest,
    ) -> Result<(), XaiTransportError> {
        let response = self
            .client
            .post(endpoint)
            .header("Authorization", format!("Bearer {}", self.api_key))
            .json(request)
            .send()
            .await
            .map_err(|_| XaiTransportError::Connection)?;

        let status = response.status();
        if !status.is_success() {
            return Err(XaiTransportError::Api {
                status: status.as_u16(),
            });
        }

        Ok(())
    }

    async fn get_batch_status(&self, endpoint: &str) -> Result<String, XaiTransportError> {
        let response = self
            .client
            .get(endpoint)
            .header("Authorization", format!("Bearer {}", self.api_key))
            .send()
            .await
            .map_err(|_| XaiTransportError::Connection)?;

        let status = response.status();
        if !status.is_success() {
            return Err(XaiTransportError::Api {
                status: status.as_u16(),
            });
        }

        response
            .text()
            .await
            .map_err(|_| XaiTransportError::Response)
    }

    async fn get_batch_results(
        &self,
        endpoint: &str,
        page_size: u32,
        pagination_token: Option<&str>,
    ) -> Result<String, XaiTransportError> {
        let mut request = self
            .client
            .get(endpoint)
            .header("Authorization", format!("Bearer {}", self.api_key))
            .query(&[("page_size", page_size.to_string())]);

        if let Some(token) = pagination_token {
            request = request.query(&[("pagination_token", token)]);
        }

        let response = request
            .send()
            .await
            .map_err(|_| XaiTransportError::Connection)?;

        let status = response.status();
        if !status.is_success() {
            return Err(XaiTransportError::Api {
                status: status.as_u16(),
            });
        }

        response
            .text()
            .await
            .map_err(|_| XaiTransportError::Response)
    }
}

pub struct XaiProviderRuntimeAdapter {
    config: XaiConfig,
    transport: Box<dyn XaiBatchTransport>,
}

impl XaiProviderRuntimeAdapter {
    pub fn new() -> Self {
        let config = XaiConfig::default();
        let transport: Box<dyn XaiBatchTransport> = match config.api_key.clone() {
            Some(api_key) => Box::new(ReqwestXaiBatchTransport::new(api_key)),
            None => Box::new(NoopXaiTransport),
        };

        Self { config, transport }
    }

    fn validate_selection_and_config(
        &self,
        selection: &ProviderSelectionDto,
    ) -> Result<(), ProviderRuntimeFailure> {
        if selection.provider_id != XAI_PROVIDER_ID {
            return Err(ExecutionControlFailureDto {
                category: ExecutionControlFailureCategoryDto::ValidationFailure,
                message: "xAI adapter cannot run the selected provider.".to_string(),
            });
        }

        if selection.execution_mode != ProviderExecutionModeDto::Batch {
            return Err(ExecutionControlFailureDto {
                category: ExecutionControlFailureCategoryDto::ValidationFailure,
                message: "xAI adapter supports only batch execution mode.".to_string(),
            });
        }

        if self.config.api_key.is_none() {
            return Err(ExecutionControlFailureDto {
                category: ExecutionControlFailureCategoryDto::ValidationFailure,
                message: "xAI API credentials are not configured.".to_string(),
            });
        }

        if self.config.batch_name.trim().is_empty() || self.config.model.trim().is_empty() {
            return Err(ExecutionControlFailureDto {
                category: ExecutionControlFailureCategoryDto::ValidationFailure,
                message: "xAI adapter configuration is incomplete.".to_string(),
            });
        }

        Ok(())
    }

    fn map_create_batch_request(&self) -> XaiCreateBatchRequest {
        XaiCreateBatchRequest {
            name: self.config.batch_name.clone(),
        }
    }

    fn map_add_requests_payload(&self) -> XaiAddRequestsRequest {
        XaiAddRequestsRequest {
            batch_requests: vec![XaiBatchRequestItem {
                batch_request_id: "provider-step-1".to_string(),
                batch_request: XaiBatchRequestPayload {
                    chat_get_completion: XaiChatGetCompletionRequest {
                        model: self.config.model.clone(),
                        messages: vec![
                            XaiMessage {
                                role: "system".to_string(),
                                content: "Run one provider step.".to_string(),
                            },
                            XaiMessage {
                                role: "user".to_string(),
                                content: "Return one concise completion.".to_string(),
                            },
                        ],
                    },
                },
            }],
        }
    }

    fn batches_endpoint(&self) -> String {
        format!("{}/v1/batches", self.config.base_url)
    }

    fn batch_requests_endpoint(&self, batch_id: &str) -> String {
        format!("{}/v1/batches/{}/requests", self.config.base_url, batch_id)
    }

    fn batch_status_endpoint(&self, batch_id: &str) -> String {
        format!("{}/v1/batches/{}", self.config.base_url, batch_id)
    }

    fn batch_results_endpoint(&self, batch_id: &str) -> String {
        format!("{}/v1/batches/{}/results", self.config.base_url, batch_id)
    }

    fn parse_create_batch_response(body: &str) -> Result<String, ProviderRuntimeFailure> {
        let response = serde_json::from_str::<XaiCreateBatchResponse>(body)
            .map_err(|_| Self::normalize_transport_failure(XaiTransportError::Response))?;

        response
            .batch_id
            .or(response.id)
            .map(|value| value.trim().to_string())
            .filter(|value| !value.is_empty())
            .ok_or_else(|| Self::normalize_transport_failure(XaiTransportError::Response))
    }

    fn parse_batch_status_response(body: &str) -> Result<XaiBatchState, ProviderRuntimeFailure> {
        let response = serde_json::from_str::<XaiBatchStatusResponse>(body)
            .map_err(|_| Self::normalize_transport_failure(XaiTransportError::Response))?;
        Ok(response.state)
    }

    fn parse_batch_results_page(
        body: &str,
    ) -> Result<XaiBatchResultsResponse, ProviderRuntimeFailure> {
        serde_json::from_str::<XaiBatchResultsResponse>(body)
            .map_err(|_| Self::normalize_transport_failure(XaiTransportError::Response))
    }

    fn extract_usable_completion(page: &XaiBatchResultsResponse) -> bool {
        page.succeeded.iter().any(|item| {
            item.response
                .choices
                .iter()
                .any(|choice| !choice.message.content.trim().is_empty())
        })
    }

    fn has_failed_items(page: &XaiBatchResultsResponse) -> bool {
        !page.failed.is_empty()
    }

    async fn poll_until_completed(
        &self,
        batch_id: &str,
    ) -> Result<XaiBatchState, ProviderRuntimeFailure> {
        let status_endpoint = self.batch_status_endpoint(batch_id);
        let max_attempts = self.config.poll_limit.max(1);

        for attempt in 0..max_attempts {
            let body = self
                .transport
                .get_batch_status(&status_endpoint)
                .await
                .map_err(Self::normalize_transport_failure)?;
            let state = Self::parse_batch_status_response(&body)?;

            if state.num_pending == 0 {
                return Ok(state);
            }

            if attempt + 1 < max_attempts && self.config.poll_interval_ms > 0 {
                tokio::time::sleep(std::time::Duration::from_millis(
                    self.config.poll_interval_ms,
                ))
                .await;
            }
        }

        Err(ExecutionControlFailureDto {
            category: ExecutionControlFailureCategoryDto::RecoverableProviderFailure,
            message: "xAI batch did not complete before the polling limit. Retry may recover."
                .to_string(),
        })
    }

    async fn fetch_results_until_end(
        &self,
        batch_id: &str,
    ) -> Result<(bool, bool), ProviderRuntimeFailure> {
        let endpoint = self.batch_results_endpoint(batch_id);
        let mut has_usable_completion = false;
        let mut has_failed_items = false;
        let mut pagination_token: Option<String> = None;

        loop {
            let body = self
                .transport
                .get_batch_results(
                    &endpoint,
                    self.config.results_page_size,
                    pagination_token.as_deref(),
                )
                .await
                .map_err(Self::normalize_transport_failure)?;

            let page = Self::parse_batch_results_page(&body)?;
            has_usable_completion = has_usable_completion || Self::extract_usable_completion(&page);
            has_failed_items = has_failed_items || Self::has_failed_items(&page);

            match page.pagination_token {
                Some(token) if !token.trim().is_empty() => pagination_token = Some(token),
                _ => break,
            }
        }

        Ok((has_usable_completion, has_failed_items))
    }

    fn normalize_transport_failure(error: XaiTransportError) -> ProviderRuntimeFailure {
        match error {
            XaiTransportError::Connection => ExecutionControlFailureDto {
                category: ExecutionControlFailureCategoryDto::RecoverableProviderFailure,
                message: "xAI provider connection failed. Retry may recover.".to_string(),
            },
            XaiTransportError::Api {
                status: 408 | 429 | 500 | 502 | 503 | 504,
            } => ExecutionControlFailureDto {
                category: ExecutionControlFailureCategoryDto::RecoverableProviderFailure,
                message: "xAI provider is temporarily unavailable. Retry may recover.".to_string(),
            },
            XaiTransportError::Api { .. } | XaiTransportError::Response => {
                ExecutionControlFailureDto {
                    category: ExecutionControlFailureCategoryDto::UnrecoverableProviderFailure,
                    message: "xAI provider returned an unusable batch response.".to_string(),
                }
            }
        }
    }

    #[cfg(test)]
    fn with_config_and_transport(config: XaiConfig, transport: Box<dyn XaiBatchTransport>) -> Self {
        Self { config, transport }
    }
}

impl Default for XaiProviderRuntimeAdapter {
    fn default() -> Self {
        Self::new()
    }
}

#[async_trait]
impl ProviderRuntimePort for XaiProviderRuntimeAdapter {
    async fn run_provider_step(
        &self,
        selection: ProviderSelectionDto,
    ) -> Result<(), ProviderRuntimeFailure> {
        self.validate_selection_and_config(&selection)?;

        let create_body = self
            .transport
            .create_batch(&self.batches_endpoint(), &self.map_create_batch_request())
            .await
            .map_err(Self::normalize_transport_failure)?;

        let batch_id = Self::parse_create_batch_response(&create_body)?;

        self.transport
            .add_requests(
                &self.batch_requests_endpoint(&batch_id),
                &self.map_add_requests_payload(),
            )
            .await
            .map_err(Self::normalize_transport_failure)?;

        let state = self.poll_until_completed(&batch_id).await?;
        let (has_usable_completion, has_failed_items) =
            self.fetch_results_until_end(&batch_id).await?;

        if !has_usable_completion {
            return Err(ExecutionControlFailureDto {
                category: ExecutionControlFailureCategoryDto::UnrecoverableProviderFailure,
                message: "xAI batch results did not contain a usable completion.".to_string(),
            });
        }

        if has_failed_items
            || state.num_success == 0
            || state.num_requests == 0
            || state.num_error > 0
            || state.num_cancelled > 0
        {
            return Err(ExecutionControlFailureDto {
                category: ExecutionControlFailureCategoryDto::UnrecoverableProviderFailure,
                message: "xAI batch completed without a fully usable result set.".to_string(),
            });
        }

        Ok(())
    }
}

struct NoopXaiTransport;

#[async_trait]
impl XaiBatchTransport for NoopXaiTransport {
    async fn create_batch(
        &self,
        _endpoint: &str,
        _request: &XaiCreateBatchRequest,
    ) -> Result<String, XaiTransportError> {
        Err(XaiTransportError::Connection)
    }

    async fn add_requests(
        &self,
        _endpoint: &str,
        _request: &XaiAddRequestsRequest,
    ) -> Result<(), XaiTransportError> {
        Err(XaiTransportError::Connection)
    }

    async fn get_batch_status(&self, _endpoint: &str) -> Result<String, XaiTransportError> {
        Err(XaiTransportError::Connection)
    }

    async fn get_batch_results(
        &self,
        _endpoint: &str,
        _page_size: u32,
        _pagination_token: Option<&str>,
    ) -> Result<String, XaiTransportError> {
        Err(XaiTransportError::Connection)
    }
}

#[cfg(test)]
mod tests;
