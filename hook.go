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
Created on 07/02/2021
*/

package hook

import (
	"context"
	"fmt"
	"net/url"
	"regexp"

	"github.com/w6d-io/hook/http"
	"github.com/w6d-io/hook/kafka"

	"github.com/w6d-io/x/logx"
)

func init() {
	AddProvider("kafka", &kafka.Kafka{})
	AddProvider("http", &http.HTTP{})
	AddProvider("https", &http.HTTP{})
}

func Send(ctx context.Context, payload interface{}, scope string) error {
	log := logx.WithName(ctx, "Hook.Send")
	log.V(1).Info("to send", "payload", payload)
	go func(ctx context.Context, payload interface{}) {
		if err := DoSend(ctx, payload, scope); err != nil {
			log.Error(err, "DoSend")
			return
		}
	}(ctx, payload)
	return nil
}

// DoSend loops into all the subscribers url. for each it get the function by the scheme and run the method/function associated
func DoSend(ctx context.Context, payload interface{}, scope string) error {
	log := logx.WithName(ctx, "Hook.DoSend")
	errc := make(chan error, len(subscribers))
	quit := make(chan struct{})
	defer close(quit)

	for _, sub := range subscribers {
		log := log.WithValues("scheme", sub.URL.Scheme)
		if !isInScope(sub, scope) {
			log.V(1).Info("skip", "sub", sub.URL.String())
			continue
		}
		go func(payload interface{}, URL *url.URL) {
			f := suppliers[URL.Scheme]
			logg := log.WithValues("url", URL)
			select {
			case errc <- f.Send(ctx, payload, URL):
				logg.Info("sent")
			case <-quit:
				logg.Info("quit")
			}
		}(payload, sub.URL)
	}
	for range subscribers {
		if err := <-errc; err != nil {
			log.Error(err, "Sent failed")
			return err
		}
	}

	return nil
}

// AddProvider adds the protocol Send function to the suppliers list
func AddProvider(name string, i Interface) {
	suppliers[name] = i
}

// DelProvider adds the protocol Send function to the suppliers list
// func DelProvider(name string) {
//     delete(suppliers, name)
// }

// Subscribe recorder the suppliers and its scope in subscribers
func Subscribe(URLRaw, scope string) error {

	log := logx.WithName(context.TODO(), "Hook.Subscribe")

	URL, err := url.Parse(URLRaw)
	if err != nil {
		log.Error(err, "URL parsing", "url", URLRaw)
		return err
	}
	s, ok := suppliers[URL.Scheme]
	if !ok {
		err := fmt.Errorf("provider %v not supported", URL.Scheme)
		log.Error(err, "check provider")
		return err
	}

	if err := s.Validate(URL); err != nil {
		log.Error(err, "validation failed")
		return err
	}
	w := subscriber{
		URL:   URL,
		Scope: scope,
	}
	subscribers = append(subscribers, w)
	return nil
}

// CleanSubscriber cleans the list of subscriber
func CleanSubscriber() {
	subscribers = []subscriber{}
}

func isInScope(s subscriber, scope string) bool {

	log := logx.WithName(context.TODO(), "Hook.isInScope")

	prefix := ""
	if s.Scope == "*" {
		prefix = "."
	}
	r, err := regexp.Compile(prefix + s.Scope)
	if err != nil {
		log.Error(err, "Match failed")
		return false
	}
	return r.MatchString(scope)
}
