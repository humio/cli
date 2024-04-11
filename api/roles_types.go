package api

import graphql "github.com/cli/shurcooL-graphql"

// ViewPermission is a type that represents the view permissions that can be assigned to a role.
type ViewPermission string

const (
	ViewPermissionChangeUserAccess                  ViewPermission = "ChangeUserAccess"
	ViewPermissionChangeTriggersAndActions          ViewPermission = "ChangeTriggersAndActions"
	ViewPermissionChangeDashboards                  ViewPermission = "ChangeDashboards"
	ViewPermissionChangeDashboardReadonlyToken      ViewPermission = "ChangeDashboardReadonlyToken"
	ViewPermissionChangeFiles                       ViewPermission = "ChangeFiles"
	ViewPermissionChangeInteractions                ViewPermission = "ChangeInteractions"
	ViewPermissionChangeParsers                     ViewPermission = "ChangeParsers"
	ViewPermissionChangeSavedQueries                ViewPermission = "ChangeSavedQueries"
	ViewPermissionConnectView                       ViewPermission = "ConnectView"
	ViewPermissionChangeDataDeletionPermissions     ViewPermission = "ChangeDataDeletionPermissions"
	ViewPermissionChangeRetention                   ViewPermission = "ChangeRetention"
	ViewPermissionChangeDefaultSearchSettings       ViewPermission = "ChangeDefaultSearchSettings"
	ViewPermissionChangeS3ArchivingSettings         ViewPermission = "ChangeS3ArchivingSettings"
	ViewPermissionDeleteDataSources                 ViewPermission = "DeleteDataSources"
	ViewPermissionDeleteRepositoryOrView            ViewPermission = "DeleteRepositoryOrView"
	ViewPermissionDeleteEvents                      ViewPermission = "DeleteEvents"
	ViewPermissionReadAccess                        ViewPermission = "ReadAccess"
	ViewPermissionChangeIngestTokens                ViewPermission = "ChangeIngestTokens"
	ViewPermissionChangePackages                    ViewPermission = "ChangePackages"
	ViewPermissionChangeViewOrRepositoryDescription ViewPermission = "ChangeViewOrRepositoryDescription"
	ViewPermissionChangeConnections                 ViewPermission = "ChangeConnections"
	ViewPermissionEventForwarding                   ViewPermission = "EventForwarding"
	ViewPermissionQueryDashboard                    ViewPermission = "QueryDashboard"
	ViewPermissionChangeViewOrRepositoryPermissions ViewPermission = "ChangeViewOrRepositoryPermissions"
	ViewPermissionChangeFdrFeeds                    ViewPermission = "ChangeFdrFeeds"
	ViewPermissionOrganizationOwnedQueries          ViewPermission = "OrganizationOwnedQueries"
	ViewPermissionReadExternalFunctions             ViewPermission = "ReadExternalFunctions"
	ViewPermissionChangeIngestFeeds                 ViewPermission = "ChangeIngestFeeds"
	ViewPermissionChangeScheduledReports            ViewPermission = "ChangeScheduledReports"
)

// Get returns the ViewPermission value for the given string, if it exists.
func (vp ViewPermission) Get(p string) (ViewPermission, bool) {
	switch ViewPermission(p) {
	case ViewPermissionChangeUserAccess,
		ViewPermissionChangeTriggersAndActions,
		ViewPermissionChangeDashboards,
		ViewPermissionChangeDashboardReadonlyToken,
		ViewPermissionChangeFiles,
		ViewPermissionChangeInteractions,
		ViewPermissionChangeParsers,
		ViewPermissionChangeSavedQueries,
		ViewPermissionConnectView,
		ViewPermissionChangeDataDeletionPermissions,
		ViewPermissionChangeRetention,
		ViewPermissionChangeDefaultSearchSettings,
		ViewPermissionChangeS3ArchivingSettings,
		ViewPermissionDeleteDataSources,
		ViewPermissionDeleteRepositoryOrView,
		ViewPermissionDeleteEvents,
		ViewPermissionReadAccess,
		ViewPermissionChangeIngestTokens,
		ViewPermissionChangePackages,
		ViewPermissionChangeViewOrRepositoryDescription,
		ViewPermissionChangeConnections,
		ViewPermissionEventForwarding,
		ViewPermissionQueryDashboard,
		ViewPermissionChangeViewOrRepositoryPermissions,
		ViewPermissionChangeFdrFeeds,
		ViewPermissionOrganizationOwnedQueries,
		ViewPermissionReadExternalFunctions,
		ViewPermissionChangeIngestFeeds,
		ViewPermissionChangeScheduledReports:
		return ViewPermission(p), true
	default:
		return "", false
	}
}

