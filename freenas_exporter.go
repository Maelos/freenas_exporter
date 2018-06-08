package main

import (
	"fmt"     //General printing for error catching and testing
	"os/exec" //for using external processes - IPMI and smartctl calls
	"strconv" //to convert a string to a float or int
	"strings" //to trim a string of spaces

	"github.com/prometheus/client_golang/prometheus" //Prometheus files
	"github.com/prometheus/common/log"               //other Prometheus files, specifically for logging
)

//a struct (like a multi field object) for the CPUs, leaves room for more later
type cpuCollector struct {
	temp *prometheus.Desc
}

func check(err error) {
	if err != nil {
		log.Error(err)
	}
}

//get the # of CPU cores and their temperatures and put them in a cpu struct, then return a slice of these structs
func getCPUtemps() (out []float64) {

	// IPMI Variables/Settings
	//useIPMI := true //I want to use IPMI, use what you like
	ipmiHost := "192.168.1.60" // IP address or DNS-resolvable hostname of IPMI server:
	ipmiUser := "ADMIN"        // IPMI username
	// IPMI password file. This is a file containing the IPMI user's password
	// on a single line and should have 0600 permissions:
	ipmiPWFile := "/root/ipmi_password" //just the file location

	//Would like to get this working with Golang standard libraries but want to try to just get it to work for now
	cmdTxt := fmt.Sprintf("/usr/local/bin/ipmitool -I lanplus -H %s -U %s -f %s sdr elist all | grep -c -i \"cpu.*temp\"", ipmiHost, ipmiUser, ipmiPWFile)

	//command to get the number of CPUs
	ipmiCmd := exec.Command("bash", "-c", cmdTxt)
	ipmiOutput, err := ipmiCmd.Output()
	check(err)

	numCPU, err := strconv.Atoi(strings.TrimSpace(string(ipmiOutput)))
	check(err)
	/*I want to get this working using Go and not relying on the grep and awk commands, but for now I will do it the ugly way
	numRegexp := regexp.Compile("(?i)cpu.*temp") //define the regular expression to use the full suite of tools
	numCPU := len(numRegexp.FindAll(ipmiOutput, -1)) //gather the length of an array that returns all occurances of the regular expression
	oneCPUregexp := regexp.Compile("(?i)")
	*/

	//go through each CPU and get the temperature
	if numCPU == 1 {
		//the single string command text after bash -c
		oneCPUcmdTxt := fmt.Sprintf("/usr/local/bin/ipmitool -I lanplus -H %s -U %s -f %s sdr elist all | grep \"CPU Temp\" | awk '{print $10}'", ipmiHost, ipmiUser, ipmiPWFile)
		oneCPUcmd := exec.Command("bash", "-c", oneCPUcmdTxt)
		oneCPUbytes, err := oneCPUcmd.Output()
		check(err)

		tempFloat, err := strconv.ParseFloat(strings.TrimSpace(string(oneCPUbytes)), 64)

		check(err) //error check

		out = append(out, tempFloat)
		fmt.Println("3. Single entry slice of float", out) //error checking

	} else {
		for i := 1; i < numCPU+1; i++ {
			//multi cpu command, though gonna have to make this work with the look
			multiCPUcmdTxt := fmt.Sprintf("/usr/local/bin/ipmitool -I lanplus -H %s -U %s -f %s sdr elist all | grep \"CPU%s Temp\" | awk '{print $10}'", ipmiHost, ipmiUser, ipmiPWFile, strconv.Itoa(i))
			multiCPUcmd := exec.Command("bash", "-c", multiCPUcmdTxt)
			multiCPUbytes, err := multiCPUcmd.Output()
			check(err)

			tempFloat, err := strconv.ParseFloat(strings.TrimSpace(string(multiCPUbytes)), 64)
			check(err) //error check

			out = append(out, tempFloat)
			fmt.Println("6. Slice of float 64s:", out) //error checking
		}
	}
	fmt.Println("7. Final slice being returned:", out) //error checking
	return out                                         // returns the slice of float64s
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

	//by taking out the last bit of the string I am now getting a working program, or at least it is pulling the temperatures correctly
	//I now need to figure out how to get Prometheus working.  I think this is progress...?
	for _, temp := range cpuTemps {
		ch <- prometheus.MustNewConstMetric(collector.temp, prometheus.GaugeValue, temp) //, strconv.Itoa(num + 1))
	}
}
