-- Created by Carlos Flores
-- Last modification date: 2021-07-19 06:40:56.311


-- Database creation must be done outside a multicommand file.
-- These commands were put in this file only as a convenience.

-- object: bird | type: ROLE --
-- DROP ROLE IF EXISTS bird;
CREATE ROLE bird WITH LOGIN;
-- ddl-end --

-- -- object: parrot | type: DATABASE --
CREATE DATABASE parrot
    ENCODING = 'UTF8'
    OWNER = bird;
-- -- ddl-end --

\c parrot

-- tables
-- Table: order
CREATE TABLE "order" (
                         id serial  NOT NULL,
                         client_name varchar(100)  NOT NULL,
                         price bigint  NOT NULL,
                         user_email varchar(100)  NOT NULL,
                         created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                         CONSTRAINT order_pk PRIMARY KEY (id)
);
ALTER TABLE public."order" OWNER TO bird;

-- Table: order_product
CREATE TABLE order_product (
                               id serial  NOT NULL,
                               amount int  NOT NULL,
                               order_id int  NOT NULL,
                               product_id int  NOT NULL,
                               CONSTRAINT order_product_pk PRIMARY KEY (id)
);
ALTER TABLE public."order_product" OWNER TO bird;

-- Table: product
CREATE TABLE product (
                         id serial  NOT NULL,
                         name varchar(255)  NOT NULL,
                         price bigint  NOT NULL,
                         description varchar(1000)  NOT NULL,
                         CONSTRAINT product_pk PRIMARY KEY (id)
);
ALTER TABLE public."product" OWNER TO bird;


-- Table: user
CREATE TABLE "user" (
                        email varchar(100)  NOT NULL,
                        full_name varchar(255)  NOT NULL,
                        password_hash varchar(100)  NOT NULL,
                        CONSTRAINT user_pk PRIMARY KEY (email)
);
ALTER TABLE public."user" OWNER TO bird;

-- foreign keys
-- Reference: order_product_order (table: order_product)
ALTER TABLE order_product ADD CONSTRAINT order_product_order
    FOREIGN KEY (order_id)
        REFERENCES "order" (id)
        NOT DEFERRABLE
            INITIALLY IMMEDIATE
;

-- Reference: order_product_product (table: order_product)
ALTER TABLE order_product ADD CONSTRAINT order_product_product
    FOREIGN KEY (product_id)
        REFERENCES product (id)
        NOT DEFERRABLE
            INITIALLY IMMEDIATE
;

-- Reference: order_user (table: order)
ALTER TABLE "order" ADD CONSTRAINT order_user
    FOREIGN KEY (user_email)
        REFERENCES "user" (email)
        NOT DEFERRABLE
            INITIALLY IMMEDIATE
;

-- End of file.