// SystemPermission is a type that represents the system permissions that can be assigned to a role.
type SystemPermission string

const (
	SystemPermissionReadHealthCheck                   SystemPermission = "ReadHealthCheck"
	SystemPermissionManageOrganizations               SystemPermission = "ManageOrganizations"
	SystemPermissionImportOrganization                SystemPermission = "ImportOrganization"
	SystemPermissionDeleteOrganizations               SystemPermission = "DeleteOrganizations"
	SystemPermissionChangeSystemPermissions           SystemPermission = "ChangeSystemPermissions"
	SystemPermissionManageCluster                     SystemPermission = "ManageCluster"
	SystemPermissionIngestAcrossAllReposWithinCluster SystemPermission = "IngestAcrossAllReposWithinCluster"
	SystemPermissionDeleteHumioOwnedRepositoryOrView  SystemPermission = "DeleteHumioOwnedRepositoryOrView"
	SystemPermissionChangeUsername                    SystemPermission = "ChangeUsername"
	SystemPermissionChangeFeatureFlags                SystemPermission = "ChangeFeatureFlags"
	SystemPermissionChangeSubdomains                  SystemPermission = "ChangeSubdomains"
	SystemPermissionListSubdomains                    SystemPermission = "ListSubdomains"
	SystemPermissionPatchGlobal                       SystemPermission = "PatchGlobal"
	SystemPermissionChangeBucketStorage               SystemPermission = "ChangeBucketStorage"
	SystemPermissionManageOrganizationLinks           SystemPermission = "ManageOrganizationLinks"
)

// Get returns the SystemPermission value for the given string, if it exists.
func (sp SystemPermission) Get(p string) (SystemPermission, bool) {
	switch SystemPermission(p) {
	case SystemPermissionReadHealthCheck,
		SystemPermissionManageOrganizations,
		SystemPermissionImportOrganization,
		SystemPermissionDeleteOrganizations,
		SystemPermissionChangeSystemPermissions,
		SystemPermissionManageCluster,
		SystemPermissionIngestAcrossAllReposWithinCluster,
		SystemPermissionDeleteHumioOwnedRepositoryOrView,
		SystemPermissionChangeUsername,
		SystemPermissionChangeFeatureFlags,
		SystemPermissionChangeSubdomains,
		SystemPermissionListSubdomains,
		SystemPermissionPatchGlobal,
		SystemPermissionChangeBucketStorage,
		SystemPermissionManageOrganizationLinks:
		return SystemPermission(p), true
	default:
		return "", false
	}
}

// OrganizationPermission is a type that represents the organization permissions that can be assigned to a role.
type OrganizationPermission string

