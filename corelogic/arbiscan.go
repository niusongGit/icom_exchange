package corelogic

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/prestonTao/keystore/crypto"
	"gitee.com/prestonTao/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/syndtr/goleveldb/leveldb/util"
	"icom_exchange/config"
	"icom_exchange/db"
	"icom_exchange/http_client"
	"icom_exchange/logger"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var WaitTransactionKey = "WaitTransactionKey_%s"
var WaitTransactionKeyPrefix = "WaitTransactionKey_"
var TransactionHashKey = "TransactionHashKey_%s"
var TransactionHashKeyPrefix = "TransactionHashKey_"
var TransactionHashesByBlockNumber = "TransactionHashesByBlockNumber_%d_%s"
var TransactionHashesByBlockNumberPrefix = "TransactionHashesByBlockNumber_"
var TransactionHashesLastBlockNumber = "TransactionHashesLastBlockNumber"
var IComTxHashIndexTransactionHash = "IComTxHashIndexTransactionHash_%s_%s"
var IComTxHashIndexTransactionHashPrefix = "IComTxHashIndexTransactionHash_"
var IAddressIndexTransactionHash = "IAddressIndexTransactionHash_%s_%s"
var IAddressIndexTransactionHashPrefix = "IAddressIndexTransactionHash_"

// APIResponse 结构体用于解析API响应
type APIResponse struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Result  []map[string]interface{} `json:"result"`
	Log     int                      `json:"log"`
}

// Transaction 结构体用于解析交易数据
type Transaction struct {
	TransactionHash string `json:"transactionHash"`
	Sender          string `json:"sender"`
	Token           string `json:"token"`
	BlockNumber     uint64 `json:"blockNumber"`
	Amount          uint64 `json:"amount"`
	ToAmount        uint64 `json:"toAmount"`
	Rate            uint64 `json:"rate"`
	IAddress        string `json:"iAddress"`
	SrcIAddress     string `json:"srcIAddress"`
	TxHash          string `json:"txHash"`
	Status          int64  `json:"status"`       // -1未转账 0转账中 1未上链 2成功上链 3上链失败
	PushIComTime    int64  `json:"pushIComTime"` //在icom上链的时间
}

func GetTransactionHashesLastBlockNumber() uint64 {
	lastBlockByte, err := db.GetDB().Find([]byte(TransactionHashesLastBlockNumber))
	if err != nil {
		return 0
	}
	return utils.BytesToUint64(*lastBlockByte)

}

