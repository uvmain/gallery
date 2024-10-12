import fs from 'node:fs'
import path from 'node:path'
import sharp from 'sharp'

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

export async function createThumbnail(filename: string) {
  try {
    const thumbnailPath = path.resolve(thumbnailsPath, filename)
    const imagePath = path.resolve(imagesDirectory, filename)

    console.info(`creating thumbnail ${thumbnailPath} for image ${imagePath}`)

    // Check if the thumbnail already exists
    try {
      await fs.promises.access(thumbnailPath, fs.constants.F_OK)
      console.info(`Thumbnail already exists: ${filename}`)
      return thumbnailPath
    }
    catch {
      // If the file doesn't exist, continue with thumbnail creation
    }

    const fileBuffer = await fs.promises.readFile(imagePath)

    const resizedImageBuffer = await sharp(fileBuffer)
      .resize({
        width: serverConfiguration.thumbnailMaxPixels,
        height: serverConfiguration.thumbnailMaxPixels,
        fit: 'inside'
      })
      .webp()
      .toBuffer()

    // Save the thumbnail as a webp file in the thumbnails directory
    console.info(`Saving thumbnail to file system: ${filename}`)
    await fs.promises.writeFile(thumbnailPath, resizedImageBuffer)

    return resizedImageBuffer
  }
  catch (error) {
    throw createError({
      statusCode: 403,
      statusMessage: `${error}`,
    })
  }
}
