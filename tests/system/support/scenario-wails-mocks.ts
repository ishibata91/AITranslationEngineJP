import type { Page } from "@playwright/test";

interface ScenarioWailsMockOptions {
  masterPersonaAISettings?: "configured" | "missing";
  /**
   * 指定した jobId のジョブは aiTargetCount=0（母数0）として振る舞う。
   * E2E-LTLE-003（境界: 母数0のとき空状態）専用オプション。
   */
  termZeroAITargetJobId?: number;
  /**
   * true のとき、jobId=14 のジョブは execution の provider/model/executionMode が空文字の
   * 未設定状態として振る舞う。
   * E2E-UC-FIX-MODEL-001/002/003（AI モデル設定未完了状態と設定固定経路）専用オプション。
   */
  termAISettingsMissing?: boolean;
  /**
   * true のとき、jobId=15 のジョブは NPC ペルソナ生成段階で execution が未設定状態として振る舞う。
   * E2E-UC-057（AI 設定保存後に pill が「固定済み」へ切り替わる NPC ペルソナ生成段階）専用オプション。
   */
  personaAISettingsMissing?: boolean;
  /**
   * true のとき、jobId=16 のジョブは本文翻訳段階で execution が未設定状態として振る舞う。
   * E2E-UC-058（AI 設定保存後に pill が「固定済み」へ切り替わる本文翻訳段階）専用オプション。
   */
  bodyAISettingsMissing?: boolean;
  /**
   * true のとき、jobId=7 の単語翻訳段階 summary は execution が未設定状態として振る舞う。
   * E2E-UC-051（AI 設定不足の開始操作で単語翻訳が未開始状態を維持する）専用オプション。
   */
  seededTermJobMissingExecution?: boolean;
  /**
   * true のとき、jobId=8 の NPC ペルソナ生成段階 summary は execution が未設定状態として振る舞う。
   * E2E-UC-052（AI 設定不足の開始操作でペルソナ生成が未開始状態を維持する）専用オプション。
   */
  seededPersonaJobMissingExecution?: boolean;
  /**
   * true のとき、jobId=9 の本文翻訳段階 summary は execution が未設定状態として振る舞う。
   * E2E-UC-053（AI 設定不足の開始操作で本文翻訳が未開始状態を維持する）専用オプション。
   */
  seededBodyJobMissingExecution?: boolean;
}

