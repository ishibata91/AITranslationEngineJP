package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"aitranslationenginejp/internal/core/termxml"
	"aitranslationenginejp/internal/model"
)

// ExportStore は engine が xTranslator 書き出しに使う中心データアクセス（使う分だけ宣言する）。
type ExportStore interface {
	ExportPlugins(ctx context.Context) ([]string, error)
	NarrationsForExport(ctx context.Context, plugin string) ([]model.Narration, error)
	LinesForExport(ctx context.Context, plugin string) ([]model.Line, error)
	ProperNounPlacementsForExport(ctx context.Context, plugin string) ([]model.ProperNounPlacement, error)
}

// ExportXTranslator は中心 DB の訳出済み翻訳結果を、plugin ごとに xTranslator 互換 XML ファイル（SSTXMLRessources 形式）へ書き出す。
// 出力先は outDir 直下で、ファイル名は "<plugin 名>_english_japanese.xml"（拡張子は除く）。
// 書き出したファイルのフルパスを返す。訳出済みの行が 1 つも無い plugin はファイルを作らない。
func (e *Engine) ExportXTranslator(ctx context.Context, outDir string) ([]string, error) {
	plugins, err := e.store.ExportPlugins(ctx)
	if err != nil {
		return nil, err
	}
	var written []string
	for _, plugin := range plugins {
		entries, err := e.exportEntries(ctx, plugin)
		if err != nil {
			return written, err
		}
		if len(entries) == 0 {
			continue // 訳出済みが無い plugin はファイルを作らない。
		}
		addon := strings.TrimSuffix(plugin, filepath.Ext(plugin))
		data, err := termxml.MarshalStrings(entries, addon)
		if err != nil {
			return written, fmt.Errorf("%s の XML 生成: %w", plugin, err)
		}
		path := filepath.Join(outDir, addon+"_english_japanese.xml")
		if err := os.WriteFile(path, data, 0o644); err != nil { //nolint:gosec // 利用者が選んだ出力フォルダ配下への翻訳結果の書き出し。読み取り可能な成果物のため 0644。
			return written, fmt.Errorf("%s の書き込み: %w", path, err)
		}
		written = append(written, path)
	}
	return written, nil
}

// exportEntries は 1 plugin の訳出済み叙述文・台詞・固有名を xTranslator の <String> 行群へ組み立てる。
// 訳状態（DB status）は partialOf で xTranslator の Partial 属性へ写す。
// 固有名は位置を持たない proper_noun を extracted_field へ結んだ ProperNounPlacement から組む。
func (e *Engine) exportEntries(ctx context.Context, plugin string) ([]termxml.StringEntry, error) {
	narrations, err := e.store.NarrationsForExport(ctx, plugin)
	if err != nil {
		return nil, err
	}
	lines, err := e.store.LinesForExport(ctx, plugin)
	if err != nil {
		return nil, err
	}
	propers, err := e.store.ProperNounPlacementsForExport(ctx, plugin)
	if err != nil {
		return nil, err
	}
	entries := make([]termxml.StringEntry, 0, len(narrations)+len(lines)+len(propers))
	for _, n := range narrations {
		entries = append(entries, termxml.StringEntry{
			EDID: n.EDID, Rec: n.Rec, Field: n.Field,
			Source: n.Source, Dest: n.Dest, Partial: partialOf(n.Status),
		})
	}
	for _, l := range lines {
		entries = append(entries, termxml.StringEntry{
			EDID: l.EDID, Rec: l.Rec, Field: l.Field,
			Source: l.Source, Dest: l.Dest, Partial: partialOf(l.Status),
		})
	}
	for _, pn := range propers {
		entries = append(entries, termxml.StringEntry{
			EDID: pn.EDID, Rec: pn.Rec, Field: pn.Field,
			Source: pn.Source, Dest: pn.Dest, Partial: partialOf(pn.Status),
		})
	}
	return entries, nil
}

// partialOf は DB の訳状態を xTranslator の Partial 属性へ写す。
// 確定訳（statusTranslated=1、権威訳の流用）は translated（属性なし＝空文字）、
// それ以外（AI 仮訳 statusProvisional=3 など）は未完了訳 "1" とする。未訳（0）は書き出し対象外のため来ない。
func partialOf(status int) string {
	if status == statusTranslated {
		return ""
	}
	return "1"
}
