// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/googlecloudplatform/ezcx"
	ezfs "github.com/yaq-cc/ezcx-cache/firestore"
)

var (
	Port       string = os.Getenv("PORT")
	ProjectID  string = os.Getenv("PROJECT_ID")
	Collection string = os.Getenv("COLLECTION")
	Document   string = os.Getenv("DOCUMENT")
)

func main() {
	ctx := context.Background()
	lg := log.Default()
	cache, close := ezfs.New[string](ctx, &ezfs.FirestoreConfig{
		ProjectID:  ProjectID,
		Collection: Collection,
		Document:   Document,
	})
	defer close()
	server := ezcx.NewServer(ctx, ":"+Port, lg)
	server.HandleCx("/route", cxRoute(cache))
	cache.Listen(ctx)
	server.ListenAndServe(ctx)
}

func cxRoute(c *ezfs.FirestoreCache[string]) ezcx.HandlerFunc {
	return func(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
		lg := req.Logger()
		params := req.GetSessionParameters()
		_lawyer, ok := params["lawyer_name"]
		if !ok {
			return errors.New("missing required parameter: lawyer_name")
		}
		lawyer, ok := _lawyer.(string)
		if !ok {
			return errors.New("lawyer_name could not be converted to string")
		}
		phNum, ok := c.Get(lawyer)
		if !ok {
			lg.Println(c.Get(lawyer))
			return errors.New("lawyer_name was not found in the cache")
		}
		log.Printf("%s : %s", lawyer, phNum)
		// Add a telephony transfer call response
		res.AddTelephonyTransferResponse(phNum)
		return nil
	}
}
