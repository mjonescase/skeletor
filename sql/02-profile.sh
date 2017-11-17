#!/bin/bash
set -e

psql -U postgres skeletor <<-EOSQL
CREATE TABLE profile (
        id SERIAL PRIMARY KEY,
        firstname varchar(255),
        lastname varchar(255),
        username varchar(255),
        email varchar(255),
        title varchar(255),
        password varchar(255),
        mobilenumber varchar(255)
);

ALTER TABLE profile ADD CONSTRAINT profile_email_unique UNIQUE (email)
EOSQL
