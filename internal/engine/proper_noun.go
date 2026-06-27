package engine

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"
	"aitranslationenginejp/internal/provider"
)

// statusTranslated は xTranslator の訳状態 1（訳済）。固有名の既訳流用（権威訳）は確定として書き戻す。
const statusTranslated = 1

// 固有名の訳の書き込み宛先テーブル名。供給源選別の不変ルールを名前で固定する。
const (
	properNounTable = "proper_noun" // AI 訳・既訳流用ともに固有名はここへ書く（実行内）。
	masterTermTable = "master_term" // 権威訳の横断永続辞書。AI 訳をここへは書かない（方針A の不変境界）。
)

// SupplyKind は固有名 1 件の供給源。既訳流用（権威訳）か、AI 訳か。
type SupplyKind int

const (
	// SupplyAuthoritative は master_term に既訳がある場合。AI 翻訳せず権威訳を使う。
	SupplyAuthoritative SupplyKind = iota
	// SupplyAITranslate は master_term に既訳が無い場合。AI 翻訳して proper_noun へ書く。
	SupplyAITranslate
)

// TermSupply は固有名 1 件の供給源判定（供給源選別の出力）。
// Dest は既訳流用のときの確定訳語。WriteTarget は AI 訳の書き込み宛先で、
// 不変ルール: AI 訳の宛先は常に proper_noun で、master_term にしない（方針A）。
type TermSupply struct {
	Kind        SupplyKind
	Dest        string
	WriteTarget string
}

// SelectSupply は固有名の既訳の有無から供給源を選ぶ純粋関数（不変ルールを 1 関数に閉じる）。
// 入力＝固有名の原語と既訳辞書（master_term の source→dest）。
// 出力＝既訳ありなら権威訳（Dest を埋め書き込み不要）、既訳なしなら AI 訳（宛先は proper_noun）。
// 不変ルール: AI 訳の WriteTarget は proper_noun で、master_term を返さない。副作用を持たない。
func SelectSupply(source string, authoritative map[string]string) TermSupply {
	if dest, ok := authoritative[source]; ok {
		return TermSupply{Kind: SupplyAuthoritative, Dest: dest}
	}
	return TermSupply{Kind: SupplyAITranslate, WriteTarget: properNounTable}
}

// ProperNounStore は engine が固有名フェーズに使う中心データアクセス（使う分だけ宣言する）。
type ProperNounStore interface {
	ListUntranslatedProperNouns(ctx context.Context) ([]model.ProperNoun, error)
	UpdateProperNounDest(ctx context.Context, id int64, dest string, status int) error
}

// translateProperNouns は本文フェーズより前に、与えられた未訳固有名 pending を確定する。
// 各固有名を SelectSupply で振り分け、既訳ありは権威訳を、既訳なしは AI 訳を proper_noun へ書く。
// authoritative は master_term の source→dest、base は base 指示、instruction は固有名 directive の指示文。
// onProcessed は固有名 1 件を処理し終えるたびに呼ぶ進捗通知（既訳流用・AI 訳の両方で呼ぶ）。
func (e *Engine) translateProperNouns(ctx context.Context, conn provider.Connection, model string,
	pending []model.ProperNoun, authoritative map[string]string, base, instruction string, onProcessed func()) error {
	for _, pn := range pending {
		sup := SelectSupply(pn.Source, authoritative)
		switch sup.Kind {
		case SupplyAuthoritative:
			// 既訳流用。AI 翻訳せず権威訳を proper_noun へ確定として書く（master_term は変えない）。
			if err := e.store.UpdateProperNounDest(ctx, pn.ID, sup.Dest, statusTranslated); err != nil {
				return fmt.Errorf("固有名（既訳流用）の書き戻し: %w", err)
			}
		default:
			// 既訳なし。固有名 directive で AI 翻訳し、仮訳として proper_noun へ書く。
			dest, err := e.provider.Translate(ctx, conn, model, ComposePrompt(base, instruction, pn.Source))
			if err != nil {
				return fmt.Errorf("固有名の翻訳: %w", err)
			}
			if err := e.store.UpdateProperNounDest(ctx, pn.ID, dest, statusProvisional); err != nil {
				return fmt.Errorf("固有名（AI 訳）の書き戻し: %w", err)
			}
		}
		if onProcessed != nil {
			onProcessed()
		}
	}
	return nil
}
