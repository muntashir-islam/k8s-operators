import kopf
import kubernetes.client
from kubernetes import client, config

# Load the Kubernetes config
config.load_kube_config()

def copy_resource(api_client, resource_type, source_namespace, target_namespace, resource_name):
    # Define API calls based on the resource type
    if resource_type == "configmap":
        api = client.CoreV1Api(api_client)
        source_resource = api.read_namespaced_config_map(name=resource_name, namespace=source_namespace)
        new_resource = client.V1ConfigMap(
            api_version="v1",
            kind="ConfigMap",
            metadata=client.V1ObjectMeta(name=resource_name, namespace=target_namespace),
            data=source_resource.data,
            binary_data=source_resource.binary_data
        )
    elif resource_type == "secret":
        api = client.CoreV1Api(api_client)
        source_resource = api.read_namespaced_secret(name=resource_name, namespace=source_namespace)
        new_resource = client.V1Secret(
            api_version="v1",
            kind="Secret",
            metadata=client.V1ObjectMeta(name=resource_name, namespace=target_namespace),
            data=source_resource.data,
            string_data=source_resource.string_data,
            type=source_resource.type
        )
    else:
        raise ValueError("Unsupported resource type")

    try:
        # Try creating the resource
        api.create_namespaced_config_map(namespace=target_namespace, body=new_resource)
    except kubernetes.client.rest.ApiException as e:
        if e.status == 409:
            # Resource already exists, update it
            api.replace_namespaced_config_map(name=resource_name, namespace=target_namespace, body=new_resource)
        else:
            raise

@kopf.on.create('apps.example.com', 'v1alpha1', 'configmapsecretsyncs')
@kopf.on.update('apps.example.com', 'v1alpha1', 'configmapsecretsyncs')
def manage_configmap_secret_sync(spec, name, namespace, logger, **kwargs):
    api_client = kubernetes.client.ApiClient()
    source_namespace = spec.get('sourceNamespace')
    target_namespaces = spec.get('targetNamespaces', [])
    configmap_names = spec.get('configMapNames', [])
    secret_names = spec.get('secretNames', [])

    # Process ConfigMaps
    for cm_name in configmap_names:
        for target_namespace in target_namespaces:
            copy_resource(api_client, 'configmap', source_namespace, target_namespace, cm_name)

    # Process Secrets
    for secret_name in secret_names:
        for target_namespace in target_namespaces:
            copy_resource(api_client, 'secret', source_namespace, target_namespace, secret_name)

    logger.info(f"Syncing ConfigMaps and Secrets for {name} in namespace {namespace}")
