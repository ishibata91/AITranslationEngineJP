package apitest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	wails "aitranslationenginejp/internal/controller/wails"
	"aitranslationenginejp/internal/service"
	"aitranslationenginejp/internal/usecase"
)

// fixedClock は golden が指定する固定時刻 2026-01-01T00:00:00Z を返す。
// MasterDictionaryCommandService と MasterDictionaryImportService の clock injection に使う。
var masterDictionaryBoundaryFixedClock = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedMasterDictionaryClock() time.Time {
	return masterDictionaryBoundaryFixedClock
}

// newMasterDictionaryBoundaryController は境界結合テスト用に bootstrap 済み controller を返す。
// SQLite は test ごとに隔離された一時 DB を使い、clock は golden と同じ固定値に揃える。
// repository から service への型変換は service.NewSQLiteMasterDictionaryRepositoryPort を経由する。
func newMasterDictionaryBoundaryController(t *testing.T) *wails.MasterDictionaryController {
	t.Helper()
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "boundary-test.sqlite3")

	// service.RepositoryPort への変換は adapter 経由で行う
	repoPort, err := service.NewSQLiteMasterDictionaryRepositoryPort(ctx, dbPath, nil)
	if err != nil {
		t.Fatalf("open master dictionary boundary test database: %v", err)
	}
	if closer := service.SQLiteMasterDictionaryRepositoryPortCloser(repoPort); closer != nil {
		t.Cleanup(func() { _ = closer(ctx) })
	}

	queryService := service.NewMasterDictionaryQueryService(repoPort)
	commandService := service.NewMasterDictionaryCommandService(repoPort, fixedMasterDictionaryClock)
	importService := service.NewMasterDictionaryImportService(
		repoPort,
		service.NewLocalMasterDictionaryXMLFilePort(),
		service.NewXMLDecoderMasterDictionaryRecordReader(),
		nil,
		fixedMasterDictionaryClock,
	)

	masterDictionaryUsecase := usecase.NewMasterDictionaryUsecase(queryService, commandService, importService, nil)
	return wails.NewMasterDictionaryController(masterDictionaryUsecase, nil)
}

// writeTempXML は t.TempDir() 配下に XML ファイルを作り、パスを返す。
// import 系の境界結合テストで使う入力 XML fixture を生成する。
func writeTempXML(t *testing.T, name string, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp xml %s: %v", name, err)
	}
	return path
}

// assertSummaryField は golden の expected フィールドと controller 応答の値が一致することを assert する。
// 失敗時に golden ファイル名、フィールド名、期待値、実際値を出力する。
func assertSummaryField(t *testing.T, goldenFile string, field string, want, got any) {
	t.Helper()
	if fmt.Sprintf("%v", want) != fmt.Sprintf("%v", got) {
		t.Errorf("BCT golden mismatch [%s] field %s: want %v, got %v", goldenFile, field, want, got)
	}
}

// --- BCT-MDC-001: ListMasterDictionaryEntries 正常系 ---

