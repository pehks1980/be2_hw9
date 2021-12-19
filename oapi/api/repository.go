package api

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	//"pehks1980/be2_hw9/oapi/api"
)

type Repo struct {
	pool *pgxpool.Pool
}

/*
type PricelistR struct {
	Id    int    `json:"id"`
	Good  string `json:"good"`
	Price int    `json:"price"`
}
*/
const (
	DDL = `
	BEGIN;
	-- Adminer 4.7.7 PostgreSQL dump
	DROP TABLE IF EXISTS "pricelists";
	DROP TABLE IF EXISTS "goods";
	DROP SEQUENCE IF EXISTS "Goods_id_seq";
	DROP SEQUENCE IF EXISTS "PriceList_id_seq";
	CREATE SEQUENCE "Goods_id_seq" INCREMENT 1 MINVALUE 100 START 100;
	CREATE SEQUENCE "PriceList_id_seq" INCREMENT 1 MINVALUE 4 START 4;

	CREATE TABLE "public"."goods" (
		"id" bigint DEFAULT nextval('"Goods_id_seq"') NOT NULL,
		"name" character varying NOT NULL,
		CONSTRAINT "goods_pkey" PRIMARY KEY ("id")
	);
	
	INSERT INTO "goods" ("id", "name") VALUES
	(1,	'bricks'),
	(2,	'wood'),
	(3,	'ceiling'),
	(4,	'furniture'),
	(5,	'tv-set');
	
	
	CREATE TABLE "public"."pricelists" (
		"id" bigint DEFAULT nextval('"PriceList_id_seq"') NOT NULL,
		"price" numeric NOT NULL,
		"good_id" bigint NOT NULL,
		"pricelist_id" bigint NOT NULL,
		CONSTRAINT "pricelist_pkey" PRIMARY KEY ("id"),
		CONSTRAINT "fk_good" FOREIGN KEY (good_id) REFERENCES goods(id) ON DELETE CASCADE NOT DEFERRABLE
	);
	
	INSERT INTO "pricelists" ("id", "price", "good_id", "pricelist_id") VALUES
	(1,	20.55,	1,	1),
	(3,	45.99,	2,	1),
	(2,	100.33,	5,	2);
	
	-- 2021-12-16 04:44:53.24713+00
	COMMIT;
	`
	// $1 pricelist id
	GetPriceListSQL = `
	SELECT goods.id::integer, goods.name::varchar, pricelists.price::integer
	FROM pricelists
	INNER JOIN goods ON pricelists.good_id = goods.id 
	WHERE pricelist_id = $1;
	`
	//$1 pricelist id $2 good $3 price
	InsertPriceListEntitySQL = `
	WITH ins1 AS (INSERT INTO "goods" (name) values ($2)
   	RETURNING id)
	INSERT INTO "pricelists" ("price", "good_id", "pricelist_id") 
	VALUES ($3, (select id from ins1), $1)
	returning id;
	`
	// $1 good_id $2 name
	UpdatePriceListEntity1SQL = `
	 UPDATE "goods"
     SET name = $2
     WHERE id = $1	
	`
	// $1 pricelist_id $2 good_id $3 price
	UpdatePriceListEntity2SQL = `
	 UPDATE "pricelists"
     SET price = $3
     WHERE good_id = $2	AND pricelist_id = $1
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
	_, err := r.pool.Exec(ctx, UpdatePriceListEntity1SQL, pricelistentity.Id, pricelistentity.Good)
	if err != nil {
		return fmt.Errorf("failed to update pricelist: %w", err)
	}
	_, err = r.pool.Exec(ctx, UpdatePriceListEntity2SQL, id, pricelistentity.Id, pricelistentity.Price)
	if err != nil {
		return fmt.Errorf("failed to update pricelist: %w", err)
	}
	return nil
}
