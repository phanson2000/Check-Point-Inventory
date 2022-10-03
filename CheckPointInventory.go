package main
import(
    "fmt"
    "strings"
    "regexp"
    "os"
    "strconv"
    "time"
    "encoding/base64"
    "bufio"
    "github.com/howeyc/gopass"
    "flag"
    api "APIFiles"
)

var (
    timeout *int
)

func extractValue(body string, key string) string {
    keystr := "\"" + key + "\":[^,;\\]}]*"
    r, _ := regexp.Compile(keystr)
    match := r.FindString(body)
    keyValMatch := strings.Split(match, ":")
    if len(keyValMatch) == 1 {return "0"}
    return strings.ReplaceAll(keyValMatch[1], "\"", "")
}

func splitstring(strline string, seperator string) (string,string) {
        split := strings.Split(strline,seperator)
        returnvalue1 := (split[0])
        returnvalue0 := (split[1])
return returnvalue0, returnvalue1
}

func createloggingenvirorment() (string,string,string,string,string,string) {
    filenamePrefix := time.Now().Format("20060102150405")
    filenameassetinfo := filenamePrefix + "-assetinfo.csv"
    filenamelicenseinfo := filenamePrefix + "-licenseinfo.csv"
    interfaceInfofilename := filenamePrefix + "-interfaceInfo.csv"
    ifconfigInfofilename := filenamePrefix + "-ifconfigInfo.csv"
    interfacelistfilename := filenamePrefix + "-interfacelist.csv"
    fwverfilename := filenamePrefix + "-fwverlist.csv"
   

    fd, err := os.Create(filenameassetinfo)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    defer fd.Close()
    fd.WriteString("Gateway name,Model,Platform,Serial Number,CPU Number, CPU Frequency, Number of Cores, Hyperthreading,Num of Disks,Disk 1 Model,Disk 1 Capacity,Disk 2 Model,Disk 2 Capacity,Total Disk Size, Total Memory,Memory 1,Memory 2,Memory 3,Memory 4,Memory 5,Memory 6,Memory 7,Memory 8,Number of Line Cards, Card Model 1, Card 1, Card Model 2, Card 2, Card Model 3, Card 3, Card model 4, Card 4, PS Name,PS Status,PS2 Name,PS Status,LOM Status,LOM Firmware \n")
    fd.Close()

    fd, err = os.Create(filenamelicenseinfo)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    defer fd.Close()
    
    fd.WriteString("Gateway name,License \n")
    fd.Close()

    fd, err = os.Create(interfaceInfofilename)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    defer fd.Close()
    
    fd.WriteString("Gateway Name,Interface,State\n")
    fd.Close()

    fd, err = os.Create(ifconfigInfofilename)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    defer fd.Close()
    
    fd.WriteString("Gateway Name,Interface,Mac 1,Mac 2 \n")
    fd.Close()

    fd, err = os.Create(interfacelistfilename)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    defer fd.Close()
    
    fd.WriteString("Gateway Name,Serial,Card,Model,Description \n")
    fd.Close()
    
    fd, err = os.Create(fwverfilename)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    defer fd.Close()
    
    fd.WriteString("Gateway Name,License \n")
    fd.Close()
    return filenameassetinfo, filenamelicenseinfo,interfaceInfofilename,ifconfigInfofilename,interfacelistfilename,fwverfilename
    
}
func getTask(client *api.ApiClient, hash string) string {
    payload := map[string]interface{}{
    "task-id" : hash,
    "details-level" : "full",
    }
    taskresponse,err2 := client.ApiCall("show-task",payload, client.GetSessionID(),false, false)
     if err2 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        return "failed"
    }
    response := taskresponse.GetData()
    taskList:= response["tasks"].([]interface{})
    taskMap := taskList[0].(map[string]interface{})
    //comment:=(taskMap["comments"])
    //status:=(taskMap["status"])
    taskDetails := taskMap["task-details"].([]interface{})
    taskdetailsmap := taskDetails[0].(map[string]interface{})
    responsemsg := taskdetailsmap["responseMessage"]
    responsestr := responsemsg.(string)
    responseDecode, err := base64.URLEncoding.DecodeString(responsestr)
    if err != nil {
        fmt.Printf("Error decoding string: %s ", err.Error())
        return "failed"
    }
    //fmt.Println(string(responseDecode))
    return string(responseDecode)
    
}