// TestBoundary_MasterDictionary_ListMasterDictionaryEntries_Normal は
// 1件以上のエントリが登録された状態で ListMasterDictionaryEntries を呼び出した時に、
// controller 応答の DTO field 値が golden の期待値と一致することを assert する。
// 証明観点: id の string 型、totalCount、page、pageSize、canonical field 値の semantic（BCT-MDC-001）。
func TestBoundary_MasterDictionary_ListMasterDictionaryEntries_Normal(t *testing.T) {
	// Arrange
	const goldenFile = "list_normal.golden.json"
	golden := MustLoadBoundaryGolden(goldenFile)

	var expectedResponse struct {
		Entries []struct {
			ID          string `json:"id"`
			Source      string `json:"source"`
			Translation string `json:"translation"`
			Category    string `json:"category"`
			Origin      string `json:"origin"`
			UpdatedAt   string `json:"updatedAt"`
		} `json:"entries"`
		TotalCount int `json:"totalCount"`
		Page       int `json:"page"`
		PageSize   int `json:"pageSize"`
	}
	if err := json.Unmarshal(golden.Expected, &expectedResponse); err != nil {
		t.Fatalf("unmarshal golden expected: %v", err)
	}

	controller := newMasterDictionaryBoundaryController(t)

	// Arrange: エントリを 1 件登録して golden の入力状態を再現する
	_, err := controller.CreateMasterDictionaryEntry(wails.CreateMasterDictionaryEntryRequestDTO{
		Payload: wails.MasterDictionaryEntryPayloadDTO{
			Source:      "Whiterun",
			Translation: "ホワイトラン",
			Category:    "地名",
			Origin:      "手動登録",
		},
	})
	if err != nil {
		t.Fatalf("setup: create entry: %v", err)
	}

	// Act
	resp, err := controller.ListMasterDictionaryEntries(wails.ListMasterDictionaryEntriesRequestDTO{
		Filters: wails.ListMasterDictionaryEntriesFiltersDTO{
			Query:    "",
			Category: "",
			Page:     1,
			PageSize: 30,
		},
	})
	if err != nil {
		t.Fatalf("ListMasterDictionaryEntries: %v", err)
	}

	// Assert: golden の expected field 値と controller 応答を突き合わせる
	assertSummaryField(t, goldenFile, "totalCount", expectedResponse.TotalCount, resp.TotalCount)
	assertSummaryField(t, goldenFile, "page", expectedResponse.Page, resp.Page)
	assertSummaryField(t, goldenFile, "pageSize", expectedResponse.PageSize, resp.PageSize)

	if len(resp.Entries) != len(expectedResponse.Entries) {
		t.Fatalf("BCT golden mismatch [%s] entries length: want %d, got %d", goldenFile, len(expectedResponse.Entries), len(resp.Entries))
	}
	if len(resp.Entries) > 0 {
		entry := resp.Entries[0]
		expected := expectedResponse.Entries[0]
		// id は string 型（int64 ではなく）で返ること
		assertSummaryField(t, goldenFile, "entries[0].id", expected.ID, entry.ID)
		assertSummaryField(t, goldenFile, "entries[0].source", expected.Source, entry.Source)
		assertSummaryField(t, goldenFile, "entries[0].translation", expected.Translation, entry.Translation)
		assertSummaryField(t, goldenFile, "entries[0].category", expected.Category, entry.Category)
		assertSummaryField(t, goldenFile, "entries[0].origin", expected.Origin, entry.Origin)
		assertSummaryField(t, goldenFile, "entries[0].updatedAt", expected.UpdatedAt, entry.UpdatedAt)
	}
}

// --- BCT-MDC-002: ListMasterDictionaryEntries 空集合 ---

// TestBoundary_MasterDictionary_ListMasterDictionaryEntries_Empty は
// エントリ 0 件の状態で ListMasterDictionaryEntries を呼び出した時に、
// entries が空配列・totalCount=0 であることを golden と突き合わせて assert する（BCT-MDC-002）。
func TestBoundary_MasterDictionary_ListMasterDictionaryEntries_Empty(t *testing.T) {
	// Arrange
	const goldenFile = "list_empty.golden.json"
	golden := MustLoadBoundaryGolden(goldenFile)

	var expectedResponse struct {
		Entries    []json.RawMessage `json:"entries"`
		TotalCount int               `json:"totalCount"`
		Page       int               `json:"page"`
		PageSize   int               `json:"pageSize"`
	}
	if err := json.Unmarshal(golden.Expected, &expectedResponse); err != nil {
		t.Fatalf("unmarshal golden expected: %v", err)
	}

	controller := newMasterDictionaryBoundaryController(t)
	// DB は空のまま（エントリ登録なし）

	// Act
	resp, err := controller.ListMasterDictionaryEntries(wails.ListMasterDictionaryEntriesRequestDTO{
		Filters: wails.ListMasterDictionaryEntriesFiltersDTO{
			Query:    "",
			Category: "",
			Page:     1,
			PageSize: 30,
		},
	})
	if err != nil {
		t.Fatalf("ListMasterDictionaryEntries: %v", err)
	}

	// Assert: 空配列と totalCount=0 の semantic を証明する
	assertSummaryField(t, goldenFile, "entries length", len(expectedResponse.Entries), len(resp.Entries))
	assertSummaryField(t, goldenFile, "totalCount", expectedResponse.TotalCount, resp.TotalCount)
	assertSummaryField(t, goldenFile, "page", expectedResponse.Page, resp.Page)
	assertSummaryField(t, goldenFile, "pageSize", expectedResponse.PageSize, resp.PageSize)
}

// --- BCT-MDC-003: GetMasterDictionaryEntry 正常系 ---

