# Blockchain-Auction
A simulation project of a car auction using blockchain network

Coordination:
 Wilson S. Melo Jr. (wsjunior@inmetro.gov.br)
 
Revised on February 14th, 2022.

# Table of Contents
1. [Introduction](#introduction)
2. [Requirements](#requirements)
3. [Usage](#usage)

# Introduction

An Auctioneer starts the auction of a car announcing it's start price and starts counting the duration of the auction.

Then the participants start calling out their bids.

There is an asset for the car and another for the bids.

The car asset only register the car ID, the start price, the winner ID and a bool for the auction existence.

The bids asset has a unique ID for all bids, it keeps the bid value and the Client ID, which is the Social Security Number of the respective participant.

 # Requirements
 
 First of all, this project uses the Inter-NMI experiment Blockchain Network, so, you should start by reading and getting familiarized with their work. 
 ```
 cd ~
 git clone https://github.com/wsmelojr/nmiblocknet
 cd nmiblocknet
 ```
 Follow their tutorial on how to set up and work with their network in: https://github.com/wsmelojr/nmiblocknet.

 # Usage
 In order to start working with the auction project we need to start up the blockchain network:
 ```
 docker-compose -f peer-orderer.yaml -f peer-inmetro.yaml up -d
 ```
 Create the channel and join the peers
 
 ``` 
 ./configchannel.sh inmetro.br -c
 ```
 Create chaincode and client aplication directories (just for the first time)
 ```
 cd ~/nmiblocknet
 mkdir auction
 mkdir auction-cli
 cp fabpki-cli/inmetro.br.json auction-cli
 mv ~/Blockchain-Auction/auction.go auction
 mv ~/Blockchain-Auction/auctionner.py auction-cli
 mv ~/Blockchain-Auction/bid_register.py auction-cli
 rm -rf ~/Blockchain-Auction
 
 ```
 The inmetro.br.json file, now in auction-cli directory, is a copy of the nmiblocknet file with
exact same name in their fabpki project. In order to make it work in the auction project it should
receive the exact same changes described in nmiblocknet repository, if you copied the file without
making the necessary changes, do it now, otherwise it won't work.

 Install and Instantiate auction chaincode
 ```
 ./configchaincode.sh install cli0 auction 1.0
 ./configchaincode.sh instantiate cli0 auction 1.0
 
 ```
 Now you can play! In order to start the auction:
 
 ```
 cd auction-cli
 python3 auctionner.py <license plate> <start value(int)>
 ```
 
 Open another terminal to call out the bids:
 
 ```
 python3 bid_register.py <license plate> <bid value(int)> <participant SSN>
 ```
 Using CouchDB
 ```
 docker ps | grep couchdb
 ```
 You can check the local copy of the ledger register in http://localhost:7984/_utils
