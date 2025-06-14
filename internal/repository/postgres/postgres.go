package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/config"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

type PostgresDb struct {
	db *sql.DB
}

func ConnectPostgresDb() (*PostgresDb, error) {
	cfg := config.NewPostgresConfig()
	ctx := context.Background()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/shopping_lists_db?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	sqlDb, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Println("Failed to connect to postgres: ", err)
		return nil, err
	}
	err = sqlDb.PingContext(ctx)
	if err != nil {
		log.Println("Failed to connect to postgres: ", err)
		return nil, err
	}
	return &PostgresDb{db: sqlDb}, nil
}

func ConnectPostgresDbWithRetries(maxAttempts int, delay time.Duration) (*PostgresDb, error) {
	var db *sql.DB
	var err error
	cfg := config.NewPostgresConfig()
	ctx := context.Background()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/shopping_lists_db?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = sql.Open("pgx", dsn)
		if err != nil {
			log.Printf("Attempt %d: connection error: %v", attempt, err)
			time.Sleep(delay)
			continue
		}

		err = db.PingContext(ctx)
		if err == nil {
			log.Printf("Successfully connected after %d attempts", attempt)
			return &PostgresDb{db: db}, nil
		}

		log.Printf("Attempt %d: ping failed: %v", attempt, err)
		db.Close()
		time.Sleep(delay)
	}

	return nil, err
}

func (p *PostgresDb) Migrate(ctx context.Context, migrate string) (err error) {
	//	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("cannot set dialect: %w", err)
	}

	if err := goose.Up(p.db, migrate); err != nil {
		return fmt.Errorf("cannot do up migration: %w", err)
	}

	return nil
}

func (p *PostgresDb) Close() error {
	return p.db.Close()
}

func (p *PostgresDb) AddShoppingList(ctx context.Context, sl *model.ShoppingList) error {
	sqlCreatedAt := sql.NullTime{Time: sl.CreatedAt.AsTime(), Valid: true}
	sqlUpdatedAt := sql.NullTime{Time: sl.CreatedAt.AsTime(), Valid: true}
	query := `INSERT INTO lists(id, title, user_id, created_at, updated_at, state)
		values($1, $2, $3, $4, $5, $6)`
	_, err := p.db.ExecContext(ctx, query, sl.Id, sl.Title, sl.UserId, sqlCreatedAt, sqlUpdatedAt, sl.State)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Inserted shopping list: ", sl.String())
	return nil
}

func (p *PostgresDb) GetSlById(ctx context.Context, id string) (*model.ShoppingList, error) {
	row := p.db.QueryRowContext(ctx, `
		SELECT id, title, user_id, created_at, updated_at, state FROM lists WHERE id = $1
	`, id)
	var sl model.ShoppingList
	var sqlCreatedAt time.Time
	var sqlUpdatedAt time.Time
	err := row.Scan(&sl.Id, &sl.Title, &sl.UserId, &sqlCreatedAt, &sqlUpdatedAt, &sl.State)

	if errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return nil, repository.ErrNotFound
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	sl.CreatedAt = timestamppb.New(sqlCreatedAt)
	sl.UpdatedAt = timestamppb.New(sqlUpdatedAt)
	log.Println("Get shopping list: ", sl)
	return &sl, nil
}

func (p *PostgresDb) GetSls(ctx context.Context) ([]*model.ShoppingList, error) {
	rows, err := p.db.QueryContext(ctx, `
		SELECT id, title, user_id, created_at, updated_at, state FROM lists
	`)
	if err != nil {
		log.Println("cannot select: ", err)
		return nil, err
	}
	defer rows.Close()

	var sl []*model.ShoppingList

	for rows.Next() {
		var s model.ShoppingList
		var sqlCreatedAt time.Time
		var sqlUpdatedAt time.Time
		if err = rows.Scan(
			&s.Id,
			&s.Title,
			&s.UserId,
			&sqlCreatedAt,
			&sqlUpdatedAt,
			&s.State,
		); err != nil {
			return nil, fmt.Errorf("cannot scan: %w", err)
		}
		s.CreatedAt = timestamppb.New(sqlCreatedAt)
		s.UpdatedAt = timestamppb.New(sqlUpdatedAt)

		sl = append(sl, &s)
	}
	log.Println("Get shopping lists: ", sl)
	return sl, rows.Err()
}

func (p *PostgresDb) UpdateSl(ctx context.Context, id string, sl *model.ShoppingList) error {
	log.Println(sl)
	sqlUpdatedAt := sql.NullTime{Time: sl.UpdatedAt.AsTime(), Valid: true}
	query := `UPDATE lists SET title=$1, user_id=$2, updated_at=$3, state=$4 WHERE id=$5`
	_, err := p.db.ExecContext(ctx, query, sl.Title, sl.UserId, &sqlUpdatedAt, sl.State, id)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return repository.ErrNotFound
	} else if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Update shopping list: ", sl.String())
	return nil
}

func (p *PostgresDb) DeleteSlById(ctx context.Context, id string) error {
	query := `DELETE FROM lists WHERE id=$1`
	_, err := p.db.ExecContext(ctx, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return repository.ErrNotFound
	} else if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Delete shopping list: ", id)
	return nil
}

