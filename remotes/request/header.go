package request

import "fmt"

func Header() map[string]string {
	return map[string]string{}
}

func InjectAuthorization(bearer string, header map[string]string) {
	header["Authorization"] = fmt.Sprintf("Bearer %s", bearer)
}

func InjectAcceptGzip(header map[string]string) {
	header["Accept-Encoding"] = "gzip"
}

func InjectSendGzip(header map[string]string) {
	header["Content-Encoding"] = "gzip"
}
