package net_http

import "fmt"

func FormattingProxy(proxy string) string {
	return fmt.Sprintf("http://%s", proxy)
}
