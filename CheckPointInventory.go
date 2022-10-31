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
    timeout,defaultport *int
    domainipv4 []string
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

func createloggingenvirorment() (string,string,string,string,string,string,string) {
    filenamePrefix := time.Now().Format("20060102150405")
    filenameassetinfo := filenamePrefix + "-assetinfo.csv"
    filenamelicenseinfo := filenamePrefix + "-licenseinfo.csv"
    interfaceInfofilename := filenamePrefix + "-interfaceInfo.csv"
    ifconfigInfofilename := filenamePrefix + "-ifconfigInfo.csv"
    interfacelistfilename := filenamePrefix + "-interfacelist.csv"
    fwverfilename := filenamePrefix + "-fwverlist.csv"
    datacollectionfilename := filenamePrefix + "-datacollection.csv"
   

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

    fd, err = os.Create(datacollectionfilename)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    defer fd.Close()
    
    fd.WriteString("Gateway Name,Patches and Versions, SecureXL, Affinity, MQ,ShowConfig,PriQ,Storage\n")
    fd.Close()


    return filenameassetinfo, filenamelicenseinfo,interfaceInfofilename,ifconfigInfofilename,interfacelistfilename,fwverfilename,datacollectionfilename
    
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
        return "OFFLINE"
    }
    time.Sleep(time.Duration(*timeout) * time.Second)
    task := showScriptTask.GetData()
        if task["tasks"] !=nil {
        taskList:= task["tasks"].([]interface{})
        taskMap := taskList[0].(map[string]interface{})
        taskhash := taskMap["task-id"]
        response:=getTask(client,taskhash.(string))
        return response
    }
    return "OFFLINE"
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

