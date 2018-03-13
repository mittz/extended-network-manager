package rancher

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mittz/extended-network-manager/config"
	"github.com/mittz/extended-network-manager/controller"
	"github.com/mittz/extended-network-manager/provider"
	rancherclient "github.com/rancher/go-rancher/client"
)

const (
	LabelExtendedNetwork = "io.rancher.container.network.extended"
)

func init() {
	enmc, err := NewExtendedNetworkManagerController()

	if err != nil {
		log.Fatalf("%v", err)
	}

	controller.RegisterController(enmc.GetName(), enmc)
}

// Event is a structure used in subscribe
type Event struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ResourceID   string `json:"resourceId"`
	ResourceType string `json:"resourceType"`
	Data         Data   `json:"data"`
}

// Data is a structure used in Event
type Data struct {
	Resource Resource `json:"resource"`
}

// Resource is a structure used in Data
type Resource struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	State      string                 `json:"state"`
	HostID     string                 `json:"hostId"`
	ExternalID string                 `json:"externalId"`
	Labels     map[string]interface{} `json:"labels"`
}

func (*ExntendedNetworkManagerController) Init(message string) {
	log.Println(message)
}

type ExntendedNetworkManagerController struct {
}

func (*ExntendedNetworkManagerController) GetName() string {
	return "rancher"
}

func NewExtendedNetworkManagerController() (*ExntendedNetworkManagerController, error) {
	enmc := &ExntendedNetworkManagerController{}
	return enmc, nil
}

func CompareMaps(a map[string]interface{}, b map[string]interface{}) bool {
	for key := range b {
		if a[key] != b[key] {
			return false
		}
	}
	return true
}

func IsFilteredEvent(event Event, resourceType string, state string, labels map[string]interface{}) bool {
	return (event.ResourceType == resourceType && event.Data.Resource.State == state && CompareMaps(event.Data.Resource.Labels, labels))
}

func GetContainersOfHost(hostID string) []Resource {
	opts := rancherclient.ClientOpts{}
	opts.Url = os.Getenv("CATTLE_URL")
	opts.AccessKey = os.Getenv("CATTLE_ACCESS_KEY")
	opts.SecretKey = os.Getenv("CATTLE_SECRET_KEY")
	rancherClient, _ := rancherclient.NewRancherClient(&opts)
	l := rancherclient.NewListOpts()
	c, _ := rancherClient.Container.List(l)

	resources := []Resource{}

	for _, v := range c.Data {
		if v.HostId == hostID {
			resources = append(resources, Resource{v.Id, v.Name, v.State, v.HostId, v.ExternalId, v.Labels})
		}
	}

	return resources
}

func SearchProviderID(hostID string) string {
	containers := GetContainersOfHost(hostID)
	var provider string
	providerStackServiceName := "extended-network-manager/provider"

	for _, c := range containers {
		v, contain := c.Labels["io.rancher.stack_service.name"]
		if contain && v == providerStackServiceName {
			provider = c.ID
		}
	}

	return provider
}

func (*ExntendedNetworkManagerController) Run(provider provider.ENMProvider) {
	accessKey := os.Getenv("CATTLE_ACCESS_KEY")
	secretKey := os.Getenv("CATTLE_SECRET_KEY")

	u, err := url.Parse(os.Getenv("CATTLE_URL"))
	host := u.Host

	if err != nil {
		log.Fatal(err)
	}

	resourceType := "container"
	state := "running"
	labels := map[string]interface{}{}

	labels[LabelExtendedNetwork] = "true"

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	h := http.Header{}
	h.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(accessKey+":"+secretKey)))

	url := "ws://" + host + "/v1/subscribe?eventNames=resource.change"
	c, _, err := websocket.DefaultDialer.Dial(url, h)
	if err != nil {
		log.Fatal("", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var event Event

			if err := json.Unmarshal(message, &event); err != nil {
				log.Fatal("JSON Unmarshal error:", err)
				return
			}

			if IsFilteredEvent(event, resourceType, state, labels) {
				log.Println(string(message))
				if event.Data.Resource.Labels[LabelExtendedNetwork] == "true" {
					var extNetwork string
					for i := 0; ; i++ {
						extNetwork = fmt.Sprintf("%s.%d", LabelExtendedNetwork, i)
						var cnc config.ContainerNetworkConfig

						if event.Data.Resource.Labels[extNetwork+".host.interface"] == nil {
							break
						}

						if event.Data.Resource.Labels[extNetwork+".interface"] == nil {
							break
						}

						if event.Data.Resource.Labels[extNetwork+".ipaddress"] == nil {
							break
						}

						if event.Data.Resource.Labels[extNetwork+".macaddress"] == nil {
							break
						}

						cnc.HostInterface = event.Data.Resource.Labels[extNetwork+".host.interface"].(string)
						cnc.Interface = event.Data.Resource.Labels[extNetwork+".interface"].(string)
						cnc.ID = event.Data.Resource.ExternalID
						cnc.IPAddress = event.Data.Resource.Labels[extNetwork+".ipaddress"].(string)
						cnc.MACAddress = event.Data.Resource.Labels[extNetwork+".macaddress"].(string)

						err := provider.ApplyConfig(SearchProviderID(event.Data.Resource.HostID))

						if err != nil {
							log.Printf("%v\n", err)
						}

						provider.AddInterface(cnc)
					}
				}
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
