package engine

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/core/mention"
	"aitranslationenginejp/internal/model"
)

// 言及語彙の供給源種別。mention.Term.Kind に入れ、検出結果を narration_mention / line_mention の
// 排他 2 列（master_term_id / proper_noun_id）へ写すときの振り分けに使う。
const (
	mentionKindMasterTerm = "master_term"
	mentionKindProperNoun = "proper_noun"
)

// MentionStore は engine が取込段の言及記録に使う中心データアクセス（使う分だけ宣言する）。
type MentionStore interface {
	ListNarrations(ctx context.Context) ([]model.Narration, error)
	ListLines(ctx context.Context) ([]model.Line, error)
	InsertNarrationMentions(ctx context.Context, rows []model.NarrationMention) (int, error)
	InsertLineMentions(ctx context.Context, rows []model.LineMention) (int, error)
	LinkNarrationDescribed(ctx context.Context) (int, error)
}

// MentionCounts は言及記録で実際に追加した件数（再実行での冪等を観測できるよう件数を返す）。
type MentionCounts struct {
	NarrationMentions int // narration_mention（e4）へ追加した件数
	LineMentions      int // line_mention（e5）へ追加した件数
	DescribedLinks    int // narration_described（e3）へ解決した件数
}

// recordMentions は叙述文・台詞の本文中の既知固有名の言及（e4/e5）と叙述文の説明対象（e3）を記録する。
// 語彙は機械置換辞書と同じ供給源（master_term ∪ proper_noun）を同じ優先順（master_term 先勝ち）で組む。
// 言及は原語（source）の出現で決まる関連のため、訳（dest）の確定状態には依存しない。
// 既存テーブルへは読みだけで、書き込み先は新規の言及テーブルに限る（既存出力の非劣化）。
func (e *Engine) recordMentions(ctx context.Context) (MentionCounts, error) {
	det, err := e.mentionDetector(ctx)
	if err != nil {
		return MentionCounts{}, err
	}

	narrations, err := e.store.ListNarrations(ctx)
	if err != nil {
		return MentionCounts{}, fmt.Errorf("言及検出の叙述文取得: %w", err)
	}
	var narrationMentions []model.NarrationMention
	for _, n := range narrations {
		for _, t := range det.Detect(n.Source) {
			narrationMentions = append(narrationMentions, model.NarrationMention{
				NarrationID:  n.ID,
				ProperNounID: mentionTermID(t, mentionKindProperNoun),
				MasterTermID: mentionTermID(t, mentionKindMasterTerm),
			})
		}
	}
	nm, err := e.store.InsertNarrationMentions(ctx, narrationMentions)
	if err != nil {
		return MentionCounts{}, fmt.Errorf("叙述文言及の投入: %w", err)
	}

	lines, err := e.store.ListLines(ctx)
	if err != nil {
		return MentionCounts{}, fmt.Errorf("言及検出の台詞取得: %w", err)
	}
	var lineMentions []model.LineMention
	for _, l := range lines {
		for _, t := range det.Detect(l.Source) {
			lineMentions = append(lineMentions, model.LineMention{
				LineID:       l.ID,
				ProperNounID: mentionTermID(t, mentionKindProperNoun),
				MasterTermID: mentionTermID(t, mentionKindMasterTerm),
			})
		}
	}
	lm, err := e.store.InsertLineMentions(ctx, lineMentions)
	if err != nil {
		return MentionCounts{}, fmt.Errorf("台詞言及の投入: %w", err)
	}

	// 説明対象（e3）は本文でなくレコード構造（同一レコードの FULL）から決まるため、store の解決 SQL に任せる。
	dl, err := e.store.LinkNarrationDescribed(ctx)
	if err != nil {
		return MentionCounts{}, fmt.Errorf("説明対象の解決: %w", err)
	}
	return MentionCounts{NarrationMentions: nm, LineMentions: lm, DescribedLinks: dl}, nil
}

// mentionDetector は言及検出器を組む。機械置換辞書（LoadDictionary）と同じ供給源・同じ順序
// （master_term → proper_noun、同一原語は先勝ち）で語彙を積み、注入と言及の相手を一致させる。
func (e *Engine) mentionDetector(ctx context.Context) (*mention.Detector, error) {
	terms, err := e.store.ListMasterTerms(ctx)
	if err != nil {
		return nil, fmt.Errorf("言及語彙（master_term）の取得: %w", err)
	}
	propers, err := e.store.ListProperNouns(ctx)
	if err != nil {
		return nil, fmt.Errorf("言及語彙（proper_noun）の取得: %w", err)
	}
	vocab := make([]mention.Term, 0, len(terms)+len(propers))
	for _, t := range terms {
		vocab = append(vocab, mention.Term{Kind: mentionKindMasterTerm, ID: t.ID, Source: t.Source})
	}
	for _, pn := range propers {
		vocab = append(vocab, mention.Term{Kind: mentionKindProperNoun, ID: pn.ID, Source: pn.Source})
	}
	return mention.NewDetector(vocab), nil
}

// mentionTermID は検出した語彙の id を、指す供給源種別の列にだけ写す（他方は 0＝指さない）。
func mentionTermID(t mention.Term, kind string) int64 {
	if t.Kind == kind {
		return t.ID
	}
	return 0
}
