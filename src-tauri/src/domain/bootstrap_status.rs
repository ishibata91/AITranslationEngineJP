#[derive(Debug, Clone, PartialEq, Eq)]
pub struct BootstrapStatus {
    pub backend_version: String,
    pub boundary_ready: bool,
    pub frontend_entry: &'static str,
}

impl BootstrapStatus {
    pub fn initial(backend_version: String) -> Self {
        Self {
            backend_version,
            boundary_ready: true,
            frontend_entry: "src/main.ts",
        }
    }
}

#[cfg(test)]
mod tests {
    use super::BootstrapStatus;

    #[test]
    fn initial_status_marks_boundary_as_ready() {
        let status = BootstrapStatus::initial("0.1.0".to_string());

        assert_eq!(status.backend_version, "0.1.0");
        assert!(status.boundary_ready);
        assert_eq!(status.frontend_entry, "src/main.ts");
    }
}
