import sqlite3 from 'sqlite3'
import path from 'node:path'
import fs from 'node:fs'

const databaseDirectory = path.resolve(path.join(serverConfiguration.dataPath, 'database'))
const databasePath = path.resolve(path.join(databaseDirectory, 'db.sqlite'))

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

export function createDatabaseDirectory() {
  if (!fs.existsSync(databaseDirectory)){
    console.info(`Creating ${databaseDirectory} directory for database`)
    fs.mkdirSync(databaseDirectory, { recursive: true })
  }
  else {
    console.info(`Database directory at ${databaseDirectory} already exists`)
  }
}

export const db = new sqlite3.Database(databasePath, (err) => {
  if (err) {
    console.error(err.message)
  }
  else
    console.info('Connected to the database.')
})

export function createMetadataTable() {
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
}