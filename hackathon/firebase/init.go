package firebase

import (
	"context"
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func InitFirebaseApp() (*firebase.App, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS_PATH"))
	return firebase.NewApp(ctx, nil, opt)
}
