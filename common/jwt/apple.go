package jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AppleJWT struct {
	ClientID string
	Secret   string
	KeyID    string
	TeamID   string
	Redirect string
}

//https://www.cnblogs.com/biwentao/p/12179321.html

//https://github.com/tptpp/sign-in-with-apple/blob/master/README.md

func NewAppleJWT(clientID, secret, keyID, teamID, redirect string) *AppleJWT {
	return &AppleJWT{
		ClientID: clientID,
		Secret:   secret,
		KeyID:    keyID,
		TeamID:   teamID,
		Redirect: redirect,
	}
}

// create client_secret
func (a *AppleJWT) GetAppleSecret() string {
	token := &jwt.Token{
		Header: map[string]interface{}{
			"alg": "ES256",
			"kid": a.KeyID,
		},
		Claims: jwt.MapClaims{
			"iss": a.TeamID,
			"iat": time.Now().Unix(),
			// constraint: exp - iat <= 180 days
			"exp": time.Now().Add(24 * time.Hour).Unix(),
			"aud": "https://appleid.apple.com",
			"sub": a.ClientID,
		},
		Method: jwt.SigningMethodES256,
	}

	ecdsaKey, _ := a.AuthKeyFromBytes([]byte(a.Secret))
	ss, _ := token.SignedString(ecdsaKey)
	return ss
}

// create private key for jwt sign
func (a *AppleJWT) AuthKeyFromBytes(key []byte) (*ecdsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errors.New("token: AuthKey must be a valid .p8 PEM file")
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
		return nil, err
	}

	var pkey *ecdsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*ecdsa.PrivateKey); !ok {
		return nil, errors.New("token: AuthKey must be of type ecdsa.PrivateKey")
	}

	return pkey, nil
}

// do http request
func (a *AppleJWT) HttpRequest(method, addr string, params map[string]string) ([]byte, int, error) {
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	var request *http.Request
	var err error
	if request, err = http.NewRequest(method, addr, strings.NewReader(form.Encode())); err != nil {
		return nil, 0, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var response *http.Response
	if response, err = http.DefaultClient.Do(request); nil != err {
		return nil, 0, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, 0, err
	}
	return data, response.StatusCode, nil
}

func (a *AppleJWT) RequestToken(code string) ([]byte, error) {
	// replace your code here

	data, status, err := a.HttpRequest("POST", "https://appleid.apple.com/auth/token", map[string]string{
		"client_id":     a.ClientID,
		"client_secret": a.GetAppleSecret(),
		"code":          code,
		"grant_type":    "authorization_code",
		"redirect_uri":  a.Redirect,
	})

	fmt.Printf("%d\n%v\n%s", status, err, data)
	return data, err
}
