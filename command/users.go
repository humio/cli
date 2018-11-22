package command

import "strings"

type simpleAccount struct {
	Username  string
	FullName  string
	IsRoot    bool
	CreatedAt string
}

func formatSimpleAccount(account simpleAccount) string {
	columns := []string{account.Username, account.FullName, yesNo(account.IsRoot), account.CreatedAt}
	return strings.Join(columns, " | ")
}
