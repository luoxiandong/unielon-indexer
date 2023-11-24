package storage

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/unielon-org/unielon-indexer/utils"
	"math/big"
)

func (c *DBClient) InstallCardinalsInfo(card *utils.Cardinals) error {
	query := "INSERT INTO cardinals_info (order_id, p, op, tick, amt, max_, lim_, dec_, burn_, func_, receive_address, fee_address, to_address, fee_tx_hash, drc20_tx_hash, block_number, block_hash, repeat_mint) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err := c.SqlDB.Exec(query, card.OrderId, card.P, card.Op, card.Tick, card.Amt.String(), card.Max.String(), card.Lim.String(), card.Dec, card.Burn, card.Func, card.ReceiveAddress, card.FeeAddress, card.ToAddress, card.FeeTxHash, card.Drc20TxHash, card.BlockNumber, card.BlockHash, card.Repeat)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (c *DBClient) InstallDrc20(max, lim *big.Int, tick, receive_address, drc20_tx_hash string) error {
	query := "INSERT INTO drc20_info (tick, `max_`, lim_, receive_address, drc20_tx_hash) VALUES (?, ?, ?, ?, ?)"
	_, err := c.SqlDB.Exec(query, tick, max.String(), lim.String(), receive_address, drc20_tx_hash)
	if err != nil {
		return err
	}
	return nil
}

func (c *DBClient) UpdateCardinalsBlockNumber(card *utils.Cardinals) error {
	query := "UPDATE cardinals_info SET block_number = ?, block_hash = ?, block_confirmations = ?, order_status = ? where order_id = ?"
	_, err := c.SqlDB.Exec(query, card.BlockNumber, card.BlockHash, card.BlockConfirmations, card.OrderStatus, card.OrderId)
	if err != nil {
		return err
	}
	return nil
}

func (c *DBClient) UpdateCardinalsInfoNewErrInfo(orderId, errInfo string) error {
	query := "update cardinals_info set err_info = ?, order_status = 1  where order_id = ?"
	_, err := c.SqlDB.Exec(query, errInfo, orderId)
	if err != nil {
		return err
	}
	return nil
}

