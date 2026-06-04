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

// Policy represents a security policy
type Policy struct {
	ID           string
	Name         string
	Description  string
	Enabled      bool
	ResultType   string
	CreateTicket bool
	UpdatedAt    time.Time
	PolicyRules  []PolicyRule
}

// PolicyRule represents a rule within a policy
type PolicyRule struct {
	ID       string
	Name     string
	Subject  string
	Operator string
	Value    string
}

// PolicyResult represents a policy evaluation result
type PolicyResult struct {
	ID         string
	PolicyID   string
	VersionID  string
	ResultType string
	Result     string
	CreatedAt  time.Time
	Policy     *Policy
	Version    *Version
}

// PoliciesResult represents the result of listing policies
type PoliciesResult struct {
	Policies    []Policy
	TotalCount  int
	HasNextPage bool
	EndCursor   string
}

// PolicyResultsResult represents the result of listing policy results
type PolicyResultsResult struct {
	PolicyResults []PolicyResult
	TotalCount    int
	HasNextPage   bool
	EndCursor     string
}

// OrganizationLicense represents a license in the organization
type OrganizationLicense struct {
	ShortID            string
	Name               string
	DerivedState       string
	CopyLeft           string
	OsiApproved        bool
	FsfLibre           bool
	Deprecated         bool
	Attribution        string
	SourceDistribution string
	Modifications      string
}

// LicensesResult represents the result of listing licenses
type LicensesResult struct {
	Licenses    []OrganizationLicense
	TotalCount  int
	HasNextPage bool
	EndCursor   string
}

// TicketingStatusInput contains parameters for fetching ticketing visibility.
type TicketingStatusInput struct {
	ProductID     string
	ProductsFirst int
	PoliciesFirst int
	TicketsFirst  int
}

// TicketingStatus contains ticketing provider configuration and policy application details.
type TicketingStatus struct {
	Connections              []TicketingConnectionStatus
	JiraVulnManagementConfig *JiraVulnManagementConfig
	Products                 []TicketingProduct
	Policies                 []TicketingPolicy
	CreatedTickets           []CreatedTicket
	ProductsTotalCount       int
	ProductsHasNextPage      bool
	PoliciesTotalCount       int
	PoliciesHasNextPage      bool
	TicketsScannedCount      int
	TicketsHasNextPage       bool
}

// TicketingConnectionStatus contains organization-level ticketing provider connection health.
type TicketingConnectionStatus struct {
	Provider          string
	ConnectionID      string
	ProviderID        string
	Enabled           bool
	URL               string
	UserName          string
	HealthCheckStatus string
	LastHealthCheckAt time.Time
	UpdatedAt         time.Time
}

// JiraVulnManagementConfig contains Jira vulnerability-management provisioning status.
type JiraVulnManagementConfig struct {
	ID                 string
	Enabled            bool
	ProvisioningStatus string
	ProvisioningStep   string
	ProvisioningErrors string
	IssueTypeID        string
	WorkflowID         string
	ScreenID           string
	UpdatedAt          time.Time
}

// TicketingProduct contains repository and environment ticketing settings for a product.
type TicketingProduct struct {
	ID           string
	Name         string
	Enabled      bool
	Repository   *ImportedRepository
	Environments []TicketingEnvironment
}

// ImportedRepository contains source-control repository import status for a product.
type ImportedRepository struct {
	Type           string
	ID             string
	Name           string
	FullName       string
	Owner          string
	DefaultBranch  string
	Slug           string
	Workspace      string
	FullPath       string
	GitlabID       string
	ImportStatus   string
	WebhookEnabled bool
}

// TicketingEnvironment contains Jira sync settings for one environment.
type TicketingEnvironment struct {
	ID                    string
	Name                  string
	Enabled               bool
	IssueTrackerSettings  []ExternalIssueTrackerSetting
	AppliedTicketPolicies []TicketingPolicySummary
}

// ExternalIssueTrackerSetting contains per-environment issue tracker configuration.
type ExternalIssueTrackerSetting struct {
	ID             string
	Provider       string
	ProjectKey     string
	IssueType      string
	Assignee       string
	Reporter       string
	Epic           string
	Components     interface{}
	TeamID         string
	StateID        string
	EnableSync     bool
	LastSyncedAt   time.Time
	LastSyncStatus string
	UpdatedAt      time.Time
}

