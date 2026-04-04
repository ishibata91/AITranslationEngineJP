use std::fs;
use std::path::Path;

use async_trait::async_trait;

use crate::application::dictionary_import::DictionaryImporter;
use crate::application::dto::{
    DictionaryImportRequestDto, DictionaryImportResultDto, ReusableDictionaryEntryDto,
};

const XTRANSLATOR_SST_MAGIC: u32 = 0x3955_5353;
const HEADER_LEN: usize = 13;
const ENTRY_PREFIX_LEN: usize = 27;

pub struct FileSystemXtranslatorImporter;

#[async_trait]
impl DictionaryImporter for FileSystemXtranslatorImporter {
    async fn import_dictionary(
        &self,
        request: &DictionaryImportRequestDto,
    ) -> Result<DictionaryImportResultDto, String> {
        let file_path = Path::new(&request.source_file_path);
        let bytes = fs::read(file_path).map_err(|error| {
            format!(
                "Failed to read xTranslator dictionary file {}: {error}",
                file_path.display()
            )
        })?;
        let entries = parse_sst_entries(&bytes, file_path)?;

        Ok(DictionaryImportResultDto {
            dictionary_name: dictionary_name_from_path(file_path),
            source_type: request.source_type.clone(),
            entries,
        })
    }
}

fn dictionary_name_from_path(path: &Path) -> String {
    path.file_stem()
        .and_then(|stem| stem.to_str())
        .filter(|stem| !stem.is_empty())
        .map(ToOwned::to_owned)
        .unwrap_or_else(|| "xtranslator-dictionary".to_string())
}

fn parse_sst_entries(
    bytes: &[u8],
    file_path: &Path,
) -> Result<Vec<ReusableDictionaryEntryDto>, String> {
    if bytes.len() < HEADER_LEN {
        return Err(format!(
            "Failed to parse xTranslator SST file {}: header is too short",
            file_path.display()
        ));
    }

    let magic = read_u32_le(bytes, 0)?;
    if magic != XTRANSLATOR_SST_MAGIC {
        return Err(format!(
            "Failed to parse xTranslator SST file {}: invalid header magic",
            file_path.display()
        ));
    }

    let mut cursor = HEADER_LEN;
    let mut entries = Vec::new();

    while cursor < bytes.len() {
        if bytes.len() - cursor < ENTRY_PREFIX_LEN {
            return Err(format!(
                "Failed to parse xTranslator SST file {}: truncated entry header",
                file_path.display()
            ));
        }

        cursor += ENTRY_PREFIX_LEN;

        let source_len = read_i32_le(bytes, cursor)?;
        cursor += 4;
        let source_text = read_utf16le_string(bytes, &mut cursor, source_len, file_path)?;

        let dest_len = read_i32_le(bytes, cursor)?;
        cursor += 4;
        let dest_text = read_utf16le_string(bytes, &mut cursor, dest_len, file_path)?;

        entries.push(ReusableDictionaryEntryDto {
            source_text,
            dest_text,
        });
    }

    Ok(entries)
}

fn read_u32_le(bytes: &[u8], offset: usize) -> Result<u32, String> {
    let slice = bytes
        .get(offset..offset + 4)
        .ok_or_else(|| "not enough bytes to read u32".to_string())?;
    Ok(u32::from_le_bytes(
        slice.try_into().expect("slice length should match"),
    ))
}

fn read_i32_le(bytes: &[u8], offset: usize) -> Result<i32, String> {
    let slice = bytes
        .get(offset..offset + 4)
        .ok_or_else(|| "not enough bytes to read i32".to_string())?;
    Ok(i32::from_le_bytes(
        slice.try_into().expect("slice length should match"),
    ))
}

fn read_utf16le_string(
    bytes: &[u8],
    cursor: &mut usize,
    byte_len: i32,
    file_path: &Path,
) -> Result<String, String> {
    if byte_len < 0 {
        return Err(format!(
            "Failed to parse xTranslator SST file {}: negative string length",
            file_path.display()
        ));
    }

    let byte_len = byte_len as usize;
    if !byte_len.is_multiple_of(2) {
        return Err(format!(
            "Failed to parse xTranslator SST file {}: UTF-16 byte length is not even",
            file_path.display()
        ));
    }

    let slice = bytes.get(*cursor..*cursor + byte_len).ok_or_else(|| {
        format!(
            "Failed to parse xTranslator SST file {}: truncated UTF-16 payload",
            file_path.display()
        )
    })?;
    *cursor += byte_len;

    let units = slice
        .chunks_exact(2)
        .map(|chunk| u16::from_le_bytes([chunk[0], chunk[1]]))
        .collect::<Vec<_>>();
    String::from_utf16(&units).map_err(|error| {
        format!(
            "Failed to parse xTranslator SST file {}: invalid UTF-16 payload: {error}",
            file_path.display()
        )
    })
}
