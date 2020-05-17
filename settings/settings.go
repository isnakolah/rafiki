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
		name = "recoin-notification"
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
		env = "recoin"
	}
	return env
}

func GetAfricasTalkingKey() string {
	env := os.Getenv("AT_KEY")
	if env == "" {
		env = "4ca8b2e12a068378cb47ae617868f2d7bcea7bc066ac2ccc739caa8bd93939bd"
	}
	return env
}
