package command

import (
	"io/ioutil"

	cli "gopkg.in/urfave/cli.v2"
)

func TokenAdd(c *cli.Context) error {
	config, _ := getServerConfig(c)

	ensureToken(config)
	ensureRepo(config)
	ensureURL(config)

	name := c.String("name")
	if name == "" {
		exit("Missing name argument")
	}
	parser := c.String("parser")

	body := ""
	if parser == "" {
		body = `{"name": "` + name + `"}`
	} else {
		body = `{"name": "` + name + `", "parser": "` + parser + `"}`
	}

	url := config.URL + "/api/v1/repositories/" + config.Repo + "/ingesttokens"

	resp, clientErr := postJSON(url, body, config.Token)

	if clientErr != nil {
		panic(clientErr)
	}
	if resp.StatusCode >= 400 {
		_, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			panic(readErr)
		}
		//		fmt.Println(resp.StatusCode)
		//		fmt.Println(string(responseData))
	}
	resp.Body.Close()

	return nil
}
