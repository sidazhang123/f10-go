package std2micro

import (
	"context"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	status_code "github.com/sidazhang123/f10-go/plugins/breaker/http"
	"math/rand"
	"net/http"
	"time"
)

var sampleFreq = 100

func init() {
	rand.Seed(time.Now().Unix())
}

func SetSamplingFrequency(n int) {
	sampleFreq = n
}

func TracerWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		md := make(map[string]string)
		spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		sp := opentracing.GlobalTracer().StartSpan(r.URL.Path, opentracing.ChildOf(spanCtx))
		defer sp.Finish()
		if err := opentracing.GlobalTracer().Inject(
			sp.Context(),
			opentracing.TextMap,
			opentracing.TextMapCarrier(md)); err != nil {
			log.Logf(err.Error())
		}
		ctx := context.TODO()
		ctx = opentracing.ContextWithSpan(ctx, sp)
		ctx = metadata.NewContext(ctx, md)
		r = r.WithContext(ctx)
		sct := &status_code.StatusCodeTracker{w, http.StatusOK}
		h.ServeHTTP(sct.WrappedResponseWriter(), r)

		ext.HTTPMethod.Set(sp, r.Method)
		ext.HTTPUrl.Set(sp, r.URL.EscapedPath())
		ext.HTTPStatusCode.Set(sp, uint16(sct.Status))
		if sct.Status >= http.StatusInternalServerError {
			ext.Error.Set(sp, true)
		} else if sampleFreq < rand.Intn(100) {
			ext.SamplingPriority.Set(sp, 0)
		}

	})
}
