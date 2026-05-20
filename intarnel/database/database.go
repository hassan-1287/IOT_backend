package database
//*هذا الملف والبكج مسوؤل عن الاتصال بقاعدة البيانات
import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)
//*داله الاتصال في قاعدة البيانات
func Connect(databaseurl string) (*pgxpool.Pool, error) {
	var ctx context.Context = context.Background()
	var config *pgxpool.Config
	var err error

	config, err = pgxpool.ParseConfig(databaseurl)
	if err != nil {
		log.Panicf("Unable to pares DATABASE_URL: %v", err)
		return nil, err
	}

	var pool *pgxpool.Pool
	pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Printf("Unable to create connction pool: %v", err)
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Printf("Unable to ping database:%v", err)
		pool.Close()
		return nil, err
	}
	return pool, nil

}
