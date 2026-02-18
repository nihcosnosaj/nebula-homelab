## 🌌 Nebula

Nebula is my attempt at a budget-friendly AWS homelab setup. 

This is a "workshop" repo of mine for personal use. Things in here will not work, or they might. I put my experiments, failures, and wins in here. 

Currently working state is I can provision AWS EC2 instances via Terraform and configure k3s on those via Ansible. Ansible installs k3s and joins the two worker nodes to the control plane. 

I'm currently using Spot Instances to keep costs to a minimum. Granted this comes with other issues, like AWS terminating my nodes with a mere two minute warning. Fair, I'm getting that compute at a deep discount but still trying to gracefully handle this. I'm a bit too stubborn at the moment to just splurge for three persistent EC2 instances. 