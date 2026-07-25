-- C#↔Go 契約の SQL schema 正本（repo-owned migration）の 6 本目。
-- 本 migration は「会話(INFO:NAM1)・本(BOOK:DESC) 以外の全 REC:FIELD へ翻訳を拡張する」task の
-- 永続層を足す。MECE モデル: プロンプト = Base 指示 ＋ その REC:FIELD に割り当てた指示文（directive）。
-- 論理 ER は docs/er.md。画面の正本は Storybook（Screens/プロンプトテンプレート）。
--
-- 足すテーブル:
--   directive             … 指示文の正本（key→instruction、変数宣言）。口調・文体・固有名・定型句を揃える。
--   record_type_master    … REC:FIELD → box・directive の割り当て正本（翻訳対象だけ載せる）。
--   extracted_field       … C# 素朴吸い出しの受け皿（箱判定を持たない原文バッファ）。
--   proper_noun           … 固有名の訳の単位（方針A・同一実行内の AI 訳。master_term へは昇格しない）。
--   extracted_info_speaker… INFO→speaker の橋渡し（line は Go 取込段が作るため安定キーで一時保持）。

-- directive: 指示文の正本。プロンプト合成の「Base ＋ directive」の directive を 1 行 1 指示文で持つ。
-- key は文体（物品説明/効果説明/世界観断片/書物体/日記体）・固有名・短文（操作名/語義）・口調。
-- 複数 REC:FIELD が 1 directive を共有する。
-- variables は実行時に埋める変数の宣言（JSON 配列 [{token, description}]）。口調だけ {traits} を持つ。
CREATE TABLE IF NOT EXISTS directive (
    key TEXT PRIMARY KEY,
    instruction TEXT NOT NULL,
    variables TEXT NOT NULL DEFAULT '[]'
);

