package storage

import (
	"context"
	"testing"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

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

	mockPool.EXPECT().Query(context.Background(), `SELECT password_hash FROM "user" WHERE email = $1`, "hellerox@gmail.com").Return(pgxRows, nil)

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

func TestDatabaseStorage_InsertOrder(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		pool pgxpoolmock.PgxPool
	}

	type args struct {
		o model.Order
	}

	o := model.Order{
		Email:      "hellerox@gmail.com",
		ClientName: "Carlos dos",
		Price:      1555,
	}

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	pgxRows := pgxpoolmock.NewRows([]string{"id"}).AddRow(1).ToPgxRows()

	mockPool.EXPECT().Query(context.Background(),
		`INSERT INTO "order" (client_name, price, user_email) VALUES ($1,$2,$3) RETURNING id`,
		o.ClientName, o.Price, o.Email).Return(pgxRows, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr error
	}{
		{"InsertOrder", fields{pool: mockPool}, args{o: o}, 1, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DatabaseStorage{
				pool: tt.fields.pool,
			}

			got, err := ds.InsertOrder(tt.args.o)
			if tt.wantErr != nil {
				assert.Equal(t, err, tt.wantErr)
				return
			}

			assert.Equal(t, got, tt.want)
		})
	}
}

func TestDatabaseStorage_UpdatePriceOrder(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		pool pgxpoolmock.PgxPool
	}

	type args struct {
		o model.Order
	}

	o := model.Order{
		ID:         11213,
		ClientName: "Carlos dos",
		Price:      15,
	}

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	mockPool.EXPECT().Exec(context.Background(), `UPDATE "order" SET price = $1 WHERE id = $2`, o.Price, o.ID)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{"UpdatePriceOrder", fields{pool: mockPool}, args{o: o}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DatabaseStorage{
				pool: tt.fields.pool,
			}

			err := ds.UpdatePriceOrder(tt.args.o)
			if tt.wantErr != nil {
				assert.Equal(t, err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseStorage_InsertProduct(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		pool pgxpoolmock.PgxPool
	}

	type args struct {
		p model.Product
	}

	p := model.Product{
		Name:        "libro",
		Price:       1234,
		Description: "Producto Libro",
	}

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	pgxRows := pgxpoolmock.NewRows([]string{"id", "name", "price", "description"}).AddRow(1, p.Name, p.Price, p.Description).ToPgxRows()

	mockPool.EXPECT().Query(context.Background(),
		gomock.Any(),
		p.Name, p.Price, p.Description).Return(pgxRows, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr error
	}{
		{"InsertProduct", fields{pool: mockPool}, args{p: p}, 1, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DatabaseStorage{
				pool: tt.fields.pool,
			}

			got, err := ds.InsertProduct(tt.args.p)
			if tt.wantErr != nil {
				assert.Equal(t, err, tt.wantErr)
				return
			}

			assert.Equal(t, got, tt.want)
		})
	}
}
