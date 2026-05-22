import { TranslateLoader } from '@ngx-translate/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export class CustomTranslateLoader implements TranslateLoader {
  constructor(private http: HttpClient) {}

  // Método para carregar o arquivo de tradução do servidor ou arquivo estático
  getTranslation(lang: string): Observable<any> {
    const url = `assets/i18n/${lang}.json`;  // Caminho dos arquivos de tradução
    return this.http.get<any>(url);
  }
}
