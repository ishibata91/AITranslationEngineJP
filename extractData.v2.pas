{ ========================================================================
   xEdit Data Extraction Script v2
   docs/conceptual-model.md (v17) に準拠して record を抽出する。
   v1（extractData.pas）からの主な変更:
     - class 分割（Item/Equipment/Consumable/Book、Magic/Enchantment/MagicEffect/Shout/Word）
     - Person/Race/Faction を分離 class として top-level に持つ
     - Dialogue.category / subtype（用途種別 / 発火種別 = DATA Category/Subtype byte）
     - Response.conditions（条件一覧 = INFO.CTDA list）
     - Response.next_ids（後続応答 = INFO.TCLT forward-link）
     - 各 carrier の Effects[] を解決して magic_effect_ids 配列を emit
     - Equipment.enchantment_id（EITM）と Shout.words[].word_id（WNAM）を emit
     - cross-file traversal: 参照解決時に Ensure* 経由で master record を自動取り込み
   速度修正は v1 検証で採用したものを全て継承:
     - dialDataMap.Sorted := True（dial lookup を O(log N) 化）
     - EscapeJSONString を run 単位の Copy 追記で実用上 O(N) 化（JvInterpreter は
       WideChar 型と文字列の index 書き込みを受け付けないため index 代入は使わない）
     - response 行に id/pid を prefix 保持し JSON 再 parse を撤廃
     - Finalize の出力構築を TStringList 行追加 + SaveToFile で O(N) 化
     - 各 Extract* で EDID / source を local キャッシュ
   出力 JSON schema は scripts/validate_extraction.py が検証する。
   ======================================================================== }

unit ExportDataV2;

var
  // 概念モデル class ごとの emit 用 buffer。
  dialogueList: TStringList;       // Dialogue
  personList: TStringList;         // Person（NPC_）、id keyed object
  raceList: TStringList;           // Race
  factionList: TStringList;        // Faction
  questList: TStringList;          // Quest（embedded stages/objectives）
  locationList: TStringList;       // Location（CELL/LCTN/WRLD）
  itemList: TStringList;           // Item（KEYM/MISC/LIGH/CONT/SLGM/DOOR/FURN）
  activatorList: TStringList;      // Activator（ACTI/FLOR/TREE、起動・採取可能オブジェクト）
  equipmentList: TStringList;      // Equipment（WEAP/ARMO/AMMO、+enchantment_id）
  consumableList: TStringList;     // Consumable（SCRL/ALCH/INGR、+magic_effect_ids）
  magicList: TStringList;          // Magic（SPEL、+magic_effect_ids）
  enchantmentList: TStringList;    // Enchantment（ENCH、+magic_effect_ids）
  magicEffectList: TStringList;    // MagicEffect（MGEF）
  voiceTypeList: TStringList;      // VoiceType（VTYP、v18 で独立 class 化）
  shoutList: TStringList;          // Shout（SHOU、embedded words）
  bookList: TStringList;           // Book（BOOK、本文を body field で保持）
  messageList: TStringList;        // Message
  loadingScreenList: TStringList;  // LoadingScreen
  perkList: TStringList;           // Perk

  // 1 つの Dialog に紐づく header + response 行。
  // Strings[0]    = Dialogue header JSON（body + responses[ まで）
  // Strings[1..N] = "id#1pid#1responseJSON"（#1 = SEP_CHAR）
  dialDataMap: TStringList;

  // 全 record 共通の処理済み id set。Ensure* が二重抽出を避ける。
  processedIds: TStringList;

  // State
  targetFileName: string;

  // Voice resolution の再帰防止（v1 から継承）
  lstRecursion: TList;
  InfoNPCID, InfoSPEAKER, InfoRACEID, InfoCONDITION: string;
  const fDebug = False;
  const SEP_CHAR = #1;

// ===== UTILITY =====

function HexFormID(elem: IInterface): string;
begin
  if Assigned(elem) then
    Result := IntToHex(FixedFormID(elem), 8)
  else
    Result := '00000000';
end;

function GetElementValue(elem: IInterface; path: string): string;
begin
  Result := '';
  if Assigned(elem) then
    Result := GetElementEditValues(elem, path);
end;

// JSON 文字列 escape を実行する。
// 文字コードの取得は v1（extractData.pas）で実績のある経路に合わせる。
// JvInterpreter は s[i] の戻り値（1 文字 string）を直接 Ord に渡すと Type mismatch に
// なる。Copy(s,i,1) で 1 文字 string を取り、その [1] を Char 変数へ代入してから Ord を
// 呼ぶと、Char として扱われ Ord が通る。WideChar 型と文字列の index 書き込みは使わない。
// 非 ASCII は \uXXXX へ escape し、出力は純 ASCII を維持する（v1 出力と同じ契約。
// app 側が読む既存 dictionaries は BOM 無し純 ASCII のため encoding を変えない）。
// escape 不要文字の連続区間は Copy で切り出してまとめて追記し、per-char 連結を避ける。
function EscapeJSONString(const s: string): string;
var
  i, len, runStart, c: Integer;
  ch: Char;
  chStr: string;
begin
  Result := '';
  len := Length(s);
  if len = 0 then Exit;
  runStart := 1;
  for i := 1 to len do begin
    chStr := Copy(s, i, 1);
    if Length(chStr) > 0 then begin
      ch := chStr[1];
      c := Ord(ch);
    end else
      c := 0;
    // escape が要る文字: " (34)、\ (92)、/ (47)、制御文字 (<32)、非 ASCII (>127)。
    if (c = 34) or (c = 92) or (c = 47) or (c < 32) or (c > 127) then begin
      // 直前までの安全区間をまとめて追記する。
      if i > runStart then
        Result := Result + Copy(s, runStart, i - runStart);
      case c of
        34: Result := Result + '\"';
        92: Result := Result + '\\';
        47: Result := Result + '\/';
         8: Result := Result + '\b';
         9: Result := Result + '\t';
        10: Result := Result + '\n';
        12: Result := Result + '\f';
        13: Result := Result + '\r';
      else
        Result := Result + '\u' + IntToHex(c, 4);
      end;
      runStart := i + 1;
    end;
  end;
  // 末尾の安全区間を追記する。
  if runStart <= len then
    Result := Result + Copy(s, runStart, len - runStart + 1);
end;

function JsonString(s: string): string;
begin
  Result := '"' + EscapeJSONString(s) + '"';
end;

function JsonField(key, value: string): string;
begin
  Result := '  "' + key + '": ' + value;
end;

function GetFileNameWithExtension(aFile: IwbFile): string;
var
  FileName, lower: string;
begin
  FileName := GetFileName(aFile);
  lower := LowerCase(FileName);
  if (Pos('.esl', lower) > 0) or (Pos('.esm', lower) > 0) or (Pos('.esp', lower) > 0) then
    Result := FileName
  else if GetIsESL(aFile) then
    Result := FileName + '.esl'
  else if GetIsESM(aFile) then
    Result := FileName + '.esm'
  else
    Result := FileName + '.esp';
end;

function GetMasterFileName(e: IInterface): string;
var
  recordFile: IwbFile;
begin
  if not Assigned(e) then begin Result := ''; Exit; end;
  recordFile := GetFile(e);
  if not Assigned(recordFile) then Result := ''
  else Result := GetFileNameWithExtension(recordFile);
end;

function GetCachedEdid(e: IInterface): string;
begin
  if Assigned(e) then
    Result := GetElementEditValues(MasterOrSelf(e), 'EDID')
  else
    Result := '';
end;

// ===== Ensure pattern (cross-file traversal) =====
// 参照解決で master record に到達した時、未処理なら適切な Extract に dispatch する。

function IsProcessed(const id: string): Boolean;
begin
  Result := processedIds.IndexOf(id) >= 0;
end;

procedure MarkProcessed(const id: string);
begin
  processedIds.Add(id);
end;

// EnsureRecord は record signature ごとの dispatch を最下端（Finalize 直前）で実装する。
// JvInterpreter は procedure 呼び出し時に名前を解決するため、定義より前で呼んでも
// forward 宣言は要らない。forward 指令自体を解釈せず parse error になるため使わない。

// ===== Race =====

procedure ExtractRace(race: IInterface);
var
  raceID, raceName, raceDesc, source, edid, entry: string;
begin
  raceID := HexFormID(race);
  if IsProcessed(raceID) then Exit;
  MarkProcessed(raceID);
  raceName := GetElementValue(race, 'FULL');
  raceDesc := GetElementValue(race, 'DESC');
  if (raceName = '') and (raceDesc = '') then Exit;
  source := GetMasterFileName(race);
  edid := GetCachedEdid(race);

  // primary type = "RACE FULL"。description は "RACE DESC" mapping（character creation 表示）。
  entry := '  {' + #13#10 +
           JsonField('id', JsonString(raceID)) + ',' + #13#10 +
           JsonField('editor_id', JsonString(edid)) + ',' + #13#10 +
           JsonField('type', JsonString('RACE FULL')) + ',' + #13#10 +
           JsonField('source', JsonString(source)) + ',' + #13#10 +
           JsonField('name', JsonString(raceName)) + ',' + #13#10 +
           JsonField('description', JsonString(raceDesc)) + #13#10 +
           '  }';
  raceList.Add(entry);
end;

// ===== Faction =====

procedure ExtractFaction(faction: IInterface);
var
  factionID, factionName, source, edid, entry: string;
begin
  factionID := HexFormID(faction);
  if IsProcessed(factionID) then Exit;
  MarkProcessed(factionID);
  factionName := GetElementValue(faction, 'FULL');
  source := GetMasterFileName(faction);
  edid := GetCachedEdid(faction);

  // 名前空でも edid を頼りに翻訳 context に効く可能性があるため emit する。
  entry := '  {' + #13#10 +
           JsonField('id', JsonString(factionID)) + ',' + #13#10 +
           JsonField('editor_id', JsonString(edid)) + ',' + #13#10 +
           JsonField('type', JsonString('FACT FULL')) + ',' + #13#10 +
           JsonField('source', JsonString(source)) + ',' + #13#10 +
           JsonField('name', JsonString(factionName)) + #13#10 +
           '  }';
  factionList.Add(entry);
end;

// ===== Speaker (NPC_ / TACT) =====
// v18: 発話主体（しゃべる人・もの）を 1 class に統合する。NPC_（固有 / 形態 / プール /
// 声型代表）と TACT（喋るオブジェクト）を種別 attribute で束ね、名前空でも emit する。

// ペルソナ名を解決する。FULL → TPLT 連鎖の本体 FULL → EditorID 役割名 の優先順。
function ResolvePersonaName(actor: IInterface): string;
var
  cur: IInterface;
  depth: Integer;
  f: string;
begin
  cur := actor;
  depth := 0;
  while Assigned(cur) and (depth < 12) do begin
    f := GetElementValue(cur, 'FULL');
    if f <> '' then begin Result := f; Exit; end;
    if not ElementExists(cur, 'TPLT - Template') then Break;
    cur := LinksTo(ElementByName(cur, 'TPLT - Template'));
    Inc(depth);
  end;
  // FULL に届かない場合は EditorID（役割名）をペルソナ名にする。
  Result := GetCachedEdid(actor);
end;