func ArbiscanTransactions(fromBlock, toBlock uint64) error {
	if fromBlock == 0 {
		fromBlock = config.DefaultArbitrumFromBlock
		lastBlock := GetTransactionHashesLastBlockNumber()
		if lastBlock > fromBlock {
			fromBlock = lastBlock + 1
		}
	}

	var toBlockAny any
	if toBlock == 0 {
		toBlockAny = "latest"
	} else {
		toBlockAny = toBlock
	}
	// API请求URL
	apiURL := fmt.Sprintf("%s?module=logs&action=getLogs&fromBlock=%d&toBlock=%v&address=%s&topic0=%s&apikey=%s",
		config.ArbitrumApi, fromBlock, toBlockAny, config.ArbitrumContractAddress, config.ArbitrumContractTopic, config.ArbitrumApiKey)

	// 发送HTTP请求
	response, err := http_client.GetRequest(apiURL)
	if err != nil {
		fmt.Printf("Failed to make API request: %s \n", err.Error())
		fmt.Printf("API request: %s \n", apiURL)
		return err
	}

	decoder := json.NewDecoder(bytes.NewBuffer(response))
	decoder.UseNumber()
	apiResponse := &APIResponse{}
	err = decoder.Decode(apiResponse)
	if err != nil {
		fmt.Printf("decoder.Decode Err: %s \n", err.Error())
		return err
	}

	if apiResponse.Status == "1" {
		for _, one := range apiResponse.Result {

			transactionHash, ok := one["transactionHash"].(string)
			if !ok {
				fmt.Printf("transactionHash.(string) Err: %s \n", err.Error())
				continue
			}

			if db.GetDB().CheckKeyExist([]byte(fmt.Sprintf(TransactionHashKey, transactionHash))) {
				continue
			}

			//处理交易记录
			blockNumber, err := strconv.ParseUint(strings.TrimPrefix(one["blockNumber"].(string), "0x"), 16, 64)
			if err != nil {
				fmt.Printf("strconv.ParseUint blockNumber Err: %s \n", err.Error())
				return err
			}
			data := one["data"].(string)
			data1 := data[2:66]
			data2 := data[66:130]
			data3 := data[130:194]
			data4 := data[194:258]
			data5 := data[258:322]
			data6 := data[322+64+64:]

			amount, _ := strconv.ParseUint(data3, 16, 64)
			rate, _ := strconv.ParseUint(data4, 16, 64)
			toAmount, _ := strconv.ParseUint(data5, 16, 64)
			iAddressBytes, _ := hex.DecodeString(data6)
			tr := &Transaction{
				TransactionHash: transactionHash,
				Sender:          "0x" + strings.ToUpper(data1),
				Token:           "0x" + strings.ToUpper(data2),
				BlockNumber:     blockNumber,
				Amount:          amount,
				Rate:            rate,
				ToAmount:        toAmount,
				IAddress:        strings.TrimRight(string(iAddressBytes), string(0)),
				SrcIAddress:     config.SrcIComAddress,
				Status:          -1,
			}
			logger.Log.Info("PullLog [ iAddress:%s toAmount:%d blockNumber:%d transactionHash:%s ]", tr.IAddress, tr.ToAmount, tr.BlockNumber, tr.TransactionHash)
			//保存到数据库
			err = tr.Save()
			if err != nil {
				fmt.Printf("SaveTransaction Err: %s  \n", err)
				return err
			}

			blockNumberByte := utils.Uint64ToBytes(tr.BlockNumber)
			err = db.GetDB().Save([]byte(TransactionHashesLastBlockNumber), &blockNumberByte)
			if err != nil {
				fmt.Printf("Save TransactionHashesLastBlockNumber Err %s \n", err.Error())
				return err
			}

			//将交易hash记录到待转账集合中
			err = db.GetDB().Save([]byte(fmt.Sprintf(WaitTransactionKey, tr.TransactionHash)), nil)
			if err != nil {
				fmt.Printf("Save WaitTransactionKey Err: %s \n", err.Error())
				return err
			}

			err = db.GetDB().Save([]byte(fmt.Sprintf(IAddressIndexTransactionHash, tr.IAddress, tr.TransactionHash)), nil)
			if err != nil {
				fmt.Printf("Save IAddressIndexTransactionHash Err: %s \n", err.Error())
				return err
			}

			//保存到数据库
			err = tr.SaveTransactionHashesByBlockNumber()
			if err != nil {
				fmt.Printf("SaveTransactionHashesByBlockNumber Err: %s \n", err.Error())
				return err
			}
			fmt.Printf("transactionHash:%s blockNumber:%d iAddress:%s toAmount:%d \n", tr.TransactionHash, tr.BlockNumber, tr.IAddress, tr.ToAmount)
		}
		fmt.Printf("成功获取 %d 条交易 \n", len(apiResponse.Result))
	} else {
		fmt.Printf("API request: %s \n", apiURL)
		fmt.Printf("API response status is not 1. Message is : %s \n", apiResponse.Message)
		return errors.New("API response status is not 1. Message is : %s" + apiResponse.Message)
	}
	return nil
}

func (t *Transaction) Save() error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return db.GetDB().Save([]byte(fmt.Sprintf(TransactionHashKey, t.TransactionHash)), &data)
}

func (t *Transaction) SaveTransactionHashesByBlockNumber() error {
	return db.GetDB().Save([]byte(fmt.Sprintf(TransactionHashesByBlockNumber, t.BlockNumber, t.TransactionHash)), nil)
}

