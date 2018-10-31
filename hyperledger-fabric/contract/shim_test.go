package contract

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type auditChainCode struct{}

func (cc *auditChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (cc *auditChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Invoke TestChainCode")
	return shim.Success(nil)
}

func TestAuditPutRecords(t *testing.T) {

	containerServiceStub := shim.NewMockStub("AuditSaft", &flowcc)
	fmt.Printf("contractname %s\n", flowcc.ContractName)

	//put_records
	args := make([][]byte, 0)
	txn := "com.auditsaft.tibco.put_records"
	args = append(args, []byte(txn))

	records := `[
			{
				"user_txn_id": "txn1",
	  			"hash_type": "hash",
	  			"hash_value": "agcv"
			},
			{
				"user_txn_id": "txn2",
	  			"hash_type": "hash",
	  			"hash_value": "agcv"
			}
		]`

	args = append(args, []byte(records))
	containerServiceStub.ChannelID = "auditsaft"

	txnId := "auditsaft"
	containerServiceStub.MockInit("txninit", [][]byte{})

	containerServiceStub.MockInvoke(txnId, args)

	containerServiceStub.MockTransactionStart(txn)

	//get_records
	reckeys := `[
		{
			"txn_id": "auditsaft",
			"sequence": 0
		},
		{
			"txn_id": "auditsaft",
			"sequence": 1
		}]`
	fmt.Printf("keys = %s\n", reckeys)

	args = make([][]byte, 0)
	args = append(args, []byte("com.auditsaft.tibco.get_records"))
	args = append(args, []byte(reckeys))
	resp := containerServiceStub.MockInvoke(txnId, args)
	fmt.Printf("resp=%s\n", string(resp.Payload))

	//query by txnId
	txn_id := "auditsaft"

	args = make([][]byte, 0)
	args = append(args, []byte("com.auditsaft.tibco.query_by_txn_id"))
	args = append(args, []byte(txn_id))
	resp = containerServiceStub.MockInvoke(txnId, args)
	fmt.Printf("query_by_txn_id resp=%s\n", string(resp.Payload))

	//query by user txnId
	user_txn_id := "txn1"

	args = make([][]byte, 0)
	args = append(args, []byte("com.auditsaft.tibco.query_by_user_txn_id"))
	args = append(args, []byte(user_txn_id))
	resp = containerServiceStub.MockInvoke(txnId, args)
	fmt.Printf("query_by__user_txn_id resp=%s\n", string(resp.Payload))

	//del_records
	reckeys = `[
		{
			"txn_id": "auditsaft",
			"sequence": 0
		}]`

	args = make([][]byte, 0)
	args = append(args, []byte("com.auditsaft.tibco.del_records"))
	args = append(args, []byte(reckeys))
	resp = containerServiceStub.MockInvoke(txnId, args)
	fmt.Printf("del_records resp=%s\n", string(resp.Payload))
	//del_records again
	resp = containerServiceStub.MockInvoke(txnId, args)
	fmt.Printf("del_records again resp=%s\n", string(resp.Payload))

	//query_by_hash_value
	hashv := "agcv"

	args = make([][]byte, 0)
	args = append(args, []byte("com.auditsaft.tibco.query_by_hash_value"))
	args = append(args, []byte(hashv))
	resp = containerServiceStub.MockInvoke(txnId, args)
	fmt.Printf("query_by_hash_value resp=%s\n", string(resp.Payload))

	//del_by_txn_id
	txn_id = "auditsaft"

	args = make([][]byte, 0)
	args = append(args, []byte("com.auditsaft.tibco.del_by_txn_id"))
	args = append(args, []byte(txn_id))
	resp = containerServiceStub.MockInvoke(txnId, args)
	fmt.Printf("del_by_txn_id resp=%s\n", string(resp.Payload))

	//query_by_hash_value
	hashv = "agcv"

	args = make([][]byte, 0)
	args = append(args, []byte("com.auditsaft.tibco.query_by_hash_value"))
	args = append(args, []byte(hashv))
	resp = containerServiceStub.MockInvoke(txnId, args)
	fmt.Printf("query_by_hash_value resp=%s\n", string(resp.Payload))
}
