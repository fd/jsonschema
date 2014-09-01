package jsonschema

import (
  "net/url"
)

type uriFormat struct{}

func (*uriFormat) IsValid(x interface{}) bool {
  s, ok := x.(string)
  if !ok {
    return true
  }

  u, err := url.Parse(s)
  if err != nil {
    return false
  }

  if u.Scheme == "" {
    return false
  }

  if u.Host == "" {
    return false
  }

  // url.Parse() already validated the path
  // for i, l := 0, len(u.Path); i < l; i++ {
  //   c := u.Path[i]
  //   if 'a' <= c && c <= 'z' {
  //     continue
  //   }
  //   if 'A' <= c && c <= 'Z' {
  //     continue
  //   }
  //   if '0' <= c && c <= '9' {
  //     continue
  //   }
  //   if '-' == c || '.' == c || '_' == c || '~' == c || '!' == c ||
  //     '$' == c || '&' == c || '\'' == c || '(' == c || ')' == c ||
  //     '*' == c || '+' == c || ',' == c || ';' == c || '=' == c ||
  //     ':' == c || '@' == c || '%' == c || '/' == c {
  //     continue
  //   }
  //   return false
  // }

  return true
}
