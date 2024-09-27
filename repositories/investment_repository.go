package repositories

import (
	"context"
	"database/sql"
	db "stt/database/postgres/sqlc"
	"stt/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type investmentRepository struct {
	queries *db.Queries
}

var pool *sql.DB

func InitInvestmentRepository(q *db.Queries) domain.IInvestmentRepository {
	return &investmentRepository{
		queries: q,
	}
}

// Create implements domain.IInvestmentRepository.
func (i *investmentRepository) Create(c context.Context, investmentParam db.CreateInvestmentParams) (db.Investment, error) {
	return i.queries.CreateInvestment(c, investmentParam)
}

// Delete implements domain.IInvestmentRepository.
func (i *investmentRepository) Delete(c context.Context, id int32) {
	panic("unimplemented")
}

// GetAll implements domain.IInvestmentRepository.
func (i *investmentRepository) GetAll(c context.Context) {
	// rows, err := i.database.Query(c, "SELECT * FROM accounts")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer rows.Close()

	// // Process rows
	// for rows.Next() {
	// 	var id int64
	// 	var owner string
	// 	var channel_name string
	// 	var balance float32
	// 	var buy_fee float32
	// 	var sell_fee float32
	// 	var currency string
	// 	var created_at time.Time
	// 	err := rows.Scan(&id, &owner, &channel_name, &balance, &buy_fee, &sell_fee, &currency, &created_at)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("ID: %d, AccountName: %s\n", id, owner)
	// }

}

// GetById implements domain.IInvestmentRepository.
func (i *investmentRepository) GetById(c context.Context, id int32) {
	panic("unimplemented")
}

// Update implements domain.IInvestmentRepository.
func (i *investmentRepository) Update(c context.Context, investment *db.Investment) {
	panic("unimplemented")
}
