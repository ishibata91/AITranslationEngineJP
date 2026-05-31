import type { Page } from "@playwright/test";

interface ScenarioWailsMockOptions {
  masterPersonaAISettings?: "configured" | "missing";
  /**
   * 指定した jobId のジョブは aiTargetCount=0（母数0）として振る舞う。
   * E2E-LTLE-003（境界: 母数0のとき空状態）専用オプション。
   */
  termZeroAITargetJobId?: number;
}

export async function installScenarioWailsMocks(
  page: Page,
  options: ScenarioWailsMockOptions = {},
): Promise<void> {
  const masterPersonaAISettings =
    options.masterPersonaAISettings ?? "configured";
  const termZeroAITargetJobId = options.termZeroAITargetJobId ?? -1;
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

  const termZeroAITargetJobId = ${JSON.stringify(termZeroAITargetJobId)};
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
    GetTermTranslationPhaseSummary: (request) => Promise.resolve(
      request.jobId === 10
        ? termSummary("completed", request.jobId)
        : request.jobId === termZeroAITargetJobId
          ? termSummary("pending", request.jobId, 0)
          : termSummary("pending", request.jobId, request.jobId === lucienJobId ? lucienTermProcessingTargets.length : 3)
    ),
    StartTermTranslationPhase: () => Promise.resolve(termSummary()),
    PauseTermTranslationPhase: () => Promise.resolve(termSummary("paused")),
    ResumeTermTranslationPhase: () => Promise.resolve(termSummary("running")),
    RetryTermTranslationPhase: () => Promise.resolve(termSummary("running")),
    // wave-3 (BE-fact-only-term) 完了後: 可否値を含まず事実状態 (phaseState) のみ返す
    GetTermTranslationNextPhaseReadiness: (request) => Promise.resolve({
      jobId: request.jobId,
      currentPhase: "term_translation",
      phaseState: request.jobId === 10 ? "completed" : "pending"
    }),
    SaveTermTranslationPhaseAISettings: (request) => Promise.resolve({ ...request, phaseId: "term", credentialStatus: "missing", modelListStatus: "credential_missing" }),
    GetPersonaGenerationPhaseSummary: () => Promise.resolve(personaSummary()),
    StartPersonaGenerationPhase: () => Promise.resolve(personaSummary()),
    PausePersonaGenerationPhase: () => Promise.resolve(personaSummary("paused")),
    ResumePersonaGenerationPhase: () => Promise.resolve(personaSummary("running")),
    RetryPersonaGenerationPhase: () => Promise.resolve(personaSummary("running")),
    CancelPersonaGenerationPhase: () => Promise.resolve(personaSummary("canceled")),
    // wave-3 (BE-fact-only-persona) 完了後: 可否値を含まず事実状態 (inputSummary) のみ返す
    GetPersonaGenerationBodyReadiness: () => Promise.resolve({ jobId: 302, currentPhase: "persona_generation", phaseState: "pending", inputSummary: { personaCount: 0, missingCount: 1, snapshotId: "persona-snapshot", snapshotDigest: "persona-digest", evidenceRefs: [] } }),
    SavePersonaGenerationPhaseAISettings: (request) => Promise.resolve({ ...request, phaseId: "persona", credentialStatus: "missing", modelListStatus: "credential_missing" }),
    GetBodyTranslationPhaseSummary: (request) => Promise.resolve(bodySummary(bodyStateForJob(request.jobId), request.jobId)),
    StartBodyTranslationPhase: (request) => Promise.resolve(bodyCommandSummary("pending", request.jobId)),
    PauseBodyTranslationPhase: (request) => Promise.resolve(bodyCommandSummary("paused", request.jobId)),
    ResumeBodyTranslationPhase: (request) => Promise.resolve(bodyCommandSummary("running", request.jobId)),
    RetryBodyTranslationPhase: (request) => Promise.resolve(bodyCommandSummary("running", request.jobId)),
    CancelBodyTranslationPhase: (request) => Promise.resolve(bodyCommandSummary("canceled", request.jobId)),
    // GetBodyTranslationOutputReadiness は wave-3 (BE-fact-only-body) で廃止済み。
    // 成果物出力確認の事実は GetBodyTranslationPhaseSummary の outputReadiness フィールドから取得する。
    SaveBodyTranslationPhaseAISettings: (request) => Promise.resolve({ ...request, phaseId: "body", credentialStatus: "missing", modelListStatus: "credential_missing" }),
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
