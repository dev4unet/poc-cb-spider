package main

import (
	"fmt"
	azdrv "github.com/cloud-barista/poc-cb-spider/cloud-driver/drivers/azure"
	idrv "github.com/cloud-barista/poc-cb-spider/cloud-driver/interfaces"
	irs "github.com/cloud-barista/poc-cb-spider/cloud-driver/interfaces/resources"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

func main() {
	testCreateVM()
}

func testCreateVM() {

	// 리소스 핸들러 로드
	var cloudDriver idrv.CloudDriver
	cloudDriver = new(azdrv.AzureDriver)

	config := readConfigFile()
	connectionInfo := idrv.ConnectionInfo{
		CredentialInfo: idrv.CredentialInfo{
			ClientId:       config.Azure.ClientId,
			ClientSecret:   config.Azure.ClientSecret,
			TenantId:       config.Azure.TenantId,
			SubscriptionId: config.Azure.SubscriptionID,
		},
		RegionInfo: idrv.RegionInfo{
			Region:        config.Azure.Location,
			ResourceGroup: config.Azure.GroupName,
		},
	}

	cloudConnection, _ := cloudDriver.ConnectCloud(connectionInfo)
	
	//imageHandler, _ := cloudConnection.CreateImageHandler()
	vNetworkHandler, _ := cloudConnection.CreateVNetworkHandler()
	securityHandler, _ := cloudConnection.CreateSecurityHandler()
	publicIPHandler, _ := cloudConnection.CreatePublicIPHandler()
	vNicHandler, _ := cloudConnection.CreateVNicHandler()
	vmHandler, _ := cloudConnection.CreateVMHandler()
	
	// 1. Virtual Network 생성
	vNetworkId := config.Azure.VNetwork.GroupName + ":" + config.Azure.VNetwork.Name
	fmt.Println("Start CreateVNetwork() ...")
	vNetReqInfo := irs.VNetworkReqInfo{Id: vNetworkId}
	_, err := vNetworkHandler.CreateVNetwork(vNetReqInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Finish CreateVNetwork()")
	
	// 2. Security Group 생성
	securityGroupId := config.Azure.Security.GroupName + ":" + config.Azure.Security.Name
	fmt.Println("Start CreateSecurity() ...")
	secReqInfo := irs.SecurityReqInfo{Id: securityGroupId}
	_, err = securityHandler.CreateSecurity(secReqInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Finish CreateSecurity()")
	
	// 3. Public IP 생성
	publicIPId := config.Azure.PublicIP.GroupName + ":" + config.Azure.PublicIP.Name
	fmt.Println("Start CreatePublicIP() ...")
	publicIPReqInfo := irs.PublicIPReqInfo{Id: publicIPId}
	_, err = publicIPHandler.CreatePublicIP(publicIPReqInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Finish CreatePublicIP()")
	
	// 4. Virtual Network Interface 생성
	vNicId := config.Azure.VNic.GroupName + ":" + config.Azure.VNic.Name
	fmt.Println("Start CreateVNic() ...")
	vNicReqInfo := irs.VNicReqInfo{Id: vNicId}
	_, err = vNicHandler.CreateVNic(vNicReqInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Finish CreateVNic()")
	
	// 5. VM 생성
	fmt.Println("Start Create VM ...")
	vmName := config.Azure.GroupName + ":" + config.Azure.VMName
	imageId := config.Azure.Image.Publisher + ":" + config.Azure.Image.Offer + ":" + config.Azure.Image.Sku + ":" + config.Azure.Image.Version
	vmReqInfo := irs.VMReqInfo{
		Name: vmName,
		ImageInfo: irs.ImageInfo{
			Id: imageId,
		},
		SpecID: config.Azure.VMSize,
		VNetworkInfo: irs.VNetworkInfo{
			Id: config.Azure.Nic.ID,
		},
		LoginInfo: irs.LoginInfo{
			AdminUsername: config.Azure.Os.AdminUsername,
			//AdminPassword: config.Azure.Os.AdminPassword,
		},
	}
	
	vm, err := vmHandler.StartVM(vmReqInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Finish Create VM")
	
	spew.Dump(vm)
}

func cleanResource() {

}

type Config struct {
	Azure struct {
		ClientId       string `yaml:"client_id"`
		ClientSecret   string `yaml:"client_secret"`
		TenantId       string `yaml:"tenant_id"`
		SubscriptionID string `yaml:"subscription_id"`

		GroupName string `yaml:"group_name"`
		VMName    string `yaml:"vm_name"`

		Location string `yaml:"location"`
		VMSize   string `yaml:"vm_size"`
		Image    struct {
			Publisher string `yaml:"publisher"`
			Offer     string `yaml:"offer"`
			Sku       string `yaml:"sku"`
			Version   string `yaml:"version"`
		} `yaml:"image"`
		Os struct {
			ComputeName   string `yaml:"compute_name"`
			AdminUsername string `yaml:"admin_username"`
			AdminPassword string `yaml:"admin_password"`
		} `yaml:"os"`
		Nic struct {
			ID string `yaml:"id"`
		} `yaml:"nic"`
		ServerId string `yaml:"server_id"`

		ImageInfo struct {
			GroupName string `yaml:"group_name"`
			Name      string `yaml:"name"`
		} `yaml:"image_info"`

		PublicIP struct {
			GroupName string `yaml:"group_name"`
			Name      string `yaml:"name"`
		} `yaml:"public_ip"`

		Security struct {
			GroupName string `yaml:"group_name"`
			Name      string `yaml:"name"`
		} `yaml:"security_group"`

		VNetwork struct {
			GroupName string `yaml:"group_name"`
			Name      string `yaml:"name"`
		} `yaml:"virtual_network"`

		VNic struct {
			GroupName string `yaml:"group_name"`
			Name      string `yaml:"name"`
		} `yaml:"network_interface"`
	} `yaml:"azure"`
}

func readConfigFile() Config {
	// Set Environment Value of Project Root Path
	rootPath := os.Getenv("CBSPIDER_PATH")
	data, err := ioutil.ReadFile(rootPath + "/config/config.yaml")
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return config
}
