apiVersion: helm.cattle.io/v1
kind: HelmChart
metadata:
  name: aws-node-termination-handler
  namespace: kube-system
spec:
  chart: aws-node-termination-handler
  repo: https://aws.github.io/eks-charts
  targetNamespace: kube-system
  set:
    enableSpotInterruptionDraining: "true"
    enableScheduledEventDraining: "true"
    awsRegion: "${region}"
    metadataCheckIntervalInSeconds: "2"