
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
boolValue=`peer chaincode query -n mycc -v 1.0 -c '{"Args":["getRandBool"]}' -C myc`
echo $boolValue

if [ "$boolValue" = "Query Result: false" ] ; then
  echo "Ticket has not been bought! EXIT!"
  exit 1
else
  act1="{'Args':['writeOrder','1','1','1','11:00']}"
  r1="1a226fad62b8bb71681ece92b75daf512b0a08776493111074d77946a94020f3"
  s1="f8d4d5314733be75fd4a30a49a4e8ae6f2229b97a2a52a91875e3634b58461da"

  echo
  echo "##########################################################"
  echo "###################### writeOrder ########################"
  echo "##########################################################"
  echo
  sleep 2
  peer chaincode invoke -n mycc  -v 1.0 -c '{"Args":["writeOrder","1","1","1","11:00","'${act1}'","'${r1}'","'${s1}'"]}' -C myc

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
  echo "################### initGuarantee ########################"
  echo "##########################################################"
  echo
  sleep 2
  peer chaincode invoke -n mycc  -v 1.0 -c '{"Args":["initGuarantee","1"]}' -C myc

  sleep 3

  #echo
  #echo "##########################################################"
  #echo "######################## timeOut #########################"
  #echo "##########################################################"
  #echo
  #sleep 2
  #peer chaincode invoke -n mycc  -v 1.0 -c '{"Args":["timeOut","11"]}' -C myc

  #sleep 1
  #echo "Client didn't deposit during the required time!"

  #exit 1

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

  #echo
  #echo "##########################################################"
  #echo "######################## timeOut #########################"
  #echo "##########################################################"
  #echo
  #sleep 2
  #peer chaincode invoke -n mycc  -v 1.0 -c '{"Args":["timeOut","11"]}' -C myc

  #sleep 1
  #echo "Insurance company didn't deposit during the required time!"

  #echo
  #echo "##########################################################"
  #echo "######################## withdraw ########################"
  #echo "##########################################################"
  #echo
  #sleep 2
  #peer chaincode invoke -n mycc  -v 1.0 -c '{"Args":["clientWithdraw","11"]}' -C myc

  #sleep 1

  #exit 1

  echo
  echo "##########################################################"
  echo "################### queryGuarantee #######################"
  echo "##########################################################"
  echo
  sleep 2
  peer chaincode query -n mycc -v 1.0 -c '{"Args":["queryByObjectType","guarantee"]}' -C myc

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
  echo "#################### queryGuarantee #######################"
  echo "##########################################################"
  echo
  sleep 2
  peer chaincode query -n mycc  -v 1.0 -c '{"Args":["queryByObjectType","guarantee"]}' -C myc

  echo
  echo "##########################################################"
  echo "###################### check_delay #######################"
  echo "##########################################################"
  echo
  sleep 2
  boolValue=`peer chaincode query -n mycc  -v 1.0 -c '{"Args":["getRandBool"]}' -C myc`
  echo $boolValue

  if [ "$boolValue" = "Query Result: true" ] ; then
    echo
    echo "##########################################################"
    echo "###################### changeOrder #######################"
    echo "##########################################################"
    echo
    sleep 2

    peer chaincode invoke -n mycc  -v 1.0 -c '{"Args":["flightDelay","1","1","DELAY"]}' -C myc

    sleep 3
    act4="{'Args':['indemnify','1','1']}"
    r4="93daaf88c6b9af4ba2b18abe795b59f6433d6e359be69e500472b8e8c56941e2"
    s4="19e9cc00056279b2a2b046b33591e0989e3316938d5318bb800754fc1d91191f"
    echo
    echo "##########################################################"
    echo "####################### indemnify ########################"
    echo "##########################################################"
    echo
    sleep 2

    peer chaincode invoke -n mycc  -v 1.0 -c '{"Args":["indemnify","1","1","'${act4}'","'${r4}'","'${s4}'"]}' -C myc

    sleep 3

    echo "Indemnify successfully"
    exit 1
  else

    Act2="{'Args':['insurCompanyWithdraw','1','1']}"
    R2="d65dc39e946181322baeef3bd1e08689cd7d383180dafee0074eae88ef714b20"
    S2="d4f8acef0fb9279c123b449f390eb724a1cf4c3084a34d15c2872bf56eaec958"

    echo
    echo "##########################################################"
    echo "############# insurance company withdraw #################"
    echo "##########################################################"
    echo
    sleep 2

    peer chaincode invoke -n mycc  -v 1.0 -c '{"Args":["insurCompanyWithdraw","1","1","'${Act2}'","'${R2}'","'${S2}'"]}' -C myc

    sleep 3

    echo "Withdraw successfully"
    exit 1
  fi

  echo
  echo "##########################################################"
  echo "###################### queryOrder ########################"
  echo "##########################################################"
  echo
  sleep 2
  peer chaincode query -n mycc  -v 1.0 -c '{"Args":["queryByObjectType","order"]}' -C myc
fi
