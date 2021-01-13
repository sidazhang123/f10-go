package breaker

import (
	"errors"
	"fmt"

	"github.com/afex/hystrix-go/hystrix"
	state_code "github.com/sidazhang123/f10-go/plugins/breaker/http"
	"net/http"
)

func BreakerWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.Method + "-" + r.RequestURI
		hystrix.Do(name, func() error {
			sct := &state_code.StatusCodeTracker{w, http.StatusOK}
			h.ServeHTTP(sct.WrappedResponseWriter(), r)

			if sct.Status >= http.StatusInternalServerError {
				str := fmt.Sprintf("status code %d", sct.Status)
				return errors.New(str)
			}
			return nil
		}, func(e error) error {
			if e == hystrix.ErrCircuitOpen {
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("[Breaker] Please try again later..."))
			}
			return e
		})
	})
}
