db = db.getSiblingDB("adService")
db.ads.createIndex({ "startat": 1, "endat": 1 })