func getHostname (client *api.ApiClient, gatewayname string) (string) {
    server := "\""+gatewayname+"\""
    payload := map[string]interface{}{
    "script-name" : "Fetching Gatway Host Name",
    "script" : "/bin/clish -s -c 'show hostname'",
    "targets" : server,
    }
    showScriptTask,err2 := client.ApiCall("run-script",payload, client.GetSessionID(),false, false)
     if err2 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        return ""
    }
    time.Sleep(time.Duration(*timeout) * time.Second)
    task := showScriptTask.GetData()
    taskList:= task["tasks"].([]interface{})
    taskMap := taskList[0].(map[string]interface{})
    taskhash := taskMap["task-id"]
    response:=getTask(client,taskhash.(string))
    return response
}

func getAssetInfo (client *api.ApiClient, gatewayname string) (string) {
    server := "\""+gatewayname+"\""
    payload := map[string]interface{}{
    "script-name" : "Show Asset Info",
    "script" : "/bin/clish -s -c 'show asset all'",
    "targets" : server,
    }
    showScriptTask,err2 := client.ApiCall("run-script",payload, client.GetSessionID(),false, false)
     if err2 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        return ""
    }
    time.Sleep(time.Duration(*timeout) * time.Second)
    task := showScriptTask.GetData()
    taskList:= task["tasks"].([]interface{})
    taskMap := taskList[0].(map[string]interface{})
    taskhash := taskMap["task-id"]
    response:=getTask(client,taskhash.(string))
   return response
}

func getLicenseInfo (client *api.ApiClient, gatewayname string) (string) {
    server := "\""+gatewayname+"\""
    payload := map[string]interface{}{
    "script-name" : "Show License Information",
    "script" : "cplic print",
    "targets" : server,
    }
    showScriptTask,err2 := client.ApiCall("run-script",payload, client.GetSessionID(),false, false)
     if err2 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        return ""
    }
    time.Sleep(time.Duration(*timeout) * time.Second)
    task := showScriptTask.GetData()
    taskList:= task["tasks"].([]interface{})
    taskMap := taskList[0].(map[string]interface{})
    taskhash := taskMap["task-id"]
    response:=getTask(client,taskhash.(string))
   return response
}

func getifconfigInfo (client *api.ApiClient, gatewayname string) (string) {
    server := "\""+gatewayname+"\""
    payload := map[string]interface{}{
    "script-name" : "Show IFCONFIG",
    "script" : "ifconfig -a",
    "targets" : server,
    }
    showScriptTask,err2 := client.ApiCall("run-script",payload, client.GetSessionID(),false, false)
     if err2 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        return ""
    }
    time.Sleep(time.Duration(*timeout) * time.Second)
    task := showScriptTask.GetData()
    taskList:= task["tasks"].([]interface{})
    taskMap := taskList[0].(map[string]interface{})
    taskhash := taskMap["task-id"]
    response:=getTask(client,taskhash.(string))
   return response
}

func getconfigInterfaceInfo (client *api.ApiClient, gatewayname string) (string) {
    server := "\""+gatewayname+"\""
    payload := map[string]interface{}{
    "script-name" : "show configuration interface",
    "script" : "/bin/clish -s -c 'show configuration interface'",
    "targets" : server,
    }
    showScriptTask,err2 := client.ApiCall("run-script",payload, client.GetSessionID(),false, false)
     if err2 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        return ""
    }
    time.Sleep(time.Duration(*timeout) * time.Second)
    task := showScriptTask.GetData()
    taskList:= task["tasks"].([]interface{})
    taskMap := taskList[0].(map[string]interface{})
    taskhash := taskMap["task-id"]
    response:=getTask(client,taskhash.(string))
    return response
}

func getShowInterfaceAllInfo (client *api.ApiClient, gatewayname string) (string) {
    server := "\""+gatewayname+"\""
    payload := map[string]interface{}{
    "script-name" : "show interfaces all",
    "script" : "/bin/clish -s -c 'show interfaces all'",
    "targets" : server,
    }
    showScriptTask,err2 := client.ApiCall("run-script",payload, client.GetSessionID(),false, false)
     if err2 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        return ""
    }
    time.Sleep(time.Duration(*timeout) * time.Second)
    task := showScriptTask.GetData()
    taskList:= task["tasks"].([]interface{})
    taskMap := taskList[0].(map[string]interface{})
    taskhash := taskMap["task-id"]
    response:=getTask(client,taskhash.(string))
    return response
}

