/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string;
  readonly VITE_ENABLE_CAPTCHA: string;
  readonly VITE_DISABLE_REGISTRATION: string;
  readonly VITE_ENABLE_EMAIL_VERIFICATION: string;
  // 更多环境变量...
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}