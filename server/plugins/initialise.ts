export default defineNitroPlugin(async() => {
  await createDatabaseDirectory()
  await createImagesDirectory()
  await createThumbnailsDirectory()
  await createMetadataTable()
  await createThumbnailsForAllImages()
  await createMetaDataForAllImages()
})