func getFirewallVersion (client *api.ApiClient, gatewayname string) (string) {
    server := "\""+gatewayname+"\""
    payload := map[string]interface{}{
    "script-name" : "fw ver -k",
    "script" : "fw ver -k",
    "targets" : server,
    }
    showScriptTask,err2 := client.ApiCall("run-script",payload, client.GetSessionID(),false, false)
     if err2 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        return ""
    }
    time.Sleep(time.Duration(*timeout) * time.Second)
    task := showScriptTask.GetData()
    taskList:= task["tasks"].([]interface{})
    taskMap := taskList[0].(map[string]interface{})
    taskhash := taskMap["task-id"]
    response:=getTask(client,taskhash.(string))
    return response
}


func getDomains(client *api.ApiClient) ([]string) {
    var ipv4address []string
    showDomains,err2 := client.ApiQuery("show-domains", "full", "objects", false, map[string]interface{}{})
    if err2 != nil {
        fmt.Println("Failed to retrieve the domains\n")
        return ipv4address
    }

    for _,sessionObj := range showDomains.GetData(){
        domainname := (sessionObj.(map[string]interface{})["name"].(string))
        payload := map[string]interface{}{
        "name":       domainname,
        }
        showDomainip,err2 := client.ApiCall("show-domain",payload, client.GetSessionID(),false, false)
        if err2 != nil {
            fmt.Println("Failed to retrieve the hosts\n")
            return ipv4address
        }
        domain := showDomainip.GetData()
        serversList:= domain["servers"].([]interface{})
        serversMap := serversList[0].(map[string]interface{})
        ipv4 := serversMap["ipv4-address"]
        ipv4address = append(ipv4address,ipv4.(string))

    }
    return ipv4address
}

func getGatewayList(client *api.ApiClient) ([]string) {
    var gwipv4address []string
    payload := map[string]interface{}{}
    showGateways,err2 := client.ApiQuery("show-gateways-and-servers", "full", "objects", false, payload)
    if err2 != nil {
        fmt.Println("Failed to retrieve the gateways\n")
    return gwipv4address
    }
    showGatewayResponse:=showGateways.GetData()
    for i := 0; i < len(showGatewayResponse); i++ {
        str := strconv.Itoa(i)
        gatewaylist:=showGatewayResponse[str]
        gatewayMap := gatewaylist.(map[string]interface{})
        objectType := gatewayMap["type"]
        gwname := gatewayMap["name"]
        if objectType == "simple-gateway" || objectType == "CpmiClusterMember" || objectType == "Member" {
            gwipv4address = append(gwipv4address,gwname.(string))
        }
    }
    return gwipv4address
}