// TicketingPolicy contains policy ticketing status and application scope.
type TicketingPolicy struct {
	ID           string
	Name         string
	Enabled      bool
	ResultType   string
	CreateTicket bool
	Inclusions   []TicketingPolicyInclusion
}

// TicketingPolicySummary contains compact policy data attached to environments.
type TicketingPolicySummary struct {
	ID         string
	Name       string
	Enabled    bool
	ResultType string
}

// TicketingPolicyInclusion contains an environment included by a policy.
type TicketingPolicyInclusion struct {
	EnvironmentID   string
	EnvironmentName string
	ProductID       string
	ProductName     string
}

// CreatedTicket contains an issue tracker link created for a component vulnerability.
type CreatedTicket struct {
	ID               string
	Provider         string
	IssueKey         string
	IssueURL         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ComponentVulnID  string
	ComponentName    string
	ComponentVersion string
	VulnID           string
	VulnerabilityID  string
	Severity         string
	VersionID        string
	Version          string
	EnvironmentID    string
	EnvironmentName  string
	ProductID        string
	ProductName      string
}

// ListPoliciesInput contains parameters for listing policies
type ListPoliciesInput struct {
	First  int
	After  string
	Search string
}

// ListPolicies fetches policies with pagination
func (c *Client) ListPolicies(ctx context.Context, input ListPoliciesInput) (*PoliciesResult, error) {
	vars := make(map[string]interface{})
	if input.First > 0 {
		vars["first"] = input.First
	} else {
		vars["first"] = 20
	}
	if input.After != "" {
		vars["after"] = input.After
	}
	if input.Search != "" {
		vars["search"] = input.Search
	}

	var result struct {
		Policies struct {
			Nodes []struct {
				ID           string    `json:"id"`
				Name         string    `json:"name"`
				Description  string    `json:"description"`
				Enabled      bool      `json:"isEnabled"`
				ResultType   string    `json:"resultType"`
				CreateTicket bool      `json:"createTicket"`
				UpdatedAt    time.Time `json:"updatedAt"`
			} `json:"nodes"`
			TotalCount int `json:"totalCount"`
			PageInfo   struct {
				HasNextPage bool   `json:"hasNextPage"`
				EndCursor   string `json:"endCursor"`
			} `json:"pageInfo"`
		} `json:"policies"`
	}

	if err := c.gql.Execute(ctx, graphql.PoliciesQuery, vars, &result); err != nil {
		return nil, err
	}

	policies := make([]Policy, len(result.Policies.Nodes))
	for i, n := range result.Policies.Nodes {
		policies[i] = Policy{
			ID:           n.ID,
			Name:         n.Name,
			Description:  n.Description,
			Enabled:      n.Enabled,
			ResultType:   n.ResultType,
			CreateTicket: n.CreateTicket,
			UpdatedAt:    n.UpdatedAt,
		}
	}

	return &PoliciesResult{
		Policies:    policies,
		TotalCount:  result.Policies.TotalCount,
		HasNextPage: result.Policies.PageInfo.HasNextPage,
		EndCursor:   result.Policies.PageInfo.EndCursor,
	}, nil
}

// GetPolicy fetches a single policy by ID with its rules
func (c *Client) GetPolicy(ctx context.Context, id string) (*Policy, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var result struct {
		Policy struct {
			ID           string    `json:"id"`
			Name         string    `json:"name"`
			Description  string    `json:"description"`
			Enabled      bool      `json:"isEnabled"`
			ResultType   string    `json:"resultType"`
			CreateTicket bool      `json:"createTicket"`
			UpdatedAt    time.Time `json:"updatedAt"`
			PolicyRules  []struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Subject  string `json:"subject"`
				Operator string `json:"operator"`
				Value    string `json:"value"`
			} `json:"policyRules"`
		} `json:"policy"`
	}

	if err := c.gql.Execute(ctx, graphql.PolicyQuery, vars, &result); err != nil {
		return nil, err
	}

	rules := make([]PolicyRule, len(result.Policy.PolicyRules))
	for i, r := range result.Policy.PolicyRules {
		rules[i] = PolicyRule{
			ID:       r.ID,
			Name:     r.Name,
			Subject:  r.Subject,
			Operator: r.Operator,
			Value:    r.Value,
		}
	}

	return &Policy{
		ID:           result.Policy.ID,
		Name:         result.Policy.Name,
		Description:  result.Policy.Description,
		Enabled:      result.Policy.Enabled,
		ResultType:   result.Policy.ResultType,
		CreateTicket: result.Policy.CreateTicket,
		UpdatedAt:    result.Policy.UpdatedAt,
		PolicyRules:  rules,
	}, nil
}

