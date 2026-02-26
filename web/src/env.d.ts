/// <reference types="vite/client" />
/// <reference types="@rotki/ui-library/vite-plugin/client" />

interface ImportMetaEnv {
  readonly VITE_RECAPTCHA_SITE_KEY: string;
  readonly VITE_SITE_URL: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}

interface Window {
  grecaptcha?: {
    render: (
      container: HTMLElement,
      parameters: {
        'sitekey': string;
        'callback': (token: string) => void;
        'expired-callback': () => void;
        'error-callback': () => void;
        'theme'?: string;
        'size'?: string;
      },
    ) => number;
    reset: (id?: number) => void;
  };
  onRecaptchaLoaded?: () => void;
}
