import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

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
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: 'en',
    debug: true,
    detection: {
      // 检测语言的顺序
      order: ['querystring', 'cookie', 'localStorage', 'sessionStorage', 'navigator', 'htmlTag'],
      // 缓存语言到localStorage
      caches: ['localStorage', 'cookie'],
      // 设置localStorage和cookie的名称
      lookupLocalStorage: 'i18nextLng',
      lookupCookie: 'i18nextLng'
    },
    interpolation: {
      escapeValue: false
    },
    // 添加语言切换时的回调
    react: {
      useSuspense: false
    }
  }, (err, t) => {
    if (err) return console.log('i18n初始化失败', err);
  });

export default i18n;