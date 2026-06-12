using Mutagen.Bethesda;
using Mutagen.Bethesda.Plugins;
using Mutagen.Bethesda.Plugins.Records;
using Mutagen.Bethesda.Skyrim;
using Mutagen.Bethesda.Strings;

namespace Extractor;

// plugin 単位列挙で翻訳対象を抽出する（docs/mutagen-migration-plan.md §4）。
// - 対象 plugin が書いた record（新規 + override）だけを列挙する。WinningOverrides は使わない。
// - 参照は FormKey で保持する。参照先が純 master の場合は ID のみ残る（dangling 許容）。
// - 話者解決（r2a〜d）は INFO の ANAM と Conditions（CTDA）から導出する。
//
// 翻訳所有（record 単位の gate。Dawnguard と xTranslator 辞書の突き合わせで確定した規則）:
//   次のいずれかを満たす record だけが「この plugin の翻訳対象」になり、その全 string を数える。
//     1. 新規 record（FormKey の origin が target 自身）
//     2. 翻訳対象 field の解決済み text 集合が master（target を除く load order の winning 版）と異なる
//     3. DIAL のみ: plugin 内に子 INFO を持つ（topic は会話の container として扱われる）
//   どれも満たさない override は参照用 stub であり、xTranslator も record 走査で辞書に
//   紐づけない（該当 string は **** 行 = 非対象になる）。string ID の振り直しや record data の
//   差分（script 並び等）は判定に関係しない。text だけを見る。
public static class PluginExtractor
{
    public static ExtractionResult Extract(PluginEnvironment env)
    {
        var mod = env.TargetMod;
        var localized = mod.ModHeader.Flags.HasFlag(SkyrimModHeader.HeaderFlag.Localized);
        var result = new ExtractionResult
        {
            TargetPlugin = mod.ModKey.FileName,
            // localized plugin の header は strings table を経由しないため翻訳対象にならない。
            HeaderAuthor = localized ? "" : mod.ModHeader.Author ?? "",
            HeaderDescription = localized ? "" : mod.ModHeader.Description ?? "",
        };

        // 参照所有: 所有された会話（topic / INFO）が参照する quest は、record 不変でも
        // この plugin の翻訳対象になる（xTranslator の実測挙動。C03 / VC01 等で確認）。
        var ownedQuests = ExtractDialogues(env, result);
        ExtractQuests(env, result, ownedQuests);
        ExtractSimpleCategories(env, result);
        ExtractSpeakers(env, result);

        return result;
    }

    // localized string を解決する。strings file に entry が無い場合は空として扱う。
    private static string S(ITranslatedStringGetter? t) => t?.String ?? "";

    private static string Edid(IMajorRecordGetter rec) => rec.EditorID ?? "";

    private static bool Any(params string[] texts) => texts.Any(t => t.Length > 0);

    // ===== Dialogue（DIAL → InfoNode → ResponseLine の 2 階層）=====

    private static HashSet<FormKey> ExtractDialogues(PluginEnvironment env, ExtractionResult result)
    {
        var ownedQuests = new HashSet<FormKey>();
        foreach (var topic in env.TargetMod.DialogTopics)
        {
            var infos = topic.Responses
                .Where(info => env.OwnsRecord(info))
                .Select(info => ExtractInfoNode(env, topic, info))
                .ToList();

            // 規則 3: 子 INFO を持つ topic は FULL を text 不変でも対象にする。
            var ownsName = infos.Count > 0
                ? S(topic.Name).Length > 0
                : env.OwnsRecord(topic) && S(topic.Name).Length > 0;

            if (!ownsName && infos.Count == 0) continue; // 参照用 stub（翻訳対象なし）

            if (!topic.Quest.IsNull) ownedQuests.Add(topic.Quest.FormKey);
            result.Dialogues.Add(new DialogueTopic
            {
                Id = topic.FormKey,
                EditorId = Edid(topic),
                Kind = "DIAL",
                Name = S(topic.Name),
                Category = topic.Category.ToString(),
                Subtype = topic.Subtype.ToString(),
                QuestId = topic.Quest.IsNull ? null : topic.Quest.FormKey,
                IdenticalToMaster = !ownsName,
                Infos = infos,
            });
        }
        return ownedQuests;
    }

