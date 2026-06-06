// 境界 API contract: MasterDictionary（pilot 6 method）
//
// 本 file は frontend と backend の境界 API「形式」の真を定義する。
// 形式とは、型、必須項目、null 許容、値域のことを指す。
//
// 本 file が持たない情報:
//   - 意味（field が業務的に何を指すか） … 画面仕様 / UC（usecase docs）が持つ
//   - 状態遷移規約 … UC が持つ
//   - 表示文言、UI 構造 … 画面設計が持つ
//
// 本 contract は frontend test と backend test の双方が import して契約検証に使う。
// Wails 生成 TypeScript binding（`frontend/wailsjs/go/...`）は実装側であり、本 contract に追従する。
//
// status: spec-as-should-be（あるべき仕様）
// 配置先: 正本反映時に `frontend/src/controller/wails/master-dictionary.contract.ts` に移す候補
// 作成日: 2026-06-06
//
// あるべき仕様で外した field:
//   - note（メモ）: 現実装は Detail に固定文言を入れる。境界 API には不要
//   - origin（由来 / 登録元）: 現実装は Summary / Detail / payload に持つ。境界 API には不要

/* ---------------------------------------------------------------------------
 * 共通型
 * ------------------------------------------------------------------------- */

/** ID は backend Go 内部で int64。境界に出る時は string に変換する。 */
export type EntryId = string;

/** RFC3339 形式の時刻文字列。UTC 前提。例: "2026-06-06T00:00:00Z" */
export type Rfc3339Timestamp = string;

/** ページ番号。1 以上の整数。 */
export type PageNumber = number;

/** 1 ページあたりの件数。1 以上の整数。default は 30。 */
export type PageSize = number;

/** 0 以上の整数件数。 */
export type NonNegativeCount = number;

/** リフレッシュ用ページ条件（optional）。 */
export interface MasterDictionaryFrontendRefresh {
  readonly page: PageNumber;
  readonly pageSize: PageSize;
}

/** ページ状態。リフレッシュ後の応答に含まれる。 */
export interface MasterDictionaryPageState {
  readonly page: PageNumber;
  readonly pageSize: PageSize;
  readonly totalCount: NonNegativeCount;
}

/* ---------------------------------------------------------------------------
 * Entry: Summary と Detail は同じ形式とする（あるべき仕様）
 * ------------------------------------------------------------------------- */

/** マスター辞書エントリ。一覧と詳細で同じ形式を使う。 */
export interface MasterDictionaryEntry {
  readonly id: EntryId;
  readonly source: string;
  readonly translation: string;
  readonly category: string;
  readonly updatedAt: Rfc3339Timestamp;
}

/** Create / Update payload。 */
export interface MasterDictionaryEntryPayload {
  readonly source: string;
  readonly translation: string;
  readonly category: string;
}

/* ---------------------------------------------------------------------------
 * 1. ListMasterDictionaryEntries
 * ------------------------------------------------------------------------- */

export interface ListMasterDictionaryEntriesRequest {
  readonly filters: {
    readonly query: string;
    readonly category: string;
    readonly page: PageNumber;
    readonly pageSize: PageSize;
  };
}

export interface ListMasterDictionaryEntriesResponse {
  readonly entries: ReadonlyArray<MasterDictionaryEntry>;
  readonly totalCount: NonNegativeCount;
  readonly page: PageNumber;
  readonly pageSize: PageSize;
}

/* ---------------------------------------------------------------------------
 * 2. GetMasterDictionaryEntry
 * ------------------------------------------------------------------------- */

export interface GetMasterDictionaryEntryRequest {
  readonly id: EntryId;
}

/** entry は不在時 null を返す（エラーにならない）。 */
export interface GetMasterDictionaryEntryResponse {
  readonly entry: MasterDictionaryEntry | null;
}

/* ---------------------------------------------------------------------------
 * 3. CreateMasterDictionaryEntry
 * ------------------------------------------------------------------------- */

export interface CreateMasterDictionaryEntryRequest {
  readonly payload: MasterDictionaryEntryPayload;
  readonly refresh?: MasterDictionaryFrontendRefresh;
}

export interface CreateMasterDictionaryEntryResponse {
  readonly entry: MasterDictionaryEntry;
  readonly refreshTargetId: EntryId;
  readonly page?: MasterDictionaryPageState;
}

/* ---------------------------------------------------------------------------
 * 4. UpdateMasterDictionaryEntry
 * ------------------------------------------------------------------------- */

export interface UpdateMasterDictionaryEntryRequest {
  readonly id: EntryId;
  readonly payload: MasterDictionaryEntryPayload;
  readonly refresh?: MasterDictionaryFrontendRefresh;
}

export interface UpdateMasterDictionaryEntryResponse {
  readonly entry: MasterDictionaryEntry;
  readonly refreshTargetId: EntryId;
  readonly page?: MasterDictionaryPageState;
}

/* ---------------------------------------------------------------------------
 * 5. DeleteMasterDictionaryEntry
 * ------------------------------------------------------------------------- */

export interface DeleteMasterDictionaryEntryRequest {
  readonly id: EntryId;
  readonly refresh?: MasterDictionaryFrontendRefresh;
}

/** nextSelectedId は次選択候補が無い場合 null になる。 */
export interface DeleteMasterDictionaryEntryResponse {
  readonly deletedId: EntryId;
  readonly nextSelectedId: EntryId | null;
  readonly page?: MasterDictionaryPageState;
}

/* ---------------------------------------------------------------------------
 * 6. ImportMasterDictionaryXml
 * ------------------------------------------------------------------------- */

export interface ImportMasterDictionaryXmlRequest {
  readonly filePath: string;
  readonly fileReference?: string;
  readonly refresh?: MasterDictionaryFrontendRefresh;
}

/** accepted は常に true（成功時）。失敗時は Promise.reject になる。 */
export interface ImportMasterDictionaryXmlResponse {
  readonly accepted: true;
  readonly summary?: MasterDictionaryImportSummary;
  readonly page?: MasterDictionaryPageState;
}

export interface MasterDictionaryImportSummary {
  readonly filePath: string;
  readonly fileName: string;
  readonly importedCount: NonNegativeCount;
  readonly updatedCount: NonNegativeCount;
  readonly skippedCount: NonNegativeCount;
  readonly lastEntryId: number;
}

/* ---------------------------------------------------------------------------
 * Wails binding 名（method 名と一致）
 * ------------------------------------------------------------------------- */

export const MASTER_DICTIONARY_BINDING_NAMES = {
  list: "ListMasterDictionaryEntries",
  get: "GetMasterDictionaryEntry",
  create: "CreateMasterDictionaryEntry",
  update: "UpdateMasterDictionaryEntry",
  delete: "DeleteMasterDictionaryEntry",
  importXml: "ImportMasterDictionaryXml",
} as const;

/* ---------------------------------------------------------------------------
 * エラー応答
 * ------------------------------------------------------------------------- */

// エラーは Wails bridge 経由で Promise.reject 相当として frontend に伝わる。
// error.Error() 文字列が message として渡る。frontend は message のパースに依存しない。
// 構造化エラー object は本 pilot 範囲では返さない。
//
// エラー種別ごとの観測形式:
//   - ID 形式不正        : Promise.reject (message: `parse id "<value>": ...`)
//   - ID 必須            : Promise.reject (message: `id is required`)
//   - ID 正整数違反      : Promise.reject (message: `id must be greater than zero`)
//   - エントリ不在       : GetMasterDictionaryEntry のみ特別扱いで { entry: null } を返す
//   - usecase / repository 失敗 : Promise.reject (message に処理名を含む)
