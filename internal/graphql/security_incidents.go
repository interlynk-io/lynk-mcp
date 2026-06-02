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

const securityIncidentFields = `
	id
	title
	slug
	summary
	severity
	status
	confidence
	recommendedActions
	sourceUrls
	firstSeenAt
	publishedAt
	lastUpdatedAt
	createdAt
	updatedAt
	markers {
		id
		markerType
		purl
		componentName
		componentVersion
		githubUrl
		active
		addedAt
		withdrawnAt
	}
	orgImpactState {
		status
		severity
		impactedProjectsCount
		impactedVersionsCount
		impactedComponentsCount
		lastEvaluatedAt
		lastNotifiedAt
	}
`

const (
	// SecurityIncidentsQuery lists security incidents visible to the current organization.
	SecurityIncidentsQuery = `
		query SecurityIncidents($status: [String!]) {
			securityIncidents(status: $status) {
				` + securityIncidentFields + `
			}
		}
	`

	// SecurityIncidentQuery fetches a single security incident by ID.
	SecurityIncidentQuery = `
		query SecurityIncident($id: Uuid!) {
			securityIncident(id: $id) {
				` + securityIncidentFields + `
			}
		}
	`

	// SecurityIncidentCreateMutation creates a draft security incident.
	SecurityIncidentCreateMutation = `
		mutation CreateSecurityIncident($title: String!, $severity: String!, $confidence: String, $summary: String, $recommendedActions: String, $sourceUrls: String) {
			createSecurityIncident(input: {
				title: $title,
				severity: $severity,
				confidence: $confidence,
				summary: $summary,
				recommendedActions: $recommendedActions,
				sourceUrls: $sourceUrls
			}) {
				securityIncident {
					` + securityIncidentFields + `
				}
				errors
			}
		}
	`

	// SecurityIncidentAddMarkersMutation adds markers to an incident.
	SecurityIncidentAddMarkersMutation = `
		mutation AddSecurityIncidentMarkers($securityIncidentId: Uuid!, $markers: [SecurityIncidentMarkerInput!]!) {
			addSecurityIncidentMarkers(input: { securityIncidentId: $securityIncidentId, markers: $markers }) {
				markers {
					id
					markerType
					purl
					componentName
					componentVersion
					githubUrl
					active
					addedAt
					withdrawnAt
				}
				errors
			}
		}
	`

	// SecurityIncidentPublishMutation publishes a draft incident and queues impact scanning.
	SecurityIncidentPublishMutation = `
		mutation PublishSecurityIncident($id: Uuid!) {
			publishSecurityIncident(input: { id: $id }) {
				securityIncident {
					` + securityIncidentFields + `
				}
				errors
			}
		}
	`

	// SecurityIncidentRerunImpactScanMutation queues impact scanning for an active or resolved incident.
	SecurityIncidentRerunImpactScanMutation = `
		mutation RerunSecurityIncidentImpactScan($id: Uuid!) {
			rerunSecurityIncidentImpactScan(input: { id: $id }) {
				securityIncident {
					` + securityIncidentFields + `
				}
				errors
			}
		}
	`

	// SecurityIncidentDryRunImpactScanMutation queues a dry-run impact scan.
	SecurityIncidentDryRunImpactScanMutation = `
		mutation DryRunSecurityIncidentImpactScan($id: Uuid!) {
			dryRunSecurityIncidentImpactScan(input: { id: $id }) {
				status
				errors
			}
		}
	`

	// SecurityIncidentDryRunResultQuery fetches dry-run scan status and org summaries.
	SecurityIncidentDryRunResultQuery = `
		query SecurityIncidentDryRunResult($incidentId: Uuid!) {
			securityIncidentDryRunResult(incidentId: $incidentId) {
				status
				error
				completedAt
				totalOrgsImpacted
				orgResults {
					organizationId
					organizationName
					impactedComponentsCount
					impactedProjectsCount
					impactedVersionsCount
				}
			}
		}
	`

	// SecurityIncidentDryRunOrgResultQuery fetches dry-run findings for a selected org.
	SecurityIncidentDryRunOrgResultQuery = `
		query SecurityIncidentDryRunOrgResult($incidentId: Uuid!, $orgId: Uuid!, $first: Int, $after: String) {
			securityIncidentDryRunResult(incidentId: $incidentId) {
				status
				error
				completedAt
				totalOrgsImpacted
				orgResults {
					organizationId
					organizationName
					impactedComponentsCount
					impactedProjectsCount
					impactedVersionsCount
				}
				org(orgId: $orgId) {
					organizationId
					organizationName
					impactedComponentsCount
					impactedProjectsCount
					impactedVersionsCount
					findings(first: $first, after: $after) {
						nodes {
							id
							organizationId
							projectId
							rootSbomId
							componentSbomId
							componentId
							componentName
							componentVersion
							componentPurl
							markerId
							markerType
							markerPurl
							markerComponentName
							markerComponentVersion
							matchMethod
							matchedFields
							rootProjectName
							rootProjectVersion
							isPartSbom
						}
						pageInfo {
							hasNextPage
							endCursor
						}
					}
				}
			}
		}
	`
)
