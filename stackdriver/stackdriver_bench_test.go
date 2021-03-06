// Copyright 2018 The zap-encoder Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stackdriver_test

import (
	"context"
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/zchee/zap-encoder/internal/testutil"
	"github.com/zchee/zap-encoder/internal/uid"
	"github.com/zchee/zap-encoder/stackdriver"
)

func BenchmarkStackdriverEncoderLogMarshalerFunc(b *testing.B) {
	ctx := context.Background()
	testProjectID := testutil.ProjectID()
	uids := uid.NewSpace(testLogIDPrefix, nil)
	testLogID := uids.New()

	lg := stackdriver.NewDefaultStackdriverClient(ctx, testProjectID, testLogID)
	enc := stackdriver.NewStackdriverEncoder(ctx, lg, stackdriver.NewStackdriverEncoderConfig())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		enc.AddObject("nested", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			enc.AddInt64("i", int64(i))
			return nil
		}))
	}
}
