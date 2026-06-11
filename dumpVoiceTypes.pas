{ ========================================================================
   Voice Type Dump
   VTYP（Voice Type）record の FormID / EditorID / FULL を一覧出力する調査用 script。
   voice_types（Response の声型一覧）が指す VTYP の実体（FormID）を確認する。
   出力は xEdit の message log のみ。
   ======================================================================== }

unit DumpVoiceTypes;

var
  count: Integer;

function Initialize: Integer;
begin
  count := 0;
  AddMessage('==== [DumpVoiceTypes] FormID  EditorID  FULL ====');
  Result := 0;
end;

function Process(e: IInterface): Integer;
begin
  if Signature(e) = 'VTYP' then begin
    Inc(count);
    AddMessage(IntToHex(FixedFormID(e), 8) + '  ' +
               EditorID(e) + '  FULL=[' + GetElementEditValues(e, 'FULL') + ']');
  end;
  Result := 0;
end;

function Finalize: Integer;
begin
  AddMessage('==== [DumpVoiceTypes] total VTYP = ' + IntToStr(count) + ' ====');
  Result := 0;
end;

end.
