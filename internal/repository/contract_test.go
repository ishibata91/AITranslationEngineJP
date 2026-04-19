// Package repository — ER v1 contract tests.
//
// このファイルは repository-contracts-er-v1 handoff の compile-time / runtime contract を確認する。
// product code (interface ファイル) がまだ存在しない場合は compile error になる。
// interface と error 型が追加されると FAIL → PASS に移行する。
package repository

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// ===== Error contract =====

// TestErrNotFoundIsExported は ErrNotFound が exported 非 nil error であることを保証する。
func TestErrNotFoundIsExported(t *testing.T) {
	if ErrNotFound == nil {
		t.Fatal("ErrNotFound must be a non-nil exported error")
	}
	if !errors.Is(ErrNotFound, ErrNotFound) {
		t.Fatal("ErrNotFound must satisfy errors.Is identity check")
	}
}

// TestErrConflictIsExported は ErrConflict が exported 非 nil error であることを保証する。
func TestErrConflictIsExported(t *testing.T) {
	if ErrConflict == nil {
		t.Fatal("ErrConflict must be a non-nil exported error")
	}
	if !errors.Is(ErrConflict, ErrConflict) {
		t.Fatal("ErrConflict must satisfy errors.Is identity check")
	}
}

// TestErrNotFoundAndErrConflictAreDistinct は 2 つのエラーが互いに一致しないことを保証する。
func TestErrNotFoundAndErrConflictAreDistinct(t *testing.T) {
	if errors.Is(ErrNotFound, ErrConflict) {
		t.Fatal("ErrNotFound must not match ErrConflict")
	}
	if errors.Is(ErrConflict, ErrNotFound) {
		t.Fatal("ErrConflict must not match ErrNotFound")
	}
}

// ===== Interface contract =====

// TestTranslationSourceRepositoryIsInterface は TranslationSourceRepository が interface 型であることを保証する。
func TestTranslationSourceRepositoryIsInterface(t *testing.T) {
	typ := reflect.TypeOf((*TranslationSourceRepository)(nil)).Elem()
	if typ.Kind() != reflect.Interface {
		t.Fatalf("TranslationSourceRepository must be an interface, got %s", typ.Kind())
	}
}

// TestFoundationDataRepositoryIsInterface は FoundationDataRepository が interface 型であることを保証する。
func TestFoundationDataRepositoryIsInterface(t *testing.T) {
	typ := reflect.TypeOf((*FoundationDataRepository)(nil)).Elem()
	if typ.Kind() != reflect.Interface {
		t.Fatalf("FoundationDataRepository must be an interface, got %s", typ.Kind())
	}
}

// TestJobLifecycleRepositoryIsInterface は JobLifecycleRepository が interface 型であることを保証する。
func TestJobLifecycleRepositoryIsInterface(t *testing.T) {
	typ := reflect.TypeOf((*JobLifecycleRepository)(nil)).Elem()
	if typ.Kind() != reflect.Interface {
		t.Fatalf("JobLifecycleRepository must be an interface, got %s", typ.Kind())
	}
}

// TestJobOutputRepositoryIsInterface は JobOutputRepository が interface 型であることを保証する。
func TestJobOutputRepositoryIsInterface(t *testing.T) {
	typ := reflect.TypeOf((*JobOutputRepository)(nil)).Elem()
	if typ.Kind() != reflect.Interface {
		t.Fatalf("JobOutputRepository must be an interface, got %s", typ.Kind())
	}
}

// TestTranslationFieldDefinitionRepositoryIsInterface は TranslationFieldDefinitionRepository が interface 型であることを保証する。
func TestTranslationFieldDefinitionRepositoryIsInterface(t *testing.T) {
	typ := reflect.TypeOf((*TranslationFieldDefinitionRepository)(nil)).Elem()
	if typ.Kind() != reflect.Interface {
		t.Fatalf("TranslationFieldDefinitionRepository must be an interface, got %s", typ.Kind())
	}
}

// ===== Transaction boundary contract =====

// TestTransactorIsInterface は Transactor が interface 型であることを保証する。
// Repository 間で同一 transaction を共有できる境界として定義される必要がある。
func TestTransactorIsInterface(t *testing.T) {
	typ := reflect.TypeOf((*Transactor)(nil)).Elem()
	if typ.Kind() != reflect.Interface {
		t.Fatalf("Transactor must be an interface, got %s", typ.Kind())
	}
}

// ===== Method count contract =====

// TestTranslationSourceRepositoryMethodCount は TranslationSourceRepository のメソッド数が期待値と一致することを保証する。
func TestTranslationSourceRepositoryMethodCount(t *testing.T) {
	typ := reflect.TypeOf((*TranslationSourceRepository)(nil)).Elem()
	const expected = 14
	if got := typ.NumMethod(); got != expected {
		t.Errorf("TranslationSourceRepository: expected %d methods, got %d", expected, got)
	}
}

