//go:build !integration

/*
Copyright 2020 WILDCARD

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 08/02/2021
*/

package hook_test

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/w6d-io/hook"

	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	zapraw "go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, " Suite")
}

var _ = BeforeSuite(func() {
	encoder := zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	opts := zap.Options{
		Encoder:         zapcore.NewConsoleEncoder(encoder),
		Development:     true,
		StacktraceLevel: zapcore.PanicLevel,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts), zap.RawZapOpts(zapraw.AddCaller())))
})

var _ = AfterSuite(func() {
	hook.CleanSubscriber()
})

type TestAllOk struct{}

func (t *TestAllOk) Init(_ context.Context, _ *url.URL) error                { return nil }
func (t *TestAllOk) Validate(_ *url.URL) error                               { return nil }
func (t *TestAllOk) Send(_ context.Context, _ interface{}, _ *url.URL) error { return nil }

type TestSendFail struct{}

func (t *TestSendFail) Init(_ context.Context, _ *url.URL) error { return nil }
func (t *TestSendFail) Send(_ context.Context, _ interface{}, _ *url.URL) error {
	return errors.New("send failed")
}
func (t *TestSendFail) Validate(_ *url.URL) error { return nil }

type TestValidateFail struct{}

func (t *TestValidateFail) Init(_ context.Context, _ *url.URL) error                { return nil }
func (t *TestValidateFail) Validate(_ *url.URL) error                               { return errors.New("validate failed") }
func (t *TestValidateFail) Send(_ context.Context, _ interface{}, _ *url.URL) error { return nil }
