# Spec: generic-tone-multi-speaker-sex

---

## R-1 複数話者の性別を汎用口調へ反映する

- R-1-1（正常系）: 複数話者を持つ INFO:NAM1 の汎用セリフの全話者が男性だけの場合に，男性口調を使うこと。
    - 前提条件: 台詞条件の性別が無く，台詞に結ばれた全話者が男性である。
    - 確かめ方: 汎用セリフの口調生成に渡る性別が男性であることを確認する。
    - 対応する実テスト: `TestLinePersonasUsesMultiSpeakerSexForGenericNam1/全男性のNAM1は汎用男性口調`
- R-1-2（対象に入る側の境界）: 複数話者を持つ INFO:NAM1 の台詞が，先頭話者の名指し話者経路ではなく汎用口調を使うこと。
    - 前提条件: INFO:NAM1 の台詞に複数話者が結ばれている。
    - 確かめ方: 汎用セリフの口調生成が台詞本文と全話者の性別を入力に使うことを確認する。
    - 対応する実テスト: `TestLinePersonasUsesMultiSpeakerSexForGenericNam1/全男性のNAM1は汎用男性口調`
- R-1-3（対象に入る側の境界）: 複数話者を持つ INFO:NAM1 の汎用セリフの全話者が女性だけの場合に，女性口調を使うこと。
    - 前提条件: 台詞条件の性別が無く，台詞に結ばれた全話者が女性である。
    - 確かめ方: 汎用セリフの口調生成に渡る性別が女性であることを確認する。
    - 対応する実テスト: `TestLinePersonasUsesMultiSpeakerSexForGenericNam1/全女性のNAM1は汎用女性口調`
- R-1-4（対象に入らない側の境界）: 複数話者を持つ INFO:NAM1 の汎用セリフに男性と女性または性別不明が含まれる場合に，性別を指定しないこと。
    - 前提条件: 台詞条件の性別が無く，台詞に結ばれた話者が男女混在または性別不明である。
    - 確かめ方: 汎用セリフの口調生成に渡る性別が空であることを確認する。
    - 対応する実テスト: `TestLinePersonasUsesMultiSpeakerSexForGenericNam1/男女混在のNAM1は性別を指定しない`，`speaker-tone-injected`
- R-1-5（対象に入らない側の境界）: INFO:RNAM の PC 発話が複数話者を持つ場合に，既存の PC 発話経路を使うこと。
    - 前提条件: INFO:RNAM の台詞に複数話者が結ばれている。
    - 確かめ方: PC 発話が汎用セリフの口調生成へ渡らないことを確認する。
    - 対応する実テスト: `TestLinePersonasUsesMultiSpeakerSexForGenericNam1/RNAMは複数話者でもPC経路`
- R-1-6（正常系）: 台詞条件の性別がある場合に，複数話者の性別集合にかかわらず台詞条件の性別を使うこと。
    - 前提条件: 汎用セリフの台詞条件に男性または女性がある。
    - 確かめ方: 汎用セリフの口調生成に渡る性別が台詞条件の性別であることを確認する。
    - 対応する実テスト: `TestLinePersonasUsesMultiSpeakerSexForGenericNam1/条件性別は話者集合より優先する`
