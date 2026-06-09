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

// Version represents a version (formerly SBOM)
type Version struct {
	ID            string
	Version       string
	Spec          string
	SpecVersion   string
	Format        string
	Lifecycle     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	EnvironmentID string
	Stats         *VersionStats
	Environment   *Environment
}

// VersionStats contains statistics for a version
type VersionStats struct {
	CompCount         int
	CompPurlCount     int
	CompCpeCount      int
	CompLicenseCount  int
	CompSupplierCount int
	VulnStats         map[string]interface{}
}

// ComponentVulnerabilitySummary contains vulnerability counts for a component.
type ComponentVulnerabilitySummary struct {
	TotalCount int
	Stats      map[string]interface{}
}

// VersionComponent represents a component in a version
type VersionComponent struct {
	ID           string
	Name         string
	Version      string
	Kind         string
	Purl         string
	Cpes         []string
	LicensesExp  string
	Group        string
	Description  string
	Scope        string
	Copyright    string
	Primary      bool
	Internal     bool
	UniqueID     string
	Notice       string
	SupportLevel string
	EndOfSupport string
	Checksums    []ComponentChecksum
	ExternalURLs []ComponentExternalURL
	VersionID    string
	UpdatedAt    time.Time
	VersionInfo  *Version
	VulnSummary  *ComponentVulnerabilitySummary
}

// VersionsResult represents the result of listing versions
type VersionsResult struct {
	Versions    []Version
	TotalCount  int
	HasNextPage bool
	EndCursor   string
}

// ComponentsResult represents the result of listing components
type ComponentsResult struct {
	Components  []VersionComponent
	TotalCount  int
	HasNextPage bool
	EndCursor   string
}

// VersionSearchInput contains parameters for searching versions across the organization.
type VersionSearchInput struct {
	Search string
	First  int
	After  string
}

// DownloadSBOMInput contains parameters for SBOM downloads.
type DownloadSBOMInput struct {
	VersionID        string
	Spec             string
	SpecVersion      string
	IncludeVulns     *bool
	RequireCompleted []string
}

// DownloadSBOMResult contains SBOM download readiness, metadata, and content.
type DownloadSBOMResult struct {
	Ready            bool
	Filename         string
	ContentType      string
	Content          string
	ProcessingStatus map[string]string
}

// VersionDiff represents a diff entry between two versions
type VersionDiff struct {
	DiffType           string
	DiffTags           []string
	SubjectComponent   *VersionComponent
	SubjectComponentID string
	TargetComponent    *VersionComponent
	TargetComponentID  string
}

