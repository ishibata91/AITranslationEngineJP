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
// 翻訳所有（record 単位。Dawnguard 等と xTranslator 辞書の突き合わせで確定した規則）:
//   次のいずれかを満たす record が「この plugin の翻訳対象」になり、その全 string を数える。
//     1. 新規 record（FormKey の origin が target 自身）
//     2. 正規化 record data が、存在するどれかの master 版と異なる（PluginEnvironment.OwnsRecord）
//     3. container 規則: DIAL は所有 INFO を、CELL は所有された子（配置 ref / navmesh / 地形）を持つ
//   どれも満たさない override は参照用 stub。record としてはモデルに出すが
//   IdenticalToMaster を立て、翻訳対象には数えない（xTranslator の **** 行 = 非対象に対応する）。
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

        ExtractDialogues(env, result);
        ExtractQuests(env, result);
        ExtractSimpleCategories(env, result);
        ExtractSpeakers(env, result);

        return result;
    }

    // localized string を解決する。strings file に entry が無い場合は空として扱う。
    private static string S(ITranslatedStringGetter? t) => t?.String ?? "";

    private static string Edid(IMajorRecordGetter rec) => rec.EditorID ?? "";

    private static bool Any(params string[] texts) => texts.Any(t => t.Length > 0);

    // ===== Dialogue（DIAL → InfoNode → ResponseLine の 2 階層）=====

    private static void ExtractDialogues(PluginEnvironment env, ExtractionResult result)
    {
        foreach (var topic in env.TargetMod.DialogTopics)
        {
            var infos = topic.Responses
                .Select(info => ExtractInfoNode(env, topic, info))
                .ToList();

            var name = S(topic.Name);
            if (name.Length == 0 && infos.Count == 0) continue; // 出すものが無い

            // container 規則: 所有 INFO を持つ topic は FULL を text 不変でも所有する。
            var ownsName = infos.Any(i => !i.IdenticalToMaster) || env.OwnsRecord(topic);
            result.Dialogues.Add(new DialogueTopic
            {
                Id = topic.FormKey,
                EditorId = Edid(topic),
                Kind = "DIAL",
                Name = name,
                Category = topic.Category.ToString(),
                Subtype = topic.Subtype.ToString(),
                QuestId = topic.Quest.IsNull ? null : topic.Quest.FormKey,
                IdenticalToMaster = !ownsName,
                Infos = infos,
            });
        }
    }

    private static InfoNode ExtractInfoNode(PluginEnvironment env, IDialogTopicGetter topic, IDialogResponsesGetter info)
    {
        // r2a: 名指し話者。ANAM、GetIsID（NPC_/TACT）、GetIsAliasRef（quest alias の ALFR/ALUA）。
        var speakerIds = new List<FormKey>();
        var factionIds = new List<FormKey>();
        var raceIds = new List<FormKey>();
        var voiceTypeIds = new List<FormKey>();

        if (!info.Speaker.IsNull)
            speakerIds.Add(info.Speaker.FormKey);

        foreach (var cond in info.Conditions)
        {
            // 否定形（GetIsID X == 0 等）は「X ではない」の主張なので話者解決に採用しない。
            if (!AssertsTrue(cond)) continue;
            switch (cond.Data)
            {
                case IGetIsIDConditionDataGetter d:
                {
                    var key = d.Object.Link.FormKey;
                    if (env.LinkCache.TryResolve<INpcGetter>(key, out _) ||
                        env.LinkCache.TryResolve<ITalkingActivatorGetter>(key, out _))
                        speakerIds.Add(key);
                    break;
                }
                case IGetIsAliasRefConditionDataGetter d:
                {
                    var actor = ResolveAliasActor(env, topic, d.ReferenceAliasIndex);
                    if (actor != null) speakerIds.Add(actor.Value);
                    break;
                }
                case IGetInFactionConditionDataGetter d:
                    factionIds.Add(d.Faction.Link.FormKey);
                    break;
                case IGetIsRaceConditionDataGetter d:
                    raceIds.Add(d.Race.Link.FormKey);
                    break;
                case IGetIsVoiceTypeConditionDataGetter d:
                {
                    // VTYP 直指定だけ採用する（FLST 指定はプール扱いで対象外、Pascal 版と同じ）。
                    var key = d.VoiceTypeOrList.Link.FormKey;
                    if (env.LinkCache.TryResolve<IVoiceTypeGetter>(key, out _))
                        voiceTypeIds.Add(key);
                    break;
                }
            }
        }

        return new InfoNode
        {
            Id = info.FormKey,
            EditorId = Edid(info),
            IdenticalToMaster = !env.OwnsRecord(info),
            Prompt = S(info.Prompt),
            Responses = info.Responses
                .Select(r => new ResponseLine(r.ResponseNumber, S(r.Text)))
                .Where(r => r.Text.Length > 0)
                .ToList(),
            // 同じ対象が ANAM と条件の両方から来ても 1 回だけ持つ（重複排除）。
            SpeakerIds = speakerIds.Distinct().ToList(),
            FactionIds = factionIds.Distinct().ToList(),
            RaceIds = raceIds.Distinct().ToList(),
            VoiceTypeIds = voiceTypeIds.Distinct().ToList(),
            PreviousId = info.PreviousDialog.IsNull ? null : info.PreviousDialog.FormKey,
            NextIds = info.LinkTo.Where(l => !l.IsNull).Select(l => l.FormKey).ToList(),
            Conditions = info.Conditions.Select(BuildCondition).ToList(),
        };
    }

    // 条件が「関数が真（=1）」を主張しているか。boolean 関数の否定は比較値で表現される
    //（GetIsID X == 0 は「X ではない」）。global 値との比較は静的に評価できないため
    // 真の主張として扱う（従来挙動の維持）。
    private static bool AssertsTrue(IConditionGetter cond)
    {
        if (cond is not IConditionFloatGetter f) return true;
        return cond.CompareOperator switch
        {
            CompareOperator.EqualTo => f.ComparisonValue != 0,
            CompareOperator.NotEqualTo => f.ComparisonValue == 0,
            CompareOperator.GreaterThan => f.ComparisonValue is >= 0 and < 1,
            CompareOperator.GreaterThanOrEqualTo => f.ComparisonValue is > 0 and <= 1,
            // LessThan / LessThanOrEqualTo は「真の時だけ通る」形にならない。
            _ => false,
        };
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
            // locale の小数点記号に依存しない表記に固定する。
            IConditionFloatGetter f => f.ComparisonValue.ToString("0.######", System.Globalization.CultureInfo.InvariantCulture),
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

    private static void ExtractQuests(PluginEnvironment env, ExtractionResult result)
    {
        foreach (var quest in env.TargetMod.Quests)
        {
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
                IdenticalToMaster = !env.OwnsRecord(quest),
                // journal 文（log entry / objective）が無いクエストの FULL は画面に出ない
                //（journal 表示には CNAM か NNAM が必要）。script / scene 制御用と推定する。
                NotPlayerFacing = stages.Count == 0 && objectives.Count == 0,
            });
        }

    }

    // ===== 単純カテゴリ（FULL / DESC 系）=====

    private static void ExtractSimpleCategories(PluginEnvironment env, ExtractionResult result)
    {
        var mod = env.TargetMod;

        foreach (var race in mod.Races)
        {
            var name = S(race.Name);
            var desc = S(race.Description);
            if (!Any(name, desc)) continue;
            result.Races.Add(new RaceEntry
            {
                Id = race.FormKey, EditorId = Edid(race), Kind = "RACE",
                Name = name, Description = desc, IdenticalToMaster = !env.OwnsRecord(race),
            });
        }

        foreach (var fact in mod.Factions)
        {
            var owned = env.OwnsRecord(fact);
            var ranks = fact.Ranks
                .Select(r => new RankTitle(S(r.Title?.Male), S(r.Title?.Female)))
                .Where(r => Any(r.Male, r.Female))
                .ToList();
            // 所有 FACT は名前空でも emit する（Pascal 版と同じ。翻訳 context に効く可能性があるため）。
            if (!owned && S(fact.Name).Length == 0 && ranks.Count == 0) continue;
            result.Factions.Add(new FactionEntry
            {
                Id = fact.FormKey, EditorId = Edid(fact), Kind = "FACT",
                Name = S(fact.Name), Ranks = ranks, IdenticalToMaster = !owned,
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
            AddActivator(result, a.FormKey, Edid(a), "ACTI", S(a.Name), S(a.ActivateTextOverride), !env.OwnsRecord(a));
        foreach (var a in mod.Florae)
            AddActivator(result, a.FormKey, Edid(a), "FLOR", S(a.Name), S(a.ActivateTextOverride), !env.OwnsRecord(a));
        foreach (var a in mod.Trees)
            AddActivator(result, a.FormKey, Edid(a), "TREE", S(a.Name), "", !env.OwnsRecord(a));

        // Equipment（FULL + DESC + EITM）: WEAP/ARMO/AMMO。
        // NonPlayable flag の装備は inventory に出ない（NPC 専用 skin 等）。
        foreach (var w in mod.Weapons)
            AddDescribed(result.Equipment, w.FormKey, Edid(w), "WEAP", S(w.Name), S(w.Description), !env.OwnsRecord(w),
                enchantmentId: w.ObjectEffect.IsNull ? null : w.ObjectEffect.FormKey,
                notPlayerFacing: w.MajorFlags.HasFlag(Weapon.MajorFlag.NonPlayable));
        foreach (var a in mod.Armors)
            AddDescribed(result.Equipment, a.FormKey, Edid(a), "ARMO", S(a.Name), S(a.Description), !env.OwnsRecord(a),
                enchantmentId: a.ObjectEffect.IsNull ? null : a.ObjectEffect.FormKey,
                notPlayerFacing: a.MajorFlags.HasFlag(Armor.MajorFlag.NonPlayable));
        foreach (var a in mod.Ammunitions)
            AddDescribed(result.Equipment, a.FormKey, Edid(a), "AMMO", S(a.Name), S(a.Description), !env.OwnsRecord(a),
                notPlayerFacing: a.Flags.HasFlag(Ammunition.Flag.NonPlayable));

        // Consumable（FULL + DESC + Effects）: ALCH/SCRL/INGR。
        // Effects は翻訳エンジン側の参照グラフ用。MGEF の翻訳所有は record data の
        // 正規化比較（OwnsRecord）だけで件数一致が成立し、参照経由の特別扱いは不要
        //（KnownDeltas 表が監視している）。
        foreach (var c in mod.Ingestibles)
            AddDescribed(result.Consumables, c.FormKey, Edid(c), "ALCH", S(c.Name), S(c.Description), !env.OwnsRecord(c), Effects(c.Effects));
        foreach (var c in mod.Scrolls)
            AddDescribed(result.Consumables, c.FormKey, Edid(c), "SCRL", S(c.Name), S(c.Description), !env.OwnsRecord(c), Effects(c.Effects));
        foreach (var c in mod.Ingredients)
            AddDescribed(result.Consumables, c.FormKey, Edid(c), "INGR", S(c.Name), "", !env.OwnsRecord(c), Effects(c.Effects));

        // Magic / Enchantment / MagicEffect
        foreach (var s in mod.Spells)
            AddDescribed(result.Magic, s.FormKey, Edid(s), "SPEL", S(s.Name), S(s.Description), !env.OwnsRecord(s), Effects(s.Effects));
        foreach (var e in mod.ObjectEffects)
        {
            var owned = env.OwnsRecord(e);
            // 所有 ENCH は名前空でも emit する（装備品付与で意味を持つ。Pascal 版と同じ）。
            if (!owned && S(e.Name).Length == 0) continue;
            result.Enchantments.Add(new DescribedEntry
            {
                Id = e.FormKey, EditorId = Edid(e), Kind = "ENCH",
                Name = S(e.Name), Description = "", MagicEffectIds = Effects(e.Effects),
                IdenticalToMaster = !owned,
            });
        }
        // HideInUI flag の MGEF は effect 名も説明も UI に出ない。
        foreach (var m in mod.MagicEffects)
            AddDescribed(result.MagicEffects, m.FormKey, Edid(m), "MGEF", S(m.Name), S(m.Description), !env.OwnsRecord(m),
                notPlayerFacing: m.Flags.HasFlag(MagicEffect.Flag.HideInUI));

        // Shout + Word of Power
        foreach (var sh in mod.Shouts)
        {
            var owned = env.OwnsRecord(sh);
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
                IdenticalToMaster = !owned,
            });
        }
        foreach (var w in mod.WordsOfPower)
        {
            var spelling = S(w.Name);          // FULL（龍語綴り、翻訳禁止）
            var translation = S(w.Translation); // TNAM（翻訳対象）
            if (!Any(spelling, translation)) continue;
            result.Words.Add(new WordOfPowerEntry
            {
                Id = w.FormKey, EditorId = Edid(w), Kind = "WOOP",
                DragonSpelling = spelling, Translation = translation,
                IdenticalToMaster = !env.OwnsRecord(w),
            });
        }

        // Book
        foreach (var b in mod.Books)
        {
            var title = S(b.Name);
            var body = S(b.BookText);        // DESC（本文）
            var author = S(b.Description);   // CNAM（SSEEdit 表示名 Description）
            if (!Any(title, body, author)) continue;
            result.Books.Add(new BookEntry
            {
                Id = b.FormKey, EditorId = Edid(b), Kind = "BOOK",
                Title = title, Body = body, Author = author,
                IdenticalToMaster = !env.OwnsRecord(b),
            });
        }

        // Location 系（FULL のみ）: LCTN/WRLD/CELL。CELL は worldspace 配下も含めて深く列挙する。
        AddNamed(result.Locations, env, mod.Locations, "LCTN", r => r.Name);
        foreach (var w in mod.Worldspaces)
        {
            var name = S(w.Name);
            if (name.Length == 0) continue;
            result.Locations.Add(new NamedEntry
            {
                Id = w.FormKey, EditorId = Edid(w), Kind = "WRLD", Name = name,
                IdenticalToMaster = !env.OwnsRecord(w),
            });
        }
        foreach (var cell in mod.EnumerateMajorRecords<ICellGetter>())
        {
            // CELL は DIAL と同様に container として働く。plugin が所有する子
            //（新規または変更された配置 ref / navmesh / 地形）を持つ override は
            // record data 不変でも対象になる。子まで master と同一の cell は対象にしない。
            var name = S(cell.Name);
            if (name.Length == 0) continue;
            var hasOwnedChildren =
                cell.Persistent.Concat(cell.Temporary).Any(env.OwnsRecord)
                || cell.NavigationMeshes.Any(env.OwnsRecord)
                || (cell.Landscape != null && env.OwnsRecord(cell.Landscape));
            result.Locations.Add(new NamedEntry
            {
                Id = cell.FormKey, EditorId = Edid(cell), Kind = "CELL", Name = name,
                IdenticalToMaster = !hasOwnedChildren && !env.OwnsRecord(cell),
            });
        }

        // Message（FULL + DESC + ITXT）
        foreach (var m in mod.Messages)
        {
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
                IdenticalToMaster = !env.OwnsRecord(m),
            });
        }


        // LoadScreen（DESC のみ）
        foreach (var l in mod.LoadScreens)
        {
            var body = S(l.Description);
            if (body.Length == 0) continue;
            result.LoadingScreens.Add(new BodyEntry
            {
                Id = l.FormKey, EditorId = Edid(l), Kind = "LSCR", Body = body,
                IdenticalToMaster = !env.OwnsRecord(l),
            });
        }

        // Perk（FULL + DESC。EPF2/EPFD は非対象）
        foreach (var p in mod.Perks)
            AddDescribed(result.Perks, p.FormKey, Edid(p), "PERK", S(p.Name), S(p.Description), !env.OwnsRecord(p));

        // VoiceType（翻訳テキスト無し、識別子のみ）
        foreach (var v in mod.VoiceTypes)
        {
            var edid = Edid(v);
            if (edid.Length == 0) continue;
            result.VoiceTypes.Add(new VoiceTypeEntry
            {
                Id = v.FormKey, EditorId = edid, Kind = "VTYP",
                IdenticalToMaster = !env.OwnsRecord(v),
            });
        }

        // Skyrim 追加 record: SNCT / EYES / REGN（docs/mutagen-migration-plan.md §8-4）
        AddNamed(result.SoundCategories, env, mod.SoundCategories, "SNCT", r => r.Name);
        AddNamed(result.Eyes, env, mod.Eyes, "EYES", r => r.Name);
        foreach (var r in mod.Regions)
        {
            var mapName = S(r.Map?.Name);
            if (mapName.Length == 0) continue;
            result.Regions.Add(new RegionEntry
            {
                Id = r.FormKey, EditorId = Edid(r), Kind = "REGN", MapName = mapName,
                IdenticalToMaster = !env.OwnsRecord(r),
            });
        }
    }

    private static void AddNamed<T>(List<NamedEntry> list, PluginEnvironment env, IEnumerable<T> records, string kind,
        Func<T, ITranslatedStringGetter?> name)
        where T : class, IMajorRecordGetter
    {
        foreach (var rec in records)
        {
            var resolved = S(name(rec));
            if (resolved.Length == 0) continue;
            list.Add(new NamedEntry
            {
                Id = rec.FormKey, EditorId = Edid(rec), Kind = kind, Name = resolved,
                IdenticalToMaster = !env.OwnsRecord(rec),
            });
        }
    }

    private static void AddActivator(ExtractionResult result, FormKey id, string edid, string kind, string name, string action,
        bool identicalToMaster)
    {
        if (!Any(name, action)) return;
        result.Activators.Add(new ActivatorEntry
        {
            Id = id, EditorId = edid, Kind = kind, Name = name, ActivateAction = action,
            IdenticalToMaster = identicalToMaster,
        });
    }

    private static void AddDescribed(List<DescribedEntry> list, FormKey id, string edid, string kind,
        string name, string desc, bool identicalToMaster, List<FormKey>? effects = null, FormKey? enchantmentId = null,
        bool notPlayerFacing = false)
    {
        if (!Any(name, desc)) return;
        list.Add(new DescribedEntry
        {
            Id = id, EditorId = edid, Kind = kind, Name = name, Description = desc,
            MagicEffectIds = effects ?? [], EnchantmentId = enchantmentId,
            IdenticalToMaster = identicalToMaster, NotPlayerFacing = notPlayerFacing,
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
            // 新規 NPC_ は名前空でも emit する（ペルソナ名が必ず解決できる、v18 仕様）。
            // override は text がある場合だけ出す。
            var name = S(npc.Name);
            var shortName = S(npc.ShortName);
            if (npc.FormKey.ModKey != mod.ModKey && !Any(name, shortName)) continue;

            var isFemale = npc.Configuration.Flags.HasFlag(NpcConfiguration.Flag.Female);
            var (personaName, speakerKind) = ResolvePersona(env, npc);
            result.Speakers.Add(new SpeakerEntry
            {
                Id = npc.FormKey,
                EditorId = Edid(npc),
                Kind = "NPC_",
                Name = name,
                PersonaName = personaName,
                SpeakerKind = speakerKind,
                ShortName = shortName,
                Sex = isFemale ? "Female" : "Male",
                VoiceTypeId = npc.Voice.IsNull ? null : npc.Voice.FormKey,
                RaceId = npc.Race.IsNull ? null : npc.Race.FormKey,
                BaseSpeakerId = npc.Template.IsNull ? null : npc.Template.FormKey,
                FactionIds = npc.Factions.Where(f => !f.Faction.IsNull).Select(f => f.Faction.FormKey).ToList(),
                IdenticalToMaster = !env.OwnsRecord(npc),
            });
        }

        foreach (var tact in mod.TalkingActivators)
        {
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
                IdenticalToMaster = !env.OwnsRecord(tact),
            });
        }
    }

    // ペルソナ名と Speaker 種別を TPLT 連鎖の 1 回の走査で解決する
    //（Pascal 版 ClassifySpeakerKind の移植 + 名前解決の統合）。
    // 名前の解決順: 自身の FULL → TPLT 連鎖の本体 FULL → EditorID（役割名）。
    // 種別: 固有（自身に FULL）/ 形態（連鎖先に FULL）/ プール（連鎖が LVLN）/
    //       声型代表（FULL 無しで声型あり）/ 役割のみ（FULL も声型も無し）。
    private static (string PersonaName, string SpeakerKind) ResolvePersona(PluginEnvironment env, INpcGetter npc)
    {
        var own = S(npc.Name);
        if (own.Length > 0) return (own, "固有");

        var curKey = npc.Template.IsNull ? (FormKey?)null : npc.Template.FormKey;
        for (var depth = 0; curKey != null && depth < 12; depth++)
        {
            if (env.LinkCache.TryResolve<ILeveledNpcGetter>(curKey.Value, out _)) return (Edid(npc), "プール");
            if (!env.LinkCache.TryResolve<INpcGetter>(curKey.Value, out var cur)) break;
            var full = S(cur.Name);
            if (full.Length > 0) return (full, "形態");
            curKey = cur.Template.IsNull ? null : cur.Template.FormKey;
        }

        return (Edid(npc), npc.Voice.IsNull ? "役割のみ" : "声型代表");
    }
}
