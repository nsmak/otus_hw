package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	if domain == "" {
		return result, nil
	}

	buf := bufio.NewReader(r)
	for {
		l, _, err := buf.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("can't read line: %w", err)
		}

		email := jsoniter.Get(l, "Email").ToString()
		if strings.HasSuffix(email, domain) {
			emailDomain := strings.ToLower(strings.SplitN(email, "@", 2)[1])
			result[emailDomain]++
		}
	}
	return result, nil
}
