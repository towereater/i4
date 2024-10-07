var res

db.createCollection("metadata")
metadata = db.collection("metadata")
res = await metadata.createIndex({ hash: 1 })
res = await metadata.createIndex({ client: 1, machine: 1, timestamp: 1 })

db.createCollection("content")
content = db.collection("content")
res = await content.createIndex({ hash: 1 })
