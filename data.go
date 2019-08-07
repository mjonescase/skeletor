package main

import (
	"database/sql"
	"log"
	"skeletor/utils"
)

func saveUserProfile(profile *Profile) {
	profile.Password = utils.HashPassword(profile.Password)
	err := session.QueryRow(`INSERT INTO profile (
		firstname,
		lastname,
		username,
		email,
		title,
		password,
		mobilenumber
	) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		profile.Firstname,
		profile.Lastname,
		profile.Username,
		profile.Email,
		profile.Title,
		profile.Password,
		profile.MobileNumber).Scan(&profile.Id)
	if err != nil {
		log.Print(err)
	}
	profile.Password = ""
}

func queryUserCredential(profile *Profile) bool {
	result := false

	err := session.QueryRow(`SELECT 
		id,
		firstname, 
		lastname, 
		username, 
		email, 
		title, 
		mobilenumber FROM profile WHERE 
		username = $1 AND password = $2`, profile.Username, profile.Password).Scan(&profile.Id,
		&profile.Firstname,
		&profile.Lastname,
		&profile.Username,
		&profile.Email,
		&profile.Title,
		&profile.MobileNumber)
	switch {
	case err == sql.ErrNoRows:
		return result
	case err != nil:
		return result
	default:
		result = true
	}

	profile.Password = ""
	return result
}

func getAllUsers() []Profile {
	profile := Profile{}
	results := []Profile{}
	rows, err := session.Query(`SELECT id, firstname, lastname, username,
title, mobilenumber FROM profile`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		profile = Profile{}
		if err := rows.Scan(&profile.Id,
			&profile.Firstname,
			&profile.Lastname,
			&profile.Username,
			&profile.Title,
			&profile.MobileNumber,
		); err != nil {
			log.Fatal(err)
		}
		results = append(results, profile)
	}
	return results
}
