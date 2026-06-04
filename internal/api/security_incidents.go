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

// SecurityIncident represents a supply-chain security incident.
type SecurityIncident struct {
	ID                 string
	Title              string
	Slug               string
	Summary            string
	Severity           string
	Status             string
	Confidence         string
	RecommendedActions string
	SourceURLs         string
	FirstSeenAt        *time.Time
	PublishedAt        *time.Time
	LastUpdatedAt      *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Markers            []SecurityIncidentMarker
	OrgImpactState     *OrganizationSecurityIncidentState
}

// SecurityIncidentMarker represents an incident indicator or affected component identity.
type SecurityIncidentMarker struct {
	ID               string
	MarkerType       string
	Purl             string
	ComponentName    string
	ComponentVersion string
	GithubURL        string
	Active           bool
	AddedAt          time.Time
	WithdrawnAt      *time.Time
}

// OrganizationSecurityIncidentState contains the current organization's impact rollup.
type OrganizationSecurityIncidentState struct {
	Status                  string
	Severity                string
	ImpactedProjectsCount   int
	ImpactedVersionsCount   int
	ImpactedComponentsCount int
	LastEvaluatedAt         *time.Time
	LastNotifiedAt          *time.Time
}

// ListSecurityIncidentsInput contains filters for listing incidents.
type ListSecurityIncidentsInput struct {
	Status []string
}

// CreateSecurityIncidentInput contains fields for creating an incident.
type CreateSecurityIncidentInput struct {
	Title              string
	Severity           string
	Confidence         *string
	Summary            *string
	RecommendedActions *string
	SourceURLs         *string
}

// UpdateSecurityIncidentInput contains editable incident fields.
type UpdateSecurityIncidentInput struct {
	ID                 string
	Title              *string
	Severity           *string
	Confidence         *string
	Summary            *string
	RecommendedActions *string
	SourceURLs         *string
}

// SecurityIncidentMarkerInput contains fields for creating a marker.
type SecurityIncidentMarkerInput struct {
	MarkerType       string `json:"markerType"`
	Purl             string `json:"purl,omitempty"`
	ComponentName    string `json:"componentName,omitempty"`
	ComponentVersion string `json:"componentVersion,omitempty"`
	GithubURL        string `json:"githubUrl,omitempty"`
}

// CreateSecurityIncidentUpdateInput contains fields for adding an incident timeline update.
type CreateSecurityIncidentUpdateInput struct {
	SecurityIncidentID string
	Title              string
	UpdateType         string
	OccurredAt         string
	Body               *string
	CustomerVisible    *bool
}

// SecurityIncidentFindingsInput contains filters for incident findings.
type SecurityIncidentFindingsInput struct {
	IncidentID string
	Statuses   []string
}

// SuppressSecurityIncidentFindingInput contains fields for suppressing a finding.
type SuppressSecurityIncidentFindingInput struct {
	FindingID string
	Reason    string
}

// CreateSecurityIncidentResult contains the create mutation result.
type CreateSecurityIncidentResult struct {
	SecurityIncident *SecurityIncident
	Errors           []string
}

// AddSecurityIncidentMarkersResult contains the marker mutation result.
type AddSecurityIncidentMarkersResult struct {
	Markers []SecurityIncidentMarker
	Errors  []string
}

// SecurityIncidentMarkersResult contains marker mutations.
type SecurityIncidentMarkersResult struct {
	Markers []SecurityIncidentMarker
	Errors  []string
}

// SecurityIncidentMutationResult contains mutations that return an incident.
type SecurityIncidentMutationResult struct {
	SecurityIncident *SecurityIncident
	Errors           []string
}

// SecurityIncidentUpdate represents a timeline update for an incident.
type SecurityIncidentUpdate struct {
	Title      string
	UpdateType string
	Body       string
	OccurredAt *time.Time
}

// CreateSecurityIncidentUpdateResult contains the timeline update mutation result.
type CreateSecurityIncidentUpdateResult struct {
	Update *SecurityIncidentUpdate
	Errors []string
}

// SecurityIncidentFinding represents a customer-facing incident finding.
type SecurityIncidentFinding struct {
	ID              string
	Status          string
	MatchMethod     string
	MatchedFields   map[string]interface{}
	FirstDetectedAt time.Time
	LastConfirmedAt *time.Time
	IsPartSbom      bool
	Component       *SecurityIncidentFindingComponent
	RootSbom        *SecurityIncidentFindingSbom
}

