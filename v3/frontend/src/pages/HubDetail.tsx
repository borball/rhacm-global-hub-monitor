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
import ClusterCard from '@/components/ClusterCard'
import NodeCard from '@/components/NodeCard'
import PolicyTable from '@/components/PolicyTable'
import { formatDate } from '@/utils/helpers'

const HubDetail: React.FC = () => {
  const { hubName } = useParams<{ hubName: string }>()
  const [activeTabKey, setActiveTabKey] = React.useState<string | number>(0)

  const { data: hub, isLoading, error } = useQuery(
    ['hub', hubName],
    () => hubsAPI.getHub(hubName!),
    { enabled: !!hubName }
  )

  if (isLoading) {
    return (
      <div className="loading-container">
        <Spinner size="xl" />
      </div>
    )
  }

  if (error || !hub) {
    return (
      <div className="error-container">
        <Alert variant="danger" title="Failed to load hub details">
          {error ? (error as Error).message : 'Hub not found'}
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
          <BreadcrumbItem isActive>{hubName}</BreadcrumbItem>
        </Breadcrumb>
        <Title headingLevel="h1" size="2xl" style={{ marginTop: '1rem' }}>
          {hub.name}
        </Title>
        <StatusLabel status={hub.status} />
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
                    <DescriptionListDescription>{hub.name}</DescriptionListDescription>
                  </DescriptionListGroup>
                  <DescriptionListGroup>
                    <DescriptionListTerm>Status</DescriptionListTerm>
                    <DescriptionListDescription>
                      <StatusLabel status={hub.status} />
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                  <DescriptionListGroup>
                    <DescriptionListTerm>Version</DescriptionListTerm>
                    <DescriptionListDescription>{hub.version || 'N/A'}</DescriptionListDescription>
                  </DescriptionListGroup>
                  <DescriptionListGroup>
                    <DescriptionListTerm>OpenShift Version</DescriptionListTerm>
                    <DescriptionListDescription>
                      {hub.clusterInfo.openshiftVersion || 'N/A'}
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                  <DescriptionListGroup>
                    <DescriptionListTerm>Platform</DescriptionListTerm>
                    <DescriptionListDescription>
                      {hub.clusterInfo.platform || 'N/A'}
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                  <DescriptionListGroup>
                    <DescriptionListTerm>Region</DescriptionListTerm>
                    <DescriptionListDescription>
                      {hub.clusterInfo.region || 'N/A'}
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                  <DescriptionListGroup>
                    <DescriptionListTerm>Console URL</DescriptionListTerm>
                    <DescriptionListDescription>
                      {hub.clusterInfo.consoleURL ? (
                        <a href={hub.clusterInfo.consoleURL} target="_blank" rel="noopener noreferrer">
                          {hub.clusterInfo.consoleURL}
                        </a>
                      ) : (
                        'N/A'
                      )}
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                  <DescriptionListGroup>
                    <DescriptionListTerm>Created</DescriptionListTerm>
                    <DescriptionListDescription>
                      {formatDate(hub.createdAt)}
                    </DescriptionListDescription>
                  </DescriptionListGroup>
                </DescriptionList>
              </CardBody>
            </Card>
          </Tab>

          <Tab eventKey={1} title={<TabTitleText>Managed Clusters ({hub.managedClusters.length})</TabTitleText>}>
            <div style={{ padding: '1rem 0' }}>
              {hub.managedClusters.length > 0 ? (
                <Grid hasGutter>
                  {hub.managedClusters.map((cluster) => (
                    <GridItem key={cluster.name} sm={12} md={6} lg={4}>
                      <ClusterCard cluster={cluster} hubName={hub.name} />
                    </GridItem>
                  ))}
                </Grid>
              ) : (
                <Card>
                  <CardBody>No managed clusters found</CardBody>
                </Card>
              )}
            </div>
          </Tab>

          <Tab eventKey={2} title={<TabTitleText>Nodes ({hub.nodesInfo.length})</TabTitleText>}>
            <div style={{ padding: '1rem 0' }}>
              {hub.nodesInfo.length > 0 ? (
                <Grid hasGutter>
                  {hub.nodesInfo.map((node) => (
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

          <Tab eventKey={3} title={<TabTitleText>Policies ({hub.policiesInfo.length})</TabTitleText>}>
            <div style={{ padding: '1rem 0' }}>
              <Card>
                <CardBody>
                  <PolicyTable policies={hub.policiesInfo} />
                </CardBody>
              </Card>
            </div>
          </Tab>
        </Tabs>
      </PageSection>
    </>
  )
}

export default HubDetail

