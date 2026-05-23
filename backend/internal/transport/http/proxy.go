package httptransport

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func (a *api) proxyToMessaging(w http.ResponseWriter, r *http.Request) bool {
	return a.proxyToService(w, r, a.deps.Config.MessagingServiceURL, "messaging")
}

func (a *api) proxyToFiling(w http.ResponseWriter, r *http.Request) bool {
	return a.proxyToService(w, r, a.deps.Config.FilingServiceURL, "filing")
}

func (a *api) proxyToFileService(w http.ResponseWriter, r *http.Request) bool {
	return a.proxyToService(w, r, a.deps.Config.FileServiceURL, "file")
}

func (a *api) proxyToService(w http.ResponseWriter, r *http.Request, rawURL, name string) bool {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return false
	}
	target, err := url.Parse(rawURL)
	if err != nil {
		writeError(w, http.StatusBadGateway, name+"_proxy_failed", name+" service url is invalid")
		return true
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Path = r.URL.Path
		req.URL.RawQuery = r.URL.RawQuery
		req.Host = target.Host
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		writeError(w, http.StatusBadGateway, name+"_unavailable", name+" service is unavailable")
	}
	proxy.ServeHTTP(w, r)
	return true
}