// SecurityIncidentFindingComponent contains the affected component in a finding.
type SecurityIncidentFindingComponent struct {
	ID          string
	Name        string
	Version     string
	Kind        string
	Purl        string
	Cpes        []string
	LicensesExp string
	Group       string
	Primary     bool
	Internal    bool
	SbomID      string
	UpdatedAt   time.Time
	Sbom        *SecurityIncidentFindingSbom
}

// SecurityIncidentFindingSbom contains root/component SBOM context for a finding.
type SecurityIncidentFindingSbom struct {
	ID             string
	ProjectVersion string
	Project        *SecurityIncidentFindingProject
}

// SecurityIncidentFindingProject contains project context for a finding.
type SecurityIncidentFindingProject struct {
	ID             string
	Name           string
	ProjectGroupID string
}

// SecurityIncidentFindingsResult contains findings for one incident.
type SecurityIncidentFindingsResult struct {
	IncidentID string
	Title      string
	Slug       string
	Status     string
	Severity   string
	Findings   []SecurityIncidentFinding
}

// SuppressSecurityIncidentFindingResult contains the suppression mutation result.
type SuppressSecurityIncidentFindingResult struct {
	Finding *SecurityIncidentFinding
	Errors  []string
}

// DryRunSecurityIncidentImpactScanResult contains a dry-run scan queue result.
type DryRunSecurityIncidentImpactScanResult struct {
	Status string
	Errors []string
}

// SecurityIncidentDryRunResultInput contains filters for reading dry-run results.
type SecurityIncidentDryRunResultInput struct {
	IncidentID string
	OrgID      string
	First      int
	After      string
}

// SecurityIncidentDryRunResult contains dry-run scan status and optional org findings.
type SecurityIncidentDryRunResult struct {
	Status            string
	Error             string
	CompletedAt       *time.Time
	TotalOrgsImpacted int
	OrgResults        []SecurityIncidentDryRunOrgResult
	Org               *SecurityIncidentDryRunOrg
}

// SecurityIncidentDryRunOrgResult contains per-organization dry-run rollup counts.
type SecurityIncidentDryRunOrgResult struct {
	OrganizationID          string
	OrganizationName        string
	ImpactedComponentsCount int
	ImpactedProjectsCount   int
	ImpactedVersionsCount   int
}

// SecurityIncidentDryRunOrg contains per-organization dry-run findings.
type SecurityIncidentDryRunOrg struct {
	SecurityIncidentDryRunOrgResult
	Findings    []SecurityIncidentDryRunFinding
	HasNextPage bool
	EndCursor   string
}

// SecurityIncidentDryRunFinding represents a single dry-run component hit.
type SecurityIncidentDryRunFinding struct {
	ID                     string
	OrganizationID         string
	ProjectID              string
	RootSbomID             string
	ComponentSbomID        string
	ComponentID            string
	ComponentName          string
	ComponentVersion       string
	ComponentPurl          string
	MarkerID               string
	MarkerType             string
	MarkerPurl             string
	MarkerComponentName    string
	MarkerComponentVersion string
	MatchMethod            string
	MatchedFields          map[string]interface{}
	RootProjectName        string
	RootProjectVersion     string
	IsPartSbom             bool
}

type securityIncidentGraphQL struct {
	ID                 string                          `json:"id"`
	Title              string                          `json:"title"`
	Slug               string                          `json:"slug"`
	Summary            string                          `json:"summary"`
	Severity           string                          `json:"severity"`
	Status             string                          `json:"status"`
	Confidence         string                          `json:"confidence"`
	RecommendedActions string                          `json:"recommendedActions"`
	SourceURLs         string                          `json:"sourceUrls"`
	FirstSeenAt        *time.Time                      `json:"firstSeenAt"`
	PublishedAt        *time.Time                      `json:"publishedAt"`
	LastUpdatedAt      *time.Time                      `json:"lastUpdatedAt"`
	CreatedAt          time.Time                       `json:"createdAt"`
	UpdatedAt          time.Time                       `json:"updatedAt"`
	Markers            []securityIncidentMarkerGraphQL `json:"markers"`
	OrgImpactState     *struct {
		Status                  string     `json:"status"`
		Severity                string     `json:"severity"`
		ImpactedProjectsCount   int        `json:"impactedProjectsCount"`
		ImpactedVersionsCount   int        `json:"impactedVersionsCount"`
		ImpactedComponentsCount int        `json:"impactedComponentsCount"`
		LastEvaluatedAt         *time.Time `json:"lastEvaluatedAt"`
		LastNotifiedAt          *time.Time `json:"lastNotifiedAt"`
	} `json:"orgImpactState"`
}

