import fs from 'node:fs'
import path from 'node:path'
import sharp from 'sharp'

export const thumbnailsPath = path.resolve(path.join(serverConfiguration.dataPath, 'thumbnails'))

export async function createThumbnailsDirectory(): Promise<void> {
  try {
    await fs.promises.access(thumbnailsPath)
    console.info(`Thumbnails directory at ${thumbnailsPath} already exists`)
  }
  catch {
    console.info(`Creating ${thumbnailsPath} directory for thumbnails`)
    await fs.promises.mkdir(thumbnailsPath, { recursive: true })
  }
}

export async function createThumbnail(filename: string) {
  try {
    const thumbnailPath = path.resolve(thumbnailsPath, toSlug(filename))
    const imagePath = path.resolve(imagesDirectory, filename)

    console.info(`Creating thumbnail ${thumbnailPath} for image ${imagePath}`)

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

export async function removeThumbnailsForRemovedFiles() {
  const slugs = await getAllSlugs() || []
  const thumbnailFiles = await toArray(ls(thumbnailsPath))
  for (const thumbnailFile of thumbnailFiles) {
    if (!slugs.includes(path.basename(thumbnailFile))) {
      const thumbnailPath = path.resolve(thumbnailsPath, thumbnailFile)
      console.info(`Removing thumbnail for ${thumbnailFile}`)
      fs.promises.unlink(thumbnailPath)
    }
  }
  return thumbnailFiles
}