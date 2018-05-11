package main

import (
	"io/ioutil" //for reading the password for IPMI
	"os/exec"   //for using external processes - IPMI and smartctl calls
	"strconv"   //used in the conversion of an in to a string in the final function which pushes the data to be a metric

	"github.com/prometheus/client_golang/prometheus"
)

//--------------------------------
// *** CPU Temperatures ***
//--------------------------------

//get the # of CPU cores and their temperatures and put them in a cpu struct, then return a slice of these structs
func getCPUtemps() (out []cpu) {

	// IPMI Variables/Settings
	//useIPMI := true //I want to use IPMI, use what you like
	ipmiHost := "IPMI ADDRESS"  // IP address or DNS-resolvable hostname of IPMI server:
	ipmiUser := "IPMI Username" // IPMI username
	// IPMI password file. This is a file containing the IPMI user's password
	// on a single line and should have 0600 permissions:
	ipmiPWfromFile, err := ioutil.ReadFile("/root/ipmi_password") //needs to find the file at location and read the line to the variable
	ipmiPW := string(ipmiPWfromFile)

	out = make([]cpu, 2)

	//define the command to get the number of CPUs and then use it
	numCPUCmd := exec.Command("/usr/local/bin/ipmitool", " -I lanplus -H ", ipmiHost, " -U ", ipmiUser, " -f ", ipmiPW, " sdr elist all | grep -c -i 'cpu.*temp")
	numCPUSoB, _ := numCPUCmd.Output() //returns a slice of bytes and an error
	numCPU := int(numCPUSoB[0])        //converts the first and hopefully only value of slice of bytes into an int

	//go through each CPU and get the temperature
	if numCPU == 1 {
		//define the command used to get the CPU temperature
		cpuTemp := exec.Command("/usr/local/bin/ipmitool", " -I lanplus -H ", ipmiHost, " -U ", ipmiUser, " -f ", ipmiPW, " sdr elist all | grep 'CPU Temp' | awk '{print $10}'")
		temp, _ := cpuTemp.Output()
		out = append(cpu, int(temp[0]))
	} else {
		for i := 0; i < numCPU; i++ {
			cpuTemp := exec.Command("/usr/local/bin/ipmitool", " -I lanplus -H ", ipmiHost, " -U ", ipmiUser, " -f ", ipmiPW, " sdr elist all | grep 'CPU", string(i), " Temp' | awk '{print $10}'")
			temp, _ := cpuTemp.Output()
			out = append(cpu, int(temp[0]))
		}
	}
	return out // returns the slice of cpu structs
}

//a struct (like a multi field object) for the CPUs, leaves room for more later
type cpu struct {
	temp *prometheus.Desc
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func newCPUCollector() *cpu {
	return &cpu{
		temp: prometheus.NewDesc("cpu_temp",
			"Displays the current CPU temperatures",
			nil, nil,
		),
	}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *cpu) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.temp
}

//Collect implements required collect function for all promehteus collectors
func (collector *cpu) Collect(ch chan<- prometheus.Metric) {

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
	for i := range getCPUtemps() {
		cpuTemp := strconv.Itoa(i)
		ch <- prometheus.MustNewConstMetric(collector.temp, prometheus.GuageValue, cpuTemp)
	}
}
