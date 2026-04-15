package state

import (
	"bytes"
	"testing"
)

func TestManifestRoundTrip(t *testing.T) {
	manifestText := []byte(`
[format]
version = "v1"

[environment]
game_version = "1.21.1"
platform = "fabric"
platform_version = "0.16.10"

[sources]
	[[sources.custom]]
	name = "mirror"
	url = "https://example.com/index.json"
	type = "modrinth"

[layout]
mods_dir = "mods"
plugins_dir = "plugins"
config_dir = "config"

[policy]
managed_roots = ["mods", "plugins"]
unmanaged_paths = ["world/**"]

[[packages]]
id = "fabric/fabric-api"
version = "compatible"
source = "auto"
role = "required"
side = "both"
optional = false
pinned = false

[[bundles]]
name = "server-config"
type = "config"
path = "config"
source = "./defaults/config"
optional = false
`)

	var manifest Manifest
	if err := manifest.Unmarshal(manifestText); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	first, err := manifest.Marshal()
	if err != nil {
		t.Fatalf("first marshal failed: %v", err)
	}

	var parsed Manifest
	if err := parsed.Unmarshal(first); err != nil {
		t.Fatalf("re-unmarshal failed: %v", err)
	}

	second, err := parsed.Marshal()
	if err != nil {
		t.Fatalf("second marshal failed: %v", err)
	}

	if !bytes.Equal(first, second) {
		t.Fatalf("round-trip produced different output:\nfirst:\n%s\nsecond:\n%s", first, second)
	}

	if err := ValidateManifest(parsed); err != nil {
		t.Fatalf("round-trip manifest should validate: %v", err)
	}
}

func TestManifestPreservesPackageSides(t *testing.T) {
	manifestText := []byte(`
[format]
version = "v1"

[environment]
game_version = "1.21.1"
platform = "fabric"
platform_version = "0.16.10"

[layout]
mods_dir = "mods"
plugins_dir = "plugins"
config_dir = "config"

[policy]
managed_roots = ["mods", "plugins"]
unmanaged_paths = []

[[packages]]
id = "fabric/lithium"
version = "compatible"
source = "auto"
role = "required"
side = "server"
optional = false
pinned = false

[[packages]]
id = "fabric/sodium"
version = "latest"
source = "modrinth"
role = "transitive"
side = "client"
optional = true
pinned = false

[[packages]]
id = "fabric/fabric-api"
version = "0.119.2+1.21.5"
source = "github"
role = "ignored"
side = "both"
optional = false
pinned = true
`)

	var manifest Manifest
	if err := manifest.Unmarshal(manifestText); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(manifest.Packages) != 3 {
		t.Fatalf("expected 3 packages, got %d", len(manifest.Packages))
	}

	gotSides := []ManifestSide{
		manifest.Packages[0].Side,
		manifest.Packages[1].Side,
		manifest.Packages[2].Side,
	}
	wantSides := []ManifestSide{SideServer, SideClient, SideBoth}
	for i := range wantSides {
		if gotSides[i] != wantSides[i] {
			t.Fatalf("package %d side mismatch: got %q want %q", i, gotSides[i], wantSides[i])
		}
	}

	if !manifest.Packages[1].Optional {
		t.Fatalf("expected client package to remain optional")
	}

	if !manifest.Packages[2].Pinned {
		t.Fatalf("expected pinned flag to remain true")
	}

	if err := ValidateManifest(manifest); err != nil {
		t.Fatalf("side-aware manifest should validate: %v", err)
	}
}

func TestManifestPreservesFuzzyVersionIntentVerbatim(t *testing.T) {
	manifestText := []byte(`
[format]
version = "v1"

[environment]
game_version = "1.21.1"
platform = "fabric"
platform_version = "0.16.10"

[layout]
mods_dir = "mods"
plugins_dir = "plugins"
config_dir = "config"

[policy]
managed_roots = ["mods", "plugins"]
unmanaged_paths = []

[[packages]]
id = "fabric/lithium"
version = ">=0.12.0 <0.13.0"
source = "auto"
role = "required"
side = "server"
optional = false
pinned = false
`)

	manifest, err := ParseManifest(manifestText)
	if err != nil {
		t.Fatalf("expected fuzzy manifest version intent to parse: %v", err)
	}

	if got := manifest.Packages[0].Version; got != ">=0.12.0 <0.13.0" {
		t.Fatalf("expected manifest to preserve fuzzy selector verbatim, got %q", got)
	}

	data, err := SerializeManifest(manifest)
	if err != nil {
		t.Fatalf("serialize manifest failed: %v", err)
	}

	reparsed, err := ParseManifest(data)
	if err != nil {
		t.Fatalf("reparse manifest failed: %v", err)
	}
	if got := reparsed.Packages[0].Version; got != ">=0.12.0 <0.13.0" {
		t.Fatalf("expected fuzzy selector to survive round trip, got %q", got)
	}
	if reparsed.Packages[0].Version == "0.12.9" {
		t.Fatal("manifest intent must not be rewritten to an exact resolved version")
	}
}

