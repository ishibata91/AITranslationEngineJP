# デザインシステム仕様: The Ethereal Archive

## 1. 概要と Creative North Star

**Creative North Star: 「Digital Alchemist's Desk」**

このデザインシステムは、無機質で平坦な SaaS の見た目から離れ、触感のある modern fantasy の世界へ移る。
単なる UI を作るのではなく、古い知恵と高密度な技術的明瞭さが同居する、洗練された「Modern Fantasy」の作業環境を構築する。

美術方向は、重厚で歴史を感じる質量感と、発光しながら透ける無重力感の緊張関係で決める。
強い backdrop blur と、残り火のように温かい配色を使い、**マスターペルソナ** と **翻訳ジョブ** を貴重な工芸品として扱う high-end editorial な体験へ寄せる。
画面構成は、縦方向の固定 sidebar よりも、横方向へほどける流れと重なり合う奥行きを優先する。

## 2. 色彩: Ember & Parchment Palette

色の物語は、揺らぐ magitech の残光で構築する。
冷たい gray は避け、深い chocolate と warm amber を基調にして、UI 全体を生きた温度で満たす。

### 2.1 Core Tonal Roles

- **Primary (`#ffba38`)**: 「Amber Glow」。重要操作と active state に使う。システムを動かす魔術的なエネルギーを表す。
- **Secondary (`#e5beb5`)**: 「Soft Bronze」。補助要素に使う、落ち着いた金属感のある色。
- **Surface (`#161311`)**: 「Obsidian Base」。glass panel が載る深い下地。

### 2.2 No-Line Rule と Surface Hierarchy

構造の区切りとして、従来の `1px solid` border は使わない。

- **Depth Stack**: 最も深い背景には `surface_container_lowest`、浮いた panel には `surface_container_highest` を使う。
- **Nesting**: **マスター辞書** の entry card に border は付けない。`surface` 背景の上に `surface_container_low` panel として置く。
- **Glass & Gradient Rule**: 浮遊 panel には半透明の `surface_container` と `backdrop-filter: blur(40px)` を使う。主要 CTA には `primary` から `primary_container` へ向かう繊細な linear gradient を使い、暖かな ember pulse を作る。

## 3. Typography: Scholarly Serif

古い巻物と現代的な洗練をつなぐために、**Noto Serif** を日本語・英語ともに使う。

- **Display (Large / Medium / Small)**: 字間はやや詰め、色は `primary` または `on_surface` を使う。**翻訳ジョブ** の開始見出しのような高位の header に使う。
- **Headline / Title**: content block 内の権威付けに使う。**各翻訳フェーズ** の title は、上質な装丁本の章見出しのように見せる。
- **Body / Label**: `on_surface_variant` (`#d8c3ae`) の高コントラストで可読性を確保する。label は少し letter-spacing を増やし、刻印のような印象を出す。

## 4. Elevation & Depth: Tonal Layering

単純な drop shadow で浮かせるのではなく、光と opacity で物理感を作る。

- **Layering Principle**: 深さは積層で表現する。`surface_container_low` の作業領域の上に、`surface_container_high` の modal を置き、色調差で lift を出す。
- **Ambient Shadows**: glass panel には、拡散の強い shadow を使う。例: `box-shadow: 0 20px 50px rgba(0, 0, 0, 0.5)`。純黒ではなく、背景の暖色を含んだ影にする。
- **Ghost Border**: **マスターペルソナ** 選択のように輪郭が必要な要素では、`outline_variant` を 20% opacity にした `0.5px` border を使う。真鍮の髪線のような premium 感を出し、密度過多は避ける。
- **Runic Accents**: 大きな container の corner には、`primary` を 30% opacity で使った小さな glyph や幾何学的な runic line を置き、fantasy の土台を支える。

## 5. Components: The Alchemical Interface

### 5.1 Navigation: Great Header

- navigation は viewport 上端の高 blur glass bar だけを使う
- sidebar は使わない
- active な nav item は `primary` を使った柔らかい発光 gradient の underline で示す。高さは `0.5px` を基準にする

### 5.2 Buttons: Pulsing Ember

- **Primary**: `primary_container` で塗り、text は `on_primary` にする。hover では外側に繊細な glow を足す
- **Secondary**: 背景は透明、`outline` を使った `0.5px` の Ghost Border を使う
- **Tertiary**: text only とし、低強度 action を `secondary` 色で見せる

### 5.3 Cards & Lists: Scroll Layout

- **No Divider Lines**: **マスター辞書** の list item は divider line で区切らない。`16px` の縦余白と、`surface_container_low` / `surface_container_lowest` の差で分節する
- **Nesting**: **NPCペルソナ生成フェーズ**、**単語翻訳フェーズ**、**本文翻訳フェーズ** は、それぞれ独立した glass panel に収め、**翻訳ジョブ** の進行を見せる

### 5.4 Inputs & Fields

- **Text Fields**: 背景は `surface_container_highest` を薄く使い、下辺だけに Ghost Border を敷く。label は `label-sm` の Serif として field 上部へ浮かせる
- **Progress Indicators**: 標準的な spinner ではなく、柔らかな amber pulse animation を使い、魔術的な処理中であることを示す

## 6. Do と Don't

### 6.1 Do

- 意図的な非対称を使う。**マスターペルソナ** card の runic accent は一角だけでもよい
- modal や pop-over では glassmorphism を積極的に使う。背面 blur は強くし、premium 感を出す
- 用語は **マスターペルソナ**、**マスター辞書**、**翻訳ジョブ**、**各翻訳フェーズ** に統一する

### 6.2 Don't

- `1px solid` の白や gray border を使わない。古層感が壊れる
- default の system font を使わない。`Noto Serif` が使えない場合も、高品質な Serif を fallback にする
- `0px` の鋭角 corner を使わない。古風であっても、**Roundedness Scale** の既定値 `0.25rem` を守り、現代的な可読性と触りやすさを保つ
- sidebar navigation を使わない。体験は巻物が横へほどけるように感じられる必要がある
