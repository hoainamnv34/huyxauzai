import csv
import os
import subprocess

def get_git_description(repo_path):
    """
    Lấy thông tin mô tả của git tags từ repository.

    :param repo_path: Đường dẫn đến repository.
    :return: Chuỗi mô tả từ lệnh `git describe --tags`, hoặc thông báo lỗi nếu không thành công.
    """
    try:
        result = subprocess.run(
            ['git', 'describe', '--tags'],
            cwd=repo_path,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            check=True
        )
        return result.stdout.strip()
    except subprocess.CalledProcessError as e:
        return f"Error: {e.stderr.strip()}"

def write_repos_to_csv(repos, csv_file):
    """
    Ghi danh sách repository và kết quả `git describe --tags` ra file CSV.

    :param repos: Danh sách đường dẫn repository.
    :param csv_file: Đường dẫn file CSV để ghi.
    """
    try:
        with open(csv_file, mode='w', newline='') as file:
            writer = csv.writer(file)
            writer.writerow(['Repository', 'Git Description'])  # Ghi header

            for repo_path in repos:
                repo_name = os.path.basename(repo_path)
                git_description = get_git_description(repo_path)
                writer.writerow([repo_name, git_description])

        print(f"Successfully written to {csv_file}.")
    except Exception as e:
        print(f"Failed to write to CSV. Error: {e}")

# Ví dụ sử dụng
def main():
    repos = [
        '/path/to/repo1',
        '/path/to/repo2',
        '/path/to/repo3'
    ]
    csv_file = '/path/to/output/repos.csv'
    write_repos_to_csv(repos, csv_file)

if __name__ == '__main__':
    main()