func TestManifestBundlesRemainSeparateFromPackages(t *testing.T) {
	manifestText := []byte(`
[format]
version = "v1"

[environment]
game_version = "1.20.6"
platform = "none"
platform_version = ""

[layout]
mods_dir = "mods"
plugins_dir = "plugins"
config_dir = "config"

[policy]
managed_roots = ["mods", "plugins"]
unmanaged_paths = []

[[packages]]
id = "none/luckperms"
version = "latest"
source = "curseforge"
role = "required"
side = "server"
optional = false
pinned = false

[[bundles]]
name = "default-config"
type = "config"
path = "config/luckperms"
source = "./overlays/luckperms"
optional = false

[[bundles]]
name = "vanilla-datapack"
type = "datapack"
path = "world/datapacks/vanilla"
source = "./overlays/vanilla"
optional = true
`)

	var manifest Manifest
	if err := manifest.Unmarshal(manifestText); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(manifest.Packages) != 1 {
		t.Fatalf("expected 1 package, got %d", len(manifest.Packages))
	}
	if len(manifest.Bundles) != 2 {
		t.Fatalf("expected 2 bundles, got %d", len(manifest.Bundles))
	}

	if manifest.Packages[0].ID == manifest.Bundles[0].Name {
		t.Fatalf("package identity space should remain separate from bundle names")
	}

	if manifest.Bundles[0].Type != BundleTypeConfig {
		t.Fatalf("expected first bundle type %q, got %q", BundleTypeConfig, manifest.Bundles[0].Type)
	}
	if manifest.Bundles[1].Type != BundleTypeDatapack {
		t.Fatalf("expected second bundle type %q, got %q", BundleTypeDatapack, manifest.Bundles[1].Type)
	}

	if err := ValidateManifest(manifest); err != nil {
		t.Fatalf("bundled manifest should validate: %v", err)
	}
}

func TestManifestDefaults(t *testing.T) {
	manifest := ManifestDefaults()

	if manifest.Format.Version != "v1" {
		t.Fatalf("expected format version v1, got %q", manifest.Format.Version)
	}
	if manifest.Environment.Platform != "none" {
		t.Fatalf("expected default platform none, got %q", manifest.Environment.Platform)
	}
	if manifest.Layout.ModsDir != "mods" {
		t.Fatalf("expected mods_dir mods, got %q", manifest.Layout.ModsDir)
	}
	if manifest.Layout.PluginsDir != "plugins" {
		t.Fatalf("expected plugins_dir plugins, got %q", manifest.Layout.PluginsDir)
	}
	if manifest.Layout.ConfigDir != "config" {
		t.Fatalf("expected config_dir config, got %q", manifest.Layout.ConfigDir)
	}
	if len(manifest.Policy.ManagedRoots) != 2 {
		t.Fatalf("expected 2 managed roots, got %d", len(manifest.Policy.ManagedRoots))
	}
	if manifest.Policy.ManagedRoots[0] != "mods" || manifest.Policy.ManagedRoots[1] != "plugins" {
		t.Fatalf("unexpected managed roots: %#v", manifest.Policy.ManagedRoots)
	}
	if len(manifest.Sources.Custom) != 0 {
		t.Fatalf("expected no custom sources by default")
	}
	if len(manifest.Packages) != 0 || len(manifest.Bundles) != 0 {
		t.Fatalf("expected empty package and bundle declarations by default")
	}

	if err := ValidateManifest(manifest); err != nil {
		t.Fatalf("default manifest should validate: %v", err)
	}
}
