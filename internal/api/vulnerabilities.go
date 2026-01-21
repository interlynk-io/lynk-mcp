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
	"time"

	"github.com/interlynk-io/lynk-mcp/internal/graphql"
)

// Vuln represents a vulnerability
type Vuln struct {
	ID             string
	VulnID         string
	Description    string
	Severity       string
	CvssScore      float64
	CvssVector     string
	Source         string
	PublishedAt    time.Time
	LastModifiedAt time.Time
	UpdatedAt      time.Time
	VulnInfo       *VulnInfo
}

// VulnInfo contains additional vulnerability information
type VulnInfo struct {
	ID             string
	CveID          string
	EpssScore      float64
	EpssPercentile float64
	Kev            bool
	Cwes           []string
	Advisories     []string
}

// ComponentVuln represents a vulnerability associated with a component
type ComponentVuln struct {
	ID               string
	ComponentID      string
	VulnID           string
	VersionID        string
	FixedIn          string
	FixedVersions    []string
	Detail           string
	Impact           string
	ActionStmt       string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Component        *VersionComponent
	Vuln             *Vuln
	VexStatus        *VexStatus
	VexJustification *VexJustification
}

// VexStatus represents a VEX status
type VexStatus struct {
	ID   string
	Name string
}

// VexJustification represents a VEX justification
type VexJustification struct {
	ID   string
	Name string
}

// ComponentVulnsResult represents the result of listing component vulnerabilities
type ComponentVulnsResult struct {
	ComponentVulns []ComponentVuln
	TotalCount     int
	HasNextPage    bool
	EndCursor      string
}

// ListVersionVulnsInput contains parameters for listing version vulnerabilities
type ListVersionVulnsInput struct {
	VersionID string
	First     int
	After     string
	Severity  []string
	Status    []string
	Kev       *bool
	Search    string
}

