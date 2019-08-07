#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 -U postgres <<-EOSQL
    CREATE USER skeletor;
    CREATE DATABASE skeletor;
    GRANT ALL PRIVILEGES ON DATABASE skeletor TO postgres;
EOSQL
