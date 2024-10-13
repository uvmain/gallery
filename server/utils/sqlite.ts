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

const selectAllMetadataSql = `SELECT * FROM metadata;`

export async function createDatabaseDirectory(): Promise<void> {
  try {
    await fs.promises.access(databaseDirectory)
    console.info(`Database directory at ${databaseDirectory} already exists`)
  }
  catch {
    console.info(`Creating ${databaseDirectory} directory for database`)
    await fs.promises.mkdir(databaseDirectory, { recursive: true })
  }
}

export const db: sqlite3.Database = new sqlite3.Database(databasePath, (err) => {
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
      metadata.dateTaken,
      metadata.dateUploaded,
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
      else if (err?.message.length)
        console.info(`Metadata insert for ${metadata.fileName} failed: ${err.message}`)
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

export async function getMetadataByFileName(fileName: string, logError = true): Promise<ImageMetadata | null> {
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
          if (logError) {
            console.warn(`No metadata found for file: ${fileName}`)
          }
          resolve(null)
        }
      }
    })
  })
}

export async function getAllMetadata(): Promise<ImageMetadata[] | null> {
  return new Promise((resolve, reject) => {
    db.all(selectAllMetadataSql, (err: Error, rows: ImageMetadata[]) => {
      if (err) {
        console.error(`Failed to retrieve all metadata; ${err.message}`)
        reject({
          statusCode: 500,
          statusMessage: `Error retrieving metadata: ${err.message}`,
        })
      }
      else {
        if (rows) {
          console.info('Retrieved all metadata from database')
          resolve(rows)
        }
        else {
          console.warn('No metadata found')
          resolve(null)
        }
      }
    })
  })
}