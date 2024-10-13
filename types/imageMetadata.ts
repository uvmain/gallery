import type { ExifDateTime } from "exiftool-vendored"

export interface ImageMetadata {
  fileName?: string
  title?: string
  dateTaken?: string | ExifDateTime
  dateUploaded?: Date
  cameraModel?: string
  lensModel?: string
  aperture?: string
  shutterSpeed?: string
  flashStatus?: string
  focusLength?: string
  iso?: string
  exposureMode?: string
  whiteBalance?: string
}