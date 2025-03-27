package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/open-policy-agent/opa/rego"
)

// Core OPA policies
var (
	//go:embed rego/authentication.rego
	opaAuthentication string
)

func main() {
	err := genToken()
	if err != nil {
		log.Fatal(err)
	}
}

func genToken() error {
	// generate a new private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generating key: %w", err)
	}

	// Generating a token requires defining a set of claims. In this applications
	// case, we only care about defining the subject and the user in question and
	// the roles they have on the database. This token will expire in a year.
	//
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)
	claims := struct {
		RegisteredClaims jwt.RegisteredClaims
		Roles            []string
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "12345678",
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"ADMIN"},
	}

	method := jwt.GetSigningMethod(jwt.SigningMethodRS256.Name)

	token := jwt.NewWithClaims(method, claims.RegisteredClaims)
	token.Header["kid"] = "346cb364-f345-4559-ae85-2739e1def422"

	str, err := token.SignedString(privateKey)
	if err != nil {
		return fmt.Errorf("signing token: %w", err)
	}

	fmt.Println("token:", str)

	// ----------------------------------------------------------------------------------

	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}))

	keyFunc := func(t *jwt.Token) (any, error) {
		return &privateKey.PublicKey, nil
	}

	var claims2 struct {
		jwt.RegisteredClaims
		Roles []string
	}

	tkn, err := parser.ParseWithClaims(str, &claims2, keyFunc)
	if err != nil {
		return fmt.Errorf("parsing token: %w", err)
	}

	if !tkn.Valid {
		return errors.New("token invalid")
	}

	fmt.Println("Signature Validated")
	fmt.Printf("Parsed claims: %#v\n", claims2)

	// -----------------------------------------------------------------------------------

	var claims3 struct {
		jwt.RegisteredClaims
		Roles []string
	}

	_, _, err = parser.ParseUnverified(str, &claims3)
	if err != nil {
		return fmt.Errorf("parsing token: %w", err)
	}

	// marshal the public key from private key to PKIX
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshaling public key: %w", err)
	}

	// construct a PEM block for the public key
	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	// write the public key to buffer
	var b bytes.Buffer
	if err := pem.Encode(&b, publicBlock); err != nil {
		return fmt.Errorf("encoding to public file: %w", err)
	}

	input := map[string]any{
		"Key":   b.String(),
		"Token": str,
	}

	if err := opaPolicyEvaluation(context.Background(), opaAuthentication, input); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	fmt.Println("VALIDATED BY REGO")
	fmt.Printf("Parsed claims: %#v\n", claims3)

	return nil
}

func opaPolicyEvaluation(ctx context.Context, opapolicy string, input any) error {
	const opaPackage string = "dibek.rego"
	const rule string = "auth"

	query := fmt.Sprintf("x = data.%s.%s", opaPackage, rule)

	q, err := rego.New(rego.Query(query), rego.Module("policy.rego", opapolicy)).PrepareForEval(ctx)
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

func genKey() (*rsa.PrivateKey, error) {
	// generate a new private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generating key: %w", err)
	}

	// crate a file for private key in PEM format
	privateFile, err := os.Create("private.pem")
	if err != nil {
		return nil, fmt.Errorf("creating private file: %w", err)
	}
	defer privateFile.Close()

	// construct a PEM block for the private key
	privateBlock := &pem.Block{
		Type:  "PRIVATE_KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(privateFile, privateBlock); err != nil {
		return nil, fmt.Errorf("encoding to private file: %w", err)
	}

	// =============================================================

	// create a file for public key in PEM format
	publicFile, err := os.Create("public.pem")
	if err != nil {
		return nil, fmt.Errorf("creating public file: %w", err)
	}
	defer publicFile.Close()

	// marshal the public key from private key to PKIX
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("marshaling public key: %w", err)
	}

	// construct a PEM block for the public key
	publicBlock := &pem.Block{
		Type:  "PUBLIC_KEY",
		Bytes: asn1Bytes,
	}

	// write the public key to the public key file
	if err := pem.Encode(publicFile, publicBlock); err != nil {
		return nil, fmt.Errorf("encoding to public file: %w", err)
	}

	fmt.Println("priavte and public key generated")

	return privateKey, nil
}
