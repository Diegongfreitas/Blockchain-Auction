package main

import (
	//the majority of the imports are trivial...
	"encoding/json"
	"fmt"
	"strconv"
	"bytes"
	//these imports are for Hyperledger Fabric interface
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

/* All the following functions are used to implement auction chaincode. This chaincode
basically works this way:
	An Auctioneer starts the auction of a car announcing it's start price and 
	starts counting the duration of the auction. Then the participants start calling out their bids. 

	There is an asset for the car and another for the bids. 
	The car asset only keeps its ID (license plate) and the start price.
	The bids asset has a unique ID for all bids, registers the car ID, the highest bid value and
	the Social Security Number of the respective participant.

	The auction ends when the time duration reaches it's limit or after 30 seconds without bids.
*/

// SmartContract defines the chaincode base structure. All the methods are implemented to
// return a SmartContrac type.
type SmartContract struct {
}

type Bid struct{
	Value, ClientID string
	
}

type Car struct {
	StartPrice string
	IsOpen bool
	WinnerID string
}

// Init method is called when the auction is instantiated.
// Best practice is to have any Ledger initialization in separate function.
// Note that chaincode upgrade also calls this function to reset
// or to migrate data, so be careful to avoid a scenario where you
// inadvertently clobber your ledger's data!
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

// Invoke function is called on each transaction invoking the chaincode. It
// follows a structure of switching calls, so each valid feature need to
// have a proper entry-point.
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	// extract the function name and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	//implements a switch for each acceptable function
	if fn == "registerCar" {
		//registers a new car into the ledger, starting a new Auction
		return s.registerCar(stub, args)
	
	} else if fn == "registerBid" {
		//registers a new bid, validates it 
		return s.registerBid(stub, args)
	
	} else if fn == "Auctionner" {
		//registers a new bid, validates it 
		return s.Auctionner(stub, args)
	}

	
	//function fn not implemented, notify error
	return shim.Error("Chaincode does not support this function.")
}

/*
	SmartContract::registerCar(...)
	Does the register of a new car into the ledger and starts a new auction.
	The car is the base of the key|value structure.
	The key constitutes the Car ID.
	- args[0] - Car ID
	- args[1] - the starting price and minimun bid value for the auction
*/
func (s *SmartContract) registerCar(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	//validate args vector lenght
	if len(args) != 2 {
		return shim.Error("It was expected the parameters: <Car id> <start price>")
	}

	//gets the parameters associated with the carID and startprice 
	carid := args[0]
	startprice := args[1]

	//creates the car record with the respective startpŕice
	var auctioncar = Car{StartPrice: startprice, IsOpen: true, WinnerID: ""}
	
	//encapsulates car in a JSON structure
	carAsBytes, _ := json.Marshal(auctioncar)
	
	//registers car in the ledger
	stub.PutState(carid, carAsBytes)

	//loging...
	fmt.Println("Registering car: ", carid, "-->", auctioncar)

	//notify procedure success
	return shim.Success(nil)
}

func (s *SmartContract) registerBid(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	//validate args vector lenght
	if len(args) != 4 {
		return shim.Error("It was expected the parameters: <Car id> <Bid id> <bid value> <SSN>")
	}

	//gets the parameters associated with the car ID, bid value and participant information
	carid := args[0]
	bidid := args[1]
	bidvalue := args[2]
	SSN := args [3]

	//loging...
	fmt.Println("Testing args: ", carid, bidid, bidvalue, SSN)

	//RETRIEVES CAR RECORD
	carAsBytes, err := stub.GetState(carid)
	if carAsBytes == nil {
		fmt.Println("There ins't an active auction for the mentioned car")
		return shim.Error("Error on retrieving car ID register")
	}
	//creates Car struct to manipulate returned bytes
	var auctioncar = Car{} 
	//loging...
	//fmt.Println("Retrieving car bytes: ", carAsBytes)
	//convert bytes into a Car object
	json.Unmarshal(carAsBytes, &auctioncar)
	
	
	//RETRIEVES BID RECORD
	bidAsBytes, err := stub.GetState(bidid)
	if err != nil {
		return shim.Error("Error on retrieving bid ID register")
	}
	//creates Bid struct to manipulate returned bytes
	var bid = Bid{}
	//loging...
	//fmt.Println("Retrieving bid bytes: ", bidAsBytes)
	//convert bytes into a Car object
	json.Unmarshal(bidAsBytes, &bid)


	//if there ins't a bid record yet, the function is being called for the first time
	if bidAsBytes == nil {
			
		//loging
		fmt.Println("This is the first bid of the auction")

		//retrieves startprice in the integer type to make comparissons
		startprice, err  := strconv.Atoi(auctioncar.StartPrice)
		if err != nil {
			// handle error
			fmt.Println(err)
			return shim.Error("Error on converting start price from string to int")
		}
		//loging...
		fmt.Println("Retrieving car after unmarshall: ", auctioncar)

		//retrieves participant bid in the integer type to make comparissons
		clientbid, err  := strconv.Atoi(bidvalue)
		if err != nil {
			// handle error
			fmt.Println(err)
			return shim.Error("Error on converting bid value from string to int")
		}

		if clientbid >= startprice {

			//creates the bid record with the respective startpŕice
			bid = Bid{Value: bidvalue, ClientID: SSN}
			
			//encapsulates bid in a JSON structure
			bidAsBytes, _ := json.Marshal(bid)
			
			//registers bid in the ledger
			stub.PutState(bidid, bidAsBytes)
	
			//loging...
			fmt.Println("Registering the first bid: ", bid)
				
	
		} else {
			fmt.Println("Your bid was rejected, minimun bid value is: ", startprice)	
			return shim.Error("Your bid was rejected") 			
		}

	} else {
		
		if auctioncar.IsOpen == false {
			fmt.Println("The auction for the mentioned car has already ended")
			return shim.Error("Auction has ended")
		}

		//creates Bid struct to manipulate returned bytes
		HighestBid := Bid{}
		
		//loging...
		//fmt.Println("Retrieving highest bid bytes: ", bidAsBytes)

		//convert bytes into a Bid object
		json.Unmarshal(bidAsBytes, &HighestBid)

		highestbid, err  := strconv.Atoi(HighestBid.Value)
		if err != nil {
			// handle error
			fmt.Println(err)
			return shim.Error("Error on converting highest bid value from string to int") 
		}
	
	
		//retrieves participant bid in the integer type to make comparissons
		clientbid, err  := strconv.Atoi(bidvalue)
		if err != nil {
			// handle error
			fmt.Println(err)
			return shim.Error("Error on converting bid value from string to int")
		}

		if clientbid > highestbid {

			//creates the bid record with the respective startpŕice
			bid = Bid{Value: bidvalue, ClientID: SSN}
			
			//encapsulates bid in a JSON structure
			bidAsBytes, _ := json.Marshal(bid)
			
			//registers bid in the ledger
			stub.PutState(bidid, bidAsBytes)

			//loging...
			fmt.Println("Registering new highest bid: ", bid)
			

		} else {
			fmt.Println("Your bid was rejected, highest bid value right now is: ", highestbid)
			return shim.Error("Your bid was rejected")
		}
	}

	//notify procedure success
	return shim.Success(nil)
}