// Speaker.種別を判定する。
function ClassifySpeakerKind(actor: IInterface): string;
var
  cur: IInterface;
  depth: Integer;
  sig: string;
begin
  if Signature(actor) = 'TACT' then begin Result := '喋るオブジェクト'; Exit; end;
  if GetElementValue(actor, 'FULL') <> '' then begin Result := '固有'; Exit; end;
  cur := actor;
  depth := 0;
  while Assigned(cur) and (depth < 12) do begin
    sig := Signature(cur);
    if sig = 'LVLN' then begin Result := 'プール'; Exit; end;
    if (sig = 'NPC_') and (GetElementValue(cur, 'FULL') <> '') then begin Result := '形態'; Exit; end;
    if not ElementExists(cur, 'TPLT - Template') then Break;
    cur := LinksTo(ElementByName(cur, 'TPLT - Template'));
    Inc(depth);
  end;
  if ElementExists(actor, 'VTCK - Voice') then Result := '声型代表'
  else Result := '役割のみ';
end;

procedure ExtractSpeaker(actor: IInterface);
var
  spkID, fullName, shortName, voiceId, sex, classID, source, edid: string;
  raceID, baseSpeakerID, personaName, kind, sig: string;
  raceRec, vtypRec, tplRec: IInterface;
  factions, factionEntry, factionRef: IInterface;
  factionIdsJson, entry: string;
  i: integer;
begin
  spkID := HexFormID(actor);
  if IsProcessed(spkID) then Exit;
  MarkProcessed(spkID);
  sig := Signature(actor);
  source := GetMasterFileName(actor);
  edid := GetCachedEdid(actor);
  fullName := GetElementValue(actor, 'FULL');
  shortName := '';
  sex := '';
  classID := '';
  raceID := '';
  baseSpeakerID := '';
  voiceId := '';
  factionIdsJson := '';

  if sig = 'TACT' then begin
    // 喋るオブジェクト。声型は VNAM。
    vtypRec := LinksTo(ElementByName(actor, 'VNAM - Voice Type'));
    if Assigned(vtypRec) then begin
      voiceId := HexFormID(vtypRec);
      EnsureRecord(vtypRec);
    end;
    personaName := fullName;
    if personaName = '' then personaName := edid;
    kind := '喋るオブジェクト';
  end else begin
    // NPC_
    shortName := GetElementValue(actor, 'SHRT');
    classID := GetElementValue(actor, 'CNAM');
    if (GetElementNativeValues(actor, 'ACBS\Flags') and 1) <> 0 then
      sex := 'Female' else sex := 'Male';

    raceRec := LinksTo(ElementByName(actor, 'RNAM - Race'));
    if Assigned(raceRec) then begin
      raceID := HexFormID(raceRec);
      EnsureRecord(raceRec);
    end;

    // rV: VTCK が指す VoiceType
    vtypRec := LinksTo(ElementByName(actor, 'VTCK - Voice'));
    if Assigned(vtypRec) then begin
      voiceId := HexFormID(vtypRec);
      EnsureRecord(vtypRec);
    end;

    // rT: TPLT が指す本体 Speaker（形態 record の本体参照）
    tplRec := LinksTo(ElementByName(actor, 'TPLT - Template'));
    if Assigned(tplRec) and (Signature(tplRec) = 'NPC_') then begin
      baseSpeakerID := HexFormID(tplRec);
      EnsureRecord(tplRec);
    end;

    // r4: Faction list
    factions := ElementByName(actor, 'Factions');
    if Assigned(factions) then
      for i := 0 to ElementCount(factions) - 1 do begin
        factionEntry := ElementByIndex(factions, i);
        factionRef := LinksTo(ElementByName(factionEntry, 'Faction'));
        if not Assigned(factionRef) then Continue;
        EnsureRecord(factionRef);
        if factionIdsJson <> '' then factionIdsJson := factionIdsJson + ', ';
        factionIdsJson := factionIdsJson + JsonString(HexFormID(factionRef));
      end;

    personaName := ResolvePersonaName(actor);
    kind := ClassifySpeakerKind(actor);
  end;

  // v18: 名前空でも emit する（ペルソナ名が必ず解決できるため）。
  entry := '  "' + spkID + '": {' + #13#10 +
           '    "id": ' + JsonString(spkID) + ',' + #13#10 +
           '    "editor_id": ' + JsonString(edid) + ',' + #13#10 +
           '    "type": ' + JsonString(sig + ' FULL') + ',' + #13#10 +
           '    "source": ' + JsonString(source) + ',' + #13#10 +
           '    "name": ' + JsonString(fullName) + ',' + #13#10 +
           '    "persona_name": ' + JsonString(personaName) + ',' + #13#10 +
           '    "kind": ' + JsonString(kind) + ',' + #13#10 +
           '    "short_name": ' + JsonString(shortName) + ',' + #13#10 +
           '    "sex": ' + JsonString(sex) + ',' + #13#10 +
           '    "voice_type_id": ' + JsonString(voiceId) + ',' + #13#10 +
           '    "class_id": ' + JsonString(classID) + ',' + #13#10 +
           '    "race_id": ' + JsonString(raceID) + ',' + #13#10 +
           '    "base_speaker_id": ' + JsonString(baseSpeakerID) + ',' + #13#10 +
           '    "faction_ids": [' + factionIdsJson + ']' + #13#10 +
           '  }';
  personList.Add(entry);
end;

// ===== Voice resolution（v1 から流用） =====

// GetRecordVoiceTypes2 / GetConditionsVoiceTypes / GetAliasVoiceTypes / InfoVoiceTypes は
// 相互に呼び合う。JvInterpreter は実行時に名前解決するため、定義順だけ整え forward は使わない。
procedure GetRecordVoiceTypes2(e: IInterface; lstVoice: TStringList);
var
  sig: string;
  i: integer;
  ent, ents: IInterface;
begin
  if lstRecursion.IndexOf(FormID(e)) <> -1 then Exit
    else lstRecursion.Add(FormID(e));

  e := WinningOverride(e);
  sig := Signature(e);
  if (sig = 'REFR') or (sig = 'ACHR') then begin
    e := WinningOverride(BaseRecord(e));
    sig := Signature(e);
  end;

  if sig = 'VTYP' then
    lstVoice.AddObject(EditorID(e), e)
  else if sig = 'NPC_' then begin
    if ElementExists(e, 'VTCK - Voice') then begin
      InfoNPCID := HexFormID(e);
      EnsureRecord(e);
      InfoRACEID := GetElementEditValues(e, 'RNAM - Race');
      ent := LinksTo(ElementByName(e, 'VTCK - Voice'));
      GetRecordVoiceTypes2(ent, lstVoice);
    end
    else if ElementExists(e, 'TPLT - Template') then begin
      ent := LinksTo(ElementByName(e, 'TPLT - Template'));
      GetRecordVoiceTypes2(ent, lstVoice);
    end;
  end
  else if sig = 'LVLN' then begin
    ents := ElementByName(e, 'Leveled List Entries');
    for i := 0 to Pred(ElementCount(ents)) do begin
      ent := ElementByIndex(ents, i);
      ent := LinksTo(ElementByPath(ent, 'LVLO\NPC'));
      GetRecordVoiceTypes2(ent, lstVoice);
    end;
  end
  else if sig = 'TACT' then begin
    ent := WinningOverride(LinksTo(ElementByName(e, 'VNAM - Voice Type')));
    if Signature(ent) = 'VTYP' then lstVoice.AddObject(EditorID(ent), ent);
  end
  else if sig = 'FLST' then begin
    ents := ElementByName(e, 'FormIDs');
    for i := 0 to Pred(ElementCount(ents)) do begin
      ent := LinksTo(ElementByIndex(ents, i));
      GetRecordVoiceTypes2(ent, lstVoice);
    end;
  end
  else if (sig = 'FACT') or (sig = 'CLAS') then begin
    for i := 0 to Pred(ReferencedByCount(e)) do begin
      ent := ReferencedByIndex(e, i);
      if Signature(ent) = 'NPC_' then
        GetRecordVoiceTypes2(ent, lstVoice);
    end;
  end;
end;

procedure GetRecordVoiceTypes(e: IInterface; lstVoice: TStringList);
begin
  lstRecursion.Clear;
  GetRecordVoiceTypes2(e, lstVoice);
end;

procedure GetAliasVoiceTypes(Quest: IInterface; aAlias: integer; lstVoice: TStringList);
var
  Aliases, Alias: IInterface;
  i, j: integer;
  lstLimit: TStringList;
begin
  Quest := WinningOverride(Quest);
  Aliases := ElementByName(Quest, 'Aliases');
  for i := 0 to Pred(ElementCount(Aliases)) do begin
    Alias := ElementByIndex(Aliases, i);
    if GetNativeValue(ElementByIndex(Alias, 0)) <> aAlias then Continue;

    if ElementExists(Alias, 'ALFR - Forced Reference') then
      GetRecordVoiceTypes(LinksTo(ElementByName(Alias, 'ALFR - Forced Reference')), lstVoice)
    else if ElementExists(Alias, 'ALUA - Unique Actor') then
      GetRecordVoiceTypes(LinksTo(ElementByName(Alias, 'ALUA - Unique Actor')), lstVoice)
    else if ElementExists(Alias, 'External Alias Reference') then
      GetAliasVoiceTypes(
        LinksTo(ElementByPath(Alias, 'External Alias Reference\ALEQ - Quest')),
        GetElementNativeValues(Alias, 'External Alias Reference\ALEA - Alias'),
        lstVoice)
    else if ElementExists(Alias, 'Conditions') then
      GetConditionsVoiceTypes(ElementByName(Alias, 'Conditions'), lstVoice);

    if GetElementNativeValues(Alias, 'VTCK - Voice Types') <> 0 then begin
      if lstVoice.Count <> 0 then begin
        lstLimit := TStringList.Create;
        GetRecordVoiceTypes(LinksTo(ElementByName(Alias, 'VTCK - Voice Types')), lstLimit);
        for j := Pred(lstVoice.Count) downto 0 do
          if lstLimit.IndexOf(lstVoice[j]) = -1 then
            lstVoice.Delete(j);
        lstLimit.Free;
      end else
        GetRecordVoiceTypes(LinksTo(ElementByName(Alias, 'VTCK - Voice Types')), lstVoice);
    end;
    Break;
  end;
end;

procedure GetConditionsVoiceTypes(Conditions: IInterface; lstVoice: TStringList);
var
  Condition, Elem, Quest: IInterface;
  ConditionFunction: string;
  lstVoiceCondition: TStringList;
  i, j, Alias: integer;
  bFactionCondition, bGetIsID: Boolean;
