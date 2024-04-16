package config

import (
	"bytes"
	"os"
	"path/filepath"
)

const (
	Path_config = "conf/config.json"
)

var (
	SrcIComAddress               = ""
	ArbitrumApi                  = "https://api-sepolia.arbiscan.io/api"
	DefaultArbitrumFromBlock     = uint64(0)
	ArbitrumContractAddress      = "0x9C607b08fE911CA20F3c661e518C13089C32Be30"
	ArbitrumContractTopic        = "0x4a6af5cfd3969baa43d215151321689f373a20d038c968f77e16eba7c5161ac3"
	ArbitrumApiKey               = "CCE52ECVS4PQK3D56ZD699Y24VJ58AZGBU"
	IComRpc                      = "http://0.0.0.0:2080/rpc"
	IComRpcUser                  = "test"
	IComRpcPassword              = "testp"
	IComWalletPwd                = "123456789"
	IComGas                      = uint64(1000)
	IComAddrPre                  = "iCom"
	IcomRequestTxidMaxNum        = 50
	CheckTransactionPushIComTime = 1
)

type Config struct {
	SrcIComAddress               string `json:"SrcIComAddress"`           //ICom转账地址
	ArbitrumApi                  string `json:"ArbitrumApi"`              //Arbitrum api url
	ArbitrumContractAddress      string `json:"ArbitrumContractAddress"`  //Arbitrum 中的合约地址
	ArbitrumContractTopic        string `json:"ArbitrumContractTopic"`    //Arbitrum 中的合约事件
	ArbitrumApiKey               string `json:"ArbitrumApiKey"`           // APIkey
	DefaultArbitrumFromBlock     uint64 `json:"DefaultArbitrumFromBlock"` //Arbitrum默认扒取数据的起始区块高度（一般为合约部署时的高度）
	IComRpc                      string `json:"IComRpc"`
	IComRpcUser                  string `json:"IComRpcUser"`
	IComRpcPassword              string `json:"IComRpcPassword"`
	IComWalletPwd                string `json:"IComWalletPwd"`
	IComGas                      uint64 `json:"IComGas"`
	IComAddrPre                  string `json:"IComAddrPre"`
	IcomRequestTxidMaxNum        int    `json:"IcomRequestTxidMaxNum"`        //一次请求最多交易数
	CheckTransactionPushIComTime int    `json:"CheckTransactionPushIComTime"` //检查icom上链交易间隔时间单位秒
}

func Step() {
	ok, err := PathExists(Path_config)
	if err != nil {
		panic("检查配置文件错误：" + err.Error())
		return
	}

	if !ok {
		panic("检查配置文件错误")
		return
	}

	bs, err := os.ReadFile(filepath.Join(Path_config))
	if err != nil {
		panic("读取配置文件错误：" + err.Error())
		return
	}

	cfi := new(Config)

	decoder := json.NewDecoder(bytes.NewBuffer(bs))
	decoder.UseNumber()
	err = decoder.Decode(cfi)
	if err != nil {
		panic("解析配置文件错误：" + err.Error())
		return
	}

	if len(cfi.SrcIComAddress) > 0 {
		SrcIComAddress = cfi.SrcIComAddress
		if SrcIComAddress == "" {
			panic("解析配置文件错误：SrcIComAddress未配置")
			return
		}
	}
	if len(cfi.ArbitrumApi) > 0 {
		ArbitrumApi = cfi.ArbitrumApi
	}
	if len(cfi.ArbitrumContractAddress) > 0 {
		ArbitrumContractAddress = cfi.ArbitrumContractAddress
	}

	if len(cfi.ArbitrumApiKey) > 0 {
		ArbitrumApiKey = cfi.ArbitrumApiKey
	}

	if cfi.DefaultArbitrumFromBlock > 0 {
		DefaultArbitrumFromBlock = cfi.DefaultArbitrumFromBlock
	}

	if len(cfi.IComRpc) > 0 {
		IComRpc = cfi.IComRpc
	}

	if len(cfi.IComRpcUser) > 0 {
		IComRpcUser = cfi.IComRpcUser
	}

	if len(cfi.IComRpcPassword) > 0 {
		IComRpcPassword = cfi.IComRpcPassword
	}

	if len(cfi.IComWalletPwd) > 0 {
		IComWalletPwd = cfi.IComWalletPwd
	}

	if len(cfi.IComAddrPre) > 0 {
		IComAddrPre = cfi.IComAddrPre
	}

	if cfi.IComGas > 0 {
		IComGas = cfi.IComGas
	}

	if cfi.IcomRequestTxidMaxNum > 0 {
		if cfi.IcomRequestTxidMaxNum > IcomRequestTxidMaxNum {
			panic("解析配置文件错误：IcomRequestTxidMaxNum最大50！！")
			return
		}
		IcomRequestTxidMaxNum = cfi.IcomRequestTxidMaxNum
	}

	if cfi.CheckTransactionPushIComTime > 0 {
		CheckTransactionPushIComTime = cfi.CheckTransactionPushIComTime
	}

}