func (t *Transaction) GetTransactionByTransactionHash(hash string) (*Transaction, error) {

	bs, err := db.GetDB().Find([]byte(fmt.Sprintf(TransactionHashKey, hash)))

	if err != nil {
		logger.Log.Error("decoder.Decode Err: %s", err.Error())
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewBuffer(*bs))
	decoder.UseNumber()
	tr := &Transaction{}
	err = decoder.Decode(tr)
	if err != nil {
		logger.Log.Error("decoder.Decode Err: %s", err.Error())
		return nil, err
	}

	return tr, nil
}

func TransactionPushICom() {
	keys, _, err := db.GetDB().FindPrefixKeyAll([]byte(WaitTransactionKeyPrefix))
	if err != nil {
		logger.Log.Error("Find TransactionHashKeyPrefix Err: %s", err.Error())
		return
	}

	var pushTxs = make([]*Transaction, 0)
	var findTxs = make([]*Transaction, 0)
	for _, k := range keys {
		transactionHash := string(k[len(WaitTransactionKeyPrefix):])

		tr := &Transaction{}
		tr, err = tr.GetTransactionByTransactionHash(transactionHash)
		if err != nil {
			logger.Log.Error("GetTransactionByTransactionHash Err %s", err.Error())
			continue
		}
		iAddrByte := crypto.AddressFromB58String(tr.IAddress)

		if iAddrByte == nil {
			logger.Log.Error("Not a valid icom address:%s TransactionHash:%s", tr.IAddress, tr.TransactionHash)
			db.GetDB().Remove(k)
			continue
		}

		if !crypto.ValidAddr(config.IComAddrPre, iAddrByte) {
			logger.Log.Error("Not a valid icom address:%s TransactionHash:%s", tr.IAddress, tr.TransactionHash)
			db.GetDB().Remove(k)
			continue
		}

		switch tr.Status {
		case -1, 3:
			pushTxs = append(pushTxs, tr)
		case 0, 1:
			findTxs = append(findTxs, tr)
		case 2:
			err = db.GetDB().Remove(k)
			if err != nil {
				logger.Log.Error("Remove err: %s", err.Error())
			}
		}
	}
	IComPushTx(pushTxs)
	IComFindTx(findTxs)
}

func IComPushTx(trs []*Transaction) {
	total := len(trs)

	if total == 0 {
		return
	}

	logger.Log.Info("---------------开始请求转帐 %d 条-----------------", total)
	successfulTotal := 0
	for i := 0; i < total; i += config.IcomRequestTxidMaxNum {
		end := i + config.IcomRequestTxidMaxNum
		if end > total {
			end = total
		}

		batch := trs[i:end]

		addresses := make([]map[string]interface{}, 0, len(batch))
		for _, one := range batch {
			addresses = append(addresses, map[string]interface{}{
				"address": one.IAddress,
				"amount":  one.ToAmount,
			})
		}

		params := map[string]interface{}{
			"method": "sendtoaddressmore",
			"params": map[string]interface{}{
				"srcaddress": config.SrcIComAddress,
				"addresses":  addresses,
				"gas":        config.IComGas,
				"pwd":        config.IComWalletPwd,
				"comment":    "exchange",
			},
		}
		res, err := http_client.PostReques(config.IComRpc, config.IComRpcUser, config.IComRpcPassword, params)
		if err != nil {
			logger.Log.Error("IComPushTx PostReques Err:%s", err.Error())
			continue
		}

		result, ok := res.(map[string]interface{})
		if !ok {
			logger.Log.Error("IComPushTx PostReques Err:  res.(map[string]interface{})")
			continue
		}

		txHash := result["hash"].(string)
		for _, one := range batch {
			if one.TxHash != "" {
				db.GetDB().Remove([]byte(fmt.Sprintf(IComTxHashIndexTransactionHash, one.TxHash, one.TransactionHash)))
			}
			db.GetDB().Save([]byte(fmt.Sprintf(IComTxHashIndexTransactionHash, txHash, one.TransactionHash)), nil)
			one.TxHash = txHash
			one.Status = 0
			one.SrcIAddress = config.SrcIComAddress
			err = one.Save()
			if err != nil {
				logger.Log.Error("IComPushTx Save Err:%s txHash:%s", err.Error(), one.TxHash)
				continue
			}
			successfulTotal++
			logger.Log.Info("IComPushTx [ iAddress:%s toAmount:%d txHash:%s blockNumber:%d transactionHash:%s ]", one.IAddress, one.ToAmount, one.TxHash, one.BlockNumber, one.TransactionHash)
		}

	}

	logger.Log.Info("---------------成功请求转帐 %d 条-----------------", successfulTotal)

}

