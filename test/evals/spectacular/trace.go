package spectaculareval

import (
	"bufio"
	"encoding/json"
	"strings"
)

func ParseTraceMetrics(trace string) TraceMetrics {
	var metrics TraceMetrics
	scanner := bufio.NewScanner(strings.NewReader(trace))
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}
		var event any
		if json.Unmarshal([]byte(line), &event) != nil {
			continue
		}
		metrics.Events++
		collectSemanticEvent(event, &metrics)
		walkTrace(event, &metrics)
	}
	return metrics
}

func collectSemanticEvent(value any, metrics *TraceMetrics) {
	event, ok := value.(map[string]any)
	if !ok {
		return
	}
	eventType, _ := event["type"].(string)
	path, _ := event["path"].(string)
	command, _ := event["command"].(string)
	switch eventType {
	case "spectacular.eval.observations":
		metrics.SemanticEvents++
		metrics.SemanticObserved = true
		metrics.ObservedFiles = append(metrics.ObservedFiles, stringList(event["files_read"])...)
		metrics.ObservedReferences = append(metrics.ObservedReferences, stringList(event["references_loaded"])...)
		metrics.ObservedCommands = append(metrics.ObservedCommands, stringList(event["commands_run"])...)
	case "spectacular.eval.file_read":
		if path != "" {
			metrics.SemanticEvents++
			metrics.ObservedFiles = append(metrics.ObservedFiles, path)
		}
	case "spectacular.eval.reference_loaded":
		if path != "" {
			metrics.SemanticEvents++
			metrics.ObservedReferences = append(metrics.ObservedReferences, path)
		}
	case "spectacular.eval.command":
		if command != "" {
			metrics.SemanticEvents++
			metrics.ObservedCommands = append(metrics.ObservedCommands, command)
		}
	}
}

func stringList(value any) []string {
	items, ok := value.([]any)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		if text, ok := item.(string); ok {
			result = append(result, text)
		}
	}
	return result
}

func walkTrace(value any, metrics *TraceMetrics) {
	switch typed := value.(type) {
	case map[string]any:
		for key, item := range typed {
			switch key {
			case "input_tokens":
				if number, ok := jsonNumber(item); ok && number > metrics.InputTokens {
					metrics.InputTokens = number
					metrics.UsageObserved = true
				}
			case "cached_input_tokens":
				if number, ok := jsonNumber(item); ok && number > metrics.CachedInputTokens {
					metrics.CachedInputTokens = number
				}
			case "output_tokens":
				if number, ok := jsonNumber(item); ok && number > metrics.OutputTokens {
					metrics.OutputTokens = number
				}
			case "type":
				if text, ok := item.(string); ok && isToolEvent(text) {
					metrics.ToolCalls++
				}
			}
			walkTrace(item, metrics)
		}
	case []any:
		for _, item := range typed {
			walkTrace(item, metrics)
		}
	}
}

func jsonNumber(value any) (int, bool) {
	number, ok := value.(float64)
	return int(number), ok
}

func isToolEvent(eventType string) bool {
	lower := strings.ToLower(eventType)
	return strings.Contains(lower, "tool_call") || strings.Contains(lower, "command_execution") || strings.Contains(lower, "mcp_call")
}
