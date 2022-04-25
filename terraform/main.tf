resource "google_compute_region_instance_group_manager" "envoy-mig" {
  name = substr("envoy-${local.environment}-group-manager-${md5(google_compute_instance_template.envoy-instance-template.name)}", 0, 63)

  base_instance_name = "envoy-${local.environment}"
  region             = var.region
  wait_for_instances = true

  version {
    instance_template = google_compute_instance_template.envoy-instance-template.id
  }

  lifecycle {
    create_before_destroy = true
  }
}