package integration

var repoShow = `
Name : %[1]s
Space Used : 0 B
Retention \(Size\) : 0 B
Retention \(Days\) : 0
`

var repoCreate = "Sucessfully created repo %[1]s\n" + repoShow

var reposList = `
NAME               | SPACE USED
-----------------------------------+-------------
humio                            | [a-zA-Z0-9_ ]*$+
humio-audit                      | 0 B
humio-metrics                    | 0 B
sandbox_[a-zA-Z0-9_ ]*$+         | 0 B
test-repo                        | 0 B
`
