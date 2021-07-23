package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	log "github.com/sirupsen/logrus"

	"github.com/hellerox/parrot/model"
)

// Connect function to start db connection
func connect(connectionString string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Fatalln("unable to create connection", err)
	}

	db, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalln("unable to create connection", err)
	}

	return db
}

// InsertUser for
func (ds *DatabaseStorage) InsertUser(u model.User) error {
	_, err := ds.pool.Exec(context.Background(),
		`INSERT INTO "user" (email, full_name, password_hash) VALUES ($1,$2,$3)`,
		u.Email, u.FullName, u.PasswordHash)
	if err != nil {
		log.Errorf("error inserting user: %s", err.Error())
		return err
	}

	return nil
}

func (ds *DatabaseStorage) GetUserHash(mail string) string {
	rows, err := ds.pool.Query(context.Background(), `SELECT password_hash FROM "user" WHERE email = $1`, mail)
	if err != nil {
		log.Errorf("error getting user hash: %s", err.Error())
		return ""
	}

	for rows.Next() {
		var hash string
		rows.Scan(&hash)
		return hash
	}

	return ""
}

// InsertOrder on DB using the given data
func (ds *DatabaseStorage) InsertOrder(o model.Order) (int, error) {
	rows, err := ds.pool.Query(context.Background(),
		`INSERT INTO "order" (client_name, price, user_email) VALUES ($1,$2,$3) RETURNING id`, o.ClientName, o.Price, o.Email)
	if err != nil {
		log.Errorf("error inserting order: %s", err.Error())
		return 0, err
	}

	for rows.Next() {
		err := rows.Scan(&o.ID)
		if err != nil {
			log.Errorf("error inserting order: %s", err.Error())
			return 0, err
		}

		return o.ID, nil
	}

	return 0, nil
}

// UpdatePriceOrder on DB using the given data
func (ds *DatabaseStorage) UpdatePriceOrder(o model.Order) error {
	_, err := ds.pool.Exec(context.Background(),
		`UPDATE "order" SET price = $1 WHERE id = $2`, o.Price, o.ID)
	if err != nil {
		log.Errorf("error updating order: %s", err.Error())
		return err
	}

	return nil
}

// InsertProduct is the function to insert products in DB, if the name exists it gets the id from DB
func (ds *DatabaseStorage) InsertProduct(p model.Product) (int, error) {
	rows, err := ds.pool.Query(context.Background(),
		`WITH sel AS
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
`,
		p.Name, p.Price, p.Description)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 0, fmt.Errorf("didn't insert data, it already exists")
		}

		log.Errorf("error inserting product: %s", err.Error())

		return 0, err
	}

	for rows.Next() {
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description)
		if err != nil {
			log.Errorf("error inserting order: %s", err.Error())
			return 0, err
		}

		return p.ID, nil
	}

	return 0, nil
}

// InsertOrderProductRelation using order id and product id
func (ds *DatabaseStorage) InsertOrderProductRelation(op model.OrderProductRelation) (int, error) {
	var opID int

	err := ds.pool.QueryRow(context.Background(),
		`INSERT INTO order_product(amount, order_id, product_id) VALUES ($1,$2,$3) RETURNING id`,
		&op.Amount, &op.OrderID, &op.ProductID).Scan(&opID)
	if err != nil {
		log.Errorf("error inserting order product relation: %s", err.Error())
		return 0, err
	}

	return opID, nil
}

func (ds *DatabaseStorage) GetReportData(r model.GenerateReportRequest) (model.GenerateReportResponse, error) {
	var response model.GenerateReportResponse

	rows, err := ds.pool.Query(context.Background(), `SELECT
   product.name, sum(order_product.amount), product.price * sum(order_product.amount) total
FROM
   public.order 
   INNER JOIN
      public."order_product" 
      ON "order".id = order_product.order_id 
   INNER JOIN
      public.product 
      ON order_product.product_id = product.id 
WHERE
   created_at >= $1
   AND created_at <  $2
GROUP BY
   product.name, "product".price
ORDER BY
   sum(order_product.amount) DESC;`, r.From, r.To)
	if err != nil {
		log.Errorf("error getting report data: %s", err.Error())
		return response, err
	}
	defer rows.Close()

	for rows.Next() {
		var data model.Data
		if err := rows.Scan(&data.Name, &data.TotalAmount, &data.TotalPrice); err != nil {
			log.Errorln(err)
		}

		response.Data = append(response.Data, data)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		log.Panicln(err)
	}

	return response, err
}