    private static InfoNode ExtractInfoNode(PluginEnvironment env, IDialogTopicGetter topic, IDialogResponsesGetter info)
    {
        var node = new InfoNode
        {
            Id = info.FormKey,
            EditorId = Edid(info),
            Prompt = S(info.Prompt),
            Responses = info.Responses
                .Select(r => new ResponseLine(r.ResponseNumber, S(r.Text)))
                .Where(r => r.Text.Length > 0)
                .ToList(),
            PreviousId = info.PreviousDialog.IsNull ? null : info.PreviousDialog.FormKey,
            NextIds = info.LinkTo.Where(l => !l.IsNull).Select(l => l.FormKey).ToList(),
            Conditions = info.Conditions.Select(BuildCondition).ToList(),
        };

        // r2a: 名指し話者。ANAM、GetIsID（NPC_/TACT）、GetIsAliasRef（quest alias の ALFR/ALUA）。
        if (!info.Speaker.IsNull)
            node.SpeakerIds.Add(info.Speaker.FormKey);

        foreach (var cond in info.Conditions)
        {
            switch (cond.Data)
            {
                case IGetIsIDConditionDataGetter d:
                {
                    var key = d.Object.Link.FormKey;
                    if (env.LinkCache.TryResolve<INpcGetter>(key, out _) ||
                        env.LinkCache.TryResolve<ITalkingActivatorGetter>(key, out _))
                        node.SpeakerIds.Add(key);
                    break;
                }
                case IGetIsAliasRefConditionDataGetter d:
                {
                    var actor = ResolveAliasActor(env, topic, d.ReferenceAliasIndex);
                    if (actor != null) node.SpeakerIds.Add(actor.Value);
                    break;
                }
                case IGetInFactionConditionDataGetter d:
                    node.FactionIds.Add(d.Faction.Link.FormKey);
                    break;
                case IGetIsRaceConditionDataGetter d:
                    node.RaceIds.Add(d.Race.Link.FormKey);
                    break;
                case IGetIsVoiceTypeConditionDataGetter d:
                {
                    // VTYP 直指定だけ採用する（FLST 指定はプール扱いで対象外、Pascal 版と同じ）。
                    var key = d.VoiceTypeOrList.Link.FormKey;
                    if (env.LinkCache.TryResolve<IVoiceTypeGetter>(key, out _))
                        node.VoiceTypeIds.Add(key);
                    break;
                }
            }
        }

        return node;
    }

    // quest alias index から ALFR（forced reference → base NPC）/ ALUA（unique actor）を解決する。
    private static FormKey? ResolveAliasActor(PluginEnvironment env, IDialogTopicGetter topic, int aliasIndex)
    {
        if (topic.Quest.IsNull) return null;
        if (!env.LinkCache.TryResolve<IQuestGetter>(topic.Quest.FormKey, out var quest)) return null;

        var alias = quest.Aliases.FirstOrDefault(a => a.ID == aliasIndex);
        if (alias == null) return null;

        if (!alias.ForcedReference.IsNull &&
            env.LinkCache.TryResolve<IPlacedNpcGetter>(alias.ForcedReference.FormKey, out var placed) &&
            !placed.Base.IsNull)
            return placed.Base.FormKey;

        if (!alias.UniqueActor.IsNull)
            return alias.UniqueActor.FormKey;

        return null;
    }

