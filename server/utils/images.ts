import { exiftool } from "exiftool-vendored"
import type { Tags } from "exiftool-vendored"
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

export async function* ls(filePath: string = imagesDirectory): AsyncGenerator<string> {
  yield filePath
  for (const dirent of await fs.promises.readdir(filePath, { withFileTypes: true })) {
    if (dirent.isDirectory()) {
      yield* ls(path.join(filePath, dirent.name))
    }
    else {
      yield path.join(filePath, dirent.name)
    }
  }
}

async function toArray<T>(iter: AsyncIterable<T>): Promise<T[]> {
  const result: T[] = []
  for await (const x of iter) {
    const ext = path.extname(x as string).toLowerCase()
    if (IMAGE_TYPES.includes(ext)) {
      const parsedFilename = `${x}`.replace(`${imagesDirectory}\\`,'')
      console.info(`Found image: ${parsedFilename}`)
      result.push(parsedFilename as T)
    }
  }
  return result
}

export async function getImageDirectoryListing(): Promise<string[]> {
  directoryListing = await toArray(ls())
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
      fileTags.dateTaken = tags.DateTimeDigitized
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
