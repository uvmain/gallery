import {
  defineConfig,
  presetAttributify,
  presetIcons,
  presetTypography,
  presetUno,
  presetWebFonts,
  presetWind,
  transformerDirectives,
  transformerVariantGroup,
} from 'unocss'

function getSafelist(): string[] {
  const base = 'prose prose-sm m-auto text-left'.split(' ')
  const unusedSafelist: string[] = []
  return [...unusedSafelist, ...base]
}

export default defineConfig({
  shortcuts: {
    text: 'text-neutral-700 dark:text-neutral-100',
    button: 'border-1 border-neutral-700 rounded-sm border-solid px-4 py-2 outline-none dark:border-neutral-300 hover:bg-neutral-300 standard hover:text-dark hover:shadow-dark hover:shadow-md hover:shadow-op-20 hover:dark:bg-neutral-700 dark:hover:text-white dark:hover:shadow-white',
    standard: 'bg-neutral-200 dark:bg-neutral-800 text',
    modal: 'bg-neutral-200 dark:bg-neutral-900 text border-solid border-1 border-neutral-400 dark:border-neutral-500 rounded-sm',
    tooltip: 'dark:border-neutral-200 border-neutral-800 border-1 border-solid rounded-sm px-2 py-1 text-sm standard invisible absolute group-hover:visible opacity-90 ml-2',
    input: 'border-1 border-neutral-800 rounded border-dashed px-2 text-xl outline-none dark:border-neutral-200 standard placeholder:text',
  },
  theme: {
    colors: {},
  },
  presets: [
    presetUno(),
    presetAttributify(),
    presetIcons(),
    presetWind(),
    presetTypography(),
    presetWebFonts({
      fonts: {
        sans: 'Open Sans',
        serif: 'DM Serif Display',
        mono: 'DM Mono',
      },
    }),
  ],
  transformers: [
    transformerDirectives(),
    transformerVariantGroup(),
  ],
  safelist: getSafelist(),
})
