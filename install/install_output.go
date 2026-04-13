package install

import (
	"fmt"

	"github.com/mclucy/lucy/logger"
	"github.com/mclucy/lucy/types"
)

func showFetchStart(id types.PackageId) {
	logger.ShowInfo(
		fmt.Sprintf(
			"fetching package metadata for %s",
			id.StringFull(),
		),
	)
}

func showFetchSuccess(p types.Package) {
	if p.Remote == nil {
		return
	}
	logger.ShowInfo(
		fmt.Sprintf(
			"package metadata fetched from %s, resolved to %s",
			p.Remote.Source.String(),
			p.Id.StringFull(),
		),
	)
}

func showDownloadStart(url string) {
	logger.ShowInfo(fmt.Sprintf("downloading from %s", url))
}

func showInstallComplete(path string) {
	logger.ShowInfo(fmt.Sprintf("installed package to %s", path))
}

func showBatchPhase(header string, ids []types.PackageId) {
	logger.ShowInfo(fmt.Sprintf("==> %s: %s", header, joinPackageNames(ids)))
}

func showBatchSummary(installed int, failed int) {
	if failed == 0 {
		logger.ShowInfo(fmt.Sprintf("%d packages installed", installed))
	} else {
		logger.ShowInfo(fmt.Sprintf("%d installed, %d failed", installed, failed))
	}
}

func joinPackageNames(ids []types.PackageId) string {
	if len(ids) == 0 {
		return ""
	}
	if len(ids) == 1 {
		return ids[0].StringFull()
	}
	if len(ids) == 2 {
		return ids[0].StringFull() + " and " + ids[1].StringFull()
	}
	result := ""
	for i := 0; i < len(ids)-1; i++ {
		result += ids[i].StringFull() + ", "
	}
	result += "and " + ids[len(ids)-1].StringFull()
	return result
}

func showRecursiveResolveStart(roots []types.PackageId) {
	logger.ShowInfo(fmt.Sprintf("resolving dependencies for %s", joinPackageNames(roots)))
}

func showRecursiveDownloadStart(count int) {
	logger.ShowInfo(fmt.Sprintf("downloading %d artifacts", count))
}

func showRecursiveVerifyStart(count int) {
	logger.ShowInfo(fmt.Sprintf("verifying %d artifacts locally", count))
}

func showRecursiveReconcileStart() {
	logger.ShowInfo("reconciling advisory and verified graphs")
}

func showRecursiveReconcileDiff(diff ReconcileDiff) {
	verbals := []string{}
	if len(diff.Missing) > 0 {
		verbals = append(verbals, fmt.Sprintf("+%d missing", len(diff.Missing)))
	}
	if len(diff.Extra) > 0 {
		verbals = append(verbals, fmt.Sprintf("-%d extra", len(diff.Extra)))
	}
	if len(diff.Tightened) > 0 {
		verbals = append(verbals, fmt.Sprintf("~%d tightened", len(diff.Tightened)))
	}
	logger.ShowInfo("reconcile: " + joinStrings(verbals))
}

func showRecursiveApplyStart(count int) {
	logger.ShowInfo(fmt.Sprintf("applying %d changes", count))
}

func showRecursiveConflict(err error) {
	logger.ShowInfo(fmt.Sprintf("conflict: %s", err.Error()))
}

func showRecursiveCompatibleInstalled(id types.PackageId, installed types.Package) {
	logger.ShowInfo(fmt.Sprintf("[recursive] compatible installed: %s (not auto-selected)", installed.Id.StringFull()))
}

func joinStrings(strs []string) string {
	if len(strs) == 0 {
		return "none"
	}
	result := ""
	for i := 0; i < len(strs)-1; i++ {
		result += strs[i] + ", "
	}
	result += strs[len(strs)-1]
	return result
}
