{ ========================================================================
   抽出取りこぼし 原因究明 dump
   3 つの調査を 1 回の Apply で行う。出力は xEdit message log のみ。

   1) INFO:NAM1 取りこぼし（原因4）
      対象 plugin の INFO について、Responses 配下の NAM1 非空件数を生で合計する。
      これと JSON 出力の INFO:NAM1 件数を比較し、取りこぼしが抽出段か整列段かを切る。

   2) field 空（原因2）
      DESC / DNAM / CNAM が JSON で空だった特定 record（TARGET_EDIDS）について、
      取得経路ごとに値が取れるかを並べる。
        - GetElementEditValues(e, FIELD)        現状の取得経路
        - GetEditValue(ElementByName(e, FIELD))
        - GetElementEditValues(WinningOverride(e), FIELD)  master+override マージ後
      どの経路で値が取れるかで pas の取得方法を直す。

   3) WOOP（叫びの言葉）
      対象 plugin の SHOU の top-level element 名と "Words" 配列の有無 / 件数を出す。
      Words が取れない理由（要素名違い / 件数 0 / WNAM 解決不可）を特定する。
   ======================================================================== }

unit DumpDiagnostics;

var
  targetFile: string;
  infoCount, nam1Count, multiInfo: integer;
  shoutShown: integer;
  probeShown: integer;

function B2S(b: boolean): string;
begin
  if b then Result := 'true' else Result := 'false';
end;

// sig ごとの「説明文」field 名を返す。
function DescFieldOf(sig: string): string;
begin
  if sig = 'MGEF' then Result := 'DNAM'
  else if sig = 'BOOK' then Result := 'CNAM'
  else Result := 'DESC';
end;

// field 空調査の対象 EditorID（小文字で照合）。
function InTargetEdids(edid: string): boolean;
var
  e: string;
begin
  e := LowerCase(edid);
  Result :=
    (e = 'dlc1ld_aetherialstaff') or
    (e = 'dlc1dawnguardruneaxe') or
    (e = 'dlc1armordawnguardcuirasslight3') or
    (e = 'vampiredrain04') or
    (e = 'dlc1vampiresleeprested') or
    (e = 'dlc1vampirechangeeffect') or
    (e = 'dlc1vqsaintpage01') or
    (e = 'dlc1dragondrainvitalityshout06');
end;

function SourceOf(e: IInterface): string;
begin
  if Assigned(GetFile(e)) then Result := GetFileName(GetFile(e))
  else Result := '';
end;

function EdidOf(e: IInterface): string;
begin
  Result := GetElementEditValues(MasterOrSelf(e), 'EDID');
end;

function Initialize: integer;
begin
  targetFile := '';
  infoCount := 0;
  nam1Count := 0;
  multiInfo := 0;
  shoutShown := 0;
  probeShown := 0;
  AddMessage('==== [DumpDiagnostics] start ====');
  Result := 0;
end;

procedure ProbeField(e: IInterface);
var
  sig, fld, edid: string;
  ele, wo: IInterface;
begin
  sig := Signature(e);
  fld := DescFieldOf(sig);
  edid := EdidOf(e);
  AddMessage('FIELD-PROBE edid=' + edid + ' sig=' + sig + ' field=' + fld +
             ' source=' + SourceOf(e) + ' formid=' + IntToHex(FixedFormID(e), 8));
  AddMessage('  ElementExists(e,fld) = ' + B2S(ElementExists(e, fld)));
  AddMessage('  GetElementEditValues(e,fld) = [' + GetElementEditValues(e, fld) + ']');
  ele := ElementByName(e, fld);
  AddMessage('  ElementByName assigned = ' + B2S(Assigned(ele)));
  if Assigned(ele) then
    AddMessage('  GetEditValue(ele) = [' + GetEditValue(ele) + ']');
  wo := WinningOverride(e);
  if Assigned(wo) then
    AddMessage('  WinningOverride GetElementEditValues = [' +
               GetElementEditValues(wo, fld) + ']');
end;

procedure DumpShoutWords(shou: IInterface);
var
  i, k: integer;
  words, wordEntry, woop: IInterface;
begin
  AddMessage('SHOU edid=' + EdidOf(shou) + ' source=' + SourceOf(shou));
  words := ElementByName(shou, 'Words of Power');
  if not Assigned(words) then begin
    AddMessage('  Words of Power 不取得');
    Exit;
  end;
  AddMessage('  Words count = ' + IntToStr(ElementCount(words)));
  for i := 0 to ElementCount(words) - 1 do begin
    wordEntry := ElementByIndex(words, i);
    // word entry の子要素名と edit 値を列挙し、WOOP への link field を特定する。
    AddMessage('  word[' + IntToStr(i) + '] child elements:');
    for k := 0 to ElementCount(wordEntry) - 1 do
      AddMessage('     [' + Name(ElementByIndex(wordEntry, k)) + '] = [' +
                 GetEditValue(ElementByIndex(wordEntry, k)) + ']');
    // 候補 link を試す。
    woop := LinksTo(ElementByName(wordEntry, 'WNAM - Word'));
    if not Assigned(woop) then woop := LinksTo(ElementByName(wordEntry, 'WNAM'));
    if not Assigned(woop) then woop := LinksTo(wordEntry);
    if Assigned(woop) then
      AddMessage('     -> WOOP=' + EdidOf(woop) +
                 ' FULL=[' + GetElementEditValues(woop, 'FULL') + ']' +
                 ' TNAM=[' + GetElementEditValues(woop, 'TNAM') + ']')
    else
      AddMessage('     -> WOOP link 解決不可');
  end;
end;

function Process(e: IInterface): integer;
var
  sig, source: string;
  resp, respItem: IInterface;
  j: integer;
  cnt: integer;
begin
  source := SourceOf(e);
  if targetFile = '' then targetFile := source;
  sig := Signature(e);

  // 1) INFO:NAM1 生件数
  if (sig = 'INFO') and (source = targetFile) then begin
    Inc(infoCount);
    resp := ElementByName(e, 'Responses');
    if Assigned(resp) then begin
      cnt := 0;
      for j := 0 to ElementCount(resp) - 1 do begin
        respItem := ElementByIndex(resp, j);
        if GetElementValue(respItem, 'NAM1') <> '' then begin
          Inc(nam1Count);
          Inc(cnt);
        end;
      end;
      if cnt > 1 then Inc(multiInfo);
    end;
  end;

  // 2) field 空調査
  if InTargetEdids(EdidOf(e)) and (probeShown < 40) then begin
    Inc(probeShown);
    ProbeField(e);
  end;

  // 3) SHOU Words
  if (sig = 'SHOU') and (source = targetFile) and (shoutShown < 3) then begin
    Inc(shoutShown);
    DumpShoutWords(e);
  end;

  Result := 0;
end;

function Finalize: integer;
begin
  AddMessage('---- INFO:NAM1 集計 ----');
  AddMessage('  targetFile = ' + targetFile);
  AddMessage('  INFO record 数 (source=target) = ' + IntToStr(infoCount));
  AddMessage('  NAM1 非空 生件数 = ' + IntToStr(nam1Count));
  AddMessage('  複数 NAM1 を持つ INFO 数 = ' + IntToStr(multiInfo));
  AddMessage('==== [DumpDiagnostics] end ====');
  Result := 0;
end;

end.
