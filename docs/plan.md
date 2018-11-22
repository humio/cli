humio users list
humio users create <username> [--name=<name>] [--country=<code>] [--email=<email>] [--root=true] ...
humio users update <username> ([-c] | [--create]) ... // Same arguments as create
humio users show <username>

humio repo list
humio repo rename <old> <new>
humio repo show <name>
humio repo rm <name>
humio repo add <name> [--retention-size=<gb>] [--retention-days=<num>] [--archiving-enabled=<bool>] [--archiving-bucket=<bucket>] [--archiving-region=<region>]
humio repo add --file=<path>
humio repo update <name> ([-c] | [--create]) // same arguments as `add`
humio repo push <path> ([-c] | [--create])
humio repo pull <repo> <output-path>

humio parsers list <repo>
humio parsers push <repo> (<file> | <url> | <github-path>) [--name=<name>] ([-u] | [--update])
humio parsers pull (<repo> <name> | <url> | <github-path>) <output-path> [--name=<name>]
humio parsers rm <repo> <name>

humio members add <view> <username> [--admin-users=<bool>] [--allow-deletion=<bool>]
humio members rm <view> <username>
humio members update <view> <username> ([-c] | [--create]) // same arguments as `add`

humio query-blacklist list [--repo=<name>]
humio query-blacklist add "<pattern>" [--repo=<name>] [--type=exact|regex]
humio query-blacklist rm <blacklist-id>
humio query-blacklist update <blacklist-id> ([-c] | [--create]) // same arguments as `add`

humio status [node]

humio view list
humio view rename <old> <new>
humio view show <name>
humio view rm <name>
humio view update --file=<path> ([-c] | [--create])

humio dashboards add <view> <name>
humio dashboards pull (<view> <name> | <url> | <github-path>) <output-path> [--name=<name>]
humio dashboards push <view> (<file> | <url> | <github-path>) [--name=<name>] ([-u] | [--update])
humio dashboards update <view> <name> [--enabled=<bool>]
humio dashboards rm <view> <name>

humio dashboards links list <view>
humio dashboards links add <view> <name>
humio dashboards links rm <view> <name>

humio alerts list <view>
humio alerts show <view> <name>
humio alerts push <view> (<file> | <url> | <github-path>) [--name=<name>] ([-u] | [--update])
humio alerts pull <view> (<url> | <github-path>) <output-path> [--name=<name>]
humio alerts update <view> <name> [--enabled=<bool>]
humio alerts rm <view>

humio notifiers list <view>
humio notifiers show <view> <name>
humio notifiers push <view> (<file> | <url> | <github-path>) [--name=<name>] ([-u] | [--update])
humio notifiers pull <view> (<url> | <github-path>) <output-path> [--name=<name>]
humio notifiers update <view> <name> [--enabled=<bool>]
humio notifiers rm <view> <webhook>
