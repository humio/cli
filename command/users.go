package command

import (
	"strings"

	"github.com/humio/cli/api"
)

func formatSimpleAccount(account api.User) string {
	columns := []string{account.Username, account.FullName, yesNo(account.IsRoot), account.CreatedAt}
	return strings.Join(columns, " | ")
}
