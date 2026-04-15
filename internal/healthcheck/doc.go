// Package healthcheck exposes a lightweight HTTP /healthz endpoint that
// reports the daemon's readiness and live metrics (scans, alerts, open
// ports, uptime). External tools such as Docker HEALTHCHECK, Kubernetes
// liveness probes, or simple shell scripts can poll this endpoint to
// determine whether portwatch is operating normally.
//
// Usage:
//
//	server := healthcheck.New(":9090", metricsProvider)
//	server.SetReady(true)
//	go server.ListenAndServe()
//	defer server.Shutdown()
package healthcheck