// TestBoundary_MasterDictionary_GetMasterDictionaryEntry_Normal は
// 登録済みエントリの ID を指定して GetMasterDictionaryEntry を呼び出した時に、
// entry フィールド（note の固定文言を含む）が golden の期待値と一致することを assert する（BCT-MDC-003）。
func TestBoundary_MasterDictionary_GetMasterDictionaryEntry_Normal(t *testing.T) {
	// Arrange: get_after_create golden を参照（GetMasterDictionaryEntry 正常系の代表）
	const goldenFile = "get_after_create.golden.json"
	golden := MustLoadBoundaryGolden(goldenFile)

	var expectedResponse struct {
		Entry struct {
			ID          string `json:"id"`
			Source      string `json:"source"`
			Translation string `json:"translation"`
			Category    string `json:"category"`
			Origin      string `json:"origin"`
			UpdatedAt   string `json:"updatedAt"`
			Note        string `json:"note"`
		} `json:"entry"`
	}
	if err := json.Unmarshal(golden.Expected, &expectedResponse); err != nil {
		t.Fatalf("unmarshal golden expected: %v", err)
	}

	controller := newMasterDictionaryBoundaryController(t)

	created, err := controller.CreateMasterDictionaryEntry(wails.CreateMasterDictionaryEntryRequestDTO{
		Payload: wails.MasterDictionaryEntryPayloadDTO{
			Source:      "Whiterun",
			Translation: "ホワイトラン",
			Category:    "地名",
			Origin:      "手動登録",
		},
	})
	if err != nil {
		t.Fatalf("setup: create entry: %v", err)
	}

	// Act
	resp, err := controller.GetMasterDictionaryEntry(wails.GetMasterDictionaryEntryRequestDTO{
		ID: created.Entry.ID,
	})
	if err != nil {
		t.Fatalf("GetMasterDictionaryEntry: %v", err)
	}

	// Assert: entry が non-null、note 固定文言、canonical field の semantic を証明する
	if resp.Entry == nil {
		t.Fatalf("BCT golden mismatch [%s] entry: want non-null, got nil", goldenFile)
	}
	assertSummaryField(t, goldenFile, "entry.id", expectedResponse.Entry.ID, resp.Entry.ID)
	assertSummaryField(t, goldenFile, "entry.source", expectedResponse.Entry.Source, resp.Entry.Source)
	assertSummaryField(t, goldenFile, "entry.translation", expectedResponse.Entry.Translation, resp.Entry.Translation)
	assertSummaryField(t, goldenFile, "entry.category", expectedResponse.Entry.Category, resp.Entry.Category)
	assertSummaryField(t, goldenFile, "entry.origin", expectedResponse.Entry.Origin, resp.Entry.Origin)
	assertSummaryField(t, goldenFile, "entry.updatedAt", expectedResponse.Entry.UpdatedAt, resp.Entry.UpdatedAt)
	// note は "マスター辞書エントリ" の固定文言が backend から返ることを証明する（BCT-MDC-003）
	assertSummaryField(t, goldenFile, "entry.note", expectedResponse.Entry.Note, resp.Entry.Note)
}

// --- BCT-MDC-004: GetMasterDictionaryEntry 不在（境界） ---

// TestBoundary_MasterDictionary_GetMasterDictionaryEntry_Absent は
// 存在しない ID を指定して GetMasterDictionaryEntry を呼び出した時に、
// entry が null を返すことを golden と突き合わせて assert する（BCT-MDC-004）。
// backend が実際に null を返すことを bootstrap 済み controller で証明する。
func TestBoundary_MasterDictionary_GetMasterDictionaryEntry_Absent(t *testing.T) {
	// Arrange
	const goldenFile = "get_absent.golden.json"
	golden := MustLoadBoundaryGolden(goldenFile)

	var expectedResponse struct {
		Entry *json.RawMessage `json:"entry"`
	}
	if err := json.Unmarshal(golden.Expected, &expectedResponse); err != nil {
		t.Fatalf("unmarshal golden expected: %v", err)
	}

	controller := newMasterDictionaryBoundaryController(t)
	// DB は空のまま（不在 ID への参照を確認する）

	// Act: golden が示す不在 ID "999999" を使う
	resp, err := controller.GetMasterDictionaryEntry(wails.GetMasterDictionaryEntryRequestDTO{
		ID: "999999",
	})
	if err != nil {
		t.Fatalf("GetMasterDictionaryEntry with absent id: %v", err)
	}

	// Assert: entry が null であることを証明する（エラーにならない）
	goldenExpectsNull := expectedResponse.Entry == nil || string(*expectedResponse.Entry) == "null"
	if !goldenExpectsNull {
		t.Fatalf("BCT golden [%s] expects non-null entry, test setup error", goldenFile)
	}
	if resp.Entry != nil {
		t.Errorf("BCT golden mismatch [%s] entry: want null, got %+v", goldenFile, resp.Entry)
	}
}

