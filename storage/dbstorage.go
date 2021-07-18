package storage

import (
	"context"
	"fmt"

	// postgres package
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"

	"github.com/hellerox/parrot/model"
)

func (ds *DatabaseStorage) PrepareStatements() error {

	_, err := ds.conn.Prepare(context.Background(), "insertUser", `INSERT INTO "user" (email, full_name, password) 
			VALUES ($1,$2,$3)`)
	if err != nil {
		log.Fatalln("error preparing statements: ", err.Error())
	}

	_, err = ds.conn.Prepare(context.Background(), "insertOrder", `INSERT INTO "order" (client_name, price, user_email)
			VALUES ($1,$2,$3) RETURNING id`)
	if err != nil {
		log.Fatalln("error preparing statements: ", err.Error())
	}

	_, err = ds.conn.Prepare(context.Background(), "updateOrder", `UPDATE "order" SET price = $1 WHERE id = $2`)
	if err != nil {
		log.Fatalln("error preparing statements: ", err.Error())
	}

	_, err = ds.conn.Prepare(context.Background(), "insertOrderProductRelation", `INSERT INTO order_product(amount, order_id, product_id) 
		VALUES ($1,$2,$3) RETURNING id`)
	if err != nil {
		log.Fatalln("error preparing statements: ", err.Error())
	}

	_, err = ds.conn.Prepare(context.Background(), "insertProduct", `WITH sel AS
(
       SELECT id,
              "name",
              "price",
              "description"
       FROM   product
       WHERE  name = $1 ), ins AS
(
            insert INTO product
                        (
                                    "name",
                                    "price",
                                    "description"
                        )
            SELECT    $1,
                      $2,$3
            WHERE     NOT EXISTS
                      (
                             SELECT 1
                             FROM   sel)
            returning id,
                      "name",
                      "price",
                      "description" )
SELECT id,
       "name",
       "price",
       "description"
FROM   ins
UNION ALL
SELECT id,
       "name",
       "price",
       "description"
FROM   sel
`)
	if err != nil {
		log.Fatalln("error preparing statements: ", err.Error())
	}

	return nil
}

// Connect function to start db connection
func connect(connectionString string) *pgx.Conn {
	db, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		log.Fatalln("unable to create connection", err)
	}

	return db
}

// InsertUser for
func (ds *DatabaseStorage) InsertUser(u model.User) error {
	_, err := ds.conn.Exec(context.Background(), "insertUser", u.Email, u.FullName, u.Password)
	if err != nil {
		log.Errorf("error inserting user: %s", err.Error())
		return err
	}

	return nil
}

// InsertOrder on DB using the given data
func (ds *DatabaseStorage) InsertOrder(o model.Order) (int, error) {
	err := ds.conn.QueryRow(context.Background(),
		"insertOrder", o.ClientName, o.Price, o.Email).Scan(&o.ID)
	if err != nil {
		log.Errorf("error inserting order: %s", err.Error())
		return 0, err
	}

	return o.ID, nil
}

// UpdateOrder on DB using the given data
func (ds *DatabaseStorage) UpdateOrder(o model.Order) error {
	_, err := ds.conn.Exec(context.Background(),
		"updateOrder", o.Price, o.ID)
	if err != nil {
		log.Errorf("error updating order: %s", err.Error())
		return err
	}

	return nil
}

// InsertProduct is the function to insert products in DB, if the name exists it gets the id from DB
func (ds *DatabaseStorage) InsertProduct(p model.Product) (int, error) {
	err := ds.conn.QueryRow(context.Background(),
		"insertProduct",
		p.Name, p.Price, p.Description).Scan(&p.ID, &p.Name, &p.Price, &p.Description)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 0, fmt.Errorf("didn't insert data, it already exists")
		}

		log.Errorf("error inserting product: %s", err.Error())

		return 0, err
	}

	return p.ID, nil
}

// InsertOrderProductRelation using order id and product id
func (ds *DatabaseStorage) InsertOrderProductRelation(cm model.OrderProductRelation) (int, error) {
	var opId int

	err := ds.conn.QueryRow(context.Background(),
		"insertOrderProductRelation",
		&cm.Amount, &cm.OrderID, &cm.ProductID).Scan(&opId)
	if err != nil {
		log.Errorf("error inserting order product relation: %s", err.Error())
		return 0, err
	}

	return opId, nil
}
