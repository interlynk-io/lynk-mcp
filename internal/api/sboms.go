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
	"time"

	"github.com/interlynk-io/lynk-mcp/internal/graphql"
)

// Sbom represents an SBOM document
type Sbom struct {
	ID             string
	ProjectVersion string
	Spec           string
	SpecVersion    string
	Format         string
	Lifecycle      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ProjectID      string
	Stats          *SbomStats
	Project        *Project
}

// SbomStats contains statistics for an SBOM
type SbomStats struct {
	CompCount         int
	CompPurlCount     int
	CompCpeCount      int
	CompLicenseCount  int
	CompSupplierCount int
	VulnStats         map[string]interface{}
}

// SbomComponent represents a component in an SBOM
type SbomComponent struct {
	ID          string
	Name        string
	Version     string
	Kind        string
	Purl        string
	Cpes        []string
	LicensesExp string
	Group       string
	Description string
	Primary     bool
	Internal    bool
	SbomID      string
	UpdatedAt   time.Time
	Sbom        *Sbom
}

// SbomsResult represents the result of listing SBOMs
type SbomsResult struct {
	Sboms       []Sbom
	TotalCount  int
	HasNextPage bool
	EndCursor   string
}

// ComponentsResult represents the result of listing components
type ComponentsResult struct {
	Components  []SbomComponent
	TotalCount  int
	HasNextPage bool
	EndCursor   string
}

// SbomDiff represents a diff entry between two SBOMs
type SbomDiff struct {
	DiffType           string
	DiffTags           []string
	SubjectComponent   *SbomComponent
	SubjectComponentID string
	TargetComponent    *SbomComponent
	TargetComponentID  string
}

// ListSbomsInput contains parameters for listing SBOMs
type ListSbomsInput struct {
	ProjectID string
	First     int
	After     string
	Lifecycle []string
}

