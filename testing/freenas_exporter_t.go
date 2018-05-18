package main

import (
	"fmt"
	"os/exec" //for using external processes - IPMI and smartctl calls
	"strconv"
	"strings"

	"github.com/prometheus/common/log"
)

func main() {

	// IPMI Variables/Settings
	//useIPMI := true //I want to use IPMI, use what you like
	ipmiHost := "192.168.1.60" // IP address or DNS-resolvable hostname of IPMI server:
	ipmiUser := "ADMIN"        // IPMI username
	// IPMI password file. This is a file containing the IPMI user's password
	// on a single lin and should have 0600 permissions:
	ipmiPWFile := "/root/ipmi_password" //just the file location

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
	fmt.Printf("numCPUString type: %T\nvalue:START\"%vEND\"\n\n", numCPUTrimmedString, numCPUTrimmedString)

	numCPUFloat, err := strconv.ParseFloat(numCPUTrimmedString, 64)
	if err != nil {
		log.Error(err)
	}
	fmt.Printf("numCPUFloat type: %T\nvalue:%v\n\n", numCPUFloat, numCPUFloat)

	numCPU := int(numCPUFloat) //converts the first and hopefully only value of slice of bytes into an int
	fmt.Printf("numCPU type: %T\nvalue:%v\n\n", numCPU, numCPU)
}

/*
//a struct (like a multi field object) for the CPUs, leaves room for more later
type cpuCollector struct {
	temp *prometheus.Desc
}

//get the # of CPU cores and their temperatures and put them in a cpu struct, then return a slice of these structs
func getCPUtemps() (out []float64) {

	// IPMI Variables/Settings
	//useIPMI := true //I want to use IPMI, use what you like
	ipmiHost := "192.168.1.60" // IP address or DNS-resolvable hostname of IPMI server:
	ipmiUser := "ADMIN"        // IPMI username
	// IPMI password file. This is a file containing the IPMI user's password
	// on a single lin and should have 0600 permissions:
	ipmiPWFile := "/root/ipmi_password" //just the file location

	//define the command to get the number of CPUs and then use it
	numCPUCmd := exec.Command("/usr/local/bin/ipmitool", " -I lanplus -H ", ipmiHost, " -U ", ipmiUser, " -f ", ipmiPWFile, " sdr elist all | grep -c -i 'cpu.*temp'")
	numCPUBytes, _ := numCPUCmd.Output() //returns a slice of bytes and an error
	numCPUBits := binary.LittleEndian.Uint64(numCPUBytes)
	numCPUFloat := math.Float64frombits(numCPUBits)
	numCPU := int(numCPUFloat) //converts the first and hopefully only value of slice of bytes into an int

	fmt.Println(numCPU) //error checking

	//go through each CPU and get the temperature
	if numCPU == 1 {
		//define the command used to get the CPU temperature
		tempCmd := exec.Command("/usr/local/bin/ipmitool", " -I lanplus -H ", ipmiHost, " -U ", ipmiUser, " -f ", ipmiPWFile, " sdr elist all | grep 'CPU Temp' | awk '{print $10}'")
		fmt.Println(tempCmd) //error checking

		out = append(out, tempFloat)
		fmt.Println("Printing out", out) //error checking
	} else {
		for i := numCPU; i > 0; i-- {
			tempCmd := exec.Command("/usr/local/bin/ipmitool", " -I lanplus -H ", ipmiHost, " -U ", ipmiUser, " -f ", ipmiPWFile, " sdr elist all | grep 'CPU", string(i), " Temp' | awk '{print $10}'")
			fmt.Println(tempCmd) //error checking
		}
	}
	fmt.Println(out) //error checking
	return out       // returns the slice of float64s
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func newCPUCollector() *cpuCollector {
	return &cpuCollector{
		temp: prometheus.NewDesc("cpu_temp_celcius",
			"Displays the current CPU temperatures in Celcius",
			nil, nil,
		),
	}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *cpuCollector) Describe(ch chan<- *prometheus.Desc) {
	//Update this section with the each metric you create for a given collector
	ch <- collector.temp
}

//Collect implements required collect function for all promehteus collectors
func (collector *cpuCollector) Collect(ch chan<- prometheus.Metric) {

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.

	//gathers the slice of cpuCollector structs, each with a temp *prometheus.Desc
	cpuTemps := getCPUtemps()

	for num, temp := range cpuTemps {
		ch <- prometheus.MustNewConstMetric(collector.temp, prometheus.GaugeValue, temp, string(num))
	}
}
*/
