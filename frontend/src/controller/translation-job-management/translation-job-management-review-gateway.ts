import type {
  TranslationJobManagementActionRequest,
  TranslationJobManagementGatewayContract,
  TranslationJobManagementJobDetail,
  TranslationJobManagementJobSummary,
  TranslationJobManagementReasonCategory
} from "@application/gateway-contract/translation-job-management"

interface ReviewJobFixture {
  detail: TranslationJobManagementJobDetail
}

const WAITING_PROMISE = new Promise<never>(() => undefined)

function createOperation(
  kind: "stop" | "resume" | "delete",
  options: {
    enabled: boolean
    label: string
    helperText: string
    reasonCategory?: TranslationJobManagementReasonCategory
    reasonText?: string
  }
) {
  return {
    kind,
    enabled: options.enabled,
    label: options.label,
    helperText: options.helperText,
    reasonCategory: options.reasonCategory,
    reasonText: options.reasonText
  } as const
}

function createFixture(
  detail: Omit<
    TranslationJobManagementJobDetail,
    "stopAvailability" | "resumeAvailability" | "deleteAvailability"
  >,
  operations: {
    stop: ReviewJobFixture["detail"]["stopAvailability"]
    resume: ReviewJobFixture["detail"]["resumeAvailability"]
    delete: ReviewJobFixture["detail"]["deleteAvailability"]
  }
): ReviewJobFixture {
  return {
    detail: {
      ...detail,
      stopAvailability: { ...operations.stop },
      resumeAvailability: { ...operations.resume },
      deleteAvailability: { ...operations.delete }
    }
  }
}

