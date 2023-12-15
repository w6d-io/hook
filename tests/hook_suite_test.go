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

package test

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/w6d-io/x/logx"
	"os"
	"os/exec"
	"regexp"
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

var (
	rpkPath   = rpkPathFinder()
	rpkPort   string
	ctx       context.Context
	kafkaHost string
)
var _ = BeforeSuite(func() {
	correlationID := uuid.New().String()
	ctx = context.Background()
	ctx = context.WithValue(ctx, logx.CorrelationID, correlationID)

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
	Expect(rpkPath).ToNot(BeEmpty())
	By(rpkPath)

	By("start redpanda")
	ps := exec.Command(rpkPath, "container", "start")
	ps.Stderr = os.Stderr
	out, err := ps.Output()
	Expect(err).ToNot(HaveOccurred())
	logx.WithName(ctx, "before_suite").Info("start RedPanda cluster", "out", out)

	By("get redpanda port")
	re := regexp.MustCompile(`.*127\.0\.0\.1:(\d+).*`)
	match := re.FindStringSubmatch(string(out))
	Expect(len(match)).To(Equal(2))
	rpkPort = match[1]
	kafkaHost = fmt.Sprintf("localhost:%s", rpkPort)

	By("set cluster config")
	ps = exec.Command("docker", "exec", "rp-node-0", "rpk", "cluster", "config", "set", "auto_create_topics_enabled", `true`)
	//ps = exec.Command(rpkPath, "cluster", "config", "set", "auto_create_topics_enabled", "true", "--brokers", kafkaHost)
	ps.Stderr = os.Stderr
	out, err = ps.Output()
	logx.WithName(ctx, "before_suite").Info("set RedPanda cluster config", "out", out)
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	if rpkPath != "" {
		ps := exec.Command(rpkPath, "container", "purge")
		ps.Stdout = os.Stdout
		ps.Stderr = os.Stderr
		Expect(ps.Start()).To(Succeed())
	}
	hook.CleanSubscriber()
})

func rpkPathFinder() string {
	if val, ok := os.LookupEnv("RPK"); ok {
		return val
	}
	return ""
}
