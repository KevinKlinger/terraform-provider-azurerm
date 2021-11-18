package clients

import (
	"context"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/common"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/features"
	advisor "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/advisor/client"
	analysisServices "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/analysisservices/client"
	apiManagement "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/apimanagement/client"
	appConfiguration "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/appconfiguration/client"
	applicationInsights "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/applicationinsights/client"
	appService "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/appservice/client"
	attestation "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/attestation/client"
	authorization "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/authorization/client"
	automation "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/automation/client"
	azureStackHCI "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/azurestackhci/client"
	batch "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/batch/client"
	blueprints "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/blueprints/client"
	bot "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/bot/client"
	cdn "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/cdn/client"
	cognitiveServices "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/cognitive/client"
	communication "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/communication/client"
	compute "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/compute/client"
	consumption "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/consumption/client"
	containerServices "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/containers/client"
	cosmosdb "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/cosmos/client"
	costmanagement "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/costmanagement/client"
	customproviders "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/customproviders/client"
	datamigration "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/databasemigration/client"
	databoxedge "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/databoxedge/client"
	databricks "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/databricks/client"
	datafactory "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/datafactory/client"
	datalake "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/datalake/client"
	dataprotection "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/dataprotection/client"
	datashare "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/datashare/client"
	desktopvirtualization "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/desktopvirtualization/client"
	devspace "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/devspace/client"
	devtestlabs "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/devtestlabs/client"
	digitaltwins "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/digitaltwins/client"
	dns "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/dns/client"
	domainservices "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/domainservices/client"
	eventgrid "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/eventgrid/client"
	eventhub "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/eventhub/client"
	firewall "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/firewall/client"
	frontdoor "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/frontdoor/client"
	hdinsight "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/hdinsight/client"
	healthcare "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/healthcare/client"
	hpccache "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/hpccache/client"
	hsm "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/hsm/client"
	iotcentral "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/iotcentral/client"
	iothub "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/iothub/client"
	timeseriesinsights "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/iottimeseriesinsights/client"
	keyvault "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/keyvault/client"
	kusto "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/kusto/client"
	lighthouse "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/lighthouse/client"
	loadbalancers "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/loadbalancer/client"
	loganalytics "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/loganalytics/client"
	logic "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/logic/client"
	logz "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/logz/client"
	machinelearning "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/machinelearning/client"
	maintenance "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/maintenance/client"
	managedapplication "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/managedapplications/client"
	managementgroup "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/managementgroup/client"
	maps "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/maps/client"
	mariadb "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/mariadb/client"
	media "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/media/client"
	mixedreality "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/mixedreality/client"
	monitor "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/monitor/client"
	msi "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/msi/client"
	mssql "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/mssql/client"
	mysql "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/mysql/client"
	netapp "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/netapp/client"
	network "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/network/client"
	notificationhub "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/notificationhub/client"
	policy "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/policy/client"
	portal "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/portal/client"
	postgres "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/postgres/client"
	powerBI "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/powerbi/client"
	privatedns "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/privatedns/client"
	purview "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/purview/client"
	recoveryServices "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/recoveryservices/client"
	redis "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/redis/client"
	redisenterprise "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/redisenterprise/client"
	relay "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/relay/client"
	resource "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/resource/client"
	search "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/search/client"
	securityCenter "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/securitycenter/client"
	sentinel "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/sentinel/client"
	serviceBus "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/servicebus/client"
	serviceFabric "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/servicefabric/client"
	serviceFabricMesh "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/servicefabricmesh/client"
	signalr "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/signalr/client"
	appPlatform "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/springcloud/client"
	sql "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/sql/client"
	storage "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/storage/client"
	streamAnalytics "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/streamanalytics/client"
	subscription "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/subscription/client"
	synapse "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/synapse/client"
	trafficManager "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/trafficmanager/client"
	videoAnalyzer "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/videoanalyzer/client"
	vmware "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/vmware/client"
	web "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/web/client"
)

