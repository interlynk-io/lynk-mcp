// Copyright 2025 Interlynk.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"context"
	"fmt"

	"github.com/interlynk-io/lynk-mcp/internal/graphql"
)

// VexScope says how a version relates to a component vulnerability for a VEX update.
type VexScope string

const (
	// VexScopeDirect means the version owns the vulnerability, so an update changes it directly.
	VexScopeDirect VexScope = "direct"
	// VexScopeParent means the version is an ancestor of the owner, so an update writes a parent-level disposition.
	VexScopeParent VexScope = "parent"
	// VexScopeNotVisible means the vulnerability is not visible from the version.
	VexScopeNotVisible VexScope = "not_visible"
)

// VexScopeResult is the scope of one component vulnerability, with the version that owns it.
type VexScopeResult struct {
	Scope          VexScope
	OwnerVersionID string
}

// ResolveVexScope classifies each component vulnerability as direct, parent, or not
// visible from the given version. Every requested ID has an entry in the result.
func (c *Client) ResolveVexScope(ctx context.Context, versionID string, componentVulnIDs []string) (map[string]VexScopeResult, error) {
	scopes := make(map[string]VexScopeResult, len(componentVulnIDs))
	if len(componentVulnIDs) == 0 {
		return scopes, nil
	}

	vars := map[string]interface{}{
		"sbomIds":          []string{versionID},
		"componentVulnIds": componentVulnIDs,
	}

	var result struct {
		SbomsCompVuln []struct {
			ID     string `json:"id"`
			IsPart bool   `json:"isPart"`
			SbomID string `json:"sbomId"`
		} `json:"sbomsCompVuln"`
	}

	if err := c.gql.Execute(ctx, graphql.SbomsCompVulnQuery, vars, &result); err != nil {
		return nil, err
	}

	for _, id := range componentVulnIDs {
		scopes[id] = VexScopeResult{Scope: VexScopeNotVisible}
	}
	for _, row := range result.SbomsCompVuln {
		if _, requested := scopes[row.ID]; !requested {
			continue
		}
		// A vulnerability seen both directly and through a part is owned by this version.
		if existing := scopes[row.ID]; existing.Scope == VexScopeDirect {
			continue
		}
		scope := VexScopeParent
		if !row.IsPart {
			scope = VexScopeDirect
		}
		scopes[row.ID] = VexScopeResult{Scope: scope, OwnerVersionID: row.SbomID}
	}

	return scopes, nil
}

// VexNotVisibleError means a component vulnerability is not visible from the requested version.
// The cause is ambiguous: the ID may be wrong, the vulnerability may be retracted or hidden, or
// the version may be neither the owner nor an ancestor of the owner.
type VexNotVisibleError struct {
	ComponentVulnID string
	VersionID       string
}

func (e *VexNotVisibleError) Error() string {
	return fmt.Sprintf("vulnerability %s is not visible from version %s: the ID may be wrong, the vulnerability may be retracted or hidden, or the version is neither the owner nor an ancestor of the owner. Use get_version_relations to check the version links.", e.ComponentVulnID, e.VersionID)
}
