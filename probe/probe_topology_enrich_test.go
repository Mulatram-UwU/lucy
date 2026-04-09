package probe

import (
	"github.com/mclucy/lucy/types"
	"testing"
)

func TestNormalizeTopology_DeduplicatesNodes(t *testing.T) {
	nodeA := makeNode("a", types.CapabilityFabricMods)
	nodeA2 := makeNode("a", types.CapabilityForgeMods) // duplicate ID, different caps
	topo := makeTopology("a", []types.RuntimeNode{nodeA, nodeA2}, nil)
	NormalizeTopology(topo)
	if len(topo.Nodes) != 1 {
		t.Errorf("expected 1 node after dedup, got %d", len(topo.Nodes))
	}
	// last-write-wins: nodeA2 should survive
	if !topo.Nodes[0].HasCapability(types.CapabilityForgeMods) {
		t.Error("last-write-wins violated: expected nodeA2 (ForgeMods) to survive")
	}
}

func TestNormalizeTopology_DeduplicatesEdges(t *testing.T) {
	e := makeEdge("a", "b", types.EdgeBridges, types.RiskHigh)
	topo := makeTopology("a",
		[]types.RuntimeNode{makeNode("a"), makeNode("b")},
		[]types.RuntimeEdge{e, e}, // duplicate
	)
	NormalizeTopology(topo)
	if len(topo.Edges) != 1 {
		t.Errorf("expected 1 edge after dedup, got %d", len(topo.Edges))
	}
}

func TestNormalizeTopology_SortsNodes(t *testing.T) {
	topo := makeTopology("a", []types.RuntimeNode{
		makeNode("z"),
		makeNode("a"),
		makeNode("m"),
	}, nil)
	NormalizeTopology(topo)
	ids := []string{string(topo.Nodes[0].ID), string(topo.Nodes[1].ID), string(topo.Nodes[2].ID)}
	if ids[0] != "a" || ids[1] != "m" || ids[2] != "z" {
		t.Errorf("nodes not sorted: %v", ids)
	}
}

func TestNormalizeTopology_SortsEdges(t *testing.T) {
	topo := makeTopology("a",
		[]types.RuntimeNode{makeNode("a"), makeNode("b"), makeNode("c")},
		[]types.RuntimeEdge{
			makeEdge("b", "c", types.EdgeBridges, 0),
			makeEdge("a", "c", types.EdgeBridges, 0),
			makeEdge("a", "b", types.EdgeBridges, 0),
		},
	)
	NormalizeTopology(topo)
	if topo.Edges[0].From != "a" || topo.Edges[0].To != "b" {
		t.Errorf("edges not sorted correctly: first edge = %+v", topo.Edges[0])
	}
}

func TestNormalizeTopology_NilSafe(t *testing.T) {
	// Should not panic
	NormalizeTopology(nil)
}

func TestMergeTopology_AddsNewNodes(t *testing.T) {
	dst := makeTopology("a", []types.RuntimeNode{makeNode("a")}, nil)
	src := makeTopology("b", []types.RuntimeNode{makeNode("b")}, nil)
	mergeTopology(dst, src)
	if len(dst.Nodes) != 2 {
		t.Errorf("expected 2 nodes after merge, got %d", len(dst.Nodes))
	}
}

func TestMergeTopology_SkipsDuplicateNodes(t *testing.T) {
	dst := makeTopology("a", []types.RuntimeNode{makeNode("a")}, nil)
	src := makeTopology("a", []types.RuntimeNode{makeNode("a")}, nil)
	mergeTopology(dst, src)
	if len(dst.Nodes) != 1 {
		t.Errorf("expected 1 node (no dup), got %d", len(dst.Nodes))
	}
}

func TestMergeTopology_AddsNewEdges(t *testing.T) {
	dst := makeTopology("a", []types.RuntimeNode{makeNode("a"), makeNode("b")}, nil)
	src := makeTopology("a",
		[]types.RuntimeNode{makeNode("a"), makeNode("b")},
		[]types.RuntimeEdge{makeEdge("a", "b", types.EdgeBridges, 0)},
	)
	mergeTopology(dst, src)
	if len(dst.Edges) != 1 {
		t.Errorf("expected 1 edge after merge, got %d", len(dst.Edges))
	}
}

func TestMergeTopology_SkipsDuplicateEdges(t *testing.T) {
	e := makeEdge("a", "b", types.EdgeBridges, 0)
	dst := makeTopology("a",
		[]types.RuntimeNode{makeNode("a"), makeNode("b")},
		[]types.RuntimeEdge{e},
	)
	src := makeTopology("a",
		[]types.RuntimeNode{makeNode("a"), makeNode("b")},
		[]types.RuntimeEdge{e},
	)
	mergeTopology(dst, src)
	if len(dst.Edges) != 1 {
		t.Errorf("expected 1 edge (no dup), got %d", len(dst.Edges))
	}
}

func TestMergeTopology_NilSafe(t *testing.T) {
	dst := makeTopology("a", []types.RuntimeNode{makeNode("a")}, nil)
	mergeTopology(dst, nil)
	mergeTopology(nil, dst)
	// no panic
}
