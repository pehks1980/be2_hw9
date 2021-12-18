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