func (p *PostgresDb) AddItem(ctx context.Context, item *model.Item) error {
	sqlCreatedAt := sql.NullTime{Time: item.CreatedAt.AsTime(), Valid: true}
	sqlUpdatedAt := sql.NullTime{Time: item.CreatedAt.AsTime(), Valid: true}
	query := `INSERT INTO items(id, title, comment, is_done, user_id, created_at, updated_at, shopping_list_id)
		values($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := p.db.ExecContext(ctx, query, item.Id, item.Title, item.Comment, item.IsDone, item.UserId, sqlCreatedAt, sqlUpdatedAt, item.ShoppingListId)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Inserted item: ", item.String())
	return nil
}

func (p *PostgresDb) GetItemById(ctx context.Context, id string) (*model.Item, error) {
	row := p.db.QueryRowContext(ctx, `
		SELECT id, title, comment, is_done, user_id, created_at, updated_at, shopping_list_id FROM items WHERE id = $1
	`, id)
	var item model.Item
	var sqlCreatedAt time.Time
	var sqlUpdatedAt time.Time
	err := row.Scan(
		&item.Id,
		&item.Title,
		&item.Comment,
		&item.IsDone,
		&item.UserId,
		&sqlCreatedAt,
		&sqlUpdatedAt,
		&item.ShoppingListId,
	)

	if errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return nil, repository.ErrNotFound
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	item.CreatedAt = timestamppb.New(sqlCreatedAt)
	item.UpdatedAt = timestamppb.New(sqlUpdatedAt)
	log.Println("Get item: ", item)
	return &item, nil
}

func (p *PostgresDb) GetItems(ctx context.Context) ([]*model.Item, error) {
	rows, err := p.db.QueryContext(ctx, `
		SELECT id, title, comment, is_done, user_id, created_at, updated_at, shopping_list_id FROM items
	`)
	if err != nil {
		log.Println("cannot select: ", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)
	var items []*model.Item
	for rows.Next() {
		var i model.Item
		var sqlCreatedAt time.Time
		var sqlUpdatedAt time.Time
		if err = rows.Scan(
			&i.Id,
			&i.Title,
			&i.Comment,
			&i.IsDone,
			&i.UserId,
			&sqlCreatedAt,
			&sqlUpdatedAt,
			&i.ShoppingListId,
		); err != nil {
			return nil, fmt.Errorf("cannot scan: %w", err)
		}
		i.CreatedAt = timestamppb.New(sqlCreatedAt)
		i.UpdatedAt = timestamppb.New(sqlUpdatedAt)
		items = append(items, &i)
	}
	log.Println("Get items: ", items)
	return items, nil
}

func (p *PostgresDb) UpdateItem(ctx context.Context, id string, item *model.Item) error {
	sqlUpdatedAt := sql.NullTime{Time: item.UpdatedAt.AsTime(), Valid: true}
	query := `UPDATE items SET title=$1, comment=$2, is_done=$3, user_id=$4, updated_at=$5  where id=$6`
	_, err := p.db.ExecContext(ctx, query, item.Title, item.Comment, item.IsDone, item.UserId, sqlUpdatedAt, id)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return repository.ErrNotFound
	} else if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Update item: ", item.String())
	return nil
}

func (p *PostgresDb) DeleteItemById(ctx context.Context, id string) error {
	query := `DELETE FROM items where id=$1`
	_, err := p.db.ExecContext(ctx, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return repository.ErrNotFound
	} else if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Delete item: ", id)
	return nil
}

func (p *PostgresDb) GetItemsBySLId(ctx context.Context, ShoppingListId string) ([]*model.Item, error) {
	rows, err := p.db.QueryContext(ctx, `
		SELECT id, title, comment, is_done, user_id, created_at, updated_at, shopping_list_id FROM items WHERE shopping_list_id=$1
	`, ShoppingListId)
	if err != nil {
		log.Println("cannot select: ", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)
	var items []*model.Item
	for rows.Next() {
		var i model.Item
		var sqlCreatedAt time.Time
		var sqlUpdatedAt time.Time
		if err = rows.Scan(
			&i.Id,
			&i.Title,
			&i.Comment,
			&i.IsDone,
			&i.UserId,
			&sqlCreatedAt,
			&sqlUpdatedAt,
			&i.ShoppingListId,
		); err != nil {
			return nil, fmt.Errorf("cannot scan: %w", err)
		}
		i.CreatedAt = timestamppb.New(sqlCreatedAt)
		i.UpdatedAt = timestamppb.New(sqlUpdatedAt)
		items = append(items, &i)
	}
	log.Println("Get items: ", items)
	return items, nil
}

func (p *PostgresDb) CreateUser(ctx context.Context, user *model.User) error {
	row := p.db.QueryRowContext(ctx, `
		SELECT id, name, password FROM users WHERE name = $1
	`, user.Name)
	var u model.User
	err := row.Scan(&u.Id, &u.Name, &u.Password, &u.State)
	if errors.Is(err, sql.ErrNoRows) {
		query := `INSERT INTO users(id, name, password, state)
		values($1, $2, $3, $4)`
		_, err := p.db.ExecContext(ctx, query, user.Id, user.Name, user.Password, user.State)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("Inserted user: ", user.Name)
		return nil
	} else if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (p *PostgresDb) Login(ctx context.Context, user *model.User) (string, error) {
	row := p.db.QueryRowContext(ctx, `
		SELECT id, name, password FROM users WHERE name = $1
	`, user.Name)
	var u []model.User
	err := row.Scan(u)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return "", repository.ErrNotFound
	} else if err != nil {
		log.Println(err)
		return "", err
	}
	if u[0].State != 2 {
		return "", errors.New("USER NOT ACTIVE")
	} else if u[0].Password == user.Password && u[0].Name == user.Name {
		secretKey := []byte(config.Cfg.Secret)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   user.Id,
			"name": user.Name,
			"exp":  time.Now().Add(time.Hour * 24).Unix(), // Срок действия — 24 часа
		})
		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			log.Println(err)
			return "", err
		}
		return tokenString, nil
	}
	log.Println("Login user: ", user.Name)
	return "", repository.ErrNotFound
}