func IComFindTx(trs []*Transaction) {
	txids := make([]string, 0)
	txidMap := make(map[string][]*Transaction)
	for _, tr := range trs {
		trarr, ok := txidMap[tr.TxHash]
		if ok {
			trarr = append(trarr, tr)
			txidMap[tr.TxHash] = trarr
		} else {
			txidMap[tr.TxHash] = []*Transaction{tr}
			txids = append(txids, tr.TxHash)
		}
	}

	total := len(txids)

	if total == 0 {
		return
	}

	for i := 0; i < total; i += config.IcomRequestTxidMaxNum {
		end := i + config.IcomRequestTxidMaxNum
		if end > total {
			end = total
		}

		batch := txids[i:end]

		params := map[string]interface{}{
			"method": "findtxs",
			"params": map[string]interface{}{
				"txids": batch,
			},
		}
		res, err := http_client.PostReques(config.IComRpc, config.IComRpcUser, config.IComRpcPassword, params)
		if err != nil {
			logger.Log.Error("IComFindTx PostReques Err: %s", err.Error())
			continue
		}
		result, ok := res.([]interface{})
		if !ok {
			logger.Log.Error("IComPushTx PostReques Err:  res.([]interface{})")
			continue
		}

		for _, v := range result {
			one, ok := v.(map[string]interface{})
			if !ok {
				logger.Log.Error("IComPushTx PostReques Err:  res.(map[string]interface{})")
				continue
			}
			txinfo := one["txinfo"].(map[string]interface{})
			txHash := txinfo["hash"].(string)

			ts := txidMap[txHash]
			upchaincode, _ := one["upchaincode"].(json.Number).Int64()
			timestamp, _ := one["timestamp"].(json.Number).Int64()
			for _, t := range ts {
				t.Status = upchaincode
				t.PushIComTime = timestamp
				err = t.Save()
				if err != nil {
					logger.Log.Error("IComFindTx Save Err:%s transactionHash:%s txHash:%s", err.Error(), t.TransactionHash, t.TxHash)
					continue
				}
				logger.Log.Info("IComFindTx [ iAddress:%s status:%d toAmount:%d txHash:%s blockNumber:%d transactionHash:%s ]", t.IAddress, t.Status, t.ToAmount, t.TxHash, t.BlockNumber, t.TransactionHash)
			}
		}
	}
}

func TransactionByBlockNumberRangePrint(minBlockNumber, maxBlockNumber uint64) {
	_, trs, err := getTransactionHashesByBlockNumberRange(minBlockNumber, maxBlockNumber)
	if err != nil {
		log.Println(err)
		return
	}
	TransactionListPrint(trs)
}

func TransactionByIComTxHashPrint(hash string) {
	prefix := IComTxHashIndexTransactionHashPrefix + hash + "_"
	keyBytes, _, err := db.GetDB().FindPrefixKeyAll([]byte(prefix))
	if err != nil {
		log.Println(err)
		return
	}

	trs := make([]*Transaction, 0, len(keyBytes))
	for _, keyByte := range keyBytes {
		key := string(keyByte)
		transactionHash := key[len(prefix):]
		tr := &Transaction{}
		tr, err = tr.GetTransactionByTransactionHash(transactionHash)
		if err != nil {
			log.Println("GetTransactionByTransactionHash Err", err)
			return
		}
		trs = append(trs, tr)
	}
	TransactionListPrint(trs)
}

