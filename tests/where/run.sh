#!/bin/sh

set -eu

run_sql "drop database if exists where_test"
run_sql "create database where_test"
export DUMPLING_TEST_DATABASE=where_test
run_sql "create table t (a int)"
run_sql "insert into t values $(seq -s, 200 | sed 's/,*$//g' | sed "s/[0-9]*/('&')/g");"

# only dump rows whose a > 100, 100 rows are expected
run_dumpling --where "a > 100"

# the WHERE condition should be written out in the dump
grep -q "WHERE (a > 100)" "$DUMPLING_OUTPUT_DIR/where_test.t.sql"

total_lines=$(grep -c "^(" "$DUMPLING_OUTPUT_DIR/where_test.t.sql")
if [ "$total_lines" != "100" ]; then
  echo "obtain record number: $total_lines, but expect: 100" && exit 1
fi

# if all rows are filtered out, a placeholder file should still be generated
run_dumpling -o "$DUMPLING_OUTPUT_DIR/all_filtered" --where "a > 10000"

file_should_exist "$DUMPLING_OUTPUT_DIR/all_filtered/where_test.t.sql"
grep -q "All data are filtered out by this condition:" "$DUMPLING_OUTPUT_DIR/all_filtered/where_test.t.sql"
grep -q "WHERE (a > 10000)" "$DUMPLING_OUTPUT_DIR/all_filtered/where_test.t.sql"

# an invalid WHERE expression should abort dumpling with an error
if run_dumpling -o "$DUMPLING_OUTPUT_DIR/bad_where" --where "not_a_column > 10000"; then
  echo "dumpling should fail with an invalid WHERE expression" && exit 1
fi
