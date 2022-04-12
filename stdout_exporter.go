// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stdoutexporter // import "github.com/fsolleza/stdoutexporter"

import (
	"context"
	"io"
	"os"
	"sync"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/model/otlp"
	"go.opentelemetry.io/collector/model/pdata"
)

// Marshaler configuration used for marhsaling Protobuf to JSON.
var tracesMarshaler = otlp.NewJSONTracesMarshaler()
var metricsMarshaler = otlp.NewJSONMetricsMarshaler()
var logsMarshaler = otlp.NewJSONLogsMarshaler()

// fileExporter is the implementation of file exporter that writes telemetry data to a file
// in Protobuf-JSON format.
type stdoutExporter struct {
	path  string
	file  io.WriteCloser
	mutex sync.Mutex
}

func (e *stdoutExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (e *stdoutExporter) ConsumeTraces(_ context.Context, td pdata.Traces) error {
	buf, err := tracesMarshaler.MarshalTraces(td)
	if err != nil {
		return err
	}
	return exportMessageAsLine(e, buf)
}

func (e *stdoutExporter) ConsumeMetrics(_ context.Context, md pdata.Metrics) error {
	buf, err := metricsMarshaler.MarshalMetrics(md)
	if err != nil {
		return err
	}
	return exportMessageAsLine(e, buf)
}

func (e *stdoutExporter) ConsumeLogs(_ context.Context, ld pdata.Logs) error {
	buf, err := logsMarshaler.MarshalLogs(ld)
	if err != nil {
		return err
	}
	return exportMessageAsLine(e, buf)
}

func exportMessageAsLine(e *stdoutExporter, buf []byte) error {
	// Ensure only one write operation happens at a time.
	e.mutex.Lock()
	defer e.mutex.Unlock()
	fmt.Println(buf)
	return nil
}

func (e *stdoutExporter) Start(context.Context, component.Host) error {
	return nil
}

// Shutdown stops the exporter and is invoked during shutdown.
func (e *stdoutExporter) Shutdown(context.Context) error {
	return nil
}
