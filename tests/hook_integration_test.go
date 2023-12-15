//go:build integration

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
Created on 19/11/2022
*/

package test

import (
	"fmt"
	"github.com/w6d-io/x/logx"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/w6d-io/hook"
)

var _ = Describe("", func() {
	Context("Test hook with RedPanda", func() {
		var payload = struct {
			Data string
		}{
			Data: "test",
		}
		var topic = "int-test"
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		It("produces successfully", func() {
			By("init kafka subscription", func() {
				Expect(hook.Subscribe(ctx, fmt.Sprintf("kafka://%s?topic=%s&protocol=SASL_SSL&mechanisms=PLAIN&messagekey=test", kafkaHost, topic), "*")).To(Succeed())
			})

			By("Send payload")
			err := hook.Send(ctx, &payload, "send succeed")
			Expect(err).ToNot(HaveOccurred())
			By("consume the message", func() {
				ps := exec.Command(rpkPath, "topic", "consume", topic, "--num", "1", "--brokers", kafkaHost)
				out, err := ps.Output()
				Expect(err).ToNot(HaveOccurred())
				logx.WithName(ctx, "int-test").Info("test", "out", out)
			})

		})
		It("misses the topic", func() {
			Expect(hook.Subscribe(ctx, fmt.Sprintf("kafka://%s", kafkaHost), "*")).ToNot(Succeed())
		})

	})
})
