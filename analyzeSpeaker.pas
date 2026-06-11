{ ========================================================================
   Speaker Resolution As-Is Analyzer
   docs/conceptual-model.md 見直しの前段として、Dialogue Response（INFO）の
   「話者がどう決まるか」を raw record から網羅集計する調査専用 script。
   抽出（extractData.v2.pas）は変更しない。出力は xEdit の message log のみ。

   集計の狙い:
     - 話者の参照先が NPC_ / LVLN / FLST / VTYP / FACT 等のどの record type か。
     - NPC_ の場合、名称（FULL/SHRT）と TPLT（形態 record の本体参照）の有無。
     - ANAM-Speaker 直接指定、conditions（CTDA）の GetIsID、quest alias
       （GetIsAliasRef）の各経路の内訳。
     - 上記を組み合わせた INFO 単位の「話者決定パターン」の分布。

   JvInterpreter 制約への対応（extractData.v2.pas の作業で判明した点）:
     - forward 宣言は使わない（呼ぶ側より前に helper を定義する）。
     - Ord(s[i]) 直接や WideChar は使わない（本 script は文字処理を持たない）。
     - Pointer() cast は使わない。カウンタは TStringList.Values（key=数値）で持つ。
     - TStringList.Values[key] への代入は Sorted=True だと例外。非ソートで保持し、
       表示直前に Sort する。
   ======================================================================== }

unit AnalyzeSpeaker;

var
  counters: TStringList;   // 'ラベル=出現数' を Values 形式で保持する集計表。
  samples: TStringList;    // 名前空話者の分類別 EditorID サンプル（役割名の傾向確認用）。
  totalInfo: Integer;      // 処理した INFO 件数。

// ===== カウンタ helper =====

// 数字文字列を整数へ。空なら 0。Values 未登録 key の既定値に使う。
function ToIntDef0(s: string): Integer;
begin
  if s = '' then Result := 0
  else Result := StrToInt(s);
end;

// ラベルの出現数を 1 増やす。Pointer cast を避けるため Values に数値文字列で持つ。
procedure IncCounter(key: string);
begin
  counters.Values[key] := IntToStr(ToIntDef0(counters.Values[key]) + 1);
end;

// 参照先 record の signature。nil は明示する。
function SafeSig(e: IInterface): string;
begin
  if Assigned(e) then Result := Signature(e)
  else Result := '(nil)';
end;

// path の値が空でないか。record が nil なら False。
function HasVal(e: IInterface; path: string): Boolean;
begin
  Result := False;
  if Assigned(e) then
    Result := GetElementEditValues(e, path) <> '';
end;

// REFR / ACHR（placed reference）なら base record（NPC_ 等）へ解決する。
function ResolveActorBase(e: IInterface): IInterface;
var sig: string;
begin
  Result := e;
  if not Assigned(e) then Exit;
  sig := Signature(e);
  if (sig = 'REFR') or (sig = 'ACHR') then
    Result := WinningOverride(BaseRecord(e));
end;

// 名前空 actor のペルソナ名解決を分類する。FULL 空の actor がどの手掛かりで
// ペルソナ名を得られるかを判定する。
//   'morph_named'  : TPLT 連鎖の末端が FULL を持つ NPC_（形態 actor、本体名がペルソナ）
//   'via_lvln'     : TPLT 連鎖の途中に LVLN（レベルドリスト、汎用プール）
//   'generic_voice': 名前に届かないが VTCK（声型）を持つ（汎用 actor、声型名がペルソナ）
//   'editor_only'  : TPLT も VTCK も無い（EditorID のみが手掛かり）
function ClassifyUnnamedActor(npc: IInterface): string;
var
  cur, tpl: IInterface;
  depth: Integer;
  sig: string;
begin
  cur := npc;
  depth := 0;
  while Assigned(cur) and (depth < 12) do begin
    sig := Signature(cur);
    if sig = 'LVLN' then begin Result := 'via_lvln'; Exit; end;
    if (sig = 'NPC_') and (GetElementEditValues(cur, 'FULL') <> '') then begin
      Result := 'morph_named'; Exit;
    end;
    if not ElementExists(cur, 'TPLT - Template') then Break;
    cur := LinksTo(ElementByName(cur, 'TPLT - Template'));
    Inc(depth);
  end;
  if HasVal(npc, 'VTCK - Voice') then Result := 'generic_voice'
  else Result := 'editor_only';
