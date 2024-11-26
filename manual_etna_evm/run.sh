#!/bin/bash

set -euo pipefail

echo -e "\n🔑 Generating keys\n"
go run ./01_generate_keys/

echo -e "\n💰 Checking balance\n" 
go run ./02_check_balance/

echo -e "\n🕸️  Creating subnet\n"
go run ./03_create_subnet/

echo -e "\n🧱 Generating genesis\n"
go run ./04_L1_genesis/

echo -e "\n⛓️  Creating chain\n"
go run ./05_create_chain/

echo -e "\n🚀 Launching nodes\n"
./06_launch_nodes/launch.sh

echo -e "\n🔄 Converting chain\n"
go run ./07_convert_chain/

echo -e "\n🔃 Restarting nodes\n"
./06_launch_nodes/launch.sh

echo -e "\n🏥 Checking subnet health\n"
go run ./09_check_subnet_health/

echo -e "\n💸 Sending some test coins\n"
go run ./10_evm_transfer/

# echo -e "\n🔄 Waiting for the transaction to be included\n"
# sleep 30

# echo -e "\n🔄 Initializing PoA\n"
# go run ./13_init_poa/
