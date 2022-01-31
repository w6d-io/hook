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
Created on 25/02/2021
*/
package hook_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/w6d-io/hook"
)

var _ = Describe("Hook", func() {
	When("provider works", func() {
		BeforeEach(func() {
			hook.AddProvider("http", &TestAllOk{})
		})
		Context("add suppliers", func() {
			It("succeed for http", func() {
				err := hook.Subscribe(context.Background(), "http://localhost:8080", "*")
				Expect(err).To(Succeed())
			})
			It("does not support this provider", func() {
				err := hook.Subscribe(context.Background(), "mongodb://login:password@localhost:27017", "end")
				Expect(err).ToNot(Succeed())
				Expect(err.Error()).To(ContainSubstring("not supported"))
			})
			It("URL is malformed", func() {
				err := hook.Subscribe(context.Background(), "test://{}", "begin")
				Expect(err).ToNot(Succeed())
				Expect(err.Error()).To(ContainSubstring("invalid character"))
			})
			It("validation failed", func() {
				hook.AddProvider("https", &TestValidateFail{})
				err := hook.Subscribe(context.Background(), "https://localhost", "*")
				Expect(err).ToNot(Succeed())
				Expect(err.Error()).To(Equal("validate failed"))
			})
			It("init failed", func() {
				err := hook.Subscribe(context.Background(), "kafka://user:pass@localhost:9092?topic=TEST&protocol=Unknown", ".*")
				Expect(err).NotTo(Succeed())
			})
		})
		Context("send a payload", func() {
			It("succeed to send", func() {
				Expect(hook.Send(context.Background(), "message", "*")).To(Succeed())
			})
			It("succeed to DoSend", func() {
				Expect(hook.DoSend(context.Background(), "message", "*")).To(Succeed())
			})
		})
	})
	When("provider Failed", func() {
		BeforeEach(func() {
			hook.AddProvider("http", &TestSendFail{})
		})
		Context("send payload", func() {
			It("with Send", func() {
				err := hook.Send(context.Background(), "message", "*")
				Expect(err).To(Succeed())
			})
			It("with bad scope", func() {
				err := hook.Subscribe(context.Background(), "http://localhost", "begin")
				Expect(err).To(Succeed())
				err = hook.Send(context.Background(), "message", "test")
				Expect(err).To(Succeed())
			})
			It("regex failed", func() {
				err := hook.Subscribe(context.Background(), "http://localhost", "[")
				Expect(err).To(Succeed())
				err = hook.DoSend(context.Background(), "message", "test")
				Expect(err).ToNot(Succeed())
				Expect(err.Error()).To(Equal("send failed"))
			})
			It("with DoSend", func() {
				err := hook.DoSend(context.Background(), "message", "*")
				Expect(err).ToNot(Succeed())
				Expect(err.Error()).To(Equal("send failed"))
			})
			It("skip", func() {

			})
		})
	})
})