type securityIncidentMarkerGraphQL struct {
	ID               string     `json:"id"`
	MarkerType       string     `json:"markerType"`
	Purl             string     `json:"purl"`
	ComponentName    string     `json:"componentName"`
	ComponentVersion string     `json:"componentVersion"`
	GithubURL        string     `json:"githubUrl"`
	Active           bool       `json:"active"`
	AddedAt          time.Time  `json:"addedAt"`
	WithdrawnAt      *time.Time `json:"withdrawnAt"`
}

type securityIncidentUpdateGraphQL struct {
	Title      string     `json:"title"`
	UpdateType string     `json:"updateType"`
	Body       string     `json:"body"`
	OccurredAt *time.Time `json:"occurredAt"`
}

type securityIncidentFindingGraphQL struct {
	ID              string                 `json:"id"`
	Status          string                 `json:"status"`
	MatchMethod     string                 `json:"matchMethod"`
	MatchedFields   map[string]interface{} `json:"matchedFields"`
	FirstDetectedAt time.Time              `json:"firstDetectedAt"`
	LastConfirmedAt *time.Time             `json:"lastConfirmedAt"`
	IsPartSbom      bool                   `json:"isPartSbom"`
	Component       *struct {
		ID          string          `json:"id"`
		Name        string          `json:"name"`
		Version     string          `json:"version"`
		Kind        string          `json:"kind"`
		Purl        string          `json:"purl"`
		Cpes        []string        `json:"cpes"`
		LicensesExp string          `json:"licensesExp"`
		Group       string          `json:"group"`
		Primary     bool            `json:"primary"`
		Internal    bool            `json:"internal"`
		SbomID      string          `json:"sbomId"`
		UpdatedAt   time.Time       `json:"updatedAt"`
		Sbom        *sbomRefGraphQL `json:"sbom"`
	} `json:"component"`
	RootSbom *sbomRefGraphQL `json:"rootSbom"`
}

type sbomRefGraphQL struct {
	ID             string             `json:"id"`
	ProjectVersion string             `json:"projectVersion"`
	Project        *projectRefGraphQL `json:"project"`
}

type projectRefGraphQL struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ProjectGroupID string `json:"projectGroupId"`
}

// ListSecurityIncidents fetches incidents visible to the current organization.
func (c *Client) ListSecurityIncidents(ctx context.Context, input ListSecurityIncidentsInput) ([]SecurityIncident, error) {
	vars := make(map[string]interface{})
	if len(input.Status) > 0 {
		vars["status"] = input.Status
	}

	var result struct {
		SecurityIncidents []securityIncidentGraphQL `json:"securityIncidents"`
	}

	if err := c.gql.Execute(ctx, graphql.SecurityIncidentsQuery, vars, &result); err != nil {
		return nil, err
	}

	incidents := make([]SecurityIncident, len(result.SecurityIncidents))
	for i, incident := range result.SecurityIncidents {
		incidents[i] = mapSecurityIncident(incident)
	}
	return incidents, nil
}

// GetSecurityIncident fetches a single incident by ID.
func (c *Client) GetSecurityIncident(ctx context.Context, id string) (*SecurityIncident, error) {
	var result struct {
		SecurityIncident *securityIncidentGraphQL `json:"securityIncident"`
	}

	if err := c.gql.Execute(ctx, graphql.SecurityIncidentQuery, map[string]interface{}{"id": id}, &result); err != nil {
		return nil, err
	}
	if result.SecurityIncident == nil {
		return nil, nil
	}
	incident := mapSecurityIncident(*result.SecurityIncident)
	return &incident, nil
}

