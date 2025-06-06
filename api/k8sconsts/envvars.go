package k8sconsts

const (
	OdigosEnvVarNamespace              = "ODIGOS_WORKLOAD_NAMESPACE"
	OdigosEnvVarContainerName          = "ODIGOS_CONTAINER_NAME"
	OdigosEnvVarPodName                = "ODIGOS_POD_NAME"
	OdigosEnvVarDistroName             = "ODIGOS_DISTRO_NAME"
	CustomContainerRuntimeSocketEnvVar = "CONTAINER_RUNTIME_SOCK"
)

func OdigosInjectedEnvVars() []string {
	return []string{
		OdigosEnvVarNamespace,
		OdigosEnvVarContainerName,
		OdigosEnvVarPodName,
		OdigosEnvVarDistroName,
	}
}