// --- BCT-MDC-005: CreateMasterDictionaryEntry 正常系 ---

// TestBoundary_MasterDictionary_CreateMasterDictionaryEntry_Normal は
// CreateMasterDictionaryEntry を呼び出した時に、
// entry フィールドと refreshTargetId が golden の期待値と一致することを assert する（BCT-MDC-005）。
func TestBoundary_MasterDictionary_CreateMasterDictionaryEntry_Normal(t *testing.T) {
	// Arrange
	const goldenFile = "create_response.golden.json"
	golden := MustLoadBoundaryGolden(goldenFile)

	var expectedResponse struct {
		Entry struct {
			ID          string `json:"id"`
			Source      string `json:"source"`
			Translation string `json:"translation"`
			Category    string `json:"category"`
			Origin      string `json:"origin"`
			UpdatedAt   string `json:"updatedAt"`
			Note        string `json:"note"`
		} `json:"entry"`
		RefreshTargetID string `json:"refreshTargetId"`
	}
	if err := json.Unmarshal(golden.Expected, &expectedResponse); err != nil {
		t.Fatalf("unmarshal golden expected: %v", err)
	}

	controller := newMasterDictionaryBoundaryController(t)

	// Act
	resp, err := controller.CreateMasterDictionaryEntry(wails.CreateMasterDictionaryEntryRequestDTO{
		Payload: wails.MasterDictionaryEntryPayloadDTO{
			Source:      "Whiterun",
			Translation: "ホワイトラン",
			Category:    "地名",
			Origin:      "手動登録",
		},
	})
	if err != nil {
		t.Fatalf("CreateMasterDictionaryEntry: %v", err)
	}

	// Assert: 作成直後の DTO field semantic を証明する
	assertSummaryField(t, goldenFile, "entry.id", expectedResponse.Entry.ID, resp.Entry.ID)
	assertSummaryField(t, goldenFile, "entry.source", expectedResponse.Entry.Source, resp.Entry.Source)
	assertSummaryField(t, goldenFile, "entry.translation", expectedResponse.Entry.Translation, resp.Entry.Translation)
	assertSummaryField(t, goldenFile, "entry.category", expectedResponse.Entry.Category, resp.Entry.Category)
	assertSummaryField(t, goldenFile, "entry.origin", expectedResponse.Entry.Origin, resp.Entry.Origin)
	assertSummaryField(t, goldenFile, "entry.updatedAt", expectedResponse.Entry.UpdatedAt, resp.Entry.UpdatedAt)
	assertSummaryField(t, goldenFile, "entry.note", expectedResponse.Entry.Note, resp.Entry.Note)
	// refreshTargetId は entry.id と同じ値であることを証明する
	assertSummaryField(t, goldenFile, "refreshTargetId", expectedResponse.RefreshTargetID, resp.RefreshTargetID)
}

// --- BCT-MDC-006: Create → Get 状態遷移連鎖 ---

