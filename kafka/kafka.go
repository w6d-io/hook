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
package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"

	"github.com/avast/retry-go"

	"github.com/w6d-io/x/kafkax"
	"github.com/w6d-io/x/logx"
)

func (k *Kafka) Send(ctx context.Context, payload interface{}, URL *url.URL) error {

	log := logx.WithName(ctx, "Kafka.Send")

	passwd, ok := URL.User.Password()
	query := URL.Query()
	topic := query["topic"][0]

	k.BootstrapServer = URL.Host
	k.Username = URL.User.Username()
	k.Password = passwd

	var (
		async      bool
		messageKey string
		protocol   = "SASL_SSL"
		mechanisms = "PLAIN"
	)
	async = len(query["async"]) > 0 && query["async"][0] == "true"
	if len(query["messagekey"]) > 0 {
		messageKey = query["messagekey"][0]
	}
	if len(query["protocol"]) > 0 {
		protocol = query["protocol"][0]
	}
	if len(query["mechanisms"]) > 0 {
		mechanisms = query["mechanisms"][0]
	}

	p, err := k.NewProducer(
		kafkax.AuthKafka(ok),
		kafkax.Async(async),
		kafkax.Protocol(protocol),
		kafkax.Mechanisms(mechanisms),
	)

	if err != nil {
		log.Error(err, "error while creating producer")
		return err
	}

	message, err := json.Marshal(&payload)
	if err != nil {
		log.Error(err, "marshal failed")
		return err
	}

	if err := retry.Do(
		func() error {
			if err := p.SetTopic(topic).Produce(messageKey, message); err != nil {
				return err
			}
			return nil
		},
		retry.Attempts(5),
	); err != nil {
		return err
	}
	//log.V(1).Info("send payload by kafka", "payload", payload,
	//	"address", URL.Host)
	return nil
}

func (k *Kafka) Validate(URL *url.URL) error {

	log := logx.WithName(context.TODO(), "Kafka.Validate")

	if URL == nil {
		return nil
	}
	values := URL.Query()
	if _, ok := values["topic"]; !ok {
		log.Error(errors.New("missing topic"), URL.Redacted())
		return errors.New("missing topic")
	}
	return nil
}
