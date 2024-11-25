package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"
    "math/rand"

    _ "github.com/go-sql-driver/mysql"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Generate a random string of length n
func randSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}


func main() {
    
    // DB conn
    db, err := sql.Open(
	"mysql", "root:5718@tcp(localhost:3306)/main")
    
    // Check errors
    if err != nil  {
	log.Fatal(err)
    } else {
	fmt.Println("Connected successfully.")
	fmt.Println()
    }
    
    // Try to ping DB
    if err := db.Ping(); err != nil {
	log.Fatal(err)
    }

    /*
     * Create the table
     */
    { 
	q := `
	CREATE TABLE IF NOT EXISTS users (
	    id INT AUTO_INCREMENT,
	    username TEXT NOT NULL,
	    password TEXT NOT NULL,
	    created_at DATETIME,
	    PRIMARY KEY (id)
	);`

	if _, err := db.Exec(q); err != nil {
	    log.Fatal(err) 
	}
    }

    /*
     * Insert into table
     */
    {
	uname 	:= randSeq(8) 
	pwd   	:= randSeq(8)
	time  	:= time.Now()

	q 	:= `
	INSERT INTO users 
	    (username, password, created_at) VALUES (?, ?, ?)`

	if res, err := db.Exec(q, uname, pwd, time); err != nil {
	    log.Fatal(err)
	} else if id, err := res.LastInsertId(); err != nil {
	    log.Fatal(err)
	} else {
	    fmt.Println("Created user with id ", id)
	    fmt.Println()
	}
     }

    /*
     * Query DB
     */
    {
	  var (
	      id	int
	      uname	string
	      pwd	string
	      createdAt	time.Time
	      timeUint	[]uint8	//time.Time
	  )
	
	  q	:= `
	  SELECT id, username, password, created_at
	  FROM users WHERE id = ?
	  `

	  if err := db.QueryRow(q, 1).Scan(&id, &uname, &pwd, &timeUint); err != nil {
	      log.Fatal(err)
	  } else {
	    createdAt, err = time.Parse("2006-01-02 15:04:05", string(timeUint));
	    if err != nil {
		log.Fatal(err)
	    }
	}

	fmt.Println("Getting first user...")
	fmt.Printf("%3d %s %s %v\n\n", id, uname, pwd, createdAt)
    }
    
    /*
     * Query db for all users and create user struct
     */
    {
	type user struct {
	    id		int
	    uname	string
	    pwd		string
	    createdAt	time.Time
	}

	var timeUint []uint8
	
	q 	:= `
	SELECT id, username, password, created_at FROM users
	`
	rows, err := db.Query(q)
	defer rows.Close()

	if err != nil {
	    log.Fatal(err)
	}

	var users []user
	    
	fmt.Println("Getting all users...")

	for rows.Next() {
	    var u user
	    if err := rows.Scan(&u.id, &u.uname, &u.pwd, &timeUint); err != nil {
		log.Fatal(err)
	    }
	    
	    u.createdAt, err = time.Parse("2006-01-02 15:04:05", string(timeUint))
	    if err != nil {
		log.Fatal(err)
	    }

	    users = append(users, u)
	}

	for _, u := range users {
	    fmt.Printf("%3d %s %s %v\n", u.id, u.uname, u.pwd, u.createdAt)
	}
    }

    /*
     * To delete we just run a SQL-query, for example:
     * _, err := db.Exec(`DELETE FROM users WHERE id = ?`, 1) 
     */
}
