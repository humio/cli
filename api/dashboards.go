package api

import (

	"github.com/shurcooL/graphql"
)

// Dashboards struct
type Dashboards struct {
	client *Client
}


// Dashboard type
type Dashboard struct {
	Name string
	TemplateYaml    string 							`yaml:",flow"`
}

// DashboardListItem for dashbooards
type DashboardListItem struct {
	Name      string
}

//Dashboards client
func (c *Client) Dashboards() *Dashboards { return &Dashboards{client: c} }


// Add , should check beforehand, force is explicit for now
func (p *Dashboards) Add(repositoryName string, dashboard *Dashboard) error {

	var mutation struct {
		CreateDashboardFromTemplate struct {
			Type string `graphql:"__typename"`
		} `graphql:"createDashboardFromTemplate(input: { overrideName: $overrideName, searchDomainName: $searchDomainName, template: $template})"`
	}

	variables := map[string]interface{}{
		"searchDomainName":     graphql.String(repositoryName),
		"template":     graphql.String(dashboard.TemplateYaml),
		"overrideName":     graphql.String(dashboard.Name),
	}

	return p.client.Mutate(&mutation, variables)
}



type DashboardQueryData struct {
	Name     string
}
//List dashboards in a repo
func (p *Dashboards) List(reposistoryName string) ([]DashboardListItem, error) {
	var q struct {
		Repository struct {
			Dashboards []DashboardListItem
		} `graphql:"repository(name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"repositoryName": graphql.String(reposistoryName),
	}

	graphqlErr := p.client.Query(&q, variables)

	var dashboards []DashboardListItem
	if graphqlErr == nil {
		dashboards = q.Repository.Dashboards
	}

	return dashboards, graphqlErr
}


type DashboardExportItem struct {
	Name      			string
	Id      				string
	TemplateYaml    string
}


// Get a dashboard to export it
func (p *Dashboards) GetAll(reposistoryName string) ([]DashboardExportItem, error) {
	var query struct {
		Repository struct {
			Dashboards [] DashboardExportItem
		} `graphql:"repository(name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"repositoryName": graphql.String(reposistoryName),
	}


	graphqlErr := p.client.Query(&query, variables)

	var dashboards[]DashboardExportItem
	if graphqlErr == nil {
		// fmt.Printf("%+v", &query.Repository.Dashboards)
		// fmt.Printf("stuf")
		dashboards = query.Repository.Dashboards

	  //  return query.Repository.Dashboards, graphqlErr
		// dashboard = Dashboard {
		// 	TemplateYaml: query.DashboardQueryData.TemplateYaml,
		// }
	}

	return dashboards, graphqlErr
}