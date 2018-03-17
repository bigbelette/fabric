/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * based on fabcar network from IBM examples
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the ticket structure, with 4 properties.  Structure tags are used by encoding/json library
type Ticket struct {
	Organisator string `json:"organisator"`
	Event       string `json:"event"`
	Date        string `json:"date"`
	Owner       string `json:"owner"`
}

/*
 * The Init method is called when the Smart Contract "testNetwork" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "testNetwork"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryTicket" {
		return s.queryTicket(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createTicket" {
		return s.createTicket(APIstub, args)
	} else if function == "queryAllTickets" {
		return s.queryAllTickets(APIstub)
	} else if function == "changeTicketOwner" {
		return s.changeTicketOwner(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ticketAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(ticketAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	tickets := []Ticket{
		Ticket{Organisator: "MadisonSquareGarden", Event: "Knicks-Bulls", Date: "20180501", Owner: "Carlos"},
		Ticket{Organisator: "MadisonSquareGarden", Event: "Knicks-Jazz", Date: "20180602", Owner: "Carlos"},
		Ticket{Organisator: "MadisonSquareGarden", Event: "Knicks-Heat", Date: "20180404", Owner: "Emad"},
		Ticket{Organisator: "BroadwayTheater", Event: "BruceSpringsteen", Date: "20180715", Owner: "Varad"},
		Ticket{Organisator: "BroadwayTheater", Event: "Madonna", Date: "201805010", Owner: "Phil"},
		Ticket{Organisator: "BroadwayTheater", Event: "CatStevens", Date: "20180519", Owner: "John"},
	}

	i := 0
	for i < len(tickets) {
		fmt.Println("i is ", i)
		ticketAsBytes, _ := json.Marshal(tickets[i])
		APIstub.PutState("Ticket"+strconv.Itoa(i), ticketAsBytes)
		fmt.Println("Added", ticket[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var ticket = Ticket{Organisator: args[1], Event: args[2], Date: args[3], Owner: args[4]}

	ticketAsBytes, _ := json.Marshal(ticket)
	APIstub.PutState(args[0], ticketAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllTickets(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "TICKET0"
	endKey := "TICKET999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllTickets:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeTicketOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ticketAsBytes, _ := APIstub.GetState(args[0])
	ticket := Ticket{}

	json.Unmarshal(ticketAsBytes, &ticket)
	ticket.Owner = args[1]

	ticketAsBytes, _ = json.Marshal(ticket)
	APIstub.PutState(args[0], ticketAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
