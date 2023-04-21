// Copyright 2017 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package ctxtags

import "net/http"

type RequestFieldExtractorFunc func(req *http.Request) map[string]interface{}
