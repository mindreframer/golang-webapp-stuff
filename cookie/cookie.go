package cookie

import (
	"github.com/gorilla/securecookie"
	"net/http"
)

type Cookie struct {
	Name   string
	Values CookieValues
	Path   string
}

type CookieValues map[string]string

func Write(w http.ResponseWriter, c *Cookie) {
	s := securecookie.New(HashKey, BlockKey)
	if encoded, err := s.Encode(c.Name, c.Values); err == nil {
		cookie := &http.Cookie{
			Name:  c.Name,
			Value: encoded,
			Path:  c.Path,
		}

		http.SetCookie(w, cookie)
	}
}

func Read(r *http.Request, name string) (*Cookie, error) {
	s := securecookie.New(HashKey, BlockKey)
	cookie, err := r.Cookie(name)
	if err != nil {
		return nil, err
	}

	values := CookieValues{}
	err = s.Decode(name, cookie.Value, &values)
	if err != nil {
		return nil, err
	}

	return &Cookie{Name: cookie.Name, Path: cookie.Path, Values: values}, nil
}

func Delete(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:   name,
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}
