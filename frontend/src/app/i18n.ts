import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

// 导入语言文件
import enTranslation from '../locales/en/translation.json';
import zhTranslation from '../locales/zh/translation.json';

const resources = {
  en: {
    translation: enTranslation
  },
  zh: {
    translation: zhTranslation
  }
};

i18n
  .use(initReactI18next) // 将i18n绑定到react-i18next
  .init({
    resources,
    lng: 'en', // 默认语言
    fallbackLng: 'en', // 备用语言
    interpolation: {
      escapeValue: false // React已经安全地处理了XSS
    }
  });

export default i18n;