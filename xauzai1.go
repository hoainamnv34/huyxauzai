import csv
import os
import subprocess
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

def get_configmap_via_sshpass(ssh_host, ssh_user, ssh_password, namespace, configmap_name):
    """
    Use sshpass to retrieve a ConfigMap in YAML format from the Kubernetes cluster master.
    :param ssh_host: SSH host for the cluster master.
    :param ssh_user: SSH username.
    :param ssh_password: SSH password.
    :param namespace: Kubernetes namespace.
    :param configmap_name: ConfigMap name.
    :return: ConfigMap data as a dictionary or None if not found.
    """
    try:
        ssh_command = f"kubectl get configmap {configmap_name} -n {namespace} -o yaml"
        result = subprocess.run(
            ['sshpass', '-p', ssh_password, 'ssh', f"{ssh_user}@{ssh_host}", ssh_command],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            check=True
        )
        return yaml.safe_load(result.stdout)
    except subprocess.CalledProcessError as e:
        if 'NotFound' in e.stderr:
            print(f"ConfigMap {configmap_name} not found in namespace {namespace}.")
        else:
            print(f"Error fetching ConfigMap {configmap_name} from namespace {namespace} via SSH: {e.stderr}")
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

def main(csv_file, output_dir, ssh_host, ssh_user, ssh_password):
    """
    Main function to process CSV and export ConfigMaps via SSH with sshpass.
    :param csv_file: Path to the CSV file.
    :param output_dir: Directory to save YAML files.
    :param ssh_host: SSH host for the cluster master.
    :param ssh_user: SSH username.
    :param ssh_password: SSH password.
    """
    # Read namespaces and services from the CSV
    data = read_csv(csv_file)
    
    # Process each namespace and service
    for namespace, service in data:
        configmap_data = get_configmap_via_sshpass(ssh_host, ssh_user, ssh_password, namespace, service)
        if configmap_data:
            save_configmap_to_yaml(configmap_data, output_dir, namespace, service)

if __name__ == "__main__":
    csv_file_path = "namespaces_services.csv"  # Replace with your CSV file path
    output_directory = "output_configmaps"  # Replace with your desired output directory
    ssh_master_host = "master-node-ip"  # Replace with your master node IP
    ssh_user = "user"  # Replace with your SSH username
    ssh_password = "password"  # Replace with your SSH password
    
    main(csv_file_path, output_directory, ssh_master_host, ssh_user, ssh_password)