// CreateSecurityIncident creates a draft incident.
func (c *Client) CreateSecurityIncident(ctx context.Context, input CreateSecurityIncidentInput) (*CreateSecurityIncidentResult, error) {
	vars := map[string]interface{}{
		"title":    input.Title,
		"severity": input.Severity,
	}
	addStringVar(vars, "confidence", input.Confidence)
	addStringVar(vars, "summary", input.Summary)
	addStringVar(vars, "recommendedActions", input.RecommendedActions)
	addStringVar(vars, "sourceUrls", input.SourceURLs)

	var result struct {
		CreateSecurityIncident struct {
			SecurityIncident *securityIncidentGraphQL `json:"securityIncident"`
			Errors           []string                 `json:"errors"`
		} `json:"createSecurityIncident"`
	}

	if err := c.gql.Execute(ctx, graphql.SecurityIncidentCreateMutation, vars, &result); err != nil {
		return nil, err
	}

	mutation := result.CreateSecurityIncident
	createResult := &CreateSecurityIncidentResult{Errors: mutation.Errors}
	if mutation.SecurityIncident != nil {
		incident := mapSecurityIncident(*mutation.SecurityIncident)
		createResult.SecurityIncident = &incident
	}
	return createResult, nil
}

// UpdateSecurityIncident updates editable incident fields.
func (c *Client) UpdateSecurityIncident(ctx context.Context, input UpdateSecurityIncidentInput) (*SecurityIncidentMutationResult, error) {
	vars := map[string]interface{}{"id": input.ID}
	addStringVar(vars, "title", input.Title)
	addStringVar(vars, "severity", input.Severity)
	addStringVar(vars, "confidence", input.Confidence)
	addStringVar(vars, "summary", input.Summary)
	addStringVar(vars, "recommendedActions", input.RecommendedActions)
	addStringVar(vars, "sourceUrls", input.SourceURLs)

	return c.securityIncidentMutationWithVars(ctx, graphql.SecurityIncidentUpdateMutation, "updateSecurityIncident", vars)
}

// GetSecurityIncidentFindings fetches customer-facing findings for the current organization.
func (c *Client) GetSecurityIncidentFindings(ctx context.Context, input SecurityIncidentFindingsInput) (*SecurityIncidentFindingsResult, error) {
	vars := map[string]interface{}{"id": input.IncidentID}
	if len(input.Statuses) > 0 {
		vars["statuses"] = input.Statuses
	}

	var result struct {
		SecurityIncident *struct {
			ID       string                           `json:"id"`
			Title    string                           `json:"title"`
			Slug     string                           `json:"slug"`
			Status   string                           `json:"status"`
			Severity string                           `json:"severity"`
			Findings []securityIncidentFindingGraphQL `json:"findings"`
		} `json:"securityIncident"`
	}

	if err := c.gql.Execute(ctx, graphql.SecurityIncidentFindingsQuery, vars, &result); err != nil {
		return nil, err
	}
	if result.SecurityIncident == nil {
		return nil, nil
	}

	findings := make([]SecurityIncidentFinding, len(result.SecurityIncident.Findings))
	for i, finding := range result.SecurityIncident.Findings {
		findings[i] = mapSecurityIncidentFinding(finding)
	}

	return &SecurityIncidentFindingsResult{
		IncidentID: result.SecurityIncident.ID,
		Title:      result.SecurityIncident.Title,
		Slug:       result.SecurityIncident.Slug,
		Status:     result.SecurityIncident.Status,
		Severity:   result.SecurityIncident.Severity,
		Findings:   findings,
	}, nil
}

// AddSecurityIncidentMarkers adds markers to an incident.
func (c *Client) AddSecurityIncidentMarkers(ctx context.Context, incidentID string, markers []SecurityIncidentMarkerInput) (*AddSecurityIncidentMarkersResult, error) {
	var result struct {
		AddSecurityIncidentMarkers struct {
			Markers []securityIncidentMarkerGraphQL `json:"markers"`
			Errors  []string                        `json:"errors"`
		} `json:"addSecurityIncidentMarkers"`
	}

	vars := map[string]interface{}{
		"securityIncidentId": incidentID,
		"markers":            markers,
	}
	if err := c.gql.Execute(ctx, graphql.SecurityIncidentAddMarkersMutation, vars, &result); err != nil {
		return nil, err
	}

	mutation := result.AddSecurityIncidentMarkers
	mapped := mapSecurityIncidentMarkers(mutation.Markers)

	return &AddSecurityIncidentMarkersResult{Markers: mapped, Errors: mutation.Errors}, nil
}

