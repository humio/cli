package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ryanuber/columnize"

	cli "gopkg.in/urfave/cli.v2"
)

type token struct {
	Name           string `json:"name"`
	Token          string `json:"token"`
	AssignedParser string `json:"parser"`
}

func TokenList(c *cli.Context) error {
	config, _ := getServerConfig(c)

	ensureToken(config)
	ensureRepo(config)
	ensureURL(config)

	url := config.URL + "api/v1/repositories/" + config.Repo + "/ingesttokens"

	resp, clientErr := getReq(url, config.Token)
	defer resp.Body.Close()

	if clientErr != nil {
		log.Fatal(clientErr)
	}

	if resp.StatusCode >= 400 {
		log.Fatal(resp)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	tokens := make([]token, 0)
	jsonErr := json.Unmarshal(body, &tokens)

	if jsonErr != nil {
		log.Fatalf("Could not parser JSON: %#v", string(body))
	}

	var output []string
	output = append(output, "Name | Token | Assigned Parser")
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		output = append(output, fmt.Sprintf("%v | %v | %v", token.Name, token.Token, valueOrEmpty(token.AssignedParser)))
	}

	table := columnize.SimpleFormat(output)

	fmt.Println()
	fmt.Println(table)
	fmt.Println()

	return nil
}

func valueOrEmpty(v string) string {
	if v == "" {
		return "-"
	}
	return v
}
