import syFullNameLight from '../assets/logo-full-name-light.svg'
import syFullNameDark from '../assets/logo-full-name-dark.svg'
import syDetailLight from '../assets/logo-detail-light.svg'
import syDetailDark from '../assets/logo-detail-dark.svg'
import sySimpleLight from '../assets/logo-simple-light.svg'
import sySimpleDark from '../assets/logo-simple-dark.svg'
import syPoweredByLight from '../assets/logo-powered-by-light.svg'
import syPoweredByDark from '../assets/logo-powered-by-dark.svg'

import dpFullNameLight from '../assets/clients/digital-parts/logo-full-name-light.svg'
import dpFullNameDark from '../assets/clients/digital-parts/logo-full-name-dark.svg'
import dpNameOnlyLight from '../assets/clients/digital-parts/logo-name-only-light.svg'
import dpNameOnlyDark from '../assets/clients/digital-parts/logo-name-only-dark.svg'
import dpDetailLight from '../assets/clients/digital-parts/logo-detail-light.svg'
import dpDetailDark from '../assets/clients/digital-parts/logo-detail-dark.svg'

export type AssetName = 'logo-full-name' | 'logo-name-only' | 'logo-detail' | 'logo-simple' | 'logo-powered-by'
type Theme = 'light' | 'dark'

const CLIENT = (import.meta.env.VITE_CLIENT as string | undefined)?.toLowerCase() ?? 'switchyard'
export const isClientOverride = CLIENT !== 'switchyard'

const assetMap: Record<string, Partial<Record<AssetName, Record<Theme, string>>>> = {
  switchyard: {
    'logo-full-name':  { light: syFullNameLight,  dark: syFullNameDark  },
    'logo-detail':     { light: syDetailLight,    dark: syDetailDark    },
    'logo-simple':     { light: sySimpleLight,    dark: sySimpleDark    },
    'logo-powered-by': { light: syPoweredByLight, dark: syPoweredByDark },
  },
  'digital-parts': {
    'logo-full-name': { light: dpFullNameLight, dark: dpFullNameDark },
    'logo-name-only': { light: dpNameOnlyLight, dark: dpNameOnlyDark },
    'logo-detail':    { light: dpDetailLight,   dark: dpDetailDark   },
    // logo-simple and logo-powered-by fall through to Switchyard defaults
  },
}

export function resolveAsset(name: AssetName, theme: Theme): string | null {
  return assetMap[CLIENT]?.[name]?.[theme]
    ?? assetMap.switchyard[name]?.[theme]
    ?? null
}