func (c *DBClient) UpdateAddressBalanceMint(tx *sql.Tx, tick string, sum1, sum2 *big.Int, address string, sub bool) error {

	var err error
	if tx == nil {
		tx, err = c.SqlDB.Begin()
		if err != nil {
			return err
		}

		defer func() {
			err = tx.Commit()
			if err != nil {
				tx.Rollback()
			}
		}()
	}

	update1 := "UPDATE drc20_info SET amt_sum=?, transactions = transactions + 1 WHERE tick = ?"
	if sub {
		update1 = "UPDATE drc20_info SET amt_sum=?, transactions = transactions - 1 WHERE tick = ?"
	}
	_, err = tx.Exec(update1, sum1.String(), tick)
	if err != nil {
		tx.Rollback()
		return err
	}

	update2 := "INSERT OR REPLACE INTO drc20_address_info (tick, receive_address, amt_sum) VALUES (?, ?, ?) "
	_, err = tx.Exec(update2, tick, address, sum2.String())
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (c *DBClient) UpdateAddressBalanceTran(tx *sql.Tx, tick string, sum1 *big.Int, address1 string, sum2 *big.Int, address2 string, sub bool) error {

	var err error
	if tx == nil {
		tx, err = c.SqlDB.Begin()
		if err != nil {
			return err
		}

		defer func() {
			err = tx.Commit()
			if err != nil {
				tx.Rollback()
			}
		}()
	}

	update1 := "UPDATE drc20_info SET transactions = transactions + 1 WHERE tick = ?"
	if sub {
		update1 = "UPDATE drc20_info SET transactions = transactions - 1 WHERE tick = ?"
	}
	_, err = tx.Exec(update1, tick)
	if err != nil {
		tx.Rollback()
		return err
	}

	update2 := "INSERT OR REPLACE INTO drc20_address_info (tick, receive_address, amt_sum) VALUES (?, ?, ?)"
	_, err = tx.Exec(update2, tick, address1, sum1.String())
	if err != nil {
		tx.Rollback()
		return err
	}

	update3 := "INSERT OR REPLACE INTO drc20_address_info (tick, receive_address, amt_sum) VALUES (?, ?, ?)"
	_, err = tx.Exec(update3, tick, address2, sum2.String())
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (c *DBClient) DelDrc20Info(tick, receive_address, drc20_tx_hash string) error {
	query := "DELETE FROM drc20_info where tick = ? and receive_address= ? and drc20_tx_hash = ?"
	_, err := c.SqlDB.Exec(query, tick, receive_address, drc20_tx_hash)
	if err != nil {
		return err
	}
	return nil
}

func (c *DBClient) FindDrc20InfoSumByTick(tick string) (*big.Int, *big.Int, *big.Int, error) {
	query := "SELECT amt_sum, max_, lim_ FROM drc20_info WHERE tick = ?"
	rows, err := c.SqlDB.Query(query, tick)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	if rows.Next() {

		var sum, max, lim string
		err := rows.Scan(&sum, &max, &lim)

		if err != nil {
			return nil, nil, nil, err
		}

		is_ok := false
		sum_big := new(big.Int)
		if sum != "" {
			sum_big, is_ok = new(big.Int).SetString(sum, 10)
			if !is_ok {
				return nil, nil, nil, fmt.Errorf("max error")
			}
		}

		max_big := new(big.Int)
		if max != "" {
			max_big, is_ok = new(big.Int).SetString(max, 10)
			if !is_ok {
				return nil, nil, nil, fmt.Errorf("max error")
			}
		}

		lim_big := new(big.Int)
		if lim != "" {
			lim_big, is_ok = new(big.Int).SetString(lim, 10)
			if !is_ok {
				return nil, nil, nil, fmt.Errorf("lim error")
			}
		}
		return sum_big, max_big, lim_big, nil
	}

	return nil, nil, nil, errors.New("not found")
}

func (c *DBClient) FindSwapDrc20InfoByTick(tx *sql.Tx, tick string) (*big.Int, *big.Int, *big.Int, error) {
	query := "SELECT amt_sum, max_, lim_ FROM drc20_info WHERE tick = ?"
	rows, err := tx.Query(query, tick)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	if rows.Next() {

		var sum, max, lim string
		err := rows.Scan(&sum, &max, &lim)

		if err != nil {
			return nil, nil, nil, err
		}

		sum_big, _ := utils.ConvetStr(sum)
		max_big, _ := utils.ConvetStr(max)
		lim_big, _ := utils.ConvetStr(lim)
		return sum_big, max_big, lim_big, nil
	}

	return nil, nil, nil, errors.New("not found")
}

func (c *DBClient) FindDrc20AddressInfoByTick(tick string, address string) (*big.Int, error) {
	query := "SELECT amt_sum  FROM drc20_address_info WHERE tick = ? and receive_address = ?"
	rows, err := c.SqlDB.Query(query, tick, address)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {

		var sum string
		err := rows.Scan(&sum)

		if err != nil {
			return nil, err
		}

		is_ok := false
		sum_big := new(big.Int)
		if sum != "" {
			sum_big, is_ok = new(big.Int).SetString(sum, 10)
			if !is_ok {
				return nil, fmt.Errorf("max error")
			}
		}

		return sum_big, nil
	}

	return nil, ErrNotFound
}

func (c *DBClient) FindSwapDrc20AddressInfoByTick(tx *sql.Tx, tick string, address string) (*big.Int, error) {
	query := "SELECT amt_sum  FROM drc20_address_info WHERE tick = ? and receive_address = ?"
	rows, err := tx.Query(query, tick, address)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {

		var sum string
		err := rows.Scan(&sum)

		if err != nil {
			return nil, err
		}

		is_ok := false
		sum_big := new(big.Int)
		if sum != "" {
			sum_big, is_ok = new(big.Int).SetString(sum, 10)
			if !is_ok {
				return nil, fmt.Errorf("max error")
			}
		}

		return sum_big, nil
	}

	return nil, ErrNotFound
}

func (c *DBClient) FindCardinalsInfoNewByNumber(number int64) ([]*utils.Cardinals, error) {
	query := "SELECT order_id, p, op, tick, amt, max_, lim_, repeat_mint,  drc20_tx_hash, block_hash, block_number, block_confirmations, receive_address,  create_date FROM cardinals_info where block_number > ? "
	rows, err := c.SqlDB.Query(query, number)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	cards := []*utils.Cardinals{}
	for rows.Next() {
		card := &utils.Cardinals{}
		var max, amt, lim string

		err := rows.Scan(&card.OrderId, &card.P, &card.Op, &card.Tick, &amt, &max, &lim, &card.Repeat, &card.Drc20TxHash, &card.BlockHash, &card.BlockNumber, &card.BlockConfirmations, &card.ReceiveAddress, &card.CreateDate)
		for err != nil {
			return nil, err
		}

		is_ok := false
		max_big := new(big.Int)
		if max != "" {
			max_big, is_ok = new(big.Int).SetString(max, 10)
			if !is_ok {
				return nil, fmt.Errorf("max error")
			}
		}
		card.Max = max_big

		amt_big := new(big.Int)
		if amt != "" {
			amt_big, is_ok = new(big.Int).SetString(amt, 10)
			if !is_ok {
				return nil, fmt.Errorf("amt error")
			}
		}
		card.Amt = amt_big

		lim_big := new(big.Int)
		if lim != "" {
			lim_big, is_ok = new(big.Int).SetString(lim, 10)
			if !is_ok {
				return nil, fmt.Errorf("lim error")
			}
		}
		card.Lim = lim_big

		cards = append(cards, card)
	}
	return cards, nil
}

func (c *DBClient) FindCardinalsInfoNewByDrc20Hash(drc20Hash string) (*utils.Cardinals, error) {
	query := "SELECT order_id, p, op, tick, amt, max_, lim_, repeat_mint,  drc20_tx_hash, block_hash, block_number, block_confirmations, receive_address,  create_date, to_address FROM cardinals_info where drc20_tx_hash = ?"
	rows, err := c.SqlDB.Query(query, drc20Hash)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		card := &utils.Cardinals{}
		var max, amt, lim string

		err := rows.Scan(&card.OrderId, &card.P, &card.Op, &card.Tick, &amt, &max, &lim, &card.Repeat, &card.Drc20TxHash, &card.BlockHash, &card.BlockNumber, &card.BlockConfirmations, &card.ReceiveAddress, &card.CreateDate, &card.ToAddress)
		for err != nil {
			return nil, err
		}

		is_ok := false
		max_big := new(big.Int)
		if max != "" {
			max_big, is_ok = new(big.Int).SetString(max, 10)
			if !is_ok {
				return nil, fmt.Errorf("max error")
			}
		}
		card.Max = max_big

		amt_big := new(big.Int)
		if amt != "" {
			amt_big, is_ok = new(big.Int).SetString(amt, 10)
			if !is_ok {
				return nil, fmt.Errorf("amt error")
			}
		}
		card.Amt = amt_big

		lim_big := new(big.Int)
		if lim != "" {
			lim_big, is_ok = new(big.Int).SetString(lim, 10)
			if !is_ok {
				return nil, fmt.Errorf("lim error")
			}
		}
		card.Lim = lim_big
		return card, nil
	}
	return nil, nil
}

func (c *DBClient) FindDrc20All(filter *utils.Drc20Params) ([]*FindDrc20AllResult, int64, error) {
	query := "SELECT di.tick AS ticker, di.amt_sum, di.max_, di.lim_, di.transactions, COUNT( ci.tick = di.tick ) AS Holders, di.create_date AS DeployTime, di.drc20_tx_hash, di.logo, di.introduction, di.is_check FROM drc20_address_info AS ci RIGHT JOIN drc20_info AS di ON ci.tick = di.tick  GROUP BY di.tick ORDER BY DeployTime DESC  LIMIT ? OFFSET ?"
	rows, err := c.SqlDB.Query(query, filter.Limit, filter.OffSet)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	var results []*FindDrc20AllResult
	for rows.Next() {

		result := &FindDrc20AllResult{}
		var max, amt, lim string
		err := rows.Scan(&result.Tick, &amt, &max, &lim, &result.Transactions, &result.Holders, &result.DeployTime, &result.Inscription, &result.Logo, &result.Introduction, &result.IsCheck)
		if err != nil {
			return nil, 0, err
		}

		is_ok := false
		max_big := new(big.Int)
		if max != "" {
			max_big, is_ok = new(big.Int).SetString(max, 10)
			if !is_ok {
				return nil, 0, fmt.Errorf("max error")
			}
		}
		result.MaxAmt = max_big

		amt_big := new(big.Int)
		if amt != "" {
			amt_big, is_ok = new(big.Int).SetString(amt, 10)
			if !is_ok {
				return nil, 0, fmt.Errorf("amt error")
			}
		}
		result.MintAmt = amt_big

		Lim, err := utils.ConvetStr(lim)
		if err != nil {
			return nil, 0, err
		}
		result.Lim = Lim

		if result.IsCheck == 0 {
			de := ""
			result.Logo = &de
			result.Introduction = &de
			result.WhitePaper = &de
			result.Official = &de
			result.Telegram = &de
			result.Discorad = &de
			result.Twitter = &de
			result.Facebook = &de
			result.Github = &de
		}

		results = append(results, result)
	}

	query1 := "SELECT COUNT(tick) AS UniqueTicks FROM drc20_info "

	rows1, err := c.SqlDB.Query(query1)
	if err != nil {
		return nil, 0, err
	}

	defer rows1.Close()
	total := int64(0)
	if rows1.Next() {
		rows1.Scan(&total)
	}
	return results, total, nil
}

func (c *DBClient) FindDrc20ByTick(tick string) (*FindDrc20AllResult, error) {
	query := "SELECT     di.tick AS ticker,     di.amt_sum,     di.max_ AS max_,     di.transactions AS Transactions,     di.update_date AS LastMintTime,     COUNT(CASE WHEN ci.tick = di.tick THEN 1 ELSE NULL END) AS Holders,     di.create_date AS DeployTime,     di.lim_ AS lim_,     di.dec_ AS dec_,     di.receive_address, di.drc20_tx_hash AS drc20_tx_hash_i0, di.logo, di.introduction, di.white_paper, di.official, di.telegram, di.discorad, di.twitter, di.facebook, di.github, di.is_check   FROM     drc20_address_info AS ci     RIGHT JOIN drc20_info AS di ON ci.tick = di.tick WHERE     di.tick = ? GROUP BY di.tick"
	rows, err := c.SqlDB.Query(query, tick)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {

		result := &FindDrc20AllResult{}
		var max, amt, lim *string
		err := rows.Scan(&result.Tick, &amt, &max, &result.Transactions, &result.LastMintTime, &result.Holders, &result.DeployTime, &lim, &result.Dec, &result.DeployBy, &result.Inscription, &result.Logo, &result.Introduction, &result.WhitePaper, &result.Official, &result.Telegram, &result.Discorad, &result.Twitter, &result.Facebook, &result.Github, &result.IsCheck)
		if err != nil {
			return nil, err
		}

		is_ok := false
		max_big := new(big.Int)
		if max != nil {
			max_big, is_ok = new(big.Int).SetString(*max, 10)
			if !is_ok {
				return nil, fmt.Errorf("max error")
			}
		}
		result.MaxAmt = max_big

		amt_big := new(big.Int)
		if amt != nil {
			amt_big, is_ok = new(big.Int).SetString(*amt, 10)
			if !is_ok {
				return nil, fmt.Errorf("amt error")
			}
		}
		result.MintAmt = amt_big

		lim_big := new(big.Int)
		if lim != nil {
			lim_big, is_ok = new(big.Int).SetString(*lim, 10)
			if !is_ok {
				return nil, fmt.Errorf("lim error")
			}
		}
		result.Lim = lim_big

		if result.IsCheck == 0 {
			de := ""
			result.Logo = &de
			result.Introduction = &de
			result.WhitePaper = &de
			result.Official = &de
			result.Telegram = &de
			result.Discorad = &de
			result.Twitter = &de
			result.Facebook = &de
			result.Github = &de
		}
		return result, nil
	}
	return nil, nil
}

func (c *DBClient) FindDrc20HoldersByTick(tick string, limit, offset int64) ([]*FindDrc20HoldersResult, int64, error) {
	query := "SELECT amt_sum, receive_address FROM drc20_address_info WHERE tick = ? ORDER BY CAST(amt_sum AS UNSIGNED) DESC LIMIT ? OFFSET ? ;"
	rows, err := c.SqlDB.Query(query, tick, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	var results []*FindDrc20HoldersResult
	var amt string
	for rows.Next() {

		result := &FindDrc20HoldersResult{}
		err := rows.Scan(&amt, &result.Address)
		if err != nil {
			return nil, 0, err
		}

		Amt, err := utils.ConvetStr(amt)
		if err != nil {
			return nil, 0, err
		}
		result.Amt = Amt
		results = append(results, result)
	}

	query1 := "SELECT count(receive_address) FROM drc20_address_info WHERE tick = ?"
	rows1, err := c.SqlDB.Query(query1, tick)
	if err != nil {
		return nil, 0, err
	}

	defer rows1.Close()
	total := int64(0)
	if rows1.Next() {
		rows1.Scan(&total)
	}

	return results, total, nil
}

func (c *DBClient) FindDrc20AllByAddress(receive_address string, limit, offset int64) ([]*FindDrc20AllByAddressResult, int64, error) {
	query := "SELECT tick, amt_sum FROM drc20_address_info where receive_address = ? and amt_sum != '0' LIMIT ? OFFSET ?;"
	rows, err := c.SqlDB.Query(query, receive_address, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	var results []*FindDrc20AllByAddressResult

	for rows.Next() {

		result := &FindDrc20AllByAddressResult{}
		var amt string

		err := rows.Scan(&result.Tick, &amt)
		if err != nil {
			return nil, 0, err
		}

		amt_big := new(big.Int)
		is_ok := false
		if amt != "" {
			amt_big, is_ok = new(big.Int).SetString(amt, 10)
			if !is_ok {
				return nil, 0, fmt.Errorf("lim error")
			}
		}
		result.Amt = amt_big

		results = append(results, result)
	}

	query1 := "SELECT count(tick) FROM drc20_address_info where receive_address = ? and amt_sum != '0' "
	rows1, err := c.SqlDB.Query(query1, receive_address)
	if err != nil {
		return nil, 0, err
	}

	defer rows1.Close()
	total := int64(0)
	if rows1.Next() {
		rows1.Scan(&total)
	}

	return results, total, nil
}

func (c *DBClient) FindDrc20AllByAddressTick(receive_address, tick string) (*FindDrc20AllByAddressResult, error) {
	query := "SELECT tick, amt_sum FROM drc20_address_info where receive_address = ? and amt_sum != '0' and tick = ?"
	rows, err := c.SqlDB.Query(query, receive_address, tick)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {

		result := &FindDrc20AllByAddressResult{}
		var amt string

		err := rows.Scan(&result.Tick, &amt)
		if err != nil {
			return nil, err
		}

		amt_big := new(big.Int)
		is_ok := false
		if amt != "" {
			amt_big, is_ok = new(big.Int).SetString(amt, 10)
			if !is_ok {
				return nil, fmt.Errorf("lim error")
			}
		}
		result.Amt = amt_big
		return result, nil
	}

	return nil, nil
}

func (c *DBClient) FindOrders(receiveAddress string, limit, offset int64) ([]*OrderResult, int64, error) {
	query := "SELECT order_id, p, op, tick, max_, lim_, amt, fee_address,receive_address,  rate_fee, fee_tx_hash,  drc20_tx_hash, block_hash, repeat_mint, block_confirmations, create_date, order_status, to_address  FROM cardinals_info where receive_address = ? or to_address = ?  order by update_date desc LIMIT ? OFFSET ?"

	rows, err := c.SqlDB.Query(query, receiveAddress, receiveAddress, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	var cards []*OrderResult
	for rows.Next() {
		card := &OrderResult{}
		var max *string
		var lim *string
		var amt *string
		var fee *string

		err := rows.Scan(&card.OrderId, &card.P, &card.Op, &card.Tick, &max, &lim, &amt, &card.FeeAddress, &card.ReceiveAddress, &fee, &card.FeeTxHash, &card.Drc20TxHash, &card.BlockHash, &card.Repeat, &card.BlockConfirmations, &card.CreateDate, &card.OrderStatus, &card.ToAddress)
		if err != nil {
			return nil, 0, err
		}

		card.Max, _ = utils.ConvetStr(*max)
		card.Amt, _ = utils.ConvetStr(*amt)
		card.Lim, _ = utils.ConvetStr(*lim)
		card.RateFee, _ = utils.ConvetStr(*fee)

		cards = append(cards, card)
	}

	query1 := "SELECT count(order_id)  FROM cardinals_info where receive_address = ? and is_del = 0"

	rows1, err := c.SqlDB.Query(query1, receiveAddress)
	if err != nil {
		return nil, 0, err
	}

	defer rows1.Close()
	total := int64(0)
	if rows1.Next() {
		rows1.Scan(&total)
	}

	return cards, total, nil
}