const (
	OrganizationPermissionExportOrganization                     OrganizationPermission = "ExportOrganization"
	OrganizationPermissionChangeOrganizationPermissions          OrganizationPermission = "ChangeOrganizationPermissions"
	OrganizationPermissionChangeIdentityProviders                OrganizationPermission = "ChangeIdentityProviders"
	OrganizationPermissionCreateRepository                       OrganizationPermission = "CreateRepository"
	OrganizationPermissionManageUsers                            OrganizationPermission = "ManageUsers"
	OrganizationPermissionViewUsage                              OrganizationPermission = "ViewUsage"
	OrganizationPermissionChangeOrganizationSettings             OrganizationPermission = "ChangeOrganizationSettings"
	OrganizationPermissionChangeIPFilters                        OrganizationPermission = "ChangeIPFilters"
	OrganizationPermissionChangeSessions                         OrganizationPermission = "ChangeSessions"
	OrganizationPermissionChangeAllViewOrRepositoryPermissions   OrganizationPermission = "ChangeAllViewOrRepositoryPermissions"
	OrganizationPermissionIngestAcrossAllReposWithinOrganization OrganizationPermission = "IngestAcrossAllReposWithinOrganization"
	OrganizationPermissionDeleteAllRepositories                  OrganizationPermission = "DeleteAllRepositories"
	OrganizationPermissionDeleteAllViews                         OrganizationPermission = "DeleteAllViews"
	OrganizationPermissionViewAllInternalNotifications           OrganizationPermission = "ViewAllInternalNotifications"
	OrganizationPermissionChangeFleetManagement                  OrganizationPermission = "ChangeFleetManagement"
	OrganizationPermissionViewFleetManagement                    OrganizationPermission = "ViewFleetManagement"
	OrganizationPermissionChangeTriggersToRunAsOtherUsers        OrganizationPermission = "ChangeTriggersToRunAsOtherUsers"
	OrganizationPermissionMonitorQueries                         OrganizationPermission = "MonitorQueries"
	OrganizationPermissionBlockQueries                           OrganizationPermission = "BlockQueries"
	OrganizationPermissionChangeSecurityPolicies                 OrganizationPermission = "ChangeSecurityPolicies"
	OrganizationPermissionChangeExternalFunctions                OrganizationPermission = "ChangeExternalFunctions"
)

// Get returns the OrganizationPermission value for the given string, if it exists.
func (op OrganizationPermission) Get(p string) (OrganizationPermission, bool) {
	switch OrganizationPermission(p) {
	case OrganizationPermissionExportOrganization,
		OrganizationPermissionChangeOrganizationPermissions,
		OrganizationPermissionChangeIdentityProviders,
		OrganizationPermissionCreateRepository,
		OrganizationPermissionManageUsers,
		OrganizationPermissionViewUsage,
		OrganizationPermissionChangeOrganizationSettings,
		OrganizationPermissionChangeIPFilters,
		OrganizationPermissionChangeSessions,
		OrganizationPermissionChangeAllViewOrRepositoryPermissions,
		OrganizationPermissionIngestAcrossAllReposWithinOrganization,
		OrganizationPermissionDeleteAllRepositories,
		OrganizationPermissionDeleteAllViews,
		OrganizationPermissionViewAllInternalNotifications,
		OrganizationPermissionChangeFleetManagement,
		OrganizationPermissionViewFleetManagement,
		OrganizationPermissionChangeTriggersToRunAsOtherUsers,
		OrganizationPermissionMonitorQueries,
		OrganizationPermissionBlockQueries,
		OrganizationPermissionChangeSecurityPolicies,
		OrganizationPermissionChangeExternalFunctions:
		return OrganizationPermission(p), true
	default:
		return "", false
	}
}

type ObjectAction string

const (
	ObjectActionUnknown             ObjectAction = "Unknown"
	ObjectActionReadOnlyAndHidden   ObjectAction = "ReadOnlyAndHidden"
	ObjectActionReadWRiteAndVisible ObjectAction = "ReadWRiteAndVisible"
)

// AddRoleInput is the struct of input variables passed to the `createRole()` API mutation.
type AddRoleInput struct {
	DisplayName       graphql.String            `json:"displayName"`
	Color             *graphql.String           `json:"color,omitempty"` // The color of the role in RGB hexadecimal, e.g. #FF0000.
	ViewPermissions   []ViewPermission          `json:"viewPermissions"`
	SystemPermissions *[]SystemPermission       `json:"systemPermissions,omitempty"`
	OrgPermissions    *[]OrganizationPermission `json:"organizationPermissions,omitempty"`
	ObjectAction      *ObjectAction             `json:"objectAction,omitempty"` // Undocumented field.

	// Oddly, the GraphQL API does not expose this field for the createRole mutation.
	// Description *graphql.String `json:"description,omitempty"`
}

