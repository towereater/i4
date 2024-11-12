var res

use("i4")

db.createCollection("metadata")
res = db.metadata.createIndex({ hash: 1 }, { name: "hash-idx" })
res = db.metadata.createIndex({ client: 1, machine: 1, timestamp: 1 }, { name: "client-idx" })

db.createCollection("content")
res = db.content.createIndex({ hash: 1 }, { name: "hash-idx" })

db.createCollection("timedata")
res = db.content.createIndex({ client: 1, machine: 1, timestamp: 1 }, { name: "client-idx" })
