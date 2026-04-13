package install

import (
	"fmt"

	"github.com/mclucy/lucy/logger"
	"github.com/mclucy/lucy/probe"
	"github.com/mclucy/lucy/types"
)

// SnapshotInstalledConstraints reads installed packages from the probe snapshot
// and converts them into fixed InstalledConstraint entries for the transaction.
// Each installed package is treated as an immutable anchor during recursive
// solving; it will never be auto-replaced by the solver.
func SnapshotInstalledConstraints(tx *RecursiveTransaction) {
	si := probe.ServerInfo()
	constraints := make([]InstalledConstraint, 0, len(si.Packages))

	for _, pkg := range si.Packages {
		if pkg.Id.Version.IsInvalid() || pkg.Id.Version == types.VersionAny {
			continue
		}
		ic := InstalledConstraint{
			Package: pkg,
			ConstraintInput: ConstraintInput{
				Requester: fmt.Sprintf("installed:%s", pkg.Id.StringFull()),
				Dependency: types.Dependency{
					Id:        pkg.Id,
					Mandatory: true,
				},
			},
		}
		constraints = append(constraints, ic)
	}

	tx.InstalledConstraints = constraints
}

// FindCompatibleInstalled searches the installed-constraint snapshot for any
// package with the same platform and name as the requested ID, returning all
// matches. Results are informational only; the solver must not auto-select them.
func FindCompatibleInstalled(tx *RecursiveTransaction, id types.PackageId) []types.Package {
	var matches []types.Package
	for _, ic := range tx.InstalledConstraints {
		pkg := ic.Package
		if pkg.Id.Platform != id.Platform {
			continue
		}
		if pkg.Id.Name != id.Name {
			continue
		}
		matches = append(matches, pkg)
	}
	return matches
}

// ReportCompatibleInstalled logs any locally-installed versions that are
// compatible with the given package ID. This is an informational-only report;
// no automatic selection occurs.
func ReportCompatibleInstalled(tx *RecursiveTransaction, id types.PackageId) {
	matches := FindCompatibleInstalled(tx, id)
	for _, pkg := range matches {
		logger.ShowInfo(fmt.Sprintf(
			"[recursive] compatible installed version found: %s (not auto-selected)",
			pkg.Id.StringFull(),
		))
	}
}
