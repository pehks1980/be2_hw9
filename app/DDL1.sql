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

-- select data
SELECT goods.id::integer, goods.name::varchar, pricelists.price::integer
FROM pricelists
INNER JOIN goods ON pricelists.id = goods.pricelists_id
WHERE pricelist_id = 1;

-- insert data

WITH ins1 AS (INSERT INTO pricelists (price, pricelist_id) values (99.99,3)
    RETURNING id)
INSERT INTO goods (name, pricelists_id)
VALUES ('MIG35', (select id from ins1))
returning id;


-- update data

WITH upd1 AS (UPDATE goods
    SET name = 'bricsss'
    WHERE id = 3
    RETURNING pricelists_id)

UPDATE pricelists
SET price = 99.99
WHERE id = (select pricelists_id from upd1)
RETURNING id;

-- delete data

DELETE FROM pricelists
WHERE pricelist_id = 3;