function createBaseFixtures(): ReviewJobFixture[] {
  return [
    createFixture(
      {
        jobId: 401,
        jobState: "Ready",
        jobStateLabel: "Ready",
        stateTone: "info",
        inputSource: {
          inputSourceId: 101,
          inputSourceLabel: "skyrim-main.esp",
          inputSourceKindLabel: "ESP ファイル",
          sourcePath:
            "/Users/example/Mods/VeryLongPath/Skyrim/Data/skyrim-main.esp",
          pluginName: "Skyrim Main Translation Pack Extended Edition",
          extractedJsonLabel: "skyrim-main.extract.json"
        },
        progress: {
          currentPhase: "term_translation",
          currentPhaseLabel: "単語翻訳の準備",
          percent: 0,
          progressLabel: "0% / 未開始",
          lastUpdatedLabel: "5 分前"
        },
        cacheState: "available",
        cacheStateLabel: "入力キャッシュあり",
        runtimeSummary: {
          providerLabel: "Gemini",
          modelLabel: "gemini-2.5-pro-preview",
          executionModeLabel: "Batch API",
          credentialState: "configured",
          credentialStateLabel: "credential 参照あり"
        },
        resumeBlockedReasons: [
          {
            category: "terminal_state",
            title: "再開できません",
            detail:
              "実行前のジョブです。一覧でジョブを選び、単語翻訳ページから開始してください。"
          }
        ],
        warnings: [],
        deleteImpactLines: [
          "削除対象はジョブ本体とジョブ配下の DB 情報です。",
          "入力データと抽出 JSON は残ります。"
        ]
      },
      {
        stop: createOperation("stop", {
          enabled: false,
          label: "停止",
          helperText: "実行中ではありません。",
          reasonText: "停止できません。実行中ではありません。"
        }),
        resume: createOperation("resume", {
          enabled: false,
          label: "再開",
          helperText: "まだ実行されていません。",
          reasonCategory: "terminal_state",
          reasonText: "再開できません。実行前のジョブです。"
        }),
        delete: createOperation("delete", {
          enabled: true,
          label: "削除",
          helperText: "ジョブの DB 情報だけを削除します。"
        })
      }
    ),
    createFixture(
      {
        jobId: 402,
        jobState: "Running",
        jobStateLabel: "Running",
        stateTone: "warning",
        inputSource: {
          inputSourceId: 102,
          inputSourceLabel: "legacy-quest.esp",
          inputSourceKindLabel: "ESP ファイル",
          sourcePath:
            "/Users/example/Mods/LegacyQuest/legacy-quest.esp",
          pluginName: "Legacy Quest Massive Expansion",
          extractedJsonLabel: "legacy-quest.extract.json"
        },
        progress: {
          currentPhase: "body_translation",
          currentPhaseLabel: "本文翻訳",
          percent: 61,
          progressLabel: "61% / 実行中",
          lastUpdatedLabel: "30 秒前"
        },
        cacheState: "available",
        cacheStateLabel: "入力キャッシュあり",
        runtimeSummary: {
          providerLabel: "xAI",
          modelLabel: "grok-4-super-long-model-name-review",
          executionModeLabel: "同期実行",
          credentialState: "configured",
          credentialStateLabel: "credential 参照あり"
        },
        resumeBlockedReasons: [
          {
            category: "terminal_state",
            title: "再開できません",
            detail: "すでに実行中です。"
          }
        ],
        warnings: [
          {
            category: "running_delete_blocked",
            title: "削除は停止後に再判定します",
            detail:
              "Running のジョブは削除できません。停止入口を使い、Paused へ収束した後に削除可否を再判定してください。"
          }
        ],
        deleteImpactLines: [
          "実行中のため削除できません。",
          "先に停止し、停止後に削除可否を再判定してください。"
        ]
      },
      {
        stop: createOperation("stop", {
          enabled: true,
          label: "停止",
          helperText: "停止要求を送信できます。"
        }),
        resume: createOperation("resume", {
          enabled: false,
          label: "再開",
          helperText: "実行中のため再開できません。",
          reasonCategory: "terminal_state",
          reasonText: "再開できません。すでに実行中です。"
        }),
        delete: createOperation("delete", {
          enabled: false,
          label: "削除",
          helperText: "Running のため削除できません。",
          reasonCategory: "running_delete_blocked",
          reasonText:
            "削除できません。実行中のため、先に停止してください。停止後に削除可否を再判定します。"
        })
      }
    ),
    createFixture(
      {
        jobId: 403,
        jobState: "Paused",
        jobStateLabel: "Paused",
        stateTone: "info",
        inputSource: {
          inputSourceId: 103,
          inputSourceLabel: "npc-pack.esm",
          inputSourceKindLabel: "ESM ファイル",
          sourcePath: "/Users/example/Mods/NpcPack/npc-pack.esm",
          pluginName: "NPC Companion Collection",
          extractedJsonLabel: "npc-pack.extract.json"
        },
        progress: {
          currentPhase: "persona_generation",
          currentPhaseLabel: "NPC ペルソナ生成",
          percent: 44,
          progressLabel: "44% / 中断中",
          lastUpdatedLabel: "2 分前"
        },
        cacheState: "available",
        cacheStateLabel: "入力キャッシュあり",
        runtimeSummary: {
          providerLabel: "Gemini",
          modelLabel: "gemini-2.0-flash",
          executionModeLabel: "同期実行",
          credentialState: "configured",
          credentialStateLabel: "credential 参照あり"
        },
        resumeBlockedReasons: [],
        warnings: [],
        deleteImpactLines: [
          "削除対象はジョブ本体とジョブ配下の DB 情報です。",
          "入力データと抽出 JSON は残ります。"
        ]
      },
      {
        stop: createOperation("stop", {
          enabled: false,
          label: "停止",
          helperText: "すでに停止済みです。",
          reasonText: "停止できません。実行中ではありません。"
        }),
        resume: createOperation("resume", {
          enabled: true,
          label: "再開",
          helperText: "中断地点から再開できます。"
        }),
        delete: createOperation("delete", {
          enabled: true,
          label: "削除",
          helperText: "ジョブの DB 情報だけを削除します。"
        })
      }
    ),
    createFixture(
      {
        jobId: 404,
        jobState: "RecoverableFailed",
        jobStateLabel: "RecoverableFailed",
        stateTone: "warning",
        inputSource: {
          inputSourceId: 104,
          inputSourceLabel: "city-overhaul.esp",
          inputSourceKindLabel: "ESP ファイル",
          sourcePath:
            "/Users/example/Mods/CityOverhaul/city-overhaul.esp",
          pluginName: "City Overhaul Plus Roads and Interiors",
          extractedJsonLabel: "city-overhaul.extract.json"
        },
        progress: {
          currentPhase: "body_translation",
          currentPhaseLabel: "本文翻訳",
          percent: 72,
          progressLabel: "72% / 再開可能な失敗",
          lastUpdatedLabel: "7 分前"
        },
        cacheState: "missing",
        cacheStateLabel: "入力キャッシュがありません",
        runtimeSummary: {
          providerLabel: "OpenAI Compatible",
          modelLabel: "custom-model-prod-build-2026-05-01",
          executionModeLabel: "同期実行",
          credentialState: "inaccessible",
          credentialStateLabel: "credential 参照に失敗"
        },
        resumeBlockedReasons: [
          {
            category: "cache_missing",
            title: "再開できません",
            detail: "入力キャッシュがありません。入力確認から再構築してください。"
          }
        ],
        warnings: [],
        deleteImpactLines: [
          "削除対象はジョブ本体とジョブ配下の DB 情報です。",
          "入力データと抽出 JSON は残ります。"
        ]
      },
      {
        stop: createOperation("stop", {
          enabled: false,
          label: "停止",
          helperText: "実行中ではありません。",
          reasonText: "停止できません。実行中ではありません。"
        }),
        resume: createOperation("resume", {
          enabled: false,
          label: "再開",
          helperText: "入力キャッシュの再構築が必要です。",
          reasonCategory: "cache_missing",
          reasonText:
            "再開できません。入力キャッシュを再構築してください。"
        }),
        delete: createOperation("delete", {
          enabled: true,
          label: "削除",
          helperText: "ジョブの DB 情報だけを削除します。"
        })
      }
    ),
    createFixture(
      {
        jobId: 405,
        jobState: "Failed",
        jobStateLabel: "Failed",
        stateTone: "danger",
        canOpenPhase: false,
        openBlockedReason: {
          category: "phase_progress_aggregation_failed",
          title: "翻訳段階を開けません",
          detail:
            "翻訳段階の進捗を確認できないため、一覧で状態を確認してください。"
        },
        inputSource: {
          inputSourceId: 105,
          inputSourceLabel: "museum-addon.esp",
          inputSourceKindLabel: "ESP ファイル",
          sourcePath:
            "/Users/example/Mods/MuseumAddon/museum-addon.esp",
          pluginName: "Museum Addon",
          extractedJsonLabel: "museum-addon.extract.json"
        },
        progress: {
          currentPhase: "term_translation",
          currentPhaseLabel: "進捗投影",
          percent: 0,
          progressLabel: "進捗を確認できません",
          lastUpdatedLabel: "10 分前"
        },
        cacheState: "available",
        cacheStateLabel: "入力キャッシュあり",
        runtimeSummary: {
          providerLabel: "Anthropic",
          modelLabel: "claude-sonnet-enterprise",
          executionModeLabel: "同期実行",
          credentialState: "configured",
          credentialStateLabel: "credential 参照あり"
        },
        resumeBlockedReasons: [
          {
            category: "state_projection_inconsistent",
            title: "再開できません",
            detail: "state projection inconsistent のため、進捗を確認できません。"
          }
        ],
        warnings: [
          {
            category: "phase_progress_aggregation_failed",
            title: "進捗を集約できません",
            detail:
              "翻訳段階の進捗集約に失敗しました。再読込しても直らない場合は backend projection を確認してください。"
          }
        ],
        deleteImpactLines: [
          "削除対象はジョブ本体とジョブ配下の DB 情報です。",
          "入力データと抽出 JSON は残ります。"
        ]
      },
      {
        stop: createOperation("stop", {
          enabled: false,
          label: "停止",
          helperText: "実行中ではありません。",
          reasonText: "停止できません。実行中ではありません。"
        }),
        resume: createOperation("resume", {
          enabled: false,
          label: "再開",
          helperText: "進捗集約を確認できません。",
          reasonCategory: "state_projection_inconsistent",
          reasonText: "再開できません。進捗を確認できません。"
        }),
        delete: createOperation("delete", {
          enabled: true,
          label: "削除",
          helperText: "ジョブの DB 情報だけを削除します。"
        })
      }
    ),
    createFixture(
      {
        jobId: 406,
        jobState: "Canceled",
        jobStateLabel: "Canceled",
        stateTone: "danger",
        inputSource: {
          inputSourceId: 106,
          inputSourceLabel: "worldspace-fixup.esp",
          inputSourceKindLabel: "ESP ファイル",
          sourcePath:
            "/Users/example/Mods/WorldspaceFix/worldspace-fixup.esp",
          pluginName: "Worldspace Fixup",
          extractedJsonLabel: "worldspace-fixup.extract.json"
        },
        progress: {
          currentPhase: "term_translation",
          currentPhaseLabel: "停止済み",
          percent: 100,
          progressLabel: "停止済み / 再開不可",
          lastUpdatedLabel: "12 分前"
        },
        cacheState: "available",
        cacheStateLabel: "入力キャッシュあり",
        runtimeSummary: {
          providerLabel: "LM Studio",
          modelLabel: "local-gguf-q4_k_m-very-very-long-name",
          executionModeLabel: "ローカル実行",
          credentialState: "configured",
          credentialStateLabel: "credential 不要"
        },
        resumeBlockedReasons: [
          {
            category: "terminal_state",
            title: "再開できません",
            detail: "Canceled は terminal state です。新しいジョブを作成してください。"
          }
        ],
        warnings: [],
        deleteImpactLines: [
          "削除対象はジョブ本体とジョブ配下の DB 情報です。",
          "入力データと抽出 JSON は残ります。"
        ]
      },
      {
        stop: createOperation("stop", {
          enabled: false,
          label: "停止",
          helperText: "すでに停止済みです。",
          reasonText: "停止できません。実行中ではありません。"
        }),
        resume: createOperation("resume", {
          enabled: false,
          label: "再開",
          helperText: "terminal state です。",
          reasonCategory: "terminal_state",
          reasonText:
            "再開できません。Canceled は terminal state です。"
        }),
        delete: createOperation("delete", {
          enabled: true,
          label: "削除",
          helperText: "ジョブの DB 情報だけを削除します。"
        })
      }
    )
  ]
}

