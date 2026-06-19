#!/usr/bin/env python3
# 評価ハーネス: 公式日本語Destを正解に、過剰置換(precision)と被覆率(recall)を測る。
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
        if rec and src and dst:
            yield rec.group(1), html.unescape(src.group(1)).strip(), html.unescape(dst.group(1)).strip()

files = sorted(glob.glob(f"{XML_DIR}/*.xml"))
def is_bg(f): return any(f.split("/")[-1].startswith(b) for b in BASE_GAME)
base = {}; npc_full_all = []; npc_full_bg = []; npc_shrt = []; dialogue = []
for f in files:
    bg = is_bg(f)
    for rec, s, d in parse(f):
        if rec.endswith(":FULL") and rec != "WOOP:FULL" and s and d and s not in base: base[s] = d
        if rec == "NPC_:FULL":
            npc_full_all.append((s, d));
            if bg: npc_full_bg.append((s, d))
        if rec == "NPC_:SHRT": npc_shrt.append((s, d))
        if rec == "INFO:NAM1": dialogue.append((s, d))

# 用法頻度(安全フィルタ用)
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

# 派生 3 種別を別々に作る(構成比較用)
shrt_d, byname_d, two_d = {}, {}, {}
for s, d in npc_shrt:
    if safe_pair(s, d) and s not in base: shrt_d.setdefault(s, d)
for s, d in npc_full_all:
    if " the " in s.lower():
        en = re.split(r"\s+the\s+", s, flags=re.I)[0].strip(); ja = trailing_kana_run(d)
        if ja and ja != d and safe_pair(en, ja) and en not in base: byname_d.setdefault(en, ja)
for s, d in npc_full_bg:
    if " the " not in s.lower():
        toks, parts = s.split(), d.split("・")
        if len(toks) == len(parts) >= 2 and not has_kanji(d):
            for en, ja in zip(toks, parts):
                en = en.strip()
                if safe_pair(en, ja.strip()) and en not in base: two_d.setdefault(en, ja.strip())

def build(*ds):
    c = dict(base)
    for d in ds:
        for k, v in d.items(): c.setdefault(k, v)
    return c
def compiled(c):
    src = sorted(c.keys(), key=lambda x: (-len(x), x))
    return re.compile(r"\b(?:" + "|".join(re.escape(s) for s in src) + r")\b")

# ---- 正解の名出現(gold): 英語に名形、公式Destに対応カナが実在 ----
# 名形の候補 en->{jp,...}: 単一FULL, SHRT, byname名, 2語FULLの各トークン(・整列)
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
# 一般語/縮約/称号は gold からも除外(名出現とみなさない)
for k in list(name_forms):
    if landmine(k) or len(k) < 4: del name_forms[k]

def eval_config(label, *ds):
    c = build(*ds); reg = compiled(c)
    gold = 0; covered = 0
    ins_total = 0; ins_unconfirmed = 0
    miss_ex, over_ex = [], []
    name_keys = set()
    for d in ds: name_keys |= set(d.keys())
    for en_src, jp_dst in dialogue:
        # gold 名出現を集める
        toks = set(re.findall(r"\b[A-Za-z'][A-Za-z'\-]*\b", en_src))
        line_gold = []
        for t in toks:
            if t in name_forms:
                for jp in name_forms[t]:
                    if jp and jp in jp_dst:
                        line_gold.append((t, jp)); break
        # この行へ辞書適用し、置換結果を得る
        reps = []  # (en, jp)
        def rep(m):
            w = m.group(0); v = c.get(w)
            if v is None: return w
            reps.append((w, v)); return v
        reg.sub(rep, en_src)
        repmap = {}
        for en, jp in reps: repmap.setdefault(en, jp)
        # 被覆: gold 名が正しいカナへ置換されたか(置換後カナが Dest にある)
        for t, jp in line_gold:
            gold += 1
            if t in repmap and repmap[t] in jp_dst: covered += 1
            elif len(miss_ex) < 25 and (t in name_keys or t in repmap): miss_ex.append((t, jp, repmap.get(t)))
        # 過剰: 派生由来の置換のうち Dest に裏が無いもの
        for en, jp in reps:
            if en in name_keys:
                ins_total += 1
                if jp not in jp_dst:
                    ins_unconfirmed += 1
                    if len(over_ex) < 25: over_ex.append((en, jp, en_src[:60]))
    print(f"\n### {label}")
    print(f"  gold名出現={gold}  被覆={covered}  被覆率={covered/gold*100:.1f}%")
    print(f"  派生名の置換={ins_total}  Dest裏なし(過剰の上限)={ins_unconfirmed}  過剰率={ins_unconfirmed/max(1,ins_total)*100:.2f}%")
    return miss_ex, over_ex

print(f"dialogue={len(dialogue)}行  base={len(base)}  shrt={len(shrt_d)} byname={len(byname_d)} two={len(two_d)}")
print(f"gold名形ユニーク={len(name_forms)}")
eval_config("baseのみ(現行)")
eval_config("base+byname+shrt", byname_d, shrt_d)
m, o = eval_config("base+byname+shrt+two", byname_d, shrt_d, two_d)
print("\n--- 過剰(Dest裏なし)の例 上位 ---")
for en, jp, ctx in o[:15]: print(f"   {en}->{jp}  | {ctx}")
print("\n--- 未被覆(miss)の例 上位 ---")
for t, jp, got in m[:15]: print(f"   {t} (正解{jp}) 置換結果={got}")
