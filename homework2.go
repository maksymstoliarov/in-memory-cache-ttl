package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {
	rand.Seed(time.Now().Unix())

	const workersCount, usersCount = 100, 100
	wg := &sync.WaitGroup{}

	users := make(chan User, usersCount)

	startTime := time.Now()

	for i := 0; i < workersCount; i++ {
		go worker(users, wg)
	}

	generateUsers(usersCount, users, wg)

	// close(users)
	wg.Wait()

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func saveUserInfo(user User, wg *sync.WaitGroup) {
	fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

	filename := fmt.Sprintf("users/uid%d.txt", user.id)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(user.getActivityInfo())
	time.Sleep(time.Second)
	wg.Done()
}

func generateUsers(count int, users chan User, wg *sync.WaitGroup) {
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			users <- User{
				id:    i + 1,
				email: fmt.Sprintf("user%d@company.com", i+1),
				logs:  generateLogs(rand.Intn(1000)),
			}
			fmt.Printf("generated user %d\n", i+1)
			time.Sleep(time.Millisecond * 100)
			wg.Done()
		}(i)
	}
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}

func worker(users chan User, wg *sync.WaitGroup) {
	for u := range users {
		wg.Add(1)
		saveUserInfo(u, wg)
	}
}
