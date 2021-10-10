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
	"net/url"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/w6d-io/hook/kafka"
)

var _ = Describe("Kafka", func() {
	Context("Validate", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		It("url is nil", func() {
			k := kafka.Kafka{}

			err := k.Validate(nil)
			Expect(err).ToNot(HaveOccurred())
			//Expect(err.Error()).To(Equal("ddd"))
		})
		It("url is missed mandatory part", func() {
			k := kafka.Kafka{}
			url, _ := url.Parse("kafka://localhost:9092")
			err := k.Validate(url)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("missing topic"))
		})
		It("url is good", func() {
			k := kafka.Kafka{}
			url, _ := url.Parse("kafka://localhost:9092?topic=TEST")
			err := k.Validate(url)
			Expect(err).To(Succeed())
		})
	})
	Context("options", func() {
		It("bad configuration", func() {
			var opts []kafka.Option
			opts = append(opts, kafka.Protocol("SASL_SSL"))
			opts = append(opts, kafka.Mechanisms("PLAIN"))
			opts = append(opts, kafka.Async(false))
			opts = append(opts, kafka.WriteTimeout(1*time.Second))
			opts = append(opts, kafka.MaxWait(1*time.Second))
			opts = append(opts, kafka.StatInterval(3*time.Second))
			opts = append(opts, kafka.NumPartitions(10))
			opts = append(opts, kafka.ReplicationFactor(3))
			opts = append(opts, kafka.AuthKafka(true))
			opts = append(opts, kafka.FullStats(false))
			opts = append(opts, kafka.Debugs([]string{"deep"}))
			opts = append(opts, kafka.SessionTimeout(10*time.Millisecond))
			opts = append(opts, kafka.MaxPollInterval(10*time.Millisecond))
			opts = append(opts, kafka.GroupInstanceID("test"))
			opts = append(opts, kafka.ConfigMapKey("test"))
			_ = kafka.NewOptions(opts...)
			k := kafka.Kafka{
				Username:        "test",
				Password:        "test",
				BootstrapServer: "localhost:9092",
			}
			err := k.Producer("TEST_ID", "message", opts...)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(Equal("Local: Invalid argument or configuration"))
		})
		It("", func() {
			var opts []kafka.Option
			opts = append(opts, kafka.Protocol("SASL_SSL"))
			opts = append(opts, kafka.Mechanisms("PLAIN"))
			opts = append(opts, kafka.Async(false))
			opts = append(opts, kafka.WriteTimeout(1*time.Second))
			opts = append(opts, kafka.MaxWait(1*time.Second))
			opts = append(opts, kafka.StatInterval(3*time.Second))
			opts = append(opts, kafka.AuthKafka(true))
			opts = append(opts, kafka.SessionTimeout(10*time.Millisecond))
			opts = append(opts, kafka.MaxPollInterval(10*time.Millisecond))
			_ = kafka.NewOptions(opts...)
			k := kafka.Kafka{
				Topic:           "TEST",
				Username:        "test",
				Password:        "test",
				BootstrapServer: "localhost:9092",
			}
			err := k.Producer("TEST_ID", "message", opts...)
			Expect(err).To(Succeed())
		})
	})
	Context("Send", func() {
		It("timeout but succeed", func() {
			k := &kafka.Kafka{}
			url, _ := url.Parse("kafka://localhost:9092?topic=TEST&messagekey=TEST_KEY&protocol=TCP&mechanisms=PLAIN")
			err := k.Send("test", url)
			Expect(err).To(Succeed())
		})
	})
})
