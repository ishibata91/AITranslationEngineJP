# 実験の方法論と出典

`research-protocol` skill が固定した判断基準の根拠を持つ。skill は判断基準だけを持ち、出典と手法の詳細と、当てはまる条件と当てはまらない条件はこの資料が持つ。

判断に迷ったとき、または規約の根拠を確かめたいときに読む。英語の題名と手法名は出典に実在する固定名として残す。

**この資料の出典は機械翻訳と自然言語処理の研究に偏っている。** そちらの分野に、繰り返し測って改善する実験の作法が最も整理されているためである。ここに書く主張は、繰り返す実験一般へ一般化して読む。各項目の「当てはまる」「当てはまらない」は、この repo で最初に使った実験（訳文の書き方を測る実験）に当てて書いた具体例であり、別の対象を測る実験では、同じ主張を自分の対象へ読み替える。

## 1 仮説を立てて落とすことの位置づけ

**反証可能性**

- 出典: Karl Popper, *The Logic of Scientific Discovery*（英訳 1959、原著 *Logik der Forschung* 1935）。<https://plato.stanford.edu/entries/popper/>
- 主張: 科学と非科学を分ける基準を、検証できることでなく反証できることに置く。あらゆる観測と両立する理論は科学的でないとみなす。
- 当てはまる: 仮説から「この操作をすればこの値がこの範囲に入らないはずだ」という予測を先に立てる根拠になる。
- 当てはまらない: 単発の実験手続きを規定しておらず、複数の要因が同時に効く実験の設計方法は示さない。

**強い推論（Strong Inference）**

- 出典: John R. Platt, "Strong Inference", *Science* 146(3642), 1964, pp. 347-353。<https://www.stat.cmu.edu/~brian/jdelaney/platt-strong-inference-science-1964.pdf>
- 主張: 3 手続きの反復を定式化する。複数の対立仮説を立てる、そのうち 1 つ以上を排除できる決定的実験を設計する、実験を実行して明確な結果を得る。
- 当てはまる: 仮説を立てて落とすことを繰り返す進め方に最も直接に対応する定式。
- 当てはまらない: 1 回の実験が複数の仮説を同時に排除できることを前提にする。1 回に数百件の翻訳を要し複数の要因が同時に効く状況では、1 回の結果がどの要因を排除したかを一意に切り分けにくい。この切り分けの難しさは項目 2 に直結する。

**多くの作業仮説を並行して持つ**

- 出典: T. C. Chamberlin, "The Method of Multiple Working Hypotheses", *Science*, 1890（再掲載は *The Journal of Geology*, 1897）。<https://www.eps.mcgill.ca/~courses/c590/Chamberlain%201897.pdf>
- 主張: 単一の仮説へ愛着を持って固執することを避けるため、複数の仮説を並行して保持する。

**候補の集合に正しい説明が無い場合（Bad Lot 反論）**

- 出典: Bas van Fraassen による「最良の説明への推論」への批判。整理は <https://link.springer.com/article/10.1007/s11098-017-0933-2>
- 主張: 候補仮説の集合に真の説明が含まれている保証はなく、候補の中で最良でも真とは限らない。仮説空間を尽くせない場合、消去法は理論上収束しない。
- 確かめられなかった: 仮説空間の網羅性を証明する一般的な方法論は確立していない。実務的な緩和（候補を定期的に見直す、「候補の中で最良」と「真」を区別して報告する、収束しないなら要因の切り分け方を疑う）は、上記文献群からの読み取りであり、単一の権威ある手続きとして定式化されているとは確かめられなかった。

## 2 予測が外れたときに何が誤ったか決まらない問題

- 出典: Pierre Duhem, *La Théorie physique*（1906）、W. V. O. Quine, "Two Dogmas of Empiricism"（1951）。総説は <https://plato.stanford.edu/entries/scientific-underdetermination/>
- 主張: 仮説は単独では経験的予測を持たず、補助仮説（測定装置の較正、統計推論の前提、簡略化の仮定）と結びついて初めて検証可能な予測を持つ。予測が外れたとき、本命の仮説と補助仮説のどちらが誤っていたかは論理的に一意に決まらない。
- 当てはまる: 翻訳品質の変化は、変えた対象だけでなく測定手段（数え方、標本の選び方、モデル、接続先）にも依存する。値が動いたとき、原因が対象側か測定手段側かを実験設計だけでは切り分けられない。
- 実務での緩和: 測定手段を実験のたびに変えない、1 実験で変える要因を絞る、補助仮説（採点基準、標本の選び方）を事前に文書化して実験間で一定に保つ。

