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
    "errors"
    "net/url"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "github.com/w6d-io/hook"
    "k8s.io/klog/klogr"
)

type TestAllOk struct {}
func (t *TestAllOk) Send(_ interface{}, _ *url.URL) error {return nil}
func (t *TestAllOk) Validate(_ *url.URL) error                      {return nil}

type TestSendFail struct {}
func (t *TestSendFail) Send(_ interface{}, _ *url.URL) error {return errors.New("send failed")}
func (t *TestSendFail) Validate(_ *url.URL) error                      {return nil}

type TestValidateFail struct {}
func (t *TestValidateFail) Send(_ interface{}, _ *url.URL) error {return nil}
func (t *TestValidateFail) Validate(_ *url.URL) error                      {return errors.New("validate failed")}

var _ = Describe("Hook", func() {
    When("provider works", func() {
        BeforeEach(func() {
            hook.AddProvider("http", &TestAllOk{})
        })
        Context("add suppliers", func() {
            AfterEach(func() {
                //hook.DelProvider("http")
            })
            It("succeed for http", func() {
                err := hook.Subscribe("http://localhost:8080", "*")
                Expect(err).To(Succeed())
            })
            It("does not support this provider", func() {
                err := hook.Subscribe("mongodb://login:password@localhost:27017", "end")
                Expect(err).ToNot(Succeed())
                Expect(err.Error()).To(ContainSubstring("not supported"))
            })
            It("URL is malformed", func() {
                err := hook.Subscribe("test://{}", "begin")
                Expect(err).ToNot(Succeed())
                Expect(err.Error()).To(ContainSubstring("invalid character"))
            })
            It("validate failed", func() {
                hook.AddProvider("https", &TestValidateFail{})
                err := hook.Subscribe("https://localhost", "*")
                Expect(err).ToNot(Succeed())
                Expect(err.Error()).To(Equal("validate failed"))
            })
        })
        Context("send a payload", func() {
            It("succeed to send", func() {
                log := klogr.New()
                Expect(hook.Send("message", log, "*")).To(Succeed())
            })
            It("succeed to DoSend", func() {
                log := klogr.New()
                Expect(hook.DoSend("message", log, "*")).To(Succeed())
            })
        })
    })
    When("provider Failed", func() {
        BeforeEach(func() {
            hook.AddProvider("http", &TestSendFail{})
        })
        Context("send payload", func() {
            It("with Send", func() {
                log := klogr.New()
                err := hook.Send("message", log, "*")
                Expect(err).To(Succeed())
            })
            It("with bad scope", func() {
                log := klogr.New()
                err := hook.Subscribe("http://localhost", "begin")
                Expect(err).To(Succeed())
                err = hook.Send("message", log, "test")
                Expect(err).To(Succeed())
            })
            It("regex failed", func() {
                log := klogr.New()
                err := hook.Subscribe("http://localhost", "[")
                Expect(err).To(Succeed())
                err = hook.DoSend("message", log, "test")
                Expect(err).ToNot(Succeed())
                Expect(err.Error()).To(Equal("send failed"))
            })
            It("with DoSend", func() {
                log := klogr.New()
                err := hook.DoSend("message", log, "*")
                Expect(err).ToNot(Succeed())
                Expect(err.Error()).To(Equal("send failed"))
            })
            It("skip", func() {

            })
        })
    })
})