    private static ConditionEntry BuildCondition(IConditionGetter cond)
    {
        var op = cond.CompareOperator switch
        {
            CompareOperator.EqualTo => "==",
            CompareOperator.NotEqualTo => "!=",
            CompareOperator.GreaterThan => ">",
            CompareOperator.GreaterThanOrEqualTo => ">=",
            CompareOperator.LessThan => "<",
            CompareOperator.LessThanOrEqualTo => "<=",
            _ => "==",
        };
        var value = cond switch
        {
            IConditionFloatGetter f => f.ComparisonValue.ToString("0.######"),
            IConditionGlobalGetter g => g.ComparisonValue.FormKey.ToString(),
            _ => "",
        };
        FormKey? refId = cond.Data switch
        {
            IGetIsIDConditionDataGetter d => d.Object.Link.FormKey,
            IGetInFactionConditionDataGetter d => d.Faction.Link.FormKey,
            IGetIsRaceConditionDataGetter d => d.Race.Link.FormKey,
            IGetIsClassConditionDataGetter d => d.Class.Link.FormKey,
            IGetIsVoiceTypeConditionDataGetter d => d.VoiceTypeOrList.Link.FormKey,
            _ => null,
        };
        return new ConditionEntry(cond.Data.Function.ToString(), op, value, refId);
    }

    // ===== Quest =====

    private static void ExtractQuests(PluginEnvironment env, ExtractionResult result, HashSet<FormKey> ownedQuests)
    {
        foreach (var quest in env.TargetMod.Quests)
        {
            if (!env.OwnsRecord(quest)) continue;

            var stages = quest.Stages
                .Select(st => new QuestStageEntry(
                    st.Index,
                    st.LogEntries.Select(le => S(le.Entry)).Where(t => t.Length > 0).ToList()))
                .Where(st => st.LogEntries.Count > 0)
                .ToList();
            var objectives = quest.Objectives
                .Select(ob => new QuestObjectiveEntry(ob.Index, S(ob.DisplayText)))
                .Where(ob => ob.DisplayText.Length > 0)
                .ToList();
            var name = S(quest.Name);
            if (!Any(name) && stages.Count == 0 && objectives.Count == 0) continue;

            result.Quests.Add(new QuestEntry
            {
                Id = quest.FormKey,
                EditorId = Edid(quest),
                Kind = "QUST",
                Name = name,
                QuestType = quest.Type.ToString(),
                Stages = stages,
                Objectives = objectives,
            });
        }

    }

    // ===== 単純カテゴリ（FULL / DESC 系）=====

