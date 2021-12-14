package main

import (
    "os"
    "log"
    "fmt"
    "context"
    "net/url"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/auth0/go-jwt-middleware/v2"
    "github.com/auth0/go-jwt-middleware/v2/validator"
)

type CustomClaims struct {
    Scope string `json:"scope"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
    return nil
}

func main() {
    handleRouts()
}

func handleRouts() {
    host := os.Getenv("HOST")
    port := os.Getenv("PORT")

    middleware := getMiddleware()

    router := mux.NewRouter()
    router.Handle("/api/public", getPublicMessage())
    router.Handle("/api/private", middleware.CheckJWT(getPrivateMessage()))

    log.Printf("server started and listening on http://%s:%s", host, port)

    http.ListenAndServe(host+":"+port, router)
}

func getMiddleware() *jwtmiddleware.JWTMiddleware {

    issuerURL, err := url.Parse("https://dev-k18jl6aj.us.auth0.com/")

    if err != nil {
        log.Fatalf("failed to parse the issuer url: %v", err)
    }

    keyFunc := func(ctx context.Context) (interface{}, error) {
        return []byte("XaJzKyFuoImgOwNUMQMiJDZsvPXv7hu2"), nil
    }

    jwtValidator, err := validator.New(
        keyFunc,
        validator.HS256,
        issuerURL.String(),
        []string{"https://dune/"},
        validator.WithCustomClaims(&CustomClaims{}),
    )

    if err != nil {
        log.Fatalf("failed to set up the validator: %v", err)
    }

    errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
        log.Printf("Encountered error while validating JWT: %v", err)
    }

    middleware := jwtmiddleware.New(
        jwtValidator.ValidateToken,
        jwtmiddleware.WithErrorHandler(errorHandler),
    )

    return middleware
}

func getPublicMessage() http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        body := fmt.Sprintf("Public message")
        w.Write([]byte(body))
    }
    return http.HandlerFunc(fn)
}

func getPrivateMessage() http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        body := fmt.Sprintf("Private fucking message")
        w.Write([]byte(body))
    }
    return http.HandlerFunc(fn)
}
