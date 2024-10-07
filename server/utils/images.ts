import fs from 'node:fs'
import path from 'node:path'
import { exiftool } from "exiftool-vendored"
import type { Tags, ExifDateTime } from "exiftool-vendored"

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

const imagesPath = path.resolve(serverConfiguration.imagePath)

if (!fs.existsSync(imagesPath)){
  console.info(`Creating ${imagesPath} directory for images`)
  fs.mkdirSync(imagesPath, { recursive: true })
}
else {
  console.info(`Image directory at ${imagesPath} already exists`)
}

function exifDateToJavascriptDate(exifDate: ExifDateTime) {
  return exifDate.toDate()
}

export async function getImageDirectoryContents() {
  const imageFiles: string[] = []
  const files = fs.readdirSync(imagesPath)
  
  for (const file of files) {
    const ext = path.extname(file).toLowerCase()
    // const fileTags: ImageMetadata = {}
    
    if (IMAGE_TYPES.includes(ext)) {
      try {
        // const tags: Tags = await exiftool.read(path.resolve(path.join(imagesPath, file)))

        // fileTags.aperture = tags.Aperture?.toString()
        // fileTags.cameraModel = `${tags.Make} ${tags.Model}`
        // fileTags.dateTaken = exifDateToJavascriptDate(tags.DateTimeOriginal as ExifDateTime)
        // fileTags.exposureMode = tags.ExposureProgram
        // fileTags.fileName = file
        // fileTags.flashStatus = tags.Flash
        // fileTags.focusLength = tags.FocalLength
        // fileTags.iso = tags.ISO?.toString()
        // fileTags.lensModel = tags.Lens
        // fileTags.shutterSpeed = tags.ShutterSpeed
        // fileTags.whiteBalance = tags.WhiteBalance

        // imageFiles.push({ file, fileTags })
        imageFiles.push(file)
      }
      catch (err) {
        console.error(`Error reading EXIF data for file ${file}:`, err)
      }
    }
    else {
      console.warn(`${ext} file ${file} is not a recognized image`)
    }
  }

  return imageFiles
}
