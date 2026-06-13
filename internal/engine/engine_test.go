package engine

import (
	"context"
	"errors"
	"testing"

	"aitranslationenginejp/internal/model"
	"aitranslationenginejp/internal/provider"
)

type fakeStore struct {
	untranslated []model.Narration
	updates      []update
}

type update struct {
	id     int64
	dest   string
	status int
}

func (f *fakeStore) ListUntranslatedNarrations(_ context.Context) ([]model.Narration, error) {
	return f.untranslated, nil
}

func (f *fakeStore) UpdateNarrationDest(_ context.Context, id int64, dest string, status int) error {
	f.updates = append(f.updates, update{id: id, dest: dest, status: status})
	return nil
}

type fakeTranslator struct {
	out      map[string]string
	gotModel string
	err      error
}

func (f *fakeTranslator) Translate(_ context.Context, _ provider.Connection, model, source string) (string, error) {
	f.gotModel = model
	if f.err != nil {
		return "", f.err
	}
	return f.out[source], nil
}

func (f *fakeTranslator) ListModels(_ context.Context, _ provider.Connection) ([]string, error) {
	return nil, nil
}

// 未訳の叙述文を順に翻訳し、訳文を status=3（仮訳）で書き戻すこと。
func TestRunTranslatesUntranslatedAsProvisional(t *testing.T) {
	store := &fakeStore{untranslated: []model.Narration{
		{ID: 1, Source: "halls"},
		{ID: 2, Source: "cairn"},
	}}
	tr := &fakeTranslator{out: map[string]string{"halls": "広間", "cairn": "ケルン"}}
	eng := New(store, tr)

	count, err := eng.Run(context.Background(), provider.Connection{Endpoint: "http://x"}, "model-x")
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
	if tr.gotModel != "model-x" {
		t.Errorf("model passed = %q, want model-x", tr.gotModel)
	}
	want := []update{{1, "広間", 3}, {2, "ケルン", 3}}
	if len(store.updates) != len(want) {
		t.Fatalf("updates = %v, want %v", store.updates, want)
	}
	for i := range want {
		if store.updates[i] != want[i] {
			t.Errorf("update[%d] = %v, want %v", i, store.updates[i], want[i])
		}
	}
}

// provider がエラーを返したら、そのエラーを返すこと（壊れた訳を書き戻さない）。
func TestRunReturnsProviderError(t *testing.T) {
	store := &fakeStore{untranslated: []model.Narration{{ID: 1, Source: "x"}}}
	tr := &fakeTranslator{err: errors.New("connection refused")}
	eng := New(store, tr)

	_, err := eng.Run(context.Background(), provider.Connection{}, "m")
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if len(store.updates) != 0 {
		t.Errorf("should not write dest on provider error, got %v", store.updates)
	}
}
