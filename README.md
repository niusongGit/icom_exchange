# icom_exchange






# arbiscan链的测试链接
## 1. 接口：https://api-sepolia.arbiscan.io/api?module=logs&action=getLogs&fromBlock=28091403&toBlock=latest&address=0x9C607b08fE911CA20F3c661e518C13089C32Be30&topic0=0x4a6af5cfd3969baa43d215151321689f373a20d038c968f77e16eba7c5161ac3&apikey=CCE52ECVS4PQK3D56ZD699Y24VJ58AZGBU
**功能**：返回合约事件日志

**返回参数**：
```json
{
  "status": "1",
  "message": "OK",
  "result": [
    {
      "address": "0x9c607b08fe911ca20f3c661e518c13089c32be30",
      "topics": [
        "0x4a6af5cfd3969baa43d215151321689f373a20d038c968f77e16eba7c5161ac3"
      ],
      "data": "0x0000000000000000000000005b27836a2c5c649724e35a56f6fb9528a8d15f60000000000000000000000000a8c2368a6f0a97d9a9639b09ceb0cc8cf4028f8200000000000000000000000000000000000000000000000000000000003567e00000000000000000000000000000000000000000000000000000000029b92700000000000000000000000000000000000000000000000000000000009208088000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000046173646600000000000000000000000000000000000000000000000000000000",
      "blockNumber": "0x1df6aa1",
      "blockHash": "0x36208fa994f78d77690ebb69c00b4bf91ac834791c2d16af22adb84c2da0c00b",
      "timeStamp": "0x66134412",
      "gasPrice": "0xc1ce421",
      "gasUsed": "0x3bc1c",
      "logIndex": "0x9",
      "transactionHash": "0xcef55a38f78636a0a2905a3c3326da8bfdfd5c55d73a8fa7d7568a7e9991ae80",
      "transactionIndex": "0x2"
    },
    {
      "address": "0x9c607b08fe911ca20f3c661e518c13089c32be30",
      "topics": [
        "0x4a6af5cfd3969baa43d215151321689f373a20d038c968f77e16eba7c5161ac3"
      ],
      "data": "0x0000000000000000000000005b27836a2c5c649724e35a56f6fb9528a8d15f600000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016bcc41e90000000000000000000000000000000000000000000000000000000000000029b92700000000000000000000000000000000000000000000000000000000000445c00000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000046173646600000000000000000000000000000000000000000000000000000000",
      "blockNumber": "0x1df6ab1",
      "blockHash": "0x1c64470a61a3c1b701613f30ce5590c34872e6f1c546bea34f80e7035a1e84c8",
      "timeStamp": "0x66134416",
      "gasPrice": "0xc108811",
      "gasUsed": "0x2fdb8",
      "logIndex": "0x5",
      "transactionHash": "0x9521ea085ac84cb76ac54c02b0b2f902d8bdb12e5da2c96f7053b75f9dd20f7f",
      "transactionIndex": "0x2"
    },
    {
      "address": "0x9c607b08fe911ca20f3c661e518c13089c32be30",
      "topics": [
        "0x4a6af5cfd3969baa43d215151321689f373a20d038c968f77e16eba7c5161ac3"
      ],
      "data": "0x0000000000000000000000005b27836a2c5c649724e35a56f6fb9528a8d15f60000000000000000000000000a8c2368a6f0a97d9a9639b09ceb0cc8cf4028f8200000000000000000000000000000000000000000000000000000000003567e00000000000000000000000000000000000000000000000000000000029b92700000000000000000000000000000000000000000000000000000000009208088000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000046173646600000000000000000000000000000000000000000000000000000000",
      "blockNumber": "0x1df6bdd",
      "blockHash": "0x8387b6fb0627390c5ead2f18fc35fcbbc2e7639e52042465c64051f11de1e7d6",
      "timeStamp": "0x66134462",
      "gasPrice": "0xe648a31",
      "gasUsed": "0x31480",
      "logIndex": "0x9",
      "transactionHash": "0x00a4716d3405ec69583f9bf2a32a59c965761addb8f7d5069f38c38ae569b08c",
      "transactionIndex": "0x2"
    },
    {
      "address": "0x9c607b08fe911ca20f3c661e518c13089c32be30",
      "topics": [
        "0x4a6af5cfd3969baa43d215151321689f373a20d038c968f77e16eba7c5161ac3"
      ],
      "data": "0x0000000000000000000000005b27836a2c5c649724e35a56f6fb9528a8d15f600000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016bcc41e90000000000000000000000000000000000000000000000000000000000000029b92700000000000000000000000000000000000000000000000000000000000445c00000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000046173646600000000000000000000000000000000000000000000000000000000",
      "blockNumber": "0x1df6bed",
      "blockHash": "0x5c02475329b5797f00fc05aea9408142c0000fd97f15bdb2ca3cf2e16777d74a",
      "timeStamp": "0x66134466",
      "gasPrice": "0xe67e591",
      "gasUsed": "0x2aa2c",
      "logIndex": "0x27",
      "transactionHash": "0x52edc90172e3ea4de193b3b43e40edad14fd9f956bf148793191cda3f22086fc",
      "transactionIndex": "0x6"
    }
  ]
}
```
**data参数说明**：
-  0x
-  0000000000000000000000005b27836a2c5c649724e35a56f6fb9528a8d15f60 //sender 发送者
-  000000000000000000000000a8c2368a6f0a97d9a9639b09ceb0cc8cf4028f82 //token 代币地址，0地址为eth
-  00000000000000000000000000000000000000000000000000000000003567e0 //amount 代币输入数量
-  0000000000000000000000000000000000000000000000000000000029b92700 //rate 兑换比例
-  0000000000000000000000000000000000000000000000000000000092080880 //toAmount 代币输出数量
-  00000000000000000000000000000000000000000000000000000000000000c0
-  0000000000000000000000000000000000000000000000000000000000000004
-  6173646600000000000000000000000000000000000000000000000000000000 //iAddress icom地址

**命令说明**：
```
 ////////////////////////////////////////////////////////命令说明//////////////////////////////////////////////////

        pull-log         拉取兑换日志，参数minBlockNumber为空时默认从上一次拉取的最高区块开始，maxBlockNumber为空时拉取到最新高度
        find-log         查询日志，参数minBlockNumber必填，maxBlockNumber必填
        find-log-thash   按TransactionsHash查找，参数TransactionsHash必填
        find-log-ihash   按ICom的TxHash查找，参数TxHash必填
        find-log-iaddr   按ICom的Address查找，参数IAddress必填
        last-block-num   获取已经拉取日志的最高区块高度，无参数
        exit             退出程序
        help             获取取命令说明
```
