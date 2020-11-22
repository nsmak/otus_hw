package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

var awaitingPrefix = []byte("\"Email\":\"")

const lenPrefix = 9

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	var builder strings.Builder
	builder.WriteString(".")
	builder.WriteString(domain)
	domainValue := builder.String()

	buf := bufio.NewReader(r)
	for {
		l, _, err := buf.ReadLine()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		startIndex := bytes.Index(l, awaitingPrefix) + lenPrefix
		for i, b := range l[startIndex:] {
			if string(b) == "\"" {
				email := string(l[startIndex:][:i])
				if strings.HasSuffix(email, domainValue) {
					emailDomain := strings.ToLower(strings.SplitN(email, "@", 2)[1])
					result[emailDomain]++
				}
				break
			}
		}
	}
	return result, nil
}
