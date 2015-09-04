package core

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/eris-ltd/eth-client/client"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

// All ethereum transactions have a common form
// NOTE: this struct isn't exported from go-ethereum/core/types :(
type Transaction struct {
	Nonce           uint64
	Price, GasLimit *big.Int
	Recipient       *common.Address `rlp:"nil"` // nil means contract creation
	Amount          *big.Int
	Data            []byte
	V               byte     // signature
	R, S            *big.Int // signature

	from *common.Address
}

func (tx *Transaction) String() string {
	var rec []byte
	if tx.Recipient != nil {
		rec = tx.Recipient.Bytes()
	}
	return fmt.Sprintf(`
	Nonce: %d,
	To: %x,
	Amount: %x,
	GasLimit: %x,
	GasPrice: %x,
	Data: %x
`, tx.Nonce, rec, tx.Amount.Bytes(), tx.GasLimit.Bytes(), tx.Price.Bytes(), tx.Data)
}

func NewTransaction(to, from *common.Address, nonce uint64, amt, gas, price *big.Int, data []byte) *Transaction {
	if len(data) > 0 {
		data = common.CopyBytes(data)
	}
	tx := &Transaction{
		Nonce:     nonce,
		Recipient: to,
		Data:      data,
		Amount:    new(big.Int),
		GasLimit:  new(big.Int),
		Price:     new(big.Int),
		R:         new(big.Int),
		S:         new(big.Int),
		from:      from,
	}
	if amt != nil {
		tx.Amount.Set(amt)
	}
	if gas != nil {
		tx.GasLimit.Set(gas)
	}
	if price != nil {
		tx.Price.Set(price)
	}
	return tx
}

// rlp encode and hash
func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

// Hash of the transaction for signing
func (tx *Transaction) SignBytes() []byte {
	h := rlpHash([]interface{}{
		tx.Nonce,
		tx.Price,
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.Data,
	})
	return h[:]
}

// Apply the signature to the transaction
func (tx *Transaction) ApplySignature(sig [65]byte) {
	tx.R = new(big.Int).SetBytes(sig[:32])
	tx.S = new(big.Int).SetBytes(sig[32:64])
	tx.V = sig[64] + 27
}

// Sign the transaction using the keys server
func (tx *Transaction) Sign(signAddr string) error {
	if tx.from == nil {
		return fmt.Errorf("from is not set")
	}
	signBytes := fmt.Sprintf("%X", tx.SignBytes())
	addrHex := fmt.Sprintf("%X", tx.from.Bytes())
	sig, err := Sign(signBytes, addrHex, signAddr)
	if err != nil {
		return err
	}
	tx.ApplySignature(sig)
	return nil
}

//------------------------------------------------------------------------------------
// core functions with string args.
// validates strings and forms transaction

func Send(nodeAddr, fromAddr, toAddr, amtS, gasS, priceS string, nonce uint64) (*Transaction, error) {
	from, nonce, amt, gas, price, err := checkCommon(nodeAddr, fromAddr, amtS, gasS, priceS, nonce)
	if err != nil {
		return nil, err
	}

	if toAddr == "" {
		return nil, fmt.Errorf("destination address must be given with --to flag")
	}

	toAddrBytes, err := hex.DecodeString(toAddr)
	if err != nil {
		return nil, fmt.Errorf("toAddr is bad hex: %v", err)
	}
	to := common.BytesToAddress(toAddrBytes)

	return NewTransaction(&to, &from, nonce, amt, gas, price, nil), nil
}

func Create(nodeAddr, fromAddr, amtS, gasS, priceS, data string, nonce uint64) (*Transaction, error) {
	from, nonce, amt, gas, price, err := checkCommon(nodeAddr, fromAddr, amtS, gasS, priceS, nonce)
	if err != nil {
		return nil, err
	}

	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("data is bad hex: %s", data)
	}
	return NewTransaction(nil, &from, nonce, amt, gas, price, dataBytes), nil
}

func Call(nodeAddr, fromAddr, toAddr, amtS, gasS, priceS, data string, nonce uint64) (*Transaction, error) {
	from, nonce, amt, gas, price, err := checkCommon(nodeAddr, fromAddr, amtS, gasS, priceS, nonce)
	if err != nil {
		return nil, err
	}

	if toAddr == "" {
		return nil, fmt.Errorf("destination address must be given with --to flag")
	}

	toAddrBytes, err := hex.DecodeString(toAddr)
	if err != nil {
		return nil, fmt.Errorf("toAddr is bad hex: %v", err)
	}
	to := common.BytesToAddress(toAddrBytes)

	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("data is bad hex: %s", data)
	}

	return NewTransaction(&to, &from, nonce, amt, gas, price, dataBytes), nil
}

//------------------------------------------------------------------------------------
// sign and broadcast

func Sign(signBytes, signAddr, signRPC string) (sig [65]byte, err error) {
	args := map[string]string{
		"hash": signBytes,
		"addr": signAddr,
	}
	b, err := json.Marshal(args)
	if err != nil {
		return
	}
	logger.Debugln("Sending request body:", string(b))
	req, err := http.NewRequest("POST", signRPC+"/sign", bytes.NewBuffer(b))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	sigS, errS, err := requestResponse(req)
	if err != nil {
		return sig, fmt.Errorf("Error calling signing daemon: %s", err.Error())
	}
	if errS != "" {
		return sig, fmt.Errorf("Error (string) calling signing daemon: %s", errS)
	}
	sigBytes, err := hex.DecodeString(sigS)
	if err != nil {
		err = fmt.Errorf("sig is bad hex:", err)
		return
	}
	copy(sig[:], sigBytes)
	return
}