export function createTranslationJobManagementReviewGateway(
  scenarioId: string
): TranslationJobManagementGatewayContract {
  if (scenarioId === "loading") {
    return {
      ListIncompleteJobs: () => WAITING_PROMISE,
      GetJobDetail: () => WAITING_PROMISE,
      RequestStop: () => WAITING_PROMISE,
      ResumeJob: () => WAITING_PROMISE,
      DeleteJob: () => WAITING_PROMISE
    }
  }

  if (scenarioId === "error") {
    return {
      ListIncompleteJobs: () =>
        Promise.reject(
          new Error("未完了ジョブの一覧取得に失敗しました。再読込してください。")
        ),
      GetJobDetail: () =>
        Promise.reject(
          new Error("選択したジョブを再度読み込めませんでした。")
        ),
      RequestStop: () =>
        Promise.reject(new Error("停止要求に失敗しました。")),
      ResumeJob: () =>
        Promise.reject(new Error("再開に失敗しました。")),
      DeleteJob: () =>
        Promise.reject(new Error("削除に失敗しました。"))
    }
  }

  const fixtures = createBaseFixtures()
  const fixtureMap = new Map(fixtures.map((fixture) => [fixture.detail.jobId, fixture]))

  if (scenarioId === "empty") {
    fixtureMap.clear()
  }

  if (scenarioId === "running") {
    for (const [jobId, fixture] of fixtureMap) {
      if (jobId !== 402) {
        fixtureMap.delete(jobId)
        continue
      }

      fixture.detail.progress.progressLabel =
        "61% / 本文翻訳を実行中です"
      fixture.detail.runtimeSummary.executionModeLabel =
        "実行中のため移動操作だけ確認可能"
      fixture.detail.resumeBlockedReasons = [
        {
          category: "terminal_state",
          title: "実行中です",
          detail:
            "本文翻訳が実行中です。再開や削除は使わず、停止するか進行状況を確認してください。"
        }
      ]
      fixture.detail.warnings = [
        {
          category: "running_delete_blocked",
          title: "実行中の制限",
          detail:
            "Running のジョブは削除できません。一覧では停止だけを次操作として確認します。"
        }
      ]
      fixture.detail.stopAvailability.helperText =
        "実行中です。停止要求を送信できます。"
      fixture.detail.resumeAvailability.helperText =
        "実行中のため再開は不要です。"
      fixture.detail.resumeAvailability.reasonText =
        "再開できません。本文翻訳が実行中です。停止または進行状況確認を行ってください。"
      fixture.detail.deleteAvailability.reasonText =
        "削除できません。本文翻訳が実行中です。停止後に削除可否を再判定します。"
    }
  }

  if (scenarioId === "config-missing") {
    for (const [jobId, fixture] of fixtureMap) {
      if (jobId !== 404) {
        fixtureMap.delete(jobId)
        continue
      }
      fixture.detail.runtimeSummary.credentialState = "missing"
      fixture.detail.runtimeSummary.credentialStateLabel =
        "credential を確認してください"
      fixture.detail.progress.currentPhaseLabel = "本文翻訳の設定確認"
      fixture.detail.progress.progressLabel =
        "設定不足 / API キー確認待ち"
      fixture.detail.cacheState = "available"
      fixture.detail.cacheStateLabel = "入力キャッシュあり"
      fixture.detail.resumeBlockedReasons = [
        {
          category: "cache_missing",
          title: "設定不足です",
          detail:
            "AI サービス設定で API キーと model を確認してください。設定が揃うまで本文翻訳へ進めません。"
        }
      ]
      fixture.detail.warnings = [
        {
          category: "cache_missing",
          title: "次操作",
          detail:
            "設定画面で credential 状態と model 選択を確認し、一覧へ戻ってジョブを選び直してください。"
        }
      ]
      fixture.detail.stopAvailability.enabled = false
      fixture.detail.stopAvailability.helperText =
        "実行中ではありません。設定確認が必要です。"
      fixture.detail.stopAvailability.reasonText =
        "停止できません。設定不足のため実行は始まっていません。"
      fixture.detail.resumeAvailability.enabled = false
      fixture.detail.resumeAvailability.helperText =
        "API キーと model の確認が必要です。"
      fixture.detail.resumeAvailability.reasonCategory = "cache_missing"
      fixture.detail.resumeAvailability.reasonText =
        "再開できません。AI サービス設定で API キーと model を確認してください。"
    }
  }

  return {
    ListIncompleteJobs: () =>
      Promise.resolve({
        jobs: [...fixtureMap.values()].map((fixture) => toSummary(fixture.detail))
      }),
    GetJobDetail: ({ jobId }) => {
      const fixture = fixtureMap.get(jobId)
      if (!fixture) {
        throw new Error(
          "選択したジョブは見つかりません。一覧を更新して選び直してください。"
        )
      }
      return Promise.resolve(cloneDetail(fixture.detail))
    },
    RequestStop: ({ jobId }: TranslationJobManagementActionRequest) => {
      const fixture = requireFixture(fixtureMap, jobId)
      if (!fixture.detail.stopAvailability.enabled) {
        return Promise.resolve({
          tone: "warning",
          message: fixture.detail.stopAvailability.reasonText ?? "停止できません。",
          detail: cloneDetail(fixture.detail),
          reasonCategory: fixture.detail.stopAvailability.reasonCategory
        })
      }

      fixture.detail.stopAvailability.enabled = false
      fixture.detail.stopAvailability.helperText = "停止要求中です。"
      fixture.detail.stopAvailability.reasonCategory = "stop_requested"
      fixture.detail.stopAvailability.reasonText =
        "停止要求中です。Paused へ収束した後に削除可否を再判定してください。"
      fixture.detail.deleteAvailability.enabled = false
      fixture.detail.deleteAvailability.reasonCategory = "running_delete_blocked"
      fixture.detail.deleteAvailability.reasonText =
        "削除できません。停止要求中です。Paused へ収束した後に再判定します。"
      fixture.detail.warnings = [
        {
          category: "stop_requested",
          title: "停止要求中です",
          detail:
            "停止要求を送信しました。停止後に削除可否を再判定してください。"
        }
      ]

      return Promise.resolve({
        tone: "info",
        message:
          "停止要求を送信しました。削除可否は停止後に再判定してください。",
        detail: cloneDetail(fixture.detail),
        reasonCategory: "stop_requested"
      })
    },
    ResumeJob: ({ jobId }: TranslationJobManagementActionRequest) => {
      const fixture = requireFixture(fixtureMap, jobId)
      if (!fixture.detail.resumeAvailability.enabled) {
        return Promise.resolve({
          tone: "warning",
          message:
            fixture.detail.resumeAvailability.reasonText ?? "再開できません。",
          detail: cloneDetail(fixture.detail),
          reasonCategory: fixture.detail.resumeAvailability.reasonCategory
        })
      }

      fixture.detail.warnings = []
      fixture.detail.progress.progressLabel = `${fixture.detail.progress.percent}% / 再開要求を受け付けました`
      fixture.detail.resumeAvailability.enabled = false
      fixture.detail.resumeAvailability.reasonText =
        "再開要求を送信しました。状態更新を待ってください。"

      return Promise.resolve({
        tone: "success",
        message: "再開要求を受け付けました。状態更新を待ってください。",
        detail: cloneDetail(fixture.detail)
      })
    },
    DeleteJob: ({ jobId }) => {
      const fixture = requireFixture(fixtureMap, jobId)
      if (!fixture.detail.deleteAvailability.enabled) {
        return Promise.resolve({
          tone: "warning",
          message:
            fixture.detail.deleteAvailability.reasonText ?? "削除できません。",
          detail: cloneDetail(fixture.detail),
          reasonCategory: fixture.detail.deleteAvailability.reasonCategory
        })
      }

      fixtureMap.delete(jobId)
      return Promise.resolve({
        tone: "success",
        message:
          "ジョブ本体と配下の DB 情報を削除しました。入力データと抽出 JSON は残ります。",
        deletedJobId: jobId
      })
    }
  }
}

