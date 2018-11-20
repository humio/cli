package command

import (
	"fmt"
	"io/ioutil"
	"strconv"

	cli "gopkg.in/urfave/cli.v2"
)

func ParserPush(c *cli.Context) error {
	config, _ := getServerConfig(c)

	ensureToken(config)
	ensureRepo(config)
	ensureURL(config)

	name := c.String("name")
	if name == "" {
		panic("Missing name argument")
	}

	query := ""

	fileNameSlices := c.StringSlice("query-file")
	if len(fileNameSlices) != 1 {
		querySlices := c.StringSlice("query")
		if len(querySlices) != 1 {
			panic("Missing query argument")
		} else {
			query = strconv.Quote(querySlices[0])
		}
	} else {
		file, readErr := ioutil.ReadFile(fileNameSlices[0])
		if readErr != nil {
			exit("Could not read file: " + fileNameSlices[0])
		}
		query = strconv.Quote(string(file))
	}

	body := `{"parser": ` + query + `, "kind": "humio", "parseKeyValues": false, "dateTimeFields": ["@timestamp"]}`
	url := config.URL + "/api/v1/repositories/" + config.Repo + "/parsers/" + name
	resp, clientErr := postJSON(url, body, config.Token)

	if clientErr != nil {
		panic(clientErr)
	}
	if resp.StatusCode == 409 && c.Bool("force") {
		resp, _ = putJSON(url, body, config.Token)
	}

	if resp.StatusCode >= 400 {
		responseData, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			panic(readErr)
		}
		fmt.Println("Error: status code =", resp.StatusCode)
		fmt.Println(string(responseData))
	}
	//fmt.Println(resp)
	resp.Body.Close()

	return nil
}