// NewAddRoleInput returns the AddRoleInput struct initialized with the given values.
func NewAddRoleInput(name string, color *string, viewPermissions, systemPermissions, orgPermissions []string) AddRoleInput {
	ari := AddRoleInput{
		DisplayName: graphql.String(name),

		// ViewPermissions is a required field, so we initialize it with an empty slice.
		ViewPermissions: make([]ViewPermission, 0, len(viewPermissions)),
	}

	if color != nil {
		ari.Color = graphql.NewString(graphql.String(*color))
	}

	for _, permission := range viewPermissions {
		if vp, ok := ViewPermission(permission).Get(permission); ok {
			ari.ViewPermissions = append(ari.ViewPermissions, vp)
		}
	}

	if len(systemPermissions) > 0 {
		if ari.SystemPermissions == nil {
			ari.SystemPermissions = &[]SystemPermission{}
		}
		for _, permission := range systemPermissions {
			if sp, ok := SystemPermission(permission).Get(permission); ok {
				*ari.SystemPermissions = append(*ari.SystemPermissions, sp)
			}
		}
	}

	if len(orgPermissions) > 0 {
		if ari.OrgPermissions == nil {
			ari.OrgPermissions = &[]OrganizationPermission{}
		}
		for _, permission := range orgPermissions {
			if op, ok := OrganizationPermission(permission).Get(permission); ok {
				*ari.OrgPermissions = append(*ari.OrgPermissions, op)
			}
		}
	}
	return ari
}

// UpdateRoleInput is the struct of input variables passed to the `updateRole()` API mutation.
type UpdateRoleInput struct {
	ID                      graphql.String            `json:"roleId"`
	DisplayName             graphql.String            `json:"displayName"`
	ViewPermissions         []ViewPermission          `json:"viewPermissions"`
	Description             *graphql.String           `json:"description,omitempty"`
	Color                   *graphql.String           `json:"color,omitempty"`
	SystemPermissions       *[]SystemPermission       `json:"systemPermissions,omitempty"`
	OrganizationPermissions *[]OrganizationPermission `json:"organizationPermissions,omitempty"`
	ObjectAction            *ObjectAction             `json:"objectAction,omitempty"`
}

// NewUpdateRoleInput returns the UpdateRoleInput struct initialized with the given values.
func NewUpdateRoleInput(id, name string, viewPermissions, systemPermissions, orgPermissions []string, color *string) UpdateRoleInput {
	uri := UpdateRoleInput{
		ID:          graphql.String(id),
		DisplayName: graphql.String(name),
	}

	if color != nil {
		uri.Color = graphql.NewString(graphql.String(*color))
	}

	viewPerms := make([]ViewPermission, 0, len(viewPermissions))
	for _, permission := range viewPermissions {
		if vp, ok := ViewPermission(permission).Get(permission); ok {
			viewPerms = append(viewPerms, vp)
		}
	}
	uri.ViewPermissions = viewPerms

	sysPerms := make([]SystemPermission, 0, len(systemPermissions))
	for _, permission := range systemPermissions {
		if sp, ok := SystemPermission(permission).Get(permission); ok {
			sysPerms = append(sysPerms, sp)
		}
	}
	uri.SystemPermissions = &sysPerms

	orgPerms := make([]OrganizationPermission, 0, len(orgPermissions))
	for _, permission := range orgPermissions {
		if op, ok := OrganizationPermission(permission).Get(permission); ok {
			orgPerms = append(orgPerms, op)
		}
	}
	uri.OrganizationPermissions = &orgPerms
	return uri
}
