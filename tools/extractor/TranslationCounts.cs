using Mutagen.Bethesda.Plugins;

namespace Extractor;

// 翻訳文字列 1 件（xTranslator の String 行に対応する単位）。
public sealed record TranslationString(string RecField, FormKey Id, string EditorId, string Text);

// 抽出結果を xTranslator の翻訳文字列単位（REC:FIELD）へフラット化する
// （scripts/compare_counts.py の flatten_json 相当。2 階層モデルのため id ユニーク化は不要）。
public static class TranslationCounts
{
    public static SortedDictionary<string, int> Flatten(ExtractionResult r)
    {
        var c = new SortedDictionary<string, int>(StringComparer.Ordinal);
        foreach (var s in Enumerate(r))
            c[s.RecField] = c.GetValueOrDefault(s.RecField) + 1;
        return c;
    }

    // 非空の翻訳文字列を全て列挙する。
    public static IEnumerable<TranslationString> Enumerate(ExtractionResult r)
    {
        var list = new List<TranslationString>();

        void Add(string recField, RecordEntry entry, string text) =>
            AddRaw(recField, entry.Id, entry.EditorId, text);

        void AddRaw(string recField, FormKey id, string editorId, string text)
        {
            if (text.Trim().Length == 0) return;
            list.Add(new TranslationString(recField, id, editorId, text));
        }

        AddRaw("TES4:CNAM", FormKey.Null, "TES4", r.HeaderAuthor);
        AddRaw("TES4:SNAM", FormKey.Null, "TES4", r.HeaderDescription);

        foreach (var dlg in r.Dialogues)
        {
            // ITM の topic は INFO の container として残るだけで、FULL は master 側の翻訳対象。
            if (!dlg.IdenticalToMaster) Add("DIAL:FULL", dlg, dlg.Name);
            foreach (var info in dlg.Infos)
            {
                AddRaw("INFO:RNAM", info.Id, info.EditorId, info.Prompt);
                foreach (var resp in info.Responses)
                    AddRaw("INFO:NAM1", info.Id, info.EditorId, resp.Text);
            }
        }

        foreach (var q in r.Quests)
        {
            Add("QUST:FULL", q, q.Name);
            foreach (var st in q.Stages)
                foreach (var log in st.LogEntries)
                    Add("QUST:CNAM", q, log);
            foreach (var ob in q.Objectives)
                Add("QUST:NNAM", q, ob.DisplayText);
        }

        foreach (var x in r.Races)
        {
            Add("RACE:FULL", x, x.Name);
            Add("RACE:DESC", x, x.Description);
        }

        foreach (var f in r.Factions)
        {
            Add("FACT:FULL", f, f.Name);
            foreach (var rank in f.Ranks)
            {
                Add("FACT:MNAM", f, rank.Male);
                Add("FACT:FNAM", f, rank.Female);
            }
        }

        foreach (var x in r.Items) Add($"{x.Kind}:FULL", x, x.Name);

        foreach (var a in r.Activators)
        {
            Add($"{a.Kind}:FULL", a, a.Name);
            Add($"{a.Kind}:RNAM", a, a.ActivateAction);
        }

        foreach (var x in r.Equipment.Concat(r.Consumables).Concat(r.Magic).Concat(r.Enchantments).Concat(r.Perks))
        {
            Add($"{x.Kind}:FULL", x, x.Name);
            Add($"{x.Kind}:DESC", x, x.Description);
        }

        foreach (var m in r.MagicEffects)
        {
            Add("MGEF:FULL", m, m.Name);
            Add("MGEF:DNAM", m, m.Description);
        }

        foreach (var sh in r.Shouts)
        {
            Add("SHOU:FULL", sh, sh.Name);
            Add("SHOU:DESC", sh, sh.Description);
        }

        // WOOP:FULL（龍語綴り）は翻訳禁止のため数えない。TNAM（意味の訳）のみ対象。
        foreach (var w in r.Words)
            Add("WOOP:TNAM", w, w.Translation);

        foreach (var b in r.Books)
        {
            Add("BOOK:FULL", b, b.Title);
            Add("BOOK:DESC", b, b.Body);
            Add("BOOK:CNAM", b, b.Author);
        }

        foreach (var x in r.Locations) Add($"{x.Kind}:FULL", x, x.Name);

        foreach (var m in r.Messages)
        {
            Add("MESG:FULL", m, m.Title);
            Add("MESG:DESC", m, m.Body);
            foreach (var b in m.Buttons)
                Add("MESG:ITXT", m, b);
        }

        foreach (var l in r.LoadingScreens) Add("LSCR:DESC", l, l.Body);

        foreach (var s in r.Speakers)
        {
            Add($"{s.Kind}:FULL", s, s.Name);
            if (s.Kind == "NPC_") Add("NPC_:SHRT", s, s.ShortName);
        }

        foreach (var x in r.SoundCategories) Add("SNCT:FULL", x, x.Name);
        foreach (var x in r.Eyes) Add("EYES:FULL", x, x.Name);
        foreach (var x in r.Regions) Add("REGN:RDMP", x, x.MapName);

        // VTYP は翻訳テキストを持たないため数えない。
        return list;
    }
}
