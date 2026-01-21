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

package graphql

import (
	"encoding/json"
	"time"
)

// FlexibleTime is a time.Time that can unmarshal from multiple date formats
// including ISO 8601 dates with and without timezone information
type FlexibleTime struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler for FlexibleTime
func (ft *FlexibleTime) UnmarshalJSON(data []byte) error {
	// Handle null
	if string(data) == "null" {
		return nil
	}

	// Remove quotes
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	if str == "" {
		return nil
	}

	// Try multiple formats
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05.999",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	var parseErr error
	for _, format := range formats {
		t, err := time.Parse(format, str)
		if err == nil {
			ft.Time = t
			return nil
		}
		parseErr = err
	}

	return parseErr
}

// MarshalJSON implements json.Marshaler for FlexibleTime
func (ft FlexibleTime) MarshalJSON() ([]byte, error) {
	if ft.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(ft.Time.Format(time.RFC3339Nano))
}

// Organization represents a Lynk organization
type Organization struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Email     string             `json:"email,omitempty"`
	URL       string             `json:"url,omitempty"`
	Status    string             `json:"status"`
	Tier      string             `json:"tier"`
	UpdatedAt time.Time          `json:"updatedAt"`
	Metrics   *OrganizationLiveMetric `json:"metrics,omitempty"`
}

// OrganizationLiveMetric contains live metrics for an organization
type OrganizationLiveMetric struct {
	ProjectGroupsCount int `json:"projectGroupsCount"`
	ProjectsCount      int `json:"projectsCount"`
	SbomsCount         int `json:"sbomsCount"`
	ComponentsCount    int `json:"componentsCount"`
	VulnsCount         int `json:"vulnsCount"`
	PoliciesCount      int `json:"policiesCount"`
}

// ProjectGroup represents a product/project group
type ProjectGroup struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Description    string     `json:"description,omitempty"`
	Enabled        bool       `json:"enabled"`
	OrganizationID string     `json:"organizationId"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	SbomsCount     int        `json:"sbomsCount,omitempty"`
	Projects       []Project  `json:"projects,omitempty"`
}

// ProjectGroupConnection represents a paginated list of project groups
type ProjectGroupConnection struct {
	Nodes      []ProjectGroup `json:"nodes"`
	TotalCount int            `json:"totalCount"`
	PageInfo   PageInfo       `json:"pageInfo"`
}

// Project represents a stream/project within a project group
type Project struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description,omitempty"`
	Enabled        bool     `json:"enabled"`
	ProjectGroupID string   `json:"projectGroupId,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt"`
	SbomsCount     int      `json:"sbomsCount,omitempty"`
}