func writetofile(hostname string, assetinfo string,licenseinfo string,ifconfiginfo string,interfaceInfo string ,interfaceallInfo string ,fwverInfo string){
    filename:=strings.Replace(hostname, "\n", "", -1)
    filename=filename+".txt"
    //fmt.Println(filename)
    f, err := os.Create(filename)
        if err != nil {
            fmt.Println(err)
            return
        }
    f.WriteString("################### ASSETS ########################" + "\n" + assetinfo)
    f.WriteString("################### licenseinfo ########################" + "\n" + licenseinfo)
    f.WriteString("################### ifconfiginfo ########################" + "\n" + ifconfiginfo)
    f.WriteString("################### interfaceInfo ########################" + "\n" + interfaceInfo)
    f.WriteString("################### interfaceallInfo ########################" + "\n" + interfaceallInfo)
    f.WriteString("################### fwverInfo ########################" + "\n" + fwverInfo)
    err = f.Close()
    if err != nil {
        fmt.Println(err)
        return
    }
}
func processassetinfo(filename string,interfacelistfilename string, hostname string, assetinfo string) {
    var interfaceslist string
    var platform string
    var models string
    var serials string
    var cpumodel string
    var cpufrequency string
    var numofcores string
    var hyperthreading string
    var numofdisks string
    var disk1model string 
    var disk1capacity string
    var disk2model string
    var disk2capacity string
    var totaldisksize string
    var totmem,mem1,mem2,mem3,mem4,mem5,mem6,mem7,mem8 string
    var ps1,ps2 string
    var ps1name,ps2name string
    var lominstalled string
    var lomversion string
    var numlinecards string
    var card1, card2, card3,card4 string
    var cardmodel1, cardmodel2,cardmodel3,cardmodel4 string
    scanner := bufio.NewScanner(strings.NewReader(assetinfo))
    for scanner.Scan() {
        if strings.Contains(scanner.Text(), ": ") {
            value2,value3 := splitstring(scanner.Text(),": ")
            value3 = strings.TrimSpace(value3)
            if value3 == "Model" {
                models = value2 }
            if value3 == "Platform" {
                platform = value2 }
           if value3 == "Serial Number" {
                serials = value2  }
           if value3 == "CPU Model" {
                cpumodel = value2 }
           if value3 == "CPU Frequency" {
                cpufrequency = value2 }
           if value3 == "Number of Cores" {
                numofcores = value2 }
           if value3 == "CPU Hyperthreading" {
                hyperthreading = value2 }
           if value3 == "Number of disks"{
                numofdisks = value2 }
           if value3 == "Disk 1 Model" {
                disk1model = value2 }
           if value3 == "Disk 1 Capacity" {
                disk1capacity = value2 }
           if value3 == "Disk 2 Model" {
                disk2capacity = value2 }
           if value3 == "Disk 2 Capacity" {
                disk2capacity = value2 }
           if value3 == "Total Disks size" {
                totaldisksize = value2 }
           if value3 == "Total Memory" {
                totmem = value2 }
           if value3 == "Number of line cards" {
                numlinecards = value2 }
           if value3 == "Line card 1 model" {
                cardmodel1 = value2 }
           if value3 == "Line card 1 type" {
                card1 = value2 }
           if value3 == "Line card 2 model" {
                cardmodel2 = value2 }
           if value3 == "Line card 2 type" {
                card2 = value2 }
           if value3 == "Line card 3 model" {
                cardmodel3 = value2 }
           if value3 == "Line card 3 type" {
                card3 = value2 }
           if value3 == "Line card 4 model" {
                cardmodel4 = value2 }
           if value3 == "Line card 4 type" {
                card4 = value2 }
           if value3 == "Power supply 1 name"  {
                ps1name = value2 }
           if value3 == "Power supply 1 status" {
                ps1 = value2  }
           if value3 == "Power supply 2 name" {
                ps2name = value2 }
           if value3 == "Power supply 2 status" {
                ps2 = value2 }
           if value3 == "LOM Status"  {
                lominstalled = value2 }
           if value3 == "LOM Firmware Revision" {
                lomversion = value2 }
           if value3 == "Memory Slot 1 Size" {
                mem1 = value2 }
           if value3 == "Memory Slot 2 Size" {
                mem2 = value2 }
           if value3 == "Memory Slot 3 Size" {
                mem3 = value2 }
           if value3 == "Memory Slot 4 Size" {
                mem4 = value2 }
           if value3 == "Memory Slot 5 Size" {
                mem5 = value2 }
           if value3 == "Memory Slot 6 Size" {
                mem6 = value2 }
           if value3 == "Memory Slot 7 Size" {
                mem7 = value2 }
           if value3 == "Memory Slot 8 Size" {
                mem8 = value2 }
           if cardmodel1 != "" {interfaceslist = interfaceslist + hostname + "," + serials + "," + "1" + "," + cardmodel1 + "," + card1 + "\n"}
           if cardmodel2 != "" {interfaceslist = interfaceslist + hostname + "," + serials + "," + "2" + "," +   cardmodel2 + "," + card2 + "\n"}
           if cardmodel3 != "" {interfaceslist = interfaceslist + hostname + "," + serials + "," + "3" + "," +   cardmodel3 + "," + card3 + "\n"}
           if cardmodel4 != "" {interfaceslist = interfaceslist + hostname + "," + serials + "," + "4" + "," +   cardmodel4 + "," + card4 + "\n"}
        }
    }
    datacsv := hostname + "," + models + "," + platform + "," + serials + "," + cpumodel + "," + cpufrequency + "," + numofcores + "," +  hyperthreading + "," +  numofdisks + "," +  disk1model + "," +  disk1capacity + "," +  disk2model + "," +  disk2capacity + "," +  totaldisksize + "," +  totmem + "," +  mem1 + "," +  mem2 + "," +  mem3 + "," +  mem4 + "," +  mem5 + "," +  mem6 + "," +  mem7 + "," + mem8 + "," +  numlinecards + "," +  cardmodel1 + "," +  card1 + "," +  cardmodel2 + "," +  card2 + "," +  cardmodel3 + "," +  card3 + "," +  cardmodel4 + "," +  card4 + "," +  ps1name + "," +  ps1 + "," +  ps2name + "," +  ps2 + "," +  lominstalled + "," +  lomversion + " " + "\n"
    

    //fmt.Println(datacsv)
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println(err)
    }
    if _, err := f.Write([]byte(datacsv)); err != nil {
        fmt.Println(err)
    }
    if err := f.Close(); err != nil {
        fmt.Println(err)
    }
    
    f, err = os.OpenFile(interfacelistfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println(err)
    }
    if _, err := f.Write([]byte(interfaceslist)); err != nil {
        fmt.Println(err)
    }
    if err := f.Close(); err != nil {
        fmt.Println(err)
    }
}