## 3 探索の段と確証の段を分ける

**探索的データ解析**

- 出典: John W. Tukey, *Exploratory Data Analysis*, Addison-Wesley, 1977。解説は <https://nightingaledvs.com/remembrances-of-things-eda/>
- 主張: データ解析を探索的段（データを見て構造を見つける、柔軟に手法を変えてよい段階）と確証的段（あらかじめ立てた仮説を厳密に検定する段階）に分ける。探索的段で見つけた構造をそのまま確証の根拠にすることを循環論法の罠とする。

**同じ評価用データを繰り返し見ることの問題**

- 出典: Rotem Dror, Gili Baumer, Segev Shlomov, Roi Reichart, "The Hitchhiker's Guide to Testing Statistical Significance in Natural Language Processing", ACL 2018。<https://aclanthology.org/P18-1128/>
- 主張: 同一の評価用データで繰り返し比較を行うことが暗黙のモデル選択になり、検定手続きの選択を事前に固定することを推奨する。
- 補助: 同一評価用データの繰り返し評価による合わせ込みの定量化は <https://arxiv.org/pdf/1905.12580>。ただし画像分類が中心で、機械翻訳固有の知見ではない。一般原則の転用として扱う。

**標本数と検出力**

- 出典: Dallas Card, Peter Henderson, Urvashi Khandelwal, Robin Jia, Kyle Mahowald, Dan Jurafsky, "With Little Power Comes Great Responsibility", EMNLP 2020。<https://aclanthology.org/2020.emnlp-main.745/>
- 主張: 検出したい差の大きさに対して標本数が足りているかを事前に確かめる。論文中の例として、2000 文の評価用データは BLEU 1 点差を検出する検出力が約 75% と報告している。
- 関連: 検出力を 0.80 に固定し、検出したい効果の大きさと有意水準から必要な標本数を逆算する作法は Cohen, J. (1988) *Statistical Power Analysis for the Behavioral Sciences* に由来する。要約は <https://www.ncbi.nlm.nih.gov/pmc/articles/PMC6736231/> と <https://stats.oarc.ucla.edu/other/mult-pkg/faq/general/effect-size-power/faqhow-is-effect-size-used-in-power-analysis/>。原著の一次確認はできなかった。
- 確かめられなかった: 翻訳品質のような評価指標に特化した検出力分析と標本数設計の確立手法は見つからなかった。評価指標が正規分布に従うか、評価者間のばらつきをどう扱うかは、上記の出典だけでは決まらない。

**評価用データの作り方**

- 出典: WMT の General Machine Translation Shared Task の findings。例として <https://aclanthology.org/2025.wmt-1.22/>、<https://www.statmt.org/wmt22/translation-task.html>
- 主張: 複数の領域から均等な件数を層別に抽出する。評価用は公開前に非公開に保ち、開発に使ったデータと切り離す。

## 4 実験の前に固定する

**結果を見てから仮説を作ること**

- 出典: Norbert L. Kerr, "HARKing: Hypothesizing After the Results are Known", *Personality and Social Psychology Review* 2(3), 1998, pp. 196-217。<https://journals.sagepub.com/doi/10.1207/s15327957pspr0203_4>
- 主張: 事後に得た結果に合わせて仮説を作り、それを事前仮説であったかのように報告する行為を問題として定義する。事前登録は、実験の合理性・手続き・分析方法・仮説を実験の前に固定する運用であり、検出は事前登録文書と事後の報告を突き合わせることで行う。
- 当てはまる: 数百件の翻訳を見てから「これが効いていそうだ」と後付けで仮説を語ることが該当する。実験の前に固定すべき対象は、仮説、判定の境界、測定手順、標本数。

**報告に含める項目**

- 出典: ACL Rolling Review, "Responsible NLP Research Checklist"。<https://aclrollingreview.org/static/responsibleNLPresearch.pdf>
- 主張（実験手続きに関わる項目）: モデルの規模と計算予算と計算基盤（C1）、探索した設定の範囲と最終的に採用した値（C2）、結果の記述統計と、報告値が最大値か平均値か 1 回の測定かの明示（C3）、既存パッケージを使った場合の実装とパラメータ設定（C4）、使ったデータの件数と分割の詳細（B6）。
- 確かめられなかった: NeurIPS 側の再現性チェックリストの原文は確認していない。

## 5 要因の振り方

**1 回に 1 つだけ変える進め方の弱点**

