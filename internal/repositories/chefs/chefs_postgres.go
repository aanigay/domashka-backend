package chefs

import (
	"context"
	dishEntity "domashka-backend/internal/entity/dishes"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"

	entity "domashka-backend/internal/entity/chefs"
	"domashka-backend/internal/utils/pointers"
	"domashka-backend/pkg/postgres"
)

type Repository struct {
	pg *postgres.Postgres
}

func New(pg *postgres.Postgres) *Repository {
	return &Repository{pg: pg}
}

func (r *Repository) GetChefByDishID(ctx context.Context, dishID int64) (*entity.Chef, error) {
	var chef entity.Chef
	err := r.pg.Pool.QueryRow(ctx, "SELECT id, name, small_image_url FROM chefs WHERE id = $1", dishID).Scan(
		&chef.ID,
		&chef.Name,
		&chef.SmallImageURL,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &chef, nil
}

func (r *Repository) GetChefByID(ctx context.Context, chefID int64) (*entity.Chef, error) {
	var chef entity.Chef
	err := r.pg.Pool.QueryRow(ctx, "SELECT id, name, image_url, small_image_url, description FROM chefs WHERE id = $1", chefID).Scan(
		&chef.ID,
		&chef.Name,
		&chef.ImageURL,
		&chef.SmallImageURL,
		&chef.Description,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &chef, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]entity.Chef, error) {
	var chefs []entity.Chef
	rows, err := r.pg.Pool.Query(ctx, "SELECT id, name, description, image_url, small_image_url FROM chefs")
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []entity.Chef{}, nil
		}
		return nil, err
	}
	for rows.Next() {
		var chef entity.Chef
		err := rows.Scan(
			&chef.ID,
			&chef.Name,
			&chef.Description,
			&chef.ImageURL,
			&chef.SmallImageURL,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		chefs = append(chefs, chef)
	}
	return chefs, nil
}

func (r *Repository) SaveChefAvatar(ctx context.Context, chefID int64, publicURL string) error {
	// Обновляем БД
	if _, err := r.pg.Pool.Exec(ctx,
		"UPDATE chefs SET image_url=$1 WHERE id=$2",
		publicURL, chefID,
	); err != nil {
		return fmt.Errorf("db update: %w", err)
	}

	return nil
}

func (r *Repository) SetSmallAvatar(ctx context.Context, chefID int64, publicURL string) error {
	// Обновляем БД
	if _, err := r.pg.Pool.Exec(ctx,
		"UPDATE chefs SET small_image_url=$1 WHERE id=$2",
		publicURL, chefID,
	); err != nil {
		return fmt.Errorf("db update: %w", err)
	}

	return nil
}

func (r *Repository) GetChefRatingByChefID(ctx context.Context, chefID int64) (*entity.Chef, error) {
	var chefRating ChefRating
	err := r.pg.Pool.QueryRow(ctx, "SELECT chef_id, rating, reviews_count FROM chef_ratings WHERE chef_id = $1", chefID).Scan(
		&chefRating.ChefID,
		&chefRating.Rating,
		&chefRating.ReviewsCount,
	)
	if err != nil {
		return nil, err
	}
	chef := &entity.Chef{Rating: pointers.To(chefRating.Rating), ReviewsCount: pointers.To(chefRating.ReviewsCount)}
	return chef, nil
}