// ListSboms fetches SBOMs for a project
func (c *Client) ListSboms(ctx context.Context, input ListSbomsInput) (*SbomsResult, error) {
	vars := map[string]interface{}{
		"projectId": input.ProjectID,
	}
	if input.First > 0 {
		vars["first"] = input.First
	} else {
		vars["first"] = 20
	}
	if input.After != "" {
		vars["after"] = input.After
	}
	if len(input.Lifecycle) > 0 {
		vars["lifestage"] = input.Lifecycle
	}

	var result struct {
		Project struct {
			SbomVersions struct {
				Nodes []struct {
					ID             string    `json:"id"`
					ProjectVersion string    `json:"projectVersion"`
					Spec           string    `json:"spec"`
					SpecVersion    string    `json:"specVersion"`
					Format         string    `json:"format"`
					Lifecycle      string    `json:"lifecycle"`
					CreatedAt      time.Time `json:"createdAt"`
					UpdatedAt      time.Time `json:"updatedAt"`
					ProjectID      string    `json:"projectId"`
					Stats          *struct {
						CompCount         int                    `json:"compCount"`
						CompPurlCount     int                    `json:"compPurlCount"`
						CompCpeCount      int                    `json:"compCpeCount"`
						CompLicenseCount  int                    `json:"compLicenseCount"`
						CompSupplierCount int                    `json:"compSupplierCount"`
						VulnStats         map[string]interface{} `json:"vulnStats"`
					} `json:"stats"`
				} `json:"nodes"`
				TotalCount int `json:"totalCount"`
				PageInfo   struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"sbomVersions"`
		} `json:"project"`
	}

	if err := c.gql.Execute(ctx, graphql.ProjectSbomsQuery, vars, &result); err != nil {
		return nil, err
	}

	sboms := make([]Sbom, len(result.Project.SbomVersions.Nodes))
	for i, n := range result.Project.SbomVersions.Nodes {
		var stats *SbomStats
		if n.Stats != nil {
			stats = &SbomStats{
				CompCount:         n.Stats.CompCount,
				CompPurlCount:     n.Stats.CompPurlCount,
				CompCpeCount:      n.Stats.CompCpeCount,
				CompLicenseCount:  n.Stats.CompLicenseCount,
				CompSupplierCount: n.Stats.CompSupplierCount,
				VulnStats:         n.Stats.VulnStats,
			}
		}
		sboms[i] = Sbom{
			ID:             n.ID,
			ProjectVersion: n.ProjectVersion,
			Spec:           n.Spec,
			SpecVersion:    n.SpecVersion,
			Format:         n.Format,
			Lifecycle:      n.Lifecycle,
			CreatedAt:      n.CreatedAt,
			UpdatedAt:      n.UpdatedAt,
			ProjectID:      n.ProjectID,
			Stats:          stats,
		}
	}

	return &SbomsResult{
		Sboms:       sboms,
		TotalCount:  result.Project.SbomVersions.TotalCount,
		HasNextPage: result.Project.SbomVersions.PageInfo.HasNextPage,
		EndCursor:   result.Project.SbomVersions.PageInfo.EndCursor,
	}, nil
}

// GetSbom fetches a single SBOM by ID
func (c *Client) GetSbom(ctx context.Context, id string) (*Sbom, error) {
	vars := map[string]interface{}{
		"sbomId": id,
	}

	var result struct {
		Sbom struct {
			ID             string    `json:"id"`
			ProjectVersion string    `json:"projectVersion"`
			Spec           string    `json:"spec"`
			SpecVersion    string    `json:"specVersion"`
			Format         string    `json:"format"`
			Lifecycle      string    `json:"lifecycle"`
			CreatedAt      time.Time `json:"createdAt"`
			UpdatedAt      time.Time `json:"updatedAt"`
			ProjectID      string    `json:"projectId"`
			Stats          *struct {
				CompCount         int                    `json:"compCount"`
				CompPurlCount     int                    `json:"compPurlCount"`
				CompCpeCount      int                    `json:"compCpeCount"`
				CompLicenseCount  int                    `json:"compLicenseCount"`
				CompSupplierCount int                    `json:"compSupplierCount"`
				VulnStats         map[string]interface{} `json:"vulnStats"`
			} `json:"stats"`
			Project struct {
				ID             string `json:"id"`
				Name           string `json:"name"`
				ProjectGroupID string `json:"projectGroupId"`
			} `json:"project"`
		} `json:"sbom"`
	}

	if err := c.gql.Execute(ctx, graphql.SbomQuery, vars, &result); err != nil {
		return nil, err
	}

	var stats *SbomStats
	if result.Sbom.Stats != nil {
		stats = &SbomStats{
			CompCount:         result.Sbom.Stats.CompCount,
			CompPurlCount:     result.Sbom.Stats.CompPurlCount,
			CompCpeCount:      result.Sbom.Stats.CompCpeCount,
			CompLicenseCount:  result.Sbom.Stats.CompLicenseCount,
			CompSupplierCount: result.Sbom.Stats.CompSupplierCount,
			VulnStats:         result.Sbom.Stats.VulnStats,
		}
	}

	return &Sbom{
		ID:             result.Sbom.ID,
		ProjectVersion: result.Sbom.ProjectVersion,
		Spec:           result.Sbom.Spec,
		SpecVersion:    result.Sbom.SpecVersion,
		Format:         result.Sbom.Format,
		Lifecycle:      result.Sbom.Lifecycle,
		CreatedAt:      result.Sbom.CreatedAt,
		UpdatedAt:      result.Sbom.UpdatedAt,
		ProjectID:      result.Sbom.ProjectID,
		Stats:          stats,
		Project: &Project{
			ID:             result.Sbom.Project.ID,
			Name:           result.Sbom.Project.Name,
			ProjectGroupID: result.Sbom.Project.ProjectGroupID,
		},
	}, nil
}

// ListComponentsInput contains parameters for listing components
type ListComponentsInput struct {
	SbomID string
	First  int
	After  string
	Search string
	Kind   []string
	Direct *bool
}

// ListComponents fetches components for an SBOM
func (c *Client) ListComponents(ctx context.Context, input ListComponentsInput) (*ComponentsResult, error) {
	vars := map[string]interface{}{
		"sbomId": input.SbomID,
	}
	if input.First > 0 {
		vars["first"] = input.First
	} else {
		vars["first"] = 50
	}
	if input.After != "" {
		vars["after"] = input.After
	}
	if input.Search != "" {
		vars["search"] = input.Search
	}
	if len(input.Kind) > 0 {
		vars["kind"] = input.Kind
	}
	if input.Direct != nil {
		vars["direct"] = *input.Direct
	}

	var result struct {
		Sbom struct {
			Components struct {
				Nodes []struct {
					ID          string    `json:"id"`
					Name        string    `json:"name"`
					Version     string    `json:"version"`
					Kind        string    `json:"kind"`
					Purl        string    `json:"purl"`
					Cpes        []string  `json:"cpes"`
					LicensesExp string    `json:"licensesExp"`
					Group       string    `json:"group"`
					Description string    `json:"description"`
					Primary     bool      `json:"primary"`
					Internal    bool      `json:"internal"`
					SbomID      string    `json:"sbomId"`
					UpdatedAt   time.Time `json:"updatedAt"`
				} `json:"nodes"`
				TotalCount int `json:"totalCount"`
				PageInfo   struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"components"`
		} `json:"sbom"`
	}

	if err := c.gql.Execute(ctx, graphql.SbomComponentsQuery, vars, &result); err != nil {
		return nil, err
	}

	components := make([]SbomComponent, len(result.Sbom.Components.Nodes))
	for i, n := range result.Sbom.Components.Nodes {
		components[i] = SbomComponent{
			ID:          n.ID,
			Name:        n.Name,
			Version:     n.Version,
			Kind:        n.Kind,
			Purl:        n.Purl,
			Cpes:        n.Cpes,
			LicensesExp: n.LicensesExp,
			Group:       n.Group,
			Description: n.Description,
			Primary:     n.Primary,
			Internal:    n.Internal,
			SbomID:      n.SbomID,
			UpdatedAt:   n.UpdatedAt,
		}
	}

	return &ComponentsResult{
		Components:  components,
		TotalCount:  result.Sbom.Components.TotalCount,
		HasNextPage: result.Sbom.Components.PageInfo.HasNextPage,
		EndCursor:   result.Sbom.Components.PageInfo.EndCursor,
	}, nil
}

// GetComponent fetches a single component by ID
func (c *Client) GetComponent(ctx context.Context, id, sbomID string) (*SbomComponent, error) {
	vars := map[string]interface{}{
		"id":     id,
		"sbomId": sbomID,
	}

	var result struct {
		Component struct {
			ID          string    `json:"id"`
			Name        string    `json:"name"`
			Version     string    `json:"version"`
			Kind        string    `json:"kind"`
			Purl        string    `json:"purl"`
			Cpes        []string  `json:"cpes"`
			LicensesExp string    `json:"licensesExp"`
			Group       string    `json:"group"`
			Description string    `json:"description"`
			Primary     bool      `json:"primary"`
			Internal    bool      `json:"internal"`
			SbomID      string    `json:"sbomId"`
			UpdatedAt   time.Time `json:"updatedAt"`
			Sbom        struct {
				ID             string `json:"id"`
				ProjectVersion string `json:"projectVersion"`
				Project        struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"project"`
			} `json:"sbom"`
		} `json:"component"`
	}

	if err := c.gql.Execute(ctx, graphql.ComponentQuery, vars, &result); err != nil {
		return nil, err
	}

	return &SbomComponent{
		ID:          result.Component.ID,
		Name:        result.Component.Name,
		Version:     result.Component.Version,
		Kind:        result.Component.Kind,
		Purl:        result.Component.Purl,
		Cpes:        result.Component.Cpes,
		LicensesExp: result.Component.LicensesExp,
		Group:       result.Component.Group,
		Description: result.Component.Description,
		Primary:     result.Component.Primary,
		Internal:    result.Component.Internal,
		SbomID:      result.Component.SbomID,
		UpdatedAt:   result.Component.UpdatedAt,
		Sbom: &Sbom{
			ID:             result.Component.Sbom.ID,
			ProjectVersion: result.Component.Sbom.ProjectVersion,
			Project: &Project{
				ID:   result.Component.Sbom.Project.ID,
				Name: result.Component.Sbom.Project.Name,
			},
		},
	}, nil
}

// CompareSboms compares two SBOMs and returns the differences
func (c *Client) CompareSboms(ctx context.Context, sourceSbomID, targetSbomID string) ([]SbomDiff, error) {
	vars := map[string]interface{}{
		"sourceSbomId": sourceSbomID,
		"targetSbomId": targetSbomID,
	}
	// Note: sourceSbomId is used for the sbom(sbomId:) query parameter

	var result struct {
		Sbom struct {
			SbomDrift []struct {
				DiffType           string   `json:"diffType"`
				DiffTags           []string `json:"diffTags"`
				SubjectComponentID string   `json:"subjectComponentId"`
				TargetComponentID  string   `json:"targetComponentId"`
				SubjectComponent   *struct {
					ID      string `json:"id"`
					Name    string `json:"name"`
					Version string `json:"version"`
					Purl    string `json:"purl"`
				} `json:"subjectComponent"`
				TargetComponent *struct {
					ID      string `json:"id"`
					Name    string `json:"name"`
					Version string `json:"version"`
					Purl    string `json:"purl"`
				} `json:"targetComponent"`
			} `json:"sbomDrift"`
		} `json:"sbom"`
	}

	if err := c.gql.Execute(ctx, graphql.SbomDriftQuery, vars, &result); err != nil {
		return nil, err
	}

	diffs := make([]SbomDiff, len(result.Sbom.SbomDrift))
	for i, d := range result.Sbom.SbomDrift {
		diff := SbomDiff{
			DiffType:           d.DiffType,
			DiffTags:           d.DiffTags,
			SubjectComponentID: d.SubjectComponentID,
			TargetComponentID:  d.TargetComponentID,
		}
		if d.SubjectComponent != nil {
			diff.SubjectComponent = &SbomComponent{
				ID:      d.SubjectComponent.ID,
				Name:    d.SubjectComponent.Name,
				Version: d.SubjectComponent.Version,
				Purl:    d.SubjectComponent.Purl,
			}
		}
		if d.TargetComponent != nil {
			diff.TargetComponent = &SbomComponent{
				ID:      d.TargetComponent.ID,
				Name:    d.TargetComponent.Name,
				Version: d.TargetComponent.Version,
				Purl:    d.TargetComponent.Purl,
			}
		}
		diffs[i] = diff
	}

	return diffs, nil
}
