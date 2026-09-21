package repository

import (
	"context"
	"fmt"

	"minisocial/internal/config"
	"minisocial/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, username string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(ctx, "INSERT INTO users (username) VALUES ($1) RETURNING id, username", username).
		Scan(&user.ID, &user.Username)
	return user, err
}

func (r *PostgresRepository) GetUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(ctx, "SELECT id, username FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Username)
	return user, err
}

func (r *PostgresRepository) FollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	_, err := r.db.Exec(ctx, "INSERT INTO follows (follower_id, following_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", followerID, followingID)
	return err
}

func (r *PostgresRepository) CreatePost(ctx context.Context, authorID uuid.UUID, content string) (domain.Post, error) {
	var post domain.Post
	err := r.db.QueryRow(ctx, `INSERT INTO posts (author_id, content) VALUES ($1, $2) 
		RETURNING id, author_id, content`, authorID, content).
		Scan(&post.ID, &post.AuthorID, &post.Content)
	return post, err
}

func (r *PostgresRepository) GetFeed(ctx context.Context, userID uuid.UUID) ([]domain.Post, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.author_id, u.username, p.content, COUNT(l.user_id) as likes_count
		FROM posts p
		JOIN follows f ON f.following_id = p.author_id
		JOIN users u ON u.id = p.author_id
		LEFT JOIN likes l ON l.post_id = p.id
		WHERE f.follower_id = $1
		GROUP BY p.id, u.username
		ORDER BY p.created_at DESC
		LIMIT 50`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var p domain.Post
		if err := rows.Scan(&p.ID, &p.AuthorID, &p.AuthorName, &p.Content, &p.LikesCount); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostgresRepository) LikePost(ctx context.Context, userID, postID uuid.UUID) error {
	_, err := r.db.Exec(ctx, "INSERT INTO likes (user_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, postID)
	return err
}

func GetDSN(cfg DBConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
}

type DBConfig = config.DBConfig
