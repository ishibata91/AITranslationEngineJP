package service

import (
	"context"
	"io"
	"strings"
)

type repositoryStub struct {
	listFunc                 func(context.Context, MasterDictionaryQuery) (MasterDictionaryListResult, error)
	getByIDFunc              func(context.Context, int64) (MasterDictionaryEntry, error)
	createFunc               func(context.Context, MasterDictionaryDraft) (MasterDictionaryEntry, error)
	updateFunc               func(context.Context, int64, MasterDictionaryDraft) (MasterDictionaryEntry, error)
	deleteFunc               func(context.Context, int64) error
	upsertBySourceAndRECFunc func(context.Context, MasterDictionaryImportRecord) (MasterDictionaryEntry, bool, error)
	listQueries              []MasterDictionaryQuery
	getByIDCalls             []int64
	createDrafts             []MasterDictionaryDraft
	updateCalls              []repositoryUpdateCall
	deleteCalls              []int64
	upsertRecords            []MasterDictionaryImportRecord
}

type repositoryUpdateCall struct {
	id    int64
	draft MasterDictionaryDraft
}

func (stub *repositoryStub) List(
	ctx context.Context,
	query MasterDictionaryQuery,
) (MasterDictionaryListResult, error) {
	stub.listQueries = append(stub.listQueries, query)
	if stub.listFunc == nil {
		return MasterDictionaryListResult{}, nil
	}
	return stub.listFunc(ctx, query)
}

func (stub *repositoryStub) GetByID(
	ctx context.Context,
	id int64,
) (MasterDictionaryEntry, error) {
	stub.getByIDCalls = append(stub.getByIDCalls, id)
	if stub.getByIDFunc == nil {
		return MasterDictionaryEntry{}, nil
	}
	return stub.getByIDFunc(ctx, id)
}

func (stub *repositoryStub) Create(
	ctx context.Context,
	draft MasterDictionaryDraft,
) (MasterDictionaryEntry, error) {
	stub.createDrafts = append(stub.createDrafts, draft)
	if stub.createFunc == nil {
		return MasterDictionaryEntry{}, nil
	}
	return stub.createFunc(ctx, draft)
}

func (stub *repositoryStub) Update(
	ctx context.Context,
	id int64,
	draft MasterDictionaryDraft,
) (MasterDictionaryEntry, error) {
	stub.updateCalls = append(stub.updateCalls, repositoryUpdateCall{id: id, draft: draft})
	if stub.updateFunc == nil {
		return MasterDictionaryEntry{}, nil
	}
	return stub.updateFunc(ctx, id, draft)
}

func (stub *repositoryStub) Delete(ctx context.Context, id int64) error {
	stub.deleteCalls = append(stub.deleteCalls, id)
	if stub.deleteFunc == nil {
		return nil
	}
	return stub.deleteFunc(ctx, id)
}

func (stub *repositoryStub) UpsertBySourceAndREC(
	ctx context.Context,
	record MasterDictionaryImportRecord,
) (MasterDictionaryEntry, bool, error) {
	stub.upsertRecords = append(stub.upsertRecords, record)
	if stub.upsertBySourceAndRECFunc == nil {
		return MasterDictionaryEntry{}, false, nil
	}
	return stub.upsertBySourceAndRECFunc(ctx, record)
}

type xmlFilePortStub struct {
	resolvePathFunc func(string) (string, error)
	openFunc        func(string) (io.ReadCloser, error)
	resolvedPaths   []string
	openedPaths     []string
}

func (stub *xmlFilePortStub) ResolvePath(rawPath string) (string, error) {
	stub.resolvedPaths = append(stub.resolvedPaths, rawPath)
	if stub.resolvePathFunc == nil {
		return rawPath, nil
	}
	return stub.resolvePathFunc(rawPath)
}

func (stub *xmlFilePortStub) Open(path string) (io.ReadCloser, error) {
	stub.openedPaths = append(stub.openedPaths, path)
	if stub.openFunc == nil {
		return io.NopCloser(strings.NewReader("")), nil
	}
	return stub.openFunc(path)
}

type xmlRecordReaderStub struct {
	countStringRecordsFunc func(io.Reader) (int, error)
	readStringRecordsFunc  func(io.Reader, func(xmlStringRecord) error) error
	countCalls             int
	readCalls              int
}

func (stub *xmlRecordReaderStub) CountStringRecords(reader io.Reader) (int, error) {
	stub.countCalls++
	if stub.countStringRecordsFunc == nil {
		return 0, nil
	}
	return stub.countStringRecordsFunc(reader)
}

func (stub *xmlRecordReaderStub) ReadStringRecords(
	reader io.Reader,
	handle func(xmlStringRecord) error,
) error {
	stub.readCalls++
	if stub.readStringRecordsFunc == nil {
		return nil
	}
	return stub.readStringRecordsFunc(reader, handle)
}

type importProgressRecorder struct {
	values []int
}

func (recorder *importProgressRecorder) EmitImportProgress(_ context.Context, progress int) {
	recorder.values = append(recorder.values, progress)
}
