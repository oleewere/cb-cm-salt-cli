package cm

import (
	"fmt"
	"os"
	"strings"
)

// Filter represents filter on agent hosts (by component / service / hosts)
type Filter struct {
	Hosts    []string
	Clusters []string
	Services []string
	Roles    []string
	Server   bool
}

// CreateFilter will make a Filter object from filter strings (hosts)
func CreateFilter(clusterFilter string, serviceFilter string, rolesFilter string, hostFilter string, cmServer bool) Filter {
	filter := Filter{}
	if len(hostFilter) > 0 {
		hosts := strings.Split(hostFilter, ",")
		filter.Hosts = hosts
	}
	if len(clusterFilter) > 0 {
		clusters := strings.Split(clusterFilter, ",")
		filter.Clusters = clusters
	}
	if len(serviceFilter) > 0 {
		services := strings.Split(serviceFilter, ",")
		filter.Services = services
	}
	if len(rolesFilter) > 0 {
		roles := strings.Split(rolesFilter, ",")
		filter.Roles = roles
	}
	filter.Server = cmServer
	return filter
}

// GetFilteredHosts obtain specific hosts based on different filters
func (c CMServer) GetFilteredHosts(filter Filter) map[string]bool {
	finalHosts := make(map[string]bool)
	hosts := make(map[string]bool) // use boolean map as a set
	if filter.Server {
		hosts[c.Hostname] = true
		finalHosts[c.Hostname] = true
	} else if len(filter.Hosts) > 0 {
		cmAgents := c.ListHosts()
		calculateAndFillFinalHosts(cmAgents, filter, hosts, finalHosts)
	} else if len(filter.Clusters) > 0 {
		cmAgents := c.ListHosts()
		for _, cluster := range filter.Clusters {
			inventory := Inventory{}
			inventory = GetInventoryWithHostsForCluster(inventory, cluster, cmAgents)
			for _, host := range inventory.Hosts {
				hosts[host.HostName] = true
			}
			calculateAndFillFinalHosts(cmAgents, filter, hosts, finalHosts)
		}
	} else if len(filter.Roles) > 0 {
		cmAgents := c.ListHosts()
		deployment := c.GetDeployment()
		inventories := CreateInventoriesFromDeploymentsAndHosts(deployment, cmAgents)
		for _, inventory := range inventories {
			if len(filter.Clusters) > 0 && !SliceContains(inventory.ClusterName, filter.Clusters) {
				continue
			}
			rolesHostsMap := inventory.ServiceRoleHostsMap[filter.Services[0]]
			for role, rolesHosts := range rolesHostsMap {
				if SliceContains(role, filter.Roles) {
					for _, host := range rolesHosts {
						hosts[host] = true
					}
				}
			}
		}
		calculateAndFillFinalHosts(cmAgents, filter, hosts, finalHosts)
	} else if len(filter.Services) > 0 {
		cmAgents := c.ListHosts()
		deployment := c.GetDeployment()
		inventories := CreateInventoriesFromDeploymentsAndHosts(deployment, cmAgents)
		for _, inventory := range inventories {
			if len(filter.Clusters) > 0 && !SliceContains(inventory.ClusterName, filter.Clusters) {
				continue
			}
			serviceHostsMap := inventory.ServiceHostsMap
			for service, serviceHosts := range serviceHostsMap {
				if SliceContains(service, filter.Services) {
					for _, host := range serviceHosts {
						hosts[host] = true
					}
				}
			}
		}
		calculateAndFillFinalHosts(cmAgents, filter, hosts, finalHosts)
	} else {
		cmAgents := c.ListHosts()
		calculateAndFillFinalHosts(cmAgents, filter, hosts, finalHosts)
	}
	return finalHosts
}

func calculateAndFillFinalHosts(cmAgents []Host, filter Filter, hosts map[string]bool, finalHosts map[string]bool) {
	for _, cmAgent := range cmAgents {
		if len(filter.Hosts) > 0 {
			filteredHosts := filter.Hosts
			containsHost := false
			for _, filteredHost := range filteredHosts {
				if filteredHost == cmAgent.HostName {
					containsHost = true
				}
			}
			if !containsHost {
				continue
			}
		}
		if len(hosts) > 0 {
			_, ok := hosts[cmAgent.HostName]
			if ok {
				finalHosts[cmAgent.IPAddress] = true
			}
			_, ok = hosts[cmAgent.IPAddress]
			if ok {
				finalHosts[cmAgent.IPAddress] = true
			}

		} else {
			finalHosts[cmAgent.IPAddress] = true
		}
	}
}

func validateFilter(filter Filter) {
	if len(filter.Roles) > 0 && len(filter.Services) != 1 {
		fmt.Println("Error: use exactly 1 service filter with roles filter!")
		os.Exit(1)
	}
}
