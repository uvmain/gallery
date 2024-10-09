export default defineNitroPlugin(async() => {
  createDatabaseDirectory()
  createMetadataTable()
  createImagesDirectory()
  createThumbnailsDirectory()
  createThumbnailsForAllImages()
  createMetaDataForAllImages()
})