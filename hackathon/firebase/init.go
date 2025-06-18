package firebase

import (
	"context"
	"os"
	"sync"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	"firebase.google.com/go/v4/auth"
)

var (
	app     *firebase.App
	authCli *auth.Client
	once    sync.Once
	errInit error
)

func InitFirebase() {
	once.Do(func() {
		ctx := context.Background()
		opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS_PATH"))
		app, errInit = firebase.NewApp(ctx, nil, opt)
		if errInit != nil {
			return
		}
		authCli, errInit = app.Auth(ctx)
	})
}

func GetAuthClient() (*auth.Client, error) {
	if authCli == nil {
		InitFirebase()
	}
	return authCli, errInit
}