func GetFirewallConfiguration (client *api.ApiClient, gatewayname string) (string,string,string,string,string,string,string,string) {
    server := "\""+gatewayname+"\""
    payload1 := map[string]interface{}{
    "script-name" : "cpinfo -y all",
    "script" : "cpinfo -y all",
    "targets" : server,
    }
    showScriptTask,err2 := client.ApiCall("run-script",payload1, client.GetSessionID(),false, false)
     if err2 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        return "","","","","","","",""
    }
    time.Sleep(time.Duration(*timeout) * time.Second)
    task := showScriptTask.GetData()
    taskList:= task["tasks"].([]interface{})
    taskMap := taskList[0].(map[string]interface{})
    taskhash := taskMap["task-id"]
    response1:=getTask(client,taskhash.(string))

    
    payload2 := map[string]interface{}{
    "script-name" : "show routed version",
    "script" : "/bin/clish -s -c 'show routed version'",
    "targets" : server,
    }

    showScriptTask,err3 := client.ApiCall("run-script",payload2, client.GetSessionID(),false, false)
     if err3 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        return "","","","","","","",""
    }
    time.Sleep(time.Duration(*timeout) * time.Second)
    task2 := showScriptTask.GetData()
    taskList2:= task2["tasks"].([]interface{})
    taskMap2 := taskList2[0].(map[string]interface{})
    taskhash2 := taskMap2["task-id"]
    response2:=getTask(client,taskhash2.(string))

    payload3 := map[string]interface{}{
    "script-name" : "fw ctl multik stat",
    "script" : "fw ctl multik stat",
    "targets" : server,
    }

    showScriptTask,err4 := client.ApiCall("run-script",payload3, client.GetSessionID(),false, false)
     if err4 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        return "","","","","","","",""
    }
    time.Sleep(time.Duration(*timeout) * time.Second)
    task3 := showScriptTask.GetData()
    taskList3:= task3["tasks"].([]interface{})
    taskMap3 := taskList3[0].(map[string]interface{})
    taskhash3 := taskMap3["task-id"]
    response3:=getTask(client,taskhash3.(string))

    payload4 := map[string]interface{}{
    "script-name" : "sim affinity –l",
    "script" : "sim affinity –l",
    "targets" : server,
    }

    showScriptTask,err5 := client.ApiCall("run-script",payload4, client.GetSessionID(),false, false)
     if err5 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        return "","","","","","","",""
    }
    time.Sleep(time.Duration(*timeout) * time.Second)
    task4 := showScriptTask.GetData()
    taskList4:= task4["tasks"].([]interface{})
    taskMap4 := taskList4[0].(map[string]interface{})
    taskhash4 := taskMap4["task-id"]
    response4:=getTask(client,taskhash4.(string))

    payload5 := map[string]interface{}{
    "script-name" : "mq_mng --show",
    "script" : "mq_mng --show",
    "targets" : server,
    }
    showScriptTask,err6 := client.ApiCall("run-script",payload5, client.GetSessionID(),false, false)
     if err6 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        return "","","","","","","",""
    }
    time.Sleep(time.Duration(*timeout) * time.Second)
    task5 := showScriptTask.GetData()
    taskList5:= task5["tasks"].([]interface{})
    taskMap5 := taskList5[0].(map[string]interface{})
    taskhash5 := taskMap5["task-id"]
    response5:=getTask(client,taskhash5.(string))

    payload6 := map[string]interface{}{
    "script-name" : "/bin/clish -s -c 'show configuration'",
    "script" : "/bin/clish -s -c 'show configuration'",
    "targets" : server,
    }

    showScriptTask,err7 := client.ApiCall("run-script",payload6, client.GetSessionID(),false, false)
     if err7 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        return "","","","","","","",""
    }
    time.Sleep(time.Duration(*timeout) * time.Second)
    task6 := showScriptTask.GetData()
    taskList6:= task6["tasks"].([]interface{})
    taskMap6 := taskList6[0].(map[string]interface{})
    taskhash6 := taskMap6["task-id"]
    response6:=getTask(client,taskhash6.(string))

    payload7 := map[string]interface{}{
    "script-name" : "fw ctl get int fwmultik_sync_processing_enabled",
    "script" : "fw ctl get int fwmultik_sync_processing_enabled",
    "targets" : server,
    }

    showScriptTask,err8 := client.ApiCall("run-script",payload7, client.GetSessionID(),false, false)
     if err8 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        return "","","","","","","",""
    }
    time.Sleep(time.Duration(*timeout) * time.Second)
    task7 := showScriptTask.GetData()
    taskList7:= task7["tasks"].([]interface{})
    taskMap7 := taskList7[0].(map[string]interface{})
    taskhash7 := taskMap7["task-id"]
    response7:=getTask(client,taskhash7.(string))


    payload8 := map[string]interface{}{
    "script-name" : "df -h",
    "script" : "df -h",
    "targets" : server,
    }

    showScriptTask,err9 := client.ApiCall("run-script",payload8, client.GetSessionID(),false, false)
     if err9 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        return "","","","","","","",""
    }
    time.Sleep(time.Duration(*timeout) * time.Second)
    task8 := showScriptTask.GetData()
    taskList8 := task8["tasks"].([]interface{})
    taskMap8 := taskList8[0].(map[string]interface{})
    taskhash8 := taskMap8["task-id"]
    response8:=getTask(client,taskhash8.(string))
    return response1,response2,response3,response4,response5,response6,response7,response8
}

