use("i4")

db.createCollection("clients")
db.clients.createIndex({ code: 1 }, { unique: true })
db.clients.createIndex({ apiKey: 1 }, { unique: true })
