package engine

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"
)

// 箱（concept-model の 4 箱）。record_type_master.box の値と一致させる。取込段の振り分けキー。
const (
	boxNarration  = "叙述文"
	boxProperNoun = "固有名"
	boxSetPhrase  = "定型句"
	boxLine       = "台詞"
)

// RecordKey は record_type_master の引きキー (rec, field)。
type RecordKey struct {
	Rec   string
	Field string
}

// Dispatched は取込段が extracted_field を箱ごとに振り分けた結果。
// Narrations は叙述文・定型句（自前で訳す本文）、ProperNouns は固有名（重複排除する訳の単位）、
// Lines は台詞（話者で口調が変わる本文）。Skipped は master に無い（翻訳しない）原文の件数。
type Dispatched struct {
	Narrations  []model.Narration
	ProperNouns []model.ProperNoun
	Lines       []model.Line
	Skipped     int
}

// Dispatch は extracted_field を record_type_master の box で narration/proper_noun/line へ振り分ける純粋関数。
// master に無い (rec, field) は翻訳対象でないため読み込まず Skipped に数える。副作用を持たない。
// narration.Style には割り当て directive のキーを入れる（由来表示用。翻訳時の正は record_type_master）。
func Dispatch(fields []model.ExtractedField, master map[RecordKey]model.RecordType) Dispatched {
	var d Dispatched
	for _, f := range fields {
		rt, ok := master[RecordKey{Rec: f.Rec, Field: f.Field}]
		if !ok {
			d.Skipped++
			continue
		}
		switch rt.Box {
		case boxProperNoun:
			d.ProperNouns = append(d.ProperNouns, model.ProperNoun{Source: f.Source, Category: f.Rec})
		case boxNarration, boxSetPhrase:
			d.Narrations = append(d.Narrations, model.Narration{
				Source: f.Source, Style: rt.Directive,
				Plugin: f.Plugin, FormID: f.FormID, EDID: f.EDID, Rec: f.Rec, Field: f.Field, Ordinal: f.Ordinal,
			})
		case boxLine:
			d.Lines = append(d.Lines, model.Line{
				Source: f.Source, ResponseOrder: f.Ordinal,
				Plugin: f.Plugin, FormID: f.FormID, EDID: f.EDID, Rec: f.Rec, Field: f.Field, Ordinal: f.Ordinal,
			})
		default:
			d.Skipped++
		}
	}
	return d
}

// recordMasterMap は record_type_master の行群を取込段の引きキー (rec, field) の map に畳む。
func recordMasterMap(rows []model.RecordType) map[RecordKey]model.RecordType {
	m := make(map[RecordKey]model.RecordType, len(rows))
	for _, r := range rows {
		m[RecordKey{Rec: r.Rec, Field: r.Field}] = r
	}
	return m
}

// IngestStore は engine が取込段に使う中心データアクセス（使う分だけ宣言する）。
type IngestStore interface {
	ListExtractedFields(ctx context.Context) ([]model.ExtractedField, error)
	ListRecordTypeMaster(ctx context.Context) ([]model.RecordType, error)
	IngestNarrations(ctx context.Context, rows []model.Narration) (int, error)
	IngestProperNouns(ctx context.Context, rows []model.ProperNoun) (int, error)
	IngestLines(ctx context.Context, rows []model.Line) (int, error)
	LinkLineSpeakersFromStaging(ctx context.Context) error
	LinkLineConditionsFromStaging(ctx context.Context) error
}

// IngestCounts は取込段で各箱へ実際に追加した件数（再取込での冪等を観測できるよう件数を返す）。
type IngestCounts struct {
	Narrations  int
	ProperNouns int
	Lines       int
	Skipped     int
}

// Ingest は extracted_field を読み、record_type_master の box で narration/proper_noun/line へ振り分けて投入する。
// 本文フェーズ・固有名フェーズより前に 1 度呼ぶ。話者連関は line 投入後に staging から解決する。
func (e *Engine) Ingest(ctx context.Context) (IngestCounts, error) {
	fields, err := e.store.ListExtractedFields(ctx)
	if err != nil {
		return IngestCounts{}, fmt.Errorf("取込段の extracted_field 取得: %w", err)
	}
	masterRows, err := e.store.ListRecordTypeMaster(ctx)
	if err != nil {
		return IngestCounts{}, fmt.Errorf("取込段の record_type_master 取得: %w", err)
	}
	d := Dispatch(fields, recordMasterMap(masterRows))

	pn, err := e.store.IngestProperNouns(ctx, d.ProperNouns)
	if err != nil {
		return IngestCounts{}, fmt.Errorf("固有名の投入: %w", err)
	}
	nr, err := e.store.IngestNarrations(ctx, d.Narrations)
	if err != nil {
		return IngestCounts{}, fmt.Errorf("叙述文・定型句の投入: %w", err)
	}
	ln, err := e.store.IngestLines(ctx, d.Lines)
	if err != nil {
		return IngestCounts{}, fmt.Errorf("台詞の投入: %w", err)
	}
	if err := e.store.LinkLineSpeakersFromStaging(ctx); err != nil {
		return IngestCounts{}, fmt.Errorf("話者連関の解決: %w", err)
	}
	// 条件由来の性別（汎用台詞の一人称・語尾の根拠）を line_condition へ解決する（話者連関と対称）。
	if err := e.store.LinkLineConditionsFromStaging(ctx); err != nil {
		return IngestCounts{}, fmt.Errorf("条件由来の性別の解決: %w", err)
	}
	return IngestCounts{Narrations: nr, ProperNouns: pn, Lines: ln, Skipped: d.Skipped}, nil
}
