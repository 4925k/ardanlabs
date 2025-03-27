package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/4925k/ardanlabs/software_design_with_kubernetes/service/foundation/logger"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/open-policy-agent/opa/rego"
)

var ErrForbidden = errors.New("attempted action is not allowed")

// Claims represents the authorization claims transmitted via a JWT
type Claims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

// KeyLookup declares a method set of behavior for looking up
type KeyLookup interface {
	PrivateKey(string) (string, error)
	PublicKey(string) (string, error)
}

// Config represents the configuration for the auth package.
type Config struct {
	Log       *logger.Logger
	KeyLookup KeyLookup
	Issuer    string
}

// Auth is used to authenticate clients. It can generate a token for a
// set of user claims and recreate the claims by parsing the token.
type Auth struct {
	log       *logger.Logger
	keyLookup KeyLookup
	method    jwt.SigningMethod
	parser    *jwt.Parser
	issuer    string
	cache     map[string]string

	mu sync.RWMutex
}

func New(cfg Config) (*Auth, error) {
	a := &Auth{
		log:       cfg.Log,
		keyLookup: cfg.KeyLookup,
		method:    jwt.SigningMethodRS256,
		parser:    jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name})),
		issuer:    cfg.Issuer,
		cache:     make(map[string]string),
	}

	return a, nil
}

// GenerateToken generates a signed JWT token string representing the user claims
func (a *Auth) GenerateToken(kid string, claims Claims) (string, error) {
	token := jwt.NewWithClaims(a.method, claims)
	token.Header["kid"] = kid

	privateKeyPEM, err := a.keyLookup.PrivateKey(kid)
	if err != nil {
		return "", fmt.Errorf("retrieving private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		return "", fmt.Errorf("parsing private pem: %w", err)
	}

	tkn, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return tkn, nil
}

func (a *Auth) Authenticate(ctx context.Context, bearerToken string) (Claims, error) {
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return Claims{}, errors.New("expected format: Bearer <token>")
	}

	var claims Claims
	token, _, err := a.parser.ParseUnverified(parts[1], &claims)
	if err != nil {
		return Claims{}, fmt.Errorf("error parsing token: %w", err)
	}

	// perform an extra level of authentication verification with OPA

	kidRaw, exists := token.Header["kid"]
	if !exists {
		return Claims{}, errors.New("kid missing from header")
	}

	kid, ok := kidRaw.(string)
	if !ok {
		return Claims{}, errors.New("kid malformed")
	}

	pem, err := a.keyLookup.PublicKey(kid)
	if err != nil {
		return Claims{}, fmt.Errorf("failed to fetch public key: %w", err)
	}

	input := map[string]any{
		"Key":   pem,
		"Token": parts[1],
		"ISS":   a.issuer,
	}

	if err := a.opaPolicyEvaluation(ctx, opaAuthentication, RuleAuthenticate, input); err != nil {
		return Claims{}, fmt.Errorf("authentication failed: %w", err)
	}

	// check database for user status
	if err := a.isUserEnabled(ctx, claims); err != nil {
		return Claims{}, fmt.Errorf("user not enabled: %w", err)
	}

	return claims, nil
}

// Authorize attempts to authorize the user with the provided input roles
func (a *Auth) Authorize(ctx context.Context, claims Claims, userID uuid.UUID, rule string) error {
	input := map[string]any{
		"Roles":   claims.Roles,
		"Subject": claims.Subject,
		"UserID":  userID,
	}

	if err := a.opaPolicyEvaluation(ctx, opaAuthorization, rule, input); err != nil {
		return fmt.Errorf("rego evaluation failed: %w", err)
	}

	return nil
}

func (a *Auth) publicKeyLookup(kid string) (string, error) {
	// cache check up
	pem, err := func() (string, error) {
		a.mu.RLock()
		defer a.mu.RUnlock()

		pem, exists := a.cache[kid]
		if !exists {
			return "", errors.New("not found")
		}

		return pem, nil
	}()
	if err == nil {
		return pem, nil
	}

	pem, err = a.keyLookup.PublicKey(kid)
	if err != nil {
		return "", fmt.Errorf("fetching public key: %w", err)
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	a.cache[kid] = pem

	return pem, nil
}

func (a *Auth) opaPolicyEvaluation(ctx context.Context, opaPolicy string, rule string, input any) error {
	query := fmt.Sprintf("x = data.%s.%s", opaPackage, rule)

	q, err := rego.New(rego.Query(query), rego.Module("policy.rego", opaPolicy)).PrepareForEval(ctx)
	if err != nil {
		return err
	}

	results, err := q.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	if len(results) == 0 {
		return errors.New("no results")
	}

	result, ok := results[0].Bindings["x"].(bool)
	if !ok || !result {
		return errors.New("no results")
	}

	return nil
}

func (a *Auth) isUserEnabled(_ context.Context, _ Claims) error {
	return nil
}
