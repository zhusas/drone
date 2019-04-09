// Copyright 2019 Drone IO, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repo

import (
	"github.com/drone/drone/core"
	"github.com/drone/go-scm/scm"
)

// convertRepository is a helper function that converts a
// repository from the source code management system to the
// local datastructure.
func convertRepository(src *scm.Repository) *core.Repository {
	return &core.Repository{
		UID:        src.ID,
		Namespace:  src.Namespace,
		Name:       src.Name,
		Slug:       scm.Join(src.Namespace, src.Name),
		HTTPURL:    src.Clone,
		SSHURL:     src.CloneSSH,
		Link:       src.Link,
		Private:    src.Private,
		Visibility: convertVisibility(src),
		Branch:     src.Branch,
	}
}

// convertVisibility is a helper function that returns the
// repository visibliity based on the privacy flag.
func convertVisibility(src *scm.Repository) string {
	switch {
	case src.Private == true:
		return core.VisibilityPrivate
	default:
		return core.VisibilityPublic
	}
}