// Sbom represents an SBOM document
type Sbom struct {
	ID             string     `json:"id"`
	ProjectVersion string     `json:"projectVersion,omitempty"`
	Spec           string     `json:"spec,omitempty"`
	SpecVersion    string     `json:"specVersion,omitempty"`
	Format         string     `json:"format,omitempty"`
	Lifecycle      string     `json:"lifecycle"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	ProjectID      string     `json:"projectId"`
	Stats          *SbomStats `json:"stats,omitempty"`
	Project        *Project   `json:"project,omitempty"`
}

// SbomStats contains statistics for an SBOM
type SbomStats struct {
	CompCount         int                    `json:"compCount"`
	CompPurlCount     int                    `json:"compPurlCount"`
	CompCpeCount      int                    `json:"compCpeCount"`
	CompLicenseCount  int                    `json:"compLicenseCount"`
	CompSupplierCount int                    `json:"compSupplierCount"`
	VulnStats         map[string]interface{} `json:"vulnStats"`
}

// SbomConnection represents a paginated list of SBOMs
type SbomConnection struct {
	Nodes      []Sbom   `json:"nodes"`
	TotalCount int      `json:"totalCount"`
	PageInfo   PageInfo `json:"pageInfo"`
}

// SbomComponent represents a component in an SBOM
type SbomComponent struct {
	ID          string   `json:"id"`
	Name        string   `json:"name,omitempty"`
	Version     string   `json:"version,omitempty"`
	Kind        string   `json:"kind,omitempty"`
	Purl        string   `json:"purl,omitempty"`
	Cpes        []string `json:"cpes,omitempty"`
	LicensesExp string   `json:"licensesExp,omitempty"`
	Group       string   `json:"group,omitempty"`
	Description string   `json:"description,omitempty"`
	Primary     bool     `json:"primary"`
	Internal    bool     `json:"internal"`
	SbomID      string   `json:"sbomId"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// SbomComponentConnection represents a paginated list of components
type SbomComponentConnection struct {
	Nodes      []SbomComponent `json:"nodes"`
	TotalCount int             `json:"totalCount"`
	PageInfo   PageInfo        `json:"pageInfo"`
}

// Vuln represents a vulnerability
type Vuln struct {
	ID             string       `json:"id"`
	VulnID         string       `json:"vulnId"`
	Description    string       `json:"desc,omitempty"`
	Severity       string       `json:"sev,omitempty"`
	CvssScore      float64      `json:"cvssScore,omitempty"`
	CvssVector     string       `json:"cvssVector,omitempty"`
	Source         string       `json:"source,omitempty"`
	PublishedAt    FlexibleTime `json:"publishedAt,omitempty"`
	LastModifiedAt FlexibleTime `json:"lastModifiedAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt"`
	VulnInfo       *VulnInfo `json:"vulnInfo,omitempty"`
}

// VulnInfo contains additional vulnerability information
type VulnInfo struct {
	ID             string   `json:"id"`
	CveID          string   `json:"cveId"`
	EpssScore      float64  `json:"epssScore"`
	EpssPercentile float64  `json:"epssPercentile,omitempty"`
	Kev            bool     `json:"kev"`
	Cwes           []string `json:"cwes,omitempty"`
	Advisories     []string `json:"advisories,omitempty"`
}

// ComponentVuln represents a vulnerability associated with a component
type ComponentVuln struct {
	ID              string        `json:"id"`
	ComponentID     string        `json:"componentId"`
	VulnID          string        `json:"vulnId"`
	SbomID          string        `json:"sbomId,omitempty"`
	FixedIn         string        `json:"fixedIn,omitempty"`
	FixedVersions   []string      `json:"fixedVersions,omitempty"`
	Detail          string        `json:"detail,omitempty"`
	Impact          string        `json:"impact,omitempty"`
	ActionStmt      string        `json:"actionStmt,omitempty"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
	Component       *SbomComponent `json:"component,omitempty"`
	Vuln            *Vuln         `json:"vuln,omitempty"`
	VexStatus       *VexStatus    `json:"vexStatus,omitempty"`
	VexJustification *VexJustification `json:"vexJustification,omitempty"`
}

// ComponentVulnConnection represents a paginated list of component vulnerabilities
type ComponentVulnConnection struct {
	Nodes      []ComponentVuln `json:"nodes"`
	TotalCount int             `json:"totalCount"`
	PageInfo   PageInfo        `json:"pageInfo"`
}

// VexStatus represents a VEX status
type VexStatus struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// VexJustification represents a VEX justification
type VexJustification struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Policy represents a security policy
type Policy struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	IsEnabled   bool         `json:"isEnabled"`
	ResultType  string       `json:"resultType,omitempty"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	PolicyRules []PolicyRule `json:"policyRules,omitempty"`
}

// PolicyRule represents a rule within a policy
type PolicyRule struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Subject  string `json:"subject"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// PolicyConnection represents a paginated list of policies
type PolicyConnection struct {
	Nodes      []Policy `json:"nodes"`
	TotalCount int      `json:"totalCount"`
	PageInfo   PageInfo `json:"pageInfo"`
}

// PolicyResult represents a policy evaluation result
type PolicyResult struct {
	ID         string    `json:"id"`
	PolicyID   string    `json:"policyId"`
	SbomID     string    `json:"sbomId"`
	ResultType string    `json:"resultType"`
	Result     string    `json:"result,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	Policy     *Policy   `json:"policy,omitempty"`
	Sbom       *Sbom     `json:"sbom,omitempty"`
}

// PolicyResultConnection represents a paginated list of policy results
type PolicyResultConnection struct {
	Nodes      []PolicyResult `json:"nodes"`
	TotalCount int            `json:"totalCount"`
	PageInfo   PageInfo       `json:"pageInfo"`
}

// OrganizationLicense represents a license in the organization
type OrganizationLicense struct {
	Content            LicenseContent `json:"content"`
	DerivedState       string         `json:"derivedState,omitempty"`
	CopyLeft           string         `json:"copyLeft,omitempty"`
	OsiApproved        bool           `json:"osiApproved,omitempty"`
	FsfLibre           bool           `json:"fsfLibre,omitempty"`
	Deprecated         bool           `json:"deprecated,omitempty"`
	Attribution        string         `json:"attribution,omitempty"`
	SourceDistribution string         `json:"sourceDistribution,omitempty"`
	Modifications      string         `json:"modifications,omitempty"`
}

// LicenseContent represents license content (either standard or custom)
type LicenseContent struct {
	ShortID string `json:"shortId,omitempty"` // For standard License
	SpdxID  string `json:"spdxId,omitempty"`  // For LicenseCustom
	Name    string `json:"name,omitempty"`
}

// LicenseConnection represents a paginated list of licenses
type LicenseConnection struct {
	Nodes      []OrganizationLicense `json:"nodes"`
	TotalCount int                   `json:"totalCount"`
	PageInfo   PageInfo              `json:"pageInfo"`
}

// SbomDiff represents a diff between two SBOMs
type SbomDiff struct {
	DiffType           string         `json:"diffType"`
	DiffTags           []string       `json:"diffTags,omitempty"`
	SubjectComponent   *SbomComponent `json:"subjectComponent,omitempty"`
	SubjectComponentID string         `json:"subjectComponentId,omitempty"`
	TargetComponent    *SbomComponent `json:"targetComponent,omitempty"`
	TargetComponentID  string         `json:"targetComponentId,omitempty"`
}

// PageInfo contains pagination information
type PageInfo struct {
	HasNextPage     bool   `json:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
	StartCursor     string `json:"startCursor,omitempty"`
	EndCursor       string `json:"endCursor,omitempty"`
}
