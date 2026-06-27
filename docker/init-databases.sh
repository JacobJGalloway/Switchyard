#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE switchyard;
    CREATE DATABASE switchyard_inventory;
    CREATE DATABASE switchyard_logistics;
EOSQL
