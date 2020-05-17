package settings

import "os"

func GetEnv() string {
	env := os.Getenv("ENV")
	if env == "" {
		env = "DEMO"
	}
	return env
}

func GetDatabaseHost() string {

	host := os.Getenv("DATABASE_HOST")
	if host == "" {
		host = "localhost"
	}
	return host
}

func GetDatabaseName() string {

	name := os.Getenv("DATABASE_NAME")
	if name == "" {
		name = "rafiki"
	}
	return name
}

func GetDatabaseUser() string {

	user := os.Getenv("DATABASE_USER")
	if user == "" {
		user = "postgres"
	}
	return user
}

func GetDatabasePassword() string {

	password := os.Getenv("DATABASE_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	return password
}

func GetAfricasTalkingUsername() string {
	env := os.Getenv("AT_USERNAME")
	if env == "" {
		env = "rafiki_app"
	}
	return env
}

func GetAfricasTalkingKey() string {
	env := os.Getenv("AT_KEY")
	if env == "" {
		env = "76a7d052969a47ff266a6a733c7db65fb86889b724f4173bf9f188f0239a23a2"
	}
	return env
}
