#!/bin/bash

# Define the old values
old_package_name="github.com/dwarvesf/go-api"
old_app_name="go-api"
old_contact_name="Andy"
old_contact_url="https://d.foundation"
old_contact_email="andy@d.foundation"
old_title="APP API DOCUMENT"
old_description="This is api document for APP API project."
old_version="v0.0.1"

# Specify the file to edit
file="./cmd/server/main.go"
script_file="./scripts/update_names.sh"

# Prompt the user for the new package name
read -p "Enter the new package name(eg: "$old_package_name"): " new_package_name
new_package_name=${new_package_name:-$old_package_name}

# Validate the package name format (you can customize this validation)
if [[ ! "$new_package_name" =~ ^[a-zA-Z0-9/._-]+$ ]]; then
  echo "Invalid package name format. Please enter a valid package name."
  exit 1
fi

read -p "Enter the new app name(eg: eg: "$old_app_name"): " new_app_name

# Prompt the user for the new values
read -p "Enter new contact name(eg: $old_contact_name): " new_contact_name
read -p "Enter new contact URL(eg: $old_contact_url): " new_contact_url
read -p "Enter new contact email(eg: $old_contact_email): " new_contact_email
read -p "Enter new title(eg: $old_title): " new_title
read -p "Enter new description(eg: $old_description): " new_description
read -p "Enter new version(eg: $old_version): " new_version

# If the value is empty, use the old value
new_app_name=${new_app_name:-$old_app_name}
new_contact_name=${new_contact_name:-$old_contact_name}
new_contact_url=${new_contact_url:-$old_contact_url}
new_contact_email=${new_contact_email:-$old_contact_email}
new_title=${new_title:-$old_title}
new_description=${new_description:-$old_description}
new_version=${new_version:-$old_version}

# show all value
echo "new_package_name: $new_package_name"
echo "new_app_name: $new_app_name"
echo "new_contact_name: $new_contact_name"
echo "new_contact_url: $new_contact_url"
echo "new_contact_email: $new_contact_email"
echo "new_title: $new_title"
echo "new_description: $new_description"
echo "new_version: $new_version"

# Use sed to replace occurrences of the old values with the new ones
sed -Ei "s|// @title\s*.*|// @title           $new_title|g" "$file"
sed -Ei "s|// @version\s*.*|// @version         $new_version|g" "$file"
sed -Ei "s|// @description\s*.*|// @description     $new_description|g" "$file"
sed -Ei "s|// @contact.name\s*.*|// @contact.name   $new_contact_name|g" "$file"
sed -Ei "s|// @contact.url\s*.*|// @contact.url    $new_contact_url|g" "$file"
sed -Ei "s|// @contact.email\s*.*|// @contact.email  $new_contact_email|g" "$file"

# Use find and sed to replace occurrences of the old package name with the new one
find . -type f ! -path "./.git/*" ! -wholename "$script_file" -exec sed -i "s@$old_package_name@$new_package_name@g" {} +
find . -type f ! -path "./.git/*" ! -wholename "$script_file" -exec sed -i "s@$old_app_name@$new_app_name@g" {} +

# Use find and sed to replace occurrences of the old package name with the new one
sed -i "s#old_package_name=\"$old_package_name\"#old_package_name=\"$new_package_name\"#" "$script_file"
sed -i "s#old_contact_name=\"$old_contact_name\"#old_contact_name=\"$new_contact_name\"#" "$script_file"
sed -i "s#old_contact_url=\"$old_contact_url\"#old_contact_url=\"$new_contact_url\"#" "$script_file"
sed -i "s#old_contact_email=\"$old_contact_email\"#old_contact_email=\"$new_contact_email\"#" "$script_file"
sed -i "s#old_title=\"$old_title\"#old_title=\"$new_title\"#" "$script_file"
sed -i "s#old_description=\"$old_description\"#old_description=\"$new_description\"#" "$script_file"
sed -i "s#old_version=\"$old_version\"#old_version=\"$new_version\"#" "$script_file"
sed -i "s#old_app_name=\"$old_app_name\"#old_app_name=\"$new_app_name\"#" "$script_file"

echo "Package name replacement complete."