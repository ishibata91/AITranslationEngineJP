package harness

// Fixture は合成入力一式。著作物を含まない自作データで、抽出子が中心 DB へ書く層
// （既訳・翻訳対象の原文・話者素材）を 1 つにまとめる。CI 恒久の非劣化回帰網の入力にする。
// 器は 2 つに分かれる。References は Data フォルダ全 plugin の英日対（既訳）、
// ExtractedFields は利用者が選んだ 1 plugin の翻訳対象の原文で、日本語を持たない。
type Fixture struct {
	PluginName      string          // 抽出対象の plugin 名（fake 抽出子は内容を fixture で固定するため表示用）
	References      []SeedReference // Data フォルダ全 plugin から集めた英日対（既訳）。横断辞書と完全一致置換の供給源。
	ExtractedFields []SeedField     // C# 抽出器が素朴吸い出しする翻訳対象の原文（箱判定なし）
	Speakers        []SeedSpeaker   // 台詞の話者素材（種族・声型・INFO 橋渡し）。口調生成経路を通すために置く。
	Emotions        []SeedEmotion   // 台詞の感情型 staging（INFO 応答→TRDT）。感情がプロンプトへ乗る合流を通すために置く。
}

// SeedReference は reference_translation の 1 行ぶんの seed 値（既訳の英日対）。
// 照合キーは (Rec, Field, Source) で、どの plugin 由来かは持たない（既訳は plugin をまたいで再利用する）。
type SeedReference struct {
	Rec    string
	Field  string
	Source string
	Dest   string
}

// SeedEmotion は extracted_info_emotion の 1 行（INFO 応答単位の感情型）。取込段が (plugin, form_id, ordinal) で line へ結ぶ。
type SeedEmotion struct {
	Plugin      string
	InfoFormID  string // 感情を持つ INFO（台詞）の form_id
	Ordinal     int    // 応答順（line の非空応答の出現順と一致させる）
	EmotionType string // Anger/Disgust/Fear/Sad/Happy/Surprise/Puzzled（Neutral は置かない）
}

