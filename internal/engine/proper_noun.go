package engine

import (
	"context"
	"fmt"
	"strings"

	"aitranslationenginejp/internal/core/prompt"
	"aitranslationenginejp/internal/core/termderive"
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
	ListUntranslatedProperNouns(ctx context.Context, plugin string) ([]model.ProperNoun, error)
	UpdateProperNounDest(ctx context.Context, id int64, dest string, status int) error
	ListConfirmedNPCNames(ctx context.Context, plugin string) ([]model.ConfirmedName, error)
	InsertDerivedProperNouns(ctx context.Context, rows []model.ProperNoun) (int, error)
}

// deriveRunProperNouns は実行内で訳が確定した対象 plugin の NPC 名から人名の部分形（名のみ・苗字のみ・短名）を
// 作り、同じ plugin の proper_noun へ書く。追記した件数を返す。
// 固有名フェーズの後・機械置換辞書を組む前に呼ぶ。同期翻訳と batch の両方が本文を組む前に通す。
// mod が追加した NPC は既訳を持たないため実行内の AI 訳でしか氏名が確定せず、抽出直後に走る
// DeriveMasterTerms（既訳が入力）では部分形を作れない。この段がその穴を埋める。
// 書き込み先は proper_noun に固定する。横断永続辞書 master_term へは昇格させない（方針A の不変境界）。
// 既出原語（横断辞書 ∪ その実行で確定済みの固有名）は派生に含めないため、同じ原語へ 2 つの訳が立たない。
// 派生の採否は termderive の語形の防御が決める（用法分布の材料は翻訳対象 plugin の台詞の英語原文）。
func (e *Engine) deriveRunProperNouns(ctx context.Context, plugin string) (int, error) {
	names, err := e.store.ListConfirmedNPCNames(ctx, plugin)
	if err != nil {
		return 0, fmt.Errorf("確定した NPC 名の取得: %w", err)
	}
	if len(names) == 0 {
		return 0, nil
	}
	var fullPairs, shrtPairs []termderive.NamePair
	for _, n := range names {
		pair := termderive.NamePair{Source: strings.TrimSpace(n.Source), Dest: strings.TrimSpace(n.Dest)}
		if pair.Source == "" || pair.Dest == "" {
			continue
		}
		switch n.Field {
		case fieldFull:
			fullPairs = append(fullPairs, pair)
		case fieldShrt:
			shrtPairs = append(shrtPairs, pair)
		}
	}

	baseSources, err := e.derivedBaseSources(ctx, names)
	if err != nil {
		return 0, err
	}
	usage, err := e.dialogueUsage(ctx)
	if err != nil {
		return 0, err
	}
	derived := termderive.DeriveTerms(fullPairs, shrtPairs, usage, baseSources, termderive.DefaultDeriveConfig())
	rows := make([]model.ProperNoun, len(derived))
	for i, d := range derived {
		// origin で翻訳対象と分ける。category は空にする。派生行は原文 record を持たないので rec の値が無い。
		// 空にすることで原文位置の解決（言及・出力位置。category を rec と突き合わせる）が派生行を拾わない。
		// UNIQUE(plugin, category, source) は source が抽出由来と重ならないため衝突しない
		// （既出原語は derivedBaseSources が派生から外す）。
		rows[i] = model.ProperNoun{
			Plugin: plugin, Source: d.Source, Dest: d.Dest,
			Status: statusProvisional, Origin: model.OriginDerived,
		}
	}
	inserted, err := e.store.InsertDerivedProperNouns(ctx, rows)
	if err != nil {
		return inserted, fmt.Errorf("派生した人名の部分形の追記: %w", err)
	}
	return inserted, nil
}

// derivedBaseSources は部分形の派生から外す既出原語の集合を作る。
// 横断辞書（master_term）の原語と、その実行で確定済みの固有名の原語の両方を入れ、同じ原語へ 2 つの訳が立たないようにする。
func (e *Engine) derivedBaseSources(ctx context.Context, names []model.ConfirmedName) (map[string]bool, error) {
	terms, err := e.store.ListMasterTerms(ctx)
	if err != nil {
		return nil, fmt.Errorf("横断辞書の取得: %w", err)
	}
	base := make(map[string]bool, len(terms)+len(names))
	for _, t := range terms {
		base[t.Source] = true
	}
	for _, n := range names {
		base[strings.TrimSpace(n.Source)] = true
	}
	return base, nil
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
			dest, err := e.provider.Translate(ctx, conn, model, prompt.ComposePrompt(base, instruction, pn.Source))
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