begin
  lstVoiceCondition := TStringList.Create;
  lstVoiceCondition.Duplicates := dupIgnore;
  lstVoiceCondition.Sorted := True;

  for i := 0 to Pred(ElementCount(Conditions)) do
    if GetElementEditValues(ElementByIndex(Conditions, i), 'CTDA\Function') = 'GetIsID' then begin
      bGetIsID := True;
      Break;
    end;

  for i := 0 to Pred(ElementCount(Conditions)) do begin
    Condition := ElementByIndex(Conditions, i);
    ConditionFunction := GetElementEditValues(Condition, 'CTDA\Function');

    if ConditionFunction = 'GetIsID' then begin
      InfoCONDITION := ConditionFunction;
      Elem := LinksTo(ElementByPath(Condition, 'CTDA\Base Object'));
      GetRecordVoiceTypes(Elem, lstVoiceCondition);
    end else
    if not bGetIsID then
    if ConditionFunction = 'GetIsVoiceType' then begin
      InfoCONDITION := ConditionFunction;
      Elem := LinksTo(ElementByPath(Condition, 'CTDA\Voice Type'));
      GetRecordVoiceTypes(Elem, lstVoiceCondition);
    end
    else if ConditionFunction = 'GetIsAliasRef' then begin
      Alias := GetElementNativeValues(Condition, 'CTDA\Alias');
      Elem := ContainingMainRecord(Conditions);
      if Signature(Elem) = 'INFO' then begin
        Elem := LinksTo(ElementByName(Elem, 'Topic'));
        Quest := LinksTo(ElementByName(Elem, 'QNAM - Quest'));
      end
      else if Signature(Elem) = 'QUST' then
        Quest := Elem;
      GetAliasVoiceTypes(Quest, Alias, lstVoiceCondition);
    end else
    if ConditionFunction = 'GetInFaction' then begin
      if not bFactionCondition then begin
        bFactionCondition := True;
        Elem := LinksTo(ElementByPath(Condition, 'CTDA\Faction'));
        GetRecordVoiceTypes(Elem, lstVoiceCondition);
      end;
    end
    else if ConditionFunction = 'GetIsClass' then begin
      Elem := LinksTo(ElementByPath(Condition, 'CTDA\Class'));
      GetRecordVoiceTypes(Elem, lstVoiceCondition);
    end;

    if lstVoiceCondition.Count = 0 then Continue;
    lstVoice.AddStrings(lstVoiceCondition);

    if GetElementNativeValues(Condition, 'CTDA\Type') and 1 > 0 then
      lstVoice.AddStrings(lstVoiceCondition)
    else begin
      for j := Pred(lstVoice.Count) downto 0 do
        if lstVoiceCondition.IndexOf(lstVoice[j]) = -1 then
          lstVoice.Delete(j);
    end;
    lstVoiceCondition.Clear;
  end;
  lstVoiceCondition.Free;
end;

procedure InfoVoiceTypes(Info: IInterface; lstVoice: TStringList; QuestConditions: IInterface);
var
  Elem, Dialogue: IInterface;
  Scene, Actions, Action: IInterface;
  i, j, Alias: integer;
  bAliasFound: Boolean;
  lstVoiceQuestConditions: TStringList;
begin
  if not Assigned(lstVoice) then Exit;

  lstVoiceQuestConditions := TStringList.Create;
  lstVoiceQuestConditions.Duplicates := dupIgnore;
  lstVoiceQuestConditions.Sorted := True;
  lstVoice.Clear;

  Elem := ElementByName(Info, 'ANAM - Speaker');
  if Assigned(Elem) then begin
    Elem := LinksTo(Elem);
    GetRecordVoiceTypes(Elem, lstVoice);
    InfoSPEAKER := HexFormID(Elem);
    EnsureRecord(Elem);
    lstVoiceQuestConditions.Free;
    Exit;
  end;

  if QuestConditions <> nil then
    GetConditionsVoiceTypes(QuestConditions, lstVoiceQuestConditions);
  if ElementExists(Info, 'Conditions') then
    GetConditionsVoiceTypes(ElementByName(Info, 'Conditions'), lstVoice);
  if (lstVoice.Count = 0) and (lstVoiceQuestConditions.Count <> 0) then
    lstVoice.AddStrings(lstVoiceQuestConditions);
  if (lstVoice.Count <> 0) and (lstVoiceQuestConditions.Count <> 0) then
    for i := Pred(lstVoice.Count) downto 0 do
      if lstVoiceQuestConditions.IndexOf(lstVoice[i]) = -1 then
        lstVoice.Delete(i);
  lstVoiceQuestConditions.Free;

  if lstVoice.Count <> 0 then Exit;

  bAliasFound := False;
  Dialogue := LinksTo(ElementByName(Info, 'Topic'));
  if GetElementEditValues(Dialogue, 'DATA\Category') = 'Scene' then
    for i := Pred(ReferencedByCount(Dialogue)) downto 0 do begin
      Scene := ReferencedByIndex(Dialogue, i);
      if Signature(Scene) <> 'SCEN' then Continue;
      Actions := ElementByName(Scene, 'Actions');
      for j := 0 to Pred(ElementCount(Actions)) do begin
        Action := ElementByIndex(Actions, j);
        if (GetElementEditValues(Action, 'ANAM - Type') = 'Dialogue') and
           (GetElementNativeValues(Action, 'DATA - Topic') = GetLoadOrderFormID(Dialogue))
        then begin
          Alias := GetElementNativeValues(Action, 'ALID - Actor ID');
          GetAliasVoiceTypes(LinksTo(ElementByName(Scene, 'PNAM - Quest')), Alias, lstVoice);
          bAliasFound := True;
          Break;
        end;
      end;
      if bAliasFound then Break;
    end;
end;

// ===== Dialogue =====

function GetOrCreateDialList(dialID: string): TStringList;
var
  idx: integer;
  lst: TStringList;
begin
  idx := dialDataMap.IndexOf(dialID);
  if idx = -1 then begin
    lst := TStringList.Create;
    lst.Add('');  // header slot 確保
    dialDataMap.AddObject(dialID, lst);
    Result := lst;
    Exit;
  end;
  if dialDataMap.Objects[idx] = nil then begin
    lst := TStringList.Create;
    lst.Add('');
    dialDataMap.Objects[idx] := lst;
    Result := lst;
    Exit;
  end;
  Result := TStringList(dialDataMap.Objects[idx]);
end;

procedure ExtractDialogue(dial: IInterface);
var
  dialID, playerText, questID, source, edid: string;
  category, subtype, servicesType: string;
  isServices, dialJSON: string;
  questRec: IInterface;
  lst: TStringList;
begin
  dialID := HexFormID(dial);
  playerText := GetElementValue(dial, 'FULL');
  source := GetMasterFileName(dial);
  edid := GetCachedEdid(dial);

  // DIAL.DATA は Category byte と Subtype byte を含む。Subtype = 発火種別。
  category := GetElementValue(dial, 'DATA\Category');
  subtype := GetElementValue(dial, 'DATA\Subtype');
  servicesType := subtype;  // v1 互換: services_type field

  // r5 Quest 解決 + Ensure
  questRec := LinksTo(ElementByName(dial, 'QNAM'));
  if Assigned(questRec) then begin
    questID := HexFormID(questRec);
    EnsureRecord(questRec);
  end else
    questID := GetElementValue(dial, 'QNAM');  // fallback raw string

  if (subtype = 'Services') or (subtype = 'Training') or (subtype = 'Barter') or (subtype = 'Misc') then
    isServices := 'true' else isServices := 'false';

  dialJSON := '  {' + #13#10 +
              '  "id": ' + JsonString(dialID) + ',' + #13#10 +
              '  "editor_id": ' + JsonString(edid) + ',' + #13#10 +
              '  "type": "DIAL FULL",' + #13#10 +
              '  "source": ' + JsonString(source) + ',' + #13#10 +
              '  "player_text": ' + JsonString(playerText) + ',' + #13#10 +
              '  "category": ' + JsonString(category) + ',' + #13#10 +
              '  "subtype": ' + JsonString(subtype) + ',' + #13#10 +
              '  "is_services_branch": ' + isServices + ',' + #13#10 +
              '  "services_type": ' + JsonString(servicesType) + ',' + #13#10 +
              '  "quest_id": ' + JsonString(questID) + ',';

  lst := GetOrCreateDialList(dialID);
  lst[0] := dialJSON;
end;

// INFO.CTDA を condition 配列 JSON に変換する。
// 各 entry は function 名 / operator / value / ref_id（あれば）を持つ。
// CTDA functions のうち代表的な ref を持つもの（GetIsID, GetInFaction, GetIsRace,
// GetIsClass, GetIsVoiceType, GetIsAliasRef）は ref_id を埋める。それ以外は ref_id 空。
function BuildConditionsJson(info: IInterface): string;
var
  Conditions, Condition, refElem: IInterface;
  i: integer;
  funcName, opStr, valStr, refID, entry: string;
  opCode, opValue: Integer;
begin
  Result := '';
  if not ElementExists(info, 'Conditions') then Exit;
  Conditions := ElementByName(info, 'Conditions');
  if not Assigned(Conditions) then Exit;

  for i := 0 to ElementCount(Conditions) - 1 do begin
    Condition := ElementByIndex(Conditions, i);
    funcName := GetElementEditValues(Condition, 'CTDA\Function');
    valStr := GetElementEditValues(Condition, 'CTDA\Comparison Value');

    // operator byte の下位 3 bit が比較演算子（==, !=, >, >=, <, <=）。
    opCode := GetElementNativeValues(Condition, 'CTDA\Type');
    opValue := (opCode shr 5) and 7;
    case opValue of
      0: opStr := '==';
      1: opStr := '!=';
      2: opStr := '>';
      3: opStr := '>=';
      4: opStr := '<';
      5: opStr := '<=';
    else
      opStr := '==';
    end;

    refID := '';
    refElem := nil;
    if funcName = 'GetIsID' then
      refElem := LinksTo(ElementByPath(Condition, 'CTDA\Base Object'))
    else if funcName = 'GetInFaction' then
      refElem := LinksTo(ElementByPath(Condition, 'CTDA\Faction'))
    else if funcName = 'GetIsRace' then
      refElem := LinksTo(ElementByPath(Condition, 'CTDA\Race'))
    else if funcName = 'GetIsClass' then
      refElem := LinksTo(ElementByPath(Condition, 'CTDA\Class'))
    else if funcName = 'GetIsVoiceType' then
      refElem := LinksTo(ElementByPath(Condition, 'CTDA\Voice Type'));
    if Assigned(refElem) then begin
      refID := HexFormID(refElem);
      EnsureRecord(refElem);
    end;

    entry := '        {' + #13#10 +
             '          "function": ' + JsonString(funcName) + ',' + #13#10 +
             '          "operator": ' + JsonString(opStr) + ',' + #13#10 +
             '          "value": ' + JsonString(valStr) + ',' + #13#10 +
             '          "ref_id": ' + JsonString(refID) + #13#10 +
             '        }';
    if Result <> '' then Result := Result + ',' + #13#10;
    Result := Result + entry;
  end;
end;

// INFO.TCLT（forward-link to next INFO）の全 entry を id list として収集。
function BuildNextIdsJson(info: IInterface): string;
var
  tclt, entry, refRec: IInterface;
  i: integer;
  refID: string;
begin
  Result := '';
  tclt := ElementByName(info, 'TCLT - Linked To');
  if not Assigned(tclt) then Exit;
  // TCLT は FormID 配列。要素数を ElementCount で取り、LinksTo で解決。
  for i := 0 to ElementCount(tclt) - 1 do begin
    entry := ElementByIndex(tclt, i);
    refRec := LinksTo(entry);
    if Assigned(refRec) then begin
      refID := HexFormID(refRec);
      if Result <> '' then Result := Result + ', ';
      Result := Result + JsonString(refID);
    end;
  end;
