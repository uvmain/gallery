import fs from 'node:fs'
import path from 'node:path'

export const thumbnailsPath = path.resolve(path.join(serverConfiguration.dataPath, 'thumbnails'))

export function createThumbnailsDirectory() {
  if (!fs.existsSync(thumbnailsPath)){
    console.info(`Creating ${thumbnailsPath} directory for images`)
    fs.mkdirSync(thumbnailsPath, { recursive: true })
  }
  else {
    console.info(`Thumbnails directory at ${thumbnailsPath} already exists`)
  }
};