func (s *SmartContract) Auctionner(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	//validate args vector lenght
	if len(args) != 3 {
		return shim.Error("It was expected 3 parameters: <carID>, <bidID> <isOpen>")
	}
	carid := args[0]
	bidid := args[1]
	isopen, _:= strconv.ParseBool(args[2])
	var Bidbuffer bytes.Buffer

	//RETRIEVES BID RECORD
	bidAsBytes, err := stub.GetState(bidid)
	if err != nil {
		return shim.Error("Error on retrieving bid ID register")
	}
	//creates Bid struct to manipulate returned bytes
	var bid = Bid{}

	//loging...
	//fmt.Println("Retrieving bid bytes: ", bidAsBytes)

	//convert bytes into a Car object
	json.Unmarshal(bidAsBytes, &bid)
	
	//RETRIEVES CAR RECORD
	carAsBytes, err := stub.GetState(carid)
	if err != nil {
		return shim.Error("Error on retrieving car ID register")
	}

	//creates Bid struct to manipulate returned bytes
	auctioncar := Car{}

	//loging...
	//fmt.Println("Retrieving car bytes: ", carAsBytes)

	//convert bytes into a Car object
	json.Unmarshal(carAsBytes, &auctioncar)


	if isopen == false {
		fmt.Println("The auction has came to an end")

		//creates Bid struct to manipulate returned bytes
		winnerbid := Bid{}

		//convert bytes into a Car object
		json.Unmarshal(bidAsBytes, &winnerbid)
		
		//notify procedure success
		fmt.Println("The auctioned car", carid, " is being sold to participant", winnerbid.ClientID, "for", winnerbid.Value, "dolars")
	
		//creates the bid record with the respective startpŕice
		auctioncar.IsOpen = isopen
		auctioncar.WinnerID = winnerbid.ClientID
		//encapsulates bid in a JSON structure
		carAsBytes, _ = json.Marshal(auctioncar)
	
		//registers bid in the ledger
		stub.PutState(carid, carAsBytes)

		return shim.Success(nil)
	}

	if bidAsBytes != nil{
		bidtime, _ := stub.GetTxTimestamp()
		strbidtime := strconv.Itoa(int(bidtime.Seconds))
		Bidbuffer.WriteString("{")
		Bidbuffer.WriteString("\"Bidtime\":")
		Bidbuffer.WriteString(strbidtime)
		Bidbuffer.WriteString(", \"Bidvalue\":")
		Bidbuffer.WriteString(bid.Value)
		Bidbuffer.WriteString("}")
		return shim.Success(Bidbuffer.Bytes())
	}
	return shim.Success(nil)
}
	


/*
 * The main function starts up the chaincode in the container during instantiate
*/
func main() {

	////////////////////////////////////////////////////////
	// USE THIS BLOCK TO COMPILE THE CHAINCODE
	if err := shim.Start(new(SmartContract)); err != nil {
		fmt.Println("Error starting SmartContract chaincode: %s\n", err)
	}
	////////////////////////////////////////////////////////
}