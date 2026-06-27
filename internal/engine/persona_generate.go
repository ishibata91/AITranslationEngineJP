package engine

import (
	"context"
	"fmt"
	"sort"

	"aitranslationenginejp/internal/engine/tone"
	"aitranslationenginejp/internal/model"
)

// PersonaGenStore は一括ペルソナ生成が使う中心データアクセス（使う分だけ宣言する）。
type PersonaGenStore interface {
	ListSpeakerLineSources(ctx context.Context) ([]model.SpeakerLineSource, error)
	GetLineAnalyses(ctx context.Context, hashes []string) (map[string]model.LineAnalysis, error)
	UpsertLineAnalysis(ctx context.Context, a model.LineAnalysis) error
	UpsertPersonaCharacter(ctx context.Context, pc model.PersonaCharacter) error
}

// PersonaGenerator は対話由来の基底口調を一括生成する。最も重い行解析は line_analysis に本文ハッシュで
// キャッシュし、本文ごとに 1 度だけ prose を回す。生成は集計するだけで安価。手修正は store 側で保護する。
type PersonaGenerator struct {
	store      PersonaGenStore
	lexicon    EmotionLexicon
	classifier *tone.Classifier
}

// NewPersonaGenerator は生成器を組む。classifier は不変ルール、lexicon は差し替え可能な感情辞書。
func NewPersonaGenerator(store PersonaGenStore, lexicon EmotionLexicon) *PersonaGenerator {
	return &PersonaGenerator{store: store, lexicon: lexicon, classifier: tone.NewClassifier()}
}

// Generate は台詞を持つ全話者の基底口調を一括生成し、persona_character へ保存する。保存した話者数を返す。
// 行解析は本文ハッシュでキャッシュし、未命中の本文だけ prose で解析する。手修正済みは store が保護する。
func (g *PersonaGenerator) Generate(ctx context.Context) (int, error) {
	rows, err := g.store.ListSpeakerLineSources(ctx)
	if err != nil {
		return 0, fmt.Errorf("生成入力の取得: %w", err)
	}
	features, err := g.ensureLineAnalyses(ctx, rows)
	if err != nil {
		return 0, err
	}

	// rows は (plugin, form_id) でソート済みのため、隣接する同一話者を 1 まとまりとして束ねる。
	count := 0
	for i := 0; i < len(rows); {
		j := i
		var feats []tone.Features
		voice := ""
		for j < len(rows) &&
			rows[j].SpeakerPlugin == rows[i].SpeakerPlugin &&
			rows[j].SpeakerFormID == rows[i].SpeakerFormID {
			feats = append(feats, features[SourceHash(rows[j].Source)])
			if voice == "" {
				voice = rows[j].VoiceEDID
			}
			j++
		}
		persona := g.classifier.Classify(feats, voice)
		if err := g.store.UpsertPersonaCharacter(ctx, model.PersonaCharacter{
			SpeakerPlugin: rows[i].SpeakerPlugin,
			SpeakerFormID: rows[i].SpeakerFormID,
			AttitudeBand:  persona.AttitudeBand,
			EmotionBand:   persona.EmotionBand,
			Marked:        persona.Marked,
			DecisionPath:  persona.DecisionPath,
		}); err != nil {
			return count, err
		}
		count++
		i = j
	}
	return count, nil
}

// ensureLineAnalyses は rows の本文を、キャッシュに無い分だけ prose で解析して保存し、
// source_hash→Features の map を返す。重複本文は 1 度だけ解析する（プール台詞も 1 回で済む）。
func (g *PersonaGenerator) ensureLineAnalyses(ctx context.Context, rows []model.SpeakerLineSource) (map[string]tone.Features, error) {
	hashToText := make(map[string]string)
	for _, r := range rows {
		hashToText[SourceHash(r.Source)] = r.Source
	}
	// hash 群は決定的な順序にする。map 反復のままだと UpsertLineAnalysis の書き込み順が実行ごとに揺れ、
	// line_analysis の採番 id が非決定になる。本文ハッシュの昇順で固定する。
	hashes := make([]string, 0, len(hashToText))
	for h := range hashToText {
		hashes = append(hashes, h)
	}
	sort.Strings(hashes)
	cached, err := g.store.GetLineAnalyses(ctx, hashes)
	if err != nil {
		return nil, fmt.Errorf("行解析キャッシュの取得: %w", err)
	}
	out := make(map[string]tone.Features, len(hashes))
	// 並びを固定した hashes でキャッシュ未命中の本文を処理し、書き込み順を決定的にする。
	for _, h := range hashes {
		text := hashToText[h]
		if row, ok := cached[h]; ok {
			out[h] = featuresFromAnalysis(row)
			continue
		}
		f := ExtractFeatures(text, g.lexicon) // prose（重い）。キャッシュ未命中の本文だけ実行する。
		if err := g.store.UpsertLineAnalysis(ctx, analysisFromFeatures(h, f)); err != nil {
			return nil, fmt.Errorf("行解析キャッシュの保存: %w", err)
		}
		out[h] = f
	}
	return out, nil
}

// featuresFromAnalysis はキャッシュ行を Classifier 入力へ写す。
func featuresFromAnalysis(a model.LineAnalysis) tone.Features {
	return tone.Features{
		Sentences:  a.SentenceCount,
		Polite:     a.PoliteCount,
		Insult:     a.InsultCount,
		Imperative: a.IsImperative,
		Exclaim:    a.ExclaimCount,
		Elong:      a.ElongCount,
		Emotion:    a.EmotionCount,
	}
}

// analysisFromFeatures は抽出特徴をキャッシュ行へ写す。
func analysisFromFeatures(hash string, f tone.Features) model.LineAnalysis {
	return model.LineAnalysis{
		SourceHash:    hash,
		SentenceCount: f.Sentences,
		PoliteCount:   f.Polite,
		InsultCount:   f.Insult,
		IsImperative:  f.Imperative,
		ExclaimCount:  f.Exclaim,
		ElongCount:    f.Elong,
		EmotionCount:  f.Emotion,
	}
}