func GetPerformanceData(client *api.ApiClient, gatewayname string,scphost string,scpusername string,scppassword string,scpdestpath string) {
    server := "\""+gatewayname+"\""
    payload := map[string]interface{}{
    "script-name" : "Export CPVIEW History",
    "script" : "cpview -s export",
    "targets" : server,
    }
    showScriptTask,err2 := client.ApiCall("run-script",payload, client.GetSessionID(),false, false)
     if err2 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        //return ""
    }
    time.Sleep(time.Duration(60) * time.Second)
    task := showScriptTask.GetData()
    taskList:= task["tasks"].([]interface{})
    taskMap := taskList[0].(map[string]interface{})
    taskhash := taskMap["task-id"]
    response:=getTask(client,taskhash.(string))
    response2,_:=splitstring(response,"Exported file to '")
    response2=strings.Replace(response2, "'", "", -1)
    var cmd string
    if scppassword != "blank" {
        cmd="sshpass -p '"+ scppassword + "' scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -r " + response2 + " " + scpusername + "@" + scphost +":" + scpdestpath + gatewayname + ".tgz"
    }
    if scppassword == "blank" {
        cmd=" scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -r " + response2 + " " + scpusername + "@" + scphost +":" + scpdestpath + gatewayname + ".tgz"
    }
     payload2 := map[string]interface{}{
    "script-name" : "Export Historical Perfomance Data",
    "script" : cmd,
    "targets" : server,
    }
    _,err3 := client.ApiCall("run-script",payload2, client.GetSessionID(),false, false)
    if err3 != nil {
        fmt.Println("Failed to retrieve the hosts\n")
        //return ""
    }

    //fmt.Println(cmd)
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
        //fmt.Println(objectType,gwipv4address,gwname.(string))
        if objectType == "simple-gateway" || objectType == "CpmiClusterMember" || objectType == "Member" || objectType == "cluster-member" {
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
    licenseinfo=strings.Replace(licenseinfo, "\"", "", -1)
    licenseinfo=strings.Replace(licenseinfo, ","," ", -1)
    licenseinfo= "\"" + licenseinfo + "\"" + "\n"
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
    fwverInfo=strings.Replace(fwverInfo, "\"", "", -1)
    fwverInfo=strings.Replace(fwverInfo, ","," ", -1)
    fwverInfo= "\"" + fwverInfo + "\"" + "\n"
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


func processcollectioninfo(fwverfilename string,hostnamestr string,patches string, routed string, securexl string, affinity string, mq string, showconfig string, priq string, storage string) {
    patches=strings.Replace(patches, "\"", "", -1)
    patches=strings.Replace(patches, ","," ", -1)
    //patches=patches + "\n"
    patches = "\"" + patches + "\""
    //fmt.Println(patches)

    routed=strings.Replace(routed, "\"", "", -1)
    routed=strings.Replace(routed, "," , " ", -1)
    //routed=routed + "\n"
    routed = "\"" + routed + "\""
    
    securexl=strings.Replace(securexl, "\"", "", -1)
    securexl=strings.Replace(securexl, "," , " ", -1)
    //securexl=securexl + "\n"
    securexl = "\"" + securexl + "\""
    
    affinity=strings.Replace(affinity, "\"", "", -1)
    affinity=strings.Replace(affinity, "," , " ", -1)
    //affinity=affinity + "\n"
    affinity = "\"" + affinity + "\""
    
    mq=strings.Replace(mq, "\"", "", -1)
    mq=strings.Replace(mq,  "," , " ", -1)
    //mq=mq + "\n"
    mq = "\"" + mq + "\""
    
    showconfig=strings.Replace(showconfig, "\"", "", -1)
    showconfig=strings.Replace(showconfig, "," , " ", -1)
    //showconfig=showconfig + "\n"
    showconfig = "\"" + showconfig + "\""
    
    priq=strings.Replace(priq, "\"", "", -1)
    priq=strings.Replace(priq, "," , " ", -1)
    //priq=priq + "\n"
    priq = "\"" + priq + "\""
    
    storage=strings.Replace(storage, "\"", "", -1)
    storage=strings.Replace(storage, "," , " ", -1)
    //storage=storage + "\n"
    storage = "\"" + storage + "\""
    
    datacsv:=hostnamestr + "," + patches + "," + routed + ","+  securexl +  ";" +  affinity +  ","  + mq +  ","  +  showconfig + ","  +  priq +  "," + storage + "\n"
    //// TODO: ADD  PARSER FOR DATA
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
    var apiServer,username,password,scphost,scpusername,scppassword,scpdestpath,domaintarget *string
    //var config,historicaldb bool
    
    timeout=flag.Int("timeout",2,"Time delay before getting results from gateway")
    apiServer=flag.String("apiserver","0.0.0.0","Check Point Management IP")
    username=flag.String("username","unknown","Domain User on MDM")
    password=flag.String("password","unknown","DOMAIN PASSWORD")
    config:=flag.Bool("config",false,"Pull Configuration for health checks")
    historicaldb:=flag.Bool("historicaldb",false,"SCP FIles")
    scphost=flag.String("scphost","192.168.1.1","SCP To transfer Historical Performance and spike detective")
    //scpusername=flag.String("scpusername", "admin","SCP username for Spike Detector and Historical DB")
    //scppassword=flag.String("scppassword", "blank","SCP Password, if blank, assumed cert")
    //scpdestpath=flag.String("spdestpath","/", "Default path of SCP")
    domaintarget=flag.String("domaintarget","NIL", "File with Domains IPs to run.  Each IP should be seperated by a ,")
    defaultport=flag.Int("defaultport",443,"Default API port")
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
    timeout2 := time.Duration(*timeout)
    fmt.Println("\n \n ######################################################################")
    args := api.APIClientArgs(*defaultport, "", "", *apiServer, "", -1, "", false, false, "deb.txt", api.WebContext, timeout2, api.SleepTime, "", "")
    client := api.APIClient(args)
    if x, _ := client.CheckFingerprint(); !x {
        print("Could not get the server's fingerprint - Check connectivity with the server.\n")
        os.Exit(1)
    }

    loginRes, err := client.Login(*username, *password, false, "", false, "")
    if err != nil {
        fmt.Println("Login error. \n", err)
        os.Exit(1)
    }

    if !loginRes.Success {
        fmt.Println("Login failed:\n" + loginRes.ErrorMsg)
        os.Exit(1)
    }
    assetinfofilename,licenseinfofilename,interfaceInfofilename,ifconfigInfofilename,interfacelistfilename,fwverfilename,datacollectionfilename:=createloggingenvirorment()
    // IF domain.targets empty, get all domains
    if *domaintarget == "NIL" {
        domainipv4=getDomains(client)
    }
    // IF domain.targets exist, read for list of domains
    if *domaintarget != "NIL" {
        file, err := os.Open(*domaintarget)
        if err != nil {
            fmt.Println(err)        }
        defer file.Close()

        var lines []string
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            lines = append(lines, scanner.Text())
        }
        domainipv4 = strings.Split(lines[0],",")
    }


    fmt.Println("Management Domain's Found at ", *apiServer, ": " , domainipv4)
    for a := 0; a < len(domainipv4); a++ {
        fmt.Println("\n")
        args2 := api.APIClientArgs(*defaultport, "", "",domainipv4[a], "", -1, "", false, false, "deb.txt", api.WebContext, timeout2, api.SleepTime, "", "")
        client2 := api.APIClient(args2)
            if x, _ := client.CheckFingerprint(); !x {
                print("Could not get the server's fingerprint - Check connectivity with the server.\n")
                os.Exit(1)
            }
        loginRes2, err := client2.Login(*username, *password, true, domainipv4[a], false, "")
            if err != nil {
                fmt.Println("Login error.\n", err)
                os.Exit(1)
            }
        if !loginRes2.Success {
            fmt.Println("Login failed:\n" + loginRes2.ErrorMsg)
            os.Exit(1)
            continue
        }
        gwaddress:=getGatewayList(client2)
        fmt.Println("Domain Address: ",domainipv4[a] , "Gateways on Domain: ", gwaddress)   
            for i := 0; i < len(gwaddress); i++ {
                fmt.Println("     Collecting from Firewall: ", gwaddress[i], "      From Domain: ", domainipv4[a])
                    hostname:=getHostname(client2,gwaddress[i])
                    if hostname != "OFFLINE" {
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
                        if *config == true {
                            version,routedver,securexl,affinity,mq,showconfig,priq,storage:=GetFirewallConfiguration(client2,gwaddress[i])
                            processcollectioninfo(datacollectionfilename,hostnamestr,version,routedver,securexl,affinity,mq,showconfig,priq,storage)
                        }
                        if *historicaldb == true {
                            GetPerformanceData(client2,gwaddress[i],*scphost,*scpusername,*scppassword,*scpdestpath)
                            //fmt.Println(scphost,scpusername,scppassword,scpdestpath)
                        }
                    }
                    if hostname == "OFFLINE" { fmt.Println("                 ",gwaddress[i], " is OFFLINE") }
                }
            }
        gwaddress:=getGatewayList(client)
        fmt.Println("\n", "Full Gateways Inventoried: " , gwaddress)
}