// TestBoundary_MasterDictionary_CreateThenGet_Chain は
// CreateMasterDictionaryEntry で作成したエントリの ID を使って GetMasterDictionaryEntry を呼び出した時に、
// Create 応答と Get 応答の entry.id・entry.source が一致することを golden と突き合わせて assert する（BCT-MDC-006）。
// 作成 → 取得の field 値 semantic 整合を証明する。
func TestBoundary_MasterDictionary_CreateThenGet_Chain(t *testing.T) {
	// Arrange: create_response.golden.json と get_after_create.golden.json を両方ロード
	const createGolden = "create_response.golden.json"
	const getGolden = "get_after_create.golden.json"

	createG := MustLoadBoundaryGolden(createGolden)
	getG := MustLoadBoundaryGolden(getGolden)

	var expectedCreate struct {
		Entry struct {
			ID     string `json:"id"`
			Source string `json:"source"`
		} `json:"entry"`
	}
	if err := json.Unmarshal(createG.Expected, &expectedCreate); err != nil {
		t.Fatalf("unmarshal create golden: %v", err)
	}

	var expectedGet struct {
		Entry struct {
			ID     string `json:"id"`
			Source string `json:"source"`
			Note   string `json:"note"`
		} `json:"entry"`
	}
	if err := json.Unmarshal(getG.Expected, &expectedGet); err != nil {
		t.Fatalf("unmarshal get golden: %v", err)
	}

	controller := newMasterDictionaryBoundaryController(t)

	// Act: Create
	createResp, err := controller.CreateMasterDictionaryEntry(wails.CreateMasterDictionaryEntryRequestDTO{
		Payload: wails.MasterDictionaryEntryPayloadDTO{
			Source:      "Whiterun",
			Translation: "ホワイトラン",
			Category:    "地名",
			Origin:      "手動登録",
		},
	})
	if err != nil {
		t.Fatalf("CreateMasterDictionaryEntry: %v", err)
	}

	// Act: Get（Create で発行された ID を使う）
	getResp, err := controller.GetMasterDictionaryEntry(wails.GetMasterDictionaryEntryRequestDTO{
		ID: createResp.Entry.ID,
	})
	if err != nil {
		t.Fatalf("GetMasterDictionaryEntry: %v", err)
	}

	// Assert: Create 応答 field と golden（create_response）を突き合わせる
	assertSummaryField(t, createGolden, "create entry.id", expectedCreate.Entry.ID, createResp.Entry.ID)
	assertSummaryField(t, createGolden, "create entry.source", expectedCreate.Entry.Source, createResp.Entry.Source)

	// Assert: Get 応答 field と golden（get_after_create）を突き合わせる
	if getResp.Entry == nil {
		t.Fatalf("BCT golden mismatch [%s] entry: want non-null after create, got nil", getGolden)
	}
	assertSummaryField(t, getGolden, "get entry.id", expectedGet.Entry.ID, getResp.Entry.ID)
	assertSummaryField(t, getGolden, "get entry.source", expectedGet.Entry.Source, getResp.Entry.Source)
	assertSummaryField(t, getGolden, "get entry.note", expectedGet.Entry.Note, getResp.Entry.Note)

	// Assert: Create と Get の entry.id・entry.source が semantic 整合することを証明する
	assertSummaryField(t, createGolden+"+"+getGolden, "create.entry.id == get.entry.id", createResp.Entry.ID, getResp.Entry.ID)
	assertSummaryField(t, createGolden+"+"+getGolden, "create.entry.source == get.entry.source", createResp.Entry.Source, getResp.Entry.Source)
}

// --- BCT-MDC-007: Update → Get 状態遷移連鎖 ---

