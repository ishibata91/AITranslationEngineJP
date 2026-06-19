#!/usr/bin/env python3
# 未被覆(miss)の内訳分類: no_entry / shadowed(長一致に食われた) / wrong_jp(同綴り異義)
import re, html, glob
from collections import Counter, defaultdict

XML_DIR = "/Users/iorishibata/Repositories/AITranslationEngineJP/dictionaries/xTranslatorXMLs"
BASE_GAME = ["Skyrim", "Dawnguard", "Dragonborn", "HearthFires", "Update"]
KATA = lambda c: ("゠" <= c <= "ヿ") or c == "ー"
def is_kana_token(s):
    s = s.strip(); return bool(s) and all(KATA(c) or c == "・" for c in s) and any(KATA(c) for c in s)
def has_kanji(s): return any("一" <= c <= "鿿" for c in s)
def trailing_kana_run(d):
    out = []
    for c in reversed(d):
        if KATA(c) or c == "・": out.append(c)
        else: break
    return "".join(reversed(out))
def parse(path):
    t = open(path, encoding="utf-8", errors="replace").read()
    for b in re.findall(r"<String\b.*?</String>", t, re.S):
        rec = re.search(r"<REC>(.*?)</REC>", b); src = re.search(r"<Source>(.*?)</Source>", b, re.S); dst = re.search(r"<Dest>(.*?)</Dest>", b, re.S)
        if rec and src and dst: yield rec.group(1), html.unescape(src.group(1)).strip(), html.unescape(dst.group(1)).strip()
files = sorted(glob.glob(f"{XML_DIR}/*.xml"))
def is_bg(f): return any(f.split("/")[-1].startswith(b) for b in BASE_GAME)
base = {}; npc_full_all = []; npc_full_bg = []; npc_shrt = []; dialogue = []
for f in files:
    bg = is_bg(f)
    for rec, s, d in parse(f):
        if rec.endswith(":FULL") and rec != "WOOP:FULL" and s and d and s not in base: base[s] = d
        if rec == "NPC_:FULL":
            npc_full_all.append((s, d))
            if bg: npc_full_bg.append((s, d))
        if rec == "NPC_:SHRT": npc_shrt.append((s, d))
        if rec == "INFO:NAM1": dialogue.append((s, d))
LC, UC = Counter(), Counter()
for s, _ in dialogue:
    for sent in re.split(r"[.!?]\s+", s):
        toks = re.findall(r"[A-Za-z'][A-Za-z']*", sent)
        for i, tok in enumerate(toks):
            if tok[0].islower():
                for w in re.findall(r"[a-z]+", tok): LC[w] += 1
            elif i > 0: UC[tok.lower().strip("'")] += 1
CONTRACT = set("aren isn wasn weren hasn haven hadn doesn didn don won couldn wouldn shouldn mustn needn daren ain shan mightn oughtn".split())
TITLES = set("master lord lady captain king queen sir madam jarl thane general commander chief elder matron scourge".split())
FACTION = set("dunmer dwemer altmer bosmer breton nord orc khajiit argonian redguard imperial falmer forsworn redoran hlaalu telvanni indoril stormcloak".split())
def landmine(en):
    st = en.lower().strip("'")
    if st in CONTRACT or st in TITLES or st in FACTION: return True
    lc, uc = LC.get(st, 0), UC.get(st, 0)
    return lc >= 3 and lc >= uc
def safe_pair(s, d):
    return s[:1].isupper() and len(s) >= 4 and s.isascii() and re.fullmatch(r"[A-Za-z'\-]+", s) and is_kana_token(d) and len(d) >= 2 and not landmine(s)
derived = {}
for s, d in npc_shrt:
    if safe_pair(s, d) and s not in base: derived.setdefault(s, d)
for s, d in npc_full_all:
    if " the " in s.lower():
        en = re.split(r"\s+the\s+", s, flags=re.I)[0].strip(); ja = trailing_kana_run(d)
        if ja and ja != d and safe_pair(en, ja) and en not in base: derived.setdefault(en, ja)