// ListPolicyResultsInput contains parameters for listing policy results
type ListPolicyResultsInput struct {
	First      int
	After      string
	PolicyID   string
	VersionID  string
	ResultType string
}

// ListPolicyResults fetches policy evaluation results
func (c *Client) ListPolicyResults(ctx context.Context, input ListPolicyResultsInput) (*PolicyResultsResult, error) {
	vars := make(map[string]interface{})
	if input.First > 0 {
		vars["first"] = input.First
	} else {
		vars["first"] = 50
	}
	if input.After != "" {
		vars["after"] = input.After
	}
	if input.PolicyID != "" {
		vars["policyId"] = []string{input.PolicyID}
	}
	if input.VersionID != "" {
		vars["sbomId"] = []string{input.VersionID}
	}
	if input.ResultType != "" {
		vars["resultType"] = []string{input.ResultType}
	}

	var result struct {
		PolicyResults struct {
			Nodes []struct {
				ID         string    `json:"id"`
				PolicyID   string    `json:"policyId"`
				SbomID     string    `json:"sbomId"`
				ResultType string    `json:"resultType"`
				Result     string    `json:"result"`
				CreatedAt  time.Time `json:"createdAt"`
				Policy     *struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"policy"`
				Sbom *struct {
					ID             string `json:"id"`
					ProjectVersion string `json:"projectVersion"`
					Project        struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"project"`
				} `json:"sbom"`
			} `json:"nodes"`
			TotalCount int `json:"totalCount"`
			PageInfo   struct {
				HasNextPage bool   `json:"hasNextPage"`
				EndCursor   string `json:"endCursor"`
			} `json:"pageInfo"`
		} `json:"policyResults"`
	}

	if err := c.gql.Execute(ctx, graphql.PolicyResultsQuery, vars, &result); err != nil {
		return nil, err
	}

	results := make([]PolicyResult, len(result.PolicyResults.Nodes))
	for i, n := range result.PolicyResults.Nodes {
		pr := PolicyResult{
			ID:         n.ID,
			PolicyID:   n.PolicyID,
			VersionID:  n.SbomID,
			ResultType: n.ResultType,
			Result:     n.Result,
			CreatedAt:  n.CreatedAt,
		}
		if n.Policy != nil {
			pr.Policy = &Policy{
				ID:   n.Policy.ID,
				Name: n.Policy.Name,
			}
		}
		if n.Sbom != nil {
			pr.Version = &Version{
				ID:      n.Sbom.ID,
				Version: n.Sbom.ProjectVersion,
				Environment: &Environment{
					ID:   n.Sbom.Project.ID,
					Name: n.Sbom.Project.Name,
				},
			}
		}
		results[i] = pr
	}

	return &PolicyResultsResult{
		PolicyResults: results,
		TotalCount:    result.PolicyResults.TotalCount,
		HasNextPage:   result.PolicyResults.PageInfo.HasNextPage,
		EndCursor:     result.PolicyResults.PageInfo.EndCursor,
	}, nil
}

