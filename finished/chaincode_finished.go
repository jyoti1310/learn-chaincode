/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type SKATEmployee struct{
        CPRNum int `json:"CPRNum"`
        VirkNum int `json:"VirkNum"`
        CPRNavn string `json:"CPRNavn"`
        DateOfWork string `json:"DOW"`
		NoOfHours int `json:"NoOfHours"`
		Comment string `json:"Comments"`
}

type SKATEmployeeRepository struct{
	EmployeeList []SKATEmployee `json:"employee_list"`
}

type searchedEmployees struct{
	SearchedEmployeeList []SKATEmployee `json:"searched_employee_list"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_Block", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if  function == "addSKATEmployee" {
		return t.addSKATEmployee(stub, args)
	} else if function == "searchSKATEmployee" {
		return t.searchSKATEmployee(stub,args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

// ============================================================================================================================
// Init Employee - create a new Employee, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) addSKATEmployee(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var key,jsonResp string
	
	// CPRNum int `json:"CPRNum"`
    // VirkNum int `json:"VirkNum"`
      //  CPRNavn string `json:"CPRNavn"`
       // DOW string `json:"DOW"`
	//NoOfHours int `json:"NoOfHours"`
	//Comment string `json:"Comments"`
	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	//input sanitation
	fmt.Println("- start init SKATEmployee")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	Employee := SKATEmployee{}
	Employee.CPRNum, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}
	Employee.VirkNum, err = strconv.Atoi(args[1])
	Employee.CPRNavn = strings.ToLower(args[2])
	Employee.DateOfWork = strings.ToLower(args[3])
	Employee.NoOfHours, err = strconv.Atoi(args[4])

	if len(args) == 6 {
	Employee.Comment = args[5]
  	}
	fmt.Println("adding employee @ " + strconv.Itoa(Employee.CPRNum) + ", " + strconv.Itoa(Employee.VirkNum) + ", " + Employee.CPRNavn);
	fmt.Println("- end add Employee 1")
	jsonAsBytes, _ := json.Marshal(Employee)

	if err != nil {
		return jsonAsBytes, err
	}
	
	 key = strconv.Itoa(Employee.CPRNum)
	err = stub.PutState(key, jsonAsBytes)	//store employee with id as key
	if err != nil {
		return nil, err
	}	
	//var  employeeRepository SKATEmployeeRepository
	//employeeRepository.EmployeeList = append(employeeRepository.EmployeeList,Employee)
	//SKATEmployeeRepository.EmployeeList = append(SKATEmployeeRepository.EmployeeList,Employee)
	repositoryJsonAsBytes, err := stub.GetState("SKATEmployeeRepository")
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + "SKATEmployeeRepository" + "\"}"
		return nil, errors.New(jsonResp)
	}
	var employeeRepository SKATEmployeeRepository
	json.Unmarshal(repositoryJsonAsBytes, &employeeRepository)	


   	employeeRepository.EmployeeList = append(employeeRepository.EmployeeList,Employee)
	//update Employee Repository
	updatedRepositoryJsonAsBytes, _  := json.Marshal(employeeRepository)
	err = stub.PutState("SKATEmployeeRepository", updatedRepositoryJsonAsBytes)	//store employee with id as key
	if err != nil {
		return nil, err
	}		
							
	fmt.Println("- end add Employee 2")
	return jsonAsBytes, nil
}

// ============================================================================================================================
// Search Employee 
// ============================================================================================================================
func (t *SimpleChaincode) searchSKATEmployee(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var cprNo, virkNo, cprForEmployee , virkForEmployee string
	var jsonResp  string
	var err error
	//var SearchedEmployeeList []SKATEmployee 

	SearchedEmployeeList := []SKATEmployee{}
	
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	cprNo = args[0]
	virkNo = args[1]
  	
	repositoryJsonAsBytes, err := stub.GetState("SKATEmployeeRepository")
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + "SKATEmployeeRepository" + "\"}"
		return nil, errors.New(jsonResp)
	}
	var employeeRepository SKATEmployeeRepository
	json.Unmarshal(repositoryJsonAsBytes, &employeeRepository)	

	for i := range employeeRepository.EmployeeList{
		cprForEmployee = strconv.Itoa(employeeRepository.EmployeeList[i].CPRNum)
																	//look for the trade
		//fmt.Println("looking at " + strconv.FormatInt(trades.OpenTrades[i].Timestamp, 10) + " for " + strconv.FormatInt(timestamp, 10))
		if 	(strings.Contains(cprForEmployee,cprNo) || strings.Contains(virkForEmployee,virkNo)){
			fmt.Println("found the employee 1");
			SearchedEmployeeList = append(SearchedEmployeeList,employeeRepository.EmployeeList[i])
			fmt.Println("found the employee 2"+ SearchedEmployeeList[i].CPRNavn );
			fmt.Println("SearchedEmployeeList[:"+ strconv.Itoa(i) +"] ==", SearchedEmployeeList[:i])		
		}
	}
	//jsonAsBytes, _ := json.Marshal(len(SearchedEmployeeList[:]))

	var interfaceSlice []interface{} = make([]interface{}, len(SearchedEmployeeList))
	for i, skatEmployee := range SearchedEmployeeList {
    interfaceSlice[i] = skatEmployee
	}
	jsonAsBytes, _ := json.Marshal(interfaceSlice)
	return jsonAsBytes, nil

}