func (r *Repository) GetTopChefs(ctx context.Context, limit int) ([]entity.Chef, error) {
	query := `
        SELECT 
            c.id,
            c.name,
            c.image_url,
            c.small_image_url,
            cr.rating,
            cr.reviews_count
        FROM 
            public.chefs c
        JOIN 
            public.chef_ratings cr ON c.id = cr.chef_id
        ORDER BY 
            cr.rating DESC
        LIMIT $1;
    `
	rows, err := r.pg.Pool.Query(ctx, query, limit)
	if errors.Is(err, pgx.ErrNoRows) {
		return []entity.Chef{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var chefs []entity.Chef
	for rows.Next() {
		var chef entity.Chef
		err := rows.Scan(
			&chef.ID,
			&chef.Name,
			&chef.ImageURL,
			&chef.SmallImageURL,
			&chef.Rating,
			&chef.ReviewsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		chefs = append(chefs, chef)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return chefs, nil
}

func (r *Repository) GetNearestChefs(ctx context.Context, lat, long float64, distance, limit int) ([]entity.Chef, error) {
	query := `
			SELECT 
				c.id,
            c.name,
            c.image_url,
            c.small_image_url,
			ST_Distance(ca.geom, ST_MakePoint($1, $2)::geography) AS distance,
    		ST_Y(ca.geom::geometry) AS lat,
			ST_X(ca.geom::geometry) AS lon
			FROM 
				chefs c
			JOIN 
				chef_addresses ca ON c.id = ca.chef_id
			WHERE 
				ST_DWithin(ca.geom, ST_MakePoint($1, $2)::geography, $3)
			ORDER BY 
				distance
			LIMIT $4;
		`

	// Превращаем в метры
	distance *= 1000

	var chefs []entity.Chef
	rows, err := r.pg.Pool.Query(ctx, query, long, lat, distance, limit)
	if errors.Is(err, pgx.ErrNoRows) {
		return []entity.Chef{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var chef entity.Chef
		var geo entity.Geo
		err := rows.Scan(
			&chef.ID,
			&chef.Name,
			&chef.ImageURL,
			&chef.SmallImageURL,
			&geo.Distance,
			&geo.Latitude,
			&geo.Longitude,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		chef.Geo = &geo
		chefs = append(chefs, chef)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return chefs, nil
}

func (r *Repository) GetChefAvatarURLByDishID(ctx context.Context, dishID int64) (string, error) {
	var chefID int64
	err := r.pg.Pool.QueryRow(ctx, "SELECT chef_id FROM dishes WHERE id = $1", dishID).Scan(
		&chefID,
	)
	if err != nil {
		return "", err
	}
	var url string
	err = r.pg.Pool.QueryRow(ctx, "SELECT image_url FROM chefs WHERE id = $1", chefID).Scan(
		&url,
	)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (r *Repository) GetChefAvatarURLByChefID(ctx context.Context, chefID int64) (string, error) {
	var url string
	if err := r.pg.Pool.QueryRow(ctx, "SELECT image_url FROM chefs WHERE id = $1", chefID).Scan(
		&url,
	); err != nil {
		return "", err
	}
	return url, nil
}
func (r *Repository) GetChefExperienceYears(ctx context.Context, chefID int64) (int, error) {
	query := "SELECT experience_years FROM chefs_experience WHERE chef_id = $1 LIMIT 1"
	log.Printf("DEBUG: Выполнение запроса: %q с chefID: %d", query, chefID)

	var expYears int
	err := r.pg.Pool.QueryRow(ctx, query, chefID).Scan(&expYears)
	if err != nil {
		// Проверяем, если ошибок, связанных с отсутствием строк, используя pgx.ErrNoRows
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("DEBUG: Для chefID %d не найдено ни одной записи", chefID)
			return 0, nil
		} else {
			log.Printf("DEBUG: Ошибка выполнения запроса: %v", err)
		}
		return 0, err
	}

	log.Printf("DEBUG: Найден опыт (experience_years): %d для chefID: %d", expYears, chefID)
	fmt.Println(expYears)
	return expYears, nil
}

// GetChefCertifications возвращает список сертификатов,
// привязанных к шеф-повару с данным chefID.
func (r *Repository) GetChefCertifications(ctx context.Context, chefID int64) ([]entity.Certification, error) {
	ctx = context.WithoutCancel(ctx)
	// Сам SQL-запрос (без лишних переводов строки в логах)
	query := `
	SELECT c.id, c.name, c.description, c.created_at
	  FROM certifications c
	  JOIN chef_certifications cc
		ON c.id = cc.certification_id
	 WHERE cc.chef_id = $1
	 ORDER BY c.name;
`
	log.Printf("DEBUG: Выполнение запроса: %q с chefID: %d", query, chefID)

	// Выполняем запрос через pgx-пул
	rows, err := r.pg.Pool.Query(ctx, query, chefID)
	if err != nil {
		// Обработка ошибок, в том числе отмены контекста
		if errors.Is(err, context.Canceled) {
			log.Printf("DEBUG: Запрос отменён контекстом: %v", err)
		} else {
			log.Printf("DEBUG: Ошибка выполнения запроса: %v", err)
		}
		return nil, err
	}
	defer rows.Close()

	var certs []entity.Certification
	for rows.Next() {
		var cert entity.Certification
		if err := rows.Scan(
			&cert.ID,
			&cert.Name,
			&cert.Description,
			&cert.CreatedAt,
		); err != nil {
			log.Printf("DEBUG: Ошибка сканирования строки: %v", err)
			return nil, err
		}
		certs = append(certs, cert)
	}

	// Проверяем, не было ли ошибок при переборе строк
	if err := rows.Err(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("DEBUG: Для chefID %d не найдено ни одного сертификата", chefID)
			return []entity.Certification{}, nil
		}
		log.Printf("DEBUG: Ошибка в итерации строк: %v", err)
		return nil, err
	}

	// Логируем результат
	if len(certs) == 0 {
		log.Printf("DEBUG: Для chefID %d не найдено ни одного сертификата", chefID)
		return []entity.Certification{}, nil
	} else {
		log.Printf("DEBUG: Найдено %d сертификатов для chefID: %d", len(certs), chefID)
	}

	return certs, nil
}
func (r *Repository) GetDishesByChefID(ctx context.Context, chefID int64) ([]dishEntity.Dish, error) {
	// Подготовим запрос без лишних переводов строк в логах
	query := `
		SELECT id, chef_id, name, description, image_url
		  FROM dishes
		 WHERE chef_id = $1
`
	log.Printf("DEBUG: Выполнение запроса: %q с chefID: %d", query, chefID)

	// Выполняем запрос через pgx-пул
	rows, err := r.pg.Pool.Query(ctx, query, chefID)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Printf("DEBUG: Запрос отменён контекстом: %v", err)
		} else {
			log.Printf("DEBUG: Ошибка выполнения запроса: %v", err)
		}
		return nil, err
	}
	defer rows.Close()

	var dishes []dishEntity.Dish
	for rows.Next() {
		var d dishEntity.Dish
		if err := rows.Scan(
			&d.ID,
			&d.ChefID,
			&d.Name,
			&d.Description,
			&d.ImageURL,
		); err != nil {
			log.Printf("DEBUG: Ошибка сканирования строки: %v", err)
			return nil, err
		}
		dishes = append(dishes, d)
	}

	// Проверяем ошибку итерации
	if err := rows.Err(); err != nil {
		log.Printf("DEBUG: Ошибка в итерации строк: %v", err)
		return nil, err
	}

	// Логируем результат
	if len(dishes) == 0 {
		log.Printf("DEBUG: Для chefID %d не найдено ни одного блюда", chefID)
	} else {
		log.Printf("DEBUG: Найдено %d блюд для chefID: %d", len(dishes), chefID)
	}

	return dishes, nil
}