end;

// REFR / ACHR を base record（NPC_ 等）へ解決する（alias の forced reference 用）。
function ResolveActorBase(e: IInterface): IInterface;
var sig: string;
begin
  Result := e;
  if not Assigned(e) then Exit;
  sig := Signature(e);
  if (sig = 'REFR') or (sig = 'ACHR') then
    Result := WinningOverride(BaseRecord(e));
end;

// INFO が属する Quest を辿る（INFO -> Topic(DIAL) -> QNAM）。
function QuestOfInfo(info: IInterface): IInterface;
var topic, q: IInterface;
begin
  Result := nil;
  topic := LinksTo(ElementByName(info, 'Topic'));
  if not Assigned(topic) then Exit;
  q := LinksTo(ElementByName(topic, 'QNAM - Quest'));
  if not Assigned(q) then q := LinksTo(ElementByName(topic, 'QNAM'));
  Result := q;
end;

// quest alias index から forced reference / unique actor を解決し actor record を返す。
// external / conditions ベースの alias は解決せず nil を返す。
function ResolveAliasActor(quest: IInterface; aliasIdx: Integer): IInterface;
var aliases, al, ref: IInterface; i: Integer;
begin
  Result := nil;
  if not Assigned(quest) then Exit;
  aliases := ElementByName(quest, 'Aliases');
  if not Assigned(aliases) then Exit;
  for i := 0 to ElementCount(aliases) - 1 do begin
    al := ElementByIndex(aliases, i);
    if GetNativeValue(ElementByIndex(al, 0)) <> aliasIdx then Continue;
    if ElementExists(al, 'ALFR - Forced Reference') then
      ref := LinksTo(ElementByName(al, 'ALFR - Forced Reference'))
    else if ElementExists(al, 'ALUA - Unique Actor') then
      ref := LinksTo(ElementByName(al, 'ALUA - Unique Actor'))
    else ref := nil;
    Result := ResolveActorBase(ref);
    Break;
  end;
end;

// JSON 配列文字列に id を追記する。空 id は無視する。
procedure AppendIdJson(var arr: string; id: string);
begin
  if id = '' then Exit;
  if arr <> '' then arr := arr + ', ';
  arr := arr + JsonString(id);
end;

// v18: 話者関連は r2a（名指し Speaker）/ r2b（Faction）/ r2c（Race）/ r2d（VoiceType）の
// 4 配列で表す。voice_types は Speaker→VoiceType の集約による導出属性として残す。
function BuildResponseJson(
    infoID, edid, source, infoText, prompt, topicText, menuDisplayText,
    speakerIdsJson, factionIdsJson, raceIdsJson, voiceTypeIdsJson,
    voiceTypesJson, previousID, conditionsJson,
    nextIdsJson, responseNum: string;
    order: integer): string;
begin
  Result := '    {' + #13#10 +
            '      "id": ' + JsonString(infoID) + ',' + #13#10 +
            '      "editor_id": ' + JsonString(edid) + ',' + #13#10 +
            '      "type": "INFO NAM1",' + #13#10 +
            '      "source": ' + JsonString(source) + ',' + #13#10 +
            '      "order": ' + IntToStr(order) + ',' + #13#10 +
            '      "text": ' + JsonString(infoText) + ',' + #13#10 +
            '      "prompt": ' + JsonString(prompt) + ',' + #13#10 +
            '      "topic_text": ' + JsonString(topicText) + ',' + #13#10 +
            '      "menu_display_text": ' + JsonString(menuDisplayText) + ',' + #13#10 +
            '      "speaker_ids": [' + speakerIdsJson + '],' + #13#10 +
            '      "faction_ids": [' + factionIdsJson + '],' + #13#10 +
            '      "race_ids": [' + raceIdsJson + '],' + #13#10 +
            '      "voice_type_ids": [' + voiceTypeIdsJson + '],' + #13#10 +
            '      "voice_types": [' + voiceTypesJson + '],' + #13#10;
  if responseNum <> '' then
    Result := Result + '      "response_number": ' + responseNum + ',' + #13#10;
  Result := Result +
            '      "previous_id": ' + JsonString(previousID) + ',' + #13#10 +
            '      "next_ids": [' + nextIdsJson + '],' + #13#10 +
            '      "conditions": [';
  if conditionsJson <> '' then
    Result := Result + #13#10 + conditionsJson + #13#10 + '      ]'
  else
    Result := Result + ']';
  Result := Result + #13#10 + '    }';
end;

procedure ExtractInfo(info: IInterface);
var
  infoID, infoText, previousID, source, edid: string;
  speakerIdsJson, factionIdsJson, raceIdsJson, voiceTypeIdsJson, fnName: string;
  condsRec, condItem, refRec, spkQuest, anamRec: IInterface;
  dialID, prompt, topicText, menuDisplay: string;
  parentRec, questRec, questConditions, responses, responseEntry: IInterface;
  lst, lstVoice: TStringList;
  i: integer;
  responseNum: string;
  voiceTypesJson: string;
  conditionsJson: string;
  nextIdsJson: string;
  responseCount: integer;
  built: string;
begin
  InfoNPCID := '';
  InfoSPEAKER := '';
  InfoRACEID := '';
  InfoCONDITION := '';

  parentRec := GetContainer(info);
  while Assigned(parentRec) and (Signature(parentRec) <> 'DIAL') and (ElementType(parentRec) <> etFile) do
    parentRec := GetContainer(parentRec);
  if not Assigned(parentRec) or (Signature(parentRec) <> 'DIAL') then begin
    parentRec := LinksTo(ElementByName(info, 'Topic'));
    if not Assigned(parentRec) or (Signature(parentRec) <> 'DIAL') then begin
      AddMessage('Debug: Info ' + HexFormID(info) + ' could not resolve DIAL parent.');
      Exit;
    end;
  end;

  dialID := HexFormID(parentRec);
  infoID := HexFormID(info);
  source := GetMasterFileName(info);
  edid := GetCachedEdid(info);

  prompt := GetElementValue(info, 'RNAM');
  topicText := GetElementValue(parentRec, 'FULL');
  menuDisplay := prompt;

  lstVoice := TStringList.Create;
  lstVoice.Duplicates := dupIgnore;
  lstVoice.Sorted := True;

  questConditions := nil;
  if ElementExists(parentRec, 'QNAM') then begin
    questRec := LinksTo(ElementByName(parentRec, 'QNAM'));
    if Assigned(questRec) and ElementExists(questRec, 'Quest Dialogue Conditions') then
      questConditions := ElementByPath(questRec, 'Quest Dialogue Conditions\Conditions');
  end;

  InfoVoiceTypes(info, lstVoice, questConditions);

  // voice_types（導出）は文字列配列の JSON として組み立てる。
  voiceTypesJson := '';
  for i := 0 to lstVoice.Count - 1 do begin
    if voiceTypesJson <> '' then voiceTypesJson := voiceTypesJson + ', ';
    voiceTypesJson := voiceTypesJson + JsonString(lstVoice[i]);
  end;
  lstVoice.Free;

  // r2a-d: 話者関連を ANAM と conditions（CTDA）から収集する（v18）。
  // 名指し（r2a）は ANAM / GetIsID(NPC_/TACT) / GetIsAliasRef、集合は GetInFaction(r2b) /
  // GetIsRace(r2c) / GetIsVoiceType(r2d)。voice 解決経由のテンプレ拾いはしない。
  spkQuest := QuestOfInfo(info);
  speakerIdsJson := '';
  factionIdsJson := '';
  raceIdsJson := '';
  voiceTypeIdsJson := '';
  anamRec := LinksTo(ElementByName(info, 'ANAM - Speaker'));
  if Assigned(anamRec) then begin
    EnsureRecord(anamRec);
    AppendIdJson(speakerIdsJson, HexFormID(anamRec));
  end;
  condsRec := ElementByName(info, 'Conditions');
  if Assigned(condsRec) then
    for i := 0 to ElementCount(condsRec) - 1 do begin
      condItem := ElementByIndex(condsRec, i);
      fnName := GetElementEditValues(condItem, 'CTDA\Function');
      if fnName = 'GetIsID' then begin
        refRec := LinksTo(ElementByPath(condItem, 'CTDA\Base Object'));
        if Assigned(refRec) and ((Signature(refRec) = 'NPC_') or (Signature(refRec) = 'TACT')) then begin
          EnsureRecord(refRec);
          AppendIdJson(speakerIdsJson, HexFormID(refRec));
        end;
      end
      else if fnName = 'GetIsAliasRef' then begin
        refRec := ResolveAliasActor(spkQuest, GetElementNativeValues(condItem, 'CTDA\Alias'));
        if Assigned(refRec) and (Signature(refRec) = 'NPC_') then begin
          EnsureRecord(refRec);
          AppendIdJson(speakerIdsJson, HexFormID(refRec));
        end;
      end
      else if fnName = 'GetInFaction' then begin
        refRec := LinksTo(ElementByPath(condItem, 'CTDA\Faction'));
        if Assigned(refRec) then begin
          EnsureRecord(refRec);
          AppendIdJson(factionIdsJson, HexFormID(refRec));
        end;
      end
      else if fnName = 'GetIsRace' then begin
        refRec := LinksTo(ElementByPath(condItem, 'CTDA\Race'));
        if Assigned(refRec) then begin
          EnsureRecord(refRec);
          AppendIdJson(raceIdsJson, HexFormID(refRec));
        end;
      end
      else if fnName = 'GetIsVoiceType' then begin
        refRec := LinksTo(ElementByPath(condItem, 'CTDA\Voice Type'));
        if Assigned(refRec) and (Signature(refRec) = 'VTYP') then begin
          EnsureRecord(refRec);
          AppendIdJson(voiceTypeIdsJson, HexFormID(refRec));
        end;
      end;
    end;

  // PNAM（previous INFO）は参照を解決して FormID で持つ。EditValue（冗長な表示
  // 文字列）を取ると ApplyInfoSorting の chain 照合も JSON の previous_id も壊れる。
  // 参照が無い root INFO は HexFormID(nil) = '00000000' になる。
  previousID := HexFormID(LinksTo(ElementByName(info, 'PNAM - Previous INFO')));
  conditionsJson := BuildConditionsJson(info);
  nextIdsJson := BuildNextIdsJson(info);
  lst := GetOrCreateDialList(dialID);

  if ElementExists(info, 'Responses') then begin
    responses := ElementByName(info, 'Responses');
    responseCount := ElementCount(responses);
    for i := 0 to responseCount - 1 do begin
      responseEntry := ElementByIndex(responses, i);
      infoText := GetElementEditValues(responseEntry, 'NAM1');
      responseNum := '';
      if ElementExists(responseEntry, 'TRDT') then
        responseNum := GetElementEditValues(responseEntry, 'TRDT\Response number');

      // prompt（RNAM）だけを持つ INFO（プレイヤー選択肢のみで NPC 応答が無い）も出す。
      if (infoText <> '') or (prompt <> '') or (speakerIdsJson <> '') or (factionIdsJson <> '')
       or (raceIdsJson <> '') or (voiceTypeIdsJson <> '') then begin
        built := BuildResponseJson(infoID, edid, source, infoText, prompt,
            topicText, menuDisplay, speakerIdsJson, factionIdsJson, raceIdsJson,
            voiceTypeIdsJson, voiceTypesJson, previousID, conditionsJson,
            nextIdsJson, responseNum, i);
        lst.Add(infoID + SEP_CHAR + previousID + SEP_CHAR + built);
      end;
    end;
  end else begin
    infoText := GetElementValue(info, 'NAM1');
    if (infoText <> '') or (prompt <> '') or (speakerIdsJson <> '') or (factionIdsJson <> '')
       or (raceIdsJson <> '') or (voiceTypeIdsJson <> '') then begin
      built := BuildResponseJson(infoID, edid, source, infoText, prompt,
          topicText, menuDisplay, speakerIdsJson, factionIdsJson, raceIdsJson,
          voiceTypeIdsJson, voiceTypesJson, previousID, conditionsJson,
          nextIdsJson, '', 0);
      lst.Add(infoID + SEP_CHAR + previousID + SEP_CHAR + built);
    end;
  end;
