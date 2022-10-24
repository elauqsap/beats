// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	conf "github.com/elastic/elastic-agent-libs/config"
)

func TestAgentInputNormalization(t *testing.T) {
	cfg, err := conf.NewConfigFrom(`
type: packet
data_stream:
  namespace: default
processors:
  - add_fields:
      target: 'elastic_agent'
      fields:
        id: agent-id
        version: 8.0.0
        snapshot: false
streams:
  - type: flow
    timeout: 10s
    period: 10s
    keep_null: false
    interface:
      device: default_route
      snaplen: 1514
      type: af_packet
      buffer_size_mb: 100
    procs:
      enabled: true
      monitored:
        - process: mysqld
          cmdline_grep: mysqld
    data_stream:
      dataset: packet.flow
      type: logs
  - type: icmp
    interface:
      device: en1
      snaplen: 1514
      type: af_packet
      buffer_size_mb: 100
    procs:
      enabled: true
      monitored:
        - process: postgresql
          cmdline_grep: postgresql
    data_stream:
      dataset: packet.icmp
      type: logs
  - type: http
    interface:
      device: en2
      snaplen: 1514
      type: af_packet
      buffer_size_mb: 100
    procs:
      enabled: true
      monitored:
        - process: curl
          cmdline_grep: curl
    data_stream:
      dataset: packet.http
      type: logs
`)
	require.NoError(t, err)
	config, err := NewAgentConfig(cfg)
	require.NoError(t, err)

	require.Equal(t, config.Flows.Timeout, "10s")
	require.Equal(t, config.Flows.Index, "logs-packet.flow-default")
	require.Len(t, config.ProtocolsList, 2)

	var protocol map[string]interface{}
	require.NoError(t, config.ProtocolsList[0].Unpack(&protocol))
	require.Len(t, protocol["processors"].([]interface{}), 3)
	require.Equal(t, "default_route", config.Interfaces[0].Device)
	require.Equal(t, "en1", config.Interfaces[1].Device)
	require.Equal(t, "en2", config.Interfaces[2].Device)
	require.Len(t, config.Procs.Monitored, 3)
}
