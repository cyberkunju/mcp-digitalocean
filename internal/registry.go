package internal

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/server"

	"mcp-digitalocean/internal/account"
	"mcp-digitalocean/internal/apps"
	"mcp-digitalocean/internal/common"
	"mcp-digitalocean/internal/dbaas"
	"mcp-digitalocean/internal/doks"
	"mcp-digitalocean/internal/droplet"
	"mcp-digitalocean/internal/insights"
	"mcp-digitalocean/internal/marketplace"
	"mcp-digitalocean/internal/networking"
	"mcp-digitalocean/internal/spaces"
)

// supportedServices is a set of services that we support in this MCP server.
var supportedServices = map[string]struct{}{
	"apps":        {},
	"networking":  {},
	"droplets":    {},
	"accounts":    {},
	"spaces":      {},
	"databases":   {},
	"marketplace": {},
	"insights":    {},
	"doks":        {},
}

// registerAppTools registers the app platform tools with the MCP server.
func registerAppTools(s *server.MCPServer, c *godo.Client) error {
	appTools, err := apps.NewAppPlatformTool(c)
	if err != nil {
		return fmt.Errorf("failed to create apps tool: %w", err)
	}

	s.AddTools(appTools.Tools()...)

	return nil
}

// registerCommonTools registers the common tools with the MCP server.
func registerCommonTools(s *server.MCPServer, c *godo.Client) error {
	s.AddTools(common.NewRegionTools(c).Tools()...)

	return nil
}

// registerDropletTools registers the droplet tools with the MCP server.
func registerDropletTools(s *server.MCPServer, c *godo.Client) error {
	s.AddTools(droplet.NewDropletTool(c).Tools()...)
	s.AddTools(droplet.NewDropletActionsTool(c).Tools()...)
	s.AddTools(droplet.NewImagesTool(c).Tools()...)
	s.AddTools(droplet.NewSizesTool(c).Tools()...)
	return nil
}

// registerNetworkingTools registers the networking tools with the MCP server.
func registerNetworkingTools(s *server.MCPServer, c *godo.Client) error {
	s.AddTools(networking.NewCertificateTool(c).Tools()...)
	s.AddTools(networking.NewDomainsTool(c).Tools()...)
	s.AddTools(networking.NewFirewallTool(c).Tools()...)
	s.AddTools(networking.NewReservedIPTool(c).Tools()...)
	// Partner attachments doesn't have much users so this has been disabled
	// s.AddTools(networking.NewPartnerAttachmentTool(c).Tools()...)
	s.AddTools(networking.NewVPCTool(c).Tools()...)
	s.AddTools(networking.NewVPCPeeringTool(c).Tools()...)
	return nil
}

// registerAccountTools registers the account tools with the MCP server.
func registerAccountTools(s *server.MCPServer, c *godo.Client) error {
	s.AddTools(account.NewAccountTools(c).Tools()...)
	s.AddTools(account.NewActionTools(c).Tools()...)
	s.AddTools(account.NewBalanceTools(c).Tools()...)
	s.AddTools(account.NewBillingTools(c).Tools()...)
	s.AddTools(account.NewInvoiceTools(c).Tools()...)
	s.AddTools(account.NewKeysTool(c).Tools()...)

	return nil
}

// registerSpacesTools registers the spaces tools and resources with the MCP server.
func registerSpacesTools(s *server.MCPServer, c *godo.Client) error {
	// Register the tools for spaces keys
	s.AddTools(spaces.NewSpacesKeysTool(c).Tools()...)
	s.AddTools(spaces.NewCDNTool(c).Tools()...)

	return nil
}

// registerMarketplaceTools registers the marketplace tools with the MCP server.
func registerMarketplaceTools(s *server.MCPServer, c *godo.Client) error {
	s.AddTools(marketplace.NewOneClickTool(c).Tools()...)

	return nil
}

func registerInsightsTools(s *server.MCPServer, c *godo.Client) error {
	s.AddTools(insights.NewUptimeTool(c).Tools()...)
	s.AddTools(insights.NewUptimeCheckAlertTool(c).Tools()...)
	s.AddTools(insights.NewAlertPolicyTool(c).Tools()...)
	return nil
}

func registerDOKSTools(s *server.MCPServer, c *godo.Client) error {
	s.AddTools(doks.NewDoksTool(c).Tools()...)

	return nil
}

func registerDatabasesTools(s *server.MCPServer, c *godo.Client) error {
	s.AddTools(dbaas.NewClusterTool(c).Tools()...)
	s.AddTools(dbaas.NewFirewallTool(c).Tools()...)
	s.AddTools(dbaas.NewKafkaTool(c).Tools()...)
	s.AddTools(dbaas.NewMongoTool(c).Tools()...)
	s.AddTools(dbaas.NewMysqlTool(c).Tools()...)
	s.AddTools(dbaas.NewOpenSearchTool(c).Tools()...)
	s.AddTools(dbaas.NewPostgreSQLTool(c).Tools()...)
	s.AddTools(dbaas.NewRedisTool(c).Tools()...)
	s.AddTools(dbaas.NewUserTool(c).Tools()...)

	return nil
}

// Register registers the set of tools for the specified services with the MCP server.
// We either register a subset of tools of the services are specified, or we register all tools if no services are specified.
func Register(logger *slog.Logger, s *server.MCPServer, c *godo.Client, servicesToActivate ...string) error {
	if len(servicesToActivate) == 0 {
		logger.Warn("no services specified, loading all supported services")
		services := make([]string, 0, len(supportedServices))
		for k := range supportedServices {
			services = append(services, k)
		}
		servicesToActivate = services
	}
	for _, svc := range servicesToActivate {
		logger.Debug(fmt.Sprintf("Registering tool and resources for service: %s", svc))
		switch svc {
		case "apps":
			if err := registerAppTools(s, c); err != nil {
				return fmt.Errorf("failed to register app tools: %w", err)
			}
		case "networking":
			if err := registerNetworkingTools(s, c); err != nil {
				return fmt.Errorf("failed to register networking tools: %w", err)
			}
		case "droplets":
			if err := registerDropletTools(s, c); err != nil {
				return fmt.Errorf("failed to register droplets tool: %w", err)
			}
		case "accounts":
			if err := registerAccountTools(s, c); err != nil {
				return fmt.Errorf("failed to register account tools: %w", err)
			}
		case "spaces":
			if err := registerSpacesTools(s, c); err != nil {
				return fmt.Errorf("failed to register spaces tools: %w", err)
			}
		case "databases":
			if err := registerDatabasesTools(s, c); err != nil {
				return fmt.Errorf("failed to register databases tools: %w", err)
			}
		case "marketplace":
			if err := registerMarketplaceTools(s, c); err != nil {
				return fmt.Errorf("failed to register marketplace tools: %w", err)
			}
		case "insights":
			if err := registerInsightsTools(s, c); err != nil {
				return fmt.Errorf("failed to register insights tools: %w", err)
			}
		case "doks":
			if err := registerDOKSTools(s, c); err != nil {
				return fmt.Errorf("failed to register DOKS tools: %w", err)
			}
		default:
			return fmt.Errorf("unsupported service: %s, supported service are: %v", svc, setToString(supportedServices))
		}
	}

	// Common tools are always registered because they provide common functionality for all services such as region resources
	if err := registerCommonTools(s, c); err != nil {
		return fmt.Errorf("failed to register common tools: %w", err)
	}

	return nil
}

func setToString(set map[string]struct{}) string {
	var result []string
	for key := range set {
		result = append(result, key)
	}

	return strings.Join(result, ",")
}
