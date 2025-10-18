export interface Condition {
  type: string
  status: string
  lastTransitionTime: string
  reason: string
  message: string
}

export interface ResourceList {
  cpu: string
  memory: string
  storage: string
  ephemeralStorage: string
  pods: string
}

export interface ClusterInfo {
  clusterID: string
  kubernetesVersion: string
  platform: string
  region: string
  openshiftVersion: string
  consoleURL: string
  apiURL: string
  networkType: string
  createdAt: string
}

export interface NodeInfo {
  name: string
  status: string
  role: string
  internalIP: string
  externalIP: string
  kernelVersion: string
  osImage: string
  containerRuntime: string
  kubeletVersion: string
  conditions: Condition[]
  capacity: ResourceList
  allocatable: ResourceList
  labels: Record<string, string>
  annotations: Record<string, string>
  createdAt: string
}

export interface PolicyInfo {
  name: string
  namespace: string
  remediationAction: string
  complianceState: string
  severity: string
  categories: string[]
  standards: string[]
  controls: string[]
  violations: number
  placementRules: string[]
  disabled: boolean
  labels: Record<string, string>
  annotations: Record<string, string>
  createdAt: string
}

export interface ManagedCluster {
  name: string
  namespace: string
  status: string
  version: string
  conditions: Condition[]
  clusterInfo: ClusterInfo
  nodesInfo: NodeInfo[]
  policiesInfo: PolicyInfo[]
  labels: Record<string, string>
  annotations: Record<string, string>
  hubName: string
  createdAt: string
}

export interface ManagedHub {
  name: string
  namespace: string
  status: string
  version: string
  conditions: Condition[]
  clusterInfo: ClusterInfo
  nodesInfo: NodeInfo[]
  policiesInfo: PolicyInfo[]
  managedClusters: ManagedCluster[]
  labels: Record<string, string>
  annotations: Record<string, string>
  createdAt: string
}

export interface APIResponse<T = any> {
  success: boolean
  message?: string
  data?: T
  error?: string
}

export interface HealthResponse {
  status: string
  version: string
  timestamp: string
}

