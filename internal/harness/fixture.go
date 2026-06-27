package harness

// Fixture は合成入力一式。著作物を含まない自作データで、抽出子が中心 DB へ書く層（extracted_field と話者素材）と、
// 固有名派生が読む xTranslator XML を 1 つにまとめる。CI 恒久の非劣化回帰網の入力にする。
type Fixture struct {
	PluginName      string        // 抽出対象の plugin 名（fake 抽出子は内容を fixture で固定するため表示用）
	ExtractedFields []SeedField   // C# 抽出器が素朴吸い出しする原文の受け皿（箱判定なし）
	Speakers        []SeedSpeaker // 台詞の話者素材（種族・声型・INFO 橋渡し）。口調生成経路を通すために置く。
	TermsXML        string        // 固有名派生（DeriveMasterTerms）が読む xTranslator 英日 XML の中身
}

// SeedField は extracted_field の 1 行ぶんの seed 値（id は採番させる）。
type SeedField struct {
	Plugin  string
	FormID  string
	EDID    string
	Rec     string
	Field   string
	Ordinal int
	Source  string
}

// SeedSpeaker は 1 話者の seed 値。種族・声型を一緒に作り、INFO（台詞の発生元）→ speaker の橋渡しを残す。
type SeedSpeaker struct {
	Plugin      string
	FormID      string // speaker（base NPC）の form_id
	EDID        string // speaker の EditorID。空だと口調生成の対象から外れるため非空にする。
	Sex         string // Female / Male
	RaceFormID  string
	RaceEDID    string
	RaceNature  string
	VoiceFormID string
	VoiceEDID   string
	VoiceKind   string
	VoiceNature string
	InfoFormID  string // この話者がしゃべる INFO（台詞）の form_id。line と staging を結ぶキー。
}