// TestFoundationDataRepositoryMethodCount は FoundationDataRepository のメソッド数が期待値と一致することを保証する。
func TestFoundationDataRepositoryMethodCount(t *testing.T) {
	typ := reflect.TypeOf((*FoundationDataRepository)(nil)).Elem()
	const expected = 12
	if got := typ.NumMethod(); got != expected {
		t.Errorf("FoundationDataRepository: expected %d methods, got %d", expected, got)
	}
}

// TestJobLifecycleRepositoryMethodCount は JobLifecycleRepository のメソッド数が期待値と一致することを保証する。
func TestJobLifecycleRepositoryMethodCount(t *testing.T) {
	typ := reflect.TypeOf((*JobLifecycleRepository)(nil)).Elem()
	const expected = 10
	if got := typ.NumMethod(); got != expected {
		t.Errorf("JobLifecycleRepository: expected %d methods, got %d", expected, got)
	}
}

// TestJobOutputRepositoryMethodCount は JobOutputRepository のメソッド数が期待値と一致することを保証する。
func TestJobOutputRepositoryMethodCount(t *testing.T) {
	typ := reflect.TypeOf((*JobOutputRepository)(nil)).Elem()
	const expected = 4
	if got := typ.NumMethod(); got != expected {
		t.Errorf("JobOutputRepository: expected %d methods, got %d", expected, got)
	}
}

// TestTranslationFieldDefinitionRepositoryMethodCount は TranslationFieldDefinitionRepository のメソッド数が期待値と一致することを保証する。
func TestTranslationFieldDefinitionRepositoryMethodCount(t *testing.T) {
	typ := reflect.TypeOf((*TranslationFieldDefinitionRepository)(nil)).Elem()
	const expected = 5
	if got := typ.NumMethod(); got != expected {
		t.Errorf("TranslationFieldDefinitionRepository: expected %d methods, got %d", expected, got)
	}
}

// TestTransactorMethodCount は Transactor のメソッド数が期待値と一致することを保証する。
func TestTransactorMethodCount(t *testing.T) {
	typ := reflect.TypeOf((*Transactor)(nil)).Elem()
	const expected = 1
	if got := typ.NumMethod(); got != expected {
		t.Errorf("Transactor: expected %d methods, got %d", expected, got)
	}
}

// ===== Transactor.WithTransaction signature contract =====

// TestTransactorWithTransactionSignature は Transactor.WithTransaction の引数型・返り値型を保証する。
func TestTransactorWithTransactionSignature(t *testing.T) {
	typ := reflect.TypeOf((*Transactor)(nil)).Elem()
	m, ok := typ.MethodByName("WithTransaction")
	if !ok {
		t.Fatal("Transactor must have WithTransaction method")
	}
	mt := m.Type
	// interface method の In は receiver を含まない: In(0)=ctx, In(1)=fn
	ctxType := reflect.TypeOf((*context.Context)(nil)).Elem()
	if mt.NumIn() != 2 || mt.In(0) != ctxType {
		t.Errorf("WithTransaction: first arg must be context.Context, got: %v (NumIn=%d)", mt.In(0), mt.NumIn())
	}
	// fn の型: func(context.Context) error
	fnType := reflect.TypeOf((*func(context.Context) error)(nil)).Elem()
	if mt.NumIn() == 2 && mt.In(1) != fnType {
		t.Errorf("WithTransaction: second arg must be func(context.Context) error, got: %v", mt.In(1))
	}
	// 返り値: error
	errType := reflect.TypeOf((*error)(nil)).Elem()
	if mt.NumOut() != 1 || mt.Out(0) != errType {
		t.Errorf("WithTransaction: return must be error, got: %v (NumOut=%d)", mt.Out(0), mt.NumOut())
	}
}

// ===== errors.Is wrapping contract =====

// TestErrNotFoundWrapsDetectably は ErrNotFound が fmt.Errorf("%w") でラップ後も errors.Is で検出できることを保証する。
func TestErrNotFoundWrapsDetectably(t *testing.T) {
	wrapped := fmt.Errorf("wrapped: %w", ErrNotFound)
	if !errors.Is(wrapped, ErrNotFound) {
		t.Error("ErrNotFound must be detectable via errors.Is after wrapping")
	}
}

// TestErrConflictWrapsDetectably は ErrConflict が fmt.Errorf("%w") でラップ後も errors.Is で検出できることを保証する。
func TestErrConflictWrapsDetectably(t *testing.T) {
	wrapped := fmt.Errorf("wrapped: %w", ErrConflict)
	if !errors.Is(wrapped, ErrConflict) {
		t.Error("ErrConflict must be detectable via errors.Is after wrapping")
	}
}

// ===== Error message contract =====

// TestErrNotFoundHasMessage は ErrNotFound が空でないメッセージを持つことを保証する。
func TestErrNotFoundHasMessage(t *testing.T) {
	if ErrNotFound.Error() == "" {
		t.Error("ErrNotFound must have a non-empty error message")
	}
}

// TestErrConflictHasMessage は ErrConflict が空でないメッセージを持つことを保証する。
func TestErrConflictHasMessage(t *testing.T) {
	if ErrConflict.Error() == "" {
		t.Error("ErrConflict must have a non-empty error message")
	}
}