type versionNode struct {
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
	Project *struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		ProjectGroupID string `json:"projectGroupId"`
		ProjectGroup   *struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"projectGroup"`
	} `json:"project"`
}

func mapVersionNode(n versionNode) Version {
	var stats *VersionStats
	if n.Stats != nil {
		stats = &VersionStats{
			CompCount:         n.Stats.CompCount,
			CompPurlCount:     n.Stats.CompPurlCount,
			CompCpeCount:      n.Stats.CompCpeCount,
			CompLicenseCount:  n.Stats.CompLicenseCount,
			CompSupplierCount: n.Stats.CompSupplierCount,
			VulnStats:         n.Stats.VulnStats,
		}
	}
	version := Version{
		ID:            n.ID,
		Version:       n.ProjectVersion,
		Spec:          n.Spec,
		SpecVersion:   n.SpecVersion,
		Format:        n.Format,
		Lifecycle:     n.Lifecycle,
		CreatedAt:     n.CreatedAt,
		UpdatedAt:     n.UpdatedAt,
		EnvironmentID: n.ProjectID,
		Stats:         stats,
	}
	if n.Project != nil {
		version.Environment = &Environment{
			ID:        n.Project.ID,
			Name:      n.Project.Name,
			ProductID: n.Project.ProjectGroupID,
		}
		if n.Project.ProjectGroup != nil {
			version.Environment.Product = &Product{
				ID:   n.Project.ProjectGroup.ID,
				Name: n.Project.ProjectGroup.Name,
			}
		}
	}
	return version
}

// ListVersionsInput contains parameters for listing versions
type ListVersionsInput struct {
	EnvironmentID string
	First         int
	After         string
	Lifecycle     []string
}

// ListVersions fetches versions for an environment
func (c *Client) ListVersions(ctx context.Context, input ListVersionsInput) (*VersionsResult, error) {
	vars := map[string]interface{}{
		"projectId": input.EnvironmentID,
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

	versions := make([]Version, len(result.Project.SbomVersions.Nodes))
	for i, n := range result.Project.SbomVersions.Nodes {
		var stats *VersionStats
		if n.Stats != nil {
			stats = &VersionStats{
				CompCount:         n.Stats.CompCount,
				CompPurlCount:     n.Stats.CompPurlCount,
				CompCpeCount:      n.Stats.CompCpeCount,
				CompLicenseCount:  n.Stats.CompLicenseCount,
				CompSupplierCount: n.Stats.CompSupplierCount,
				VulnStats:         n.Stats.VulnStats,
			}
		}
		versions[i] = Version{
			ID:            n.ID,
			Version:       n.ProjectVersion,
			Spec:          n.Spec,
			SpecVersion:   n.SpecVersion,
			Format:        n.Format,
			Lifecycle:     n.Lifecycle,
			CreatedAt:     n.CreatedAt,
			UpdatedAt:     n.UpdatedAt,
			EnvironmentID: n.ProjectID,
			Stats:         stats,
		}
	}

	return &VersionsResult{
		Versions:    versions,
		TotalCount:  result.Project.SbomVersions.TotalCount,
		HasNextPage: result.Project.SbomVersions.PageInfo.HasNextPage,
		EndCursor:   result.Project.SbomVersions.PageInfo.EndCursor,
	}, nil
}

// SearchVersions searches versions across the organization by version string.
func (c *Client) SearchVersions(ctx context.Context, input VersionSearchInput) (*VersionsResult, error) {
	vars := map[string]interface{}{
		"search": input.Search,
	}
	if input.First > 0 {
		vars["first"] = input.First
	} else {
		vars["first"] = 20
	}
	if input.After != "" {
		vars["after"] = input.After
	}

	var result struct {
		ProjectVersions struct {
			Nodes      []versionNode `json:"nodes"`
			TotalCount int           `json:"totalCount"`
			PageInfo   struct {
				HasNextPage bool   `json:"hasNextPage"`
				EndCursor   string `json:"endCursor"`
			} `json:"pageInfo"`
		} `json:"projectVersions"`
	}

	if err := c.gql.Execute(ctx, graphql.ProjectVersionsSearchQuery, vars, &result); err != nil {
		return nil, err
	}

	versions := make([]Version, len(result.ProjectVersions.Nodes))
	for i, n := range result.ProjectVersions.Nodes {
		versions[i] = mapVersionNode(n)
	}

	return &VersionsResult{
		Versions:    versions,
		TotalCount:  result.ProjectVersions.TotalCount,
		HasNextPage: result.ProjectVersions.PageInfo.HasNextPage,
		EndCursor:   result.ProjectVersions.PageInfo.EndCursor,
	}, nil
}

// GetVersion fetches a single version by ID
func (c *Client) GetVersion(ctx context.Context, id string) (*Version, error) {
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

	var stats *VersionStats
	if result.Sbom.Stats != nil {
		stats = &VersionStats{
			CompCount:         result.Sbom.Stats.CompCount,
			CompPurlCount:     result.Sbom.Stats.CompPurlCount,
			CompCpeCount:      result.Sbom.Stats.CompCpeCount,
			CompLicenseCount:  result.Sbom.Stats.CompLicenseCount,
			CompSupplierCount: result.Sbom.Stats.CompSupplierCount,
			VulnStats:         result.Sbom.Stats.VulnStats,
		}
	}

	return &Version{
		ID:            result.Sbom.ID,
		Version:       result.Sbom.ProjectVersion,
		Spec:          result.Sbom.Spec,
		SpecVersion:   result.Sbom.SpecVersion,
		Format:        result.Sbom.Format,
		Lifecycle:     result.Sbom.Lifecycle,
		CreatedAt:     result.Sbom.CreatedAt,
		UpdatedAt:     result.Sbom.UpdatedAt,
		EnvironmentID: result.Sbom.ProjectID,
		Stats:         stats,
		Environment: &Environment{
			ID:        result.Sbom.Project.ID,
			Name:      result.Sbom.Project.Name,
			ProductID: result.Sbom.Project.ProjectGroupID,
		},
	}, nil
}

// DownloadSBOM fetches SBOM download readiness, metadata, and content.
func (c *Client) DownloadSBOM(ctx context.Context, input DownloadSBOMInput) (*DownloadSBOMResult, error) {
	vars := map[string]interface{}{
		"sbomId": input.VersionID,
	}
	if input.Spec != "" {
		vars["spec"] = input.Spec
	}
	if input.SpecVersion != "" {
		vars["specVersion"] = input.SpecVersion
	}
	if input.IncludeVulns != nil {
		vars["includeVulns"] = *input.IncludeVulns
	}
	if len(input.RequireCompleted) > 0 {
		vars["requireCompleted"] = input.RequireCompleted
	}

	var result struct {
		Sbom struct {
			Download *struct {
				Ready            bool   `json:"ready"`
				Filename         string `json:"filename"`
				ContentType      string `json:"contentType"`
				Content          string `json:"content"`
				ProcessingStatus *struct {
					Automation string `json:"automation"`
					PolicyScan string `json:"policyScan"`
					VulnScan   string `json:"vulnScan"`
				} `json:"processingStatus"`
			} `json:"download"`
		} `json:"sbom"`
	}

	if err := c.gql.Execute(ctx, graphql.SbomDownloadQuery, vars, &result); err != nil {
		return nil, err
	}
	if result.Sbom.Download == nil {
		return &DownloadSBOMResult{}, nil
	}

	download := &DownloadSBOMResult{
		Ready:       result.Sbom.Download.Ready,
		Filename:    result.Sbom.Download.Filename,
		ContentType: result.Sbom.Download.ContentType,
		Content:     result.Sbom.Download.Content,
	}
	if result.Sbom.Download.ProcessingStatus != nil {
		download.ProcessingStatus = map[string]string{
			"automation": result.Sbom.Download.ProcessingStatus.Automation,
			"policyScan": result.Sbom.Download.ProcessingStatus.PolicyScan,
			"vulnScan":   result.Sbom.Download.ProcessingStatus.VulnScan,
		}
	}

	return download, nil
}

// ListComponentsInput contains parameters for listing components
type ListComponentsInput struct {
	VersionID string
	First     int
	After     string
	Search    string
	Kind      []string
	Direct    *bool
}

// ListComponents fetches components for a version
func (c *Client) ListComponents(ctx context.Context, input ListComponentsInput) (*ComponentsResult, error) {
	vars := map[string]interface{}{
		"sbomId": input.VersionID,
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
					ID          string   `json:"id"`
					Name        string   `json:"name"`
					Version     string   `json:"version"`
					Kind        string   `json:"kind"`
					Purl        string   `json:"purl"`
					Cpes        []string `json:"cpes"`
					LicensesExp string   `json:"licensesExp"`
					Group       string   `json:"group"`
					Description string   `json:"description"`
					Primary     bool     `json:"primary"`
					Internal    bool     `json:"internal"`
					SbomID      string   `json:"sbomId"`
					Stats       *struct {
						VulnStats      map[string]interface{} `json:"vulnStats"`
						VulnTotalCount int                    `json:"vulnTotalCount"`
					} `json:"stats"`
					UpdatedAt time.Time `json:"updatedAt"`
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

	components := make([]VersionComponent, len(result.Sbom.Components.Nodes))
	for i, n := range result.Sbom.Components.Nodes {
		components[i] = VersionComponent{
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
			VersionID:   n.SbomID,
			UpdatedAt:   n.UpdatedAt,
		}
		if n.Stats != nil {
			components[i].VulnSummary = &ComponentVulnerabilitySummary{
				TotalCount: n.Stats.VulnTotalCount,
				Stats:      n.Stats.VulnStats,
			}
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
func (c *Client) GetComponent(ctx context.Context, id, versionID string) (*VersionComponent, error) {
	vars := map[string]interface{}{
		"id":     id,
		"sbomId": versionID,
	}

	var result struct {
		Component struct {
			ID          string   `json:"id"`
			Name        string   `json:"name"`
			Version     string   `json:"version"`
			Kind        string   `json:"kind"`
			Purl        string   `json:"purl"`
			Cpes        []string `json:"cpes"`
			LicensesExp string   `json:"licensesExp"`
			Group       string   `json:"group"`
			Description string   `json:"description"`
			Primary     bool     `json:"primary"`
			Internal    bool     `json:"internal"`
			SbomID      string   `json:"sbomId"`
			Stats       *struct {
				VulnStats      map[string]interface{} `json:"vulnStats"`
				VulnTotalCount int                    `json:"vulnTotalCount"`
			} `json:"stats"`
			UpdatedAt time.Time `json:"updatedAt"`
			Sbom      struct {
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

	component := &VersionComponent{
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
		VersionID:   result.Component.SbomID,
		UpdatedAt:   result.Component.UpdatedAt,
		VersionInfo: &Version{
			ID:      result.Component.Sbom.ID,
			Version: result.Component.Sbom.ProjectVersion,
			Environment: &Environment{
				ID:   result.Component.Sbom.Project.ID,
				Name: result.Component.Sbom.Project.Name,
			},
		},
	}
	if result.Component.Stats != nil {
		component.VulnSummary = &ComponentVulnerabilitySummary{
			TotalCount: result.Component.Stats.VulnTotalCount,
			Stats:      result.Component.Stats.VulnStats,
		}
	}
	return component, nil
}

// CompareVersions compares two versions and returns the differences
func (c *Client) CompareVersions(ctx context.Context, sourceVersionID, targetVersionID string) ([]VersionDiff, error) {
	vars := map[string]interface{}{
		"sourceSbomId": sourceVersionID,
		"targetSbomId": targetVersionID,
	}

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

	diffs := make([]VersionDiff, len(result.Sbom.SbomDrift))
	for i, d := range result.Sbom.SbomDrift {
		diff := VersionDiff{
			DiffType:           d.DiffType,
			DiffTags:           d.DiffTags,
			SubjectComponentID: d.SubjectComponentID,
			TargetComponentID:  d.TargetComponentID,
		}
		if d.SubjectComponent != nil {
			diff.SubjectComponent = &VersionComponent{
				ID:      d.SubjectComponent.ID,
				Name:    d.SubjectComponent.Name,
				Version: d.SubjectComponent.Version,
				Purl:    d.SubjectComponent.Purl,
			}
		}
		if d.TargetComponent != nil {
			diff.TargetComponent = &VersionComponent{
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
