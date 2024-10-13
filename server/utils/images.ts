import fs from 'node:fs'
import path from 'node:path'
import { exiftool } from "exiftool-vendored"
import type { ExifDateTime, Tags } from "exiftool-vendored"

export const imagesDirectory = path.resolve(serverConfiguration.imagePath)

let directoryListing: string[] = []

export function getCachedDirectoryListing(): string[] {
  return directoryListing
}

export function createImagesDirectory() {
  if (!fs.existsSync(imagesDirectory)){
    console.info(`Creating ${imagesDirectory} directory for images`)
    fs.mkdirSync(imagesDirectory, { recursive: true })
  }
  else {
    console.info(`Image directory at ${imagesDirectory} already exists`)
  }
}

export function getMimeType(filename: string): string {
  const ext = path.extname(filename).toLowerCase()
  switch (ext) {
    case '.jpg':
    case '.jpeg':
      return 'image/jpeg'
    default:
      return `image/${ext}`
  }
}

export async function getImageDirectoryContents(): Promise<string[]> {  
  try {
    const files = await fs.promises.readdir(imagesDirectory, {recursive: false})
    console.info('directory has been listed')
    const filteredFiles = await Promise.all(
      files.map(async (file) => {
        const ext = path.extname(file).toLowerCase()
        if (IMAGE_TYPES.includes(ext)) {
          console.info(`Found image: ${file}`)
          return file
        }
        else {
          console.info(`File ${file} has been filtered out`)
        }
      })
    )

    return filteredFiles.filter(Boolean) as string[]
  }
 catch (err) {
    console.error('Error reading directory:', err)
    return []
  }
}

export async function getImageDirectoryListing(): Promise<string[]> {
  directoryListing = await toArray(ls(imagesDirectory))
  return directoryListing
}

export async function getExifForImage(imagePath: string): Promise<ImageMetadata> {
  const imageMetaData = await getMetadataByFileName(imagePath, false)

  if (imageMetaData) {
    return imageMetaData
  }
  else {
    const fileTags: ImageMetadata = {}
    try {
      const tags: Tags = await exiftool.read(path.resolve(path.join(imagesDirectory, imagePath)))
      fileTags.aperture = tags.Aperture?.toString()
      fileTags.cameraModel = `${tags.Make} ${tags.Model}`
      fileTags.dateTaken = ((tags.DigitalCreationDateTime || tags.DateTimeDigitized || tags.FileCreateDate || tags.DateTimeOriginal || tags.DateTimeCreated) as ExifDateTime).toISOString()
      fileTags.exposureMode = tags.ExposureProgram
      fileTags.fileName = imagePath
      fileTags.flashStatus = tags.Flash
      fileTags.focusLength = tags.FocalLength
      fileTags.iso = tags.ISO?.toString()
      fileTags.lensModel = tags.Lens
      fileTags.shutterSpeed = tags.ShutterSpeed
      fileTags.whiteBalance = tags.WhiteBalance
      insertMetadata(fileTags)
      return fileTags
    }
    catch (error) {
      throw createError({
        statusCode: 404,
        statusMessage: `${error}`,
      })
    }
  }
}

export async function createThumbnailsForAllImages(): Promise<void> {
  return new Promise((resolve, reject) => {
    // Fetch filenames from the database
    db.all('SELECT filename FROM metadata;', async (err: Error, rows: { fileName: string }[]) => {
      if (err) {
        console.error(`Failed to retrieve all filenames; ${err.message}`)
        reject({
          statusCode: 500,
          statusMessage: `Error retrieving filenames: ${err.message}`,
        })
      }
      else {
        if (rows && rows.length > 0) {
          console.info('Retrieved all filenames from database')
          // Create thumbnail for each file
          try {
            const promises = rows.map(row => createThumbnail(row.fileName))
            // Wait for all thumbnails to be created
            await Promise.all(promises)
            console.info('Thumbnails created successfully')
            resolve()
          }
          catch (error) {
            console.error(`Error creating thumbnails: ${error}`)
            reject({
              statusCode: 500,
              statusMessage: `Error creating thumbnails: ${error}`,
            })
          }
        }
        else {
          console.warn('No metadata found')
          resolve()
        }
      }
    })
  })
}

export async function createMetaDataForAllImages() {
  const filenames = getCachedDirectoryListing()
  for (const filename of filenames) {
    console.info(`Getting exif for ${filename}`)
    await getExifForImage(filename)
  }
}

export async function removeMetadataForRemovedFiles() {
  const filenames = getCachedDirectoryListing()

  db.all('SELECT filename FROM metadata;', (err: Error, rows: { fileName: string }[]) => {
    if (err) {
      console.error(`Failed to retrieve all filenames; ${err.message}`)
    }
    else {
      const metadataFilenames = rows.map(row => row.fileName)

      for (const filename of metadataFilenames) {
        if (!filenames.includes(filename)) {
          console.info(`Deleting metadata for removed image ${filename}`)
          db.run('DELETE FROM metadata WHERE filename = ?', [filename], (deleteErr: Error) => {
            if (deleteErr) {
              console.error(`Failed to delete metadata for file ${filename}; ${deleteErr.message}`)
            }
          })
        }
      }
    }
  })
}