function requireFixture(
  fixtureMap: Map<number, ReviewJobFixture>,
  jobId: number
): ReviewJobFixture {
  const fixture = fixtureMap.get(jobId)
  if (!fixture) {
    throw new Error("選択したジョブは見つかりません。")
  }
  return fixture
}

function toSummary(detail: TranslationJobManagementJobDetail): TranslationJobManagementJobSummary {
  return {
    jobId: detail.jobId,
    jobState: detail.jobState,
    jobStateLabel: detail.jobStateLabel,
    stateTone: detail.stateTone,
    canOpenPhase: detail.canOpenPhase,
    openBlockedReason: detail.openBlockedReason
      ? { ...detail.openBlockedReason }
      : undefined,
    inputSource: { ...detail.inputSource },
    progress: { ...detail.progress },
    stopAvailability: { ...detail.stopAvailability },
    resumeAvailability: { ...detail.resumeAvailability },
    deleteAvailability: { ...detail.deleteAvailability }
  }
}

function cloneDetail(detail: TranslationJobManagementJobDetail): TranslationJobManagementJobDetail {
  return {
    ...toSummary(detail),
    cacheState: detail.cacheState,
    cacheStateLabel: detail.cacheStateLabel,
    runtimeSummary: { ...detail.runtimeSummary },
    resumeBlockedReasons: detail.resumeBlockedReasons.map((reason) => ({
      ...reason
    })),
    warnings: detail.warnings.map((reason) => ({ ...reason })),
    deleteImpactLines: [...detail.deleteImpactLines]
  }
}
