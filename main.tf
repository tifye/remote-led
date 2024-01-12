variable "project_id" {
  type = string
  sensitive = true
}

variable "artifact_repository_id" {
  type = string
  sensitive = true
}

variable "artifact_name" {
  type = string
  sensitive = true
}

variable "region" {
  type = string
}

provider "google" {
  project = var.project_id
}

resource "google_cloud_run_v2_service" "default" {
  name = "main"
  location = var.region
  client = "terraform"

  template {
    scaling {
      max_instance_count = 1
    }

    containers {
      image = format("%s-docker.pkg.dev/%s/%s/%s", var.region, var.project_id, var.artifact_repository_id, var.artifact_name)
      ports {
        container_port = 9000
      }
      env {
        name = "LED_SERVER_URL"
        value = "${LED_SERVER_URL}"
      }
    }
  }

  traffic {
    percent = 100
    type = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
  }
}

resource "google_cloud_run_v2_service_iam_member" "noauth" {
  location = google_cloud_run_v2_service.default.location
  name = google_cloud_run_v2_service.default.name
  role = "roles/run.invoker"
  member = "allUsers"
}