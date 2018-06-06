package main

import (
	"fmt"
	"os/exec"
	"strconv"
)

//for using external processes - IPMI and smartctl calls

func main() {

	// IPMI Variables/Settings
	//useIPMI := true //I want to use IPMI, use what you like
	ipmiHost := "192.168.1.60" // IP address or DNS-resolvable hostname of IPMI server:
	ipmiUser := "ADMIN"        // IPMI username
	// IPMI password file. This is a file containing the IPMI user's password
	// on a single line and should have 0600 permissions:
	ipmiPWFile := "/root/ipmi_password" //just the file location

	commandText := []string{"bash", "-c", "/usr/local/bin/ipmitool -I lanplus -H", ipmiHost, "-U", ipmiUser, "-f", ipmiPWFile, "sdr elist all | grep -c -i \"cpu.*temp\""}
	fmt.Printf("commandText is of type: %T\nWith a value of: %v\nAnd a length of: %v\n\n", commandText, commandText, len(commandText))
	for i, s := range commandText {
		fmt.Printf("Value at %v is %v\n\n", i, s)
	}

	//define the command to get the number of CPUs and then use it
	numCPUCmd := exec.Command(commandText[0], commandText[1:]...)

	fmt.Println("Path:", numCPUCmd.Path)
	fmt.Println("Args:", numCPUCmd.Args)

	//command to get the temperature of multiple CPUs.  Take careful note of index 9 or [9] as this is what we will be changing to iterate through the CPUs, and later drives (with smartclr)
	multiTempCPUcmdTxt := []string{"bash", "-c", "/usr/local/bin/ipmitool -I lanplus -H", ipmiHost, "-U", ipmiUser, "-f", ipmiPWFile, "sdr elist all | grep \"'CPU", "tempString", " Temp\" | awk '{print $10}'"}
	for n, s := range multiTempCPUcmdTxt {
		fmt.Printf("Value at %v is %v\n\n", n, s)
	}

	multiTempCPUcmdTxt[9] = strconv.FormatInt(9000, 10)

	for n, s := range multiTempCPUcmdTxt {
		fmt.Printf("Value at %v is %v\n\n", n, s)
	}
	/*
		//define the command to get the number of CPUs and then use it.  Every , between arguments is a space
		fmt.Println("Command Args:", "/usr/local/bin/ipmitool", "-I lanplus -H", ipmiHost, "-U", ipmiUser, "-f", ipmiPWFile, "sdr elist all | grep -c -i \"cpu.*temp\"")

		numCPUCmd := exec.Command("/usr/local/bin/ipmitool", "-I lanplus -H", ipmiHost, "-U", ipmiUser, "-f", ipmiPWFile, "sdr elist all | grep -c -i \"cpu.*temp\"")
		fmt.Println("numCPUCmd.Path:", numCPUCmd.Path)
		fmt.Printf("numCPUCmd.Args:%v\n\n", numCPUCmd.Args)

		//numCPUBytes, _ := numCPUCmd.Output() //returns a slice of bytes and an error
		numCPUBytes := []byte{50, 10} //setting this manually as my current machine does not have access to the IPMI desired
		fmt.Printf("(Manually set) numCPUBytes type: %T\nvalue:%v\n\n", numCPUBytes, numCPUBytes)

		numCPUString := string(numCPUBytes)
		fmt.Printf("numCPUString type: %T\nvalue:START\"%vEND\"\n\n", numCPUString, numCPUString)

		numCPUTrimmedString := strings.TrimSpace(numCPUString)
		fmt.Printf("numCPUTrimmedString type: %T\nvalue:START\"%vEND\"\n\n", numCPUTrimmedString, numCPUTrimmedString)

		numCPUFloat, err := strconv.ParseFloat(numCPUTrimmedString, 64)
		if err != nil {
			log.Error(err)
		}
		fmt.Printf("numCPUFloat type: %T\nvalue:%v\n\n", numCPUFloat, numCPUFloat)

		numCPU := int(numCPUFloat) //converts the first and hopefully only value of slice of bytes into an int
		fmt.Printf("numCPU type: %T\nvalue:%v\n\n", numCPU, numCPU)
	*/
}
