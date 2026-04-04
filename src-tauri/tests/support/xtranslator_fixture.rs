use std::fs;
use std::path::PathBuf;
use std::sync::atomic::{AtomicU64, Ordering};

pub struct FixtureFile {
    dir_path: PathBuf,
    file_path: PathBuf,
}

impl FixtureFile {
    pub fn from_bytes(file_name: &str, contents: &[u8]) -> Self {
        let dir_path = std::env::temp_dir().join(format!(
            "ai-translation-engine-jp-xtranslator-{file_name}-{}",
            next_unique_test_suffix()
        ));
        let file_path = dir_path.join(file_name);

        fs::create_dir_all(&dir_path).expect("fixture directory should be created");
        fs::write(&file_path, contents).expect("fixture file should be written");

        Self {
            dir_path,
            file_path,
        }
    }

    pub fn path_string(&self) -> String {
        self.file_path.to_string_lossy().into_owned()
    }
}

impl Drop for FixtureFile {
    fn drop(&mut self) {
        let _ = fs::remove_dir_all(&self.dir_path);
    }
}

pub fn shared_contract_fixture_file() -> FixtureFile {
    FixtureFile::from_bytes(
        "Skyrim Base Terms.sst",
        include_bytes!(
            "../validation/dictionary-rebuild/fixtures/xtranslator-shared-reusable-entry.sst"
        ),
    )
}

#[allow(dead_code)]
pub fn build_xtranslator_sst_bytes(entries: &[(&str, &str)]) -> Vec<u8> {
    let mut bytes = Vec::new();

    // xTranslator v8 SST layout with zeroed metadata keeps the fixture focused on source/dest pairs.
    bytes.extend_from_slice(&0x3955_5353u32.to_le_bytes());
    bytes.push(0);
    bytes.extend_from_slice(&0i32.to_le_bytes());
    bytes.extend_from_slice(&0i32.to_le_bytes());

    for (source_text, dest_text) in entries {
        bytes.push(0);
        bytes.extend_from_slice(&0i32.to_le_bytes());
        bytes.extend_from_slice(&0u32.to_le_bytes());
        bytes.extend_from_slice(&[0; 4]);
        bytes.extend_from_slice(&[0; 4]);
        bytes.extend_from_slice(&0u16.to_le_bytes());
        bytes.extend_from_slice(&0u16.to_le_bytes());
        bytes.extend_from_slice(&0u32.to_le_bytes());
        bytes.push(0);
        bytes.push(1);

        let source_bytes = encode_utf16le(source_text);
        bytes.extend_from_slice(&(source_bytes.len() as i32).to_le_bytes());
        bytes.extend_from_slice(&source_bytes);

        let dest_bytes = encode_utf16le(dest_text);
        bytes.extend_from_slice(&(dest_bytes.len() as i32).to_le_bytes());
        bytes.extend_from_slice(&dest_bytes);
    }

    bytes
}

#[allow(dead_code)]
fn encode_utf16le(value: &str) -> Vec<u8> {
    value
        .encode_utf16()
        .flat_map(|unit| unit.to_le_bytes())
        .collect()
}

fn next_unique_test_suffix() -> u64 {
    static COUNTER: AtomicU64 = AtomicU64::new(0);
    COUNTER.fetch_add(1, Ordering::Relaxed) + 1
}
