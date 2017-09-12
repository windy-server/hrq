package hrq

import "net/url"

func mapStringList(ms map[string]string) map[string][]string {
	result := map[string][]string{}
	for k, v := range ms {
		result[k] = []string{v}
	}
	return result
}

// MakeURL makes url
func MakeURL(baseURL string, params map[string]string) string {
	msl := mapStringList(params)
	v := url.Values(msl)
	u := baseURL + "?" + v.Encode()
	return u
}
