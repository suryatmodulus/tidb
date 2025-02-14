// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"

	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/sessionctx/variable"
	"github.com/pingcap/tidb/util/testbridge"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

var isCoverageServer string

func TestMain(m *testing.M) {
	testbridge.WorkaroundGoCheckFlags()
	opts := []goleak.Option{
		goleak.IgnoreTopFunction("go.etcd.io/etcd/pkg/logutil.(*MergeLogger).outputLoop"),
		goleak.IgnoreTopFunction("go.opencensus.io/stats/view.(*worker).start"),
	}
	goleak.VerifyTestMain(m, opts...)
}

// TestRunMain is a dummy test case, which contains only the main function of tidb-server,
// and it is used to generate coverage_server.
func TestRunMain(t *testing.T) {
	if isCoverageServer == "1" {
		main()
	}
}

func TestSetGlobalVars(t *testing.T) {
	require.Equal(t, "tikv,tiflash,tidb", variable.GetSysVar(variable.TiDBIsolationReadEngines).Value)
	require.Equal(t, "1073741824", variable.GetSysVar(variable.TiDBMemQuotaQuery).Value)

	config.UpdateGlobal(func(conf *config.Config) {
		conf.IsolationRead.Engines = []string{"tikv", "tidb"}
		conf.MemQuotaQuery = 9999999
	})
	setGlobalVars()

	require.Equal(t, "tikv,tidb", variable.GetSysVar(variable.TiDBIsolationReadEngines).Value)
	require.Equal(t, "9999999", variable.GetSysVar(variable.TiDBMemQuotaQuery).Value)
}
