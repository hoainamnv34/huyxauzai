import os
import subprocess

# Đường dẫn đến file chứa danh sách repo và thư mục đích
file_path = 'labels.txt'  # Thay bằng tên file chứa danh sách repo
destination_folder = '/path/repo'  # Thay bằng đường dẫn thư mục đích

# Tạo thư mục đích nếu chưa tồn tại
os.makedirs(destination_folder, exist_ok=True)

# Đọc danh sách repo từ file
with open(file_path, 'r') as file:
    repos = [line.strip() for line in file if line.strip()]

# Clone từng repo
for repo in repos:
    repo_url = f"https://github.com/{repo}.git"
    repo_name = os.path.basename(repo)
    target_path = os.path.join(destination_folder, repo_name)
    
    try:
        print(f"Cloning {repo_url} into {target_path}...")
        subprocess.run(['git', 'clone', repo_url, target_path], check=True)
        print(f"Successfully cloned {repo}.")
    except subprocess.CalledProcessError:
        print(f"Failed to clone {repo}. Please check the repository URL or your network connection.")
