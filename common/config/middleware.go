package config

import (
	"errors"

	"github.com/odigos-io/odigos/common"
)

const (
	target = "MW_TARGET"
)

type Middleware struct{}

func (m *Middleware) DestType() common.DestinationType {
	return common.MiddlewareDestinationType
}

func (m *Middleware) ModifyConfig(dest ExporterConfigurer, currentConfig *Config) ([]string, error) {
	if !isTracingEnabled(dest) && !isMetricsEnabled(dest) && !isLoggingEnabled(dest) {
		return nil, errors.New("Middleware is not enabled for any supported signals, skipping")
	}

	_, exists := dest.GetConfig()[target]
	if !exists {
		return nil, errors.New("Middleware target not specified, gateway will not be configured for Middleware")
	}

	exporterName := "otlp/middleware-" + dest.GetID()
	currentConfig.Exporters[exporterName] = GenericMap{
		"endpoint": "${MW_TARGET}",
		"headers": GenericMap{
			"authorization": "${MW_API_KEY}",
		},
	}
	var pipelineNames []string
	if isTracingEnabled(dest) {
		tracesPipelineName := "traces/middleware-" + dest.GetID()
		currentConfig.Service.Pipelines[tracesPipelineName] = Pipeline{
			Exporters: []string{exporterName},
		}
		pipelineNames = append(pipelineNames, tracesPipelineName)
	}

	if isMetricsEnabled(dest) {
		metricsPipelineName := "metrics/middleware-" + dest.GetID()
		currentConfig.Service.Pipelines[metricsPipelineName] = Pipeline{
			Exporters: []string{exporterName},
		}
		pipelineNames = append(pipelineNames, metricsPipelineName)
	}

	if isLoggingEnabled(dest) {
		logsPipelineName := "logs/middleware-" + dest.GetID()
		currentConfig.Service.Pipelines[logsPipelineName] = Pipeline{
			Exporters: []string{exporterName},
		}
		pipelineNames = append(pipelineNames, logsPipelineName)
	}

	return pipelineNames, nil
}