    private static void ExtractSimpleCategories(PluginEnvironment env, ExtractionResult result)
    {
        var mod = env.TargetMod;

        foreach (var race in mod.Races)
        {
            if (!env.OwnsRecord(race)) continue;
            var name = S(race.Name);
            var desc = S(race.Description);
            if (!Any(name, desc)) continue;
            result.Races.Add(new RaceEntry
            {
                Id = race.FormKey, EditorId = Edid(race), Kind = "RACE",
                Name = name, Description = desc,
            });
        }

        foreach (var fact in mod.Factions)
        {
            if (!env.OwnsRecord(fact)) continue;
            var ranks = fact.Ranks
                .Select(r => new RankTitle(S(r.Title?.Male), S(r.Title?.Female)))
                .Where(r => Any(r.Male, r.Female))
                .ToList();
            // 新規 FACT は名前空でも emit する（Pascal 版と同じ。翻訳 context に効く可能性があるため）。
            result.Factions.Add(new FactionEntry
            {
                Id = fact.FormKey, EditorId = Edid(fact), Kind = "FACT",
                Name = S(fact.Name), Ranks = ranks,
            });
        }


        // Item 系（FULL のみ）: KEYM/MISC/LIGH/CONT/SLGM/DOOR/FURN/APPA
        AddNamed(result.Items, env, mod.Keys, "KEYM", r => r.Name);
        AddNamed(result.Items, env, mod.MiscItems, "MISC", r => r.Name);
        AddNamed(result.Items, env, mod.Lights, "LIGH", r => r.Name);
        AddNamed(result.Items, env, mod.Containers, "CONT", r => r.Name);
        AddNamed(result.Items, env, mod.SoulGems, "SLGM", r => r.Name);
        AddNamed(result.Items, env, mod.Doors, "DOOR", r => r.Name);
        AddNamed(result.Items, env, mod.Furniture, "FURN", r => r.Name);
        AddNamed(result.Items, env, mod.AlchemicalApparatuses, "APPA", r => r.Name);
        AddNamed(result.Items, env, mod.Hazards, "HAZD", r => r.Name);

        // Activator 系（FULL + RNAM）: ACTI/FLOR/TREE
        foreach (var a in mod.Activators)
        {
            if (!env.OwnsRecord(a)) continue;
            AddActivator(result, a.FormKey, Edid(a), "ACTI", S(a.Name), S(a.ActivateTextOverride));
        }
        foreach (var a in mod.Florae)
        {
            if (!env.OwnsRecord(a)) continue;
            AddActivator(result, a.FormKey, Edid(a), "FLOR", S(a.Name), S(a.ActivateTextOverride));
        }
        foreach (var a in mod.Trees)
        {
            if (!env.OwnsRecord(a)) continue;
            AddActivator(result, a.FormKey, Edid(a), "TREE", S(a.Name), "");
        }

        // Equipment（FULL + DESC + EITM）: WEAP/ARMO/AMMO
        foreach (var w in mod.Weapons)
        {
            if (!env.OwnsRecord(w)) continue;
            AddDescribed(result.Equipment, w.FormKey, Edid(w), "WEAP", S(w.Name), S(w.Description),
                enchantmentId: w.ObjectEffect.IsNull ? null : w.ObjectEffect.FormKey);
        }
        foreach (var a in mod.Armors)
        {
            if (!env.OwnsRecord(a)) continue;
            AddDescribed(result.Equipment, a.FormKey, Edid(a), "ARMO", S(a.Name), S(a.Description),
                enchantmentId: a.ObjectEffect.IsNull ? null : a.ObjectEffect.FormKey);
        }
        foreach (var a in mod.Ammunitions)
        {
            if (!env.OwnsRecord(a)) continue;
            AddDescribed(result.Equipment, a.FormKey, Edid(a), "AMMO", S(a.Name), S(a.Description));
        }

        // 参照所有: 所有された SPEL / ENCH / 消耗品の Effects が参照する MGEF は
        // record 不変でも対象になる（xTranslator の実測挙動。AbFXDwarvenSpider で確認）。
        var ownedMgefs = new HashSet<FormKey>();

        // Consumable（FULL + DESC + Effects）: ALCH/SCRL/INGR
        foreach (var c in mod.Ingestibles)
        {
            if (!env.OwnsRecord(c)) continue;
            var effects = Effects(c.Effects);
            ownedMgefs.UnionWith(effects);
            AddDescribed(result.Consumables, c.FormKey, Edid(c), "ALCH", S(c.Name), S(c.Description), effects);
        }
        foreach (var c in mod.Scrolls)
        {
            if (!env.OwnsRecord(c)) continue;
            var effects = Effects(c.Effects);
            ownedMgefs.UnionWith(effects);
            AddDescribed(result.Consumables, c.FormKey, Edid(c), "SCRL", S(c.Name), S(c.Description), effects);
        }
        foreach (var c in mod.Ingredients)
        {
            if (!env.OwnsRecord(c)) continue;
            var effects = Effects(c.Effects);
            ownedMgefs.UnionWith(effects);
            AddDescribed(result.Consumables, c.FormKey, Edid(c), "INGR", S(c.Name), "", effects);
        }

        // Magic / Enchantment / MagicEffect
        foreach (var s in mod.Spells)
        {
            if (!env.OwnsRecord(s)) continue;
            var effects = Effects(s.Effects);
            ownedMgefs.UnionWith(effects);
            AddDescribed(result.Magic, s.FormKey, Edid(s), "SPEL", S(s.Name), S(s.Description), effects);
        }
        foreach (var e in mod.ObjectEffects)
        {
            if (!env.OwnsRecord(e)) continue;
            var effects = Effects(e.Effects);
            ownedMgefs.UnionWith(effects);
            // 新規 ENCH は名前空でも emit する（装備品付与で意味を持つ。Pascal 版と同じ）。
            result.Enchantments.Add(new DescribedEntry
            {
                Id = e.FormKey, EditorId = Edid(e), Kind = "ENCH",
                Name = S(e.Name), Description = "", MagicEffectIds = effects,
            });
        }
        foreach (var m in mod.MagicEffects)
        {
            if (!env.OwnsRecord(m)) continue;
            AddDescribed(result.MagicEffects, m.FormKey, Edid(m), "MGEF", S(m.Name), S(m.Description));
        }

        // Shout + Word of Power
        foreach (var sh in mod.Shouts)
        {
            if (!env.OwnsRecord(sh)) continue;
            var name = S(sh.Name);
            var desc = S(sh.Description);
            var wordIds = sh.WordsOfPower
                .Where(w => !w.Word.IsNull)
                .Select(w => w.Word.FormKey)
                .ToList();
            var isNew = sh.FormKey.ModKey == mod.ModKey;
            if (!Any(name, desc) && !isNew) continue;
            result.Shouts.Add(new ShoutEntry
            {
                Id = sh.FormKey, EditorId = Edid(sh), Kind = "SHOU",
                Name = name, Description = desc, WordIds = wordIds,
            });
        }
        foreach (var w in mod.WordsOfPower)
        {
            if (!env.OwnsRecord(w)) continue;
            var spelling = S(w.Name);          // FULL（龍語綴り、翻訳禁止）
            var translation = S(w.Translation); // TNAM（翻訳対象）
            if (!Any(spelling, translation)) continue;
            result.Words.Add(new WordOfPowerEntry
            {
                Id = w.FormKey, EditorId = Edid(w), Kind = "WOOP",
                DragonSpelling = spelling, Translation = translation,
            });
        }

        // Book
        foreach (var b in mod.Books)
        {
            if (!env.OwnsRecord(b)) continue;
            var title = S(b.Name);
            var body = S(b.BookText);        // DESC（本文）
            var author = S(b.Description);   // CNAM（SSEEdit 表示名 Description）
            if (!Any(title, body, author)) continue;
            result.Books.Add(new BookEntry
            {
                Id = b.FormKey, EditorId = Edid(b), Kind = "BOOK",
                Title = title, Body = body, Author = author,
            });
        }

        // Location 系（FULL のみ）: LCTN/WRLD/CELL。CELL は worldspace 配下も含めて深く列挙する。
        AddNamed(result.Locations, env, mod.Locations, "LCTN", r => r.Name);
        foreach (var w in mod.Worldspaces)
        {
            if (!env.OwnsRecord(w)) continue;
            var name = S(w.Name);
            if (name.Length == 0) continue;
            result.Locations.Add(new NamedEntry { Id = w.FormKey, EditorId = Edid(w), Kind = "WRLD", Name = name });
        }
        foreach (var cell in mod.EnumerateMajorRecords<ICellGetter>())
        {
            // CELL は DIAL と同様に container として働く。plugin が所有する配置 ref
            //（新規または変更された REFR/ACHR）を持つ override は record data 不変でも対象になる。
            // navmesh（NAVM/LAND）だけの編集や、子まで master と同一の cell は対象にしない。
            var hasOwnedChildren = cell.Persistent.Concat(cell.Temporary).Any(env.OwnsRecord);
            if (!hasOwnedChildren && !env.OwnsRecord(cell)) continue;
            var name = S(cell.Name);
            if (name.Length == 0) continue;
            result.Locations.Add(new NamedEntry { Id = cell.FormKey, EditorId = Edid(cell), Kind = "CELL", Name = name });
        }

        // Message（FULL + DESC + ITXT）
        foreach (var m in mod.Messages)
        {
            if (!env.OwnsRecord(m)) continue;
            var title = S(m.Name);
            var body = S(m.Description);
            var buttons = m.MenuButtons.Select(b => S(b.Text)).Where(t => t.Length > 0).ToList();
            if (!Any(title, body) && buttons.Count == 0) continue;
            result.Messages.Add(new MessageEntry
            {
                Id = m.FormKey, EditorId = Edid(m), Kind = "MESG",
                Title = title, Body = body,
                QuestId = m.Quest.IsNull ? null : m.Quest.FormKey,
                Buttons = buttons,
            });
        }


        // LoadScreen（DESC のみ）
        foreach (var l in mod.LoadScreens)
        {
            if (!env.OwnsRecord(l)) continue;
            var body = S(l.Description);
            if (body.Length == 0) continue;
            result.LoadingScreens.Add(new BodyEntry { Id = l.FormKey, EditorId = Edid(l), Kind = "LSCR", Body = body });
        }

        // Perk（FULL + DESC。EPF2/EPFD は非対象）
        foreach (var p in mod.Perks)
        {
            if (!env.OwnsRecord(p)) continue;
            AddDescribed(result.Perks, p.FormKey, Edid(p), "PERK", S(p.Name), S(p.Description));
        }

        // VoiceType（翻訳テキスト無し、識別子のみ。新規だけ emit する）
        foreach (var v in mod.VoiceTypes)
        {
            if (v.FormKey.ModKey != mod.ModKey) continue;
            var edid = Edid(v);
            if (edid.Length == 0) continue;
            result.VoiceTypes.Add(new VoiceTypeEntry { Id = v.FormKey, EditorId = edid, Kind = "VTYP", Identifier = edid });
        }

        // Skyrim 追加 record: SNCT / EYES / REGN（docs/mutagen-migration-plan.md §8-4）
        AddNamed(result.SoundCategories, env, mod.SoundCategories, "SNCT", r => r.Name);
        AddNamed(result.Eyes, env, mod.Eyes, "EYES", r => r.Name);
        foreach (var r in mod.Regions)
        {
            if (!env.OwnsRecord(r)) continue;
            var mapName = S(r.Map?.Name);
            if (mapName.Length == 0) continue;
            result.Regions.Add(new RegionEntry { Id = r.FormKey, EditorId = Edid(r), Kind = "REGN", MapName = mapName });
        }
    }

