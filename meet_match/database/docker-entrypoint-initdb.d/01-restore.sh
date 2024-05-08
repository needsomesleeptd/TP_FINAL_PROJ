#!/bin/sh

file="/docker-entrypoint-initdb.d/dump.pgdata"
dbname=meetmatch_db

echo "Restoring DB using $file"
pg_restore -U any1 --dbname=$dbname --verbose --single-transaction < "$file" || exit 1