import sqlite3 from 'sqlite3'
import process from 'node:process'
import path from 'node:path'
import fs from 'node:fs'

if (!fs.existsSync(serverConfiguration.databasePath)){
  console.info(`Creating ${serverConfiguration.databasePath} directory for database`)
  fs.mkdirSync(serverConfiguration.databasePath, { recursive: true });
}
else {
  console.info(`Database directory at ${serverConfiguration.databasePath} already exists`)
}

const databasePath = path.join(process.cwd(), serverConfiguration.databasePath, 'db.sqlite')

export const db = new sqlite3.Database(databasePath, (err) => {
  if (err) {
    console.error(err.message)
  }
  else
    console.info('Connected to the database.')
})

const createMetadataTableSql = `CREATE TABLE metadata (
  fileName TEXT PRIMARY KEY,
  title TEXT,
  dateTaken DATETIME,
  dateUploaded DATETIME,
  cameraModel TEXT,
  lensModel TEXT,
  aperture TEXT,
  shutterSpeed TEXT,
  flashStatus TEXT,
  focusLength TEXT,
  iso TEXT,
  exposureMode TEXT,
  whiteBalance TEXT
);`

db.all(createMetadataTableSql, [], (err) => {
  if (err) {
    if (err.message == "SQLITE_ERROR: table metadata already exists") {
      console.info('Metadata table already exists.')
      return { outcome: 'success' }
    }
    console.warn(`Failed to create metadata table; ${err.message}`)
    return { outcome: err }
  }
  console.info('Created metadata table.')
  return { outcome: 'success' }
})