// GetTicketingStatus fetches Jira and ticketing policy application status.
func (c *Client) GetTicketingStatus(ctx context.Context, input TicketingStatusInput) (*TicketingStatus, error) {
	vars := map[string]interface{}{
		"policiesFirst": defaultPositive(input.PoliciesFirst, 50),
		"ticketsFirst":  defaultPositive(input.TicketsFirst, 500),
	}
	query := graphql.TicketingStatusQuery
	if input.ProductID != "" {
		query = graphql.ProductTicketingStatusQuery
		vars["productId"] = input.ProductID
	} else {
		vars["productsFirst"] = defaultPositive(input.ProductsFirst, 20)
	}

	var result struct {
		Organization struct {
			Connections struct {
				Nodes []ticketingConnectionNode `json:"nodes"`
			} `json:"connections"`
			ProjectGroups struct {
				Nodes      []ticketingProductNode `json:"nodes"`
				TotalCount int                    `json:"totalCount"`
				PageInfo   struct {
					HasNextPage bool `json:"hasNextPage"`
				} `json:"pageInfo"`
			} `json:"projectGroups"`
		} `json:"organization"`
		JiraVulnManagementConfig *ticketingJiraConfigNode `json:"jiraVulnManagementConfig"`
		ProjectGroup             *ticketingProductNode    `json:"projectGroup"`
		ComponentVulns           ticketingComponentVulns  `json:"componentVulns"`
		Policies                 struct {
			Nodes      []ticketingPolicyNode `json:"nodes"`
			TotalCount int                   `json:"totalCount"`
			PageInfo   struct {
				HasNextPage bool `json:"hasNextPage"`
			} `json:"pageInfo"`
		} `json:"policies"`
	}

	if err := c.gql.Execute(ctx, query, vars, &result); err != nil {
		return nil, err
	}

	status := &TicketingStatus{
		Connections:         mapTicketingConnections(result.Organization.Connections.Nodes),
		ProductsTotalCount:  result.Organization.ProjectGroups.TotalCount,
		ProductsHasNextPage: result.Organization.ProjectGroups.PageInfo.HasNextPage,
		PoliciesTotalCount:  result.Policies.TotalCount,
		PoliciesHasNextPage: result.Policies.PageInfo.HasNextPage,
	}
	if result.JiraVulnManagementConfig != nil {
		status.JiraVulnManagementConfig = mapJiraConfig(*result.JiraVulnManagementConfig)
	}

	productNodes := result.Organization.ProjectGroups.Nodes
	if input.ProductID != "" && result.ProjectGroup != nil {
		productNodes = []ticketingProductNode{*result.ProjectGroup}
		status.ProductsTotalCount = len(productNodes)
		status.ProductsHasNextPage = false
		status.CreatedTickets = mapCreatedTickets(result.ProjectGroup.ComponentVulns.Nodes)
		status.TicketsScannedCount = result.ProjectGroup.ComponentVulns.TotalCount
		status.TicketsHasNextPage = result.ProjectGroup.ComponentVulns.PageInfo.HasNextPage
	} else {
		status.CreatedTickets = mapCreatedTickets(result.ComponentVulns.Nodes)
		status.TicketsScannedCount = result.ComponentVulns.TotalCount
		status.TicketsHasNextPage = result.ComponentVulns.PageInfo.HasNextPage
	}

	productIDs := make(map[string]bool, len(productNodes))
	environmentPolicies := make(map[string][]TicketingPolicySummary)
	for _, product := range productNodes {
		productIDs[product.ID] = true
	}

	status.Policies = mapTicketingPolicies(result.Policies.Nodes, productIDs, input.ProductID != "", environmentPolicies)
	status.Products = mapTicketingProducts(productNodes, environmentPolicies)

	return status, nil
}