func Broadcast(tx *Transaction, broadcastRPC string) (interface{}, error) {
	w := new(bytes.Buffer)
	if err := rlp.Encode(w, tx); err != nil {
		return nil, err
	}
	fmt.Println("Tx Serialized")
	txHex := fmt.Sprintf("%X", w.Bytes())
	fmt.Println(txHex)
	r, err := client.RequestResponse("eth", "sendRawTransaction", txHex)
	if err != nil {
		return nil, err
	}
	return r, nil
}

//------------------------------------------------------------------------------------
// utils for talking to the key server

type HTTPResponse struct {
	Response string
	Error    string
}

func requestResponse(req *http.Request) (string, string, error) {
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode >= 400 {
		return "", "", fmt.Errorf(resp.Status)
	}
	return unpackResponse(resp)
}

func unpackResponse(resp *http.Response) (string, string, error) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	r := new(HTTPResponse)
	if err := json.Unmarshal(b, r); err != nil {
		return "", "", err
	}
	return r.Response, r.Error, nil
}

//------------------------------------------------------------------------------------
// sign and broadcast convenience
type TxResult struct {
	BlockHash []byte // all txs get in a block
	Hash      []byte // all txs get a hash

	// only CallTx
	Address   []byte // only for new contracts
	Return    []byte
	Exception string

	//TODO: make Broadcast() errors more responsive so we
	// can differentiate mempool errors from other
}

func SignAndBroadcast(nodeAddr, signAddr string, tx *Transaction, sign, broadcast, wait bool) (txid string, err error) {
	if sign {
		if err = tx.Sign(signAddr); err != nil {
			return
		}
	}

	if broadcast {
		/*
			if wait {
				var ch chan Msg
				ch, err = subscribeAndWait(tx, chainID, nodeAddr, inputAddr)
				if err != nil {
					return nil, err
				} else {
					defer func() {
						if err != nil {
							// if broadcast threw an error, just return
							return
						}
						logger.Debugln("Waiting for tx to be committed ...")
						msg := <-ch
						if msg.Error != nil {
							logger.Infof("Encountered error waiting for event: %v\n", msg.Error)
							err = msg.Error
						} else {
							txResult.BlockHash = msg.BlockHash
							txResult.Return = msg.Value
							txResult.Exception = msg.Exception
						}
					}()
				}
			}*/
		var r interface{}
		r, err = Broadcast(tx, nodeAddr)
		if err != nil {
			return "", err
		}
		return r.(string), nil
		/*
			txResult = &TxResult{
				Hash: receipt.TxHash,
			}
			if tx_, ok := tx.(*types.CallTx); ok {
				if len(tx_.Address) == 0 {
					txResult.Address = types.NewContractAddress(tx_.Input.Address, tx_.Input.Sequence)
				}
			}*/
	}
	return
}

//------------------------------------------------------------------------------------
// convenience function

func stringToBig(s string) (*big.Int, error) {
	if len(s) > 2 && s[:2] == "0x" {
		s = s[2:]
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(b), nil
}

// if the nonce is given, the nodeAddr and addr are not needed
func checkCommon(nodeAddr, addr, amtS, gasS, priceS string, seq uint64) (from common.Address, nonce uint64, amount, gas, price *big.Int, err error) {
	// resolve the big ints
	if amount, err = stringToBig(amtS); err != nil {
		err = fmt.Errorf("amt %s is bad hex: %v", amtS, err)
		return
	}
	if gas, err = stringToBig(gasS); err != nil {
		err = fmt.Errorf("gas %s is bad hex: %v", gasS, err)
		return
	}
	if price, err = stringToBig(priceS); err != nil {
		err = fmt.Errorf("price %s is bad hex: %v", priceS, err)
		return
	}

	// resolve the address
	if addr == "" {
		err = fmt.Errorf("--addr must be given")
		return
	}
	var addrBytes []byte
	addrBytes, err = hex.DecodeString(addr)
	if err != nil {
		err = fmt.Errorf("addr is bad hex: %v", err)
		return
	}
	from = common.BytesToAddress(addrBytes)

	// resolve the nonce (or fetch it)
	if seq == 0 {
		if nodeAddr == "" {
			err = fmt.Errorf("input must specify a nonce with the --nonce flag or use --node-addr (or MINTX_NODE_ADDR) to fetch the nonce from a node")
			return
		}

		var r interface{}
		// fetch block num
		r, err = client.RequestResponse("eth", "blockNumber")
		if err != nil {
			err = fmt.Errorf("Error fetching block number", err)
			return
		}
		// NOTE: both block num and account nonces are hex. (why?!)
		blockNum := client.HexToInt(r.(string))

		r, err = client.RequestResponse("eth", "getTransactionCount", addr, blockNum)
		if err != nil {
			err = fmt.Errorf("Error fetching account nonce", err)
			return
		}

		nonce = uint64(client.HexToInt(r.(string))) + 1
	} else {
		nonce = seq
	}

	return
}
