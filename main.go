package main

import (
    "os"
    "log"
    "fmt"
    "time"
    "context"
    "net/url"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/auth0/go-jwt-middleware/v2"
    "github.com/auth0/go-jwt-middleware/v2/jwks"
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

    issuerURL, err := url.Parse("https://dev-k18jl6aj.us.auth0.com/")

    if err != nil {
        log.Fatalf("failed to parse the issuer url: %v", err)
    }

    provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

    jwtValidator, err := validator.New(
        provider.KeyFunc,
        validator.RS256,
        issuerURL.String(),
        []string{"https://ruben/"},
        // validator.WithCustomClaims(&CustomClaims{}),
        // validator.WithAllowedClockSkew(time.Minute),
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

    // Gateway API starts here
    host := os.Getenv("HOST")
    port := os.Getenv("PORT")

    router := mux.NewRouter()
    router.Handle("/api/public", getPublicMessage())
    router.Handle("/api/private", middleware.CheckJWT(getPrivateMessage()))

    log.Printf("server started and listening on http://%s:%s", host, port)

    http.ListenAndServe(host+":"+port, router)
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
        body := fmt.Sprintf("Public message")
        w.Write([]byte(body))
    }
    return http.HandlerFunc(fn)
}
