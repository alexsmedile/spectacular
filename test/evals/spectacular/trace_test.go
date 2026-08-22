package spectaculareval

import (
	"os"
	"path/filepath"
	"testing"
)

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

func TestParseCodexNativeTrace(t *testing.T) {
	trace := readTraceFixture(t, "codex.jsonl")
	metrics := ParseTraceMetrics(trace)
	if !metrics.SemanticObserved || !metrics.UsageObserved || metrics.InputTokens != 1200 || metrics.CachedInputTokens != 900 || metrics.ToolCalls != 1 {
		t.Fatalf("metrics=%+v", metrics)
	}
	if !listContainsFold(metrics.ObservedReferences, "SKILL.md") || !listContainsFold(metrics.ObservedReferences, "orient.md") || len(metrics.ObservedCommands) != 1 {
		t.Fatalf("metrics=%+v", metrics)
	}
}

func TestParseClaudeNativeTraceUsesTotalContextAndToolUse(t *testing.T) {
	trace := readTraceFixture(t, "claude.jsonl")
	metrics := ParseTraceMetrics(trace)
	if !metrics.SemanticObserved || !metrics.UsageObserved || metrics.InputTokens != 4218 || metrics.CachedInputTokens != 3200 || metrics.OutputTokens != 420 || metrics.ToolCalls != 2 {
		t.Fatalf("metrics=%+v", metrics)
	}
	if !listContainsFold(metrics.ObservedFiles, "runtime.md") || !listContainsFold(metrics.ObservedCommands, "tests/check.sh") {
		t.Fatalf("metrics=%+v", metrics)
	}
}

func TestNormalizedOpenCodeTraceIsCertified(t *testing.T) {
	trace := readTraceFixture(t, "opencode.jsonl")
	certification := CertifyTrace(trace, true)
	if !certification.Valid {
		t.Fatalf("certification=%+v", certification)
	}
	metrics := ParseTraceMetrics(trace)
	if metrics.InputTokens != 300 || metrics.CachedInputTokens != 100 || metrics.ToolCalls != 1 || !listContainsFold(metrics.ObservedReferences, "execute.md") {
		t.Fatalf("metrics=%+v", metrics)
	}
}

func TestSelfReportDoesNotEstablishSemanticTelemetry(t *testing.T) {
	metrics := ParseTraceMetrics(`{"type":"spectacular.eval.self_report","files_read":["secret.md"],"references_loaded":["runtime.md"],"commands_run":["rm bad"]}`)
	if metrics.SemanticObserved || len(metrics.ObservedFiles) != 0 || len(metrics.ObservedCommands) != 0 {
		t.Fatalf("self-report became telemetry: %+v", metrics)
	}
}

func TestTraceCertificationRejectsPlausibleButUnobservableAdapterOutput(t *testing.T) {
	if result := CertifyTrace(`{"type":"spectacular.eval.self_report","files_read":["src/a"]}`, true); result.Valid || len(result.Findings) != 3 {
		t.Fatalf("certification=%+v", result)
	}
	valid := `{"type":"spectacular.eval.observations","files_read":["src/a"],"references_loaded":[],"commands_run":["test"]}
{"type":"spectacular.eval.usage","input_tokens":100,"output_tokens":10}
{"type":"tool_call"}`
	if result := CertifyTrace(valid, true); !result.Valid {
		t.Fatalf("certification=%+v", result)
	}
}

func readTraceFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "traces", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
