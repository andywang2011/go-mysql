package mysql

import (
	"context"
	"testing"

	"gorm.io/gorm"
)

func TestMysqlDB_GORM_nilReceiver(t *testing.T) {
	var d *MysqlDB
	if d.GORM() != nil {
		t.Fatal("expected nil")
	}
}

func TestMysqlDB_implementsDB(t *testing.T) {
	var _ DB = (*MysqlDB)(nil)
}

func TestMysqlDB_GORM_nonNilReceiver(t *testing.T) {
	g := &gorm.DB{}
	d := &MysqlDB{gormDB: g}
	if d.GORM() != g {
		t.Fatal("expected gorm pointer")
	}
	if d.Gorm() != g {
		t.Fatal("expected gorm pointer from Gorm()")
	}
}

func TestMysqlDB_Gorm_nilReceiver(t *testing.T) {
	var d *MysqlDB
	if d.Gorm() != nil || d.GORM() != nil {
		t.Fatal("expected nil")
	}
}

func TestMysqlDB_Ping_nil(t *testing.T) {
	var d *MysqlDB
	if err := d.Ping(context.Background()); err == nil {
		t.Fatal("expected error")
	}
}

func TestFakeDB_Ping(t *testing.T) {
	f := &fakeDB{}
	if err := f.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}
}