type Client struct {
	// StopContext is used for propagating control from Terraform Core (e.g. Ctrl/Cmd+C)
	StopContext context.Context

	Account  *ResourceManagerAccount
	Features features.UserFeatures

	Advisor               *advisor.Client
	AnalysisServices      *analysisServices.Client
	ApiManagement         *apiManagement.Client
	AppConfiguration      *appConfiguration.Client
	AppInsights           *applicationInsights.Client
	AppPlatform           *appPlatform.Client
	AppService            *appService.Client
	Attestation           *attestation.Client
	Authorization         *authorization.Client
	Automation            *automation.Client
	AzureStackHCI         *azureStackHCI.Client
	Batch                 *batch.Client
	Blueprints            *blueprints.Client
	Bot                   *bot.Client
	Cdn                   *cdn.Client
	Cognitive             *cognitiveServices.Client
	Communication         *communication.Client
	Compute               *compute.Client
	Consumption           *consumption.Client
	Containers            *containerServices.Client
	Cosmos                *cosmosdb.Client
	CostManagement        *costmanagement.Client
	CustomProviders       *customproviders.Client
	DatabaseMigration     *datamigration.Client
	DataBricks            *databricks.Client
	DataboxEdge           *databoxedge.Client
	DataFactory           *datafactory.Client
	Datalake              *datalake.Client
	DataProtection        *dataprotection.Client
	DataShare             *datashare.Client
	DesktopVirtualization *desktopvirtualization.Client
	DevSpace              *devspace.Client
	DevTestLabs           *devtestlabs.Client
	DigitalTwins          *digitaltwins.Client
	Dns                   *dns.Client
	DomainServices        *domainservices.Client
	EventGrid             *eventgrid.Client
	Eventhub              *eventhub.Client
	Firewall              *firewall.Client
	Frontdoor             *frontdoor.Client
	HPCCache              *hpccache.Client
	HSM                   *hsm.Client
	HDInsight             *hdinsight.Client
	HealthCare            *healthcare.Client
	IoTCentral            *iotcentral.Client
	IoTHub                *iothub.Client
	IoTTimeSeriesInsights *timeseriesinsights.Client
	KeyVault              *keyvault.Client
	Kusto                 *kusto.Client
	Lighthouse            *lighthouse.Client
	LoadBalancers         *loadbalancers.Client
	LogAnalytics          *loganalytics.Client
	Logic                 *logic.Client
	Logz                  *logz.Client
	MachineLearning       *machinelearning.Client
	Maintenance           *maintenance.Client
	ManagedApplication    *managedapplication.Client
	ManagementGroups      *managementgroup.Client
	Maps                  *maps.Client
	MariaDB               *mariadb.Client
	Media                 *media.Client
	MixedReality          *mixedreality.Client
	Monitor               *monitor.Client
	MSI                   *msi.Client
	MSSQL                 *mssql.Client
	MySQL                 *mysql.Client
	NetApp                *netapp.Client
	Network               *network.Client
	NotificationHubs      *notificationhub.Client
	Policy                *policy.Client
	Portal                *portal.Client
	Postgres              *postgres.Client
	PowerBI               *powerBI.Client
	PrivateDns            *privatedns.Client
	Purview               *purview.Client
	RecoveryServices      *recoveryServices.Client
	Redis                 *redis.Client
	RedisEnterprise       *redisenterprise.Client
	Relay                 *relay.Client
	Resource              *resource.Client
	Search                *search.Client
	SecurityCenter        *securityCenter.Client
	Sentinel              *sentinel.Client
	ServiceBus            *serviceBus.Client
	ServiceFabric         *serviceFabric.Client
	ServiceFabricMesh     *serviceFabricMesh.Client
	SignalR               *signalr.Client
	Storage               *storage.Client
	StreamAnalytics       *streamAnalytics.Client
	Subscription          *subscription.Client
	Sql                   *sql.Client
	Synapse               *synapse.Client
	TrafficManager        *trafficManager.Client
	VideoAnalyzer         *videoAnalyzer.Client
	Vmware                *vmware.Client
	Web                   *web.Client
}

// NOTE: it should be possible for this method to become Private once the top level Client's removed

