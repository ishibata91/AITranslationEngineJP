using Microsoft.Data.Sqlite;
using System.Xml;

namespace Extractor;

// xTranslator 英日辞書 XML（SSTXMLRessources）から固有名（FULL）の「原語 → 確定訳語」を読み、
// 中心 DB の master_term へ書く。叙述文・台詞の本文へ固有名を機械置換するための辞書。
// 公式日本語版の既訳流用（既訳流用主軸）の供給源として、既にある英日対訳 XML をそのまま使う。
public static class MasterTermXmlWriter
{
    // master_term へ入れる固有名 1 件。
    private readonly record struct Term(string Source, string Dest, string Category);

    // 与えた XML 群から固有名の対を読み、master_term へ UPSERT する。書いた行数を返す。
    public static int Write(string dbPath, string schemaSql, IEnumerable<string> xmlPaths)
    {
        using var conn = new SqliteConnection($"Data Source={dbPath}");
        conn.Open();

        using (var ddl = conn.CreateCommand())
        {
            ddl.CommandText = schemaSql;
            ddl.ExecuteNonQuery();
        }

        using var tx = conn.BeginTransaction();
        using var cmd = conn.CreateCommand();
        cmd.Transaction = tx;
        cmd.CommandText =
            """
            INSERT OR IGNORE INTO master_term (source, dest, category)
            VALUES ($source, $dest, $category)
            """;
        var source = cmd.Parameters.Add("$source", SqliteType.Text);
        var dest = cmd.Parameters.Add("$dest", SqliteType.Text);
        var category = cmd.Parameters.Add("$category", SqliteType.Text);

        var written = 0;
        foreach (var xmlPath in xmlPaths)
        {
            foreach (var term in ReadTerms(xmlPath))
            {
                source.Value = term.Source;
                dest.Value = term.Dest;
                category.Value = term.Category;
                written += cmd.ExecuteNonQuery();
            }
        }
        tx.Commit();
        return written;
    }

    // 1 つの XML から固有名（FULL 系）の対を読む。REC が ":FULL" で終わる行を固有名とし、
    // 龍語綴り（WOOP:FULL、翻訳禁止）は除く。原語・確定訳語が空、または同一の行は捨てる。
    private static IEnumerable<Term> ReadTerms(string xmlPath)
    {
        using var reader = XmlReader.Create(xmlPath, new XmlReaderSettings { DtdProcessing = DtdProcessing.Ignore });

        string rec = "", src = "", dst = "";
        var inString = false;

        while (reader.Read())
        {
            if (reader.NodeType == XmlNodeType.Element && reader.Name == "String")
            {
                inString = true;
                rec = src = dst = "";
            }
            else if (inString && reader.NodeType == XmlNodeType.Element)
            {
                switch (reader.Name)
                {
                    case "REC": rec = reader.ReadElementContentAsString().Trim().ToUpperInvariant(); break;
                    case "Source": src = reader.ReadElementContentAsString(); break;
                    case "Dest": dst = reader.ReadElementContentAsString(); break;
                }
            }
            else if (inString && reader.NodeType == XmlNodeType.EndElement && reader.Name == "String")
            {
                inString = false;
                if (!IsProperNounFull(rec)) continue;
                var s = src.Trim();
                var d = dst.Trim();
                if (s.Length == 0 || d.Length == 0 || s == d) continue;
                yield return new Term(s, d, rec.Split(':', 2)[0]);
            }
        }
    }

    // 固有名の FULL かどうか。REC:FIELD が ":FULL" で終わり、龍語綴り（翻訳禁止）でない行を固有名とする。
    private static bool IsProperNounFull(string recField)
    {
        if (!recField.EndsWith(":FULL", StringComparison.Ordinal)) return false;
        return recField != "WOOP:FULL"; // 力の言葉の龍語綴りは翻訳禁止のため辞書にしない。
    }
}