// TestBoundary_MasterDictionary_UpdateThenGet_Chain は
// UpdateMasterDictionaryEntry で translation を変更した後に GetMasterDictionaryEntry を呼び出した時に、
// Update 応答と Get 応答の translation が新しい値と一致し、source が変化しないことを
// golden と突き合わせて assert する（BCT-MDC-007）。
func TestBoundary_MasterDictionary_UpdateThenGet_Chain(t *testing.T) {
	// Arrange: update_response.golden.json と get_after_update.golden.json を両方ロード
	const updateGolden = "update_response.golden.json"
	const getGolden = "get_after_update.golden.json"

	updateG := MustLoadBoundaryGolden(updateGolden)
	getG := MustLoadBoundaryGolden(getGolden)

	var expectedUpdate struct {
		Entry struct {
			ID          string `json:"id"`
			Source      string `json:"source"`
			Translation string `json:"translation"`
			UpdatedAt   string `json:"updatedAt"`
		} `json:"entry"`
		RefreshTargetID string `json:"refreshTargetId"`
	}
	if err := json.Unmarshal(updateG.Expected, &expectedUpdate); err != nil {
		t.Fatalf("unmarshal update golden: %v", err)
	}

	var expectedGet struct {
		Entry struct {
			ID          string `json:"id"`
			Source      string `json:"source"`
			Translation string `json:"translation"`
		} `json:"entry"`
	}
	if err := json.Unmarshal(getG.Expected, &expectedGet); err != nil {
		t.Fatalf("unmarshal get golden: %v", err)
	}

	controller := newMasterDictionaryBoundaryController(t)

	// Arrange: エントリを 1 件登録する
	created, err := controller.CreateMasterDictionaryEntry(wails.CreateMasterDictionaryEntryRequestDTO{
		Payload: wails.MasterDictionaryEntryPayloadDTO{
			Source:      "Whiterun",
			Translation: "ホワイトラン",
			Category:    "地名",
			Origin:      "手動登録",
		},
	})
	if err != nil {
		t.Fatalf("setup: create entry: %v", err)
	}

	// Act: Update（translation だけ変更）
	updateResp, err := controller.UpdateMasterDictionaryEntry(wails.UpdateMasterDictionaryEntryRequestDTO{
		ID: created.Entry.ID,
		Payload: wails.MasterDictionaryEntryPayloadDTO{
			Source:      "Whiterun",
			Translation: "白の城下町",
			Category:    "地名",
			Origin:      "手動登録",
		},
	})
	if err != nil {
		t.Fatalf("UpdateMasterDictionaryEntry: %v", err)
	}

	// Act: Get（更新後の同一 ID）
	getResp, err := controller.GetMasterDictionaryEntry(wails.GetMasterDictionaryEntryRequestDTO{
		ID: created.Entry.ID,
	})
	if err != nil {
		t.Fatalf("GetMasterDictionaryEntry: %v", err)
	}

	// Assert: Update 応答と golden（update_response）を突き合わせる
	assertSummaryField(t, updateGolden, "update entry.id", expectedUpdate.Entry.ID, updateResp.Entry.ID)
	assertSummaryField(t, updateGolden, "update entry.source", expectedUpdate.Entry.Source, updateResp.Entry.Source)
	assertSummaryField(t, updateGolden, "update entry.translation", expectedUpdate.Entry.Translation, updateResp.Entry.Translation)
	assertSummaryField(t, updateGolden, "update entry.updatedAt", expectedUpdate.Entry.UpdatedAt, updateResp.Entry.UpdatedAt)
	assertSummaryField(t, updateGolden, "update refreshTargetId", expectedUpdate.RefreshTargetID, updateResp.RefreshTargetID)

	// Assert: Get 応答と golden（get_after_update）を突き合わせる
	if getResp.Entry == nil {
		t.Fatalf("BCT golden mismatch [%s] entry: want non-null after update, got nil", getGolden)
	}
	assertSummaryField(t, getGolden, "get entry.translation", expectedGet.Entry.Translation, getResp.Entry.Translation)
	assertSummaryField(t, getGolden, "get entry.source", expectedGet.Entry.Source, getResp.Entry.Source)

	// Assert: 更新前後で source が変化しない、translation が新しい値に変化することを証明する
	assertSummaryField(t, updateGolden+"+"+getGolden, "update.source == get.source（不変）", updateResp.Entry.Source, getResp.Entry.Source)
	assertSummaryField(t, updateGolden+"+"+getGolden, "update.translation == get.translation（変化後の値）", updateResp.Entry.Translation, getResp.Entry.Translation)
}

// --- BCT-MDC-008: Delete → Get 状態遷移連鎖 ---

