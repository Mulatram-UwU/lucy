package types

type RuntimeNodeID string

const RuntimeNodeUnknown RuntimeNodeID = ""

type RuntimeRole string

const (
	RuntimeRoleModLoader      RuntimeRole = "mod_loader"
	RuntimeRolePluginCore     RuntimeRole = "plugin_core"
	RuntimeRoleHybrid         RuntimeRole = "hybrid"
	RuntimeRoleProxy          RuntimeRole = "proxy"
	RuntimeRoleBridge         RuntimeRole = "bridge"
	RuntimeRoleProtocolBridge RuntimeRole = "protocol_bridge"
	RuntimeRoleVanilla        RuntimeRole = "vanilla"
	RuntimeRoleUnknown        RuntimeRole = ""
)

type RuntimeCapability string

const (
	CapabilityFabricMods     RuntimeCapability = "fabric_mods"
	CapabilityForgeMods      RuntimeCapability = "forge_mods"
	CapabilityNeoforgeMods   RuntimeCapability = "neoforge_mods"
	CapabilityBukkitPlugins  RuntimeCapability = "bukkit_plugins"
	CapabilityMCDRPlugins    RuntimeCapability = "mcdr_plugins"
	CapabilityProxying       RuntimeCapability = "proxying"
	CapabilityProtocolBridge RuntimeCapability = "protocol_bridge"
)

type RuntimeRiskLevel int

const (
	RiskNone     RuntimeRiskLevel = 0
	RiskLow      RuntimeRiskLevel = 1
	RiskMedium   RuntimeRiskLevel = 2
	RiskHigh     RuntimeRiskLevel = 3
	RiskCritical RuntimeRiskLevel = 4
)

type CompatVerdict string

const (
	CompatCompatible   CompatVerdict = "compatible"
	CompatDegraded     CompatVerdict = "degraded"
	CompatIncompatible CompatVerdict = "incompatible"
	CompatUnresolved   CompatVerdict = "unresolved"
)

type CompatResult struct {
	Verdict   CompatVerdict
	Reason    string
	Detail    string
	RiskLevel RuntimeRiskLevel
}

type RuntimeNode struct {
	ID               RuntimeNodeID
	Role             RuntimeRole
	IdentityPlatform Platform
	Capabilities     []RuntimeCapability
	RiskLevel        RuntimeRiskLevel
}

func (n RuntimeNode) HasCapability(c RuntimeCapability) bool {
	for _, capability := range n.Capabilities {
		if capability == c {
			return true
		}
	}

	return false
}

type RuntimeEdgeKind string

const (
	EdgeHosts   RuntimeEdgeKind = "hosts"
	EdgeBridges RuntimeEdgeKind = "bridges"
	EdgeRoutes  RuntimeEdgeKind = "routes"
	EdgeAdapts  RuntimeEdgeKind = "adapts"
)

type RuntimeEdge struct {
	From RuntimeNodeID
	To   RuntimeNodeID
	Kind RuntimeEdgeKind
	Risk RuntimeRiskLevel
}

type RuntimeTopology struct {
	PrimaryNode RuntimeNodeID
	Nodes       []RuntimeNode
	Edges       []RuntimeEdge
}

func (t *RuntimeTopology) Resolved() bool {
	return t != nil && t.PrimaryNode != RuntimeNodeUnknown && len(t.Nodes) > 0
}

func (t *RuntimeTopology) FindNode(id RuntimeNodeID) (RuntimeNode, bool) {
	if t == nil {
		return RuntimeNode{}, false
	}

	for _, node := range t.Nodes {
		if node.ID == id {
			return node, true
		}
	}

	return RuntimeNode{}, false
}

func (t *RuntimeTopology) HasCapability(c RuntimeCapability) bool {
	if t == nil {
		return false
	}

	for _, node := range t.Nodes {
		if node.HasCapability(c) {
			return true
		}
	}

	return false
}

func (t *RuntimeTopology) PrimaryNodeData() (RuntimeNode, bool) {
	if t == nil {
		return RuntimeNode{}, false
	}

	return t.FindNode(t.PrimaryNode)
}
