package api

type Roles struct {
	client *Client
}

type Role struct {
	DisplayName             string
	ViewPermissions         []string
	SystemPermissions       []string
	OrganizationPermissions []string
}

type SearchDomainRole struct {
	DomainName *string
	Role       *Role
}

//TODO: This whole block is reserved for future development

// func (r *Roles) Create(displayName string, changeset GroupChangeSet) error {
// 	// var mutation struct {
// 	// 	RemoveUsersFromGroup struct {
// 	// 		Group struct {
// 	// 			ID string
// 	// 		}
// 	// 	} `graphql:"removeUsersFromGroup(input:{users:[$userID], groupId: $groupID})"`
// 	// }

// 	// variables := map[string]interface{}{
// 	// 	"userID":  graphql.String(userID),
// 	// 	"groupID": graphql.String(groupID),
// 	// }

// 	// return g.client.Mutate(&mutation, variables)
// }

// func (g *Roles) Delete(displayName string) error {
// 	// var mutation struct {
// 	// 	RemoveUsersFromGroup struct {
// 	// 		Group struct {
// 	// 			ID string
// 	// 		}
// 	// 	} `graphql:"removeUsersFromGroup(input:{users:[$userID], groupId: $groupID})"`
// 	// }

// 	// variables := map[string]interface{}{
// 	// 	"userID":  graphql.String(userID),
// 	// 	"groupID": graphql.String(groupID),
// 	// }

// 	// return g.client.Mutate(&mutation, variables)
// }

// func (g *Roles) Update(displayName string, changeset GroupChangeSet) error {
// 	// var mutation struct {
// 	// 	RemoveUsersFromGroup struct {
// 	// 		Group struct {
// 	// 			ID string
// 	// 		}
// 	// 	} `graphql:"removeUsersFromGroup(input:{users:[$userID], groupId: $groupID})"`
// 	// }

// 	// variables := map[string]interface{}{
// 	// 	"userID":  graphql.String(userID),
// 	// 	"groupID": graphql.String(groupID),
// 	// }

// 	// return g.client.Mutate(&mutation, variables)
// }
