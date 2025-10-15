#!/bin/bash

# MongoDB credentials
USERNAME="admin"
PASSWORD="secret"
AUTH_DB="admin"

# Backup directory
BACKUP_DIR="/home/ec2-user/bkpMongoDbProd/flytura"

# Loop through all .bson files in the backup directory
for bson_file in "$BACKUP_DIR"/*.bson; do
    # Extract collection name from the file name
    collection_name=$(basename "$bson_file" .bson)

    echo "Restoring collection: $collection_name"

    # Run mongorestore for each collection
    mongorestore --username "$USERNAME" --password "$PASSWORD" \
                 --authenticationDatabase "$AUTH_DB" \
                 --drop --db flytura --collection "$collection_name" "$bson_file"
done

echo "MongoDB restoration completed."