// WithdrawSecurityIncidentMarkers withdraws active markers from an incident.
func (c *Client) WithdrawSecurityIncidentMarkers(ctx context.Context, incidentID string, markerIDs []string) (*SecurityIncidentMarkersResult, error) {
	var result struct {
		WithdrawSecurityIncidentMarkers struct {
			Markers []securityIncidentMarkerGraphQL `json:"markers"`
			Errors  []string                        `json:"errors"`
		} `json:"withdrawSecurityIncidentMarkers"`
	}

	vars := map[string]interface{}{
		"securityIncidentId": incidentID,
		"markerIds":          markerIDs,
	}
	if err := c.gql.Execute(ctx, graphql.SecurityIncidentWithdrawMarkersMutation, vars, &result); err != nil {
		return nil, err
	}

	mutation := result.WithdrawSecurityIncidentMarkers
	return &SecurityIncidentMarkersResult{
		Markers: mapSecurityIncidentMarkers(mutation.Markers),
		Errors:  mutation.Errors,
	}, nil
}

// PublishSecurityIncident publishes a draft incident and queues impact scanning.
func (c *Client) PublishSecurityIncident(ctx context.Context, id string) (*SecurityIncidentMutationResult, error) {
	return c.securityIncidentMutation(ctx, graphql.SecurityIncidentPublishMutation, "publishSecurityIncident", id)
}

// ResolveSecurityIncident resolves an active incident.
func (c *Client) ResolveSecurityIncident(ctx context.Context, id string) (*SecurityIncidentMutationResult, error) {
	return c.securityIncidentMutation(ctx, graphql.SecurityIncidentResolveMutation, "resolveSecurityIncident", id)
}

// ArchiveSecurityIncident archives an incident.
func (c *Client) ArchiveSecurityIncident(ctx context.Context, id string) (*SecurityIncidentMutationResult, error) {
	return c.securityIncidentMutation(ctx, graphql.SecurityIncidentArchiveMutation, "archiveSecurityIncident", id)
}

// CreateSecurityIncidentUpdate adds a timeline update to an incident.
func (c *Client) CreateSecurityIncidentUpdate(ctx context.Context, input CreateSecurityIncidentUpdateInput) (*CreateSecurityIncidentUpdateResult, error) {
	vars := map[string]interface{}{
		"securityIncidentId": input.SecurityIncidentID,
		"title":              input.Title,
		"updateType":         input.UpdateType,
		"occurredAt":         input.OccurredAt,
	}
	addStringVar(vars, "body", input.Body)
	addBoolVar(vars, "customerVisible", input.CustomerVisible)

	var result struct {
		CreateSecurityIncidentUpdate struct {
			Update *securityIncidentUpdateGraphQL `json:"securityIncidentUpdate"`
			Errors []string                       `json:"errors"`
		} `json:"createSecurityIncidentUpdate"`
	}

	if err := c.gql.Execute(ctx, graphql.SecurityIncidentCreateUpdateMutation, vars, &result); err != nil {
		return nil, err
	}

	mutation := result.CreateSecurityIncidentUpdate
	updateResult := &CreateSecurityIncidentUpdateResult{Errors: mutation.Errors}
	if mutation.Update != nil {
		update := mapSecurityIncidentUpdate(*mutation.Update)
		updateResult.Update = &update
	}
	return updateResult, nil
}

// SuppressSecurityIncidentFinding suppresses a finding for the current organization.
func (c *Client) SuppressSecurityIncidentFinding(ctx context.Context, input SuppressSecurityIncidentFindingInput) (*SuppressSecurityIncidentFindingResult, error) {
	var result struct {
		SuppressSecurityIncidentFinding struct {
			Finding *securityIncidentFindingGraphQL `json:"finding"`
			Errors  []string                        `json:"errors"`
		} `json:"suppressSecurityIncidentFinding"`
	}

	vars := map[string]interface{}{
		"findingId": input.FindingID,
		"reason":    input.Reason,
	}
	if err := c.gql.Execute(ctx, graphql.SecurityIncidentSuppressFindingMutation, vars, &result); err != nil {
		return nil, err
	}

	mutation := result.SuppressSecurityIncidentFinding
	suppressResult := &SuppressSecurityIncidentFindingResult{Errors: mutation.Errors}
	if mutation.Finding != nil {
		finding := mapSecurityIncidentFinding(*mutation.Finding)
		suppressResult.Finding = &finding
	}
	return suppressResult, nil
}

// RerunSecurityIncidentImpactScan queues impact scanning for an active or resolved incident.
func (c *Client) RerunSecurityIncidentImpactScan(ctx context.Context, id string) (*SecurityIncidentMutationResult, error) {
	return c.securityIncidentMutation(ctx, graphql.SecurityIncidentRerunImpactScanMutation, "rerunSecurityIncidentImpactScan", id)
}

