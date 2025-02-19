#!/bin/bash

echo "Waiting for MongoDB to start..."
sleep 10  # Wait for MongoDB to be ready

echo "Initiating replica set..."
mongosh --host mongo-primary:27017 <<EOF
rs.initiate({
  _id: "rs0",
  members: [
    { _id: 0, host: "mongo-primary:27017" },
    { _id: 1, host: "mongo-secondary-1:27017" },
    { _id: 2, host: "mongo-secondary-2:27017" }
  ]
});
EOF

echo "Replica set initialized."