// ListVersionVulns fetches vulnerabilities for a version
func (c *Client) ListVersionVulns(ctx context.Context, input ListVersionVulnsInput) (*ComponentVulnsResult, error) {
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
	if len(input.Severity) > 0 {
		vars["severity"] = input.Severity
	}
	if len(input.Status) > 0 {
		vars["status"] = input.Status
	}
	if input.Kev != nil {
		vars["kev"] = *input.Kev
	}
	if input.Search != "" {
		vars["search"] = input.Search
	}

	var result struct {
		Sbom struct {
			Vulns struct {
				Nodes []struct {
					ID            string    `json:"id"`
					ComponentID   string    `json:"componentId"`
					VulnID        string    `json:"vulnId"`
					SbomID        string    `json:"sbomId"`
					FixedIn       string    `json:"fixedIn"`
					FixedVersions []string  `json:"fixedVersions"`
					Detail        string    `json:"detail"`
					Impact        string    `json:"impact"`
					ActionStmt    string    `json:"actionStmt"`
					CreatedAt     time.Time `json:"createdAt"`
					UpdatedAt     time.Time `json:"updatedAt"`
					Component     *struct {
						ID      string `json:"id"`
						Name    string `json:"name"`
						Version string `json:"version"`
						Purl    string `json:"purl"`
					} `json:"component"`
					Vuln *struct {
						ID             string    `json:"id"`
						VulnID         string    `json:"vulnId"`
						Desc           string    `json:"desc"`
						Sev            string    `json:"sev"`
						CvssScore      float64   `json:"cvssScore"`
						CvssVector     string    `json:"cvssVector"`
						Source         string    `json:"source"`
						PublishedAt    time.Time `json:"publishedAt"`
						LastModifiedAt time.Time `json:"lastModifiedAt"`
						VulnInfo       *struct {
							CveID          string   `json:"cveId"`
							EpssScore      float64  `json:"epssScore"`
							EpssPercentile float64  `json:"epssPercentile"`
							Kev            bool     `json:"kev"`
							Cwes           []string `json:"cwes"`
						} `json:"vulnInfo"`
					} `json:"vuln"`
					VexStatus *struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"vexStatus"`
					VexJustification *struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"vexJustification"`
				} `json:"nodes"`
				TotalCount int `json:"totalCount"`
				PageInfo   struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"vulns"`
		} `json:"sbom"`
	}

	if err := c.gql.Execute(ctx, graphql.SbomVulnsQuery, vars, &result); err != nil {
		return nil, err
	}

	vulns := make([]ComponentVuln, len(result.Sbom.Vulns.Nodes))
	for i, n := range result.Sbom.Vulns.Nodes {
		cv := ComponentVuln{
			ID:            n.ID,
			ComponentID:   n.ComponentID,
			VulnID:        n.VulnID,
			VersionID:     n.SbomID,
			FixedIn:       n.FixedIn,
			FixedVersions: n.FixedVersions,
			Detail:        n.Detail,
			Impact:        n.Impact,
			ActionStmt:    n.ActionStmt,
			CreatedAt:     n.CreatedAt,
			UpdatedAt:     n.UpdatedAt,
		}
		if n.Component != nil {
			cv.Component = &VersionComponent{
				ID:      n.Component.ID,
				Name:    n.Component.Name,
				Version: n.Component.Version,
				Purl:    n.Component.Purl,
			}
		}
		if n.Vuln != nil {
			cv.Vuln = &Vuln{
				ID:             n.Vuln.ID,
				VulnID:         n.Vuln.VulnID,
				Description:    n.Vuln.Desc,
				Severity:       n.Vuln.Sev,
				CvssScore:      n.Vuln.CvssScore,
				CvssVector:     n.Vuln.CvssVector,
				Source:         n.Vuln.Source,
				PublishedAt:    n.Vuln.PublishedAt,
				LastModifiedAt: n.Vuln.LastModifiedAt,
			}
			if n.Vuln.VulnInfo != nil {
				cv.Vuln.VulnInfo = &VulnInfo{
					CveID:          n.Vuln.VulnInfo.CveID,
					EpssScore:      n.Vuln.VulnInfo.EpssScore,
					EpssPercentile: n.Vuln.VulnInfo.EpssPercentile,
					Kev:            n.Vuln.VulnInfo.Kev,
					Cwes:           n.Vuln.VulnInfo.Cwes,
				}
			}
		}
		if n.VexStatus != nil {
			cv.VexStatus = &VexStatus{
				ID:   n.VexStatus.ID,
				Name: n.VexStatus.Name,
			}
		}
		if n.VexJustification != nil {
			cv.VexJustification = &VexJustification{
				ID:   n.VexJustification.ID,
				Name: n.VexJustification.Name,
			}
		}
		vulns[i] = cv
	}

	return &ComponentVulnsResult{
		ComponentVulns: vulns,
		TotalCount:     result.Sbom.Vulns.TotalCount,
		HasNextPage:    result.Sbom.Vulns.PageInfo.HasNextPage,
		EndCursor:      result.Sbom.Vulns.PageInfo.EndCursor,
	}, nil
}

// GetVuln fetches a single vulnerability by ID or CVE ID
func (c *Client) GetVuln(ctx context.Context, id, vulnID string) (*Vuln, error) {
	// If CVE ID is provided (starts with CVE- or similar), use CveLookup
	if vulnID != "" {
		return c.lookupByCveID(ctx, vulnID)
	}

	// Otherwise use the internal UUID lookup
	if id != "" {
		return c.lookupByUUID(ctx, id)
	}

	return nil, fmt.Errorf("either id or vulnID must be provided")
}

// lookupByUUID fetches vulnerability by internal UUID
func (c *Client) lookupByUUID(ctx context.Context, id string) (*Vuln, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var result struct {
		Vuln struct {
			ID             string    `json:"id"`
			VulnID         string    `json:"vulnId"`
			Desc           string    `json:"desc"`
			Sev            string    `json:"sev"`
			CvssScore      float64   `json:"cvssScore"`
			CvssVector     string    `json:"cvssVector"`
			Source         string    `json:"source"`
			PublishedAt    time.Time `json:"publishedAt"`
			LastModifiedAt time.Time `json:"lastModifiedAt"`
			UpdatedAt      time.Time `json:"updatedAt"`
			VulnInfo       *struct {
				ID             string   `json:"id"`
				CveID          string   `json:"cveId"`
				EpssScore      float64  `json:"epssScore"`
				EpssPercentile float64  `json:"epssPercentile"`
				Kev            bool     `json:"kev"`
				Cwes           []string `json:"cwes"`
				Advisories     []string `json:"advisories"`
			} `json:"vulnInfo"`
		} `json:"vuln"`
	}

	if err := c.gql.Execute(ctx, graphql.VulnQuery, vars, &result); err != nil {
		return nil, err
	}

	v := &Vuln{
		ID:             result.Vuln.ID,
		VulnID:         result.Vuln.VulnID,
		Description:    result.Vuln.Desc,
		Severity:       result.Vuln.Sev,
		CvssScore:      result.Vuln.CvssScore,
		CvssVector:     result.Vuln.CvssVector,
		Source:         result.Vuln.Source,
		PublishedAt:    result.Vuln.PublishedAt,
		LastModifiedAt: result.Vuln.LastModifiedAt,
		UpdatedAt:      result.Vuln.UpdatedAt,
	}

	if result.Vuln.VulnInfo != nil {
		v.VulnInfo = &VulnInfo{
			ID:             result.Vuln.VulnInfo.ID,
			CveID:          result.Vuln.VulnInfo.CveID,
			EpssScore:      result.Vuln.VulnInfo.EpssScore,
			EpssPercentile: result.Vuln.VulnInfo.EpssPercentile,
			Kev:            result.Vuln.VulnInfo.Kev,
			Cwes:           result.Vuln.VulnInfo.Cwes,
			Advisories:     result.Vuln.VulnInfo.Advisories,
		}
	}

	return v, nil
}

// lookupByCveID fetches vulnerability by CVE ID using cveLookup query
func (c *Client) lookupByCveID(ctx context.Context, vulnID string) (*Vuln, error) {
	vars := map[string]interface{}{
		"vulnId": vulnID,
	}

	var result struct {
		CveLookup *struct {
			VulnID       string    `json:"vulnId"`
			Description  string    `json:"description"`
			Severity     string    `json:"severity"`
			Published    time.Time `json:"published"`
			LastModified time.Time `json:"lastModified"`
			CvssScore    float64   `json:"cvssScore"`
			CvssVector   string    `json:"cvssVector"`
			Cwes         []string  `json:"cwes"`
			Advisories   []string  `json:"advisories"`
		} `json:"cveLookup"`
	}

	if err := c.gql.Execute(ctx, graphql.CveLookupQuery, vars, &result); err != nil {
		return nil, err
	}

	if result.CveLookup == nil {
		return nil, fmt.Errorf("vulnerability not found: %s", vulnID)
	}

	return &Vuln{
		VulnID:         result.CveLookup.VulnID,
		Description:    result.CveLookup.Description,
		Severity:       result.CveLookup.Severity,
		CvssScore:      result.CveLookup.CvssScore,
		CvssVector:     result.CveLookup.CvssVector,
		PublishedAt:    result.CveLookup.Published,
		LastModifiedAt: result.CveLookup.LastModified,
		VulnInfo: &VulnInfo{
			Cwes:       result.CveLookup.Cwes,
			Advisories: result.CveLookup.Advisories,
		},
	}, nil
}

// ListComponentVulnsInput contains parameters for listing component vulnerabilities
type ListComponentVulnsInput struct {
	First          int
	After          string
	Severity       []string
	Status         []string
	Kev            *bool
	Search         string
	EnvironmentIDs []string
	ProductIDs     []string
}

// ListComponentVulns fetches component vulnerabilities across the organization
func (c *Client) ListComponentVulns(ctx context.Context, input ListComponentVulnsInput) (*ComponentVulnsResult, error) {
	vars := make(map[string]interface{})
	if input.First > 0 {
		vars["first"] = input.First
	} else {
		vars["first"] = 50
	}
	if input.After != "" {
		vars["after"] = input.After
	}
	if len(input.Severity) > 0 {
		vars["severity"] = input.Severity
	}
	if len(input.Status) > 0 {
		vars["status"] = input.Status
	}
	if input.Kev != nil {
		vars["kev"] = *input.Kev
	}
	if input.Search != "" {
		vars["search"] = input.Search
	}
	if len(input.EnvironmentIDs) > 0 {
		vars["projectIds"] = input.EnvironmentIDs
	}
	if len(input.ProductIDs) > 0 {
		vars["projectGroupIds"] = input.ProductIDs
	}

	var result struct {
		ComponentVulns struct {
			Nodes []struct {
				ID            string    `json:"id"`
				ComponentID   string    `json:"componentId"`
				VulnID        string    `json:"vulnId"`
				SbomID        string    `json:"sbomId"`
				FixedIn       string    `json:"fixedIn"`
				FixedVersions []string  `json:"fixedVersions"`
				CreatedAt     time.Time `json:"createdAt"`
				UpdatedAt     time.Time `json:"updatedAt"`
				Component     *struct {
					ID      string `json:"id"`
					Name    string `json:"name"`
					Version string `json:"version"`
					Purl    string `json:"purl"`
					SbomID  string `json:"sbomId"`
				} `json:"component"`
				Vuln *struct {
					ID        string  `json:"id"`
					VulnID    string  `json:"vulnId"`
					Desc      string  `json:"desc"`
					Sev       string  `json:"sev"`
					CvssScore float64 `json:"cvssScore"`
					Source    string  `json:"source"`
					VulnInfo  *struct {
						EpssScore      float64 `json:"epssScore"`
						EpssPercentile float64 `json:"epssPercentile"`
						Kev            bool    `json:"kev"`
					} `json:"vulnInfo"`
				} `json:"vuln"`
				VexStatus *struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"vexStatus"`
			} `json:"nodes"`
			TotalCount int `json:"totalCount"`
			PageInfo   struct {
				HasNextPage bool   `json:"hasNextPage"`
				EndCursor   string `json:"endCursor"`
			} `json:"pageInfo"`
		} `json:"componentVulns"`
	}

	if err := c.gql.Execute(ctx, graphql.ComponentVulnsQuery, vars, &result); err != nil {
		return nil, err
	}

	vulns := make([]ComponentVuln, len(result.ComponentVulns.Nodes))
	for i, n := range result.ComponentVulns.Nodes {
		cv := ComponentVuln{
			ID:            n.ID,
			ComponentID:   n.ComponentID,
			VulnID:        n.VulnID,
			VersionID:     n.SbomID,
			FixedIn:       n.FixedIn,
			FixedVersions: n.FixedVersions,
			CreatedAt:     n.CreatedAt,
			UpdatedAt:     n.UpdatedAt,
		}
		if n.Component != nil {
			cv.Component = &VersionComponent{
				ID:        n.Component.ID,
				Name:      n.Component.Name,
				Version:   n.Component.Version,
				Purl:      n.Component.Purl,
				VersionID: n.Component.SbomID,
			}
		}
		if n.Vuln != nil {
			cv.Vuln = &Vuln{
				ID:          n.Vuln.ID,
				VulnID:      n.Vuln.VulnID,
				Description: n.Vuln.Desc,
				Severity:    n.Vuln.Sev,
				CvssScore:   n.Vuln.CvssScore,
				Source:      n.Vuln.Source,
			}
			if n.Vuln.VulnInfo != nil {
				cv.Vuln.VulnInfo = &VulnInfo{
					EpssScore:      n.Vuln.VulnInfo.EpssScore,
					EpssPercentile: n.Vuln.VulnInfo.EpssPercentile,
					Kev:            n.Vuln.VulnInfo.Kev,
				}
			}
		}
		if n.VexStatus != nil {
			cv.VexStatus = &VexStatus{
				ID:   n.VexStatus.ID,
				Name: n.VexStatus.Name,
			}
		}
		vulns[i] = cv
	}

	return &ComponentVulnsResult{
		ComponentVulns: vulns,
		TotalCount:     result.ComponentVulns.TotalCount,
		HasNextPage:    result.ComponentVulns.PageInfo.HasNextPage,
		EndCursor:      result.ComponentVulns.PageInfo.EndCursor,
	}, nil
}

// GetVexStatuses fetches all VEX statuses
func (c *Client) GetVexStatuses(ctx context.Context) ([]VexStatus, error) {
	var result struct {
		VexStatuses []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"vexStatuses"`
	}

	if err := c.gql.Execute(ctx, graphql.VexStatusesQuery, nil, &result); err != nil {
		return nil, err
	}

	statuses := make([]VexStatus, len(result.VexStatuses))
	for i, s := range result.VexStatuses {
		statuses[i] = VexStatus{
			ID:   s.ID,
			Name: s.Name,
		}
	}

	return statuses, nil
}

// GetVexJustifications fetches all VEX justifications
func (c *Client) GetVexJustifications(ctx context.Context) ([]VexJustification, error) {
	var result struct {
		VexJustifications []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"vexJustifications"`
	}

	if err := c.gql.Execute(ctx, graphql.VexJustificationsQuery, nil, &result); err != nil {
		return nil, err
	}

	justifications := make([]VexJustification, len(result.VexJustifications))
	for i, j := range result.VexJustifications {
		justifications[i] = VexJustification{
			ID:   j.ID,
			Name: j.Name,
		}
	}

	return justifications, nil
}