// DryRunSecurityIncidentImpactScan queues a dry-run impact scan.
func (c *Client) DryRunSecurityIncidentImpactScan(ctx context.Context, id string) (*DryRunSecurityIncidentImpactScanResult, error) {
	var result struct {
		DryRunSecurityIncidentImpactScan struct {
			Status string   `json:"status"`
			Errors []string `json:"errors"`
		} `json:"dryRunSecurityIncidentImpactScan"`
	}

	if err := c.gql.Execute(ctx, graphql.SecurityIncidentDryRunImpactScanMutation, map[string]interface{}{"id": id}, &result); err != nil {
		return nil, err
	}

	mutation := result.DryRunSecurityIncidentImpactScan
	return &DryRunSecurityIncidentImpactScanResult{Status: mutation.Status, Errors: mutation.Errors}, nil
}

// GetSecurityIncidentDryRunResult fetches the latest dry-run scan result.
func (c *Client) GetSecurityIncidentDryRunResult(ctx context.Context, input SecurityIncidentDryRunResultInput) (*SecurityIncidentDryRunResult, error) {
	vars := map[string]interface{}{"incidentId": input.IncidentID}
	query := graphql.SecurityIncidentDryRunResultQuery
	if input.OrgID != "" {
		query = graphql.SecurityIncidentDryRunOrgResultQuery
		vars["orgId"] = input.OrgID
		if input.First > 0 {
			vars["first"] = input.First
		} else {
			vars["first"] = 50
		}
		if input.After != "" {
			vars["after"] = input.After
		}
	}

	var result struct {
		SecurityIncidentDryRunResult dryRunResultGraphQL `json:"securityIncidentDryRunResult"`
	}

	if err := c.gql.Execute(ctx, query, vars, &result); err != nil {
		return nil, err
	}

	mapped := mapSecurityIncidentDryRunResult(result.SecurityIncidentDryRunResult)
	return &mapped, nil
}

func (c *Client) securityIncidentMutation(ctx context.Context, query, fieldName, id string) (*SecurityIncidentMutationResult, error) {
	return c.securityIncidentMutationWithVars(ctx, query, fieldName, map[string]interface{}{"id": id})
}

func (c *Client) securityIncidentMutationWithVars(ctx context.Context, query, fieldName string, vars map[string]interface{}) (*SecurityIncidentMutationResult, error) {
	var result map[string]struct {
		SecurityIncident *securityIncidentGraphQL `json:"securityIncident"`
		Errors           []string                 `json:"errors"`
	}

	if err := c.gql.Execute(ctx, query, vars, &result); err != nil {
		return nil, err
	}

	mutation := result[fieldName]
	mutationResult := &SecurityIncidentMutationResult{Errors: mutation.Errors}
	if mutation.SecurityIncident != nil {
		incident := mapSecurityIncident(*mutation.SecurityIncident)
		mutationResult.SecurityIncident = &incident
	}
	return mutationResult, nil
}

type dryRunResultGraphQL struct {
	Status            string                   `json:"status"`
	Error             string                   `json:"error"`
	CompletedAt       *time.Time               `json:"completedAt"`
	TotalOrgsImpacted int                      `json:"totalOrgsImpacted"`
	OrgResults        []dryRunOrgResultGraphQL `json:"orgResults"`
	Org               *dryRunOrgGraphQL        `json:"org"`
}

type dryRunOrgResultGraphQL struct {
	OrganizationID          string `json:"organizationId"`
	OrganizationName        string `json:"organizationName"`
	ImpactedComponentsCount int    `json:"impactedComponentsCount"`
	ImpactedProjectsCount   int    `json:"impactedProjectsCount"`
	ImpactedVersionsCount   int    `json:"impactedVersionsCount"`
}

