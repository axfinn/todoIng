import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector'; // 添加语言检测器

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
  .use(LanguageDetector) // 添加语言检测器
  .use(initReactI18next) // 将i18n绑定到react-i18next
  .init({
    resources,
    fallbackLng: 'en', // 备用语言
    detection: {
      // 检测语言的顺序
      order: ['querystring', 'cookie', 'localStorage', 'sessionStorage', 'navigator', 'htmlTag'],
      // 缓存语言到localStorage
      caches: ['localStorage', 'cookie']
    },
    interpolation: {
      escapeValue: false // React已经安全地处理了XSS
    }
  });

export default i18n;