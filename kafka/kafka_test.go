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
Created on 26/02/2021
*/

package kafka_test

import (
	"context"
	"errors"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/w6d-io/hook/kafka"
	"github.com/w6d-io/x/kafkax"
)

var _ = Describe("Kafka", func() {
	Context("Validate", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		It("url is nil", func() {
			k := &kafka.Kafka{
				Producer: &kafkax.Producer{
					ClientProducerAPI: &kafkax.MockClientProducer{},
				},
			}
			err := k.Validate(nil)
			Expect(err).ToNot(HaveOccurred())
			//Expect(err.Error()).To(Equal("ddd"))
		})
		It("url is missed mandatory part", func() {
			k := &kafka.Kafka{
				Producer: &kafkax.Producer{
					ClientProducerAPI: &kafkax.MockClientProducer{},
				},
			}
			url, _ := url.Parse("kafka://localhost:9092")
			err := k.Validate(url)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("missing topic"))
		})
		It("url is good", func() {
			k := &kafka.Kafka{
				Producer: &kafkax.Producer{
					ClientProducerAPI: &kafkax.MockClientProducer{},
				},
			}
			url, _ := url.Parse("kafka://localhost:9092?topic=TEST")
			err := k.Validate(url)
			Expect(err).To(Succeed())
		})
	})
	Context("Send", func() {
		It("init success", func() {
			k := &kafka.Kafka{
				Producer: &kafkax.Producer{
					ClientProducerAPI: &kafkax.MockClientProducer{},
				},
			}
			url, _ := url.Parse("kafka://localhost:9092?topic=TEST")
			err := k.Init(context.Background(), url)
			Expect(err).To(Succeed())
		})
		It("wrong option", func() {
			k := &kafka.Kafka{
				Producer: &kafkax.Producer{
					ClientProducerAPI: &kafkax.MockClientProducer{},
				},
			}
			url, _ := url.Parse("kafka://user:pass@localhost:9092?topic=TEST&mechanisms=PLAIN&protocol=Unknown")
			err := k.Init(context.Background(), url)
			Expect(err).NotTo(Succeed())
		})
		It("wrong payload", func() {
			k := &kafka.Kafka{
				Producer: &kafkax.Producer{
					ClientProducerAPI: &kafkax.MockClientProducer{},
				},
			}
			url, _ := url.Parse("kafka://localhost:9092?topic=TEST&messagekey=TEST_KEY")
			err := k.Send(context.Background(), make(chan int), url)
			Expect(err).NotTo(Succeed())
		})
		It("timeout but succeed", func() {
			k := &kafka.Kafka{
				Producer: &kafkax.Producer{
					ClientProducerAPI: &kafkax.MockClientProducer{},
				},
			}
			url, _ := url.Parse("kafka://localhost:9092?topic=TEST&messagekey=TEST_KEY")
			err := k.Send(context.Background(), "test", url)
			Expect(err).To(Succeed())
		})
		It("error while producing", func() {
			k := &kafka.Kafka{
				Producer: &kafkax.Producer{
					ClientProducerAPI: &kafkax.MockClientProducer{
						ErrProduce: errors.New("error while producing"),
					},
				},
			}
			url, _ := url.Parse("kafka://localhost:9092?topic=TEST&messagekey=TEST_KEY")
			err := k.Send(context.Background(), "test", url)
			Expect(err).NotTo(Succeed())
		})
	})
})