    private static void AddNamed<T>(List<NamedEntry> list, PluginEnvironment env, IEnumerable<T> records, string kind,
        Func<T, ITranslatedStringGetter?> name)
        where T : class, IMajorRecordGetter
    {
        foreach (var rec in records)
        {
            if (!env.OwnsRecord(rec)) continue;
            var resolved = S(name(rec));
            if (resolved.Length == 0) continue;
            list.Add(new NamedEntry { Id = rec.FormKey, EditorId = Edid(rec), Kind = kind, Name = resolved });
        }
    }

    private static void AddActivator(ExtractionResult result, FormKey id, string edid, string kind, string name, string action)
    {
        if (!Any(name, action)) return;
        result.Activators.Add(new ActivatorEntry
        {
            Id = id, EditorId = edid, Kind = kind, Name = name, ActivateAction = action,
        });
    }

    private static void AddDescribed(List<DescribedEntry> list, FormKey id, string edid, string kind,
        string name, string desc, List<FormKey>? effects = null, FormKey? enchantmentId = null)
    {
        if (!Any(name, desc)) return;
        list.Add(new DescribedEntry
        {
            Id = id, EditorId = edid, Kind = kind, Name = name, Description = desc,
            MagicEffectIds = effects ?? [], EnchantmentId = enchantmentId,
        });
    }

