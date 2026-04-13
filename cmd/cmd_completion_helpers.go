package cmd

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/mclucy/lucy/types"
	"github.com/urfave/cli/v3"
)

// CompletionCandidate holds a value and optional description for shell completion.
type CompletionCandidate struct {
	Value       string
	Description string
}

type CompletionRequest struct {
	Current            string
	Previous           string
	CompletingFlagName bool
	FlagValueName      string
	FlagValuePrefix    string
}

const (
	shellCompletionTrigger       = "--generate-shell-completion"
	completionDoubleDashSentinel = "__lucy_complete_double_dash__"
)

// NormalizeCompletionArgs preserves the user's intent when the current token is
// a literal "--".
//
// urfave/cli treats "--" as the end-of-options marker, so a raw invocation like
// `lucy search -- --generate-shell-completion` never enters shell-completion
// mode. We replace only that exact shape with an internal sentinel so the
// completion request survives parsing, then map it back to "--" below.
func NormalizeCompletionArgs(args []string) []string {
	normalized := slices.Clone(args)
	if len(normalized) >= 2 && normalized[len(normalized)-1] == shellCompletionTrigger && normalized[len(normalized)-2] == "--" {
		normalized[len(normalized)-2] = completionDoubleDashSentinel
	}
	return normalized
}

// PrintCandidates outputs plain completion values, one per line.
//
// urfave/cli's generated shell scripts do not expose a reliable shell identifier
// to custom ShellComplete handlers. Descriptions are therefore omitted here to
// keep completion output portable across bash, zsh, fish, and pwsh.
func PrintCandidates(candidates []CompletionCandidate) {
	for _, c := range candidates {
		fmt.Println(c.Value)
	}
}

// FilterByPrefix returns candidates whose Value starts with prefix (case-insensitive).
func FilterByPrefix(candidates []CompletionCandidate, prefix string) []CompletionCandidate {
	if prefix == "" {
		return candidates
	}
	lower := strings.ToLower(prefix)
	var out []CompletionCandidate
	for _, c := range candidates {
		if strings.HasPrefix(strings.ToLower(c.Value), lower) {
			out = append(out, c)
		}
	}
	return out
}

// StaticPlatformCandidates returns completion candidates for all user-facing platforms.
func StaticPlatformCandidates() []CompletionCandidate {
	return []CompletionCandidate{
		{Value: types.PlatformMinecraft.String(), Description: ""},
		{Value: types.PlatformFabric.String(), Description: ""},
		{Value: types.PlatformForge.String(), Description: ""},
		{Value: types.PlatformNeoforge.String(), Description: ""},
		{Value: types.PlatformMCDR.String(), Description: ""},
	}
}

// StaticSourceCandidates returns completion candidates for concrete upstream sources.
func StaticSourceCandidates() []CompletionCandidate {
	return []CompletionCandidate{
		{Value: "curseforge", Description: ""},
		{Value: types.SourceModrinth.String(), Description: ""},
		{Value: types.SourceGitHub.String(), Description: ""},
		{Value: types.SourceMCDR.String(), Description: ""},
	}
}

// StaticSortCandidates returns completion candidates for search sort options.
func StaticSortCandidates() []CompletionCandidate {
	return []CompletionCandidate{
		{Value: string(types.SearchSortRelevance), Description: "Sort by relevance"},
		{Value: string(types.SearchSortDownloads), Description: "Sort by download count"},
		{Value: string(types.SearchSortNewest), Description: "Sort by newest"},
	}
}

// ParseCompletionToken parses a partial "platform/name@version" token for shell completion.
// Returns parsed components and the active segment ("platform", "name", or "version").
//
// Uses manual string splitting instead of syntax.Parse which panics on partial input.
func ParseCompletionToken(token string) (platform, name, version, segment string) {
	if before, after, ok := strings.Cut(token, "@"); ok {
		version = after
		beforeAt := before
		if slashIdx := strings.Index(beforeAt, "/"); slashIdx >= 0 {
			platform = beforeAt[:slashIdx]
			name = beforeAt[slashIdx+1:]
		} else {
			name = beforeAt
		}
		segment = "version"
		return
	}

	if before, after, ok := strings.Cut(token, "/"); ok {
		platform = before
		name = after
		segment = "name"
		return
	}

	platform = token
	segment = "platform"
	return
}

