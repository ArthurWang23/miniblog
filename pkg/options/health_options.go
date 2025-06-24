package options

import (
	"net/http"
	"net/http/pprof"

	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
)

var _ IOptions = (*HealthOptions)(nil)

// defines options for redis cluster.
type HealthOptions struct {
	// Enable debugging by exposing profiling information.
	HTTPProfile        bool   `json:"enable-http-profile" mapstructure:"enable-http-profile"`
	HealthCheckPath    string `json:"check-path" mapstructure:"check-path"`
	HealthCheckAddress string `json:"check-address" mapstructure:"check-address"`
}

func NewHealthOptions() *HealthOptions {
	return &HealthOptions{
		HTTPProfile:        false,
		HealthCheckPath:    "/healthz",
		HealthCheckAddress: "0.0.0.0:20250",
	}
}

func (o *HealthOptions) Validate() []error {
	errs := []error{}
	return errs
}

func (o *HealthOptions) AddFlags(fs *pflag.FlagSet, prefixes ...string) {
	fs.BoolVar(&o.HTTPProfile, "health.enable-http-profiler", o.HTTPProfile, "Expose runtime profiling data via HTTP.")
	fs.StringVar(&o.HealthCheckPath, "health.check-path", o.HealthCheckPath, "Specifies liveness health check request path.")
	fs.StringVar(&o.HealthCheckAddress, "health.check-address", o.HealthCheckAddress, "Specifies liveness health check bind address.")
}

func (o *HealthOptions) ServeHealthCheck() {
	r := mux.NewRouter()
	r.HandleFunc(o.HealthCheckPath, handler).Methods(http.MethodGet)
	if o.HTTPProfile {
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/{_:.*}", pprof.Index)
	}
	log.Infow("Starting health check server", "path", o.HealthCheckPath, "addr", o.HealthCheckAddress)
	if err := http.ListenAndServe(o.HealthCheckAddress, r); err != nil {
		log.Fatalw("Error serving health check endpoint:%v", err)
	}
}

func handler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`{"status":"ok"}`))
}
