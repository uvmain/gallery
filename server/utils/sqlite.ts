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

const insertMetadataSql = `INSERT INTO metadata (
  fileName, title, dateTaken, dateUploaded, cameraModel, lensModel, 
  aperture, shutterSpeed, flashStatus, focusLength, iso, exposureMode, whiteBalance
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

const selectMetadataSql = `SELECT * FROM metadata WHERE fileName = ?;`

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

export async function createMetadataTable() {
  return new Promise((resolve) => {
    db.all(createMetadataTableSql, [], (err) => {
      if (err) {
        if (err.message === "SQLITE_ERROR: table metadata already exists") {
          console.info('Metadata table already exists.')
          resolve({ outcome: 'success' })
        }
        else {
          console.warn(`Failed to create metadata table; ${err.message}`)
          resolve({ outcome: err })
        }
      }
      else {
        console.info('Created metadata table.')
        resolve({ outcome: 'success' })
      }
    })
  })
}

export async function insertMetadata(metadata: ImageMetadata) {
  try {
    await db.run(insertMetadataSql, [
      metadata.fileName,
      metadata.title,
      metadata.dateTaken?.toISOString(),
      metadata.dateUploaded?.toISOString(),
      metadata.cameraModel,
      metadata.lensModel,
      metadata.aperture,
      metadata.shutterSpeed,
      metadata.flashStatus,
      metadata.focusLength,
      metadata.iso,
      metadata.exposureMode,
      metadata.whiteBalance
    ], (err) => {
      if (`${err}` == "SQLITE_CONSTRAINT: UNIQUE constraint failed: metadata.fileName")
        console.info(`Metadata already exists for file: ${metadata.fileName}`)
      else 
      console.info(`Metadata insert for ${metadata.fileName} failed: ${err}`)
    })
    console.info(`Inserted metadata for file: ${metadata.fileName}`)
  }
  catch (error) {
    throw createError({
      statusCode: 500,
      statusMessage: `${error}`,
    })
  }
}

export async function getMetadataByFileName(fileName: string): Promise<ImageMetadata | null> {
  return new Promise((resolve, reject) => {
    db.get(selectMetadataSql, [fileName], (err: Error, row: ImageMetadata) => {
      if (err) {
        console.error(`Failed to retrieve metadata for file ${fileName}; ${err.message}`)
        reject({
          statusCode: 500,
          statusMessage: `Error retrieving metadata: ${err.message}`,
        })
      }
      else {
        if (row) {
          console.info(`Retrieved metadata from database for file: ${fileName}`)
          resolve(row)
        }
        else {
          console.warn(`No metadata found for file: ${fileName}`)
          resolve(null)
        }
      }
    })
  })
}