for s, d in npc_full_bg:
    if " the " not in s.lower():
        toks, parts = s.split(), d.split("・")
        if len(toks) == len(parts) >= 2 and not has_kanji(d):
            for en, ja in zip(toks, parts):
                if safe_pair(en.strip(), ja.strip()) and en.strip() not in base: derived.setdefault(en.strip(), ja.strip())
combined = dict(base)
for k, v in derived.items(): combined.setdefault(k, v)
src = sorted(combined.keys(), key=lambda x: (-len(x), x))
reg = re.compile(r"\b(?:" + "|".join(re.escape(s) for s in src) + r")\b")
# gold
name_forms = defaultdict(set)
for s, d in npc_full_all:
    if " " not in s and "-" not in s and is_kana_token(d): name_forms[s].add(d)
for s, d in npc_shrt:
    if is_kana_token(d): name_forms[s].add(d)
for s, d in npc_full_all:
    if " the " in s.lower():
        en = re.split(r"\s+the\s+", s, flags=re.I)[0].strip(); ja = trailing_kana_run(d)
        if ja and ja != d: name_forms[en].add(ja)
    else:
        toks, parts = s.split(), d.split("・")
        if len(toks) == len(parts) >= 2 and not has_kanji(d):
            for en, ja in zip(toks, parts):
                if en[:1].isupper(): name_forms[en.strip()].add(ja.strip())
for k in list(name_forms):
    if landmine(k) or len(k) < 4: del name_forms[k]

cat = Counter(); examples = defaultdict(list)
gold = covered = covered2 = 0
true_miss = []
for en_src, jp_dst in dialogue:
    toks = set(re.findall(r"\b[A-Za-z'][A-Za-z'\-]*\b", en_src))
    line_gold = []
    for t in toks:
        if t in name_forms:
            for jp in name_forms[t]:
                if jp and jp in jp_dst: line_gold.append((t, jp)); break
    reps = []
    def rep(m):
        w = m.group(0); v = combined.get(w)
        if v is None: return w
        reps.append((w, v)); return v
    applied = reg.sub(rep, en_src)   # 適用後テキスト(長一致含む)
    repmap = {}
    for en, jp in reps: repmap.setdefault(en, jp)
    for t, jp in line_gold:
        gold += 1
        # 正しい被覆: 最終出力に正解カナが含まれる(長一致でも可)
        if jp in applied: covered2 += 1
        elif len(true_miss) < 20: true_miss.append((t, jp, en_src[:60]))
        if t in repmap and repmap[t] in jp_dst:
            covered += 1; continue
        # miss 分類
        if t not in combined:
            c = "no_entry"
        elif t in repmap and repmap[t] != jp and repmap[t] not in jp_dst:
            c = "wrong_jp(同綴り異義)"
        elif t not in repmap:
            c = "shadowed(長一致に食われた)"
        else:
            c = "other"
        cat[c] += 1
        if len(examples[c]) < 12: examples[c].append((t, jp, repmap.get(t), en_src[:55]))
print(f"[単独一致のみ] gold={gold} covered={covered} 被覆率={covered/gold*100:.1f}%")
print(f"[長一致も算入(正しい指標)] covered={covered2} 被覆率={covered2/gold*100:.1f}%  真の未被覆={gold-covered2}")
print("真の未被覆(最終出力にも正解カナ無し)の例:")
for t, jp, ctx in true_miss[:20]:
    print(f"   {t}(正解{jp}) | {ctx}")
print("\n[参考] 単独一致 miss の内訳:")
for c, n in cat.most_common():
    print(f"  {n:5d} ({n/(gold-covered)*100:4.1f}% of miss)  {c}")
for c in cat:
    print(f"\n--- {c} 例 ---")
    for t, jp, got, ctx in examples[c][:8]:
        print(f"   {t}(正解{jp}) 結果={got} | {ctx}")
