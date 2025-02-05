package utils

import (
	"math/big"
)

// Config
type ServerConfig struct {
	Port      string `json:"port"`
	FromBlock int64  `json:"from_block"`
}

type SqliteConfig struct {
	Database string `json:"database"`
}

type ChainConfig struct {
	ChainName string `json:"chain_name"`
	Rpc       string `json:"rpc"`
	UserName  string `json:"user_name"`
	PassWord  string `json:"pass_word"`
}

// RouterResult
type HttpResult struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
}

type BaseParams struct {
	P  string `json:"p"`
	Op string `json:"op"`
}

type NewParams struct {
	P              string `json:"p"`
	Op             string `json:"op"`
	Tick           string `json:"tick"`
	Max            string `json:"max"`
	Amt            string `json:"amt"`
	Lim            string `json:"lim"`
	Dec            int64  `json:"dec"`
	Burn           string `json:"burn"`
	Func           string `json:"func"`
	ReceiveAddress string `json:"receive_address"`
	ToAddress      string `json:"to_address"`
	RateFee        string `json:"rate_fee"`
	Repeat         int64  `json:"repeat"`
}

type Drc20Params struct {
	Tick      string `json:"tick"`
	Limit     uint64 `json:"limit"`
	OffSet    uint64 `json:"offset"`
	Completed uint64 `json:"completed"`
}

type SwapParams struct {
	Op            string `json:"op"`
	Tick0         string `json:"tick0"`
	Tick1         string `json:"tick1"`
	Amt0          string `json:"amt0"`
	Amt1          string `json:"amt1"`
	Amt0Min       string `json:"amt0_min"`
	Amt1Min       string `json:"amt1_min"`
	Liquidity     string `json:"liquidity"`
	Path          string `json:"path"`
	HolderAddress string `json:"holder_address"`
}

type WDogeParams struct {
	Op            string `json:"op"`
	Tick          string `json:"tick"`
	Amt           string `json:"amt"`
	HolderAddress string `json:"holder_address"`
}

// Models
type Cardinals struct {
	OrderId            string   `json:"order_id"`
	P                  string   `json:"p"`
	Op                 string   `json:"op"`
	Tick               string   `json:"tick"`
	Amt                *big.Int `json:"amt"`
	Max                *big.Int `json:"max"`
	Lim                *big.Int `json:"lim"`
	Dec                int64    `json:"dec"`
	Burn               string   `json:"burn"`
	Func               string   `json:"func"`
	RateFee            *big.Int `json:"rate_fee"`
	Repeat             int64    `json:"repeat"`
	FeeTxHash          string   `json:"fee_tx_hash"`
	FeeTxIndex         uint32   `json:"fee_tx_index"`
	FeeTxRaw           string   `json:"fee_tx_raw"`
	Drc20TxHash        string   `json:"drc20_tx_hash"`
	Drc20TxRaw         string   `json:"drc20_tx_raw"`
	BlockNumber        int64    `json:"block_number"`
	BlockHash          string   `json:"block_hash"`
	BlockConfirmations uint64   `json:"block_confirmations"`
	ReceiveAddress     string   `json:"receive_address"`
	ToAddress          string   `json:"to_address"`
	AdminAddress       string   `json:"admin_address"`
	FeeAddress         string   `json:"fee_address"`
	OrderStatus        int64    `json:"order_status"`
	ErrInfo            string   `json:"err_info"`
	CreateDate         string   `json:"create_date"`
}

// SWAP
type SwapInfo struct {
	OrderId         string   `json:"order_id"`
	Op              string   `json:"op"`
	Tick            string   `json:"tick"`
	Tick0           string   `json:"tick0"`
	Tick1           string   `json:"tick1"`
	Amt0            *big.Int `json:"amt0"`
	Amt1            *big.Int `json:"amt1"`
	Amt0Min         *big.Int `json:"amt0_min"`
	Amt1Min         *big.Int `json:"amt1_min"`
	Amt0Out         *big.Int `json:"amt0_out"`
	Amt1Out         *big.Int `json:"amt1_out"`
	Path            []string `json:"path"`
	Liquidity       *big.Int `json:"liquidity"`
	HolderAddress   string   `json:"holder_address"`
	FeeAddress      string   `json:"fee_address"`
	FeeTxHash       string   `json:"fee_tx_hash"`
	FeeTxIndex      uint32   `json:"fee_tx_index"`
	FeeTxRaw        *string  `json:"fee_tx_raw"`
	FeeBlockNumber  int64    `json:"fee_block_number"`
	FeeBlockHash    string   `json:"fee_block_hash"`
	SwapTxHash      string   `json:"swap_tx_hash"`
	SwapTxRaw       *string  `json:"swap_tx_raw"`
	SwapBlockNumber int64    `json:"swap_block_number"`
	SwapBlockHash   string   `json:"swap_block_hash"`
	OrderStatus     int64    `json:"order_status"`
	UpdateDate      string   `json:"update_date"`
	CreateDate      string   `json:"create_date"`
}

// swap_liquidity
type SwapLiquidity struct {
	Tick            string   `json:"tick"`
	Tick0           string   `json:"tick0"`
	Tick1           string   `json:"tick1"`
	Amt0            *big.Int `json:"amt0"`
	Amt1            *big.Int `json:"amt1"`
	Path            string   `json:"path"`
	LiquidityTotal  *big.Int `json:"liquidity_total"`
	ReservesAddress string   `json:"reserves_address"`
	HolderAddress   string   `json:"holder_address"`
}

// swap_revert
type SwapRevert struct {
	Tick        string   `json:"tick"`
	FromAddress string   `json:"from_address"`
	ToAddress   string   `json:"to_address"`
	Amt         *big.Int `json:"amt"`
	BlockNumber int64    `json:"block_number"`
}

// WDOGE
type WDogeInfo struct {
	OrderId          string   `json:"order_id"`
	Op               string   `json:"op"`
	Tick             string   `json:"tick"`
	Amt              *big.Int `json:"amt"`
	HolderAddress    string   `json:"holder_address"`
	FeeAddress       string   `json:"fee_address"`
	FeeTxHash        string   `json:"fee_tx_hash"`
	FeeTxIndex       uint32   `json:"fee_tx_index"`
	FeeTxRaw         *string  `json:"fee_tx_raw"`
	FeeBlockNumber   int64    `json:"fee_block_number"`
	FeeBlockHash     string   `json:"fee_block_hash"`
	WDogeTxHash      string   `json:"wdoge_tx_hash"`
	WDogeTxRaw       *string  `json:"wdoge_tx_raw"`
	WDogeBlockNumber int64    `json:"wdoge_block_number"`
	WDogeBlockHash   string   `json:"wdoge_block_hash"`
	UpdateDate       string   `json:"update_date"`
	CreateDate       string   `json:"create_date"`
}