func processlicenseinfo(licenseinfofilename string,hostnamestr string,licenseinfo string) {
    licenseinfo=strings.Replace(licenseinfo, "\n", "", -1)
    licenseinfo=licenseinfo + "\n"
    datacsv:=hostnamestr+","+licenseinfo
    //// TODO: ADD NEW PARSER FOR LICENS
    f, err := os.OpenFile(licenseinfofilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println(err)
    }

    if _, err := f.Write([]byte(datacsv)); err != nil {
        fmt.Println(err)
    }
    if err := f.Close(); err != nil {
        fmt.Println(err)
    }
}

func processinterfaceInfo(interfaceInfofilename string,hostnamestr string,interfaceInfo string) {
    var intname4,state,intlinks string
    scanner := bufio.NewScanner(strings.NewReader(interfaceInfo))
    for scanner.Scan() {
        if strings.Contains(scanner.Text(), "set interface") {
             if strings.Contains(scanner.Text(),"state ") {
                tmp := strings.Replace(scanner.Text(), "set interface ","",-1)       
                value4,value5 := splitstring(tmp,"state ")
                intname4  = value5
                state = value4
                intlinks = intlinks + hostnamestr + "," + intname4 + "," + state + "\n"
            }       
        } 
    }
    f, err := os.OpenFile(interfaceInfofilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println(err)
    }

    if _, err := f.Write([]byte(intlinks)); err != nil {
        fmt.Println(err)
    }
    if err := f.Close(); err != nil {
        fmt.Println(err)
    }
}

func processifconfiginfo(ifconfigInfofilename string,hostnamestr string,ifconfiginfo string) {
    var mac1,mac2,interfaces string
    scanner := bufio.NewScanner(strings.NewReader(ifconfiginfo))
    for scanner.Scan() {
        if strings.Contains(scanner.Text(),"   Link encap:Ethernet  HWaddr ") {
            value2,value3 := splitstring(scanner.Text(),"   Link encap:Ethernet  HWaddr ")
            intname := value3
            mac1 = value2
            mac2 = strings.Replace(mac1, ":","",-1)
            mac1 = strings.Replace(mac1, " ","",-1)
            mac2 = strings.Replace(mac2, " ","",-1)
            intname = strings.Replace(intname, " ","",-1)
            interfaces = interfaces + hostnamestr + "," + intname + "," + mac1 + "," + mac2 + "\n"
        } 
    }
    f, err := os.OpenFile(ifconfigInfofilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println(err)
    }

    if _, err := f.Write([]byte(interfaces)); err != nil {
        fmt.Println(err)
    }
    if err := f.Close(); err != nil {
        fmt.Println(err)
    }
}
func processfwverinfo(fwverfilename string,hostnamestr string,fwverInfo string) {
    fwverInfo=strings.Replace(fwverInfo, "\n", "", -1)
    fwverInfo=fwverInfo + "\n"
    datacsv:=hostnamestr+","+fwverInfo
    //// TODO: ADD NEW PARSER FOR LICENS
    f, err := os.OpenFile(fwverfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println(err)
    }

    if _, err := f.Write([]byte(datacsv)); err != nil {
        fmt.Println(err)
    }
    if err := f.Close(); err != nil {
        fmt.Println(err)
    }
}

