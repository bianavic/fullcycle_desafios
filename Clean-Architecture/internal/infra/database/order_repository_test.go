package database

import (
	"database/sql"
	"testing"

	"github.com/bianavic/fullcycle_clean-architecture/internal/entity"

	// sqlite3
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
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
	_, err := suite.Db.Exec("CREATE TABLE IF NOT EXISTS orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.NoError(err)
}

func (suite *OrderRepositoryTestSuite) TearDownSuite() {
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestSave() {
	order, err := entity.NewOrder("123", 10.0, 2.0)
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

func (suite *OrderRepositoryTestSuite) TestListOrders() {

	suite.Run("Listar pedidos com sucesso", func() {
		_, err := suite.Db.Exec("DELETE FROM orders")
		suite.NoError(err)

		order1, _ := entity.NewOrder("123", 10.0, 2.0)
		order1.CalculateFinalPrice()
		order2, _ := entity.NewOrder("456", 20.0, 3.0)
		order2.CalculateFinalPrice()

		repo := NewOrderRepository(suite.Db)
		repo.Save(order1)
		repo.Save(order2)

		orders, err := repo.ListOrders()

		suite.NoError(err)
		suite.Len(orders, 2)
		suite.Equal(order1.ID, orders[0].ID)
		suite.Equal(order2.ID, orders[1].ID)
	})

	suite.Run("Erro ao listar quando tabela não existe", func() {
		_, err := suite.Db.Exec("DROP TABLE IF EXISTS orders")
		suite.NoError(err)

		repo := NewOrderRepository(suite.Db)
		_, err = repo.ListOrders()

		suite.Error(err)

		_, err = suite.Db.Exec(`CREATE TABLE orders (
            id varchar(255) NOT NULL PRIMARY KEY,
            price float NOT NULL,
            tax float NOT NULL,
            final_price float NOT NULL
        )`)
		suite.NoError(err)
	})

	suite.Run("Erro ao escanear dados inválidos", func() {
		_, err := suite.Db.Exec(
			"INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)",
			"invalid", "not-a-float", "not-a-float", "not-a-float")
		suite.NoError(err)

		repo := NewOrderRepository(suite.Db)
		_, err = repo.ListOrders()

		suite.Error(err)

		_, err = suite.Db.Exec("DELETE FROM orders WHERE id = ?", "invalid")
		suite.NoError(err)
	})
}
