{ ========================================================================
   INFO Elements Dump
   INFO record の top-level element 名と Responses 配列の要素数を出力する調査用 script。
   ExtractInfo が複数 NAM1（response）を展開できているか、Responses の element 名が
   何かを確認する。出力は xEdit の message log のみ。
   ======================================================================== }

unit DumpInfoElements;

var
  shown: Integer;

function Initialize: Integer;
begin
  shown := 0;
  AddMessage('==== [DumpInfoElements] start ====');
  Result := 0;
end;

function Process(e: IInterface): Integer;
var
  i: Integer;
  resp: IInterface;
begin
  // 複数 NAM1 を持ちそうな INFO を 5 件まで観察する。
  if (Signature(e) = 'INFO') and (shown < 5) then begin
    Inc(shown);
    AddMessage('INFO ' + IntToHex(FixedFormID(e), 8) + ' top elements:');
    for i := 0 to ElementCount(e) - 1 do
      AddMessage('  [' + Name(ElementByIndex(e, i)) + ']');
    if ElementExists(e, 'Responses') then begin
      resp := ElementByName(e, 'Responses');
      AddMessage('  -> ElementByName(Responses) count = ' + IntToStr(ElementCount(resp)));
    end else
      AddMessage('  -> no element named "Responses"');
  end;
  Result := 0;
end;

function Finalize: Integer;
begin
  AddMessage('==== [DumpInfoElements] end ====');
  Result := 0;
end;

end.
