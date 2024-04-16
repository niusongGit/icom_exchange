package main

import (
	"bufio"
	"context"
	"fmt"
	"icom_exchange/config"
	"icom_exchange/corelogic"
	"icom_exchange/db"
	"icom_exchange/logger"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	wg sync.WaitGroup // 等待所有任务完成
)

func main() {
	config.Step()
	l := logger.InitLog("./logs/")
	defer l.Close()
	//本地数据库
	db.InitDB("./leveldb_database")

	ctx, cancel := context.WithCancel(context.Background())
	// 创建定时器，每2秒执行一次 TransactionPushICom 方法
	ticker := time.NewTicker(time.Duration(config.CheckTransactionPushIComTime) * time.Second)
	go func(ctx context.Context) {
		defer wg.Done()
		wg.Add(1)

		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				corelogic.TransactionPushICom()
			}
		}
	}(ctx)

	// 监听系统信号，例如Ctrl+C
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		helpStr := "\n\n        ////////////////////////////////////////////////////////命令说明//////////////////////////////////////////////////\n\n        pull-log         拉取兑换日志，参数minBlockNumber为空时默认从上一次拉取的最高区块开始，maxBlockNumber为空时拉取到最新高度\n        find-log         查询日志，参数minBlockNumber必填，maxBlockNumber必填\n        find-log-thash   按TransactionsHash查找，参数TransactionsHash必填\n        find-log-ihash   按ICom的TxHash查找，参数TxHash必填\n        find-log-iaddr   按ICom的Address查找，参数IAddress必填\n        last-block-num   获取已经拉取日志的最高区块高度，无参数\n        exit             退出程序\n        help             获取取命令说明\n\n"
		fmt.Println(helpStr)

		fmt.Println("*************************请输入指令：")

		// 创建命令行输入监听器
		scanner := bufio.NewScanner(os.Stdin)

		// 监听命令行输入
		for scanner.Scan() {
			command := scanner.Text()

			switch strings.ToLower(command) {
			case "pull-log":
				// 提示用户输入 minBlockNumber 和 maxBlockNumber
				fmt.Println("当前已拉取到最高高度：", corelogic.GetTransactionHashesLastBlockNumber())
				fmt.Print("Enter minBlockNumber: ")
				scanner.Scan()
				minBlockNumberStr := scanner.Text()
				minBlockNumber := uint64(0)
				var err error
				if minBlockNumberStr != "" {
					minBlockNumber, err = strconv.ParseUint(minBlockNumberStr, 10, 64)
					if err != nil {
						fmt.Println("Invalid minBlockNumber")
						fmt.Println("\n\n\n*************************请输入指令：")
						continue
					}
				}

				fmt.Print("Enter maxBlockNumber: ")
				scanner.Scan()
				maxBlockNumberStr := scanner.Text()
				maxBlockNumber := uint64(0)
				if maxBlockNumberStr != "" {
					maxBlockNumber, err = strconv.ParseUint(maxBlockNumberStr, 10, 64)
					if err != nil {
						fmt.Println("Invalid maxBlockNumber")
						fmt.Println("\n\n\n*************************请输入指令：")
						continue
					}
				}

				corelogic.ArbiscanTransactions(minBlockNumber, maxBlockNumber)
			case "find-log":
				fmt.Println("当前已拉取到最高高度：", corelogic.GetTransactionHashesLastBlockNumber())
				// 提示用户输入 minBlockNumber 和 maxBlockNumber
				fmt.Print("Enter minBlockNumber: ")
				scanner.Scan()
				minBlockNumberStr := scanner.Text()
				minBlockNumber, err := strconv.ParseUint(minBlockNumberStr, 10, 64)
				if err != nil {
					fmt.Println("Invalid minBlockNumber")
					fmt.Println("\n\n\n*************************请输入指令：")
					continue
				}

				fmt.Print("Enter maxBlockNumber: ")
				scanner.Scan()
				maxBlockNumberStr := scanner.Text()
				maxBlockNumber, err := strconv.ParseUint(maxBlockNumberStr, 10, 64)
				if err != nil {
					fmt.Println("Invalid maxBlockNumber")
					fmt.Println("\n\n\n*************************请输入指令：")
					continue
				}
				corelogic.TransactionByBlockNumberRangePrint(minBlockNumber, maxBlockNumber)
			case "find-log-thash":
				// 提示用户输入 minBlockNumber 和 maxBlockNumber
				fmt.Print("Enter TransactionHash: ")
				scanner.Scan()
				transactionHash := scanner.Text()
				if transactionHash == "" {
					fmt.Println("Invalid TransactionHash")
					fmt.Println("\n\n\n*************************请输入指令：")
					continue
				}
				corelogic.TransactionByTransactionHashPrint(transactionHash)
			case "find-log-ihash":
				// 提示用户输入 minBlockNumber 和 maxBlockNumber
				fmt.Print("Enter IComTxHash: ")
				scanner.Scan()
				iComTxHash := scanner.Text()
				if iComTxHash == "" {
					fmt.Println("Invalid IComTxHash")
					fmt.Println("\n\n\n*************************请输入指令：")
					continue
				}
				corelogic.TransactionByIComTxHashPrint(iComTxHash)
			case "find-log-iaddr":
				// 提示用户输入 minBlockNumber 和 maxBlockNumber
				fmt.Print("Enter IAddress: ")
				scanner.Scan()
				IAddress := scanner.Text()
				if IAddress == "" {
					fmt.Println("Invalid IAddress")
					fmt.Println("\n\n\n*************************请输入指令：")
					continue
				}
				corelogic.TransactionByIAddressPrint(IAddress)
			case "last-block-num":
				fmt.Println(corelogic.GetTransactionHashesLastBlockNumber())
			case "exit":
				// 优雅地退出程序
				cancel()
				fmt.Println("正在退出...")
				return
			case "help":
				fmt.Println(helpStr)
			default:
				fmt.Println("Invalid command")
			}
			fmt.Println("\n\n\n*************************请输入指令：")
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading standard input:", err)
		}
	}()
	// 等待所有任务完成或接收到系统信号
	for {
		select {
		case <-signals:
			cancel()
			fmt.Println("正在退出...")
			break
		case <-ctx.Done():
			break
		}
		break
	}
	wg.Wait()

	fmt.Println("Exiting the program...")
}
