package mcdr

import (
	"fmt"

	"github.com/mclucy/lucy/logger"
	"github.com/mclucy/lucy/probe"
	"github.com/mclucy/lucy/syntax"
	"github.com/mclucy/lucy/types"
	"github.com/mclucy/lucy/upstream"
)

type provider struct{}

func (s provider) Id() types.SourceId {
	return types.SourceMCDR
}

var Provider provider

// Just a trivial type to implement the SearchResults interface
type mcdrSearchResult []string

func (m mcdrSearchResult) ToSearchResults() upstream.SearchResponse {
	var res upstream.SearchResponse
	for _, id := range m {
		res.Projects = append(res.Projects, syntax.ToProjectName(id))
	}
	res.Source = types.SourceMCDR
	return res
}

// TODO: handle search options

func (s provider) SearchLegacy(
	query string,
	options types.SearchOptions,
) (res upstream.RawSearchResults, err error) {
	if options.FilterPlatform != types.PlatformMCDR && options.FilterPlatform != types.PlatformAny {
		return nil, fmt.Errorf(
			"invalid search platform: expected %s, got %s",
			types.PlatformMCDR,
			options.FilterPlatform,
		)
	}
	res, err = search(query)
	return
}

func (s provider) Fetch(id types.VersionedPackageRef) (
	rem upstream.RawPackageRemote,
	err error,
) {
	rem, err = getRelease(id.Name.Pep8String(), id.Version)
	return
}

func (s provider) Metadata(name types.BarePackageName) (
	info upstream.RawProjectInformation,
	err error,
) {
	plugin, err := getInfo(name.Pep8String())
	if err != nil {
		return nil, err
	}
	meta, err := getMeta(name.Pep8String())
	if err != nil {
		return nil, err
	}
	repo, err := getRepository(name.Pep8String())
	if err != nil {
		return nil, err
	}

	info = rawProjectInformation{
		Info:       plugin,
		Meta:       meta,
		Repository: repo,
	}

	return info, nil
}

func (s provider) Dependencies(id types.VersionedPackageRef) (
	upstream.RawPackageDependencies,
	error,
) {
	// TODO implement me
	panic("implement me")
}

func (s provider) Support(name types.BarePackageName) (
	supports upstream.RawProjectSupport,
	err error,
) {
	// TODO implement me
	panic("implement me")
}

func (s provider) ParseAmbiguousId(id types.VersionedPackageRef) (
	parsed types.VersionedPackageRef,
	err error,
) {
	var rel *release
	switch id.Version {
	case types.VersionCompatible:
		serverInfo := probe.ServerInfo()
		rel, err = getLatestCompatibleRelease(
			id.Name.Pep8String(),
			serverInfo.Environments.Mcdr.Version,
		)
	case types.VersionLatest, types.VersionAny:
		rel, err = getLatestRelease(id.Name.Pep8String())
		if err != nil {
			return id, err
		}
	default:
		return id, fmt.Errorf(
			"cannot parse version %s for package %s",
			id.Version,
			id.Name,
		)
	}
	if err != nil {
		return id, err
	}
	parsed = types.VersionedPackageRef{
		Platform: types.PlatformMCDR,
		Name:     id.Name,
		Version:  types.BareVersion(rel.Meta.Version),
	}
	logger.Debug("parsed from" + id.StringFull() + " to " + parsed.StringFull())
	return parsed, nil
}