end;

// ===== Quest =====
// stages（INDX + Log Entries の CNAM list）と objectives（QOBJ + NNAM）を embedded。

procedure ExtractQuest(quest: IInterface);
var
  questID, questName, questType, source, edid: string;
  stagesJson, objectivesJson, questEntry: string;
  stages, stageItem, logs, logItem, objectives, objItem: IInterface;
  i, j: integer;
  stageIdx, stageLog, objIdx, objText: string;
  logsJson, stageEntry, objEntry: string;
begin
  questID := HexFormID(quest);
  if IsProcessed(questID) then Exit;
  MarkProcessed(questID);
  questName := GetElementValue(quest, 'FULL');
  questType := GetElementValue(quest, 'DNAM\Quest Type');
  source := GetMasterFileName(quest);
  edid := GetCachedEdid(quest);

  // r13 stages（QuestStage 配列）。各 stage は log_entries list を持ち、各 log は
  // 独立に xTranslator 書き戻し可能（type = "QUST CNAM"）。
  stagesJson := '';
  if ElementExists(quest, 'Stages') then begin
    stages := ElementByName(quest, 'Stages');
    for i := 0 to ElementCount(stages) - 1 do begin
      stageItem := ElementByIndex(stages, i);
      logs := ElementByName(stageItem, 'INDX - Stage Index');
      if Assigned(logs) and (ElementCount(logs) > 0) then
        stageIdx := IntToStr(Integer(GetNativeValue(ElementByIndex(logs, 0))))
      else if Assigned(logs) then
        stageIdx := IntToStr(Integer(GetNativeValue(logs)))
      else
        stageIdx := '0';

      logsJson := '';
      if ElementExists(stageItem, 'Log Entries') then begin
        logs := ElementByName(stageItem, 'Log Entries');
        for j := 0 to ElementCount(logs) - 1 do begin
          logItem := ElementByIndex(logs, j);
          stageLog := GetElementValue(logItem, 'CNAM');
          if stageLog = '' then Continue;
          if logsJson <> '' then logsJson := logsJson + ',' + #13#10;
          logsJson := logsJson +
                      '          {' + #13#10 +
                      '            "log_index": ' + IntToStr(j) + ',' + #13#10 +
                      '            "type": "QUST CNAM",' + #13#10 +
                      '            "text": ' + JsonString(stageLog) + #13#10 +
                      '          }';
        end;
      end;
      if logsJson = '' then Continue;

      stageEntry := '      {' + #13#10 +
                    '        "stage_index": ' + stageIdx + ',' + #13#10 +
                    '        "log_entries": [' + #13#10 + logsJson + #13#10 +
                    '        ]' + #13#10 +
                    '      }';
      if stagesJson <> '' then stagesJson := stagesJson + ',' + #13#10;
      stagesJson := stagesJson + stageEntry;
    end;
  end;

  // r14 objectives（QuestObjective 配列、type = "QUST NNAM"）
  objectivesJson := '';
  if ElementExists(quest, 'Objectives') then begin
    objectives := ElementByName(quest, 'Objectives');
    for i := 0 to ElementCount(objectives) - 1 do begin
      objItem := ElementByIndex(objectives, i);
      objIdx := GetElementValue(objItem, 'QOBJ');
      objText := GetElementValue(objItem, 'NNAM');
      if objText = '' then Continue;
      objEntry := '      {' + #13#10 +
                  '        "objective_index": ' + JsonString(objIdx) + ',' + #13#10 +
                  '        "type": "QUST NNAM",' + #13#10 +
                  '        "display_text": ' + JsonString(objText) + #13#10 +
                  '      }';
      if objectivesJson <> '' then objectivesJson := objectivesJson + ',' + #13#10;
      objectivesJson := objectivesJson + objEntry;
    end;
  end;

  // 名前が無くても stage / objective を持つ quest（NNAM / CNAM のみ）は出力する。
  if (questName = '') and (stagesJson = '') and (objectivesJson = '') then Exit;

  questEntry := '  {' + #13#10 +
                JsonField('id', JsonString(questID)) + ',' + #13#10 +
                JsonField('editor_id', JsonString(edid)) + ',' + #13#10 +
                JsonField('type', JsonString('QUST FULL')) + ',' + #13#10 +
                JsonField('source', JsonString(source)) + ',' + #13#10 +
                JsonField('name', JsonString(questName)) + ',' + #13#10 +
                JsonField('quest_type', JsonString(questType)) + ',' + #13#10 +
                '    "stages": [';
  if stagesJson <> '' then
    questEntry := questEntry + #13#10 + stagesJson + #13#10 + '    ]'
  else
    questEntry := questEntry + ']';
  questEntry := questEntry + ',' + #13#10 + '    "objectives": [';
  if objectivesJson <> '' then
    questEntry := questEntry + #13#10 + objectivesJson + #13#10 + '    ]'
  else
    questEntry := questEntry + ']';
  questEntry := questEntry + #13#10 + '  }';
  questList.Add(questEntry);
end;

// ===== VoiceType (VTYP) =====
// v18 で独立 class 化。VTYP は FULL を持たず、EditorID が識別子兼ペルソナ手掛かり。

// EditorID 規約から声型カテゴリを判定する。
function VoiceCategory(edid: string): string;
begin
  if Copy(edid, 1, 2) = 'Cr' then Result := '生物'
  else if Pos('Unique', edid) > 0 then Result := 'キャラ専用'
  else Result := '汎用人間';
end;

procedure ExtractVoiceType(vtyp: IInterface);
var
  vID, edid, source, entry: string;
begin
  vID := HexFormID(vtyp);
  if IsProcessed(vID) then Exit;
  MarkProcessed(vID);
  edid := GetCachedEdid(vtyp);
  if edid = '' then Exit;
  source := GetMasterFileName(vtyp);

  // VTYP は翻訳対象テキストを持たない。identifier は EditorID。
  entry := '  {' + #13#10 +
           JsonField('id', JsonString(vID)) + ',' + #13#10 +
           JsonField('editor_id', JsonString(edid)) + ',' + #13#10 +
           JsonField('type', JsonString('VTYP EDID')) + ',' + #13#10 +
           JsonField('source', JsonString(source)) + ',' + #13#10 +
           JsonField('identifier', JsonString(edid)) + ',' + #13#10 +
           JsonField('category', JsonString(VoiceCategory(edid))) + #13#10 +
           '  }';
  voiceTypeList.Add(entry);
end;

// ===== MagicEffect =====

procedure ExtractMagicEffect(mgef: IInterface);
var
  mgefID, mgefName, mgefDesc, source, edid, entry: string;
begin
  mgefID := HexFormID(mgef);
  if IsProcessed(mgefID) then Exit;
  MarkProcessed(mgefID);
  mgefName := GetElementValue(mgef, 'FULL');
  mgefDesc := GetElementValue(mgef, 'DNAM');
  if mgefName = '' then Exit;
  source := GetMasterFileName(mgef);
  edid := GetCachedEdid(mgef);

  // primary type = "MGEF FULL"。description は MGEF DNAM mapping。
  entry := '  {' + #13#10 +
           JsonField('id', JsonString(mgefID)) + ',' + #13#10 +
           JsonField('editor_id', JsonString(edid)) + ',' + #13#10 +
           JsonField('type', JsonString('MGEF FULL')) + ',' + #13#10 +
           JsonField('source', JsonString(source)) + ',' + #13#10 +
           JsonField('name', JsonString(mgefName)) + ',' + #13#10 +
           JsonField('description', JsonString(mgefDesc)) + #13#10 +
           '  }';
  magicEffectList.Add(entry);
end;

// Effects[] 配列を辿り、各 entry の EFID で MGEF を解決。
// Ensure と magic_effect_ids 配列構築を返す。
function ResolveEffectsToIds(carrier: IInterface): string;
var
  effects, effectEntry, mgefRec: IInterface;
  i: integer;
  mgefID: string;
begin
  Result := '';
  if not ElementExists(carrier, 'Effects') then Exit;
  effects := ElementByName(carrier, 'Effects');
  for i := 0 to ElementCount(effects) - 1 do begin
    effectEntry := ElementByIndex(effects, i);
    mgefRec := LinksTo(ElementByPath(effectEntry, 'EFID'));
    if not Assigned(mgefRec) then Continue;
    mgefID := HexFormID(mgefRec);
    EnsureRecord(mgefRec);
    if Result <> '' then Result := Result + ', ';
    Result := Result + JsonString(mgefID);
  end;
end;

// ===== Magic (SPEL) =====

procedure ExtractMagic(spel: IInterface);
var
  spelID, spelName, spelDesc, source, edid, mgefIdsJson, entry: string;
begin
  spelID := HexFormID(spel);
  if IsProcessed(spelID) then Exit;
  MarkProcessed(spelID);
  spelName := GetElementValue(spel, 'FULL');
  spelDesc := GetElementValue(spel, 'DESC');
  if spelDesc = '' then spelDesc := GetElementValue(spel, 'DNAM');
  if spelName = '' then Exit;
  source := GetMasterFileName(spel);
  edid := GetCachedEdid(spel);
  mgefIdsJson := ResolveEffectsToIds(spel);

  entry := '  {' + #13#10 +
           JsonField('id', JsonString(spelID)) + ',' + #13#10 +
           JsonField('editor_id', JsonString(edid)) + ',' + #13#10 +
           JsonField('type', JsonString('SPEL FULL')) + ',' + #13#10 +
           JsonField('source', JsonString(source)) + ',' + #13#10 +
           JsonField('name', JsonString(spelName)) + ',' + #13#10 +
           JsonField('description', JsonString(spelDesc)) + ',' + #13#10 +
           '    "magic_effect_ids": [' + mgefIdsJson + ']' + #13#10 +
           '  }';
  magicList.Add(entry);
end;

// ===== Enchantment (ENCH) =====

procedure ExtractEnchantment(ench: IInterface);
var
  enchID, enchName, enchDesc, source, edid, mgefIdsJson, entry: string;