-- 9 指示文を seed する。既に行があれば書き換えない（編集結果を保つ）。
-- 文面は Skyrim.esm の公式日本語既訳を種別ごとに集計して決めた（文末の形・長さ・タグの扱い・句読点）。
-- 粒度は要求される文体で割る。物品説明（物の用途）と効果説明（数値・実行時タグを含む効果）は別の文体を要求し、
-- 操作名（動詞）と語義（語釈）も出力すべき品詞が違うため、それぞれ 1 指示文へ分ける。
-- 口調の instruction は prompt_template.persona_template を畳んだもの（{traits} を持つ）。役割語は戯画的にしない。
INSERT OR IGNORE INTO directive (key, instruction, variables) VALUES
    ('物品説明', 'これは武器・防具・薬・巻物などの品物の説明文、またはゲームの操作説明です。用途と効果を正確に保ち、簡潔に訳すこと。常体（だ・である調）で書き、文は終止形で言い切って句点で終えること。原文の [ ] で囲まれた操作名は訳さずそのまま残すこと。', '[]'),
    ('効果説明', 'これは呪文・付呪・特典・シャウト・魔法効果の説明文です。数値と <> で囲まれた実行時タグは半角のまま残し、増減や順序を変えないこと。タグの前後には半角空白を置くこと。常体で「〜する」「〜を与える」のように終止形で言い切り、句点で終えること。体言止めにしないこと。', '[]'),
    ('世界観断片', 'これはロード画面や種族の解説です。作品世界を語る地の文として、常体で訳し、終止形で言い切って句点で終えること。1 文か 2 文にまとめ、説明を足さないこと。', '[]'),
    ('書物体', 'これは書物の本文です。文章の格調と語り口を保ち、読み物として自然な日本語へ訳すこと。', '[]'),
    ('日記体', 'これはクエストの進行ログ（CNAM）または目標（NNAM）です。ログは主人公の視点の記録として、常体の過去形「〜した。」で書き句点で終えること。目標は今からすることを指す短い句として「〜する」「〜へ行く」の形で書き、句点を付けないこと。', '[]'),
    ('固有名', 'これは名前です。人名・地名・種族名など、原語自体が固有の名はカタカナで音写すること。普通名詞を含む名（勢力名・クエスト名・施設名・役職名など）は意味の分かる日本語へ訳し、音写で済ませないこと。ゲーム UI に収まるよう、原語の文字数を超えないこと。既存の確定訳語があればそれに合わせること。', '[]'),
    ('操作名', 'これは調べる・採取するなどの操作名、またはボタンの文言です。操作を表すものは動詞の終止形（開ける・盗む・取る）、状態や結果を表すものは名詞（完了・鉱山）で訳すこと。2 〜 8 文字に収め、句点を付けないこと。', '[]'),
    ('語義', 'これは龍語など、語そのものの意味を述べる語釈です。語の意味を名詞で短く言い切ること（有限・馬鹿者・定命の者）。2 〜 6 文字に収め、句点を付けず、説明を足さないこと。', '[]'),
    ('口調', 'この台詞の話者の人物像:
{traits}
この人物像に合う口調と人称で訳すこと。
台詞は文末に句点を打たない。原文が疑問符・感嘆符で終わる時だけ ？ ！ を置く。
1 つの台詞に文が 2 つ以上ある時は、句点でつなげず全角空白で区切る。
一人称は漢字で書く（私・俺・僕）。ひらがなの一人称（わたし・おれ・ぼく）は使わない。', '[{"token":"{traits}","description":"話者の性質（生成済みの基底口調と種族訛り）"}]');

-- record_type_master: REC:FIELD → box・directive の割り当て正本。rec/field は Skyrim 仕様の固定行。
-- 翻訳対象だけを載せる（無訳片 WOOP:FULL や header TES4 は載せない。翻訳しない種別は取込段が読み込まない）。
-- logical_name は REC:FIELD だけでは分からない人間向けの種別名（画面のレコード別タブの対象一覧で出す）。
CREATE TABLE IF NOT EXISTS record_type_master (
    rec TEXT NOT NULL,
    field TEXT NOT NULL,
    box TEXT NOT NULL,           -- 叙述文/固有名/定型句/台詞
    directive TEXT NOT NULL REFERENCES directive(key),
    logical_name TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (rec, field)
);

-- 全 translatable REC:FIELD を 1 つの directive へ割り当てる（排他・網羅は engine の seed 整合テストで担保する）。
INSERT OR IGNORE INTO record_type_master (rec, field, box, directive, logical_name) VALUES
    -- 叙述文（物品説明）: 装備・消耗品・巻物・素材・メッセージ本文の説明文。物の用途と特徴を述べる。
    ('WEAP', 'DESC', '叙述文', '物品説明', '武器の説明'),
    ('ARMO', 'DESC', '叙述文', '物品説明', '防具の説明'),
    ('AMMO', 'DESC', '叙述文', '物品説明', '矢弾の説明'),
    ('ALCH', 'DESC', '叙述文', '物品説明', '薬・食料の説明'),
    ('SCRL', 'DESC', '叙述文', '物品説明', '巻物の説明'),
    ('INGR', 'DESC', '叙述文', '物品説明', '錬金素材の説明'),
    ('MESG', 'DESC', '叙述文', '物品説明', 'メッセージ本文'),
    -- 叙述文（効果説明）: 呪文・付呪・特典・シャウト・魔法効果。数値と実行時タグを含む効果の記述。
    ('SPEL', 'DESC', '叙述文', '効果説明', '呪文の説明'),
    ('ENCH', 'DESC', '叙述文', '効果説明', '付呪効果の説明'),
    ('PERK', 'DESC', '叙述文', '効果説明', '特典の説明'),
    ('SHOU', 'DESC', '叙述文', '効果説明', 'シャウトの説明'),
    ('MGEF', 'DNAM', '叙述文', '効果説明', '魔法効果の説明'),
    -- 叙述文（書物体）: 本の本文。
    ('BOOK', 'DESC', '叙述文', '書物体', '本の本文'),
    -- 叙述文（日記体）: クエストログと目標。
    ('QUST', 'CNAM', '叙述文', '日記体', 'クエストログ'),
    ('QUST', 'NNAM', '叙述文', '日記体', 'クエスト目標'),
    -- 叙述文（世界観断片）: ロード画面と種族の解説。作品世界を語る地の文。
    ('LSCR', 'DESC', '叙述文', '世界観断片', 'ロード画面の解説'),
    ('RACE', 'DESC', '叙述文', '世界観断片', '種族の説明'),
    -- 固有名（名称 FULL ほか）: 物・場所・人物・組織の名前。文脈に依存しないため重複排除する。
    ('WEAP', 'FULL', '固有名', '固有名', '武器の名前'),
    ('ARMO', 'FULL', '固有名', '固有名', '防具の名前'),
    ('AMMO', 'FULL', '固有名', '固有名', '矢弾の名前'),
    ('ALCH', 'FULL', '固有名', '固有名', '薬・食料の名前'),
    ('SCRL', 'FULL', '固有名', '固有名', '巻物の名前'),
    ('INGR', 'FULL', '固有名', '固有名', '錬金素材の名前'),
    ('SPEL', 'FULL', '固有名', '固有名', '呪文の名前'),
    ('ENCH', 'FULL', '固有名', '固有名', '付呪の名前'),
    ('MGEF', 'FULL', '固有名', '固有名', '魔法効果の名前'),
    ('PERK', 'FULL', '固有名', '固有名', '特典の名前'),
    ('SHOU', 'FULL', '固有名', '固有名', 'シャウトの名前'),
    ('BOOK', 'FULL', '固有名', '固有名', '本の題名'),
    ('BOOK', 'CNAM', '固有名', '固有名', '本の著者名'),
    ('KEYM', 'FULL', '固有名', '固有名', '鍵の名前'),
    ('MISC', 'FULL', '固有名', '固有名', '雑貨の名前'),
    ('LIGH', 'FULL', '固有名', '固有名', '光源品の名前'),
    ('CONT', 'FULL', '固有名', '固有名', '容器の名前'),
    ('SLGM', 'FULL', '固有名', '固有名', '魂石の名前'),
    ('DOOR', 'FULL', '固有名', '固有名', '扉の名前'),
    ('FURN', 'FULL', '固有名', '固有名', '家具の名前'),
    ('APPA', 'FULL', '固有名', '固有名', '錬金器具の名前'),
    ('HAZD', 'FULL', '固有名', '固有名', '危険物の名前'),
    ('ACTI', 'FULL', '固有名', '固有名', 'オブジェクトの名前'),
    ('FLOR', 'FULL', '固有名', '固有名', '採取植物の名前'),
    ('TREE', 'FULL', '固有名', '固有名', '樹木の名前'),
    ('LCTN', 'FULL', '固有名', '固有名', 'ロケーションの名前'),
    ('WRLD', 'FULL', '固有名', '固有名', 'ワールドの名前'),
    ('CELL', 'FULL', '固有名', '固有名', 'セルの名前'),
    ('QUST', 'FULL', '固有名', '固有名', 'クエスト名'),
    ('RACE', 'FULL', '固有名', '固有名', '種族名'),
    ('FACT', 'FULL', '固有名', '固有名', '勢力名'),
    ('FACT', 'MNAM', '固有名', '固有名', '勢力の階級称号（男性）'),
    ('FACT', 'FNAM', '固有名', '固有名', '勢力の階級称号（女性）'),
    ('MESG', 'FULL', '固有名', '固有名', 'メッセージの題名'),
    ('NPC_', 'FULL', '固有名', '固有名', 'NPC の氏名'),
    ('NPC_', 'SHRT', '固有名', '固有名', 'NPC の短縮名'),
    ('TACT', 'FULL', '固有名', '固有名', '会話オブジェクトの名前'),
    ('SNCT', 'FULL', '固有名', '固有名', '音響カテゴリ名'),
    ('EYES', 'FULL', '固有名', '固有名', '瞳の名前'),
    ('REGN', 'RDMP', '固有名', '固有名', '地域のマップ名'),
    -- 定型句（操作名）: オブジェクト操作とボタン文言。動詞または短い名詞句で、UI の表示幅に収める。
    ('ACTI', 'RNAM', '定型句', '操作名', 'オブジェクトの操作名'),
    ('FLOR', 'RNAM', '定型句', '操作名', '採取植物の操作名'),
    ('TREE', 'RNAM', '定型句', '操作名', '樹木の操作名'),
    ('MESG', 'ITXT', '定型句', '操作名', 'メッセージのボタン'),
    -- 定型句（語義）: 龍語の語釈。語そのものの意味を短く言い切る。
    ('WOOP', 'TNAM', '定型句', '語義', '龍語の語義'),
    -- 台詞: 話者で口調が変わるため重複排除しない。口調 directive の {traits} へ話者の性質を埋める。
    ('INFO', 'NAM1', '台詞', '口調', 'NPC の返答'),
    ('INFO', 'RNAM', '台詞', '口調', '選択肢の条件別上書き'),
    ('DIAL', 'FULL', '台詞', '口調', '選択肢の既定文');

-- extracted_field: C# 抽出器が素朴吸い出しした原文の受け皿。箱・directive の判定を持たない。
-- Go 取込段が record_type_master を引いて narration/proper_noun/line へ振り分ける。
-- レコード識別キーは他の抽出テーブルと同じ (plugin, form_id, rec, field, ordinal)。
CREATE TABLE IF NOT EXISTS extracted_field (
    id INTEGER PRIMARY KEY,
    plugin TEXT NOT NULL,
    form_id TEXT NOT NULL,
    edid TEXT NOT NULL,
    rec TEXT NOT NULL,
    field TEXT NOT NULL,
    ordinal INTEGER NOT NULL DEFAULT 0,
    source TEXT NOT NULL,
    UNIQUE (plugin, form_id, rec, field, ordinal)
);

-- proper_noun: 固有名の訳の単位（方針A）。同一実行内で AI 訳を留め、横断永続辞書 master_term へは昇格しない。
-- source（原語）が本文機械置換の照合キー。category（種別＝rec）は同綴り異義の区別用（concept-model 弱点1）。
CREATE TABLE IF NOT EXISTS proper_noun (
    id INTEGER PRIMARY KEY,
    source TEXT NOT NULL,
    category TEXT NOT NULL DEFAULT '',
    dest TEXT NOT NULL DEFAULT '',
    status INTEGER NOT NULL DEFAULT 0,
    UNIQUE (category, source)
);

-- extracted_info_speaker: INFO（台詞の発生元）→ speaker の橋渡し。
-- line 行は Go 取込段が extracted_field から作るため C# 抽出時には存在しない。
-- そこで C# は安定キー（INFO の plugin・form_id）と upsert 済み speaker.id を一時保持し、
-- Go 取込段が line を作った後に line_speaker（e6）へ解決する。
CREATE TABLE IF NOT EXISTS extracted_info_speaker (
    info_plugin TEXT NOT NULL,
    info_form_id TEXT NOT NULL,
    speaker_id INTEGER NOT NULL REFERENCES speaker(id),
    PRIMARY KEY (info_plugin, info_form_id, speaker_id)
);