func TransactionByIAddressPrint(iAddress string) {
	prefix := IAddressIndexTransactionHashPrefix + iAddress + "_"
	keyBytes, _, err := db.GetDB().FindPrefixKeyAll([]byte(prefix))
	if err != nil {
		log.Println(err)
		return
	}

	trs := make([]*Transaction, 0, len(keyBytes))
	for _, keyByte := range keyBytes {
		key := string(keyByte)
		transactionHash := key[len(prefix):]
		tr := &Transaction{}
		tr, err = tr.GetTransactionByTransactionHash(transactionHash)
		if err != nil {
			log.Println("GetTransactionByTransactionHash Err", err)
			return
		}
		trs = append(trs, tr)
	}
	TransactionListPrint(trs)
}

func TransactionByTransactionHashPrint(hash string) {
	tr := &Transaction{}
	tr, err := tr.GetTransactionByTransactionHash(hash)
	if err != nil {
		log.Println("GetTransactionByTransactionHash Err", err)
		return
	}
	TransactionListPrint([]*Transaction{tr})
}

func TransactionListPrint(trs []*Transaction) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"Transaction Hash",
		//"Sender",
		//"Token",
		"Block Number",
		//"Amount",
		//"Rate",
		"to Amount",
		"ICom Address",
		"ICom Tx Hash",
		"Status",
		"Push ICom Time",
	})

	for _, transaction := range trs {
		status := ""
		switch transaction.Status {
		case -1:
			status = "未转账"
		case 0:
			status = "转账中"
		case 1:
			status = "未上链"
		case 2:
			status = "成功上链"
		case 3:
			status = "上链失败"
		}

		timestr := ""
		if transaction.PushIComTime != 0 {
			timestr = time.Unix(transaction.PushIComTime, 0).Format("2006/01/02 15:04:05")
		}

		row := table.Row{
			transaction.TransactionHash,
			//transaction.Sender,
			//transaction.Token,
			strconv.FormatUint(transaction.BlockNumber, 10),
			//strconv.FormatUint(transaction.Amount, 10),
			//strconv.FormatUint(transaction.Rate, 10),
			strconv.FormatUint(transaction.ToAmount, 10),
			transaction.IAddress,
			transaction.TxHash,
			status,
			timestr,
		}
		t.AppendRow(row)
	}

	t.Render()
}

func getTransactionHashesByBlockNumberRange(minBlockNumber, maxBlockNumber uint64) ([]string, []*Transaction, error) {
	iter := db.GetDB().GetLevelDB().NewIterator(util.BytesPrefix([]byte(TransactionHashesByBlockNumberPrefix)), nil)
	defer iter.Release()

	hashes := make([]string, 0)
	trs := make([]*Transaction, 0)
	for iter.Next() {
		key := string(iter.Key())
		if key[:len([]byte(TransactionHashesByBlockNumberPrefix))] == TransactionHashesByBlockNumberPrefix {
			parts := strings.Split(key, "_")
			if len(parts) >= 3 {
				blockNumber, err := strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					log.Println("ParseUint err: ", err)
					continue
				}

				hash := parts[2]
				if blockNumber >= minBlockNumber && blockNumber <= maxBlockNumber {
					hashes = append(hashes, hash)
					tr := &Transaction{}
					tr, err = tr.GetTransactionByTransactionHash(hash)
					if err != nil {
						log.Fatal("GetTransactionByTransactionHash Err:", err)
						return hashes, trs, err
					}
					trs = append(trs, tr)
				}
			}

		}
	}

	if err := iter.Error(); err != nil {
		return hashes, trs, err
	}

	return hashes, trs, nil
}