type dryRunOrgGraphQL struct {
	OrganizationID          string `json:"organizationId"`
	OrganizationName        string `json:"organizationName"`
	ImpactedComponentsCount int    `json:"impactedComponentsCount"`
	ImpactedProjectsCount   int    `json:"impactedProjectsCount"`
	ImpactedVersionsCount   int    `json:"impactedVersionsCount"`
	Findings                struct {
		Nodes []struct {
			ID                     string                 `json:"id"`
			OrganizationID         string                 `json:"organizationId"`
			ProjectID              string                 `json:"projectId"`
			RootSbomID             string                 `json:"rootSbomId"`
			ComponentSbomID        string                 `json:"componentSbomId"`
			ComponentID            string                 `json:"componentId"`
			ComponentName          string                 `json:"componentName"`
			ComponentVersion       string                 `json:"componentVersion"`
			ComponentPurl          string                 `json:"componentPurl"`
			MarkerID               string                 `json:"markerId"`
			MarkerType             string                 `json:"markerType"`
			MarkerPurl             string                 `json:"markerPurl"`
			MarkerComponentName    string                 `json:"markerComponentName"`
			MarkerComponentVersion string                 `json:"markerComponentVersion"`
			MatchMethod            string                 `json:"matchMethod"`
			MatchedFields          map[string]interface{} `json:"matchedFields"`
			RootProjectName        string                 `json:"rootProjectName"`
			RootProjectVersion     string                 `json:"rootProjectVersion"`
			IsPartSbom             bool                   `json:"isPartSbom"`
		} `json:"nodes"`
		PageInfo struct {
			HasNextPage bool   `json:"hasNextPage"`
			EndCursor   string `json:"endCursor"`
		} `json:"pageInfo"`
	} `json:"findings"`
}

func mapSecurityIncident(input securityIncidentGraphQL) SecurityIncident {
	markers := mapSecurityIncidentMarkers(input.Markers)

	var orgImpactState *OrganizationSecurityIncidentState
	if input.OrgImpactState != nil {
		orgImpactState = &OrganizationSecurityIncidentState{
			Status:                  input.OrgImpactState.Status,
			Severity:                input.OrgImpactState.Severity,
			ImpactedProjectsCount:   input.OrgImpactState.ImpactedProjectsCount,
			ImpactedVersionsCount:   input.OrgImpactState.ImpactedVersionsCount,
			ImpactedComponentsCount: input.OrgImpactState.ImpactedComponentsCount,
			LastEvaluatedAt:         input.OrgImpactState.LastEvaluatedAt,
			LastNotifiedAt:          input.OrgImpactState.LastNotifiedAt,
		}
	}

	return SecurityIncident{
		ID:                 input.ID,
		Title:              input.Title,
		Slug:               input.Slug,
		Summary:            input.Summary,
		Severity:           input.Severity,
		Status:             input.Status,
		Confidence:         input.Confidence,
		RecommendedActions: input.RecommendedActions,
		SourceURLs:         input.SourceURLs,
		FirstSeenAt:        input.FirstSeenAt,
		PublishedAt:        input.PublishedAt,
		LastUpdatedAt:      input.LastUpdatedAt,
		CreatedAt:          input.CreatedAt,
		UpdatedAt:          input.UpdatedAt,
		Markers:            markers,
		OrgImpactState:     orgImpactState,
	}
}

func mapSecurityIncidentMarkers(input []securityIncidentMarkerGraphQL) []SecurityIncidentMarker {
	markers := make([]SecurityIncidentMarker, len(input))
	for i, marker := range input {
		markers[i] = SecurityIncidentMarker{
			ID:               marker.ID,
			MarkerType:       marker.MarkerType,
			Purl:             marker.Purl,
			ComponentName:    marker.ComponentName,
			ComponentVersion: marker.ComponentVersion,
			GithubURL:        marker.GithubURL,
			Active:           marker.Active,
			AddedAt:          marker.AddedAt,
			WithdrawnAt:      marker.WithdrawnAt,
		}
	}
	return markers
}

func mapSecurityIncidentUpdate(input securityIncidentUpdateGraphQL) SecurityIncidentUpdate {
	return SecurityIncidentUpdate{
		Title:      input.Title,
		UpdateType: input.UpdateType,
		Body:       input.Body,
		OccurredAt: input.OccurredAt,
	}
}

func mapSecurityIncidentFinding(input securityIncidentFindingGraphQL) SecurityIncidentFinding {
	finding := SecurityIncidentFinding{
		ID:              input.ID,
		Status:          input.Status,
		MatchMethod:     input.MatchMethod,
		MatchedFields:   input.MatchedFields,
		FirstDetectedAt: input.FirstDetectedAt,
		LastConfirmedAt: input.LastConfirmedAt,
		IsPartSbom:      input.IsPartSbom,
		RootSbom:        mapSecurityIncidentFindingSbom(input.RootSbom),
	}

	if input.Component != nil {
		finding.Component = &SecurityIncidentFindingComponent{
			ID:          input.Component.ID,
			Name:        input.Component.Name,
			Version:     input.Component.Version,
			Kind:        input.Component.Kind,
			Purl:        input.Component.Purl,
			Cpes:        input.Component.Cpes,
			LicensesExp: input.Component.LicensesExp,
			Group:       input.Component.Group,
			Primary:     input.Component.Primary,
			Internal:    input.Component.Internal,
			SbomID:      input.Component.SbomID,
			UpdatedAt:   input.Component.UpdatedAt,
			Sbom:        mapSecurityIncidentFindingSbom(input.Component.Sbom),
		}
	}

	return finding
}

