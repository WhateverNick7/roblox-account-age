package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type AccountAge struct {
	Usernames []string
	Answer    []map[string]string
}

type User struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type UserBirthdate struct {
	Created string `json:"created"`
}

func NewAccountAge(usernames []string) *AccountAge {
	uidList := getUIDs(usernames)
	answer := make([]map[string]string, 0)

	for _, id := range uidList {
		birthdate := getUserBirthdate(id.ID)
		answer = append(answer, map[string]string{
			"name": id.Name,
			"age":  format(birthdate.Created),
		})
	}

	return &AccountAge{Usernames: usernames, Answer: answer}
}

func format(isoDatetime string) string {
	// Parse the ISO datetime
	t, _ := time.Parse(time.RFC3339, isoDatetime)
	formattedDate := t.Format("January 02, 2006")

	now := time.Now().UTC()
	diff := now.Sub(t)

	totalDays := int(diff.Hours() / 24)

	years := totalDays / 365
	months := (totalDays % 365) / 30
	weeks := (totalDays % 365 % 30) / 7
	days := totalDays % 7

	parts := []string{}
	if years > 0 {
		parts = append(parts, fmt.Sprintf("%d year%s", years, pluralize(years)))
	}
	if months > 0 {
		parts = append(parts, fmt.Sprintf("%d month%s", months, pluralize(months)))
	}
	if weeks > 0 {
		parts = append(parts, fmt.Sprintf("%d week%s", weeks, pluralize(weeks)))
	}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d day%s", days, pluralize(days)))
	}

	ageDesc := strings.Join(parts, ", ")
	if ageDesc == "" {
		ageDesc = "0 days"
	}

	return fmt.Sprintf("%s : %s old", formattedDate, ageDesc)
}

func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

func getUserBirthdate(uid int) UserBirthdate {
	headers := map[string]string{
		"accept":        "application/json",
		"Content-Type":  "application/json",
		"Origin":        "https://www.roblox.com",
		"Sec-GPC":       "1",
		"Connection":    "keep-alive",
		"Referer":       "https://www.roblox.com/",
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("https://users.roblox.com/v1/users/%d", uid), nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var userBirthdate UserBirthdate
	json.Unmarshal(body, &userBirthdate)

	return userBirthdate
}

func getUIDs(usernames []string) []User {
	headers := map[string]string{
		"accept":        "application/json",
		"Content-Type":  "application/json",
		"Origin":        "https://www.roblox.com",
		"Sec-GPC":       "1",
		"Connection":    "keep-alive",
		"Referer":       "https://www.roblox.com/",
	}

	jsonData, _ := json.Marshal(map[string]interface{}{
		"usernames":           usernames,
		"excludeBannedUsers": true,
	})

	req, _ := http.NewRequest("POST", "https://users.roblox.com/v1/usernames/users", bytes.NewBuffer(jsonData))
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var response struct {
		Data []User `json:"data"`
	}
	json.Unmarshal(body, &response)

	return response.Data
}

func getUsers(file string) []string {
	answer := []string{}
	data, _ := ioutil.ReadFile(file)
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if line != "" {
			if strings.Contains(line, ":") {
				answer = append(answer, strings.Split(line, ":")[0])
			} else {
				answer = append(answer, line)
			}
		}
	}
	return answer
}

func formatEntries(entries []map[string]string, minSpaces int) string {
	var maxNameLen, maxAgeLen int
	for _, entry := range entries {
		if len(entry["name"]) > maxNameLen {
			maxNameLen = len(entry["name"])
		}
		if len(entry["age"]) > maxAgeLen {
			maxAgeLen = len(entry["age"])
		}
	}

	var lines []string
	for _, entry := range entries {
		nameSpaces := maxNameLen - len(entry["name"]) + minSpaces
		ageSpaces := maxAgeLen - len(entry["age"])

		line := fmt.Sprintf("%s%s%s%s", entry["name"], strings.Repeat(" ", nameSpaces), entry["age"], strings.Repeat(" ", ageSpaces))
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func main() {
	fmt.Println("Enter path to accounts. (user:pass format or just user)")
	var path string
	fmt.Scanln(&path)

	answer := NewAccountAge(getUsers(path)).Answer
	fmt.Println(formatEntries(answer, 3))
	fmt.Println("Press Enter to exit...")
    fmt.Scanln()
}