type ticketingConnectionNode struct {
	ID         string    `json:"id"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Connection struct {
		Type              string    `json:"__typename"`
		ID                string    `json:"id"`
		URL               string    `json:"url"`
		UserName          string    `json:"userName"`
		HealthCheckStatus string    `json:"healthCheckStatus"`
		LastHealthCheckAt time.Time `json:"lastHealthCheckAt"`
		UpdatedAt         time.Time `json:"updatedAt"`
	} `json:"connection"`
}

type ticketingJiraConfigNode struct {
	ID                 string    `json:"id"`
	Enabled            bool      `json:"enabled"`
	ProvisioningStatus string    `json:"provisioningStatus"`
	ProvisioningStep   string    `json:"provisioningStep"`
	ProvisioningErrors string    `json:"provisioningErrors"`
	IssueTypeID        string    `json:"issueTypeId"`
	WorkflowID         string    `json:"workflowId"`
	ScreenID           string    `json:"screenId"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type ticketingProductNode struct {
	ID                 string                   `json:"id"`
	Name               string                   `json:"name"`
	Enabled            bool                     `json:"enabled"`
	ImportedRepository *ticketingRepositoryNode `json:"importedRepository"`
	Projects           struct {
		Nodes []ticketingEnvironmentNode `json:"nodes"`
	} `json:"projects"`
	ComponentVulns ticketingComponentVulns `json:"componentVulns"`
}

type ticketingComponentVulns struct {
	Nodes      []ticketingComponentVulnNode `json:"nodes"`
	TotalCount int                          `json:"totalCount"`
	PageInfo   struct {
		HasNextPage bool `json:"hasNextPage"`
	} `json:"pageInfo"`
}

type ticketingComponentVulnNode struct {
	ID        string `json:"id"`
	Component *struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		SbomID  string `json:"sbomId"`
		Sbom    *struct {
			ID             string `json:"id"`
			ProjectVersion string `json:"projectVersion"`
			Project        *struct {
				ID           string `json:"id"`
				Name         string `json:"name"`
				ProjectGroup *struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"projectGroup"`
			} `json:"project"`
		} `json:"sbom"`
	} `json:"component"`
	Vuln *struct {
		ID     string `json:"id"`
		VulnID string `json:"vulnId"`
		Sev    string `json:"sev"`
	} `json:"vuln"`
	ExternalIssueTrackerLinks []struct {
		ID        string    `json:"id"`
		Provider  string    `json:"provider"`
		IssueKey  string    `json:"issueKey"`
		IssueURL  string    `json:"issueUrl"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	} `json:"externalIssueTrackerLinks"`
}

type ticketingRepositoryNode struct {
	Type           string `json:"__typename"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	FullName       string `json:"fullName"`
	Owner          string `json:"owner"`
	DefaultBranch  string `json:"defaultBranch"`
	Slug           string `json:"slug"`
	Workspace      string `json:"workspace"`
	FullPath       string `json:"fullPath"`
	GitlabID       string `json:"gitlabId"`
	ImportStatus   string `json:"importStatus"`
	WebhookEnabled bool   `json:"webhookEnabled"`
}

type ticketingEnvironmentNode struct {
	ID                           string `json:"id"`
	Name                         string `json:"name"`
	Enabled                      bool   `json:"enabled"`
	ExternalIssueTrackerSettings []struct {
		ID             string      `json:"id"`
		Provider       string      `json:"provider"`
		ProjectKey     string      `json:"projectKey"`
		IssueType      string      `json:"issueType"`
		Assignee       string      `json:"assignee"`
		Reporter       string      `json:"reporter"`
		Epic           string      `json:"epic"`
		Components     interface{} `json:"components"`
		TeamID         string      `json:"teamId"`
		StateID        string      `json:"stateId"`
		EnableJiraSync bool        `json:"enableJiraSync"`
		LastSyncedAt   time.Time   `json:"lastSyncedAt"`
		LastSyncStatus string      `json:"lastSyncStatus"`
		UpdatedAt      time.Time   `json:"updatedAt"`
	} `json:"externalIssueTrackerSettings"`
}

type ticketingPolicyNode struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Enabled          bool   `json:"isEnabled"`
	ResultType       string `json:"resultType"`
	CreateTicket     bool   `json:"createTicket"`
	PolicyInclusions []struct {
		ProjectID string `json:"projectId"`
		Project   struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			ProjectGroup struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"projectGroup"`
		} `json:"project"`
	} `json:"policyInclusions"`
}

func defaultPositive(value, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}

func mapTicketingConnections(nodes []ticketingConnectionNode) []TicketingConnectionStatus {
	connections := make([]TicketingConnectionStatus, 0, len(nodes))
	for _, node := range nodes {
		provider := ticketingProviderFromType(node.Connection.Type)
		if provider == "" {
			continue
		}
		connections = append(connections, TicketingConnectionStatus{
			Provider:          provider,
			ConnectionID:      node.ID,
			ProviderID:        node.Connection.ID,
			Enabled:           node.Enabled,
			URL:               node.Connection.URL,
			UserName:          node.Connection.UserName,
			HealthCheckStatus: node.Connection.HealthCheckStatus,
			LastHealthCheckAt: node.Connection.LastHealthCheckAt,
			UpdatedAt:         node.Connection.UpdatedAt,
		})
	}
	return connections
}

func ticketingProviderFromType(typeName string) string {
	switch typeName {
	case "JiraConnection":
		return "jira"
	case "LinearConnection":
		return "linear"
	default:
		return ""
	}
}

