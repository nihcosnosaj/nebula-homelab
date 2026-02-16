## 🌌 Nebula

Nebula is my attempt at a budget-friendly AWS homelab setup. 

### Motivation

I'm still trying to get my hands on a NUC or two, or some tiny-mini-micro PCs to build an actual, bare-metal homelab. In the meantime, I've been exploring other options. 

I first built out `kindling` which just uses `kind` (kubernetes-in-docker) for some ephemeral k8s clusters that I can boot up and tear down quickly. This worked for some time, but just felt a bit underwhelming. Plus, I wanted to sharpen my AWS and Terraform skills. 

As for the name, I give all my non-work machines some sort of celestial body or NASA mission name, hence `nebula`. 

### Features

As of now, this deploys 16GB RAM ARM nodes via AWS Spot Instances to keep the per-hour compute time to a minimum. I've added a price calculator that grabs the price of Spot Instances in real-time and calculates what you should expect to spend when you bring down the nodes. 

I also built out a manager in Go to track these costs as well as to ensure we terminate the nodes properly. 

I also configured SSO with the IAM Identity Center on AWS so we don't have to worry about static keys anywhere. 

### Architecture

Terraform does most of the leg work by building the VPC, public and private subnets, security groups, and actually deploying the EC2 instances. 

I used the AWS SDK in Go to do some automation of price counting and node teardown (if i wanted to keep the VPC and other config stuff up and running but save on the compute bill).