// TestBoundary_MasterDictionary_DeleteThenGet_Chain は
// DeleteMasterDictionaryEntry で削除した後に GetMasterDictionaryEntry を呼び出した時に、
// deletedId が削除した ID と一致し、削除後の Get が entry=null を返すことを
// golden と突き合わせて assert する（BCT-MDC-008）。
func TestBoundary_MasterDictionary_DeleteThenGet_Chain(t *testing.T) {
	// Arrange: delete_response.golden.json と get_after_delete.golden.json を両方ロード
	const deleteGolden = "delete_response.golden.json"
	const getGolden = "get_after_delete.golden.json"

	deleteG := MustLoadBoundaryGolden(deleteGolden)
	getG := MustLoadBoundaryGolden(getGolden)

	var expectedDelete struct {
		DeletedID      string  `json:"deletedId"`
		NextSelectedID *string `json:"nextSelectedId"`
	}
	if err := json.Unmarshal(deleteG.Expected, &expectedDelete); err != nil {
		t.Fatalf("unmarshal delete golden: %v", err)
	}

	var expectedGet struct {
		Entry *json.RawMessage `json:"entry"`
	}
	if err := json.Unmarshal(getG.Expected, &expectedGet); err != nil {
		t.Fatalf("unmarshal get-after-delete golden: %v", err)
	}

	controller := newMasterDictionaryBoundaryController(t)

	// Arrange: エントリを 1 件登録する
	created, err := controller.CreateMasterDictionaryEntry(wails.CreateMasterDictionaryEntryRequestDTO{
		Payload: wails.MasterDictionaryEntryPayloadDTO{
			Source:      "Whiterun",
			Translation: "ホワイトラン",
			Category:    "地名",
			Origin:      "手動登録",
		},
	})
	if err != nil {
		t.Fatalf("setup: create entry: %v", err)
	}
	entryID := created.Entry.ID

	// Act: Delete
	deleteResp, err := controller.DeleteMasterDictionaryEntry(wails.DeleteMasterDictionaryEntryRequestDTO{
		ID: entryID,
	})
	if err != nil {
		t.Fatalf("DeleteMasterDictionaryEntry: %v", err)
	}

	// Act: Get（削除後の同一 ID）
	getResp, err := controller.GetMasterDictionaryEntry(wails.GetMasterDictionaryEntryRequestDTO{
		ID: entryID,
	})
	if err != nil {
		t.Fatalf("GetMasterDictionaryEntry after delete: %v", err)
	}

	// Assert: Delete 応答と golden（delete_response）を突き合わせる
	// deletedId は削除した ID と同じ string 値であることを証明する
	assertSummaryField(t, deleteGolden, "deletedId", expectedDelete.DeletedID, deleteResp.DeletedID)

	// nextSelectedId: golden は null（1 件削除後 0 件のため）
	goldenNextIsNull := expectedDelete.NextSelectedID == nil
	gotNextIsNull := deleteResp.NextSelectedID == nil
	if goldenNextIsNull != gotNextIsNull {
		t.Errorf("BCT golden mismatch [%s] nextSelectedId null: want %v, got %v", deleteGolden, goldenNextIsNull, gotNextIsNull)
	}

	// Assert: 削除後の Get で entry=null となる状態遷移 semantic を証明する
	goldenExpectsNull := expectedGet.Entry == nil || string(*expectedGet.Entry) == "null"
	if !goldenExpectsNull {
		t.Fatalf("BCT golden [%s] expects non-null entry after delete, test setup error", getGolden)
	}
	if getResp.Entry != nil {
		t.Errorf("BCT golden mismatch [%s] entry after delete: want null, got %+v", getGolden, getResp.Entry)
	}
}

// --- BCT-MDC-009: ImportMasterDictionaryXml 13種別内 REC（正常） ---

// TestBoundary_MasterDictionary_ImportMasterDictionaryXml_RecInScope は
// 13 種別内 REC（NPC_:FULL）のみを含む XML を取り込んだ時に、
// importedCount>0・skippedCount=0 を golden と突き合わせて assert する（BCT-MDC-009）。
// IsTermTarget の判定ロジック自体ではなく、応答 DTO の skippedCount semantic を証明する。
func TestBoundary_MasterDictionary_ImportMasterDictionaryXml_RecInScope(t *testing.T) {
	// Arrange
	const goldenFile = "import_rec_in_scope.golden.json"
	golden := MustLoadBoundaryGolden(goldenFile)

	var expectedResponse struct {
		Accepted bool `json:"accepted"`
		Summary  struct {
			ImportedCount int `json:"importedCount"`
			UpdatedCount  int `json:"updatedCount"`
			SkippedCount  int `json:"skippedCount"`
		} `json:"summary"`
	}
	if err := json.Unmarshal(golden.Expected, &expectedResponse); err != nil {
		t.Fatalf("unmarshal golden expected: %v", err)
	}

	// 13 種別内 REC（NPC_:FULL）を含む XML fixture を生成する
	inScopeXMLContent := strings.Join([]string{
		`<?xml version="1.0" encoding="utf-8"?>`,
		`<SSETranslator>`,
		`  <String>`,
		`    <EDID>Actor_Whiterun</EDID>`,
		`    <REC>NPC_:FULL</REC>`,
		`    <Source>Whiterun Guard</Source>`,
		`    <Dest>ホワイトランの衛兵</Dest>`,
		`  </String>`,
		`</SSETranslator>`,
	}, "\n")
	xmlPath := writeTempXML(t, "test_in_scope.xml", inScopeXMLContent)

	controller := newMasterDictionaryBoundaryController(t)

	// Act
	resp, err := controller.ImportMasterDictionaryXml(wails.ImportMasterDictionaryXMLRequestDTO{
		FilePath: xmlPath,
	})
	if err != nil {
		t.Fatalf("ImportMasterDictionaryXml (in-scope): %v", err)
	}

	// Assert: accepted=true、importedCount>0、skippedCount=0 の semantic を証明する
	assertSummaryField(t, goldenFile, "accepted", expectedResponse.Accepted, resp.Accepted)
	if resp.Summary == nil {
		t.Fatalf("BCT golden mismatch [%s] summary: want non-null, got nil", goldenFile)
	}
	// golden が importedCount=1 を指定しているため、取り込み件数が 1 件以上であることを確認する
	if resp.Summary.ImportedCount < expectedResponse.Summary.ImportedCount {
		t.Errorf("BCT golden mismatch [%s] summary.importedCount: want >= %d, got %d",
			goldenFile, expectedResponse.Summary.ImportedCount, resp.Summary.ImportedCount)
	}
	assertSummaryField(t, goldenFile, "summary.updatedCount", expectedResponse.Summary.UpdatedCount, resp.Summary.UpdatedCount)
	assertSummaryField(t, goldenFile, "summary.skippedCount", expectedResponse.Summary.SkippedCount, resp.Summary.SkippedCount)
}

