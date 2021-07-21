package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/hellerox/parrot/model"
)

func TestDatabaseStorage_InsertUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		pool pgxpoolmock.PgxPool
	}

	type args struct {
		u model.User
	}

	uSuccess := model.User{
		Email:        "hellerox@gmail.com",
		FullName:     "Carlos Flores",
		PasswordHash: "123",
	}

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	mockPool.EXPECT().Exec(gomock.Any(), `INSERT INTO "user" (email, full_name, password_hash) VALUES ($1,$2,$3)`, "hellerox@gmail.com", "Carlos Flores", "123")

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{"Insert", fields{pool: mockPool}, args{u: uSuccess}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DatabaseStorage{
				pool: tt.fields.pool,
			}

			err := ds.InsertUser(tt.args.u)
			if tt.wantErr != nil {
				assert.Equal(t, err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseStorage_GetUserHash(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		pool pgxpoolmock.PgxPool
	}

	type args struct {
		mail string
	}

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	pgxRows := pgxpoolmock.NewRows([]string{"password_hash"}).AddRow("123").ToPgxRows()

	mockPool.EXPECT().QueryRow(context.Background(), `SELECT password_hash FROM "user" WHERE email = $1`, "hellerox@gmail.com").Return(pgxRows)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"getHash", fields{pool: mockPool}, args{mail: "hellerox@gmail.com"}, "123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DatabaseStorage{
				pool: tt.fields.pool,
			}

			got := ds.GetUserHash(tt.args.mail)
			assert.Equal(t, got, tt.want)
		})
	}
}