export async function installScenarioWailsMocks(
  page: Page,
  options: ScenarioWailsMockOptions = {},
): Promise<void> {
  const masterPersonaAISettings =
    options.masterPersonaAISettings ?? "configured";
  const termZeroAITargetJobId = options.termZeroAITargetJobId ?? -1;
  const termAISettingsMissing = options.termAISettingsMissing ?? false;
  const personaAISettingsMissing = options.personaAISettingsMissing ?? false;
  const bodyAISettingsMissing = options.bodyAISettingsMissing ?? false;
  const seededTermJobMissingExecution = options.seededTermJobMissingExecution ?? false;
  const seededPersonaJobMissingExecution = options.seededPersonaJobMissingExecution ?? false;
  const seededBodyJobMissingExecution = options.seededBodyJobMissingExecution ?? false;
  await page.addInitScript({
    content: `
(() => {
  const personaItems = [
    {
      identityKey: "Skyrim.esm|000A2C8E|NPC_",
      targetPlugin: "Skyrim.esm",
      formId: "000A2C8E",
      recordType: "NPC_",
      editorId: "Lydia",
      displayName: "Lydia",
      voiceType: "FemaleEvenToned",
      className: "Warrior",
      sourcePlugin: "Skyrim.esm",
      personaSummary: "忠実な従士",
      speechStyle: "丁寧な口調",
      personaBody: "ホワイトランの従士として忠実に話す。",
      updatedAt: "2026-05-25T09:00:00Z"
    },
    {
      identityKey: "Skyrim.esm|00013482|NPC_",
      targetPlugin: "Skyrim.esm",
      formId: "00013482",
      recordType: "NPC_",
      editorId: "Hadvar",
      displayName: "Hadvar",
      voiceType: "MaleEvenToned",
      className: "Soldier",
      sourcePlugin: "Skyrim.esm",
      personaSummary: "帝国軍兵士",
      speechStyle: "落ち着いた口調",
      personaBody: "帝国軍兵士として冷静に話す。",
      updatedAt: "2026-05-25T09:00:00Z"
    },
    {
      identityKey: "Dawnguard.esm|02002B74|NPC_",
      targetPlugin: "Dawnguard.esm",
      formId: "02002B74",
      recordType: "NPC_",
      editorId: "Serana",
      displayName: "Serana",
      voiceType: "FemaleUniqueSerana",
      className: "Vampire",
      sourcePlugin: "Dawnguard.esm",
      personaSummary: "吸血鬼の同行者",
      speechStyle: "慎重な口調",
      personaBody: "長い眠りから覚めた吸血鬼として話す。",
      updatedAt: "2026-05-25T09:00:00Z"
    }
  ];

  const pageFromPersonaItems = (request = {}) => {
    const refresh = request.refresh || {};
    const keyword = String(refresh.keyword || "").toLowerCase();
    const pluginFilter = String(refresh.pluginFilter || "");
    const filtered = personaItems.filter((item) => {
      const matchesKeyword =
        keyword === "" ||
        item.displayName.toLowerCase().includes(keyword) ||
        item.editorId.toLowerCase().includes(keyword) ||
        item.targetPlugin.toLowerCase().includes(keyword);
      const matchesPlugin = pluginFilter === "" || item.targetPlugin === pluginFilter;
      return matchesKeyword && matchesPlugin;
    });
    return {
      page: {
        items: filtered,
        pluginGroups: [
          { targetPlugin: "Skyrim.esm", count: 2 },
          { targetPlugin: "Dawnguard.esm", count: 1 }
        ],
        totalCount: filtered.length,
        page: refresh.page || 1,
        pageSize: refresh.pageSize || 30,
        selectedIdentityKey: filtered[0]?.identityKey
      }
    };
  };

  const configuredAISettings = {
    aiSettings: {
      provider: "gemini",
      model: "gemini-test",
      executionMethod: "single_request"
    },
    providerOptions: [
      { value: "gemini", label: "Gemini", credentialStatus: "configured" },
      { value: "lm_studio", label: "LM Studio", credentialStatus: "not_required" }
    ],
    modelList: {
      provider: "gemini",
      credentialStatus: "configured",
      status: "success",
      models: [{ modelId: "gemini-test", label: "gemini-test" }]
    }
  };

  const missingAISettings = {
    aiSettings: { provider: "", model: "", executionMethod: "single_request" },
    providerOptions: [
      { value: "gemini", label: "Gemini", credentialStatus: "missing" },
      { value: "lm_studio", label: "LM Studio", credentialStatus: "not_required" }
    ],
    modelList: {
      provider: "",
      credentialStatus: "missing",
      status: "credential_missing",
      models: []
    }
  };

  let masterPersonaSettings = ${JSON.stringify(masterPersonaAISettings)} === "missing" ? missingAISettings : configuredAISettings;
  let masterPersonaRunStatus = {
    runState: "入力待ち",
    targetPlugin: "",
    processedCount: 0,
    successCount: 0,
    existingSkipCount: 0,
    currentActorLabel: "",
    message: "入力ファイルを選ぶと状態を表示します。"
  };

  const commonProgress = (processed, total, step) => ({
    percent: total === 0 ? 0 : Math.round((processed / total) * 100),
    processedCount: processed,
    totalCount: total,
    aiTargetCount: total,
    targetCount: total,
    translatedCount: processed,
    skippedCount: 0,
    currentStep: step
  });
  const processingTargetsByPhase = {
    term_translation: [
      {
        id: "term-target-1",
        name: "Dragonborn",
        detail: "共通辞書対象外の固有名詞。AI 翻訳対象語として扱う。",
        titleParts: [
          { text: "原語: Dragonborn" },
          { text: "訳語: ドラゴンボーン" },
          { text: "レコード種別: NPC_:FULL" }
        ],
        metadata: [
          { label: "FormID", value: "0001A001" },
          { label: "原文", value: "Dragonborn" },
          { label: "種別", value: "固有名詞" },
          { label: "レコード種別", value: "NPC_:FULL" }
        ]
      },
      {
        id: "term-target-2",
        name: "Whiterun Guard",
        detail: "共通辞書対象外の用語。AI 翻訳対象語として扱う。",
        titleParts: [
          { text: "原語: Whiterun Guard" },
          { text: "訳語: ホワイトラン衛兵" },
          { text: "レコード種別: NPC_:FULL" }
        ],
        metadata: [
          { label: "FormID", value: "0001A002" },
          { label: "原文", value: "Whiterun Guard" },
          { label: "種別", value: "用語" },
          { label: "レコード種別", value: "NPC_:FULL" }
        ]
      },
      {
        id: "term-target-3",
        name: "Riverwood Trader",
        detail: "辞書にない店舗名。AI 翻訳対象語として扱う。",
        titleParts: [
          { text: "原語: Riverwood Trader" },
          { text: "訳語: リバーウッド・トレーダー" },
          { text: "レコード種別: NPC_:FULL" }
        ],
        metadata: [
          { label: "FormID", value: "0001A003" },
          { label: "原文", value: "Riverwood Trader" },
          { label: "種別", value: "固有名詞" },
          { label: "レコード種別", value: "NPC_:FULL" }
        ]
      }
    ],
    persona_generation: [
      {
        id: "persona-target-1",
        name: "Lydia",
        detail: "名前、FormID、EditorID、NPC 属性で検索できる生成対象。",
        titleParts: [
          { text: "NPC: Lydia" },
          { text: "属性: FemaleEvenToned Warrior" }
        ],
        metadata: [
          { label: "FormID", value: "000A2C8E" },
          { label: "EditorID", value: "Lydia" },
          { label: "属性", value: "FemaleEvenToned Warrior" }
        ]
      },
      {
        id: "persona-target-2",
        name: "Hadvar",
        detail: "名前、FormID、EditorID、NPC 属性で検索できる生成対象。",
        titleParts: [
          { text: "NPC: Hadvar" },
          { text: "属性: MaleEvenToned Soldier" }
        ],
        metadata: [
          { label: "FormID", value: "00013482" },
          { label: "EditorID", value: "Hadvar" },
          { label: "属性", value: "MaleEvenToned Soldier" }
        ]
      }
    ],
    body_translation: [
      {
        id: "body-target-1",
        name: "Lydia burden line",
        detail: "辞書置換対象外の本文翻訳項目。AI 送信対象として扱う。",
        titleParts: [
          { text: "原文: I am sworn to carry your burdens." },
          { text: "訳文: あなたの荷物を背負うと誓いました。" }
        ],
        metadata: [
          { label: "FormID", value: "000A2C8E" },
          { label: "EditorID", value: "LydiaLine" },
          { label: "名前", value: "Lydia burden line" }
        ]
      },
      {
        id: "body-target-2",
        name: "Hadvar opening line",
        detail: "辞書置換対象外の本文翻訳項目。AI 送信対象として扱う。",
        titleParts: [
          { text: "原文: You are finally awake." },
          { text: "訳文: ようやく目が覚めたな。" }
        ],
        metadata: [
          { label: "FormID", value: "00013482" },
          { label: "EditorID", value: "HadvarLine" },
          { label: "名前", value: "Hadvar opening line" }
        ]
      },
      {
        id: "body-target-3",
        name: "Serana caution line",
        detail: "辞書置換対象外の本文翻訳項目。AI 送信対象として扱う。",
        titleParts: [
          { text: "原文: We should be careful." },
          { text: "訳文: 慎重に進むべきね。" }
        ],
        metadata: [
          { label: "FormID", value: "02002B74" },
          { label: "EditorID", value: "SeranaLine" },
          { label: "名前", value: "Serana caution line" }
        ]
      },
      {
        id: "body-target-4",
        name: "Whiterun gate line",
        detail: "辞書置換対象外の本文翻訳項目。AI 送信対象として扱う。",
        titleParts: [
          { text: "原文: The gate is closed." },
          { text: "訳文: 門は閉まっている。" }
        ],
        metadata: [
          { label: "FormID", value: "0001A004" },
          { label: "EditorID", value: "WhiterunGateLine" },
          { label: "名前", value: "Whiterun gate line" }
        ]
      }
    ],
    translation_complete: [
      {
        id: "translation-complete-target-1",
        name: "Lydia dialogue",
        detail: "本文翻訳で保持された訳文として出力管理へ進む前に確認する訳文。",
        titleParts: [
          { text: "原文: I am sworn to carry your burdens." },
          { text: "訳文: あなたの荷物を背負うと誓いました。" }
        ],
        metadata: [
          { label: "FormID", value: "000A2C8E" },
          { label: "EditorID", value: "LydiaLine" },
          { label: "出力状態", value: "ready" }
        ]
      }
    ]
  };
  const processingTargetPhaseLabels = {
    term_translation: "単語翻訳",
    persona_generation: "NPC ペルソナ生成",
    body_translation: "本文翻訳",
    translation_complete: "翻訳完了"
  };
  const lucienInputFileName = "Lucien.esp_Export.json";
  const lucienInputId = 1401;
  const lucienJobId = 1401;
  const lucienTermProcessingTargets = [
    {
      id: "lucien-term-target-1",
      name: "Lucien",
      detail: "データロードから登録した Lucien の単語翻訳対象。",
      titleParts: [
        { text: "原語: Lucien" },
        { text: "" },
        { text: "レコード種別: NPC_:FULL" }
      ],
      metadata: [
        { label: "FormID", value: "0001B001" },
        { label: "原文", value: "Lucien" },
        { label: "種別", value: "固有名詞" },
        { label: "レコード種別", value: "NPC_:FULL" }
      ]
    },
    {
      id: "lucien-term-target-2",
      name: "Dumzbthar",
      detail: "データロードから登録した Lucien の単語翻訳対象。",
      titleParts: [
        { text: "原語: Dumzbthar" },
        { text: "" },
        { text: "レコード種別: MISC:FULL" }
      ],
      metadata: [
        { label: "FormID", value: "0001B002" },
        { label: "原文", value: "Dumzbthar" },
        { label: "種別", value: "固有名詞" },
        { label: "レコード種別", value: "MISC:FULL" }
      ]
    }
  ];
  const processingTargetSearchText = (item) => [
    item.name,
    item.detail,
    ...(item.titleParts || []).map((part) => part.text),
    ...(item.metadata || []).flatMap((entry) => [entry.label, entry.value])
  ].join(" ").toLowerCase();
  const getProcessingTargets = (request = {}) => {
    const phase = String(request.phase || "term_translation");
    const sourceItems =
      Number(request.jobId) === termZeroAITargetJobId && phase === "term_translation"
        ? []
        : Number(request.jobId) === lucienJobId && phase === "term_translation"
          ? lucienTermProcessingTargets
          : processingTargetsByPhase[phase] || processingTargetsByPhase.term_translation;
    const searchQuery = String(request.searchQuery || "");
    const normalizedSearchQuery = searchQuery.trim().toLowerCase();
    const page = Math.max(1, Number(request.page) || 1);
    const pageSize = Math.max(1, Number(request.pageSize) || 30);
    const filteredItems = normalizedSearchQuery === ""
      ? sourceItems
      : sourceItems.filter((item) =>
          processingTargetSearchText(item).includes(normalizedSearchQuery)
        );
    const startIndex = (page - 1) * pageSize;
    const items = filteredItems.slice(startIndex, startIndex + pageSize);
    return {
      items,
      metadata: [
        { label: "段階", value: processingTargetPhaseLabels[phase] || phase },
        { label: "対象件数", value: String(filteredItems.length) },
        { label: "要求ページ", value: String(page) },
        { label: "要求ページサイズ", value: String(pageSize) },
        { label: "検索語", value: searchQuery }
      ],
      page,
      pageSize,
      totalCount: filteredItems.length,
      searchQuery
    };
  };
  const execution = {
    credentialRef: "-",
    provider: "-",
    model: "-",
    executionMode: "batch",
    snapshotDigest: "system-test-digest",
    snapshotVersion: "1"
  };
  const targetSummary = {
    targetCount: 2,
    commonPersonaHitCount: 0,
    commonPersonaMissCount: 2,
    skippedCount: 0,
    skippedReasons: [],
    targetSnapshotId: "persona-snapshot",
    targetSnapshotDigest: "persona-digest"
  };
  const personaExecution = {
    credentialRef: "-",
    provider: "-",
    model: "-",
    executionMode: "batch",
    promptDigest: "prompt-digest",
    inputCount: 1,
    outputCount: 0,
    evidenceRefs: []
  };
  const inputSummary = {
    targetCount: 4,
    skippedReasons: [],
    inputSnapshotRef: "input-snapshot",
    dictionaryDigest: "dictionary-digest",
    personaDigest: "persona-digest",
    metadataDigest: "metadata-digest",
    promptDigest: "prompt-digest"
  };
  const bodyExecution = {
    credentialRef: "-",
    provider: "-",
    model: "-",
    executionMode: "batch",
    requestUnitCount: 4,
    outputCount: 0
  };
  const requestSummary = {
    providerTargetCount: 4,
    exactDictionaryExclusionCount: 0,
    partialDictionaryConstraintCount: 0
  };

  const blockedStartError = {
    errorKind: "secret_redacted",
    reason: "AI 設定未完了。認証状態を確認してください。",
    retryable: true,
    isRedacted: true
  };

  const termSummary = (state = "pending", jobId = 301, targetCount = 3) => ({
    jobId,
    currentPhase: "term_translation",
    phaseState: state,
    phaseRunId: 1,
    progress: state === "completed" ? commonProgress(targetCount, targetCount, "完了") : commonProgress(0, targetCount, "未開始"),
    totalTermCount: targetCount,
    dictionaryHitCount: 0,
    aiTargetCount: targetCount,
    execution,
    // wave-3 (BE-fact-only-term) 完了後: frontend が resultSummary.confirmedCount から次段階開始可否を導出する
    // state=completed では confirmedCount=targetCount (全件確認済み) として次段階開始可を導く
    resultSummary: state === "completed" ? { confirmedCount: targetCount, jobDictionaryAppliedCount: 0, replacementTargetCount: 0, unmatchedCount: 0, providerSkipped: false } : undefined,
    errorSummary: blockedStartError,
    actionEnablement: {
      canStart: state === "pending",
      startBlockedReason: "AI 設定未完了。認証状態を確認してください。",
      canPause: false,
      pauseBlockedReason: "未開始です。",
      canResume: false,
      resumeBlockedReason: "未開始です。",
      canRetry: false,
      retryBlockedReason: "未開始です。"
    }
  });
  // E2E-UC-FIX-MODEL-001/002: execution の provider/model が空文字の未設定状態。
  // presenter の isExecutionConfigured が false を返し、状態 pill が「設定未完了」になる。
  // executionMode は "batch" をデフォルト値として保持し、saveAISettings の executionMode 空チェックを回避する。
  const missingExecution = {
    credentialRef: "",
    provider: "",
    model: "",
    executionMode: "batch",
    snapshotDigest: "",
    snapshotVersion: ""
  };
  const termSummaryMissing = (jobId = 14, targetCount = 3) => ({
    jobId,
    currentPhase: "term_translation",
    phaseState: "pending",
    phaseRunId: 1,
    progress: commonProgress(0, targetCount, "未開始"),
    totalTermCount: targetCount,
    dictionaryHitCount: 0,
    aiTargetCount: targetCount,
    execution: missingExecution,
    resultSummary: undefined,
    // AI 設定未完了は retry ではなく設定不足のため retryable=false にする。
    // retryable=true にすると presenter が isRecoverableFailed=true と判断し canResume=true になり
    // selectPhaseProgressActions で resume が優先されて start が除外される。
    errorSummary: { ...blockedStartError, retryable: false },
    actionEnablement: {
      canStart: false,
      startBlockedReason: "実行設定が未構成のため開始できません。",
      canPause: false,
      pauseBlockedReason: "未開始です。",
      canResume: false,
      resumeBlockedReason: "未開始です。",
      canRetry: false,
      retryBlockedReason: "未開始です。"
    }
  });
  // E2E-UC-FIX-MODEL-003: SaveTermTranslationPhaseAISettings 後に configured execution を持つ summary。
  const configuredExecution = {
    credentialRef: "-",
    provider: "gemini",
    model: "gemini-test",
    executionMode: "single_request",
    snapshotDigest: "system-test-digest",
    snapshotVersion: "1"
  };
  const termSummaryConfigured = (jobId = 14, targetCount = 3) => ({
    ...termSummaryMissing(jobId, targetCount),
    execution: configuredExecution,
    actionEnablement: {
      canStart: true,
      startBlockedReason: undefined,
      canPause: false,
      pauseBlockedReason: "未開始です。",
      canResume: false,
      resumeBlockedReason: "未開始です。",
      canRetry: false,
      retryBlockedReason: "未開始です。"
    }
  });

  // E2E-UC-057: NPC ペルソナ生成段階の AI 設定未完了シナリオ用 summary。
  // executionMode は "batch" をデフォルト値として保持し、saveAISettings の executionMode 空チェックを回避する。
  const missingPersonaExecution = {
    credentialRef: "",
    provider: "",
    model: "",
    executionMode: "batch",
    promptDigest: "",
    inputCount: 0,
    outputCount: 0,
    evidenceRefs: []
  };
  const configuredPersonaExecution = {
    credentialRef: "-",
    provider: "gemini",
    model: "gemini-test",
    executionMode: "single_request",
    promptDigest: "prompt-digest",
    inputCount: 0,
    outputCount: 0,
    evidenceRefs: []
  };
  const personaSummaryMissing = (jobId = 15) => ({
    jobId,
    currentPhase: "persona_generation",
    phaseState: "pending",
    phaseRunId: 2,
    progress: commonProgress(0, 2, "未開始"),
    targetSummary,
    execution: missingPersonaExecution,
    // AI 設定未完了は retry ではなく設定不足のため retryable=false にする。
    // retryable=true にすると presenter が isRecoverableFailed=true と判断し canResume=true になり
    // selectPhaseProgressActions で resume が優先されて start が除外される。
    errorSummary: { ...blockedStartError, retryable: false },
    actionEnablement: {
      canStart: false,
      startBlockedReason: "実行設定が未構成のため開始できません。",
      canPause: false,
      pauseBlockedReason: "未開始です。",
      canResume: false,
      resumeBlockedReason: "未開始です。",
      canRetry: false,
      retryBlockedReason: "未開始です。",
      canCancel: false,
      cancelBlockedReason: "実行中ではありません。"
    }
  });
  const personaSummaryConfigured = (jobId = 15) => ({
    ...personaSummaryMissing(jobId),
    execution: configuredPersonaExecution,
    actionEnablement: {
      canStart: true,
      startBlockedReason: undefined,
      canPause: false,
      pauseBlockedReason: "未開始です。",
      canResume: false,
      resumeBlockedReason: "未開始です。",
      canRetry: false,
      retryBlockedReason: "未開始です。",
      canCancel: false,
      cancelBlockedReason: "実行中ではありません。"
    }
  });
  // E2E-UC-058: 本文翻訳段階の AI 設定未完了シナリオ用 summary。
  // executionMode は "batch" をデフォルト値として保持し、saveAISettings の executionMode 空チェックを回避する。
  const missingBodyExecution = {
    credentialRef: "",
    provider: "",
    model: "",
    executionMode: "batch",
    requestUnitCount: 0,
    outputCount: 0
  };
  const configuredBodyExecution = {
    credentialRef: "-",
    provider: "gemini",
    model: "gemini-test",
    executionMode: "single_request",
    requestUnitCount: 4,
    outputCount: 0
  };
  const bodySummaryMissing = (jobId = 16) => ({
    jobId,
    currentPhase: "body_translation",
    phaseState: "pending",
    phaseRunId: 3,
    progress: commonProgress(0, 4, "未開始"),
    inputSummary,
    requestSummary,
    execution: missingBodyExecution,
    fieldResults: [],
    resultSummary: undefined,
    errorSummary: undefined,
    actionEnablement: {
      canStart: false,
      startBlockedReason: "実行設定が未構成のため開始できません。",
      canPause: false,
      pauseBlockedReason: "実行中ではありません。",
      canResume: false,
      resumeBlockedReason: "再開対象ではありません。",
      canRetry: false,
      retryBlockedReason: "失敗状態ではありません。",
      canCancel: false,
      cancelBlockedReason: "実行中ではありません。"
    },
    outputReadiness: bodyOutputReadiness("pending")
  });
  const bodySummaryConfigured = (jobId = 16) => ({
    ...bodySummaryMissing(jobId),
    execution: configuredBodyExecution,
    actionEnablement: {
      canStart: true,
      startBlockedReason: undefined,
      canPause: false,
      pauseBlockedReason: "実行中ではありません。",
      canResume: false,
      resumeBlockedReason: "再開対象ではありません。",
      canRetry: false,
      retryBlockedReason: "失敗状態ではありません。",
      canCancel: false,
      cancelBlockedReason: "実行中ではありません。"
    }
  });

  const personaSummary = (state = "pending") => ({
    jobId: 302,
    currentPhase: "persona_generation",
    phaseState: state,
    phaseRunId: 2,
    progress: commonProgress(0, 2, "未開始"),
    targetSummary,
    execution: personaExecution,
    errorSummary: blockedStartError,
    actionEnablement: {
      canStart: true,
      startBlockedReason: "AI 設定未完了。認証状態を確認してください。",
      canPause: false,
      pauseBlockedReason: "未開始です。",
      canResume: false,
      resumeBlockedReason: "未開始です。",
      canRetry: state === "failed",
      retryBlockedReason: state === "failed" ? "" : "未開始です。",
      canCancel: state === "running",
      cancelBlockedReason: state === "running" ? "" : "実行中ではありません。"
      // wave-3 (BE-fact-only-persona) 完了後: canStartBodyPhase は frontend が
      // derivePersonaCanStartBodyPhase で summary.resultSummary.bodyReadiness から導出するため除去済み。
    }
  });
  // BodyTranslationOutputReadinessSummary: completedFieldCount / statusConsistent / outputCount
  // wave-3 (BE-fact-only-body) 完了後は ready / blockedReason を DTO が含まない。
  // ただし body gateway の hasBodyPhaseBase バリデーションが通るよう型フィールドのみ返す。
  const bodyOutputReadiness = (state = "pending") => ({
    completedFieldCount: state === "completed" ? 1 : 0,
    statusConsistent: true,
    outputCount: state === "completed" ? 1 : 0
  });
  const bodySummary = (state = "pending", jobId = 303) => ({
    jobId,
    currentPhase: "body_translation",
    phaseState: state,
    phaseRunId: 3,
    progress: state === "completed" ? commonProgress(4, 4, "完了") : commonProgress(0, 4, "未開始"),
    inputSummary,
    requestSummary,
    execution: bodyExecution,
    fieldResults: state === "completed" ? [{
      fieldId: 1,
      fieldLabel: "FULL",
      sourceExcerpt: "I am sworn to carry your burdens.",
      translatedText: "あなたの荷物を背負うと誓いました。",
      outputStatus: "ready",
      protectionValidationResult: "passed",
      protectionValidationSummary: "ok",
      retryCount: 0
    }] : [],
    resultSummary: state === "completed" ? {
      translatedCount: 1,
      failedCount: 0,
      skippedCount: 0,
      protectionFailedCount: 0,
      outputReadyCount: 1,
      outputCount: 1,
      fieldResults: []
    } : undefined,
    errorSummary: state === "failed" ? {
      errorKind: "provider_failure",
      reason: "system-test retryable failure",
      retryable: true,
      isRedacted: true
    } : undefined,
    actionEnablement: {
      canStart: state === "pending",
      startBlockedReason: "AI 設定未完了。認証状態を確認してください。",
      canPause: state === "running",
      pauseBlockedReason: state === "running" ? "" : "実行中ではありません。",
      canResume: false,
      resumeBlockedReason: "再開対象ではありません。",
      canRetry: state === "failed",
      retryBlockedReason: state === "failed" ? "" : "失敗状態ではありません。",
      canCancel: state === "running",
      cancelBlockedReason: state === "running" ? "" : "実行中ではありません。"
      // wave-3 (BE-fact-only-body) 完了後: canCheckOutputReadiness は frontend が
      // deriveBodyOutputReadinessReady で summary.outputReadiness から導出するため除去済み。
    },
    outputReadiness: bodyOutputReadiness(state)
  });
  const bodyCommandSummary = (state = "pending", jobId = 303) => ({
    ...bodySummary(state, jobId),
    inputSnapshotDigest: "system-test-input-digest",
    retryable: state === "failed",
    errorSummary: state === "pending" ? blockedStartError : bodySummary(state, jobId).errorSummary
  });
  const bodyStateForJob = (jobId) => {
    if (jobId === 11) return "completed";
    if (jobId === 13) return "running";
    if (jobId === 9) return "pending";
    return "failed";
  };

  // E2E-UC-051: jobId=7 用。execution を未設定にしつつ errorSummary.retryable=false にする。
  // pill「設定未完了」を表示させる。retryable=true のままだと presenter が isRecoverableFailed=true と判断し
  // canResume=true になって selectPhaseProgressActions で resume が優先され start が非表示になる。
  // retryable=false にすることで canResume=false を維持し、start ボタン（disabled）が表示される。
  // Playwright は disabled ボタンもクリックできるため、開始操作を実行して状態が変わらないことを証明できる。
  const termSummarySeededMissing = (jobId = 7, targetCount = 3) => ({
    ...termSummary("pending", jobId, targetCount),
    execution: missingExecution,
    errorSummary: { ...blockedStartError, retryable: false }
  });
  // E2E-UC-052: jobId=8 用。execution を未設定にしつつ errorSummary.retryable=false にする。
  // termSummarySeededMissing と同じ理由: retryable=true だと resume が表示されて start が非表示になる。
  const personaSummarySeededMissing = (jobId = 8) => ({
    ...personaSummary("pending"),
    jobId,
    execution: missingPersonaExecution,
    errorSummary: { ...blockedStartError, retryable: false }
  });
  // E2E-UC-053: jobId=9 用。execution だけ未設定にしつつ canStart=true を維持する。
  const bodySummarySeededMissing = (jobId = 9) => ({
    ...bodySummary("pending", jobId),
    execution: missingBodyExecution
  });

  const termZeroAITargetJobId = ${JSON.stringify(termZeroAITargetJobId)};
  const termAISettingsMissingJobId = ${JSON.stringify(termAISettingsMissing)} ? 14 : -1;
  // E2E-UC-FIX-MODEL-001/002/003: AI 設定未完了ジョブの execution 状態をステートフルに管理する。
  // false = 未設定（空文字 execution）、true = 設定済み（configured execution）
  let termAIMissingJobConfigured = false;
  // E2E-UC-051/052/053: seeded phase jobs（jobId=7/8/9）の execution を未設定状態にする。
  // デフォルトの execution は provider="-" を持ち isExecutionConfigured が truthy になるため、
  // テストが期待する「設定未完了」pill を表示するには missingExecution が必要。
  const seededTermJobMissingExecution = ${JSON.stringify(seededTermJobMissingExecution)};
  const seededPersonaJobMissingExecution = ${JSON.stringify(seededPersonaJobMissingExecution)};
  const seededBodyJobMissingExecution = ${JSON.stringify(seededBodyJobMissingExecution)};
  const personaAISettingsMissingJobId = ${JSON.stringify(personaAISettingsMissing)} ? 15 : -1;
  // E2E-UC-057: NPC ペルソナ生成段階で AI 設定未完了ジョブの execution 状態をステートフルに管理する。
  let personaAIMissingJobConfigured = false;
  const bodyAISettingsMissingJobId = ${JSON.stringify(bodyAISettingsMissing)} ? 16 : -1;
  // E2E-UC-058: 本文翻訳段階で AI 設定未完了ジョブの execution 状態をステートフルに管理する。
  let bodyAIMissingJobConfigured = false;
  const seededPhaseJobs = [
    { jobId: 7, label: "system-test-term", state: "Ready", currentPhase: "term_translation", progressPercent: 0 },
    { jobId: 8, label: "system-test-persona", state: "Ready", currentPhase: "persona_generation", progressPercent: 0 },
    { jobId: 9, label: "system-test-body-pending", state: "Ready", currentPhase: "body_translation", progressPercent: 0 },
    { jobId: 10, label: "system-test-completed-term", state: "Ready", currentPhase: "term_translation", progressPercent: 100 },
    { jobId: 11, label: "system-test-body-ready-for-completion", state: "Ready", currentPhase: "body_translation", progressPercent: 100 },
    { jobId: 12, label: "system-test-body-failed", state: "Failed", currentPhase: "body_translation", progressPercent: 68 },
    { jobId: 13, label: "system-test-body-running", state: "Running", currentPhase: "body_translation", progressPercent: 48 }
  ];
  if (termZeroAITargetJobId >= 0) {
    seededPhaseJobs.push({ jobId: termZeroAITargetJobId, label: "system-test-term-zero-ai-target", state: "Ready", currentPhase: "term_translation", progressPercent: 0 });
  }
  if (termAISettingsMissingJobId >= 0) {
    seededPhaseJobs.push({ jobId: termAISettingsMissingJobId, label: "system-test-term-ai-missing", state: "Ready", currentPhase: "term_translation", progressPercent: 0 });
  }
  if (personaAISettingsMissingJobId >= 0) {
    seededPhaseJobs.push({ jobId: personaAISettingsMissingJobId, label: "system-test-persona-ai-missing", state: "Ready", currentPhase: "persona_generation", progressPercent: 0 });
  }
  if (bodyAISettingsMissingJobId >= 0) {
    seededPhaseJobs.push({ jobId: bodyAISettingsMissingJobId, label: "system-test-body-ai-missing", state: "Ready", currentPhase: "body_translation", progressPercent: 0 });
  }

  const job = (jobId, label, state, currentPhase, progressPercent = 0) => ({
    jobId,
    jobState: state,
    jobStateLabel: state,
    stateTone: state === "Failed" ? "danger" : "info",
    canOpenPhase: true,
    inputSource: {
      inputSourceId: jobId,
      inputSourceLabel: label,
      inputSourceKindLabel: "JSON",
      sourcePath: label + ".json",
      pluginName: label + ".esm",
      extractedJsonLabel: label + ".json"
    },
    progress: {
      currentPhase,
      currentPhaseLabel:
        currentPhase === "term_translation" ? "単語翻訳" :
        currentPhase === "persona_generation" ? "NPC ペルソナ生成" : "本文翻訳",
      percent: progressPercent,
      progressLabel: progressPercent + "%",
      lastUpdatedLabel: "2026-05-25 09:00"
    },
    stopAvailability: { kind: "stop", enabled: state === "Running", label: "停止", helperText: "", reasonText: state === "Running" ? "" : "実行中ではありません。" },
    resumeAvailability: { kind: "resume", enabled: false, label: "再開", helperText: "", reasonText: "再開対象ではありません。" },
    deleteAvailability: { kind: "delete", enabled: true, label: "削除", helperText: "" }
  });

  const jobDetail = (jobId) => {
    const jobsById = Object.fromEntries(
      seededPhaseJobs.map((seededJob) => [
        seededJob.jobId,
        job(
          seededJob.jobId,
          seededJob.label,
          seededJob.state,
          seededJob.currentPhase,
          seededJob.progressPercent
        )
      ])
    );
    return jobsById[jobId] || job(jobId, "system-test-detail", "Ready", "term_translation");
  };

  const lucienImportSummary = (request = {}) => ({
    accepted: true,
    summary: {
      input: {
        id: lucienInputId,
        sourceFilePath: request.filePath || request.fileName || lucienInputFileName,
        sourceTool: "xEdit",
        targetPluginName: "Lucien.esp",
        targetPluginType: "ESP",
        recordCount: 2,
        importedAt: "2026-05-25T09:00:00Z"
      },
      translationRecordCount: 2,
      translationFieldCount: lucienTermProcessingTargets.length,
      categories: [
        { category: "NPC_", recordCount: 1, fieldCount: 1 },
        { category: "DIAL", recordCount: 1, fieldCount: 1 }
      ],
      sampleFields: [
        {
          recordType: "NPC_",
          subrecordType: "FULL",
          formId: "0001B001",
          editorId: "Lucien",
          sourceText: "Lucien",
          translatable: true
        }
      ],
      warnings: []
    },
    warnings: []
  });

  const outputReview = (selectedJobId = 401) => ({
    completedJobs: [
      { jobId: 401, jobStatus: "completed", artifactStatus: "none", outputReady: true, translatedCount: 2, outputStatusDistribution: { ready: 2 } },
      { jobId: 402, jobStatus: "completed", artifactStatus: "current", outputReady: true, translatedCount: 2, outputStatusDistribution: { ready: 2 } },
      { jobId: 403, jobStatus: "completed", artifactStatus: "current", outputReady: true, translatedCount: 0, outputStatusDistribution: {} }
    ],
    hasSelectedJob: true,
    selectedJob: {
      jobId: selectedJobId,
      jobStatus: "completed",
      bodyPhaseStatus: "completed",
      outputReady: selectedJobId !== 403,
      resultSummary: {
        translatedCount: selectedJobId === 403 ? 0 : 2,
        rowCount: selectedJobId === 403 ? 0 : 2,
        inputProvenance: { inputSnapshotDigest: "input-digest", sourceFileDigest: "source-digest" }
      }
    },
    outputReadiness: {
      ready: selectedJobId !== 403,
      retryable: selectedJobId === 401,
      rejectionKind: selectedJobId === 403 ? "status_mismatch" : undefined
    },
    artifactStatus: {
      artifactId: selectedJobId === 401 ? 0 : selectedJobId,
      status: selectedJobId === 401 ? "none" : "current",
      rowCount: selectedJobId === 403 ? 0 : 2,
      currentVersion: selectedJobId !== 401
    },
    rejectionReasons: selectedJobId === 403 ? [{
      errorKind: "status_mismatch",
      reason: "差分なし",
      retryable: false,
      isRedacted: true
    }] : []
  });

  const scenarioAppController = {
    MasterPersonaGetPage: (request) => Promise.resolve(pageFromPersonaItems(request)),
    MasterPersonaGetDetail: (request) => Promise.resolve({ entry: personaItems.find((item) => item.identityKey === request.identityKey) || personaItems[0] }),
    MasterPersonaLoadAISettings: () => Promise.resolve(masterPersonaSettings),
    MasterPersonaListProviderModels: (request) => Promise.resolve({
      provider: request.provider,
      credentialStatus: request.provider === "gemini" && masterPersonaSettings === missingAISettings ? "missing" : "configured",
      status: request.provider === "gemini" && masterPersonaSettings === missingAISettings ? "credential_missing" : "success",
      models: request.provider === "gemini" && masterPersonaSettings === missingAISettings ? [] : [{ modelId: "gemini-test", label: "gemini-test" }]
    }),
    MasterPersonaSaveAISettings: (request) => {
      masterPersonaSettings = {
        ...configuredAISettings,
        aiSettings: request,
        modelList: { ...configuredAISettings.modelList, provider: request.provider, models: [{ modelId: request.model, label: request.model }] }
      };
      return Promise.resolve(request);
    },
    MasterPersonaPreviewGeneration: (request) => Promise.resolve({
      fileName: String(request.filePath || "system-test-persona.json").split("/").pop(),
      targetPlugin: "Skyrim.esm",
      candidateCount: 1,
      newlyAddableCount: masterPersonaSettings === missingAISettings ? 0 : 1,
      existingCount: 0,
      status: masterPersonaSettings === missingAISettings ? "settings_incomplete" : "生成可能"
    }),
    MasterPersonaExecuteGeneration: () => {
      masterPersonaRunStatus = {
        runState: "完了",
        targetPlugin: "Skyrim.esm",
        processedCount: 1,
        successCount: 1,
        existingSkipCount: 0,
        currentActorLabel: "Lydia",
        message: "生成完了"
      };
      return Promise.resolve(masterPersonaRunStatus);
    },
    MasterPersonaGetRunStatus: () => Promise.resolve(masterPersonaRunStatus),
    MasterPersonaUpdate: (request) => {
      const item = personaItems.find((candidate) => candidate.identityKey === request.identityKey) || personaItems[0];
      item.personaSummary = request.entry.personaSummary || item.personaSummary;
      item.speechStyle = request.entry.speechStyle || item.speechStyle;
      item.personaBody = request.entry.personaBody || item.personaBody;
      return Promise.resolve({ ...pageFromPersonaItems(request), changedEntry: item });
    },
    MasterPersonaDelete: (request) => {
      const index = personaItems.findIndex((candidate) => candidate.identityKey === request.identityKey);
      if (index >= 0) personaItems.splice(index, 1);
      return Promise.resolve({ ...pageFromPersonaItems(request), deletedEntryId: request.identityKey });
    },
    ImportTranslationInput: (request) => Promise.resolve(lucienImportSummary(request)),
    RebuildTranslationInputCache: () => Promise.resolve(lucienImportSummary({ filePath: lucienInputFileName })),
    CreateTranslationJobFromInput: () => Promise.resolve({
      accepted: true,
      jobId: lucienJobId,
      jobState: "Ready",
      currentPhase: "term_translation"
    }),
    ListIncompleteJobs: () => Promise.resolve({
      jobs: seededPhaseJobs.map((seededJob) =>
        job(
          seededJob.jobId,
          seededJob.label,
          seededJob.state,
          seededJob.currentPhase,
          seededJob.progressPercent
        )
      )
    }),
    GetJobDetail: (request) => Promise.resolve({ ...jobDetail(request.jobId), cacheState: "available", cacheStateLabel: "available", runtimeSummary: { providerLabel: "-", modelLabel: "-", executionModeLabel: "batch", credentialState: "missing", credentialStateLabel: "設定未完了" }, resumeBlockedReasons: [], warnings: [], deleteImpactLines: [] }),
    GetProcessingTargetList: (request) => Promise.resolve(getProcessingTargets(request)),
    GetTermTranslationPhaseSummary: (request) => {
      if (request.jobId === 10) {
        return Promise.resolve(termSummary("completed", request.jobId));
      }
      if (request.jobId === termZeroAITargetJobId) {
        return Promise.resolve(termSummary("pending", request.jobId, 0));
      }
      // E2E-UC-FIX-MODEL-001/002/003: AI 設定未完了ジョブ。
      // SaveTermTranslationPhaseAISettings 呼び出し後は termAIMissingJobConfigured=true に切り替わり、
      // configured execution を持つ summary を返す。
      if (request.jobId === termAISettingsMissingJobId) {
        return termAIMissingJobConfigured
          ? Promise.resolve(termSummaryConfigured(request.jobId))
          : Promise.resolve(termSummaryMissing(request.jobId));
      }
      // E2E-UC-051: jobId=7 を未設定 execution で返す。
      // デフォルト execution の provider="-" は isExecutionConfigured が truthy になるため
      // pill「固定済み」になってしまう。missingExecution で「設定未完了」状態を実現する。
      // canStart=true を維持して start ボタンを enabled にし、start 操作が状態を変えないことを証明する。
      if (request.jobId === 7 && seededTermJobMissingExecution) {
        return Promise.resolve(termSummarySeededMissing(request.jobId));
      }
      return Promise.resolve(termSummary("pending", request.jobId, request.jobId === lucienJobId ? lucienTermProcessingTargets.length : 3));
    },
    StartTermTranslationPhase: (request) => {
      // E2E-UC-051: seededTermJobMissingExecution 有効時は jobId=7 の start 後も未設定 execution を返す。
      // AI 設定不足のため開始操作が状態を変えないシナリオを実現する。
      if (request && request.jobId === 7 && seededTermJobMissingExecution) {
        return Promise.resolve(termSummarySeededMissing(request.jobId));
      }
      return Promise.resolve(termSummary());
    },
    PauseTermTranslationPhase: () => Promise.resolve(termSummary("paused")),
    ResumeTermTranslationPhase: () => Promise.resolve(termSummary("running")),
    RetryTermTranslationPhase: () => Promise.resolve(termSummary("running")),
    // wave-3 (BE-fact-only-term) 完了後: 可否値を含まず事実状態 (phaseState) のみ返す
    GetTermTranslationNextPhaseReadiness: (request) => Promise.resolve({
      jobId: request.jobId,
      currentPhase: "term_translation",
      phaseState: request.jobId === 10 ? "completed" : "pending"
    }),
    SaveTermTranslationPhaseAISettings: (request) => {
      // E2E-UC-FIX-MODEL-003: termAISettingsMissing オプション有効時（jobId=14 専用シナリオ）、
      // 保存が呼ばれたら configured 状態に切り替える。
      // GetTermTranslationPhaseSummary の次回呼び出しで configured execution を返すようにする。
      if (termAISettingsMissingJobId >= 0) {
        termAIMissingJobConfigured = true;
      }
      // isTermTranslationAISettingsResponseDto が phaseType を必須フィールドとして要求するため含める。
      return Promise.resolve({ ...request, phaseType: "term_translation", phaseId: "term", credentialStatus: "configured", modelListStatus: "success" });
    },
    GetPersonaGenerationPhaseSummary: (request) => {
      // E2E-UC-057: AI 設定未完了ジョブ（jobId=15）は SavePersonaGenerationPhaseAISettings 後に configured summary を返す。
      if (request && request.jobId === personaAISettingsMissingJobId && personaAISettingsMissingJobId >= 0) {
        return personaAIMissingJobConfigured
          ? Promise.resolve(personaSummaryConfigured(request.jobId))
          : Promise.resolve(personaSummaryMissing(request.jobId));
      }
      // E2E-UC-052: jobId=8 を未設定 execution で返す。
      // デフォルト personaExecution の provider="-" は isExecutionConfigured が truthy になるため
      // pill「固定済み」になってしまう。missingPersonaExecution で「設定未完了」状態を実現する。
      // canStart=true を維持して start ボタンを enabled にし、start 操作が状態を変えないことを証明する。
      if (request && request.jobId === 8 && seededPersonaJobMissingExecution) {
        return Promise.resolve(personaSummarySeededMissing(request.jobId));
      }
      return Promise.resolve(personaSummary());
    },
    StartPersonaGenerationPhase: (request) => {
      // E2E-UC-052: seededPersonaJobMissingExecution 有効時は jobId=8 の start 後も未設定 execution を返す。
      // AI 設定不足のため開始操作が状態を変えないシナリオを実現する。
      if (request && request.jobId === 8 && seededPersonaJobMissingExecution) {
        return Promise.resolve(personaSummarySeededMissing(request.jobId));
      }
      return Promise.resolve(personaSummary());
    },
    PausePersonaGenerationPhase: () => Promise.resolve(personaSummary("paused")),
    ResumePersonaGenerationPhase: () => Promise.resolve(personaSummary("running")),
    RetryPersonaGenerationPhase: () => Promise.resolve(personaSummary("running")),
    CancelPersonaGenerationPhase: () => Promise.resolve(personaSummary("canceled")),
    // wave-3 (BE-fact-only-persona) 完了後: 可否値を含まず事実状態 (inputSummary) のみ返す
    GetPersonaGenerationBodyReadiness: () => Promise.resolve({ jobId: 302, currentPhase: "persona_generation", phaseState: "pending", inputSummary: { personaCount: 0, missingCount: 1, snapshotId: "persona-snapshot", snapshotDigest: "persona-digest", evidenceRefs: [] } }),
    SavePersonaGenerationPhaseAISettings: (request) => {
      // E2E-UC-057: personaAISettingsMissing オプション有効時（jobId=15 専用シナリオ）、
      // 保存が呼ばれたら configured 状態に切り替える。
      if (personaAISettingsMissingJobId >= 0) {
        personaAIMissingJobConfigured = true;
      }
      return Promise.resolve({ ...request, phaseId: "persona", credentialStatus: "configured", modelListStatus: "success" });
    },
    GetBodyTranslationPhaseSummary: (request) => {
      // E2E-UC-058: AI 設定未完了ジョブ（jobId=16）は SaveBodyTranslationPhaseAISettings 後に configured summary を返す。
      if (request && request.jobId === bodyAISettingsMissingJobId && bodyAISettingsMissingJobId >= 0) {
        return bodyAIMissingJobConfigured
          ? Promise.resolve(bodySummaryConfigured(request.jobId))
          : Promise.resolve(bodySummaryMissing(request.jobId));
      }
      // E2E-UC-053: jobId=9 を未設定 execution で返す。
      // デフォルト bodyExecution の provider="-" は isExecutionConfigured が truthy になるため
      // pill「固定済み」になってしまう。missingBodyExecution で「設定未完了」状態を実現する。
      // canStart=true を維持して start ボタンを enabled にし、start 操作が状態を変えないことを証明する。
      if (request && request.jobId === 9 && seededBodyJobMissingExecution) {
        return Promise.resolve(bodySummarySeededMissing(request.jobId));
      }
      return Promise.resolve(bodySummary(bodyStateForJob(request.jobId), request.jobId));
    },
    StartBodyTranslationPhase: (request) => {
      // E2E-UC-053: seededBodyJobMissingExecution 有効時は jobId=9 の start 後も未設定 execution を返す。
      // AI 設定不足のため開始操作が状態を変えないシナリオを実現する。
      if (request && request.jobId === 9 && seededBodyJobMissingExecution) {
        return Promise.resolve(bodySummarySeededMissing(request.jobId));
      }
      return Promise.resolve(bodyCommandSummary("pending", request.jobId));
    },
    PauseBodyTranslationPhase: (request) => Promise.resolve(bodyCommandSummary("paused", request.jobId)),
    ResumeBodyTranslationPhase: (request) => Promise.resolve(bodyCommandSummary("running", request.jobId)),
    RetryBodyTranslationPhase: (request) => Promise.resolve(bodyCommandSummary("running", request.jobId)),
    CancelBodyTranslationPhase: (request) => Promise.resolve(bodyCommandSummary("canceled", request.jobId)),
    // GetBodyTranslationOutputReadiness は wave-3 (BE-fact-only-body) で廃止済み。
    // 成果物出力確認の事実は GetBodyTranslationPhaseSummary の outputReadiness フィールドから取得する。
    SaveBodyTranslationPhaseAISettings: (request) => {
      // E2E-UC-058: bodyAISettingsMissing オプション有効時（jobId=16 専用シナリオ）、
      // 保存が呼ばれたら configured 状態に切り替える。
      if (bodyAISettingsMissingJobId >= 0) {
        bodyAIMissingJobConfigured = true;
      }
      // isBodyTranslationAISettingsResponseDto が phaseType を必須フィールドとして要求するため含める。
      return Promise.resolve({ ...request, phaseType: "body_translation", phaseId: "body", credentialStatus: "configured", modelListStatus: "success" });
    },
    // E2E-UC-056/057/058: 「モデル一覧を更新」ボタンが ListTranslationJobSetupProviderModels を呼ぶ。
    // gemini provider に対して gemini-test モデルを返す。
    ListTranslationJobSetupProviderModels: (request) => Promise.resolve({
      provider: request.provider,
      credentialStatus: "configured",
      status: "success",
      models: request.provider === "gemini" ? [{ modelId: "gemini-test", label: "gemini-test" }] : []
    }),
    // E2E-UC-056/057/058: AI サービス provider 一覧を返す。
    // createTermTranslationPhaseScreenControllerFactory が ListProviderSettings を呼ぶ。
    // isListProviderSettingsResponseDto バリデーションを通過するために route と各 provider の
    // credentialState / validationState / savedState フィールドが必要。
    ListProviderSettings: () => Promise.resolve({
      route: {
        routeId: "translation-job-management",
        label: "翻訳管理",
        currentRouteState: "idle",
        dashboardEntryId: "translation-management"
      },
      providers: [
        {
          providerId: "gemini",
          label: "Gemini",
          endpoint: undefined,
          credentialState: "configured",
          credentialReferenceId: undefined,
          validationState: "validated",
          savedState: "configured",
          requestToken: undefined
        },
        {
          providerId: "lm_studio",
          label: "LM Studio",
          endpoint: "http://localhost:1234",
          credentialState: "not_required",
          credentialReferenceId: undefined,
          validationState: "validated",
          savedState: "configured",
          requestToken: undefined
        }
      ]
    }),
    GetTranslationOutputReview: (request) => Promise.resolve(outputReview(request.selectedJobId || 401)),
    GetTranslationOutputDiffPreview: (request) => Promise.resolve({
      jobId: request.jobId,
      artifactId: request.artifactId,
      rows: request.jobId === 403 ? [] : [{
        fieldId: 1,
        rowDigest: "row-digest-1",
        edid: "LydiaLine",
        rec: "INFO",
        field: "NAM1",
        formId: "000A2C8E",
        sourceExcerpt: "I am sworn to carry your burdens.",
        destExcerpt: "あなたの荷物を背負うと誓いました。",
        xTranslatorStatus: 0,
        internalOutputStatus: "ready",
        rowReflectionSummary: "差分あり",
        canRegenerate: true
      }],
      compatibilitySummary: { passed: true, warningCount: 0, rejectCount: 0 }
    }),
    GenerateXTranslatorOutputArtifact: (request) => Promise.resolve({
      jobId: request.jobId,
      artifactId: 501,
      artifactStatus: "current",
      rowCount: 2,
      filePath: request.outputPath,
      targetGame: request.targetGame,
      operationSummary: { operationKind: "generate", replacedArtifactId: 0, duplicateRowCreated: false }
    }),
    RegenerateXTranslatorOutputArtifact: (request) => Promise.resolve({
      jobId: request.jobId,
      artifactId: request.artifactId,
      artifactStatus: "current",
      rowCount: 2,
      filePath: request.outputPath,
      targetGame: request.targetGame,
      operationSummary: { operationKind: "regenerate", replacedArtifactId: request.artifactId, duplicateRowCreated: false }
    })
  };

  let scenarioGo = globalThis.go || {};
  let latestAppController = {};

  const installScenarioControllers = () => {
    const currentGo = globalThis.go || {};
    if (currentGo !== scenarioGo) {
      scenarioGo = currentGo;
    }
    const wails = scenarioGo.wails || {};
    scenarioGo.wails = wails;
    latestAppController = {
      ...latestAppController,
      ...scenarioAppController
    };
    Object.defineProperty(wails, "AppController", {
      configurable: true,
      get: () => latestAppController,
      set: (nextController) => {
        latestAppController = {
          ...(nextController || {}),
          ...scenarioAppController
        };
      }
    });
    wails.MasterPersonaController = wails.AppController;
    wails.TranslationInputController = wails.AppController;
    wails.TranslationOutputArtifactController = wails.AppController;
    wails.PersonaGenerationPhaseController = wails.AppController;
    wails.ProcessingTargetController = wails.AppController;
    Object.defineProperty(globalThis, "go", {
      configurable: true,
      get: () => scenarioGo,
      set: (nextGo) => {
        scenarioGo = nextGo || {};
        installScenarioControllers();
      }
    });
  };

  Object.defineProperty(globalThis, "go", {
    configurable: true,
    get: () => scenarioGo,
    set: (nextGo) => {
      scenarioGo = nextGo || {};
      installScenarioControllers();
    }
  });
  installScenarioControllers();

  globalThis.setInterval(installScenarioControllers, 20);
})();
`,
  });
}
