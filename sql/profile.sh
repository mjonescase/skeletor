#!/bin/bash
set -e

psql -U postgres skeletor <<-EOSQL
CREATE TABLE profile (
        id SERIAL,
        firstname varchar(255),
        lastname varchar(255),
        username varchar(255),
        title varchar(255),
        password varchar(255),
        mobilenumber varchar(255)
);
EOSQL
