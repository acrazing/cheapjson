#!/usr/bin/env bash
#
# profile.sh
# @author acrazing
# @since 2017-07-22 14:39:45
# @desc profile.sh
#

if [ ! -f ./data/normal.json ]; then
  mkdir -p ./data
  go test -v -run NONE ./parser_test.go
fi

go test -v -cpuprofile cpu.prof -memprofile mem.prof -run Profile parser_test.go
go tool pprof -pdf -focus testing.tRunner -output prof_cpu.pdf cheapjson.test cpu.prof
# Both is ignored...
go tool pprof -pdf -focus testing.tRunner -output prof_mem.pdf cheapjson.test mem.prof
go tool pprof -pdf -output full_mem.pdf cheapjson.test mem.prof

rm -rf *.gif *.prof *.test