func mapJiraConfig(node ticketingJiraConfigNode) *JiraVulnManagementConfig {
	return &JiraVulnManagementConfig{
		ID:                 node.ID,
		Enabled:            node.Enabled,
		ProvisioningStatus: node.ProvisioningStatus,
		ProvisioningStep:   node.ProvisioningStep,
		ProvisioningErrors: node.ProvisioningErrors,
		IssueTypeID:        node.IssueTypeID,
		WorkflowID:         node.WorkflowID,
		ScreenID:           node.ScreenID,
		UpdatedAt:          node.UpdatedAt,
	}
}

func mapTicketingPolicies(nodes []ticketingPolicyNode, productIDs map[string]bool, filterByProduct bool, environmentPolicies map[string][]TicketingPolicySummary) []TicketingPolicy {
	policies := make([]TicketingPolicy, 0, len(nodes))
	for _, node := range nodes {
		policy := TicketingPolicy{
			ID:           node.ID,
			Name:         node.Name,
			Enabled:      node.Enabled,
			ResultType:   node.ResultType,
			CreateTicket: node.CreateTicket,
		}
		summary := TicketingPolicySummary{
			ID:         node.ID,
			Name:       node.Name,
			Enabled:    node.Enabled,
			ResultType: node.ResultType,
		}

		for _, inclusion := range node.PolicyInclusions {
			productID := inclusion.Project.ProjectGroup.ID
			if filterByProduct && !productIDs[productID] {
				continue
			}

			policy.Inclusions = append(policy.Inclusions, TicketingPolicyInclusion{
				EnvironmentID:   inclusion.ProjectID,
				EnvironmentName: inclusion.Project.Name,
				ProductID:       productID,
				ProductName:     inclusion.Project.ProjectGroup.Name,
			})
			if node.CreateTicket {
				environmentPolicies[inclusion.ProjectID] = append(environmentPolicies[inclusion.ProjectID], summary)
			}
		}

		if !filterByProduct || len(policy.Inclusions) > 0 {
			policies = append(policies, policy)
		}
	}
	return policies
}

func mapTicketingProducts(nodes []ticketingProductNode, environmentPolicies map[string][]TicketingPolicySummary) []TicketingProduct {
	products := make([]TicketingProduct, len(nodes))
	for i, node := range nodes {
		product := TicketingProduct{
			ID:      node.ID,
			Name:    node.Name,
			Enabled: node.Enabled,
		}
		if node.ImportedRepository != nil {
			product.Repository = &ImportedRepository{
				Type:           node.ImportedRepository.Type,
				ID:             node.ImportedRepository.ID,
				Name:           node.ImportedRepository.Name,
				FullName:       node.ImportedRepository.FullName,
				Owner:          node.ImportedRepository.Owner,
				DefaultBranch:  node.ImportedRepository.DefaultBranch,
				Slug:           node.ImportedRepository.Slug,
				Workspace:      node.ImportedRepository.Workspace,
				FullPath:       node.ImportedRepository.FullPath,
				GitlabID:       node.ImportedRepository.GitlabID,
				ImportStatus:   node.ImportedRepository.ImportStatus,
				WebhookEnabled: node.ImportedRepository.WebhookEnabled,
			}
		}

		product.Environments = make([]TicketingEnvironment, len(node.Projects.Nodes))
		for j, envNode := range node.Projects.Nodes {
			env := TicketingEnvironment{
				ID:                    envNode.ID,
				Name:                  envNode.Name,
				Enabled:               envNode.Enabled,
				AppliedTicketPolicies: environmentPolicies[envNode.ID],
			}
			for _, setting := range envNode.ExternalIssueTrackerSettings {
				env.IssueTrackerSettings = append(env.IssueTrackerSettings, ExternalIssueTrackerSetting{
					ID:             setting.ID,
					Provider:       setting.Provider,
					ProjectKey:     setting.ProjectKey,
					IssueType:      setting.IssueType,
					Assignee:       setting.Assignee,
					Reporter:       setting.Reporter,
					Epic:           setting.Epic,
					Components:     setting.Components,
					TeamID:         setting.TeamID,
					StateID:        setting.StateID,
					EnableSync:     setting.EnableJiraSync,
					LastSyncedAt:   setting.LastSyncedAt,
					LastSyncStatus: setting.LastSyncStatus,
					UpdatedAt:      setting.UpdatedAt,
				})
			}
			product.Environments[j] = env
		}
		products[i] = product
	}
	return products
}

