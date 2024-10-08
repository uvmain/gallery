import path from 'node:path'
import { exiftool } from "exiftool-vendored"
import type { Tags, ExifDateTime } from "exiftool-vendored"

export default defineEventHandler(async (event) => {
  const slug = event.context.params?.slug

  if (!slug) {
    throw createError({
      statusCode: 400,
      statusMessage: 'Missing slug parameter',
    })
  }

  try {
    const exifData = await getExifForImage(slug)
    return exifData
  }
  catch (error) {
    throw createError({
      statusCode: 404,
      statusMessage: `${error}`,
    })
  }
})

function exifDateToJavascriptDate(exifDate: ExifDateTime) {
  return exifDate.toDate()
}

export async function getExifForImage(imagePath: string) {
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
    return fileTags
  }
  catch (error) {
    throw createError({
      statusCode: 404,
      statusMessage: `${error}`,
    })
  }
}