end;

// 分類ごとに最大 10 件まで EditorID を控える。役割名の傾向を目視するため。
procedure AddSample(cls, edid: string);
var k: string;
begin
  k := 'samplecount.' + cls;
  if ToIntDef0(counters.Values[k]) < 10 then begin
    samples.Add('  [' + cls + '] ' + edid);
    counters.Values[k] := IntToStr(ToIntDef0(counters.Values[k]) + 1);
  end;
end;

// NPC_ の名称・TPLT 状態を prefix 付きで記録する。
// 形態 record（FULL 空で TPLT が本体を指す）の割合をここで可視化する。
procedure RecordNpcState(prefix: string; npc: IInterface);
var
  tpl: IInterface;
  tsig, cls: string;
begin
  if HasVal(npc, 'FULL') then IncCounter(prefix + '.hasFULL')
  else begin
    IncCounter(prefix + '.noFULL');
    // 名前空 actor のペルソナ名解決を分類する（morph / lvln / voice / editor）。
    cls := ClassifyUnnamedActor(npc);
    IncCounter(prefix + '.unnamed.' + cls);
    AddSample(cls, EditorID(npc));
  end;
  if HasVal(npc, 'SHRT') then IncCounter(prefix + '.hasSHRT');

  if ElementExists(npc, 'TPLT - Template') then begin
    IncCounter(prefix + '.hasTPLT');
    tpl := LinksTo(ElementByName(npc, 'TPLT - Template'));
    tsig := SafeSig(tpl);
    IncCounter(prefix + '.tplt.sig.' + tsig);
    if tsig = 'NPC_' then begin
      if HasVal(tpl, 'FULL') then IncCounter(prefix + '.tplt.npc.hasFULL')
      else IncCounter(prefix + '.tplt.npc.noFULL');
    end;
  end else
    IncCounter(prefix + '.noTPLT');
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

// quest alias index から forced reference / unique actor を解決し、
// REFR は base まで辿って actor record を返す。external / conditions ベースの
// alias は本 helper では解決しない（その分は呼び出し側で unresolved に数える）。
function ResolveAliasActor(quest: IInterface; aliasIdx: Integer): IInterface;
var
  aliases, al, ref: IInterface;
  i: Integer;
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
    else
      ref := nil;
    Result := ResolveActorBase(ref);
    Break;
  end;
end;

// ===== 1 INFO の話者解析 =====

procedure AnalyzeInfo(info: IInterface);
var
  anam, spk, conds, c, ref, quest: IInterface;
  i, aliasIdx: Integer;
  asig, fn, refsig: string;
  anamPresent, anamNpcNamed, anamNpcUnnamed, anamLvln, anamOther: Boolean;
  hasGetIsID, getIsIDNamed: Boolean;
  hasAliasRef, aliasNamed: Boolean;
  hasFaction, hasVoiceType: Boolean;
