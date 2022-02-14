import sys
from hfc.fabric import Client as client_fabric
import asyncio
import time
import json 

domain = "inmetro.br" 
channel_name = "nmi-channel"
cc_name = "auction"
cc_version = "1.0"

if __name__ == "__main__":

    #test if the license plate and start price were informed as argument
    if len(sys.argv) != 3:
        print("Usage:",sys.argv[0],"<license plate>", "<start price>")
        exit(1)

    #get the carID and the start price
    car_id = sys.argv[1]
    start_price = sys.argv[2]

    #feedback to the user
    print("Generating car asset",  car_id)
    
    #creating file to save starting price for the corresponding car
    start_price_file = car_id + ".val"

    #write starting price in its file 
    with open(start_price_file, "wt") as f:
        f.write(start_price)

    #shows license plate and start price save file
    print("The start price was saved into", start_price_file)

    #shows the car start price
    print("Starting auction with the following bid value:\n", start_price)

    #creates a loop object to manage async transactions
    loop = asyncio.get_event_loop()

    #instantiate the hyperledeger fabric client
    c_hlf = client_fabric(net_profile=(domain + ".json"))

    #get access to Fabric as Admin user
    admin = c_hlf.get_user(domain, 'Admin')
    callpeer = "peer0." + domain

    #query peer installed chaincodes, make sure the chaincode is installed
    print("Checking if the chaincode auction is properly installed:")
    response = loop.run_until_complete(c_hlf.query_installed_chaincodes(
        requestor=admin,
        peers=[callpeer]
    ))
    print(response)

    #the Fabric Python SDK do not read the channel configuration, we need to add it mannually'''
    c_hlf.new_channel(channel_name)

    #invoke the chaincode to register the car
    response = loop.run_until_complete(c_hlf.chaincode_invoke(
        requestor=admin, 
        channel_name=channel_name, 
        peers=[callpeer],
        cc_name=cc_name, 
        cc_version=cc_version,
        fcn='registerCar', 
        args=[car_id, start_price], 
        cc_pattern=None))
    time.sleep(5)
    #so far, so good
    print("Success on registering car and start price!")
    
    #Generates BidId for respective car auction and opens auction
    bid_id = car_id + "-B" 
    is_open = "True"
    start_time = int(time.time()) 
    highest_bid_value = 0
    while True:
        current_time = int(time.time())
        if current_time < start_time + 120:
            #invoke the chaincode to register info
            response = loop.run_until_complete(c_hlf.chaincode_invoke(
                requestor=admin, 
                channel_name=channel_name, 
                peers=[callpeer],
                cc_name=cc_name, 
                cc_version=cc_version,
                fcn='Auctionner', 
                args=[car_id, bid_id, is_open], 
                cc_pattern=None))
            
            if response != "":
                #Geting Bid Timestamp
                print("The last bid ocurred in:\n", response)
                #Retrieving last bid
                bid_info = json.loads(response)
                bid_time = bid_info.get('Bidtime')
                print("Bidtime: ", bid_time)
                bid_value = bid_info.get('Bidvalue')
                print("Bidvalue: ", bid_value)
                if bid_value == highest_bid_value:
                    print("One minute passed without new bids, the car is being sold, check log to see the winner")
                    is_open = "False"
                    response = loop.run_until_complete(c_hlf.chaincode_invoke(
                    requestor=admin, 
                    channel_name=channel_name, 
                    peers=[callpeer],
                    cc_name=cc_name, 
                    cc_version=cc_version,
                    fcn='Auctionner', 
                    args=[car_id, bid_id, is_open], 
                    cc_pattern=None))
                    exit(1)
                else:
                    highest_bid_value = bid_value
                time.sleep(60)
            
            else:
                print("The auction started in:\n", start_time)
                highest_bid_value = 0
                time.sleep(10)
        else:
            print("Auction has reached it's time limit, check log to see the winner")
            is_open = "False"
            #invoke the chaincode to register info
            response = loop.run_until_complete(c_hlf.chaincode_invoke(
                requestor=admin, 
                channel_name=channel_name, 
                peers=[callpeer],
                cc_name=cc_name, 
                cc_version=cc_version,
                fcn='Auctionner', 
                args=[car_id, bid_id, is_open], 
                cc_pattern=None))
            exit(1)