// --- BCT-MDC-010: ImportMasterDictionaryXml 13種別外 REC（境界） ---

// TestBoundary_MasterDictionary_ImportMasterDictionaryXml_RecOutOfScope は
// 13 種別外 REC（DOOR:FULL）のみを含む XML を取り込んだ時に、
// importedCount=0・skippedCount>0 を golden と突き合わせて assert する（BCT-MDC-010）。
// 応答 DTO の skippedCount フィールドに件数が乗ることの semantic を証明する（Q-007 確定）。
func TestBoundary_MasterDictionary_ImportMasterDictionaryXml_RecOutOfScope(t *testing.T) {
	// Arrange
	const goldenFile = "import_rec_out_of_scope.golden.json"
	golden := MustLoadBoundaryGolden(goldenFile)

	var expectedResponse struct {
		Accepted bool `json:"accepted"`
		Summary  struct {
			ImportedCount int `json:"importedCount"`
			UpdatedCount  int `json:"updatedCount"`
			SkippedCount  int `json:"skippedCount"`
		} `json:"summary"`
	}
	if err := json.Unmarshal(golden.Expected, &expectedResponse); err != nil {
		t.Fatalf("unmarshal golden expected: %v", err)
	}

	// 13 種別外 REC（DOOR:FULL）のみを含む XML fixture を生成する
	outOfScopeXMLContent := strings.Join([]string{
		`<?xml version="1.0" encoding="utf-8"?>`,
		`<SSETranslator>`,
		`  <String>`,
		`    <EDID>Door_Whiterun_Main</EDID>`,
		`    <REC>DOOR:FULL</REC>`,
		`    <Source>Whiterun City Gate</Source>`,
		`    <Dest>ホワイトランの城門</Dest>`,
		`  </String>`,
		`</SSETranslator>`,
	}, "\n")
	xmlPath := writeTempXML(t, "test_out_of_scope.xml", outOfScopeXMLContent)

	controller := newMasterDictionaryBoundaryController(t)

	// Act
	resp, err := controller.ImportMasterDictionaryXml(wails.ImportMasterDictionaryXMLRequestDTO{
		FilePath: xmlPath,
	})
	if err != nil {
		t.Fatalf("ImportMasterDictionaryXml (out-of-scope): %v", err)
	}

	// Assert: accepted=true、importedCount=0、skippedCount>0 の semantic を証明する
	assertSummaryField(t, goldenFile, "accepted", expectedResponse.Accepted, resp.Accepted)
	if resp.Summary == nil {
		t.Fatalf("BCT golden mismatch [%s] summary: want non-null, got nil", goldenFile)
	}
	assertSummaryField(t, goldenFile, "summary.importedCount", expectedResponse.Summary.ImportedCount, resp.Summary.ImportedCount)
	assertSummaryField(t, goldenFile, "summary.updatedCount", expectedResponse.Summary.UpdatedCount, resp.Summary.UpdatedCount)
	// golden が skippedCount=1 を指定しているため、skippedCount が 1 件以上であることを確認する
	if resp.Summary.SkippedCount < expectedResponse.Summary.SkippedCount {
		t.Errorf("BCT golden mismatch [%s] summary.skippedCount: want >= %d, got %d",
			goldenFile, expectedResponse.Summary.SkippedCount, resp.Summary.SkippedCount)
	}
}
