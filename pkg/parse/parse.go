package parse

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/carapace-sh/carapace-spec/pkg/command"
	"github.com/neurosnap/sentences/english"
)

type bazelFlag struct {
	negated     bool
	longhand    string
	shorthand   string
	definition  string
	description string
}

func (f bazelFlag) Usage() string {
	tokenizer, err := english.NewSentenceTokenizer(nil)
	if err != nil {
		return "" // TODO handle error
	}

	if tokens := tokenizer.Tokenize(f.description); len(tokens) > 0 {
		usage := tokens[0].Text
		usage = strings.TrimSuffix(usage, ".")
		usage = strings.TrimSpace(usage)
		return usage
	}
	return ""
}

func (f bazelFlag) Bool() bool {
	return strings.Contains(f.definition, "a boolean")
}

func (f bazelFlag) Repeatable() bool {
	return strings.Contains(f.definition, "may be used multiple times")
}

func (f bazelFlag) ToFlags() (flags []command.Flag) {
	flags = append(flags, command.Flag{
		Longhand:   "--" + f.longhand,
		Usage:      f.Usage(),
		Repeatable: f.Repeatable(),
		Value:      !f.Bool(),
	})
	if f.shorthand != "" {
		flags[0].Shorthand = "-" + f.shorthand
	}

	if f.negated {
		nf := flags[0]
		nf.Longhand = "--no-" + f.longhand
		nf.Shorthand = ""
		flags = append(flags, nf)
	}

	return flags
}

// Bazel parses `bazel help --long <command>` output.
func Bazel(name, description string, reader io.Reader) command.Command {
	r := regexp.MustCompile(`^  --(?P<negated>\[no\])?(?P<longhand>[^ ]+)( \[-(?P<shorthand>.)\])? \((?P<definition>.*)\)$`)
	command := command.Command{
		Name:        name,
		Description: description,
	}

	var f *bazelFlag
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if f != nil && strings.HasPrefix(scanner.Text(), "    ") {
			if f.description != "" {
				f.description += " "
			}
			f.description += strings.TrimSpace(scanner.Text())
			continue
		}

		if matches := r.FindStringSubmatch(scanner.Text()); matches != nil {
			if f != nil {
				for _, flag := range f.ToFlags() {
					command.AddFlag(flag)
				}
			}

			f = &bazelFlag{
				negated:    matches[1] != "",
				longhand:   matches[2],
				shorthand:  matches[4],
				definition: matches[5],
			}
		}
	}
	if f != nil {
		for _, flag := range f.ToFlags() {
			command.AddFlag(flag)
		}
	}
	return command
}
