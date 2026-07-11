package store

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"
)

// ExportPlugins は書き出し対象になりうる plugin 名を昇順で返す。
// narration・line・extracted_field に現れる plugin を集める（訳出済みかは問わない列挙）。
func (s *Store) ExportPlugins(ctx context.Context) ([]string, error) {
	var plugins []string
	if err := s.db.SelectContext(ctx, &plugins,
		`SELECT plugin FROM narration
		 UNION SELECT plugin FROM line
		 UNION SELECT plugin FROM extracted_field
		 ORDER BY plugin`); err != nil {
		return nil, fmt.Errorf("書き出し対象 plugin の取得: %w", err)
	}
	return plugins, nil
}

// NarrationsForExport は指定 plugin の訳出済み（status != 0）叙述文・定型句を、レコード位置順で返す。
// 未訳（status = 0）は書き出し対象外のため除く。
func (s *Store) NarrationsForExport(ctx context.Context, plugin string) ([]model.Narration, error) {
	return s.queryNarrations(ctx,
		`SELECT `+narrationColumns+` FROM narration
		 WHERE plugin = ? AND status != 0 ORDER BY form_id, rec, field, ordinal`, plugin)
}

// LinesForExport は指定 plugin の訳出済み（status != 0）台詞を、レコード位置順で返す。
// 未訳（status = 0）は書き出し対象外のため除く。
func (s *Store) LinesForExport(ctx context.Context, plugin string) ([]model.Line, error) {
	return s.queryLines(ctx,
		`SELECT `+lineColumns+` FROM line
		 WHERE plugin = ? AND status != 0 ORDER BY form_id, rec, field, ordinal`, plugin)
}

// ProperNounPlacementsForExport は指定 plugin の固有名を、原文出現位置（extracted_field）へ確定訳（proper_noun）を
// 結んで返す。box='固有名' の行だけを対象にし、確定訳（status != 0）のあるものだけ返す。
// 結合キーは「plugin 一致かつ category = rec かつ source 一致」。proper_noun が plugin スコープの非共有に
// なったため、pn.plugin = ef.plugin を足して当該 plugin の固有名だけを位置解決する（別 plugin の同綴り訳を混ぜない）。
// category を rec で絞るため同綴り異義（人名と地名など）を取り違えない。box 判定は record_type_master を rec・field で引く。
func (s *Store) ProperNounPlacementsForExport(ctx context.Context, plugin string) ([]model.ProperNounPlacement, error) {
	var rows []model.ProperNounPlacement
	if err := s.db.SelectContext(ctx, &rows,
		`SELECT ef.plugin, ef.form_id, ef.edid, ef.rec, ef.field, ef.ordinal, ef.source, pn.dest, pn.status
		 FROM extracted_field ef
		 JOIN record_type_master rtm
		   ON rtm.rec = ef.rec AND rtm.field = ef.field AND rtm.box = '固有名'
		 JOIN proper_noun pn
		   ON pn.plugin = ef.plugin AND pn.category = ef.rec AND pn.source = ef.source
		 WHERE ef.plugin = ? AND pn.status != 0
		 ORDER BY ef.form_id, ef.rec, ef.field, ef.ordinal`, plugin); err != nil {
		return nil, fmt.Errorf("書き出し用 固有名の位置解決: %w", err)
	}
	return rows, nil
}