func (client *Client) Build(ctx context.Context, o *common.ClientOptions) error {
	autorest.Count429AsRetry = false
	// Disable the Azure SDK for Go's validation since it's unhelpful for our use-case
	validation.Disabled = true

	client.Features = o.Features
	client.StopContext = ctx

	client.Advisor = advisor.NewClient(o)
	client.AnalysisServices = analysisServices.NewClient(o)
	client.ApiManagement = apiManagement.NewClient(o)
	client.AppConfiguration = appConfiguration.NewClient(o)
	client.AppInsights = applicationInsights.NewClient(o)
	client.AppPlatform = appPlatform.NewClient(o)
	client.AppService = appService.NewClient(o)
	client.Attestation = attestation.NewClient(o)
	client.Authorization = authorization.NewClient(o)
	client.Automation = automation.NewClient(o)
	client.AzureStackHCI = azureStackHCI.NewClient(o)
	client.Batch = batch.NewClient(o)
	client.Blueprints = blueprints.NewClient(o)
	client.Bot = bot.NewClient(o)
	client.Cdn = cdn.NewClient(o)
	client.Cognitive = cognitiveServices.NewClient(o)
	client.Communication = communication.NewClient(o)
	client.Compute = compute.NewClient(o)
	client.Consumption = consumption.NewClient(o)
	client.Containers = containerServices.NewClient(o)
	client.Cosmos = cosmosdb.NewClient(o)
	client.CostManagement = costmanagement.NewClient(o)
	client.CustomProviders = customproviders.NewClient(o)
	client.DatabaseMigration = datamigration.NewClient(o)
	client.DataBricks = databricks.NewClient(o)
	client.DataboxEdge = databoxedge.NewClient(o)
	client.DataFactory = datafactory.NewClient(o)
	client.Datalake = datalake.NewClient(o)
	client.DataProtection = dataprotection.NewClient(o)
	client.DataShare = datashare.NewClient(o)
	client.DesktopVirtualization = desktopvirtualization.NewClient(o)
	client.DevSpace = devspace.NewClient(o)
	client.DevTestLabs = devtestlabs.NewClient(o)
	client.DigitalTwins = digitaltwins.NewClient(o)
	client.Dns = dns.NewClient(o)
	client.DomainServices = domainservices.NewClient(o)
	client.EventGrid = eventgrid.NewClient(o)
	client.Eventhub = eventhub.NewClient(o)
	client.Firewall = firewall.NewClient(o)
	client.Frontdoor = frontdoor.NewClient(o)
	client.HPCCache = hpccache.NewClient(o)
	client.HSM = hsm.NewClient(o)
	client.HDInsight = hdinsight.NewClient(o)
	client.HealthCare = healthcare.NewClient(o)
	client.IoTCentral = iotcentral.NewClient(o)
	client.IoTHub = iothub.NewClient(o)
	client.IoTTimeSeriesInsights = timeseriesinsights.NewClient(o)
	client.KeyVault = keyvault.NewClient(o)
	client.Kusto = kusto.NewClient(o)
	client.Lighthouse = lighthouse.NewClient(o)
	client.LogAnalytics = loganalytics.NewClient(o)
	client.LoadBalancers = loadbalancers.NewClient(o)
	client.Logic = logic.NewClient(o)
	client.Logz = logz.NewClient(o)
	client.MachineLearning = machinelearning.NewClient(o)
	client.Maintenance = maintenance.NewClient(o)
	client.ManagedApplication = managedapplication.NewClient(o)
	client.ManagementGroups = managementgroup.NewClient(o)
	client.Maps = maps.NewClient(o)
	client.MariaDB = mariadb.NewClient(o)
	client.Media = media.NewClient(o)
	client.MixedReality = mixedreality.NewClient(o)
	client.Monitor = monitor.NewClient(o)
	client.MSI = msi.NewClient(o)
	client.MSSQL = mssql.NewClient(o)
	client.MySQL = mysql.NewClient(o)
	client.NetApp = netapp.NewClient(o)
	client.Network = network.NewClient(o)
	client.NotificationHubs = notificationhub.NewClient(o)
	client.Policy = policy.NewClient(o)
	client.Portal = portal.NewClient(o)
	client.Postgres = postgres.NewClient(o)
	client.PowerBI = powerBI.NewClient(o)
	client.PrivateDns = privatedns.NewClient(o)
	client.Purview = purview.NewClient(o)
	client.RecoveryServices = recoveryServices.NewClient(o)
	client.Redis = redis.NewClient(o)
	client.RedisEnterprise = redisenterprise.NewClient(o)
	client.Relay = relay.NewClient(o)
	client.Resource = resource.NewClient(o)
	client.Search = search.NewClient(o)
	client.SecurityCenter = securityCenter.NewClient(o)
	client.Sentinel = sentinel.NewClient(o)
	client.ServiceBus = serviceBus.NewClient(o)
	client.ServiceFabric = serviceFabric.NewClient(o)
	client.ServiceFabricMesh = serviceFabricMesh.NewClient(o)
	client.SignalR = signalr.NewClient(o)
	client.Sql = sql.NewClient(o)
	client.Storage = storage.NewClient(o)
	client.StreamAnalytics = streamAnalytics.NewClient(o)
	client.Subscription = subscription.NewClient(o)
	client.Synapse = synapse.NewClient(o)
	client.TrafficManager = trafficManager.NewClient(o)
	client.VideoAnalyzer = videoAnalyzer.NewClient(o)
	client.Vmware = vmware.NewClient(o)
	client.Web = web.NewClient(o)

	return nil
}
