CREATE DATABASE mercrawl;

CREATE TABLE "items"
(
id char(20),
name text,
photo1 text,
photo2 text,
photo3 text,
photo4 text,
status char(40),
price integer,
shippingFee char(20),
description text,
url text
);

ALTER TABLE items ADD PRIMARY KEY (id);

ALTER TABLE items ADD COLUMN sent bool;

CREATE INDEX sent_idx ON items (sent);

ALTER TABLE items ALTER COLUMN sent SET DEFAULT FALSE;

UPDATE items SET sent = false;

ALTER TABLE items ALTER COLUMN sent SET NOT NULL;

ALTER TABLE items ALTER COLUMN id TYPE varchar(20)
ALTER TABLE items ALTER COLUMN status TYPE varchar(40)
ALTER TABLE items ALTER COLUMN shippingFee TYPE varchar(20)