    private static List<FormKey> Effects(IReadOnlyList<IEffectGetter> effects) =>
        effects.Where(e => !e.BaseEffect.IsNull).Select(e => e.BaseEffect.FormKey).ToList();

    // ===== Speaker（NPC_ / TACT）=====

    private static void ExtractSpeakers(PluginEnvironment env, ExtractionResult result)
    {
        var mod = env.TargetMod;

        foreach (var npc in mod.Npcs)
        {
            if (!env.OwnsRecord(npc)) continue;
            // 新規 NPC_ は名前空でも emit する（ペルソナ名が必ず解決できる、v18 仕様）。
            // override は所有 text がある場合だけ出す。
            var name = S(npc.Name);
            var shortName = S(npc.ShortName);
            if (npc.FormKey.ModKey != mod.ModKey && !Any(name, shortName)) continue;

            var isFemale = npc.Configuration.Flags.HasFlag(NpcConfiguration.Flag.Female);
            result.Speakers.Add(new SpeakerEntry
            {
                Id = npc.FormKey,
                EditorId = Edid(npc),
                Kind = "NPC_",
                Name = name,
                PersonaName = ResolvePersonaName(env, npc),
                SpeakerKind = ClassifySpeakerKind(env, npc),
                ShortName = shortName,
                Sex = isFemale ? "Female" : "Male",
                VoiceTypeId = npc.Voice.IsNull ? null : npc.Voice.FormKey,
                RaceId = npc.Race.IsNull ? null : npc.Race.FormKey,
                BaseSpeakerId = npc.Template.IsNull ? null : npc.Template.FormKey,
                FactionIds = npc.Factions.Where(f => !f.Faction.IsNull).Select(f => f.Faction.FormKey).ToList(),
            });
        }

        foreach (var tact in mod.TalkingActivators)
        {
            if (!env.OwnsRecord(tact)) continue;
            var name = S(tact.Name);
            if (tact.FormKey.ModKey != mod.ModKey && !Any(name)) continue;
            result.Speakers.Add(new SpeakerEntry
            {
                Id = tact.FormKey,
                EditorId = Edid(tact),
                Kind = "TACT",
                Name = name,
                PersonaName = name.Length > 0 ? name : Edid(tact),
                SpeakerKind = "喋るオブジェクト",
                ShortName = "",
                Sex = "",
                VoiceTypeId = tact.Voice.IsNull ? null : tact.Voice.FormKey,
            });
        }
    }

