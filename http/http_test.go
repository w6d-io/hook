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
Created on 26/02/2021
*/
package http_test

import (
	"context"
	"net/url"

	//"os"

	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/ginkgo/tmp"
	. "github.com/onsi/gomega"

	//"github.com/onsi/gomega/ghttp"
	"github.com/w6d-io/hook/http"
)

var _ = Describe("HTTP", func() {
	Context("Validate", func() {
		It("nothing to do", func() {
			h := &http.HTTP{}
			URL, err := url.Parse("http://localhost:1234")
			err = h.Validate(URL)
			Expect(err).To(Succeed())
		})
	})
	Context("Send", func() {
		It("init success", func() {
			h := http.HTTP{}
			url, _ := url.Parse("http://localhost:1234")
			err := h.Init(context.Background(), url)
			Expect(err).To(Succeed())
		})
		It("all attempts failed", func() {
			h := http.HTTP{}
			URL, err := url.Parse("http://localhost:1234")
			Ω(err).To(Succeed())
			err = h.Send(context.Background(), "message", URL)
			Ω(err).ToNot(Succeed())
			Ω(err.Error()).To(ContainSubstring("All attempts fail"))
		})
		It("test timeout", func() {
			h := http.HTTP{}
			URL, err := url.Parse("http://localhost:1234?timeout=120")
			Ω(err).To(Succeed())
			err = h.Send(context.Background(), "message", URL)
			Ω(err).ToNot(Succeed())
			Ω(err.Error()).To(ContainSubstring("All attempts fail"))
		})
		It("test bad timeout", func() {
			h := http.HTTP{}
			URL, err := url.Parse("http://localhost:1234?timeout=s0")
			Ω(err).To(Succeed())
			err = h.Send(context.Background(), "message", URL)
			Ω(err).ToNot(Succeed())
			Ω(err.Error()).To(ContainSubstring("invalid syntax"))
		})
	})
})
