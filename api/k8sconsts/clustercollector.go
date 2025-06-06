package k8sconsts

const (
	OdigosClusterCollectorImage     = "odigos-collector"
	OdigosClusterCollectorImageUBI9 = "odigos-collector-ubi9"

	OdigosClusterCollectorDeploymentName = "odigos-gateway"
	OdigosClusterCollectorConfigMapName  = OdigosClusterCollectorDeploymentName
	OdigosClusterCollectorServiceName    = OdigosClusterCollectorDeploymentName

	OdigosClusterCollectorCollectorGroupName = OdigosClusterCollectorDeploymentName
	OdigosClusterCollectorConfigMapKey       = "collector-conf"

	OdigosClusterCollectorServiceAccountName = "odigos-gateway"
	OdigosClusterCollectorRoleName           = "odigos-gateway"
	OdigosClusterCollectorRoleBindingName    = "odigos-gateway"

	// The cluster gateway collector runs as a deployment and the pod is exposed as a service.
	// Thus it cannot collide with other ports on the same node, and we can use an handy default port.
	OdigosClusterCollectorOwnTelemetryPortDefault = int32(8888)
)