- 出典: NIST/SEMATECH e-Handbook of Statistical Methods §5.3.3。<https://www.itl.nist.gov/div898/handbook/pri/section3/pri3347.htm>、<https://en.wikipedia.org/wiki/One-factor-at-a-time_method>
- 主張: 要因の間で効果が変わる関係を推定できない。同じ精度の個別効果の推定に、組み合わせで振る計画より多くの実験回数を要する（3 要因の例で 16 回に対し 8 回）。ある要因を固定した状態で最適点を探すため、固定した水準の選び方に結果が依存する。
- 当てはまる: 指示文が句点を禁じても訳例が全て句点で終わっていると訳例が勝つ、という観測は、文献が記述する要因の間の効果の見落としと一致する。

**組み合わせで振る計画と、一部だけ実施する計画**

- 出典: NIST/SEMATECH e-Handbook §5.3.3.4.5 <https://itl.nist.gov/div898/handbook/pri/section3/pri3345.htm>、§5.3.3.4.7 <https://itl.nist.gov/div898/handbook/pri/section3/pri3347.htm>
- 主張: 要因が 4 つで各 2 水準の場合、全組み合わせは 16 回。resolution IV の半分実施計画は 8 回で、個別効果が 2 要因間の効果と混同されない。全ての 2 要因間の効果を個別に分離するには resolution V 以上か全組み合わせが要る。
- 当てはまる: 4 要因（訳例、指示文、注記、モデル）で、まず 8 回で効く要因を絞る形が現実的な出発点になる。
- 当てはまらない: resolution IV では、どの 2 要因間の効果が効いているかまでは分離できない。

**段階に配分する**

- 出典: Box, G. E. P. and Wilson, K. B. (1951) "On the Experimental Attainment of Optimum Conditions", *Journal of the Royal Statistical Society, Series B*。要旨は <https://en.wikipedia.org/wiki/Response_surface_methodology>
- 主張: まず少ない回数で重要な要因を絞り、その後に絞れた要因について回数を足して要因の間の効果と曲率を探る。全予算を 1 回で使い切らない。
- 当てはまらない: 連続量の水準を前提とした古典的な形は、訳例の書き方・指示文の文言・注記の有無・モデルの選択のようなカテゴリの要因にはそのまま適用しにくい。2 段階に配分する考え方だけを援用する。

**要因を 1 つの塊として扱う**

- 出典: Penn State STAT 502, Lesson 5.1 "Factorial or Crossed Treatment Design"。<https://online.stat.psu.edu/stat502/lesson/5/5.1>
- 主張: 2 つの要因の組み合わせを 1 つの多水準要因として扱う設計は、組み合わせで振る計画の代替として認められた方法である。
- 当てはまらない: 塊にすると各要因の個別の効果は推定できなくなる。どちらを直せば改善するのかという切り分けには答えが出ない。

**実験回数を減らす手法**

- 出典: ベイズ最適化のサーベイ <https://distill.pub/2020/bayesian-optimization/>、<https://arxiv.org/pdf/2107.05847>。最良の選択肢の同定は <https://jmlr.org/papers/volume17/kaufman16a/kaufman16a.pdf>
- 主張: 評価コストが高い関数を代理モデルで近似して次に評価する点を選ぶ、または最も良い選択肢を最少の試行で見つける枠組み。
- 当てはまる: 1 回のコストが高いという前提そのものは一致する。
- 当てはまらない: どちらも最良の組み合わせを当てることが目的で、要因の間の構造を説明する目的には向かない。カテゴリで水準数が少ない離散の要因は、連続パラメータの探索と性質が違う。逐次的な停止規則は、要因を絞った後の個別比較には使える。
- 確かめられなかった: カテゴリの要因に対するベイズ最適化の適用例は、確立した手法の有無を確かめられなかった。

## 6 値の差を改善と呼ぶ手続き

**標本の取り直しで揺れを測る**

- 出典: Philipp Koehn, "Statistical Significance Tests for Machine Translation Evaluation", EMNLP 2004。<https://aclanthology.org/W04-3250/>
- 主張: 同一の評価用データから重複ありの復元抽出で疑似データを大量（通常 1000 回）作り、各回で 2 つの系のスコア差を計算し、差が 0 より大きい割合を p 値として扱う。
- 実装: sacreBLEU の `--paired-bs`（既定 1000 回、乱数種 12345 固定）。<https://github.com/mjpost/sacrebleu>。Python への依存が増える。
- 当てはまる: 同じ原文で対にできる領域ではそのまま使える。人間の訳が無い領域では、指標を原文と訳文だけから測る形へ差し替えたうえで同じ枠組みを適用する。

