package store

import (
	"context"
	"fmt"
	"strings"
)

// speakerChunkSize は IN 句に渡す host parameter の上限を避けるための分割サイズ。
// SQLite の SQLITE_MAX_VARIABLE_NUMBER を十分に下回る値にする。
const speakerChunkSize = 900

// count は COUNT(*) の単一値クエリを実行して件数を返す。
func (s *Store) count(ctx context.Context, query string) (int, error) {
	var n int
	if err := s.db.GetContext(ctx, &n, query); err != nil {
		return 0, fmt.Errorf("件数の取得: %w", err)
	}
	return n, nil
}

// placeholders は id 群を IN 句用の "?,?,..." と引数列へ変換する。
func placeholders(ids []int64) (string, []any) {
	marks := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	return marks, args
}

// chunkInt64 は id 群を size ごとの塊へ分割する。IN 句の host parameter 上限回避に使う。
func chunkInt64(ids []int64, size int) [][]int64 {
	var chunks [][]int64
	for start := 0; start < len(ids); start += size {
		end := min(start+size, len(ids))
		chunks = append(chunks, ids[start:end])
	}
	return chunks
}
