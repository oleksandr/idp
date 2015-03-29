#!/bin/bash -e

# Reset DB
#rm -rf db.sqlite3 && sqlite3 db.sqlite3 < sql/sqlite3.sql
#psql -f sql/postgres.sql idp_dev
#mysql -u root idp_dev < sql/mysql.sql
idp-cli db drop --please
idp-cli db create

#
# Environment
#
. env.sh

#
# Domains
#
echo -n "Creating domains... "

# Create various test domains
output=$(idp-cli domains add --name=domain1.com --description="Test domain #1")
ID1=$(echo "$output" | awk '{ print $2}')

output=$(idp-cli domains add --name=domain2.com --description="Test domain #2")
ID2=$(echo "$output" | awk '{ print $2}')

output=$(idp-cli domains add --name=domain3.com --description="Test domain #3")
ID3=$(echo "$output" | awk '{ print $2}')

output=$(idp-cli domains add --name=ugly.org --description="Disabled when created" --disable)
ID4=$(echo "$output" | awk '{ print $2}')

output=$(idp-cli domains add --name=delete.me --description="Domain to delete" --disable)
ID5=$(echo "$output" | awk '{ print $2}')

echo "ok"

#
# Users
#
echo -n "Creating users... "

# Add users to domain1.com
output=$(idp-cli users add \
    --domain="$ID1" \
    --name=user1 \
    --password=pass1)
U1=$(echo "$output" | awk '{ print $2}')
output=$(idp-cli users add \
    --domain="$ID1" \
    --name=user2 \
    --password=pass2)
U2=$(echo "$output" | awk '{ print $2}')
output=$(idp-cli users add \
    --domain="$ID1" \
    --name=user3 \
    --password=pass3)
U3=$(echo "$output" | awk '{ print $2}')

# Add users to domain1.com, domain2.com, domain3.com
output=$(idp-cli users add \
    --domain="$ID1" \
    --domain="$ID2" \
    --domain="$ID3" \
    --name=manager1 \
    --password=pass1)
U4=$(echo "$output" | awk '{ print $2}')

# Create disabled user
output=$(idp-cli users add \
    --domain="$ID4" \
    --domain="$ID5" \
    --name=tester2 \
    --password=tester2 \
    --disable)
U5=$(echo "$output" | awk '{ print $2}')

echo "ok"

#
# Roles
#
echo "Creating roles... "

idp-cli roles add --name=admin --description="System administrator"
idp-cli roles add --name=manager --description="Basic manager"
idp-cli roles add --name=moderator --description="Content moderator"
idp-cli roles add --name=user --description="Basic user"
idp-cli roles add --name=tester --description="Tester (disabled by default)" --disable

#
# Permissions
#
echo "Creating permissions... "

idp-cli permissions add --name="*" --description="Allow all"
idp-cli permissions add --name="login" --description="Allow to sign in"
idp-cli permissions add --name="roles.*" --description="Manage roles"
idp-cli permissions add --name="permissions.*" --description="Manage permissions"
idp-cli permissions add --name="users.*" --description="Manage users"
idp-cli permissions add --name="domains.create" --description="Create domain"
idp-cli permissions add --name="domains.read" --description="Read access to domains"
idp-cli permissions add --name="posts.create" --description="Allow to create new post"
idp-cli permissions add --name="posts.delete" --description="Allow to delete any post"
idp-cli permissions add --name="dummy" --description="Disabled dummy permission" --disable

echo "Assigning permissions to roles... "

idp-cli roles update --add "*" admin

idp-cli roles update \
    --add "login" \
    --add "users.*" \
    --add "domains.create" \
    --add "domains.read" \
    --add "posts.create" \
    --add "posts.delete" \
    manager
idp-cli roles update \
    --add "login" \
    --add "domains.read" \
    --add "posts.create" \
    --add "posts.delete" \
    moderator
idp-cli roles update \
    --add "login" \
    --add "domains.read" \
    user
idp-cli roles update --add "dummy" tester

# List domains
echo ""
echo "DOMAINS"
echo "----------------"
idp-cli domains list

# Show users
echo ""
echo "USERS"
echo "----------------"

idp-cli users list

# Show users
echo ""
echo "ROLES"
echo "----------------"

idp-cli roles list

# Show users
echo ""
echo "PERMISSIONS"
echo "----------------"

idp-cli permissions list

exit 0