**出力を入れ替えて比べる**

- 出典: Stefan Riezler, John T. Maxwell III, "On Some Pitfalls in Automatic Evaluation and Significance Testing for MT", ACL Workshop on Intrinsic and Extrinsic Evaluation Measures for MT, 2005。
- 確かめられなかった: 原論文の全文の入手先を確認できなかった（著者名・題名・年のみ確認）。標本の取り直しとどちらが分野の標準かは文献間で結論が割れており、一意に決められなかった。実装が広く使われているのは標本の取り直しの側。

**同時に多く比べるときの補正**

- 出典: Dror ら, ACL 2018（前掲）。<https://aclanthology.org/P18-1128/>、複数データセットの扱いは <https://arxiv.org/pdf/1709.09500>
- 主張: 複数の比較を同時に行う場合、Bonferroni 補正は検出力を落としすぎるため、Holm 法や Benjamini-Yekutieli（従属する仮説向けの偽発見率の補正）を使う。
- 当てはまる: 複数の指標や複数の条件を同一の評価用データで繰り返し比べる場面。
- 当てはまらない: 単一の指標を単一の前後比較にしか使わない場面では補正は不要。

## 7 自動で測った値を目標に据えることの危険

- 出典: Chris Callison-Burch, Miles Osborne, Philipp Koehn, "Re-evaluating the Role of Bleu in Machine Translation Research", EACL 2006。<https://aclanthology.org/E06-1032/>
  - 主張: BLEU の改善は翻訳品質の実改善に対して必要条件でも十分条件でもないことを、反例 2 件で示す。
- 出典: Markus Freitag ら, "Results of WMT22 Metrics Shared Task: Stop Using BLEU – Neural Metrics Are Better and More Robust"。<https://aclanthology.org/2022.wmt-1.2/>
  - 主張: BLEU は人が読む評価との相関がニューラル指標より弱く、BLEU 上位でも人が読む評価で劣る系が観測される。
- 出典: Benjamin Marie, Atsushi Fujita, Raphael Rubino, "Scientific Credibility of Machine Translation Research: A Meta-Evaluation of 769 Papers", ACL 2021。<https://aclanthology.org/2021.acl-long.566/>
  - 主張: 769 本の調査で、BLEU のみに依存し有意性の確認も人が読む評価も行わない論文が増加傾向にある。
- 出典: Markus Freitag ら, "Are LLMs Breaking MT Metrics? Results of the WMT24 Metrics Shared Task"。<https://aclanthology.org/2024.wmt-1.2/>
- 確かめられなかった: 指標を狙って調整した系が人が読む評価で負けた事例について、具体的な系の名前とスコア差は findings 本文の該当箇所まで確認していない。

## 8 人が読む評価の枠組み

- 出典: Arle Lommel, Aljoscha Burchardt, Hans Uszkoreit, "Multidimensional Quality Metrics (MQM): A Framework for Declaring and Describing Translation Quality Metrics", 2014。WMT21 以降の公式手続きとしての採用は <https://aclanthology.org/2021.wmt-1.73/>、<https://aclanthology.org/2023.wmt-1.51/>
- 主張: 訳文中の各誤りを部分文字列として特定し、種類（Accuracy、Fluency、Terminology、Style、Locale、Non-translation）と重さ（major / minor / neutral）を付ける。文単位のスコアでなく誤りの積み上げでスコアを出す。
- 主張（相関の確かめ方）: 系ごとの相関と文ごとの相関の両方を、人が付けた誤りを基準に計算する。WMT21 以降は不特定多数による評価でなくこの枠組みを基準に据えている。
- 当てはまらない: 専門の評価者を前提にした重い手続きで、繰り返す実験サイクルには向かない。抜き取った標本へ適用する形になる。

## 9 人間の訳が無い領域

- 出典: WMT Quality Estimation Shared Task <https://wmt-qe-task.github.io/wmt-qe-2022/>、"CometKiwi: IST-Unbabel 2022 Submission for the Quality Estimation Shared Task" <https://aclanthology.org/2022.wmt-1.60/>
- 主張: 原文と訳文だけからモデルが品質を推定する。人が付けた品質評価との相関で妥当性を確かめる。
- 要る資源: 学習済みモデルと、推論に使う機械。
- 当てはまらない: 参照する訳を使う指標に比べて絶対値の解釈が難しく、何点差を改善と呼ぶかの根拠が整備されていない。揺れを測る枠組み自体は文ごとに同様に適用できる。
