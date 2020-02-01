package api

import (

	"github.com/shurcooL/graphql"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"bytes"
)

// Queries client
type Queries struct {
	client *Client
}



// QueryListItem for dashbooards
type QueryListItem struct {
	Name      string
}

// QueryJSON is the json object structure from rest, differs I think
type QueryJSON struct { 
	ID 													string `json:"id,omitempty"`
	Name 												string `json:"name,omitempty"`
	Options struct {
		Visualisation struct {
			Clr 				string 		`json:"clr,omitempty"`
			Cly 				string    `json:"cly,omitempty"`
			Col 				string		`json:"col,omitempty"`
			Cur 				string		`json:"cur,omitempty"`
			Dp 					string		`json:"dp,omitempty"`
			WidgetType 	string		`json:"widgetType,omitempty"`
		} `json:"visualisation"`
	} `json:"options"`
	Query struct {
		QueryString 									string	`json:"queryString,omitempty"`
		IsInteractive 								bool		`json:"isInteractive,omitempty"`
		IncludeDeletedEvents 					bool 		`json:"includeDeletedEvents,omitempty"`
		End 													string	`json:"end,omitempty"`
		ShowQueryEventDistribution 		bool 		`json:"showQueryEventDistribution,omitempty"`
		IsLive 												bool 		`json:"isLive,omitempty"`
		Start 												string	`json:"start,omitempty"`
		NoResultUntilDone 						bool 		`json:"noResultUntilDone,omitempty"`
		TimeZoneOffsetMinutes  				int			`json:"timeZoneOffsetMinutes,omitempty"`

	} `json:"query"`
}


//Queries client
func (c *Client) Queries() *Queries { return &Queries{client: c} }


// Add query
func (p *Queries) Add(repositoryName string, query *QueryJSON) error {
	url := "api/v1/repositories/"+repositoryName+"/savedqueries"
	var jsonOut []byte
	var jerr error
	jsonOut, jerr = json.Marshal(query)
	if jerr != nil {
		fmt.Println((fmt.Errorf("error while sending data: %v", jerr)))

	}
	_, err := p.client.HttpPOST(url, bytes.NewBuffer(jsonOut))
	if err != nil {
		fmt.Println((fmt.Errorf("error while sending data: %v", err)))
		return err
	}
	return err
}


//List queries in a repo
func (p *Queries) List(reposistoryName string) ([]QueryListItem, error) {
	var q struct {
		Repository struct {
			Queries []QueryListItem
		} `graphql:"repository(name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"repositoryName": graphql.String(reposistoryName),
	}

	graphqlErr := p.client.Query(&q, variables)

	var queries[]QueryListItem
	if graphqlErr == nil {
		queries = q.Repository.Queries
	}

	return queries, graphqlErr
}


// using REST version for now
// GetAll the queries to export
func (p *Queries) GetAll(repositoryName string) ([]QueryJSON, error) {
	url := "api/v1/repositories/"+repositoryName+"/savedqueries"
	resp, err := p.client.httpGET(url)
	if resp.StatusCode >= 400 {
		fmt.Println(fmt.Errorf("Error getting query from %s: %s", repositoryName,err ))
	}
  
	jsonData, errorthing := ioutil.ReadAll(resp.Body)
	var queries []QueryJSON
	jerr := json.Unmarshal(jsonData, &queries)
	// fmt.Println(fmt.Sprintf("%v", &queries))
	// fmt.Println(jerr)
	if jerr != nil {
		fmt.Println(fmt.Errorf("Error getting query from %s: %s", repositoryName,jerr ))

	}
	return queries, errorthing

}

// Get all the queries to export
// graphql version works, but is different than the rest API
// func (p *Queries) GetAll(reposistoryName string) ([]Query, error) {
// 	var query struct {
// 		Repository struct {
// 			Queries [] Query
// 		} `graphql:"repository(name: $repositoryName)"`
// 	}

// 	variables := map[string]interface{}{
// 		"repositoryName": graphql.String(reposistoryName),
// 	}

// 	graphqlErr := p.client.Query(&query, variables)

// 	var queries[]Query
// 	if graphqlErr == nil {
// 		queries = query.Repository.Queries

// 	}

// 	return queries, graphqlErr
// }