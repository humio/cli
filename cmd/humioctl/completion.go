/*
Copyright 2018 Humio Ltd.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

This file is derived of source code from the Helm Project.
You can find the original source here: https://github.com/helm/helm

The project is released under the same Apache 2.0 license, the original
work is copyright of the Helm Authorer.
*/

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

const completionDesc = `
Generate autocompletions script for Humio for the specified shell (bash or zsh).

This command can generate shell autocompletions. e.g.

	$ humioctl completion bash

Can be sourced as such

	$ source <(humio completion bash)
`

var (
	completionShells = map[string]func(out io.Writer, cmd *cobra.Command) error{
		"bash": runCompletionBash,
		"zsh":  runCompletionZsh,
	}
)

func newCompletionCmd() *cobra.Command {
	var shells []string
	for s := range completionShells {
		shells = append(shells, s)
	}

	cmd := &cobra.Command{
		Use:   "completion SHELL",
		Short: "Generate autocompletions script for the specified shell (bash or zsh)",
		Long:  completionDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompletion(os.Stdout, cmd, args)
		},
		ValidArgs: shells,
	}

	return cmd
}

func runCompletion(out io.Writer, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("shell not specified")
	}
	if len(args) > 1 {
		return fmt.Errorf("too many arguments, expected only the shell type")
	}
	run, found := completionShells[args[0]]
	if !found {
		return fmt.Errorf("unsupported shell type %q", args[0])
	}

	return run(out, cmd)
}

func runCompletionBash(out io.Writer, cmd *cobra.Command) error {
	return cmd.Root().GenBashCompletion(out)
}

func runCompletionZsh(out io.Writer, cmd *cobra.Command) error {
	zshInitialization := `#compdef humio

__humio_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand
	source "$@"
}
__humio_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift
		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__humio_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}
__humio_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?
	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}
__humio_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}
__humio_declare() {
	if [ "$1" == "-F" ]; then
		whence -w "$@"
	else
		builtin declare "$@"
	fi
}
__humio_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}
__humio_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}
__humio_filedir() {
	local RET OLD_IFS w qw
	__debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi
	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"
	IFS="," __debug "RET=${RET[@]} len=${#RET[@]}"
	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__humio_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}
__humio_quote() {
	if [[ $1 == \'* || $1 == \"* ]]; then
		# Leave out first character
		printf %q "${1:1}"
	else
		printf %q "$1"
	fi
}
autoload -U +X bashcompinit && bashcompinit
# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi
__humio_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__humio_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__humio_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__humio_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__humio_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__humio_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/__humio_declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__humio_type/g" \
	<<'BASH_COMPLETION_EOF'
`
	_, _ = out.Write([]byte(zshInitialization))

	buf := new(bytes.Buffer)
	_ = cmd.Root().GenBashCompletion(buf)
	_, _ = out.Write(buf.Bytes())

	zshTail := `
BASH_COMPLETION_EOF
}
__humio_bash_source <(__humio_convert_bash_to_zsh)
`
	_, _ = out.Write([]byte(zshTail))
	return nil
}