// SeedField は extracted_field の 1 行ぶんの seed 値（id は採番させる）。翻訳対象の原文だけを持ち、日本語を持たない。
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
// 翻訳対象 `Synthetic.esm` は日本語 Strings を持たない mod を模し、既訳は同じ Data フォルダの
// 別 plugin（日本語 Strings を持つ公式相当）から集めた References が担う。
// 叙述文（物品説明・書物体）・固有名・定型句・話者付き台詞・話者なし台詞・翻訳対象外（skip）を 1 通り網羅し、
// 取込段の振り分け、固有名→本文の機械置換、口調生成、プロンプト合成を端から端まで通す。
func SyntheticFixture() Fixture {
	const plugin = "Synthetic.esm"
	return Fixture{
		PluginName: plugin,
		// 既訳: Data フォルダにある日本語 Strings 付き plugin から集めた英日対。
		// 翻訳対象 Synthetic.esm 自身は日本語を持たないので、供給はすべて別 plugin 由来になる。
		References: []SeedReference{
			// 固有名の英日対。横断辞書の完全形（FULL）と、姓名分割（two）・二つ名前部（byname）派生の供給源。
			{Rec: "NPC_", Field: "FULL", Source: "Aventus Aretino", Dest: "アベンタス・アレティノ"},
			{Rec: "NPC_", Field: "FULL", Source: "Grelod the Kind", Dest: "親切者のグレロッド"},
			// 台詞の英日対。原文が完全一致する台詞は AI を呼ばず既訳を流用する（known-issues 項目7）。
			{Rec: "INFO", Field: "NAM1", Source: "Well met, traveler.", Dest: "ようこそ、旅の方。"},
		},
		ExtractedFields: []SeedField{
			// 固有名: 本文より先に確定し、機械置換辞書へ載る訳の単位。既訳（References）に同じ原語があるため
			// 固有名フェーズは AI を呼ばず横断辞書の訳を流用する。
			{Plugin: plugin, FormID: "0x100", EDID: "AventusNPC", Rec: "NPC_", Field: "FULL", Source: "Aventus Aretino"},
			// 固有名（二つ名）: 既訳あり。二つ名前部（byname）派生 Grelod→グレロッド の供給源になる。
			{Plugin: plugin, FormID: "0x110", EDID: "GrelodNPC", Rec: "NPC_", Field: "FULL", Source: "Grelod the Kind"},
			// 固有名（mod 追加・既訳なし・空白 2 語）: 固有名フェーズで AI 訳が確定し、そこから姓名分割で
			// 名のみ（Sorine）・苗字のみ（Trueblade）の部分形が派生する経路を観測する。
			{Plugin: plugin, FormID: "0x120", EDID: "SorineNPC", Rec: "NPC_", Field: "FULL", Source: "Sorine Trueblade"},
			// 固有名（mod 追加・既訳なし・二つ名）: AI 訳から二つ名前部（Ulfarr）が派生する経路を観測する。
			{Plugin: plugin, FormID: "0x130", EDID: "UlfarrNPC", Rec: "NPC_", Field: "FULL", Source: "Ulfarr the Grim"},
			// 固有名（mod 追加・既訳なし・名が一般語）: 名 Hunter は台詞に小文字で 3 回出るため用法分布が
			// 一般語と判定し、部分形として辞書へ載らない（破壊置換の防止）。苗字側は派生する。
			{Plugin: plugin, FormID: "0x140", EDID: "HunterNPC", Rec: "NPC_", Field: "FULL", Source: "Hunter Greycloak"},
			// 叙述文（物品説明）: 本文中に固有名を含み、機械置換が当たることを観測する。
			{Plugin: plugin, FormID: "0x200", EDID: "TestSword", Rec: "WEAP", Field: "DESC", Source: "A blade once held by Aventus Aretino."},
			// 叙述文（書物体）: 文体 directive が物品説明と分かれることを観測する。
			// 末尾に Grelod を含め、二つ名前部（byname）派生語 Grelod→グレロッド が本文機械置換にかかる経路を観測する。
			{Plugin: plugin, FormID: "0x300", EDID: "TestBook", Rec: "BOOK", Field: "DESC", Source: "A long account of the old kings of the north, told by Grelod."},
			// 定型句: オブジェクト操作の短文。固有名でも台詞でもない箱。
			{Plugin: plugin, FormID: "0x400", EDID: "TestActi", Rec: "ACTI", Field: "RNAM", Source: "Open"},
			// 台詞（話者あり・複数話者）: 口調生成と注入を通す。後段で同 INFO に 2 話者を結び、複数話者の解決も観測する。
			{Plugin: plugin, FormID: "0x500", EDID: "GuardGreet", Rec: "INFO", Field: "NAM1", Source: "Have you come to help me with the trouble in town?"},
			// 台詞（話者あり・感情語＋部分形）: 強感情語 fear で emotion_band を、部分形 Aventus で派生→台詞機械置換を観測する。
			// TRDT 感情型（Fear）を Emotions で付け、感情がこの台詞のプロンプトへ乗る合流（line-emotion-injected）を通す。
			{Plugin: plugin, FormID: "0x510", EDID: "GuardWarn", Rec: "INFO", Field: "NAM1", Source: "I fear Aventus will not come back to town."},
			// 台詞（フルネームの固有名を含む）: 叙述文（WEAP:DESC）と同じ固有名 Aventus Aretino を本文に持ち、
			// 抽出した固有名が叙述文と台詞で同一訳になる合流（proper-noun-consistent）を観測する。
			{Plugin: plugin, FormID: "0x520", EDID: "GuardAsk", Rec: "INFO", Field: "NAM1", Source: "Have you seen Aventus Aretino today?"},
			// 台詞（話者なし）: 話者を解決できない台詞は口調指示なしで訳されることを観測する。
			// Grelod を含め、byname 派生語が台詞（line）本文の機械置換にもかかる経路を観測する（叙述文と同じ辞書 Apply）。
			{Plugin: plugin, FormID: "0x600", EDID: "OrphanGreet", Rec: "INFO", Field: "NAM1", Source: "Grelod will not return here, of course."},
			// 台詞（mod NPC の名のみ・苗字のみ）: 実行内で確定した氏名から派生した部分形が本文へ当たることを観測する。
			{Plugin: plugin, FormID: "0x610", EDID: "GateTalk", Rec: "INFO", Field: "NAM1", Source: "Sorine has left, and Trueblade now guards the gate."},
			// 台詞（mod NPC の二つ名の前部）: 二つ名を除いた名だけが本文へ当たることを観測する。
			{Plugin: plugin, FormID: "0x620", EDID: "ShrineTalk", Rec: "INFO", Field: "NAM1", Source: "Ulfarr keeps the old shrine."},
			// 台詞（一般語 hunter の小文字出現 3 件）: 用法分布の LC を積み、名 Hunter を一般語と判定させる。
			{Plugin: plugin, FormID: "0x630", EDID: "RidgeTalk", Rec: "INFO", Field: "NAM1", Source: "I saw a hunter near the ridge."},
			{Plugin: plugin, FormID: "0x631", EDID: "DawnTalk", Rec: "INFO", Field: "NAM1", Source: "The hunter went north at dawn."},
			{Plugin: plugin, FormID: "0x632", EDID: "RoadTalk", Rec: "INFO", Field: "NAM1", Source: "Ask the hunter about the road."},
			// 台詞（文頭の大文字 Hunter）: 一般語が文頭で大文字になる形。部分形が辞書へ載っていれば
			// ここが壊れる（実データで観測した破壊置換の形）。原文のまま残ることを観測する。
			{Plugin: plugin, FormID: "0x640", EDID: "TrackTalk", Rec: "INFO", Field: "NAM1", Source: "Hunter tracks are fresh on the road."},
			// 翻訳対象外: record_type_master に無い REC:FIELD は取込段が読み込まず skip する。
			{Plugin: plugin, FormID: "0x700", EDID: "Ignored", Rec: "ZZZZ", Field: "FULL", Source: "should be skipped"},
			// 一般語と同綴りの管理用文字列: Mod 制作では勢力の階級称号（FACT:MNAM）を対話状態の
			// 内部フラグに使う慣行があり（実例 inigo.esp）、固有名 box へ Yes/No が入る。
			// 固有名フェーズで AI 訳は付くが、stoplist が供給を選別し置換も言及も起きないことを golden で凍結する。
			{Plugin: plugin, FormID: "0x800", EDID: "OptionFaction", Rec: "FACT", Field: "MNAM", Ordinal: 0, Source: "Yes"},
			{Plugin: plugin, FormID: "0x800", EDID: "OptionFaction", Rec: "FACT", Field: "MNAM", Ordinal: 1, Source: "No"},
			// 台詞（話者なし・文頭 Yes）: stoplist 語が本文で置換されないことをプロンプトで観測する。
			{Plugin: plugin, FormID: "0x820", EDID: "ScoutReply", Rec: "INFO", Field: "NAM1", Source: "Yes, the road is clear now."},
			// 叙述文（文頭 No）: 実データで観測した誤爆（No matter → 訳NNN matter）の形を再現し、
			// stoplist 選別後は原文のままプロンプトへ渡ることを観測する。
			{Plugin: plugin, FormID: "0x830", EDID: "TestDagger", Rec: "WEAP", Field: "DESC", Source: "No matter what the weather, this blade holds its edge."},
			// 叙述文（実行時タグ入り）: 本文に <Alias=...> を含む。退避→送信→復元でタグが原形保持されることを観測する
			// （known-issues 項目8）。既訳を置かないため、AI 翻訳経路（退避）を通る。
			{Plugin: plugin, FormID: "0x900", EDID: "TagBook", Rec: "BOOK", Field: "DESC", Source: "Deliver this letter to <Alias=Player> at once."},
			// 台詞（既訳と完全一致）: References に同じ原文の英日対があり、参照訳経由で AI を呼ばず
			// 既訳が確定訳で流用されることを観測する（known-issues 項目7）。
			{Plugin: plugin, FormID: "0x910", EDID: "KnownGreet", Rec: "INFO", Field: "NAM1", Source: "Well met, traveler."},
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
		// 台詞の感情型 staging。GuardWarn（0x510）の応答 0 に Fear を付け、感情がプロンプトへ乗る合流を通す。
		Emotions: []SeedEmotion{
			{Plugin: plugin, InfoFormID: "0x510", Ordinal: 0, EmotionType: "Fear"},
		},
	}
}

// SyntheticAITranslations は合成入力で fake provider が返す固定訳。
// 実行内で確定した氏名から人名の部分形が派生する経路を通すため、mod NPC の氏名には
// 中黒区切りのカタカナ訳（姓名分割が成立する形）と、二つ名を含む訳を与える。
// ここに無い原文は連番訳（訳NNN）へ落ちる。
func SyntheticAITranslations() map[string]string {
	return map[string]string{
		"Sorine Trueblade": "ソリーヌ・トゥルーブレイド",
		"Ulfarr the Grim":  "厳格なウルファール",
		"Hunter Greycloak": "ハンター・グレイクローク",
	}
}
