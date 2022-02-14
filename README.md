# Blockchain-Auction
A project of a car Auction using blockchain network
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
 
 First of all, this project uses the Inter-NMI experiment Blokchain Network, so, you should start by reading and getting familiarized with their work.
 Follow their tutorial on how to setting up and working with their network in: https://github.com/wsmelojr/nmiblocknet.
