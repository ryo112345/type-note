package account

import "strings"

type Email string

func ParseEmail(raw string) (Email, error) {
	e := Email(strings.TrimSpace(raw))
	if !strings.Contains(string(e), "@") {
		return "", ErrInvalidEmail
	}
	return e, nil
}

func (e Email) String() string {
	return string(e)
}

func (e Email) Domain() string {
	parts := strings.Split(string(e), "@")
	return parts[1]
}

func (e Email) Local() string {
	parts := strings.Split(string(e), "@")
	return parts[0]
}

func (e Email) Mask() string {
	local := e.Local()
	domain := e.Domain()

	if len(local) <= 3 {
		return local + "******@" + domain
	}
	return local[:3] + "******@" + domain
}
