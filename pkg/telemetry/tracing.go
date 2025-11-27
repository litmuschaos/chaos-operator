package telemetry

import (
	"context"
	"encoding/json"
	"github.com/litmuschaos/chaos-operator/api/litmuschaos/v1alpha1"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

const (
	// TracerName is the name of the tracer
	TracerName = "litmuschaos.io/chaos-operator"
)

// InitSpanContext initialize tracing by creating the root span and injecting the
// spanContext is propagated through annotations in the CR
func InitSpanContext(ctx context.Context, engine *v1alpha1.ChaosEngine) context.Context {
	tracerProvider := otel.GetTracerProvider()
	pro := otel.GetTextMapPropagator()

	spanContext := make(map[string]string)

	// Create a new root span since there was no parent spanContext provided through annotations
	ctxWithTrace, span := tracerProvider.Tracer(TracerName).Start(ctx, "ChaosEngine:Reconciler")
	defer span.End()
	span.SetAttributes(attribute.String("chaos_engine.name", engine.Name), attribute.String("chaos_engine.namespace", engine.Namespace))

	pro.Inject(ctxWithTrace, propagation.MapCarrier(spanContext))

	log.Debug("got tracing carrier", spanContext)
	if len(spanContext) == 0 {
		log.Debug("tracerProvider doesn't provide a traceId, tracing is disabled")
		return ctx
	}
	return ctxWithTrace
}

// GetMarshalledSpanFromContext Extract spanContext from the context and return it as json encoded string
func GetMarshalledSpanFromContext(ctx context.Context) string {
	carrier := make(map[string]string)
	pro := otel.GetTextMapPropagator()

	pro.Inject(ctx, propagation.MapCarrier(carrier))

	if len(carrier) == 0 {
		log.Error("spanContext not present in the context, unable to marshall")
		return ""
	}

	marshalled, err := json.Marshal(carrier)
	if err != nil {
		log.Error("unable to marshal span context, err: %s", err)
		return ""
	}
	if len(marshalled) >= 1024 {
		log.Error("marshalled span context is too large, unable to marshall")
		return ""
	}
	return string(marshalled)
}
