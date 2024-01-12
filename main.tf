provider "google" {
  project = "proj-delphinium"
}

resource "google_cloud_run_v2_service" "default" {
  name = "main"
  location = "europe-west1"
  client = "terraform"

  template {
    scaling {
      max_instance_count = 1
    }

    containers {
      image = "europe-west1-docker.pkg.dev/proj-delphinium/proj-delphinium/delphinium"
      ports {
        container_port = 9000
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