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

export function getMimeType(filename: string) {
  const ext = path.extname(filename).toLowerCase()
  switch (ext) {
    case '.jpg':
    case '.jpeg':
      return 'image/jpeg'
    default:
      return `image/${ext}`
  }
}

export async function getImageDirectoryContents() {
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
