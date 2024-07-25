## Config-Secret-Sync-Operator pythonic way
For reference purpose I have developed Config-Secret-Sync-Operator using KOFPs framework

## Description
Often, we need to deploy the same ConfigMaps and Secrets across multiple namespaces. This can lead to inconsistencies, especially if we forget to update all namespaces simultaneously. Instead of deploying the same ConfigMaps and Secrets manifests multiple times in separate namespaces, our operator can simplify this process. It will automatically sync ConfigMaps and Secrets to multiple namespaces, ensuring consistency and reducing manual effort.

## RBAC not integrated Yet
TBA
