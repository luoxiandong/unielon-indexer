package storage

import (
	"fmt"
	"github.com/unielon-org/unielon-indexer/utils"
)

func (c *DBClient) InstallWDogeInfo(wdoge *utils.WDogeInfo) error {
	query := "INSERT INTO wdoge_info (order_id, op, tick, amt, fee_address, holder_address, wdoge_tx_hash ) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := c.SqlDB.Exec(query, wdoge.OrderId, wdoge.Op, wdoge.Tick, wdoge.Amt.String(), wdoge.FeeAddress, wdoge.HolderAddress, wdoge.WDogeTxHash)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (c *DBClient) UpdateWDogeInfo(swap *utils.WDogeInfo) error {
	query := "update wdoge_info set fee_tx_hash = ?, fee_tx_raw = ?, fee_tx_hash= ?, fee_tx_index = ?, fee_block_hash = ?, fee_block_number = ?, wdoge_tx_hash = ?, wdoge_tx_raw = ? where order_id = ?"
	_, err := c.SqlDB.Exec(query, swap.FeeTxHash, swap.FeeTxRaw, swap.FeeTxHash, swap.FeeTxIndex, swap.FeeBlockHash, swap.FeeBlockNumber, swap.WDogeTxHash, swap.WDogeTxRaw, swap.OrderId)
	if err != nil {
		return err
	}
	return nil
}

func (c *DBClient) UpdateWDogeInfoFork(height int64) error {
	query := "update swap_info set wdoge_block_number = '0', wdoge_block_hash = '' where wdoge_block_number > ?"
	_, err := c.SqlDB.Exec(query, height)
	if err != nil {
		return err
	}
	return nil
}

func (c *DBClient) FindWDogeInfoByTxHash(WdogeTXHash string) (*utils.WDogeInfo, error) {
	query := "SELECT  order_id, op, tick, amt, fee_tx_hash, fee_tx_index, fee_block_hash, fee_block_number, wdoge_tx_hash, wdoge_tx_raw, wdoge_block_hash, wdoge_block_number, fee_address, holder_address, update_date, create_date FROM wdoge_info where wdoge_tx_hash = ?"
	rows, err := c.SqlDB.Query(query, WdogeTXHash)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		wdoge := &utils.WDogeInfo{}
		var amt string
		err := rows.Scan(&wdoge.OrderId, &wdoge.Op, &wdoge.Tick, &amt, &wdoge.FeeTxHash, &wdoge.FeeTxIndex, &wdoge.FeeBlockHash, &wdoge.FeeBlockNumber, &wdoge.WDogeTxHash, &wdoge.WDogeTxRaw, &wdoge.WDogeBlockHash, &wdoge.WDogeBlockNumber, &wdoge.FeeAddress, &wdoge.HolderAddress, &wdoge.UpdateDate, &wdoge.CreateDate)
		if err != nil {
			return nil, err
		}
		wdoge.Amt, _ = utils.ConvetStr(amt)
		return wdoge, nil
	}
	return nil, nil
}

func (c *DBClient) FindWDogeInfo(op, holder_address string) ([]*utils.WDogeInfo, int64, error) {
	query := "SELECT  order_id, op, tick, amt, fee_tx_hash, fee_tx_index, fee_block_hash, fee_block_number, wdoge_tx_hash, wdoge_block_hash, wdoge_block_number, fee_address, holder_address,  update_date, create_date FROM wdoge_info  "

	where := "where"
	whereAges := []any{}

	if op != "" {
		where += "  op = ? "
		whereAges = append(whereAges, op)
	}

	if holder_address != "" {
		if where != "where" {
			where += " and "
		}
		where += "  holder_address = ? "
		whereAges = append(whereAges, holder_address)
	}

	order := " order by update_date desc "

	rows, err := c.SqlDB.Query(query+where+order, whereAges...)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	wdoges := make([]*utils.WDogeInfo, 0)
	for rows.Next() {
		wdoge := &utils.WDogeInfo{}
		var amt string
		rows.Scan(&wdoge.OrderId, &wdoge.Op, &wdoge.Tick, &amt, &wdoge.FeeTxHash, &wdoge.FeeTxIndex, &wdoge.FeeBlockHash, &wdoge.FeeBlockNumber, &wdoge.WDogeTxHash, &wdoge.WDogeBlockHash, &wdoge.WDogeBlockNumber, &wdoge.FeeAddress, &wdoge.HolderAddress, &wdoge.UpdateDate, &wdoge.CreateDate)

		if err != nil {
			return nil, 0, err
		}

		wdoge.Amt, _ = utils.ConvetStr(amt)
		wdoges = append(wdoges, wdoge)
	}

	query1 := "SELECT count(order_id)  FROM wdoge_info "
	rows1, err := c.SqlDB.Query(query1+where, whereAges...)
	if err != nil {
		return nil, 0, err
	}

	defer rows1.Close()
	total := int64(0)
	if rows1.Next() {
		rows1.Scan(&total)
	}

	return wdoges, total, nil
}
