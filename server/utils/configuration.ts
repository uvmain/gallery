import { env } from 'node:process'

export const serverConfiguration = {
  adminUser: env.ADMIN_USER || 'admin',
  adminPassword: env.ADMIN_PASSWORD || 'password',
  dataPath: env.DATA_PATH || './data',
  sessionMaxAge: (env.SESSION_MAX_AGE || 86400) as number,
  imagePath: env.IMAGE_PATH || './data/images'
}