func mapSecurityIncidentFindingSbom(input *sbomRefGraphQL) *SecurityIncidentFindingSbom {
	if input == nil {
		return nil
	}
	return &SecurityIncidentFindingSbom{
		ID:             input.ID,
		ProjectVersion: input.ProjectVersion,
		Project:        mapSecurityIncidentFindingProject(input.Project),
	}
}

func mapSecurityIncidentFindingProject(input *projectRefGraphQL) *SecurityIncidentFindingProject {
	if input == nil {
		return nil
	}
	return &SecurityIncidentFindingProject{
		ID:             input.ID,
		Name:           input.Name,
		ProjectGroupID: input.ProjectGroupID,
	}
}

func mapSecurityIncidentDryRunResult(input dryRunResultGraphQL) SecurityIncidentDryRunResult {
	orgResults := make([]SecurityIncidentDryRunOrgResult, len(input.OrgResults))
	for i, orgResult := range input.OrgResults {
		orgResults[i] = mapSecurityIncidentDryRunOrgResult(orgResult)
	}

	var org *SecurityIncidentDryRunOrg
	if input.Org != nil {
		findings := make([]SecurityIncidentDryRunFinding, len(input.Org.Findings.Nodes))
		for i, finding := range input.Org.Findings.Nodes {
			findings[i] = SecurityIncidentDryRunFinding{
				ID:                     finding.ID,
				OrganizationID:         finding.OrganizationID,
				ProjectID:              finding.ProjectID,
				RootSbomID:             finding.RootSbomID,
				ComponentSbomID:        finding.ComponentSbomID,
				ComponentID:            finding.ComponentID,
				ComponentName:          finding.ComponentName,
				ComponentVersion:       finding.ComponentVersion,
				ComponentPurl:          finding.ComponentPurl,
				MarkerID:               finding.MarkerID,
				MarkerType:             finding.MarkerType,
				MarkerPurl:             finding.MarkerPurl,
				MarkerComponentName:    finding.MarkerComponentName,
				MarkerComponentVersion: finding.MarkerComponentVersion,
				MatchMethod:            finding.MatchMethod,
				MatchedFields:          finding.MatchedFields,
				RootProjectName:        finding.RootProjectName,
				RootProjectVersion:     finding.RootProjectVersion,
				IsPartSbom:             finding.IsPartSbom,
			}
		}
		org = &SecurityIncidentDryRunOrg{
			SecurityIncidentDryRunOrgResult: SecurityIncidentDryRunOrgResult{
				OrganizationID:          input.Org.OrganizationID,
				OrganizationName:        input.Org.OrganizationName,
				ImpactedComponentsCount: input.Org.ImpactedComponentsCount,
				ImpactedProjectsCount:   input.Org.ImpactedProjectsCount,
				ImpactedVersionsCount:   input.Org.ImpactedVersionsCount,
			},
			Findings:    findings,
			HasNextPage: input.Org.Findings.PageInfo.HasNextPage,
			EndCursor:   input.Org.Findings.PageInfo.EndCursor,
		}
	}

	return SecurityIncidentDryRunResult{
		Status:            input.Status,
		Error:             input.Error,
		CompletedAt:       input.CompletedAt,
		TotalOrgsImpacted: input.TotalOrgsImpacted,
		OrgResults:        orgResults,
		Org:               org,
	}
}

func mapSecurityIncidentDryRunOrgResult(input dryRunOrgResultGraphQL) SecurityIncidentDryRunOrgResult {
	return SecurityIncidentDryRunOrgResult{
		OrganizationID:          input.OrganizationID,
		OrganizationName:        input.OrganizationName,
		ImpactedComponentsCount: input.ImpactedComponentsCount,
		ImpactedProjectsCount:   input.ImpactedProjectsCount,
		ImpactedVersionsCount:   input.ImpactedVersionsCount,
	}
}
