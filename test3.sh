
echo
echo "##########################################################"
echo "################### chaincode install ####################"
echo "##########################################################"
echo
sleep 2
peer chaincode install -p chaincodedev/chaincode/contract -n mycc -v 0

echo
echo "##########################################################"
echo "############### chaincode instantiate ####################"
echo "##########################################################"
echo
sleep 2
peer chaincode instantiate -n mycc -v 0 -c '{"Args":["Init"]}' -C myc

sleep 3

cX="b72b32932795f37880940dcc172e206ec1ad113a171d600df8998ed87c120dc5"
cY="f6a60d055f03b65a353e8962706b9e3615cb312ccf81f4bc98e7c1166306f8bf"

echo
echo "##########################################################"
echo "##################### client regist ######################"
echo "##########################################################"
echo
sleep 2
peer chaincode invoke -n mycc -v 1.0 -c '{"Args":["clientRegist", "1", "'${cX}'", "'${cY}'"]}' -C myc

sleep 3

iX="70abfb4f27f570df417234600aa29303def0fadeff2bf4c36ca52800ae0dd418"
iY="3e0387b4c93d2c4ef9e103f2006e15b5950370a5b977202f45cd375566b0292c"

echo
echo "##########################################################"
echo "############## insurance company regist ##################"
echo "##########################################################"
echo
sleep 2
peer chaincode invoke -n mycc -v 1.0 -c '{"Args":["insurCompanyRegist", "1", "'${iX}'", "'${iY}'"]}' -C myc

sleep 3

echo
echo "##########################################################"
echo "################### check_buyTicket ######################"
echo "##########################################################"
echo
sleep 2
boolValue=`peer chaincode query -n mycc -v 1.0 -c '{"Args":["check_buyTicket"]}' -C myc`
echo $boolValue

if [ "$boolValue" = "Query Result: false" ] ; then
  echo "Ticket has not been bought! EXIT!"
  exit 1
