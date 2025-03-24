package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/bianavic/fullcycle_clean-architecture/internal/entity"

	"github.com/stretchr/testify/suite"
	// sqlite3
	_ "github.com/mattn/go-sqlite3"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)
	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) SetupTest() {
	_, err := suite.Db.Exec("CREATE TABLE IF NOT EXISTS orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id))")
	suite.NoError(err)
}

func (suite *OrderRepositoryTestSuite) TearDownTest() {
	if suite.Db != nil {
		_, err := suite.Db.Exec("DELETE FROM orders")
		suite.NoError(err)
	}
}

func (suite *OrderRepositoryTestSuite) TearDownSuite() {
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10.0, 2.0, time.Now())
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.QueryRow("Select id, price, tax, final_price from orders where id = ?", order.ID).
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	suite.NoError(err)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}

func (suite *OrderRepositoryTestSuite) TestList() {
	order1, err := entity.NewOrder("123", 10.0, 2.0, time.Now())
	suite.NoError(err)
	suite.NoError(order1.CalculateFinalPrice())

	order2, err := entity.NewOrder("456", 20.0, 3.0, time.Now())
	suite.NoError(err)
	suite.NoError(order2.CalculateFinalPrice())

	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order1)
	suite.NoError(err)
	err = repo.Save(order2)
	suite.NoError(err)

	// Test List method
	orders, err := repo.List()
	suite.NoError(err)
	suite.Len(orders, 2)

	suite.Equal(order1.ID, orders[0].ID)
	suite.Equal(order1.Price, orders[0].Price)
	suite.Equal(order1.Tax, orders[0].Tax)
	suite.Equal(order1.FinalPrice, orders[0].FinalPrice)

	suite.Equal(order2.ID, orders[1].ID)
	suite.Equal(order2.Price, orders[1].Price)
	suite.Equal(order2.Tax, orders[1].Tax)
	suite.Equal(order2.FinalPrice, orders[1].FinalPrice)
}

func (suite *OrderRepositoryTestSuite) TestListQueryError() {
	_, err := suite.Db.Exec("DROP TABLE orders")
	suite.NoError(err)

	repo := NewOrderRepository(suite.Db)
	_, err = repo.List()

	suite.Error(err)

	_, err = suite.Db.Exec("CREATE TABLE IF NOT EXISTS orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id))")
	suite.NoError(err)
}

func (suite *OrderRepositoryTestSuite) TestListScanError() {
	_, err := suite.Db.Exec("INSERT INTO orders (id, price, tax, final_price, created_at) VALUES (?, ?, ?, ?, ?)", "invalid", "invalid", "invalid", "invalid", time.Now())
	suite.NoError(err)

	repo := NewOrderRepository(suite.Db)
	_, err = repo.List()

	suite.Error(err)

	_, err = suite.Db.Exec("DELETE FROM orders WHERE id = ?", "invalid")
	suite.NoError(err)
}
