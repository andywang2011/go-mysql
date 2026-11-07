package mysql

import (
	"context"
	"testing"

	"gorm.io/gorm"
)

type fakeDB struct {
	closed bool
}

func (f *fakeDB) GORM() *gorm.DB             { return nil }
func (f *fakeDB) Ping(context.Context) error { return nil }
func (f *fakeDB) Close() error {
	f.closed = true
	return nil
}

var _ DB = (*fakeDB)(nil)

func TestFakeDB_Close(t *testing.T) {
	f := &fakeDB{}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	if !f.closed {
		t.Fatal("expected closed")
	}
}
