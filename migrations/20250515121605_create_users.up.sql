CREATE TABLE users
(
    id bigserial NOT NULL PRIMARY KEY,
    email varchar NOT NULL UNIQUE,
    enc_password varchar NOT NULL
);
