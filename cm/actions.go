package cm

import "fmt"

// ListClusters get all the registered clusters
func (c CMServer) ListClusters() []Cluster {
	var cmItems CMItems
	var uri = "clusters"
	if c.UseGateway {
		curlCommand := c.CreateGatewayCurlGetCommand(uri)
		cmItems = ProcessCMItemsFromSSHResponse(c.RunGatewayCMCommand(curlCommand, false, true))
	} else {
		request := c.CreateGetRequest(uri)
		cmItems = ProcessCMItems(request)
	}
	return cmItems.ConvertClustersResponse()
}

// ListHosts get all the registered hosts
func (c CMServer) ListHosts() []Host {
	var cmItems CMItems
	var uri = "hosts"
	if c.UseGateway {
		curlCommand := c.CreateGatewayCurlGetCommand(uri)
		cmItems = ProcessCMItemsFromSSHResponse(c.RunGatewayCMCommand(curlCommand, false, true))
	} else {
		request := c.CreateGetRequest(uri)
		cmItems = ProcessCMItems(request)
	}
	return cmItems.ConvertHostsResponse()
}

// ListServices get all the registered services per cluster
func (c CMServer) ListServices(cluster string) []Service {
	var cmItems CMItems
	var uri = fmt.Sprintf("clusters/%v/services", cluster)
	if c.UseGateway {
		curlCommand := c.CreateGatewayCurlGetCommand(uri)
		cmItems = ProcessCMItemsFromSSHResponse(c.RunGatewayCMCommand(curlCommand, false, true))
	} else {
		request := c.CreateGetRequest(uri)
		cmItems = ProcessCMItems(request)
	}
	return cmItems.ConvertServicesResponse(cluster)
}

// GetDeployment get full deployement data for CM server
func (c CMServer) GetDeployment() Deployment {
	var deploymentMap map[string]interface{}
	var uri = fmt.Sprintf("cm/deployment")
	if c.UseGateway {
		curlCommand := c.CreateGatewayCurlGetCommand(uri)
		deploymentMap = ProcessAsMapFromSSHResponse(c.RunGatewayCMCommand(curlCommand, false, true))
	} else {
		request := c.CreateGetRequest(uri)
		deploymentMap = ProcessAsMap(request)
	}
	return ConvertDeploymentResponse(deploymentMap)
}

// ExportClusterTemplate exporting template for a specific cluster
func (c CMServer) ExportClusterTemplate(cluster string) []byte {
	var uri = fmt.Sprintf("clusters/%v/export", cluster)
	if c.UseGateway {
		curlCommand := c.CreateGatewayCurlGetCommand(uri)
		return []byte(c.RunGatewayCMCommand(curlCommand, false, true).StdOut)
	}
	request := c.CreateGetRequest(uri)
	return ProcessRequest(request)
}

// GetUsers returns a list of the user names configured in the system.
func (c CMServer) GetUsers() []User {
	var cmItems CMItems
	var uri = "users"
	if c.UseGateway {
		curlCommand := c.CreateGatewayCurlGetCommand(uri)
		cmItems = ProcessCMItemsFromSSHResponse(c.RunGatewayCMCommand(curlCommand, false, true))
	} else {
		request := c.CreateGetRequest(uri)
		cmItems = ProcessCMItems(request)
	}
	return cmItems.ConvertUsersResponse()
}
