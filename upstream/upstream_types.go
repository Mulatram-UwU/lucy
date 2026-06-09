package upstream

import "github.com/mclucy/lucy/types"

// Provider is the inversion boundary between core upstream orchestration and
// concrete upstream integrations.
//
// Provider is an executable capability boundary: implementations perform actual
// native API calls and adapt upstream-specific data into raw contracts.
//
// Provider is intentionally not the same concept as types.Source:
//   - types.Source is a stable semantic identifier exposed to users and storage.
//   - Provider is the runtime executor selected by routing logic.
//
// Rules:
//   - Core code depends on this interface, never on concrete provider packages.
//   - Provider packages implement this interface and perform upstream-specific
//     API/data handling.
//   - Source selection/fallback policy is handled by dedicated resolver logic
//     outside this file.
type Provider interface {
	SearchLegacy(query string, options types.SearchOptions) (
		res RawSearchResults,
		err error,
	)
	Fetch(id types.VersionedPackageRef) (
		remote RawPackageRemote,
		err error,
	)
	Metadata(name types.BarePackageName) (
		info RawProjectInformation,
		err error,
	)
	Dependencies(id types.VersionedPackageRef) (
		deps RawPackageDependencies,
		err error,
	)
	Support(name types.BarePackageName) (
		supports RawProjectSupport,
		err error,
	)
	ParseAmbiguousId(id types.VersionedPackageRef) (
		parsed types.VersionedPackageRef,
		err error,
	)
	// Id returns the semantic source identity represented by this provider.
	Id() types.SourceId
}

type FetchResult struct {
	ResolvedID types.VersionedPackageRef
	Remote     types.PackageRemote
}

// Raw interfaces are internal conversion contracts returned by providers before
// being normalized into types.* structures.

type (
	RawProjectSupport interface {
		ToProjectSupport() types.PlatformSupport
	}
	RawProjectInformation interface {
		ToProjectInformation() types.Metadata
	}
	RawPackageRemote interface {
		ToPackageRemote() types.PackageRemote
	}
	RawPackageDependencies interface {
		ToPackageDependencies() types.PackageDependencies
	}

	// TODO: Consider make SortBy a method on the RawSearchResults interface

	RawSearchResults interface {
		ToSearchResults() types.SearchResults
	}
)

type RemotePackageName struct {
	RemoteName string
	Source     types.SourceId
}
