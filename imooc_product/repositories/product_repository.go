package repositories

import (
	"database/sql"
	"imooc_product/common"
	"imooc_product/datamodels"
	"strconv"
)

// 第一步，先开发对应的接口
// 第二步，实现定义的接口

type IProduct interface {
	// 连接数据
	Conn() (error)
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) (error)
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
	SubProductNum(productID int64) error
}

type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{table, db}
}

// 数据库连接
func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table == "" {
		p.table = "product"
	}
	return
}

// 插入
func (p *ProductManager) Insert(product *datamodels.Product) (productId int64, err error) {
	// 1. 判断连接是否存在
	if err = p.Conn(); err != nil {
		return
	}
	// 2. 准备sql
	sql := "INSERT product SET productName=?,productNum=?,productImage=?,productUrl=?"
	stmt, errSql := p.mysqlConn.Prepare(sql)
	if errSql != nil {
		return 0, errSql
	}
	// 3. 传入参数
	result, errStmt := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if errStmt != nil {
		return 0, errStmt
	}
	return result.LastInsertId()
}

// 商品删除
func (p *ProductManager) Delete(productID int64) bool {
	// 1. 判断连接是否存在
	if err := p.Conn(); err != nil {
		return false
	}
	// 2. 准备sql
	sql := "DELETE FROM product WHERE ID=?"
	stmt, errSql := p.mysqlConn.Prepare(sql)
	if errSql != nil {
		return false
	}
	// 3. 传入参数
	_, errStmt := stmt.Exec(productID)
	if errStmt != nil {
		return false
	}
	return true
}

// 商品更新
func (p *ProductManager) Update(product *datamodels.Product) error {
	// 1. 判断连接是否存在
	if err := p.Conn(); err != nil {
		return err
	}
	// 2. 准备sql
	sql := "UPDATE product SET productName=?,productNum=?,productImage=?,productUrl=? WHERE ID=" + strconv.FormatInt(product.ID, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	// 3. 传入参数
	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return err
	}
	return nil
}

// 根据商品ID查询对应商品
func (p *ProductManager) SelectByKey(productID int64) (productResult *datamodels.Product, err error) {
	// 1. 判断连接师傅存在
	if err = p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}
	sql := "select * from " + p.table + " where ID=" + strconv.FormatInt(productID, 10)
	row, errRow := p.mysqlConn.Query(sql)
	if errRow != nil {
		return &datamodels.Product{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}
	productResult = &datamodels.Product{}
	common.DataToStructByTagSql(result, productResult)
	return
}

// 获取所有商品
func (p *ProductManager) SelectAll() (ProductArray []*datamodels.Product, errProduct error) {
	// 1. 判断连接是否存在
	if err := p.Conn(); err != nil {
		return nil, err
	}
	sql := "select * from " + p.table
	rows, err := p.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}
	for _, v := range result {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		ProductArray = append(ProductArray, product)
	}
	return
}

func (p *ProductManager) SubProductNum(productID int64) error {
	if err := p.Conn(); err != nil {
		return err
	}
	sql := "update " + p.table + " set " + "productNum=productNum-1 where ID=" + strconv.FormatInt(productID, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}