begin
  enchID := HexFormID(ench);
  if IsProcessed(enchID) then Exit;
  MarkProcessed(enchID);
  enchName := GetElementValue(ench, 'FULL');
  enchDesc := GetElementValue(ench, 'DESC');
  if enchDesc = '' then enchDesc := GetElementValue(ench, 'DNAM');
  source := GetMasterFileName(ench);
  edid := GetCachedEdid(ench);
  mgefIdsJson := ResolveEffectsToIds(ench);

  // ENCH は名称空でも装備品付与で意味を持つので名無しでも emit する。
  entry := '  {' + #13#10 +
           JsonField('id', JsonString(enchID)) + ',' + #13#10 +
           JsonField('editor_id', JsonString(edid)) + ',' + #13#10 +
           JsonField('type', JsonString('ENCH FULL')) + ',' + #13#10 +
           JsonField('source', JsonString(source)) + ',' + #13#10 +
           JsonField('name', JsonString(enchName)) + ',' + #13#10 +
           JsonField('description', JsonString(enchDesc)) + ',' + #13#10 +
           '    "magic_effect_ids": [' + mgefIdsJson + ']' + #13#10 +
           '  }';
  enchantmentList.Add(entry);
end;

// ===== Consumable (SCRL / ALCH / INGR) =====

procedure ExtractConsumable(c: IInterface);
var
  cID, cName, cDesc, sig, source, edid, mgefIdsJson, entry: string;
begin
  cID := HexFormID(c);
  if IsProcessed(cID) then Exit;
  MarkProcessed(cID);
  sig := Signature(c);
  cName := GetElementValue(c, 'FULL');
  cDesc := GetElementValue(c, 'DESC');
  if cName = '' then Exit;
  source := GetMasterFileName(c);
  edid := GetCachedEdid(c);
  mgefIdsJson := ResolveEffectsToIds(c);

  // primary type = "<SIG> FULL"。description は "<SIG> DESC" mapping。
  entry := '  {' + #13#10 +
           JsonField('id', JsonString(cID)) + ',' + #13#10 +
           JsonField('editor_id', JsonString(edid)) + ',' + #13#10 +
           JsonField('type', JsonString(sig + ' FULL')) + ',' + #13#10 +
           JsonField('source', JsonString(source)) + ',' + #13#10 +
           JsonField('name', JsonString(cName)) + ',' + #13#10 +
           JsonField('description', JsonString(cDesc)) + ',' + #13#10 +
           JsonField('kind', JsonString(sig)) + ',' + #13#10 +
           '    "magic_effect_ids": [' + mgefIdsJson + ']' + #13#10 +
           '  }';
  consumableList.Add(entry);
end;

// ===== Equipment (WEAP / ARMO / AMMO) =====

procedure ExtractEquipment(e: IInterface);
var
  eqID, eqName, eqDesc, sig, source, edid, enchID, entry: string;
  enchRec: IInterface;
begin
  eqID := HexFormID(e);
  if IsProcessed(eqID) then Exit;
  MarkProcessed(eqID);
  sig := Signature(e);
  eqName := GetElementValue(e, 'FULL');
  eqDesc := GetElementValue(e, 'DESC');
  if eqName = '' then Exit;
  source := GetMasterFileName(e);
  edid := GetCachedEdid(e);

  // r12: EITM が指す ENCH を解決
  enchID := '';
  enchRec := LinksTo(ElementByName(e, 'EITM'));
  if Assigned(enchRec) then begin
    enchID := HexFormID(enchRec);
    EnsureRecord(enchRec);
  end;

  // primary type = "<SIG> FULL"。description は "<SIG> DESC" mapping（多くは空）。
  entry := '  {' + #13#10 +
           JsonField('id', JsonString(eqID)) + ',' + #13#10 +
           JsonField('editor_id', JsonString(edid)) + ',' + #13#10 +
           JsonField('type', JsonString(sig + ' FULL')) + ',' + #13#10 +
           JsonField('source', JsonString(source)) + ',' + #13#10 +
           JsonField('name', JsonString(eqName)) + ',' + #13#10 +
           JsonField('description', JsonString(eqDesc)) + ',' + #13#10 +
           JsonField('kind', JsonString(sig)) + ',' + #13#10 +
           JsonField('enchantment_id', JsonString(enchID)) + #13#10 +
           '  }';
  equipmentList.Add(entry);
end;

// ===== Item =====
// KEYM / MISC / LIGH / CONT / SLGM / DOOR / FURN
procedure ExtractItem(i: IInterface);
var
  iID, iName, sig, source, edid, entry: string;
begin
  iID := HexFormID(i);
  if IsProcessed(iID) then Exit;
  MarkProcessed(iID);
  sig := Signature(i);
  iName := GetElementValue(i, 'FULL');
  if iName = '' then Exit;
  source := GetMasterFileName(i);
  edid := GetCachedEdid(i);

  entry := '  {' + #13#10 +
           JsonField('id', JsonString(iID)) + ',' + #13#10 +
           JsonField('editor_id', JsonString(edid)) + ',' + #13#10 +
           JsonField('type', JsonString(sig + ' FULL')) + ',' + #13#10 +
           JsonField('source', JsonString(source)) + ',' + #13#10 +
           JsonField('name', JsonString(iName)) + ',' + #13#10 +
           JsonField('kind', JsonString(sig)) + #13#10 +
           '  }';
  itemList.Add(entry);
end;

// ===== Activator =====
// ACTI / FLOR / TREE。起動・採取可能なオブジェクト（item ではない）。
// FULL = クロスヘア表示名、RNAM = 起動動作テキスト（"開ける"/"採取"等）。種別は signature。
procedure ExtractActivator(a: IInterface);
var
  aID, aName, aAction, sig, source, edid, entry: string;
begin
  aID := HexFormID(a);
  if IsProcessed(aID) then Exit;
  MarkProcessed(aID);
  sig := Signature(a);
  aName := GetElementValue(a, 'FULL');
  aAction := GetElementValue(a, 'RNAM');
  if (aName = '') and (aAction = '') then Exit;
  source := GetMasterFileName(a);
  edid := GetCachedEdid(a);

  // primary type = "<SIG> FULL"。activate_action は "<SIG> RNAM" mapping。
  entry := '  {' + #13#10 +
           JsonField('id', JsonString(aID)) + ',' + #13#10 +
           JsonField('editor_id', JsonString(edid)) + ',' + #13#10 +
           JsonField('type', JsonString(sig + ' FULL')) + ',' + #13#10 +
           JsonField('source', JsonString(source)) + ',' + #13#10 +
           JsonField('name', JsonString(aName)) + ',' + #13#10 +
           JsonField('activate_action', JsonString(aAction)) + ',' + #13#10 +
           JsonField('kind', JsonString(sig)) + #13#10 +
           '  }';
  activatorList.Add(entry);
end;

// ===== Book =====

procedure ExtractBook(book: IInterface);
var
  bookID, title, body, author, source, edid, entry: string;
begin
  bookID := HexFormID(book);
  if IsProcessed(bookID) then Exit;
  MarkProcessed(bookID);
  title := GetElementValue(book, 'FULL');
  body := GetElementValue(book, 'DESC');
  author := GetElementValue(book, 'CNAM');
  if (title = '') and (body = '') then Exit;
  source := GetMasterFileName(book);
  edid := GetCachedEdid(book);

  // primary type = "BOOK FULL"。body は "BOOK DESC"、author は "BOOK CNAM" mapping。
  entry := '  {' + #13#10 +
           JsonField('id', JsonString(bookID)) + ',' + #13#10 +
           JsonField('editor_id', JsonString(edid)) + ',' + #13#10 +
           JsonField('type', JsonString('BOOK FULL')) + ',' + #13#10 +
           JsonField('source', JsonString(source)) + ',' + #13#10 +
           JsonField('title', JsonString(title)) + ',' + #13#10 +
           JsonField('body', JsonString(body)) + ',' + #13#10 +
           JsonField('author', JsonString(author)) + #13#10 +
           '  }';
  bookList.Add(entry);
end;

// ===== Word (WOOP) =====
// Word は Shout の composition として embedded で emit する（r8）。
// 個別 WOOP record の参照経路は現状無いため、standalone Extract は持たない。

// Word を Shout 内で embedded として組み立てる helper
function BuildShoutWordsJson(shou: IInterface): string;
var
  words, wordEntry, woopRec: IInterface;
  i: integer;
  woopID, dragonSpelling, meaning, source, edid, entry: string;
begin
  Result := '';
  words := ElementByName(shou, 'Words of Power');  // SSEEdit の配列名（'Words' ではない）
  if not Assigned(words) then Exit;
  // 配列要素は SNAM。各要素の 'Word' field が WOOP への link（'WNAM' では取れない）。
  for i := 0 to ElementCount(words) - 1 do begin
    wordEntry := ElementByIndex(words, i);
    woopRec := LinksTo(ElementByName(wordEntry, 'Word'));
    if not Assigned(woopRec) then Continue;

    woopID := HexFormID(woopRec);
    dragonSpelling := GetElementValue(woopRec, 'FULL');  // 龍語綴り、翻訳禁止
    meaning := GetElementValue(woopRec, 'TNAM');          // 意味の訳、翻訳対象
    source := GetMasterFileName(woopRec);
    edid := GetCachedEdid(woopRec);
    MarkProcessed(woopID);

    // Word の primary translation target = "WOOP TNAM"（意味の訳）。
    // dragon_spelling（WOOP FULL）は龍文字表示用で翻訳禁止のため field は出すが対象にしない。
    entry := '      {' + #13#10 +
             '        "id": ' + JsonString(woopID) + ',' + #13#10 +
             '        "editor_id": ' + JsonString(edid) + ',' + #13#10 +
             '        "type": "WOOP TNAM",' + #13#10 +
             '        "source": ' + JsonString(source) + ',' + #13#10 +
             '        "dragon_spelling": ' + JsonString(dragonSpelling) + ',' + #13#10 +
             '        "translation": ' + JsonString(meaning) + #13#10 +
             '      }';
    if Result <> '' then Result := Result + ',' + #13#10;
    Result := Result + entry;
  end;
end;

// ===== Shout (SHOU) =====

procedure ExtractShout(shou: IInterface);
var
  shoutID, shoutName, shoutDesc, source, edid, wordsJson, entry: string;
begin
  shoutID := HexFormID(shou);
  if IsProcessed(shoutID) then Exit;
  MarkProcessed(shoutID);
  shoutName := GetElementValue(shou, 'FULL');
  shoutDesc := GetElementValue(shou, 'DESC');
  if shoutName = '' then Exit;
  source := GetMasterFileName(shou);
  edid := GetCachedEdid(shou);
  wordsJson := BuildShoutWordsJson(shou);

  entry := '  {' + #13#10 +
           JsonField('id', JsonString(shoutID)) + ',' + #13#10 +
           JsonField('editor_id', JsonString(edid)) + ',' + #13#10 +
           JsonField('type', JsonString('SHOU FULL')) + ',' + #13#10 +
           JsonField('source', JsonString(source)) + ',' + #13#10 +
           JsonField('name', JsonString(shoutName)) + ',' + #13#10 +
           JsonField('description', JsonString(shoutDesc)) + ',' + #13#10 +
           '    "words": [';
  if wordsJson <> '' then
    entry := entry + #13#10 + wordsJson + #13#10 + '    ]'
  else
    entry := entry + ']';
  entry := entry + #13#10 + '  }';
  shoutList.Add(entry);
