package spectaculareval

import "testing"

func TestParseTraceMetricsUsesCumulativeMaximum(t *testing.T) {
	trace := `{"type":"item.command_execution","usage":{"input_tokens":100,"cached_input_tokens":20,"output_tokens":10}}
{"type":"turn.completed","usage":{"input_tokens":250,"cached_input_tokens":80,"output_tokens":40}}
ADAPTER_STDERR
ignored`
	metrics := ParseTraceMetrics(trace)
	if !metrics.UsageObserved || metrics.InputTokens != 250 || metrics.CachedInputTokens != 80 || metrics.OutputTokens != 40 {
		t.Fatalf("metrics=%+v", metrics)
	}
	if metrics.ToolCalls != 1 || metrics.Events != 2 {
		t.Fatalf("metrics=%+v", metrics)
	}
}

func TestMissingUsageRemainsUnsupported(t *testing.T) {
	metrics := ParseTraceMetrics(`{"type":"turn.completed"}`)
	if metrics.UsageObserved || metrics.InputTokens != 0 {
		t.Fatalf("metrics=%+v", metrics)
	}
}

func TestOutputOnlyUsageDoesNotClaimInputCoverage(t *testing.T) {
	metrics := ParseTraceMetrics(`{"type":"turn.completed","usage":{"output_tokens":25}}`)
	if metrics.UsageObserved || metrics.OutputTokens != 25 {
		t.Fatalf("metrics=%+v", metrics)
	}
}

func TestParseNormalizedSemanticObservations(t *testing.T) {
	metrics := ParseTraceMetrics(`{"type":"spectacular.eval.observations","files_read":["src/a"],"references_loaded":["references/execute.md"],"commands_run":["go test ./..."]}`)
	if !metrics.SemanticObserved || metrics.SemanticEvents != 1 || len(metrics.ObservedFiles) != 1 || len(metrics.ObservedReferences) != 1 || len(metrics.ObservedCommands) != 1 {
		t.Fatalf("metrics=%+v", metrics)
	}
}
