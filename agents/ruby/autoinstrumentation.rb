require 'bundler'

require 'opentelemetry-api'
require 'opentelemetry-sdk'
require 'opentelemetry-exporter-otlp'
require 'opentelemetry-instrumentation-all'

require 'opentelemetry-resource_detectors'
require 'opentelemetry-resource-detector-container'
require 'opentelemetry-resource-detector-google_cloud_platform'
require 'opentelemetry-resource-detector-azure'

require 'opentelemetry-helpers-mysql'
require 'opentelemetry-helpers-sql'
require 'opentelemetry-helpers-sql-obfuscation'

module OTelBundlerPatch
  def require(*_groups)
    super
    require_otel
  end

  def require_resources
    env = ENV['OTEL_RUBY_RESOURCE_DETECTORS'].to_s
    additional_resource = ::OpenTelemetry::SDK::Resources::Resource.create({})

    additional_resource = additional_resource.merge(::OpenTelemetry::Resource::Detector::Container.detect) if defined? ::OpenTelemetry::Resource::Detector::Container
    additional_resource = additional_resource.merge(::OpenTelemetry::Resource::Detector::GoogleCloudPlatform.detect) if defined? ::OpenTelemetry::Resource::Detector::GoogleCloudPlatform
    additional_resource = additional_resource.merge(::OpenTelemetry::Resource::Detector::Azure.detect) if defined? ::OpenTelemetry::Resource::Detector::Azure

    additional_resource
  end

  def require_otel
    lib = File.expand_path('..', __dir__)
    $LOAD_PATH.reject! { |path| path.include?('zero-code-instrumentation') }
    $LOAD_PATH.unshift(lib)

    begin
      OpenTelemetry::SDK.configure do |c|
        c.resource = require_resources
        c.use_all # enables all instrumentation!
      end
      OpenTelemetry.logger.info { 'Auto-instrumentation initialized' }
    rescue StandardError => e
      OpenTelemetry.logger.info { "Auto-instrumentation failed to initialize. Error: #{e.message}" }
    end
  end
end

additional_gem_path = Gem.dir # ??? Switch to "/var/odigos/ruby"
puts "Loading additional gems from path #{additional_gem_path}"

# google-protobuf is used for otel trace exporter
Dir.glob("#{additional_gem_path}/gems/*").each do |file|
  if file.include?('opentelemetry') || file.include?('google')
    puts "Unshift #{file.inspect}"
    $LOAD_PATH.unshift("#{file}/lib")
  end
end

Bundler::Runtime.prepend(OTelBundlerPatch)
Bundler.require if ENV['REQUIRE_BUNDLER'].to_s == 'true'