func ParseCompletionRequest(cmd *cli.Command) CompletionRequest {
	args := completionInvocationArgs(cmd)
	for len(args) > 0 && args[len(args)-1] == shellCompletionTrigger {
		args = args[:len(args)-1]
	}

	request := CompletionRequest{}
	if len(args) > 0 {
		request.Current = args[len(args)-1]
	}
	if len(args) > 1 {
		request.Previous = args[len(args)-2]
	}

	// The sentinel means the user was actually typing "--" and wants long-flag
	// completion, not a real end-of-options terminator.
	if request.Current == completionDoubleDashSentinel {
		request.Current = "--"
		request.CompletingFlagName = true
		return request
	}

	if request.Current == "--" {
		request.CompletingFlagName = true
		return request
	}

	if flag, ok := flagByToken(cmd, request.Current); ok && flagTakesValue(flag) {
		request.FlagValueName = primaryFlagName(flag)
		return request
	}

	if strings.HasPrefix(request.Current, "-") {
		request.CompletingFlagName = true
		return request
	}

	if flag, ok := flagByToken(cmd, request.Previous); ok && flagTakesValue(flag) {
		request.FlagValueName = primaryFlagName(flag)
		request.FlagValuePrefix = request.Current
	}

	return request
}

func completionInvocationArgs(cmd *cli.Command) []string {
	args := os.Args[1:]
	if cmd == nil {
		return slices.Clone(args)
	}

	for i, arg := range args {
		if arg == cmd.Name {
			return slices.Clone(args[i+1:])
		}
	}

	return slices.Clone(args)
}

func CompleteFlagNames(cmd *cli.Command, prefix string) {
	PrintCandidates(FilterByPrefix(flagNameCompletionCandidates(cmd), prefix))
}

func CompleteFlagNameIfRequested(request CompletionRequest, cmd *cli.Command) bool {
	if !request.CompletingFlagName {
		return false
	}

	CompleteFlagNames(cmd, request.Current)
	return true
}

func CompleteFlagValueIfRequested(
	request CompletionRequest,
	candidatesByFlag map[string][]CompletionCandidate,
) bool {
	if request.FlagValueName == "" {
		return false
	}

	candidates, ok := candidatesByFlag[request.FlagValueName]
	if !ok {
		return false
	}

	PrintCandidates(FilterByPrefix(candidates, request.FlagValuePrefix))
	return true
}

func CompletePackageIDIfRequested(ctx context.Context, cmd *cli.Command, request CompletionRequest) {
	CompletePackageIDSuggestions(ctx, cmd, request.Current)
}

func flagNameCompletionCandidates(cmd *cli.Command) []CompletionCandidate {
	flags := completionFlags(cmd)

	seen := make(map[string]struct{})
	out := make([]CompletionCandidate, 0, len(flags)*2)

	for _, flag := range flags {
		usage := ""
		if docFlag, ok := flag.(cli.DocGenerationFlag); ok {
			usage = docFlag.GetUsage()
		}

		for _, name := range flag.Names() {
			if name == "" {
				continue
			}
			prefix := "--"
			if len(name) == 1 {
				prefix = "-"
			}

			value := prefix + name
			if _, ok := seen[value]; ok {
				continue
			}
			seen[value] = struct{}{}

			out = append(out, CompletionCandidate{Value: value, Description: usage})
		}
	}

	return out
}

func completionFlags(cmd *cli.Command) []cli.Flag {
	flags := append([]cli.Flag{}, cmd.VisibleFlags()...)
	flags = append(flags, cmd.VisiblePersistentFlags()...)
	return flags
}

func flagByToken(cmd *cli.Command, token string) (cli.Flag, bool) {
	if !strings.HasPrefix(token, "-") {
		return nil, false
	}

	name := strings.TrimLeft(token, "-")
	for _, flag := range completionFlags(cmd) {
		for _, candidate := range flag.Names() {
			if candidate == name {
				return flag, true
			}
		}
	}

	return nil, false
}

func primaryFlagName(flag cli.Flag) string {
	names := flag.Names()
	if len(names) == 0 {
		return ""
	}
	return names[0]
}

func flagTakesValue(flag cli.Flag) bool {
	docFlag, ok := flag.(cli.DocGenerationFlag)
	return ok && docFlag.TakesValue()
}
