package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

var ctx = context.Background()

func main() {
	fmt.Println(sql.ErrNoRows)
	dbpool := connToDB()
	defer dbpool.Close()

	album := Album{Title: "American", Artist: "Joshua", Price: 15.42}
	albums, err := addAlbum(dbpool, album)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("albums %v\n", albums)
}

func connToDB() *pgxpool.Pool {
	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to create connection pool: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("DB connected successfully")

	return dbpool
}

func albumsByArtist(dbpool *pgxpool.Pool, name string) ([]Album, error) {
	var albums []Album

	rows, err := dbpool.Query(ctx, "SELECT id, price, artist FROM album WHERE artist = $1", name)

	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	defer rows.Close()

	for rows.Next() {
		var alb Album

		if err = rows.Scan(&alb.ID, &alb.Price, &alb.Artist); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}

		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	return albums, nil
}

func albumsByID(dbpool *pgxpool.Pool, id int64) (Album, error) {
	var album Album

	row := dbpool.QueryRow(ctx, "SELECT * FROM album WHERE id = $1", id)

	if err := row.Scan(&album.ID, &album.Artist, &album.Title, &album.Price); err != nil {
		// You could return a status 404 code here
		if err == sql.ErrNoRows {
			return album, fmt.Errorf("albumsById %d: no such album", id)
		}

		return album, fmt.Errorf("albumByID %d: %v", id, err)
	}

	return album, nil
}

func addAlbum(dbpool *pgxpool.Pool, alb Album) (int64, error) {
	result, err := dbpool.Exec(ctx, "INSERT INTO album (artist, title, price) VALUES ($1, $2, $3) RETURNING id", alb.Artist, alb.Title, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: err", err)
	}

	fmt.Print("result: %v", result)
	
	return 0, nil
}
