
import { TranslateService } from '@ngx-translate/core';

export function appInitFactory(translate: TranslateService) {
  return () => {
    const stored = localStorage.getItem('language');
    const chosen = stored ?? (translate.getBrowserLang() || 'es');

    translate.setFallbackLang('es');
    // Se você carrega traduções por HTTP, pode retornar uma Promise:
    // return translate.use(chosen).toPromise();

    translate.use(chosen);
  };

}
