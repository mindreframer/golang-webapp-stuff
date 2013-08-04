package repository

import (
	"github.com/jingweno/octokat"
	"github.com/jingweno/progmob/database"
	"labix.org/v2/mgo/bson"
)

var (
	db             *database.Database
	userRepository *user
	repoRepository *repo
	mobRepository  *mob
)

func init() {
	userRepository = &user{}
	repoRepository = &repo{}
	mobRepository = &mob{}
}

type Query bson.M

func Setup(database *database.Database) error {
	db = database
	err := userRepository.ensureIndex()
	if err != nil {
		return err
	}

	err = repoRepository.ensureIndex()
	if err != nil {
		return err
	}

	err = mobRepository.ensureIndex()

	return err
}

func User() *user {
	return userRepository
}

func Repo() *repo {
	return repoRepository
}

func Mob() *mob {
	return mobRepository
}

func GitHub(token string) *gitHub {
	client := octokat.NewClient().WithToken(token)

	return &gitHub{client}
}
