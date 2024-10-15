// Create the tasks collection
db.createCollection("tasks");

// Create indexes on specified fields
db.tasks.createIndex({ "created_at": 1 });
db.tasks.createIndex({ "title": 1 });
db.tasks.createIndex({ "status": 1 });
