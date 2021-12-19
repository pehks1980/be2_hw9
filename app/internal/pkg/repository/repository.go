package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Repo struct {
	pool *pgxpool.Pool
}

// этот репозиторий не проверяет консистентность и наличие всех данных для совершения действа с бд!
// be aware!!!
type Pricelist struct {
	Id    int    `json:"id"`
	Good  string `json:"good"`
	Price int    `json:"price"`
}

const (
	DDL = `
	BEGIN;
	-- Adminer 4.7.7 PostgreSQL dump
	DROP TABLE IF EXISTS "pricelists" CASCADE ;
	DROP TABLE IF EXISTS "goods";
	DROP SEQUENCE IF EXISTS "Goods_id_seq";
	DROP SEQUENCE IF EXISTS "PriceList_id_seq";
	CREATE SEQUENCE "Goods_id_seq" INCREMENT 1 MINVALUE 6 START 6;
	CREATE SEQUENCE "PriceList_id_seq" INCREMENT 1 MINVALUE 6 START 6;
	
	CREATE TABLE "public"."pricelists" (
		   "id" bigint DEFAULT nextval('"PriceList_id_seq"') NOT NULL,
		   "price" numeric NOT NULL,
		   "pricelist_id" bigint NOT NULL,
		   CONSTRAINT "pricelists_pkey" PRIMARY KEY ("id")
	);
	
	CREATE TABLE "public"."goods" (
		  "id" bigint DEFAULT nextval('"Goods_id_seq"') NOT NULL,
		  "name" character varying NOT NULL,
		  "pricelists_id" bigint NOT NULL,
		  CONSTRAINT "goods_pkey" PRIMARY KEY ("id"),
		  CONSTRAINT "fk_pricelists" FOREIGN KEY (pricelists_id) REFERENCES pricelists(id) ON DELETE CASCADE NOT DEFERRABLE
	);
	
	INSERT INTO "pricelists" ("id", "price", "pricelist_id") VALUES
			(1,	20.55,	1),
			(2,	45.99,	1),
			(3,	100.33,	2),
			(4,	10.33,	2),
			(5,	1000.33, 3);
	
	INSERT INTO "goods" ("id", "name", "pricelists_id") VALUES
										   (1,	'bricks', 1),
										   (2,	'wood', 2),
										   (3,	'ceiling', 3 ),
										   (4,	'furniture', 4),
										   (5,	'tv-set', 5);
	
	
	
	-- 2021-12-16 04:44:53.24713+00
	COMMIT;
	`
	// $1 pricelist id
	GetPriceListSQL = `
	SELECT goods.id::integer, goods.name::varchar, pricelists.price::integer
	FROM pricelists
	INNER JOIN goods ON pricelists.id = goods.pricelists_id
	WHERE pricelist_id = $1;
	`
	//$1 pricelist id $2 good $3 price
	InsertPriceListEntitySQL = `
	WITH ins1 AS (INSERT INTO pricelists (price, pricelist_id) values ($3, $1)
	RETURNING id)
	INSERT INTO goods (name, pricelists_id)
	VALUES ( $2, (select id from ins1))
	returning id;
	`
	// $1 good_id $2 name $3 price
	UpdatePriceListEntitySQL = `
	WITH upd1 AS (UPDATE goods
    SET name = $2
    WHERE id = $1
    RETURNING pricelists_id)
	UPDATE pricelists
	SET price = $3
	WHERE id = (select pricelists_id from upd1)
	RETURNING id;
	`
	DeletePriceListSQL = `
	DELETE FROM "pricelists" 
	WHERE "pricelist_id" = $1
	`
)

func NewRepository(ctx context.Context, dbconn string) (*Repo, error) {

	config, err := pgxpool.ParseConfig(dbconn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse conn string (%s): %w", dbconn, err)
	}
	config.ConnConfig.LogLevel = pgx.LogLevelDebug
	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &Repo{pool: pool}, nil
}

func (r *Repo) InitSchema(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, DDL)
	return err
}

func (r *Repo) GetPriceList(ctx context.Context, id int) ([]Pricelist, error) {

	rows, _ := r.pool.Query(ctx, GetPriceListSQL, id)
	var pricelist []Pricelist

	for rows.Next() {
		var pricelistrow Pricelist

		if err := rows.Scan(&pricelistrow.Id, &pricelistrow.Good, &pricelistrow.Price); err != nil {
			log.Printf("error get from sql: %v", err)
			return nil, err
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}
		pricelist = append(pricelist, pricelistrow)
	}

	return pricelist, nil
}

func (r *Repo) AddPriceListEntity(ctx context.Context, pricelistentity Pricelist) (int, error) {
	var rid int
	err := r.pool.QueryRow(ctx, InsertPriceListEntitySQL,
		pricelistentity.Id,
		pricelistentity.Good,
		pricelistentity.Price,
	).Scan(&rid)
	if err != nil {
		return 0, fmt.Errorf("failed to add pricelistentity: %w", err)
	}
	return rid, nil
}

func (r *Repo) DelPriceList(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, DeletePriceListSQL, id)
	if err != nil {
		return fmt.Errorf("failed to del pricelist: %w", err)
	}
	return nil
}

func (r *Repo) UpdatePriceList(ctx context.Context, id int, pricelistentity Pricelist) error {
	_, err := r.pool.Exec(ctx, UpdatePriceListEntitySQL, pricelistentity.Id, pricelistentity.Good, pricelistentity.Price)
	if err != nil {
		return fmt.Errorf("failed to update pricelist: %w", err)
	}
	return nil
}
