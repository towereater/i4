use("i4")

db.createCollection("clients")
db.clients.createIndex({ code: 1 }, { name: "code-idx", unique: true })
db.clients.createIndex({ apiKey: 1 }, { name: "apiKey-idx", unique: true })
db.clients.insertOne({
    code: "0000",
    name: "admin",
    apiKey: "cj9c023kkx9rjc493j4cjeprc",
})

db.createCollection("uplmeta")
db.uplmeta.createIndex({ hash: 1 }, { name: "hash-idx" })
db.uplmeta.createIndex({ client: 1, machine: 1, ts: 1 }, { name: "client-idx" })

db.createCollection("uplcont")
db.uplcont.createIndex({ hash: 1 }, { name: "hash-idx" })

use("i4-0001")

db.createCollection("datagau")
db.datagau.createIndex({ machine: 1, key: 1, value: 1, ts: 1 }, { name: "client-idx" })

db.createCollection("dataint")
db.dataint.createIndex({ machine: 1, key: 1, value: 1, start: 1, end: 1 }, { name: "client-idx" })
