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

const securityIncidentUpdateFields = `
	title
	updateType
	body
	occurredAt
`

const securityIncidentFindingFields = `
	id
	status
	matchMethod
	matchedFields
	firstDetectedAt
	lastConfirmedAt
	isPartSbom
	component {
		id
		name
		version
		kind
		purl
		cpes
		licensesExp
		group
		primary
		internal
		sbomId
		updatedAt
		sbom {
			id
			projectVersion
			project {
				id
				name
				projectGroupId
			}
		}
	}
	rootSbom {
		id
		projectVersion
		project {
			id
			name
			projectGroupId
		}
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

	// SecurityIncidentFindingsQuery fetches customer-visible findings for the current organization.
	SecurityIncidentFindingsQuery = `
		query SecurityIncidentFindings($id: Uuid!, $statuses: [SecurityIncidentFindingStatusEnum!]) {
			securityIncident(id: $id) {
				id
				title
				slug
				status
				severity
				findings(statuses: $statuses) {
					` + securityIncidentFindingFields + `
				}
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

	// SecurityIncidentUpdateMutation updates editable incident fields.
	SecurityIncidentUpdateMutation = `
		mutation UpdateSecurityIncident($id: Uuid!, $title: String, $severity: String, $confidence: String, $summary: String, $recommendedActions: String, $sourceUrls: String) {
			updateSecurityIncident(input: {
				id: $id,
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

	// SecurityIncidentWithdrawMarkersMutation withdraws active markers from an incident.
	SecurityIncidentWithdrawMarkersMutation = `
		mutation WithdrawSecurityIncidentMarkers($securityIncidentId: Uuid!, $markerIds: [Uuid!]!) {
			withdrawSecurityIncidentMarkers(input: { securityIncidentId: $securityIncidentId, markerIds: $markerIds }) {
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

	// SecurityIncidentResolveMutation resolves an active incident.
	SecurityIncidentResolveMutation = `
		mutation ResolveSecurityIncident($id: Uuid!) {
			resolveSecurityIncident(input: { id: $id }) {
				securityIncident {
					` + securityIncidentFields + `
				}
				errors
			}
		}
	`

	// SecurityIncidentArchiveMutation archives an incident.
	SecurityIncidentArchiveMutation = `
		mutation ArchiveSecurityIncident($id: Uuid!) {
			archiveSecurityIncident(input: { id: $id }) {
				securityIncident {
					` + securityIncidentFields + `
				}
				errors
			}
		}
	`

	// SecurityIncidentCreateUpdateMutation adds a timeline update to an incident.
	SecurityIncidentCreateUpdateMutation = `
		mutation CreateSecurityIncidentUpdate($securityIncidentId: Uuid!, $title: String!, $updateType: String!, $occurredAt: ISO8601DateTime!, $body: String, $customerVisible: Boolean = false) {
			createSecurityIncidentUpdate(input: {
				securityIncidentId: $securityIncidentId,
				title: $title,
				updateType: $updateType,
				occurredAt: $occurredAt,
				body: $body,
				customerVisible: $customerVisible
			}) {
				securityIncidentUpdate {
					` + securityIncidentUpdateFields + `
				}
				errors
			}
		}
	`

	// SecurityIncidentSuppressFindingMutation suppresses a finding for the current organization.
	SecurityIncidentSuppressFindingMutation = `
		mutation SuppressSecurityIncidentFinding($findingId: Uuid!, $reason: String) {
			suppressSecurityIncidentFinding(input: { findingId: $findingId, reason: $reason }) {
				finding {
					` + securityIncidentFindingFields + `
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

	// SecurityIncidentEnabledOrganizationsQuery lists organizations with incident scans enabled.
	SecurityIncidentEnabledOrganizationsQuery = `
		query SecurityIncidentEnabledOrganizations {
			securityIncidentEnabledOrganizations {
				id
				name
				securityIncidentsEnabled
			}
		}
	`
)
