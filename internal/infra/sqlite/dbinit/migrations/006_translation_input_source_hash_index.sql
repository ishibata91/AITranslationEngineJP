CREATE UNIQUE INDEX IF NOT EXISTS idx_x_edit_extracted_data_source_content_hash
  ON X_EDIT_EXTRACTED_DATA(source_content_hash)
  WHERE source_content_hash <> '';