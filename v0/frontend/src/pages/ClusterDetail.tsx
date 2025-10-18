import React from 'react'
import { useParams, Link } from 'react-router-dom'
import { useQuery } from 'react-query'
import {
  PageSection,
  Title,
  Spinner,
  Alert,
  Tabs,
  Tab,
  TabTitleText,
  Card,
  CardBody,
  Grid,
  GridItem,
  DescriptionList,
  DescriptionListGroup,
  DescriptionListTerm,
  DescriptionListDescription,
  Breadcrumb,
  BreadcrumbItem,
} from '@patternfly/react-core'
import { hubsAPI } from '@/services/api'
import StatusLabel from '@/components/StatusLabel'
import NodeCard from '@/components/NodeCard'
import PolicyTable from '@/components/PolicyTable'
import { formatDate } from '@/utils/helpers'

const ClusterDetail: React.FC = () => {
  const { hubName, clusterName } = useParams<{ hubName: string; clusterName: string }>()
  const [activeTabKey, setActiveTabKey] = React.useState<string | number>(0)

  const { data: clusters, isLoading, error } = useQuery(
    ['hubClusters', hubName],
    () => hubsAPI.getHubClusters(hubName!),
    { enabled: !!hubName }
  )

  const cluster = clusters?.find((c) => c.name === clusterName)

  if (isLoading) {
    return (
      <div className="loading-container">
        <Spinner size="xl" />
      </div>
    )
  }

  if (error || !cluster) {
    return (
      <div className="error-container">
        <Alert variant="danger" title="Failed to load cluster details">
          {error ? (error as Error).message : 'Cluster not found'}
        </Alert>
      </div>
    )
  }

  return (
    <>
      <PageSection>
        <Breadcrumb>
          <BreadcrumbItem to="/hubs" component={Link}>
            Managed Hubs
          </BreadcrumbItem>
          <BreadcrumbItem to={`/hubs/${hubName}`} component={Link}>
            {hubName}
          </BreadcrumbItem>
          <BreadcrumbItem isActive>{clusterName}</BreadcrumbItem>
        </Breadcrumb>
        <Title headingLevel="h1" size="2xl" style={{ marginTop: '1rem' }}>
          {cluster.name}
        </Title>
        <StatusLabel status={cluster.status} />
      </PageSection>

      <PageSection>
        <Tabs
          activeKey={activeTabKey}
          onSelect={(_, tabIndex) => setActiveTabKey(tabIndex)}
        >
          <Tab eventKey={0} title={<TabTitleText>Overview</TabTitleText>}>
            <Card>
              <CardBody>
                <Title headingLevel="h3" size="lg" style={{ marginBottom: '1rem' }}>
                  Cluster Information
                </Title>
                <DescriptionList isHorizontal>
                  <DescriptionListGroup>
                    <DescriptionListTerm>Name</DescriptionListTerm>
                    <DescriptionListDescription>{cluster.name}</DescriptionListDescription>
                  </DescriptionListGroup>
                  <DescriptionListGroup>
                    <DescriptionListTerm>Status</DescriptionListTerm>
                    <DescriptionListDescription>
                      <StatusLabel status={cluster.status} />
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                  <DescriptionListGroup>
                    <DescriptionListTerm>Version</DescriptionListTerm>
                    <DescriptionListDescription>
                      {cluster.version || 'N/A'}
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                  <DescriptionListGroup>
                    <DescriptionListTerm>OpenShift Version</DescriptionListTerm>
                    <DescriptionListDescription>
                      {cluster.clusterInfo.openshiftVersion || 'N/A'}
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                  <DescriptionListGroup>
                    <DescriptionListTerm>Platform</DescriptionListTerm>
                    <DescriptionListDescription>
                      {cluster.clusterInfo.platform || 'N/A'}
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                  <DescriptionListGroup>
                    <DescriptionListTerm>Managed By</DescriptionListTerm>
                    <DescriptionListDescription>
                      <Link to={`/hubs/${cluster.hubName}`}>{cluster.hubName}</Link>
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                  <DescriptionListGroup>
                    <DescriptionListTerm>Console URL</DescriptionListTerm>
                    <DescriptionListDescription>
                      {cluster.clusterInfo.consoleURL ? (
                        <a href={cluster.clusterInfo.consoleURL} target="_blank" rel="noopener noreferrer">
                          {cluster.clusterInfo.consoleURL}
                        </a>
                      ) : (
                        'N/A'
                      )}
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                  <DescriptionListGroup>
                    <DescriptionListTerm>Created</DescriptionListTerm>
                    <DescriptionListDescription>
                      {formatDate(cluster.createdAt)}
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                </DescriptionList>
              </CardBody>
            </Card>
          </Tab>

          <Tab eventKey={1} title={<TabTitleText>Nodes ({cluster.nodesInfo.length})</TabTitleText>}>
            <div style={{ padding: '1rem 0' }}>
              {cluster.nodesInfo.length > 0 ? (
                <Grid hasGutter>
                  {cluster.nodesInfo.map((node) => (
                    <GridItem key={node.name} sm={12} md={6} lg={4}>
                      <NodeCard node={node} />
                    </GridItem>
                  ))}
                </Grid>
              ) : (
                <Card>
                  <CardBody>No node information available</CardBody>
                </Card>
              )}
            </div>
          </Tab>

          <Tab eventKey={2} title={<TabTitleText>Policies ({cluster.policiesInfo.length})</TabTitleText>}>
            <div style={{ padding: '1rem 0' }}>
              <Card>
                <CardBody>
                  <PolicyTable policies={cluster.policiesInfo} />
                </CardBody>
              </Card>
            </div>
          </Tab>
        </Tabs>
      </PageSection>
    </>
  )
}

export default ClusterDetail