    // ペルソナ名の解決順: 自身の FULL → TPLT 連鎖の本体 FULL → EditorID（役割名）。
    private static string ResolvePersonaName(PluginEnvironment env, INpcGetter npc)
    {
        INpcGetter? cur = npc;
        for (var depth = 0; cur != null && depth < 12; depth++)
        {
            var full = S(cur.Name);
            if (full.Length > 0) return full;
            cur = ResolveTemplateNpc(env, cur);
        }
        return Edid(npc);
    }

    // Speaker 種別判定（Pascal 版 ClassifySpeakerKind の移植）。
    private static string ClassifySpeakerKind(PluginEnvironment env, INpcGetter npc)
    {
        if (S(npc.Name).Length > 0) return "固有";

        var curKey = npc.Template.IsNull ? (FormKey?)null : npc.Template.FormKey;
        for (var depth = 0; curKey != null && depth < 12; depth++)
        {
            if (env.LinkCache.TryResolve<ILeveledNpcGetter>(curKey.Value, out _)) return "プール";
            if (!env.LinkCache.TryResolve<INpcGetter>(curKey.Value, out var cur)) break;
            if (S(cur.Name).Length > 0) return "形態";
            curKey = cur.Template.IsNull ? null : cur.Template.FormKey;
        }

        return npc.Voice.IsNull ? "役割のみ" : "声型代表";
    }

    private static INpcGetter? ResolveTemplateNpc(PluginEnvironment env, INpcGetter npc)
    {
        if (npc.Template.IsNull) return null;
        return env.LinkCache.TryResolve<INpcGetter>(npc.Template.FormKey, out var resolved) ? resolved : null;
    }
}