// SyntheticFixture は著作物を含まない決定的な合成入力を返す。
// 叙述文（説明体・書物体）・固有名・定型句・話者付き台詞・話者なし台詞・翻訳対象外（skip）を 1 通り網羅し、
// 取込段の振り分け、固有名→本文の機械置換、口調生成、プロンプト合成を端から端まで通す。
func SyntheticFixture() Fixture {
	const plugin = "Synthetic.esm"
	return Fixture{
		PluginName: plugin,
		ExtractedFields: []SeedField{
			// 固有名: 本文より先に確定し、機械置換辞書へ載る訳の単位。
			{Plugin: plugin, FormID: "0x100", EDID: "AventusNPC", Rec: "NPC_", Field: "FULL", Source: "Aventus Aretino"},
			// 叙述文（説明体）: 本文中に固有名を含み、機械置換が当たることを観測する。
			{Plugin: plugin, FormID: "0x200", EDID: "TestSword", Rec: "WEAP", Field: "DESC", Source: "A blade once held by Aventus Aretino."},
			// 叙述文（書物体）: 文体 directive が説明体と分かれることを観測する。
			// 末尾に Grelod を含め、二つ名前部（byname）派生語 Grelod→グレロッド が本文機械置換にかかる経路を観測する。
			{Plugin: plugin, FormID: "0x300", EDID: "TestBook", Rec: "BOOK", Field: "DESC", Source: "A long account of the old kings of the north, told by Grelod."},
			// 定型句: オブジェクト操作の短文。固有名でも台詞でもない箱。
			{Plugin: plugin, FormID: "0x400", EDID: "TestActi", Rec: "ACTI", Field: "RNAM", Source: "Open"},
			// 台詞（話者あり・複数話者）: 口調生成と注入を通す。後段で同 INFO に 2 話者を結び、複数話者の解決も観測する。
			{Plugin: plugin, FormID: "0x500", EDID: "GuardGreet", Rec: "INFO", Field: "NAM1", Source: "Have you come to help me with the trouble in town?"},
			// 台詞（話者あり・感情語＋部分形）: 強感情語 fear で emotion_band を、部分形 Aventus で派生→台詞機械置換を観測する。
			{Plugin: plugin, FormID: "0x510", EDID: "GuardWarn", Rec: "INFO", Field: "NAM1", Source: "I fear Aventus will not come back to town."},
			// 台詞（話者なし）: 話者を解決できない台詞は口調指示なしで訳されることを観測する。
			// Grelod を含め、byname 派生語が台詞（line）本文の機械置換にもかかる経路を観測する（叙述文と同じ辞書 Apply）。
			{Plugin: plugin, FormID: "0x600", EDID: "OrphanGreet", Rec: "INFO", Field: "NAM1", Source: "Grelod will not return here, of course."},
			// 翻訳対象外: record_type_master に無い REC:FIELD は取込段が読み込まず skip する。
			{Plugin: plugin, FormID: "0x700", EDID: "Ignored", Rec: "ZZZZ", Field: "FULL", Source: "should be skipped"},
		},
		// 各話者は声型を 1 つだけ持つ（speaker.voice_type_id は 1 対 1 の FK）。複数声型の話者は作らない。
		Speakers: []SeedSpeaker{
			// 0x500 の 1 人目。seed 順で speaker.id が一番若く、複数話者の「先頭話者（id 昇順）」採用を観測する。
			// form_id を 0x090 にして、insert 順（id 昇順で TownGuard が先）と form_id 辞書順（MarketWoman 0x060 が先）を
			// わざと食い違わせる。これで「先頭採用キーを s.id から s.form_id へ変える」誤リファクタも golden 差分で検出できる。
			{
				Plugin: plugin, FormID: "0x090", EDID: "TownGuard", Sex: "Male",
				RaceFormID: "0x010", RaceEDID: "NordRace", RaceNature: "頑健で実直",
				VoiceFormID: "0x020", VoiceEDID: "MaleNord", VoiceKind: "成人男性", VoiceNature: "落ち着いた低い声",
				InfoFormID: "0x500",
			},
			// 0x500 の 2 人目。同一 INFO に結び、line_speaker が 2 行になる複数話者を観測する。
			// 声型を Condescending（対人 prior 尊大）にして 1 人目（中立）と persona を変え、台詞へ注入される口調が
			// 「先頭話者（id 昇順）」の TownGuard 由来であることを golden に出す（先頭採用の逆転を検出可能にする）。
			{
				Plugin: plugin, FormID: "0x060", EDID: "MarketWoman", Sex: "Female",
				RaceFormID: "0x011", RaceEDID: "ImperialRace", RaceNature: "如才ない",
				VoiceFormID: "0x021", VoiceEDID: "FemaleCondescending", VoiceKind: "成人女性", VoiceNature: "見下す調子",
				InfoFormID: "0x500",
			},
			// 0x510（感情語台詞）の話者。
			{
				Plugin: plugin, FormID: "0x070", EDID: "GateWatch", Sex: "Male",
				RaceFormID: "0x010", RaceEDID: "NordRace", RaceNature: "頑健で実直",
				VoiceFormID: "0x020", VoiceEDID: "MaleNord", VoiceKind: "成人男性", VoiceNature: "落ち着いた低い声",
				InfoFormID: "0x510",
			},
		},
		// xTranslator 英日 XML。固有名派生（DeriveMasterTerms）が NPC_:FULL から名の部分形を派生して master_term へ載せる。
		// Aventus Aretino は姓名分割（two）、Grelod the Kind は二つ名前部（byname）の派生経路を観測するために置く。
		TermsXML: `<?xml version="1.0" encoding="utf-8"?>
<SSTXMLRessources>
  <Content>
    <String><REC>NPC_:FULL</REC><Source>Aventus Aretino</Source><Dest>アベンタス・アレティノ</Dest></String>
    <String><REC>NPC_:FULL</REC><Source>Grelod the Kind</Source><Dest>親切者のグレロッド</Dest></String>
  </Content>
</SSTXMLRessources>
`,
	}
}
