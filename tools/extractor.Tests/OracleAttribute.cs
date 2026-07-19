namespace Extractor.Tests;

// テスト関数へ対応する spec id を紐づける（オラクルテストの規約: オラクル ID を引ける）。
// 網羅番人はこの属性の id 集合と specs.json（自段・非委任）の一致を見る。
[AttributeUsage(AttributeTargets.Method, AllowMultiple = false)]
public sealed class OracleAttribute : Attribute
{
    public string Id { get; }
    public OracleAttribute(string id) => Id = id;
}
