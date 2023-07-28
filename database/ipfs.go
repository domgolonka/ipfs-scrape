package database

import (
	"database/sql"
	"domgolonka/blockparty/model"
	"fmt"
	"strings"
)

// InsertIpfs inserts a new ipfs' into the database
func (c *Client) InsertIpfs(ipfs []*model.Ipfs) error {
	sqlStr := "INSERT INTO ipfs(image, description, name, cid) VALUES "
	vals := []interface{}{}
	for i, row := range ipfs {
		sqlStr += fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4) + ","
		vals = append(vals, row.Image, row.Description, row.Name, row.Cid)
	}
	//trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	//prepare the statement
	p, err := c.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer p.Close()
	_, err = p.Exec(vals...)
	if err != nil {
		return err
	}
	return nil
}

// SelectIpfs returns all the ipfs
func (c *Client) SelectIpfs() ([]model.Ipfs, error) {
	data := make([]model.Ipfs, 0)
	// We assign the result to 'rows'
	allRows, err := c.Query("SELECT * FROM ipfs")
	if err != nil {
		return nil, err
	}
	defer allRows.Close()

	return c.allRows(allRows, data)
}

// SelectIpfsByCid returns all the ipfs by cid
func (c *Client) SelectIpfsByCid(cid string) ([]model.Ipfs, error) {
	data := make([]model.Ipfs, 0)
	// We assign the result to 'rows'
	allRows, err := c.Query("SELECT * FROM ipfs WHERE = $1", cid)
	if err != nil {
		return nil, err
	}
	defer allRows.Close()

	return c.allRows(allRows, data)
}

func (c *Client) allRows(rows *sql.Rows, data []model.Ipfs) ([]model.Ipfs, error) {
	// we loop through the values of rows
	for rows.Next() {
		ipfs := model.Ipfs{}
		err := rows.Scan(&ipfs.Description, &ipfs.Name, &ipfs.Image)
		if err != nil {
			return data, err
		}
		data = append(data, ipfs)
	}
	if err := rows.Err(); err != nil {
		return data, err
	}
	return data, nil
}
