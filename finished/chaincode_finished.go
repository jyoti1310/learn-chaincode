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

var employeeLogBog map[string]SKATEmployee

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
	} else if  function == "addToLogBog" {
		return t.addSKATEmployee(stub, args)
	} else if  function == "updateLogBog" {
		return t.updateSKATEmployee(stub, args)
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
	} else if function == "searchLogBog" {
		return t.searchSKATEmployee(stub,args)
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
	
	key = strconv.Itoa(Employee.CPRNum) + "_" + strconv.Itoa(Employee.VirkNum) + "_" + Employee.DateOfWork
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
	fmt.Println("len(employeeRepository) in Add:"+ strconv.Itoa(len(employeeRepository.EmployeeList)));						
	fmt.Println("- end add Employee 2")
	return jsonAsBytes, nil
}

// ============================================================================================================================
// Search Employees
// ============================================================================================================================
func (t *SimpleChaincode) searchSKATEmployee(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var cprNo, virkNo, cprForEmployee , virkForEmployee string
	var jsonResp  string
	var err error
	//var SearchedEmployeeList []SKATEmployee 

	SearchedEmployeeList := []SKATEmployee{}
	
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting CPRNum , VirkNum as input")
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
	
	fmt.Println("len(employeeRepository) in search:"+ strconv.Itoa(len(employeeRepository.EmployeeList)));
	/*for i := range employeeRepository.EmployeeList{
		cprForEmployee = strconv.Itoa(employeeRepository.EmployeeList[i].CPRNum)
		fmt.Println("matching record:"+ strconv.Itoa(i))
		//fmt.Println("looking at " + strconv.FormatInt(trades.OpenTrades[i].Timestamp, 10) + " for " + strconv.FormatInt(timestamp, 10))
		if 	(strings.Contains(cprForEmployee,cprNo) || strings.Contains(virkForEmployee,virkNo)){
			fmt.Println("found the employee 1");
			SearchedEmployeeList = append(SearchedEmployeeList,employeeRepository.EmployeeList[i])
			fmt.Println("found the employee 2"+ SearchedEmployeeList[i].CPRNavn );
			fmt.Println("SearchedEmployeeList[:"+ strconv.Itoa(i) +"] ==", SearchedEmployeeList[i])		
		}
	}*/
	//jsonAsBytes, _ := json.Marshal(len(SearchedEmployeeList[:]))
	result := "["
	
	for _,skatEmployee := range employeeRepository.EmployeeList{
		cprForEmployee = strconv.Itoa(skatEmployee.CPRNum)
		virkForEmployee= strconv.Itoa(skatEmployee.VirkNum)
		fmt.Println("matching record")
		//fmt.Println("looking at " + strconv.FormatInt(trades.OpenTrades[i].Timestamp, 10) + " for " + strconv.FormatInt(timestamp, 10))
		if 	(strings.Contains(cprForEmployee,cprNo) || strings.Contains(virkForEmployee,virkNo)){
			fmt.Println("found the employee 1");
			SearchedEmployeeList = append(SearchedEmployeeList,skatEmployee)
			fmt.Println("found the employee 2"+ skatEmployee.CPRNavn );
			temp, err := json.Marshal(skatEmployee)
			if err == nil {
			result += string(temp) + ","
			}
		}
	
	}
	if len(result) == 1 {
		result = "[]"
	} else {
		result = result[:len(result)-1] + "]"
	}	
	fmt.Println("SearchedEmployeeList:", SearchedEmployeeList[0:])
/*	var interfaceSlice []interface{} = make([]interface{}, len(SearchedEmployeeList))
	for i, skatEmployee := range SearchedEmployeeList {
    interfaceSlice[i] = skatEmployee
	jsonAsBytes, _ := json.Marshal(interfaceSlice)
	}*/
		
	return []byte(result), nil

}
// ============================================================================================================================
// Update Employee - Update Employee with Comments, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) updateSKATEmployee(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
var cprNum,comment, key,jsonResp string
var employee SKATEmployee
var err error


if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting CPRNum , Comments as input")
	}

comment=args[1]
employee , err= t.getEmployee(stub,cprNum)
employee.Comment = comment
key = strconv.Itoa(employee.CPRNum) + "_" + strconv.Itoa(employee.VirkNum) + "_" + employee.DateOfWork
jsonAsBytes, _ := json.Marshal(employee)
err = stub.PutState(key, []byte(jsonAsBytes)) //write the variable into the chaincode state

if err != nil {
		jsonResp = "{\"Error\":\"Failed to update Employee" + "\"}"
		return nil, errors.New(jsonResp)
	}
_, err = t.updateEmployeeRepository(stub,key,employee)
if err != nil {
		jsonResp = "{\"Error\":\"Failed to update SKAT Repository" + "\"}"
		return nil, errors.New(jsonResp)
	}

return jsonAsBytes,nil
}

//==================================================================================================================================

//===================================================================================================================================

func (t *SimpleChaincode) updateEmployeeRepository(stub shim.ChaincodeStubInterface,  key string,employee SKATEmployee) (bool, error){

var jsonResp string 
repositoryJsonAsBytes, err := stub.GetState("SKATEmployeeRepository")
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + "SKATEmployeeRepository" + "\"}"
		return false, errors.New(jsonResp)
	}
	var employeeRepository SKATEmployeeRepository
	json.Unmarshal(repositoryJsonAsBytes, &employeeRepository)	


   	employeeRepository.EmployeeList = append(employeeRepository.EmployeeList,employee)
	//update Employee Repository
	updatedRepositoryJsonAsBytes, _  := json.Marshal(employeeRepository)
	err = stub.PutState("SKATEmployeeRepository", updatedRepositoryJsonAsBytes)	//store employee with id as key
	if err != nil {
		return false, err
	}		
	return true, nil
}
// ============================================================================================================================
// Get single Employee
// ============================================================================================================================
func (t *SimpleChaincode) getEmployee(stub shim.ChaincodeStubInterface, cprNum string) (SKATEmployee, error) {

	var jsonResp, cprForEmployee  string
	var employee SKATEmployee
	
	repositoryJsonAsBytes, err := stub.GetState("SKATEmployeeRepository")
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + "SKATEmployeeRepository" + "\"}"
		return employee, errors.New(jsonResp)
	}
	var employeeRepository SKATEmployeeRepository
	json.Unmarshal(repositoryJsonAsBytes, &employeeRepository)	
	
	for _,skatEmployee := range employeeRepository.EmployeeList{
		cprForEmployee = strconv.Itoa(skatEmployee.CPRNum)
		
		fmt.Println("matching record")
		//fmt.Println("looking at " + strconv.FormatInt(trades.OpenTrades[i].Timestamp, 10) + " for " + strconv.FormatInt(timestamp, 10))
		if 	(strings.Contains(cprForEmployee,cprNum)){
			fmt.Println("found the employee 1" + skatEmployee.CPRNavn );
			employee=skatEmployee
			}
		
		}
			return employee, nil
	
	}
	
