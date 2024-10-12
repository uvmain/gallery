import { exiftool } from "exiftool-vendored"
import type { Tags, ExifDateTime } from "exiftool-vendored"
import fs from 'node:fs'
import path from 'node:path'

const IMAGE_TYPES = [
  '.avif',
  '.bmp',
  '.gif',
  '.jpg',
  '.jpeg',
  '.png',
  '.webp'
]

export const imagesDirectory = path.resolve(serverConfiguration.imagePath)

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

export async function getExifForImage(imagePath: string): Promise<ImageMetadata> {
  const imageMetaData = await getMetadataByFileName(imagePath)

  if (imageMetaData) {
    return imageMetaData
  }
  else {
    const fileTags: ImageMetadata = {}
    try {
      const tags: Tags = await exiftool.read(path.resolve(path.join(imagesDirectory, imagePath)))
      fileTags.aperture = tags.Aperture?.toString()
      fileTags.cameraModel = `${tags.Make} ${tags.Model}`
      fileTags.dateTaken = exifDateToJavascriptDate(tags.DateTimeOriginal as ExifDateTime)
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

function exifDateToJavascriptDate(exifDate: ExifDateTime): Date {
  return exifDate.toDate()
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
          console.info(rows)
          // Create thumbnail for each file
          try {
            const promises = rows.map(row => createThumbnail(`/api/thumbnail/${row.fileName}`))
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
  const filenames = await getImageDirectoryContents()
  for (const filename of filenames) {
    await getExifForImage(filename)
  }
}