func main() {
    var uname2,api2 string
    //var timeout *int
    var apiServer,username,password *string
    
    timeout=flag.Int("timeout",2,"Time delay before getting results from gateway")
    apiServer=flag.String("apiserver","0.0.0.0","Check Point Management IP")
    username=flag.String("username","unknown","Domain User to User")
    password=flag.String("password","unknown","DOMAIN PASSWORD")
    flag.Parse()
    //fmt.Println(*flag.Args())
    
    if *apiServer == "0.0.0.0" { 
        fmt.Printf("Enter server IP address or hostname for Check Point Management/MDM: ")
        fmt.Scanf("%s \n",&api2)
        *apiServer=api2
    }
    if *username == "unknown" {
        fmt.Printf("Enter username: ")
        fmt.Scanf("%s \n",&uname2)
        *username=uname2
    }
    if *password=="unknown"{
        fmt.Printf("Enter password: ")
        silentpassword, err := gopass.GetPasswd()
         if err != nil {
            fmt.Println("Login error.\n")
            os.Exit(1)
        }
        *password=string(silentpassword)
    }
    ////  hard coded for testing.  Would not utilize this for security reasons
    //apiServer = "192.168.1.30"
    //username = "admin"
    //password = "vpn123"

    fmt.Println("\n \n ######################################################################")
    args := api.APIClientArgs(api.DefaultPort, "", "", *apiServer, "", -1, "", false, false, "deb.txt", api.WebContext, api.TimeOut, api.SleepTime, "", "")
    client := api.APIClient(args)
    if x, _ := client.CheckFingerprint(); !x {
        print("Could not get the server's fingerprint - Check connectivity with the server.\n")
        os.Exit(1)
    }

    loginRes, err := client.Login(*username, *password, false, "", false, "")
    if err != nil {
        fmt.Println("Login error.\n")
        os.Exit(1)
    }

    if !loginRes.Success {
        fmt.Println("Login failed:\n" + loginRes.ErrorMsg)
        os.Exit(1)
    }
    assetinfofilename,licenseinfofilename,interfaceInfofilename,ifconfigInfofilename,interfacelistfilename,fwverfilename:=createloggingenvirorment()
    domainipv4:=getDomains(client)
    fmt.Println("Management Domain's Found at ", *apiServer, ": " , domainipv4)
    for i := 0; i < len(domainipv4); i++ {
        fmt.Println("\n")
        args2 := api.APIClientArgs(api.DefaultPort, "", "",domainipv4[i], "", -1, "", false, false, "deb.txt", api.WebContext, api.TimeOut, api.SleepTime, "", "")
        client2 := api.APIClient(args2)
            if x, _ := client.CheckFingerprint(); !x {
                print("Could not get the server's fingerprint - Check connectivity with the server.\n")
                os.Exit(1)
            }
        loginRes2, err := client2.Login(*username, *password, false, domainipv4[i], false, "")
            if err != nil {
                fmt.Println("Login error.\n")
                os.Exit(1)
            }
        if !loginRes2.Success {
            fmt.Println("Login failed:\n" + loginRes.ErrorMsg)
            os.Exit(1)
            continue
        }
        gwaddress:=getGatewayList(client2)
        fmt.Println("Domain Address: ",domainipv4[i] , "Gateways on Domain: ", gwaddress)   
            for i := 0; i < len(gwaddress); i++ {
                fmt.Println("     Collecting from Firewall: ", gwaddress[i], "      From Domain: ", domainipv4[i])
                hostname:=getHostname(client2,gwaddress[i])
                hostnamestr:=strings.Replace(hostname, "\n", "", -1)
                assetinfo:=getAssetInfo(client2,gwaddress[i])
                licenseinfo:=getLicenseInfo(client2,gwaddress[i])
                ifconfiginfo:=getifconfigInfo(client2,gwaddress[i])
                interfaceInfo:=getconfigInterfaceInfo(client2,gwaddress[i])
                interfaceallInfo:=getShowInterfaceAllInfo(client2,gwaddress[i])
                fwverInfo:=getFirewallVersion(client2,gwaddress[i])
                processassetinfo(assetinfofilename,interfacelistfilename,hostnamestr,assetinfo)
                processlicenseinfo(licenseinfofilename,hostnamestr,licenseinfo)
                processifconfiginfo(ifconfigInfofilename,hostnamestr,ifconfiginfo)
                processinterfaceInfo(interfaceInfofilename,hostnamestr,interfaceInfo)
                processfwverinfo(fwverfilename,hostnamestr,fwverInfo)
                writetofile(hostname,assetinfo,licenseinfo,ifconfiginfo,interfaceInfo,interfaceallInfo,fwverInfo)
            }
        }
    gwaddress:=getGatewayList(client)
    fmt.Println("\n", "Full Gateways Inventoried: " , gwaddress)
}
