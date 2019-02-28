/*
Copyright (c) 2019 TriggerMesh, Inc
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseData(t *testing.T) {
	data := "eyJzdGF0dXMiOiJXT1JLSU5HIiwic291cmNlIjp7InJlcG9Tb3VyY2UiOnsicHJvamVjdElkIjoiZm9vIiwicmVwb05hbWUiOiJnaXRodWItYmFyIiwiYnJhbmNoTmFtZSI6Im1hc3RlciJ9fX0="
	payload, err := parseData(data)
	assert.NoError(t, err)
	assert.Equal(t, "foo", payload.Source.RepoSource.ProjectID)
	assert.Equal(t, "WORKING", payload.Status)
	assert.Equal(t, "github-bar", payload.Source.RepoSource.RepoName)
	assert.Equal(t, "master", payload.Source.RepoSource.BranchName)
}
