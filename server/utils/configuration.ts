import { env } from 'node:process'

export const serverConfiguration = {
  adminUser: env.ADMIN_USER || 'admin',
  adminPassword: env.ADMIN_PASSWORD || 'password',
  dataPath: env.DATA_PATH || './data',
  sessionMaxAge: Number(env.SESSION_MAX_AGE) || 86400,
  imagePath: env.IMAGE_PATH || './data/images',
  thumbnailMaxPixels: Number(env.THUMBNAIL_MAX_PIXELS) || 400,
  imageMaxPixels: Number(env.IMAGE_MAX_PIXELS) || 1080
}