#!/usr/bin/env python3
# 汎用性テスト: フィルタ基準を Skyrim 本文だけで作り、held-out の Dragonborn 本文で検証。
#  (1) フィルタ判定が corpus 間で揺れるか(過学習の有無)
#  (2) held-out corpus で過剰置換が増えないか
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
def fname(f): return f.split("/")[-1]
def is_bg(f): return any(fname(f).startswith(b) for b in BASE_GAME)

base = {}; npc_full_all = []; npc_full_bg = []; npc_shrt = []
dia_skyrim = []; dia_dragonborn = []
for f in files:
    bg = is_bg(f)
    for rec, s, d in parse(f):
        if rec.endswith(":FULL") and rec != "WOOP:FULL" and s and d and s not in base: base[s] = d
        if rec == "NPC_:FULL":
            npc_full_all.append((s, d))
            if bg: npc_full_bg.append((s, d))
        if rec == "NPC_:SHRT": npc_shrt.append((s, d))
        if rec == "INFO:NAM1":
            if fname(f).startswith("Skyrim"): dia_skyrim.append((s, d))
            if fname(f).startswith("Dragonborn"): dia_dragonborn.append((s, d))

def build_lcuc(dia):
    LC, UC = Counter(), Counter()
    for s, _ in dia:
        for sent in re.split(r"[.!?]\s+", s):
            toks = re.findall(r"[A-Za-z'][A-Za-z']*", sent)
            for i, tok in enumerate(toks):
                if tok[0].islower():
                    for w in re.findall(r"[a-z]+", tok): LC[w] += 1
                elif i > 0: UC[tok.lower().strip("'")] += 1
    return LC, UC
LC_sky, UC_sky = build_lcuc(dia_skyrim)         # フィルタ基準(訓練)
LC_db, UC_db = build_lcuc(dia_dragonborn)       # 別corpus(比較用)

CONTRACT = set("aren isn wasn weren hasn haven hadn doesn didn don won couldn wouldn shouldn mustn needn daren ain shan mightn oughtn".split())
TITLES = set("master lord lady captain king queen sir madam jarl thane general commander chief elder matron scourge".split())
FACTION = set("dunmer dwemer altmer bosmer breton nord orc khajiit argonian redguard imperial falmer forsworn redoran hlaalu telvanni indoril stormcloak".split())
def landmine(en, LC, UC):
    st = en.lower().strip("'")
    if st in CONTRACT or st in TITLES or st in FACTION: return True
    lc, uc = LC.get(st, 0), UC.get(st, 0)
    return lc >= 3 and lc >= uc
def safe_pair(s, d, LC, UC):
    return s[:1].isupper() and len(s) >= 4 and s.isascii() and re.fullmatch(r"[A-Za-z'\-]+", s) and is_kana_token(d) and len(d) >= 2 and not landmine(s, LC, UC)

def derive(LC, UC):
    out = {}
    for s, d in npc_shrt:
        if safe_pair(s, d, LC, UC) and s not in base: out.setdefault(s, d)
    for s, d in npc_full_all:
        if " the " in s.lower():
            en = re.split(r"\s+the\s+", s, flags=re.I)[0].strip(); ja = trailing_kana_run(d)
            if ja and ja != d and safe_pair(en, ja, LC, UC) and en not in base: out.setdefault(en, ja)
    for s, d in npc_full_bg:
        if " the " not in s.lower():
            toks, parts = s.split(), d.split("・")
            if len(toks) == len(parts) >= 2 and not has_kanji(d):
                for en, ja in zip(toks, parts):
                    if safe_pair(en.strip(), ja.strip(), LC, UC) and en.strip() not in base: out.setdefault(en.strip(), ja.strip())
    return out

# 訓練=Skyrim基準で派生
derived = derive(LC_sky, UC_sky)
print(f"派生(Skyrim基準)={len(derived)}")

# (1) フィルタ判定の corpus 間の揺れ: Skyrimでsafeだが Dragonborn基準なら landmine になる key
flip = []
for en in derived:
    if landmine(en, LC_db, UC_db):  # 別corpusでは一般語扱いになる
        flip.append(en)
print(f"\n(1) フィルタ判定の揺れ: Skyrimでは採用だが Dragonborn corpus では地雷判定 = {len(flip)}件 / {len(derived)}")
print("    例:", ", ".join(f"{e}(sky lc{LC_sky.get(e.lower(),0)}/uc{UC_sky.get(e.lower(),0)} db lc{LC_db.get(e.lower(),0)}/uc{UC_db.get(e.lower(),0)})" for e in flip[:10]))

# (2) held-out 過剰置換: Skyrim基準で作った辞書を Dragonborn 本文へ適用
combined = dict(base)
for k, v in derived.items(): combined.setdefault(k, v)
src = sorted(combined.keys(), key=lambda x: (-len(x), x))
reg = re.compile(r"\b(?:" + "|".join(re.escape(s) for s in src) + r")\b")
dset = set(derived)
def measure(dia, label):
    ins = unconf = 0; ex = []
    for en_src, jp_dst in dia:
        for m in reg.finditer(en_src):
            w = m.group(0)
            if w in dset:
                ins += 1
                if combined[w] not in jp_dst:
                    unconf += 1
                    if len(ex) < 15: ex.append((w, combined[w], en_src[:55]))
    print(f"\n(2) {label}: 派生置換={ins} Dest裏なし={unconf} 過剰上限率={unconf/max(1,ins)*100:.2f}%")
    for w, jp, ctx in ex[:10]: print(f"    {w}->{jp} | {ctx}")
measure(dia_skyrim, "in-sample (Skyrim, 訓練と同じ)")
measure(dia_dragonborn, "held-out (Dragonborn, 未使用corpus)")