else
  act1="{'Args':['buyTicket','1','1','1','11:00']}"
  r1="f4ad4aac6c5b78405d69b1dcbdb75ee9ef21c68e6342f37305ddb6d2742f968c"
  s1="45fd34af8d2f93131bbfa42bb7423844c6dc9ff3365f4590ac4e979c64650a9d"

  echo
  echo "##########################################################"
  echo "####################### buyTicket ########################"
  echo "##########################################################"
  echo
  sleep 2
  peer chaincode invoke -n mycc  -v 1.0 -c '{"Args":["buyTicket","1","1","1","11:00","'${act1}'","'${r1}'","'${s1}'"]}' -C myc

  sleep 3

  echo
  echo "##########################################################"
  echo "###################### queryOrder ########################"
  echo "##########################################################"
  echo
  sleep 2
  peer chaincode query -n mycc  -v 1.0 -c '{"Args":["queryByObjectType","order"]}' -C myc

  echo
  echo "##########################################################"
  echo "##################### initPolicy #########################"
  echo "##########################################################"
  echo
  sleep 2
  peer chaincode invoke -n mycc  -v 1.0 -c '{"Args":["initPolicy","1"]}' -C myc

  sleep 3

  act2="{'Args':['clientDeposit','1','1','1','1','12:00']}"
  r2="2da67e2fe36b96be6f042c3e308396caaa6f83058e0a344f57a3bf725cd93aac"
  s2="58a15127bd9ce4d051e89a54be12d4ccfb80fdff16d34b8ebdc85d3da56348bf"

  echo
  echo "##########################################################"
  echo "#################### clientDeposit #######################"
  echo "##########################################################"
  echo
  sleep 2
  peer chaincode invoke -n mycc  -v 1.0 -c '{"Args":["clientDeposit","1","1","1","1","12:00","'${act2}'","'${r2}'","'${s2}'"]}' -C myc

  sleep 3

  echo
  echo "##########################################################"
  echo "######################## timeOut #########################"
  echo "##########################################################"
  echo
  sleep 2
  peer chaincode invoke -n mycc  -v 1.0 -c '{"Args":["timeOut","11"]}' -C myc

  sleep 1
  echo "Insurance company didn't deposit during the required time!"

  act3="{'Args':['clientRefund','1','1']}"
  r3="3b3fb6716cb7089a4edbb5a24e00a8d5b96ace91df54fa647614e52674279698"
  s3="7a7a02e6869e4c7817de567c4795f9e6f85c964c5edf9dcb665924e11cceff35"

  echo
  echo "##########################################################"
  echo "####################### clientRefund #####################"
  echo "##########################################################"
  echo
  sleep 2
  peer chaincode invoke -n mycc  -v 1.0 -c '{"Args":["clientRefund","1","1","'${act3}'","'${r3}'","'${s3}'"]}' -C myc

  sleep 1

  exit 1

  echo
  echo "##########################################################"
  echo "##################### queryPolicy ########################"
  echo "##########################################################"
  echo
  sleep 2
  peer chaincode query -n mycc -v 1.0 -c '{"Args":["queryByObjectType","policy"]}' -C myc

  Act1="{'Args':['insurCompanyDeposit','1','1','1','1','12:00']}"
  R1="ab2dedda006dbfe57c10e900ddefa6ce59aa9896b9545884c2a781f85841f50f"
  S1="c34eb1003aebc60c6068867bb1c91c7a693790fb2d2f6a9918276c80ad5cb6cd"

  echo
  echo "##########################################################"
  echo "################# insurCompanyDeposit ####################"
  echo "##########################################################"
  echo
  sleep 2
  peer chaincode invoke -n mycc -v 1.0 -c '{"Args":["insurCompanyDeposit","1","1","1","1","12:00","'${Act1}'","'${R1}'","'${S1}'"]}' -C myc
  sleep 3

  echo
  echo "##########################################################"
  echo "##################### queryPolicy ########################"
  echo "##########################################################"
  echo
  sleep 2
  peer chaincode query -n mycc  -v 1.0 -c '{"Args":["queryByObjectType","policy"]}' -C myc

  echo
  echo "##########################################################"
  echo "################## check_flightDelay #####################"
  echo "##########################################################"
  echo
  sleep 2
  boolValue=`peer chaincode query -n mycc  -v 1.0 -c '{"Args":["check_flightDelay"]}' -C myc`
  echo $boolValue

  if [ "$boolValue" = "Query Result: true" ] ; then
    echo
    echo "##########################################################"
    echo "###################### flightDelay #######################"
    echo "##########################################################"
    echo
    sleep 2

    peer chaincode invoke -n mycc  -v 1.0 -c '{"Args":["flightDelay","1","1","DELAY"]}' -C myc

    sleep 3

    act4="{'Args':['compensate','1','1']}"
    r4="94ed96c29fdc2d9878c7086fcbd9f2746ff5010f9a657cb14ebf7c23134e1df7"
    s4="88026c3364b7c14c072654710fb6e9a1a9a8155a397f85ae9edc1a0c30b7c68d"

    echo
    echo "##########################################################"
    echo "###################### compensate ########################"
    echo "##########################################################"
    echo
    sleep 2

    peer chaincode invoke -n mycc  -v 1.0 -c '{"Args":["compensate","1","1","'${act4}'","'${r4}'","'${s4}'"]}' -C myc

    sleep 3

    echo "Compensate successfully"
    #exit 1
  else

    Act2="{'Args':['insurCompanyRefund','1','1']}"
    R2="6d3d277c774632f5fa699762336a1581ebf54b823543ed7408a1add9cc4cea8"
    S2="c290ecdf5f97bbe412ed60c5d167790abf81662f07333240928b56955b37c558"

    echo
    echo "##########################################################"
    echo "############### insurance company refund #################"
    echo "##########################################################"
    echo
    sleep 2

    peer chaincode invoke -n mycc  -v 1.0 -c '{"Args":["insurCompanyRefund","1","1","'${Act2}'","'${R2}'","'${S2}'"]}' -C myc

    sleep 3

    echo "Refund successfully"
    #exit 1
  fi

  echo
  echo "##########################################################"
  echo "###################### queryOrder ########################"
  echo "##########################################################"
  echo
  sleep 2
  peer chaincode query -n mycc  -v 1.0 -c '{"Args":["queryByObjectType","order"]}' -C myc
fi
