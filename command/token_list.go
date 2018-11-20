package command

import (
	"fmt"
	"io/ioutil"

	cli "gopkg.in/urfave/cli.v2"
)

func TokenList(c *cli.Context) error {
	config, _ := getServerConfig(c)

	ensureToken(config)
	ensureRepo(config)
	ensureURL(config)

	url := config.URL + "/api/v1/repositories/" + config.Repo + "/ingesttokens"

	resp, clientErr := getReq(url, config.Token)
	fmt.Println(resp.Body)
	if clientErr != nil {
		panic(clientErr)
	}
	if resp.StatusCode >= 400 {
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			panic(readErr)
		}
		fmt.Println(body)
	}
	resp.Body.Close()

	return nil
}
