import { exiftool } from "exiftool-vendored"
import type { Tags, ExifDateTime } from "exiftool-vendored"
import fs from 'node:fs'
import path from 'node:path'

const IMAGE_TYPES = [
  'avif',
  '.bmp',
  '.gif',
  '.jpg',
  '.jpeg',
  '.png',
  'tiff',
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
  const imageFiles: string[] = []
  const files = fs.readdirSync(imagesDirectory)
  
  for (const file of files) {
    const ext = path.extname(file).toLowerCase() 
    if (IMAGE_TYPES.includes(ext)) {
      try {
        imageFiles.push(file)
      }
      catch (err) {
        console.error(`Error reading EXIF data for file ${file}:`, err)
      }
    }
  }
  return imageFiles
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

export async function createThumbnailsForAllImages() {
  const filenames = await getImageDirectoryContents()
  for (const filename of filenames) {
    createThumbnail(filename)
  }
}

export async function createMetaDataForAllImages() {
  const filenames = await getImageDirectoryContents()
  for (const filename of filenames) {
    await getExifForImage(filename)
  }
}