func mapCreatedTickets(nodes []ticketingComponentVulnNode) []CreatedTicket {
	tickets := []CreatedTicket{}
	for _, node := range nodes {
		if len(node.ExternalIssueTrackerLinks) == 0 {
			continue
		}

		base := CreatedTicket{
			ComponentVulnID: node.ID,
		}
		if node.Component != nil {
			base.ComponentName = node.Component.Name
			base.ComponentVersion = node.Component.Version
			base.VersionID = node.Component.SbomID
			if node.Component.Sbom != nil {
				base.VersionID = node.Component.Sbom.ID
				base.Version = node.Component.Sbom.ProjectVersion
				if node.Component.Sbom.Project != nil {
					base.EnvironmentID = node.Component.Sbom.Project.ID
					base.EnvironmentName = node.Component.Sbom.Project.Name
					if node.Component.Sbom.Project.ProjectGroup != nil {
						base.ProductID = node.Component.Sbom.Project.ProjectGroup.ID
						base.ProductName = node.Component.Sbom.Project.ProjectGroup.Name
					}
				}
			}
		}
		if node.Vuln != nil {
			base.VulnerabilityID = node.Vuln.ID
			base.VulnID = node.Vuln.VulnID
			base.Severity = node.Vuln.Sev
		}

		for _, link := range node.ExternalIssueTrackerLinks {
			ticket := base
			ticket.ID = link.ID
			ticket.Provider = link.Provider
			ticket.IssueKey = link.IssueKey
			ticket.IssueURL = link.IssueURL
			ticket.CreatedAt = link.CreatedAt
			ticket.UpdatedAt = link.UpdatedAt
			tickets = append(tickets, ticket)
		}
	}
	return tickets
}

// ListLicensesInput contains parameters for listing licenses
type ListLicensesInput struct {
	First  int
	After  string
	Status string
	Search string
}

// ListLicenses fetches licenses with pagination
func (c *Client) ListLicenses(ctx context.Context, input ListLicensesInput) (*LicensesResult, error) {
	vars := make(map[string]interface{})
	if input.First > 0 {
		vars["first"] = input.First
	} else {
		vars["first"] = 50
	}
	if input.After != "" {
		vars["after"] = input.After
	}
	if input.Status != "" {
		vars["status"] = []string{input.Status}
	}
	if input.Search != "" {
		vars["search"] = input.Search
	}

	var result struct {
		Organization struct {
			Licenses struct {
				Nodes []struct {
					ID      string `json:"id"`
					Content struct {
						ShortID string `json:"shortId"` // For License type
						SpdxID  string `json:"spdxId"`  // For LicenseCustom type
						Name    string `json:"name"`
					} `json:"content"`
					State              string `json:"state"`
					CopyLeft           string `json:"copyLeft"`
					OsiApproved        bool   `json:"osiApproved"`
					FsfLibre           bool   `json:"fsfLibre"`
					Deprecated         bool   `json:"deprecated"`
					Attribution        string `json:"attribution"`
					SourceDistribution string `json:"sourceDistribution"`
					Modifications      string `json:"modifications"`
				} `json:"nodes"`
				TotalCount int `json:"totalCount"`
				PageInfo   struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"licenses"`
		} `json:"organization"`
	}

	if err := c.gql.Execute(ctx, graphql.LicensesQuery, vars, &result); err != nil {
		return nil, err
	}

	licenses := make([]OrganizationLicense, len(result.Organization.Licenses.Nodes))
	for i, n := range result.Organization.Licenses.Nodes {
		// Use shortId if available, otherwise use spdxId
		shortID := n.Content.ShortID
		if shortID == "" {
			shortID = n.Content.SpdxID
		}
		licenses[i] = OrganizationLicense{
			ShortID:            shortID,
			Name:               n.Content.Name,
			DerivedState:       n.State,
			CopyLeft:           n.CopyLeft,
			OsiApproved:        n.OsiApproved,
			FsfLibre:           n.FsfLibre,
			Deprecated:         n.Deprecated,
			Attribution:        n.Attribution,
			SourceDistribution: n.SourceDistribution,
			Modifications:      n.Modifications,
		}
	}

	return &LicensesResult{
		Licenses:    licenses,
		TotalCount:  result.Organization.Licenses.TotalCount,
		HasNextPage: result.Organization.Licenses.PageInfo.HasNextPage,
		EndCursor:   result.Organization.Licenses.PageInfo.EndCursor,
	}, nil
}
