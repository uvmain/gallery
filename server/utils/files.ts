import fs from 'node:fs'
import path from 'node:path'

export const IMAGE_TYPES = [
  '.avif',
  '.bmp',
  '.gif',
  '.jpg',
  '.jpeg',
  '.png',
  '.webp'
]

export async function* ls(filePath: string): AsyncGenerator<string> {
  yield filePath
  for (const dirent of await fs.promises.readdir(filePath, { withFileTypes: true })) {
    if (dirent.isDirectory()) {
      yield* ls(path.join(filePath, dirent.name))
    }
    else {
      yield path.join(filePath, dirent.name)
    }
  }
}

export async function toArray<T>(iter: AsyncIterable<T>): Promise<T[]> {
  const result: T[] = []
  for await (const x of iter) {
    const ext = path.extname(x as string).toLowerCase()
    if (IMAGE_TYPES.includes(ext)) {
      const parsedFilename = `${x}`.replace(`${imagesDirectory}\\`,'')
      console.info(`Found image: ${parsedFilename}`)
      result.push(parsedFilename as T)
    }
  }
  return result
}

export function toSlug(filename?: string): string {
  return filename?.replaceAll('\\','-').replaceAll('/','-') || ''
}