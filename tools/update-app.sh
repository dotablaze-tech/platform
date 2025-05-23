#!/bin/bash

# Script: update-app.sh
# Description:
#   Updates the appVersion field in a Kubernetes Helm chart or manifest file
#   within the deployments Git repository. This script clones the repo, updates
#   the specified file, commits the change, and pushes it to the remote.
#
# Usage: ./update-app.sh <file_path> <version_number> <project_name>
#
# Example:
#   ./update-app.sh charts/my-app/Chart.yaml 1.2.3 my-app
#
# Requirements:
#   - SSH access to the deployments repo
#   - Git must be installed and configured
#
# Arguments:
#   <file_path>      Relative path to the file to update (within the repo)
#   <version_number> New version to set (e.g., 1.2.3)
#   <project_name>   Project name used in the commit message

git_repository="git@github.com:dotablaze-tech/deployments.git"
temp_dir=$(mktemp -d)

clone_repository() {
  git clone "${git_repository}" "${temp_dir}"
  cd "${temp_dir}" || exit
}

clean_repository() {
  rm -rf "${temp_dir}"
}

# Function to check if the repository is a Git repository
check_git_repository() {
  if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo "Error: Not inside a Git repository."
    clean_repository
    exit 1
  fi
}

# Function to check if a file exists and is tracked by Git
check_file() {
  local file_path=${1}
  if [ ! -f "${file_path}" ]; then
    echo "Error: File ${file_path} does not exist."
    clean_repository
    exit 1
  fi

  if ! git ls-files --error-unmatch "${file_path}" >/dev/null 2>&1; then
    echo "Error: File ${file_path} is not tracked by Git."
    clean_repository
    exit 1
  fi
}

# Function to check if there are uncommitted changes in the repository
check_uncommitted_changes() {
  if ! git diff-index --quiet HEAD --; then
    echo "Error: There are uncommitted changes in the repository. Please commit or stash them first."
    clean_repository
    exit 1
  fi
}

# Function to push committed changes in the repository to remote
push_changes() {
  if ! git push; then
    echo "Error: Failed to push changes to remote repository."
    clean_repository
    exit 1
  fi
}


# Function to update a file in a Git repository
update_file() {
  local file_path=${1}
  local new_version=${2}
  local project_name=${3}

  clone_repository
  check_git_repository
  check_file "${file_path}"
  check_uncommitted_changes

  git checkout HEAD "${file_path}"

  # Replace the app version line in the file
  sed -i "s/^appVersion: .*/appVersion: \"${new_version}\"/" "${file_path}"

  git add "${file_path}"
  git commit -m "chore(${project_name}): update app version to version ${new_version}"
  push_changes

  echo "[${project_name}]: appVersion in ${file_path} updated to ${new_version}."
  clean_repository
}

# Check if correct number of arguments are provided
if [ "${#}" -lt 3 ]; then
  echo "Usage: ${0} <file_path> <version_number> <project_name>"
  clean_repository
  exit 1
fi

update_file "${1}" "${2}" "${3}"
