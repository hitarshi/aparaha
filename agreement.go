/*/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
"errors"
"fmt"
"strconv"
"encoding/json"
"time"

"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ManageLoan example simple Chaincode implementation
type ManageLoan struct {
}

var LoanIndexStr = "_LoanIndex"				//name for the key/value that will store a list of all known Agreement

type Agreement struct{							// Attributes of a Agreement 
	AgreeementID string `json:"agreement_id"`					
	BorrowerName string `json:"borrower_name"`
	LenderName string `json:"lender_name"`					
	AgreementDate string `json:"agreement_date"`
	AgreementStatus string `json:"agreement_status"`
	LoanAmount string `json:"loan_amount"`
	InterestRate string `json:"interest_rate"`
	LoanDuration string `json:"loan_duration"`
	RepaymentDate string `json:"repayment_date"`
	BorrowerSigned string `json:"borrower_signed"`
	LenderSigned string `json:"lender_signed"`
	Comments string `json:"comments"`
}
// ============================================================================================================================
// Main - start the chaincode for Agreement management
// ============================================================================================================================
func main() {			
	err := shim.Start(new(ManageLoan))
	if err != nil {
		fmt.Printf("Error starting Agreement management chaincode: %s", err)
	}
}
// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *ManageLoan) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var msg string
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	// Initialize the chaincode
	msg = args[0]
	fmt.Println("ManageLoan chaincode is deployed successfully.");
	
	// Write the state to the ledger
	err = stub.PutState("abc", []byte(msg))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(LoanIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}
// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
	func (t *ManageLoan) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
		fmt.Println("run is running " + function)
		return t.Invoke(stub, function, args)
	}
// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
	func (t *ManageLoan) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
		fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "create_agreement" {											//create a new Agreement
		return t.create_agreement(stub, args)
	}else if function == "delete_po" {									// delete a Agreement
		return t.delete_po(stub, args)
	}else if function == "update_po" {									//update a Agreement
		return t.update_po(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)					//error
	return nil, errors.New("Received unknown function invocation")
}
// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *ManageLoan) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getAgreement_byID" {													//Read a Agreement by agreeMent_id
		return t.getAgreement_byID(stub, args)
	} else if function == "getAgreement_byBuyer" {													//Read a Agreement by Buyer's name
		return t.getAgreement_byBuyer(stub, args)
	} else if function == "getAgreement_bySeller" {													//Read a Agreement by Seller's name
		return t.getAgreement_bySeller(stub, args)
	} else if function == "get_AllAgreement" {													//Read all Agreements
		return t.get_AllAgreement(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error
	return nil, errors.New("Received unknown function query")
}
// ============================================================================================================================
// getAgreement_byID - get Agreement details for a specific ID from chaincode state
// ============================================================================================================================
func (t *ManageLoan) getAgreement_byID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var agreement_id, jsonResp string
	var err error
	fmt.Println("start getAgreement_byID")
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting ID of the var to query")
	}
	// set agreement_id
	agreement_id = args[0]
	valAsbytes, err := stub.GetState(agreement_id)									//get the agreement_id from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + agreement_id + "\"}"
		return nil, errors.New(jsonResp)
	}
	//fmt.Print("valAsbytes : ")
	//fmt.Println(valAsbytes)
	fmt.Println("end getAgreement_byID")
	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
//  getAgreement_byBuyer - get Agreement details by buyer's name from chaincode state
// ============================================================================================================================
func (t *ManageLoan) getAgreement_byBuyer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, lender_name, errResp string
	var poIndex []string
	var valIndex Agreement
	fmt.Println("start getAgreement_byBuyer")
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 argument")
	}
	// set buyer's name
	lender_name = args[0]
	//fmt.Println("lender_name" + lender_name)
	poAsBytes, err := stub.GetState(LoanIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Agreement index string")
	}
	//fmt.Print("poAsBytes : ")
	//fmt.Println(poAsBytes)
	json.Unmarshal(poAsBytes, &poIndex)								//un stringify it aka JSON.parse()
	//fmt.Print("poIndex : ")
	//fmt.Println(poIndex)
	//fmt.Println("len(poIndex) : ")
	//fmt.Println(len(poIndex))
	jsonResp = "{"
	for i,val := range poIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getAgreement_byBuyer")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		//fmt.Print("valueAsBytes : ")
		//fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		//fmt.Print("valIndex: ")
		//fmt.Print(valIndex)
		if valIndex.LenderName == lender_name{
			fmt.Println("Buyer found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			//fmt.Println("jsonResp inside if")
			//fmt.Println(jsonResp)
			if i < len(poIndex)-1 {
				jsonResp = jsonResp + ","
			}
		}
		
	}
	jsonResp = jsonResp + "}"
	//fmt.Println("jsonResp : " + jsonResp)
	//fmt.Print("jsonResp in bytes : ")
	//fmt.Println([]byte(jsonResp))
	fmt.Println("end getAgreement_byBuyer")
	return []byte(jsonResp), nil											//send it onward
}

// ============================================================================================================================
//  getAgreement_bySeller - get Agreement details for a specific Seller from chaincode state
// ============================================================================================================================
func (t *ManageLoan) getAgreement_bySeller(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, borrower_name, errResp string
	var poIndex []string
	var valIndex Agreement
	fmt.Println("start getAgreement_bySeller")
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 argument")
	}
	// set seller name
	borrower_name = args[0]
	//fmt.Println("lender_name" + borrower_name)
	poAsBytes, err := stub.GetState(LoanIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Agreement index")
	}
	//fmt.Print("poAsBytes : ")
	//fmt.Println(poAsBytes)
	json.Unmarshal(poAsBytes, &poIndex)								//un stringify it aka JSON.parse()
	//fmt.Print("poIndex : ")
	//fmt.Println(poIndex)
	//fmt.Println("len(poIndex) : ")
	//fmt.Println(len(poIndex))
	jsonResp = "{"
	for i,val := range poIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getting borrower_name")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		//fmt.Print("valueAsBytes : ")
		//fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		//fmt.Print("valIndex: ")
		//fmt.Print(valIndex)
		if valIndex.BorrowerName == borrower_name{
			fmt.Println("Seller found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			//fmt.Println("jsonResp inside if")
			//fmt.Println(jsonResp)
			if i < len(poIndex)-1 {
				jsonResp = jsonResp + ","
			}
		}
		
	}
	
	jsonResp = jsonResp + "}"
	//fmt.Println("jsonResp : " + jsonResp)
	//fmt.Print("jsonResp in bytes : ")
	//fmt.Println([]byte(jsonResp))
	fmt.Println("end getAgreement_bySeller")
	return []byte(jsonResp), nil											//send it onward
}
// ============================================================================================================================
//  get_AllAgreement- get details of all Agreement from chaincode state
// ============================================================================================================================
func (t *ManageLoan) get_AllAgreement(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var poIndex []string
	fmt.Println("start get_AllAgreement")
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 argument")
	}
	poAsBytes, err := stub.GetState(LoanIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Agreement index")
	}
	//fmt.Print("poAsBytes : ")
	//fmt.Println(poAsBytes)
	json.Unmarshal(poAsBytes, &poIndex)								//un stringify it aka JSON.parse()
	//fmt.Print("poIndex : ")
	//fmt.Println(poIndex)
	jsonResp = "{"
	for i,val := range poIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all Agreement")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		//fmt.Print("valueAsBytes : ")
		//fmt.Println(valueAsBytes)
		jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
		if i < len(poIndex)-1 {
			jsonResp = jsonResp + ","
		}
	}
	//fmt.Println("len(poIndex) : ")
	//fmt.Println(len(poIndex))
	jsonResp = jsonResp + "}"
	//fmt.Println("jsonResp : " + jsonResp)
	//fmt.Print("jsonResp in bytes : ")
	//fmt.Println([]byte(jsonResp))
	fmt.Println("end get_AllAgreement")
	return []byte(jsonResp), nil
											//send it onward
}
// ============================================================================================================================
// Delete - remove a Agreement from chain
// ============================================================================================================================
func (t *ManageLoan) delete_po(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	// set agreement_id
	agreement_id := args[0]
	err := stub.DelState(agreement_id)													//remove the Agreement from chaincode
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	//get the Agreement index
	poAsBytes, err := stub.GetState(LoanIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Agreement index")
	}
	//fmt.Println("poAsBytes in delete po")
	//fmt.Println(poAsBytes);
	var poIndex []string
	json.Unmarshal(poAsBytes, &poIndex)								//un stringify it aka JSON.parse()
	//fmt.Println("poIndex in delete po")
	//fmt.Println(poIndex);
	//remove marble from index
	for i,val := range poIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + agreement_id)
		if val == agreement_id{															//find the correct Agreement
			fmt.Println("found Agreement with matching agreement_id")
			poIndex = append(poIndex[:i], poIndex[i+1:]...)			//remove it
			for x:= range poIndex{											//debug prints...
				fmt.Println(string(x) + " - " + poIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(poIndex)									//save new index
	err = stub.PutState(LoanIndexStr, jsonAsBytes)
	return nil, nil
}
// ============================================================================================================================
// Write - update Agreement into chaincode state
// ============================================================================================================================
func (t *ManageLoan) update_po(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error
	fmt.Println("start update_po")
	if len(args) != 12 {
		return nil, errors.New("Incorrect number of arguments. Expecting 12.")
	}
	// set agreement_id
	agreement_id := args[0]
	poAsBytes, err := stub.GetState(agreement_id)									//get the Agreement for the specified agreement_id from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + agreement_id + "\"}"
		return nil, errors.New(jsonResp)
	}
	//fmt.Print("poAsBytes in update po")
	//fmt.Println(poAsBytes);
	res := Agreement{}
	json.Unmarshal(poAsBytes, &res)
	if res.AgreeementID == agreement_id{
		fmt.Println("Agreement found with agreement_id : " + agreement_id)
		//fmt.Println(res);
		res.BorrowerName = args[1]
		res.LenderName = args[2]
		res.AgreementDate = args[3]
		res.LoanAmount = args[4]
		res.AgreementStatus = args[5]
		res.InterestRate = args[6]
		res.LoanDuration = args[7]
		res.RepaymentDate = args[8]
		res.BorrowerSigned = args[9]
		res.LenderSigned = args[10]
		res.Comments = args[10]		
	}
	
	//build the Agreement json string manually
	order := 	`{`+
		`"agreement_id": "` + res.AgreeementID + `" , `+
		`"borrower_name": "` + res.BorrowerName + `" , `+
		`"lender_name": "` + res.LenderName + `" , `+
		`"agreement_date": "` + res.AgreementDate + `" , `+ 
		`"loan_amount": "` + res.LoanAmount + `" , `+ 
		`"agreement_status": "` + res.AgreementStatus + `" , `+ 
		`"interest_rate": "` + res.InterestRate + `" , `+ 
		`"loan_duration": "` + res.LoanDuration + `" , `+ 
		`"repayment_date": "` + res.RepaymentDate + `" , `+ 
		`"borrower_signed": "` + res.BorrowerSigned + `" , `+ 
		`"lender_signed": "` +  res.LenderSigned + `" , `+ 
		`"comments": "` +  res.Comments + `" `+ 
		`}`
	err = stub.PutState(agreement_id, []byte(order))									//store Agreement with id as key
	if err != nil {
		return nil, err
	}
	return nil, nil
}
// ============================================================================================================================
// create Agreement - create a new Agreement, store into chaincode state
// ============================================================================================================================
func (t *ManageLoan) create_agreement(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 12 {
		return nil, errors.New("Incorrect number of arguments. Expecting 12")
	}
	fmt.Println("start create_agreement")

	agreement_id := args[0]
	borrower_name := args[1]
	lender_name := args[2]
	agreement_date := args[3]
	loan_amount := args[4]
	agreement_status := args[5]
	interest_rate := args[6]
	loan_duration := args[7]
	repayment_date := args[8]
	borrower_signed := args[9]
	lender_signed := args[10]
	comments := args[11]
	
	poAsBytes, err := stub.GetState(agreement_id)
	if err != nil {
		return nil, errors.New("Failed to get Agreement transID")
	}
	//fmt.Print("poAsBytes: ")
	//fmt.Println(poAsBytes)
	res := Agreement{}
	json.Unmarshal(poAsBytes, &res)
	//fmt.Print("res: ")
	//fmt.Println(res)
	if res.AgreeementID == agreement_id{
		//fmt.Println("This Agreement arleady exists: " + agreement_id)
		//fmt.Println(res);
		return nil, errors.New("This Agreement arleady exists")				//all stop a Agreement by this name exists
	}
	
	//build the Agreement json string manually
	order := 	`{`+
		`"agreement_id": "` + agreement_id + `" , `+
		`"borrower_name": "` + borrower_name + `" , `+
		`"lender_name": "` + lender_name + `" , `+
		`"agreement_date": "` + agreement_date + `" , `+ 
		`"loan_amount": "` + loan_amount + `" , `+ 
		`"agreement_status": "` + agreement_status + `" , `+ 
		`"interest_rate": "` + interest_rate + `" , `+ 
		`"loan_duration": "` + loan_duration + `" , `+ 
		`"repayment_date": "` + repayment_date + `" , `+ 
		`"borrower_signed": "` + borrower_signed + `" , `+ 
		`"lender_signed": "` +  lender_signed + `" , `+
		`"comments": "` +  comments + `" `+ 
		`}`
		//fmt.Println("order: " + order)
		fmt.Print("order in bytes array: ")
		fmt.Println([]byte(order))
	err = stub.PutState(agreement_id, []byte(order))									//store Agreement with agreement_id as key
	if err != nil {
		return nil, err
	}
	//get the Agreement index
	poIndexAsBytes, err := stub.GetState(LoanIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Agreement index")
	}
	var poIndex []string
	//fmt.Print("poIndexAsBytes: ")
	//fmt.Println(poIndexAsBytes)
	
	json.Unmarshal(poIndexAsBytes, &poIndex)							//un stringify it aka JSON.parse()
	//fmt.Print("poIndex after unmarshal..before append: ")
	//fmt.Println(poIndex)
	//append
	poIndex = append(poIndex, agreement_id)									//add Agreement transID to index list
	//fmt.Println("! Agreement index after appending agreement_id: ", poIndex)
	jsonAsBytes, _ := json.Marshal(poIndex)
	//fmt.Print("jsonAsBytes: ")
	//fmt.Println(jsonAsBytes)
	err = stub.PutState(LoanIndexStr, jsonAsBytes)						//store name of Agreement
	if err != nil {
		return nil, err
	}
	fmt.Println("end create_agreement")
	
	fmt.Println("start timer")
	
	timer := NewTimer(10, func() {
    		fmt.Printf("Congratulations! Your %d second timer finished.", 10)
  	})
  	timer.Stop()
	
	fmt.Println("end timer")

	
	return nil, nil
}

func NewTimer(seconds int, action func()) {
  timer := time.NewTimer(time.Seconds * time.Duration(seconds))
  
  go func() {
    <-timer.C
    action()
  }
  
  return timer
}

func action () {
	fmt.Println("check agreement")	
}