end;

// ===== Location =====

procedure ExtractLocation(loc: IInterface);
var
  locID, locName, parentID, sig, source, edid, entry: string;
begin
  locID := HexFormID(loc);
  if IsProcessed(locID) then Exit;
  MarkProcessed(locID);
  sig := Signature(loc);
  locName := GetElementValue(loc, 'FULL');
  if locName = '' then Exit;
  if sig = 'LCTN' then parentID := GetElementValue(loc, 'PNAM')
  else if sig = 'WRLD' then parentID := GetElementValue(loc, 'WNAM')
  else if sig = 'CELL' then parentID := GetElementValue(loc, 'XLCN');
  source := GetMasterFileName(loc);
  edid := GetCachedEdid(loc);

  entry := '  {' + #13#10 +
           JsonField('id', JsonString(locID)) + ',' + #13#10 +
           JsonField('editor_id', JsonString(edid)) + ',' + #13#10 +
           JsonField('type', JsonString(sig + ' FULL')) + ',' + #13#10 +
           JsonField('source', JsonString(source)) + ',' + #13#10 +
           JsonField('name', JsonString(locName)) + ',' + #13#10 +
           JsonField('kind', JsonString(sig)) + ',' + #13#10 +
           JsonField('parent_id', JsonString(parentID)) + #13#10 +
           '  }';
  locationList.Add(entry);
end;

// ===== Message =====

procedure ExtractMessage(msg: IInterface);
var
  msgID, title, body, questID, source, edid, entry, buttonsJson, btnText: string;
  questRec, buttons, btnItem: IInterface;
  i: integer;
begin
  msgID := HexFormID(msg);
  if IsProcessed(msgID) then Exit;
  MarkProcessed(msgID);
  body := GetElementValue(msg, 'DESC');
  title := GetElementValue(msg, 'FULL');
  source := GetMasterFileName(msg);
  edid := GetCachedEdid(msg);
  questRec := LinksTo(ElementByName(msg, 'QNAM'));
  if Assigned(questRec) then begin
    questID := HexFormID(questRec);
    EnsureRecord(questRec);
  end else
    questID := GetElementValue(msg, 'QNAM');

  // 選択肢ボタン（ITXT）。各 button は独立に書き戻し可能（type = "MESG ITXT"）。
  buttonsJson := '';
  if ElementExists(msg, 'Menu Buttons') then begin
    buttons := ElementByName(msg, 'Menu Buttons');
    for i := 0 to ElementCount(buttons) - 1 do begin
      btnItem := ElementByIndex(buttons, i);
      btnText := GetElementValue(btnItem, 'ITXT');
      if btnText = '' then Continue;
      if buttonsJson <> '' then buttonsJson := buttonsJson + ',' + #13#10;
      buttonsJson := buttonsJson +
                     '      {' + #13#10 +
                     '        "button_index": ' + IntToStr(i) + ',' + #13#10 +
                     '        "type": "MESG ITXT",' + #13#10 +
                     '        "text": ' + JsonString(btnText) + #13#10 +
                     '      }';
    end;
  end;

  // 翻訳対象テキスト（本文 / 題名 / 選択肢）が全て無ければ出力しない。
  if (body = '') and (title = '') and (buttonsJson = '') then Exit;

  // primary type = "MESG DESC"（v1 互換、body が主翻訳対象）。title は "MESG FULL" mapping。
  entry := '  {' + #13#10 +
           JsonField('id', JsonString(msgID)) + ',' + #13#10 +
           JsonField('editor_id', JsonString(edid)) + ',' + #13#10 +
           JsonField('type', JsonString('MESG DESC')) + ',' + #13#10 +
           JsonField('source', JsonString(source)) + ',' + #13#10 +
           JsonField('title', JsonString(title)) + ',' + #13#10 +
           JsonField('body', JsonString(body)) + ',' + #13#10 +
           JsonField('quest_id', JsonString(questID)) + ',' + #13#10 +
           '    "buttons": [';
  if buttonsJson <> '' then
    entry := entry + #13#10 + buttonsJson + #13#10 + '    ]'
  else
    entry := entry + ']';
  entry := entry + #13#10 + '  }';
  messageList.Add(entry);
end;

// ===== LoadingScreen =====

procedure ExtractLoadingScreen(lscr: IInterface);
var
  lscrID, body, source, edid, entry: string;
begin
  lscrID := HexFormID(lscr);
  if IsProcessed(lscrID) then Exit;
  MarkProcessed(lscrID);
  body := GetElementValue(lscr, 'DESC');
  if body = '' then Exit;
  source := GetMasterFileName(lscr);
  edid := GetCachedEdid(lscr);

  entry := '  {' + #13#10 +
           JsonField('id', JsonString(lscrID)) + ',' + #13#10 +
           JsonField('editor_id', JsonString(edid)) + ',' + #13#10 +
           JsonField('type', JsonString('LSCR DESC')) + ',' + #13#10 +
           JsonField('source', JsonString(source)) + ',' + #13#10 +
           JsonField('body', JsonString(body)) + #13#10 +
           '  }';
  loadingScreenList.Add(entry);
end;

// ===== Perk =====

procedure ExtractPerk(perk: IInterface);
var
  perkID, perkName, perkDesc, source, edid, entry: string;
begin
  perkID := HexFormID(perk);
  if IsProcessed(perkID) then Exit;
  MarkProcessed(perkID);
  perkName := GetElementValue(perk, 'FULL');
  perkDesc := GetElementValue(perk, 'DESC');
  // 名前が無くても説明（DESC）を持つ perk は出力する。
  if (perkName = '') and (perkDesc = '') then Exit;
  source := GetMasterFileName(perk);
  edid := GetCachedEdid(perk);

  // primary type = "PERK FULL"。description は "PERK DESC" mapping。
  entry := '  {' + #13#10 +
           JsonField('id', JsonString(perkID)) + ',' + #13#10 +
           JsonField('editor_id', JsonString(edid)) + ',' + #13#10 +
           JsonField('type', JsonString('PERK FULL')) + ',' + #13#10 +
           JsonField('source', JsonString(source)) + ',' + #13#10 +
           JsonField('name', JsonString(perkName)) + ',' + #13#10 +
           JsonField('description', JsonString(perkDesc)) + #13#10 +
           '  }';
  perkList.Add(entry);
end;

// ===== EnsureRecord 実装 =====
// 参照解決時に未処理 record を自動 dispatch する。
// cross-file（master plugin）の record もこの経路で取り込まれる。

procedure EnsureRecord(e: IInterface);
var
  sig: string;
begin
  if not Assigned(e) then Exit;
  // 対象 plugin が新規追加 / 上書き していない参照先 record（純 master）は翻訳対象外。
  // GetMasterFileName は GetFile ベースで、override 要素なら対象 plugin を返す。
  // 参照先が対象 plugin の外にしか存在しない場合（vanilla 定義そのまま）は出力しない。
  // ID 参照（speaker_ids 等）は文字列として残り、翻訳時に元 plugin の辞書で解決する。
  if GetMasterFileName(e) <> targetFileName then Exit;
  if IsProcessed(HexFormID(e)) then Exit;
  sig := Signature(e);

  if sig = 'NPC_' then ExtractSpeaker(e)
  else if sig = 'RACE' then ExtractRace(e)
  else if sig = 'FACT' then ExtractFaction(e)
  else if sig = 'QUST' then ExtractQuest(e)
  else if sig = 'MGEF' then ExtractMagicEffect(e)
  else if sig = 'SPEL' then ExtractMagic(e)
  else if sig = 'ENCH' then ExtractEnchantment(e)
  else if sig = 'SHOU' then ExtractShout(e)
  else if (sig = 'SCRL') or (sig = 'ALCH') or (sig = 'INGR') then ExtractConsumable(e)
  else if (sig = 'WEAP') or (sig = 'ARMO') or (sig = 'AMMO') then ExtractEquipment(e)
  else if sig = 'BOOK' then ExtractBook(e)
  else if (sig = 'KEYM') or (sig = 'MISC') or (sig = 'LIGH') or (sig = 'CONT') or
          (sig = 'SLGM') or (sig = 'DOOR') or (sig = 'FURN') then ExtractItem(e)
  else if (sig = 'ACTI') or (sig = 'FLOR') or (sig = 'TREE') then ExtractActivator(e)
  else if (sig = 'CELL') or (sig = 'LCTN') or (sig = 'WRLD') then ExtractLocation(e)
  else if sig = 'MESG' then ExtractMessage(e)
  else if sig = 'LSCR' then ExtractLoadingScreen(e)
  else if sig = 'PERK' then ExtractPerk(e)
  else if sig = 'VTYP' then ExtractVoiceType(e)
  else if sig = 'TACT' then ExtractSpeaker(e);
  // DIAL / INFO は file iteration で必ず到達するため Ensure 経路には乗せない。
end;

// ===== Process / Initialize / Finalize =====

function Initialize: integer;
begin
  dialogueList := TStringList.Create;
  personList := TStringList.Create;
  raceList := TStringList.Create;
  factionList := TStringList.Create;
  questList := TStringList.Create;
  locationList := TStringList.Create;
  itemList := TStringList.Create;
  activatorList := TStringList.Create;
  equipmentList := TStringList.Create;
  consumableList := TStringList.Create;
  magicList := TStringList.Create;
  enchantmentList := TStringList.Create;
  magicEffectList := TStringList.Create;
  shoutList := TStringList.Create;
  bookList := TStringList.Create;
  messageList := TStringList.Create;
  loadingScreenList := TStringList.Create;
  perkList := TStringList.Create;
  voiceTypeList := TStringList.Create;

  dialDataMap := TStringList.Create;
  dialDataMap.Sorted := True;
  dialDataMap.Duplicates := dupError;

  processedIds := TStringList.Create;
  processedIds.Sorted := True;
  processedIds.Duplicates := dupIgnore;

  lstRecursion := TList.Create;
  targetFileName := '';
  AddMessage('[ExportDataV2] Initialized.');
  Result := 0;
end;

function Process(e: IInterface): integer;
var
  sig: string;
begin
  if targetFileName = '' then
    targetFileName := GetFileNameWithExtension(GetFile(e));

  sig := Signature(e);
  if sig = 'DIAL' then ExtractDialogue(e)
  else if sig = 'INFO' then ExtractInfo(e)
  else if sig = 'NPC_' then ExtractSpeaker(e)
  else if sig = 'RACE' then ExtractRace(e)
  else if sig = 'FACT' then ExtractFaction(e)
  else if sig = 'QUST' then ExtractQuest(e)
  else if sig = 'MGEF' then ExtractMagicEffect(e)
  else if sig = 'SPEL' then ExtractMagic(e)
  else if sig = 'ENCH' then ExtractEnchantment(e)
  else if sig = 'SHOU' then ExtractShout(e)
  else if (sig = 'SCRL') or (sig = 'ALCH') or (sig = 'INGR') then ExtractConsumable(e)
  else if (sig = 'WEAP') or (sig = 'ARMO') or (sig = 'AMMO') then ExtractEquipment(e)
  else if sig = 'BOOK' then ExtractBook(e)
  else if (sig = 'KEYM') or (sig = 'MISC') or (sig = 'LIGH') or (sig = 'CONT') or
          (sig = 'SLGM') or (sig = 'DOOR') or (sig = 'FURN') then ExtractItem(e)
  else if (sig = 'ACTI') or (sig = 'FLOR') or (sig = 'TREE') then ExtractActivator(e)
  else if (sig = 'CELL') or (sig = 'LCTN') or (sig = 'WRLD') then ExtractLocation(e)
  else if sig = 'MESG' then ExtractMessage(e)
  else if sig = 'LSCR' then ExtractLoadingScreen(e)
  else if sig = 'PERK' then ExtractPerk(e)
  else if sig = 'VTYP' then ExtractVoiceType(e)
  else if sig = 'TACT' then ExtractSpeaker(e);

  Result := 0;
