// Copyright 2018 The zap-encoder Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stackdriver_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/zchee/zap-encoder/stackdriver"
)

const (
	testProjectID = "testProjectID"
	testLogID     = "testLogID"
)

// TestJSONEncodeEntry is an more "integrated" test that makes it easier to get
// coverage on the json encoder (e.g. putJSONEncoder, let alone EncodeEntry
// itself) than the tests in json_encoder_impl_test.go; it needs to be in the
// zapcore_test package, so that it can import the toplevel zap package for
// field constructor convenience.
func TestStackdriverEncodeEntry(t *testing.T) {
	type bar struct {
		Key string  `json:"key"`
		Val float64 `json:"val"`
	}

	type foo struct {
		A string  `json:"aee"`
		B int     `json:"bee"`
		C float64 `json:"cee"`
		D []bar   `json:"dee"`
	}

	tests := []struct {
		name     string
		expected string
		ent      zapcore.Entry
		fields   []zapcore.Field
	}{
		{
			name: "info entry with some fields",
			expected: `{
				"eventTime": "2018-06-19T16:33:42.000Z",
				"severity": "info",
				"logger": "bob",
				"message": "lob law",
				"so": "passes",
				"answer": 42,
				"common_pie": 3.14,
				"such": {
					"aee": "lol",
					"bee": 123,
					"cee": 0.9999,
					"dee": [
						{"key": "pi", "val": 3.141592653589793},
						{"key": "tau", "val": 6.283185307179586}
					]
				}
			}`,
			ent: zapcore.Entry{
				Level:      zapcore.InfoLevel,
				Time:       time.Date(2018, 6, 19, 16, 33, 42, 99, time.UTC),
				LoggerName: "bob",
				Message:    "lob law",
			},
			fields: []zapcore.Field{
				zap.String("so", "passes"),
				zap.Int("answer", 42),
				zap.Float64("common_pie", 3.14),
				zap.Reflect("such", foo{
					A: "lol",
					B: 123,
					C: 0.9999,
					D: []bar{
						{"pi", 3.141592653589793},
						{"tau", 6.283185307179586},
					},
				}),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := stackdriver.NewStackdriverEncoderConfig()
			enc := stackdriver.NewStackdriverEncoder(context.Background(), cfg, testProjectID, testLogID)
			buf, err := enc.EncodeEntry(tt.ent, tt.fields)
			if err != nil {
				t.Errorf("Unexpected JSON encoding error: %+v", err)
				return
			}

			var expectedJSONAsInterface, actualJSONAsInterface interface{}
			if err := json.Unmarshal([]byte(tt.expected), &expectedJSONAsInterface); err != nil {
				t.Errorf(fmt.Sprintf("Expected value (%q) is not valid json.\nJSON parsing error: %+v", tt.expected, err))
				return
			}
			if err := json.Unmarshal([]byte(buf.String()), &actualJSONAsInterface); err != nil {
				t.Errorf(fmt.Sprintf("Actual value (%q) is not valid json.\nJSON parsing error: %+v", buf.String(), err))
				return
			}

			if diff := cmp.Diff(&expectedJSONAsInterface, &actualJSONAsInterface); diff != "" {
				t.Errorf("%s: Incorrect encoded JSON entry: (-got, +want)\n%s\n", tt.name, diff)
			}
			// if !cmp.Equal(expectedJSONAsInterface, actualJSONAsInterface) {
			// 	t.Error("Incorrect encoded JSON entry")
			// }
			// buf.Free()
		})
	}
}