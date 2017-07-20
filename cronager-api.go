package main

import (
	"bytes"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	DB_USER := os.Getenv("CRONDBUSER")
	DB_PASS := os.Getenv("CRONDBPASS")
	DB_HOST := os.Getenv("DBHOST")

	dsn := DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ":3306)/cronager?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()
	// make sure our connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}
	type Cronjob struct {
		Id          int
		Name        string
		Command     string
		Cron_def    string
		Description string
		Active      bool
		Logtime     time.Time
	}

	router := gin.Default()
	// Add API handlers here

	// GET a cronjob
	router.GET("/job/:id", func(c *gin.Context) {
		var (
			cronjob Cronjob
			result  gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("select id, name, crondef, command, description, active, logtime from jobs where id = ?;", id)
		err = row.Scan(&cronjob.Id, &cronjob.Name, &cronjob.Cron_def, &cronjob.Command, &cronjob.Description, &cronjob.Active, &cronjob.Logtime)
		if err != nil {
			// if no results, send null
			result = gin.H{
				"result": nil,
				"count":  0,
			}
		} else {
			result = gin.H{
				"result": cronjob,
				"count":  1,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	// GET all cronjobs
	router.GET("/jobs", func(c *gin.Context) {
		var (
			cronjob  Cronjob
			cronjobs []Cronjob
		)

		rows, err := db.Query("SELECT id, name, crondef, command, description, active, logtime from jobs;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&cronjob.Id, &cronjob.Name, &cronjob.Cron_def, &cronjob.Command, &cronjob.Description, &cronjob.Active, &cronjob.Logtime)
			cronjobs = append(cronjobs, cronjob)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"result": cronjobs,
			"count":  len(cronjobs),
		})
	})

	// POST new cronjob
	router.POST("/job", func(c *gin.Context) {
		var buffer bytes.Buffer
		name := c.PostForm("name")
		crondef := c.PostForm("crondef")
		command := c.PostForm("command")
		description := c.PostForm("description")
		active := c.PostForm("active")

		stmt, err := db.Prepare("insert into jobs (name, crondef, command, description, active) values(?,?,?,?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}

		_, err = stmt.Exec(name, crondef, command, description, active)

		if err != nil {
			fmt.Print(err.Error())
		}

		// Append strings
		buffer.WriteString(name)
		buffer.WriteString(" ")
		buffer.WriteString(crondef)
		buffer.WriteString(" ")
		buffer.WriteString(command)
		buffer.WriteString(" ")
		buffer.WriteString(description)
		buffer.WriteString(" ")
		buffer.WriteString(active)
		defer stmt.Close()
		cronjob_string := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s successfully created", cronjob_string),
		})
	})

	router.PUT("/job/:id", func(c *gin.Context) {
		var buffer bytes.Buffer
		id := c.Param("id")
		name := c.PostForm("name")
		crondef := c.PostForm("crondef")
		command := c.PostForm("command")
		description := c.PostForm("description")
		active := c.PostForm("active")
		stmt, err := db.Prepare("update jobs set name = ?, crondef = ?, command = ?, description = ?, active = ? where id = ?;")

		if err != nil {
			fmt.Print(err.Error())
		}

		_, err = stmt.Exec(name, crondef, command, description, active, id)
		if err != nil {
			fmt.Print(err.Error())
		}

		// Append strings
		buffer.WriteString(name)
		buffer.WriteString(" ")
		buffer.WriteString(crondef)
		buffer.WriteString(" ")
		buffer.WriteString(command)
		buffer.WriteString(" ")
		buffer.WriteString(description)
		buffer.WriteString(" ")
		buffer.WriteString(active)
		defer stmt.Close()
		cronjob_string := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s successfully updated", cronjob_string),
		})
	})

	router.DELETE("/job/:id", func(c *gin.Context) {
		id := c.Param("id")
		stmt, err := db.Prepare("delete from jobs where id=?;")
		if err != nil {
			fmt.Print(err.Error())
		}

		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Print(err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully deleted job with id: %s", id),
		})

	})

	router.Run(":3000")
}
