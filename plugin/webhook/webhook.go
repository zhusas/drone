// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

// +build !oss

package webhook

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/drone/drone/core"

	"github.com/99designs/httpsignatures-go"
)

// required http headers
var headers = []string{
	"date",
	"digest",
}

var signer = httpsignatures.NewSigner(
	httpsignatures.AlgorithmHmacSha256,
	headers...,
)

// New returns a new Webhook sender.
func New(endpoints []string, secret string) core.WebhookSender {
	return &sender{
		Endpoints: endpoints,
		Secret:    secret,
	}
}

type sender struct {
	Client    *http.Client
	Endpoints []string
	Secret    string
}

// Send sends the JSON encoded webhook to the global
// HTTP endpoints.
func (s *sender) Send(ctx context.Context, payload *core.WebhookData) error {
	if len(s.Endpoints) == 0 {
		return nil
	}

	data, _ := json.Marshal(payload)
	for _, endpoint := range s.Endpoints {
		err := s.send(endpoint, s.Secret, payload.Event, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sender) send(endpoint, secret, event string, data []byte) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	buf := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", endpoint, buf)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	req.Header.Add("X-Drone-Event", event)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Digest", "SHA-256="+digest(data))
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))
	err = signer.SignRequest("hmac-key", s.Secret, req)
	res, err := s.client().Do(req)
	if res != nil {
		res.Body.Close()
	}
	return err
}

func (s *sender) client() *http.Client {
	if s.Client == nil {
		return http.DefaultClient
	}
	return s.Client
}

func digest(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
