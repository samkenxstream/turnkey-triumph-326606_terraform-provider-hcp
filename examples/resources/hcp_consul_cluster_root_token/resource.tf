resource "hcp_consul_cluster_root_token" "example" {
  cluster_id = var.cluster_id
  project_id = var.project_id
}