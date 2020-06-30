output "service" {
  description = "The Generated Serverless Service"
  value       = twilio_serverless_service.service
}

output "function" {
  description = "The Generated Serverless Function"
  value       = twilio_serverless_function.function
}

output "function_version" {
  description = "The Generated Serverless Function Version"
  value       = twilio_serverless_function_version.function_version
}

output "asset" {
  description = "The Generated Serverless Asset"
  value       = twilio_serverless_asset.asset
}

output "asset_version" {
  description = "The Generated Serverless Asset Version"
  value       = twilio_serverless_asset_version.asset_version
}

output "build" {
  description = "The Generated Serverless Build"
  value       = twilio_serverless_build.build
}

output "environment" {
  description = "The Generated Serverless Environment"
  value       = twilio_serverless_environment.environment
}

output "deployment" {
  description = "The Generated Serverless Deployment"
  value       = twilio_serverless_deployment.deployment
}

