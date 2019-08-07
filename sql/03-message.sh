#!/bin/bash
set -e

psql -U postgres skeletor <<-EOSQL

CREATE TABLE conversation (
       id              SERIAL PRIMARY KEY,
       createddate     TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
       modifieddate    TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

CREATE TABLE conversation_profile (
       id              SERIAL PRIMARY KEY,
       createddate     TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
       modifieddate    TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
       conversation_id INTEGER NOT NULL REFERENCES conversation (id),
       profile_id      INTEGER NOT NULL REFERENCES profile (id)
);

CREATE TABLE message (
        id   	                SERIAL PRIMARY KEY,
	createddate             TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    	modifieddate    	TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
	conversation_profile_id INTEGER NOT NULL REFERENCES conversation_profile (id),
	content                 VARCHAR(10000)
);
EOSQL