begin
  Inc(totalInfo);

  anamPresent := False;
  anamNpcNamed := False;
  anamNpcUnnamed := False;
  anamLvln := False;
  anamOther := False;
  hasGetIsID := False;
  getIsIDNamed := False;
  hasAliasRef := False;
  aliasNamed := False;
  hasFaction := False;
  hasVoiceType := False;

  quest := QuestOfInfo(info);

  // --- ANAM-Speaker（特定話者の直接指定）---
  anam := ElementByName(info, 'ANAM - Speaker');
  if not Assigned(anam) then
    IncCounter('anam.absent')
  else begin
    spk := LinksTo(anam);
    if not Assigned(spk) then
      IncCounter('anam.present_but_unresolved')
    else begin
      anamPresent := True;
      asig := SafeSig(spk);
      IncCounter('anam.sig.' + asig);
      if asig = 'NPC_' then begin
        RecordNpcState('anam.npc', spk);
        if HasVal(spk, 'FULL') then anamNpcNamed := True
        else anamNpcUnnamed := True;
      end
      else if asig = 'LVLN' then anamLvln := True
      else anamOther := True;
    end;
  end;

  // --- conditions（CTDA）経由の話者特定 ---
  conds := ElementByName(info, 'Conditions');
  if Assigned(conds) then
    for i := 0 to ElementCount(conds) - 1 do begin
      c := ElementByIndex(conds, i);
      fn := GetElementEditValues(c, 'CTDA\Function');
      IncCounter('cond.func.' + fn);

      if fn = 'GetIsID' then begin
        hasGetIsID := True;
        ref := LinksTo(ElementByPath(c, 'CTDA\Base Object'));
        refsig := SafeSig(ref);
        IncCounter('cond.getisid.sig.' + refsig);
        if refsig = 'NPC_' then begin
          RecordNpcState('cond.getisid.npc', ref);
          if HasVal(ref, 'FULL') then getIsIDNamed := True;
        end;
      end
      else if fn = 'GetIsAliasRef' then begin
        hasAliasRef := True;
        aliasIdx := GetElementNativeValues(c, 'CTDA\Alias');
        ref := ResolveAliasActor(quest, aliasIdx);
        refsig := SafeSig(ref);
        IncCounter('cond.aliasref.sig.' + refsig);
        if refsig = 'NPC_' then begin
          RecordNpcState('cond.aliasref.npc', ref);
          if HasVal(ref, 'FULL') then aliasNamed := True;
        end
        else if not Assigned(ref) then
          IncCounter('cond.aliasref.unresolved');
      end
      else if fn = 'GetInFaction' then begin
        hasFaction := True;
        IncCounter('cond.getinfaction.sig.' + SafeSig(LinksTo(ElementByPath(c, 'CTDA\Faction'))));
      end
      else if fn = 'GetIsVoiceType' then begin
        hasVoiceType := True;
        IncCounter('cond.getisvoicetype.sig.' + SafeSig(LinksTo(ElementByPath(c, 'CTDA\Voice Type'))));
      end;
    end;

  // --- INFO 単位の話者決定パターン分類（1 INFO を 1 パターンへ）---
  // 優先順位は「より特定的な話者指定」を上位に置く。
  if anamPresent then begin
    if anamNpcNamed then IncCounter('pattern.anam_named_npc')
    else if anamNpcUnnamed then IncCounter('pattern.anam_unnamed_npc')
    else if anamLvln then IncCounter('pattern.anam_lvln')
    else if anamOther then IncCounter('pattern.anam_other')
    else IncCounter('pattern.anam_unknown');
  end else begin
    if getIsIDNamed then IncCounter('pattern.cond_getisid_named')
    else if hasGetIsID then IncCounter('pattern.cond_getisid_unnamed_or_nonnpc')
    else if aliasNamed then IncCounter('pattern.cond_aliasref_named')
    else if hasAliasRef then IncCounter('pattern.cond_aliasref_other')
    else if hasFaction then IncCounter('pattern.cond_faction')
    else if hasVoiceType then IncCounter('pattern.cond_voicetype')
    else IncCounter('pattern.voice_or_generic');
  end;
end;

// ===== xEdit エントリ =====

function Initialize: Integer;
begin
  counters := TStringList.Create;
  // Values[key] への代入は Sorted=True だと「Operation not allowed on sorted list」に
  // なるため非ソートで持つ。表示順は Finalize で Sort して prefix ごとにまとめる。
  counters.Sorted := False;
  samples := TStringList.Create;
  samples.Sorted := False;
  totalInfo := 0;
  AddMessage('[AnalyzeSpeaker] Initialized.');
  Result := 0;
end;

function Process(e: IInterface): Integer;
begin
  if Signature(e) = 'INFO' then
    AnalyzeInfo(e);
  Result := 0;
end;

function Finalize: Integer;
var
  i: Integer;
begin
  AddMessage('==== [AnalyzeSpeaker] result ====');
  AddMessage('total INFO = ' + IntToStr(totalInfo));
  AddMessage('---- counters (label = count) ----');
  counters.Sort;   // prefix ごとに並ぶよう表示直前に整列する。
  for i := 0 to counters.Count - 1 do
    AddMessage('  ' + counters[i]);
  AddMessage('---- unnamed-actor EditorID samples ----');
  samples.Sort;
  for i := 0 to samples.Count - 1 do
    AddMessage(samples[i]);
  AddMessage('==== [AnalyzeSpeaker] end ====');
  counters.Free;
  samples.Free;
  Result := 0;
end;

end.
