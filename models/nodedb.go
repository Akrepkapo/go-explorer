package models

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/IBAX-io/go-explorer/conf"
	"github.com/IBAX-io/go-explorer/storage"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var nodeErr = errors.New("node params error")

func GetNodeRows(node int64, tableName string) (int64, error) {
	var count int64

	db := GetNodedb(node)
	if db != nil {
		err := db.Table(tableName).Count(&count).Error
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}

func GetNodeTables(node int64) (int64, error) {
	var count int64
	db := GetNodedb(node)
	if db != nil {
		err := db.Raw("SELECT count(*) FROM pg_tables WHERE schemaname='public'").Count(&count).Error
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}

func GetNodeAllColumnTypes(node int64, tblname string) ([]map[string]string, error) {
	return GetAll(node, `SELECT column_name, data_type,column_default 
		FROM information_schema.columns
		WHERE table_name = ?
		ORDER BY ordinal_position ASC`, -1, tblname)
}

func GetNodeTableOrder(order, tablname string) string {
	defaultorder := " asc"
	if order != "" {
		if strings.Contains(order, "desc") || strings.Contains(order, "DESC") {
			defaultorder = "desc"
		}
	}

	switch tablname {
	case "confirmations", "info_block":
		order = "block_id " + defaultorder
		break
	case "stop_daemons":
		order = "stop_time " + defaultorder
		break
	case "install":
		order = "progress " + defaultorder
		break
	case "transactions", "queue_tx", "log_transactions", "queue_blocks", "transactions_status", "transactions_attempts":
		order = "hash " + defaultorder
		break
	}

	return order
}

func GetNodedb(node int64) *gorm.DB {
	FullNodedb := conf.GetFullNodesDbConn()
	dlen := len(FullNodedb)
	for i := 0; i < dlen; i++ {
		if int64(dlen) == node {
			node = node - 1
		}
		if node == FullNodedb[i].NodePosition {
			if FullNodedb[i].Enable {
				return FullNodedb[i].DBConn
			} else {
					log.Info("node dblink ok", FullNodedb[i].NodeName, FullNodedb[i].NodePosition)
					FullNodedb[i].Enable = true
					FullNodedb[i].DBConn = db
					return FullNodedb[i].DBConn
				}
			}
		}
	}
	return nil
}

func GetNodeALLTable(node int64, ids, icount int, order, table string) (int64, []map[string]string, error) {
	var count int64
	db := GetNodedb(node)
	ns := "%" + table + "%"
	err := db.Table("pg_tables").Where("schemaname='public' and tablename like ?", ns).Count(&count).Error
	if err == gorm.ErrRecordNotFound {
		return count, nil, nil
	}

	if ids < 1 || icount < 1 {
		return count, nil, nodeErr
	}

	rs, err := GetAll(node, fmt.Sprintf(`SELECT tablename FROM pg_tables WHERE schemaname='public' and tablename like %s order by %s offset %d`, "'%"+table+"%'", order, (ids-1)*icount), int(icount))
	return count, rs, err
}

// GetAll returns all transaction
func GetAll(node int64, query string, countRows int, args ...interface{}) ([]map[string]string, error) {
	db := GetNodedb(node)
	if db != nil {
		return GetAllTransaction(db, query, countRows, args)
	}

	return nil, nodeErr
}

// GetAllTransaction is retrieve all query result rows
func GetAllTransaction(db *gorm.DB, query string, countRows int, args ...interface{}) ([]map[string]string, error) {
	var result []map[string]string
	rows, err := db.Raw(query, args...).Rows()
	if err != nil {
		return result, fmt.Errorf("%s in query %s %s", err, query, args)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	//columntypes, err1 := rows.ColumnTypes();
	if err != nil {
		return result, fmt.Errorf("%s in query %s %s", err, query, args)
	}

	columntypes, err1 := rows.ColumnTypes()
	if err1 != nil {
		return result, fmt.Errorf("%s in query %s %s", err1, query, args)
	}
	// Make a slice for the values
	values := make([][]byte /*sql.RawBytes*/, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	r := 0
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return result, fmt.Errorf("%s in query %s %s", err, query, args)
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		rez := make(map[string]string)
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				if columntypes[i].DatabaseTypeName() == "BYTEA" {
					value = hex.EncodeToString(col)
				} else {
					value = string(col)
				}
			}
			rez[columns[i]] = value
		}
		result = append(result, rez)
		r++
		if countRows != -1 && r >= countRows {
			break
		}
	}
	if err = rows.Err(); err != nil {
		return result, fmt.Errorf("%s in query %s %s", err, query, args)
	}
	return result, nil
}

func GetQueryexec(db *gorm.DB, query string, countRows int, args ...interface{}) (*[]map[string]string, error) {
	var result []map[string]string
	rows, err := db.Raw(query, args...).Rows()
	if err != nil {
		return &result, fmt.Errorf("%s in query %s %s", err, query, args)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return &result, fmt.Errorf("%s in query %s %s", err, query, args)
	}

	columntypes, err1 := rows.ColumnTypes()
	if err1 != nil {
		return &result, fmt.Errorf("%s in query %s %s", err1, query, args)
	}
	// Make a slice for the values
	values := make([][]byte /*sql.RawBytes*/, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	r := 0
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return &result, fmt.Errorf("%s in query %s %s", err, query, args)
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		rez := make(map[string]string)
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				if columntypes[i].DatabaseTypeName() == "BYTEA" {
					value = hex.EncodeToString(col)
				} else {
					value = string(col)
				}
				//value = string(col)
			}
			rez[columns[i]] = value
		}
		result = append(result, rez)
		r++
		if countRows != -1 && r >= countRows {
			break
		}
	}
	if err = rows.Err(); err != nil {
		return &result, fmt.Errorf("%s in query %s %s", err, query, args)
	}
	return &result, nil
}

// GetDB is returning gorm.DB
func GetDB(db *DbTransaction) *gorm.DB {
	if db != nil && db.conn != nil {
		return db.conn
	}
	return conf.GetDbConn().Conn()
}

func GetsqliteALLTable() (*[]map[string]string, error) {
	return GetQueryexec(conf.GetDbConn().Conn(), `select name from sqlite_master where type='table' order by name`, -1)
}

func GetsqliteblcokALLTable() (*[]map[string]string, error) {
	return GetQueryexec(conf.GetDbConn().Conn(), `select name from sqlite_master where type='table' order by name`, -1)
}
