import csv
import os
from kubernetes import client, config
import yaml

def read_csv(file_path):
    """
    Read namespace and service information from a CSV file.
    :param file_path: Path to the CSV file.
    :return: List of tuples (namespace, service).
    """
    data = []
    try:
        with open(file_path, mode='r') as file:
            reader = csv.DictReader(file)
            for row in reader:
                namespace = row['Namespace'].strip()
                service = row['Service'].strip()
                if namespace and service:
                    data.append((namespace, service))
    except Exception as e:
        print(f"Error reading CSV file: {e}")
    return data

def get_configmap(namespace, configmap_name):
    """
    Retrieve a ConfigMap by name from a specific namespace.
    :param namespace: Kubernetes namespace.
    :param configmap_name: ConfigMap name.
    :return: ConfigMap data as a dictionary or None if not found.
    """
    try:
        v1 = client.CoreV1Api()
        configmap = v1.read_namespaced_config_map(name=configmap_name, namespace=namespace)
        return configmap.to_dict()
    except client.exceptions.ApiException as e:
        if e.status == 404:
            print(f"ConfigMap {configmap_name} not found in namespace {namespace}.")
        else:
            print(f"Error fetching ConfigMap {configmap_name} from namespace {namespace}: {e}")
    return None

def save_configmap_to_yaml(configmap_data, output_dir, namespace, configmap_name):
    """
    Save a ConfigMap as a YAML file.
    :param configmap_data: ConfigMap data as a dictionary.
    :param output_dir: Directory to save the YAML file.
    :param namespace: Kubernetes namespace of the ConfigMap.
    :param configmap_name: Name of the ConfigMap.
    """
    try:
        namespace_dir = os.path.join(output_dir, namespace)
        os.makedirs(namespace_dir, exist_ok=True)
        file_path = os.path.join(namespace_dir, f"{configmap_name}.yaml")
        
        with open(file_path, 'w') as yaml_file:
            yaml.dump(configmap_data, yaml_file, default_flow_style=False)
        
        print(f"Saved ConfigMap {configmap_name} from namespace {namespace} to {file_path}.")
    except Exception as e:
        print(f"Error saving ConfigMap {configmap_name} to YAML: {e}")

def main(csv_file, output_dir):
    """
    Main function to process CSV and export ConfigMaps.
    :param csv_file: Path to the CSV file.
    :param output_dir: Directory to save YAML files.
    """
    # Load Kubernetes config (default from kubeconfig or in-cluster)
    try:
        config.load_kube_config()  # Use load_incluster_config() if running inside a cluster
    except Exception as e:
        print(f"Error loading Kubernetes config: {e}")
        return

    # Read namespaces and services from the CSV
    data = read_csv(csv_file)
    
    # Process each namespace and service
    for namespace, service in data:
        configmap_data = get_configmap(namespace, service)
        if configmap_data:
            save_configmap_to_yaml(configmap_data, output_dir, namespace, service)

if __name__ == "__main__":
    csv_file_path = "namespaces_services.csv"  # Replace with your CSV file path
    output_directory = "output_configmaps"  # Replace with your desired output directory
    main(csv_file_path, output_directory)
