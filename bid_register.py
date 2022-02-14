import sys
from hfc.fabric import Client as client_fabric
import asyncio

domain = "inmetro.br" 
channel_name = "nmi-channel"
cc_name = "auction"
cc_version = "1.0"

if __name__ == "__main__":
    #test if the license plate and start bid value were informed as argument
    if len(sys.argv) != 4:
        print("Usage:",sys.argv[0],"<car license plate>", "<your bid>", "<your Social Security Number>")
        exit(1)
    
    car_id = sys.argv[1]
    bid_value = sys.argv[2]
    cli_ssn = sys.argv[3]

    bid_id = car_id + "-B"
    
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
        fcn='registerBid', 
        args=[car_id, bid_id, bid_value, cli_ssn], 
        cc_pattern=None))

    #so far, so good
    print("Your bid will be analized!")