end;

// PNAM chain で response をソートする。id/pid は行 prefix から取り JSON 再 parse はしない。
// 格納は TStringList ベース（v1 実績）。JvInterpreter で未実証の TList / Pointer cast は使わない。
// children/queue には prefix 付きの行そのものを入れる。Finalize 側で prefix を剥がすため整合する。
procedure ApplyInfoSorting(var respList: TStringList);
var
  i, j, k, p1, p2: Integer;
  pidMap: TStringList;
  sortedList, queue, children, visited, expandedInfo: TStringList;
  curLine, curId, curPid, rest: string;
begin
  if respList.Count <= 2 then Exit;

  pidMap := TStringList.Create;
  pidMap.Sorted := True;
  pidMap.Duplicates := dupIgnore;

  // previous_id -> その pid を先行に持つ response 行（prefix 付き）の list を作る。
  for i := 1 to respList.Count - 1 do begin
    curLine := respList[i];
    p1 := Pos(SEP_CHAR, curLine);
    if p1 <= 0 then Continue;
    rest := Copy(curLine, p1 + 1, Length(curLine));
    p2 := Pos(SEP_CHAR, rest);
    if p2 <= 0 then Continue;
    curPid := Copy(rest, 1, p2 - 1);

    j := pidMap.IndexOf(curPid);
    if j < 0 then begin
      children := TStringList.Create;
      pidMap.AddObject(curPid, children);
    end else
      children := TStringList(pidMap.Objects[j]);
    children.Add(curLine);
  end;

  sortedList := TStringList.Create;
  sortedList.Add(respList[0]);  // header を保持する。

  queue := TStringList.Create;
  visited := TStringList.Create;
  visited.Sorted := True;
  visited.Duplicates := dupIgnore;
  // 後続展開を INFO（infoID）単位で 1 回に抑えるための集合。
  expandedInfo := TStringList.Create;
  expandedInfo.Sorted := True;
  expandedInfo.Duplicates := dupIgnore;

  // root = previous_id が空、または 00000000 の response。
  j := pidMap.IndexOf('');
  if j >= 0 then queue.AddStrings(TStringList(pidMap.Objects[j]));
  j := pidMap.IndexOf('00000000');
  if j >= 0 then queue.AddStrings(TStringList(pidMap.Objects[j]));

  // BFS で chain を辿る。visited で循環と二重追加を防ぐ。
  i := 0;
  while i < queue.Count do begin
    curLine := queue[i];
    p1 := Pos(SEP_CHAR, curLine);
    if p1 > 0 then begin
      curId := Copy(curLine, 1, p1 - 1);
      // visited は行全体で判定する。同一 INFO の複数 NAM1（response）は order が違い
      // 行が異なるため、すべて残す。infoID 単位で判定すると 2 つ目以降を落とす。
      if visited.IndexOf(curLine) < 0 then begin
        sortedList.Add(curLine);
        visited.Add(curLine);
        // 後続展開（PNAM chain）は INFO 単位で 1 回だけ行う。
        if expandedInfo.IndexOf(curId) < 0 then begin
          expandedInfo.Add(curId);
          j := pidMap.IndexOf(curId);
          if j >= 0 then begin
            children := TStringList(pidMap.Objects[j]);
            for k := 0 to children.Count - 1 do
              queue.Add(children[k]);
          end;
        end;
      end;
    end;
    Inc(i);
  end;

  // chain に未接続の response を元順で末尾に足す。
  if sortedList.Count < respList.Count then
    for i := 1 to respList.Count - 1 do begin
      curLine := respList[i];
      if visited.IndexOf(curLine) < 0 then begin
        sortedList.Add(curLine);
        visited.Add(curLine);
      end;
    end;

  // 反転はしない。BFS が root（PNAM 空）→ 後続の自然順で積み、同一 INFO 内の
  // 複数 NAM1 も order 昇順で並ぶため、そのまま採用する。
  respList.Assign(sortedList);

  queue.Free;
  visited.Free;
  expandedInfo.Free;
  sortedList.Free;
  for i := 0 to pidMap.Count - 1 do
    if pidMap.Objects[i] <> nil then
      TStringList(pidMap.Objects[i]).Free;
  pidMap.Free;
end;

procedure EmitArrayInto(target, source: TStringList);
var
  i: Integer;
begin
  for i := 0 to source.Count - 1 do begin
    if i < source.Count - 1 then
      target.Add(source[i] + ',')
    else
      target.Add(source[i]);
  end;
end;

function Finalize: integer;
var
  filename, filepath: string;
  outputFile: TStringList;
  dlg: TSaveDialog;
  i, j, p1, p2: integer;
  dialHeader, finalDialEntry, line, jsonOnly: string;
  innerList: TStringList;
begin
  if targetFileName = '' then begin
    AddMessage('[ExportDataV2] No records processed.');
    Result := 1;
    Exit;
  end;

  AddMessage('[ExportDataV2] Finalizing...');
  AddMessage('[ExportDataV2] Merging ' + IntToStr(dialDataMap.Count) + ' dialogue groups...');

  for i := 0 to dialDataMap.Count - 1 do begin
    if dialDataMap.Objects[i] = nil then Continue;
    innerList := TStringList(dialDataMap.Objects[i]);
    if innerList.Count = 0 then Continue;

    dialHeader := innerList[0];
    if dialHeader = '' then Continue;

    ApplyInfoSorting(innerList);

    finalDialEntry := dialHeader + #13#10 + '  "responses": [';
    if innerList.Count > 1 then begin
      for j := 1 to innerList.Count - 1 do begin
        line := innerList[j];
        p1 := Pos(SEP_CHAR, line);
        if p1 <= 0 then jsonOnly := line
        else begin
          p2 := Pos(SEP_CHAR, Copy(line, p1 + 1, Length(line)));
          if p2 <= 0 then jsonOnly := line
          else jsonOnly := Copy(line, p1 + p2 + 1, Length(line));
        end;
        finalDialEntry := finalDialEntry + #13#10 + jsonOnly;
        if j < innerList.Count - 1 then
          finalDialEntry := finalDialEntry + ',';
      end;
    end;
    finalDialEntry := finalDialEntry + #13#10 + '  ]' + #13#10 + '  }';
    dialogueList.Add(finalDialEntry);
  end;

  filename := targetFileName + '_Export.json';

  dlg := TSaveDialog.Create(nil);
  try
    dlg.Filter := 'JSON files (*.json)|*.json|All files (*.*)|*.*';
    dlg.DefaultExt := 'json';
    dlg.FileName := filename;
    dlg.InitialDir := ProgramPath + 'Edit Scripts\';
    if dlg.Execute then
      filepath := dlg.FileName
    else begin
      AddMessage('[ExportDataV2] Export cancelled.');
      Result := 1;
      Exit;
    end;
  finally
    dlg.Free;
  end;

  outputFile := TStringList.Create;
  try
    outputFile.Add('{');
    outputFile.Add('  "target_plugin": ' + JsonString(targetFileName) + ',');

    outputFile.Add('  "dialogues": [');
    EmitArrayInto(outputFile, dialogueList);
    outputFile.Add('  ],');

    outputFile.Add('  "quests": [');
    EmitArrayInto(outputFile, questList);
    outputFile.Add('  ],');

    outputFile.Add('  "races": [');
    EmitArrayInto(outputFile, raceList);
    outputFile.Add('  ],');

    outputFile.Add('  "factions": [');
    EmitArrayInto(outputFile, factionList);
    outputFile.Add('  ],');

    outputFile.Add('  "items": [');
    EmitArrayInto(outputFile, itemList);
    outputFile.Add('  ],');

    outputFile.Add('  "activators": [');
    EmitArrayInto(outputFile, activatorList);
    outputFile.Add('  ],');

    outputFile.Add('  "equipment": [');
    EmitArrayInto(outputFile, equipmentList);
    outputFile.Add('  ],');

    outputFile.Add('  "consumables": [');
    EmitArrayInto(outputFile, consumableList);
    outputFile.Add('  ],');

    outputFile.Add('  "magic": [');
    EmitArrayInto(outputFile, magicList);
    outputFile.Add('  ],');

    outputFile.Add('  "enchantments": [');
    EmitArrayInto(outputFile, enchantmentList);
    outputFile.Add('  ],');

    outputFile.Add('  "magic_effects": [');
    EmitArrayInto(outputFile, magicEffectList);
    outputFile.Add('  ],');

    outputFile.Add('  "shouts": [');
    EmitArrayInto(outputFile, shoutList);
    outputFile.Add('  ],');

    outputFile.Add('  "books": [');
    EmitArrayInto(outputFile, bookList);
    outputFile.Add('  ],');

    outputFile.Add('  "locations": [');
    EmitArrayInto(outputFile, locationList);
    outputFile.Add('  ],');

    outputFile.Add('  "messages": [');
    EmitArrayInto(outputFile, messageList);
    outputFile.Add('  ],');

    outputFile.Add('  "loading_screens": [');
    EmitArrayInto(outputFile, loadingScreenList);
    outputFile.Add('  ],');

    outputFile.Add('  "perks": [');
    EmitArrayInto(outputFile, perkList);
    outputFile.Add('  ],');

    outputFile.Add('  "voice_types": [');
    EmitArrayInto(outputFile, voiceTypeList);
    outputFile.Add('  ],');

    outputFile.Add('  "speakers": {');
    EmitArrayInto(outputFile, personList);
    outputFile.Add('  }');

    outputFile.Add('}');
    outputFile.SaveToFile(filepath);
    AddMessage('[ExportDataV2] Export complete: ' + filename);
  finally
    outputFile.Free;
    for i := 0 to dialDataMap.Count - 1 do
      if dialDataMap.Objects[i] <> nil then
        TObject(dialDataMap.Objects[i]).Free;
    dialDataMap.Free;

    dialogueList.Free;
    personList.Free;
    raceList.Free;
    factionList.Free;
    questList.Free;
    locationList.Free;
    itemList.Free;
    activatorList.Free;
    equipmentList.Free;
    consumableList.Free;
    magicList.Free;
    enchantmentList.Free;
    magicEffectList.Free;
    shoutList.Free;
    bookList.Free;
    messageList.Free;
    loadingScreenList.Free;
    perkList.Free;
    voiceTypeList.Free;
    processedIds.Free;
    lstRecursion.Free;
  end;

  Result := 0;
end;

end.
