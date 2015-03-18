#!/bin/bash

if [ $# -ne 3 ]; then
  echo "Usage: $0 HOST USERNAME PASSWORD" >&2
  echo >&2
  echo "Example: $0 mysql_db.com username password" >&2
  exit 1
fi

set -e

HOST=$1
USERNAME=$2
PASSWORD=$3

#Don't select mysql system DBs or column title
export CF_DATABASES=$(mysql -h $HOST -u $USERNAME "-p$PASSWORD" -e "show databases" | grep -v Database \
  | grep -v mysql | grep -v performance_schema | grep -v information_schema  | grep -v innodb)
echo "Databases to be dropped and recreated:
$CF_DATABASES"

for DB in $CF_DATABASES
do
  echo "Dropping and recreating database: $DB"
  $(mysql -h $HOST -u $USERNAME -p$PASSWORD -e "DROP DATABASE $DB; CREATE DATABASE $DB");
done
