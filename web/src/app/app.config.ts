import { ApplicationConfig } from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { HttpClient, provideHttpClient, withInterceptors } from '@angular/common/http';
import { AuthInterceptor } from './auth.interceptor';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { BrowserModule } from '@angular/platform-browser';
import { provideMomentDateAdapter } from '@angular/material-moment-adapter';
import { MAT_DATE_LOCALE, MAT_DATE_FORMATS } from '@angular/material/core';

import { registerLocaleData } from '@angular/common';
import localePt from '@angular/common/locales/pt';
import { DateFormatService } from './services/date-format.service.service';
import { TranslateLoader, TranslateModule } from '@ngx-translate/core';
import { TranslateHttpLoader } from '@ngx-translate/http-loader';
import { CustomTranslateLoader } from './services/customer-translate.service';
registerLocaleData(localePt);

export const MY_DATE_FORMATS = {
  parse: {
    dateInput: 'DD/MM/YYYY', // Formato para parsing da data (entrada)
  },
  display: {
    dateInput: 'DD/MM/YYYY', // Formato para exibição no campo de input
    monthYearLabel: 'MM/YYYY', // Formato para exibir mês e ano
    dateA11yLabel: 'DD/MM/YYYY',
    monthYearA11yLabel: 'MM/YYYY',
  },
};




const translateProviders = TranslateModule.forRoot({
  defaultLanguage: 'pt',
  loader: {
    provide: TranslateLoader,
    useFactory: (http: HttpClient) => new CustomTranslateLoader(http),
    deps: [HttpClient]
  }
}).providers ?? []; // se for undefined, usa array vazio

export const appConfig: ApplicationConfig = {
  providers: [
    provideHttpClient(),
    provideHttpClient(withInterceptors([AuthInterceptor])),
    provideRouter(routes),
    provideMomentDateAdapter(),
    { provide: MAT_DATE_LOCALE, useValue: 'pt-BR' },
    {
      provide: MAT_DATE_FORMATS,
      useFactory: (dateFormatService: DateFormatService) => {
        let dateFormat;
        dateFormatService.currentDateFormat.subscribe(format => dateFormat = format);
        return dateFormat;
      },
      deps: [DateFormatService]
    },
    ...translateProviders
  ]
};


